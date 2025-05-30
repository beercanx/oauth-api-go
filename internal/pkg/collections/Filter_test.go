package collections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	t.Run("filters integers based on predicate", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		isEven := func(n int) bool {
			return n%2 == 0
		}

		// Act
		result := Filter(numbers, isEven)

		// Assert
		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("filters strings based on predicate", func(t *testing.T) {
		// Arrange
		words := []string{"apple", "banana", "cherry", "date", "elderberry"}
		startsWithB := func(s string) bool {
			return len(s) > 0 && s[0] == 'b'
		}

		// Act
		result := Filter(words, startsWithB)

		// Assert
		assert.Equal(t, []string{"banana"}, result)
	})

	t.Run("returns empty slice when no elements match predicate", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 3, 5, 7, 9}
		isEven := func(n int) bool {
			return n%2 == 0
		}

		// Act
		result := Filter(numbers, isEven)

		// Assert
		assert.Empty(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("returns all elements when all match predicate", func(t *testing.T) {
		// Arrange
		numbers := []int{2, 4, 6, 8, 10}
		isEven := func(n int) bool {
			return n%2 == 0
		}

		// Act
		result := Filter(numbers, isEven)

		// Assert
		assert.Equal(t, numbers, result)
	})

	t.Run("returns empty slice when input is empty", func(t *testing.T) {
		// Arrange
		var numbers []int
		isEven := func(n int) bool {
			return n%2 == 0
		}

		// Act
		result := Filter(numbers, isEven)

		// Assert
		assert.Empty(t, result)
		assert.Equal(t, 0, len(result))
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

		isAdult := func(p Person) bool {
			return p.Age >= 18
		}

		// Act
		result := Filter(people, isAdult)

		// Assert
		assert.Equal(t, []Person{{"Alice", 25}, {"Charlie", 30}}, result)
		assert.Equal(t, 2, len(result))
	})
}