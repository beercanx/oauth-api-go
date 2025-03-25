package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct { // TODO - Consider locking down access to getters.
	Value     uuid.UUID
	Username  user.AuthenticatedUsername
	ClientId  client.Id
	Scopes    []scope.Scope
	IssuedAt  time.Time
	ExpiresAt time.Time
	NotBefore time.Time
}

// assert RefreshToken implements Token
var _ Token = (*RefreshToken)(nil)

func (token RefreshToken) GetValue() uuid.UUID {
	return token.Value
}

func (token RefreshToken) GetUsername() user.AuthenticatedUsername {
	return token.Username
}

func (token RefreshToken) GetClientId() client.Id {
	return token.ClientId
}

func (token RefreshToken) HasExpired() bool {
	return time.Now().After(token.ExpiresAt)
}

func (token RefreshToken) IsBefore() bool {
	return time.Now().Before(token.NotBefore)
}
