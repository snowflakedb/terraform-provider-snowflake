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
	{
		Name:       "comma inside double-quoted identifier",
		Value:      `TESTDB.TESTSCHEMA.IT1,TESTDB.TESTSCHEMA."I,T2"`,
		TrimQuotes: false,
		Result:     []string{`TESTDB.TESTSCHEMA.IT1`, `TESTDB.TESTSCHEMA."I,T2"`},
	},
	{
		Name:       "comma inside double-quoted identifier - multiple",
		Value:      `A."x,y",B."z,w"`,
		TrimQuotes: false,
		Result:     []string{`A."x,y"`, `B."z,w"`},
	},
	{
		// doubled quotes ("") escape a quote inside a quoted identifier; the comma inside stays protected
		Name:       "comma inside doubled-quoted identifier with escaped quote",
		Value:      `A.B."a""b,c",D.E.F`,
		TrimQuotes: false,
		Result:     []string{`A.B."a""b,c"`, `D.E.F`},
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
		{
			Name:       "comma inside double-quoted value alongside a nested list",
			Value:      `["a,b", [1, 2]]`,
			TrimQuotes: true,
			Result:     []string{"a,b", "[1, 2]"},
		},
		{
			Name:       "opening bracket inside a double-quoted value does not increase depth",
			Value:      `["a[b", "c"]`,
			TrimQuotes: true,
			Result:     []string{"a[b", "c"},
		},
		{
			Name:       "closing bracket inside a double-quoted value does not decrease depth",
			Value:      `["a]b", "c,d"]`,
			TrimQuotes: true,
			Result:     []string{"a]b", "c,d"},
		},
		{
			Name:       "comma and brackets both inside a double-quoted value",
			Value:      `["a,[b]", "c"]`,
			TrimQuotes: true,
			Result:     []string{"a,[b]", "c"},
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

func TestParseDescribeSignature(t *testing.T) {
	testCases := []struct {
		Name      string
		Signature string
		Result    []NormalizedArgument
	}{
		{
			Name:      "no arguments",
			Signature: "()",
			Result:    []NormalizedArgument{},
		},
		{
			Name:      "empty string",
			Signature: "",
			Result:    []NormalizedArgument{},
		},
		{
			Name:      "single argument",
			Signature: "(A VARCHAR)",
			Result:    []NormalizedArgument{{"A", dataTypeVarchar}},
		},
		{
			Name:      "without wrapping parentheses",
			Signature: "A VARCHAR, my col NUMBER, V VECTOR(FLOAT, 10)",
			Result: []NormalizedArgument{
				{"A", dataTypeVarchar},
				{"my col", dataTypeNumber},
				{"V", dataTypeVectorFloat_10},
			},
		},
		{
			Name:      "multiple arguments",
			Signature: "(A VARCHAR, B BOOLEAN)",
			Result:    []NormalizedArgument{{"A", dataTypeVarchar}, {"B", dataTypeBoolean}},
		},
		{
			Name:      "name with spaces",
			Signature: "(my col VARCHAR)",
			Result:    []NormalizedArgument{{"my col", dataTypeVarchar}},
		},
		{
			Name:      "empty name",
			Signature: "( VARCHAR)",
			Result:    []NormalizedArgument{{"", dataTypeVarchar}},
		},
		{
			Name:      "empty name - not the first argument",
			Signature: "(A VARCHAR,  BOOLEAN)",
			Result:    []NormalizedArgument{{"A", dataTypeVarchar}, {"", dataTypeBoolean}},
		},
		{
			Name:      "whitespace-only names",
			Signature: "(  VARCHAR,    BOOLEAN)",
			Result:    []NormalizedArgument{{" ", dataTypeVarchar}, {"  ", dataTypeBoolean}},
		},
		{
			Name:      "name with a leading space",
			Signature: "( a VARCHAR,  b BOOLEAN)",
			Result:    []NormalizedArgument{{" a", dataTypeVarchar}, {" b", dataTypeBoolean}},
		},
		{
			Name:      "name with a trailing space",
			Signature: "(a  VARCHAR, b  BOOLEAN)",
			Result:    []NormalizedArgument{{"a ", dataTypeVarchar}, {"b ", dataTypeBoolean}},
		},
		{
			Name:      "name with lowercase letters",
			Signature: "(lower NUMBER)",
			Result:    []NormalizedArgument{{"lower", dataTypeNumber}},
		},
		{
			Name:      "name with a double quote",
			Signature: `(x"y VARCHAR, B BOOLEAN)`,
			Result:    []NormalizedArgument{{`x"y`, dataTypeVarchar}, {"B", dataTypeBoolean}},
		},
		{
			Name:      "name DEFAULT",
			Signature: "(DEFAULT NUMBER)",
			Result:    []NormalizedArgument{{"DEFAULT", dataTypeNumber}},
		},
		{
			Name:      "vector",
			Signature: "(V VECTOR(FLOAT, 10))",
			Result:    []NormalizedArgument{{"V", dataTypeVectorFloat_10}},
		},
		{
			Name:      "vector with a name with spaces",
			Signature: "(my vec VECTOR(FLOAT, 10), A VARCHAR)",
			Result:    []NormalizedArgument{{"my vec", dataTypeVectorFloat_10}, {"A", dataTypeVarchar}},
		},
		{
			Name:      "structured types",
			Signature: "(O OBJECT(a VARCHAR, b NUMBER NOT NULL), R ARRAY(NUMBER), M MAP(VARCHAR, NUMBER))",
			Result: []NormalizedArgument{
				{"O", dataTypeObjectStructured},
				{"R", dataTypeArrayNumber},
				{"M", dataTypeMapVarcharNumber},
			},
		},
		{
			Name:      "structured types with names with spaces",
			Signature: "(my obj OBJECT(a VARCHAR, b NUMBER NOT NULL), my arr ARRAY(NUMBER), my map MAP(VARCHAR, NUMBER))",
			Result: []NormalizedArgument{
				{"my obj", dataTypeObjectStructured},
				{"my arr", dataTypeArrayNumber},
				{"my map", dataTypeMapVarcharNumber},
			},
		},
		{
			Name:      "function signature returned by Snowflake",
			Signature: `(my col VARCHAR, lower NUMBER, DEFAULT NUMBER, x"y VARCHAR, V VECTOR(INT, 10), D VARCHAR)`,
			Result: []NormalizedArgument{
				{"my col", dataTypeVarchar},
				{"lower", dataTypeNumber},
				{"DEFAULT", dataTypeNumber},
				{`x"y`, dataTypeVarchar},
				{"V", dataTypeVectorInt_10},
				{"D", dataTypeVarchar},
			},
		},
		{
			Name:      "policy signature returned by Snowflake",
			Signature: `(my col VARCHAR, lower NUMBER, PLAIN VARCHAR, DEFAULT NUMBER, x"y VARCHAR, V VECTOR(FLOAT, 10))`,
			Result: []NormalizedArgument{
				{"my col", dataTypeVarchar},
				{"lower", dataTypeNumber},
				{"PLAIN", dataTypeVarchar},
				{"DEFAULT", dataTypeNumber},
				{`x"y`, dataTypeVarchar},
				{"V", dataTypeVectorFloat_10},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := ParseDescribeSignature(tc.Signature)
			require.NoError(t, err)
			require.Equal(t, tc.Result, result)
		})
	}
}

func TestParseDescribeSignature_Invalid(t *testing.T) {
	testCases := []struct {
		Name      string
		Signature string
		Error     string
	}{
		{
			Name:      "missing argument name",
			Signature: "(VARCHAR)",
			Error:     "expected argument name and type, got VARCHAR",
		},
		{
			Name:      "missing argument name - one of many",
			Signature: "(A VARCHAR, BOOLEAN)",
			Error:     "expected argument name and type, got BOOLEAN",
		},
		{
			Name:      "empty argument",
			Signature: "(A VARCHAR, )",
			Error:     "empty argument",
		},
		{
			Name:      "empty argument - in the middle",
			Signature: "(A VARCHAR, , B BOOLEAN)",
			Error:     "empty argument",
		},
		{
			Name:      "invalid data type",
			Signature: "(A INVALID_TYPE)",
			Error:     "invalid data type",
		},
		{
			Name:      "missing data type",
			Signature: "(A )",
			Error:     "invalid data type",
		},
		{
			Name:      "unsupported name with a comma",
			Signature: "(a,b VARCHAR)",
			Error:     "expected argument name and type, got a",
		},
		{
			Name:      "unsupported name with a parenthesis",
			Signature: "(p(x VARCHAR)",
			Error:     "expected argument name and type, got p(x VARCHAR",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			_, err := ParseDescribeSignature(tc.Signature)
			require.ErrorContains(t, err, tc.Error)
		})
	}
}
