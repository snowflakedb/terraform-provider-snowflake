package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsvUnescape(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no escaping needed",
			input:    "simple text",
			expected: "simple text",
		},
		{
			name:     "escaped newline",
			input:    "line1\\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "escaped carriage return",
			input:    "line1\\rline2",
			expected: "line1\rline2",
		},
		{
			name:     "escaped backslash",
			input:    "path\\\\to\\\\file",
			expected: "path\\to\\file",
		},
		{
			name:     "multiple newlines",
			input:    "line1\\nline2\\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "RSA key format",
			input:    "-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A\\nMIIBCgKCAQEA...\\n-----END PUBLIC KEY-----",
			expected: "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A\nMIIBCgKCAQEA...\n-----END PUBLIC KEY-----",
		},
		{
			name:     "mixed escape sequences",
			input:    "text\\nwith\\nnewlines\\rand\\\\backslash",
			expected: "text\nwith\nnewlines\rand\\backslash",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "backslash followed by n (not newline)",
			input:    "\\\\n",
			expected: "\\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := csvUnescape(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}
