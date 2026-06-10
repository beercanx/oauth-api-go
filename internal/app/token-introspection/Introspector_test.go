package token_introspection

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/token"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errorNoDatabase = errors.New("no database")

func TestIntrospector(t *testing.T) {
	t.Parallel()

	t.Run("when token repository errors", func(t *testing.T) {
		t.Parallel()

		accessTokenRepository := token.NewMockAuthenticator[token.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			Authenticate(mock.AnythingOfType("uuid.UUID")).
			Return(token.AccessToken{}, errorNoDatabase).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: uuid.New()})
		require.Error(t, err)
		require.ErrorIs(t, err, errorNoDatabase)
		assert.Zero(t, result)
		assert.IsType(t, response{}, result)
	})

	for name, authenticatorError := range map[string]error{
		"does not exist":   token.ErrNoSuchToken,
		"has expired":      token.ErrTokenHasExpired,
		"is not yet valid": token.ErrTokenIsBefore,
	} {
		t.Run(fmt.Sprintf("when token %s", name), func(t *testing.T) {
			t.Parallel()

			authenticator := token.NewMockAuthenticator[token.AccessToken](t)
			authenticator.
				EXPECT().
				Authenticate(mock.AnythingOfType("uuid.UUID")).
				Return(token.AccessToken{}, authenticatorError).
				Once()

			underTest := NewIntrospector(authenticator)

			result, err := underTest.introspect(request{token: uuid.New()})
			require.NoError(t, err)
			assert.Equal(t, response{Active: false}, result)
		})
	}

	t.Run("when token is just right", func(t *testing.T) {
		t.Parallel()

		now := time.Now()

		accessToken := token.AccessToken{
			ID:        uuid.New(),
			Username:  user.AuthenticatedUsername("aardvark"),
			Scopes:    scope.Scopes{"basic"},
			ClientID:  client.Id("v"),
			IssuedAt:  now,
			ExpiresAt: now.Add(time.Minute),
			NotBefore: now.Add(-time.Minute),
		}

		accessTokenRepository := token.NewMockAuthenticator[token.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			Authenticate(accessToken.ID).
			Return(accessToken, nil).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: accessToken.ID})
		require.NoError(t, err)
		assert.NotZero(t, result)
		assert.Equal(t, response{
			Active:         true,
			Scope:          scope.Scopes{"basic"},
			Subject:        "aardvark",
			Username:       "aardvark",
			ClientId:       "v",
			TokenType:      token.Bearer,
			IssuedAt:       now.Unix(),
			NotBefore:      now.Add(-time.Minute).Unix(),
			ExpirationTime: now.Add(time.Minute).Unix(),
		}, result)
	})
}
