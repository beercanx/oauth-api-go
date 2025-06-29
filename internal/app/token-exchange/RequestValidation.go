package token_exchange

import (
	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"github.com/gin-gonic/gin"
)

func validateRequest(scopeService *scope.Service, context *gin.Context) (Valid, *Invalid) {

	switch grantType := context.PostForm("grant_type"); grantType {

	case "":
		return nil, &Invalid{Err: InvalidRequest, Description: "missing parameter: grant_type"}

	case string(grant.Password):
		return validatePasswordRequest(scopeService, context)

	default:
		return nil, &Invalid{Err: UnsupportedGrantType, Description: "unsupported: " + grantType}
	}
}
