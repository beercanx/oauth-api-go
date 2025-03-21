package token_exchange

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/grant"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
)

func Route(
	scopeService *scope.Service,
	passwordGrant Grant[PasswordRequest],
) gin.HandlerFunc {

	// TODO - Open to do other setup, no stored state of course...

	return func(context *gin.Context) {

		// TODO - Confidential and Public client based authentication
		principal := client.Principal{
			Id:                client.Id{Value: "aardvark"},
			Type:              client.Confidential,
			AllowedScopes:     []scope.Scope{scope.Basic},
			AllowedGrantTypes: []grant.Type{grant.Password},
		}

		// TODO - Look at replacing with some middleware, this endpoint only supports URL encoded form posts, others may also.
		if contentType := context.ContentType(); contentType != "application/x-www-form-urlencoded" {
			context.JSON(http.StatusUnsupportedMediaType, Failed{Error: InvalidRequest, Description: "Unsupported Media Type: " + contentType})
			return
		}

		request, invalid := validateRequest(scopeService, principal, context)
		if invalid != nil {
			// TODO - Can we combine Valid and Invalid to one return type and add a type case in the switch below?
			context.JSON(http.StatusBadRequest, Failed(*invalid))
			return
		}

		var result Response

		switch valid := request.(type) {
		// TODO - Add support for other grant types
		case *PasswordRequest:
			result = passwordGrant.Exchange(*valid)
		default:
			result = &Failed{Error: UnsupportedGrantType, Description: reflect.TypeOf(valid).Name()}
		}

		switch response := result.(type) {
		case Failed:
			context.JSON(http.StatusBadRequest, response)
		case Success:
			context.JSON(http.StatusOK, response)
		default:
			log.Println("[ERROR] No success or failure to reply with...", reflect.TypeOf(response).Name())
			context.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
