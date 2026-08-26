package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

func dataMetricFunctionColumnDef() *g.QueryStruct {
	return g.NewQueryStruct("DataMetricFunctionColumn").
		Text("Name", g.KeywordOptions().DoubleQuotes().Required()).
		PredefinedQueryStructField("DataType", "datatypes.DataType", g.ParameterOptions().NoEquals().Required())
}

func dataMetricFunctionArgumentDef() *g.QueryStruct {
	return g.NewQueryStruct("DataMetricFunctionArgument").
		Text("Name", g.KeywordOptions().DoubleQuotes().Required()).
		SQL("TABLE").
		ListQueryStructField("Columns", dataMetricFunctionColumnDef(), g.ParameterOptions().Parentheses().NoEquals().Required())
}

var dataMetricFunctionsDef = g.NewInterface(
	"DataMetricFunctions",
	"DataMetricFunction",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CustomOperation(
	"Create",
	"https://docs.snowflake.com/en/sql-reference/sql/create-data-metric-function",
	g.NewQueryStruct("CreateDataMetricFunction").
		Create().
		SQL("DATA METRIC FUNCTION").
		Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField("Arguments", dataMetricFunctionArgumentDef(), g.ListOptions().MustParentheses().Required()).
		SQL("RETURNS NUMBER").
		PredefinedQueryStructField("Body", "string", g.ParameterOptions().NoEquals().SQL("AS").Required()).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "Name").
		WithValidation(g.ValidateValueSet, "Arguments").
		WithValidation(g.ValidateValueSet, "Body"),
).CustomOperation(
	"DropWithSignature",
	"https://docs.snowflake.com/en/sql-reference/sql/drop-data-metric-function",
	g.NewQueryStruct("DropDataMetricFunction").
		Drop().
		SQL("DATA METRIC FUNCTION").
		OptionalSQL("IF EXISTS").
		Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField("Arguments", dataMetricFunctionArgumentDef(), g.ListOptions().MustParentheses().Required()).
		WithValidation(g.ValidIdentifier, "Name").
		WithValidation(g.ValidateValueSet, "Arguments"),
)
