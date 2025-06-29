package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RequireUrlEncodedForm(context *gin.Context) {
	if contentType := context.ContentType(); contentType != "application/x-www-form-urlencoded" {
		context.AbortWithStatusJSON(http.StatusUnsupportedMediaType, ErrorBody{
			ErrorType:   "invalid_request",
			Description: "Content-Type must be application/x-www-form-urlencoded",
		})
	}
}
