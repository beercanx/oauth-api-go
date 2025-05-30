package collections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// FuzzMap tests the Map function with random inputs
func FuzzMap(f *testing.F) {
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

		// Test double map
		double := func(n int) int {
			return n * 2
		}

		// This should never panic
		result := Map(numbers, double)

		// Verify result length
		assert.Equal(t, len(numbers), len(result), "Expected %d elements, got %d", len(numbers), len(result))

		// Verify each element is correctly transformed
		for i, n := range numbers {
			expected := double(n)
			assert.Equal(t, expected, result[i], "Expected %d at index %d, got %d", expected, i, result[i])
		}
	})
}
