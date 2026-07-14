package user

import (
	"testing"

	"baconi.co.uk/oauth/internal/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCredentialRepository(t *testing.T) {
	t.Parallel()

	database, databaseError := db.Connect("file:user_credential_repository_integration_tests?mode=memory&cache=shared")
	require.NoError(t, databaseError)
	require.NoError(t, db.RunMigrations(database, "file:../../../sql/migrations"))

	statusRepository := NewStatusRepository(t.Context(), database)
	underTest := NewCredentialRepository(t.Context(), database)

	// #nosec G101 -- It's test code with a dummy credential, so it's fine.
	t.Run("insert credential", func(t *testing.T) {
		require.NoError(t, statusRepository.insert(Status{"insert", false}))
		require.NoError(t, underTest.insert(Credential{username: "insert", hashedSecret: "$argon2id$v=19$m=19456,t=2,p=1$ABCDEFGHIJKLMNOPQRS"}))
	})

	// #nosec G101 -- It's test code with a dummy credential, so it's fine.
	t.Run("findByUsername with valid credential", func(t *testing.T) {
		credential, err := underTest.FindByUsername("insert")
		require.NoError(t, err)
		assert.Equal(t, Credential{username: "insert", hashedSecret: "$argon2id$v=19$m=19456,t=2,p=1$ABCDEFGHIJKLMNOPQRS"}, credential)
	})

	t.Run("findByUsername with no credential for username", func(t *testing.T) {
		credential, err := underTest.FindByUsername("no-such-username")
		require.ErrorIs(t, err, ErrNoSuchUserCredential)
		assert.Zero(t, credential)
	})
}
