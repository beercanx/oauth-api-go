package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/grant"
	"github.com/gin-gonic/gin"
	"strings"
)

func validatePasswordRequest(principal client.Principal, context *gin.Context) (*PasswordRequest, *Invalid) {

	username, usernameOk := context.GetPostForm("username")
	password, passwordOk := context.GetPostForm("password")
	scope, scopeOk := context.GetPostForm("scope")
	rawScopes := strings.Split(scope, " ")
	scopes := rawScopes // TODO - Perform a check to filter down to only valid scopes

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

	default:
		return &PasswordRequest{Principal: principal, Scopes: scopes, Username: username, Password: password}, nil
	}
}
