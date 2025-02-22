package token

import "baconi.co.uk/oauth/internal/pkg/client"

type PasswordRequest struct { // TODO - Remove json tags its only for debug during build
	Principal client.Principal `json:"principal"` // TODO - Work out how we can change this to only return say client.Id
	Scopes    []string         `json:"scopes,omitempty"`
	Username  string           `json:"username,omitempty"`
	Password  string           `json:"password,omitempty"`
}

func (request PasswordRequest) getPrincipal() client.Principal {
	return request.Principal
}
