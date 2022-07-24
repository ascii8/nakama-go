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
	return cl.Do(ctx, "GET", "healthcheck", nil, nil, nil)
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
	if err := cl.Do(ctx, "GET", "v2/account", nil, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AccountResponse is the account repsonse.
type AccountResponse = nkapi.Account

// UpdateAccountRequest is a UpdateAccount request.
type UpdateAccountRequest struct {
	req *nkapi.UpdateAccountRequest
}

// UpdateAccount creates a new UpdateAccount request.
func UpdateAccount() *UpdateAccountRequest {
	return &UpdateAccountRequest{
		req: &nkapi.UpdateAccountRequest{},
	}
}

// WithUsername sets the username on the request.
func (req *UpdateAccountRequest) WithUsername(username string) *UpdateAccountRequest {
	req.req.Username = wrapperspb.String(username)
	return req
}

// WithDisplayName sets the displayName on the request.
func (req *UpdateAccountRequest) WithDisplayName(displayName string) *UpdateAccountRequest {
	req.req.DisplayName = wrapperspb.String(displayName)
	return req
}

// WithAvatarUrl sets the avatarUrl on the request.
func (req *UpdateAccountRequest) WithAvatarUrl(avatarUrl string) *UpdateAccountRequest {
	req.req.AvatarUrl = wrapperspb.String(avatarUrl)
	return req
}

// WithLangTag sets the langTag on the request.
func (req *UpdateAccountRequest) WithLangTag(langTag string) *UpdateAccountRequest {
	req.req.LangTag = wrapperspb.String(langTag)
	return req
}

// WithLocation sets the location on the request.
func (req *UpdateAccountRequest) WithLocation(location string) *UpdateAccountRequest {
	req.req.Location = wrapperspb.String(location)
	return req
}

// WithTimezone sets the timezone on the request.
func (req *UpdateAccountRequest) WithTimezone(timezone string) *UpdateAccountRequest {
	req.req.Timezone = wrapperspb.String(timezone)
	return req
}

// Do executes the request against the context and client.
func (req *UpdateAccountRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "PUT", "v2/account", nil, req.req, nil)
}

// SessionResponse is the authenticate repsonse.
type SessionResponse = nkapi.Session

// AuthenticateAppleRequest is a AuthenticateApple request.
type AuthenticateAppleRequest struct {
	req *nkapi.AuthenticateAppleRequest
}

// AuthenticateApple creates a new AuthenticateApple request.
func AuthenticateApple() *AuthenticateAppleRequest {
	return &AuthenticateAppleRequest{
		req: &nkapi.AuthenticateAppleRequest{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateAppleRequest) WithCreate(create bool) *AuthenticateAppleRequest {
	req.req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateAppleRequest) WithUsername(username string) *AuthenticateAppleRequest {
	req.req.Username = username
	return req
}

// WithToken sets the token on the request.
func (req *AuthenticateAppleRequest) WithToken(token string) *AuthenticateAppleRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountApple{}
	}
	req.req.Account.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateAppleRequest) WithVars(vars map[string]string) *AuthenticateAppleRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountApple{}
	}
	req.req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateAppleRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.req.Create != nil {
		query.Set("create", strconv.FormatBool(req.req.Create.Value))
	}
	if req.req.Username != "" {
		query.Set("username", req.req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/apple", query, req.req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateCustomRequest is a AuthenticateCustom request.
type AuthenticateCustomRequest struct {
	req *nkapi.AuthenticateCustomRequest
}

// AuthenticateCustom creates a new AuthenticateCustom request.
func AuthenticateCustom() *AuthenticateCustomRequest {
	return &AuthenticateCustomRequest{
		req: &nkapi.AuthenticateCustomRequest{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateCustomRequest) WithCreate(create bool) *AuthenticateCustomRequest {
	req.req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateCustomRequest) WithUsername(username string) *AuthenticateCustomRequest {
	req.req.Username = username
	return req
}

// WithId sets the id on the request.
func (req *AuthenticateCustomRequest) WithId(id string) *AuthenticateCustomRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountCustom{}
	}
	req.req.Account.Id = id
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateCustomRequest) WithVars(vars map[string]string) *AuthenticateCustomRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountCustom{}
	}
	req.req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateCustomRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.req.Create != nil {
		query.Set("create", strconv.FormatBool(req.req.Create.Value))
	}
	if req.req.Username != "" {
		query.Set("username", req.req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/custom", query, req.req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateDeviceRequest is a AuthenticateDevice request.
type AuthenticateDeviceRequest struct {
	req *nkapi.AuthenticateDeviceRequest
}

// AuthenticateDevice creates a new AuthenticateDevice request.
func AuthenticateDevice() *AuthenticateDeviceRequest {
	return &AuthenticateDeviceRequest{
		req: &nkapi.AuthenticateDeviceRequest{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateDeviceRequest) WithCreate(create bool) *AuthenticateDeviceRequest {
	req.req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateDeviceRequest) WithUsername(username string) *AuthenticateDeviceRequest {
	req.req.Username = username
	return req
}

// WithId sets the id on the request.
func (req *AuthenticateDeviceRequest) WithId(id string) *AuthenticateDeviceRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountDevice{}
	}
	req.req.Account.Id = id
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateDeviceRequest) WithVars(vars map[string]string) *AuthenticateDeviceRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountDevice{}
	}
	req.req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateDeviceRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.req.Create != nil {
		query.Set("create", strconv.FormatBool(req.req.Create.Value))
	}
	if req.req.Username != "" {
		query.Set("username", req.req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/device", query, req.req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateEmailRequest is a AuthenticateEmail request.
type AuthenticateEmailRequest struct {
	req *nkapi.AuthenticateEmailRequest
}

// AuthenticateEmail creates a new AuthenticateEmail request.
func AuthenticateEmail() *AuthenticateEmailRequest {
	return &AuthenticateEmailRequest{
		req: &nkapi.AuthenticateEmailRequest{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateEmailRequest) WithCreate(create bool) *AuthenticateEmailRequest {
	req.req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateEmailRequest) WithUsername(username string) *AuthenticateEmailRequest {
	req.req.Username = username
	return req
}

// WithEmail sets the email on the request.
func (req *AuthenticateEmailRequest) WithEmail(email string) *AuthenticateEmailRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountEmail{}
	}
	req.req.Account.Email = email
	return req
}

// WithPassword sets the password on the request.
func (req *AuthenticateEmailRequest) WithPassword(password string) *AuthenticateEmailRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountEmail{}
	}
	req.req.Account.Password = password
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateEmailRequest) WithVars(vars map[string]string) *AuthenticateEmailRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountEmail{}
	}
	req.req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateEmailRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.req.Create != nil {
		query.Set("create", strconv.FormatBool(req.req.Create.Value))
	}
	if req.req.Username != "" {
		query.Set("username", req.req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/email", query, req.req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateFacebookRequest is a AuthenticateFacebook request.
type AuthenticateFacebookRequest struct {
	req *nkapi.AuthenticateFacebookRequest
}

// AuthenticateFacebook creates a new AuthenticateFacebook request.
func AuthenticateFacebook() *AuthenticateFacebookRequest {
	return &AuthenticateFacebookRequest{
		req: &nkapi.AuthenticateFacebookRequest{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateFacebookRequest) WithCreate(create bool) *AuthenticateFacebookRequest {
	req.req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateFacebookRequest) WithUsername(username string) *AuthenticateFacebookRequest {
	req.req.Username = username
	return req
}

// WithSync sets the sync on the request.
func (req *AuthenticateFacebookRequest) WithSync(sync bool) *AuthenticateFacebookRequest {
	req.req.Sync = wrapperspb.Bool(sync)
	return req
}

// WithToken sets the token on the request.
func (req *AuthenticateFacebookRequest) WithToken(token string) *AuthenticateFacebookRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountFacebook{}
	}
	req.req.Account.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateFacebookRequest) WithVars(vars map[string]string) *AuthenticateFacebookRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountFacebook{}
	}
	req.req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateFacebookRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.req.Create != nil {
		query.Set("create", strconv.FormatBool(req.req.Create.Value))
	}
	if req.req.Username != "" {
		query.Set("username", req.req.Username)
	}
	if req.req.Sync != nil {
		query.Set("sync", strconv.FormatBool(req.req.Sync.Value))
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/facebook", query, req.req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateFacebookInstantGameRequest is a AuthenticateFacebookInstantGame request.
type AuthenticateFacebookInstantGameRequest struct {
	req *nkapi.AuthenticateFacebookInstantGameRequest
}

// AuthenticateFacebookInstantGame creates a new AuthenticateFacebookInstantGame request.
func AuthenticateFacebookInstantGame() *AuthenticateFacebookInstantGameRequest {
	return &AuthenticateFacebookInstantGameRequest{
		req: &nkapi.AuthenticateFacebookInstantGameRequest{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateFacebookInstantGameRequest) WithCreate(create bool) *AuthenticateFacebookInstantGameRequest {
	req.req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateFacebookInstantGameRequest) WithUsername(username string) *AuthenticateFacebookInstantGameRequest {
	req.req.Username = username
	return req
}

// WithSignedPlayerInfo sets the signedPlayerInfo on the request.
func (req *AuthenticateFacebookInstantGameRequest) WithSignedPlayerInfo(signedPlayerInfo string) *AuthenticateFacebookInstantGameRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountFacebookInstantGame{}
	}
	req.req.Account.SignedPlayerInfo = signedPlayerInfo
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateFacebookInstantGameRequest) WithVars(vars map[string]string) *AuthenticateFacebookInstantGameRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountFacebookInstantGame{}
	}
	req.req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateFacebookInstantGameRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.req.Create != nil {
		query.Set("create", strconv.FormatBool(req.req.Create.Value))
	}
	if req.req.Username != "" {
		query.Set("username", req.req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/facebookinstantgame", query, req.req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateGameCenterRequest is a AuthenticateGameCenter request.
type AuthenticateGameCenterRequest struct {
	req *nkapi.AuthenticateGameCenterRequest
}

// AuthenticateGameCenter creates a new AuthenticateGameCenter request.
func AuthenticateGameCenter() *AuthenticateGameCenterRequest {
	return &AuthenticateGameCenterRequest{
		req: &nkapi.AuthenticateGameCenterRequest{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateGameCenterRequest) WithCreate(create bool) *AuthenticateGameCenterRequest {
	req.req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateGameCenterRequest) WithUsername(username string) *AuthenticateGameCenterRequest {
	req.req.Username = username
	return req
}

// WithPlayerId sets the playerId on the request.
func (req *AuthenticateGameCenterRequest) WithPlayerId(playerId string) *AuthenticateGameCenterRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountGameCenter{}
	}
	req.req.Account.PlayerId = playerId
	return req
}

// WithBundleId sets the bundleId on the request.
func (req *AuthenticateGameCenterRequest) WithBundleId(bundleId string) *AuthenticateGameCenterRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountGameCenter{}
	}
	req.req.Account.BundleId = bundleId
	return req
}

// WithTimestampSeconds sets the timestampSeconds on the request.
func (req *AuthenticateGameCenterRequest) WithTimestampSeconds(timestampSeconds int64) *AuthenticateGameCenterRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountGameCenter{}
	}
	req.req.Account.TimestampSeconds = timestampSeconds
	return req
}

// WithSalt sets the salt on the request.
func (req *AuthenticateGameCenterRequest) WithSalt(salt string) *AuthenticateGameCenterRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountGameCenter{}
	}
	req.req.Account.Salt = salt
	return req
}

// WithSignature sets the signature on the request.
func (req *AuthenticateGameCenterRequest) WithSignature(signature string) *AuthenticateGameCenterRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountGameCenter{}
	}
	req.req.Account.Signature = signature
	return req
}

// WithPublicKeyUrl sets the publicKeyUrl on the request.
func (req *AuthenticateGameCenterRequest) WithPublicKeyUrl(publicKeyUrl string) *AuthenticateGameCenterRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountGameCenter{}
	}
	req.req.Account.PublicKeyUrl = publicKeyUrl
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateGameCenterRequest) WithVars(vars map[string]string) *AuthenticateGameCenterRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountGameCenter{}
	}
	req.req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateGameCenterRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.req.Create != nil {
		query.Set("create", strconv.FormatBool(req.req.Create.Value))
	}
	if req.req.Username != "" {
		query.Set("username", req.req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/gamecenter", query, req.req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateGoogleRequest is a AuthenticateGoogle request.
type AuthenticateGoogleRequest struct {
	req *nkapi.AuthenticateGoogleRequest
}

// AuthenticateGoogle creates a new AuthenticateGoogle request.
func AuthenticateGoogle() *AuthenticateGoogleRequest {
	return &AuthenticateGoogleRequest{
		req: &nkapi.AuthenticateGoogleRequest{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateGoogleRequest) WithCreate(create bool) *AuthenticateGoogleRequest {
	req.req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateGoogleRequest) WithUsername(username string) *AuthenticateGoogleRequest {
	req.req.Username = username
	return req
}

// WithToken sets the token on the request.
func (req *AuthenticateGoogleRequest) WithToken(token string) *AuthenticateGoogleRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountGoogle{}
	}
	req.req.Account.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateGoogleRequest) WithVars(vars map[string]string) *AuthenticateGoogleRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountGoogle{}
	}
	req.req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateGoogleRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.req.Create != nil {
		query.Set("create", strconv.FormatBool(req.req.Create.Value))
	}
	if req.req.Username != "" {
		query.Set("username", req.req.Username)
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/google", query, req.req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// AuthenticateSteamRequest is a AuthenticateSteam request.
type AuthenticateSteamRequest struct {
	req *nkapi.AuthenticateSteamRequest
}

// AuthenticateSteam creates a new AuthenticateSteam request.
func AuthenticateSteam() *AuthenticateSteamRequest {
	return &AuthenticateSteamRequest{
		req: &nkapi.AuthenticateSteamRequest{},
	}
}

// WithCreate sets the create on the request.
func (req *AuthenticateSteamRequest) WithCreate(create bool) *AuthenticateSteamRequest {
	req.req.Create = wrapperspb.Bool(create)
	return req
}

// WithUsername sets the username on the request.
func (req *AuthenticateSteamRequest) WithUsername(username string) *AuthenticateSteamRequest {
	req.req.Username = username
	return req
}

// WithSync sets the sync on the request.
func (req *AuthenticateSteamRequest) WithSync(sync bool) *AuthenticateSteamRequest {
	req.req.Sync = wrapperspb.Bool(sync)
	return req
}

// WithToken sets the token on the request.
func (req *AuthenticateSteamRequest) WithToken(token string) *AuthenticateSteamRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountSteam{}
	}
	req.req.Account.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *AuthenticateSteamRequest) WithVars(vars map[string]string) *AuthenticateSteamRequest {
	if req.req.Account == nil {
		req.req.Account = &nkapi.AccountSteam{}
	}
	req.req.Account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *AuthenticateSteamRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	query := url.Values{}
	if req.req.Create != nil {
		query.Set("create", strconv.FormatBool(req.req.Create.Value))
	}
	if req.req.Username != "" {
		query.Set("username", req.req.Username)
	}
	if req.req.Sync != nil {
		query.Set("sync", strconv.FormatBool(req.req.Sync.Value))
	}
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/authenticate/steam", query, req.req.Account, res); err != nil {
		return nil, err
	}
	return res, nil
}

// LinkAppleRequest is a LinkApple request.
type LinkAppleRequest struct {
	req *nkapi.AccountApple
}

// LinkApple creates a new LinkApple request.
func LinkApple() *LinkAppleRequest {
	return &LinkAppleRequest{
		req: &nkapi.AccountApple{},
	}
}

// WithToken sets the token on the request.
func (req *LinkAppleRequest) WithToken(token string) *LinkAppleRequest {
	if req.req == nil {
		req.req = &nkapi.AccountApple{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *LinkAppleRequest) WithVars(vars map[string]string) *LinkAppleRequest {
	if req.req == nil {
		req.req = &nkapi.AccountApple{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkAppleRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/apple", nil, req.req, nil)
}

// LinkCustomRequest is a LinkCustom request.
type LinkCustomRequest struct {
	req *nkapi.AccountCustom
}

// LinkCustom creates a new LinkCustom request.
func LinkCustom() *LinkCustomRequest {
	return &LinkCustomRequest{
		req: &nkapi.AccountCustom{},
	}
}

// WithId sets the id on the request.
func (req *LinkCustomRequest) WithId(id string) *LinkCustomRequest {
	if req.req == nil {
		req.req = &nkapi.AccountCustom{}
	}
	req.req.Id = id
	return req
}

// WithVars sets the vars on the request.
func (req *LinkCustomRequest) WithVars(vars map[string]string) *LinkCustomRequest {
	if req.req == nil {
		req.req = &nkapi.AccountCustom{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkCustomRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/custom", nil, req.req, nil)
}

// LinkDeviceRequest is a LinkDevice request.
type LinkDeviceRequest struct {
	req *nkapi.AccountDevice
}

// LinkDevice creates a new LinkDevice request.
func LinkDevice() *LinkDeviceRequest {
	return &LinkDeviceRequest{
		req: &nkapi.AccountDevice{},
	}
}

// WithId sets the id on the request.
func (req *LinkDeviceRequest) WithId(id string) *LinkDeviceRequest {
	if req.req == nil {
		req.req = &nkapi.AccountDevice{}
	}
	req.req.Id = id
	return req
}

// WithVars sets the vars on the request.
func (req *LinkDeviceRequest) WithVars(vars map[string]string) *LinkDeviceRequest {
	if req.req == nil {
		req.req = &nkapi.AccountDevice{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkDeviceRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/device", nil, req.req, nil)
}

// LinkEmailRequest is a LinkEmail request.
type LinkEmailRequest struct {
	req *nkapi.AccountEmail
}

// LinkEmail creates a new LinkEmail request.
func LinkEmail() *LinkEmailRequest {
	return &LinkEmailRequest{
		req: &nkapi.AccountEmail{},
	}
}

// WithEmail sets the email on the request.
func (req *LinkEmailRequest) WithEmail(email string) *LinkEmailRequest {
	if req.req == nil {
		req.req = &nkapi.AccountEmail{}
	}
	req.req.Email = email
	return req
}

// WithPassword sets the password on the request.
func (req *LinkEmailRequest) WithPassword(password string) *LinkEmailRequest {
	if req.req == nil {
		req.req = &nkapi.AccountEmail{}
	}
	req.req.Password = password
	return req
}

// WithVars sets the vars on the request.
func (req *LinkEmailRequest) WithVars(vars map[string]string) *LinkEmailRequest {
	if req.req == nil {
		req.req = &nkapi.AccountEmail{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkEmailRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/email", nil, req.req, nil)
}

// LinkFacebookRequest is a LinkFacebook request.
type LinkFacebookRequest struct {
	sync *wrapperspb.BoolValue
	req  *nkapi.AccountFacebook
}

// LinkFacebook creates a new LinkFacebook request.
func LinkFacebook() *LinkFacebookRequest {
	return &LinkFacebookRequest{
		req: &nkapi.AccountFacebook{},
	}
}

// WithSync sets the sync on the request.
func (req *LinkFacebookRequest) WithSync(sync bool) *LinkFacebookRequest {
	req.sync = wrapperspb.Bool(sync)
	return req
}

// WithToken sets the token on the request.
func (req *LinkFacebookRequest) WithToken(token string) *LinkFacebookRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebook{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *LinkFacebookRequest) WithVars(vars map[string]string) *LinkFacebookRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebook{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkFacebookRequest) Do(ctx context.Context, cl *Client) error {
	query := url.Values{}
	if req.sync != nil {
		query.Set("sync", strconv.FormatBool(req.sync.Value))
	}
	return cl.Do(ctx, "POST", "v2/account/link/facebook", query, req.req, nil)
}

// LinkFacebookInstantGameRequest is a LinkFacebookInstantGame request.
type LinkFacebookInstantGameRequest struct {
	req *nkapi.AccountFacebookInstantGame
}

// LinkFacebookInstantGame creates a new LinkFacebookInstantGame request.
func LinkFacebookInstantGame() *LinkFacebookInstantGameRequest {
	return &LinkFacebookInstantGameRequest{
		req: &nkapi.AccountFacebookInstantGame{},
	}
}

// WithSignedPlayerInfo sets the signedPlayerInfo on the request.
func (req *LinkFacebookInstantGameRequest) WithSignedPlayerInfo(signedPlayerInfo string) *LinkFacebookInstantGameRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebookInstantGame{}
	}
	req.req.SignedPlayerInfo = signedPlayerInfo
	return req
}

// WithVars sets the vars on the request.
func (req *LinkFacebookInstantGameRequest) WithVars(vars map[string]string) *LinkFacebookInstantGameRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebookInstantGame{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkFacebookInstantGameRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/facebookinstantgame", nil, req.req, nil)
}

// LinkGameCenterRequest is a LinkGameCenter request.
type LinkGameCenterRequest struct {
	req *nkapi.AccountGameCenter
}

// LinkGameCenter creates a new LinkGameCenter request.
func LinkGameCenter() *LinkGameCenterRequest {
	return &LinkGameCenterRequest{
		req: &nkapi.AccountGameCenter{},
	}
}

// WithPlayerId sets the playerId on the request.
func (req *LinkGameCenterRequest) WithPlayerId(playerId string) *LinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.PlayerId = playerId
	return req
}

// WithBundleId sets the bundleId on the request.
func (req *LinkGameCenterRequest) WithBundleId(bundleId string) *LinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.BundleId = bundleId
	return req
}

// WithTimestampSeconds sets the timestampSeconds on the request.
func (req *LinkGameCenterRequest) WithTimestampSeconds(timestampSeconds int64) *LinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.TimestampSeconds = timestampSeconds
	return req
}

// WithSalt sets the salt on the request.
func (req *LinkGameCenterRequest) WithSalt(salt string) *LinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.Salt = salt
	return req
}

// WithSignature sets the signature on the request.
func (req *LinkGameCenterRequest) WithSignature(signature string) *LinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.Signature = signature
	return req
}

// WithPublicKeyUrl sets the publicKeyUrl on the request.
func (req *LinkGameCenterRequest) WithPublicKeyUrl(publicKeyUrl string) *LinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.PublicKeyUrl = publicKeyUrl
	return req
}

// WithVars sets the vars on the request.
func (req *LinkGameCenterRequest) WithVars(vars map[string]string) *LinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkGameCenterRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/gamecenter", nil, req.req, nil)
}

// LinkGoogleRequest is a LinkGoogle request.
type LinkGoogleRequest struct {
	req *nkapi.AccountGoogle
}

// LinkGoogle creates a new LinkGoogle request.
func LinkGoogle() *LinkGoogleRequest {
	return &LinkGoogleRequest{
		req: &nkapi.AccountGoogle{},
	}
}

// WithToken sets the token on the request.
func (req *LinkGoogleRequest) WithToken(token string) *LinkGoogleRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGoogle{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *LinkGoogleRequest) WithVars(vars map[string]string) *LinkGoogleRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGoogle{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkGoogleRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/link/google", nil, req.req, nil)
}

// LinkSteamRequest is a LinkSteam request.
type LinkSteamRequest struct {
	account *nkapi.AccountSteam
	sync    *wrapperspb.BoolValue
}

// LinkSteam creates a new LinkSteam request.
func LinkSteam() *LinkSteamRequest {
	return &LinkSteamRequest{
		account: &nkapi.AccountSteam{},
	}
}

// WithSync sets the sync on the request.
func (req *LinkSteamRequest) WithSync(sync bool) *LinkSteamRequest {
	req.sync = wrapperspb.Bool(sync)
	return req
}

// WithToken sets the token on the request.
func (req *LinkSteamRequest) WithToken(token string) *LinkSteamRequest {
	if req.account == nil {
		req.account = &nkapi.AccountSteam{}
	}
	req.account.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *LinkSteamRequest) WithVars(vars map[string]string) *LinkSteamRequest {
	if req.account == nil {
		req.account = &nkapi.AccountSteam{}
	}
	req.account.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *LinkSteamRequest) Do(ctx context.Context, cl *Client) error {
	params := make(map[string]interface{})
	if req.account != nil {
		params["account"] = req.account
	}
	if req.sync != nil {
		params["sync"] = req.sync.Value
	}
	return cl.Do(ctx, "POST", "v2/account/link/steam", nil, params, nil)
}

// SessionRefreshRequest is a SessionRefresh request.
type SessionRefreshRequest struct {
	req *nkapi.SessionRefreshRequest
}

// SessionRefresh creates a new SessionRefresh request.
func SessionRefresh() *SessionRefreshRequest {
	return &SessionRefreshRequest{
		req: &nkapi.SessionRefreshRequest{},
	}
}

// WithToken sets the token on the request.
func (req *SessionRefreshRequest) WithToken(token string) *SessionRefreshRequest {
	if req.req == nil {
		req.req = &nkapi.SessionRefreshRequest{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *SessionRefreshRequest) WithVars(vars map[string]string) *SessionRefreshRequest {
	if req.req == nil {
		req.req = &nkapi.SessionRefreshRequest{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *SessionRefreshRequest) Do(ctx context.Context, cl *Client) (*SessionResponse, error) {
	res := new(SessionResponse)
	if err := cl.Do(ctx, "POST", "v2/account/session/refresh", nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// UnlinkAppleRequest is a UnlinkApple request.
type UnlinkAppleRequest struct {
	req *nkapi.AccountApple
}

// UnlinkApple creates a new UnlinkApple request.
func UnlinkApple() *UnlinkAppleRequest {
	return &UnlinkAppleRequest{
		req: &nkapi.AccountApple{},
	}
}

// WithToken sets the token on the request.
func (req *UnlinkAppleRequest) WithToken(token string) *UnlinkAppleRequest {
	if req.req == nil {
		req.req = &nkapi.AccountApple{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkAppleRequest) WithVars(vars map[string]string) *UnlinkAppleRequest {
	if req.req == nil {
		req.req = &nkapi.AccountApple{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkAppleRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/apple", nil, req.req, nil)
}

// UnlinkCustomRequest is a UnlinkCustom request.
type UnlinkCustomRequest struct {
	req *nkapi.AccountCustom
}

// UnlinkCustom creates a new UnlinkCustom request.
func UnlinkCustom() *UnlinkCustomRequest {
	return &UnlinkCustomRequest{
		req: &nkapi.AccountCustom{},
	}
}

// WithId sets the id on the request.
func (req *UnlinkCustomRequest) WithId(id string) *UnlinkCustomRequest {
	if req.req == nil {
		req.req = &nkapi.AccountCustom{}
	}
	req.req.Id = id
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkCustomRequest) WithVars(vars map[string]string) *UnlinkCustomRequest {
	if req.req == nil {
		req.req = &nkapi.AccountCustom{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkCustomRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/custom", nil, req.req, nil)
}

// UnlinkDeviceRequest is a UnlinkDevice request.
type UnlinkDeviceRequest struct {
	req *nkapi.AccountDevice
}

// UnlinkDevice creates a new UnlinkDevice request.
func UnlinkDevice() *UnlinkDeviceRequest {
	return &UnlinkDeviceRequest{
		req: &nkapi.AccountDevice{},
	}
}

// WithId sets the id on the request.
func (req *UnlinkDeviceRequest) WithId(id string) *UnlinkDeviceRequest {
	if req.req == nil {
		req.req = &nkapi.AccountDevice{}
	}
	req.req.Id = id
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkDeviceRequest) WithVars(vars map[string]string) *UnlinkDeviceRequest {
	if req.req == nil {
		req.req = &nkapi.AccountDevice{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkDeviceRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/device", nil, req.req, nil)
}

// UnlinkEmailRequest is a UnlinkEmail request.
type UnlinkEmailRequest struct {
	req *nkapi.AccountEmail
}

// UnlinkEmail creates a new UnlinkEmail request.
func UnlinkEmail() *UnlinkEmailRequest {
	return &UnlinkEmailRequest{
		req: &nkapi.AccountEmail{},
	}
}

// WithEmail sets the email on the request.
func (req *UnlinkEmailRequest) WithEmail(email string) *UnlinkEmailRequest {
	if req.req == nil {
		req.req = &nkapi.AccountEmail{}
	}
	req.req.Email = email
	return req
}

// WithPassword sets the password on the request.
func (req *UnlinkEmailRequest) WithPassword(password string) *UnlinkEmailRequest {
	if req.req == nil {
		req.req = &nkapi.AccountEmail{}
	}
	req.req.Password = password
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkEmailRequest) WithVars(vars map[string]string) *UnlinkEmailRequest {
	if req.req == nil {
		req.req = &nkapi.AccountEmail{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkEmailRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/email", nil, req.req, nil)
}

// UnlinkFacebookRequest is a UnlinkFacebook request.
type UnlinkFacebookRequest struct {
	req *nkapi.AccountFacebook
}

// UnlinkFacebook creates a new UnlinkFacebook request.
func UnlinkFacebook() *UnlinkFacebookRequest {
	return &UnlinkFacebookRequest{
		req: &nkapi.AccountFacebook{},
	}
}

// WithToken sets the token on the request.
func (req *UnlinkFacebookRequest) WithToken(token string) *UnlinkFacebookRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebook{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkFacebookRequest) WithVars(vars map[string]string) *UnlinkFacebookRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebook{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkFacebookRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/facebook", nil, req.req, nil)
}

// UnlinkFacebookInstantGameRequest is a UnlinkFacebookInstantGame request.
type UnlinkFacebookInstantGameRequest struct {
	req *nkapi.AccountFacebookInstantGame
}

// UnlinkFacebookInstantGame creates a new UnlinkFacebookInstantGame request.
func UnlinkFacebookInstantGame() *UnlinkFacebookInstantGameRequest {
	return &UnlinkFacebookInstantGameRequest{
		req: &nkapi.AccountFacebookInstantGame{},
	}
}

// WithSignedPlayerInfo sets the signedPlayerInfo on the request.
func (req *UnlinkFacebookInstantGameRequest) WithSignedPlayerInfo(signedPlayerInfo string) *UnlinkFacebookInstantGameRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebookInstantGame{}
	}
	req.req.SignedPlayerInfo = signedPlayerInfo
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkFacebookInstantGameRequest) WithVars(vars map[string]string) *UnlinkFacebookInstantGameRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebookInstantGame{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkFacebookInstantGameRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/facebookinstantgame", nil, req.req, nil)
}

// UnlinkGameCenterRequest is a UnlinkGameCenter request.
type UnlinkGameCenterRequest struct {
	req *nkapi.AccountGameCenter
}

// UnlinkGameCenter creates a new UnlinkGameCenter request.
func UnlinkGameCenter() *UnlinkGameCenterRequest {
	return &UnlinkGameCenterRequest{
		req: &nkapi.AccountGameCenter{},
	}
}

// WithPlayerId sets the playerId on the request.
func (req *UnlinkGameCenterRequest) WithPlayerId(playerId string) *UnlinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.PlayerId = playerId
	return req
}

// WithBundleId sets the bundleId on the request.
func (req *UnlinkGameCenterRequest) WithBundleId(bundleId string) *UnlinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.BundleId = bundleId
	return req
}

// WithTimestampSeconds sets the timestampSeconds on the request.
func (req *UnlinkGameCenterRequest) WithTimestampSeconds(timestampSeconds int64) *UnlinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.TimestampSeconds = timestampSeconds
	return req
}

// WithSalt sets the salt on the request.
func (req *UnlinkGameCenterRequest) WithSalt(salt string) *UnlinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.Salt = salt
	return req
}

// WithSignature sets the signature on the request.
func (req *UnlinkGameCenterRequest) WithSignature(signature string) *UnlinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.Signature = signature
	return req
}

// WithPublicKeyUrl sets the publicKeyUrl on the request.
func (req *UnlinkGameCenterRequest) WithPublicKeyUrl(publicKeyUrl string) *UnlinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.PublicKeyUrl = publicKeyUrl
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkGameCenterRequest) WithVars(vars map[string]string) *UnlinkGameCenterRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGameCenter{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkGameCenterRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/gamecenter", nil, req.req, nil)
}

// UnlinkGoogleRequest is a UnlinkGoogle request.
type UnlinkGoogleRequest struct {
	req *nkapi.AccountGoogle
}

// UnlinkGoogle creates a new UnlinkGoogle request.
func UnlinkGoogle() *UnlinkGoogleRequest {
	return &UnlinkGoogleRequest{
		req: &nkapi.AccountGoogle{},
	}
}

// WithToken sets the token on the request.
func (req *UnlinkGoogleRequest) WithToken(token string) *UnlinkGoogleRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGoogle{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkGoogleRequest) WithVars(vars map[string]string) *UnlinkGoogleRequest {
	if req.req == nil {
		req.req = &nkapi.AccountGoogle{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkGoogleRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/google", nil, req.req, nil)
}

// UnlinkSteamRequest is a UnlinkSteam request.
type UnlinkSteamRequest struct {
	req *nkapi.AccountSteam
}

// UnlinkSteam creates a new UnlinkSteam request.
func UnlinkSteam() *UnlinkSteamRequest {
	return &UnlinkSteamRequest{
		req: &nkapi.AccountSteam{},
	}
}

// WithToken sets the token on the request.
func (req *UnlinkSteamRequest) WithToken(token string) *UnlinkSteamRequest {
	if req.req == nil {
		req.req = &nkapi.AccountSteam{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *UnlinkSteamRequest) WithVars(vars map[string]string) *UnlinkSteamRequest {
	if req.req == nil {
		req.req = &nkapi.AccountSteam{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *UnlinkSteamRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/account/unlink/steam", nil, req.req, nil)
}

// ListChannelMessagesRequest is a ListChannelMessages request.
type ListChannelMessagesRequest struct {
	req *nkapi.ListChannelMessagesRequest
}

// ListChannelMessages creates a new ListChannelMessages request.
func ListChannelMessages(channelId string) *ListChannelMessagesRequest {
	return &ListChannelMessagesRequest{
		req: &nkapi.ListChannelMessagesRequest{
			ChannelId: channelId,
			Limit:     wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListChannelMessagesRequest) WithLimit(limit int) *ListChannelMessagesRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithForward sets the forward on the request.
func (req *ListChannelMessagesRequest) WithForward(forward bool) *ListChannelMessagesRequest {
	req.req.Forward = wrapperspb.Bool(forward)
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListChannelMessagesRequest) WithCursor(cursor string) *ListChannelMessagesRequest {
	req.req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListChannelMessagesRequest) Do(ctx context.Context, cl *Client) (*ListChannelMessagesResponse, error) {
	query := url.Values{}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.Forward != nil {
		query.Set("forward", strconv.FormatBool(req.req.Forward.Value))
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	res := new(ListChannelMessagesResponse)
	if err := cl.Do(ctx, "GET", "v2/channel/"+req.req.ChannelId, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListChannelMessagesResponse is the ListChannelMessages response.
type ListChannelMessagesResponse = nkapi.ChannelMessageList

// EventRequest is a Event request.
type EventRequest struct {
	req *nkapi.Event
}

// Event creates a new Event request.
func Event() *EventRequest {
	return &EventRequest{
		req: &nkapi.Event{},
	}
}

// WithName sets the name on the request.
func (req *EventRequest) WithName(name string) *EventRequest {
	req.req.Name = name
	return req
}

// WithProperties sets the properties on the request.
func (req *EventRequest) WithProperties(properties map[string]string) *EventRequest {
	req.req.Properties = properties
	return req
}

// WithTimestamp sets the timestamp on the request.
func (req *EventRequest) WithTimestamp(t time.Time) *EventRequest {
	req.req.Timestamp = timestamppb.New(t)
	return req
}

// WithExternal sets the external on the request.
func (req *EventRequest) WithExternal(external bool) *EventRequest {
	req.req.External = external
	return req
}

// Do executes the request against the context and client.
func (req *EventRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/event", nil, req.req, nil)
}

// ListFriendsRequest is a ListFriends request.
type ListFriendsRequest struct {
	req *nkapi.ListFriendsRequest
}

// ListFriends creates a new ListFriends request.
func ListFriends() *ListFriendsRequest {
	return &ListFriendsRequest{
		req: &nkapi.ListFriendsRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListFriendsRequest) WithLimit(limit int) *ListFriendsRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithState sets the state on the request.
func (req *ListFriendsRequest) WithState(state int) *ListFriendsRequest {
	req.req.State = wrapperspb.Int32(int32(state))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListFriendsRequest) WithCursor(cursor string) *ListFriendsRequest {
	req.req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListFriendsRequest) Do(ctx context.Context, cl *Client) (*ListFriendsResponse, error) {
	query := url.Values{}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.State != nil {
		query.Set("state", strconv.FormatInt(int64(req.req.State.Value), 10))
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	res := new(ListFriendsResponse)
	if err := cl.Do(ctx, "GET", "v2/friend", query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListFriendsResponse is the ListFriends response.
type ListFriendsResponse = nkapi.FriendList

// DeleteFriendsRequest is a DeleteFriends request.
type DeleteFriendsRequest struct {
	req *nkapi.DeleteFriendsRequest
}

// DeleteFriends creates a new DeleteFriends request.
func DeleteFriends() *DeleteFriendsRequest {
	return &DeleteFriendsRequest{
		req: &nkapi.DeleteFriendsRequest{},
	}
}

// WithIds sets the Ids on the request.
func (req *DeleteFriendsRequest) WithIds(ids ...string) *DeleteFriendsRequest {
	req.req.Ids = ids
	return req
}

// WithUsernames sets the Usernames on the request.
func (req *DeleteFriendsRequest) WithUsernames(usernames ...string) *DeleteFriendsRequest {
	req.req.Usernames = usernames
	return req
}

// Do executes the request against the context and client.
func (req *DeleteFriendsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/friend", nil, req.req, nil)
}

// AddFriendsRequest is a AddFriends request.
type AddFriendsRequest struct {
	req *nkapi.AddFriendsRequest
}

// AddFriends creates a new AddFriends request.
func AddFriends() *AddFriendsRequest {
	return &AddFriendsRequest{
		req: &nkapi.AddFriendsRequest{},
	}
}

// WithIds sets the Ids on the request.
func (req *AddFriendsRequest) WithIds(ids ...string) *AddFriendsRequest {
	req.req.Ids = ids
	return req
}

// WithUsernames sets the Usernames on the request.
func (req *AddFriendsRequest) WithUsernames(usernames ...string) *AddFriendsRequest {
	req.req.Usernames = usernames
	return req
}

// Do executes the request against the context and client.
func (req *AddFriendsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/friend", nil, req.req, nil)
}

// BlockFriendsRequest is a BlockFriends request.
type BlockFriendsRequest struct {
	req *nkapi.BlockFriendsRequest
}

// BlockFriends creates a new BlockFriends request.
func BlockFriends() *BlockFriendsRequest {
	return &BlockFriendsRequest{
		req: &nkapi.BlockFriendsRequest{},
	}
}

// WithIds sets the Ids on the request.
func (req *BlockFriendsRequest) WithIds(ids ...string) *BlockFriendsRequest {
	req.req.Ids = ids
	return req
}

// WithUsernames sets the Usernames on the request.
func (req *BlockFriendsRequest) WithUsernames(usernames ...string) *BlockFriendsRequest {
	req.req.Usernames = usernames
	return req
}

// Do executes the request against the context and client.
func (req *BlockFriendsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/friend/block", nil, req.req, nil)
}

// ImportFacebookFriendsRequest is a ImportFacebookFriends request.
type ImportFacebookFriendsRequest struct {
	reset bool
	req   *nkapi.AccountFacebook
}

// ImportFacebookFriends creates a new ImportFacebookFriends request.
func ImportFacebookFriends() *ImportFacebookFriendsRequest {
	return &ImportFacebookFriendsRequest{}
}

// WithReset sets the reset on the request.
func (req *ImportFacebookFriendsRequest) WithReset(reset bool) *ImportFacebookFriendsRequest {
	req.reset = reset
	return req
}

// WithToken sets the token on the request.
func (req *ImportFacebookFriendsRequest) WithToken(token string) *ImportFacebookFriendsRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebook{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *ImportFacebookFriendsRequest) WithVars(vars map[string]string) *ImportFacebookFriendsRequest {
	if req.req == nil {
		req.req = &nkapi.AccountFacebook{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *ImportFacebookFriendsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/friend/facebook", nil, req.req, nil)
}

// ImportSteamFriendsRequest is a ImportSteamFriends request.
type ImportSteamFriendsRequest struct {
	reset bool
	req   *nkapi.AccountSteam
}

// ImportSteamFriends creates a new ImportSteamFriends request.
func ImportSteamFriends() *ImportSteamFriendsRequest {
	return &ImportSteamFriendsRequest{}
}

// WithReset sets the reset on the request.
func (req *ImportSteamFriendsRequest) WithReset(reset bool) *ImportSteamFriendsRequest {
	req.reset = reset
	return req
}

// WithToken sets the token on the request.
func (req *ImportSteamFriendsRequest) WithToken(token string) *ImportSteamFriendsRequest {
	if req.req == nil {
		req.req = &nkapi.AccountSteam{}
	}
	req.req.Token = token
	return req
}

// WithVars sets the vars on the request.
func (req *ImportSteamFriendsRequest) WithVars(vars map[string]string) *ImportSteamFriendsRequest {
	if req.req == nil {
		req.req = &nkapi.AccountSteam{}
	}
	req.req.Vars = vars
	return req
}

// Do executes the request against the context and client.
func (req *ImportSteamFriendsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/friend/steam", nil, req.req, nil)
}

// ListGroupsRequest is a ListGroups request.
type ListGroupsRequest struct {
	req *nkapi.ListGroupsRequest
}

// ListGroups creates a new ListGroups request.
func ListGroups() *ListGroupsRequest {
	return &ListGroupsRequest{
		req: &nkapi.ListGroupsRequest{},
	}
}

// WithName sets the name on the request.
func (req *ListGroupsRequest) WithName(name string) *ListGroupsRequest {
	req.req.Name = name
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListGroupsRequest) WithCursor(cursor string) *ListGroupsRequest {
	req.req.Cursor = cursor
	return req
}

// WithLimit sets the limit on the request.
func (req *ListGroupsRequest) WithLimit(limit int) *ListGroupsRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithLangTag sets the langTag on the request.
func (req *ListGroupsRequest) WithLangTag(langTag string) *ListGroupsRequest {
	req.req.LangTag = langTag
	return req
}

// WithMembers sets the members on the request.
func (req *ListGroupsRequest) WithMembers(members int) *ListGroupsRequest {
	req.req.Members = wrapperspb.Int32(int32(members))
	return req
}

// WithOpen sets the open on the request.
func (req *ListGroupsRequest) WithOpen(open bool) *ListGroupsRequest {
	req.req.Open = wrapperspb.Bool(open)
	return req
}

// Do executes the request against the context and client.
func (req *ListGroupsRequest) Do(ctx context.Context, cl *Client) (*ListGroupsResponse, error) {
	query := url.Values{}
	if req.req.Name != "" {
		query.Set("name", req.req.Name)
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.LangTag != "" {
		query.Set("langTag", req.req.LangTag)
	}
	if req.req.Members != nil {
		query.Set("members", strconv.FormatInt(int64(req.req.Members.Value), 10))
	}
	if req.req.Open != nil {
		query.Set("open", strconv.FormatBool(req.req.Open.Value))
	}
	res := new(ListGroupsResponse)
	if err := cl.Do(ctx, "GET", "v2/group", query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListGroupsResponse is the ListGroups response.
type ListGroupsResponse = nkapi.GroupList

// CreateGroupRequest is a CreateGroup request.
type CreateGroupRequest struct {
	req *nkapi.CreateGroupRequest
}

// CreateGroup creates a new CreateGroup request.
func CreateGroup() *CreateGroupRequest {
	return &CreateGroupRequest{
		req: &nkapi.CreateGroupRequest{},
	}
}

// WithName sets the name on the request.
func (req *CreateGroupRequest) WithName(name string) *CreateGroupRequest {
	req.req.Name = name
	return req
}

// WithDescription sets the description on the request.
func (req *CreateGroupRequest) WithDescription(description string) *CreateGroupRequest {
	req.req.Description = description
	return req
}

// WithLangTag sets the langTag on the request.
func (req *CreateGroupRequest) WithLangTag(langTag string) *CreateGroupRequest {
	req.req.LangTag = langTag
	return req
}

// WithAvatarUrl sets the avatarUrl on the request.
func (req *CreateGroupRequest) WithAvatarUrl(avatarUrl string) *CreateGroupRequest {
	req.req.AvatarUrl = avatarUrl
	return req
}

// WithOpen sets the open on the request.
func (req *CreateGroupRequest) WithOpen(open bool) *CreateGroupRequest {
	req.req.Open = open
	return req
}

// WithMaxCount sets the maxCount on the request.
func (req *CreateGroupRequest) WithMaxCount(maxCount int) *CreateGroupRequest {
	req.req.MaxCount = int32(maxCount)
	return req
}

// Do executes the request against the context and client.
func (req *CreateGroupRequest) Do(ctx context.Context, cl *Client) (*nkapi.Group, error) {
	res := new(nkapi.Group)
	if err := cl.Do(ctx, "POST", "v2/group", nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteGroupRequest is a DeleteGroup request.
type DeleteGroupRequest struct {
	groupId string
}

// DeleteGroup creates a new DeleteGroup request.
func DeleteGroup(groupId string) *DeleteGroupRequest {
	return &DeleteGroupRequest{
		groupId: groupId,
	}
}

// Do executes the request against the context and client.
func (req *DeleteGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/group/"+req.groupId, nil, nil, nil)
}

// UpdateGroupRequest is a UpdateGroup request.
type UpdateGroupRequest struct {
	req *nkapi.UpdateGroupRequest
}

// UpdateGroup creates a new UpdateGroup request.
func UpdateGroup(groupId string) *UpdateGroupRequest {
	return &UpdateGroupRequest{
		req: &nkapi.UpdateGroupRequest{
			GroupId: groupId,
		},
	}
}

// WithName sets the name on the request.
func (req *UpdateGroupRequest) WithName(name string) *UpdateGroupRequest {
	req.req.Name = wrapperspb.String(name)
	return req
}

// WithDescription sets the description on the request.
func (req *UpdateGroupRequest) WithDescription(description string) *UpdateGroupRequest {
	req.req.Description = wrapperspb.String(description)
	return req
}

// WithLangTag sets the langTag on the request.
func (req *UpdateGroupRequest) WithLangTag(langTag string) *UpdateGroupRequest {
	req.req.LangTag = wrapperspb.String(langTag)
	return req
}

// WithAvatarUrl sets the avatarUrl on the request.
func (req *UpdateGroupRequest) WithAvatarUrl(avatarUrl string) *UpdateGroupRequest {
	req.req.AvatarUrl = wrapperspb.String(avatarUrl)
	return req
}

// WithOpen sets the open on the request.
func (req *UpdateGroupRequest) WithOpen(open bool) *UpdateGroupRequest {
	req.req.Open = wrapperspb.Bool(open)
	return req
}

// Do executes the request against the context and client.
func (req *UpdateGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "PUT", "v2/group/"+req.req.GroupId, nil, req.req, nil)
}

// AddGroupUsersRequest is a AddGroupUsers request.
type AddGroupUsersRequest struct {
	req *nkapi.AddGroupUsersRequest
}

// AddGroupUsers creates a new AddGroupUsers request.
func AddGroupUsers(groupId string) *AddGroupUsersRequest {
	return &AddGroupUsersRequest{
		req: &nkapi.AddGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *AddGroupUsersRequest) WithUserIds(userIds ...string) *AddGroupUsersRequest {
	req.req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *AddGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	params := make(map[string]interface{})
	if len(req.req.UserIds) != 0 {
		params["userIds"] = req.req.UserIds
	}
	return cl.Do(ctx, "POST", "v2/group/"+req.req.GroupId+"/add", nil, params, nil)
}

// BanGroupUsersRequest is a BanGroupUsers request.
type BanGroupUsersRequest struct {
	req *nkapi.BanGroupUsersRequest
}

// BanGroupUsers creates a new BanGroupUsers request.
func BanGroupUsers(groupId string) *BanGroupUsersRequest {
	return &BanGroupUsersRequest{
		req: &nkapi.BanGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *BanGroupUsersRequest) WithUserIds(userIds ...string) *BanGroupUsersRequest {
	req.req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *BanGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	params := make(map[string]interface{})
	if len(req.req.UserIds) != 0 {
		params["userIds"] = req.req.UserIds
	}
	return cl.Do(ctx, "POST", "v2/group/"+req.req.GroupId+"/ban", nil, params, nil)
}

// DemoteGroupUsersRequest is a DemoteGroupUsers request.
type DemoteGroupUsersRequest struct {
	req *nkapi.DemoteGroupUsersRequest
}

// DemoteGroupUsers creates a new DemoteGroupUsers request.
func DemoteGroupUsers(groupId string) *DemoteGroupUsersRequest {
	return &DemoteGroupUsersRequest{
		req: &nkapi.DemoteGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *DemoteGroupUsersRequest) WithUserIds(userIds ...string) *DemoteGroupUsersRequest {
	req.req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *DemoteGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	params := make(map[string]interface{})
	if len(req.req.UserIds) != 0 {
		params["userIds"] = req.req.UserIds
	}
	return cl.Do(ctx, "POST", "v2/group/"+req.req.GroupId+"/demote", nil, params, nil)
}

// JoinGroupRequest is a JoinGroup request.
type JoinGroupRequest struct {
	req *nkapi.JoinGroupRequest
}

// JoinGroup creates a new JoinGroup request.
func JoinGroup(groupId string) *JoinGroupRequest {
	return &JoinGroupRequest{
		req: &nkapi.JoinGroupRequest{
			GroupId: groupId,
		},
	}
}

// Do executes the request against the context and client.
func (req *JoinGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.req.GroupId+"/join", nil, nil, nil)
}

// KickGroupUsersRequest is a KickGroupUsers request.
type KickGroupUsersRequest struct {
	req *nkapi.KickGroupUsersRequest
}

// KickGroupUsers creates a new KickGroupUsers request.
func KickGroupUsers(groupId string) *KickGroupUsersRequest {
	return &KickGroupUsersRequest{
		req: &nkapi.KickGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *KickGroupUsersRequest) WithUserIds(userIds ...string) *KickGroupUsersRequest {
	req.req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *KickGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	params := make(map[string]interface{})
	if len(req.req.UserIds) != 0 {
		params["userIds"] = req.req.UserIds
	}
	return cl.Do(ctx, "POST", "v2/group/"+req.req.GroupId+"/kick", nil, params, nil)
}

// LeaveGroupRequest is a LeaveGroup request.
type LeaveGroupRequest struct {
	req *nkapi.LeaveGroupRequest
}

// LeaveGroup creates a new LeaveGroup request.
func LeaveGroup(groupId string) *LeaveGroupRequest {
	return &LeaveGroupRequest{
		req: &nkapi.LeaveGroupRequest{
			GroupId: groupId,
		},
	}
}

// Do executes the request against the context and client.
func (req *LeaveGroupRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/group/"+req.req.GroupId+"/leave", nil, nil, nil)
}

// PromoteGroupUsersRequest is a PromoteGroupUsers request.
type PromoteGroupUsersRequest struct {
	req *nkapi.PromoteGroupUsersRequest
}

// PromoteGroupUsers creates a new PromoteGroupUsers request.
func PromoteGroupUsers(groupId string) *PromoteGroupUsersRequest {
	return &PromoteGroupUsersRequest{
		req: &nkapi.PromoteGroupUsersRequest{
			GroupId: groupId,
		},
	}
}

// WithUserIds sets the userIds on the request.
func (req *PromoteGroupUsersRequest) WithUserIds(userIds ...string) *PromoteGroupUsersRequest {
	req.req.UserIds = userIds
	return req
}

// Do executes the request against the context and client.
func (req *PromoteGroupUsersRequest) Do(ctx context.Context, cl *Client) error {
	params := make(map[string]interface{})
	if len(req.req.UserIds) != 0 {
		params["userIds"] = req.req.UserIds
	}
	return cl.Do(ctx, "POST", "v2/group/"+req.req.GroupId+"/promote", nil, params, nil)
}

// ListGroupUsersRequest is a ListGroupUsers request.
type ListGroupUsersRequest struct {
	groupId string
	req     *nkapi.ListGroupUsersRequest
}

// ListGroupUsers creates a new ListGroupUsers request.
func ListGroupUsers(groupId string) *ListGroupUsersRequest {
	return &ListGroupUsersRequest{
		groupId: groupId,
		req: &nkapi.ListGroupUsersRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListGroupUsersRequest) WithLimit(limit int) *ListGroupUsersRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithState sets the state on the request.
func (req *ListGroupUsersRequest) WithState(state int) *ListGroupUsersRequest {
	req.req.State = wrapperspb.Int32(int32(state))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListGroupUsersRequest) WithCursor(cursor string) *ListGroupUsersRequest {
	req.req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListGroupUsersRequest) Do(ctx context.Context, cl *Client) (*ListGroupUsersResponse, error) {
	query := url.Values{}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.State != nil {
		query.Set("state", strconv.FormatInt(int64(req.req.State.Value), 10))
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	res := new(ListGroupUsersResponse)
	if err := cl.Do(ctx, "GET", "v2/friend", query, nil, res); err != nil {
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
	req *nkapi.ValidatePurchaseAppleRequest
}

// ValidatePurchaseApple creates a new ValidatePurchaseApple request.
func ValidatePurchaseApple() *ValidatePurchaseAppleRequest {
	return &ValidatePurchaseAppleRequest{
		req: &nkapi.ValidatePurchaseAppleRequest{},
	}
}

// WithReceipt sets the receipt on the request.
func (req *ValidatePurchaseAppleRequest) WithReceipt(receipt string) *ValidatePurchaseAppleRequest {
	req.req.Receipt = receipt
	return req
}

// WithPersist sets the persist on the request.
func (req *ValidatePurchaseAppleRequest) WithPersist(persist bool) *ValidatePurchaseAppleRequest {
	req.req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidatePurchaseAppleRequest) Do(ctx context.Context, cl *Client) (*ValidatePurchaseResponse, error) {
	res := new(ValidatePurchaseResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/purchase/apple", nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ValidatePurchaseGoogleRequest is a ValidatePurchaseGoogle request.
type ValidatePurchaseGoogleRequest struct {
	req *nkapi.ValidatePurchaseGoogleRequest
}

// ValidatePurchaseGoogle creates a new ValidatePurchaseGoogle request.
func ValidatePurchaseGoogle() *ValidatePurchaseGoogleRequest {
	return &ValidatePurchaseGoogleRequest{
		req: &nkapi.ValidatePurchaseGoogleRequest{},
	}
}

// WithPurchase sets the purchase on the request.
func (req *ValidatePurchaseGoogleRequest) WithPurchase(purchase string) *ValidatePurchaseGoogleRequest {
	req.req.Purchase = purchase
	return req
}

// WithPersist sets the persist on the request.
func (req *ValidatePurchaseGoogleRequest) WithPersist(persist bool) *ValidatePurchaseGoogleRequest {
	req.req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidatePurchaseGoogleRequest) Do(ctx context.Context, cl *Client) (*ValidatePurchaseResponse, error) {
	res := new(ValidatePurchaseResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/purchase/google", nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ValidatePurchaseHuaweiRequest is a ValidatePurchaseHuawei request.
type ValidatePurchaseHuaweiRequest struct {
	req *nkapi.ValidatePurchaseHuaweiRequest
}

// ValidatePurchaseHuawei creates a new ValidatePurchaseHuawei request.
func ValidatePurchaseHuawei() *ValidatePurchaseHuaweiRequest {
	return &ValidatePurchaseHuaweiRequest{
		req: &nkapi.ValidatePurchaseHuaweiRequest{},
	}
}

// WithPurchase sets the purchase on the request.
func (req *ValidatePurchaseHuaweiRequest) WithPurchase(purchase string) *ValidatePurchaseHuaweiRequest {
	req.req.Purchase = purchase
	return req
}

// WithSignature sets the signature on the request.
func (req *ValidatePurchaseHuaweiRequest) WithSignature(signature string) *ValidatePurchaseHuaweiRequest {
	req.req.Signature = signature
	return req
}

// WithPersist sets the persist on the request.
func (req *ValidatePurchaseHuaweiRequest) WithPersist(persist bool) *ValidatePurchaseHuaweiRequest {
	req.req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidatePurchaseHuaweiRequest) Do(ctx context.Context, cl *Client) (*ValidatePurchaseResponse, error) {
	res := new(ValidatePurchaseResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/purchase/huawei", nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

/*
// ListSubscriptionsRequest is a ListSubscriptions request.
type ListSubscriptionsRequest struct {
	req *nkapi.ListSubscriptionsRequest
}

// ListSubscriptions creates a new ListSubscriptions request.
func ListSubscriptions(groupId string) *ListSubscriptionsRequest {
	return &ListSubscriptionsRequest{
		req: &nkapi.ListSubscriptionsRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListSubscriptionsRequest) WithLimit(limit int) *ListSubscriptionsRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListSubscriptionsRequest) WithCursor(cursor string) *ListSubscriptionsRequest {
	req.req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListSubscriptionsRequest) Do(ctx context.Context, cl *Client) (*ListSubscriptionsResponse, error) {
	query := url.Values{}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	res := new(ListSubscriptionsResponse)
	if err := cl.Do(ctx, "GET", "v2/iap/subscription", nil, req.req, res); err != nil {
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
	req *nkapi.ValidateSubscriptionAppleRequest
}

// ValidateSubscriptionApple creates a new ValidateSubscriptionApple request.
func ValidateSubscriptionApple() *ValidateSubscriptionAppleRequest {
	return &ValidateSubscriptionAppleRequest{
		req: &nkapi.ValidateSubscriptionAppleRequest{},
	}
}

// WithReceipt sets the receipt on the request.
func (req *ValidateSubscriptionAppleRequest) WithReceipt(receipt string) *ValidateSubscriptionAppleRequest {
	req.req.Receipt = receipt
	return req
}

// WithPersist sets the persist on the request.
func (req *ValidateSubscriptionAppleRequest) WithPersist(persist bool) *ValidateSubscriptionAppleRequest {
	req.req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidateSubscriptionAppleRequest) Do(ctx context.Context, cl *Client) (*ValidateSubscriptionResponse, error) {
	res := new(ValidateSubscriptionResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/subscription/apple", nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ValidateSubscriptionGoogleRequest is a ValidateSubscriptionGoogle request.
type ValidateSubscriptionGoogleRequest struct {
	req *nkapi.ValidateSubscriptionGoogleRequest
}

// ValidateSubscriptionGoogle creates a new ValidateSubscriptionGoogle request.
func ValidateSubscriptionGoogle() *ValidateSubscriptionGoogleRequest {
	return &ValidateSubscriptionGoogleRequest{
		req: &nkapi.ValidateSubscriptionGoogleRequest{},
	}
}

// WithReceipt sets the receipt on the request.
func (req *ValidateSubscriptionGoogleRequest) WithReceipt(receipt string) *ValidateSubscriptionGoogleRequest {
	req.req.Receipt = receipt
	return req
}

// WithPersist sets the persist on the request.
func (req *ValidateSubscriptionGoogleRequest) WithPersist(persist bool) *ValidateSubscriptionGoogleRequest {
	req.req.Persist = wrapperspb.Bool(persist)
	return req
}

// Do executes the request against the context and client.
func (req *ValidateSubscriptionGoogleRequest) Do(ctx context.Context, cl *Client) (*ValidateSubscriptionResponse, error) {
	res := new(ValidateSubscriptionResponse)
	if err := cl.Do(ctx, "POST", "v2/iap/subscription/google", nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}
*/

/*
// SubscriptionRequest is a Subscription request.
type SubscriptionRequest struct {
	req *nkapi.GetSubscriptionRequest
}

// Subscription creates a new Subscription request.
func Subscription(productId string) *SubscriptionRequest {
	return &SubscriptionRequest{
		req: &nkapi.GetSubscriptionRequest{
			ProductId: productId,
		},
	}
}

// Do executes the request against the context and client.
func (req *SubscriptionRequest) Do(ctx context.Context, cl *Client) (*SubscriptionResponse, error) {
	res := new(SubscriptionResponse)
	if err := cl.Do(ctx, "GET", "v2/iap/subscription/"+req.req.ProductId, nil, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// SubscriptionResponse is a Subscription response.
type SubscriptionResponse = nkapi.ValidatedSubscription

*/

// ListLeaderboardRecordsRequest is a ListLeaderboardRecords request.
type ListLeaderboardRecordsRequest struct {
	req *nkapi.ListLeaderboardRecordsRequest
}

// ListLeaderboardRecords creates a new ListLeaderboardRecords request.
func ListLeaderboardRecords(leaderboardId string) *ListLeaderboardRecordsRequest {
	return &ListLeaderboardRecordsRequest{
		req: &nkapi.ListLeaderboardRecordsRequest{
			LeaderboardId: leaderboardId,
			Limit:         wrapperspb.Int32(100),
		},
	}
}

// WithOwnerIds sets the ownerIds on the request.
func (req *ListLeaderboardRecordsRequest) WithOwnerIds(ownerIds ...string) *ListLeaderboardRecordsRequest {
	req.req.OwnerIds = ownerIds
	return req
}

// WithLimit sets the limit on the request.
func (req *ListLeaderboardRecordsRequest) WithLimit(limit int) *ListLeaderboardRecordsRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListLeaderboardRecordsRequest) WithCursor(cursor string) *ListLeaderboardRecordsRequest {
	req.req.Cursor = cursor
	return req
}

// WithExpiry sets the expiry on the request.
func (req *ListLeaderboardRecordsRequest) WithExpiry(expiry int) *ListLeaderboardRecordsRequest {
	req.req.Expiry = wrapperspb.Int64(int64(expiry))
	return req
}

// Do executes the request against the context and client.
func (req *ListLeaderboardRecordsRequest) Do(ctx context.Context, cl *Client) (*ListLeaderboardRecordsResponse, error) {
	query := url.Values{}
	if req.req.OwnerIds != nil {
		query.Set("ownerIds", strings.Join(req.req.OwnerIds, ","))
	}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	if req.req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.req.Expiry.Value), 10))
	}
	res := new(ListLeaderboardRecordsResponse)
	if err := cl.Do(ctx, "GET", "v2/leaderboard/"+req.req.LeaderboardId, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListLeaderboardRecordsResponse is the ListLeaderboardRecords response.
type ListLeaderboardRecordsResponse = nkapi.LeaderboardRecordList

// DeleteLeaderboardRecordRequest is a DeleteLeaderboardRecord request.
type DeleteLeaderboardRecordRequest struct {
	leaderboardId string
}

// DeleteLeaderboardRecord creates a new DeleteLeaderboardRecord request.
func DeleteLeaderboardRecord(leaderboardId string) *DeleteLeaderboardRecordRequest {
	return &DeleteLeaderboardRecordRequest{
		leaderboardId: leaderboardId,
	}
}

// Do executes the request against the context and client.
func (req *DeleteLeaderboardRecordRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/leaderboard/"+req.leaderboardId, nil, nil, nil)
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
	leaderboardId string
	req           *nkapi.WriteLeaderboardRecordRequest_LeaderboardRecordWrite
}

// WriteLeaderboardRecord creates a new WriteLeaderboardRecord request.
func WriteLeaderboardRecord(leaderboardId string) *WriteLeaderboardRecordRequest {
	return &WriteLeaderboardRecordRequest{
		leaderboardId: leaderboardId,
		req:           &nkapi.WriteLeaderboardRecordRequest_LeaderboardRecordWrite{},
	}
}

// WithScore sets the score on the request.
func (req *WriteLeaderboardRecordRequest) WithScore(score int64) *WriteLeaderboardRecordRequest {
	req.req.Score = score
	return req
}

// WithSubscore sets the subscore on the request.
func (req *WriteLeaderboardRecordRequest) WithSubscore(subscore int64) *WriteLeaderboardRecordRequest {
	req.req.Subscore = subscore
	return req
}

// WithMetadata sets the metadata on the request.
func (req *WriteLeaderboardRecordRequest) WithMetadata(metadata string) *WriteLeaderboardRecordRequest {
	req.req.Metadata = metadata
	return req
}

// WithOperator sets the operator on the request.
func (req *WriteLeaderboardRecordRequest) WithOperator(operator Operator) *WriteLeaderboardRecordRequest {
	req.req.Operator = operator
	return req
}

// Do executes the request against the context and client.
func (req *WriteLeaderboardRecordRequest) Do(ctx context.Context, cl *Client) (*WriteLeaderboardRecordResponse, error) {
	res := new(WriteLeaderboardRecordResponse)
	if err := cl.Do(ctx, "POST", "v2/leaderboard/"+req.leaderboardId, nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// WriteLeaderboardRecordResponse is the WriteLeaderboardRecord response.
type WriteLeaderboardRecordResponse = nkapi.LeaderboardRecord

// ListLeaderboardRecordsAroundOwnerRequest is a ListLeaderboardRecordsAroundOwner request.
type ListLeaderboardRecordsAroundOwnerRequest struct {
	req *nkapi.ListLeaderboardRecordsAroundOwnerRequest
}

// ListLeaderboardRecordsAroundOwner creates a new ListLeaderboardRecordsAroundOwner request.
func ListLeaderboardRecordsAroundOwner(leaderboardId, ownerId string) *ListLeaderboardRecordsAroundOwnerRequest {
	return &ListLeaderboardRecordsAroundOwnerRequest{
		req: &nkapi.ListLeaderboardRecordsAroundOwnerRequest{
			LeaderboardId: leaderboardId,
			OwnerId:       ownerId,
			Limit:         wrapperspb.UInt32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListLeaderboardRecordsAroundOwnerRequest) WithLimit(limit int) *ListLeaderboardRecordsAroundOwnerRequest {
	req.req.Limit = wrapperspb.UInt32(uint32(limit))
	return req
}

// WithExpiry sets the expiry on the request.
func (req *ListLeaderboardRecordsAroundOwnerRequest) WithExpiry(expiry int) *ListLeaderboardRecordsAroundOwnerRequest {
	req.req.Expiry = wrapperspb.Int64(int64(expiry))
	return req
}

// Do executes the request against the context and client.
func (req *ListLeaderboardRecordsAroundOwnerRequest) Do(ctx context.Context, cl *Client) (*ListLeaderboardRecordsAroundOwnerResponse, error) {
	query := url.Values{}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.req.Expiry.Value), 10))
	}
	/*
		if req.req.Cursor != "" {
			query.Set("cursor", req.req.Cursor)
		}
	*/
	res := new(ListLeaderboardRecordsAroundOwnerResponse)
	if err := cl.Do(ctx, "GET", "v2/leaderboard/"+req.req.LeaderboardId+"/owner/"+req.req.OwnerId, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListLeaderboardRecordsAroundOwnerResponse is the ListLeaderboardRecordsAroundOwner response.
type ListLeaderboardRecordsAroundOwnerResponse = nkapi.LeaderboardRecordList

// ListMatchesRequest is a ListMatches request.
type ListMatchesRequest struct {
	req *nkapi.ListMatchesRequest
}

// ListMatches creates a new ListMatches request.
func ListMatches() *ListMatchesRequest {
	return &ListMatchesRequest{
		req: &nkapi.ListMatchesRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListMatchesRequest) WithLimit(limit int) *ListMatchesRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithAuthoritative sets the authoritative on the request.
func (req *ListMatchesRequest) WithAuthoritative(authoritative bool) *ListMatchesRequest {
	req.req.Authoritative = wrapperspb.Bool(authoritative)
	return req
}

// WithLabel sets the label on the request.
func (req *ListMatchesRequest) WithLabel(label string) *ListMatchesRequest {
	req.req.Label = wrapperspb.String(label)
	return req
}

// WithMinSize sets the minSize on the request.
func (req *ListMatchesRequest) WithMinSize(minSize int) *ListMatchesRequest {
	req.req.MinSize = wrapperspb.Int32(int32(minSize))
	return req
}

// WithMaxSize sets the maxSize on the request.
func (req *ListMatchesRequest) WithMaxSize(maxSize int) *ListMatchesRequest {
	req.req.MaxSize = wrapperspb.Int32(int32(maxSize))
	return req
}

// WithQuery sets the query on the request.
func (req *ListMatchesRequest) WithQuery(query string) *ListMatchesRequest {
	req.req.Query = wrapperspb.String(query)
	return req
}

// Do executes the request against the context and client.
func (req *ListMatchesRequest) Do(ctx context.Context, cl *Client) (*ListMatchesResponse, error) {
	query := url.Values{}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.Authoritative != nil {
		query.Set("authoritative", strconv.FormatBool(req.req.Authoritative.Value))
	}
	if req.req.Label != nil {
		query.Set("label", req.req.Label.Value)
	}
	if req.req.MinSize != nil {
		query.Set("minSize", strconv.FormatInt(int64(req.req.MinSize.Value), 10))
	}
	if req.req.MaxSize != nil {
		query.Set("maxSize", strconv.FormatInt(int64(req.req.MaxSize.Value), 10))
	}
	if req.req.Query != nil {
		query.Set("query", req.req.Query.Value)
	}
	res := new(ListMatchesResponse)
	if err := cl.Do(ctx, "GET", "v2/match", query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListMatchesResponse is the ListMatches response.
type ListMatchesResponse = nkapi.MatchList

// ListNotificationsRequest is a ListNotifications request.
type ListNotificationsRequest struct {
	req *nkapi.ListNotificationsRequest
}

// ListNotifications creates a new ListNotifications request.
func ListNotifications() *ListNotificationsRequest {
	return &ListNotificationsRequest{
		req: &nkapi.ListNotificationsRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListNotificationsRequest) WithLimit(limit int) *ListNotificationsRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCacheableCursor sets the cacheableCursor on the request.
func (req *ListNotificationsRequest) WithCacheableCursor(cacheableCursor string) *ListNotificationsRequest {
	req.req.CacheableCursor = cacheableCursor
	return req
}

// Do executes the request against the context and client.
func (req *ListNotificationsRequest) Do(ctx context.Context, cl *Client) (*ListNotificationsResponse, error) {
	query := url.Values{}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.CacheableCursor != "" {
		query.Set("cacheableCursor", req.req.CacheableCursor)
	}
	res := new(ListNotificationsResponse)
	if err := cl.Do(ctx, "GET", "v2/notifications", query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListNotificationsResponse is the ListNotifications response.
type ListNotificationsResponse = nkapi.NotificationList

// DeleteNotificationsRequest is a DeleteNotifications request.
type DeleteNotificationsRequest struct {
	req *nkapi.DeleteNotificationsRequest
}

// DeleteNotifications creates a new DeleteNotifications request.
func DeleteNotifications() *DeleteNotificationsRequest {
	return &DeleteNotificationsRequest{
		req: &nkapi.DeleteNotificationsRequest{},
	}
}

// WithIds sets the Ids on the request.
func (req *DeleteNotificationsRequest) WithIds(ids ...string) *DeleteNotificationsRequest {
	req.req.Ids = ids
	return req
}

// Do executes the request against the context and client.
func (req *DeleteNotificationsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "DELETE", "v2/notification", nil, req.req, nil)
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
	return cl.Do(ctx, "POST", "v2/rpc/"+req.id, query, req.payload, v)
}

// SessionLogoutRequest is a SessionLogout request.
type SessionLogoutRequest struct {
	req *nkapi.SessionLogoutRequest
}

// SessionLogout creates a new SessionLogout request.
func SessionLogout() *SessionLogoutRequest {
	return &SessionLogoutRequest{
		req: &nkapi.SessionLogoutRequest{},
	}
}

// WithToken sets the token on the request.
func (req *SessionLogoutRequest) WithToken(token string) *SessionLogoutRequest {
	req.req.Token = token
	return req
}

// WithRefreshToken sets the refreshToken on the request.
func (req *SessionLogoutRequest) WithRefreshToken(refreshToken string) *SessionLogoutRequest {
	req.req.RefreshToken = refreshToken
	return req
}

// Do executes the request against the context and client.
func (req *SessionLogoutRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/session/logout", nil, req.req, nil)
}

// WriteStorageObject is the write storage object.
type WriteStorageObject = nkapi.WriteStorageObject

// ReadStorageObjectsRequest is a ReadStorageObjects request.
type ReadStorageObjectsRequest struct {
	req *nkapi.ReadStorageObjectsRequest
}

// ReadStorageObjects creates a new ReadStorageObjects request.
func ReadStorageObjects() *ReadStorageObjectsRequest {
	return &ReadStorageObjectsRequest{
		req: &nkapi.ReadStorageObjectsRequest{},
	}
}

// WithObjectId sets the objectId on the request.
func (req *ReadStorageObjectsRequest) WithObjectId(collection, key, userId string) *ReadStorageObjectsRequest {
	req.req.ObjectIds = append(req.req.ObjectIds, &nkapi.ReadStorageObjectId{
		Collection: collection,
		Key:        key,
		UserId:     userId,
	})
	return req
}

// Do executes the request against the context and client.
func (req *ReadStorageObjectsRequest) Do(ctx context.Context, cl *Client) (*ReadStorageObjectsResponse, error) {
	res := new(ReadStorageObjectsResponse)
	if err := cl.Do(ctx, "POST", "v2/storage", nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ReadStorageObjectsResponse is the ReadStorageObjects response.
type ReadStorageObjectsResponse = nkapi.StorageObjects

// WriteStorageObjectsRequest is a WriteStorageObjects request.
type WriteStorageObjectsRequest struct {
	req *nkapi.WriteStorageObjectsRequest
}

// WriteStorageObjects creates a new WriteStorageObjects request.
func WriteStorageObjects() *WriteStorageObjectsRequest {
	return &WriteStorageObjectsRequest{
		req: &nkapi.WriteStorageObjectsRequest{},
	}
}

// WithObject sets the object on the request.
func (req *WriteStorageObjectsRequest) WithObject(object *WriteStorageObject) *WriteStorageObjectsRequest {
	req.req.Objects = append(req.req.Objects, object)
	return req
}

// Do executes the request against the context and client.
func (req *WriteStorageObjectsRequest) Do(ctx context.Context, cl *Client) (*WriteStorageObjectsResponse, error) {
	res := new(WriteStorageObjectsResponse)
	if err := cl.Do(ctx, "PUT", "v2/storage", nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// WriteStorageObjectsResponse is the WriteStorageObjects response.
type WriteStorageObjectsResponse = nkapi.StorageObjectAcks

// DeleteStorageObjectsRequest is a DeleteStorageObjects request.
type DeleteStorageObjectsRequest struct {
	req *nkapi.DeleteStorageObjectsRequest
}

// DeleteStorageObjects creates a new DeleteStorageObjects request.
func DeleteStorageObjects() *DeleteStorageObjectsRequest {
	return &DeleteStorageObjectsRequest{
		req: &nkapi.DeleteStorageObjectsRequest{},
	}
}

// WithObjectId sets the objectId on the request.
func (req *DeleteStorageObjectsRequest) WithObjectId(collection, key, version string) *DeleteStorageObjectsRequest {
	req.req.ObjectIds = append(req.req.ObjectIds, &nkapi.DeleteStorageObjectId{
		Collection: collection,
		Key:        key,
		Version:    version,
	})
	return req
}

// Do executes the request against the context and client.
func (req *DeleteStorageObjectsRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "PUT", "v2/storage/delete", nil, req.req, nil)
}

// ListStorageObjectsRequest is a ListStorageObjects request.
type ListStorageObjectsRequest struct {
	req *nkapi.ListStorageObjectsRequest
}

// ListStorageObjects creates a new ListStorageObjects request.
func ListStorageObjects(collection string) *ListStorageObjectsRequest {
	return &ListStorageObjectsRequest{
		req: &nkapi.ListStorageObjectsRequest{
			Collection: collection,
			Limit:      wrapperspb.Int32(100),
		},
	}
}

// WithUserId sets the userId on the request.
func (req *ListStorageObjectsRequest) WithUserId(userId string) *ListStorageObjectsRequest {
	req.req.UserId = userId
	return req
}

// WithLimit sets the limit on the request.
func (req *ListStorageObjectsRequest) WithLimit(limit int) *ListStorageObjectsRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListStorageObjectsRequest) WithCursor(cursor string) *ListStorageObjectsRequest {
	req.req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListStorageObjectsRequest) Do(ctx context.Context, cl *Client) (*ListStorageObjectsResponse, error) {
	query := url.Values{}
	if req.req.UserId != "" {
		query.Set("userId", req.req.UserId)
	}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	res := new(ListStorageObjectsResponse)
	if err := cl.Do(ctx, "GET", "v2/storage/"+req.req.Collection, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListStorageObjectsResponse is the ListStorageObjects response.
type ListStorageObjectsResponse = nkapi.StorageObjectList

// ListTournamentsRequest is a ListTournaments request.
type ListTournamentsRequest struct {
	req *nkapi.ListTournamentsRequest
}

// ListTournaments creates a new ListTournaments request.
func ListTournaments() *ListTournamentsRequest {
	return &ListTournamentsRequest{
		req: &nkapi.ListTournamentsRequest{
			Limit: wrapperspb.Int32(100),
		},
	}
}

// WithCategoryStart sets the categoryStart on the request.
func (req *ListTournamentsRequest) WithCategoryStart(categoryStart uint32) *ListTournamentsRequest {
	req.req.CategoryStart = wrapperspb.UInt32(categoryStart)
	return req
}

// WithCategoryEnd sets the categoryEnd on the request.
func (req *ListTournamentsRequest) WithCategoryEnd(categoryEnd uint32) *ListTournamentsRequest {
	req.req.CategoryEnd = wrapperspb.UInt32(categoryEnd)
	return req
}

// WithLimit sets the limit on the request.
func (req *ListTournamentsRequest) WithLimit(limit int) *ListTournamentsRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithStartTime sets the startTime on the request.
func (req *ListTournamentsRequest) WithStartTime(startTime uint32) *ListTournamentsRequest {
	req.req.StartTime = wrapperspb.UInt32(startTime)
	return req
}

// WithEndTime sets the endTime on the request.
func (req *ListTournamentsRequest) WithEndTime(endTime uint32) *ListTournamentsRequest {
	req.req.EndTime = wrapperspb.UInt32(endTime)
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListTournamentsRequest) WithCursor(cursor string) *ListTournamentsRequest {
	req.req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListTournamentsRequest) Do(ctx context.Context, cl *Client) (*ListTournamentsResponse, error) {
	query := url.Values{}
	if req.req.CategoryStart != nil {
		query.Set("categoryStart", strconv.FormatUint(uint64(req.req.CategoryStart.Value), 10))
	}
	if req.req.CategoryEnd != nil {
		query.Set("categoryEnd", strconv.FormatUint(uint64(req.req.CategoryEnd.Value), 10))
	}
	if req.req.StartTime != nil {
		query.Set("startTime", strconv.FormatUint(uint64(req.req.StartTime.Value), 10))
	}
	if req.req.EndTime != nil {
		query.Set("endTime", strconv.FormatUint(uint64(req.req.EndTime.Value), 10))
	}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	res := new(ListTournamentsResponse)
	if err := cl.Do(ctx, "GET", "v2/tournament", query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListTournamentsResponse is the ListTournaments response.
type ListTournamentsResponse = nkapi.TournamentList

// ListTournamentRecordsRequest is a ListTournamentRecords request.
type ListTournamentRecordsRequest struct {
	req *nkapi.ListTournamentRecordsRequest
}

// ListTournamentRecords creates a new ListTournamentRecords request.
func ListTournamentRecords(tournamentId string) *ListTournamentRecordsRequest {
	return &ListTournamentRecordsRequest{
		req: &nkapi.ListTournamentRecordsRequest{
			TournamentId: tournamentId,
			Limit:        wrapperspb.Int32(100),
		},
	}
}

// WithOwnerIds sets the ownerIds on the request.
func (req *ListTournamentRecordsRequest) WithOwnerIds(ownerIds ...string) *ListTournamentRecordsRequest {
	req.req.OwnerIds = ownerIds
	return req
}

// WithLimit sets the limit on the request.
func (req *ListTournamentRecordsRequest) WithLimit(limit int) *ListTournamentRecordsRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithExpiry sets the expiry on the request.
func (req *ListTournamentRecordsRequest) WithExpiry(expiry int64) *ListTournamentRecordsRequest {
	req.req.Expiry = wrapperspb.Int64(expiry)
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListTournamentRecordsRequest) WithCursor(cursor string) *ListTournamentRecordsRequest {
	req.req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListTournamentRecordsRequest) Do(ctx context.Context, cl *Client) (*ListTournamentRecordsResponse, error) {
	query := url.Values{}
	if req.req.OwnerIds != nil {
		query.Set("ownerIds", strings.Join(req.req.OwnerIds, ","))
	}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	if req.req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.req.Expiry.Value), 10))
	}
	res := new(ListTournamentRecordsResponse)
	if err := cl.Do(ctx, "GET", "v2/tournament/"+req.req.TournamentId, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListTournamentRecordsResponse is the ListTournamentRecords response.
type ListTournamentRecordsResponse = nkapi.TournamentRecordList

// WriteTournamentRecordRequest is a WriteTournamentRecord request.
type WriteTournamentRecordRequest struct {
	tournamentId string
	req          *nkapi.WriteTournamentRecordRequest_TournamentRecordWrite
}

// WriteTournamentRecord creates a new WriteTournamentRecord request.
func WriteTournamentRecord(tournamentId string) *WriteTournamentRecordRequest {
	return &WriteTournamentRecordRequest{
		tournamentId: tournamentId,
		req:          &nkapi.WriteTournamentRecordRequest_TournamentRecordWrite{},
	}
}

// WithScore sets the score on the request.
func (req *WriteTournamentRecordRequest) WithScore(score int64) *WriteTournamentRecordRequest {
	req.req.Score = score
	return req
}

// WithSubscore sets the subscore on the request.
func (req *WriteTournamentRecordRequest) WithSubscore(subscore int64) *WriteTournamentRecordRequest {
	req.req.Subscore = subscore
	return req
}

// WithMetadata sets the metadata on the request.
func (req *WriteTournamentRecordRequest) WithMetadata(metadata string) *WriteTournamentRecordRequest {
	req.req.Metadata = metadata
	return req
}

// WithOperator sets the operator on the request.
func (req *WriteTournamentRecordRequest) WithOperator(operator Operator) *WriteTournamentRecordRequest {
	req.req.Operator = operator
	return req
}

// Do executes the request against the context and client.
func (req *WriteTournamentRecordRequest) Do(ctx context.Context, cl *Client) (*WriteTournamentRecordResponse, error) {
	res := new(WriteTournamentRecordResponse)
	if err := cl.Do(ctx, "POST", "v2/tournament/"+req.tournamentId, nil, req.req, res); err != nil {
		return nil, err
	}
	return res, nil
}

// WriteTournamentRecordResponse is the WriteTournamentRecord response.
type WriteTournamentRecordResponse = nkapi.LeaderboardRecord

// JoinTournamentRequest is a JoinTournament request.
type JoinTournamentRequest struct {
	req *nkapi.JoinTournamentRequest
}

// JoinTournament creates a new JoinTournament request.
func JoinTournament(tournamentId string) *JoinTournamentRequest {
	return &JoinTournamentRequest{
		req: &nkapi.JoinTournamentRequest{
			TournamentId: tournamentId,
		},
	}
}

// Do executes the request against the context and client.
func (req *JoinTournamentRequest) Do(ctx context.Context, cl *Client) error {
	return cl.Do(ctx, "POST", "v2/tournament/"+req.req.TournamentId+"/join", nil, nil, nil)
}

// ListTournamentRecordsAroundOwnerRequest is a ListTournamentRecordsAroundOwner request.
type ListTournamentRecordsAroundOwnerRequest struct {
	req *nkapi.ListTournamentRecordsAroundOwnerRequest
}

// ListTournamentRecordsAroundOwner creates a new ListTournamentRecordsAroundOwner request.
func ListTournamentRecordsAroundOwner(tournamentId, ownerId string) *ListTournamentRecordsAroundOwnerRequest {
	return &ListTournamentRecordsAroundOwnerRequest{
		req: &nkapi.ListTournamentRecordsAroundOwnerRequest{
			TournamentId: tournamentId,
			OwnerId:      ownerId,
			Limit:        wrapperspb.UInt32(100),
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListTournamentRecordsAroundOwnerRequest) WithLimit(limit int) *ListTournamentRecordsAroundOwnerRequest {
	req.req.Limit = wrapperspb.UInt32(uint32(limit))
	return req
}

// WithExpiry sets the expiry on the request.
func (req *ListTournamentRecordsAroundOwnerRequest) WithExpiry(expiry int) *ListTournamentRecordsAroundOwnerRequest {
	req.req.Expiry = wrapperspb.Int64(int64(expiry))
	return req
}

// Do executes the request against the context and client.
func (req *ListTournamentRecordsAroundOwnerRequest) Do(ctx context.Context, cl *Client) (*ListTournamentRecordsAroundOwnerResponse, error) {
	query := url.Values{}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.Expiry != nil {
		query.Set("expiry", strconv.FormatInt(int64(req.req.Expiry.Value), 10))
	}
	/*
		if req.req.Cursor != "" {
			query.Set("cursor", req.req.Cursor)
		}
	*/
	res := new(ListTournamentRecordsAroundOwnerResponse)
	if err := cl.Do(ctx, "GET", "v2/tournament/"+req.req.TournamentId+"/owner/"+req.req.OwnerId, query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListTournamentRecordsAroundOwnerResponse is the ListTournamentRecordsAroundOwner response.
type ListTournamentRecordsAroundOwnerResponse = nkapi.TournamentRecordList

// UsersRequest is a Users request.
type UsersRequest struct {
	req *nkapi.GetUsersRequest
}

// Users creates a new Users request.
func Users() *UsersRequest {
	return &UsersRequest{
		req: &nkapi.GetUsersRequest{},
	}
}

// WithIds sets the ids on the request.
func (req *UsersRequest) WithIds(ids ...string) *UsersRequest {
	req.req.Ids = ids
	return req
}

// WithUsernames sets the usernames on the request.
func (req *UsersRequest) WithUsernames(usernames ...string) *UsersRequest {
	req.req.Usernames = usernames
	return req
}

// WithFacebookIds sets the facebookIds on the request.
func (req *UsersRequest) WithFacebookIds(facebookIds ...string) *UsersRequest {
	req.req.FacebookIds = facebookIds
	return req
}

// Do executes the request against the context and client.
func (req *UsersRequest) Do(ctx context.Context, cl *Client) (*UsersResponse, error) {
	query := url.Values{}
	if len(req.req.Ids) != 0 {
		query.Set("ids", strings.Join(req.req.Ids, ","))
	}
	if len(req.req.Usernames) != 0 {
		query.Set("usernames", strings.Join(req.req.Usernames, ","))
	}
	if len(req.req.FacebookIds) != 0 {
		query.Set("facebookIds", strings.Join(req.req.FacebookIds, ","))
	}
	res := new(UsersResponse)
	if err := cl.Do(ctx, "GET", "v2/user", query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// UsersResponse is the Users response.
type UsersResponse = nkapi.Users

// ListUserGroupsRequest is a ListUserGroups request.
type ListUserGroupsRequest struct {
	req *nkapi.ListUserGroupsRequest
}

// ListUserGroups creates a new ListUserGroups request.
func ListUserGroups(userId string) *ListUserGroupsRequest {
	return &ListUserGroupsRequest{
		req: &nkapi.ListUserGroupsRequest{
			UserId: userId,
		},
	}
}

// WithLimit sets the limit on the request.
func (req *ListUserGroupsRequest) WithLimit(limit int) *ListUserGroupsRequest {
	req.req.Limit = wrapperspb.Int32(int32(limit))
	return req
}

// WithState sets the state on the request.
func (req *ListUserGroupsRequest) WithState(state int) *ListUserGroupsRequest {
	req.req.State = wrapperspb.Int32(int32(state))
	return req
}

// WithCursor sets the cursor on the request.
func (req *ListUserGroupsRequest) WithCursor(cursor string) *ListUserGroupsRequest {
	req.req.Cursor = cursor
	return req
}

// Do executes the request against the context and client.
func (req *ListUserGroupsRequest) Do(ctx context.Context, cl *Client) (*ListUserGroupsResponse, error) {
	query := url.Values{}
	if req.req.Limit != nil {
		query.Set("limit", strconv.FormatInt(int64(req.req.Limit.Value), 10))
	}
	if req.req.State != nil {
		query.Set("state", strconv.FormatInt(int64(req.req.State.Value), 10))
	}
	if req.req.Cursor != "" {
		query.Set("cursor", req.req.Cursor)
	}
	res := new(ListUserGroupsResponse)
	if err := cl.Do(ctx, "GET", "v2/user/"+req.req.UserId+"/group", query, nil, res); err != nil {
		return nil, err
	}
	return res, nil
}

// ListUserGroupsResponse is the ListUserGroups response.
type ListUserGroupsResponse = nkapi.UserGroupList
