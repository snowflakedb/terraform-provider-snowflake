package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
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
	OptionalAssignment(
		"REF_ENTITY_DOMAIN",
		g.KindOfT[sdkcommons.DataMetricFunctionRefEntityDomainOption](),
		g.ParameterOptions().SingleQuotes().ArrowEquals().Required(),
	).WithValidation(g.ValidateValueSet, "RefEntityDomain").
	WithValidation(g.ValidateValueSet, "refEntityName")

var DataMetricFunctionReferenceDef = g.NewInterface(
	"DataMetricFunctionReferences",
	"DataMetricFunctionReference",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CustomOperation(
	"GetForEntity",
	"https://docs.snowflake.com/en/sql-reference/functions/data_metric_function_references",
	g.NewQueryStruct("GetForEntity").
		SQLWithCustomFieldName("selectEverythingFrom", "SELECT * FROM TABLE").
		OptionalQueryStructField(
			"parameters",
			dataMetricFunctionReferenceParametersDef,
			g.ListOptions().Parentheses().NoComma().Required(),
		).WithValidation(g.ValidateValueSet, "parameters"),
	dataMetricFunctionReferenceParametersDef,
	dataMetricFunctionReferenceFunctionArgumentsDef,
	g.DbStruct("dataMetricFunctionReferencesRow").
		Text("METRIC_DATABASE_NAME").
		Text("METRIC_SCHEMA_NAME").
		Text("METRIC_NAME").
		Text("METRIC_SIGNATURE").
		Text("METRIC_DATA_TYPE").
		Text("REF_ENTITY_DATABASE_NAME").
		Text("REF_ENTITY_SCHEMA_NAME").
		Text("REF_ENTITY_NAME").
		Text("REF_ENTITY_DOMAIN").
		Text("REF_ARGUMENTS").
		Text("REF_ID").
		Text("SCHEDULE").
		Text("SCHEDULE_STATUS"),
	g.PlainStruct("DataMetricFunctionReference").
		Text("MetricDatabaseName").
		Text("MetricSchemaName").
		Text("MetricName").
		Text("ArgumentSignature").
		Text("DataType").
		Text("RefEntityDatabaseName").
		Text("RefEntitySchemaName").
		Text("RefEntityName").
		Text("RefEntityDomain").
		Field("RefArguments", "[]DataMetricFunctionRefArgument").
		Text("RefId").
		Text("Schedule").
		Text("ScheduleStatus"),
)
