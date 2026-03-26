package user

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticatedUsername(t *testing.T) {

	t.Run("when printing as a string", func(t *testing.T) {
		assert.Equal(t, "aardvark", AuthenticatedUsername{"aardvark"}.String())
	})

	t.Run("when marshalling to JSON", func(t *testing.T) {
		result, err := json.Marshal(AuthenticatedUsername{"AARDVARK"})
		assert.Nil(t, err)
		assert.Equal(t, `"AARDVARK"`, string(result))
	})

	t.Run("when unmarshalling from JSON", func(t *testing.T) {
		var underTest AuthenticatedUsername
		err := json.Unmarshal([]byte(`"badger"`), &underTest)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, errors.ErrUnsupported)
		assert.Zero(t, underTest)
	})
}
