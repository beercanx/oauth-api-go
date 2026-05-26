package token_introspection

import (
	"errors"
	"testing"
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/db"
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

		accessTokenRepository := token.NewMockRepository[db.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(mock.AnythingOfType("uuid.UUID")).
			Return(db.AccessToken{}, errorNoDatabase).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: uuid.New()})
		require.Error(t, err)
		require.ErrorIs(t, err, errorNoDatabase)
		assert.Zero(t, result)
		assert.IsType(t, response{}, result)
	})

	t.Run("when token does not exist", func(t *testing.T) {
		t.Parallel()

		accessTokenRepository := token.NewMockRepository[db.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(mock.AnythingOfType("uuid.UUID")).
			Return(db.AccessToken{}, token.ErrNoSuchToken).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: uuid.New()})
		require.NoError(t, err)
		assert.Equal(t, response{Active: false}, result)
	})

	t.Run("when token has expired", func(t *testing.T) {
		t.Parallel()

		now := time.Now()

		accessTokenRepository := token.NewMockRepository[db.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(mock.AnythingOfType("uuid.UUID")).
			Return(db.AccessToken{IssuedAt: now, ExpiresAt: now.Add(-time.Hour), NotBefore: now.Add(-time.Hour)}, nil).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: uuid.New()})
		require.NoError(t, err)
		assert.Equal(t, response{Active: false}, result)
	})

	t.Run("when token is not yet valid", func(t *testing.T) {
		t.Parallel()

		now := time.Now()

		accessTokenRepository := token.NewMockRepository[db.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(mock.AnythingOfType("uuid.UUID")).
			Return(db.AccessToken{IssuedAt: now, ExpiresAt: now.Add(time.Minute), NotBefore: now.Add(time.Minute)}, nil).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: uuid.New()})
		require.NoError(t, err)
		assert.Equal(t, response{Active: false}, result)
	})

	t.Run("when token is just right", func(t *testing.T) {
		t.Parallel()

		now := time.Now()

		accessToken := db.AccessToken{
			ID:        uuid.New(),
			Username:  user.AuthenticatedUsername("aardvark"),
			Scopes:    scope.Scopes{"basic"},
			ClientID:  client.Id("v"),
			IssuedAt:  now,
			ExpiresAt: now.Add(time.Minute),
			NotBefore: now.Add(-time.Minute),
		}

		accessTokenRepository := token.NewMockRepository[db.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(accessToken.ID).
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
