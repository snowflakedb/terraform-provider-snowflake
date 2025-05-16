package example

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/poc/generator"
)

//go:generate go run ../main.go

var ToOptsOptionalExample = g.NewInterface(
	"ToOptsOptionalExamples",
	"ToOptsOptionalExample",
	g.KindOfT[DatabaseObjectIdentifier](),
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
		),
)
