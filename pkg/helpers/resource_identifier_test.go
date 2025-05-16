package helpers

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func Test_Encoding_And_Parsing_Of_ResourceIdentifier_String(t *testing.T) {
	testCases := []struct {
		Input                 []string
		Expected              string
		ExpectedAfterDecoding []string
	}{
		{Input: []string{sdk.NewSchemaObjectIdentifier("a", "b", "c").FullyQualifiedName(), "info"}, Expected: `"a"."b"."c"|info`},
		{Input: []string{}, Expected: ``},
		{Input: []string{"", "", ""}, Expected: `||`},
		{Input: []string{"a", "b", "c"}, Expected: `a|b|c`},
		// If one of the parts contains a separator sign (pipe in this case),
		// we can end up with more parts than we started with.
		{Input: []string{"a", "b", "c|d"}, Expected: `a|b|c|d`, ExpectedAfterDecoding: []string{"a", "b", "c", "d"}},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Encoding and parsing %s resource identifier`, testCase.Input), func(t *testing.T) {
			encodedIdentifier := EncodeResourceIdentifier(testCase.Input...)
			assert.Equal(t, testCase.Expected, encodedIdentifier)

			parsedIdentifier := ParseResourceIdentifier(encodedIdentifier)
			if testCase.ExpectedAfterDecoding != nil {
				assert.Equal(t, testCase.ExpectedAfterDecoding, parsedIdentifier)
			} else {
				assert.Equal(t, testCase.Input, parsedIdentifier)
			}
		})
	}
}

func Test_Encoding_And_Parsing_Of_ResourceIdentifier_Identifier(t *testing.T) {
	testCases := []struct {
		Input                 any
		Expected              string
		ExpectedAfterDecoding []string
	}{
		{Input: []sdk.AccountObjectIdentifier{sdk.NewAccountObjectIdentifier("a"), sdk.NewAccountObjectIdentifier("b")}, Expected: `a|b`},
		{Input: []sdk.DatabaseObjectIdentifier{sdk.NewDatabaseObjectIdentifier("a", "b"), sdk.NewDatabaseObjectIdentifier("b", "c")}, Expected: `"a"."b"|"b"."c"`},
		{Input: []sdk.SchemaObjectIdentifier{sdk.NewSchemaObjectIdentifier("a", "b", "c"), sdk.NewSchemaObjectIdentifier("c", "b", "a")}, Expected: `"a"."b"."c"|"c"."b"."a"`},
		{Input: []sdk.SchemaObjectIdentifierWithArguments{sdk.NewSchemaObjectIdentifierWithArguments("a", "b", "c", sdk.DataTypeFloat), sdk.NewSchemaObjectIdentifierWithArguments("c", "b", "a", sdk.DataTypeInt)}, Expected: `"a"."b"."c"(FLOAT)|"c"."b"."a"(NUMBER)`},
		{Input: []sdk.SchemaObjectIdentifierWithArguments{sdk.NewSchemaObjectIdentifierWithArguments("a", "b", "c"), sdk.NewSchemaObjectIdentifierWithArguments("c", "b", "a")}, Expected: `"a"."b"."c"()|"c"."b"."a"()`},
		{Input: []sdk.TableColumnIdentifier{sdk.NewTableColumnIdentifier("a", "b", "c", "d"), sdk.NewTableColumnIdentifier("c", "b", "a", "f")}, Expected: `"a"."b"."c"."d"|"c"."b"."a"."f"`},
		{Input: []sdk.ExternalObjectIdentifier{sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier("o", "a"), sdk.NewAccountObjectIdentifier("ob")), sdk.NewExternalObjectIdentifier(sdk.NewAccountIdentifier("o2", "a2"), sdk.NewAccountObjectIdentifier("ob2"))}, Expected: `"o"."a"."ob"|"o2"."a2"."ob2"`},
		{Input: []sdk.AccountIdentifier{sdk.NewAccountIdentifier("a", "b"), sdk.NewAccountIdentifier("b", "c")}, Expected: `"a"."b"|"b"."c"`},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf(`Encoding and parsing %s resource identifier`, testCase.Input), func(t *testing.T) {
			switch typedInput := testCase.Input.(type) {
			case []sdk.AccountObjectIdentifier:
				encodedIdentifier := EncodeResourceIdentifier(typedInput...)
				assert.Equal(t, testCase.Expected, encodedIdentifier)
			case []sdk.DatabaseObjectIdentifier:
				encodedIdentifier := EncodeResourceIdentifier(typedInput...)
				assert.Equal(t, testCase.Expected, encodedIdentifier)
			case []sdk.SchemaObjectIdentifier:
				encodedIdentifier := EncodeResourceIdentifier(typedInput...)
				assert.Equal(t, testCase.Expected, encodedIdentifier)
			case []sdk.SchemaObjectIdentifierWithArguments:
				encodedIdentifier := EncodeResourceIdentifier(typedInput...)
				assert.Equal(t, testCase.Expected, encodedIdentifier)
			case []sdk.TableColumnIdentifier:
				encodedIdentifier := EncodeResourceIdentifier(typedInput...)
				assert.Equal(t, testCase.Expected, encodedIdentifier)
			case []sdk.ExternalObjectIdentifier:
				encodedIdentifier := EncodeResourceIdentifier(typedInput...)
				assert.Equal(t, testCase.Expected, encodedIdentifier)
			case []sdk.AccountIdentifier:
				encodedIdentifier := EncodeResourceIdentifier(typedInput...)
				assert.Equal(t, testCase.Expected, encodedIdentifier)
			default:
				t.Errorf("Unexpected type %T", typedInput)
			}
		})
	}
}
