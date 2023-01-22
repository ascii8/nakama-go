package nakama

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

// ClientHandler is the interface for connection handlers.
type ClientHandler interface {
	HttpClient() *http.Client
	SocketURL() (string, error)
	Token(context.Context) (string, error)
	Logf(string, ...interface{})
	Errf(string, ...interface{})
}

// ConnHandler is an empty interface that provides a clean, extensible way to
// register a type as a handler that is future-proofed. As used by WithConnHandler,
// a type that supports the any of the following smuggled interfaces:
//
//	ConnectHandler(context.Context)
//	DisconnectHandler(context.Context, error)
//	ErrorHandler(context.Context, *nakama.ErrorMsg)
//	ChannelMessageHandler(context.Context, *nakama.ChannelMessageMsg)
//	ChannelPresenceEventHandler(context.Context, *nakama.ChannelPresenceEventMsg)
//	MatchDataHandler(context.Context, *nakama.MatchDataMsg)
//	MatchPresenceEventHandler(context.Context, *nakama.MatchPresenceEventMsg)
//	MatchmakerMatchedHandler(context.Context, *nakama.MatchmakerMatchedMsg)
//	NotificationsHandler(context.Context, *nakama.NotificationsMsg)
//	StatusPresenceEventHandler(context.Context, *nakama.StatusPresenceEventMsg)
//	StreamDataHandler(context.Context, *nakama.StreamDataMsg)
//	StreamPresenceEventHandler(context.Context, *nakama.StreamPresenceEventMsg)
//
// Will have its method added to Conn as its respective <MessageType>Handler.
//
// For an overview of Go interface smuggling as a concept, see:
//
// https://utcc.utoronto.ca/~cks/space/blog/programming/GoInterfaceSmuggling
type ConnHandler interface{}

// Conn is a nakama realtime websocket connection.
type Conn struct {
	h                 ClientHandler
	url               string
	token             string
	binary            bool
	query             url.Values
	persist           bool
	backoffMax        time.Duration
	backoffMin        time.Duration
	backoffMultiplier float64

	ctx    context.Context
	ws     *websocket.Conn
	cancel func()
	stop   bool

	id  uint64
	out chan *res
	m   map[string]*res

	ConnectHandler              func(context.Context)
	DisconnectHandler           func(context.Context, error)
	ErrorHandler                func(context.Context, *ErrorMsg)
	ChannelMessageHandler       func(context.Context, *ChannelMessageMsg)
	ChannelPresenceEventHandler func(context.Context, *ChannelPresenceEventMsg)
	MatchDataHandler            func(context.Context, *MatchDataMsg)
	MatchPresenceEventHandler   func(context.Context, *MatchPresenceEventMsg)
	MatchmakerMatchedHandler    func(context.Context, *MatchmakerMatchedMsg)
	NotificationsHandler        func(context.Context, *NotificationsMsg)
	StatusPresenceEventHandler  func(context.Context, *StatusPresenceEventMsg)
	StreamDataHandler           func(context.Context, *StreamDataMsg)
	StreamPresenceEventHandler  func(context.Context, *StreamPresenceEventMsg)

	rw sync.RWMutex
}

// NewConn creates a new nakama realtime websocket connection.
func NewConn(ctx context.Context, opts ...ConnOption) (*Conn, error) {
	conn := &Conn{
		binary:            true,
		query:             url.Values{},
		backoffMin:        20 * time.Millisecond,
		backoffMax:        3 * time.Second,
		backoffMultiplier: 1.2,
		out:               make(chan *res),
		m:                 make(map[string]*res),
		stop:              true,
	}
	for _, o := range opts {
		o(conn)
	}
	if err := conn.Open(ctx); err != nil {
		return nil, err
	}
	return conn, nil
}

// Open opens and persists (when enabled) the websocket connection to the
// Nakama server.
func (conn *Conn) Open(ctx context.Context) error {
	if conn.Connected() {
		return nil
	}
	conn.rw.Lock()
	conn.stop = false
	conn.rw.Unlock()
	if !conn.persist {
		return conn.open(ctx)
	}
	go conn.run(ctx)
	return nil
}

// run keeps open the websocket connection to the Nakama server when persist is
// enabled.
func (conn *Conn) run(ctx context.Context) {
	for d, last := conn.backoffMin, true; !conn.stop; d = min(time.Duration(float64(d)*conn.backoffMultiplier), conn.backoffMax) {
		connected := conn.Connected()
		if last != connected {
			d = conn.backoffMin
		}
		last = connected
		if connected {
			select {
			case <-ctx.Done():
				return
			case <-time.After(d):
				continue
			}
		}
		if err := conn.open(ctx); err == nil {
			continue
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(d):
		}
	}
}

// open opens the websocket connection to the Nakama server.
func (conn *Conn) open(ctx context.Context) error {
	ws, err := conn.dial(ctx)
	if err != nil {
		return err
	}
	conn.rw.Lock()
	defer conn.rw.Unlock()
	ctx, cancel := context.WithCancel(ctx)
	conn.ctx, conn.ws, conn.cancel = ctx, ws, cancel
	if conn.ConnectHandler != nil {
		go conn.ConnectHandler(conn.ctx)
	}
	// incoming
	go func() {
		for {
			_, r, err := ws.Reader(ctx)
			if err != nil {
				_ = conn.CloseWithErr(err)
				return
			}
			buf, err := io.ReadAll(r)
			if err != nil {
				_ = conn.CloseWithErr(err)
				return
			}
			if buf == nil {
				_ = conn.CloseWithErr(ErrConnReadEmptyMessage)
				return
			}
			if err := conn.recv(ctx, buf); err != nil {
				conn.h.Errf("unable to dispatch incoming message: %v", err)
			}
		}
	}()
	// outgoing
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case m := <-conn.out:
				id, err := conn.send(ctx, ws, m.msg)
				if err != nil {
					if !errors.Is(err, context.Canceled) {
						conn.h.Errf("unable to send message: %v", err)
					}
					m.err <- fmt.Errorf("unable to send message: %w", err)
					close(m.err)
					continue
				}
				if m.v == nil || id == "" {
					close(m.err)
					continue
				}
				conn.rw.Lock()
				conn.m[id] = m
				conn.rw.Unlock()
			}
		}
	}()
	return nil
}

// send marshals the message and writes it to the websocket connection.
func (conn *Conn) send(ctx context.Context, ws *websocket.Conn, msg EnvelopeBuilder) (string, error) {
	env := msg.BuildEnvelope()
	env.Cid = strconv.FormatUint(atomic.AddUint64(&conn.id, 1), 10)
	buf, err := conn.marshal(env)
	if err != nil {
		return "", err
	}
	typ := websocket.MessageBinary
	if !conn.binary {
		typ = websocket.MessageText
	}
	if err := ws.Write(ctx, typ, buf); err != nil {
		_ = conn.CloseWithErr(err)
		return "", err
	}
	return env.Cid, nil
}

// recv unmarshals buf, dispatching the message.
func (conn *Conn) recv(ctx context.Context, buf []byte) error {
	env, err := conn.unmarshal(buf)
	switch {
	case err != nil:
		return fmt.Errorf("unable to unmarshal: %w", err)
	case env.Cid == "":
		return conn.recvNotify(ctx, env)
	}
	return conn.recvResponse(env)
}

// recvNotify dispaches events and received updates.
func (conn *Conn) recvNotify(ctx context.Context, env *Envelope) error {
	switch v := env.Message.(type) {
	case *Envelope_Error:
		if conn.ErrorHandler != nil {
			go conn.ErrorHandler(ctx, v.Error)
		}
		return v.Error
	case *Envelope_ChannelMessage:
		if conn.ChannelMessageHandler != nil {
			go conn.ChannelMessageHandler(ctx, v.ChannelMessage)
		}
		return nil
	case *Envelope_ChannelPresenceEvent:
		if conn.ChannelPresenceEventHandler != nil {
			go conn.ChannelPresenceEventHandler(ctx, v.ChannelPresenceEvent)
		}
		return nil
	case *Envelope_MatchData:
		if conn.MatchDataHandler != nil {
			go conn.MatchDataHandler(ctx, v.MatchData)
		}
		return nil
	case *Envelope_MatchPresenceEvent:
		if conn.MatchPresenceEventHandler != nil {
			go conn.MatchPresenceEventHandler(ctx, v.MatchPresenceEvent)
		}
		return nil
	case *Envelope_MatchmakerMatched:
		if conn.MatchmakerMatchedHandler != nil {
			go conn.MatchmakerMatchedHandler(ctx, v.MatchmakerMatched)
		}
		return nil
	case *Envelope_Notifications:
		if conn.NotificationsHandler != nil {
			go conn.NotificationsHandler(ctx, v.Notifications)
		}
		return nil
	case *Envelope_StatusPresenceEvent:
		if conn.StatusPresenceEventHandler != nil {
			go conn.StatusPresenceEventHandler(ctx, v.StatusPresenceEvent)
		}
		return nil
	case *Envelope_StreamData:
		if conn.StreamDataHandler != nil {
			go conn.StreamDataHandler(ctx, v.StreamData)
		}
		return nil
	case *Envelope_StreamPresenceEvent:
		if conn.StreamPresenceEventHandler != nil {
			go conn.StreamPresenceEventHandler(ctx, v.StreamPresenceEvent)
		}
		return nil
	}
	return fmt.Errorf("unknown type %T", env.Message)
}

// recvResponse dispatches a received response (messages with cid != "").
func (conn *Conn) recvResponse(env *Envelope) error {
	conn.rw.RLock()
	m, ok := conn.m[env.Cid]
	conn.rw.RUnlock()
	if !ok || m == nil {
		return fmt.Errorf("no callback id %s (%T)", env.Cid, env.Message)
	}
	// remove and close
	defer func() {
		close(m.err)
		conn.rw.Lock()
		delete(conn.m, env.Cid)
		conn.rw.Unlock()
	}()
	// check error
	if err, ok := env.Message.(*Envelope_Error); ok {
		conn.h.Errf("realtime error: %v", err.Error)
		m.err <- err.Error
		return nil
	}
	// ignore response for RPC
	if m.v == nil {
		return nil
	}
	// merge
	proto.Merge(m.v.BuildEnvelope(), env)
	return nil
}

// Send sends a message.
func (conn *Conn) Send(ctx context.Context, msg, v EnvelopeBuilder) error {
	m := &res{
		msg: msg,
		v:   v,
		err: make(chan error, 1),
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case conn.out <- m:
	}
	var err error
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-m.err:
	}
	return err
}

// Connected returns true when the websocket connection is connected to the
// Nakama server.
func (conn *Conn) Connected() bool {
	ws := conn.ws
	return ws != nil
}

// CloseWithErr closes the websocket connection with an error.
func (conn *Conn) CloseWithErr(err error) error {
	conn.rw.Lock()
	defer conn.rw.Unlock()
	if conn.ws != nil {
		defer conn.ws.Close(websocket.StatusGoingAway, "going away")
		defer conn.cancel()
		for k := range conn.m {
			delete(conn.m, k)
		}
		if conn.DisconnectHandler != nil {
			go conn.DisconnectHandler(conn.ctx, err)
		}
		conn.stop, conn.ctx, conn.ws, conn.cancel = true, nil, nil, nil
	}
	return nil
}

// Close closes the websocket connection.
func (conn *Conn) Close() error {
	return conn.CloseWithErr(nil)
}

// dial creates a new websocket connection to the Nakama server.
func (conn *Conn) dial(ctx context.Context) (*websocket.Conn, error) {
	urlstr, opts, err := conn.dialParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create dial params: %w", err)
	}
	ws, _, err := websocket.Dial(ctx, urlstr, opts)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %s: %w", urlstr, err)
	}
	return ws, nil
}

// dialParams builds the dial parameters for the nakama server.
func (conn *Conn) dialParams(ctx context.Context) (string, *websocket.DialOptions, error) {
	// build url
	urlstr := conn.url
	if urlstr == "" && conn.h != nil {
		var err error
		if urlstr, err = conn.h.SocketURL(); err != nil {
			return "", nil, err
		}
	}
	// build token
	token := conn.token
	if token == "" && conn.h != nil {
		var err error
		if token, err = conn.h.Token(ctx); err != nil {
			return "", nil, err
		}
	}
	// build query
	query := url.Values{}
	for k, v := range conn.query {
		query[k] = v
	}
	query.Set("token", token)
	format := "protobuf"
	if !conn.binary {
		format = "json"
	}
	query.Set("format", format)
	httpClient := http.DefaultClient
	if conn.h != nil {
		httpClient = conn.h.HttpClient()
	}
	return urlstr + "?" + query.Encode(), buildWsOptions(httpClient), nil
}

// marshal marshals the message. If the format set on the connection is json,
// then the message will be marshaled using json encoding.
func (conn *Conn) marshal(env *Envelope) ([]byte, error) {
	f := proto.Marshal
	if !conn.binary {
		f = protojson.Marshal
	}
	return f(env)
}

// unmarshal unmarshals the message. If the format set on the connection is
// json, then v will be unmarshaled using json encoding.
func (conn *Conn) unmarshal(buf []byte) (*Envelope, error) {
	f := proto.Unmarshal
	if !conn.binary {
		f = protojson.Unmarshal
	}
	env := new(Envelope)
	if err := f(buf, env); err != nil {
		return nil, err
	}
	return env, nil
}

// ChannelJoin sends a message to join a chat channel.
func (conn *Conn) ChannelJoin(ctx context.Context, target string, typ ChannelType, persistence, hidden bool) (*ChannelMsg, error) {
	return ChannelJoin(target, typ).
		WithPersistence(persistence).
		WithHidden(hidden).
		Send(ctx, conn)
}

// ChannelJoinAsync sends a message to join a chat channel.
func (conn *Conn) ChannelJoinAsync(ctx context.Context, target string, typ ChannelType, persistence, hidden bool, f func(*ChannelMsg, error)) {
	ChannelJoin(target, typ).
		WithPersistence(persistence).
		WithHidden(hidden).
		Async(ctx, conn, f)
}

// ChannelLeave sends a message to leave a chat channel.
func (conn *Conn) ChannelLeave(ctx context.Context, channelId string) error {
	return ChannelLeave(channelId).Send(ctx, conn)
}

// ChannelLeaveAsync sends a message to leave a chat channel.
func (conn *Conn) ChannelLeaveAsync(ctx context.Context, channelId string, f func(error)) {
	ChannelLeave(channelId).Async(ctx, conn, f)
}

// ChannelMessageRemove sends a message to remove a message from a channel.
func (conn *Conn) ChannelMessageRemove(ctx context.Context, channelId, messageId string) (*ChannelMessageAckMsg, error) {
	return ChannelMessageRemove(channelId, messageId).Send(ctx, conn)
}

// ChannelMessageRemoveAsync sends a message to remove a message from a channel.
func (conn *Conn) ChannelMessageRemoveAsync(ctx context.Context, channelId, messageId string, f func(*ChannelMessageAckMsg, error)) {
	ChannelMessageRemove(channelId, messageId).Async(ctx, conn, f)
}

// ChannelMessageSend sends a message on a channel.
func (conn *Conn) ChannelMessageSend(ctx context.Context, channelId, content string) (*ChannelMessageAckMsg, error) {
	return ChannelMessageSend(channelId, content).Send(ctx, conn)
}

// ChannelMessageSendAsync sends a message on a channel.
func (conn *Conn) ChannelMessageSendAsync(ctx context.Context, channelId, content string, f func(*ChannelMessageAckMsg, error)) {
	ChannelMessageSend(channelId, content).Async(ctx, conn, f)
}

// ChannelMessageUpdate sends a message to update a message on a channel.
func (conn *Conn) ChannelMessageUpdate(ctx context.Context, channelId, messageId, content string) (*ChannelMessageAckMsg, error) {
	return ChannelMessageUpdate(channelId, messageId, content).Send(ctx, conn)
}

// ChannelMessageUpdateAsync sends a message to update a message on a channel.
func (conn *Conn) ChannelMessageUpdateAsync(ctx context.Context, channelId, messageId, content string, f func(*ChannelMessageAckMsg, error)) {
	ChannelMessageUpdate(channelId, messageId, content).Async(ctx, conn, f)
}

// MatchCreate sends a message to create a multiplayer match.
func (conn *Conn) MatchCreate(ctx context.Context, name string) (*MatchMsg, error) {
	return MatchCreate(name).Send(ctx, conn)
}

// MatchCreateAsync sends a message to create a multiplayer match.
func (conn *Conn) MatchCreateAsync(ctx context.Context, name string, f func(*MatchMsg, error)) {
	MatchCreate(name).Async(ctx, conn, f)
}

// MatchJoin sends a message to join a match.
func (conn *Conn) MatchJoin(ctx context.Context, matchId string, metadata map[string]string) (*MatchMsg, error) {
	return MatchJoin(matchId).
		WithMetadata(metadata).
		Send(ctx, conn)
}

// MatchJoinAsync sends a message to join a match.
func (conn *Conn) MatchJoinAsync(ctx context.Context, matchId string, metadata map[string]string, f func(*MatchMsg, error)) {
	MatchJoin(matchId).
		WithMetadata(metadata).
		Async(ctx, conn, f)
}

// MatchJoinToken sends a message to join a match with a token.
func (conn *Conn) MatchJoinToken(ctx context.Context, token string, metadata map[string]string) (*MatchMsg, error) {
	return MatchJoinToken(token).
		WithMetadata(metadata).
		Send(ctx, conn)
}

// MatchJoinTokenAsync sends a message to join a match with a token.
func (conn *Conn) MatchJoinTokenAsync(ctx context.Context, token string, metadata map[string]string, f func(*MatchMsg, error)) {
	MatchJoinToken(token).
		WithMetadata(metadata).
		Async(ctx, conn, f)
}

// MatchLeave sends a message to leave a multiplayer match.
func (conn *Conn) MatchLeave(ctx context.Context, matchId string) error {
	return MatchLeave(matchId).Send(ctx, conn)
}

// MatchLeaveAsync sends a message to leave a multiplayer match.
func (conn *Conn) MatchLeaveAsync(ctx context.Context, matchId string, f func(error)) {
	MatchLeave(matchId).Async(ctx, conn, f)
}

// MatchmakerAdd sends a message to join the matchmaker pool and search for opponents on the server.
func (conn *Conn) MatchmakerAdd(ctx context.Context, msg *MatchmakerAddMsg) (*MatchmakerTicketMsg, error) {
	return msg.Send(ctx, conn)
}

// MatchmakerAddAsync sends a message to join the matchmaker pool and search for opponents on the server.
func (conn *Conn) MatchmakerAddAsync(ctx context.Context, msg *MatchmakerAddMsg, f func(*MatchmakerTicketMsg, error)) {
	msg.Async(ctx, conn, f)
}

// MatchmakerRemove sends a message to leave the matchmaker pool for a ticket.
func (conn *Conn) MatchmakerRemove(ctx context.Context, ticket string) error {
	return MatchmakerRemove(ticket).Send(ctx, conn)
}

// MatchmakerRemoveAsync sends a message to leave the matchmaker pool for a ticket.
func (conn *Conn) MatchmakerRemoveAsync(ctx context.Context, ticket string, f func(error)) {
	MatchmakerRemove(ticket).Async(ctx, conn, f)
}

// MatchDataSend sends a message to send input to a multiplayer match.
func (conn *Conn) MatchDataSend(ctx context.Context, matchId string, opCode int64, data []byte, reliable bool, presences ...*UserPresenceMsg) error {
	return MatchDataSend(matchId, opCode, data).
		WithPresences(presences...).
		WithReliable(reliable).
		Send(ctx, conn)
}

// MatchDataSendAsync sends a message to send input to a multiplayer match.
func (conn *Conn) MatchDataSendAsync(ctx context.Context, matchId string, opCode int64, data []byte, reliable bool, presences []*UserPresenceMsg, f func(error)) {
	MatchDataSend(matchId, opCode, data).
		WithPresences(presences...).
		WithReliable(reliable).
		Async(ctx, conn, f)
}

// PartyAccept sends a message to accept a party member.
func (conn *Conn) PartyAccept(ctx context.Context, partyId string, presence *UserPresenceMsg) error {
	return PartyAccept(partyId, presence).Send(ctx, conn)
}

// PartyAcceptAsync sends a message to accept a party member.
func (conn *Conn) PartyAcceptAsync(ctx context.Context, partyId string, presence *UserPresenceMsg, f func(error)) {
	PartyAccept(partyId, presence).Async(ctx, conn, f)
}

// PartyClose sends a message closes a party, kicking all party members.
func (conn *Conn) PartyClose(ctx context.Context, partyId string) error {
	return PartyClose(partyId).Send(ctx, conn)
}

// PartyCloseAsync sends a message closes a party, kicking all party members.
func (conn *Conn) PartyCloseAsync(ctx context.Context, partyId string, f func(error)) {
	PartyClose(partyId).Async(ctx, conn, f)
}

// PartyCreate sends a message to create a party.
func (conn *Conn) PartyCreate(ctx context.Context, open bool, maxSize int) (*PartyMsg, error) {
	return PartyCreate(open, maxSize).Send(ctx, conn)
}

// PartyCreateAsync sends a message to create a party.
func (conn *Conn) PartyCreateAsync(ctx context.Context, open bool, maxSize int, f func(*PartyMsg, error)) {
	PartyCreate(open, maxSize).Async(ctx, conn, f)
}

// PartyDataSend sends a message to send input to a multiplayer party.
func (conn *Conn) PartyDataSend(ctx context.Context, partyId string, opCode OpType, data []byte, reliable bool, presences ...*UserPresenceMsg) error {
	return PartyDataSend(partyId, opCode, data).Send(ctx, conn)
}

// PartyDataSendAsync sends a message to send input to a multiplayer party.
func (conn *Conn) PartyDataSendAsync(ctx context.Context, partyId string, opCode OpType, data []byte, reliable bool, presences []*UserPresenceMsg, f func(error)) {
	PartyDataSend(partyId, opCode, data).Async(ctx, conn, f)
}

// PartyJoin sends a message to join a party.
func (conn *Conn) PartyJoin(ctx context.Context, partyId string) error {
	return PartyJoin(partyId).Send(ctx, conn)
}

// PartyJoinAsync sends a message to join a party.
func (conn *Conn) PartyJoinAsync(ctx context.Context, partyId string, f func(error)) {
	PartyJoin(partyId).Async(ctx, conn, f)
}

// PartyJoinRequests sends a message to request the list of pending join requests for a party.
func (conn *Conn) PartyJoinRequests(ctx context.Context, partyId string) (*PartyJoinRequestMsg, error) {
	return PartyJoinRequests(partyId).Send(ctx, conn)
}

// PartyJoinRequestsAsync sends a message to request the list of pending join requests for a party.
func (conn *Conn) PartyJoinRequestsAsync(ctx context.Context, partyId string, f func(*PartyJoinRequestMsg, error)) {
	PartyJoinRequests(partyId).Async(ctx, conn, f)
}

// PartyLeave sends a message to leave a party.
func (conn *Conn) PartyLeave(ctx context.Context, partyId string) error {
	return PartyLeave(partyId).Send(ctx, conn)
}

// PartyLeaveAsync sends a message to leave a party.
func (conn *Conn) PartyLeaveAsync(ctx context.Context, partyId string, f func(error)) {
	PartyLeave(partyId).Async(ctx, conn, f)
}

// PartyMatchmakerAdd sends a message to begin matchmaking as a party.
func (conn *Conn) PartyMatchmakerAdd(ctx context.Context, partyId, query string, minCount, maxCount int) (*PartyMatchmakerTicketMsg, error) {
	return PartyMatchmakerAdd(partyId, query, minCount, maxCount).Send(ctx, conn)
}

// PartyMatchmakerAddAsync sends a message to begin matchmaking as a party.
func (conn *Conn) PartyMatchmakerAddAsync(ctx context.Context, partyId, query string, minCount, maxCount int, f func(*PartyMatchmakerTicketMsg, error)) {
	PartyMatchmakerAdd(partyId, query, minCount, maxCount).Async(ctx, conn, f)
}

// PartyMatchmakerRemove sends a message to cancel a party matchmaking process for a ticket.
func (conn *Conn) PartyMatchmakerRemove(ctx context.Context, partyId, ticket string) error {
	return PartyMatchmakerRemove(partyId, ticket).Send(ctx, conn)
}

// PartyMatchmakerRemoveAsync sends a message to cancel a party matchmaking process for a ticket.
func (conn *Conn) PartyMatchmakerRemoveAsync(ctx context.Context, partyId, ticket string, f func(error)) {
	PartyMatchmakerRemove(partyId, ticket).Async(ctx, conn, f)
}

// PartyPromote sends a message to promote a new party leader.
func (conn *Conn) PartyPromote(ctx context.Context, partyId string, presence *UserPresenceMsg) (*PartyLeaderMsg, error) {
	return PartyPromote(partyId, presence).Send(ctx, conn)
}

// PartyPromoteAsync sends a message to promote a new party leader.
func (conn *Conn) PartyPromoteAsync(ctx context.Context, partyId string, presence *UserPresenceMsg, f func(*PartyLeaderMsg, error)) {
	PartyPromote(partyId, presence).Async(ctx, conn, f)
}

// PartyRemove sends a message to kick a party member or decline a request to join.
func (conn *Conn) PartyRemove(ctx context.Context, partyId string, presence *UserPresenceMsg) error {
	return PartyRemove(partyId, presence).Send(ctx, conn)
}

// PartyRemoveAsync sends a message to kick a party member or decline a request to join.
func (conn *Conn) PartyRemoveAsync(ctx context.Context, partyId string, presence *UserPresenceMsg, f func(error)) {
	PartyRemove(partyId, presence).Async(ctx, conn, f)
}

// Ping sends a message to do a ping.
func (conn *Conn) Ping(ctx context.Context) error {
	return Ping().Send(ctx, conn)
}

// PingAsync sends a message to do a ping.
func (conn *Conn) PingAsync(ctx context.Context, f func(error)) {
	Ping().Async(ctx, conn, f)
}

// Rpc sends a message to execute a remote procedure call.
func (conn *Conn) Rpc(ctx context.Context, id string, payload, v interface{}) error {
	return Rpc(id, payload, v).Send(ctx, conn)
}

// RpcAsync sends a message to execute a remote procedure call.
func (conn *Conn) RpcAsync(ctx context.Context, id string, payload, v interface{}, f func(error)) {
	Rpc(id, payload, v).SendAsync(ctx, conn, f)
}

// StatusFollow sends a message to subscribe to user status updates.
func (conn *Conn) StatusFollow(ctx context.Context, userIds ...string) (*StatusMsg, error) {
	return StatusFollow(userIds...).Send(ctx, conn)
}

// StatusFollowAsync sends a message to subscribe to user status updates.
func (conn *Conn) StatusFollowAsync(ctx context.Context, userIds []string, f func(*StatusMsg, error)) {
	StatusFollow(userIds...).Async(ctx, conn, f)
}

// StatusUnfollow sends a message to unfollow user's status updates.
func (conn *Conn) StatusUnfollow(ctx context.Context, userIds ...string) error {
	return StatusUnfollow(userIds...).Send(ctx, conn)
}

// StatusUnfollowAsync sends a message to unfollow user's status updates.
func (conn *Conn) StatusUnfollowAsync(ctx context.Context, userIds []string, f func(error)) {
	StatusUnfollow(userIds...).Async(ctx, conn, f)
}

// StatusUpdate sends a message to update the user's status.
func (conn *Conn) StatusUpdate(ctx context.Context, status string) error {
	return StatusUpdate().
		WithStatus(status).
		Send(ctx, conn)
}

// StatusUpdateAsync sends a message to update the user's status.
func (conn *Conn) StatusUpdateAsync(ctx context.Context, status string, f func(error)) {
	StatusUpdate().
		WithStatus(status).
		Async(ctx, conn, f)
}

// res wraps a request and results.
type res struct {
	msg EnvelopeBuilder
	v   EnvelopeBuilder
	err chan error
}

// ConnOption is a nakama realtime websocket connection option.
type ConnOption func(*Conn)

// WithConnClientHandler is a nakama websocket connection option to set the
// ClientHandler used.
func WithConnClientHandler(h ClientHandler) ConnOption {
	return func(conn *Conn) {
		conn.h = h
	}
}

// WithConnUrl is a nakama websocket connection option to set the websocket
// URL.
func WithConnUrl(urlstr string) ConnOption {
	return func(conn *Conn) {
		conn.url = urlstr
	}
}

// WithConnToken is a nakama websocket connection option to set the auth token
// for the websocket.
func WithConnToken(token string) ConnOption {
	return func(conn *Conn) {
		conn.token = token
	}
}

// WithConnFormat is a nakama websocket connection option to set the message
// encoding format (either "json" or "protobuf").
func WithConnFormat(format string) ConnOption {
	return func(conn *Conn) {
		switch s := strings.ToLower(format); s {
		case "protobuf":
			conn.binary = true
		case "json":
			conn.binary = false
		default:
			panic(fmt.Sprintf("invalid websocket format %q", format))
		}
	}
}

// WithConnQuery is a nakama websocket connection option to add an additional
// key/value query param on the websocket URL.
//
// Note: this should not be used to set "token" or "format". Use WithConnToken
// and WithConnFormat, respectively, to change the token and format query
// params.
func WithConnQuery(key, value string) ConnOption {
	return func(conn *Conn) {
		conn.query.Set(key, value)
	}
}

// WithConnLang is a nakama websocket connection option to set the lang query
// param on the websocket URL.
func WithConnLang(lang string) ConnOption {
	return func(conn *Conn) {
		conn.query.Set("lang", lang)
	}
}

// WithConnCreateStatus is a nakama websocket connection option to set the
// status query param on the websocket URL.
func WithConnCreateStatus(status bool) ConnOption {
	return func(conn *Conn) {
		conn.query.Set("status", strconv.FormatBool(status))
	}
}

// WithConnPersist is a nakama websocket connection option to enable keeping
// open a persistent connection to the Nakama server.
func WithConnPersist(persist bool) ConnOption {
	return func(conn *Conn) {
		conn.persist = persist
	}
}

// WithConnBackoff is a nakama websocket connection option to set the
// connection backoff (retry) settings.
func WithConnBackoff(backoffMin, backoffMax time.Duration, backoffMultiplier float64) ConnOption {
	return func(conn *Conn) {
		conn.backoffMin, conn.backoffMax, conn.backoffMultiplier = backoffMin, backoffMax, backoffMultiplier
	}
}

// WithConnHandler is a nakama websocket connection option to set the
// connection's message handlers. See the ConnHandler type for documentation on
// supported interfaces.
//
// WithConnHandler works by "smuggling" interfaces. That is, WithConnHandler
// checks via a type cast when the underlying type supports methods of the
// following format:
//
//	interface{
//		<MessageType>Handler(context.Context, *<MessageType>Msg)
//	}
//
// If the ConnHandler's underlying type supports the above, then the
// ConnHandler's <MessageType>Handler method will be set as
// Conn.<MessageType>Handler. For example, given the following:
//
//	type MyClient struct{}
//
//	func (cl *MyClient) MatchDataHandler(context.Context, *nakama.MatchDataMsg) {}
//	func (cl *MyClient) NotificationsHandler(context.Context, *nakama.NotificationsMsg) {}
//
// The following:
//
//	cl := nakama.New(/* ... */)
//	myClient := &MyClient{}
//	conn, err := cl.NewConn(ctx, nakama.WithConnHandler(myClient))
//
// Is equivalent to:
//
//	cl := nakama.New(/* ... */)
//	myClient := &MyClient{}
//	conn, err := cl.NewConn(ctx)
//	conn.MatchDataHandler = myClient.MatchDataHandler
//	conn.NotificationsHandler = myClient.NotificationsHandler
//
// For an overview of Go interface smuggling as a concept, see:
//
// https://utcc.utoronto.ca/~cks/space/blog/programming/GoInterfaceSmuggling
func WithConnHandler(handler ConnHandler) ConnOption {
	return func(conn *Conn) {
		if x, ok := handler.(interface {
			ConnectHandler(context.Context)
		}); ok {
			conn.ConnectHandler = x.ConnectHandler
		}
		if x, ok := handler.(interface {
			DisconnectHandler(context.Context, error)
		}); ok {
			conn.DisconnectHandler = x.DisconnectHandler
		}
		if x, ok := handler.(interface {
			ErrorHandler(context.Context, *ErrorMsg)
		}); ok {
			conn.ErrorHandler = x.ErrorHandler
		}
		if x, ok := handler.(interface {
			ChannelMessageHandler(context.Context, *ChannelMessageMsg)
		}); ok {
			conn.ChannelMessageHandler = x.ChannelMessageHandler
		}
		if x, ok := handler.(interface {
			ChannelPresenceEventHandler(context.Context, *ChannelPresenceEventMsg)
		}); ok {
			conn.ChannelPresenceEventHandler = x.ChannelPresenceEventHandler
		}
		if x, ok := handler.(interface {
			MatchDataHandler(context.Context, *MatchDataMsg)
		}); ok {
			conn.MatchDataHandler = x.MatchDataHandler
		}
		if x, ok := handler.(interface {
			MatchPresenceEventHandler(context.Context, *MatchPresenceEventMsg)
		}); ok {
			conn.MatchPresenceEventHandler = x.MatchPresenceEventHandler
		}
		if x, ok := handler.(interface {
			MatchmakerMatchedHandler(context.Context, *MatchmakerMatchedMsg)
		}); ok {
			conn.MatchmakerMatchedHandler = x.MatchmakerMatchedHandler
		}
		if x, ok := handler.(interface {
			NotificationsHandler(context.Context, *NotificationsMsg)
		}); ok {
			conn.NotificationsHandler = x.NotificationsHandler
		}
		if x, ok := handler.(interface {
			StatusPresenceEventHandler(context.Context, *StatusPresenceEventMsg)
		}); ok {
			conn.StatusPresenceEventHandler = x.StatusPresenceEventHandler
		}
		if x, ok := handler.(interface {
			StreamDataHandler(context.Context, *StreamDataMsg)
		}); ok {
			conn.StreamDataHandler = x.StreamDataHandler
		}
		if x, ok := handler.(interface {
			StreamPresenceEventHandler(context.Context, *StreamPresenceEventMsg)
		}); ok {
			conn.StreamPresenceEventHandler = x.StreamPresenceEventHandler
		}
	}
}

// ConnError is a websocket connection error.
type ConnError string

const (
	// ErrConnAlreadyOpen is the conn already open error.
	ErrConnAlreadyOpen ConnError = "conn already open"
	// ErrConnReadEmptyMessage is the conn read empty message error.
	ErrConnReadEmptyMessage ConnError = "conn read empty message"
)

// Error satisfies the error interface.
func (err ConnError) Error() string {
	return string(err)
}

// min returns the minimum of a, b.
func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
