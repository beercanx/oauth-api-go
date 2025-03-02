package token_exchange

import (
	"baconi.co.uk/oauth/internal/pkg/token"
	"github.com/google/uuid"
)

type Response interface {
	success() bool
}

// Success https://www.rfc-editor.org/rfc/rfc6749#section-5.1
type Success struct {

	// AccessToken The access token issued by the authorization server.
	AccessToken uuid.UUID `json:"access_token"`

	// TokenType The type of the token issued as described in
	// https://www.rfc-editor.org/rfc/rfc6749#section-7.1
	TokenType token.Type `json:"token_type"`

	// ExpiresIn The lifetime in seconds of the access token. For example, the value
	// "3600" denotes that the access token will expire in one hour from the time the
	// response was generated. If omitted, the authorization server SHOULD provide
	// the expiration time via other means or document the default value.
	ExpiresIn int64 `json:"expires_in"`

	// RefreshToken OPTIONAL. The refresh token, which can be used to obtain new
	// access tokens using the same authorization grant as described in
	// https://www.rfc-editor.org/rfc/rfc6749#section-6
	RefreshToken *uuid.UUID `json:"refresh_token,omitempty"`

	// Scope OPTIONAL if identical to the scope requested by the client; otherwise,
	// REQUIRED. The scope of the access token as described by
	// https://www.rfc-editor.org/rfc/rfc6749#section-3.3
	Scope *string `json:"scope,omitempty"`

	// State REQUIRED if the "state" parameter was present in the client
	// authorization request. The exact value received from the client.
	State *string `json:"state,omitempty"`
}

func (s Success) success() bool {
	return true
}

// Failed https://www.rfc-editor.org/rfc/rfc6749#section-5.2
type Failed struct {

	// Error A single ASCII error code from the defined list.
	Error ErrorType `json:"error"`

	// Description Human-readable ASCII text providing additional information, used
	// to assist the client developer in understanding the error that occurred.
	Description string `json:"error_description"`
}

func (f Failed) success() bool {
	return false
}
