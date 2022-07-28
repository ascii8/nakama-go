package nakama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"nhooyr.io/websocket"
)

// DefaultWsPath is the default websocket path.
var DefaultWsPath = "/ws"

// Conn wraps a websocket connection.
type Conn struct {
	cl     *http.Client
	url    string
	query  url.Values
	header http.Header
	conn   *websocket.Conn
}

// Dial creates a new websocket connection.
func Dial(ctx context.Context, opts ...DialOption) (*Conn, error) {
	conn := new(Conn)
	for _, o := range opts {
		if err := o(conn); err != nil {
			return nil, err
		}
	}
	urlstr := conn.BuildUrl()
	c, _, err := websocket.Dial(ctx, urlstr, &websocket.DialOptions{
		HTTPClient: conn.cl,
		HTTPHeader: conn.header,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to dial %s: %w", urlstr, err)
	}
	return &Conn{conn: c}, nil
}

// BuildUrl builds the url.
func (conn *Conn) BuildUrl() string {
	urlstr := conn.url
	if len(conn.query) != 0 {
		urlstr += "?" + conn.query.Encode()
	}
	return urlstr
}

// Close closes the websocket connection.
func (conn *Conn) Close() error {
	return conn.conn.Close(websocket.StatusGoingAway, "going away")
}

// DialOption is a nakama websocket dial option.
type DialOption func(*Conn) error

// FromClient is a nakama websocket dial option to copy settings from a nakama
// client.
func FromClient(cl *Client) DialOption {
	return func(conn *Conn) error {
		u, err := url.Parse(cl.url)
		switch {
		case err != nil:
			return err
		case cl.serverKey != "":
			u.User = url.UserPassword(cl.serverKey, "")
		case cl.username != "":
			u.User = url.UserPassword(cl.username, cl.password)
		}
		switch strings.ToLower(u.Scheme) {
		case "http":
			u.Scheme = "ws"
		case "https":
			u.Scheme = "wss"
		default:
			return fmt.Errorf("invalid scheme %q", u.Scheme)
		}
		conn.url = u.String() + DefaultWsPath
		conn.cl = cl.cl
		if cl.token != "" {
			conn.header = make(http.Header)
			conn.header.Add("Authorization", "Bearer "+cl.token)
		}
		if cl.httpKey != "" {
			if conn.query == nil {
				conn.query = url.Values{}
			}
			conn.query.Set("http_key", cl.httpKey)
		}
		return nil
	}
}

// WithDialUrl is a nakama websocket dial option to set the url.
func WithDialUrl(url string) DialOption {
	return func(conn *Conn) error {
		conn.url = url
		return nil
	}
}

// WithFormat is a nakama websocket dial option to set the http_key on the dial
// url (<url>?http_key=<httpKey>).
func WithDialHttpKey(httpKey string) DialOption {
	return func(conn *Conn) error {
		if conn.query == nil {
			conn.query = url.Values{}
		}
		conn.query.Set("http_key", httpKey)
		return nil
	}
}

// WithDialToken is a nakama websocket dial option to set the token on the dial
// url (<url>?token=<token>).
func WithDialToken(token string) DialOption {
	return func(conn *Conn) error {
		if conn.query == nil {
			conn.query = url.Values{}
		}
		conn.query.Set("token", token)
		return nil
	}
}

// WithDialFormat is a nakama websocket dial option to set the format on the
// dial url (<url>?format=<format>).
func WithDialFormat(format string) DialOption {
	return func(conn *Conn) error {
		if conn.query == nil {
			conn.query = url.Values{}
		}
		conn.query.Set("format", format)
		return nil
	}
}

// WithDialCreateStatus is a nakama websocket dial option to set the status on
// the dial url (<url>?status=<true/false>).
func WithDialCreateStatus(create bool) DialOption {
	return func(conn *Conn) error {
		if conn.query == nil {
			conn.query = url.Values{}
		}
		conn.query.Set("status", strconv.FormatBool(create))
		return nil
	}
}
