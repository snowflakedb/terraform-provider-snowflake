package gen_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/stretchr/testify/require"
)

func fieldWithKind(kind string) *gen.Field {
	return &gen.Field{Kind: kind, Tags: map[string][]string{}}
}

func identifierField(kind string) *gen.Field {
	return &gen.Field{Kind: kind, Tags: map[string][]string{"ddl": {"identifier"}}}
}

func structFieldWithChildren(kind string) *gen.Field {
	return &gen.Field{Kind: kind, Tags: map[string][]string{}, Fields: []gen.Field{{Name: "x", Kind: "string", Tags: map[string][]string{}}}}
}

func Test_ZeroValueFor(t *testing.T) {
	tests := []struct {
		name          string
		field         *gen.Field
		expectedValue string
	}{
		{name: "pointer type", field: fieldWithKind("*SomeType"), expectedValue: "nil"},
		{name: "slice type", field: fieldWithKind("[]SomeType"), expectedValue: "nil"},
		{name: "known identifier value type", field: identifierField("SchemaObjectIdentifier"), expectedValue: "emptySchemaObjectIdentifier"},
		{name: "known identifier pointer type", field: identifierField("*SchemaObjectIdentifier"), expectedValue: "nil"}, // pointer wins before identifier check
		{name: "struct value", field: structFieldWithChildren("SomeStruct"), expectedValue: "SomeStruct{}"},
		{name: "bool", field: fieldWithKind("bool"), expectedValue: "false"},
		{name: "int", field: fieldWithKind("int"), expectedValue: "0"},
		{name: "any", field: fieldWithKind("any"), expectedValue: "nil"},
		{name: "package-qualified type (interface, e.g. datatypes.DataType)", field: fieldWithKind("datatypes.DataType"), expectedValue: "nil"},
		{name: "string", field: fieldWithKind("string"), expectedValue: `""`},
		{name: "named string kind (e.g. enum)", field: fieldWithKind("DataType"), expectedValue: `""`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedValue, gen.ZeroValueFor(tt.field))
		})
	}
}

func Test_NonZeroValueFor(t *testing.T) {
	tests := []struct {
		name          string
		field         *gen.Field
		expectedValue string
		expectedOk    bool
	}{
		{
			name:       "slice — not derivable",
			field:      fieldWithKind("[]SomeType"),
			expectedOk: false,
		},
		{
			name:          "pointer to struct",
			field:         structFieldWithChildren("*SomeStruct"),
			expectedValue: "&SomeStruct{}",
			expectedOk:    true,
		},
		{
			name:          "struct value",
			field:         structFieldWithChildren("SomeStruct"),
			expectedValue: "SomeStruct{}",
			expectedOk:    true,
		},
		{
			name:          "*bool",
			field:         fieldWithKind("*bool"),
			expectedValue: "new(true)",
			expectedOk:    true,
		},
		{
			name:          "*int",
			field:         fieldWithKind("*int"),
			expectedValue: "new(1)",
			expectedOk:    true,
		},
		{
			name:          "*string",
			field:         fieldWithKind("*string"),
			expectedValue: `new("foo")`,
			expectedOk:    true,
		},
		{
			name:          "*SchemaObjectIdentifier (pointer to known identifier)",
			field:         identifierField("*SchemaObjectIdentifier"),
			expectedValue: "new(randomSchemaObjectIdentifier())",
			expectedOk:    true,
		},
		{
			name:          "SchemaObjectIdentifier (value identifier)",
			field:         identifierField("SchemaObjectIdentifier"),
			expectedValue: "randomSchemaObjectIdentifier()",
			expectedOk:    true,
		},
		{
			name:          "bool",
			field:         fieldWithKind("bool"),
			expectedValue: "true",
			expectedOk:    true,
		},
		{
			name:          "int",
			field:         fieldWithKind("int"),
			expectedValue: "1",
			expectedOk:    true,
		},
		{
			name:          "string",
			field:         fieldWithKind("string"),
			expectedValue: `"foo"`,
			expectedOk:    true,
		},
		{
			name:       "named non-primitive (e.g. enum, DataType) — not derivable",
			field:      fieldWithKind("DataType"),
			expectedOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := gen.NonZeroValueFor(tt.field)
			require.Equal(t, tt.expectedOk, ok)
			if tt.expectedOk {
				require.Equal(t, tt.expectedValue, got)
			}
		})
	}
}
