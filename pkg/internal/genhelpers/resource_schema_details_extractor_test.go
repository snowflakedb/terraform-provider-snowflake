package genhelpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_OrderByNameWithPreferredFront(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "all identity fields scrambled with other attributes",
			input:    []string{"min_nodes", "name", "comment", "schema", "database"},
			expected: []string{"database", "schema", "name", "comment", "min_nodes"},
		},
		{
			name:     "same names already sorted still put identity first",
			input:    []string{"comment", "database", "min_nodes", "name", "schema"},
			expected: []string{"database", "schema", "name", "comment", "min_nodes"},
		},
		{
			name:     "only database and name",
			input:    []string{"name", "database"},
			expected: []string{"database", "name"},
		},
		{
			name:     "database and name with other attributes",
			input:    []string{"comment", "name", "zzz", "database"},
			expected: []string{"database", "name", "comment", "zzz"},
		},
		{
			name:     "only name",
			input:    []string{"name"},
			expected: []string{"name"},
		},
		{
			name:     "only name with other attributes",
			input:    []string{"min_nodes", "comment", "name"},
			expected: []string{"name", "comment", "min_nodes"},
		},
		{
			name:     "no identity fields",
			input:    []string{"b", "a"},
			expected: []string{"a", "b"},
		},
		{
			name:     "empty",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, OrderByNameWithPreferredFront(tt.input, func(name string) string { return name }))
		})
	}
}
