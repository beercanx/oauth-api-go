package collections

func Fold[T, R any](array []T, initial R, combine func(R, T) R) R {
	for _, item := range array {
		initial = combine(initial, item)
	}
	return initial
}
