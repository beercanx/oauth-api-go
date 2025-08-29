package token

import (
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type Token interface {
	GetValue() uuid.UUID
	GetUsername() user.AuthenticatedUsername
	GetClientId() client.Id
	GetScopes() scope.Scopes
	GetIssuedAt() time.Time
	GetExpiresAt() time.Time
	GetNotBefore() time.Time
}

func HasExpired(token Token) bool {
	return time.Now().After(token.GetExpiresAt())
}

func IsBefore(token Token) bool {
	return time.Now().Before(token.GetNotBefore())
}
