package nakama

import (
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
