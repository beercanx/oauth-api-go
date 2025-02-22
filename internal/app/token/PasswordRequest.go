package token

import "baconi.co.uk/oauth/internal/pkg/client"

type PasswordRequest struct {
	Principal client.Principal
	Scopes    []string
	Username  string
	Password  string
}

func (request PasswordRequest) getPrincipal() client.Principal {
	return request.Principal
}
