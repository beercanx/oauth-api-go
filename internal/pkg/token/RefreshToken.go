package token

import (
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type RefreshToken struct {
	value     uuid.UUID                  `db:"id"`
	username  user.AuthenticatedUsername `db:"username"`
	clientId  client.Id                  `db:"client_id"`
	scopes    scope.Scopes               `db:"scopes"`
	issuedAt  time.Time                  `db:"issued_at"`
	expiresAt time.Time                  `db:"expires_at"`
	notBefore time.Time                  `db:"not_before"`
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

func (token RefreshToken) GetScopes() scope.Scopes {
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
