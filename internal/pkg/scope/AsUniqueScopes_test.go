package scope

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	input    []string
	expected Scopes
}

func TestAsUniqueScopes(t *testing.T) {
	t.Parallel()

	t.Run("should return unique scopes", func(t *testing.T) {
		t.Parallel()

		for _, data := range []testData{
			{
				input:    []string{},
				expected: Scopes{},
			},
			{
				input:    []string{""},
				expected: Scopes{""},
			},
			{
				input:    []string{"", "-"},
				expected: Scopes{"", "-"},
			},
			{
				input:    []string{"aardvark", "aardvark"},
				expected: Scopes{"aardvark"},
			},
			{
				input:    []string{"aardvark", "badger", "cicada", "dodo", "dodo"},
				expected: Scopes{"aardvark", "badger", "cicada", "dodo"},
			},
			{
				input:    []string{"aardvark", "badger", "cicada", "dodo"},
				expected: Scopes{"aardvark", "badger", "cicada", "dodo"},
			},
		} {
			assert.Equal(t, data.expected, AsUniqueScopes(data.input))
		}
	})
}
