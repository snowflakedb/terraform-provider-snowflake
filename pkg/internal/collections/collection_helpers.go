package collections

import (
	"errors"
	"strings"
)

var ErrObjectNotFound = errors.New("object does not exist")

func FindFirst[T any](collection []T, condition func(T) bool) (*T, error) {
	for _, o := range collection {
		if condition(o) {
			return &o, nil
		}
	}
	return nil, ErrObjectNotFound
}

func Map[T any, R any](collection []T, mapper func(T) R) []R {
	result := make([]R, len(collection))
	for i, elem := range collection {
		result[i] = mapper(elem)
	}
	return result
}

func MapErr[T any, R any](collection []T, mapper func(T) (R, error)) ([]R, error) {
	result := make([]R, len(collection))
	errs := make([]error, 0)
	for i, elem := range collection {
		value, err := mapper(elem)
		if err != nil {
			errs = append(errs, err)
		}
		result[i] = value
	}
	return result, errors.Join(errs...)
}

// TODO(SNOW-1479870): Test
// MergeMaps takes any number of maps (of the same type) and concatenates them.
// In case of key collision, the value will be selected from the map that is provided
// later in the src function parameter.
func MergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)
	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}

func JoinStrings[S ~string](stringCollection []S, separator string) string {
	mappedCollection := Map(stringCollection, func(stringValue S) string { return string(stringValue) })
	return strings.Join(mappedCollection, separator)
}

func GroupByProperty[T any, K comparable](items []T, getProperty func(T) K) map[K][]T {
	grouped := make(map[K][]T)

	for _, item := range items {
		key := getProperty(item)
		grouped[key] = append(grouped[key], item)
	}

	return grouped
}
