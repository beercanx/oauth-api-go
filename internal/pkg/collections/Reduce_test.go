package collections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReduce(t *testing.T) {
	t.Run("reduces integers to calculate sum", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		sum := func(a, b int) int {
			return a + b
		}

		// Act
		result := Reduce(numbers, sum)

		// Assert
		assert.Equal(t, 15, result)
	})

	t.Run("reduces integers to calculate product", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		product := func(a, b int) int {
			return a * b
		}

		// Act
		result := Reduce(numbers, product)

		// Assert
		assert.Equal(t, 120, result)
	})

	t.Run("reduces strings to concatenate", func(t *testing.T) {
		// Arrange
		words := []string{"Hello", " ", "World", "!"}
		concat := func(a, b string) string {
			return a + b
		}

		// Act
		result := Reduce(words, concat)

		// Assert
		assert.Equal(t, "Hello World!", result)
	})

	t.Run("reduces to find maximum value", func(t *testing.T) {
		// Arrange
		numbers := []int{3, 1, 7, 4, 2}
		calculateMax := func(a, b int) int {
			if a > b {
				return a
			}
			return b
		}

		// Act
		result := Reduce(numbers, calculateMax)

		// Assert
		assert.Equal(t, 7, result)
	})

	t.Run("reduces to find minimum value", func(t *testing.T) {
		// Arrange
		numbers := []int{3, 1, 7, 4, 2}
		calculateMin := func(a, b int) int {
			if a < b {
				return a
			}
			return b
		}

		// Act
		result := Reduce(numbers, calculateMin)

		// Assert
		assert.Equal(t, 1, result)
	})

	t.Run("returns zero value when array is empty", func(t *testing.T) {
		// Arrange
		var numbers []int
		sum := func(a, b int) int {
			return a + b
		}

		// Act
		result := Reduce(numbers, sum)

		// Assert
		assert.Equal(t, 0, result)
	})

	t.Run("returns first element when array has only one element", func(t *testing.T) {
		// Arrange
		numbers := []int{42}
		sum := func(a, b int) int {
			return a + b
		}

		// Act
		result := Reduce(numbers, sum)

		// Assert
		assert.Equal(t, 42, result)
	})

	t.Run("works with custom struct types", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}

		people := []Person{
			{"Alice", 25},
			{"Bob", 17},
			{"Charlie", 30},
			{"David", 16},
		}

		findOldest := func(a, b Person) Person {
			if a.Age > b.Age {
				return a
			}
			return b
		}

		// Act
		result := Reduce(people, findOldest)

		// Assert
		assert.Equal(t, Person{"Charlie", 30}, result)
	})

	t.Run("reduces with boolean operations", func(t *testing.T) {
		// Arrange
		conditions := []bool{true, false, true, true, false}

		// Test AND operation
		and := func(a, b bool) bool {
			return a && b
		}

		// Test OR operation
		or := func(a, b bool) bool {
			return a || b
		}

		// Act
		resultAnd := Reduce(conditions, and)
		resultOr := Reduce(conditions, or)

		// Assert
		assert.False(t, resultAnd)
		assert.True(t, resultOr)
	})

	t.Run("reduces with floating point numbers", func(t *testing.T) {
		// Arrange
		numbers := []float64{1.5, 2.5, 3.5, 4.5, 5.5}
		sum := func(a, b float64) float64 {
			return a + b
		}

		// Act
		result := Reduce(numbers, sum)

		// Assert
		assert.Equal(t, 17.5, result)
	})
}
