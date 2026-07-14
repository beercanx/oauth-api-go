package token

import (
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type AccessToken struct {
	ID        uuid.UUID                  `db:"id"`
	Username  user.AuthenticatedUsername `db:"username"`
	ClientID  client.Id                  `db:"client_id"`
	Scopes    scope.Scopes               `db:"scopes"`
	IssuedAt  time.Time                  `db:"issued_at"`
	ExpiresAt time.Time                  `db:"expires_at"`
	NotBefore time.Time                  `db:"not_before"`
}

func (token AccessToken) HasExpired() bool {
	return time.Now().After(token.ExpiresAt)
}

func (token AccessToken) IsBefore() bool {
	return time.Now().Before(token.NotBefore)
}
