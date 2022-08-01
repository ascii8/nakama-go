package nakama

import (
	"context"

	nkapi "github.com/heroiclabs/nakama-common/api"
	rtapi "github.com/heroiclabs/nakama-common/rtapi"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// EnvelopeBuilder is the shared interface for realtime messages.
type EnvelopeBuilder interface {
	BuildEnvelope() *rtapi.Envelope
}

// ChannelJoinType is the channel join type.
type ChannelJoinType = rtapi.ChannelJoin_Type

// ChannelJoinType values.
const (
	// Default case. Assumed as ROOM type.
	ChannelJoinUnspecified ChannelJoinType = rtapi.ChannelJoin_TYPE_UNSPECIFIED
	// A room which anyone can join to chat.
	ChannelJoinRoom ChannelJoinType = rtapi.ChannelJoin_ROOM
	// A private channel for 1-on-1 chat.
	ChannelJoinDirectMessage ChannelJoinType = rtapi.ChannelJoin_DIRECT_MESSAGE
	// A channel for group chat.
	ChannelJoinGroup ChannelJoinType = rtapi.ChannelJoin_GROUP
)

// ErrorCode is the error code type.
type ErrorCode = rtapi.Error_Code

// ErrorCode values.
const (
	// An unexpected result from the server.
	ErrRuntimeException ErrorCode = rtapi.Error_RUNTIME_EXCEPTION
	// The server received a message which is not recognised.
	ErrUnrecognizedPlayload ErrorCode = rtapi.Error_UNRECOGNIZED_PAYLOAD
	// A message was expected but contains no content.
	ErrMissingPayload ErrorCode = rtapi.Error_MISSING_PAYLOAD
	// Fields in the message have an invalid format.
	ErrBadInput ErrorCode = rtapi.Error_BAD_INPUT
	// The match id was not found.
	ErrMatchNotFound ErrorCode = rtapi.Error_MATCH_NOT_FOUND
	// The match join was rejected.
	ErrMatchJoinRejected ErrorCode = rtapi.Error_MATCH_JOIN_REJECTED
	// The runtime function does not exist on the server.
	ErrRuntimeFunctionNotFound ErrorCode = rtapi.Error_RUNTIME_FUNCTION_NOT_FOUND
	// The runtime function executed with an error.
	ErrRuntimeFunctionException ErrorCode = rtapi.Error_RUNTIME_FUNCTION_EXCEPTION
)

// ChannelMsg is a realtime channel message.
type ChannelMsg struct {
	rtapi.Channel
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_Channel{
			Channel: &msg.Channel,
		},
	}
}

// ChannelJoinMsg is a realtime message to join a chat channel.
type ChannelJoinMsg struct {
	rtapi.ChannelJoin
}

// ChannelJoin creates a realtime message to join a chat channel.
func ChannelJoin(target string, typ ChannelJoinType) *ChannelJoinMsg {
	return &ChannelJoinMsg{
		ChannelJoin: rtapi.ChannelJoin{
			Target: target,
			Type:   int32(typ),
		},
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
func (msg *ChannelJoinMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelJoin{
			ChannelJoin: &msg.ChannelJoin,
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
		f(msg.Send(ctx, conn))
	}()
}

// ChannelLeaveMsg is a realtime message to leave a chat channel.
type ChannelLeaveMsg struct {
	rtapi.ChannelLeave
}

// ChannelLeave creates a realtime message to leave a chat channel.
func ChannelLeave(channelId string) *ChannelLeaveMsg {
	return &ChannelLeaveMsg{
		ChannelLeave: rtapi.ChannelLeave{
			ChannelId: channelId,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelLeaveMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelLeave{
			ChannelLeave: &msg.ChannelLeave,
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
		f(msg.Send(ctx, conn))
	}()
}

// ChannelMessageMsg is a realtime channel message message.
type ChannelMessageMsg struct {
	nkapi.ChannelMessage
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessage{
			ChannelMessage: &msg.ChannelMessage,
		},
	}
}

// ChannelMessageAckMsg is a realtime channel message ack message.
type ChannelMessageAckMsg struct {
	rtapi.ChannelMessageAck
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageAckMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessageAck{
			ChannelMessageAck: &msg.ChannelMessageAck,
		},
	}
}

// ChannelMessageRemoveMsg is a realtime message to remove a message from a channel.
type ChannelMessageRemoveMsg struct {
	rtapi.ChannelMessageRemove
}

// ChannelMessageRemove creates a realtime message to remove a message from a channel.
func ChannelMessageRemove(channelId, messageId string) *ChannelMessageRemoveMsg {
	return &ChannelMessageRemoveMsg{
		ChannelMessageRemove: rtapi.ChannelMessageRemove{
			ChannelId: channelId,
			MessageId: messageId,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageRemoveMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessageRemove{
			ChannelMessageRemove: &msg.ChannelMessageRemove,
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
		f(msg.Send(ctx, conn))
	}()
}

// ChannelMessageSendMsg is a realtime message to send a message on a channel.
type ChannelMessageSendMsg struct {
	rtapi.ChannelMessageSend
}

// ChannelMessageSend creates a realtime message to send a message on a channel.
func ChannelMessageSend(channelId, content string) *ChannelMessageSendMsg {
	return &ChannelMessageSendMsg{
		ChannelMessageSend: rtapi.ChannelMessageSend{
			ChannelId: channelId,
			Content:   content,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageSendMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessageSend{
			ChannelMessageSend: &msg.ChannelMessageSend,
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
		f(msg.Send(ctx, conn))
	}()
}

// ChannelMessageUpdateMsg is a realtime message to update a message on a channel.
type ChannelMessageUpdateMsg struct {
	rtapi.ChannelMessageUpdate
}

// ChannelMessageUpdate creates a realtime message to update a message on a channel.
func ChannelMessageUpdate(channelId, messageId, content string) *ChannelMessageUpdateMsg {
	return &ChannelMessageUpdateMsg{
		ChannelMessageUpdate: rtapi.ChannelMessageUpdate{
			ChannelId: channelId,
			MessageId: messageId,
			Content:   content,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelMessageUpdateMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessageUpdate{
			ChannelMessageUpdate: &msg.ChannelMessageUpdate,
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
		f(msg.Send(ctx, conn))
	}()
}

// ChannelPresenceEventMsg is a realtime channel presence event message.
type ChannelPresenceEventMsg struct {
	rtapi.ChannelPresenceEvent
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ChannelPresenceEventMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelPresenceEvent{
			ChannelPresenceEvent: &msg.ChannelPresenceEvent,
		},
	}
}

// ErrorMsg is a realtime error message.
type ErrorMsg struct {
	rtapi.Error
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *ErrorMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_Error{
			Error: &msg.Error,
		},
	}
}

// MatchMsg is a realtime match message.
type MatchMsg struct {
	rtapi.Match
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_Match{
			Match: &msg.Match,
		},
	}
}

// MatchCreateMsg is a realtime message to create a multiplayer match.
type MatchCreateMsg struct {
	rtapi.MatchCreate
}

// MatchCreate creates a realtime message to create a multiplayer match.
func MatchCreate(name string) *MatchCreateMsg {
	return &MatchCreateMsg{
		MatchCreate: rtapi.MatchCreate{
			Name: name,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchCreateMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchCreate{
			MatchCreate: &msg.MatchCreate,
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
		f(msg.Send(ctx, conn))
	}()
}

// MatchDataMsg is a realtime match data message.
type MatchDataMsg struct {
	rtapi.MatchData
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchDataMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchData{
			MatchData: &msg.MatchData,
		},
	}
}

// MatchDataSendMsg is a realtime message to send input to a multiplayer match.
type MatchDataSendMsg struct {
	rtapi.MatchDataSend
}

// MatchDataSend creates a realtime message to send input to a multiplayer match.
func MatchDataSend(matchId string, opCode OpType, data []byte) *MatchDataSendMsg {
	return &MatchDataSendMsg{
		MatchDataSend: rtapi.MatchDataSend{
			MatchId: matchId,
			OpCode:  int64(opCode),
			Data:    data,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchDataSendMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchDataSend{
			MatchDataSend: &msg.MatchDataSend,
		},
	}
}

// WithPresences sets the presences on the message.
func (msg *MatchDataSendMsg) WithPresences(presences ...*UserPresenceMsg) *MatchDataSendMsg {
	p := make([]*rtapi.UserPresence, len(presences))
	for i, presence := range presences {
		p[i] = &presence.UserPresence
	}
	msg.Presences = p
	return msg
}

// WithReliable sets the reliable on the message.
func (msg *MatchDataSendMsg) WithReliable(reliable bool) *MatchDataSendMsg {
	msg.Reliable = reliable
	return msg
}

// Send sends the message to the connection.
func (msg *MatchDataSendMsg) Send(ctx context.Context, conn *Conn) error {
	return conn.Send(ctx, msg, empty())
}

// Async sends the message to the connection.
func (msg *MatchDataSendMsg) Async(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		f(msg.Send(ctx, conn))
	}()
}

// MatchJoinMsg is a realtime message to join a match.
type MatchJoinMsg struct {
	rtapi.MatchJoin
}

// MatchJoin creates a realtime message to join a match.
func MatchJoin(token string) *MatchJoinMsg {
	return &MatchJoinMsg{
		MatchJoin: rtapi.MatchJoin{
			Id: &rtapi.MatchJoin_Token{
				Token: token,
			},
		},
	}
}

// MatchJoinToken creates a new realtime to join a match with a token.
func MatchJoinToken(token string) *MatchJoinMsg {
	return &MatchJoinMsg{
		MatchJoin: rtapi.MatchJoin{
			Id: &rtapi.MatchJoin_Token{
				Token: token,
			},
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchJoinMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchJoin{
			MatchJoin: &msg.MatchJoin,
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
		f(msg.Send(ctx, conn))
	}()
}

// MatchLeaveMsg is a realtime message to leave a multiplayer match.
type MatchLeaveMsg struct {
	rtapi.MatchLeave
}

// MatchLeave creates a realtime message to leave a multiplayer match.
func MatchLeave(matchId string) *MatchLeaveMsg {
	return &MatchLeaveMsg{
		MatchLeave: rtapi.MatchLeave{
			MatchId: matchId,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchLeaveMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchLeave{
			MatchLeave: &msg.MatchLeave,
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
		f(msg.Send(ctx, conn))
	}()
}

// MatchPresenceEventMsg is a realtime match presence event message.
type MatchPresenceEventMsg struct {
	rtapi.MatchPresenceEvent
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchPresenceEventMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchPresenceEvent{
			MatchPresenceEvent: &msg.MatchPresenceEvent,
		},
	}
}

// MatchmakerAddMsg is a realtime message to join the matchmaker pool and search for opponents on the server.
type MatchmakerAddMsg struct {
	rtapi.MatchmakerAdd
}

// MatchmakerAdd creates a realtime message to join the matchmaker pool and search for opponents on the server.
func MatchmakerAdd(query string, minCount, maxCount int) *MatchmakerAddMsg {
	return &MatchmakerAddMsg{
		MatchmakerAdd: rtapi.MatchmakerAdd{
			Query:    query,
			MinCount: int32(minCount),
			MaxCount: int32(maxCount),
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchmakerAddMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchmakerAdd{
			MatchmakerAdd: &msg.MatchmakerAdd,
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
		f(msg.Send(ctx, conn))
	}()
}

// MatchmakerMatchedMsg is a realtime matchmaker matched message.
type MatchmakerMatchedMsg struct {
	rtapi.MatchmakerMatched
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchmakerMatchedMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchmakerMatched{
			MatchmakerMatched: &msg.MatchmakerMatched,
		},
	}
}

// MatchmakerRemoveMsg is a realtime message to leave the matchmaker pool for a ticket.
type MatchmakerRemoveMsg struct {
	rtapi.MatchmakerRemove
}

// MatchmakerRemove creates a realtime message to leave the matchmaker pool for a ticket.
func MatchmakerRemove(ticket string) *MatchmakerRemoveMsg {
	return &MatchmakerRemoveMsg{
		MatchmakerRemove: rtapi.MatchmakerRemove{
			Ticket: ticket,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchmakerRemoveMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchmakerRemove{
			MatchmakerRemove: &msg.MatchmakerRemove,
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
		f(msg.Send(ctx, conn))
	}()
}

// MatchmakerTicketMsg is a realtime matchmaker ticket message.
type MatchmakerTicketMsg struct {
	rtapi.MatchmakerTicket
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *MatchmakerTicketMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchmakerTicket{
			MatchmakerTicket: &msg.MatchmakerTicket,
		},
	}
}

// NotificationsMsg is a realtime notifications message.
type NotificationsMsg struct {
	rtapi.Notifications
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *NotificationsMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_Notifications{
			Notifications: &msg.Notifications,
		},
	}
}

// PartyMsg is a realtime party message.
type PartyMsg struct {
	rtapi.Party
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_Party{
			Party: &msg.Party,
		},
	}
}

// PartyAcceptMsg is a realtime message to accept a party member.
type PartyAcceptMsg struct {
	rtapi.PartyAccept
}

// PartyAccept creates a realtime message to accept a party member.
func PartyAccept(partyId string, presence *UserPresenceMsg) *PartyAcceptMsg {
	return &PartyAcceptMsg{
		PartyAccept: rtapi.PartyAccept{
			PartyId:  partyId,
			Presence: &presence.UserPresence,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyAcceptMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyAccept{
			PartyAccept: &msg.PartyAccept,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyCloseMsg is a realtime message to close a party, kicking all party members.
type PartyCloseMsg struct {
	rtapi.PartyClose
}

// PartyClose creates a realtime message to close a party, kicking all party members.
func PartyClose(partyId string) *PartyCloseMsg {
	return &PartyCloseMsg{
		PartyClose: rtapi.PartyClose{
			PartyId: partyId,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyCloseMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyClose{
			PartyClose: &msg.PartyClose,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyCreateMsg is a realtime message to create a party.
type PartyCreateMsg struct {
	rtapi.PartyCreate
}

// PartyCreate creates a realtime message to create a party.
func PartyCreate(open bool, maxSize int) *PartyCreateMsg {
	return &PartyCreateMsg{
		PartyCreate: rtapi.PartyCreate{
			Open:    open,
			MaxSize: int32(maxSize),
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyCreateMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyCreate{
			PartyCreate: &msg.PartyCreate,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyDataSendMsg is a realtime message to send data to a party.
type PartyDataSendMsg struct {
	rtapi.PartyDataSend
}

// PartyDataSend creates a realtime message to send data to a party.
func PartyDataSend(partyId string, opCode OpType, data []byte) *PartyDataSendMsg {
	return &PartyDataSendMsg{
		PartyDataSend: rtapi.PartyDataSend{
			PartyId: partyId,
			OpCode:  int64(opCode),
			Data:    data,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyDataSendMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyDataSend{
			PartyDataSend: &msg.PartyDataSend,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyJoinMsg is a realtime message to join a party.
type PartyJoinMsg struct {
	rtapi.PartyJoin
}

// PartyJoin creates a realtime message to join a party.
func PartyJoin(partyId string) *PartyJoinMsg {
	return &PartyJoinMsg{
		PartyJoin: rtapi.PartyJoin{
			PartyId: partyId,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyJoinMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyJoin{
			PartyJoin: &msg.PartyJoin,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyJoinRequestsMsg is a realtime message to request the list of pending join requests for a party.
type PartyJoinRequestsMsg struct {
	rtapi.PartyJoinRequestList
}

// PartyJoinRequests creates a realtime message to request the list of pending join requests for a party.
func PartyJoinRequests(partyId string) *PartyJoinRequestsMsg {
	return &PartyJoinRequestsMsg{
		PartyJoinRequestList: rtapi.PartyJoinRequestList{
			PartyId: partyId,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyJoinRequestsMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyJoinRequestList{
			PartyJoinRequestList: &msg.PartyJoinRequestList,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyJoinRequestMsg is a realtime party join request message.
type PartyJoinRequestMsg struct {
	rtapi.PartyJoinRequest
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyJoinRequestMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyJoinRequest{
			PartyJoinRequest: &msg.PartyJoinRequest,
		},
	}
}

// PartyLeaderMsg is a realtime party leader message.
type PartyLeaderMsg struct {
	rtapi.PartyLeader
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyLeaderMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyLeader{
			PartyLeader: &msg.PartyLeader,
		},
	}
}

// PartyLeaveMsg is a realtime message to leave a party.
type PartyLeaveMsg struct {
	rtapi.PartyLeave
}

// PartyLeave creates a realtime message to leave a party.
func PartyLeave(partyId string) *PartyLeaveMsg {
	return &PartyLeaveMsg{
		PartyLeave: rtapi.PartyLeave{
			PartyId: partyId,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyLeaveMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyLeave{
			PartyLeave: &msg.PartyLeave,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyMatchmakerAddMsg is a realtime message to begin matchmaking as a party.
type PartyMatchmakerAddMsg struct {
	rtapi.PartyMatchmakerAdd
}

// PartyMatchmakerAdd creates a realtime message to begin matchmaking as a party.
func PartyMatchmakerAdd(partyId, query string, minCount, maxCount int) *PartyMatchmakerAddMsg {
	return &PartyMatchmakerAddMsg{
		PartyMatchmakerAdd: rtapi.PartyMatchmakerAdd{
			PartyId:  partyId,
			Query:    query,
			MinCount: int32(minCount),
			MaxCount: int32(maxCount),
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyMatchmakerAddMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyMatchmakerAdd{
			PartyMatchmakerAdd: &msg.PartyMatchmakerAdd,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyMatchmakerRemoveMsg is a realtime message to cancel a party matchmaking process for a ticket.
type PartyMatchmakerRemoveMsg struct {
	rtapi.PartyMatchmakerRemove
}

// PartyMatchmakerRemove creates a realtime message to cancel a party matchmaking process for a ticket.
func PartyMatchmakerRemove(partyId, ticket string) *PartyMatchmakerRemoveMsg {
	return &PartyMatchmakerRemoveMsg{
		PartyMatchmakerRemove: rtapi.PartyMatchmakerRemove{
			PartyId: partyId,
			Ticket:  ticket,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyMatchmakerRemoveMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyMatchmakerRemove{
			PartyMatchmakerRemove: &msg.PartyMatchmakerRemove,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyMatchmakerTicketMsg is a realtime party matchmaker ticket message.
type PartyMatchmakerTicketMsg struct {
	rtapi.PartyMatchmakerTicket
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyMatchmakerTicketMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyMatchmakerTicket{
			PartyMatchmakerTicket: &msg.PartyMatchmakerTicket,
		},
	}
}

// PartyPromoteMsg is a realtime message to promote a new party leader.
type PartyPromoteMsg struct {
	rtapi.PartyPromote
}

// PartyPromote creates a realtime message to promote a new party leader.
func PartyPromote(partyId string, presence *UserPresenceMsg) *PartyPromoteMsg {
	return &PartyPromoteMsg{
		PartyPromote: rtapi.PartyPromote{
			PartyId:  partyId,
			Presence: &presence.UserPresence,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyPromoteMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyPromote{
			PartyPromote: &msg.PartyPromote,
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
		f(msg.Send(ctx, conn))
	}()
}

// PartyRemoveMsg is a realtime message to kick a party member or decline a request to join.
type PartyRemoveMsg struct {
	rtapi.PartyRemove
}

// PartyRemove creates a realtime message to kick a party member or decline a request to join.
func PartyRemove(partyId string, presence *UserPresenceMsg) *PartyRemoveMsg {
	return &PartyRemoveMsg{
		PartyRemove: rtapi.PartyRemove{
			PartyId:  partyId,
			Presence: &presence.UserPresence,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PartyRemoveMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyRemove{
			PartyRemove: &msg.PartyRemove,
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
		f(msg.Send(ctx, conn))
	}()
}

// PingMsg is a realtime message to do a ping.
type PingMsg struct {
	rtapi.Ping
}

// Ping creates a realtime message to do a ping.
func Ping() *PingMsg {
	return &PingMsg{}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *PingMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_Ping{
			Ping: &msg.Ping,
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
		f(msg.Send(ctx, conn))
	}()
}

// rpcMsg is a realtime rpc message.
type rpcMsg struct {
	nkapi.Rpc
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *rpcMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_Rpc{
			Rpc: &msg.Rpc,
		},
	}
}

// StatusMsg is a realtime status message.
type StatusMsg struct {
	rtapi.Status
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_Status{
			Status: &msg.Status,
		},
	}
}

// StatusFollowMsg is a realtime message to subscribe to user status updates.
type StatusFollowMsg struct {
	rtapi.StatusFollow
}

// StatusFollow creates a realtime message to subscribe to user status updates.
func StatusFollow(userIds ...string) *StatusFollowMsg {
	return &StatusFollowMsg{
		StatusFollow: rtapi.StatusFollow{
			UserIds: userIds,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusFollowMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_StatusFollow{
			StatusFollow: &msg.StatusFollow,
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
		f(msg.Send(ctx, conn))
	}()
}

// StatusPresenceEventMsg is a realtime statusPresenceEvent message.
type StatusPresenceEventMsg struct {
	rtapi.StatusPresenceEvent
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusPresenceEventMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_StatusPresenceEvent{
			StatusPresenceEvent: &msg.StatusPresenceEvent,
		},
	}
}

// StatusUnfollowMsg is a realtime message to unfollow user's status updates.
type StatusUnfollowMsg struct {
	rtapi.StatusUnfollow
}

// StatusUnfollow creates a realtime message to unfollow user's status updates.
func StatusUnfollow(userIds ...string) *StatusUnfollowMsg {
	return &StatusUnfollowMsg{
		StatusUnfollow: rtapi.StatusUnfollow{
			UserIds: userIds,
		},
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusUnfollowMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_StatusUnfollow{
			StatusUnfollow: &msg.StatusUnfollow,
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
		f(msg.Send(ctx, conn))
	}()
}

// StatusUpdateMsg is a realtime message to update the user's status.
type StatusUpdateMsg struct {
	rtapi.StatusUpdate
}

// StatusUpdate creates a realtime message to update the user's status.
func StatusUpdate() *StatusUpdateMsg {
	return &StatusUpdateMsg{}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StatusUpdateMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_StatusUpdate{
			StatusUpdate: &msg.StatusUpdate,
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
		f(msg.Send(ctx, conn))
	}()
}

// StreamDataMsg is a realtime streamData message.
type StreamDataMsg struct {
	rtapi.StreamData
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StreamDataMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_StreamData{
			StreamData: &msg.StreamData,
		},
	}
}

// StreamPresenceEventMsg is a realtime streamPresenceEvent message.
type StreamPresenceEventMsg struct {
	rtapi.StreamPresenceEvent
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (msg *StreamPresenceEventMsg) BuildEnvelope() *rtapi.Envelope {
	return &rtapi.Envelope{
		Message: &rtapi.Envelope_StreamPresenceEvent{
			StreamPresenceEvent: &msg.StreamPresenceEvent,
		},
	}
}

// UserPresenceMsg is a realtime user presence message.
type UserPresenceMsg struct {
	rtapi.UserPresence
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
func (emptyMsg) BuildEnvelope() *rtapi.Envelope {
	return new(rtapi.Envelope)
}
