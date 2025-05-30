package collections

import (
	"fmt"
	"testing"
)

// Benchmark_Map measures the performance of the Map function
func Benchmark_Map(b *testing.B) {
	// Prepare test data
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	double := func(n int) int {
		return n * 2
	}

	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		Map(numbers, double)
	}
}

// Benchmark different sizes for Map function
func Benchmark_Map_Size(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// Prepare test data
			numbers := make([]int, size)
			for i := range numbers {
				numbers[i] = i
			}

			double := func(n int) int {
				return n * 2
			}

			b.ResetTimer()

			// Run the benchmark
			for i := 0; i < b.N; i++ {
				Map(numbers, double)
			}
		})
	}
}