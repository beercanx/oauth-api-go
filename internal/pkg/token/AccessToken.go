package token

import (
	"baconi.co.uk/oauth/internal/pkg/authentication"
	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/google/uuid"
	"time"
)

type AccessToken struct {
	Value     uuid.UUID
	Username  authentication.AuthenticatedUsername
	ClientId  client.Id
	Scopes    []string
	IssuedAt  time.Time
	ExpiresAt time.Time
	NotBefore time.Time
}

func (token AccessToken) GetValue() uuid.UUID {
	return token.Value
}

func (token AccessToken) GetUsername() authentication.AuthenticatedUsername {
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
