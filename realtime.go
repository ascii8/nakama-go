package nakama

import (
	"context"

	rtapi "github.com/heroiclabs/nakama-common/rtapi"
)

// rtapi.ChannelJoin_Type values.
const (
	// Default case. Assumed as ROOM type.
	ChannelJoinTypeUnspecified = rtapi.ChannelJoin_TYPE_UNSPECIFIED
	// A room which anyone can join to chat.
	ChannelJoinRoom = rtapi.ChannelJoin_ROOM
	// A private channel for 1-on-1 chat.
	ChannelJoinDirectMessage = rtapi.ChannelJoin_DIRECT_MESSAGE
	// A channel for group chat.
	ChannelJoinGroup = rtapi.ChannelJoin_GROUP
)

// rtapi.Error_Code values.
const (
	// An unexpected result from the server.
	ErrRuntimeException = rtapi.Error_RUNTIME_EXCEPTION
	// The server received a message which is not recognised.
	ErrUnrecognizedPlayload = rtapi.Error_UNRECOGNIZED_PAYLOAD
	// A message was expected but contains no content.
	ErrMissingPayload = rtapi.Error_MISSING_PAYLOAD
	// Fields in the message have an invalid format.
	ErrBadInput = rtapi.Error_BAD_INPUT
	// The match id was not found.
	ErrMatchNotFound = rtapi.Error_MATCH_NOT_FOUND
	// The match join was rejected.
	ErrMatchJoinRejected = rtapi.Error_MATCH_JOIN_REJECTED
	// The runtime function does not exist on the server.
	ErrRuntimeFunctionNotFound = rtapi.Error_RUNTIME_FUNCTION_NOT_FOUND
	// The runtime function executed with an error.
	ErrRuntimeFunctionException = rtapi.Error_RUNTIME_FUNCTION_EXCEPTION
)

// Message is the interface for
type Message interface {
	BuildEnv(string) *rtapi.Envelope
}

// Ping is a realtime ping message.
type PingMessage struct {
	rtapi.Ping
}

// Ping creates a new realtime ping message.
func Ping() *PingMessage {
	return &PingMessage{}
}

// BuildEnv builds an envelope message.
func (msg *PingMessage) BuildEnv(id string) *rtapi.Envelope {
	return &rtapi.Envelope{
		Cid: id,
		Message: &rtapi.Envelope_Ping{
			Ping: &msg.Ping,
		},
	}
}

// Do sends the message to the connection.
func (msg *PingMessage) Do(ctx context.Context, conn *Conn) error {
	return conn.Do(ctx, msg, new(PongMessage))
}

// Async sends the message to the connection.
func (msg *PingMessage) Async(ctx context.Context, conn *Conn, f func(err error)) {
	go func() {
		f(conn.Do(ctx, msg, new(PongMessage)))
	}()
}

// PingResponse is a realtime pong message.
type PongMessage struct {
	rtapi.Pong
}

// BuildEnv builds an envelope message.
func (msg *PongMessage) BuildEnv(id string) *rtapi.Envelope {
	return &rtapi.Envelope{
		Cid: id,
		Message: &rtapi.Envelope_Pong{
			Pong: &msg.Pong,
		},
	}
}
