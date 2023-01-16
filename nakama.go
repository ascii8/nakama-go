// Package nakama is a nakama http and realtime websocket client.
package nakama

//go:generate protoc -I. --go_out=. --go_opt=paths=source_relative nakama.proto realtime.proto

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// Healthcheck creates a new healthcheck request.
func Healthcheck() *HealthcheckRequest {
	return &HealthcheckRequest{}
}

// Do executes the healthcheck request against the context and client.
func (req *HealthcheckRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "GET", "healthcheck", false, nil, nil, nil)
}

// Async executes the request against the context and client.
func (req *HealthcheckRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// Account creates a request to retrieve the user's account.
func Account() *AccountRequest {
	return &AccountRequest{}
}

// Do executes the request against the context and client.
func (req *AccountRequest) Do(ctx context.Context, cl *Client) (*AccountResponse, error) {
	res := new(AccountResponse)
	if err := cl.Do(ctx, "GET", "v2/account", true, nil, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AccountRequest) Async(ctx context.Context, cl *Client, f func(*AccountResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// UpdateAccount creates a request to update the user's account.
func UpdateAccount() *UpdateAccountRequest {
	return &UpdateAccountRequest{}
}

// WithUsername sets the username on the request.
func (req *UpdateAccountRequest) WithUsername(username string) *UpdateAccountRequest {
	req.Username = wrapperspb.String(username)
	return req
}

// WithDisplayName sets the displayName on the request.
func (req *UpdateAccountRequest) WithDisplayName(displayName string) *UpdateAccountRequest {
	req.DisplayName = wrapperspb.String(displayName)
	return req
}

// WithAvatarUrl sets the avatarUrl on the request.
func (req *UpdateAccountRequest) WithAvatarUrl(avatarUrl string) *UpdateAccountRequest {
	req.AvatarUrl = wrapperspb.String(avatarUrl)
	return req
}

// WithLangTag sets the langTag on the request.
func (req *UpdateAccountRequest) WithLangTag(langTag string) *UpdateAccountRequest {
	req.LangTag = wrapperspb.String(langTag)
	return req
}

// WithLocation sets the location on the request.
func (req *UpdateAccountRequest) WithLocation(location string) *UpdateAccountRequest {
	req.Location = wrapperspb.String(location)
	return req
}

// WithTimezone sets the timezone on the request.
func (req *UpdateAccountRequest) WithTimezone(timezone string) *UpdateAccountRequest {
	req.Timezone = wrapperspb.String(timezone)
	return req
}

// Do executes the request against the context and client.
func (req *UpdateAccountRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "PUT", "v2/account", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UpdateAccountRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// AuthenticateApple creates a request to authenticate a user with an Apple
// token.
func AuthenticateApple(token string) *AuthenticateAppleRequest {
	return &AuthenticateAppleRequest{
		Account: &AccountApple{
			Token: token,
		},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateAppleRequest) WithCreate(create bool) *AuthenticateAppleRequest {
	req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateAppleRequest) WithUsername(username string) *AuthenticateAppleRequest {
	req.Username = username
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateAppleRequest) WithVars(vars map[string]string) *AuthenticateAppleRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateAppleRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.Create != nil {
		query.Set("create", strconv.FormatBool(req.Create.Value))
	}
	if req.Username != "" {
		query.Set("username", req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/apple", false, query, req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AuthenticateAppleRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// AuthenticateCustom creates a request to authenicate a user id against the
// server.
func AuthenticateCustom(id string) *AuthenticateCustomRequest {
	return &AuthenticateCustomRequest{
		Account: &AccountCustom{
			Id: id,
		},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateCustomRequest) WithCreate(create bool) *AuthenticateCustomRequest {
	req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateCustomRequest) WithUsername(username string) *AuthenticateCustomRequest {
	req.Username = username
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateCustomRequest) WithVars(vars map[string]string) *AuthenticateCustomRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateCustomRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.Create != nil {
		query.Set("create", strconv.FormatBool(req.Create.Value))
	}
	if req.Username != "" {
		query.Set("username", req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/custom", false, query, req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AuthenticateCustomRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// AuthenticateDevice creates a request to authenticate a user with a device
// id.
func AuthenticateDevice(id string) *AuthenticateDeviceRequest {
	return &AuthenticateDeviceRequest{
		Account: &AccountDevice{
			Id: id,
		},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateDeviceRequest) WithCreate(create bool) *AuthenticateDeviceRequest {
	req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateDeviceRequest) WithUsername(username string) *AuthenticateDeviceRequest {
	req.Username = username
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateDeviceRequest) WithVars(vars map[string]string) *AuthenticateDeviceRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateDeviceRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.Create != nil {
		query.Set("create", strconv.FormatBool(req.Create.Value))
	}
	if req.Username != "" {
		query.Set("username", req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/device", false, query, req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AuthenticateDeviceRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// AuthenticateEmail creates a request to authenticate a user with an email and
// password.
func AuthenticateEmail(email, password string) *AuthenticateEmailRequest {
	return &AuthenticateEmailRequest{
		Account: &AccountEmail{
			Email:    email,
			Password: password,
		},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateEmailRequest) WithCreate(create bool) *AuthenticateEmailRequest {
	req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateEmailRequest) WithUsername(username string) *AuthenticateEmailRequest {
	req.Username = username
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateEmailRequest) WithVars(vars map[string]string) *AuthenticateEmailRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateEmailRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.Create != nil {
		query.Set("create", strconv.FormatBool(req.Create.Value))
	}
	if req.Username != "" {
		query.Set("username", req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/email", false, query, req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AuthenticateEmailRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// AuthenticateFacebook creates a request to authenticate a user with a
// Facebook token.
func AuthenticateFacebook(token string) *AuthenticateFacebookRequest {
	return &AuthenticateFacebookRequest{
		Account: &AccountFacebook{
			Token: token,
		},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateFacebookRequest) WithCreate(create bool) *AuthenticateFacebookRequest {
	req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateFacebookRequest) WithUsername(username string) *AuthenticateFacebookRequest {
	req.Username = username
	return req
}

// WithSync sets the sync on the request.
func (req *AuthenticateFacebookRequest) WithSync(sync bool) *AuthenticateFacebookRequest {
	req.Sync = wrapperspb.Bool(sync)
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateFacebookRequest) WithVars(vars map[string]string) *AuthenticateFacebookRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateFacebookRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.Create != nil {
		query.Set("create", strconv.FormatBool(req.Create.Value))
	}
	if req.Username != "" {
		query.Set("username", req.Username)
	}
	if req.Sync != nil {
		query.Set("sync", strconv.FormatBool(req.Sync.Value))
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/facebook", false, query, req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AuthenticateFacebookRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// AuthenticateFacebookInstantGame creates a request to authenticate a user
// with a Facebook Instant Game token.
func AuthenticateFacebookInstantGame(signedPlayerInfo string) *AuthenticateFacebookInstantGameRequest {
	return &AuthenticateFacebookInstantGameRequest{
		Account: &AccountFacebookInstantGame{
			SignedPlayerInfo: signedPlayerInfo,
		},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateFacebookInstantGameRequest) WithCreate(create bool) *AuthenticateFacebookInstantGameRequest {
	req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateFacebookInstantGameRequest) WithUsername(username string) *AuthenticateFacebookInstantGameRequest {
	req.Username = username
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateFacebookInstantGameRequest) WithVars(vars map[string]string) *AuthenticateFacebookInstantGameRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateFacebookInstantGameRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.Create != nil {
		query.Set("create", strconv.FormatBool(req.Create.Value))
	}
	if req.Username != "" {
		query.Set("username", req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/facebookinstantgame", false, query, req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AuthenticateFacebookInstantGameRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// AuthenticateGameCenter creates a request to authenticate a user with a Apple
// GameCenter token.
func AuthenticateGameCenter() *AuthenticateGameCenterRequest {
	return &AuthenticateGameCenterRequest{
		Account: &AccountGameCenter{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateGameCenterRequest) WithCreate(create bool) *AuthenticateGameCenterRequest {
	req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateGameCenterRequest) WithUsername(username string) *AuthenticateGameCenterRequest {
	req.Username = username
	return req
}

// WithPlayerId sets the playerId on the request.
func (req *AuthenticateGameCenterRequest) WithPlayerId(playerId string) *AuthenticateGameCenterRequest {
	req.Account.PlayerId = playerId
	return req
}

// WithBundleId sets the bundleId on the request.
func (req *AuthenticateGameCenterRequest) WithBundleId(bundleId string) *AuthenticateGameCenterRequest {
	req.Account.BundleId = bundleId
	return req
}

// WithTimestampSeconds sets the timestampSeconds on the request.
func (req *AuthenticateGameCenterRequest) WithTimestampSeconds(timestampSeconds int64) *AuthenticateGameCenterRequest {
	req.Account.TimestampSeconds = timestampSeconds
	return req
}

// WithSalt sets the salt on the request.
func (req *AuthenticateGameCenterRequest) WithSalt(salt string) *AuthenticateGameCenterRequest {
	req.Account.Salt = salt
	return req
}

// WithSignature sets the signature on the request.
func (req *AuthenticateGameCenterRequest) WithSignature(signature string) *AuthenticateGameCenterRequest {
	req.Account.Signature = signature
	return req
}

// WithPublicKeyUrl sets the publicKeyUrl on the request.
func (req *AuthenticateGameCenterRequest) WithPublicKeyUrl(publicKeyUrl string) *AuthenticateGameCenterRequest {
	req.Account.PublicKeyUrl = publicKeyUrl
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateGameCenterRequest) WithVars(vars map[string]string) *AuthenticateGameCenterRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateGameCenterRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.Create != nil {
		query.Set("create", strconv.FormatBool(req.Create.Value))
	}
	if req.Username != "" {
		query.Set("username", req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/gamecenter", false, query, req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AuthenticateGameCenterRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// AuthenticateGoogle creates a request to authenicate a user with a Google
// token.
func AuthenticateGoogle(token string) *AuthenticateGoogleRequest {
	return &AuthenticateGoogleRequest{
		Account: &AccountGoogle{
			Token: token,
		},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateGoogleRequest) WithCreate(create bool) *AuthenticateGoogleRequest {
	req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateGoogleRequest) WithUsername(username string) *AuthenticateGoogleRequest {
	req.Username = username
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateGoogleRequest) WithVars(vars map[string]string) *AuthenticateGoogleRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateGoogleRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.Create != nil {
		query.Set("create", strconv.FormatBool(req.Create.Value))
	}
	if req.Username != "" {
		query.Set("username", req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/google", false, query, req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AuthenticateGoogleRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// AuthenticateSteam creates a request to authenticate a user with a Steam
// token.
func AuthenticateSteam(token string) *AuthenticateSteamRequest {
	return &AuthenticateSteamRequest{
		Account: &AccountSteam{
			Token: token,
		},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateSteamRequest) WithCreate(create bool) *AuthenticateSteamRequest {
	req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateSteamRequest) WithUsername(username string) *AuthenticateSteamRequest {
	req.Username = username
	return req
}

// WithSync sets the sync on the request.
func (req *AuthenticateSteamRequest) WithSync(sync bool) *AuthenticateSteamRequest {
	req.Sync = wrapperspb.Bool(sync)
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateSteamRequest) WithVars(vars map[string]string) *AuthenticateSteamRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateSteamRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.Create != nil {
		query.Set("create", strconv.FormatBool(req.Create.Value))
	}
	if req.Username != "" {
		query.Set("username", req.Username)
	}
	if req.Sync != nil {
		query.Set("sync", strconv.FormatBool(req.Sync.Value))
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/steam", false, query, req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *AuthenticateSteamRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// LinkAppleRequest is a request to add a Apple token to the user's account.
type LinkAppleRequest struct {
	AccountApple
}

// LinkApple creates a request to add a Apple token to the user's account.
func LinkApple(token string) *LinkAppleRequest {
	return &LinkAppleRequest{
		AccountApple: AccountApple{
			Token: token,
		},
	}
}

// WithToken sets the token on the request.
func (req *LinkAppleRequest) WithToken(token string) *LinkAppleRequest {
	req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *LinkAppleRequest) WithVars(vars map[string]string) *LinkAppleRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkAppleRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/apple", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *LinkAppleRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// LinkCustomRequest is a request to add a custom id to the user's account.
type LinkCustomRequest struct {
	AccountCustom
}

// LinkCustom creates a request to add a custom id to the user's account.
func LinkCustom(id string) *LinkCustomRequest {
	return &LinkCustomRequest{
		AccountCustom: AccountCustom{
			Id: id,
		},
	}
}

// WithVars sets the vars on the request.
func (req *LinkCustomRequest) WithVars(vars map[string]string) *LinkCustomRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkCustomRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/custom", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *LinkCustomRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// LinkDeviceRequest is a request to add a device id to a user's account.
type LinkDeviceRequest struct {
	AccountDevice
}

// LinkDevice creates a request to add a device id to a user's account.
func LinkDevice(id string) *LinkDeviceRequest {
	return &LinkDeviceRequest{
		AccountDevice: AccountDevice{
			Id: id,
		},
	}
}

// WithVars sets the vars on the request.
func (req *LinkDeviceRequest) WithVars(vars map[string]string) *LinkDeviceRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkDeviceRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/device", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *LinkDeviceRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// LinkEmailRequest is a request to add a email/password to the user's account.
type LinkEmailRequest struct {
	AccountEmail
}

// LinkEmail creates a request to add a email/password to the user's account.
func LinkEmail(email, password string) *LinkEmailRequest {
	return &LinkEmailRequest{
		AccountEmail: AccountEmail{
			Email:    email,
			Password: password,
		},
	}
}

// WithVars sets the vars on the request.
func (req *LinkEmailRequest) WithVars(vars map[string]string) *LinkEmailRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkEmailRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/email", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *LinkEmailRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// LinkFacebook creates a request to add a Facebook token to the user's
// account.
func LinkFacebook(token string) *LinkFacebookRequest {
	return &LinkFacebookRequest{
		Account: &AccountFacebook{
			Token: token,
		},
	}
}

// WithSync sets the sync on the request.
func (req *LinkFacebookRequest) WithSync(sync bool) *LinkFacebookRequest {
	req.Sync = wrapperspb.Bool(sync)
	return req
}

// WithVars sets the vars on the request.
func (req *LinkFacebookRequest) WithVars(vars map[string]string) *LinkFacebookRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkFacebookRequest) Do(ctx context.Context, cl *Client) error {
	query := url.Values{}
	if req.Sync != nil {
		query.Set("sync", strconv.FormatBool(req.Sync.Value))
	}
	return cl.Do(ctx, "POST", "v2/account/link/facebook", true, query, req.Account, nil)
}

// Async executes the request against the context and client.
func (req *LinkFacebookRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// LinkFacebookInstantGameRequest is a request to add Facebook Instant Game
// token to the user's account.
type LinkFacebookInstantGameRequest struct {
	AccountFacebookInstantGame
}

// LinkFacebookInstantGame creates a request to add Facebook Instant Game token
// to the user's account.
func LinkFacebookInstantGame(signedPlayerInfo string) *LinkFacebookInstantGameRequest {
	return &LinkFacebookInstantGameRequest{
		AccountFacebookInstantGame: AccountFacebookInstantGame{
			SignedPlayerInfo: signedPlayerInfo,
		},
	}
}

// WithVars sets the vars on the request.
func (req *LinkFacebookInstantGameRequest) WithVars(vars map[string]string) *LinkFacebookInstantGameRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkFacebookInstantGameRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/facebookinstantgame", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *LinkFacebookInstantGameRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// LinkGameCenterRequest is a request to add a Apple GameCenter token to a
// user's account.
type LinkGameCenterRequest struct {
	AccountGameCenter
}

// LinkGameCenter creates a request to add a Apple GameCenter token to a user's
// account.
func LinkGameCenter() *LinkGameCenterRequest {
	return &LinkGameCenterRequest{}
}

// WithPlayerId sets the playerId on the request.
func (req *LinkGameCenterRequest) WithPlayerId(playerId string) *LinkGameCenterRequest {
	req.PlayerId = playerId
	return req
}

// WithBundleId sets the bundleId on the request.
func (req *LinkGameCenterRequest) WithBundleId(bundleId string) *LinkGameCenterRequest {
	req.BundleId = bundleId
	return req
}

// WithTimestampSeconds sets the timestampSeconds on the request.
func (req *LinkGameCenterRequest) WithTimestampSeconds(timestampSeconds int64) *LinkGameCenterRequest {
	req.TimestampSeconds = timestampSeconds
	return req
}

// WithSalt sets the salt on the request.
func (req *LinkGameCenterRequest) WithSalt(salt string) *LinkGameCenterRequest {
	req.Salt = salt
	return req
}

// WithSignature sets the signature on the request.
func (req *LinkGameCenterRequest) WithSignature(signature string) *LinkGameCenterRequest {
	req.Signature = signature
	return req
}

// WithPublicKeyUrl sets the publicKeyUrl on the request.
func (req *LinkGameCenterRequest) WithPublicKeyUrl(publicKeyUrl string) *LinkGameCenterRequest {
	req.PublicKeyUrl = publicKeyUrl
	return req
}

// WithVars sets the vars on the request.
func (req *LinkGameCenterRequest) WithVars(vars map[string]string) *LinkGameCenterRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkGameCenterRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/gamecenter", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *LinkGameCenterRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// LinkGoogleRequest is a request to add a Google token to a user's account.
type LinkGoogleRequest struct {
	AccountGoogle
}

// LinkGoogle creates a request to add a Google token to a user's account.
func LinkGoogle(token string) *LinkGoogleRequest {
	return &LinkGoogleRequest{
		AccountGoogle: AccountGoogle{
			Token: token,
		},
	}
}

// WithVars sets the vars on the request.
func (req *LinkGoogleRequest) WithVars(vars map[string]string) *LinkGoogleRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkGoogleRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/google", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *LinkGoogleRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// LinkSteam creates a request adds a Steam token to a user's account.
func LinkSteam(token string) *LinkSteamRequest {
	return &LinkSteamRequest{
		Account: &AccountSteam{
			Token: token,
		},
	}
}

// WithSync sets the sync on the request.
func (req *LinkSteamRequest) WithSync(sync bool) *LinkSteamRequest {
	req.Sync = wrapperspb.Bool(sync)
	return req
}

// WithVars sets the vars on the request.
func (req *LinkSteamRequest) WithVars(vars map[string]string) *LinkSteamRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkSteamRequest) Do(ctx context.Context, cl *Client) error {
	query := url.Values{}
	if req.Sync != nil {
		query.Set("sync", strconv.FormatBool(req.Sync.Value))
	}
	return cl.Do(ctx, "POST", "v2/account/link/steam", true, query, req.Account, nil)
}

// Async executes the request against the context and client.
func (req *LinkSteamRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// SessionRefresh creates a request to refresh the session token.
func SessionRefresh(refreshToken string) *SessionRefreshRequest {
	return &SessionRefreshRequest{
		Token: refreshToken,
	}
}

// WithVars sets the vars on the request.
func (req *SessionRefreshRequest) WithVars(vars map[string]string) *SessionRefreshRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *SessionRefreshRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/session/refresh", false, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *SessionRefreshRequest) Async(ctx context.Context, cl *Client, f func(*SessionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// UnlinkAppleRequest is a request to remove a Apple token from a user's account.
type UnlinkAppleRequest struct {
	AccountApple
}

// UnlinkApple creates a request to remove a Apple token from a user's account.
func UnlinkApple(token string) *UnlinkAppleRequest {
	return &UnlinkAppleRequest{
		AccountApple: AccountApple{
			Token: token,
		},
	}
}

// WithVars sets the vars on the request.
func (req *UnlinkAppleRequest) WithVars(vars map[string]string) *UnlinkAppleRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkAppleRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/apple", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UnlinkAppleRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// UnlinkCustomRequest is a request to remove a custom id from the user's account.
type UnlinkCustomRequest struct {
	AccountCustom
}

// UnlinkCustom creates a request to remove a custom id from the user's account.
func UnlinkCustom(id string) *UnlinkCustomRequest {
	return &UnlinkCustomRequest{
		AccountCustom: AccountCustom{
			Id: id,
		},
	}
}

// WithVars sets the vars on the request.
func (req *UnlinkCustomRequest) WithVars(vars map[string]string) *UnlinkCustomRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkCustomRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/custom", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UnlinkCustomRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// UnlinkDeviceRequest is a request to remove a device id from a user's account.
type UnlinkDeviceRequest struct {
	AccountDevice
}

// UnlinkDevice creates a request to remove a device id from a user's account.
func UnlinkDevice(id string) *UnlinkDeviceRequest {
	return &UnlinkDeviceRequest{
		AccountDevice: AccountDevice{
			Id: id,
		},
	}
}

// WithVars sets the vars on the request.
func (req *UnlinkDeviceRequest) WithVars(vars map[string]string) *UnlinkDeviceRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkDeviceRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/device", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UnlinkDeviceRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// UnlinkEmailRequest is a request to remove a email/password from a user's account.
type UnlinkEmailRequest struct {
	AccountEmail
}

// UnlinkEmail creates a request to remove a email/password from a user's account.
func UnlinkEmail(email, password string) *UnlinkEmailRequest {
	return &UnlinkEmailRequest{
		AccountEmail: AccountEmail{
			Email:    email,
			Password: password,
		},
	}
}

// WithVars sets the vars on the request.
func (req *UnlinkEmailRequest) WithVars(vars map[string]string) *UnlinkEmailRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkEmailRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/email", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UnlinkEmailRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// UnlinkFacebookRequest is a request to remove a Facebook token from a user's account.
type UnlinkFacebookRequest struct {
	AccountFacebook
}

// UnlinkFacebook creates a request to remove a Facebook token from a user's account.
func UnlinkFacebook(token string) *UnlinkFacebookRequest {
	return &UnlinkFacebookRequest{
		AccountFacebook: AccountFacebook{
			Token: token,
		},
	}
}

// WithVars sets the vars on the request.
func (req *UnlinkFacebookRequest) WithVars(vars map[string]string) *UnlinkFacebookRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkFacebookRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/facebook", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UnlinkFacebookRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// UnlinkFacebookInstantGameRequest is a request to remove Facebook Instant
// Game signedPlayerInfo from the user's account.
type UnlinkFacebookInstantGameRequest struct {
	AccountFacebookInstantGame
}

// UnlinkFacebookInstantGame creates a request to remove Facebook Instant Game
// signedPlayerInfo from the user's account.
func UnlinkFacebookInstantGame(signedPlayerInfo string) *UnlinkFacebookInstantGameRequest {
	return &UnlinkFacebookInstantGameRequest{
		AccountFacebookInstantGame: AccountFacebookInstantGame{
			SignedPlayerInfo: signedPlayerInfo,
		},
	}
}

// WithVars sets the vars on the request.
func (req *UnlinkFacebookInstantGameRequest) WithVars(vars map[string]string) *UnlinkFacebookInstantGameRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkFacebookInstantGameRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/facebookinstantgame", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UnlinkFacebookInstantGameRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// UnlinkGameCenterRequest is a request to remove a Apple GameCenter token from
// a user's account.
type UnlinkGameCenterRequest struct {
	AccountGameCenter
}

// UnlinkGameCenter creates a request to remove a Apple GameCenter token from a
// user's account.
func UnlinkGameCenter() *UnlinkGameCenterRequest {
	return &UnlinkGameCenterRequest{}
}

// WithPlayerId sets the playerId on the request.
func (req *UnlinkGameCenterRequest) WithPlayerId(playerId string) *UnlinkGameCenterRequest {
	req.PlayerId = playerId
	return req
}

// WithBundleId sets the bundleId on the request.
func (req *UnlinkGameCenterRequest) WithBundleId(bundleId string) *UnlinkGameCenterRequest {
	req.BundleId = bundleId
	return req
}

// WithTimestampSeconds sets the timestampSeconds on the request.
func (req *UnlinkGameCenterRequest) WithTimestampSeconds(timestampSeconds int64) *UnlinkGameCenterRequest {
	req.TimestampSeconds = timestampSeconds
	return req
}

// WithSalt sets the salt on the request.
func (req *UnlinkGameCenterRequest) WithSalt(salt string) *UnlinkGameCenterRequest {
	req.Salt = salt
	return req
}

// WithSignature sets the signature on the request.
func (req *UnlinkGameCenterRequest) WithSignature(signature string) *UnlinkGameCenterRequest {
	req.Signature = signature
	return req
}

// WithPublicKeyUrl sets the publicKeyUrl on the request.
func (req *UnlinkGameCenterRequest) WithPublicKeyUrl(publicKeyUrl string) *UnlinkGameCenterRequest {
	req.PublicKeyUrl = publicKeyUrl
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkGameCenterRequest) WithVars(vars map[string]string) *UnlinkGameCenterRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkGameCenterRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/gamecenter", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UnlinkGameCenterRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// UnlinkGoogleRequest is a request to remove a Google token from a user's account.
type UnlinkGoogleRequest struct {
	AccountGoogle
}

// UnlinkGoogle creates a request to remove a Google token from a user's account.
func UnlinkGoogle(token string) *UnlinkGoogleRequest {
	return &UnlinkGoogleRequest{
		AccountGoogle: AccountGoogle{
			Token: token,
		},
	}
}

// WithVars sets the vars on the request.
func (req *UnlinkGoogleRequest) WithVars(vars map[string]string) *UnlinkGoogleRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkGoogleRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/google", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UnlinkGoogleRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// UnlinkSteamRequest is a request to remove a Steam token from a user's account.
type UnlinkSteamRequest struct {
	AccountSteam
}

// UnlinkSteam creates a request to remove a Steam token from a user's account.
func UnlinkSteam(token string) *UnlinkSteamRequest {
	return &UnlinkSteamRequest{
		AccountSteam: AccountSteam{
			Token: token,
		},
	}
}

// WithVars sets the vars on the request.
func (req *UnlinkSteamRequest) WithVars(vars map[string]string) *UnlinkSteamRequest {
	req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkSteamRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/steam", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UnlinkSteamRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// ChannelMessages creates a request to retrieve a channel's messages.
func ChannelMessages(channelId string) *ChannelMessagesRequest {
	return &ChannelMessagesRequest{
		ChannelId: channelId,
		Limit:     wrapperspb.Int32(100),
	}
}

// WithLimit sets the limit on the request.
func (req *ChannelMessagesRequest) WithLimit(limit int) *ChannelMessagesRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithForward sets the forward on the request.
func (req *ChannelMessagesRequest) WithForward(forward bool) *ChannelMessagesRequest {
	req.Forward = wrapperspb.Bool(forward)
	return req
}

// WithCursor sets the cursor on the request.
func (req *ChannelMessagesRequest) WithCursor(cursor string) *ChannelMessagesRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ChannelMessagesRequest) Do(ctx context.Context, cl *Client) (*ChannelMessagesResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Forward != nil {
		query.Set("forward", strconv.FormatBool(req.Forward.Value))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	res := new(ChannelMessagesResponse)
	if err := cl.Do(ctx, "GET", "v2/channel/"+req.ChannelId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *ChannelMessagesRequest) Async(ctx context.Context, cl *Client, f func(*ChannelMessagesResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// Event creates a request to send an event.
func Event(name string) *EventRequest {
	return &EventRequest{
		Name: name,
	}
}

// WithProperties sets the properties on the request.
func (req *EventRequest) WithProperties(properties map[string]string) *EventRequest {
	req.Properties = properties
	return req
}

// WithTimestamp sets the timestamp on the request.
func (req *EventRequest) WithTimestamp(t time.Time) *EventRequest {
	req.Timestamp = timestamppb.New(t)
	return req
}

// WithExternal sets the external on the request.
func (req *EventRequest) WithExternal(external bool) *EventRequest {
	req.External = external
	return req
}

// Do executes the request against the context and client.
func (req *EventRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/event", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *EventRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// Friends creates a request to retrieve friends.
func Friends() *FriendsRequest {
	return &FriendsRequest{
		Limit: wrapperspb.Int32(100),
	}
}

// WithLimit sets the limit on the request.
func (req *FriendsRequest) WithLimit(limit int) *FriendsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithState sets the state on the request.
func (req *FriendsRequest) WithState(state FriendState) *FriendsRequest {
	req.State = wrapperspb.Int32(int32(state))
	return req
}

// WithCursor sets the cursor on the request.
func (req *FriendsRequest) WithCursor(cursor string) *FriendsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *FriendsRequest) Do(ctx context.Context, cl *Client) (*FriendsResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.State != nil {
		query.Set("state", strconv.FormatInt(int64(req.State.Value), 10))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	res := new(FriendsResponse)
	if err := cl.Do(ctx, "GET", "v2/friend", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *FriendsRequest) Async(ctx context.Context, cl *Client, f func(*FriendsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// DeleteFriends creates a request to delete friends by ID or username.
func DeleteFriends(ids ...string) *DeleteFriendsRequest {
	return &DeleteFriendsRequest{
		Ids: ids,
	}
}

// WithUsernames sets the Usernames on the request.
func (req *DeleteFriendsRequest) WithUsernames(usernames ...string) *DeleteFriendsRequest {
	req.Usernames = usernames
	return req
}

// Do executes the request against the context and client.
func (req *DeleteFriendsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/friend", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *DeleteFriendsRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// AddFriends creates a new request to add friends by ID or username.
func AddFriends(ids ...string) *AddFriendsRequest {
	return &AddFriendsRequest{
		Ids: ids,
	}
}

// WithUsernames sets the Usernames on the request.
func (req *AddFriendsRequest) WithUsernames(usernames ...string) *AddFriendsRequest {
	req.Usernames = usernames
	return req
}

// Do executes the request against the context and client.
func (req *AddFriendsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/friend", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *AddFriendsRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// BlockFriends creates a request to block friends by ID or username.
func BlockFriends(ids ...string) *BlockFriendsRequest {
	return &BlockFriendsRequest{
		Ids: ids,
	}
}

// WithUsernames sets the Usernames on the request.
func (req *BlockFriendsRequest) WithUsernames(usernames ...string) *BlockFriendsRequest {
	req.Usernames = usernames
	return req
}

// Do executes the request against the context and client.
func (req *BlockFriendsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/friend/block", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *BlockFriendsRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// ImportFacebookFriends creates a request to import Facebook friends.
func ImportFacebookFriends(token string) *ImportFacebookFriendsRequest {
	return &ImportFacebookFriendsRequest{
		Account: &AccountFacebook{
			Token: token,
		},
	}
}

// WithReset sets the reset on the request.
func (req *ImportFacebookFriendsRequest) WithReset(reset bool) *ImportFacebookFriendsRequest {
	req.Reset_ = wrapperspb.Bool(reset)
	return req
}

// WithVars sets the vars on the request.
func (req *ImportFacebookFriendsRequest) WithVars(vars map[string]string) *ImportFacebookFriendsRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *ImportFacebookFriendsRequest) Do(ctx context.Context, cl *Client) error {
	query := url.Values{}
	if req.Reset_ != nil {
		query.Set("reset", strconv.FormatBool(req.Reset_.Value))
	}
	return cl.Do(ctx, "POST", "v2/friend/facebook", true, query, req.Account, nil)
}

// Async executes the request against the context and client.
func (req *ImportFacebookFriendsRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// ImportSteamFriends creates a request to import Steam friends.
func ImportSteamFriends(token string) *ImportSteamFriendsRequest {
	return &ImportSteamFriendsRequest{
		Account: &AccountSteam{
			Token: token,
		},
	}
}

// WithReset sets the reset on the request.
func (req *ImportSteamFriendsRequest) WithReset(reset bool) *ImportSteamFriendsRequest {
	req.Reset_ = wrapperspb.Bool(reset)
	return req
}

// WithVars sets the vars on the request.
func (req *ImportSteamFriendsRequest) WithVars(vars map[string]string) *ImportSteamFriendsRequest {
	req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *ImportSteamFriendsRequest) Do(ctx context.Context, cl *Client) error {
	query := url.Values{}
	if req.Reset_ != nil {
		query.Set("reset", strconv.FormatBool(req.Reset_.Value))
	}
	return cl.Do(ctx, "POST", "v2/friend/steam", true, query, req.Account, nil)
}

// Async executes the request against the context and client.
func (req *ImportSteamFriendsRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// Groups creates a request to retrieve groups.
func Groups() *GroupsRequest {
	return &GroupsRequest{}
}

// WithName sets the name on the request.
func (req *GroupsRequest) WithName(name string) *GroupsRequest {
	req.Name = name
	return req
}

// WithCursor sets the cursor on the request.
func (req *GroupsRequest) WithCursor(cursor string) *GroupsRequest {
	req.Cursor = cursor
	return req
}

// WithLimit sets the limit on the request.
func (req *GroupsRequest) WithLimit(limit int) *GroupsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithLangTag sets the langTag on the request.
func (req *GroupsRequest) WithLangTag(langTag string) *GroupsRequest {
	req.LangTag = langTag
	return req
}

// WithMembers sets the members on the request.
func (req *GroupsRequest) WithMembers(members int) *GroupsRequest {
	req.Members = wrapperspb.Int32(int32(members))
	return req
}

// WithOpen sets the open on the request.
func (req *GroupsRequest) WithOpen(open bool) *GroupsRequest {
	req.Open = wrapperspb.Bool(open)
	return req
}

// Do executes the request against the context and client.
func (req *GroupsRequest) Do(ctx context.Context, cl *Client) (*GroupsResponse, error) {
	query := url.Values{}
	if req.Name != "" {
		query.Set("name", req.Name)
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.LangTag != "" {
		query.Set("langTag", req.LangTag)
	}
	if req.Members != nil {
		query.Set("members", strconv.FormatInt(int64(req.Members.Value), 10))
	}
	if req.Open != nil {
		query.Set("open", strconv.FormatBool(req.Open.Value))
	}
	res := new(GroupsResponse)
	if err := cl.Do(ctx, "GET", "v2/group", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *GroupsRequest) Async(ctx context.Context, cl *Client, f func(*GroupsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// CreateGroup creates a request to create a new group.
func CreateGroup() *CreateGroupRequest {
	return &CreateGroupRequest{}
}

// WithName sets the name on the request.
func (req *CreateGroupRequest) WithName(name string) *CreateGroupRequest {
	req.Name = name
	return req
}

// WithDescription sets the description on the request.
func (req *CreateGroupRequest) WithDescription(description string) *CreateGroupRequest {
	req.Description = description
	return req
}

// WithLangTag sets the langTag on the request.
func (req *CreateGroupRequest) WithLangTag(langTag string) *CreateGroupRequest {
	req.LangTag = langTag
	return req
}

// WithAvatarUrl sets the avatarUrl on the request.
func (req *CreateGroupRequest) WithAvatarUrl(avatarUrl string) *CreateGroupRequest {
	req.AvatarUrl = avatarUrl
	return req
}

// WithOpen sets the open on the request.
func (req *CreateGroupRequest) WithOpen(open bool) *CreateGroupRequest {
	req.Open = open
	return req
}

// WithMaxCount sets the maxCount on the request.
func (req *CreateGroupRequest) WithMaxCount(maxCount int) *CreateGroupRequest {
	req.MaxCount = int32(maxCount)
	return req
}

// Do executes the request against the context and client.
func (req *CreateGroupRequest) Do(ctx context.Context, cl *Client) (*Group, error) {
	res := new(CreateGroupResponse)
	if err := cl.Do(ctx, "POST", "v2/group", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *CreateGroupRequest) Async(ctx context.Context, cl *Client, f func(*CreateGroupResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// CreateGroupResponse is the create group response.
type CreateGroupResponse = Group

// DeleteGroup creates a request to delete a group.
func DeleteGroup(groupId string) *DeleteGroupRequest {
	return &DeleteGroupRequest{
		GroupId: groupId,
	}
}

// Do executes the request against the context and client.
func (req *DeleteGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/group/"+req.GroupId, true, nil, nil, nil)
}

// Async executes the request against the context and client.
func (req *DeleteGroupRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// UpdateGroup creates a request to update a group.
func UpdateGroup(groupId string) *UpdateGroupRequest {
	return &UpdateGroupRequest{
		GroupId: groupId,
	}
}

// WithName sets the name on the request.
func (req *UpdateGroupRequest) WithName(name string) *UpdateGroupRequest {
	req.Name = wrapperspb.String(name)
	return req
}

// WithDescription sets the description on the request.
func (req *UpdateGroupRequest) WithDescription(description string) *UpdateGroupRequest {
	req.Description = wrapperspb.String(description)
	return req
}

// WithLangTag sets the langTag on the request.
func (req *UpdateGroupRequest) WithLangTag(langTag string) *UpdateGroupRequest {
	req.LangTag = wrapperspb.String(langTag)
	return req
}

// WithAvatarUrl sets the avatarUrl on the request.
func (req *UpdateGroupRequest) WithAvatarUrl(avatarUrl string) *UpdateGroupRequest {
	req.AvatarUrl = wrapperspb.String(avatarUrl)
	return req
}

// WithOpen sets the open on the request.
func (req *UpdateGroupRequest) WithOpen(open bool) *UpdateGroupRequest {
	req.Open = wrapperspb.Bool(open)
	return req
}

// Do executes the request against the context and client.
func (req *UpdateGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "PUT", "v2/group/"+req.GroupId, true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *UpdateGroupRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// AddGroupUsers creates a new request to add users to a group or accepts their
// join request.
func AddGroupUsers(groupId string, userIds ...string) *AddGroupUsersRequest {
	return &AddGroupUsersRequest{
		GroupId: groupId,
		UserIds: userIds,
	}
}

// Do executes the request against the context and client.
func (req *AddGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/add", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *AddGroupUsersRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// BanGroupUsers creates a request to ban users from a group.
func BanGroupUsers(groupId string, userIds ...string) *BanGroupUsersRequest {
	return &BanGroupUsersRequest{
		GroupId: groupId,
		UserIds: userIds,
	}
}

// Do executes the request against the context and client.
func (req *BanGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/ban", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *BanGroupUsersRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// DemoteGroupUsers creates a request to demote group users.
func DemoteGroupUsers(groupId string, userIds ...string) *DemoteGroupUsersRequest {
	return &DemoteGroupUsersRequest{
		GroupId: groupId,
		UserIds: userIds,
	}
}

// Do executes the request against the context and client.
func (req *DemoteGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/demote", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *DemoteGroupUsersRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// JoinGroup creates a request to join a group.
func JoinGroup(groupId string) *JoinGroupRequest {
	return &JoinGroupRequest{
		GroupId: groupId,
	}
}

// Do executes the request against the context and client.
func (req *JoinGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/join", true, nil, nil, nil)
}

// Async executes the request against the context and client.
func (req *JoinGroupRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// KickGroupUsers creates a request to kick users from a group or decline their
// join request.
func KickGroupUsers(groupId string, userIds ...string) *KickGroupUsersRequest {
	return &KickGroupUsersRequest{
		GroupId: groupId,
		UserIds: userIds,
	}
}

// Do executes the request against the context and client.
func (req *KickGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/kick", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *KickGroupUsersRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// LeaveGroup creates a request to leave a group.
func LeaveGroup(groupId string) *LeaveGroupRequest {
	return &LeaveGroupRequest{
		GroupId: groupId,
	}
}

// Do executes the request against the context and client.
func (req *LeaveGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/leave", true, nil, nil, nil)
}

// Async executes the request against the context and client.
func (req *LeaveGroupRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// PromoteGroupUsers creates a request to promote users in a group to the next
// role up.
func PromoteGroupUsers(groupId string, userIds ...string) *PromoteGroupUsersRequest {
	return &PromoteGroupUsersRequest{
		GroupId: groupId,
		UserIds: userIds,
	}
}

// Do executes the request against the context and client.
func (req *PromoteGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/promote", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *PromoteGroupUsersRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// GroupUsers creates a request to retrieve a group's users.
func GroupUsers(groupId string) *GroupUsersRequest {
	return &GroupUsersRequest{
		GroupId: groupId,
		Limit:   wrapperspb.Int32(100),
	}
}

// WithLimit sets the limit on the request.
func (req *GroupUsersRequest) WithLimit(limit int) *GroupUsersRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithState sets the state on the request.
func (req *GroupUsersRequest) WithState(state UserRoleState) *GroupUsersRequest {
	req.State = wrapperspb.Int32(int32(state))
	return req
}

// WithCursor sets the cursor on the request.
func (req *GroupUsersRequest) WithCursor(cursor string) *GroupUsersRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *GroupUsersRequest) Do(ctx context.Context, cl *Client) (*GroupUsersResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.State != nil {
		query.Set("state", strconv.FormatInt(int64(req.State.Value), 10))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	res := new(GroupUsersResponse)
	if err := cl.Do(ctx, "GET", "v2/group/"+req.GroupId+"/user", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *GroupUsersRequest) Async(ctx context.Context, cl *Client, f func(*GroupUsersResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// ValidatePurchaseApple creates a request to validate a Apple purchase.
func ValidatePurchaseApple(receipt string) *ValidatePurchaseAppleRequest {
	return &ValidatePurchaseAppleRequest{
		Receipt: receipt,
	}
}

// WithPersist sets the persist on the request.
func (req *ValidatePurchaseAppleRequest) WithPersist(persist bool) *ValidatePurchaseAppleRequest {
	req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidatePurchaseAppleRequest) Do(ctx context.Context, cl *Client) (*ValidatePurchaseResponse, error) {
	res := new(ValidatePurchaseResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/purchase/apple", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *ValidatePurchaseAppleRequest) Async(ctx context.Context, cl *Client, f func(*ValidatePurchaseResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// ValidatePurchaseGoogle creates a request to validate a Google purchase.
func ValidatePurchaseGoogle(purchase string) *ValidatePurchaseGoogleRequest {
	return &ValidatePurchaseGoogleRequest{
		Purchase: purchase,
	}
}

// WithPersist sets the persist on the request.
func (req *ValidatePurchaseGoogleRequest) WithPersist(persist bool) *ValidatePurchaseGoogleRequest {
	req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidatePurchaseGoogleRequest) Do(ctx context.Context, cl *Client) (*ValidatePurchaseResponse, error) {
	res := new(ValidatePurchaseResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/purchase/google", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *ValidatePurchaseGoogleRequest) Async(ctx context.Context, cl *Client, f func(*ValidatePurchaseResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// ValidatePurchaseHuawei creates a request to validate a Huawei purchase.
func ValidatePurchaseHuawei(purchase, signature string) *ValidatePurchaseHuaweiRequest {
	return &ValidatePurchaseHuaweiRequest{
		Purchase:  purchase,
		Signature: signature,
	}
}

// WithPersist sets the persist on the request.
func (req *ValidatePurchaseHuaweiRequest) WithPersist(persist bool) *ValidatePurchaseHuaweiRequest {
	req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidatePurchaseHuaweiRequest) Do(ctx context.Context, cl *Client) (*ValidatePurchaseResponse, error) {
	res := new(ValidatePurchaseResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/purchase/huawei", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *ValidatePurchaseHuaweiRequest) Async(ctx context.Context, cl *Client, f func(*ValidatePurchaseResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// Subscriptions creates a request to retrieve subscriptions.
func Subscriptions(groupId string) *SubscriptionsRequest {
	return &SubscriptionsRequest{
		Limit: wrapperspb.Int32(100),
	}
}

// WithLimit sets the limit on the request.
func (req *SubscriptionsRequest) WithLimit(limit int) *SubscriptionsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCursor sets the cursor on the request.
func (req *SubscriptionsRequest) WithCursor(cursor string) *SubscriptionsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *SubscriptionsRequest) Do(ctx context.Context, cl *Client) (*SubscriptionsResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	res := new(SubscriptionsResponse)
	if err := cl.Do(ctx, "GET", "v2/iap/subscription", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *SubscriptionsRequest) Async(ctx context.Context, cl *Client, f func(*SubscriptionsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// ValidateSubscriptionApple creates a request to validate Apple subscriptions.
func ValidateSubscriptionApple(receipt string) *ValidateSubscriptionAppleRequest {
	return &ValidateSubscriptionAppleRequest{
		Receipt: receipt,
	}
}

// WithPersist sets the persist on the request.
func (req *ValidateSubscriptionAppleRequest) WithPersist(persist bool) *ValidateSubscriptionAppleRequest {
	req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidateSubscriptionAppleRequest) Do(ctx context.Context, cl *Client) (*ValidateSubscriptionResponse, error) {
	res := new(ValidateSubscriptionResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/subscription/apple", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *ValidateSubscriptionAppleRequest) Async(ctx context.Context, cl *Client, f func(*ValidateSubscriptionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// ValidateSubscriptionGoogle creates a request to validate a Google subscription.
func ValidateSubscriptionGoogle(receipt string) *ValidateSubscriptionGoogleRequest {
	return &ValidateSubscriptionGoogleRequest{
		Receipt: receipt,
	}
}

// WithPersist sets the persist on the request.
func (req *ValidateSubscriptionGoogleRequest) WithPersist(persist bool) *ValidateSubscriptionGoogleRequest {
	req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidateSubscriptionGoogleRequest) Do(ctx context.Context, cl *Client) (*ValidateSubscriptionResponse, error) {
	res := new(ValidateSubscriptionResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/subscription/google", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *ValidateSubscriptionGoogleRequest) Async(ctx context.Context, cl *Client, f func(*ValidateSubscriptionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

func Subscription(productId string) *SubscriptionRequest {
	return &SubscriptionRequest{
		ProductId: productId,
	}
}

// Do executes the request against the context and client.
func (req *SubscriptionRequest) Do(ctx context.Context, cl *Client) (*SubscriptionResponse, error) {
	res := new(SubscriptionResponse)
	if err := cl.Do(ctx, "GET", "v2/iap/subscription/"+req.ProductId, true, nil, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *SubscriptionRequest) Async(ctx context.Context, cl *Client, f func(*SubscriptionResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// SubscriptionResponse is a Subscription response.
type SubscriptionResponse = ValidatedSubscription

// LeaderboardRecords creates a request to retrieve the leaderboard records.
func LeaderboardRecords(leaderboardId string) *LeaderboardRecordsRequest {
	return &LeaderboardRecordsRequest{
		LeaderboardId: leaderboardId,
		Limit:         wrapperspb.Int32(100),
	}
}

// WithOwnerIds sets the ownerIds on the request.
func (req *LeaderboardRecordsRequest) WithOwnerIds(ownerIds ...string) *LeaderboardRecordsRequest {
	req.OwnerIds = ownerIds
	return req
}

// WithLimit sets the limit on the request.
func (req *LeaderboardRecordsRequest) WithLimit(limit int) *LeaderboardRecordsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCursor sets the cursor on the request.
func (req *LeaderboardRecordsRequest) WithCursor(cursor string) *LeaderboardRecordsRequest {
	req.Cursor = cursor
	return req
}

// WithExpiry sets the expiry on the request.
func (req *LeaderboardRecordsRequest) WithExpiry(expiry int) *LeaderboardRecordsRequest {
	req.Expiry = wrapperspb.Int64(int64(expiry))
	return req
}

// Do executes the request against the context and client.
func (req *LeaderboardRecordsRequest) Do(ctx context.Context, cl *Client) (*LeaderboardRecordsResponse, error) {
	query := url.Values{}
	if req.OwnerIds != nil {
		query.Set("ownerIds", strings.Join(req.OwnerIds, ","))
	}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	if req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.Expiry.Value), 10))
	}
	res := new(LeaderboardRecordsResponse)
	if err := cl.Do(ctx, "GET", "v2/leaderboard/"+req.LeaderboardId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *LeaderboardRecordsRequest) Async(ctx context.Context, cl *Client, f func(*LeaderboardRecordsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// DeleteLeaderboardRecord creates a request to delete a leaderboard.
func DeleteLeaderboardRecord(leaderboardId string) *DeleteLeaderboardRecordRequest {
	return &DeleteLeaderboardRecordRequest{
		LeaderboardId: leaderboardId,
	}
}

// Do executes the request against the context and client.
func (req *DeleteLeaderboardRecordRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/leaderboard/"+req.LeaderboardId, true, nil, nil, nil)
}

// Async executes the request against the context and client.
func (req *DeleteLeaderboardRecordRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// WriteLeaderboardRecord creates a request to write a leaderboard record.
func WriteLeaderboardRecord(leaderboardId string) *WriteLeaderboardRecordRequest {
	return &WriteLeaderboardRecordRequest{
		LeaderboardId: leaderboardId,
		Record:        &LeaderboardRecordWrite{},
	}
}

// WithScore sets the score on the request.
func (req *WriteLeaderboardRecordRequest) WithScore(score int64) *WriteLeaderboardRecordRequest {
	req.Record.Score = score
	return req
}

// WithSubscore sets the subscore on the request.
func (req *WriteLeaderboardRecordRequest) WithSubscore(subscore int64) *WriteLeaderboardRecordRequest {
	req.Record.Subscore = subscore
	return req
}

// WithMetadata sets the metadata on the request.
func (req *WriteLeaderboardRecordRequest) WithMetadata(metadata string) *WriteLeaderboardRecordRequest {
	req.Record.Metadata = metadata
	return req
}

// WithOperator sets the operator on the request.
func (req *WriteLeaderboardRecordRequest) WithOperator(operator OpType) *WriteLeaderboardRecordRequest {
	req.Record.Operator = operator
	return req
}

// Do executes the request against the context and client.
func (req *WriteLeaderboardRecordRequest) Do(ctx context.Context, cl *Client) (*WriteLeaderboardRecordResponse, error) {
	res := new(WriteLeaderboardRecordResponse)
	if err := cl.Do(ctx, "POST", "v2/leaderboard/"+req.LeaderboardId, true, nil, req.Record, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *WriteLeaderboardRecordRequest) Async(ctx context.Context, cl *Client, f func(*WriteLeaderboardRecordResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// WriteLeaderboardRecordResponse is the WriteLeaderboardRecord response.
type WriteLeaderboardRecordResponse = LeaderboardRecord

// LeaderboardRecordsAroundOwner creates a request to retrieve leaderboard
// records around owner.
func LeaderboardRecordsAroundOwner(leaderboardId, ownerId string) *LeaderboardRecordsAroundOwnerRequest {
	return &LeaderboardRecordsAroundOwnerRequest{
		LeaderboardId: leaderboardId,
		OwnerId:       ownerId,
		Limit:         wrapperspb.UInt32(100),
	}
}

// WithLimit sets the limit on the request.
func (req *LeaderboardRecordsAroundOwnerRequest) WithLimit(limit int) *LeaderboardRecordsAroundOwnerRequest {
	req.Limit = wrapperspb.UInt32(uint32(limit))
	return req
}

// WithExpiry sets the expiry on the request.
func (req *LeaderboardRecordsAroundOwnerRequest) WithExpiry(expiry int) *LeaderboardRecordsAroundOwnerRequest {
	req.Expiry = wrapperspb.Int64(int64(expiry))
	return req
}

// Do executes the request against the context and client.
func (req *LeaderboardRecordsAroundOwnerRequest) Do(ctx context.Context, cl *Client) (*LeaderboardRecordsResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.Expiry.Value), 10))
	}
	res := new(LeaderboardRecordsResponse)
	if err := cl.Do(ctx, "GET", "v2/leaderboard/"+req.LeaderboardId+"/owner/"+req.OwnerId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *LeaderboardRecordsAroundOwnerRequest) Async(ctx context.Context, cl *Client, f func(*LeaderboardRecordsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// Matches creates a request to retrieve matches.
func Matches() *MatchesRequest {
	return &MatchesRequest{
		Limit: wrapperspb.Int32(100),
	}
}

// WithLimit sets the limit on the request.
func (req *MatchesRequest) WithLimit(limit int) *MatchesRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithAuthoritative sets the authoritative on the request.
func (req *MatchesRequest) WithAuthoritative(authoritative bool) *MatchesRequest {
	req.Authoritative = wrapperspb.Bool(authoritative)
	return req
}

// WithLabel sets the label on the request.
func (req *MatchesRequest) WithLabel(label string) *MatchesRequest {
	req.Label = wrapperspb.String(label)
	return req
}

// WithMinSize sets the minSize on the request.
func (req *MatchesRequest) WithMinSize(minSize int) *MatchesRequest {
	req.MinSize = wrapperspb.Int32(int32(minSize))
	return req
}

// WithMaxSize sets the maxSize on the request.
func (req *MatchesRequest) WithMaxSize(maxSize int) *MatchesRequest {
	req.MaxSize = wrapperspb.Int32(int32(maxSize))
	return req
}

// WithQuery sets the query on the request.
func (req *MatchesRequest) WithQuery(query string) *MatchesRequest {
	req.Query = wrapperspb.String(query)
	return req
}

// Do executes the request against the context and client.
func (req *MatchesRequest) Do(ctx context.Context, cl *Client) (*MatchesResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Authoritative != nil {
		query.Set("authoritative", strconv.FormatBool(req.Authoritative.Value))
	}
	if req.Label != nil {
		query.Set("label", req.Label.Value)
	}
	if req.MinSize != nil {
		query.Set("minSize", strconv.FormatInt(int64(req.MinSize.Value), 10))
	}
	if req.MaxSize != nil {
		query.Set("maxSize", strconv.FormatInt(int64(req.MaxSize.Value), 10))
	}
	if req.Query != nil {
		query.Set("query", req.Query.Value)
	}
	res := new(MatchesResponse)
	if err := cl.Do(ctx, "GET", "v2/match", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *MatchesRequest) Async(ctx context.Context, cl *Client, f func(*MatchesResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// Notifications creates a request to retrieve notifications.
func Notifications() *NotificationsRequest {
	return &NotificationsRequest{
		Limit: wrapperspb.Int32(100),
	}
}

// WithLimit sets the limit on the request to retrieve notifications.
func (req *NotificationsRequest) WithLimit(limit int) *NotificationsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCacheableCursor sets the cacheableCursor on the request.
func (req *NotificationsRequest) WithCacheableCursor(cacheableCursor string) *NotificationsRequest {
	req.CacheableCursor = cacheableCursor
	return req
}

// Do executes the request against the context and client.
func (req *NotificationsRequest) Do(ctx context.Context, cl *Client) (*NotificationsResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.CacheableCursor != "" {
		query.Set("cacheableCursor", req.CacheableCursor)
	}
	res := new(NotificationsResponse)
	if err := cl.Do(ctx, "GET", "v2/notifications", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *NotificationsRequest) Async(ctx context.Context, cl *Client, f func(*NotificationsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// DeleteNotifications creates a request to delete notifications.
func DeleteNotifications(ids ...string) *DeleteNotificationsRequest {
	return &DeleteNotificationsRequest{
		Ids: ids,
	}
}

// Do executes the request against the context and client.
func (req *DeleteNotificationsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/notification", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *DeleteNotificationsRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// RpcRequest is a request/message to execute a remote procedure call.
type RpcRequest struct {
	id      string
	payload interface{}
	v       interface{}
	httpKey string
	proto   bool
	buf     []byte
	mutex   sync.Mutex
}

// Rpc creates a request to execute a remote procedure call.
func Rpc(id string, payload, v interface{}) *RpcRequest {
	return &RpcRequest{
		id:      id,
		payload: payload,
		v:       v,
	}
}

// BuildEnvelope satisfies the EnvelopeBuilder interface.
func (req *RpcRequest) BuildEnvelope() *Envelope {
	return &Envelope{
		Message: &Envelope_Rpc{
			Rpc: &RpcMsg{
				Id:      req.id,
				Payload: string(req.buf),
			},
		},
	}
}

// WithHttpKey sets the httpKey on the request.
func (req *RpcRequest) WithHttpKey(httpKey string) *RpcRequest {
	req.httpKey = httpKey
	return req
}

// WithProto sets the Protobuf encoding toggle for the realtime message.
func (req *RpcRequest) WithProto(proto bool) *RpcRequest {
	req.proto = proto
	return req
}

// Do executes the request against the context and client.
func (req *RpcRequest) Do(ctx context.Context, cl *Client) error {
	query := url.Values{}
	query.Set("unwrap", "true")
	if req.httpKey != "" {
		query.Set("http_key", req.httpKey)
	}
	return cl.Do(ctx, "POST", "v2/rpc/"+req.id, req.httpKey == "", query, req.payload, req.v)
}

// Async executes the request against the context and client.
func (req *RpcRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// Send sends the message to the connection.
func (req *RpcRequest) Send(ctx context.Context, conn *Conn) error {
	if err := req.marshal(); err != nil {
		return err
	}
	res := new(RpcMsg)
	if err := conn.Send(ctx, req, res); err != nil {
		return err
	}
	return req.unmarshal(res)
}

// SendAsync sends the message to the connection.
func (req *RpcRequest) SendAsync(ctx context.Context, conn *Conn, f func(error)) {
	go func() {
		f(req.Send(ctx, conn))
	}()
}

// marshal marshals the request.
func (req *RpcRequest) marshal() error {
	req.mutex.Lock()
	defer req.mutex.Unlock()
	if req.buf != nil {
		return nil
	}
	// protobuf encode
	if req.proto {
		msg, ok := req.payload.(proto.Message)
		if !ok {
			return fmt.Errorf("payload type %T is not a proto.Message", req.payload)
		}
		buf, err := proto.Marshal(msg)
		if err != nil {
			return err
		}
		req.buf = buf
		return nil
	}
	// json encode
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	if err := enc.Encode(req.payload); err != nil {
		return err
	}
	req.buf = buf.Bytes()
	return nil
}

// unmarshal unmarshals the response.
func (req *RpcRequest) unmarshal(msg *RpcMsg) error {
	if msg.Payload == "" {
		return nil
	}
	// protobuf decode
	if req.proto {
		v, ok := req.v.(proto.Message)
		if !ok {
			return fmt.Errorf("payload type %T is not a proto.Message", req.v)
		}
		return proto.Unmarshal([]byte(msg.Payload), v)
	}
	// json decode
	dec := json.NewDecoder(strings.NewReader(msg.Payload))
	dec.DisallowUnknownFields()
	return dec.Decode(req.v)
}

// SessionLogout creates a request to logout of the session.
func SessionLogout(token, refreshToken string) *SessionLogoutRequest {
	return &SessionLogoutRequest{
		Token:        token,
		RefreshToken: refreshToken,
	}
}

// Do executes the request against the context and client.
func (req *SessionLogoutRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/session/logout", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *SessionLogoutRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// ReadStorageObjects creates a request to read storage objects.
func ReadStorageObjects() *ReadStorageObjectsRequest {
	return &ReadStorageObjectsRequest{}
}

// WithObjectId sets the objectId on the request.
func (req *ReadStorageObjectsRequest) WithObjectId(collection, key, userId string) *ReadStorageObjectsRequest {
	req.ObjectIds = append(req.ObjectIds, &ReadStorageObjectId{
		Collection: collection,
		Key:        key,
		UserId:     userId,
	})
	return req
}

// Do executes the request against the context and client.
func (req *ReadStorageObjectsRequest) Do(ctx context.Context, cl *Client) (*ReadStorageObjectsResponse, error) {
	res := new(ReadStorageObjectsResponse)
	if err := cl.Do(ctx, "POST", "v2/storage", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *ReadStorageObjectsRequest) Async(ctx context.Context, cl *Client, f func(*ReadStorageObjectsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// WriteStorageObjects creates a request to write storage objects.
func WriteStorageObjects() *WriteStorageObjectsRequest {
	return &WriteStorageObjectsRequest{}
}

// WithObject sets the object on the request.
func (req *WriteStorageObjectsRequest) WithObject(object *WriteStorageObject) *WriteStorageObjectsRequest {
	req.Objects = append(req.Objects, object)
	return req
}

// Do executes the request against the context and client.
func (req *WriteStorageObjectsRequest) Do(ctx context.Context, cl *Client) (*WriteStorageObjectsResponse, error) {
	res := new(WriteStorageObjectsResponse)
	if err := cl.Do(ctx, "PUT", "v2/storage", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *WriteStorageObjectsRequest) Async(ctx context.Context, cl *Client, f func(*WriteStorageObjectsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// DeleteStorageObjects creates a request to delete storage objects.
func DeleteStorageObjects() *DeleteStorageObjectsRequest {
	return &DeleteStorageObjectsRequest{}
}

// WithObjectId sets the objectId on the request.
func (req *DeleteStorageObjectsRequest) WithObjectId(collection, key, version string) *DeleteStorageObjectsRequest {
	req.ObjectIds = append(req.ObjectIds, &DeleteStorageObjectId{
		Collection: collection,
		Key:        key,
		Version:    version,
	})
	return req
}

// Do executes the request against the context and client.
func (req *DeleteStorageObjectsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "PUT", "v2/storage/delete", true, nil, req, nil)
}

// Async executes the request against the context and client.
func (req *DeleteStorageObjectsRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// StorageObjects creates a request to retrieve storage objects.
func StorageObjects(collection string) *StorageObjectsRequest {
	return &StorageObjectsRequest{
		Collection: collection,
		Limit:      wrapperspb.Int32(100),
	}
}

// WithUserId sets the userId on the request.
func (req *StorageObjectsRequest) WithUserId(userId string) *StorageObjectsRequest {
	req.UserId = userId
	return req
}

// WithLimit sets the limit on the request.
func (req *StorageObjectsRequest) WithLimit(limit int) *StorageObjectsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCursor sets the cursor on the request.
func (req *StorageObjectsRequest) WithCursor(cursor string) *StorageObjectsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *StorageObjectsRequest) Do(ctx context.Context, cl *Client) (*StorageObjectsResponse, error) {
	query := url.Values{}
	if req.UserId != "" {
		query.Set("userId", req.UserId)
	}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	res := new(StorageObjectsResponse)
	if err := cl.Do(ctx, "GET", "v2/storage/"+req.Collection, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *StorageObjectsRequest) Async(ctx context.Context, cl *Client, f func(*StorageObjectsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// Tournaments creates a request to retrieve tournaments.
func Tournaments() *TournamentsRequest {
	return &TournamentsRequest{
		Limit: wrapperspb.Int32(100),
	}
}

// WithCategoryStart sets the categoryStart on the request.
func (req *TournamentsRequest) WithCategoryStart(categoryStart uint32) *TournamentsRequest {
	req.CategoryStart = wrapperspb.UInt32(categoryStart)
	return req
}

// WithCategoryEnd sets the categoryEnd on the request.
func (req *TournamentsRequest) WithCategoryEnd(categoryEnd uint32) *TournamentsRequest {
	req.CategoryEnd = wrapperspb.UInt32(categoryEnd)
	return req
}

// WithLimit sets the limit on the request.
func (req *TournamentsRequest) WithLimit(limit int) *TournamentsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithStartTime sets the startTime on the request.
func (req *TournamentsRequest) WithStartTime(startTime uint32) *TournamentsRequest {
	req.StartTime = wrapperspb.UInt32(startTime)
	return req
}

// WithEndTime sets the endTime on the request.
func (req *TournamentsRequest) WithEndTime(endTime uint32) *TournamentsRequest {
	req.EndTime = wrapperspb.UInt32(endTime)
	return req
}

// WithCursor sets the cursor on the request.
func (req *TournamentsRequest) WithCursor(cursor string) *TournamentsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *TournamentsRequest) Do(ctx context.Context, cl *Client) (*TournamentsResponse, error) {
	query := url.Values{}
	if req.CategoryStart != nil {
		query.Set("categoryStart", strconv.FormatUint(uint64(req.CategoryStart.Value), 10))
	}
	if req.CategoryEnd != nil {
		query.Set("categoryEnd", strconv.FormatUint(uint64(req.CategoryEnd.Value), 10))
	}
	if req.StartTime != nil {
		query.Set("startTime", strconv.FormatUint(uint64(req.StartTime.Value), 10))
	}
	if req.EndTime != nil {
		query.Set("endTime", strconv.FormatUint(uint64(req.EndTime.Value), 10))
	}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	res := new(TournamentsResponse)
	if err := cl.Do(ctx, "GET", "v2/tournament", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *TournamentsRequest) Async(ctx context.Context, cl *Client, f func(*TournamentsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// TournamentRecords creates a request to retrieve tournament records.
func TournamentRecords(tournamentId string) *TournamentRecordsRequest {
	return &TournamentRecordsRequest{
		TournamentId: tournamentId,
		Limit:        wrapperspb.Int32(100),
	}
}

// WithOwnerIds sets the ownerIds on the request.
func (req *TournamentRecordsRequest) WithOwnerIds(ownerIds ...string) *TournamentRecordsRequest {
	req.OwnerIds = ownerIds
	return req
}

// WithLimit sets the limit on the request.
func (req *TournamentRecordsRequest) WithLimit(limit int) *TournamentRecordsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithExpiry sets the expiry on the request.
func (req *TournamentRecordsRequest) WithExpiry(expiry int64) *TournamentRecordsRequest {
	req.Expiry = wrapperspb.Int64(expiry)
	return req
}

// WithCursor sets the cursor on the request.
func (req *TournamentRecordsRequest) WithCursor(cursor string) *TournamentRecordsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *TournamentRecordsRequest) Do(ctx context.Context, cl *Client) (*TournamentRecordsResponse, error) {
	query := url.Values{}
	if req.OwnerIds != nil {
		query.Set("ownerIds", strings.Join(req.OwnerIds, ","))
	}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	if req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.Expiry.Value), 10))
	}
	res := new(TournamentRecordsResponse)
	if err := cl.Do(ctx, "GET", "v2/tournament/"+req.TournamentId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *TournamentRecordsRequest) Async(ctx context.Context, cl *Client, f func(*TournamentRecordsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// WriteTournamentRecord creates a request to write a tournament record.
func WriteTournamentRecord(tournamentId string) *WriteTournamentRecordRequest {
	return &WriteTournamentRecordRequest{
		TournamentId: tournamentId,
		Record:       &TournamentRecordWrite{},
	}
}

// WithScore sets the score on the request.
func (req *WriteTournamentRecordRequest) WithScore(score int64) *WriteTournamentRecordRequest {
	req.Record.Score = score
	return req
}

// WithSubscore sets the subscore on the request.
func (req *WriteTournamentRecordRequest) WithSubscore(subscore int64) *WriteTournamentRecordRequest {
	req.Record.Subscore = subscore
	return req
}

// WithMetadata sets the metadata on the request.
func (req *WriteTournamentRecordRequest) WithMetadata(metadata string) *WriteTournamentRecordRequest {
	req.Record.Metadata = metadata
	return req
}

// WithOperator sets the operator on the request.
func (req *WriteTournamentRecordRequest) WithOperator(operator OpType) *WriteTournamentRecordRequest {
	req.Record.Operator = operator
	return req
}

// Do executes the request against the context and client.
func (req *WriteTournamentRecordRequest) Do(ctx context.Context, cl *Client) (*WriteTournamentRecordResponse, error) {
	res := new(WriteTournamentRecordResponse)
	if err := cl.Do(ctx, "POST", "v2/tournament/"+req.TournamentId, true, nil, req.Record, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *WriteTournamentRecordRequest) Async(ctx context.Context, cl *Client, f func(*WriteTournamentRecordResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// WriteTournamentRecordResponse is the WriteTournamentRecord response.
type WriteTournamentRecordResponse = LeaderboardRecord

// JoinTournament creates a request to join a tournament.
func JoinTournament(tournamentId string) *JoinTournamentRequest {
	return &JoinTournamentRequest{
		TournamentId: tournamentId,
	}
}

// Do executes the request against the context and client.
func (req *JoinTournamentRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/tournament/"+req.TournamentId+"/join", true, nil, nil, nil)
}

// Async executes the request against the context and client.
func (req *JoinTournamentRequest) Async(ctx context.Context, cl *Client, f func(error)) {
	go func() {
		if err := req.Do(ctx, cl); f != nil {
			f(err)
		}
	}()
}

// TournamentRecordsAroundOwner creates a request to retrieve tournament
// records around owner.
func TournamentRecordsAroundOwner(tournamentId, ownerId string) *TournamentRecordsAroundOwnerRequest {
	return &TournamentRecordsAroundOwnerRequest{
		TournamentId: tournamentId,
		OwnerId:      ownerId,
		Limit:        wrapperspb.UInt32(100),
	}
}

// WithLimit sets the limit on the request.
func (req *TournamentRecordsAroundOwnerRequest) WithLimit(limit int) *TournamentRecordsAroundOwnerRequest {
	req.Limit = wrapperspb.UInt32(uint32(limit))
	return req
}

// WithExpiry sets the expiry on the request.
func (req *TournamentRecordsAroundOwnerRequest) WithExpiry(expiry int) *TournamentRecordsAroundOwnerRequest {
	req.Expiry = wrapperspb.Int64(int64(expiry))
	return req
}

// Do executes the request against the context and client.
func (req *TournamentRecordsAroundOwnerRequest) Do(ctx context.Context, cl *Client) (*TournamentRecordsResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.Expiry.Value), 10))
	}
	res := new(TournamentRecordsResponse)
	if err := cl.Do(ctx, "GET", "v2/tournament/"+req.TournamentId+"/owner/"+req.OwnerId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *TournamentRecordsAroundOwnerRequest) Async(ctx context.Context, cl *Client, f func(*TournamentRecordsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// Users creates a request to retrieve users.
func Users(ids ...string) *UsersRequest {
	return &UsersRequest{
		Ids: ids,
	}
}

// WithUsernames sets the usernames on the request.
func (req *UsersRequest) WithUsernames(usernames ...string) *UsersRequest {
	req.Usernames = usernames
	return req
}

// WithFacebookIds sets the facebookIds on the request.
func (req *UsersRequest) WithFacebookIds(facebookIds ...string) *UsersRequest {
	req.FacebookIds = facebookIds
	return req
}

// Do executes the request against the context and client.
func (req *UsersRequest) Do(ctx context.Context, cl *Client) (*UsersResponse, error) {
	query := url.Values{}
	if len(req.Ids) != 0 {
		query.Set("ids", strings.Join(req.Ids, ","))
	}
	if len(req.Usernames) != 0 {
		query.Set("usernames", strings.Join(req.Usernames, ","))
	}
	if len(req.FacebookIds) != 0 {
		query.Set("facebookIds", strings.Join(req.FacebookIds, ","))
	}
	res := new(UsersResponse)
	if err := cl.Do(ctx, "GET", "v2/user", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *UsersRequest) Async(ctx context.Context, cl *Client, f func(*UsersResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}

// UserGroups creates a request to retrieve a user's groups.
func UserGroups(userId string) *UserGroupsRequest {
	return &UserGroupsRequest{
		UserId: userId,
	}
}

// WithLimit sets the limit on the request.
func (req *UserGroupsRequest) WithLimit(limit int) *UserGroupsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithState sets the state on the request.
func (req *UserGroupsRequest) WithState(state UserRoleState) *UserGroupsRequest {
	req.State = wrapperspb.Int32(int32(state))
	return req
}

// WithCursor sets the cursor on the request.
func (req *UserGroupsRequest) WithCursor(cursor string) *UserGroupsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *UserGroupsRequest) Do(ctx context.Context, cl *Client) (*UserGroupsResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.State != nil {
		query.Set("state", strconv.FormatInt(int64(req.State.Value), 10))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	res := new(UserGroupsResponse)
	if err := cl.Do(ctx, "GET", "v2/user/"+req.UserId+"/group", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// Async executes the request against the context and client.
func (req *UserGroupsRequest) Async(ctx context.Context, cl *Client, f func(*UserGroupsResponse, error)) {
	go func() {
		if res, err := req.Do(ctx, cl); f != nil {
			f(res, err)
		}
	}()
}
