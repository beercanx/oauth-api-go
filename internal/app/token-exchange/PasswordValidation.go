package token_exchange

import (
	"fmt"
	"strings"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"github.com/gin-gonic/gin"
)

func validatePasswordRequest(scopeService *scope.Service, context *gin.Context) (*PasswordRequest, *Invalid) {

	username, usernameOk := context.GetPostForm("username")
	password, passwordOk := context.GetPostForm("password")
	scopeP, scopeOk := context.GetPostForm("scope")
	rawScopes := strings.Split(scopeP, " ")
	scopes := scopeService.Validate(rawScopes)
	state, stateOk := context.GetPostForm("state")

	principal := context.MustGet(client.AuthClientKey).(client.Principal)

	switch {
	case !principal.IsConfidential(), !principal.CanBeGranted(grant.Password):
		return nil, &Invalid{Err: UnauthorizedClient, Description: fmt.Sprintf("%s is not authorized to: %s", principal.Id, grant.Password)}

	case !usernameOk:
		return nil, &Invalid{Err: InvalidRequest, Description: "missing parameter: username"}
	case username == "":
		return nil, &Invalid{Err: InvalidRequest, Description: "invalid parameter: username"}

	// As long as the password field is present, we should not restrict what it contains.
	case !passwordOk:
		return nil, &Invalid{Err: InvalidRequest, Description: "missing parameter: password"}

	// The requested scope is invalid, unknown, or malformed.
	case scopeOk && len(scopes.Value) == 0:
		return nil, &Invalid{Err: InvalidScope, Description: "invalid parameter: scope"}
	case len(rawScopes) != len(scopes.Value):
		return nil, &Invalid{Err: InvalidScope, Description: "invalid parameter: scope"}
	case !principal.CanBeIssued(scopes.Value):
		return nil, &Invalid{Err: InvalidScope, Description: "invalid parameter: scope"}

	// TODO - Enforce unique scopes requested

	// If state is provided it cannot be a zero string.
	case stateOk && state == "":
		return nil, &Invalid{Err: InvalidScope, Description: "invalid parameter: state"}

	default:
		return &PasswordRequest{Principal: principal, Scopes: scopes, Username: username, Password: password, State: state}, nil
	}
}
