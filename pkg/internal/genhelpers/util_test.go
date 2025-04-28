package genhelpers

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToSnakeCase(t *testing.T) {
	type test struct {
		input    string
		expected string
	}

	testCases := []test{
		{input: "CamelCase", expected: "camel_case"},
		{input: "ACamelCase", expected: "a_camel_case"},
		{input: "URLParser", expected: "url_parser"},
		{input: "Camel1Case", expected: "camel1_case"},
		{input: "camelCase", expected: "camel_case"},
		{input: "camelURL", expected: "camel_url"},
		{input: "camelURLSomething", expected: "camel_url_something"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s=>%s", tc.input, tc.expected), func(t *testing.T) {
			result := ToSnakeCase(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}

func Test_ColumnOutput(t *testing.T) {
	spaces := func(count int) string {
		return strings.Repeat(" ", count)
	}
	chars := func(count int) string {
		return strings.Repeat("a", count)
	}
	ten := chars(10)
	twenty := chars(20)

	type test struct {
		name           string
		columnWidth    int
		columns        []string
		expectedOutput string
	}

	testCases := []test{
		{name: "no columns", columnWidth: 16, columns: []string{}, expectedOutput: ""},
		{name: "one column, shorter than width", columnWidth: 16, columns: []string{ten}, expectedOutput: ten},
		{name: "one column, longer than width", columnWidth: 16, columns: []string{twenty}, expectedOutput: twenty},
		{name: "two column, shorter than width", columnWidth: 16, columns: []string{ten, ten}, expectedOutput: ten + spaces(6) + ten},
		{name: "two column, longer than width", columnWidth: 16, columns: []string{twenty, ten}, expectedOutput: twenty + spaces(1) + ten},
		{name: "zero width", columnWidth: 0, columns: []string{ten, ten}, expectedOutput: ten + spaces(1) + ten},
		{name: "negative width", columnWidth: -10, columns: []string{ten, ten}, expectedOutput: ten + spaces(1) + ten},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s - column width [%d]", tc.name, tc.columnWidth), func(t *testing.T) {
			result := ColumnOutput(tc.columnWidth, tc.columns...)
			require.Equal(t, tc.expectedOutput, result)
		})
	}
}

func Test_SanitizeAttributeName(t *testing.T) {
	type test struct {
		input    string
		expected string
	}

	testCases := []test{
		{input: "type", expected: "type" + forbiddenAttributeNameSuffix},
		{input: "abc", expected: "abc"},
		{input: "type_abc", expected: "type_abc"},
		{input: "abc_type", expected: "abc_type"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s=>%s", tc.input, tc.expected), func(t *testing.T) {
			result := SanitizeAttributeName(tc.input)
			require.Equal(t, tc.expected, result)
		})
	}
}
