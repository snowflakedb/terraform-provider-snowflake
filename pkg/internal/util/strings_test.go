package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_TrimAllPrefixes(t *testing.T) {
	type test struct {
		input    string
		prefixes []string
		expected string
	}

	tests := []test{
		{input: "VARCHAR(30)", prefixes: []string{"VARCHAR", "TEXT"}, expected: "(30)"},
		{input: "VARCHAR  (30) ", prefixes: []string{"VARCHAR", "TEXT"}, expected: "  (30) "},
		{input: "VARCHAR(30)", prefixes: []string{"VARCHAR"}, expected: "(30)"},
		{input: "VARCHAR(30)", prefixes: []string{}, expected: "VARCHAR(30)"},
		{input: "VARCHARVARCHAR(30)", prefixes: []string{"VARCHAR"}, expected: "VARCHAR(30)"},
		{input: "VARCHAR(30)", prefixes: []string{"NUMBER"}, expected: "VARCHAR(30)"},
		{input: "VARCHARTEXT(30)", prefixes: []string{"VARCHAR", "TEXT"}, expected: "(30)"},
		{input: "TEXTVARCHAR(30)", prefixes: []string{"VARCHAR", "TEXT"}, expected: "VARCHAR(30)"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			output := TrimAllPrefixes(tc.input, tc.prefixes...)
			require.Equal(t, tc.expected, output)
		})
	}
}
