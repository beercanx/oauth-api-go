package collections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// FuzzFold tests the Fold function with random inputs
func FuzzFold(f *testing.F) {
	// Add seed corpus
	f.Add([]byte{1, 2, 3, 4, 5})
	f.Add([]byte{10, 20, 30})
	f.Add([]byte{})
	f.Add([]byte{42})

	// Fuzz test
	f.Fuzz(func(t *testing.T, data []byte) {
		// Convert bytes to ints
		numbers := make([]int, len(data))
		for i, b := range data {
			numbers[i] = int(b)
		}

		// Test sum fold
		sum := func(acc, item int) int {
			return acc + item
		}

		// This should never panic
		result := Fold(numbers, 0, sum)

		// Verify result
		expected := 0
		for _, n := range numbers {
			expected += n
		}
		assert.Equal(t, expected, result, "Expected sum %d, got %d", expected, result)

		// Test with non-zero initial value
		initialValue := 10
		result = Fold(numbers, initialValue, sum)
		assert.Equal(t, expected + initialValue, result, "Expected sum with initial value %d, got %d", expected + initialValue, result)
	})
}
