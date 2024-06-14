package slice

func Get[T any](slice []T, index int, _default T) T {
	if index >= len(slice) || index < 0 {
		return _default
	}
	return slice[index]
}
