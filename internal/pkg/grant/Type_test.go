package grant

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGrantTypes_Scanner(t *testing.T) {
	t.Parallel()

	jsonGrantTypes := `["authorization_code","password","refresh_token","urn:ietf:params:oauth:grant-type:jwt-bearer"]`

	t.Run("should be able to scan a string", func(t *testing.T) {
		t.Parallel()

		var grantTypes Types
		require.NoError(t, grantTypes.Scan(jsonGrantTypes))
		assert.Equal(t, Types{AuthorisationCode, Password, RefreshToken, Assertion}, grantTypes)
	})

	t.Run("should be able to scan bytes", func(t *testing.T) {
		t.Parallel()

		var grantTypes Types
		require.NoError(t, grantTypes.Scan([]byte(jsonGrantTypes)))
		assert.Equal(t, Types{AuthorisationCode, Password, RefreshToken, Assertion}, grantTypes)
	})

	t.Run("should reject invalid scan types", func(t *testing.T) {
		t.Parallel()

		var grantTypes Types
		require.ErrorIs(t, grantTypes.Scan(123), ErrUnsupportedGrantTypeSource)
	})
}
