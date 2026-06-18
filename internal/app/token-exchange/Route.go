package token_exchange

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/server"
	"github.com/gin-gonic/gin"
)

func Route(
	engine *gin.Engine,
	clientAuthenticator client.Authenticator,
	passwordGrant Grant[PasswordRequest],
) {

	engine.POST("/token",

		client.AuthenticateConfidentialClient(clientAuthenticator),
		client.AuthenticatePublicClient(clientAuthenticator),
		client.RequireClientAuthentication,

		server.RequireUrlEncodedForm,

		func(context *gin.Context) {

			request, invalid := validateRequest(context)
			if invalid != nil {
				context.JSON(http.StatusBadRequest, Failed(*invalid))
				return
			}

			var result Success
			var exchangeError error

			switch valid := request.(type) {
			// TODO - Add support for other grant types
			case *PasswordRequest:
				result, exchangeError = passwordGrant.Exchange(valid)
			default:
				exchangeError = Failed{Err: UnsupportedGrantType, Description: fmt.Sprintf("unsupported grant type: %T", valid)}
			}

			var failed Failed
			switch {
			case exchangeError != nil && errors.As(exchangeError, &failed):
				context.JSON(http.StatusBadRequest, failed)
			case exchangeError != nil:
				log.Printf("[ERROR][token_exchange.Route] Some kind of error bubbled up... %T: %v\n", exchangeError, exchangeError)
				context.AbortWithStatus(http.StatusInternalServerError)
			default:
				context.JSON(http.StatusOK, result)
			}
		},
	)
}
