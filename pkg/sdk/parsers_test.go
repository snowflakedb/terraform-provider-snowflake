package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var commaSeparatedStringArrayTestCases = []struct {
	Name       string
	Value      string
	TrimQuotes bool
	Result     []string
}{
	{
		Name:   "empty list",
		Value:  "[]",
		Result: []string{},
	},
	{
		Name:   "empty string",
		Value:  "",
		Result: []string{},
	},
	{
		Name:   "one element in list",
		Value:  "[one]",
		Result: []string{"one"},
	},
	{
		Name:       "one element in list - with quotes",
		Value:      "['one']",
		TrimQuotes: true,
		Result:     []string{"one"},
	},
	{
		Name:       "multiple elements in list - with quotes",
		Value:      "['one', 'two', 'three']",
		TrimQuotes: true,
		Result:     []string{"one", "two", "three"},
	},
	{
		Name:   "multiple elements in list",
		Value:  "[one, two, three]",
		Result: []string{"one", "two", "three"},
	},
	{
		Name:   "multiple elements in list - packed",
		Value:  "[one,two,three]",
		Result: []string{"one", "two", "three"},
	},
	{
		Name:   "multiple elements in list - additional spaces",
		Value:  "[one    ,          two  ,three]",
		Result: []string{"one", "two", "three"},
	},
	{
		Name:   "list without brackets",
		Value:  "one,two,three",
		Result: []string{"one", "two", "three"},
	},
	{
		Name:       "list without brackets - with single quotes",
		Value:      "'one','two','three'",
		TrimQuotes: true,
		Result:     []string{"one", "two", "three"},
	},
	{
		Name:       "list without brackets - with double quotes",
		Value:      `"one","two","three"`,
		TrimQuotes: true,
		Result:     []string{"one", "two", "three"},
	},
	{
		Name:       "list with brackets - with double quotes",
		Value:      `"one","two","three"`,
		TrimQuotes: true,
		Result:     []string{"one", "two", "three"},
	},
	{
		Name:       "list with brackets - with double quotes, no trimming",
		Value:      `"one","two","three"`,
		TrimQuotes: false,
		Result:     []string{"\"one\"", "\"two\"", "\"three\""},
	},
	{
		Name:       "list with brackets - with double quotes",
		Value:      `["one","two","three"]`,
		TrimQuotes: true,
		Result:     []string{"one", "two", "three"},
	},
	{
		Name:       "multiple quote types",
		Value:      `['"'one'"',"'two",'"three'"]`,
		TrimQuotes: true,
		Result:     []string{"one", "two", "three"},
	},
}

func TestParseCommaSeparatedStringArray(t *testing.T) {
	for _, tc := range commaSeparatedStringArrayTestCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Result, ParseCommaSeparatedStringArray(tc.Value, tc.TrimQuotes))
		})
	}
}

func TestParseOuterCommaSeparatedStringArray(t *testing.T) {
	testCases := commaSeparatedStringArrayTestCases
	testCases = append(testCases, []struct {
		Name       string
		Value      string
		TrimQuotes bool
		Result     []string
	}{
		{
			Name:       "list with nested lists",
			Value:      "['one', [['two'], 'three', 'four'], ['five']]",
			TrimQuotes: true,
			Result:     []string{"one", "[['two'], 'three', 'four']", "['five']"},
		},
		{
			Name:       "list with single nested list",
			Value:      "[['one', 'two', 'three']]",
			TrimQuotes: true,
			Result:     []string{"['one', 'two', 'three']"},
		},
		{
			Name:       "opening bracket in the last element",
			Value:      "['one', 'two', 'thr[ee']",
			TrimQuotes: true,
			Result:     []string{"one", "two", "thr[ee"},
		},
		{
			Name:       "opening bracket in the middle element",
			Value:      "['one', 'tw[o', 'three']",
			TrimQuotes: true,
			Result:     []string{"one", "tw[o', 'three"},
		},
		{
			Name:       "closing brackets in the middle element",
			Value:      "['one', 'tw]]]o', 'three']",
			TrimQuotes: true,
			Result:     []string{"one", "tw]]]o", "three"},
		},
	}...)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			assert.Equal(t, tc.Result, ParseOuterCommaSeparatedStringArray(tc.Value, tc.TrimQuotes))
		})
	}
}

func TestParseCommaSeparatedSchemaObjectIdentifierArray(t *testing.T) {
	testCases := []struct {
		Name   string
		Value  string
		Result []SchemaObjectIdentifier
	}{
		{
			Name:   "empty list",
			Value:  "[]",
			Result: []SchemaObjectIdentifier{},
		},
		{
			Name:   "empty string",
			Value:  "",
			Result: []SchemaObjectIdentifier{},
		},
		{
			Name:   "one element in list",
			Value:  "[A.B.C]",
			Result: []SchemaObjectIdentifier{NewSchemaObjectIdentifier("A", "B", "C")},
		},
		{
			Name:   "one element in list - with mixed cases",
			Value:  `[A."b".C]`,
			Result: []SchemaObjectIdentifier{NewSchemaObjectIdentifier("A", "b", "C")},
		},
		{
			Name:   "multiple elements in list",
			Value:  "[A.B.C, D.E.F]",
			Result: []SchemaObjectIdentifier{NewSchemaObjectIdentifier("A", "B", "C"), NewSchemaObjectIdentifier("D", "E", "F")},
		},
		{
			Name:   "multiple elements in list - with mixed cases",
			Value:  `[A."b".C, "d"."e"."f"]`,
			Result: []SchemaObjectIdentifier{NewSchemaObjectIdentifier("A", "b", "C"), NewSchemaObjectIdentifier("d", "e", "f")},
		},
		{
			Name:   "multiple elements in list - packed",
			Value:  "[A.B.C,D.E.F]",
			Result: []SchemaObjectIdentifier{NewSchemaObjectIdentifier("A", "B", "C"), NewSchemaObjectIdentifier("D", "E", "F")},
		},
		{
			Name:   "multiple elements in list - additional spaces",
			Value:  "[A.B.C,     	 D.E.F]",
			Result: []SchemaObjectIdentifier{NewSchemaObjectIdentifier("A", "B", "C"), NewSchemaObjectIdentifier("D", "E", "F")},
		},
		{
			Name:   "list without brackets",
			Value:  "A.B.C, D.E.F",
			Result: []SchemaObjectIdentifier{NewSchemaObjectIdentifier("A", "B", "C"), NewSchemaObjectIdentifier("D", "E", "F")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ids, err := ParseCommaSeparatedSchemaObjectIdentifierArray(tc.Value)
			require.NoError(t, err)
			require.Equal(t, tc.Result, ids)
		})
	}
}

func TestParseCommaSeparatedSchemaObjectIdentifierArray_Invalid(t *testing.T) {
	testCases := []struct {
		Name  string
		Value string
		Error string
	}{
		{
			Name:  "bad quotes",
			Value: `["a]`,
			Error: "unable to read identifier: \"a, err = parse error on line 1, column 3: extraneous or missing \" in quoted-field",
		},
		{
			Name:  "missing parts",
			Value: "[a.b.c, a.b]",
			Error: "unexpected number of parts 2 in identifier a.b, expected 3 in a form of \"<database_name>.<schema_name>.<schema_object_name>\"",
		},
		{
			Name:  "too many parts",
			Value: "[a.b.c, a.b.c.d]",
			Error: "unexpected number of parts 4 in identifier a.b.c.d, expected 3 in a form of \"<database_name>.<schema_name>.<schema_object_name>\"",
		},
		{
			Name:  "missing parts - empty id",
			Value: "[a.b.c, ]",
			Error: "incompatible identifier",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := ParseCommaSeparatedSchemaObjectIdentifierArray(tc.Value)
			require.ErrorContains(t, err, tc.Error)
		})
	}
}

func TestParseCommaSeparatedAccountIdentifierArray(t *testing.T) {
	testCases := []struct {
		Name   string
		Value  string
		Result []AccountIdentifier
	}{
		{
			Name:   "empty list",
			Value:  "[]",
			Result: []AccountIdentifier{},
		},
		{
			Name:   "empty string",
			Value:  "",
			Result: []AccountIdentifier{},
		},
		{
			Name:   "one element in list",
			Value:  "[A.B]",
			Result: []AccountIdentifier{NewAccountIdentifier("A", "B")},
		},
		{
			Name:   "one element in list - with mixed cases",
			Value:  `[A."b"]`,
			Result: []AccountIdentifier{NewAccountIdentifier("A", "b")},
		},
		{
			Name:   "multiple elements in list",
			Value:  "[A.B, C.D]",
			Result: []AccountIdentifier{NewAccountIdentifier("A", "B"), NewAccountIdentifier("C", "D")},
		},
		{
			Name:   "multiple elements in list - with mixed cases",
			Value:  `[A."b", "c"."d"]`,
			Result: []AccountIdentifier{NewAccountIdentifier("A", "b"), NewAccountIdentifier("c", "d")},
		},
		{
			Name:   "multiple elements in list - packed",
			Value:  "[A.B,C.D]",
			Result: []AccountIdentifier{NewAccountIdentifier("A", "B"), NewAccountIdentifier("C", "D")},
		},
		{
			Name:   "multiple elements in list - additional spaces",
			Value:  "[A.B,     	 C.D]",
			Result: []AccountIdentifier{NewAccountIdentifier("A", "B"), NewAccountIdentifier("C", "D")},
		},
		{
			Name:   "list without brackets",
			Value:  "A.B, C.D",
			Result: []AccountIdentifier{NewAccountIdentifier("A", "B"), NewAccountIdentifier("C", "D")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ids, err := ParseCommaSeparatedAccountIdentifierArray(tc.Value)
			require.NoError(t, err)
			require.Equal(t, tc.Result, ids)
		})
	}
}

func TestParseCommaSeparatedAccountIdentifierArray_Invalid(t *testing.T) {
	testCases := []struct {
		Name  string
		Value string
		Error string
	}{
		{
			Name:  "invalid qoutes",
			Value: `["a.b]`,
			Error: "unable to read identifier: \"a.b, err = parse error on line 1, column 5: extraneous or missing \" in quoted-field",
		},
		{
			Name:  "missing parts",
			Value: "[a.b, a]",
			Error: "unexpected number of parts 1 in identifier a, expected 2 in a form of \"<organization_name>.<account_name>\"",
		},
		{
			Name:  "too many parts",
			Value: "[a.b, a.b.c]",
			Error: "unexpected number of parts 3 in identifier a.b.c, expected 2 in a form of \"<organization_name>.<account_name>\"",
		},
		{
			Name:  "missing parts - empty id",
			Value: "[a.b, ]",
			Error: "incompatible identifier",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := ParseCommaSeparatedAccountIdentifierArray(tc.Value)
			require.ErrorContains(t, err, tc.Error)
		})
	}
}
