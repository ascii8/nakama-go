package nakama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	nkapi "github.com/heroiclabs/nakama-common/api"
	"golang.org/x/net/publicsuffix"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// Client is a nakama client.
type Client struct {
	cl          *http.Client
	url         string
	token       string
	httpKey     string
	serverKey   string
	username    string
	password    string
	jar         http.CookieJar
	transport   http.RoundTripper
	marshaler   *protojson.MarshalOptions
	unmarshaler *protojson.UnmarshalOptions
}

// New creates a new nakama client.
func New(opts ...Option) *Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	cl := &Client{
		url: "http://127.0.0.1:7350",
		jar: jar,
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
	if cl.cl == nil {
		cl.cl = &http.Client{
			Transport: cl.transport,
			Jar:       cl.jar,
		}
	}
	return cl
}

// BuildRequest builds a request.
func (cl *Client) BuildRequest(ctx context.Context, method, typ string, query url.Values, body io.Reader) (*http.Request, error) {
	urlstr := cl.url + "/" + typ
	if query.Get("http_key") == "" && cl.httpKey != "" {
		query.Set("http_key", cl.httpKey)
	}
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
	if cl.token != "" {
		req.Header.Add("Authorization", "Bearer "+cl.token)
	}
	return req, nil
}

// Exec executes the request.
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

// Do executes a request of method and type, encoding params to the request and
// decoding results to v.
func (cl *Client) Do(ctx context.Context, method, typ string, query url.Values, msg, v proto.Message) error {
	// encode body
	var body io.Reader
	if msg != nil {
		buf, err := cl.marshaler.Marshal(msg)
		if err != nil {
			return err
		}
		body = bytes.NewReader(buf)
	}
	// create request
	req, err := cl.BuildRequest(ctx, method, typ, query, body)
	if err != nil {
		return err
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
	// read and unmarshal
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return cl.unmarshaler.Unmarshal(buf, v)
}

// DoRaw executes a raw json encode/decode.
func (cl *Client) DoRaw(ctx context.Context, method, typ string, query url.Values, params, v interface{}) error {
	var body io.Reader
	if params != nil {
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		if err := enc.Encode(params); err != nil {
			return err
		}
		body = buf
	}
	// create request
	req, err := cl.BuildRequest(ctx, method, typ, query, body)
	if err != nil {
		return err
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
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// SetToken sets the auth token.
func (cl *Client) SetToken(token string) {
	cl.token = token
}

// Healthcheck checks the health of the remote nakama server.
func (cl *Client) Healthcheck(ctx context.Context) error {
	return Healthcheck().Do(ctx, cl)
}

// Account retrieves the account of the current user.
func (cl *Client) Account(ctx context.Context) (*nkapi.Account, error) {
	return Account().Do(ctx, cl)
}

// Dial opens the websocket connection.
func (cl *Client) Dial(ctx context.Context, opts ...DialOption) (*Conn, error) {
	return Dial(ctx, append([]DialOption{FromClient(cl)}, opts...)...)
}

// Option is a nakama client option.
type Option func(*Client)

// WithURL is a nakama client option to set the url used.
func WithURL(urlstr string) Option {
	return func(cl *Client) {
		cl.url = urlstr
	}
}

// WithToken is a nakama client option to set the token used.
func WithToken(token string) Option {
	return func(cl *Client) {
		cl.token = token
	}
}

// WithHttpKey is a nakama client option to set the http key used.
func WithHttpKey(httpKey string) Option {
	return func(cl *Client) {
		cl.httpKey = httpKey
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

// WithServerKey is a nakama client option to set the server key used.
func WithServerKey(serverKey string) Option {
	return func(cl *Client) {
		cl.serverKey = serverKey
	}
}

// WithHttpClient is a nakama client option to set the underlying http.Client used for requests.
func WithHttpClient(httpClient *http.Client) Option {
	return func(cl *Client) {
		cl.cl = httpClient
	}
}

// WithJar is a nakama client option to set a cookie jar used by the underlying
// http.Client.
func WithJar(jar http.CookieJar) Option {
	return func(cl *Client) {
		cl.jar = jar
	}
}

// WithTransport is a nakama client option to set the http transport.
func WithTransport(transport http.RoundTripper) Option {
	return func(cl *Client) {
		cl.transport = transport
	}
}
