package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
)

type Issuer[T Token] interface {
	Issue(username user.AuthenticatedUsername, clientId client.Id, scopes []scope.Scope) (T, error)
}
