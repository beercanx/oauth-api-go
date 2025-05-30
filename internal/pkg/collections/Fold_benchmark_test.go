package collections

import (
	"testing"
)

// Benchmark_Fold measures the performance of the Fold function
func Benchmark_Fold(b *testing.B) {
	// Prepare test data
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	sum := func(acc, item int) int {
		return acc + item
	}

	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		Fold(numbers, 0, sum)
	}
}