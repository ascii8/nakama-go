package nakama

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
	"google.golang.org/grpc/codes"
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
	refreshAuto bool
	expiryGrace time.Duration

	session             *SessionResponse
	expiry              time.Time
	expiryGraced        time.Time
	expiryRefresh       time.Time
	expiryRefreshGraced time.Time

	marshaler   *protojson.MarshalOptions
	unmarshaler *protojson.UnmarshalOptions

	logf func(string, ...interface{})

	rw sync.RWMutex
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
		url:         "http://127.0.0.1:7350",
		refreshAuto: true,
		expiryGrace: 5 * time.Second,
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

// Logf satisfies the Handler interface.
func (cl *Client) Logf(s string, v ...interface{}) {
	if cl.logf != nil {
		cl.logf(s, v...)
	}
}

// Errf satisfies the handler interface.
func (cl *Client) Errf(s string, v ...interface{}) {
	cl.Logf("ERROR: "+s, v...)
}

// HttpClient satisfies the handler interface.
func (cl *Client) HttpClient() *http.Client {
	return cl.cl
}

// SocketURL satisfies the Handler interface.
func (cl *Client) SocketURL() (string, error) {
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

// Token returns the current session token. Satisfies the Handler interface.
func (cl *Client) Token(ctx context.Context) (string, error) {
	if err := cl.SessionRefresh(ctx); err != nil {
		return "", err
	}
	return cl.session.Token, nil
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
		return nil, NewClientErrorFromReader(res.StatusCode, res.Body)
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
func (cl *Client) Do(ctx context.Context, method, typ string, session bool, query url.Values, msg, v interface{}) error {
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
	if session && cl.refreshAuto {
		if err := cl.SessionRefresh(ctx); err != nil {
			return err
		}
	}
	// check active session
	switch {
	case session && cl.session == nil:
		// error here ?
	case session && cl.session != nil:
		// add auth token
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
	msg, ok := v.(proto.Message)
	if ok {
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
		buf, err := io.ReadAll(r)
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

/*
// MarshalBytes marshals v. If v is a proto.Message, will use Protobuf's
// google.golang.org/protobuf/encoding/protojson package to encode the message,
// otherwise uses Go's encoding/json package.
func (cl *Client) MarshalBytes(v interface{}) ([]byte, error) {
	// protojson encode
	if msg, ok := v.(proto.Message); ok {
		if msg != nil {
			return cl.marshaler.Marshal(msg)
		}
		return nil, nil
	}
	// json encode
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBytes unmarshals r to v. If v is a proto.Message, will use
// Protobuf's google.golang.org/protobuf/encoding/protojson package to decode
// the message, otherwise uses Go's encoding/json package.
func (cl *Client) UnmarshalBytes(buf []byte, v interface{}) error {
	// protojson decode
	if msg, ok := v.(proto.Message); ok {
		return cl.unmarshaler.Unmarshal(buf, msg)
	}
	// json decode
	dec := json.NewDecoder(bytes.NewReader(buf))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
*/

// SessionStart starts a session.
func (cl *Client) SessionStart(session *SessionResponse) error {
	expiry, expiryGraced, err := ParseTokenExpiry(session.Token, "session", cl.expiryGrace)
	if err != nil {
		return fmt.Errorf("unable to start session: %w", err)
	}
	expiryRefresh, expiryRefreshGraced, err := ParseTokenExpiry(session.RefreshToken, "refresh", cl.expiryGrace)
	if err != nil {
		return fmt.Errorf("unable to start session: %w", err)
	}
	cl.rw.Lock()
	defer cl.rw.Unlock()
	cl.session, cl.expiry, cl.expiryGraced, cl.expiryRefresh, cl.expiryRefreshGraced = session, expiry, expiryGraced, expiryRefresh, expiryRefreshGraced
	return nil
}

// SessionRefresh refreshes auth token for the session.
func (cl *Client) SessionRefresh(ctx context.Context) error {
	switch {
	case cl.session == nil:
		return fmt.Errorf("unable to refresh session: no active session")
	case !cl.SessionExpired():
		return nil
	case cl.SessionRefreshExpired():
		return fmt.Errorf("unable to refresh session: refresh token expired")
	}
	res, err := SessionRefresh(cl.session.RefreshToken).Do(ctx, cl)
	if err != nil {
		return fmt.Errorf("unable to refresh session: %w", err)
	}
	if err := cl.SessionStart(res); err != nil {
		return fmt.Errorf("unable to refresh session: %w", err)
	}
	return nil
}

// SessionLogout logs out the session.
func (cl *Client) SessionLogout(ctx context.Context) error {
	cl.rw.Lock()
	defer cl.rw.Unlock()
	if cl.session == nil {
		return nil
	}
	_ = SessionLogout(cl.session.Token, cl.session.RefreshToken).Do(ctx, cl)
	cl.session, cl.expiry, cl.expiryGraced, cl.expiryRefresh, cl.expiryRefreshGraced = nil, time.Time{}, time.Time{}, time.Time{}, time.Time{}
	return nil
}

// SessionToken returns the session token.
func (cl *Client) SessionToken() string {
	cl.rw.RLock()
	defer cl.rw.RUnlock()
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

// SessionRefreshExpiry returns the session refresh expiry time.
func (cl *Client) SessionRefreshExpiry() time.Time {
	return cl.expiryRefresh
}

// SessionExpired returns whether or not the session is expired.
func (cl *Client) SessionExpired() bool {
	return cl.session == nil || cl.expiry.IsZero() || time.Now().After(cl.expiryGraced)
}

// SessionRefreshExpired returns whether or not the session refresh token is expired.
func (cl *Client) SessionRefreshExpired() bool {
	return cl.session == nil || cl.expiryRefresh.IsZero() || time.Now().After(cl.expiryRefreshGraced)
}

// NewConn creates a new a nakama realtime websocket connection, and runs until
// the context is closed.
func (cl *Client) NewConn(ctx context.Context, opts ...ConnOption) (*Conn, error) {
	return NewConn(ctx, append([]ConnOption{WithConnClientHandler(cl)}, opts...)...)
}

// Account retrieves the user's account.
func (cl *Client) Account(ctx context.Context) (*AccountResponse, error) {
	return Account().Do(ctx, cl)
}

// AccountAsync retrieves the user's account.
func (cl *Client) AccountAsync(ctx context.Context, f func(*AccountResponse, error)) {
	Account().Async(ctx, cl, f)
}

// Healthcheck checks the health of the server.
func (cl *Client) Healthcheck(ctx context.Context) error {
	return Healthcheck().Do(ctx, cl)
}

// Healthcheck checks the health of the server.
func (cl *Client) HealthcheckAsync(ctx context.Context, f func(error)) {
	Healthcheck().Async(ctx, cl, f)
}

// AddGroupUsers adds users to a group, or accepts their join requests.
func (cl *Client) AddGroupUsers(ctx context.Context, groupId string, userIds ...string) error {
	return AddGroupUsers(groupId, userIds...).Do(ctx, cl)
}

// AddGroupUsersAsync adds users to a group or accepts their join requests.
func (cl *Client) AddGroupUsersAsync(ctx context.Context, groupId string, userIds []string, f func(error)) {
	AddGroupUsers(groupId, userIds...).Async(ctx, cl, f)
}

// AddFriends adds friends by id.
func (cl *Client) AddFriends(ctx context.Context, ids ...string) error {
	return AddFriends(ids...).Do(ctx, cl)
}

// AddFriendsAsync adds friends by id.
func (cl *Client) AddFriendsAsync(ctx context.Context, ids []string, f func(error)) {
	AddFriends(ids...).Async(ctx, cl, f)
}

// AddFriendsUsernames adds friends by username.
func (cl *Client) AddFriendsUsernames(ctx context.Context, usernames ...string) error {
	return AddFriends().WithUsernames(usernames...).Do(ctx, cl)
}

// AddFriendsUsernamesAsync adds friends by username.
func (cl *Client) AddFriendsUsernamesAsync(ctx context.Context, usernames []string, f func(error)) {
	AddFriends().WithUsernames(usernames...).Async(ctx, cl, f)
}

// AuthenticateApple authenticates a user with a Apple token.
func (cl *Client) AuthenticateApple(ctx context.Context, token string, create bool, username string) error {
	res, err := AuthenticateApple(token).
		WithCreate(create).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateAppleAsync authenticates a user with a Apple token.
func (cl *Client) AuthenticateAppleAsync(ctx context.Context, token string, create bool, username string, f func(err error)) {
	AuthenticateApple(token).
		WithCreate(create).
		WithUsername(username).
		Async(ctx, cl, func(res *SessionResponse, err error) {
			if err == nil {
				err = cl.SessionStart(res)
			}
			f(err)
		})
}

// AuthenticateCustom authenticates a user with a id.
func (cl *Client) AuthenticateCustom(ctx context.Context, id string, create bool, username string) error {
	res, err := AuthenticateCustom(id).
		WithCreate(create).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateCustomAsync authenticates a user with a id.
func (cl *Client) AuthenticateCustomAsync(ctx context.Context, id string, create bool, username string, f func(err error)) {
	AuthenticateCustom(id).
		WithCreate(create).
		WithUsername(username).
		Async(ctx, cl, func(res *SessionResponse, err error) {
			if err == nil {
				err = cl.SessionStart(res)
			}
			f(err)
		})
}

// AuthenticateDevice authenticates a user with a device id.
func (cl *Client) AuthenticateDevice(ctx context.Context, id string, create bool, username string) error {
	res, err := AuthenticateDevice(id).
		WithCreate(create).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateDeviceAsync authenticates a user with a device id.
func (cl *Client) AuthenticateDeviceAsync(ctx context.Context, id string, create bool, username string, f func(err error)) {
	AuthenticateDevice(id).
		WithCreate(create).
		WithUsername(username).
		Async(ctx, cl, func(res *SessionResponse, err error) {
			if err == nil {
				err = cl.SessionStart(res)
			}
			f(err)
		})
}

// AuthenticateEmail authenticates a user with a email/password.
func (cl *Client) AuthenticateEmail(ctx context.Context, email, password string, create bool, username string) error {
	res, err := AuthenticateEmail(email, password).
		WithCreate(create).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateEmailAsync authenticates a user with a email/password.
func (cl *Client) AuthenticateEmailAsync(ctx context.Context, email, password string, create bool, username string, f func(err error)) {
	AuthenticateEmail(email, password).
		WithCreate(create).
		WithUsername(username).
		Async(ctx, cl, func(res *SessionResponse, err error) {
			if err == nil {
				err = cl.SessionStart(res)
			}
			f(err)
		})
}

// AuthenticateFacebook authenticates a user with a Facebook token.
func (cl *Client) AuthenticateFacebook(ctx context.Context, token string, create bool, username string, sync bool) error {
	res, err := AuthenticateFacebook(token).
		WithCreate(create).
		WithUsername(username).
		WithSync(sync).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateFacebookAsync authenticates a user with a Facebook token.
func (cl *Client) AuthenticateFacebookAsync(ctx context.Context, token string, create bool, username string, sync bool, f func(err error)) {
	AuthenticateFacebook(token).
		WithCreate(create).
		WithUsername(username).
		WithSync(sync).
		Async(ctx, cl, func(res *SessionResponse, err error) {
			if err == nil {
				err = cl.SessionStart(res)
			}
			f(err)
		})
}

// AuthenticateFacebookInstantGame authenticates a user with a Facebook Instant Game token.
func (cl *Client) AuthenticateFacebookInstantGame(ctx context.Context, token string, create bool, username string) error {
	res, err := AuthenticateFacebookInstantGame(token).
		WithCreate(create).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateFacebookInstantGameAsync authenticates a user with a Facebook Instant Game token.
func (cl *Client) AuthenticateFacebookInstantGameAsync(ctx context.Context, signedPlayerInfo string, create bool, username string, f func(err error)) {
	AuthenticateFacebookInstantGame(signedPlayerInfo).
		WithCreate(create).
		WithUsername(username).
		Async(ctx, cl, func(res *SessionResponse, err error) {
			if err == nil {
				err = cl.SessionStart(res)
			}
			f(err)
		})
}

// AuthenticateGoogle authenticates a user with a Google token.
func (cl *Client) AuthenticateGoogle(ctx context.Context, token string, create bool, username string) error {
	res, err := AuthenticateGoogle(token).
		WithCreate(create).
		WithUsername(username).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateGoogleAsync authenticates a user with a Google token.
func (cl *Client) AuthenticateGoogleAsync(ctx context.Context, token string, create bool, username string, f func(err error)) {
	AuthenticateGoogle(token).
		WithCreate(create).
		WithUsername(username).
		Async(ctx, cl, func(res *SessionResponse, err error) {
			if err == nil {
				err = cl.SessionStart(res)
			}
			f(err)
		})
}

// AuthenticateGameCenter authenticates a user with a Apple GameCenter token.
func (cl *Client) AuthenticateGameCenter(ctx context.Context, req *AuthenticateGameCenterRequest) error {
	res, err := req.Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateGameCenterAsync authenticates a user with a Apple GameCenter token.
func (cl *Client) AuthenticateGameCenterAsync(ctx context.Context, req *AuthenticateGameCenterRequest, f func(err error)) {
	req.Async(ctx, cl, func(res *SessionResponse, err error) {
		if err == nil {
			err = cl.SessionStart(res)
		}
		f(err)
	})
}

// AuthenticateSteam authenticates a user with a Steam token.
func (cl *Client) AuthenticateSteam(ctx context.Context, token string, create bool, username string, sync bool) error {
	res, err := AuthenticateSteam(token).
		WithCreate(create).
		WithUsername(username).
		WithSync(sync).
		Do(ctx, cl)
	if err != nil {
		return err
	}
	return cl.SessionStart(res)
}

// AuthenticateSteamAsync authenticates a user with a Steam token.
func (cl *Client) AuthenticateSteamAsync(ctx context.Context, token string, create bool, username string, sync bool, f func(err error)) {
	AuthenticateSteam(token).
		WithCreate(create).
		WithUsername(username).
		WithSync(sync).
		Async(ctx, cl, func(res *SessionResponse, err error) {
			if err == nil {
				err = cl.SessionStart(res)
			}
			f(err)
		})
}

// BanGroupUsers bans users from a group.
func (cl *Client) BanGroupUsers(ctx context.Context, groupId string, userIds ...string) error {
	return BanGroupUsers(groupId, userIds...).Do(ctx, cl)
}

// BanGroupUsersAsync bans users from a group.
func (cl *Client) BanGroupUsersAsync(ctx context.Context, groupId string, userIds []string, f func(error)) {
	BanGroupUsers(groupId, userIds...).Async(ctx, cl, f)
}

// BlockFriends blocks friends by id.
func (cl *Client) BlockFriends(ctx context.Context, ids ...string) error {
	return BlockFriends(ids...).Do(ctx, cl)
}

// BlockFriendsAsync blocks friends by id.
func (cl *Client) BlockFriendsAsync(ctx context.Context, ids []string, f func(error)) {
	BlockFriends(ids...).Async(ctx, cl, f)
}

// BlockFriendsUsernames blocks friends by username.
func (cl *Client) BlockFriendsUsernames(ctx context.Context, usernames ...string) error {
	return BlockFriends().WithUsernames(usernames...).Do(ctx, cl)
}

// BlockFriendsUsernamesAsync blocks friends by username.
func (cl *Client) BlockFriendsUsernamesAsync(ctx context.Context, usernames []string, f func(error)) {
	BlockFriends().WithUsernames(usernames...).Async(ctx, cl, f)
}

// CreateGroup creates a new group.
func (cl *Client) CreateGroup(ctx context.Context, req *CreateGroupRequest) (*CreateGroupResponse, error) {
	return req.Do(ctx, cl)
}

// CreateGroupAsync creates a new group.
func (cl *Client) CreateGroupAsync(ctx context.Context, req *CreateGroupRequest, f func(*CreateGroupResponse, error)) {
	req.Async(ctx, cl, f)
}

// DeleteFriends deletes friends by id.
func (cl *Client) DeleteFriends(ctx context.Context, ids ...string) error {
	return DeleteFriends(ids...).Do(ctx, cl)
}

// DeleteFriendsAsync deletes friends by id.
func (cl *Client) DeleteFriendsAsync(ctx context.Context, ids []string, f func(error)) {
	DeleteFriends(ids...).Async(ctx, cl, f)
}

// DeleteFriendsUsernames deletes friends by username.
func (cl *Client) DeleteFriendsUsernames(ctx context.Context, usernames ...string) error {
	return DeleteFriends().WithUsernames(usernames...).Do(ctx, cl)
}

// DeleteFriendsUsernamesAsync deletes friends by username.
func (cl *Client) DeleteFriendsUsernamesAsync(ctx context.Context, usernames []string, f func(error)) {
	DeleteFriends().WithUsernames(usernames...).Async(ctx, cl, f)
}

// DeleteGroup deletes a group.
func (cl *Client) DeleteGroup(ctx context.Context, groupId string) error {
	return DeleteGroup(groupId).Do(ctx, cl)
}

// DeleteGroupAsync deletes a group.
func (cl *Client) DeleteGroupAsync(ctx context.Context, groupId string, f func(error)) {
	DeleteGroup(groupId).Async(ctx, cl, f)
}

// DeleteLeaderboardRecord deletes a leaderboardRecord.
func (cl *Client) DeleteLeaderboardRecord(ctx context.Context, leaderboardId string) error {
	return DeleteLeaderboardRecord(leaderboardId).Do(ctx, cl)
}

// DeleteLeaderboardRecordAsync deletes a leaderboardRecord.
func (cl *Client) DeleteLeaderboardRecordAsync(ctx context.Context, leaderboardRecordId string, f func(error)) {
	DeleteLeaderboardRecord(leaderboardRecordId).Async(ctx, cl, f)
}

// DeleteNotifications deletes notifications.
func (cl *Client) DeleteNotifications(ctx context.Context, ids ...string) error {
	return DeleteNotifications(ids...).Do(ctx, cl)
}

// DeleteNotificationsAsync deletes notifications.
func (cl *Client) DeleteNotificationsAsync(ctx context.Context, ids []string, f func(error)) {
	DeleteNotifications(ids...).Async(ctx, cl, f)
}

// DeleteStorageObjects deletes storage objects.
func (cl *Client) DeleteStorageObjects(ctx context.Context, req *DeleteStorageObjectsRequest) error {
	return req.Do(ctx, cl)
}

// DeleteStorageObjectsAsync deletes storage objects.
func (cl *Client) DeleteStorageObjectsAsync(ctx context.Context, req *DeleteStorageObjectsRequest, f func(error)) {
	req.Async(ctx, cl, f)
}

// DemoteGroupUsers demotes users from a group.
func (cl *Client) DemoteGroupUsers(ctx context.Context, groupId string, userIds ...string) error {
	return DemoteGroupUsers(groupId, userIds...).Do(ctx, cl)
}

// DemoteGroupUsersAsync demotes users from a group.
func (cl *Client) DemoteGroupUsersAsync(ctx context.Context, groupId string, userIds []string, f func(error)) {
	DemoteGroupUsers(groupId, userIds...).Async(ctx, cl, f)
}

// Event sends an event.
func (cl *Client) Event(ctx context.Context, req *EventRequest) error {
	return req.Do(ctx, cl)
}

// EventAsync sends an event.
func (cl *Client) EventAsync(ctx context.Context, req *EventRequest, f func(error)) {
	req.Async(ctx, cl, f)
}

// ImportFacebookFriends imports Facebook friends.
func (cl *Client) ImportFacebookFriends(ctx context.Context, token string, reset bool) error {
	return ImportFacebookFriends(token).WithReset(reset).Do(ctx, cl)
}

// ImportFacebookFriendsAsync imports Facebook friends.
func (cl *Client) ImportFacebookFriendsAsync(ctx context.Context, token string, reset bool, f func(error)) {
	ImportFacebookFriends(token).WithReset(reset).Async(ctx, cl, f)
}

// ImportSteamFriends imports Steam friends.
func (cl *Client) ImportSteamFriends(ctx context.Context, token string, reset bool) error {
	return ImportSteamFriends(token).WithReset(reset).Do(ctx, cl)
}

// ImportSteamFriendsAsync imports Steam friends.
func (cl *Client) ImportSteamFriendsAsync(ctx context.Context, token string, reset bool, f func(error)) {
	ImportSteamFriends(token).WithReset(reset).Async(ctx, cl, f)
}

// Users retrieves users by id.
func (cl *Client) Users(ctx context.Context, ids ...string) (*UsersResponse, error) {
	return Users(ids...).Do(ctx, cl)
}

// UsersAsync retrieves users by id.
func (cl *Client) UsersAsync(ctx context.Context, ids []string, f func(*UsersResponse, error)) {
	Users(ids...).Async(ctx, cl, f)
}

// UsersUsernames retrieves users by username.
func (cl *Client) UsersUsernames(ctx context.Context, usernames ...string) (*UsersResponse, error) {
	return Users().WithUsernames(usernames...).Do(ctx, cl)
}

// UsersUsernamesAsync retrieves users by username.
func (cl *Client) UsersUsernamesAsync(ctx context.Context, usernames []string, f func(*UsersResponse, error)) {
	Users().WithUsernames(usernames...).Async(ctx, cl, f)
}

// JoinGroup joins a group.
func (cl *Client) JoinGroup(ctx context.Context, groupId string) error {
	return JoinGroup(groupId).Do(ctx, cl)
}

// JoinGroupAsync joins a group.
func (cl *Client) JoinGroupAsync(ctx context.Context, groupId string, f func(error)) {
	JoinGroup(groupId).Async(ctx, cl, f)
}

// JoinTournament joins a tournament.
func (cl *Client) JoinTournament(ctx context.Context, tournamentId string) error {
	return JoinTournament(tournamentId).Do(ctx, cl)
}

// JoinTournamentAsync joins a tournament.
func (cl *Client) JoinTournamentAsync(ctx context.Context, tournamentId string, f func(error)) {
	JoinTournament(tournamentId).Async(ctx, cl, f)
}

// KickGroupUsers kicks users from a group or decline their join request.
func (cl *Client) KickGroupUsers(ctx context.Context, groupId string, userIds ...string) error {
	return KickGroupUsers(groupId, userIds...).Do(ctx, cl)
}

// KickGroupUsersAsync kicks users froum a group or decline their join request.
func (cl *Client) KickGroupUsersAsync(ctx context.Context, groupId string, userIds []string, f func(error)) {
	KickGroupUsers(groupId, userIds...).Async(ctx, cl, f)
}

// LeaveGroup leaves a group.
func (cl *Client) LeaveGroup(ctx context.Context, groupId string) error {
	return LeaveGroup(groupId).Do(ctx, cl)
}

// LeaveGroupAsync leaves a group.
func (cl *Client) LeaveGroupAsync(ctx context.Context, groupId string, f func(error)) {
	LeaveGroup(groupId).Async(ctx, cl, f)
}

// ChannelMessages retrieves a channel's messages.
func (cl *Client) ChannelMessages(ctx context.Context, req *ChannelMessagesRequest) (*ChannelMessagesResponse, error) {
	return req.Do(ctx, cl)
}

// ChannelMessagesAsync retrieves a channel's messages.
func (cl *Client) ChannelMessagesAsync(ctx context.Context, req *ChannelMessagesRequest, f func(*ChannelMessagesResponse, error)) {
	req.Async(ctx, cl, f)
}

// GroupUsers retrieves a group's users.
func (cl *Client) GroupUsers(ctx context.Context, req *GroupUsersRequest) (*GroupUsersResponse, error) {
	return req.Do(ctx, cl)
}

// GroupUsersAsync retrieves a group's users.
func (cl *Client) GroupUsersAsync(ctx context.Context, req *GroupUsersRequest, f func(*GroupUsersResponse, error)) {
	req.Async(ctx, cl, f)
}

// UserGroups retrieves a user's groups.
func (cl *Client) UserGroups(ctx context.Context, userId string) (*UserGroupsResponse, error) {
	return UserGroups(userId).Do(ctx, cl)
}

// UserGroupsAsync retrieves a user's groups.
func (cl *Client) UserGroupsAsync(ctx context.Context, userId string, f func(*UserGroupsResponse, error)) {
	UserGroups(userId).Async(ctx, cl, f)
}

// Groups retrieves groups.
func (cl *Client) Groups(ctx context.Context, req *GroupsRequest) (*GroupsResponse, error) {
	return req.Do(ctx, cl)
}

// GroupsAsync retrieves groups.
func (cl *Client) GroupsAsync(ctx context.Context, req *GroupsRequest, f func(*GroupsResponse, error)) {
	req.Async(ctx, cl, f)
}

// LinkApple adds a Apple token to the user's account.
func (cl *Client) LinkApple(ctx context.Context, token string) error {
	return LinkApple(token).Do(ctx, cl)
}

// LinkApple adds a Apple token to the user's account.
func (cl *Client) LinkAppleAsync(ctx context.Context, token string, f func(error)) {
	LinkApple(token).Async(ctx, cl, f)
}

// LinkCustom adds a custom id to the user's account.
func (cl *Client) LinkCustom(ctx context.Context, id string) error {
	return LinkCustom(id).Do(ctx, cl)
}

// LinkCustom adds a custom id to the user's account.
func (cl *Client) LinkCustomAsync(ctx context.Context, id string, f func(error)) {
	LinkCustom(id).Async(ctx, cl, f)
}

// LinkDevice adds a device id to the user's account.
func (cl *Client) LinkDevice(ctx context.Context, id string) error {
	return LinkDevice(id).Do(ctx, cl)
}

// LinkDevice adds a device id to the user's account.
func (cl *Client) LinkDeviceAsync(ctx context.Context, id string, f func(error)) {
	LinkDevice(id).Async(ctx, cl, f)
}

// LinkEmail adds a email/password to the user's account.
func (cl *Client) LinkEmail(ctx context.Context, email, password string) error {
	return LinkEmail(email, password).Do(ctx, cl)
}

// LinkEmail adds a email/password to the user's account.
func (cl *Client) LinkEmailAsync(ctx context.Context, email, password string, f func(error)) {
	LinkEmail(email, password).Async(ctx, cl, f)
}

// LinkFacebook adds a Facebook token to the user's account.
func (cl *Client) LinkFacebook(ctx context.Context, token string, sync bool) error {
	return LinkFacebook(token).WithSync(sync).Do(ctx, cl)
}

// LinkFacebook adds a Facebook token to the user's account.
func (cl *Client) LinkFacebookAsync(ctx context.Context, token string, sync bool, f func(error)) {
	LinkFacebook(token).WithSync(sync).Async(ctx, cl, f)
}

// LinkFacebookInstantGame adds a Facebook Instant Game signedPlayerInfo to the
// user's account.
func (cl *Client) LinkFacebookInstantGame(ctx context.Context, signedPlayerInfo string) error {
	return LinkFacebookInstantGame(signedPlayerInfo).Do(ctx, cl)
}

// LinkFacebookInstantGame adds a Facebook Instant Game signedPlayerInfo to the
// user's account.
func (cl *Client) LinkFacebookInstantGameAsync(ctx context.Context, signedPlayerInfo string, f func(error)) {
	LinkFacebookInstantGame(signedPlayerInfo).Async(ctx, cl, f)
}

// LinkGoogle adds a Google token to the user's account.
func (cl *Client) LinkGoogle(ctx context.Context, token string) error {
	return LinkGoogle(token).Do(ctx, cl)
}

// LinkGoogle adds a Google token to the user's account.
func (cl *Client) LinkGoogleAsync(ctx context.Context, token string, f func(error)) {
	LinkGoogle(token).Async(ctx, cl, f)
}

// LinkGameCenter adds a Apple GameCenter token to the user's account.
func (cl *Client) LinkGameCenter(ctx context.Context, req *LinkGameCenterRequest) error {
	return req.Do(ctx, cl)
}

// LinkGameCenter adds a Apple GameCenter token to the user's account.
func (cl *Client) LinkGameCenterAsync(ctx context.Context, req *LinkGameCenterRequest, f func(error)) {
	req.Async(ctx, cl, f)
}

// LinkSteam adds a Steam token to the user's account.
func (cl *Client) LinkSteam(ctx context.Context, token string, sync bool) error {
	return LinkSteam(token).
		WithSync(sync).
		Do(ctx, cl)
}

// LinkSteam adds a Steam token to the user's account.
func (cl *Client) LinkSteamAsync(ctx context.Context, token string, sync bool, f func(error)) {
	LinkSteam(token).
		WithSync(sync).
		Async(ctx, cl, f)
}

// Friends retrieves friends.
func (cl *Client) Friends(ctx context.Context, req *FriendsRequest) (*FriendsResponse, error) {
	return req.Do(ctx, cl)
}

// FriendsAsync retrieves friends.
func (cl *Client) FriendsAsync(ctx context.Context, req *FriendsRequest, f func(*FriendsResponse, error)) {
	req.Async(ctx, cl, f)
}

// LeaderboardRecords retrieves leaderboard records.
func (cl *Client) LeaderboardRecords(ctx context.Context, req *LeaderboardRecordsRequest) (*LeaderboardRecordsResponse, error) {
	return req.Do(ctx, cl)
}

// LeaderboardRecordsAsync retrieves leaderboardRecords.
func (cl *Client) LeaderboardRecordsAsync(ctx context.Context, req *LeaderboardRecordsRequest, f func(*LeaderboardRecordsResponse, error)) {
	req.Async(ctx, cl, f)
}

// LeaderboardRecordsAroundOwner retrieves leaderboard records around owner.
func (cl *Client) LeaderboardRecordsAroundOwner(ctx context.Context, req *LeaderboardRecordsAroundOwnerRequest) (*LeaderboardRecordsResponse, error) {
	return req.Do(ctx, cl)
}

// LeaderboardRecordsAroundOwnerAsync retrieves leaderboard records around
// owner.
func (cl *Client) LeaderboardRecordsAroundOwnerAsync(ctx context.Context, req *LeaderboardRecordsAroundOwnerRequest, f func(*LeaderboardRecordsResponse, error)) {
	req.Async(ctx, cl, f)
}

// Matches retrieves matches.
func (cl *Client) Matches(ctx context.Context, req *MatchesRequest) (*MatchesResponse, error) {
	return req.Do(ctx, cl)
}

// MatchesAsync retrieves matches.
func (cl *Client) MatchesAsync(ctx context.Context, req *MatchesRequest, f func(*MatchesResponse, error)) {
	req.Async(ctx, cl, f)
}

// Notifications retrieves notifications.
func (cl *Client) Notifications(ctx context.Context, req *NotificationsRequest) (*NotificationsResponse, error) {
	return req.Do(ctx, cl)
}

// NotificationsAsync retrieves notifications.
func (cl *Client) NotificationsAsync(ctx context.Context, req *NotificationsRequest, f func(*NotificationsResponse, error)) {
	req.Async(ctx, cl, f)
}

// StorageObjects retrieves storage objects.
func (cl *Client) StorageObjects(ctx context.Context, req *StorageObjectsRequest) (*StorageObjectsResponse, error) {
	return req.Do(ctx, cl)
}

// StorageObjectsAsync retrieves storage objects.
func (cl *Client) StorageObjectsAsync(ctx context.Context, req *StorageObjectsRequest, f func(*StorageObjectsResponse, error)) {
	req.Async(ctx, cl, f)
}

// Subscription retrieves subscription by product id.
func (cl *Client) Subscription(ctx context.Context, productId string) (*SubscriptionResponse, error) {
	return Subscription(productId).Do(ctx, cl)
}

// SubscriptionAsync retrieves subscription by product id.
func (cl *Client) SubscriptionAsync(ctx context.Context, productId string, f func(*SubscriptionResponse, error)) {
	Subscription(productId).Async(ctx, cl, f)
}

// Subscriptions retrieves subscriptions.
func (cl *Client) Subscriptions(ctx context.Context, req *SubscriptionsRequest) (*SubscriptionsResponse, error) {
	return req.Do(ctx, cl)
}

// SubscriptionsAsync retrieves subscriptions.
func (cl *Client) SubscriptionsAsync(ctx context.Context, req *SubscriptionsRequest, f func(*SubscriptionsResponse, error)) {
	req.Async(ctx, cl, f)
}

// Tournaments retrieves tournaments.
func (cl *Client) Tournaments(ctx context.Context, req *TournamentsRequest) (*TournamentsResponse, error) {
	return req.Do(ctx, cl)
}

// TournamentsAsync retrieves tournaments.
func (cl *Client) TournamentsAsync(ctx context.Context, req *TournamentsRequest, f func(*TournamentsResponse, error)) {
	req.Async(ctx, cl, f)
}

// TournamentRecords retrieves tournament records.
func (cl *Client) TournamentRecords(ctx context.Context, req *TournamentRecordsRequest) (*TournamentRecordsResponse, error) {
	return req.Do(ctx, cl)
}

// TournamentRecordsAsync retrieves tournament records.
func (cl *Client) TournamentRecordsAsync(ctx context.Context, req *TournamentRecordsRequest, f func(*TournamentRecordsResponse, error)) {
	req.Async(ctx, cl, f)
}

// TournamentRecordsAroundOwner retrieves tournament records around owner.
func (cl *Client) TournamentRecordsAroundOwner(ctx context.Context, req *TournamentRecordsAroundOwnerRequest) (*TournamentRecordsResponse, error) {
	return req.Do(ctx, cl)
}

// TournamentRecordsAroundOwnerAsync retrieves tournament records around owner.
func (cl *Client) TournamentRecordsAroundOwnerAsync(ctx context.Context, req *TournamentRecordsAroundOwnerRequest, f func(*TournamentRecordsResponse, error)) {
	req.Async(ctx, cl, f)
}

// PromoteGroupUsers promotes users from a group.
func (cl *Client) PromoteGroupUsers(ctx context.Context, groupId string, userIds ...string) error {
	return PromoteGroupUsers(groupId, userIds...).Do(ctx, cl)
}

// PromoteGroupUsersAsync promotes users from a group.
func (cl *Client) PromoteGroupUsersAsync(ctx context.Context, groupId string, userIds []string, f func(error)) {
	PromoteGroupUsers(groupId, userIds...).Async(ctx, cl, f)
}

// ReadStorageObjects reads storage objects.
func (cl *Client) ReadStorageObjects(ctx context.Context, req *ReadStorageObjectsRequest) (*ReadStorageObjectsResponse, error) {
	return req.Do(ctx, cl)
}

// ReadStorageObjectsAsync reads storage objects.
func (cl *Client) ReadStorageObjectsAsync(ctx context.Context, req *ReadStorageObjectsRequest, f func(*ReadStorageObjectsResponse, error)) {
	req.Async(ctx, cl, f)
}

// Rpc executes a remote procedure call.
func (cl *Client) Rpc(ctx context.Context, id string, payload, v interface{}) error {
	return Rpc(id, payload, v).Do(ctx, cl)
}

// RpcAsync executes a remote procedure call.
func (cl *Client) RpcAsync(ctx context.Context, id string, payload, v interface{}, f func(error)) {
	Rpc(id, payload, v).Async(ctx, cl, f)
}

// UnlinkApple removes a Apple token from the user's account.
func (cl *Client) UnlinkApple(ctx context.Context, token string) error {
	return UnlinkApple(token).Do(ctx, cl)
}

// UnlinkApple removes a Apple token from the user's account.
func (cl *Client) UnlinkAppleAsync(ctx context.Context, token string, f func(error)) {
	UnlinkApple(token).Async(ctx, cl, f)
}

// UnlinkCustom removes a custom id from the user's account.
func (cl *Client) UnlinkCustom(ctx context.Context, id string) error {
	return UnlinkCustom(id).Do(ctx, cl)
}

// UnlinkCustom removes a custom id from the user's account.
func (cl *Client) UnlinkCustomAsync(ctx context.Context, id string, f func(error)) {
	UnlinkCustom(id).Async(ctx, cl, f)
}

// UnlinkDevice removes a device id from the user's account.
func (cl *Client) UnlinkDevice(ctx context.Context, id string) error {
	return UnlinkDevice(id).Do(ctx, cl)
}

// UnlinkDevice removes a device id from the user's account.
func (cl *Client) UnlinkDeviceAsync(ctx context.Context, id string, f func(error)) {
	UnlinkDevice(id).Async(ctx, cl, f)
}

// UnlinkEmail removes a email/password from the user's account.
func (cl *Client) UnlinkEmail(ctx context.Context, email, password string) error {
	return UnlinkEmail(email, password).Do(ctx, cl)
}

// UnlinkEmail removes a email/password from the user's account.
func (cl *Client) UnlinkEmailAsync(ctx context.Context, email, password string, f func(error)) {
	UnlinkEmail(email, password).Async(ctx, cl, f)
}

// UnlinkFacebook removes a Facebook token from the user's account.
func (cl *Client) UnlinkFacebook(ctx context.Context, token string, sync bool) error {
	return UnlinkFacebook(token).Do(ctx, cl)
}

// UnlinkFacebook removes a Facebook token from the user's account.
func (cl *Client) UnlinkFacebookAsync(ctx context.Context, token string, sync bool, f func(error)) {
	UnlinkFacebook(token).Async(ctx, cl, f)
}

// UnlinkFacebookInstantGame removes a Facebook Instant Game signedPlayerInfo to the
// user's account.
func (cl *Client) UnlinkFacebookInstantGame(ctx context.Context, signedPlayerInfo string) error {
	return UnlinkFacebookInstantGame(signedPlayerInfo).Do(ctx, cl)
}

// UnlinkFacebookInstantGame removes a Facebook Instant Game signedPlayerInfo to the
// user's account.
func (cl *Client) UnlinkFacebookInstantGameAsync(ctx context.Context, signedPlayerInfo string, f func(error)) {
	UnlinkFacebookInstantGame(signedPlayerInfo).Async(ctx, cl, f)
}

// UnlinkGameCenter removes a Apple GameCenter token from the user's account.
func (cl *Client) UnlinkGameCenter(ctx context.Context, req *UnlinkGameCenterRequest) error {
	return req.Do(ctx, cl)
}

// UnlinkGameCenter removes a Apple GameCenter token from the user's account.
func (cl *Client) UnlinkGameCenterAsync(ctx context.Context, req *UnlinkGameCenterRequest, f func(error)) {
	req.Async(ctx, cl, f)
}

// UnlinkGoogle removes a Google token from the user's account.
func (cl *Client) UnlinkGoogle(ctx context.Context, token string) error {
	return UnlinkGoogle(token).Do(ctx, cl)
}

// UnlinkGoogle removes a Google token from the user's account.
func (cl *Client) UnlinkGoogleAsync(ctx context.Context, token string, f func(error)) {
	UnlinkGoogle(token).Async(ctx, cl, f)
}

// UnlinkSteam removes a Steam token from the user's account.
func (cl *Client) UnlinkSteam(ctx context.Context, token string, sync bool) error {
	return UnlinkSteam(token).Do(ctx, cl)
}

// UnlinkSteam removes a Steam token from the user's account.
func (cl *Client) UnlinkSteamAsync(ctx context.Context, token string, sync bool, f func(error)) {
	UnlinkSteam(token).Async(ctx, cl, f)
}

// UpdateAccount updates the user's account.
func (cl *Client) UpdateAccount(ctx context.Context, req *UpdateAccountRequest) error {
	return req.Do(ctx, cl)
}

// UpdateAccountAsync updates the user's account.
func (cl *Client) UpdateAccountAsync(ctx context.Context, req *UpdateAccountRequest, f func(error)) {
	req.Async(ctx, cl, f)
}

// UpdateGroup updates a group.
func (cl *Client) UpdateGroup(ctx context.Context, req *UpdateGroupRequest) error {
	return req.Do(ctx, cl)
}

// UpdateGroupAsync updates a group.
func (cl *Client) UpdateGroupAsync(ctx context.Context, req *UpdateGroupRequest, f func(error)) {
	req.Async(ctx, cl, f)
}

// ValidatePurchaseApple validates a Apple purchase.
func (cl *Client) ValidatePurchaseApple(ctx context.Context, receipt string, persist bool) (*ValidatePurchaseResponse, error) {
	return ValidatePurchaseApple(receipt).WithPersist(persist).Do(ctx, cl)
}

// ValidatePurchaseApple validates a Apple purchase.
func (cl *Client) ValidatePurchaseAppleAsync(ctx context.Context, receipt string, persist bool, f func(*ValidatePurchaseResponse, error)) {
	ValidatePurchaseApple(receipt).WithPersist(persist).Async(ctx, cl, f)
}

// ValidatePurchaseGoogle validates a Google purchase.
func (cl *Client) ValidatePurchaseGoogle(ctx context.Context, receipt string, persist bool) (*ValidatePurchaseResponse, error) {
	return ValidatePurchaseGoogle(receipt).WithPersist(persist).Do(ctx, cl)
}

// ValidatePurchaseGoogle validates a Google purchase.
func (cl *Client) ValidatePurchaseGoogleAsync(ctx context.Context, receipt string, persist bool, f func(*ValidatePurchaseResponse, error)) {
	ValidatePurchaseGoogle(receipt).WithPersist(persist).Async(ctx, cl, f)
}

// ValidatePurchaseHuawei validates a Huawei purchase.
func (cl *Client) ValidatePurchaseHuawei(ctx context.Context, purchase, signature string, persist bool) (*ValidatePurchaseResponse, error) {
	return ValidatePurchaseHuawei(purchase, signature).WithPersist(persist).Do(ctx, cl)
}

// ValidatePurchaseHuawei validates a Huawei purchase.
func (cl *Client) ValidatePurchaseHuaweiAsync(ctx context.Context, purchase, signature string, persist bool, f func(*ValidatePurchaseResponse, error)) {
	ValidatePurchaseHuawei(purchase, signature).WithPersist(persist).Async(ctx, cl, f)
}

// ValidateSubscriptionApple validates a Apple subscription.
func (cl *Client) ValidateSubscriptionApple(ctx context.Context, receipt string, persist bool) (*ValidateSubscriptionResponse, error) {
	return ValidateSubscriptionApple(receipt).WithPersist(persist).Do(ctx, cl)
}

// ValidateSubscriptionApple validates a Apple subscription.
func (cl *Client) ValidateSubscriptionAppleAsync(ctx context.Context, receipt string, persist bool, f func(*ValidateSubscriptionResponse, error)) {
	ValidateSubscriptionApple(receipt).WithPersist(persist).Async(ctx, cl, f)
}

// ValidateSubscriptionGoogle validates a Google subscription.
func (cl *Client) ValidateSubscriptionGoogle(ctx context.Context, receipt string, persist bool) (*ValidateSubscriptionResponse, error) {
	return ValidateSubscriptionGoogle(receipt).WithPersist(persist).Do(ctx, cl)
}

// ValidateSubscriptionGoogle validates a Google subscription.
func (cl *Client) ValidateSubscriptionGoogleAsync(ctx context.Context, receipt string, persist bool, f func(*ValidateSubscriptionResponse, error)) {
	ValidateSubscriptionGoogle(receipt).WithPersist(persist).Async(ctx, cl, f)
}

// WriteLeaderboardRecord writes a leaderboard record.
func (cl *Client) WriteLeaderboardRecord(ctx context.Context, req *WriteLeaderboardRecordRequest) (*WriteLeaderboardRecordResponse, error) {
	return req.Do(ctx, cl)
}

// WriteLeaderboardRecordAsync writes a leaderboard record.
func (cl *Client) WriteLeaderboardRecordAsync(ctx context.Context, req *WriteLeaderboardRecordRequest, f func(*WriteLeaderboardRecordResponse, error)) {
	req.Async(ctx, cl, f)
}

// WriteStorageObjects writes a storage objects.
func (cl *Client) WriteStorageObjects(ctx context.Context, req *WriteStorageObjectsRequest) (*WriteStorageObjectsResponse, error) {
	return req.Do(ctx, cl)
}

// WriteStorageObjectsAsync writes a storage objects.
func (cl *Client) WriteStorageObjectsAsync(ctx context.Context, req *WriteStorageObjectsRequest, f func(*WriteStorageObjectsResponse, error)) {
	req.Async(ctx, cl, f)
}

// WriteTournamentRecord writes a tournament record.
func (cl *Client) WriteTournamentRecord(ctx context.Context, req *WriteTournamentRecordRequest) (*WriteTournamentRecordResponse, error) {
	return req.Do(ctx, cl)
}

// WriteTournamentRecordAsync writes a tournament record.
func (cl *Client) WriteTournamentRecordAsync(ctx context.Context, req *WriteTournamentRecordRequest, f func(*WriteTournamentRecordResponse, error)) {
	req.Async(ctx, cl, f)
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

// WithRefreshAuto is a nakama client option to set whether or not to
// automatically refresh the session.
func WithRefreshAuto(refreshAuto bool) Option {
	return func(cl *Client) {
		cl.refreshAuto = refreshAuto
	}
}

// WithExpiryGrace is a nakama client option to set the expiry grace used for
// session refresh.
func WithExpiryGrace(expiryGrace time.Duration) Option {
	return func(cl *Client) {
		cl.expiryGrace = expiryGrace
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

// WithLogger is a nakama client option to set a logger.
func WithLogger(f func(string, ...interface{})) Option {
	return func(cl *Client) {
		cl.logf = f
	}
}

// ParseTokenExpiry parse the exp field on a jwt token.
func ParseTokenExpiry(tokenstr, typ string, grace time.Duration) (time.Time, time.Time, error) {
	if tokenstr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("empty %s token", typ)
	}
	// split
	token := strings.Split(tokenstr, ".")
	if len(token) != 3 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid %s token jwt encoding", typ)
	}
	// decode
	buf, err := base64.RawStdEncoding.DecodeString(token[1])
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid %s token encoding: %w", typ, err)
	}
	// unmarshal
	var v struct {
		Exp int64 `json:"exp"`
	}
	switch err := json.NewDecoder(bytes.NewReader(buf)).Decode(&v); {
	case err != nil:
		return time.Time{}, time.Time{}, fmt.Errorf("cannot decode %s token: %w", typ, err)
	case v.Exp == 0:
		return time.Time{}, time.Time{}, fmt.Errorf("%s token expiry cannot be 0", typ)
	}
	// check
	expiry := time.Unix(v.Exp, 0)
	expiryGraced := expiry.Add(-grace)
	now := time.Now()
	switch {
	case now.After(expiry):
		return time.Time{}, time.Time{}, fmt.Errorf("%s token expiry (%s [%d]) is in the past", typ, expiry, v.Exp)
	case grace != 0 && now.After(expiryGraced):
		return time.Time{}, time.Time{}, fmt.Errorf("%s token expiry (%s [%d]) is after the grace expiry (%s)", typ, expiry, v.Exp, grace)
	}
	return expiry, expiryGraced, nil
}

// ClientError is a client error.
type ClientError struct {
	StatusCode int
	Code       codes.Code `json:"code"`
	Message    string     `json:"message"`
}

// NewClientErrorFromReader reads a client error from a reader.
func NewClientErrorFromReader(statusCode int, r io.Reader) error {
	dec := json.NewDecoder(r)
	err := &ClientError{
		StatusCode: statusCode,
	}
	if e := dec.Decode(err); e != nil {
		return fmt.Errorf("status %d != 200 (and unable to decode error: %w)", statusCode, e)
	}
	return err
}

// Error satisfies the error interface.
func (err *ClientError) Error() string {
	return fmt.Sprintf("http status %d != 200: %s: %s", err.StatusCode, err.Code, err.Message)
}
