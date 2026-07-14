package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActions_Scanner(t *testing.T) {
	t.Parallel()

	t.Run("should be able to scan a string", func(t *testing.T) {
		t.Parallel()

		var actions Actions
		require.NoError(t, actions.Scan(`["authorise","pkce"]`))
		assert.Equal(t, Actions{Authorise, ProofKeyForCodeExchange}, actions)
	})

	t.Run("should be able to scan bytes", func(t *testing.T) {
		t.Parallel()

		var actions Actions
		require.NoError(t, actions.Scan([]byte(`["authorise","pkce"]`)))
		assert.Equal(t, Actions{Authorise, ProofKeyForCodeExchange}, actions)
	})

	t.Run("should reject invalid scan types", func(t *testing.T) {
		t.Parallel()

		var actions Actions
		require.ErrorIs(t, actions.Scan(123), ErrUnsupportedActionsSource)
	})
}