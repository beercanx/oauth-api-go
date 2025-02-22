package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/grant"
	"github.com/gin-gonic/gin"
)

func validateRequest(principal client.Principal, context *gin.Context) (Valid, *Invalid) {

	switch grantType := context.PostForm("grant_type"); grantType {

	case "":
		return nil, &Invalid{Error: InvalidRequest, Description: "missing parameter: grant_type"}

	case string(grant.Password):
		return validatePasswordRequest(principal, context)

	default:
		return nil, &Invalid{Error: UnsupportedGrantType, Description: "unsupported: " + grantType}
	}
}
