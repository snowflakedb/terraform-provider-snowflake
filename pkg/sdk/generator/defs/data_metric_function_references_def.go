package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var DataMetricFunctionRefEntityDomainOptionEnumDef = g.NewEnum(
	"DataMetricFunctionRefEntityDomainOption", "DataMetricFunctionRefEntityDomainOptions",
	"VIEW",
)

var dataMetricFunctionReferenceParametersDef = g.NewQueryStruct("dataMetricFunctionReferenceParameters").
	SQLWithCustomFieldName("functionFullyQualifiedName", "SNOWFLAKE.INFORMATION_SCHEMA.DATA_METRIC_FUNCTION_REFERENCES").
	OptionalQueryStructField(
		"arguments",
		dataMetricFunctionReferenceFunctionArgumentsDef,
		g.ListOptions().Parentheses().Required(),
	).WithValidation(g.ValidateValueSet, "arguments")

var dataMetricFunctionReferenceFunctionArgumentsDef = g.NewQueryStruct("dataMetricFunctionReferenceFunctionArguments").
	PredefinedQueryStructField("refEntityName", "[]ObjectIdentifier", g.ParameterOptions().ArrowEquals().SingleQuotes().SQL("REF_ENTITY_NAME").Required()).
	OptionalEnumAssignment("REF_ENTITY_DOMAIN", DataMetricFunctionRefEntityDomainOptionEnumDef, g.ParameterOptions().SingleQuotes().ArrowEquals().Required()).
	WithValidation(g.ValidateValueSet, "RefEntityDomain").
	WithValidation(g.ValidateValueSet, "refEntityName")

var dataMetricFunctionReferencePairs = g.StructPair("dataMetricFunctionReferencesRow", "DataMetricFunctionReference").
	Text("METRIC_DATABASE_NAME", g.WithManualConvert()).
	Text("METRIC_SCHEMA_NAME", g.WithManualConvert()).
	Text("METRIC_NAME", g.WithManualConvert()).
	Text("METRIC_SIGNATURE", g.WithPlainFieldName("ArgumentSignature")).
	Text("METRIC_DATA_TYPE", g.WithPlainFieldName("DataType")).
	Text("REF_ENTITY_DATABASE_NAME", g.WithManualConvert()).
	Text("REF_ENTITY_SCHEMA_NAME", g.WithManualConvert()).
	Text("REF_ENTITY_NAME", g.WithManualConvert()).
	Text("REF_ENTITY_DOMAIN").
	JsonField("REF_ARGUMENTS", "[]DataMetricFunctionRefArgument").
	Text("REF_ID").
	Text("SCHEDULE").
	Text("SCHEDULE_STATUS")

var dataMetricFunctionReferencesDef = g.NewInterface(
	"DataMetricFunctionReferences",
	"DataMetricFunctionReference",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CustomShowOperationWithPairedStructs(
	"GetForEntity",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/functions/data_metric_function_references",
	dataMetricFunctionReferencePairs,
	g.NewQueryStruct("GetForEntity").
		SQLWithCustomFieldName("selectEverythingFrom", "SELECT * FROM TABLE").
		OptionalQueryStructField(
			"parameters",
			dataMetricFunctionReferenceParametersDef,
			g.ListOptions().Parentheses().NoComma().Required(),
		).WithValidation(g.ValidateValueSet, "parameters"),
	dataMetricFunctionReferenceFunctionArgumentsDef,
).WithEnums(
	DataMetricFunctionRefEntityDomainOptionEnumDef,
)
