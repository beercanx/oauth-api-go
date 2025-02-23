package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/grant"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"reflect"
)

func Route(
	passwordGrant *PasswordGrant,
) gin.HandlerFunc {

	// TODO - Open to do other setup, no stored state of course...

	return func(context *gin.Context) {

		// TODO - Confidential and Public client based authentication
		principal := client.Principal{
			Id:                client.Id{Value: "aardvark"},
			Type:              client.Confidential,
			AllowedScopes:     []string{"basic"},
			AllowedGrantTypes: []grant.Type{grant.Password},
		}

		// TODO - Look at replacing with some middleware, this endpoint only supports URL encoded form posts, others may also.
		if contentType := context.ContentType(); contentType != "application/x-www-form-urlencoded" {
			context.JSON(http.StatusUnsupportedMediaType, Failed{Error: InvalidRequest, Description: "Unsupported Media Type: " + contentType})
			return
		}

		request, invalid := validateRequest(principal, context)
		if invalid != nil {
			// TODO - Can we combine Valid and Invalid to one return type and add a type case in the switch below?
			context.JSON(http.StatusBadRequest, Failed(*invalid))
			return
		}

		var success *Success
		var failed *Failed
		var err error

		switch valid := request.(type) {
		// TODO - Add support for other grant types
		case *PasswordRequest:
			success, failed, err = passwordGrant.Exchange(*valid)
		default:
			failed = &Failed{Error: UnsupportedGrantType, Description: reflect.TypeOf(valid).Name()}
		}

		switch {
		case err != nil:
			log.Println("[ERROR] Failed to exchange the grant:", err.Error())
			context.AbortWithStatus(http.StatusInternalServerError)
		case failed != nil:
			context.JSON(http.StatusBadRequest, failed)
		case success != nil:
			context.JSON(http.StatusOK, success)
		default:
			log.Println("[ERROR] No success or failure to reply with...")
			context.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
