package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryStatusRepository(t *testing.T) {

	underTest := NewInMemoryStatusRepository()

	t.Run("insert", func(t *testing.T) {
		assert.NoError(t, underTest.Insert(Status{username: "insert", locked: false}))
	})

	t.Run("findByUsername with valid status", func(t *testing.T) {
		assert.NoError(t, underTest.Insert(Status{username: "findByUsername", locked: true}))
		status, err := underTest.FindByUsername("findByUsername")
		assert.NoError(t, err)
		assert.Equal(t, Status{username: "findByUsername", locked: true}, status)
	})

	t.Run("findByUsername with no status for username", func(t *testing.T) {
		status, err := underTest.FindByUsername("no-such-username")
		assert.ErrorIs(t, err, ErrNoSuchStatus)
		assert.Zero(t, status)
	})
}
