package client

import (
	"testing"

	"baconi.co.uk/oauth/internal/pkg/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSecretRepository(t *testing.T) {
	t.Parallel()

	database, databaseError := db.Connect("file:client_secret_integration_tests?mode=memory&cache=shared")
	require.NoError(t, databaseError)
	require.NoError(t, db.RunMigrations(database, "file:../../../sql/migrations"))

	underTest := NewSecretRepository(t.Context(), database)

	// #nosec G101 -- It's test code with an actual credential, so it's fine.
	validSecret := Secret{
		id:           SecretId(uuid.MustParse("508795e7-bde8-44c5-acf0-ca17e01d6a7a")),
		clientId:     "dodo",
		hashedSecret: "$argon2id$v=19$m=19456,t=2,p=1$YllYTFgxbFZNVk9wVWp6NQ$1uEMXAYYbW7w+/mwlsS/vG1Kmg",
	}

	t.Run("find client secret by id", func(t *testing.T) {

		t.Run("should return no such client secret error if no client secret is found", func(t *testing.T) {
			result, err := underTest.FindById(SecretId(uuid.New()))
			require.ErrorIs(t, err, ErrNoSuchClientSecret)
			assert.Zero(t, result)
		})

		t.Run("should return client secret if found", func(t *testing.T) {
			result, err := underTest.FindById(validSecret.id)
			require.NoError(t, err)
			assert.Equal(t, validSecret, result)
		})
	})

	t.Run("find client secret by client", func(t *testing.T) {

		t.Run("should return empty secrets array if no client secret is found", func(t *testing.T) {
			results, err := underTest.FindByClient("no-such-client")
			require.NoError(t, err)
			assert.Empty(t, results)
		})

		t.Run("should return client secret if found", func(t *testing.T) {
			results, err := underTest.FindByClient(validSecret.clientId)
			require.NoError(t, err)
			assert.Contains(t, results, validSecret)
		})
	})

	t.Run("find client secret by client id string", func(t *testing.T) {

		t.Run("should return empty secrets array if no client secret is found", func(t *testing.T) {
			results, err := underTest.FindByClientId("no-such-client")
			require.NoError(t, err)
			assert.Empty(t, results)
		})

		t.Run("should return client secret if found", func(t *testing.T) {
			results, err := underTest.FindByClientId(string(validSecret.clientId))
			require.NoError(t, err)
			assert.Contains(t, results, validSecret)
		})
	})
}
