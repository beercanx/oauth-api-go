package token_exchange

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
)

type PasswordRequest struct {
	Principal client.Principal
	Scopes    []scope.Scope
	Username  string
	Password  string
}

func (request PasswordRequest) getPrincipal() client.Principal {
	return request.Principal
}
