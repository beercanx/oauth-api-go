package token

import (
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type AccessToken struct {
	Value     uuid.UUID
	Username  user.AuthenticatedUsername
	ClientId  client.Id
	Scopes    scope.Scopes
	IssuedAt  time.Time
	ExpiresAt time.Time
	NotBefore time.Time
}

// assert AccessToken implements Token
var _ Token = (*AccessToken)(nil)

func (token AccessToken) GetValue() uuid.UUID {
	return token.Value
}

func (token AccessToken) GetUsername() user.AuthenticatedUsername {
	return token.Username
}

func (token AccessToken) GetClientId() client.Id {
	return token.ClientId
}

func (token AccessToken) GetScopes() scope.Scopes {
	return token.Scopes
}

func (token AccessToken) GetIssuedAt() time.Time {
	return token.IssuedAt
}

func (token AccessToken) GetExpiresAt() time.Time {
	return token.ExpiresAt
}

func (token AccessToken) GetNotBefore() time.Time {
	return token.NotBefore
}
