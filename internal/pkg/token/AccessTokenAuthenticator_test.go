package token

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errorNoDatabase = errors.New("no database")

func TestAccessTokenAuthenticator_Authenticate(t *testing.T) {
	t.Parallel()

	t.Run("should return error is repository returns error", func(t *testing.T) {
		t.Parallel()

		repository := NewMockRepositoryReadDelete[AccessToken](t)
		underTest := NewAccessTokenAuthenticator(repository)

		repository.EXPECT().FindById(mock.AnythingOfType("uuid.UUID")).Return(AccessToken{}, errorNoDatabase).Once()

		token, err := underTest.Authenticate(uuid.New())

		require.ErrorIs(t, err, errorNoDatabase)
		assert.Zero(t, token)
	})

	t.Run("should return error if access token is expired", func(t *testing.T) {
		t.Parallel()

		repository := NewMockRepositoryReadDelete[AccessToken](t)
		underTest := NewAccessTokenAuthenticator(repository)

		token := uuid.New()
		dbToken := AccessToken{ID: token, ExpiresAt: time.Now().Add(-time.Hour)}

		repository.EXPECT().FindById(token).Return(dbToken, nil).Once()
		repository.EXPECT().DeleteByRecord(dbToken).Return(nil).Once()

		accessToken, err := underTest.Authenticate(token)

		require.ErrorIs(t, err, ErrTokenHasExpired)
		assert.Zero(t, accessToken)
		assert.NotEqual(t, dbToken, accessToken)
	})

	t.Run("should return error if access token is before now", func(t *testing.T) {
		t.Parallel()

		repository := NewMockRepositoryReadDelete[AccessToken](t)
		underTest := NewAccessTokenAuthenticator(repository)

		token := uuid.New()
		dbToken := AccessToken{ID: token, ExpiresAt: time.Now().Add(time.Hour), NotBefore: time.Now().Add(time.Minute)}

		repository.EXPECT().FindById(token).Return(dbToken, nil).Once()

		accessToken, err := underTest.Authenticate(token)

		require.ErrorIs(t, err, ErrTokenIsBefore)
		assert.Zero(t, accessToken)
		assert.NotEqual(t, dbToken, accessToken)
	})

	t.Run("should return access token if it is valid", func(t *testing.T) {
		t.Parallel()

		repository := NewMockRepositoryReadDelete[AccessToken](t)
		underTest := NewAccessTokenAuthenticator(repository)

		token := uuid.New()
		dbToken := AccessToken{ID: token, ExpiresAt: time.Now().Add(time.Hour), NotBefore: time.Now().Add(-time.Minute)}

		repository.EXPECT().FindById(token).Return(dbToken, nil).Once()

		accessToken, err := underTest.Authenticate(token)

		require.NoError(t, err)
		assert.Equal(t, dbToken, accessToken)
	})
}