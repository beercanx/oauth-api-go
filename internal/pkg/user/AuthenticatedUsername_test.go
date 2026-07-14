package user

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthenticatedUsername(t *testing.T) {
	t.Parallel()

	t.Run("when printing as a string", func(t *testing.T) {
		t.Parallel()
		assert.Equal(t, "aardvark", AuthenticatedUsername("aardvark").String())
	})

	t.Run("when marshalling to JSON", func(t *testing.T) {
		t.Parallel()
		result, err := json.Marshal(AuthenticatedUsername("AARDVARK"))
		require.NoError(t, err)
		assert.Equal(t, `"AARDVARK"`, string(result))
	})

	t.Run("when unmarshalling from JSON", func(t *testing.T) {
		t.Parallel()
		var underTest AuthenticatedUsername
		err := json.Unmarshal([]byte(`"badger"`), &underTest)
		require.Error(t, err)
		require.ErrorIs(t, err, errors.ErrUnsupported)
		assert.Zero(t, underTest)
	})
}
