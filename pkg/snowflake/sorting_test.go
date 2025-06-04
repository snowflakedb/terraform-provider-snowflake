package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSortStrings(t *testing.T) {
	r := require.New(t)

	ss := []string{"a", "b", "c"}

	r.Equal(ss, sortStrings(map[string]string{"c": "", "b": "", "a": ""}))
	r.Equal(ss, sortStringList(map[string][]string{"c": {}, "b": {}, "a": {}}))
	r.Equal(ss, sortStringsInt(map[string]int{"c": 0, "b": 1, "a": 2}))
	r.Equal(ss, sortStringsFloat(map[string]float64{"c": 0, "b": 1, "a": 2}))
	r.Equal(ss, sortStringsBool(map[string]bool{"c": true, "b": false, "a": true}))
}
