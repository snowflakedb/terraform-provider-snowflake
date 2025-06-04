package sdk

import (
	"strconv"
)

// String returns a pointer to the given string.
func String(s string) *string {
	return &s
}

// Bool returns a pointer to the given bool.
func Bool(b bool) *bool {
	return &b
}

// ToBool converts a string to a bool.
func ToBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}
	return b
}

// Int returns a pointer to the given int.
func Int(i int) *int {
	return &i
}

// ToInt converts a string to an int.
func ToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

// Float64 returns a pointer to the given float64.
func Float64(f float64) *float64 {
	return &f
}

// ToFloat64 converts a string to a float64.
func ToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// Pointer is a generic function that returns a pointer to a given value.
func Pointer[K any](v K) *K {
	return &v
}
