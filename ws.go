//go:build !js

package nakama

import (
	"net/http"

	"nhooyr.io/websocket"
)

// buildWsOptions builds the websocket dial options.
func buildWsOptions(httpClient *http.Client) *websocket.DialOptions {
	return &websocket.DialOptions{
		HTTPClient: httpClient,
	}
}
