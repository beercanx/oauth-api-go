package token

import (
	"baconi.co.uk/oauth/internal/pkg/authentication"
	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/google/uuid"
	"time"
)

type RefreshToken struct {
	Value     uuid.UUID
	Username  authentication.AuthenticatedUsername
	ClientId  client.Id
	Scopes    []string
	IssuedAt  time.Time
	ExpiresAt time.Time
	NotBefore time.Time
}

func (token RefreshToken) GetValue() uuid.UUID {
	return token.Value
}

func (token RefreshToken) GetUsername() authentication.AuthenticatedUsername {
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
