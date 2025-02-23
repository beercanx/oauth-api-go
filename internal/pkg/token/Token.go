package token

import (
	"baconi.co.uk/oauth/internal/pkg/authentication"
	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/google/uuid"
)

type Token interface {
	GetValue() uuid.UUID
	GetUsername() authentication.AuthenticatedUsername
	GetClientId() client.Id
	HasExpired() bool
	IsBefore() bool
}
