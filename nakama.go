// Package nakama is a nakama http and realtime websocket client.
package nakama

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	nkapi "github.com/heroiclabs/nakama-common/api"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// HealthcheckRequest is a healthcheck request.
type HealthcheckRequest struct{}

// Healthcheck creates a new healthcheck request.
func Healthcheck() *HealthcheckRequest {
	return &HealthcheckRequest{}
}

// Do executes the healthcheck request against the context and client.
func (req *HealthcheckRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "GET", "healthcheck", false, nil, nil, nil)
}

// AccountRequest is a account request.
type AccountRequest struct{}

// Account creates a new account request.
func Account() *AccountRequest {
	return &AccountRequest{}
}

// Do executes the request against the context and client.
func (req *AccountRequest) Do(ctx context.Context, cl *Client) (*AccountResponse, error) {
	res := new(nkapi.Account)
	if err := cl.Do(ctx, "GET", "v2/account", true, nil, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AccountResponse is the account repsonse.
type AccountResponse = nkapi.Account

// UpdateAccountRequest is a UpdateAccount request.
type UpdateAccountRequest struct {
	nkapi.UpdateAccountRequest
}

// UpdateAccount creates a new UpdateAccount request.
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

// SessionResponse is the authenticate repsonse.
type SessionResponse = nkapi.Session

// AuthenticateAppleRequest is a AuthenticateApple request.
type AuthenticateAppleRequest struct {
	nkapi.AuthenticateAppleRequest
}

// AuthenticateApple creates a new AuthenticateApple request.
func AuthenticateApple() *AuthenticateAppleRequest {
	return &AuthenticateAppleRequest{
		AuthenticateAppleRequest: nkapi.AuthenticateAppleRequest{
			Account: &nkapi.AccountApple{},
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

// WithToken sets the token on the request.
func (req *AuthenticateAppleRequest) WithToken(token string) *AuthenticateAppleRequest {
	req.Account.Token = token
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

// AuthenticateCustomRequest is a AuthenticateCustom request.
type AuthenticateCustomRequest struct {
	nkapi.AuthenticateCustomRequest
}

// AuthenticateCustom creates a new AuthenticateCustom request.
func AuthenticateCustom() *AuthenticateCustomRequest {
	return &AuthenticateCustomRequest{
		AuthenticateCustomRequest: nkapi.AuthenticateCustomRequest{
			Account: &nkapi.AccountCustom{},
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

// WithId sets the id on the request.
func (req *AuthenticateCustomRequest) WithId(id string) *AuthenticateCustomRequest {
	req.Account.Id = id
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

// AuthenticateDeviceRequest is a AuthenticateDevice request.
type AuthenticateDeviceRequest struct {
	nkapi.AuthenticateDeviceRequest
}

// AuthenticateDevice creates a new AuthenticateDevice request.
func AuthenticateDevice() *AuthenticateDeviceRequest {
	return &AuthenticateDeviceRequest{
		AuthenticateDeviceRequest: nkapi.AuthenticateDeviceRequest{
			Account: &nkapi.AccountDevice{},
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

// WithId sets the id on the request.
func (req *AuthenticateDeviceRequest) WithId(id string) *AuthenticateDeviceRequest {
	req.Account.Id = id
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

// AuthenticateEmailRequest is a AuthenticateEmail request.
type AuthenticateEmailRequest struct {
	nkapi.AuthenticateEmailRequest
}

// AuthenticateEmail creates a new AuthenticateEmail request.
func AuthenticateEmail() *AuthenticateEmailRequest {
	return &AuthenticateEmailRequest{
		AuthenticateEmailRequest: nkapi.AuthenticateEmailRequest{
			Account: &nkapi.AccountEmail{},
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

// WithEmail sets the email on the request.
func (req *AuthenticateEmailRequest) WithEmail(email string) *AuthenticateEmailRequest {
	req.Account.Email = email
	return req
}

// WithPassword sets the password on the request.
func (req *AuthenticateEmailRequest) WithPassword(password string) *AuthenticateEmailRequest {
	req.Account.Password = password
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

// AuthenticateFacebookRequest is a AuthenticateFacebook request.
type AuthenticateFacebookRequest struct {
	nkapi.AuthenticateFacebookRequest
}

// AuthenticateFacebook creates a new AuthenticateFacebook request.
func AuthenticateFacebook() *AuthenticateFacebookRequest {
	return &AuthenticateFacebookRequest{
		AuthenticateFacebookRequest: nkapi.AuthenticateFacebookRequest{
			Account: &nkapi.AccountFacebook{},
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

// WithToken sets the token on the request.
func (req *AuthenticateFacebookRequest) WithToken(token string) *AuthenticateFacebookRequest {
	req.Account.Token = token
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

// AuthenticateFacebookInstantGameRequest is a AuthenticateFacebookInstantGame request.
type AuthenticateFacebookInstantGameRequest struct {
	nkapi.AuthenticateFacebookInstantGameRequest
}

// AuthenticateFacebookInstantGame creates a new AuthenticateFacebookInstantGame request.
func AuthenticateFacebookInstantGame() *AuthenticateFacebookInstantGameRequest {
	return &AuthenticateFacebookInstantGameRequest{
		AuthenticateFacebookInstantGameRequest: nkapi.AuthenticateFacebookInstantGameRequest{
			Account: &nkapi.AccountFacebookInstantGame{},
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

// WithSignedPlayerInfo sets the signedPlayerInfo on the request.
func (req *AuthenticateFacebookInstantGameRequest) WithSignedPlayerInfo(signedPlayerInfo string) *AuthenticateFacebookInstantGameRequest {
	req.Account.SignedPlayerInfo = signedPlayerInfo
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

// AuthenticateGameCenterRequest is a AuthenticateGameCenter request.
type AuthenticateGameCenterRequest struct {
	nkapi.AuthenticateGameCenterRequest
}

// AuthenticateGameCenter creates a new AuthenticateGameCenter request.
func AuthenticateGameCenter() *AuthenticateGameCenterRequest {
	return &AuthenticateGameCenterRequest{
		AuthenticateGameCenterRequest: nkapi.AuthenticateGameCenterRequest{
			Account: &nkapi.AccountGameCenter{},
		},
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

// AuthenticateGoogleRequest is a AuthenticateGoogle request.
type AuthenticateGoogleRequest struct {
	nkapi.AuthenticateGoogleRequest
}

// AuthenticateGoogle creates a new AuthenticateGoogle request.
func AuthenticateGoogle() *AuthenticateGoogleRequest {
	return &AuthenticateGoogleRequest{
		AuthenticateGoogleRequest: nkapi.AuthenticateGoogleRequest{
			Account: &nkapi.AccountGoogle{},
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

// WithToken sets the token on the request.
func (req *AuthenticateGoogleRequest) WithToken(token string) *AuthenticateGoogleRequest {
	req.Account.Token = token
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

// AuthenticateSteamRequest is a AuthenticateSteam request.
type AuthenticateSteamRequest struct {
	nkapi.AuthenticateSteamRequest
}

// AuthenticateSteam creates a new AuthenticateSteam request.
func AuthenticateSteam() *AuthenticateSteamRequest {
	return &AuthenticateSteamRequest{
		AuthenticateSteamRequest: nkapi.AuthenticateSteamRequest{
			Account: &nkapi.AccountSteam{},
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

// WithToken sets the token on the request.
func (req *AuthenticateSteamRequest) WithToken(token string) *AuthenticateSteamRequest {
	req.Account.Token = token
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

// LinkAppleRequest is a LinkApple request.
type LinkAppleRequest struct {
	nkapi.AccountApple
}

// LinkApple creates a new LinkApple request.
func LinkApple() *LinkAppleRequest {
	return &LinkAppleRequest{}
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

// LinkCustomRequest is a LinkCustom request.
type LinkCustomRequest struct {
	nkapi.AccountCustom
}

// LinkCustom creates a new LinkCustom request.
func LinkCustom() *LinkCustomRequest {
	return &LinkCustomRequest{}
}

// WithId sets the id on the request.
func (req *LinkCustomRequest) WithId(id string) *LinkCustomRequest {
	req.Id = id
	return req
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

// LinkDeviceRequest is a LinkDevice request.
type LinkDeviceRequest struct {
	nkapi.AccountDevice
}

// LinkDevice creates a new LinkDevice request.
func LinkDevice() *LinkDeviceRequest {
	return &LinkDeviceRequest{}
}

// WithId sets the id on the request.
func (req *LinkDeviceRequest) WithId(id string) *LinkDeviceRequest {
	req.Id = id
	return req
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

// LinkEmailRequest is a LinkEmail request.
type LinkEmailRequest struct {
	nkapi.AccountEmail
}

// LinkEmail creates a new LinkEmail request.
func LinkEmail() *LinkEmailRequest {
	return &LinkEmailRequest{}
}

// WithEmail sets the email on the request.
func (req *LinkEmailRequest) WithEmail(email string) *LinkEmailRequest {
	req.Email = email
	return req
}

// WithPassword sets the password on the request.
func (req *LinkEmailRequest) WithPassword(password string) *LinkEmailRequest {
	req.Password = password
	return req
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

// LinkFacebookRequest is a LinkFacebook request.
type LinkFacebookRequest struct {
	nkapi.LinkFacebookRequest
}

// LinkFacebook creates a new LinkFacebook request.
func LinkFacebook() *LinkFacebookRequest {
	return &LinkFacebookRequest{
		LinkFacebookRequest: nkapi.LinkFacebookRequest{
			Account: &nkapi.AccountFacebook{},
		},
	}
}

// WithSync sets the sync on the request.
func (req *LinkFacebookRequest) WithSync(sync bool) *LinkFacebookRequest {
	req.Sync = wrapperspb.Bool(sync)
	return req
}

// WithToken sets the token on the request.
func (req *LinkFacebookRequest) WithToken(token string) *LinkFacebookRequest {
	req.Account.Token = token
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

// LinkFacebookInstantGameRequest is a LinkFacebookInstantGame request.
type LinkFacebookInstantGameRequest struct {
	nkapi.AccountFacebookInstantGame
}

// LinkFacebookInstantGame creates a new LinkFacebookInstantGame request.
func LinkFacebookInstantGame() *LinkFacebookInstantGameRequest {
	return &LinkFacebookInstantGameRequest{}
}

// WithSignedPlayerInfo sets the signedPlayerInfo on the request.
func (req *LinkFacebookInstantGameRequest) WithSignedPlayerInfo(signedPlayerInfo string) *LinkFacebookInstantGameRequest {
	req.SignedPlayerInfo = signedPlayerInfo
	return req
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

// LinkGameCenterRequest is a LinkGameCenter request.
type LinkGameCenterRequest struct {
	nkapi.AccountGameCenter
}

// LinkGameCenter creates a new LinkGameCenter request.
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

// LinkGoogleRequest is a LinkGoogle request.
type LinkGoogleRequest struct {
	nkapi.AccountGoogle
}

// LinkGoogle creates a new LinkGoogle request.
func LinkGoogle() *LinkGoogleRequest {
	return &LinkGoogleRequest{}
}

// WithToken sets the token on the request.
func (req *LinkGoogleRequest) WithToken(token string) *LinkGoogleRequest {
	req.Token = token
	return req
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

// LinkSteamRequest is a LinkSteam request.
type LinkSteamRequest struct {
	nkapi.LinkSteamRequest
}

// LinkSteam creates a new LinkSteam request.
func LinkSteam() *LinkSteamRequest {
	return &LinkSteamRequest{
		LinkSteamRequest: nkapi.LinkSteamRequest{
			Account: &nkapi.AccountSteam{},
		},
	}
}

// WithSync sets the sync on the request.
func (req *LinkSteamRequest) WithSync(sync bool) *LinkSteamRequest {
	req.Sync = wrapperspb.Bool(sync)
	return req
}

// WithToken sets the token on the request.
func (req *LinkSteamRequest) WithToken(token string) *LinkSteamRequest {
	req.Account.Token = token
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

// SessionRefreshRequest is a SessionRefresh request.
type SessionRefreshRequest struct {
	nkapi.SessionRefreshRequest
}

// SessionRefresh creates a new SessionRefresh request.
func SessionRefresh(refreshToken string) *SessionRefreshRequest {
	return &SessionRefreshRequest{
		SessionRefreshRequest: nkapi.SessionRefreshRequest{
			Token: refreshToken,
		},
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

// UnlinkAppleRequest is a UnlinkApple request.
type UnlinkAppleRequest struct {
	nkapi.AccountApple
}

// UnlinkApple creates a new UnlinkApple request.
func UnlinkApple() *UnlinkAppleRequest {
	return &UnlinkAppleRequest{}
}

// WithToken sets the token on the request.
func (req *UnlinkAppleRequest) WithToken(token string) *UnlinkAppleRequest {
	req.Token = token
	return req
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

// UnlinkCustomRequest is a UnlinkCustom request.
type UnlinkCustomRequest struct {
	nkapi.AccountCustom
}

// UnlinkCustom creates a new UnlinkCustom request.
func UnlinkCustom() *UnlinkCustomRequest {
	return &UnlinkCustomRequest{}
}

// WithId sets the id on the request.
func (req *UnlinkCustomRequest) WithId(id string) *UnlinkCustomRequest {
	req.Id = id
	return req
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

// UnlinkDeviceRequest is a UnlinkDevice request.
type UnlinkDeviceRequest struct {
	nkapi.AccountDevice
}

// UnlinkDevice creates a new UnlinkDevice request.
func UnlinkDevice() *UnlinkDeviceRequest {
	return &UnlinkDeviceRequest{}
}

// WithId sets the id on the request.
func (req *UnlinkDeviceRequest) WithId(id string) *UnlinkDeviceRequest {
	req.Id = id
	return req
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

// UnlinkEmailRequest is a UnlinkEmail request.
type UnlinkEmailRequest struct {
	nkapi.AccountEmail
}

// UnlinkEmail creates a new UnlinkEmail request.
func UnlinkEmail() *UnlinkEmailRequest {
	return &UnlinkEmailRequest{}
}

// WithEmail sets the email on the request.
func (req *UnlinkEmailRequest) WithEmail(email string) *UnlinkEmailRequest {
	req.Email = email
	return req
}

// WithPassword sets the password on the request.
func (req *UnlinkEmailRequest) WithPassword(password string) *UnlinkEmailRequest {
	req.Password = password
	return req
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

// UnlinkFacebookRequest is a UnlinkFacebook request.
type UnlinkFacebookRequest struct {
	nkapi.AccountFacebook
}

// UnlinkFacebook creates a new UnlinkFacebook request.
func UnlinkFacebook() *UnlinkFacebookRequest {
	return &UnlinkFacebookRequest{}
}

// WithToken sets the token on the request.
func (req *UnlinkFacebookRequest) WithToken(token string) *UnlinkFacebookRequest {
	req.Token = token
	return req
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

// UnlinkFacebookInstantGameRequest is a UnlinkFacebookInstantGame request.
type UnlinkFacebookInstantGameRequest struct {
	nkapi.AccountFacebookInstantGame
}

// UnlinkFacebookInstantGame creates a new UnlinkFacebookInstantGame request.
func UnlinkFacebookInstantGame() *UnlinkFacebookInstantGameRequest {
	return &UnlinkFacebookInstantGameRequest{}
}

// WithSignedPlayerInfo sets the signedPlayerInfo on the request.
func (req *UnlinkFacebookInstantGameRequest) WithSignedPlayerInfo(signedPlayerInfo string) *UnlinkFacebookInstantGameRequest {
	req.SignedPlayerInfo = signedPlayerInfo
	return req
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

// UnlinkGameCenterRequest is a UnlinkGameCenter request.
type UnlinkGameCenterRequest struct {
	nkapi.AccountGameCenter
}

// UnlinkGameCenter creates a new UnlinkGameCenter request.
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

// UnlinkGoogleRequest is a UnlinkGoogle request.
type UnlinkGoogleRequest struct {
	nkapi.AccountGoogle
}

// UnlinkGoogle creates a new UnlinkGoogle request.
func UnlinkGoogle() *UnlinkGoogleRequest {
	return &UnlinkGoogleRequest{}
}

// WithToken sets the token on the request.
func (req *UnlinkGoogleRequest) WithToken(token string) *UnlinkGoogleRequest {
	req.Token = token
	return req
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

// UnlinkSteamRequest is a UnlinkSteam request.
type UnlinkSteamRequest struct {
	nkapi.AccountSteam
}

// UnlinkSteam creates a new UnlinkSteam request.
func UnlinkSteam() *UnlinkSteamRequest {
	return &UnlinkSteamRequest{}
}

// WithToken sets the token on the request.
func (req *UnlinkSteamRequest) WithToken(token string) *UnlinkSteamRequest {
	req.Token = token
	return req
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

// ListChannelMessagesRequest is a ListChannelMessages request.
type ListChannelMessagesRequest struct {
	nkapi.ListChannelMessagesRequest
}

// ListChannelMessages creates a new ListChannelMessages request.
func ListChannelMessages(channelId string) *ListChannelMessagesRequest {
	return &ListChannelMessagesRequest{
		ListChannelMessagesRequest: nkapi.ListChannelMessagesRequest{
			ChannelId: channelId,
			Limit:     wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListChannelMessagesRequest) WithLimit(limit int) *ListChannelMessagesRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithForward sets the forward on the request.
func (req *ListChannelMessagesRequest) WithForward(forward bool) *ListChannelMessagesRequest {
	req.Forward = wrapperspb.Bool(forward)
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListChannelMessagesRequest) WithCursor(cursor string) *ListChannelMessagesRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListChannelMessagesRequest) Do(ctx context.Context, cl *Client) (*ListChannelMessagesResponse, error) {
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
	res := new(ListChannelMessagesResponse)
	if err := cl.Do(ctx, "GET", "v2/channel/"+req.ChannelId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListChannelMessagesResponse is the ListChannelMessages response.
type ListChannelMessagesResponse = nkapi.ChannelMessageList

// EventRequest is a Event request.
type EventRequest struct {
	nkapi.Event
}

// Event creates a new Event request.
func Event() *EventRequest {
	return &EventRequest{}
}

// WithName sets the name on the request.
func (req *EventRequest) WithName(name string) *EventRequest {
	req.Name = name
	return req
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

// ListFriendsRequest is a ListFriends request.
type ListFriendsRequest struct {
	nkapi.ListFriendsRequest
}

// ListFriends creates a new ListFriends request.
func ListFriends() *ListFriendsRequest {
	return &ListFriendsRequest{
		ListFriendsRequest: nkapi.ListFriendsRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListFriendsRequest) WithLimit(limit int) *ListFriendsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithState sets the state on the request.
func (req *ListFriendsRequest) WithState(state int) *ListFriendsRequest {
	req.State = wrapperspb.Int32(int32(state))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListFriendsRequest) WithCursor(cursor string) *ListFriendsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListFriendsRequest) Do(ctx context.Context, cl *Client) (*ListFriendsResponse, error) {
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
	res := new(ListFriendsResponse)
	if err := cl.Do(ctx, "GET", "v2/friend", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListFriendsResponse is the ListFriends response.
type ListFriendsResponse = nkapi.FriendList

// DeleteFriendsRequest is a DeleteFriends request.
type DeleteFriendsRequest struct {
	nkapi.DeleteFriendsRequest
}

// DeleteFriends creates a new DeleteFriends request.
func DeleteFriends() *DeleteFriendsRequest {
	return &DeleteFriendsRequest{}
}

// WithIds sets the Ids on the request.
func (req *DeleteFriendsRequest) WithIds(ids ...string) *DeleteFriendsRequest {
	req.Ids = ids
	return req
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

// AddFriendsRequest is a AddFriends request.
type AddFriendsRequest struct {
	nkapi.AddFriendsRequest
}

// AddFriends creates a new AddFriends request.
func AddFriends() *AddFriendsRequest {
	return &AddFriendsRequest{}
}

// WithIds sets the Ids on the request.
func (req *AddFriendsRequest) WithIds(ids ...string) *AddFriendsRequest {
	req.Ids = ids
	return req
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

// BlockFriendsRequest is a BlockFriends request.
type BlockFriendsRequest struct {
	nkapi.BlockFriendsRequest
}

// BlockFriends creates a new BlockFriends request.
func BlockFriends() *BlockFriendsRequest {
	return &BlockFriendsRequest{}
}

// WithIds sets the Ids on the request.
func (req *BlockFriendsRequest) WithIds(ids ...string) *BlockFriendsRequest {
	req.Ids = ids
	return req
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

// ImportFacebookFriendsRequest is a ImportFacebookFriends request.
type ImportFacebookFriendsRequest struct {
	nkapi.ImportFacebookFriendsRequest
}

// ImportFacebookFriends creates a new ImportFacebookFriends request.
func ImportFacebookFriends() *ImportFacebookFriendsRequest {
	return &ImportFacebookFriendsRequest{
		ImportFacebookFriendsRequest: nkapi.ImportFacebookFriendsRequest{
			Account: &nkapi.AccountFacebook{},
		},
	}
}

// WithReset sets the reset on the request.
func (req *ImportFacebookFriendsRequest) WithReset(reset bool) *ImportFacebookFriendsRequest {
	req.Reset_ = wrapperspb.Bool(reset)
	return req
}

// WithToken sets the token on the request.
func (req *ImportFacebookFriendsRequest) WithToken(token string) *ImportFacebookFriendsRequest {
	req.Account.Token = token
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

// ImportSteamFriendsRequest is a ImportSteamFriends request.
type ImportSteamFriendsRequest struct {
	nkapi.ImportSteamFriendsRequest
}

// ImportSteamFriends creates a new ImportSteamFriends request.
func ImportSteamFriends() *ImportSteamFriendsRequest {
	return &ImportSteamFriendsRequest{
		ImportSteamFriendsRequest: nkapi.ImportSteamFriendsRequest{
			Account: &nkapi.AccountSteam{},
		},
	}
}

// WithReset sets the reset on the request.
func (req *ImportSteamFriendsRequest) WithReset(reset bool) *ImportSteamFriendsRequest {
	req.Reset_ = wrapperspb.Bool(reset)
	return req
}

// WithToken sets the token on the request.
func (req *ImportSteamFriendsRequest) WithToken(token string) *ImportSteamFriendsRequest {
	req.Account.Token = token
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

// ListGroupsRequest is a ListGroups request.
type ListGroupsRequest struct {
	nkapi.ListGroupsRequest
}

// ListGroups creates a new ListGroups request.
func ListGroups() *ListGroupsRequest {
	return &ListGroupsRequest{}
}

// WithName sets the name on the request.
func (req *ListGroupsRequest) WithName(name string) *ListGroupsRequest {
	req.Name = name
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListGroupsRequest) WithCursor(cursor string) *ListGroupsRequest {
	req.Cursor = cursor
	return req
}

// WithLimit sets the limit on the request.
func (req *ListGroupsRequest) WithLimit(limit int) *ListGroupsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithLangTag sets the langTag on the request.
func (req *ListGroupsRequest) WithLangTag(langTag string) *ListGroupsRequest {
	req.LangTag = langTag
	return req
}

// WithMembers sets the members on the request.
func (req *ListGroupsRequest) WithMembers(members int) *ListGroupsRequest {
	req.Members = wrapperspb.Int32(int32(members))
	return req
}

// WithOpen sets the open on the request.
func (req *ListGroupsRequest) WithOpen(open bool) *ListGroupsRequest {
	req.Open = wrapperspb.Bool(open)
	return req
}

// Do executes the request against the context and client.
func (req *ListGroupsRequest) Do(ctx context.Context, cl *Client) (*ListGroupsResponse, error) {
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
	res := new(ListGroupsResponse)
	if err := cl.Do(ctx, "GET", "v2/group", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListGroupsResponse is the ListGroups response.
type ListGroupsResponse = nkapi.GroupList

// CreateGroupRequest is a CreateGroup request.
type CreateGroupRequest struct {
	nkapi.CreateGroupRequest
}

// CreateGroup creates a new CreateGroup request.
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
func (req *CreateGroupRequest) Do(ctx context.Context, cl *Client) (*nkapi.Group, error) {
	res := new(nkapi.Group)
	if err := cl.Do(ctx, "POST", "v2/group", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteGroupRequest is a DeleteGroup request.
type DeleteGroupRequest struct {
	nkapi.DeleteGroupRequest
}

// DeleteGroup creates a new DeleteGroup request.
func DeleteGroup(groupId string) *DeleteGroupRequest {
	return &DeleteGroupRequest{
		DeleteGroupRequest: nkapi.DeleteGroupRequest{
			GroupId: groupId,
		},
	}
}

// Do executes the request against the context and client.
func (req *DeleteGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/group/"+req.GroupId, true, nil, nil, nil)
}

// UpdateGroupRequest is a UpdateGroup request.
type UpdateGroupRequest struct {
	nkapi.UpdateGroupRequest
}

// UpdateGroup creates a new UpdateGroup request.
func UpdateGroup(groupId string) *UpdateGroupRequest {
	return &UpdateGroupRequest{
		UpdateGroupRequest: nkapi.UpdateGroupRequest{
			GroupId: groupId,
		},
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

// AddGroupUsersRequest is a AddGroupUsers request.
type AddGroupUsersRequest struct {
	nkapi.AddGroupUsersRequest
}

// AddGroupUsers creates a new AddGroupUsers request.
func AddGroupUsers(groupId string) *AddGroupUsersRequest {
	return &AddGroupUsersRequest{
		AddGroupUsersRequest: nkapi.AddGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *AddGroupUsersRequest) WithUserIds(userIds ...string) *AddGroupUsersRequest {
	req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *AddGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/add", true, nil, req, nil)
}

// BanGroupUsersRequest is a BanGroupUsers request.
type BanGroupUsersRequest struct {
	nkapi.BanGroupUsersRequest
}

// BanGroupUsers creates a new BanGroupUsers request.
func BanGroupUsers(groupId string) *BanGroupUsersRequest {
	return &BanGroupUsersRequest{
		BanGroupUsersRequest: nkapi.BanGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *BanGroupUsersRequest) WithUserIds(userIds ...string) *BanGroupUsersRequest {
	req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *BanGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/ban", true, nil, req, nil)
}

// DemoteGroupUsersRequest is a DemoteGroupUsers request.
type DemoteGroupUsersRequest struct {
	nkapi.DemoteGroupUsersRequest
}

// DemoteGroupUsers creates a new DemoteGroupUsers request.
func DemoteGroupUsers(groupId string) *DemoteGroupUsersRequest {
	return &DemoteGroupUsersRequest{
		DemoteGroupUsersRequest: nkapi.DemoteGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *DemoteGroupUsersRequest) WithUserIds(userIds ...string) *DemoteGroupUsersRequest {
	req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *DemoteGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/demote", true, nil, req, nil)
}

// JoinGroupRequest is a JoinGroup request.
type JoinGroupRequest struct {
	nkapi.JoinGroupRequest
}

// JoinGroup creates a new JoinGroup request.
func JoinGroup(groupId string) *JoinGroupRequest {
	return &JoinGroupRequest{
		JoinGroupRequest: nkapi.JoinGroupRequest{
			GroupId: groupId,
		},
	}
}

// Do executes the request against the context and client.
func (req *JoinGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/join", true, nil, nil, nil)
}

// KickGroupUsersRequest is a KickGroupUsers request.
type KickGroupUsersRequest struct {
	nkapi.KickGroupUsersRequest
}

// KickGroupUsers creates a new KickGroupUsers request.
func KickGroupUsers(groupId string) *KickGroupUsersRequest {
	return &KickGroupUsersRequest{
		KickGroupUsersRequest: nkapi.KickGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *KickGroupUsersRequest) WithUserIds(userIds ...string) *KickGroupUsersRequest {
	req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *KickGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/kick", true, nil, req, nil)
}

// LeaveGroupRequest is a LeaveGroup request.
type LeaveGroupRequest struct {
	nkapi.LeaveGroupRequest
}

// LeaveGroup creates a new LeaveGroup request.
func LeaveGroup(groupId string) *LeaveGroupRequest {
	return &LeaveGroupRequest{
		LeaveGroupRequest: nkapi.LeaveGroupRequest{
			GroupId: groupId,
		},
	}
}

// Do executes the request against the context and client.
func (req *LeaveGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/leave", true, nil, nil, nil)
}

// PromoteGroupUsersRequest is a PromoteGroupUsers request.
type PromoteGroupUsersRequest struct {
	nkapi.PromoteGroupUsersRequest
}

// PromoteGroupUsers creates a new PromoteGroupUsers request.
func PromoteGroupUsers(groupId string) *PromoteGroupUsersRequest {
	return &PromoteGroupUsersRequest{
		PromoteGroupUsersRequest: nkapi.PromoteGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *PromoteGroupUsersRequest) WithUserIds(userIds ...string) *PromoteGroupUsersRequest {
	req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *PromoteGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.GroupId+"/promote", true, nil, req, nil)
}

// ListGroupUsersRequest is a ListGroupUsers request.
type ListGroupUsersRequest struct {
	nkapi.ListGroupUsersRequest
}

// ListGroupUsers creates a new ListGroupUsers request.
func ListGroupUsers(groupId string) *ListGroupUsersRequest {
	return &ListGroupUsersRequest{
		ListGroupUsersRequest: nkapi.ListGroupUsersRequest{
			GroupId: groupId,
			Limit:   wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListGroupUsersRequest) WithLimit(limit int) *ListGroupUsersRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithState sets the state on the request.
func (req *ListGroupUsersRequest) WithState(state int) *ListGroupUsersRequest {
	req.State = wrapperspb.Int32(int32(state))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListGroupUsersRequest) WithCursor(cursor string) *ListGroupUsersRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListGroupUsersRequest) Do(ctx context.Context, cl *Client) (*ListGroupUsersResponse, error) {
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
	res := new(ListGroupUsersResponse)
	if err := cl.Do(ctx, "GET", "v2/group/"+req.GroupId+"/user", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListGroupUsersResponse is the ListGroupUsers response.
type ListGroupUsersResponse = nkapi.GroupUserList

// ValidatePurchaseResponse is the ValidatePurchaseApple response.
type ValidatePurchaseResponse = nkapi.ValidatePurchaseResponse

// ValidatePurchaseAppleRequest is a ValidatePurchaseApple request.
type ValidatePurchaseAppleRequest struct {
	nkapi.ValidatePurchaseAppleRequest
}

// ValidatePurchaseApple creates a new ValidatePurchaseApple request.
func ValidatePurchaseApple() *ValidatePurchaseAppleRequest {
	return &ValidatePurchaseAppleRequest{}
}

// WithReceipt sets the receipt on the request.
func (req *ValidatePurchaseAppleRequest) WithReceipt(receipt string) *ValidatePurchaseAppleRequest {
	req.Receipt = receipt
	return req
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

// ValidatePurchaseGoogleRequest is a ValidatePurchaseGoogle request.
type ValidatePurchaseGoogleRequest struct {
	nkapi.ValidatePurchaseGoogleRequest
}

// ValidatePurchaseGoogle creates a new ValidatePurchaseGoogle request.
func ValidatePurchaseGoogle() *ValidatePurchaseGoogleRequest {
	return &ValidatePurchaseGoogleRequest{}
}

// WithPurchase sets the purchase on the request.
func (req *ValidatePurchaseGoogleRequest) WithPurchase(purchase string) *ValidatePurchaseGoogleRequest {
	req.Purchase = purchase
	return req
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

// ValidatePurchaseHuaweiRequest is a ValidatePurchaseHuawei request.
type ValidatePurchaseHuaweiRequest struct {
	nkapi.ValidatePurchaseHuaweiRequest
}

// ValidatePurchaseHuawei creates a new ValidatePurchaseHuawei request.
func ValidatePurchaseHuawei() *ValidatePurchaseHuaweiRequest {
	return &ValidatePurchaseHuaweiRequest{}
}

// WithPurchase sets the purchase on the request.
func (req *ValidatePurchaseHuaweiRequest) WithPurchase(purchase string) *ValidatePurchaseHuaweiRequest {
	req.Purchase = purchase
	return req
}

// WithSignature sets the signature on the request.
func (req *ValidatePurchaseHuaweiRequest) WithSignature(signature string) *ValidatePurchaseHuaweiRequest {
	req.Signature = signature
	return req
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

/*
// ListSubscriptionsRequest is a ListSubscriptions request.
type ListSubscriptionsRequest struct {
	nkapi.ListSubscriptionsRequest
}

// ListSubscriptions creates a new ListSubscriptions request.
func ListSubscriptions(groupId string) *ListSubscriptionsRequest {
	return &ListSubscriptionsRequest{
		Limit: wrapperspb.Int32(100),
	}
}

// WithLimit sets the limit on the request.
func (req *ListSubscriptionsRequest) WithLimit(limit int) *ListSubscriptionsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListSubscriptionsRequest) WithCursor(cursor string) *ListSubscriptionsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListSubscriptionsRequest) Do(ctx context.Context, cl *Client) (*ListSubscriptionsResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Cursor != "" {
		query.Set("cursor", req.Cursor)
	}
	res := new(ListSubscriptionsResponse)
	if err := cl.Do(ctx, "GET", "v2/iap/subscription", true, nil, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListSubscriptionsResponse is the ListSubscriptions response.
type ListSubscriptionsResponse = nkapi.SubscriptionList

// ValidateSubscriptionResponse is the ValidateSubscriptionApple response.
type ValidateSubscriptionResponse = nkapi.ValidateSubscriptionResponse

// ValidateSubscriptionAppleRequest is a ValidateSubscriptionApple request.
type ValidateSubscriptionAppleRequest struct {
	nkapi.ValidateSubscriptionAppleRequest
}

// ValidateSubscriptionApple creates a new ValidateSubscriptionApple request.
func ValidateSubscriptionApple() *ValidateSubscriptionAppleRequest {
	return &ValidateSubscriptionAppleRequest{
	}
}

// WithReceipt sets the receipt on the request.
func (req *ValidateSubscriptionAppleRequest) WithReceipt(receipt string) *ValidateSubscriptionAppleRequest {
	req.Receipt = receipt
	return req
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

// ValidateSubscriptionGoogleRequest is a ValidateSubscriptionGoogle request.
type ValidateSubscriptionGoogleRequest struct {
	nkapi.ValidateSubscriptionGoogleRequest
}

// ValidateSubscriptionGoogle creates a new ValidateSubscriptionGoogle request.
func ValidateSubscriptionGoogle() *ValidateSubscriptionGoogleRequest {
	return &ValidateSubscriptionGoogleRequest{
	}
}

// WithReceipt sets the receipt on the request.
func (req *ValidateSubscriptionGoogleRequest) WithReceipt(receipt string) *ValidateSubscriptionGoogleRequest {
	req.Receipt = receipt
	return req
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
*/

/*
// SubscriptionRequest is a Subscription request.
type SubscriptionRequest struct {
	nkapi.GetSubscriptionRequest
}

// Subscription creates a new Subscription request.
func Subscription(productId string) *SubscriptionRequest {
	return &SubscriptionRequest{
		ProductId: productId,
	}
}

// Do executes the request against the context and client.
func (req *SubscriptionRequest) Do(ctx context.Context, cl *Client) (*SubscriptionResponse, error) {
	res := new(SubscriptionResponse)
	if err := cl.Do(ctx, "GET", "v2/iap/subscription/"+req.ProductId, nil, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// SubscriptionResponse is a Subscription response.
type SubscriptionResponse = nkapi.ValidatedSubscription

*/

// ListLeaderboardRecordsRequest is a ListLeaderboardRecords request.
type ListLeaderboardRecordsRequest struct {
	nkapi.ListLeaderboardRecordsRequest
}

// ListLeaderboardRecords creates a new ListLeaderboardRecords request.
func ListLeaderboardRecords(leaderboardId string) *ListLeaderboardRecordsRequest {
	return &ListLeaderboardRecordsRequest{
		ListLeaderboardRecordsRequest: nkapi.ListLeaderboardRecordsRequest{
			LeaderboardId: leaderboardId,
			Limit:         wrapperspb.Int32(100),
		},
	}
}

// WithOwnerIds sets the ownerIds on the request.
func (req *ListLeaderboardRecordsRequest) WithOwnerIds(ownerIds ...string) *ListLeaderboardRecordsRequest {
	req.OwnerIds = ownerIds
	return req
}

// WithLimit sets the limit on the request.
func (req *ListLeaderboardRecordsRequest) WithLimit(limit int) *ListLeaderboardRecordsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListLeaderboardRecordsRequest) WithCursor(cursor string) *ListLeaderboardRecordsRequest {
	req.Cursor = cursor
	return req
}

// WithExpiry sets the expiry on the request.
func (req *ListLeaderboardRecordsRequest) WithExpiry(expiry int) *ListLeaderboardRecordsRequest {
	req.Expiry = wrapperspb.Int64(int64(expiry))
	return req
}

// Do executes the request against the context and client.
func (req *ListLeaderboardRecordsRequest) Do(ctx context.Context, cl *Client) (*ListLeaderboardRecordsResponse, error) {
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
	res := new(ListLeaderboardRecordsResponse)
	if err := cl.Do(ctx, "GET", "v2/leaderboard/"+req.LeaderboardId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListLeaderboardRecordsResponse is the ListLeaderboardRecords response.
type ListLeaderboardRecordsResponse = nkapi.LeaderboardRecordList

// DeleteLeaderboardRecordRequest is a DeleteLeaderboardRecord request.
type DeleteLeaderboardRecordRequest struct {
	nkapi.DeleteLeaderboardRecordRequest
}

// DeleteLeaderboardRecord creates a new DeleteLeaderboardRecord request.
func DeleteLeaderboardRecord(leaderboardId string) *DeleteLeaderboardRecordRequest {
	return &DeleteLeaderboardRecordRequest{
		DeleteLeaderboardRecordRequest: nkapi.DeleteLeaderboardRecordRequest{
			LeaderboardId: leaderboardId,
		},
	}
}

// Do executes the request against the context and client.
func (req *DeleteLeaderboardRecordRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/leaderboard/"+req.LeaderboardId, true, nil, nil, nil)
}

// Operator is the operator.
type Operator = nkapi.Operator

// Operators.
const (
	OpNoOverride Operator = 0
	OpBest       Operator = 1
	OpSet        Operator = 2
	OpIncrement  Operator = 3
	OpDecrement  Operator = 4
)

// WriteLeaderboardRecordRequest is a WriteLeaderboardRecord request.
type WriteLeaderboardRecordRequest struct {
	nkapi.WriteLeaderboardRecordRequest
}

// WriteLeaderboardRecord creates a new WriteLeaderboardRecord request.
func WriteLeaderboardRecord(leaderboardId string) *WriteLeaderboardRecordRequest {
	return &WriteLeaderboardRecordRequest{
		WriteLeaderboardRecordRequest: nkapi.WriteLeaderboardRecordRequest{
			LeaderboardId: leaderboardId,
			Record:        &nkapi.WriteLeaderboardRecordRequest_LeaderboardRecordWrite{},
		},
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
func (req *WriteLeaderboardRecordRequest) WithOperator(operator Operator) *WriteLeaderboardRecordRequest {
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

// WriteLeaderboardRecordResponse is the WriteLeaderboardRecord response.
type WriteLeaderboardRecordResponse = nkapi.LeaderboardRecord

// ListLeaderboardRecordsAroundOwnerRequest is a ListLeaderboardRecordsAroundOwner request.
type ListLeaderboardRecordsAroundOwnerRequest struct {
	nkapi.ListLeaderboardRecordsAroundOwnerRequest
}

// ListLeaderboardRecordsAroundOwner creates a new ListLeaderboardRecordsAroundOwner request.
func ListLeaderboardRecordsAroundOwner(leaderboardId, ownerId string) *ListLeaderboardRecordsAroundOwnerRequest {
	return &ListLeaderboardRecordsAroundOwnerRequest{
		ListLeaderboardRecordsAroundOwnerRequest: nkapi.ListLeaderboardRecordsAroundOwnerRequest{
			LeaderboardId: leaderboardId,
			OwnerId:       ownerId,
			Limit:         wrapperspb.UInt32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListLeaderboardRecordsAroundOwnerRequest) WithLimit(limit int) *ListLeaderboardRecordsAroundOwnerRequest {
	req.Limit = wrapperspb.UInt32(uint32(limit))
	return req
}

// WithExpiry sets the expiry on the request.
func (req *ListLeaderboardRecordsAroundOwnerRequest) WithExpiry(expiry int) *ListLeaderboardRecordsAroundOwnerRequest {
	req.Expiry = wrapperspb.Int64(int64(expiry))
	return req
}

// Do executes the request against the context and client.
func (req *ListLeaderboardRecordsAroundOwnerRequest) Do(ctx context.Context, cl *Client) (*ListLeaderboardRecordsAroundOwnerResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.Expiry.Value), 10))
	}
	/*
		if req.Cursor != "" {
			query.Set("cursor", req.Cursor)
		}
	*/
	res := new(ListLeaderboardRecordsAroundOwnerResponse)
	if err := cl.Do(ctx, "GET", "v2/leaderboard/"+req.LeaderboardId+"/owner/"+req.OwnerId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListLeaderboardRecordsAroundOwnerResponse is the ListLeaderboardRecordsAroundOwner response.
type ListLeaderboardRecordsAroundOwnerResponse = nkapi.LeaderboardRecordList

// ListMatchesRequest is a ListMatches request.
type ListMatchesRequest struct {
	nkapi.ListMatchesRequest
}

// ListMatches creates a new ListMatches request.
func ListMatches() *ListMatchesRequest {
	return &ListMatchesRequest{
		ListMatchesRequest: nkapi.ListMatchesRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListMatchesRequest) WithLimit(limit int) *ListMatchesRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithAuthoritative sets the authoritative on the request.
func (req *ListMatchesRequest) WithAuthoritative(authoritative bool) *ListMatchesRequest {
	req.Authoritative = wrapperspb.Bool(authoritative)
	return req
}

// WithLabel sets the label on the request.
func (req *ListMatchesRequest) WithLabel(label string) *ListMatchesRequest {
	req.Label = wrapperspb.String(label)
	return req
}

// WithMinSize sets the minSize on the request.
func (req *ListMatchesRequest) WithMinSize(minSize int) *ListMatchesRequest {
	req.MinSize = wrapperspb.Int32(int32(minSize))
	return req
}

// WithMaxSize sets the maxSize on the request.
func (req *ListMatchesRequest) WithMaxSize(maxSize int) *ListMatchesRequest {
	req.MaxSize = wrapperspb.Int32(int32(maxSize))
	return req
}

// WithQuery sets the query on the request.
func (req *ListMatchesRequest) WithQuery(query string) *ListMatchesRequest {
	req.Query = wrapperspb.String(query)
	return req
}

// Do executes the request against the context and client.
func (req *ListMatchesRequest) Do(ctx context.Context, cl *Client) (*ListMatchesResponse, error) {
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
	res := new(ListMatchesResponse)
	if err := cl.Do(ctx, "GET", "v2/match", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListMatchesResponse is the ListMatches response.
type ListMatchesResponse = nkapi.MatchList

// ListNotificationsRequest is a ListNotifications request.
type ListNotificationsRequest struct {
	nkapi.ListNotificationsRequest
}

// ListNotifications creates a new ListNotifications request.
func ListNotifications() *ListNotificationsRequest {
	return &ListNotificationsRequest{
		ListNotificationsRequest: nkapi.ListNotificationsRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListNotificationsRequest) WithLimit(limit int) *ListNotificationsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCacheableCursor sets the cacheableCursor on the request.
func (req *ListNotificationsRequest) WithCacheableCursor(cacheableCursor string) *ListNotificationsRequest {
	req.CacheableCursor = cacheableCursor
	return req
}

// Do executes the request against the context and client.
func (req *ListNotificationsRequest) Do(ctx context.Context, cl *Client) (*ListNotificationsResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.CacheableCursor != "" {
		query.Set("cacheableCursor", req.CacheableCursor)
	}
	res := new(ListNotificationsResponse)
	if err := cl.Do(ctx, "GET", "v2/notifications", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListNotificationsResponse is the ListNotifications response.
type ListNotificationsResponse = nkapi.NotificationList

// DeleteNotificationsRequest is a DeleteNotifications request.
type DeleteNotificationsRequest struct {
	nkapi.DeleteNotificationsRequest
}

// DeleteNotifications creates a new DeleteNotifications request.
func DeleteNotifications() *DeleteNotificationsRequest {
	return &DeleteNotificationsRequest{}
}

// WithIds sets the Ids on the request.
func (req *DeleteNotificationsRequest) WithIds(ids ...string) *DeleteNotificationsRequest {
	req.Ids = ids
	return req
}

// Do executes the request against the context and client.
func (req *DeleteNotificationsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/notification", true, nil, req, nil)
}

// RpcRequest is a Rpc request.
type RpcRequest struct {
	id      string
	payload interface{}
	httpKey string
}

// Rpc creates a new Rpc request.
func Rpc(id string) *RpcRequest {
	return &RpcRequest{
		id: id,
	}
}

// WithPayload sets the payload on the request.
func (req *RpcRequest) WithPayload(payload interface{}) *RpcRequest {
	req.payload = payload
	return req
}

// WithHttpKey sets the httpKey on the request.
func (req *RpcRequest) WithHttpKey(httpKey string) *RpcRequest {
	req.httpKey = httpKey
	return req
}

// Do executes the request against the context and client.
func (req *RpcRequest) Do(ctx context.Context, cl *Client, v interface{}) error {
	query := url.Values{}
	query.Set("unwrap", "true")
	if req.httpKey != "" {
		query.Set("http_key", req.httpKey)
	}
	return cl.Do(ctx, "POST", "v2/rpc/"+req.id, req.httpKey == "", query, req.payload, v)
}

// SessionLogoutRequest is a SessionLogout request.
type SessionLogoutRequest struct {
	nkapi.SessionLogoutRequest
}

// SessionLogout creates a new SessionLogout request.
func SessionLogout() *SessionLogoutRequest {
	return &SessionLogoutRequest{}
}

// WithToken sets the token on the request.
func (req *SessionLogoutRequest) WithToken(token string) *SessionLogoutRequest {
	req.Token = token
	return req
}

// WithRefreshToken sets the refreshToken on the request.
func (req *SessionLogoutRequest) WithRefreshToken(refreshToken string) *SessionLogoutRequest {
	req.RefreshToken = refreshToken
	return req
}

// Do executes the request against the context and client.
func (req *SessionLogoutRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/session/logout", true, nil, req, nil)
}

// WriteStorageObject is the write storage object.
type WriteStorageObject = nkapi.WriteStorageObject

// ReadStorageObjectsRequest is a ReadStorageObjects request.
type ReadStorageObjectsRequest struct {
	nkapi.ReadStorageObjectsRequest
}

// ReadStorageObjects creates a new ReadStorageObjects request.
func ReadStorageObjects() *ReadStorageObjectsRequest {
	return &ReadStorageObjectsRequest{}
}

// WithObjectId sets the objectId on the request.
func (req *ReadStorageObjectsRequest) WithObjectId(collection, key, userId string) *ReadStorageObjectsRequest {
	req.ObjectIds = append(req.ObjectIds, &nkapi.ReadStorageObjectId{
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

// ReadStorageObjectsResponse is the ReadStorageObjects response.
type ReadStorageObjectsResponse = nkapi.StorageObjects

// WriteStorageObjectsRequest is a WriteStorageObjects request.
type WriteStorageObjectsRequest struct {
	nkapi.WriteStorageObjectsRequest
}

// WriteStorageObjects creates a new WriteStorageObjects request.
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

// WriteStorageObjectsResponse is the WriteStorageObjects response.
type WriteStorageObjectsResponse = nkapi.StorageObjectAcks

// DeleteStorageObjectsRequest is a DeleteStorageObjects request.
type DeleteStorageObjectsRequest struct {
	nkapi.DeleteStorageObjectsRequest
}

// DeleteStorageObjects creates a new DeleteStorageObjects request.
func DeleteStorageObjects() *DeleteStorageObjectsRequest {
	return &DeleteStorageObjectsRequest{}
}

// WithObjectId sets the objectId on the request.
func (req *DeleteStorageObjectsRequest) WithObjectId(collection, key, version string) *DeleteStorageObjectsRequest {
	req.ObjectIds = append(req.ObjectIds, &nkapi.DeleteStorageObjectId{
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

// ListStorageObjectsRequest is a ListStorageObjects request.
type ListStorageObjectsRequest struct {
	nkapi.ListStorageObjectsRequest
}

// ListStorageObjects creates a new ListStorageObjects request.
func ListStorageObjects(collection string) *ListStorageObjectsRequest {
	return &ListStorageObjectsRequest{
		ListStorageObjectsRequest: nkapi.ListStorageObjectsRequest{
			Collection: collection,
			Limit:      wrapperspb.Int32(100),
		},
	}
}

// WithUserId sets the userId on the request.
func (req *ListStorageObjectsRequest) WithUserId(userId string) *ListStorageObjectsRequest {
	req.UserId = userId
	return req
}

// WithLimit sets the limit on the request.
func (req *ListStorageObjectsRequest) WithLimit(limit int) *ListStorageObjectsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListStorageObjectsRequest) WithCursor(cursor string) *ListStorageObjectsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListStorageObjectsRequest) Do(ctx context.Context, cl *Client) (*ListStorageObjectsResponse, error) {
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
	res := new(ListStorageObjectsResponse)
	if err := cl.Do(ctx, "GET", "v2/storage/"+req.Collection, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListStorageObjectsResponse is the ListStorageObjects response.
type ListStorageObjectsResponse = nkapi.StorageObjectList

// ListTournamentsRequest is a ListTournaments request.
type ListTournamentsRequest struct {
	nkapi.ListTournamentsRequest
}

// ListTournaments creates a new ListTournaments request.
func ListTournaments() *ListTournamentsRequest {
	return &ListTournamentsRequest{
		ListTournamentsRequest: nkapi.ListTournamentsRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithCategoryStart sets the categoryStart on the request.
func (req *ListTournamentsRequest) WithCategoryStart(categoryStart uint32) *ListTournamentsRequest {
	req.CategoryStart = wrapperspb.UInt32(categoryStart)
	return req
}

// WithCategoryEnd sets the categoryEnd on the request.
func (req *ListTournamentsRequest) WithCategoryEnd(categoryEnd uint32) *ListTournamentsRequest {
	req.CategoryEnd = wrapperspb.UInt32(categoryEnd)
	return req
}

// WithLimit sets the limit on the request.
func (req *ListTournamentsRequest) WithLimit(limit int) *ListTournamentsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithStartTime sets the startTime on the request.
func (req *ListTournamentsRequest) WithStartTime(startTime uint32) *ListTournamentsRequest {
	req.StartTime = wrapperspb.UInt32(startTime)
	return req
}

// WithEndTime sets the endTime on the request.
func (req *ListTournamentsRequest) WithEndTime(endTime uint32) *ListTournamentsRequest {
	req.EndTime = wrapperspb.UInt32(endTime)
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListTournamentsRequest) WithCursor(cursor string) *ListTournamentsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListTournamentsRequest) Do(ctx context.Context, cl *Client) (*ListTournamentsResponse, error) {
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
	res := new(ListTournamentsResponse)
	if err := cl.Do(ctx, "GET", "v2/tournament", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListTournamentsResponse is the ListTournaments response.
type ListTournamentsResponse = nkapi.TournamentList

// ListTournamentRecordsRequest is a ListTournamentRecords request.
type ListTournamentRecordsRequest struct {
	nkapi.ListTournamentRecordsRequest
}

// ListTournamentRecords creates a new ListTournamentRecords request.
func ListTournamentRecords(tournamentId string) *ListTournamentRecordsRequest {
	return &ListTournamentRecordsRequest{
		ListTournamentRecordsRequest: nkapi.ListTournamentRecordsRequest{
			TournamentId: tournamentId,
			Limit:        wrapperspb.Int32(100),
		},
	}
}

// WithOwnerIds sets the ownerIds on the request.
func (req *ListTournamentRecordsRequest) WithOwnerIds(ownerIds ...string) *ListTournamentRecordsRequest {
	req.OwnerIds = ownerIds
	return req
}

// WithLimit sets the limit on the request.
func (req *ListTournamentRecordsRequest) WithLimit(limit int) *ListTournamentRecordsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithExpiry sets the expiry on the request.
func (req *ListTournamentRecordsRequest) WithExpiry(expiry int64) *ListTournamentRecordsRequest {
	req.Expiry = wrapperspb.Int64(expiry)
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListTournamentRecordsRequest) WithCursor(cursor string) *ListTournamentRecordsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListTournamentRecordsRequest) Do(ctx context.Context, cl *Client) (*ListTournamentRecordsResponse, error) {
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
	res := new(ListTournamentRecordsResponse)
	if err := cl.Do(ctx, "GET", "v2/tournament/"+req.TournamentId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListTournamentRecordsResponse is the ListTournamentRecords response.
type ListTournamentRecordsResponse = nkapi.TournamentRecordList

// WriteTournamentRecordRequest is a WriteTournamentRecord request.
type WriteTournamentRecordRequest struct {
	nkapi.WriteTournamentRecordRequest
}

// WriteTournamentRecord creates a new WriteTournamentRecord request.
func WriteTournamentRecord(tournamentId string) *WriteTournamentRecordRequest {
	return &WriteTournamentRecordRequest{
		WriteTournamentRecordRequest: nkapi.WriteTournamentRecordRequest{
			TournamentId: tournamentId,
			Record:       &nkapi.WriteTournamentRecordRequest_TournamentRecordWrite{},
		},
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
func (req *WriteTournamentRecordRequest) WithOperator(operator Operator) *WriteTournamentRecordRequest {
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

// WriteTournamentRecordResponse is the WriteTournamentRecord response.
type WriteTournamentRecordResponse = nkapi.LeaderboardRecord

// JoinTournamentRequest is a JoinTournament request.
type JoinTournamentRequest struct {
	nkapi.JoinTournamentRequest
}

// JoinTournament creates a new JoinTournament request.
func JoinTournament(tournamentId string) *JoinTournamentRequest {
	return &JoinTournamentRequest{
		JoinTournamentRequest: nkapi.JoinTournamentRequest{
			TournamentId: tournamentId,
		},
	}
}

// Do executes the request against the context and client.
func (req *JoinTournamentRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/tournament/"+req.TournamentId+"/join", true, nil, nil, nil)
}

// ListTournamentRecordsAroundOwnerRequest is a ListTournamentRecordsAroundOwner request.
type ListTournamentRecordsAroundOwnerRequest struct {
	nkapi.ListTournamentRecordsAroundOwnerRequest
}

// ListTournamentRecordsAroundOwner creates a new ListTournamentRecordsAroundOwner request.
func ListTournamentRecordsAroundOwner(tournamentId, ownerId string) *ListTournamentRecordsAroundOwnerRequest {
	return &ListTournamentRecordsAroundOwnerRequest{
		ListTournamentRecordsAroundOwnerRequest: nkapi.ListTournamentRecordsAroundOwnerRequest{
			TournamentId: tournamentId,
			OwnerId:      ownerId,
			Limit:        wrapperspb.UInt32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListTournamentRecordsAroundOwnerRequest) WithLimit(limit int) *ListTournamentRecordsAroundOwnerRequest {
	req.Limit = wrapperspb.UInt32(uint32(limit))
	return req
}

// WithExpiry sets the expiry on the request.
func (req *ListTournamentRecordsAroundOwnerRequest) WithExpiry(expiry int) *ListTournamentRecordsAroundOwnerRequest {
	req.Expiry = wrapperspb.Int64(int64(expiry))
	return req
}

// Do executes the request against the context and client.
func (req *ListTournamentRecordsAroundOwnerRequest) Do(ctx context.Context, cl *Client) (*ListTournamentRecordsAroundOwnerResponse, error) {
	query := url.Values{}
	if req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.Limit.Value), 10))
	}
	if req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.Expiry.Value), 10))
	}
	/*
		if req.Cursor != "" {
			query.Set("cursor", req.Cursor)
		}
	*/
	res := new(ListTournamentRecordsAroundOwnerResponse)
	if err := cl.Do(ctx, "GET", "v2/tournament/"+req.TournamentId+"/owner/"+req.OwnerId, true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListTournamentRecordsAroundOwnerResponse is the ListTournamentRecordsAroundOwner response.
type ListTournamentRecordsAroundOwnerResponse = nkapi.TournamentRecordList

// UsersRequest is a Users request.
type UsersRequest struct {
	nkapi.GetUsersRequest
}

// Users creates a new Users request.
func Users() *UsersRequest {
	return &UsersRequest{}
}

// WithIds sets the ids on the request.
func (req *UsersRequest) WithIds(ids ...string) *UsersRequest {
	req.Ids = ids
	return req
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

// UsersResponse is the Users response.
type UsersResponse = nkapi.Users

// ListUserGroupsRequest is a ListUserGroups request.
type ListUserGroupsRequest struct {
	nkapi.ListUserGroupsRequest
}

// ListUserGroups creates a new ListUserGroups request.
func ListUserGroups(userId string) *ListUserGroupsRequest {
	return &ListUserGroupsRequest{
		ListUserGroupsRequest: nkapi.ListUserGroupsRequest{
			UserId: userId,
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListUserGroupsRequest) WithLimit(limit int) *ListUserGroupsRequest {
	req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithState sets the state on the request.
func (req *ListUserGroupsRequest) WithState(state int) *ListUserGroupsRequest {
	req.State = wrapperspb.Int32(int32(state))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListUserGroupsRequest) WithCursor(cursor string) *ListUserGroupsRequest {
	req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListUserGroupsRequest) Do(ctx context.Context, cl *Client) (*ListUserGroupsResponse, error) {
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
	res := new(ListUserGroupsResponse)
	if err := cl.Do(ctx, "GET", "v2/user/"+req.UserId+"/group", true, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListUserGroupsResponse is the ListUserGroups response.
type ListUserGroupsResponse = nkapi.UserGroupList
