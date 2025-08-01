package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldRotateToken(t *testing.T) {
	tests := []struct {
		name     string
		old      string
		new      string
		isKnown  bool
		expected bool
	}{
		// Cases where old is empty
		{
			name:     "old empty, new empty, known",
			old:      "",
			new:      "",
			isKnown:  true,
			expected: false,
		},
		{
			name:     "old empty, new empty, unknown",
			old:      "",
			new:      "",
			isKnown:  false,
			expected: false,
		},
		// Cases where old is empty and the value was added to the config.
		{
			name:     "old empty, new non-empty, known",
			old:      "",
			new:      "new_value",
			isKnown:  true,
			expected: false,
		},
		{
			name:     "old empty, new non-empty, unknown",
			old:      "",
			new:      "new_value",
			isKnown:  false,
			expected: false,
		},

		// Cases where old is non-empty and new is empty (the value is removed from the config)
		{
			name:     "old non-empty, new empty, known",
			old:      "old_value",
			new:      "",
			isKnown:  true,
			expected: false,
		},
		{
			name:     "old non-empty, new empty, unknown",
			old:      "old_value",
			new:      "",
			isKnown:  false,
			expected: true,
		},

		// Cases where old and new are the same non-empty values
		{
			name:     "old and new same non-empty, known",
			old:      "same_value",
			new:      "same_value",
			isKnown:  true,
			expected: false,
		},
		{
			name:     "old and new same non-empty, unknown",
			old:      "same_value",
			new:      "same_value",
			isKnown:  false,
			expected: true,
		},

		// Cases where old and new are different non-empty values
		{
			name:     "old and new different non-empty, known",
			old:      "old_value",
			new:      "new_value",
			isKnown:  true,
			expected: true,
		},
		{
			name:     "old and new different non-empty, unknown",
			old:      "old_value",
			new:      "new_value",
			isKnown:  false,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldRotateToken(tt.old, tt.new, tt.isKnown)
			assert.Equal(t, tt.expected, result)
		})
	}
}
