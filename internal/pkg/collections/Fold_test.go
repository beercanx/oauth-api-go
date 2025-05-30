package collections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFold(t *testing.T) {
	t.Run("folds integers to calculate sum", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		sum := func(acc int, item int) int {
			return acc + item
		}

		// Act
		result := Fold(numbers, 0, sum)

		// Assert
		assert.Equal(t, 15, result)
	})

	t.Run("folds integers with non-zero initial value", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		sum := func(acc int, item int) int {
			return acc + item
		}

		// Act
		result := Fold(numbers, 10, sum)

		// Assert
		assert.Equal(t, 25, result)
	})

	t.Run("folds strings to concatenate", func(t *testing.T) {
		// Arrange
		words := []string{"Hello", " ", "World", "!"}
		concat := func(acc string, item string) string {
			return acc + item
		}

		// Act
		result := Fold(words, "", concat)

		// Assert
		assert.Equal(t, "Hello World!", result)
	})

	t.Run("folds with different input and accumulator types", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		toString := func(acc string, item int) string {
			if acc == "" {
				return string(rune('0' + item))
			}
			return acc + "," + string(rune('0'+item))
		}

		// Act
		result := Fold(numbers, "", toString)

		// Assert
		assert.Equal(t, "1,2,3,4,5", result)
	})

	t.Run("returns initial value when array is empty", func(t *testing.T) {
		// Arrange
		var numbers []int
		sum := func(acc int, item int) int {
			return acc + item
		}

		// Act
		result := Fold(numbers, 42, sum)

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

		calculateTotalAge := func(acc int, person Person) int {
			return acc + person.Age
		}

		// Act
		result := Fold(people, 0, calculateTotalAge)

		// Assert
		assert.Equal(t, 88, result)
	})

	t.Run("folds to find maximum value", func(t *testing.T) {
		// Arrange
		numbers := []int{3, 1, 7, 4, 2}
		calculateMax := func(acc int, item int) int {
			if item > acc {
				return item
			}
			return acc
		}

		// Act
		result := Fold(numbers, numbers[0], calculateMax)

		// Assert
		assert.Equal(t, 7, result)
	})

	t.Run("folds to count occurrences", func(t *testing.T) {
		// Arrange
		words := []string{"apple", "banana", "apple", "cherry", "apple", "date"}

		countMap := make(map[string]int)
		counter := func(acc map[string]int, item string) map[string]int {
			acc[item]++
			return acc
		}

		// Act
		result := Fold(words, countMap, counter)

		// Assert
		assert.Equal(t, 3, result["apple"])
		assert.Equal(t, 1, result["banana"])
		assert.Equal(t, 1, result["cherry"])
		assert.Equal(t, 1, result["date"])
	})
}
