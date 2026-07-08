package sdk

import (
	"strconv"
)

// String returns a pointer to the given string.
//
//go:fix inline
func String(s string) *string {
	return new(s)
}

// Bool returns a pointer to the given bool.
//
//go:fix inline
func Bool(b bool) *bool {
	return new(b)
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
//
//go:fix inline
func Int(i int) *int {
	return new(i)
}

// ToInt converts a string to an int.
func ToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func ToIntWithDefault(s string, defaultValue int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return i
}

// Float64 returns a pointer to the given float64.
//
//go:fix inline
func Float64(f float64) *float64 {
	return new(f)
}

// ToFloat64 converts a string to a float64.
func ToFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// Pointer is a generic function that returns a pointer to a given value.
//
//go:fix inline
func Pointer[K any](v K) *K {
	return new(v)
}
