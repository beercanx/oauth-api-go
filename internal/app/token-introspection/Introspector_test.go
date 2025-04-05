package token_introspection

import (
	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/token"
	"baconi.co.uk/oauth/internal/pkg/user"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestIntrospector(t *testing.T) {

	t.Run("when token repository errors", func(t *testing.T) {

		noDatabase := errors.New("no database")

		accessTokenRepository := token.NewMockRepository[token.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(mock.AnythingOfType("uuid.UUID")).
			Return(token.AccessToken{}, noDatabase).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: uuid.New()})
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, noDatabase)
		assert.Zero(t, result)
		assert.IsType(t, response{}, result)
	})

	t.Run("when token does not exist", func(t *testing.T) {

		accessTokenRepository := token.NewMockRepository[token.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(mock.AnythingOfType("uuid.UUID")).
			Return(token.AccessToken{}, token.ErrNoSuchToken).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: uuid.New()})
		assert.Nil(t, err)
		assert.Equal(t, response{Active: false}, result)
	})

	t.Run("when token has expired", func(t *testing.T) {

		now := time.Now()

		accessTokenRepository := token.NewMockRepository[token.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(mock.AnythingOfType("uuid.UUID")).
			Return(token.AccessToken{IssuedAt: now, ExpiresAt: now.Add(-time.Hour), NotBefore: now.Add(-time.Hour)}, nil).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: uuid.New()})
		assert.Nil(t, err)
		assert.Equal(t, response{Active: false}, result)
	})

	t.Run("when token is not yet valid", func(t *testing.T) {

		now := time.Now()

		accessTokenRepository := token.NewMockRepository[token.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(mock.AnythingOfType("uuid.UUID")).
			Return(token.AccessToken{IssuedAt: now, ExpiresAt: now.Add(time.Minute), NotBefore: now.Add(time.Minute)}, nil).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: uuid.New()})
		assert.Nil(t, err)
		assert.Equal(t, response{Active: false}, result)
	})

	t.Run("when token is just right", func(t *testing.T) {

		now := time.Now()

		accessToken := token.AccessToken{
			Value:     uuid.New(),
			Username:  user.AuthenticatedUsername{Value: "aardvark"},
			Scopes:    scope.Scopes{Value: []scope.Scope{{Value: "basic"}}},
			ClientId:  client.Id{Value: "v"},
			IssuedAt:  now,
			ExpiresAt: now.Add(time.Minute),
			NotBefore: now.Add(-time.Minute),
		}

		accessTokenRepository := token.NewMockRepository[token.AccessToken](t)
		accessTokenRepository.
			EXPECT().
			FindById(accessToken.Value).
			Return(accessToken, nil).
			Once()

		underTest := NewIntrospector(accessTokenRepository)

		result, err := underTest.introspect(request{token: accessToken.Value})
		assert.Nil(t, err)
		assert.NotZero(t, result)
		assert.Equal(t, response{
			Active:         true,
			Scope:          scope.Scopes{Value: []scope.Scope{{Value: "basic"}}},
			Subject:        user.AuthenticatedUsername{Value: "aardvark"},
			Username:       user.AuthenticatedUsername{Value: "aardvark"},
			ClientId:       client.Id{Value: "v"},
			TokenType:      token.Bearer,
			IssuedAt:       now.Unix(),
			NotBefore:      now.Add(-time.Minute).Unix(),
			ExpirationTime: now.Add(time.Minute).Unix(),
		}, result)
	})
}
