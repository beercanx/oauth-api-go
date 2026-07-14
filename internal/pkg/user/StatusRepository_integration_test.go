package user

import (
	"testing"

	"baconi.co.uk/oauth/internal/pkg/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusRepository(t *testing.T) {
	t.Parallel()

	database, databaseError := db.Connect("file:user_status_repository_integration_tests?mode=memory&cache=shared")
	require.NoError(t, databaseError)
	require.NoError(t, db.RunMigrations(database, "file:../../../sql/migrations"))

	underTest := NewStatusRepository(t.Context(), database)

	t.Run("insert", func(t *testing.T) {
		require.NoError(t, underTest.insert(Status{username: "insert", locked: false}))
	})

	t.Run("findByUsername with valid status", func(t *testing.T) {
		require.NoError(t, underTest.insert(Status{username: "findByUsername", locked: true}))
		status, err := underTest.FindByUsername("findByUsername")
		require.NoError(t, err)
		assert.Equal(t, Status{username: "findByUsername", locked: true}, status)
	})

	t.Run("findByUsername with no status for username", func(t *testing.T) {
		status, err := underTest.FindByUsername("no-such-username")
		require.ErrorIs(t, err, ErrNoSuchUserStatus)
		assert.Zero(t, status)
	})
}
