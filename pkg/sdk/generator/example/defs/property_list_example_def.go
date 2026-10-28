package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

// PropertyListExample demonstrates a DESCRIBE that returns a generic property list, which the SDK parses into a
// details struct. It combines:
//   - DescribeOperationWithPairedStructs for the raw rows (property, property_value), with the Plain details struct
//     passed as a helper so the struct is generated;
//   - WithDescribeParse, which generates ParsePropertyListExampleDescribed from the property lists;
//   - WithDescribeDetails, which generates the DescribeDetails method (Describe, parse, set Id).
//
// The hand-written parts are in property_list_examples_ext.go: the validate function and the element parser of
// the prefix list.
var PropertyListExample = g.NewInterface(
	"PropertyListExamples",
	"PropertyListExample",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts(g.PartDefault, g.PartDto, g.PartDtoBuilders, g.PartImpl, g.PartValidations).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://example.com",
		g.StructPair("propertyListExampleRow", "PropertyListExampleProperty").
			Text("property", g.WithPlainFieldName("Name")).
			Text("property_value", g.WithPlainFieldName("Value")),
		g.NewQueryStruct("DescPropertyListExample").
			Describe().
			SQL("WHATEVER").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
		g.PlainStruct("PropertyListExampleDetails").
			AccountObjectIdentifier().
			Text("Comment").
			Text("Active").
			Field("Items", "[]PropertyListExampleItemDetails"),
	).
	WithDescribeParse(g.DescribeParseConfig{
		RowsType:      "PropertyListExampleProperty",
		PropertyField: "Name",
		ValueField:    "Value",
		DetailsType:   "PropertyListExampleDetails",
		ErrorLabel:    "property list example",
		ValidateFunc:  "validatePropertyListExampleDescribed",
		Properties: []g.DescribeParseProperty{
			{Property: "COMMENT", Field: "Comment"},
			{Property: "ACTIVE", Field: "Active"},
		},
		Prefixes: []g.DescribeParsePrefix{{
			Prefix:        "ITEM_",
			Field:         "Items",
			Var:           "items",
			ElementType:   "PropertyListExampleItemDetails",
			ElementParser: "ParsePropertyListExampleItemDescribed",
		}},
	}).
	WithDescribeDetails("PropertyListExampleDetails", "ParsePropertyListExampleDescribed")
