package scope

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScopes_Marshaler(t *testing.T) {
	t.Parallel()

	t.Run("should be able to marshal empty scopes as empty string", func(t *testing.T) {
		t.Parallel()

		result, err := Scopes{}.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`""`), result)
	})

	t.Run("should be able to marshal scopes as space deliminated string", func(t *testing.T) {
		t.Parallel()

		result, err := Scopes{"basic","read","write"}.MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"basic read write"`), result)
	})
}

func TestScopes_Unmarshaler(t *testing.T) {
	t.Parallel()

	t.Run("should reject any unmarshalling attempts", func(t *testing.T) {
		t.Parallel()
		require.ErrorIs(t, Scopes{}.UnmarshalJSON([]byte(`["basic"]`)), errors.ErrUnsupported)
	})
}

func TestScopes_Valuer(t *testing.T) {
	t.Parallel()

	t.Run("should be able to value empty scopes", func(t *testing.T) {
		t.Parallel()

		result, err := Scopes{}.Value()
		require.NoError(t, err)
		assert.Equal(t, []byte(`[]`), result)
	})

	t.Run("should be able to value scopes", func(t *testing.T) {
		t.Parallel()

		result, err := Scopes{"basic","read","write"}.Value()
		require.NoError(t, err)
		assert.Equal(t, []byte(`["basic","read","write"]`), result)
	})
}

func TestScopes_Scanner(t *testing.T) {
	t.Parallel()

	t.Run("should be able to scan a string", func(t *testing.T) {
		t.Parallel()

		var scopes Scopes
		require.NoError(t, scopes.Scan(`["basic","read","write"]`))
		assert.Equal(t, Scopes{"basic","read","write"}, scopes)
	})

	t.Run("should be able to scan bytes", func(t *testing.T) {
		t.Parallel()

		var scopes Scopes
		require.NoError(t, scopes.Scan([]byte(`["basic","read","write"]`)))
		assert.Equal(t, Scopes{"basic","read","write"}, scopes)
	})

	t.Run("should reject invalid JSON", func(t *testing.T) {
		t.Parallel()

		var scopes Scopes
		require.ErrorContains(t, scopes.Scan(`["basic",`), "unexpected end of JSON input")
	})

	t.Run("should reject invalid scan types", func(t *testing.T) {
		t.Parallel()

		var scopes Scopes
		require.ErrorIs(t, scopes.Scan(123), ErrUnsupportedScopesType)
	})
}