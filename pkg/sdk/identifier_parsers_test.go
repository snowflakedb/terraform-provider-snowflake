package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseIdentifierString(t *testing.T) {
	containsAll := func(t *testing.T, parts, expectedParts []string) {
		t.Helper()
		require.Len(t, parts, len(expectedParts))
		for _, part := range expectedParts {
			require.Contains(t, parts, part)
		}
	}

	t.Run("returns read error", func(t *testing.T) {
		input := `ab"c`

		_, err := ParseIdentifierString(input)

		require.ErrorContains(t, err, "unable to read identifier")
		require.ErrorContains(t, err, `bare " in non-quoted-field`)
	})

	t.Run("returns error for empty input", func(t *testing.T) {
		input := ""

		_, err := ParseIdentifierString(input)

		require.ErrorContains(t, err, "incompatible identifier")
	})

	t.Run("returns error for multiple lines", func(t *testing.T) {
		input := "abc\ndef"

		_, err := ParseIdentifierString(input)

		require.ErrorContains(t, err, "incompatible identifier")
	})

	t.Run("returns parts correctly without quoting", func(t *testing.T) {
		input := "abc.def"
		expected := []string{"abc", "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly with quoting", func(t *testing.T) {
		input := `"abc"."def"`
		expected := []string{"abc", "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly with mixed quoting", func(t *testing.T) {
		input := `"abc".def."ghi"`
		expected := []string{"abc", "def", "ghi"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	// Quote inside must have a preceding quote (https://docs.snowflake.com/en/sql-reference/identifiers-syntax).
	t.Run("returns parts correctly with quote inside", func(t *testing.T) {
		input := `"ab""c".def`
		_, err := ParseIdentifierString(input)

		require.ErrorContains(t, err, `unable to parse identifier: "ab""c".def, currently identifiers containing double quotes are not supported in the provider`)
	})

	t.Run("returns parts correctly when identifier contains opening parenthesis", func(t *testing.T) {
		input := `"ab(c".def`
		expected := []string{"ab(c", "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly when identifier contains closing parenthesis", func(t *testing.T) {
		input := `"ab)c".def`
		expected := []string{"ab)c", "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly when identifier contains opening and closing parentheses", func(t *testing.T) {
		input := `"ab()c".def`
		expected := []string{"ab()c", "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("returns parts correctly with dots inside", func(t *testing.T) {
		input := `"ab.c".def`
		expected := []string{`ab.c`, "def"}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("empty identifier", func(t *testing.T) {
		input := `""`
		expected := []string{""}

		parts, err := ParseIdentifierString(input)

		require.NoError(t, err)
		containsAll(t, parts, expected)
	})

	t.Run("handled correctly double quotes", func(t *testing.T) {
		input := `""."."".".".""."".""."".""."".""."""""`

		_, err := ParseIdentifierString(input)
		require.ErrorContains(t, err, `unable to parse identifier: "".".""."."."".""."".""."".""."".""""", currently identifiers containing double quotes are not supported in the provider`)
	})
}

func Test_IdentifierParsers(t *testing.T) {
	testCases := []struct {
		IdentifierType string
		Input          string
		Expected       ObjectIdentifier
		Error          string
	}{
		{IdentifierType: "AccountObjectIdentifier", Input: ``, Expected: NewAccountObjectIdentifier("")},
		{IdentifierType: "AccountObjectIdentifier", Input: "a\nb", Expected: NewAccountObjectIdentifier("a\nb")},
		{IdentifierType: "AccountObjectIdentifier", Input: `a"b`, Error: `unable to read identifier: "a"b", err = parse error on line 1, column 3: extraneous or missing " in quoted-field`},
		{IdentifierType: "AccountObjectIdentifier", Input: `abc.cde`, Expected: NewAccountObjectIdentifier("abc.cde")},
		{IdentifierType: "AccountObjectIdentifier", Input: `""""`, Error: `unable to parse identifier: """", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "AccountObjectIdentifier", Input: `"a""bc"`, Error: `unable to parse identifier: "a""bc", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "AccountObjectIdentifier", Input: `""`, Expected: NewAccountObjectIdentifier(``)},
		{IdentifierType: "AccountObjectIdentifier", Input: `abc`, Expected: NewAccountObjectIdentifier(`abc`)},
		{IdentifierType: "AccountObjectIdentifier", Input: `"abc"`, Expected: NewAccountObjectIdentifier(`abc`)},
		{IdentifierType: "AccountObjectIdentifier", Input: `"ab.c"`, Expected: NewAccountObjectIdentifier(`ab.c`)},

		{IdentifierType: "DatabaseObjectIdentifier", Input: ``, Error: "incompatible identifier: "},
		{IdentifierType: "DatabaseObjectIdentifier", Input: "a\nb.cde", Error: "unable to read identifier: a\nb.cde, err = record on line 2: wrong number of fields"},
		{IdentifierType: "DatabaseObjectIdentifier", Input: `a"b.cde`, Error: "unable to read identifier: a\"b.cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{IdentifierType: "DatabaseObjectIdentifier", Input: `abc.cde.efg`, Error: `unexpected number of parts 3 in identifier abc.cde.efg, expected 2 in a form of "<database_name>.<database_object_name>"`},
		{IdentifierType: "DatabaseObjectIdentifier", Input: `abc`, Error: `unexpected number of parts 1 in identifier abc, expected 2 in a form of "<database_name>.<database_object_name>"`},
		{IdentifierType: "DatabaseObjectIdentifier", Input: `"""".""""`, Error: `unable to parse identifier: """"."""", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "DatabaseObjectIdentifier", Input: `"a""bc"."cd""e"`, Error: `unable to parse identifier: "a""bc"."cd""e", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "DatabaseObjectIdentifier", Input: `"".""`, Expected: NewDatabaseObjectIdentifier(``, ``)},
		{IdentifierType: "DatabaseObjectIdentifier", Input: `abc.cde`, Expected: NewDatabaseObjectIdentifier(`abc`, `cde`)},
		{IdentifierType: "DatabaseObjectIdentifier", Input: `"abc"."cde"`, Expected: NewDatabaseObjectIdentifier(`abc`, `cde`)},
		{IdentifierType: "DatabaseObjectIdentifier", Input: `"ab.c"."cd.e"`, Expected: NewDatabaseObjectIdentifier(`ab.c`, `cd.e`)},

		{IdentifierType: "SchemaObjectIdentifier", Input: ``, Error: "incompatible identifier: "},
		{IdentifierType: "SchemaObjectIdentifier", Input: "a\nb.cde.efg", Error: "unable to read identifier: a\nb.cde.efg, err = record on line 2: wrong number of fields"},
		{IdentifierType: "SchemaObjectIdentifier", Input: `a"b.cde.efg`, Error: "unable to read identifier: a\"b.cde.efg, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{IdentifierType: "SchemaObjectIdentifier", Input: `abc.cde.efg.ghi`, Error: `unexpected number of parts 4 in identifier abc.cde.efg.ghi, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`},
		{IdentifierType: "SchemaObjectIdentifier", Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 3 in a form of "<database_name>.<schema_name>.<schema_object_name>"`},
		{IdentifierType: "SchemaObjectIdentifier", Input: `""""."""".""""`, Error: `unable to parse identifier: """".""""."""", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "SchemaObjectIdentifier", Input: `"a""bc"."cd""e"."ef""g"`, Error: `unable to parse identifier: "a""bc"."cd""e"."ef""g", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "SchemaObjectIdentifier", Input: `""."".""`, Expected: NewSchemaObjectIdentifier(``, ``, ``)},
		{IdentifierType: "SchemaObjectIdentifier", Input: `abc.cde.efg`, Expected: NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
		{IdentifierType: "SchemaObjectIdentifier", Input: `"abc"."cde"."efg"`, Expected: NewSchemaObjectIdentifier(`abc`, `cde`, `efg`)},
		{IdentifierType: "SchemaObjectIdentifier", Input: `"ab.c"."cd.e"."ef.g"`, Expected: NewSchemaObjectIdentifier(`ab.c`, `cd.e`, `ef.g`)},

		{IdentifierType: "TableColumnIdentifier", Input: ``, Error: "incompatible identifier: "},
		{IdentifierType: "TableColumnIdentifier", Input: "a\nb.cde.efg.ghi", Error: "unable to read identifier: a\nb.cde.efg.ghi, err = record on line 2: wrong number of fields"},
		{IdentifierType: "TableColumnIdentifier", Input: `a"b.cde.efg.ghi`, Error: "unable to read identifier: a\"b.cde.efg.ghi, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{IdentifierType: "TableColumnIdentifier", Input: `abc.cde.efg.ghi.ijk`, Error: `unexpected number of parts 5 in identifier abc.cde.efg.ghi.ijk, expected 4 in a form of "<database_name>.<schema_name>.<table_name>.<table_column_name>"`},
		{IdentifierType: "TableColumnIdentifier", Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 4 in a form of "<database_name>.<schema_name>.<table_name>.<table_column_name>"`},
		{IdentifierType: "TableColumnIdentifier", Input: `"""".""""."""".""""`, Error: `unable to parse identifier: """"."""".""""."""", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "TableColumnIdentifier", Input: `"a""bc"."cd""e"."ef""g"."gh""i"`, Error: `unable to parse identifier: "a""bc"."cd""e"."ef""g"."gh""i", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "TableColumnIdentifier", Input: `"".""."".""`, Expected: NewTableColumnIdentifier(``, ``, ``, ``)},
		{IdentifierType: "TableColumnIdentifier", Input: `abc.cde.efg.ghi`, Expected: NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
		{IdentifierType: "TableColumnIdentifier", Input: `"abc"."cde"."efg"."ghi"`, Expected: NewTableColumnIdentifier(`abc`, `cde`, `efg`, `ghi`)},
		{IdentifierType: "TableColumnIdentifier", Input: `"ab.c"."cd.e"."ef.g"."gh.i"`, Expected: NewTableColumnIdentifier(`ab.c`, `cd.e`, `ef.g`, `gh.i`)},

		{IdentifierType: "AccountIdentifier", Input: ``, Error: "incompatible identifier: "},
		{IdentifierType: "AccountIdentifier", Input: "a\nb.cde", Error: "unable to read identifier: a\nb.cde, err = record on line 2: wrong number of fields"},
		{IdentifierType: "AccountIdentifier", Input: `a"b.cde`, Error: "unable to read identifier: a\"b.cde, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{IdentifierType: "AccountIdentifier", Input: `abc.cde.efg`, Error: `unexpected number of parts 3 in identifier abc.cde.efg, expected 2 in a form of "<organization_name>.<account_name>"`},
		{IdentifierType: "AccountIdentifier", Input: `abc`, Error: `unexpected number of parts 1 in identifier abc, expected 2 in a form of "<organization_name>.<account_name>"`},
		{IdentifierType: "AccountIdentifier", Input: `"""".""""`, Error: `unable to parse identifier: """"."""", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "AccountIdentifier", Input: `"a""bc"."cd""e"`, Error: `unable to parse identifier: "a""bc"."cd""e", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "AccountIdentifier", Input: `"".""`, Expected: NewAccountIdentifier(``, ``)},
		{IdentifierType: "AccountIdentifier", Input: `abc.cde`, Expected: NewAccountIdentifier(`abc`, `cde`)},
		{IdentifierType: "AccountIdentifier", Input: `"abc"."cde"`, Expected: NewAccountIdentifier(`abc`, `cde`)},
		{IdentifierType: "AccountIdentifier", Input: `"ab.c"."cd.e"`, Expected: NewAccountIdentifier(`ab.c`, `cd.e`)},

		{IdentifierType: "ExternalObjectIdentifier", Input: ``, Error: "incompatible identifier: "},
		{IdentifierType: "ExternalObjectIdentifier", Input: "a\nb.cde.efg", Error: "unable to read identifier: a\nb.cde.efg, err = record on line 2: wrong number of fields"},
		{IdentifierType: "ExternalObjectIdentifier", Input: `a"b.cde.efg`, Error: "unable to read identifier: a\"b.cde.efg, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{IdentifierType: "ExternalObjectIdentifier", Input: `abc.cde.efg.ghi`, Error: `unexpected number of parts 4 in identifier abc.cde.efg.ghi, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`},
		{IdentifierType: "ExternalObjectIdentifier", Input: `abc.cde`, Error: `unexpected number of parts 2 in identifier abc.cde, expected 3 in a form of "<organization_name>.<account_name>.<external_object_name>"`},
		{IdentifierType: "ExternalObjectIdentifier", Input: `""""."""".""""`, Error: `unable to parse identifier: """".""""."""", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "ExternalObjectIdentifier", Input: `"a""bc"."cd""e"."ef""g"`, Error: `unable to parse identifier: "a""bc"."cd""e"."ef""g", currently identifiers containing double quotes are not supported in the provider`},
		{IdentifierType: "ExternalObjectIdentifier", Input: `""."".""`, Expected: NewExternalObjectIdentifier(NewAccountIdentifier(``, ``), NewAccountObjectIdentifier(``))},
		{IdentifierType: "ExternalObjectIdentifier", Input: `abc.cde.efg`, Expected: NewExternalObjectIdentifier(NewAccountIdentifier(`abc`, `cde`), NewAccountObjectIdentifier(`efg`))},
		{IdentifierType: "ExternalObjectIdentifier", Input: `"abc"."cde"."efg"`, Expected: NewExternalObjectIdentifier(NewAccountIdentifier(`abc`, `cde`), NewAccountObjectIdentifier(`efg`))},
		{IdentifierType: "ExternalObjectIdentifier", Input: `"ab.c"."cd.e"."ef.g"`, Expected: NewExternalObjectIdentifier(NewAccountIdentifier(`ab.c`, `cd.e`), NewAccountObjectIdentifier(`ef.g`))},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Parsing %s with input: "%s"`, testCase.IdentifierType, testCase.Input), func(t *testing.T) {
			var id ObjectIdentifier
			var err error

			switch testCase.IdentifierType {
			case "AccountObjectIdentifier":
				id, err = ParseAccountObjectIdentifier(testCase.Input)
			case "DatabaseObjectIdentifier":
				id, err = ParseDatabaseObjectIdentifier(testCase.Input)
			case "SchemaObjectIdentifier":
				id, err = ParseSchemaObjectIdentifier(testCase.Input)
			case "TableColumnIdentifier":
				id, err = ParseTableColumnIdentifier(testCase.Input)
			case "AccountIdentifier":
				id, err = ParseAccountIdentifier(testCase.Input)
			case "ExternalObjectIdentifier":
				id, err = ParseExternalObjectIdentifier(testCase.Input)
			default:
				t.Fatalf("unknown identifier type: %s", testCase.IdentifierType)
			}

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.Equal(t, testCase.Expected, id)
				assert.NoError(t, err)
			}
		})
	}
}

func Test_ParseObjectIdentifierString(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected ObjectIdentifier
		Error    string
	}{
		{Input: `to.many.parts.for.identifier`, Error: "unsupported identifier: to.many.parts.for.identifier (number of parts: 5)"},
		{Input: "a\nb.cde.efg", Error: "unable to read identifier: a\nb.cde.efg, err = record on line 2: wrong number of fields"},
		{Input: `a"b.cde.efg`, Error: "unable to read identifier: a\"b.cde.efg, err = parse error on line 1, column 2: bare \" in non-quoted-field"},
		{Input: ``, Error: "incompatible identifier: "},
		{Input: `abc`, Expected: NewAccountObjectIdentifier(`abc`)},
		{Input: `abc.def`, Expected: NewDatabaseObjectIdentifier(`abc`, `def`)},
		{Input: `abc.def.ghi`, Expected: NewSchemaObjectIdentifier(`abc`, `def`, `ghi`)},
		{Input: `abc."d.e.f".ghi`, Expected: NewSchemaObjectIdentifier(`abc`, `d.e.f`, `ghi`)},
		{Input: `abc."d""e""f".ghi`, Expected: NewSchemaObjectIdentifier(`abc`, `d"e"f`, `ghi`), Error: `unable to parse identifier: abc."d""e""f".ghi, currently identifiers containing double quotes are not supported in the provider`},
		{Input: `abc.def.ghi.jkl`, Expected: NewTableColumnIdentifier(`abc`, `def`, `ghi`, `jkl`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("ParseObjectIdentifierString for input %s", testCase.Input), func(t *testing.T) {
			id, err := ParseObjectIdentifierString(testCase.Input)

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.Equal(t, testCase.Expected, id)
				assert.NoError(t, err)
			}
		})
	}
}

func TestNewSchemaObjectIdentifierWithArgumentsFromFullyQualifiedName(t *testing.T) {
	testCases := []struct {
		Input SchemaObjectIdentifierWithArguments
		Error string
	}{
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, DataTypeNumber, DataTypeTimestampTZ)},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)")},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, "VECTOR(INT, 20)", DataTypeFloat)},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)", "VECTOR(INT, 10)")},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeTime, "VECTOR(INT, 20)", "VECTOR(FLOAT, 10)", DataTypeFloat)},
		// TODO(SNOW-1571674): Won't work, because of the assumption that identifiers are not containing '(' and ')' parentheses (unfortunately, we're not able to produce meaningful errors for those cases)
		{Input: NewSchemaObjectIdentifierWithArguments(`ab()c`, `def()`, `()ghi`, DataTypeTime, "VECTOR(INT, 20)", "VECTOR(FLOAT, 10)", DataTypeFloat), Error: `unable to read identifier: "ab`},
		{Input: NewSchemaObjectIdentifierWithArguments(`ab(,)c`, `,def()`, `()ghi,`, DataTypeTime, "VECTOR(INT, 20)", "VECTOR(FLOAT, 10)", DataTypeFloat), Error: `unable to read identifier: "ab`},
		{Input: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("processing %s", testCase.Input.FullyQualifiedName()), func(t *testing.T) {
			id, err := ParseSchemaObjectIdentifierWithArguments(testCase.Input.FullyQualifiedName())

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.Input.FullyQualifiedName(), id.FullyQualifiedName())
			}
		})
	}
}

func TestNewSchemaObjectIdentifierWithArgumentsFromFullyQualifiedName_WithRawInput(t *testing.T) {
	testCases := []struct {
		RawInput                    string
		ExpectedIdentifierStructure SchemaObjectIdentifierWithArguments
		Error                       string
	}{
		{RawInput: `abc.def.ghi()`, ExpectedIdentifierStructure: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`)},
		{RawInput: `abc.def.ghi(FLOAT, VECTOR(INT, 20))`, ExpectedIdentifierStructure: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)")},
		// TODO(SNOW-1571674): Won't work, because of the assumption that identifiers are not containing '(' and ')' parentheses (unfortunately, we're not able to produce meaningful errors for those cases)
		{RawInput: `abc."(ef".ghi(FLOAT, VECTOR(INT, 20))`, Error: `unable to read identifier: abc."`},
		{RawInput: `abc.def.ghi`, Error: `unable to parse identifier: '(' not present`},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("processing %s", testCase.ExpectedIdentifierStructure.FullyQualifiedName()), func(t *testing.T) {
			id, err := ParseSchemaObjectIdentifierWithArguments(testCase.RawInput)

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.ExpectedIdentifierStructure.FullyQualifiedName(), id.FullyQualifiedName())
			}
		})
	}
}

func TestNewSchemaObjectIdentifierWithArgumentsAndReturnTypeFromFullyQualifiedName_WithRawInput(t *testing.T) {
	testCases := []struct {
		RawInput                    string
		ExpectedIdentifierStructure SchemaObjectIdentifierWithArguments
		Error                       string
	}{
		{RawInput: `abc.def.ghi()`, ExpectedIdentifierStructure: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`)},
		{RawInput: `abc.def.ghi(FLOAT, VECTOR(INT, 20))`, ExpectedIdentifierStructure: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)")},
		{RawInput: `abc.def.ghi():FLOAT`, ExpectedIdentifierStructure: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`)},
		{RawInput: `abc.def."ghi(FLOAT, VECTOR(INT, 20)):NUMBER(10,2)"`, ExpectedIdentifierStructure: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)")},
		{RawInput: `abc.def."ghi(FLOAT, VECTOR(INT, 20)):NUMBER"`, ExpectedIdentifierStructure: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)")},
		{RawInput: `abc.def."ghi(ab FLOAT, VECTOR VECTOR(INT, 20), FLOAT):NUMBER"`, ExpectedIdentifierStructure: NewSchemaObjectIdentifierWithArguments(`abc`, `def`, `ghi`, DataTypeFloat, "VECTOR(INT, 20)", DataTypeFloat)},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("processing %s", testCase.ExpectedIdentifierStructure.FullyQualifiedName()), func(t *testing.T) {
			id, err := ParseSchemaObjectIdentifierWithArgumentsAndReturnType(testCase.RawInput)

			if testCase.Error != "" {
				assert.ErrorContains(t, err, testCase.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.ExpectedIdentifierStructure.FullyQualifiedName(), id.FullyQualifiedName())
			}
		})
	}
}
