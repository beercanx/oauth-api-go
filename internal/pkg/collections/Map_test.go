package collections

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("maps integers to their squares", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		square := func(n int) int {
			return n * n
		}

		// Act
		result := Map(numbers, square)

		// Assert
		assert.Equal(t, []int{1, 4, 9, 16, 25}, result)
	})

	t.Run("maps strings to their lengths", func(t *testing.T) {
		// Arrange
		words := []string{"apple", "banana", "cherry", "date"}
		length := func(s string) int {
			return len(s)
		}

		// Act
		result := Map(words, length)

		// Assert
		assert.Equal(t, []int{5, 6, 6, 4}, result)
	})

	t.Run("maps integers to strings", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		toString := func(n int) string {
			return strconv.Itoa(n)
		}

		// Act
		result := Map(numbers, toString)

		// Assert
		assert.Equal(t, []string{"1", "2", "3", "4", "5"}, result)
	})

	t.Run("maps strings to uppercase", func(t *testing.T) {
		// Arrange
		words := []string{"apple", "banana", "cherry", "date"}
		toUpper := func(s string) string {
			return strings.ToUpper(s)
		}

		// Act
		result := Map(words, toUpper)

		// Assert
		assert.Equal(t, []string{"APPLE", "BANANA", "CHERRY", "DATE"}, result)
	})

	t.Run("returns empty slice when input is empty", func(t *testing.T) {
		// Arrange
		var numbers []int
		square := func(n int) int {
			return n * n
		}

		// Act
		result := Map(numbers, square)

		// Assert
		assert.Empty(t, result)
		assert.Equal(t, 0, len(result))
	})

	t.Run("preserves the length of the input array", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		double := func(n int) int {
			return n * 2
		}

		// Act
		result := Map(numbers, double)

		// Assert
		assert.Equal(t, len(numbers), len(result))
	})

	t.Run("works with custom struct types", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}

		type PersonSummary struct {
			FullName string
			IsAdult  bool
		}

		people := []Person{
			{"Alice", 25},
			{"Bob", 17},
			{"Charlie", 30},
			{"David", 16},
		}

		toSummary := func(p Person) PersonSummary {
			return PersonSummary{
				FullName: p.Name,
				IsAdult:  p.Age >= 18,
			}
		}

		// Act
		result := Map(people, toSummary)

		// Assert
		expected := []PersonSummary{
			{"Alice", true},
			{"Bob", false},
			{"Charlie", true},
			{"David", false},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("maps to a completely different type", func(t *testing.T) {
		// Arrange
		type Item struct {
			ID    int
			Value string
		}

		items := []Item{
			{1, "one"},
			{2, "two"},
			{3, "three"},
		}

		toMap := func(item Item) map[int]string {
			return map[int]string{item.ID: item.Value}
		}

		// Act
		result := Map(items, toMap)

		// Assert
		assert.Equal(t, 3, len(result))
		assert.Equal(t, "one", result[0][1])
		assert.Equal(t, "two", result[1][2])
		assert.Equal(t, "three", result[2][3])
	})
}