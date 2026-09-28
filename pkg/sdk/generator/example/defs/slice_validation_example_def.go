package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// Slice-validation case matrix. Nested QueryStructs that themselves have struct children
// are function wrappers — sharing them as a var panics (`Field already has a parent`).
// SliceElemVar is CutSuffix(name, "s"); list field names must stay distinct (Items / SubItems).

// Case 1: slice with two WithValidations (previously emitted invalid `} for` Go).
func sliceValidationDualCheck() *g.QueryStruct {
	return g.NewQueryStruct("DualCheckItem").
		OptionalText("A", g.KeywordOptions()).
		OptionalText("B", g.KeywordOptions()).
		OptionalText("C", g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "A", "B").
		WithValidation(g.ConflictingFields, "B", "C")
}

// Case 2: nested struct under a slice with no own slice validation (previously skipped).
func sliceValidationNestedNoOwn() *g.QueryStruct {
	return g.NewQueryStruct("NestedNoOwn").
		OptionalText("X", g.KeywordOptions()).
		OptionalText("Y", g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "X", "Y")
}

func sliceValidationPlainItem() *g.QueryStruct {
	return g.NewQueryStruct("PlainItem").
		OptionalQueryStructField("Nested", sliceValidationNestedNoOwn(), g.KeywordOptions())
}

// Case 3: nested struct under a slice that also has its own element validation
// (previously nested fields used opts.Slice.Child instead of the loop var).
func sliceValidationNestedWithOwn() *g.QueryStruct {
	return g.NewQueryStruct("NestedWithOwn").
		OptionalText("P", g.KeywordOptions()).
		OptionalText("Q", g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "P", "Q")
}

func sliceValidationCheckedItem() *g.QueryStruct {
	return g.NewQueryStruct("CheckedItem").
		OptionalText("Left", g.KeywordOptions()).
		OptionalText("Right", g.KeywordOptions()).
		OptionalQueryStructField("Nested", sliceValidationNestedWithOwn(), g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "Left", "Right")
}

// Case 4: nested ListQueryStructField (list of objects containing a list of objects).
func sliceValidationSubItem() *g.QueryStruct {
	return g.NewQueryStruct("SubItem").
		OptionalText("Name", g.KeywordOptions()).
		OptionalText("Alias", g.KeywordOptions()).
		WithValidation(g.ExactlyOneValueSet, "Name", "Alias")
}

func sliceValidationItem() *g.QueryStruct {
	return g.NewQueryStruct("NestedListItem").
		ListQueryStructField("SubItems", sliceValidationSubItem(), g.KeywordOptions().SQL("SUB_ITEMS"))
}

var SliceValidationExample = g.NewInterface(
	"SliceValidationExamples",
	"SliceValidationExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://example.com",
	g.NewQueryStruct("CreateSliceValidationExample").
		Create().
		SQL("SLICE VALIDATION EXAMPLE").
		Name().
		ListQueryStructField("DualChecks", sliceValidationDualCheck(), g.KeywordOptions().SQL("DUAL_CHECKS")).
		ListQueryStructField("PlainItems", sliceValidationPlainItem(), g.KeywordOptions().SQL("PLAIN_ITEMS")).
		ListQueryStructField("CheckedItems", sliceValidationCheckedItem(), g.KeywordOptions().SQL("CHECKED_ITEMS")).
		ListQueryStructField("Items", sliceValidationItem(), g.KeywordOptions().SQL("ITEMS")).
		WithValidation(g.ValidIdentifier, "name"),
)
