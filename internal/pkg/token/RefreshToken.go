package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	value     uuid.UUID
	username  user.AuthenticatedUsername
	clientId  client.Id
	scopes    []scope.Scope
	issuedAt  time.Time
	expiresAt time.Time
	notBefore time.Time
}

// assert RefreshToken implements Token
var _ Token = (*RefreshToken)(nil)

func (token RefreshToken) GetValue() uuid.UUID {
	return token.value
}

func (token RefreshToken) GetUsername() user.AuthenticatedUsername {
	return token.username
}

func (token RefreshToken) GetClientId() client.Id {
	return token.clientId
}

func (token RefreshToken) GetScopes() []scope.Scope {
	return token.scopes
}

func (token RefreshToken) GetIssuedAt() time.Time {
	return token.issuedAt
}

func (token RefreshToken) GetExpiresAt() time.Time {
	return token.expiresAt
}

func (token RefreshToken) GetNotBefore() time.Time {
	return token.notBefore
}
