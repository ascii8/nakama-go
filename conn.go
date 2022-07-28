package nakama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"google.golang.org/protobuf/proto"
	"nhooyr.io/websocket"
)

// Handler is the interface for connection handlers.
type Handler interface {
	HttpClient() *http.Client
	Token() (string, string, error)
	Query() url.Values
	Marshal(proto.Message) ([]byte, error)
	Unmarshal([]byte, bool, proto.Message) error
}

// Conn is a nakama realtime websocket connection.
type Conn struct {
	url  string
	conn *websocket.Conn
}

// NewConn creates a new nakama realtime websocket connection.
func NewConn(ctx context.Context, h Handler, opts ...ConnOption) (*Conn, error) {
	conn := new(Conn)
	for _, o := range opts {
		o(conn)
	}
	urlstr := conn.url
	if q := h.Query(); len(q) != 0 {
		urlstr += "?" + q.Encode()
	}
	c, _, err := websocket.Dial(ctx, urlstr, &websocket.DialOptions{
		HTTPClient: h.HttpClient(),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to dial %s: %w", urlstr, err)
	}
	return &Conn{conn: c}, nil
}

// Close closes the websocket connection.
func (conn *Conn) Close() error {
	return conn.conn.Close(websocket.StatusGoingAway, "going away")
}

// ConnOption is a nakama realtime websocket connection option.
type ConnOption func(*Conn)

/*
	u, err := url.Parse(cl.url)
	switch {
	case err != nil:
		return err
	}
	scheme := "ws"
	switch strings.ToLower(u.Scheme) {
	case "http":
	case "https":
		scheme = "wss"
	default:
		return fmt.Errorf("invalid scheme %q", u.Scheme)
	}
	conn.cl = cl.cl
	conn.url = scheme + "://" + u.Host + DefaultWsPath
	conn.query.Set("token", cl.token)
*/

// WithConnUrl is a nakama websocket dial option to set the dial url.
func WithConnUrl(urlstr string) ConnOption {
	return func(conn *Conn) {
		conn.url = urlstr
	}
}
