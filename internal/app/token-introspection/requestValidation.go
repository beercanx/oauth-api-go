package token_introspection

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func validateRequest(context *gin.Context) (request, error) {

	principal := context.MustGet(client.AuthClientConfidentialKey).(client.Principal)
	token, tokenOk := context.GetPostForm("token")
	tokenUuid, tokenUuidError := uuid.Parse(token)

	switch {

	case !principal.IsConfidential():
		return request{}, invalid{UnauthorizedClient, "client is not allowed to introspect"}
	case !principal.CanPerformAction(client.Introspect):
		return request{}, invalid{UnauthorizedClient, "client is not allowed to introspect"}

	case !tokenOk:
		return request{}, invalid{InvalidRequest, "missing parameter: token"}
	case tokenUuidError != nil:
		return request{}, invalid{InvalidRequest, "invalid parameter: token"}

	default:
		return request{principal, tokenUuid}, nil
	}
}
