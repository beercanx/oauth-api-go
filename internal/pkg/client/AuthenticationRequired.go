package client

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const AuthClientKey = "auth:client"
const AuthClientConfidentialKey = "auth:client:confidential"
const AuthClientPublicKey = "auth:client:public"

func AuthenticateConfidentialClient(authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if clientId, clientSecret, hadBasicAuth := c.Request.BasicAuth(); hadBasicAuth {
			if principal, principalOk := authenticator.AuthenticateAsConfidential(clientId, clientSecret); principalOk {
				c.Set(AuthClientKey, principal)
				c.Set(AuthClientConfidentialKey, principal)
			}
		}
	}
}

func AuthenticatePublicClient(authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		if clientId, clientIdOk := c.GetPostForm("client_id"); clientIdOk {
			if principal, principalOk := authenticator.AuthenticateAsPublic(clientId); principalOk {
				c.Set(AuthClientKey, principal)
				c.Set(AuthClientPublicKey, principal)
			}
		}
	}
}

func RequireClientAuthentication(c *gin.Context) {

	_, clientOk := c.Get(AuthClientKey)
	_, confidentialOk := c.Get(AuthClientConfidentialKey)
	_, publicOk := c.Get(AuthClientPublicKey)

	switch { // TODO - Do we need any response headers?
	case !clientOk:
		c.AbortWithStatus(http.StatusUnauthorized) // Cannot have none
	case !confidentialOk && !publicOk:
		c.AbortWithStatus(http.StatusUnauthorized) // Cannot have none
	case confidentialOk && publicOk:
		c.AbortWithStatus(http.StatusUnauthorized) // Cannot have both
	}
}

func RequireConfidentialClientAuthentication(c *gin.Context) {

	_, clientOk := c.Get(AuthClientKey)
	_, confidentialOk := c.Get(AuthClientConfidentialKey)

	if !clientOk || !confidentialOk {
		// TODO - Do we need any response headers?
		c.AbortWithStatus(http.StatusUnauthorized) // Cannot have none
	}
}
