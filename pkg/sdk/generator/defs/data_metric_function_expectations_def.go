package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var dataMetricFunctionExpectationParametersDef = g.NewQueryStruct("dataMetricFunctionExpectationParameters").
	SQLWithCustomFieldName("functionFullyQualifiedName", "SNOWFLAKE.INFORMATION_SCHEMA.DATA_METRIC_FUNCTION_EXPECTATIONS").
	OptionalQueryStructField(
		"arguments",
		dataMetricFunctionExpectationFunctionArgumentsDef,
		g.ListOptions().Parentheses().Required(),
	).WithValidation(g.ValidateValueSet, "arguments")

var dataMetricFunctionExpectationFunctionArgumentsDef = g.NewQueryStruct("dataMetricFunctionExpectationFunctionArguments").
	PredefinedQueryStructField("refEntityName", "[]ObjectIdentifier", g.ParameterOptions().ArrowEquals().SingleQuotes().SQL("REF_ENTITY_NAME").Required()).
	OptionalEnumAssignment("REF_ENTITY_DOMAIN", DataMetricFunctionRefEntityDomainOptionEnumDef, g.ParameterOptions().SingleQuotes().ArrowEquals().Required()).
	WithValidation(g.ValidateValueSet, "RefEntityDomain").
	WithValidation(g.ValidateValueSet, "refEntityName")

var dataMetricFunctionExpectationPairs = g.StructPair("dataMetricFunctionExpectationsRow", "DataMetricFunctionExpectation").
	Text("METRIC_DATABASE_NAME", g.WithManualConvert()).
	Text("METRIC_SCHEMA_NAME", g.WithManualConvert()).
	Text("METRIC_NAME", g.WithManualConvert()).
	Text("METRIC_SIGNATURE", g.WithPlainFieldName("ArgumentSignature")).
	Text("EXPECTATION_NAME").
	Text("EXPECTATION_EXPRESSION").
	Text("EXPECTATION_ID").
	Text("REF_ID")

var dataMetricFunctionExpectationsDef = g.NewInterface(
	"DataMetricFunctionExpectations",
	"DataMetricFunctionExpectation",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CustomShowOperationWithPairedStructs(
	"GetForEntity",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/functions/data_metric_function_expectations",
	dataMetricFunctionExpectationPairs,
	g.NewQueryStruct("GetForEntity").
		SQLWithCustomFieldName("selectEverythingFrom", "SELECT * FROM TABLE").
		OptionalQueryStructField(
			"parameters",
			dataMetricFunctionExpectationParametersDef,
			g.ListOptions().Parentheses().NoComma().Required(),
		).WithValidation(g.ValidateValueSet, "parameters"),
	dataMetricFunctionExpectationFunctionArgumentsDef,
)
