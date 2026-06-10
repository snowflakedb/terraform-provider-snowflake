package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var toOptsSliceItemDef = g.NewQueryStruct("ToOptsOptionalExampleSliceItem").
	Text("NAME", g.KeywordOptions()).
	OptionalIdentifier("Parent", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions())

var ToOptsOptionalExample = g.NewInterface(
	"ToOptsOptionalExamples",
	"ToOptsOptionalExample",
	g.KindOfT[sdkcommons.DatabaseObjectIdentifier](),
).CreateOperation("https://example.com",
	g.NewQueryStruct("Alter").
		Alter().
		IfExists().
		Name(),
).AlterOperation("https://example.com",
	g.NewQueryStruct("Alter").
		Alter().
		IfExists().
		Name().
		OptionalQueryStructField(
			"OptionalField",
			g.NewQueryStruct("OptionalField").
				List("SomeList", "DatabaseObjectIdentifier", g.ListOptions()),
			g.KeywordOptions(),
		).
		QueryStructField(
			"RequiredField",
			g.NewQueryStruct("RequiredField").
				List("SomeRequiredList", "DatabaseObjectIdentifier", g.ListOptions().Required()),
			g.KeywordOptions().Required(),
		).
		// Demonstrates ListQueryStructField — generates a for-loop in toOpts using v/s[i] pattern.
		ListQueryStructField("ItemList", toOptsSliceItemDef, g.KeywordOptions().SQL("ITEMS")),
)
