package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
)

type Token interface {
	GetValue() uuid.UUID
	GetUsername() user.AuthenticatedUsername
	GetClientId() client.Id
	HasExpired() bool
	IsBefore() bool
}
