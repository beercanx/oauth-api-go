package token_introspection

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"baconi.co.uk/oauth/internal/pkg/client"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func testContext(t *testing.T) *gin.Context {
	t.Helper()
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Request = &http.Request{}
	ctx.Request.PostForm = make(map[string][]string)
	return ctx
}

func TestRequestValidation(t *testing.T) {
	t.Parallel()

	validIntrospectPrincipal := client.Principal{
		ClientId: "aardvark",
		ClientType: client.Confidential,
		AllowedActions: client.Actions{client.Introspect},
	}

	t.Run("should reject if client principal is not confidential", func(t *testing.T) {
		ctx := testContext(t)
		ctx.Set(client.AuthClientConfidentialKey, client.Principal{ClientType: client.Public})

		valid, failed := validateRequest(ctx)

		assert.Zero(t, valid)
		assert.NotNil(t, failed)
		assert.Equal(t, UnauthorizedClient, failed.ErrorType)
		assert.Equal(t, "client is not allowed to introspect", failed.Description)
	})

	t.Run("should reject if client principal does not have introspection action", func(t *testing.T) {
		ctx := testContext(t)
		ctx.Set(client.AuthClientConfidentialKey, client.Principal{ClientType: client.Confidential, AllowedActions: client.Actions{}})

		valid, failed := validateRequest(ctx)

		assert.Zero(t, valid)
		assert.NotNil(t, failed)
		assert.Equal(t, UnauthorizedClient, failed.ErrorType)
		assert.Equal(t, "client is not allowed to introspect", failed.Description)
	})

	t.Run("should reject if token is not present", func(t *testing.T) {
		ctx := testContext(t)
		ctx.Set(client.AuthClientConfidentialKey, validIntrospectPrincipal)

		valid, failed := validateRequest(ctx)

		assert.Zero(t, valid)
		assert.NotNil(t, failed)
		assert.Equal(t, InvalidRequest, failed.ErrorType)
		assert.Equal(t, "missing parameter: token", failed.Description)
	})

	t.Run("should reject if token is not a UUID", func(t *testing.T) {
		ctx := testContext(t)
		ctx.Set(client.AuthClientConfidentialKey, validIntrospectPrincipal)
		ctx.Request.PostForm.Set("token", "not-a-uuid")

		valid, failed := validateRequest(ctx)

		assert.Zero(t, valid)
		assert.NotNil(t, failed)
		assert.Equal(t, InvalidRequest, failed.ErrorType)
		assert.Equal(t, "invalid parameter: token", failed.Description)
	})

	t.Run("should return a valid request", func(t *testing.T) {
		ctx := testContext(t)
		ctx.Set(client.AuthClientConfidentialKey, validIntrospectPrincipal)
		ctx.Request.PostForm.Set("token", "94efe4d7-7dbe-455f-b974-46656fd8d035")

		valid, failed := validateRequest(ctx)

		assert.Nil(t, failed)
		assert.NotZero(t, valid)
		assert.Equal(t, validIntrospectPrincipal, valid.principal)
		assert.Equal(t, uuid.MustParse("94efe4d7-7dbe-455f-b974-46656fd8d035"), valid.token)
	})
}