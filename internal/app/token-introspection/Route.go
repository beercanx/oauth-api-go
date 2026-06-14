package token_introspection

import (
	"log"
	"net/http"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/server"
	"github.com/gin-gonic/gin"
)

func Route(engine *gin.Engine, clientAuthenticator client.Authenticator, introspector Introspector) {

	engine.POST("/introspect",
		client.AuthenticateConfidentialClient(clientAuthenticator),
		client.RequireConfidentialClientAuthentication,

		server.RequireUrlEncodedForm,

		func(context *gin.Context) {

			validated, validationError := validateRequest(context)

			if validationError != nil {
				switch validationError.ErrorType {
				case UnauthorizedClient:
					context.AbortWithStatusJSON(http.StatusForbidden, validationError)
				default:
					context.AbortWithStatusJSON(http.StatusBadRequest, validationError)
				}
				return
			}

			introspected, introspectionError := introspector.introspect(validated)

			if introspectionError != nil {
				log.Println("[ERROR] Unexpected introspection error:", introspectionError)
				context.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			context.JSON(http.StatusOK, introspected)
		},
	)
}
