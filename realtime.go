package nakama

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

// EnvelopeBuilder is the shared interface for realtime messages.
type EnvelopeBuilder interface {
	BuildEnvelope() *Envelope
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_Channel{
			Channel: msg,
		},
	}
}

// ChannelJoin creates a realtime message to join a chat channel.
func ChannelJoin(target string, typ ChannelType) *ChannelJoinMsg {
	return &ChannelJoinMsg{
		Target: target,
		Type:   typ,
	}
}

// WithPersistence sets the persistence on the message.
func (msg *ChannelJoinMsg) WithPersistence(persistence bool) *ChannelJoinMsg {
	msg.Persistence = wrapperspb.Bool(persistence)
	return msg
}

// WithHidden sets the hidden on the message.
func (msg *ChannelJoinMsg) WithHidden(hidden bool) *ChannelJoinMsg {
	msg.Hidden = wrapperspb.Bool(hidden)
	return msg
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelJoinMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_ChannelJoin{
			ChannelJoin: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *ChannelJoinMsg) Send(ctx context.Context, conn *Conn) (*ChannelMsg, error) {
	res := new(ChannelMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *ChannelJoinMsg) Async(ctx context.Context, conn *Conn, f func(*ChannelMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// ChannelLeave creates a realtime message to leave a chat channel.
func ChannelLeave(channelId string) *ChannelLeaveMsg {
	return &ChannelLeaveMsg{
		ChannelId: channelId,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelLeaveMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_ChannelLeave{
			ChannelLeave: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *ChannelLeaveMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *ChannelLeaveMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// ChannelMessageMsg is a realtime channel message message.
type ChannelMessageMsg = ChannelMessage

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_ChannelMessage{
			ChannelMessage: msg,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageAckMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_ChannelMessageAck{
			ChannelMessageAck: msg,
		},
	}
}

// ChannelMessageRemove creates a realtime message to remove a message from a channel.
func ChannelMessageRemove(channelId, messageId string) *ChannelMessageRemoveMsg {
	return &ChannelMessageRemoveMsg{
		ChannelId: channelId,
		MessageId: messageId,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageRemoveMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_ChannelMessageRemove{
			ChannelMessageRemove: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *ChannelMessageRemoveMsg) Send(ctx context.Context, conn *Conn) (*ChannelMessageAckMsg, error) {
	res := new(ChannelMessageAckMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *ChannelMessageRemoveMsg) Async(ctx context.Context, conn *Conn, f func(*ChannelMessageAckMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// ChannelMessageSend creates a realtime message to send a message on a channel.
func ChannelMessageSend(channelId, content string) *ChannelMessageSendMsg {
	return &ChannelMessageSendMsg{
		ChannelId: channelId,
		Content:   content,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageSendMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_ChannelMessageSend{
			ChannelMessageSend: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *ChannelMessageSendMsg) Send(ctx context.Context, conn *Conn) (*ChannelMessageAckMsg, error) {
	res := new(ChannelMessageAckMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *ChannelMessageSendMsg) Async(ctx context.Context, conn *Conn, f func(*ChannelMessageAckMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// ChannelMessageUpdate creates a realtime message to update a message on a channel.
func ChannelMessageUpdate(channelId, messageId, content string) *ChannelMessageUpdateMsg {
	return &ChannelMessageUpdateMsg{
		ChannelId: channelId,
		MessageId: messageId,
		Content:   content,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageUpdateMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_ChannelMessageUpdate{
			ChannelMessageUpdate: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *ChannelMessageUpdateMsg) Send(ctx context.Context, conn *Conn) (*ChannelMessageAckMsg, error) {
	res := new(ChannelMessageAckMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *ChannelMessageUpdateMsg) Async(ctx context.Context, conn *Conn, f func(*ChannelMessageAckMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelPresenceEventMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_ChannelPresenceEvent{
			ChannelPresenceEvent: msg,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ErrorMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_Error{
			Error: msg,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_Match{
			Match: msg,
		},
	}
}

// Error satisfies the error interface.
func (err *ErrorMsg) Error() string {
	var keys []string
	for key := range err.Context {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var s []string
	for _, k := range keys {
		s = append(s, k+":"+err.Context[k])
	}
	var extra string
	if len(s) != 0 {
		extra = " <" + strings.Join(s, " ") + ">"
	}
	return fmt.Sprintf("realtime socket error %s (%d): %s%s", err.Code, err.Code, err.Message, extra)
}

// MatchCreate creates a realtime message to create a multiplayer match.
func MatchCreate(name string) *MatchCreateMsg {
	return &MatchCreateMsg{
		Name: name,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchCreateMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchCreate{
			MatchCreate: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *MatchCreateMsg) Send(ctx context.Context, conn *Conn) (*MatchMsg, error) {
	res := new(MatchMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *MatchCreateMsg) Async(ctx context.Context, conn *Conn, f func(*MatchMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchDataMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchData{
			MatchData: msg,
		},
	}
}

// MatchDataSend creates a realtime message to send input to a multiplayer match.
func MatchDataSend(matchId string, opCode int64, data []byte) *MatchDataSendMsg {
	return &MatchDataSendMsg{
		MatchId: matchId,
		OpCode:  opCode,
		Data:    data,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchDataSendMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchDataSend{
			MatchDataSend: msg,
		},
	}
}

// WithPresences sets the presences on the message.
func (msg *MatchDataSendMsg) WithPresences(presences ...*UserPresenceMsg) *MatchDataSendMsg {
	msg.Presences = presences
	return msg
}

// WithReliable sets the reliable on the message.
func (msg *MatchDataSendMsg) WithReliable(reliable bool) *MatchDataSendMsg {
	msg.Reliable = reliable
	return msg
}

// Send sends the message to the connection.
func (msg *MatchDataSendMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, nil)
}

// Async sends the message to the connection.
func (msg *MatchDataSendMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// MatchJoin creates a realtime message to join a match.
func MatchJoin(matchId string) *MatchJoinMsg {
	return &MatchJoinMsg{
		Id: &MatchJoinMsg_MatchId{
			MatchId: matchId,
		},
	}
}

// MatchJoinToken creates a new realtime to join a match with a token.
func MatchJoinToken(token string) *MatchJoinMsg {
	return &MatchJoinMsg{
		Id: &MatchJoinMsg_Token{
			Token: token,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchJoinMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchJoin{
			MatchJoin: msg,
		},
	}
}

// WithMetadata sets the metadata on the message.
func (msg *MatchJoinMsg) WithMetadata(metadata map[string]string) *MatchJoinMsg {
	msg.Metadata = metadata
	return msg
}

// Send sends the message to the connection.
func (msg *MatchJoinMsg) Send(ctx context.Context, conn *Conn) (*MatchMsg, error) {
	res := new(MatchMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *MatchJoinMsg) Async(ctx context.Context, conn *Conn, f func(*MatchMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// MatchLeave creates a realtime message to leave a multiplayer match.
func MatchLeave(matchId string) *MatchLeaveMsg {
	return &MatchLeaveMsg{
		MatchId: matchId,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchLeaveMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchLeave{
			MatchLeave: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *MatchLeaveMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *MatchLeaveMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchPresenceEventMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchPresenceEvent{
			MatchPresenceEvent: msg,
		},
	}
}

// MatchmakerAdd creates a realtime message to join the matchmaker pool and search for opponents on the server.
func MatchmakerAdd(query string, minCount, maxCount int) *MatchmakerAddMsg {
	return &MatchmakerAddMsg{
		Query:    query,
		MinCount: int32(minCount),
		MaxCount: int32(maxCount),
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchmakerAddMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchmakerAdd{
			MatchmakerAdd: msg,
		},
	}
}

// WithStringProperties sets the stringProperties on the message.
func (msg *MatchmakerAddMsg) WithStringProperties(stringProperties map[string]string) *MatchmakerAddMsg {
	msg.StringProperties = stringProperties
	return msg
}

// WithNumericProperties sets the stringProperties on the message.
func (msg *MatchmakerAddMsg) WithNumericProperties(numericProperties map[string]float64) *MatchmakerAddMsg {
	msg.NumericProperties = numericProperties
	return msg
}

// WithCountMultiple sets the stringProperties on the message.
func (msg *MatchmakerAddMsg) WithCountMultiple(countMultiple int) *MatchmakerAddMsg {
	msg.CountMultiple = wrapperspb.Int32(int32(countMultiple))
	return msg
}

// Send sends the message to the connection.
func (msg *MatchmakerAddMsg) Send(ctx context.Context, conn *Conn) (*MatchmakerTicketMsg, error) {
	res := new(MatchmakerTicketMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *MatchmakerAddMsg) Async(ctx context.Context, conn *Conn, f func(*MatchmakerTicketMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchmakerMatchedMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchmakerMatched{
			MatchmakerMatched: msg,
		},
	}
}

// MatchmakerRemove creates a realtime message to leave the matchmaker pool for a ticket.
func MatchmakerRemove(ticket string) *MatchmakerRemoveMsg {
	return &MatchmakerRemoveMsg{
		Ticket: ticket,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchmakerRemoveMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchmakerRemove{
			MatchmakerRemove: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *MatchmakerRemoveMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *MatchmakerRemoveMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchmakerTicketMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_MatchmakerTicket{
			MatchmakerTicket: msg,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *NotificationsMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_Notifications{
			Notifications: msg,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_Party{
			Party: msg,
		},
	}
}

// PartyAccept creates a realtime message to accept a party member.
func PartyAccept(partyId string, presence *UserPresenceMsg) *PartyAcceptMsg {
	return &PartyAcceptMsg{
		PartyId:  partyId,
		Presence: presence,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyAcceptMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyAccept{
			PartyAccept: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyAcceptMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *PartyAcceptMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// PartyClose creates a realtime message to close a party, kicking all party members.
func PartyClose(partyId string) *PartyCloseMsg {
	return &PartyCloseMsg{
		PartyId: partyId,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyCloseMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyClose{
			PartyClose: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyCloseMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *PartyCloseMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// PartyCreate creates a realtime message to create a party.
func PartyCreate(open bool, maxSize int) *PartyCreateMsg {
	return &PartyCreateMsg{
		Open:    open,
		MaxSize: int32(maxSize),
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyCreateMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyCreate{
			PartyCreate: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyCreateMsg) Send(ctx context.Context, conn *Conn) (*PartyMsg, error) {
	res := new(PartyMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *PartyCreateMsg) Async(ctx context.Context, conn *Conn, f func(*PartyMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// PartyDataSend creates a realtime message to send data to a party.
func PartyDataSend(partyId string, opCode OpType, data []byte) *PartyDataSendMsg {
	return &PartyDataSendMsg{
		PartyId: partyId,
		OpCode:  int64(opCode),
		Data:    data,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyDataSendMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyDataSend{
			PartyDataSend: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyDataSendMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *PartyDataSendMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// PartyJoin creates a realtime message to join a party.
func PartyJoin(partyId string) *PartyJoinMsg {
	return &PartyJoinMsg{
		PartyId: partyId,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyJoinMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyJoin{
			PartyJoin: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyJoinMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *PartyJoinMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// PartyJoinRequests creates a realtime message to request the list of pending join requests for a party.
func PartyJoinRequests(partyId string) *PartyJoinRequestsMsg {
	return &PartyJoinRequestsMsg{
		PartyId: partyId,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyJoinRequestsMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyJoinRequestList{
			PartyJoinRequestList: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyJoinRequestsMsg) Send(ctx context.Context, conn *Conn) (*PartyJoinRequestMsg, error) {
	res := new(PartyJoinRequestMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *PartyJoinRequestsMsg) Async(ctx context.Context, conn *Conn, f func(*PartyJoinRequestMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyJoinRequestMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyJoinRequest{
			PartyJoinRequest: msg,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyLeaderMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyLeader{
			PartyLeader: msg,
		},
	}
}

// PartyLeave creates a realtime message to leave a party.
func PartyLeave(partyId string) *PartyLeaveMsg {
	return &PartyLeaveMsg{
		PartyId: partyId,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyLeaveMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyLeave{
			PartyLeave: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyLeaveMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *PartyLeaveMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// PartyMatchmakerAdd creates a realtime message to begin matchmaking as a party.
func PartyMatchmakerAdd(partyId, query string, minCount, maxCount int) *PartyMatchmakerAddMsg {
	return &PartyMatchmakerAddMsg{
		PartyId:  partyId,
		Query:    query,
		MinCount: int32(minCount),
		MaxCount: int32(maxCount),
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyMatchmakerAddMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyMatchmakerAdd{
			PartyMatchmakerAdd: msg,
		},
	}
}

// WithStringProperties sets the stringProperties on the message.
func (msg *PartyMatchmakerAddMsg) WithStringProperties(stringProperties map[string]string) *PartyMatchmakerAddMsg {
	msg.StringProperties = stringProperties
	return msg
}

// WithNumericProperties sets the stringProperties on the message.
func (msg *PartyMatchmakerAddMsg) WithNumericProperties(numericProperties map[string]float64) *PartyMatchmakerAddMsg {
	msg.NumericProperties = numericProperties
	return msg
}

// WithCountMultiple sets the stringProperties on the message.
func (msg *PartyMatchmakerAddMsg) WithCountMultiple(countMultiple int) *PartyMatchmakerAddMsg {
	msg.CountMultiple = wrapperspb.Int32(int32(countMultiple))
	return msg
}

// Send sends the message to the connection.
func (msg *PartyMatchmakerAddMsg) Send(ctx context.Context, conn *Conn) (*PartyMatchmakerTicketMsg, error) {
	res := new(PartyMatchmakerTicketMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *PartyMatchmakerAddMsg) Async(ctx context.Context, conn *Conn, f func(*PartyMatchmakerTicketMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// PartyMatchmakerRemove creates a realtime message to cancel a party matchmaking process for a ticket.
func PartyMatchmakerRemove(partyId, ticket string) *PartyMatchmakerRemoveMsg {
	return &PartyMatchmakerRemoveMsg{
		PartyId: partyId,
		Ticket:  ticket,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyMatchmakerRemoveMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyMatchmakerRemove{
			PartyMatchmakerRemove: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyMatchmakerRemoveMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *PartyMatchmakerRemoveMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyMatchmakerTicketMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyMatchmakerTicket{
			PartyMatchmakerTicket: msg,
		},
	}
}

// PartyPromote creates a realtime message to promote a new party leader.
func PartyPromote(partyId string, presence *UserPresenceMsg) *PartyPromoteMsg {
	return &PartyPromoteMsg{
		PartyId:  partyId,
		Presence: presence,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyPromoteMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyPromote{
			PartyPromote: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyPromoteMsg) Send(ctx context.Context, conn *Conn) (*PartyLeaderMsg, error) {
	res := new(PartyLeaderMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *PartyPromoteMsg) Async(ctx context.Context, conn *Conn, f func(*PartyLeaderMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// PartyRemove creates a realtime message to kick a party member or decline a request to join.
func PartyRemove(partyId string, presence *UserPresenceMsg) *PartyRemoveMsg {
	return &PartyRemoveMsg{
		PartyId:  partyId,
		Presence: presence,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyRemoveMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_PartyRemove{
			PartyRemove: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PartyRemoveMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *PartyRemoveMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// Ping creates a realtime message to do a ping.
func Ping() *PingMsg {
	return &PingMsg{}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PingMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_Ping{
			Ping: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *PingMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *PingMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *RpcMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_Rpc{
			Rpc: msg,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_Status{
			Status: msg,
		},
	}
}

// StatusFollow creates a realtime message to subscribe to user status updates.
func StatusFollow(userIds ...string) *StatusFollowMsg {
	return &StatusFollowMsg{
		UserIds: userIds,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusFollowMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_StatusFollow{
			StatusFollow: msg,
		},
	}
}

// WithUsernames sets the usernames on the message.
func (msg *StatusFollowMsg) WithUsernames(usernames ...string) *StatusFollowMsg {
	msg.Usernames = usernames
	return msg
}

// Send sends the message to the connection.
func (msg *StatusFollowMsg) Send(ctx context.Context, conn *Conn) (*StatusMsg, error) {
	res := new(StatusMsg)
	if err := conn.Send(ctx, msg, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async sends the message to the connection.
func (msg *StatusFollowMsg) Async(ctx context.Context, conn *Conn, f func(*StatusMsg, error)) {
	go func() {
		if res, err := msg.Send(ctx, conn); f != nil {
			f(res, err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusPresenceEventMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_StatusPresenceEvent{
			StatusPresenceEvent: msg,
		},
	}
}

// StatusUnfollow creates a realtime message to unfollow user's status updates.
func StatusUnfollow(userIds ...string) *StatusUnfollowMsg {
	return &StatusUnfollowMsg{
		UserIds: userIds,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusUnfollowMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_StatusUnfollow{
			StatusUnfollow: msg,
		},
	}
}

// Send sends the message to the connection.
func (msg *StatusUnfollowMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *StatusUnfollowMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// StatusUpdate creates a realtime message to update the user's status.
func StatusUpdate() *StatusUpdateMsg {
	return &StatusUpdateMsg{}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusUpdateMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_StatusUpdate{
			StatusUpdate: msg,
		},
	}
}

// WithStatus sets the status on the message.
func (msg *StatusUpdateMsg) WithStatus(status string) *StatusUpdateMsg {
	msg.Status = wrapperspb.String(status)
	return msg
}

// Send sends the message to the connection.
func (msg *StatusUpdateMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *StatusUpdateMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		if err := msg.Send(ctx, conn); f != nil {
			f(err)
		}
	}()
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StreamDataMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_StreamData{
			StreamData: msg,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StreamPresenceEventMsg) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_StreamPresenceEvent{
			StreamPresenceEvent: msg,
		},
	}
}

// UserPresence creates a new realtime user presence message.
func UserPresence() *UserPresenceMsg {
	return &UserPresenceMsg{}
}

// WithUserId sets the user id on the message.
func (msg *UserPresenceMsg) WithUserId(userId string) *UserPresenceMsg {
	msg.UserId = userId
	return msg
}

// WithSessionId sets the session id on the message.
func (msg *UserPresenceMsg) WithSessionId(sessionId string) *UserPresenceMsg {
	msg.SessionId = sessionId
	return msg
}

// WithUsername sets the username on the message.
func (msg *UserPresenceMsg) WithUsername(username string) *UserPresenceMsg {
	msg.Username = username
	return msg
}

// WithPersistence sets the persistence on the message.
func (msg *UserPresenceMsg) WithPersistence(persistence bool) *UserPresenceMsg {
	msg.Persistence = persistence
	return msg
}

// WithStatus sets the status on the message.
func (msg *UserPresenceMsg) WithStatus(status string) *UserPresenceMsg {
	msg.Status = wrapperspb.String(status)
	return msg
}

// emptyMsg is an empty message.
type emptyMsg struct{}

// empty creates a new empty message.
func empty() emptyMsg {
	return emptyMsg{}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (emptyMsg) BuildEnvelope() *Envelope {
	return new(Envelope)
}
