package sdk

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommaSeparatedEnumMap(t *testing.T) {
	testCases := []struct {
		Name   string
		Value  string
		Result []string
	}{
		{
			Name:   "empty enum map",
			Value:  "{}",
			Result: []string{},
		},
		{
			Name:   "empty string",
			Value:  "",
			Result: []string{},
		},
		{
			Name:   "multiple elements",
			Value:  "{KEY=value, KEY2=value2}",
			Result: []string{"KEY=value", "KEY2=value2"},
		},
		{
			Name:   "multiple elements without curly braces",
			Value:  "KEY=value, KEY2=value2",
			Result: []string{"KEY=value", "KEY2=value2"},
		},
		{
			Name:   "multiple elements with nested arrays",
			Value:  "{KEY=value, KEY2=[INNER_KEY=value2, INNER_KEY2=value3]}",
			Result: []string{"KEY=value", "KEY2=[INNER_KEY=value2, INNER_KEY2=value3]"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actual := parseCommaSeparatedEnumMap(CatalogIntegrationProperty{Value: tc.Value})
			require.Equal(t, tc.Result, actual)
		})
	}
}
