package collections

func Filter[T any](array []T, predicate func(T) bool) []T {
	results := make([]T, 0)
	for _, item := range array {
		if predicate(item) {
			results = append(results, item)
		}
	}
	return results
}
