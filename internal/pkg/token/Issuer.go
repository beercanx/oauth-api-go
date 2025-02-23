package token

import (
	"baconi.co.uk/oauth/internal/pkg/authentication"
	"baconi.co.uk/oauth/internal/pkg/client"
)

type Issuer[T Token] interface {
	Issue(username authentication.AuthenticatedUsername, clientId client.Id, scopes []string) (*T, error)
}
