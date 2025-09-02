package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapResourceId(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{`"abc123"`, `snowflake_generated_abc123`},
		{`-abc123-`, `snowflake_generated_-abc123-`},
		{`_abc123_`, `snowflake_generated__abc123_`},
		{`"abc_-123"`, `snowflake_generated_abc_-123`},
		{`"1a1"."2b2"."3c3"`, `snowflake_generated_1a1_2b2_3c3`},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := NormalizeResourceId(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
