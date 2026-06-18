package token

import (
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID                  `db:"id"`
	Username  user.AuthenticatedUsername `db:"username"`
	ClientID  client.Id                  `db:"client_id"`
	Scopes    scope.Scopes               `db:"scopes"`
	IssuedAt  time.Time                  `db:"issued_at"`
	ExpiresAt time.Time                  `db:"expires_at"`
	NotBefore time.Time                  `db:"not_before"`
	// TODO - Decide if we need to store a database link to the access token issued with it,
	// 				so that once a refresh is performed, both this and its access token are revoked when issuing the new tokens.
}

func (token RefreshToken) HasExpired() bool {
	return time.Now().After(token.ExpiresAt)
}

func (token RefreshToken) IsBefore() bool {
	return time.Now().Before(token.NotBefore)
}
