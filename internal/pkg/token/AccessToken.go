package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
	"time"
)

type AccessToken struct { // TODO - Consider locking down access to getters.
	Value     uuid.UUID
	Username  user.AuthenticatedUsername
	ClientId  client.Id
	Scopes    []scope.Scope
	IssuedAt  time.Time
	ExpiresAt time.Time
	NotBefore time.Time
}

// assert AccessToken implements Token
var _ Token = &AccessToken{}

func (token AccessToken) GetValue() uuid.UUID {
	return token.Value
}

func (token AccessToken) GetUsername() user.AuthenticatedUsername {
	return token.Username
}

func (token AccessToken) GetClientId() client.Id {
	return token.ClientId
}

func (token AccessToken) HasExpired() bool {
	return time.Now().After(token.ExpiresAt)
}

func (token AccessToken) IsBefore() bool {
	return time.Now().Before(token.NotBefore)
}
