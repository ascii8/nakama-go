package nakama

import (
	"context"
	"reflect"

	nkapi "github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/rtapi"
)

func addHandler[T any](conn *Conn, handler func(context.Context, *Conn, T)) func() {
	conn.rw.Lock()
	defer conn.rw.Unlock()
	h := reflect.ValueOf(handler)
	msgType := h.Type().In(2)
	i := len(conn.notify[msgType])
	conn.notify[msgType] = append(conn.notify[msgType], h)
	return func() {
		conn.rw.Lock()
		defer conn.rw.Unlock()
		conn.notify[msgType][i] = reflect.Value{}
	}
}

func notifyHandler(conn *Conn, ctx context.Context, msg any) {
	msgType := reflect.TypeOf(msg)
	connValue := reflect.ValueOf(conn)
	msgValue := reflect.ValueOf(msg)
	ctxValue := reflect.ValueOf(ctx)
	conn.rw.RLock()
	handlers := conn.notify[msgType]
	conn.rw.RUnlock()
	for _, v := range handlers {
		if !v.IsValid() {
			continue
		}
		go v.Call([]reflect.Value{ctxValue, connValue, msgValue})
	}
}

func (conn *Conn) HandleError(handler func(context.Context, *Conn, *ErrorMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) HandleChannelMessage(handler func(context.Context, *Conn, *ChannelMessageMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) HandleChannelPresenceEvent(handler func(context.Context, *Conn, *ChannelPresenceEventMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) HandleMatchData(handler func(context.Context, *Conn, *MatchDataMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) HandleMatchPresenceEvent(handler func(context.Context, *Conn, *MatchPresenceEventMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) HandleMatchmakerMatched(handler func(context.Context, *Conn, *MatchmakerMatchedMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) HandleNotifications(handler func(context.Context, *Conn, *NotificationsMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) HandleStatusPresenceEvent(handler func(context.Context, *Conn, *StatusPresenceEventMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) HandleStreamData(handler func(context.Context, *Conn, *StreamDataMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) HandleStreamPresenceEvent(handler func(context.Context, *Conn, StreamPresenceEventMsg)) func() {
	return addHandler(conn, handler)
}

func (conn *Conn) notifyError(ctx context.Context, msg *rtapi.Error) {
	notifyHandler(conn, ctx, &ErrorMsg{*msg})
}

func (conn *Conn) notifyChannelMessage(ctx context.Context, msg *nkapi.ChannelMessage) {
	notifyHandler(conn, ctx, &ChannelMessageMsg{*msg})
}

func (conn *Conn) notifyChannelPresenceEvent(ctx context.Context, msg *rtapi.ChannelPresenceEvent) {
	notifyHandler(conn, ctx, &ChannelPresenceEventMsg{*msg})
}

func (conn *Conn) notifyMatchData(ctx context.Context, msg *rtapi.MatchData) {
	notifyHandler(conn, ctx, &MatchDataMsg{*msg})
}

func (conn *Conn) notifyMatchPresenceEvent(ctx context.Context, msg *rtapi.MatchPresenceEvent) {
	notifyHandler(conn, ctx, &MatchPresenceEventMsg{*msg})
}

func (conn *Conn) notifyMatchmakerMatched(ctx context.Context, msg *rtapi.MatchmakerMatched) {
	notifyHandler(conn, ctx, &MatchmakerMatchedMsg{*msg})
}

func (conn *Conn) notifyNotifications(ctx context.Context, msg *rtapi.Notifications) {
	notifyHandler(conn, ctx, &NotificationsMsg{*msg})
}

func (conn *Conn) notifyStatusPresenceEvent(ctx context.Context, msg *rtapi.StatusPresenceEvent) {
	notifyHandler(conn, ctx, &StatusPresenceEventMsg{*msg})
}

func (conn *Conn) notifyStreamData(ctx context.Context, msg *rtapi.StreamData) {
	notifyHandler(conn, ctx, &StreamDataMsg{*msg})
}

func (conn *Conn) notifyStreamPresenceEvent(ctx context.Context, msg *rtapi.StreamPresenceEvent) {
	notifyHandler(conn, ctx, &StreamPresenceEventMsg{*msg})
}
