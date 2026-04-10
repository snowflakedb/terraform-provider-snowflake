package collections

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func Test_FindFirst(t *testing.T) {
	stringSlice := []string{"1", "22", "333", "334"}

	t.Run("basic find", func(t *testing.T) {
		result, resultErr := FindFirst(stringSlice, func(s string) bool { return s == "22" })

		require.Equal(t, "22", *result)
		require.NoError(t, resultErr)
	})

	t.Run("two matching, first returned", func(t *testing.T) {
		result, resultErr := FindFirst(stringSlice, func(s string) bool { return strings.HasPrefix(s, "33") })

		require.Equal(t, "333", *result)
		require.NoError(t, resultErr)
	})

	t.Run("no item", func(t *testing.T) {
		result, resultErr := FindFirst(stringSlice, func(s string) bool { return s == "4444" })

		require.Nil(t, result)
		require.ErrorIs(t, resultErr, ErrObjectNotFound)
	})
}

func Test_Map(t *testing.T) {
	t.Run("basic mapping", func(t *testing.T) {
		stringSlice := []string{"1", "22", "333"}
		stringLenSlice := Map(stringSlice, func(s string) int { return len(s) })
		require.Equal(t, []int{1, 2, 3}, stringLenSlice)
	})

	t.Run("validation: empty slice", func(t *testing.T) {
		stringSlice := make([]string, 0)
		stringLenSlice := Map(stringSlice, func(s string) int { return len(s) })
		require.Equal(t, []int{}, stringLenSlice)
	})

	t.Run("validation: nil slice", func(t *testing.T) {
		var stringSlice []string = nil
		stringLenSlice := Map(stringSlice, func(s string) int { return len(s) })
		require.Equal(t, []int{}, stringLenSlice)
	})

	t.Run("validation: nil mapping function", func(t *testing.T) {
		require.PanicsWithError(t, "runtime error: invalid memory address or nil pointer dereference", func() {
			stringSlice := []string{"1", "22", "333"}
			_ = Map[string, int](stringSlice, nil)
		})
	})
}

func Test_Filter(t *testing.T) {
	stringSlice := []string{"1", "22", "333"}

	t.Run("all matches", func(t *testing.T) {
		allMatches := func(s string) bool { return true }
		require.Equal(t, stringSlice, Filter(stringSlice, allMatches))
	})

	t.Run("no matches", func(t *testing.T) {
		noMatches := func(s string) bool { return false }
		require.Equal(t, []string{}, Filter(stringSlice, noMatches))
	})

	t.Run("some matches", func(t *testing.T) {
		someMatches := func(s string) bool { return s == "22" }
		require.Equal(t, []string{"22"}, Filter(stringSlice, someMatches))
	})
}

func Test_MapErr(t *testing.T) {
	t.Run("basic mapping", func(t *testing.T) {
		stringSlice := []string{"1", "22", "333"}
		stringLenSlice, err := MapErr(stringSlice, func(s string) (int, error) { return len(s), nil })
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, stringLenSlice)
	})

	t.Run("basic mapping - multiple errors", func(t *testing.T) {
		stringSlice := []string{"1", "22", "333"}
		stringLenSlice, err := MapErr(stringSlice, func(s string) (int, error) {
			if s == "1" {
				return -1, fmt.Errorf("error: 1")
			}
			if s == "22" {
				return -1, fmt.Errorf("error: 22")
			}
			return len(s), nil
		})
		assert.Equal(t, []int{-1, -1, 3}, stringLenSlice)
		assert.ErrorContains(t, err, errors.Join(fmt.Errorf("error: 1"), fmt.Errorf("error: 22")).Error())
	})

	t.Run("validation: empty slice", func(t *testing.T) {
		stringSlice := make([]string, 0)
		stringLenSlice, err := MapErr(stringSlice, func(s string) (int, error) { return len(s), nil })
		assert.NoError(t, err)
		assert.Equal(t, []int{}, stringLenSlice)
	})

	t.Run("validation: nil slice", func(t *testing.T) {
		var stringSlice []string = nil
		stringLenSlice, err := MapErr(stringSlice, func(s string) (int, error) { return len(s), nil })
		assert.NoError(t, err)
		assert.Equal(t, []int{}, stringLenSlice)
	})

	t.Run("validation: nil mapping function", func(t *testing.T) {
		assert.PanicsWithError(t, "runtime error: invalid memory address or nil pointer dereference", func() {
			stringSlice := []string{"1", "22", "333"}
			_, _ = MapErr[string, int](stringSlice, nil)
		})
	})
}

func Test_GroupByProperty(t *testing.T) {
	type Item struct {
		Name     string
		Category string
		Number   int
	}

	t.Run("validation: empty list", func(t *testing.T) {
		var items []Item
		groups := GroupByProperty(items, func(item Item) string {
			return item.Category
		})
		require.Empty(t, groups)
	})

	t.Run("basic grouping", func(t *testing.T) {
		items := []Item{
			{Name: "Item1", Category: "A", Number: 1},
			{Name: "Item2", Category: "B", Number: 2},
			{Name: "Item3", Category: "B", Number: 3},
			{Name: "Item4", Category: "A", Number: 4},
		}

		groups := GroupByProperty(items, func(item Item) string {
			return item.Category
		})
		require.Len(t, groups, 2)

		assert.Len(t, groups["A"], 2)
		assert.Contains(t, groups["A"], items[0])
		assert.Contains(t, groups["A"], items[3])

		assert.Len(t, groups["B"], 2)
		assert.Contains(t, groups["B"], items[1])
		assert.Contains(t, groups["B"], items[2])
	})

	t.Run("multi property grouping", func(t *testing.T) {
		items := []Item{
			{Name: "Item1", Category: "A", Number: 1},
			{Name: "Item2", Category: "B", Number: 2},
			{Name: "Item3", Category: "B", Number: 3},
			{Name: "Item4", Category: "A", Number: 4},
			{Name: "Item5", Category: "A", Number: 1},
		}

		groups := GroupByProperty(items, func(item Item) string {
			return fmt.Sprintf("%s_%d", item.Category, item.Number)
		})
		require.Len(t, groups, 4)

		assert.Len(t, groups["A_1"], 2)
		assert.Contains(t, groups["A_1"], items[0])
		assert.Contains(t, groups["A_1"], items[4])
	})
}

func Test_CommonPrefixLastIndex(t *testing.T) {
	testCases := []struct {
		name     string
		a        []int
		b        []int
		expected int
	}{
		{name: "nil slices", a: nil, b: nil, expected: -1},
		{name: "first nil second non-empty", a: nil, b: []int{1}, expected: -1},
		{name: "two empty lists", a: []int{}, b: []int{}, expected: -1},
		{name: "first list empty", a: []int{}, b: []int{1}, expected: -1},
		{name: "second list empty", a: []int{1}, b: []int{}, expected: -1},
		{name: "no common prefix - length 1", a: []int{1}, b: []int{2}, expected: -1},
		{name: "no common prefix - length 2", a: []int{1, 2}, b: []int{3, 4}, expected: -1},
		{name: "identical lists - length 1", a: []int{1}, b: []int{1}, expected: 0},
		{name: "identical lists - length 2", a: []int{1, 2}, b: []int{1, 2}, expected: 1},
		{name: "common prefix up to index 1 out of 3", a: []int{1, 2, 3}, b: []int{1, 2, 4}, expected: 1},
		{name: "common prefix up to index 2 out of 4", a: []int{1, 2, 3, 4}, b: []int{1, 2, 3, 5}, expected: 2},
		{name: "common prefix up to index 1 out of 4", a: []int{1, 2, 3, 4}, b: []int{1, 2, 5, 4}, expected: 1},
		{name: "different lengths - matching up to last index of shorter list", a: []int{1, 2, 3}, b: []int{1, 2}, expected: 1},
		{name: "different lengths - matching up to index 2 out of longer list", a: []int{1, 2, 3, 4, 5, 6}, b: []int{1, 2, 3, 7, 8}, expected: 2},
	}

	intEqual := func(a, b int) bool { return a == b }
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, CommonPrefixLastIndex(tc.a, tc.b, intEqual))
		})
	}

	t.Run("custom comparator - compare by struct field", func(t *testing.T) {
		type item struct {
			key   string
			value int
		}
		cmpByKey := func(a, b item) bool { return a.key == b.key }

		a := []item{{key: "a", value: 1}, {key: "b", value: 2}, {key: "c", value: 3}}
		b := []item{{key: "a", value: 99}, {key: "b", value: 88}, {key: "d", value: 3}}

		require.Equal(t, 1, CommonPrefixLastIndex(a, b, cmpByKey))
	})
}
