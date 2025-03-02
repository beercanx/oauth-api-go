package token_exchange

import "baconi.co.uk/oauth/internal/pkg/client"

type Valid interface {
	getPrincipal() client.Principal
}

type Invalid struct {
	Error       ErrorType
	Description string
}
