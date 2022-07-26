package nakama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	nkapi "github.com/heroiclabs/nakama-common/api"
	"golang.org/x/net/publicsuffix"
)

// Client is a nakama client.
type Client struct {
	cl        *http.Client
	url       string
	token     string
	username  string
	password  string
	jar       http.CookieJar
	transport http.RoundTripper
}

// New creates a new nakama client.
func New(opts ...Option) *Client {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	cl := &Client{
		url: "http://127.0.0.1:7350",
		jar: jar,
	}
	for _, o := range opts {
		o(cl)
	}
	cl.url = strings.TrimSuffix(cl.url, "/")
	cl.cl = &http.Client{
		Transport: cl.transport,
	}
	return cl
}

// Do executes a request of method and type, encoding params to the request and
// decoding results to v.
func (cl *Client) Do(ctx context.Context, method, typ string, query url.Values, params, v interface{}) error {
	urlstr := cl.url + "/" + typ
	if len(query) != 0 {
		urlstr += "?" + query.Encode()
	}
	u, err := url.Parse(urlstr)
	switch {
	case err != nil:
		return err
	case cl.username != "":
		u.User = url.UserPassword(cl.username, cl.password)
	}
	// encode body
	var body io.Reader
	if params != nil {
		buf := new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		if err := enc.Encode(params); err != nil {
			return err
		}
		body = buf
	}
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	if cl.token != "" {
		req.Header.Add("Authorization", "Bearer "+cl.token)
	}
	res, err := cl.cl.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status %d != 200", res.StatusCode)
	}
	if v == nil {
		return nil
	}
	// decode result
	dec := json.NewDecoder(res.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// Healthcheck checks the health of the remote nakama server.
func (cl *Client) Healthcheck(ctx context.Context) error {
	return Healthcheck().Do(ctx, cl)
}

// Account retrieves the account of the current user.
func (cl *Client) Account(ctx context.Context) (*nkapi.Account, error) {
	return Account().Do(ctx, cl)
}

// Option is a nakama client option.
type Option func(*Client)

// WithURL is a nakama client option to set the url used.
func WithURL(urlstr string) Option {
	return func(cl *Client) {
		cl.url = urlstr
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

// WithJar is a nakama client option to set the jar used.
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
