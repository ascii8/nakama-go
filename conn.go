package nakama

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"nhooyr.io/websocket"
)

// Handler is the interface for connection handlers.
type Handler interface {
	HttpClient() *http.Client
	WebsocketURL() (string, error)
	Token(context.Context) (string, error)
	Marshal(interface{}) (io.Reader, error)
	Unmarshal(io.Reader, interface{}) error
}

// Conn is a nakama realtime websocket connection.
type Conn struct {
	url    string
	token  string
	query  url.Values
	conn   *websocket.Conn
	cancel func()
}

// NewConn creates a new nakama realtime websocket connection.
func NewConn(ctx context.Context, h Handler, opts ...ConnOption) (*Conn, error) {
	conn := &Conn{
		query: url.Values{},
	}
	for _, o := range opts {
		o(conn)
	}
	// build url
	urlstr := conn.url
	if urlstr == "" {
		var err error
		if urlstr, err = h.WebsocketURL(); err != nil {
			return nil, err
		}
	}
	// build token
	token := conn.token
	if token == "" {
		var err error
		if token, err = h.Token(ctx); err != nil {
			return nil, err
		}
	}
	// build query
	query := url.Values{}
	for k, v := range conn.query {
		query[k] = v
	}
	query.Set("token", token)
	c, _, err := websocket.Dial(ctx, urlstr+"?"+query.Encode(), &websocket.DialOptions{
		HTTPClient: h.HttpClient(),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to dial %s: %w", urlstr, err)
	}
	ctx, conn.cancel = context.WithCancel(ctx)
	go conn.run(ctx)
	return &Conn{conn: c}, nil
}

// run handles incoming and outgoing websocket messages.
func (conn *Conn) run(ctx context.Context) {
	select {
	case <-ctx.Done():
	}
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

// ConnOption is a nakama realtime websocket connection option.
type ConnOption func(*Conn)

// WithConnUrl is a nakama websocket connection option to set the dial url.
func WithConnUrl(urlstr string) ConnOption {
	return func(conn *Conn) {
		conn.url = urlstr
	}
}

// WithConnToken is a nakama websocket connection option to set the dial token.
func WithConnToken(token string) ConnOption {
	return func(conn *Conn) {
		conn.token = token
	}
}
