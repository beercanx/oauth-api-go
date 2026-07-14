package token_exchange

import (
	"fmt"

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

var _ fmt.Stringer = (*PasswordRequest)(nil)

func (r PasswordRequest) String() string {
	return fmt.Sprintf("PasswordRequest{Principal: %s, Scopes: %s, Username: %s, State: %s}", r.Principal, r.Scopes, r.Username, r.State)
}
