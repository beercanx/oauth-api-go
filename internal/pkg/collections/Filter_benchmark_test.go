package collections

import (
	"testing"
)

// Benchmark_Filter measures the performance of the Filter function
func Benchmark_Filter(b *testing.B) {
	// Prepare test data
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	isEven := func(n int) bool {
		return n%2 == 0
	}

	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		Filter(numbers, isEven)
	}
}