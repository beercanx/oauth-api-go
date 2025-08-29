package token_exchange

import (
	"errors"
	"log"
	"net/http"
	"reflect"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/server"
	"github.com/gin-gonic/gin"
)

func Route(
	engine *gin.Engine,
	clientAuthenticator client.Authenticator,
	scopeService *scope.Service,
	passwordGrant Grant[PasswordRequest],
) {

	engine.POST("/token",

		client.AuthenticateConfidentialClient(clientAuthenticator),
		client.AuthenticatePublicClient(clientAuthenticator),
		client.RequireClientAuthentication,

		server.RequireUrlEncodedForm,

		func(context *gin.Context) {

			request, invalid := validateRequest(scopeService, context)
			if invalid != nil {
				context.JSON(http.StatusBadRequest, Failed(*invalid))
				return
			}

			var result Success
			var err error = nil

			switch valid := request.(type) {
			// TODO - Add support for other grant types
			case *PasswordRequest:
				result, err = passwordGrant.Exchange(*valid)
			default:
				err = Failed{Err: UnsupportedGrantType, Description: reflect.TypeOf(valid).Name()}
			}

			var failed Failed
			switch {
			case err != nil && errors.As(err, &failed):
				context.JSON(http.StatusBadRequest, failed)
			case err != nil:
				log.Println("[ERROR] Some kind of error bubbled up...", reflect.TypeOf(err).Name(), err)
				context.AbortWithStatus(http.StatusInternalServerError)
			default:
				context.JSON(http.StatusOK, result)
			}
		},
	)
}
