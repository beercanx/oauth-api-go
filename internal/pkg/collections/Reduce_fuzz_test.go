package collections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// FuzzReduce tests the Reduce function with random inputs
func FuzzReduce(f *testing.F) {
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

		// Test sum reduction
		sum := func(a, b int) int {
			return a + b
		}

		// This should never panic
		result := Reduce(numbers, sum)

		// Verify result for empty slice
		if len(numbers) == 0 {
			assert.Equal(t, 0, result, "Expected 0 for empty slice")
		}

		// Verify result for single element slice
		if len(numbers) == 1 {
			assert.Equal(t, numbers[0], result, "Expected first element for single element slice")
		}
	})
}
