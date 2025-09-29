// TODO - Add build tag? go:build integration
package token_introspection

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTokenIntrospectionRequests(t *testing.T) {

	accessTokenRepository := token.NewInMemoryRepository[token.AccessToken]()

	clientSecretRepository := client.NewInMemorySecretRepository()
	clientPrincipalRepository := client.NewInMemoryPrincipalRepository()
	clientAuthenticationService := client.NewAuthenticationService(clientSecretRepository, clientPrincipalRepository)

	tokenIntrospector := NewIntrospector(accessTokenRepository)

	router := gin.New(func(engine *gin.Engine) {
		engine.HandleMethodNotAllowed = true
		engine.Use(gin.Logger(), gin.Recovery())
	})

	Route(router, clientAuthenticationService, tokenIntrospector)

	for _, invalidMethod := range []string{
		http.MethodGet,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPut,
		http.MethodTrace,
	} {
		t.Run(fmt.Sprintf("should allow only post requests %s", invalidMethod), func(t *testing.T) {

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(invalidMethod, "/introspect", nil)

			router.ServeHTTP(recorder, request)

			assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
			assert.NotEmpty(t, recorder.Header())
			assert.Equal(t, http.MethodPost, recorder.Header().Get("Allow"))
			assert.Zero(t, recorder.Header().Get("WWW-Authenticate"))
		})
	}

	// must allow only authorised requests
	// - reject missing authentication
	// - reject invalid basic authentication
	// - reject public client authentication
	// - reject a valid client that is missing the introspection allowed action
	// - accept a valid client using basic authentication

	// should allow only url encoded form requests
	// - reject JSON body requests
	// - reject XML body requests
	// - reject non url encoded form posts
	// - accept URL encoded form body requests

}
