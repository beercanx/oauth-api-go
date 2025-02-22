package token

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/grant"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

func Route(context *gin.Context) {

	// TODO - Confidential and Public client based authentication
	principal := client.Principal{
		Id:                client.Id{Value: "aardvark"},
		Type:              client.Confidential,
		AllowedScopes:     []string{"basic"},
		AllowedGrantTypes: []grant.Type{grant.Password},
	}

	// TODO - Look at replacing with some middleware, this endpoint only supports URL encoded form posts, others may also.
	if contentType := context.ContentType(); contentType != "application/x-www-form-urlencoded" {
		context.JSON(http.StatusUnsupportedMediaType, ErrorResponse{Error: InvalidRequest, Description: "Unsupported Media Type: " + contentType})
		return
	}

	request, invalid := validateRequest(principal, context)
	if invalid != nil {
		// TODO - Can we combine Valid and Invalid to one return type and add a type case in the switch below?
		context.JSON(http.StatusBadRequest, ErrorResponse(*invalid))
		return
	}

	// TODO - Route requests to specific Grant Type handler
	switch valid := request.(type) {
	case *PasswordRequest:
		context.JSON(http.StatusOK, valid)
	default:
		context.JSON(http.StatusInternalServerError, ErrorResponse{Error: UnsupportedGrantType, Description: reflect.TypeOf(valid).Name()})
	}
}
