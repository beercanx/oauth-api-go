package token

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Route(c *gin.Context) {

	// TODO - Confidential and Public client based authentication

	// TODO - Validate contentType
	// TODO - Validate formPost

	_, invalid := validateRequest(c)
	if invalid != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: string(invalid.Error), Description: invalid.Description})
	}

	// TODO - Route requests to specific Grant Type handler
}

type ErrorType string

const (
	InvalidRequest       ErrorType = "invalid_request"
	InvalidClient        ErrorType = "invalid_client"
	InvalidGrant         ErrorType = "invalid_grant"
	InvalidScope         ErrorType = "invalid_scope"
	UnauthorizedClient   ErrorType = "unauthorized_client"
	UnsupportedGrantType ErrorType = "unsupported_grant_type"
)

type Valid interface { // TODO - Improve or change this approach
	//isValid() bool
}

type Invalid struct {
	Error       ErrorType
	Description string
}

func validateRequest(c *gin.Context) (Valid, *Invalid) {

	grantType := c.PostForm("grant_type")

	if grantType == "" {
		return nil, &Invalid{Error: InvalidRequest, Description: "missing parameter: grant_type"}
	}

	switch grantType {
	case "password":
		return validatePasswordRequest(c)
	default:
		return nil, &Invalid{Error: UnsupportedGrantType, Description: "unsupported_grant_type: " + grantType}
	}
}

type PasswordRequest struct {
	// TODO client string
	scopes   []string
	username string
	password string
}

func (r PasswordRequest) isValid() bool {
	return true
}

func validatePasswordRequest(c *gin.Context) (*PasswordRequest, *Invalid) {

	// TODO - Check client is a confidential client
	// TODO - Check client is allowed to perform grant type

	username, usernameOk := c.GetPostForm("username")
	password, passwordOk := c.GetPostForm("password")
	scope, scopeOk := c.GetPostForm("scope")
	rawScopes := strings.Split(scope, " ")
	scopes := rawScopes // TODO - Perform a check to filter down to only valid scopes

	switch {

	case !usernameOk:
		return nil, &Invalid{Error: InvalidRequest, Description: "missing parameter: username"}
	case username == "":
		return nil, &Invalid{Error: InvalidRequest, Description: "invalid parameter: username"}

	case !passwordOk:
		return nil, &Invalid{Error: InvalidRequest, Description: "missing parameter: password"}

	case scopeOk && len(scopes) == 0:
		return nil, &Invalid{Error: InvalidScope, Description: "invalid parameter: scope"}

	case len(rawScopes) != len(scopes):
		return nil, &Invalid{Error: InvalidScope, Description: "invalid parameter: scope"}

		// TODO - Check that client can be issued all these scopes
	}

	return &PasswordRequest{scopes: scopes, username: username, password: password}, nil
}
