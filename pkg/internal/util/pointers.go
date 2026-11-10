package util

// DerefOrZero returns the value behind ptr, or the zero value of T when ptr is nil.
func DerefOrZero[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}
