package collections

import (
	"testing"
)

// Benchmark_Reduce measures the performance of the Reduce function
func Benchmark_Reduce(b *testing.B) {
	// Prepare test data
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	sum := func(a, b int) int {
		return a + b
	}

	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		Reduce(numbers, sum)
	}
}