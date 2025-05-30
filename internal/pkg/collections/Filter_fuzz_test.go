package collections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// FuzzFilter tests the Filter function with random inputs
func FuzzFilter(f *testing.F) {
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

		// Test even filter
		isEven := func(n int) bool {
			return n%2 == 0
		}

		// This should never panic
		result := Filter(numbers, isEven)

		// Verify all results satisfy the predicate
		for _, n := range result {
			assert.True(t, isEven(n), "Filter returned %d which doesn't satisfy the predicate", n)
		}

		// Verify result length is correct
		expectedCount := 0
		for _, n := range numbers {
			if isEven(n) {
				expectedCount++
			}
		}
		assert.Equal(t, expectedCount, len(result), "Expected %d elements, got %d", expectedCount, len(result))
	})
}
