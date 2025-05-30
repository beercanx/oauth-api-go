package collections

func Reduce[T any](array []T, combine func(T, T) T) T {
	if len(array) == 0 {
		var zero T
		return zero
	}

	result := array[0]
	for _, item := range array[1:] {
		result = combine(result, item)
	}
	return result
}
