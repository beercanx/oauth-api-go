package token

import (
	"strings"
	"testing"
	"time"

	"baconi.co.uk/oauth/internal/pkg/client"
	"baconi.co.uk/oauth/internal/pkg/db"
	"baconi.co.uk/oauth/internal/pkg/scope"
	"baconi.co.uk/oauth/internal/pkg/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertAccessTokensEqual compares two access tokens for equality.
// This helper exists because of https://github.com/stretchr/testify/issues/984
func assertAccessTokensEqual(t *testing.T, expected db.AccessToken, actual db.AccessToken) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Username, actual.Username)
	assert.Equal(t, expected.ClientID, actual.ClientID)
	assert.Equal(t, expected.Scopes, actual.Scopes)
	assert.WithinDuration(t, expected.IssuedAt, actual.IssuedAt, time.Millisecond)
	assert.WithinDuration(t, expected.ExpiresAt, actual.ExpiresAt, time.Millisecond)
	assert.WithinDuration(t, expected.NotBefore, actual.NotBefore, time.Millisecond)
}

func TestAccessTokenRepository(t *testing.T) {
	t.Parallel()

	database, databaseError := db.Connect("file:access_tokens_integration_tests?mode=memory&cache=shared")
	require.NoError(t, databaseError)
	require.NoError(t, db.RunMigrations(database, "file:../../../sqlc/migrations"))

	underTest := NewAccessTokenRepository(t.Context(), database)

	validAccessToken := db.CreateAccessTokenParams{
		ID:        uuid.MustParse("dad063b2-bf86-4aed-a505-a8329535c0a8"),
		Username:  "aardvark",
		ClientID:  "badger",
		Scopes:    scope.Scopes{"basic", "read", "write"},
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(time.Hour),
		NotBefore: time.Now().Add(-time.Minute),
	}

	t.Run("insert access token", func(t *testing.T) {

		t.Run("should reject if access token is zero", func(t *testing.T) {
			result := underTest.Insert(db.CreateAccessTokenParams{})
			require.ErrorContains(t, result, "sqlite3: constraint failed")
		})

		t.Run("should reject if ID is default", func(t *testing.T) {
			missingId := validAccessToken
			missingId.ID = uuid.UUID{}
			require.ErrorContains(t, underTest.Insert(missingId), "CHECK constraint failed: id != '00000000-0000-0000-0000-000000000000'")
		})

		t.Run("should reject if username is blank", func(t *testing.T) {
			blankUsername := validAccessToken
			blankUsername.Username = ""
			require.ErrorContains(t, underTest.Insert(blankUsername), "CHECK constraint failed: LENGTH(username) > 0")
		})

		t.Run("should reject if username is too long", func(t *testing.T) {
			tooLongUsername := validAccessToken
			tooLongUsername.Username = user.AuthenticatedUsername(strings.Repeat("a", 65))
			require.ErrorContains(t, underTest.Insert(tooLongUsername), "CHECK constraint failed: LENGTH(username) <= 64")
		})

		t.Run("should reject if client ID is blank", func(t *testing.T) {
			blankClientID := validAccessToken
			blankClientID.ClientID = ""
			require.ErrorContains(t, underTest.Insert(blankClientID), "FOREIGN KEY constraint failed")
		})

		t.Run("should reject if client ID is too long", func(t *testing.T) {
			tooLongClientID := validAccessToken
			tooLongClientID.ClientID = client.Id(strings.Repeat("a", 65))
			require.ErrorContains(t, underTest.Insert(tooLongClientID), "FOREIGN KEY constraint failed")
		})

		t.Run("should reject if issued at is default", func(t *testing.T) {
			defaultIssuedAt := validAccessToken
			defaultIssuedAt.IssuedAt = time.Time{}
			require.ErrorContains(t, underTest.Insert(defaultIssuedAt), "CHECK constraint failed: issued_at != '0001-01-01T00:00:00Z'")
		})

		t.Run("should reject if expired at is default", func(t *testing.T) {
			defaultExpiresAt := validAccessToken
			defaultExpiresAt.ExpiresAt = time.Time{}
			require.ErrorContains(t, underTest.Insert(defaultExpiresAt), "CHECK constraint failed: expires_at != '0001-01-01T00:00:00Z'")
		})

		t.Run("should reject if not before is default", func(t *testing.T) {
			defaultNotBefore := validAccessToken
			defaultNotBefore.NotBefore = time.Time{}
			require.ErrorContains(t, underTest.Insert(defaultNotBefore), "CHECK constraint failed: not_before != '0001-01-01T00:00:00Z'")
		})

		t.Run("should allow only one insertion of an access token", func(t *testing.T) {
			require.NoError(t, underTest.Insert(validAccessToken))
			require.ErrorContains(t, underTest.Insert(validAccessToken), "UNIQUE constraint failed: access_tokens.id")
		})
	})

	t.Run("find access token by id", func(t *testing.T) {

		validAccessToken := validAccessToken
		validAccessToken.ID = uuid.New()
		require.NoError(t, underTest.Insert(validAccessToken))

		t.Run("should return no such token error if no access token is found", func(t *testing.T) {
			result, err := underTest.FindById(uuid.New())
			require.ErrorIs(t, err, ErrNoSuchToken)
			assert.Zero(t, result)
		})

		t.Run("should return access token if found", func(t *testing.T) {
			result, err := underTest.FindById(validAccessToken.ID)
			require.NoError(t, err)
			assertAccessTokensEqual(t, db.AccessToken(validAccessToken), result)
		})
	})

	t.Run("delete access token by id", func(t *testing.T) {

		validAccessToken := validAccessToken
		validAccessToken.ID = uuid.New()
		require.NoError(t, underTest.Insert(validAccessToken))

		t.Run("should return no error if no access token is found", func(t *testing.T) {
			require.NoError(t, underTest.DeleteById(uuid.New()))
		})

		t.Run("should return no error if access token is found", func(t *testing.T) {
			require.NoError(t, underTest.DeleteById(validAccessToken.ID))
			_, err := underTest.FindById(validAccessToken.ID)
			require.ErrorIs(t, err, ErrNoSuchToken)
		})
	})

	t.Run("delete access token by record", func(t *testing.T) {

		validAccessToken := validAccessToken
		validAccessToken.ID = uuid.New()
		require.NoError(t, underTest.Insert(validAccessToken))

		t.Run("should return no error if no access token is found", func(t *testing.T) {
			require.NoError(t, underTest.DeleteByRecord(db.AccessToken{ID: uuid.New()}))
		})

		t.Run("should return no error if access token is found", func(t *testing.T) {
			require.NoError(t, underTest.DeleteByRecord(db.AccessToken(validAccessToken)))
			_, err := underTest.FindById(validAccessToken.ID)
			require.ErrorIs(t, err, ErrNoSuchToken)
		})
	})

	t.Run("delete all expired access tokens", func(t *testing.T) {

		expiredAccessToken := validAccessToken
		expiredAccessToken.ID = uuid.New()
		expiredAccessToken.ExpiresAt = time.Now().Add(-time.Minute)
		require.NoError(t, underTest.Insert(expiredAccessToken))

		t.Run("should delete all expired access tokens", func(t *testing.T) {
			require.NoError(t, underTest.DeletedExpired())

			_, findExpiredErr := underTest.FindById(expiredAccessToken.ID)
			require.ErrorIs(t, findExpiredErr, ErrNoSuchToken)

			_, findValidErr := underTest.FindById(validAccessToken.ID)
			require.NoError(t, findValidErr)
		})
	})
}
