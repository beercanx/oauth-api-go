package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
	"time"
)

type AccessToken struct {
	value     uuid.UUID
	username  user.AuthenticatedUsername
	clientId  client.Id
	scopes    []scope.Scope
	issuedAt  time.Time
	expiresAt time.Time
	notBefore time.Time
}

// assert AccessToken implements Token
var _ Token = (*AccessToken)(nil)

func (token AccessToken) GetValue() uuid.UUID {
	return token.value
}

func (token AccessToken) GetUsername() user.AuthenticatedUsername {
	return token.username
}

func (token AccessToken) GetClientId() client.Id {
	return token.clientId
}

func (token AccessToken) GetScopes() []scope.Scope {
	return token.scopes
}

func (token AccessToken) GetIssuedAt() time.Time {
	return token.issuedAt
}

func (token AccessToken) GetExpiresAt() time.Time {
	return token.expiresAt
}

func (token AccessToken) GetNotBefore() time.Time {
	return token.notBefore
}
