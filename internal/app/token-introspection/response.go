package token_introspection

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/token"
	"baconi.co.uk/oauth/internal/pkg/user"
)

type response struct {
	Active         bool                       `json:"active"`
	Scope          scope.Scopes               `json:"scope,omitzero"`
	ClientId       client.Id                  `json:"client_id,omitzero"`
	Username       user.AuthenticatedUsername `json:"username,omitzero"`
	TokenType      token.Type                 `json:"token_type,omitzero"`
	ExpirationTime int64                      `json:"expiration_time,omitzero"`
	IssuedAt       int64                      `json:"issued_at,omitzero"`
	NotBefore      int64                      `json:"not_before,omitzero"`
	Subject        user.AuthenticatedUsername `json:"sub,omitzero"`
}
