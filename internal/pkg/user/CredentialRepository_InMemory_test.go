package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInMemoryCredentialRepository(t *testing.T) {

	underTest := NewInMemoryCredentialRepository()

	t.Run("insert credential", func(t *testing.T) {
		require.NoError(t, underTest.Insert(Credential{username: "insert", hashedSecret: "hash"}))
	})

	t.Run("findByUsername with valid credential", func(t *testing.T) {
		credential, err := underTest.FindByUsername("insert")
		require.NoError(t, err)
		assert.Equal(t, Credential{username: "insert", hashedSecret: "hash"}, credential)
	})

	t.Run("findByUsername with no credential for username", func(t *testing.T) {
		credential, err := underTest.FindByUsername("no-such-username")
		require.ErrorIs(t, err, ErrNoSuchCredential)
		assert.Zero(t, credential)
	})
}
