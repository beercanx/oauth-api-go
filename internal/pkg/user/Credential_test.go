package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryCredentialRepository(t *testing.T) {

	underTest := NewInMemoryCredentialRepository()

	t.Run("insert credential", func(t *testing.T) {
		assert.NoError(t, underTest.Insert(Credential{username: "insert", hashedSecret: "hash"}))
	})

	t.Run("findByUsername with valid credential", func(t *testing.T) {
		credential, err := underTest.FindByUsername("insert")
		assert.NoError(t, err)
		assert.Equal(t, Credential{username: "insert", hashedSecret: "hash"}, credential)
	})

	t.Run("findByUsername with no credential for username", func(t *testing.T) {
		credential, err := underTest.FindByUsername("no-such-username")
		assert.ErrorIs(t, err, ErrNoSuchCredential)
		assert.Zero(t, credential)
	})
}
