package stringhelpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnakeCaseToCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"snake", "Snake"},
		{"snake_case", "SnakeCase"},
		{"s", "S"},
		{"", ""},
		{"multiple__underscores", "MultipleUnderscores"},
		{"_leading", "Leading"},
		{"trailing_", "Trailing_"},
		{"ALL_CAPS", "AllCaps"},
		{"with123_numbers", "With123Numbers"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SnakeCaseToCamel(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}
