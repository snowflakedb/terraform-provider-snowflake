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

func namedIdentifierField(name, kind string) *gen.Field {
	return &gen.Field{Name: name, Kind: kind, Tags: map[string][]string{"ddl": {"identifier"}}}
}

func structFieldWithChildren(kind string) *gen.Field {
	return &gen.Field{Kind: kind, Tags: map[string][]string{}, Fields: []gen.Field{{Name: "x", Kind: "string", Tags: map[string][]string{}}}}
}

// buildTree wires Parent pointers through root's subtree (mirrors the generator's own setParent call in 0_defs.go)
// so Path/IndexedPath/AncestorsFromRoot behave the same way they do on real definitions.
func buildTree(root *gen.Field) *gen.Field {
	gen.SetParent(root)
	return root
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

func Test_Field_IndexedPath_IndexedElemPath(t *testing.T) {
	tests := []struct {
		name                string
		root                *gen.Field
		target              func(root *gen.Field) *gen.Field
		expectedIndexedPath string
		expectedElemPath    string
	}{
		{
			name: "no slice ancestor — equals Path",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "A", Kind: "AStruct", Fields: []gen.Field{
					{Name: "B", Kind: "string"},
				}},
			}},
			target:              func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			expectedIndexedPath: ".A.B",
			expectedElemPath:    ".A.B",
		},
		{
			name: "field itself is a slice — IndexedElemPath indexes it, IndexedPath does not",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument"},
			}},
			target:              func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			expectedIndexedPath: ".Arguments",
			expectedElemPath:    ".Arguments[0]",
		},
		{
			name: "one slice ancestor — child addressed through the single primed element",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument", Fields: []gen.Field{
					{Name: "ArgDataType", Kind: "string"},
				}},
			}},
			target:              func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			expectedIndexedPath: ".Arguments[0].ArgDataType",
			expectedElemPath:    ".Arguments[0].ArgDataType",
		},
		{
			name: "slice of slice — every slice ancestor is indexed",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Outer", Kind: "[]Outer", Fields: []gen.Field{
					{Name: "Inner", Kind: "[]Inner", Fields: []gen.Field{
						{Name: "Leaf", Kind: "string"},
					}},
				}},
			}},
			target:              func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			expectedIndexedPath: ".Outer[0].Inner[0].Leaf",
			expectedElemPath:    ".Outer[0].Inner[0].Leaf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.root)
			target := tt.target(root)
			require.Equal(t, tt.expectedIndexedPath, target.IndexedPath())
			require.Equal(t, tt.expectedElemPath, target.IndexedElemPath())
		})
	}
}

func Test_Validation_TestExpectedError(t *testing.T) {
	tests := []struct {
		name           string
		validation     *gen.Validation
		field          *gen.Field
		expectedOk     bool
		expectedErrSub string // substring the returned expression must contain, when expectedOk
	}{
		{
			name:       "ValidateValue — never derivable, delegates to a nested struct's own validate()",
			validation: gen.NewValidation(gen.ValidateValue, "SessionParameters"),
			field:      &gen.Field{Name: "opts", Kind: "AlterSessionOptions"},
			expectedOk: false,
		},
		{
			name:           "ValidIdentifier — delegates to ReturnedError",
			validation:     gen.NewValidation(gen.ValidIdentifier, "name"),
			field:          &gen.Field{Name: "opts", Kind: "CreateFooOptions"},
			expectedOk:     true,
			expectedErrSub: "ErrInvalidObjectIdentifier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line, ok := tt.validation.TestExpectedError(tt.field)
			require.Equal(t, tt.expectedOk, ok)
			if tt.expectedOk {
				require.Equal(t, tt.expectedErrSub, line)
			} else {
				require.Empty(t, line)
			}
		})
	}
}

func Test_DefaultOptsFieldFor(t *testing.T) {
	tests := []struct {
		name          string
		nameField     *gen.Field
		idVarRef      string
		expectedValue string
	}{
		{
			name:          "value identifier (e.g. name on a Create op)",
			nameField:     namedIdentifierField("name", "AccountObjectIdentifier"),
			idVarRef:      "accountsTestIdAccountObjectIdentifier",
			expectedValue: "name: accountsTestIdAccountObjectIdentifier",
		},
		{
			name:          "pointer identifier (e.g. optional Name on an Alter op)",
			nameField:     namedIdentifierField("Name", "*AccountObjectIdentifier"),
			idVarRef:      "organizationAccountsTestIdAccountObjectIdentifier",
			expectedValue: "Name: new(organizationAccountsTestIdAccountObjectIdentifier)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedValue, gen.DefaultOptsFieldFor(tt.nameField, tt.idVarRef))
		})
	}
}
