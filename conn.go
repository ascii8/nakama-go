package nakama

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/heroiclabs/nakama-common/rtapi"
	"golang.org/x/exp/maps"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

// Handler is the interface for connection handlers.
type Handler interface {
	HttpClient() *http.Client
	SocketURL() (string, error)
	Token(context.Context) (string, error)
	Logf(string, ...interface{})
	Errf(string, ...interface{})
}

// Conn is a nakama realtime websocket connection.
type Conn struct {
	h      Handler
	url    string
	token  string
	binary bool
	query  url.Values
	conn   *websocket.Conn
	cancel func()
	out    chan *req
	in     chan []byte
	l      map[string]*req
	rw     sync.RWMutex
	id     uint64
}

// NewConn creates a new nakama realtime websocket connection.
func NewConn(ctx context.Context, opts ...ConnOption) (*Conn, error) {
	conn := &Conn{
		binary: true,
		query:  url.Values{},
		out:    make(chan *req),
		in:     make(chan []byte),
		l:      make(map[string]*req),
	}
	for _, o := range opts {
		o(conn)
	}
	// build url
	urlstr := conn.url
	if urlstr == "" && conn.h != nil {
		var err error
		if urlstr, err = conn.h.SocketURL(); err != nil {
			return nil, err
		}
	}
	// build token
	token := conn.token
	if token == "" && conn.h != nil {
		var err error
		if token, err = conn.h.Token(ctx); err != nil {
			return nil, err
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
	// open socket
	var err error
	conn.conn, _, err = websocket.Dial(ctx, urlstr+"?"+query.Encode(), &websocket.DialOptions{
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to open nakama websocket %s: %w", urlstr, err)
	}
	// run
	ctx, conn.cancel = context.WithCancel(ctx)
	go conn.run(ctx)
	return conn, nil
}

// Next returns the next id.
func (conn *Conn) Next() string {
	return strconv.FormatUint(atomic.AddUint64(&conn.id, 1), 10)
}

// Marshal marshals the message. If the format set on the connection is json,
// then the message will be marshaled using json encoding.
func (conn *Conn) Marshal(id string, msg Message) ([]byte, error) {
	f := proto.Marshal
	if !conn.binary {
		f = protojson.Marshal
	}
	return f(msg.BuildEnv(id))
}

// Unmarshal unmarshals the message. If the format set on the connection is
// json, then v will be unmarshaled using json encoding.
func (conn *Conn) Unmarshal(buf []byte, v proto.Message) error {
	f := proto.Unmarshal
	if !conn.binary {
		f = protojson.Unmarshal
	}
	return f(buf, v)
}

// run handles incoming and outgoing websocket messages.
func (conn *Conn) run(ctx context.Context) {
	// read incoming
	go func() {
		for {
			select {
			case <-ctx.Done():
			default:
			}
			_, r, err := conn.conn.Reader(ctx)
			switch {
			case err != nil && (errors.Is(err, context.Canceled) || errors.As(err, &websocket.CloseError{})):
				return
			case err != nil:
				conn.h.Errf("reader error: %v", err)
				continue
			}
			buf, err := ioutil.ReadAll(r)
			if err != nil {
				conn.h.Errf("unable to read message: %v", err)
				continue
			}
			conn.in <- buf
		}
	}()
	// dispatch outgoing/incoming
	for {
		select {
		case <-ctx.Done():
			return
		case m := <-conn.out:
			if m == nil {
				continue
			}
			id, err := conn.Send(ctx, m.msg)
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
			conn.l[id] = m
			conn.rw.Unlock()
		case buf := <-conn.in:
			if buf == nil {
				continue
			}
			if err := conn.Recv(buf); err != nil {
				conn.h.Errf("unable to dispatch message: %v", err)
				continue
			}
		}
	}
}

// Send marshals the message and writes it to the websocket connection.
func (conn *Conn) Send(ctx context.Context, msg Message) (string, error) {
	id := conn.Next()
	buf, err := conn.Marshal(id, msg)
	if err != nil {
		return "", err
	}
	typ := websocket.MessageBinary
	if !conn.binary {
		typ = websocket.MessageText
	}
	if err := conn.conn.Write(ctx, typ, buf); err != nil {
		return "", err
	}
	return id, nil
}

// Recv unmarshals buf, dispatching the message.
func (conn *Conn) Recv(buf []byte) error {
	env := new(rtapi.Envelope)
	if err := conn.Unmarshal(buf, env); err != nil {
		return err
	}
	if env.Cid != "" {
		return conn.RecvResponse(env)
	}
	return nil
}

// RecvResponse dispatches a received a response.
func (conn *Conn) RecvResponse(env *rtapi.Envelope) error {
	conn.rw.RLock()
	req, ok := conn.l[env.Cid]
	conn.rw.RUnlock()
	if !ok || req == nil {
		return fmt.Errorf("no callback id %s", env.Cid)
	}
	defer func() {
		close(req.err)
		conn.rw.Lock()
		delete(conn.l, env.Cid)
		conn.rw.Unlock()
	}()
	// check error
	if v, ok := env.Message.(*rtapi.Envelope_Error); ok {
		err := NewRealtimeError(v.Error)
		conn.h.Errf("received %s error: %v", env.Cid, err)
		req.err <- err
		return nil
	}
	// merge
	proto.Merge(req.v.BuildEnv(""), env)
	return nil
}

// Do sends a message to the websocket connection, blocking until results are
// received and decoded to v, or an error is encountered. Returns immediately
// after writing the message when v is nil.
func (conn *Conn) Do(ctx context.Context, msg, v Message) error {
	m := &req{
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

// Close closes the websocket connection.
func (conn *Conn) Close() error {
	if conn.cancel != nil {
		defer conn.cancel()
	}
	if conn.conn != nil {
		return conn.conn.Close(websocket.StatusGoingAway, "going away")
	}
	return nil
}

// Ping sends a ping message.
func (conn *Conn) Ping(ctx context.Context) error {
	return Ping().Do(ctx, conn)
}

// PingAsync sends a ping message.
func (conn *Conn) PingAsync(ctx context.Context, f func(error)) {
	Ping().Async(ctx, conn, f)
}

// req wraps a request and results.
type req struct {
	msg Message
	v   Message
	err chan error
}

// RealtimeError wraps a nakama realtime websocket error.
type RealtimeError struct {
	Code    rtapi.Error_Code
	Message string
	Context map[string]string
}

// NewRealtimeError creates a nakama realtime websocket error from an error
// message.
func NewRealtimeError(err *rtapi.Error) error {
	return &RealtimeError{
		Code:    rtapi.Error_Code(err.Code),
		Message: err.Message,
		Context: err.Context,
	}
}

// Error satisfies the error interface.
func (err *RealtimeError) Error() string {
	var s []string
	keys := maps.Keys(err.Context)
	sort.Strings(keys)
	for _, k := range keys {
		s = append(s, k+":"+err.Context[k])
	}
	var extra string
	if len(s) != 0 {
		extra = " <" + strings.Join(s, " ") + ">"
	}
	return fmt.Sprintf("%s (%d): %s%s", err.Code, err.Code, err.Message, extra)
}

// ConnOption is a nakama realtime websocket connection option.
type ConnOption func(*Conn)

// WithConnHandler is a nakama websocket connection option to set the Handler
// used.
func WithConnHandler(h Handler) ConnOption {
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
		case "json":
			conn.binary = false
		default:
			panic(fmt.Sprintf("invalid websocket format %q", format))
		}
	}
}

// WithConnQuery is a nakama websocket connection option to add an additional key/value
// query param on the websocket URL.
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
