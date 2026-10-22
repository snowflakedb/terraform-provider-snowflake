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

func Test_Field_AccessExpr(t *testing.T) {
	tests := []struct {
		name               string
		root               *gen.Field
		target             func(root *gen.Field) *gen.Field
		expectedAccessExpr string
	}{
		{
			name: "no slice ancestor — opts + Path",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "A", Kind: "AStruct", Fields: []gen.Field{
					{Name: "B", Kind: "string"},
				}},
			}},
			target:             func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			expectedAccessExpr: "opts.A.B",
		},
		{
			name: "field directly under a slice — elemVar + .Child",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument", Fields: []gen.Field{
					{Name: "ArgDataType", Kind: "string"},
				}},
			}},
			target:             func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			expectedAccessExpr: "argument.ArgDataType",
		},
		{
			name: "nested struct under a slice — elemVar + .Child.Grandchild",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument", Fields: []gen.Field{
					{Name: "Nested", Kind: "NestedStruct", Fields: []gen.Field{
						{Name: "Leaf", Kind: "string"},
					}},
				}},
			}},
			target:             func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			expectedAccessExpr: "argument.Nested.Leaf",
		},
		{
			name: "nested list under a slice — inner slice field itself is outerElem.Inner",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Outer", Kind: "[]Outer", Fields: []gen.Field{
					{Name: "Inner", Kind: "[]Inner"},
				}},
			}},
			target:             func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			expectedAccessExpr: "outer.Inner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.root)
			target := tt.target(root)
			require.Equal(t, tt.expectedAccessExpr, target.AccessExpr())
		})
	}
}

func Test_Field_SliceIndexVar(t *testing.T) {
	tests := []struct {
		name     string
		field    *gen.Field
		expected string
	}{
		{name: "Arguments", field: &gen.Field{Name: "Arguments"}, expected: "argumentIdx"},
		{name: "Columns", field: &gen.Field{Name: "Columns"}, expected: "columnIdx"},
		{name: "On", field: &gen.Field{Name: "On"}, expected: "onIdx"},
		{name: "LeafItems", field: &gen.Field{Name: "LeafItems"}, expected: "leafItemIdx"},
		{name: "DualChecks", field: &gen.Field{Name: "DualChecks"}, expected: "dualCheckIdx"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.field.SliceIndexVar())
		})
	}
}

func Test_Field_NeedsSliceIndexVar(t *testing.T) {
	tests := []struct {
		name     string
		root     *gen.Field
		target   func(root *gen.Field) *gen.Field
		expected bool
	}{
		{
			name:     "slice with only ValidIdentifier in subtree — index is referenced in errInvalidIdentifier",
			root:     identifierUnderSliceTree(),
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			expected: true,
		},
		{
			name:     "identifier-element slice with ValidIdentifier on the slice itself",
			root:     identifierElementSliceTree(),
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			expected: true,
		},
		{
			name:     "tag-association slice with ValidIdentifier on the slice itself",
			root:     tagAssociationSliceTree(),
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			expected: true,
		},
		{
			name: "slice with its own ExactlyOneValueSet",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument", Validations: []*gen.Validation{
					gen.NewValidation(gen.ExactlyOneValueSet, "ArgDataTypeOld", "ArgDataType"),
				}, Fields: []gen.Field{
					{Name: "ArgDataTypeOld", Kind: "DataType"},
					{Name: "ArgDataType", Kind: "datatypes.DataType"},
				}},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			expected: true,
		},
		{
			name: "slice with nested path-bearing validation — outer index is referenced",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Items", Kind: "[]Item", Fields: []gen.Field{
					{Name: "SubItems", Kind: "[]SubItem", Validations: []*gen.Validation{
						gen.NewValidation(gen.ExactlyOneValueSet, "Name", "Alias"),
					}, Fields: []gen.Field{
						{Name: "Name", Kind: "*string"},
						{Name: "Alias", Kind: "*string"},
					}},
				}},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.root)
			require.Equal(t, tt.expected, tt.target(root).NeedsSliceIndexVar())
		})
	}
}

func Test_Validation_ReturnedError_usesPathWithRootExpr(t *testing.T) {
	noSlice := buildTree(&gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "Set", Kind: "*Set", Fields: []gen.Field{
			{Name: "A", Kind: "*string"},
			{Name: "B", Kind: "*string"},
		}},
	}})
	slice := buildTree(&gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "Arguments", Kind: "[]Argument", Fields: []gen.Field{
			{Name: "ArgDataTypeOld", Kind: "DataType"},
			{Name: "ArgDataType", Kind: "datatypes.DataType"},
		}},
	}})
	nested := buildTree(threeLevelSliceTree())
	identifierSlice := buildTree(identifierUnderSliceTree())

	v := gen.NewValidation(gen.ExactlyOneValueSet, "ArgDataTypeOld", "ArgDataType")
	require.Equal(t, `errExactlyOneOf("RootOptions.Set", "A","B")`,
		gen.NewValidation(gen.ExactlyOneValueSet, "A", "B").ReturnedError(&noSlice.Fields[0]))
	require.Equal(t, `errExactlyOneOf(fmt.Sprintf("RootOptions.Arguments[%d]", argumentIdx), "ArgDataTypeOld","ArgDataType")`,
		v.ReturnedError(&slice.Fields[0]))
	require.Equal(t, `errExactlyOneOf(fmt.Sprintf("RootOptions.Items[%d].SubItems[%d].LeafItems[%d]", itemIdx, subItemIdx, leafItemIdx), "Name","Alias")`,
		gen.NewValidation(gen.ExactlyOneValueSet, "Name", "Alias").ReturnedError(&nested.Fields[0].Fields[0].Fields[0]))
	require.Equal(t, `errInvalidIdentifier("RootOptions", "name")`,
		gen.NewValidation(gen.ValidIdentifier, "name").ReturnedError(noSlice))
	require.Equal(t, `errInvalidIdentifier("RootOptions", "name")`,
		gen.NewValidation(gen.ValidIdentifierIfSet, "name").ReturnedError(noSlice))
	require.Equal(t, `errInvalidIdentifier(fmt.Sprintf("RootOptions.Columns[%d].MaskingPolicy", columnIdx), "MaskingPolicy")`,
		gen.NewValidation(gen.ValidIdentifier, "MaskingPolicy").ReturnedError(&identifierSlice.Fields[0].Fields[0]))
	identifierElemSlice := buildTree(identifierElementSliceTree())
	require.Equal(t, `errInvalidIdentifier(fmt.Sprintf("RootOptions.ExternalAccessIntegrations[%d]", externalAccessIntegrationIdx), "ExternalAccessIntegrations")`,
		gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations").ReturnedError(&identifierElemSlice.Fields[0]))
	tagSlice := buildTree(tagAssociationSliceTree())
	require.Equal(t, `errInvalidIdentifier(fmt.Sprintf("RootOptions.Tag[%d]", tagIdx), "Name")`,
		gen.NewValidation(gen.ValidIdentifier, "Tag").ReturnedError(&tagSlice.Fields[0]))
}

func identifierUnderSliceTree() *gen.Field {
	return &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "Columns", Kind: "[]Column", Fields: []gen.Field{
			{Name: "MaskingPolicy", Kind: "*MaskingPolicy", Validations: []*gen.Validation{
				gen.NewValidation(gen.ValidIdentifier, "MaskingPolicy"),
			}, Fields: []gen.Field{
				{Name: "MaskingPolicy", Kind: "SchemaObjectIdentifier"},
			}},
		}},
	}}
}

// identifierElementSliceTree is Gap B: a flat []AccountObjectIdentifier with ValidIdentifier
// already on the slice field (post-relocate).
func identifierElementSliceTree() *gen.Field {
	return &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "ExternalAccessIntegrations", Kind: "[]AccountObjectIdentifier", Validations: []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations"),
		}},
	}}
}

func identifierElementSliceUnderSetTree() *gen.Field {
	return &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "Set", Kind: "*Set", Fields: []gen.Field{
			{Name: "ExternalAccessIntegrations", Kind: "[]AccountObjectIdentifier", Validations: []*gen.Validation{
				gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations"),
			}},
		}},
	}}
}

func tagAssociationSliceTree() *gen.Field {
	return &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "Tag", Kind: "[]TagAssociation", Validations: []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "Tag"),
		}},
	}}
}

func tagAssociationSliceUnderSetTree() *gen.Field {
	return &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "SetTags", Kind: "[]TagAssociation", Validations: []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "SetTags"),
		}},
	}}
}

func threeLevelSliceTree() *gen.Field {
	return &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "Items", Kind: "[]Item", Fields: []gen.Field{
			{Name: "SubItems", Kind: "[]SubItem", Fields: []gen.Field{
				{Name: "LeafItems", Kind: "[]LeafItem", Fields: []gen.Field{
					{Name: "Name", Kind: "string"},
				}},
			}},
		}},
	}}
}

func Test_Field_PathWithRootExpr(t *testing.T) {
	tests := []struct {
		name     string
		root     *gen.Field
		target   func(root *gen.Field) *gen.Field
		expected string
	}{
		{
			name: "no slice — quoted PathWithRoot",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "A", Kind: "AStruct", Fields: []gen.Field{
					{Name: "B", Kind: "string"},
				}},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			expected: `"RootOptions.A.B"`,
		},
		{
			name: "one slice — the slice field itself",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument"},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			expected: `fmt.Sprintf("RootOptions.Arguments[%d]", argumentIdx)`,
		},
		{
			name: "one slice — child of the element",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument", Fields: []gen.Field{
					{Name: "ArgDataType", Kind: "string"},
				}},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			expected: `fmt.Sprintf("RootOptions.Arguments[%d].ArgDataType", argumentIdx)`,
		},
		{
			name: "nested struct under a slice",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument", Fields: []gen.Field{
					{Name: "Nested", Kind: "NestedStruct", Fields: []gen.Field{
						{Name: "Leaf", Kind: "string"},
					}},
				}},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			expected: `fmt.Sprintf("RootOptions.Arguments[%d].Nested.Leaf", argumentIdx)`,
		},
		{
			name: "two nested slices — inner slice field",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Items", Kind: "[]Item", Fields: []gen.Field{
					{Name: "SubItems", Kind: "[]SubItem"},
				}},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			expected: `fmt.Sprintf("RootOptions.Items[%d].SubItems[%d]", itemIdx, subItemIdx)`,
		},
		{
			name:     "three nested slices — innermost slice field",
			root:     threeLevelSliceTree(),
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			expected: `fmt.Sprintf("RootOptions.Items[%d].SubItems[%d].LeafItems[%d]", itemIdx, subItemIdx, leafItemIdx)`,
		},
		{
			name:     "three nested slices — leaf child of innermost element",
			root:     threeLevelSliceTree(),
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0].Fields[0] },
			expected: `fmt.Sprintf("RootOptions.Items[%d].SubItems[%d].LeafItems[%d].Name", itemIdx, subItemIdx, leafItemIdx)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.root)
			target := tt.target(root)
			require.Equal(t, tt.expected, target.PathWithRootExpr())
			// PathWithRoot stays unindexed — it is used in generator panics, not validate() bodies.
			require.NotContains(t, target.PathWithRoot(), "[")
		})
	}
}

func Test_Field_PathWithRootForTest(t *testing.T) {
	tests := []struct {
		name         string
		root         *gen.Field
		target       func(root *gen.Field) *gen.Field
		failingSlice func(root *gen.Field) *gen.Field
		failingIndex int
		expected     string
	}{
		{
			name: "no slice — equals PathWithRoot",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "A", Kind: "AStruct", Fields: []gen.Field{
					{Name: "B", Kind: "string"},
				}},
			}},
			target:       func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			failingSlice: func(root *gen.Field) *gen.Field { return nil },
			expected:     "RootOptions.A.B",
		},
		{
			name: "one slice — all [0] when failingSlice is nil",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument", Fields: []gen.Field{
					{Name: "ArgDataType", Kind: "string"},
				}},
			}},
			target:       func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			failingSlice: func(root *gen.Field) *gen.Field { return nil },
			expected:     "RootOptions.Arguments[0].ArgDataType",
		},
		{
			name: "one slice — OneValidOneInvalid uses [1] on the slice",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument"},
			}},
			target:       func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			failingSlice: func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			failingIndex: 1,
			expected:     "RootOptions.Arguments[1]",
		},
		{
			name:         "three nested slices — default all [0]",
			root:         threeLevelSliceTree(),
			target:       func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			failingSlice: func(root *gen.Field) *gen.Field { return nil },
			expected:     "RootOptions.Items[0].SubItems[0].LeafItems[0]",
		},
		{
			name:         "three nested slices — [1] only on the innermost (OneValidOneInvalid / BothInvalid[1])",
			root:         threeLevelSliceTree(),
			target:       func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			failingSlice: func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			failingIndex: 1,
			expected:     "RootOptions.Items[0].SubItems[0].LeafItems[1]",
		},
		{
			name:         "three nested slices — BothInvalid[0] on the innermost",
			root:         threeLevelSliceTree(),
			target:       func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			failingSlice: func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			failingIndex: 0,
			expected:     "RootOptions.Items[0].SubItems[0].LeafItems[0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.root)
			target := tt.target(root)
			require.Equal(t, tt.expected, target.PathWithRootForTest(tt.failingSlice(root), tt.failingIndex))
		})
	}
}

func Test_Field_IndexedElemPathAt(t *testing.T) {
	tests := []struct {
		name     string
		root     *gen.Field
		target   func(root *gen.Field) *gen.Field
		index    int
		expected string
	}{
		{
			name: "non-slice — equals IndexedPath regardless of i",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "A", Kind: "string"},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			index:    1,
			expected: ".A",
		},
		{
			name: "slice at 0 — equals IndexedElemPath",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument"},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			index:    0,
			expected: ".Arguments[0]",
		},
		{
			name: "slice at 1",
			root: &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Arguments", Kind: "[]Argument"},
			}},
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			index:    1,
			expected: ".Arguments[1]",
		},
		{
			name:     "nested slice at 1 — ancestors stay [0]",
			root:     threeLevelSliceTree(),
			target:   func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0].Fields[0] },
			index:    1,
			expected: ".Items[0].SubItems[0].LeafItems[1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.root)
			target := tt.target(root)
			require.Equal(t, tt.expected, target.IndexedElemPathAt(tt.index))
			if tt.index == 0 && target.IsSlice() {
				require.Equal(t, target.IndexedElemPath(), target.IndexedElemPathAt(0))
			}
		})
	}
}

func Test_Validation_DeriveModify_usesIndexedPath(t *testing.T) {
	tests := []struct {
		name          string
		root          *gen.Field
		field         func(root *gen.Field) *gen.Field
		validation    *gen.Validation
		expectedLast  string
		expectedPrime []string
	}{
		{
			name: "no slice ancestor — Path-style assignment",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "MaskingPolicy", Kind: "*ColumnMaskingPolicy", Fields: []gen.Field{
					{Name: "Name", Kind: "SchemaObjectIdentifier"},
				}},
			}},
			field:        func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			validation:   gen.NewValidation(gen.ValidIdentifier, "Name"),
			expectedLast: "opts.MaskingPolicy.Name = emptySchemaObjectIdentifier",
			expectedPrime: []string{
				"opts.MaskingPolicy = &ColumnMaskingPolicy{}",
			},
		},
		{
			name: "slice ancestor — assignment goes through [0]",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "Columns", Kind: "[]TableColumn", Fields: []gen.Field{
					{Name: "MaskingPolicy", Kind: "*ColumnMaskingPolicy", Fields: []gen.Field{
						{Name: "Name", Kind: "SchemaObjectIdentifier"},
					}},
				}},
			}},
			field:        func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			validation:   gen.NewValidation(gen.ValidIdentifier, "Name"),
			expectedLast: "opts.Columns[0].MaskingPolicy.Name = emptySchemaObjectIdentifier",
			expectedPrime: []string{
				"opts.Columns = []TableColumn{{}}",
				"opts.Columns[0].MaskingPolicy = &ColumnMaskingPolicy{}",
			},
		},
		{
			name: "validation on a slice field — assignment goes through the element",
			root: &gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
				{Name: "OutOfLineConstraint", Kind: "[]OutOfLineConstraint", Fields: []gen.Field{
					{Name: "Columns", Kind: "[]string"},
				}},
			}},
			field:        func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			validation:   gen.NewValidation(gen.ValidateValueSet, "Columns"),
			expectedLast: "opts.OutOfLineConstraint[0].Columns = nil",
			expectedPrime: []string{
				"opts.OutOfLineConstraint = []OutOfLineConstraint{{}}",
			},
		},
		{
			name:         "identifier-element slice — assign a one-element slice of empty identifiers",
			root:         identifierElementSliceTree(),
			field:        func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			validation:   gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations"),
			expectedLast: "opts.ExternalAccessIntegrations = []AccountObjectIdentifier{emptyAccountObjectIdentifier}",
		},
		{
			name:         "identifier-element slice under a pointer — prime the parent then assign the slice",
			root:         identifierElementSliceUnderSetTree(),
			field:        func(root *gen.Field) *gen.Field { return &root.Fields[0].Fields[0] },
			validation:   gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations"),
			expectedLast: "opts.Set.ExternalAccessIntegrations = []AccountObjectIdentifier{emptyAccountObjectIdentifier}",
			expectedPrime: []string{
				"opts.Set = &Set{}",
			},
		},
		{
			name:         "tag-association slice — assign a one-element slice with an empty Name",
			root:         tagAssociationSliceTree(),
			field:        func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			validation:   gen.NewValidation(gen.ValidIdentifier, "Tag"),
			expectedLast: "opts.Tag = []TagAssociation{{Name: emptyAccountObjectIdentifier}}",
		},
		{
			name:         "tag-association SetTags — assign a one-element slice with an empty Name",
			root:         tagAssociationSliceUnderSetTree(),
			field:        func(root *gen.Field) *gen.Field { return &root.Fields[0] },
			validation:   gen.NewValidation(gen.ValidIdentifier, "SetTags"),
			expectedLast: "opts.SetTags = []TagAssociation{{Name: emptyAccountObjectIdentifier}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := buildTree(tt.root)
			field := tt.field(root)
			lines, ok := tt.validation.DeriveModify(field)
			require.True(t, ok)
			require.Equal(t, append(append([]string{}, tt.expectedPrime...), tt.expectedLast), lines)
		})
	}
}

func Test_deriveConflictingFieldsModify_primesSlice(t *testing.T) {
	root := buildTree(&gen.Field{Name: "Root", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "OutOfLineConstraint", Kind: "[]OutOfLineConstraint", Fields: []gen.Field{
			{Name: "Enforced", Kind: "*bool"},
			{Name: "NotEnforced", Kind: "*bool"},
		}},
	}})
	field := &root.Fields[0]
	v := gen.NewValidation(gen.ConflictingFields, "Enforced", "NotEnforced")
	lines, ok := gen.DeriveConflictingFieldsModify(v, field)
	require.True(t, ok)
	require.Equal(t, []string{
		"opts.OutOfLineConstraint = []OutOfLineConstraint{{}}",
		"opts.OutOfLineConstraint[0].Enforced = new(true)",
		"opts.OutOfLineConstraint[0].NotEnforced = new(true)",
	}, lines)
}

func Test_BuildMultiFieldValidationCases_ConflictingFields_BothInvalid(t *testing.T) {
	root := buildTree(&gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
		{Name: "DualChecks", Kind: "[]DualCheckItem", Fields: []gen.Field{
			{Name: "B", Kind: "*string"},
			{Name: "C", Kind: "*string"},
		}},
	}})
	field := &root.Fields[0]
	v := gen.NewValidation(gen.ConflictingFields, "B", "C")
	cases := gen.BuildMultiFieldValidationCases(v, field, "Create", false)

	bySuffix := map[string]*gen.ValidationTestCase{}
	for _, c := range cases {
		bySuffix[c.Name] = c
	}
	both := bySuffix["validation_Create_opts_DualChecks_ConflictingFields_BothInvalid"]
	require.NotNil(t, both)
	require.Equal(t, []string{
		`errOneOf("RootOptions.DualChecks[0]", "B","C")`,
		`errOneOf("RootOptions.DualChecks[1]", "B","C")`,
	}, both.ExpectedErrLines)
	require.True(t, both.HasModify)
	require.Equal(t, []string{
		"opts.DualChecks = []DualCheckItem{{}, {}}",
		`opts.DualChecks[0].B = new("foo")`,
		`opts.DualChecks[0].C = new("foo")`,
		`opts.DualChecks[1].B = new("foo")`,
		`opts.DualChecks[1].C = new("foo")`,
	}, both.ModifyLines)
}

func Test_Validation_TestExpectedError(t *testing.T) {
	sliceTree := buildTree(&gen.Field{Name: "CreateFooOptions", Kind: "CreateFooOptions", Fields: []gen.Field{
		{Name: "Arguments", Kind: "[]Argument", Fields: []gen.Field{
			{Name: "ArgDataTypeOld", Kind: "DataType"},
			{Name: "ArgDataType", Kind: "datatypes.DataType"},
		}},
	}})
	arguments := &sliceTree.Fields[0]

	nestedTree := buildTree(threeLevelSliceTree())
	leafItems := &nestedTree.Fields[0].Fields[0].Fields[0]
	identifierRoot := buildTree(&gen.Field{Name: "RootOptions", Kind: "RootOptions"})
	identifierSlice := buildTree(identifierUnderSliceTree())
	maskingPolicy := &identifierSlice.Fields[0].Fields[0]
	identifierElemSlice := buildTree(identifierElementSliceTree())
	eaiSlice := &identifierElemSlice.Fields[0]
	tagAssocSlice := buildTree(tagAssociationSliceTree())
	tagSlice := &tagAssocSlice.Fields[0]

	tests := []struct {
		name         string
		validation   *gen.Validation
		field        *gen.Field
		failingSlice *gen.Field
		failingIndex int
		expectedOk   bool
		expectedErr  string
	}{
		{
			name:       "ValidateValue — never derivable, delegates to a nested struct's own validate()",
			validation: gen.NewValidation(gen.ValidateValue, "SessionParameters"),
			field:      &gen.Field{Name: "opts", Kind: "AlterSessionOptions"},
			expectedOk: false,
		},
		{
			name:        "ValidIdentifier — root name uses quoted PathWithRootForTest",
			validation:  gen.NewValidation(gen.ValidIdentifier, "name"),
			field:       identifierRoot,
			expectedOk:  true,
			expectedErr: `errInvalidIdentifier("RootOptions", "name")`,
		},
		{
			name:        "ValidIdentifierIfSet — same helper as ValidIdentifier",
			validation:  gen.NewValidation(gen.ValidIdentifierIfSet, "name"),
			field:       identifierRoot,
			expectedOk:  true,
			expectedErr: `errInvalidIdentifier("RootOptions", "name")`,
		},
		{
			name:        "ValidIdentifier under a slice — default all [0]",
			validation:  gen.NewValidation(gen.ValidIdentifier, "MaskingPolicy"),
			field:       maskingPolicy,
			expectedOk:  true,
			expectedErr: `errInvalidIdentifier("RootOptions.Columns[0].MaskingPolicy", "MaskingPolicy")`,
		},
		{
			name:        "ValidIdentifier on an identifier-element slice — default [0]",
			validation:  gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations"),
			field:       eaiSlice,
			expectedOk:  true,
			expectedErr: `errInvalidIdentifier("RootOptions.ExternalAccessIntegrations[0]", "ExternalAccessIntegrations")`,
		},
		{
			name:        "ValidIdentifier on a tag-association slice — default [0], field Name",
			validation:  gen.NewValidation(gen.ValidIdentifier, "Tag"),
			field:       tagSlice,
			expectedOk:  true,
			expectedErr: `errInvalidIdentifier("RootOptions.Tag[0]", "Name")`,
		},
		{
			name:        "ExactlyOneValueSet on a slice — default all [0]",
			validation:  gen.NewValidation(gen.ExactlyOneValueSet, "ArgDataTypeOld", "ArgDataType"),
			field:       arguments,
			expectedOk:  true,
			expectedErr: `errExactlyOneOf("CreateFooOptions.Arguments[0]", "ArgDataTypeOld","ArgDataType")`,
		},
		{
			name:         "ExactlyOneValueSet on a slice — OneValidOneInvalid uses [1]",
			validation:   gen.NewValidation(gen.ExactlyOneValueSet, "ArgDataTypeOld", "ArgDataType"),
			field:        arguments,
			failingSlice: arguments,
			failingIndex: 1,
			expectedOk:   true,
			expectedErr:  `errExactlyOneOf("CreateFooOptions.Arguments[1]", "ArgDataTypeOld","ArgDataType")`,
		},
		{
			name:        "three nested slices — default all [0]",
			validation:  gen.NewValidation(gen.ExactlyOneValueSet, "Name", "Alias"),
			field:       leafItems,
			expectedOk:  true,
			expectedErr: `errExactlyOneOf("RootOptions.Items[0].SubItems[0].LeafItems[0]", "Name","Alias")`,
		},
		{
			name:         "three nested slices — [1] only on the validated innermost slice",
			validation:   gen.NewValidation(gen.ExactlyOneValueSet, "Name", "Alias"),
			field:        leafItems,
			failingSlice: leafItems,
			failingIndex: 1,
			expectedOk:   true,
			expectedErr:  `errExactlyOneOf("RootOptions.Items[0].SubItems[0].LeafItems[1]", "Name","Alias")`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line, ok := tt.validation.TestExpectedError(tt.field, tt.failingSlice, tt.failingIndex)
			require.Equal(t, tt.expectedOk, ok)
			if tt.expectedOk {
				require.Equal(t, tt.expectedErr, line)
			} else {
				require.Empty(t, line)
			}
		})
	}
}

func Test_BuildMultiFieldValidationCases_OneValidOneInvalid_usesIndexOne(t *testing.T) {
	root := buildTree(&gen.Field{Name: "CreateFooOptions", Kind: "CreateFooOptions", Fields: []gen.Field{
		{Name: "Arguments", Kind: "[]Argument", Fields: []gen.Field{
			{Name: "ArgDataTypeOld", Kind: "DataType"},
			{Name: "ArgDataType", Kind: "datatypes.DataType"},
		}},
	}})
	arguments := &root.Fields[0]
	v := gen.NewValidation(gen.ExactlyOneValueSet, "ArgDataTypeOld", "ArgDataType")

	cases := gen.BuildMultiFieldValidationCases(v, arguments, "CreateForJava", false)

	bySuffix := map[string]*gen.ValidationTestCase{}
	for _, c := range cases {
		bySuffix[c.Name] = c
	}
	require.Equal(t, `errExactlyOneOf("CreateFooOptions.Arguments[0]", "ArgDataTypeOld","ArgDataType")`, bySuffix["validation_CreateForJava_opts_Arguments_ExactlyOneValueSet_NoneSet"].ExpectedErrLine)
	require.Equal(t, `errExactlyOneOf("CreateFooOptions.Arguments[0]", "ArgDataTypeOld","ArgDataType")`, bySuffix["validation_CreateForJava_opts_Arguments_ExactlyOneValueSet_MoreThanOneSet"].ExpectedErrLine)
	require.Equal(t, `errExactlyOneOf("CreateFooOptions.Arguments[1]", "ArgDataTypeOld","ArgDataType")`, bySuffix["validation_CreateForJava_opts_Arguments_ExactlyOneValueSet_OneValidOneInvalid"].ExpectedErrLine)

	both := bySuffix["validation_CreateForJava_opts_Arguments_ExactlyOneValueSet_BothInvalid"]
	require.NotNil(t, both)
	require.True(t, both.HasExpectedErrs())
	require.Equal(t, []string{
		`errExactlyOneOf("CreateFooOptions.Arguments[0]", "ArgDataTypeOld","ArgDataType")`,
		`errExactlyOneOf("CreateFooOptions.Arguments[1]", "ArgDataTypeOld","ArgDataType")`,
	}, both.ExpectedErrLines)
	require.True(t, both.HasModify)
	require.Equal(t, []string{
		"opts.Arguments = []Argument{{}, {}}",
	}, both.ModifyLines)
}

func Test_BuildMultiFieldValidationCases_BothInvalid_threeLevelPriming(t *testing.T) {
	root := buildTree(threeLevelSliceTree())
	leafItems := &root.Fields[0].Fields[0].Fields[0]
	v := gen.NewValidation(gen.ExactlyOneValueSet, "Name", "Alias")
	cases := gen.BuildMultiFieldValidationCases(v, leafItems, "Create", false)

	bySuffix := map[string]*gen.ValidationTestCase{}
	for _, c := range cases {
		bySuffix[c.Name] = c
	}
	both := bySuffix["validation_Create_opts_Items_SubItems_LeafItems_ExactlyOneValueSet_BothInvalid"]
	require.NotNil(t, both)
	require.Equal(t, []string{
		`errExactlyOneOf("RootOptions.Items[0].SubItems[0].LeafItems[0]", "Name","Alias")`,
		`errExactlyOneOf("RootOptions.Items[0].SubItems[0].LeafItems[1]", "Name","Alias")`,
	}, both.ExpectedErrLines)
	require.Equal(t, []string{
		"opts.Items = []Item{{}}",
		"opts.Items[0].SubItems = []SubItem{{}}",
		"opts.Items[0].SubItems[0].LeafItems = []LeafItem{{}, {}}",
	}, both.ModifyLines)
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

func Test_Field_IsIdentifierElementSlice(t *testing.T) {
	tests := []struct {
		name     string
		field    *gen.Field
		expected bool
	}{
		{
			name:     "[]AccountObjectIdentifier with no children",
			field:    &gen.Field{Name: "ExternalAccessIntegrations", Kind: "[]AccountObjectIdentifier"},
			expected: true,
		},
		{
			name:     "[]SchemaObjectIdentifier with no children",
			field:    &gen.Field{Name: "Tables", Kind: "[]SchemaObjectIdentifier"},
			expected: true,
		},
		{
			name:     "[]string is not an identifier element slice",
			field:    &gen.Field{Name: "Values", Kind: "[]string"},
			expected: false,
		},
		{
			name: "[]Column with QueryStruct children is not an identifier element slice",
			field: &gen.Field{Name: "Columns", Kind: "[]Column", Fields: []gen.Field{
				{Name: "Name", Kind: "string"},
			}},
			expected: false,
		},
		{
			name: "[]AccountObjectIdentifier with children is not treated as a flat identifier list",
			field: &gen.Field{Name: "Integrations", Kind: "[]AccountObjectIdentifier", Fields: []gen.Field{
				{Name: "Name", Kind: "string"},
			}},
			expected: false,
		},
		{
			name:     "non-slice identifier",
			field:    &gen.Field{Name: "name", Kind: "AccountObjectIdentifier"},
			expected: false,
		},
		{
			name:     "[]TagAssociation is not an identifier element slice",
			field:    &gen.Field{Name: "Tag", Kind: "[]TagAssociation"},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.field.IsIdentifierElementSlice())
		})
	}
}

func Test_Field_IsTagAssociationSlice(t *testing.T) {
	tests := []struct {
		name     string
		field    *gen.Field
		expected bool
	}{
		{
			name:     "[]TagAssociation with no children",
			field:    &gen.Field{Name: "Tag", Kind: "[]TagAssociation"},
			expected: true,
		},
		{
			name:     "[]TagAssociation SetTags with no children",
			field:    &gen.Field{Name: "SetTags", Kind: "[]TagAssociation"},
			expected: true,
		},
		{
			name:     "[]AccountObjectIdentifier is not a tag-association slice",
			field:    &gen.Field{Name: "ExternalAccessIntegrations", Kind: "[]AccountObjectIdentifier"},
			expected: false,
		},
		{
			name: "[]TagAssociation with QueryStruct children is not treated as flat",
			field: &gen.Field{Name: "Tag", Kind: "[]TagAssociation", Fields: []gen.Field{
				{Name: "Name", Kind: "ObjectIdentifier"},
			}},
			expected: false,
		},
		{
			name:     "non-slice TagAssociation",
			field:    &gen.Field{Name: "Tag", Kind: "TagAssociation"},
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, tt.field.IsTagAssociationSlice())
		})
	}
}

func Test_RelocateIdentifierElementSliceValidations(t *testing.T) {
	t.Run("moves ValidIdentifier from the container onto the identifier-element slice", func(t *testing.T) {
		root := &gen.Field{Name: "RootOptions", Kind: "RootOptions", Validations: []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "name"),
			gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations"),
			gen.NewValidation(gen.ConflictingFields, "OrReplace", "IfNotExists"),
		}, Fields: []gen.Field{
			{Name: "name", Kind: "SchemaObjectIdentifier"},
			{Name: "ExternalAccessIntegrations", Kind: "[]AccountObjectIdentifier"},
			{Name: "OrReplace", Kind: "*bool"},
			{Name: "IfNotExists", Kind: "*bool"},
		}}
		gen.SetParent(root)
		gen.RelocateIdentifierElementSliceValidations(root)

		require.Equal(t, []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "name"),
			gen.NewValidation(gen.ConflictingFields, "OrReplace", "IfNotExists"),
		}, root.Validations)
		require.Equal(t, []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations"),
		}, root.Fields[1].Validations)
	})

	t.Run("does not move ValidIdentifier for a struct-element slice", func(t *testing.T) {
		root := identifierUnderSliceTree()
		root.Validations = []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "Columns"),
		}
		gen.SetParent(root)
		gen.RelocateIdentifierElementSliceValidations(root)

		require.Equal(t, []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "Columns"),
		}, root.Validations)
		require.Empty(t, root.Fields[0].Validations)
	})

	t.Run("does not move ValidIdentifier for a []string list", func(t *testing.T) {
		root := &gen.Field{Name: "RootOptions", Kind: "RootOptions", Validations: []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "Values"),
		}, Fields: []gen.Field{
			{Name: "Values", Kind: "[]string"},
		}}
		gen.SetParent(root)
		gen.RelocateIdentifierElementSliceValidations(root)

		require.Equal(t, []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "Values"),
		}, root.Validations)
		require.Empty(t, root.Fields[0].Validations)
	})

	t.Run("recurses into nested containers", func(t *testing.T) {
		root := &gen.Field{Name: "RootOptions", Kind: "RootOptions", Fields: []gen.Field{
			{Name: "Set", Kind: "*Set", Validations: []*gen.Validation{
				gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations"),
			}, Fields: []gen.Field{
				{Name: "ExternalAccessIntegrations", Kind: "[]AccountObjectIdentifier"},
			}},
		}}
		gen.SetParent(root)
		gen.RelocateIdentifierElementSliceValidations(root)

		require.Empty(t, root.Fields[0].Validations)
		require.Equal(t, []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations"),
		}, root.Fields[0].Fields[0].Validations)
	})

	t.Run("moves ValidIdentifier from the container onto a tag-association slice", func(t *testing.T) {
		root := &gen.Field{Name: "RootOptions", Kind: "RootOptions", Validations: []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "name"),
			gen.NewValidation(gen.ValidIdentifier, "Tag"),
		}, Fields: []gen.Field{
			{Name: "name", Kind: "SchemaObjectIdentifier"},
			{Name: "Tag", Kind: "[]TagAssociation"},
		}}
		gen.SetParent(root)
		gen.RelocateIdentifierElementSliceValidations(root)

		require.Equal(t, []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "name"),
		}, root.Validations)
		require.Equal(t, []*gen.Validation{
			gen.NewValidation(gen.ValidIdentifier, "Tag"),
		}, root.Fields[1].Validations)
	})
}

func Test_Validation_Condition_identifierElementSlice(t *testing.T) {
	root := buildTree(identifierElementSliceTree())
	slice := &root.Fields[0]
	v := gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations")
	require.Equal(t, "!ValidObjectIdentifier(externalAccessIntegration)", v.Condition(slice))

	require.Panics(t, func() {
		gen.NewValidation(gen.ValidIdentifierIfSet, "ExternalAccessIntegrations").Condition(slice)
	})
}

func Test_Validation_Condition_tagAssociationSlice(t *testing.T) {
	root := buildTree(tagAssociationSliceTree())
	slice := &root.Fields[0]
	v := gen.NewValidation(gen.ValidIdentifier, "Tag")
	require.Equal(t, "!ValidObjectIdentifier(tag.Name)", v.Condition(slice))

	require.Panics(t, func() {
		gen.NewValidation(gen.ValidIdentifierIfSet, "Tag").Condition(slice)
	})
}

func Test_BuildSingleFieldValidationCase_identifierElementSliceSlug(t *testing.T) {
	t.Run("does not double the slice name", func(t *testing.T) {
		root := buildTree(identifierElementSliceTree())
		slice := &root.Fields[0]
		v := gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations")
		tc := gen.BuildSingleFieldValidationCase(v, slice, "Create", `errInvalidIdentifier("RootOptions.ExternalAccessIntegrations[0]", "ExternalAccessIntegrations")`, true)
		require.Equal(t, "validation_Create_ExternalAccessIntegrations_ValidIdentifier", tc.Name)
		require.Equal(t, []string{
			"opts.ExternalAccessIntegrations = []AccountObjectIdentifier{emptyAccountObjectIdentifier}",
		}, tc.ModifyLines)
	})

	t.Run("keeps nested container path without doubling", func(t *testing.T) {
		root := buildTree(identifierElementSliceUnderSetTree())
		slice := &root.Fields[0].Fields[0]
		v := gen.NewValidation(gen.ValidIdentifier, "ExternalAccessIntegrations")
		tc := gen.BuildSingleFieldValidationCase(v, slice, "Alter", "", true)
		require.Equal(t, "validation_Alter_Set_ExternalAccessIntegrations_ValidIdentifier", tc.Name)
	})
}

func Test_BuildSingleFieldValidationCase_tagAssociationSliceSlug(t *testing.T) {
	t.Run("does not double the slice name", func(t *testing.T) {
		root := buildTree(tagAssociationSliceTree())
		slice := &root.Fields[0]
		v := gen.NewValidation(gen.ValidIdentifier, "Tag")
		tc := gen.BuildSingleFieldValidationCase(v, slice, "Create", `errInvalidIdentifier("RootOptions.Tag[0]", "Name")`, true)
		require.Equal(t, "validation_Create_Tag_ValidIdentifier", tc.Name)
		require.Equal(t, []string{
			"opts.Tag = []TagAssociation{{Name: emptyAccountObjectIdentifier}}",
		}, tc.ModifyLines)
	})

	t.Run("SetTags does not double the slice name", func(t *testing.T) {
		root := buildTree(tagAssociationSliceUnderSetTree())
		slice := &root.Fields[0]
		v := gen.NewValidation(gen.ValidIdentifier, "SetTags")
		tc := gen.BuildSingleFieldValidationCase(v, slice, "Alter", "", true)
		require.Equal(t, "validation_Alter_SetTags_ValidIdentifier", tc.Name)
	})
}
