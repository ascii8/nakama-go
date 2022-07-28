package nakama

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	nkapi "github.com/heroiclabs/nakama-common/api"
	"golang.org/x/net/publicsuffix"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// DefaultWsPath is the default websocket path.
var DefaultWsPath = "/ws"

// Client is a nakama client.
type Client struct {
	cl          *http.Client
	url         string
	serverKey   string
	username    string
	password    string
	session     *SessionResponse
	expiry      time.Time
	marshaler   *protojson.MarshalOptions
	unmarshaler *protojson.UnmarshalOptions
}

// New creates a new nakama client.
func New(opts ...Option) *Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	cl := &Client{
		cl: &http.Client{
			Jar: jar,
		},
		url: "http://127.0.0.1:7350",
		marshaler: &protojson.MarshalOptions{
			UseProtoNames:  true,
			UseEnumNumbers: true,
		},
		unmarshaler: &protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	for _, o := range opts {
		o(cl)
	}
	cl.url = strings.TrimSuffix(cl.url, "/")
	return cl
}

// HttpClient satisfies the handler interface.
func (cl *Client) HttpClient() *http.Client {
	return cl.cl
}

// WebsocketURL satisfies the Handler interface.
func (cl *Client) WebsocketURL() (string, error) {
	u, err := url.Parse(cl.url)
	if err != nil {
		return "", err
	}
	scheme := "ws"
	switch strings.ToLower(u.Scheme) {
	case "http":
	case "https":
		scheme = "wss"
	default:
		return "", fmt.Errorf("invalid scheme %q", u.Scheme)
	}
	return scheme + "://" + u.Host + DefaultWsPath, nil
}

// BuildRequest builds a http request.
func (cl *Client) BuildRequest(ctx context.Context, method, typ string, query url.Values, body io.Reader) (*http.Request, error) {
	// build url
	urlstr := cl.url + "/" + typ
	if len(query) != 0 {
		urlstr += "?" + query.Encode()
	}
	u, err := url.Parse(urlstr)
	switch {
	case err != nil:
		return nil, err
	case cl.serverKey != "" && (strings.Contains(typ, "authenticate") || strings.Contains(typ, "refresh")):
		u.User = url.UserPassword(cl.serverKey, "")
	case cl.username != "":
		u.User = url.UserPassword(cl.username, cl.password)
	}
	// create request
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	return req, nil
}

// Exec executes the request http request.
func (cl *Client) Exec(req *http.Request) (*http.Response, error) {
	res, err := cl.cl.Do(req)
	if err != nil {
		return nil, err
	}
	switch {
	case res.StatusCode != http.StatusOK:
		defer res.Body.Close()
		return nil, fmt.Errorf("status %d != 200", res.StatusCode)
	}
	return res, nil
}

// Do executes a http request with method, type and url query values, passing
// msg as the request body (when not nil), and decoding the response body to v
// (when not nil). Will attempt to refresh the session token if the session is
// expired and refresh is true.
//
// Uses Protobuf's google.golang.org/protobuf/encoding/protojson package to
// encode/decode msg and v when msg/v are a proto.Message. Otherwise uses Go's
// encoding/json package to encode/decode.
//
// See: Marshal and Unmarshal.
func (cl *Client) Do(ctx context.Context, method, typ string, refresh bool, query url.Values, msg, v interface{}) error {
	// marshal
	var body io.Reader
	if msg != nil {
		var err error
		if body, err = cl.Marshal(msg); err != nil {
			return err
		}
	}
	// build request
	req, err := cl.BuildRequest(ctx, method, typ, query, body)
	if err != nil {
		return err
	}
	// refresh
	if refresh {
		if err := cl.SessionRefresh(ctx); err != nil {
			return err
		}
	}
	// add auth token
	if refresh && cl.session != nil {
		req.Header.Set("Authorization", "Bearer "+cl.session.Token)
	}
	// exec
	res, err := cl.Exec(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if v == nil {
		return nil
	}
	// unmarshal
	return cl.Unmarshal(res.Body, v)
}

// Marshal marshals v. If v is a proto.Message, will use Protobuf's
// google.golang.org/protobuf/encoding/protojson package to encode the message,
// otherwise uses Go's encoding/json package.
func (cl *Client) Marshal(v interface{}) (io.Reader, error) {
	// protojson encode
	if msg, ok := v.(proto.Message); ok {
		if msg != nil {
			buf, err := cl.marshaler.Marshal(msg)
			if err != nil {
				return nil, err
			}
			return bytes.NewReader(buf), nil
		}
		return nil, nil
	}
	// json encode
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf, nil
}

// Unmarshal unmarshals r to v. If v is a proto.Message, will use Protobuf's
// google.golang.org/protobuf/encoding/protojson package to decode the message,
// otherwise uses Go's encoding/json package.
func (cl *Client) Unmarshal(r io.Reader, v interface{}) error {
	// protojson decode
	if msg, ok := v.(proto.Message); ok {
		buf, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		return cl.unmarshaler.Unmarshal(buf, msg)
	}
	// json decode
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// SessionStart starts a session.
func (cl *Client) SessionStart(session *SessionResponse) error {
	if session.Token == "" {
		return fmt.Errorf("unable to start session: empty token")
	}
	token := strings.Split(session.Token, ".")
	if len(token) != 3 {
		return fmt.Errorf("unable to start session: token is not jwt token")
	}
	// base64 decode
	buf, err := base64.RawStdEncoding.DecodeString(token[1])
	if err != nil {
		return fmt.Errorf("unable to start session: invalid encoding: %w", err)
	}
	// unmarshal
	var v struct {
		Exp int64 `json:"exp"`
	}
	switch err := json.NewDecoder(bytes.NewReader(buf)).Decode(&v); {
	case err != nil:
		return fmt.Errorf("unable to start session: cannot decode token: %w", err)
	case v.Exp == 0:
		return fmt.Errorf("unable to start session: expiry cannot be 0")
	}
	// check expiry
	expiry := time.Unix(v.Exp, 0)
	if time.Now().After(expiry) {
		return fmt.Errorf("unable to start session: %s (%d) is in the past", expiry, v.Exp)
	}
	cl.session, cl.expiry = session, expiry
	return nil
}

// SessionRefresh refreshes auth token for the session.
func (cl *Client) SessionRefresh(ctx context.Context) error {
	switch {
	case cl.session == nil:
		return fmt.Errorf("unable to refresh token: no active session")
	case !cl.SessionExpired():
		return nil
	}
	res, err := SessionRefresh(cl.session.RefreshToken).Do(ctx, cl)
	if err != nil {
		return fmt.Errorf("unable to refresh token: %w", err)
	}
	if err := cl.SessionStart(res); err != nil {
		return fmt.Errorf("unable to refresh token: %w", err)
	}
	return nil
}

// SessionToken returns the session token.
func (cl *Client) SessionToken() string {
	if cl.session != nil {
		return cl.session.Token
	}
	return ""
}

// SessionRefreshToken returns the session refresh token.
func (cl *Client) SessionRefreshToken() string {
	if cl.session != nil {
		return cl.session.RefreshToken
	}
	return ""
}

// SessionExpiry returns the session expiry time.
func (cl *Client) SessionExpiry() time.Time {
	return cl.expiry
}

// SessionExpired returns whether or not the session is expired.
func (cl *Client) SessionExpired() bool {
	return cl.session == nil || cl.expiry.IsZero() || time.Now().After(cl.expiry)
}

// SessionLogout logs out the session.
func (cl *Client) SessionLogout(ctx context.Context) error {
	if cl.session == nil {
		return nil
	}
	_ = SessionLogout().
		WithToken(cl.session.Token).
		WithRefreshToken(cl.session.RefreshToken).
		Do(ctx, cl)
	cl.session = nil
	return nil
}

// Token returns the current session token. Attempts to
func (cl *Client) Token(ctx context.Context) (string, error) {
	if err := cl.SessionRefresh(ctx); err != nil {
		return "", err
	}
	return cl.session.Token, nil
}

// Healthcheck checks the health of the remote nakama server.
func (cl *Client) Healthcheck(ctx context.Context) error {
	return Healthcheck().Do(ctx, cl)
}

// Account retrieves the account of the current user.
func (cl *Client) Account(ctx context.Context) (*nkapi.Account, error) {
	return Account().Do(ctx, cl)
}

// AuthenticateApple authenticates the apple token with the nakama server.
func (cl *Client) AuthenticateApple(ctx context.Context, create bool, token, username string) error {
	res, err := AuthenticateApple().
		WithCreate(create).
		WithToken(token).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateCustom authenticates the custom token with the nakama server.
func (cl *Client) AuthenticateCustom(ctx context.Context, create bool, id, username string) error {
	res, err := AuthenticateCustom().
		WithCreate(create).
		WithId(id).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateDevice authenticates the device id with the nakama server.
func (cl *Client) AuthenticateDevice(ctx context.Context, create bool, deviceId, username string) error {
	res, err := AuthenticateDevice().
		WithCreate(create).
		WithId(deviceId).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateEmail authenticates the email/password with the nakama server.
func (cl *Client) AuthenticateEmail(ctx context.Context, create bool, email, password, username string) error {
	res, err := AuthenticateEmail().
		WithCreate(create).
		WithEmail(email).
		WithPassword(password).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateFacebook authenticates the facebook token with the nakama server.
func (cl *Client) AuthenticateFacebook(ctx context.Context, create, sync bool, token, username string) error {
	res, err := AuthenticateFacebook().
		WithCreate(create).
		WithSync(sync).
		WithToken(token).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateFacebookInstantGame authenticates the facebookInstantGame token with the nakama server.
func (cl *Client) AuthenticateFacebookInstantGame(ctx context.Context, create bool, username string) error {
	res, err := AuthenticateFacebookInstantGame().
		WithCreate(create).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateGameCenter authenticates the gameCenter token with the nakama server.
func (cl *Client) AuthenticateGameCenter(ctx context.Context, create bool, username string) error {
	res, err := AuthenticateGameCenter().
		WithCreate(create).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateGoogle authenticates the google token with the nakama server.
func (cl *Client) AuthenticateGoogle(ctx context.Context, create bool, token, username string) error {
	res, err := AuthenticateGoogle().
		WithCreate(create).
		WithToken(token).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateSteam authenticates the steam token with the nakama server.
func (cl *Client) AuthenticateSteam(ctx context.Context, create, sync bool, token, username string) error {
	res, err := AuthenticateSteam().
		WithCreate(create).
		WithSync(sync).
		WithToken(token).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// NewConn creates a new a nakama realtime websocket connection, and runs until
// the context is closed.
func (cl *Client) NewConn(ctx context.Context, opts ...ConnOption) (*Conn, error) {
	return NewConn(ctx, cl, opts...)
}

// Option is a nakama client option.
type Option func(*Client)

// WithURL is a nakama client option to set the url used.
func WithURL(urlstr string) Option {
	return func(cl *Client) {
		cl.url = urlstr
	}
}

// WithServerKey is a nakama client option to set the server key used.
func WithServerKey(serverKey string) Option {
	return func(cl *Client) {
		cl.serverKey = serverKey
	}
}

// WithUsername is a nakama client option to set the username used.
func WithUsername(username string) Option {
	return func(cl *Client) {
		cl.username = username
	}
}

// WithPassword is a nakama client option to set the password used.
func WithPassword(password string) Option {
	return func(cl *Client) {
		cl.password = password
	}
}

// WithHttpClient is a nakama client option to set the underlying http.Client
// used for requests.
func WithHttpClient(httpClient *http.Client) Option {
	return func(cl *Client) {
		cl.cl = httpClient
	}
}

// WithJar is a nakama client option to set the cookie jar used by the underlying
// http.Client.
func WithJar(jar http.CookieJar) Option {
	return func(cl *Client) {
		cl.cl.Jar = jar
	}
}

// WithTransport is a nakama client option to set the transport used by the
// underlying http.Client.
func WithTransport(transport http.RoundTripper) Option {
	return func(cl *Client) {
		cl.cl.Transport = transport
	}
}
