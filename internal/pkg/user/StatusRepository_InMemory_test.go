package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryStatusRepository(t *testing.T) {
	t.Parallel()

	underTest := NewInMemoryStatusRepository()

	t.Run("insert", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, underTest.Insert(Status{username: "insert", locked: false}))
	})

	t.Run("findByUsername with valid status", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, underTest.Insert(Status{username: "findByUsername", locked: true}))
		status, err := underTest.FindByUsername("findByUsername")
		require.NoError(t, err)
		assert.Equal(t, Status{username: "findByUsername", locked: true}, status)
	})

	t.Run("findByUsername with no status for username", func(t *testing.T) {
		t.Parallel()
		status, err := underTest.FindByUsername("no-such-username")
		require.ErrorIs(t, err, ErrNoSuchStatus)
		assert.Zero(t, status)
	})
}
