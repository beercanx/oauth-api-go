package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectUris_Scanner(t *testing.T) {
	t.Parallel()

	t.Run("should be able to scan a string", func(t *testing.T) {
		t.Parallel()

		var redirectUris RedirectUris
		require.NoError(t, redirectUris.Scan(`["https://cicada.baconi.co.uk/callback"]`))
		assert.Equal(t, RedirectUris{"https://cicada.baconi.co.uk/callback"}, redirectUris)
	})

	t.Run("should be able to scan bytes", func(t *testing.T) {
		t.Parallel()

		var redirectUris RedirectUris
		require.NoError(t, redirectUris.Scan([]byte(`["https://cicada.baconi.co.uk/callback"]`)))
		assert.Equal(t, RedirectUris{"https://cicada.baconi.co.uk/callback"}, redirectUris)
	})

	t.Run("should reject invalid scan types", func(t *testing.T) {
		t.Parallel()

		var redirectUris RedirectUris
		require.ErrorIs(t, redirectUris.Scan(123), ErrUnsupportedRedirectUrisSource)
	})
}
