package token_introspection

import (
	"log/slog"
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
				case InvalidRequest:
					fallthrough
				default:
					context.AbortWithStatusJSON(http.StatusBadRequest, validationError)
				}
				return
			}

			introspected, introspectionError := introspector.introspect(validated)

			if introspectionError != nil {
				slog.Error("Unexpected error during introspection", slog.Any("error", introspectionError))
				context.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			context.JSON(http.StatusOK, introspected)
		},
	)
}
