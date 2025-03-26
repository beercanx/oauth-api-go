package token_exchange

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"github.com/gin-gonic/gin"
	"strings"
)

func validatePasswordRequest(scopeService *scope.Service, principal client.Principal, context *gin.Context) (*PasswordRequest, *Invalid) {

	username, usernameOk := context.GetPostForm("username")
	password, passwordOk := context.GetPostForm("password")
	scopeP, scopeOk := context.GetPostForm("scope")
	rawScopes := strings.Split(scopeP, " ")
	scopes := scopeService.Validate(rawScopes)
	state, stateOk := context.GetPostForm("state")

	switch {
	case !principal.IsConfidential():
		return nil, &Invalid{Error: UnauthorizedClient, Description: string("not authorized to: " + grant.Password)}
	case !principal.Can(grant.Password):
		return nil, &Invalid{Error: UnauthorizedClient, Description: string("not authorized to: " + grant.Password)}

	case !usernameOk:
		return nil, &Invalid{Error: InvalidRequest, Description: "missing parameter: username"}
	case username == "":
		return nil, &Invalid{Error: InvalidRequest, Description: "invalid parameter: username"}

	// As long as the password field is present we should not restrict what it contains.
	case !passwordOk:
		return nil, &Invalid{Error: InvalidRequest, Description: "missing parameter: password"}

	// The requested scope is invalid, unknown, or malformed.
	case scopeOk && len(scopes) == 0:
		return nil, &Invalid{Error: InvalidScope, Description: "invalid parameter: scope"}
	case len(rawScopes) != len(scopes):
		return nil, &Invalid{Error: InvalidScope, Description: "invalid parameter: scope"}
	case !principal.CanBeIssued(scopes):
		return nil, &Invalid{Error: InvalidScope, Description: "invalid parameter: scope"}

	// If state is provided it cannot be a zero string.
	case stateOk && state == "":
		return nil, &Invalid{Error: InvalidScope, Description: "invalid parameter: state"}

	default:
		return &PasswordRequest{Principal: principal, Scopes: scopes, Username: username, Password: password, State: state}, nil
	}
}
