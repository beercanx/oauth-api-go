package token_exchange

import (
	"baconi.co.uk/oauth/internal/pkg/grant"
	"github.com/gin-gonic/gin"
)

func validateRequest(context *gin.Context) (any, *Invalid) {

	switch grantType := context.PostForm("grant_type"); grantType {

	case "":
		return nil, &Invalid{Err: InvalidRequest, Description: "missing parameter: grant_type"}

	case string(grant.Password):
		return validatePasswordRequest(context)

	default:
		return nil, &Invalid{Err: UnsupportedGrantType, Description: "unsupported: " + grantType}
	}
}
