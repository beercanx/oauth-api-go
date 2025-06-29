package token_introspection

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/server"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Route(engine *gin.Engine, clientAuthenticator client.Authenticator, introspector Introspector) {

	engine.POST("/introspect",
		client.AuthenticateConfidentialClient(clientAuthenticator),
		client.RequireConfidentialClientAuthentication,

		server.RequireUrlEncodedForm,

		func(context *gin.Context) {

			validated, validationError := validateRequest(context)

			if validationError != nil {

				var failedValidation invalid
				if errors.As(validationError, &failedValidation) {

					switch failedValidation.ErrorType {

					case InvalidRequest:
						context.AbortWithStatusJSON(http.StatusBadRequest, failedValidation)
						return

					case UnauthorizedClient:
						context.AbortWithStatusJSON(http.StatusForbidden, failedValidation)
						return
					}
				}

				log.Println("[ERROR] Unexpected introspection validation error:", validationError)
				context.AbortWithStatus(http.StatusInternalServerError)
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
