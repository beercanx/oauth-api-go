// TODO - Add build tag? go:build integration
package token_introspection

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/db"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/token"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenIntrospectionRequests(t *testing.T) {
	t.Parallel()

	database, databaseError := db.Connect("file:token_introspection_route_integration_tests?mode=memory&cache=shared")
	require.NoError(t, databaseError)
	require.NoError(t, db.RunMigrations(database, "file:../../../sqlc/migrations"))

	accessTokenRepository := token.NewAccessTokenRepository(t.Context(), database)
	accessTokenAuthenticator := token.NewAccessTokenAuthenticator(accessTokenRepository)

	clientSecretRepository := client.NewInMemorySecretRepository()
	clientPrincipalRepository := client.NewInMemoryPrincipalRepository()
	clientAuthenticationService := client.NewAuthenticationService(clientSecretRepository, clientPrincipalRepository)

	tokenIntrospector := NewIntrospector(accessTokenAuthenticator)

	router := gin.New(func(engine *gin.Engine) {
		engine.HandleMethodNotAllowed = true
		engine.Use(gin.Logger(), gin.Recovery())
	})

	Route(router, clientAuthenticationService, tokenIntrospector)

	t.Run("should allow only post requests", func(t *testing.T) {
		t.Parallel()

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
			t.Run(fmt.Sprintf("reject %s", invalidMethod), func(t *testing.T) {
				t.Parallel()

				testRequest := httptest.NewRequestWithContext(t.Context(), invalidMethod, "/introspect", nil)

				recorder := httptest.NewRecorder()

				router.ServeHTTP(recorder, testRequest)

				assert.Equal(t, http.StatusMethodNotAllowed, recorder.Code)
				assert.NotEmpty(t, recorder.Header())
				assert.Equal(t, http.MethodPost, recorder.Header().Get("Allow"))
				assert.Empty(t, recorder.Header().Get("WWW-Authenticate"))
			})
		}
	})

	t.Run("must allow only authorized requests", func(t *testing.T) {
		t.Parallel()

		t.Run("reject missing authentication", func(t *testing.T) {
			t.Parallel()

			testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", nil)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, testRequest)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			assert.Empty(t, recorder.Header())
			assert.Empty(t, recorder.Body.String())
		})

		t.Run("reject invalid basic authentication", func(t *testing.T) {
			t.Parallel()

			testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", nil)
			testRequest.SetBasicAuth("invalid", "invalid")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, testRequest)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			assert.Empty(t, recorder.Header())
			assert.Empty(t, recorder.Body.String())
		})

		t.Run("reject public client authentication", func(t *testing.T) {
			t.Parallel()

			formBody := url.Values{"client_id": {"cicada"}, "token": {"a"}}
			testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", strings.NewReader(formBody.Encode()))
			testRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, testRequest)

			assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			assert.Empty(t, recorder.Header())
			assert.Empty(t, recorder.Body.String())
		})

		t.Run("reject a valid client missing the introspection allowed action", func(t *testing.T) {
			t.Parallel()

			formBody := url.Values{"token": {"a"}}
			testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", strings.NewReader(formBody.Encode()))
			testRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			testRequest.SetBasicAuth("dodo", "echidna")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, testRequest)

			assert.Equal(t, http.StatusForbidden, recorder.Code)
			assert.NotEmpty(t, recorder.Header())
			assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))

			var result invalid
			err := json.Unmarshal(recorder.Body.Bytes(), &result)
			require.NoError(t, err)
			assert.Equal(t, UnauthorizedClient, result.ErrorType)
			assert.Equal(t, "client is not allowed to introspect", result.Description)
		})
	})

	t.Run("should allow only url encoded form requests", func(t *testing.T) {
		t.Parallel()

		for contentType, contentBody := range map[string]string{
			"json": `{"token":"94efe4d7-7dbe-455f-b974-46656fd8d035"}`,
			"xml":  `<request><token>94efe4d7-7dbe-455f-b974-46656fd8d035</token></request>`,
		} {
			t.Run(fmt.Sprintf("reject %s body requests", contentType), func(t *testing.T) {
				t.Parallel()

				testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", strings.NewReader(contentBody))
				testRequest.Header.Set("Content-Type", fmt.Sprintf("application/%s", contentType))
				testRequest.SetBasicAuth("aardvark", "badger")

				recorder := httptest.NewRecorder()

				router.ServeHTTP(recorder, testRequest)

				assert.Equal(t, http.StatusUnsupportedMediaType, recorder.Code)

				var result invalid
				err := json.Unmarshal(recorder.Body.Bytes(), &result)
				require.NoError(t, err)
				assert.Equal(t, InvalidRequest, result.ErrorType)
				assert.Equal(t, "Content-Type must be application/x-www-form-urlencoded", result.Description)
			})
		}
	})

	t.Run("should handle invalid body content", func(t *testing.T) {
		t.Parallel()

		for _, data := range []struct {
			state    string
			body     url.Values
			expected string
		}{
			{"missing", url.Values{}, "missing parameter: token"},
			{"blank", url.Values{"token": {""}}, "invalid parameter: token"},
			{"non uuid", url.Values{"token": {"aardvark"}}, "invalid parameter: token"},
		} {
			t.Run(fmt.Sprintf("return invalid request on %s token", data.state), func(t *testing.T) {
				t.Parallel()

				testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", strings.NewReader(data.body.Encode()))
				testRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				testRequest.SetBasicAuth("aardvark", "badger")

				recorder := httptest.NewRecorder()

				router.ServeHTTP(recorder, testRequest)

				assert.Equal(t, http.StatusBadRequest, recorder.Code)

				var result invalid
				err := json.Unmarshal(recorder.Body.Bytes(), &result)
				require.NoError(t, err)
				assert.Equal(t, InvalidRequest, result.ErrorType)
				assert.Equal(t, data.expected, result.Description)
			})
		}
	})

	t.Run("should handle various token states", func(t *testing.T) {

		t.Run("return an inactive response for an access token that does not exist", func(t *testing.T) {

			formBody := url.Values{"token": {"94efe4d7-7dbe-455f-b974-46656fd8d035"}}
			testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", strings.NewReader(formBody.Encode()))
			testRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			testRequest.SetBasicAuth("aardvark", "badger")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, testRequest)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, `{"active":false}`, recorder.Body.String())
		})

		t.Run("return an inactive response for an access token that has expired", func(t *testing.T) {

			accessToken := db.CreateAccessTokenParams{
				ID:        uuid.New(),
				Username:  "expired",
				ClientID:  "badger",
				Scopes:    scope.Scopes{},
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(-(10 * time.Minute)),
				NotBefore: time.Now().Add(-(20 * time.Minute)),
			}

			require.NoError(t, accessTokenRepository.Insert(accessToken))

			formBody := url.Values{"token": {accessToken.ID.String()}}
			testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", strings.NewReader(formBody.Encode()))
			testRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			testRequest.SetBasicAuth("aardvark", "badger")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, testRequest)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, `{"active":false}`, recorder.Body.String())
		})

		t.Run("return an inactive response for an access token that is in the future", func(t *testing.T) {

			accessToken := db.CreateAccessTokenParams{
				ID:        uuid.New(),
				Username:  "future",
				ClientID:  "badger",
				Scopes:    scope.Scopes{},
				IssuedAt:  time.Now(),
				ExpiresAt: time.Now().Add(20 * time.Minute),
				NotBefore: time.Now().Add(10 * time.Minute),
			}

			require.NoError(t, accessTokenRepository.Insert(accessToken))

			formBody := url.Values{"token": {accessToken.ID.String()}}
			testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", strings.NewReader(formBody.Encode()))
			testRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			testRequest.SetBasicAuth("aardvark", "badger")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, testRequest)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.Equal(t, `{"active":false}`, recorder.Body.String())
		})

		t.Run("return an active response for an access token that is active", func(t *testing.T) {

			accessTokenIssuer := token.NewAccessTokenIssuer(accessTokenRepository)

			username := user.AuthenticatedUsername("ant")
			clientId := client.Id("badger")
			scopes := scope.Scopes{"basic"}
			accessToken, issueError := accessTokenIssuer.Issue(username, clientId, scopes)
			require.NoError(t, issueError)

			formBody := url.Values{"token": {accessToken.ID.String()}}
			testRequest := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/introspect", strings.NewReader(formBody.Encode()))
			testRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			testRequest.SetBasicAuth("aardvark", "badger")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, testRequest)

			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.NotEmpty(t, recorder.Header())
			assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))

			var result map[string]any
			unmarshalError := json.Unmarshal(recorder.Body.Bytes(), &result)
			require.NoError(t, unmarshalError)
			assert.Len(t, result, 9)
			assert.Equal(t, true, result["active"])
			assert.Equal(t, "badger", result["client_id"])
			assert.Contains(t, result, "expiration_time")
			assert.Contains(t, result, "issued_at")
			assert.Contains(t, result, "not_before")
			assert.Equal(t, "basic", result["scope"])
			assert.Equal(t, "ant", result["sub"])
			assert.Equal(t, "bearer", result["token_type"])
			assert.Equal(t, "ant", result["username"])
		})
	})
}
