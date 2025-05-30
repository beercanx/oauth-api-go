package collections

func Map[T, R any](array []T, transform func(T) R) []R {
	results := make([]R, len(array))
	for index, item := range array {
		results[index] = transform(item)
	}
	return results
}
