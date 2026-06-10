package token_exchange

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
)

type PasswordRequest struct {
	Principal client.Principal
	Scopes    scope.Scopes
	Username  string
	Password  string
	State     string
}
