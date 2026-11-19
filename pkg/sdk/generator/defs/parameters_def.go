package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
)

const parametersDoc = "https://docs.snowflake.com/en/sql-reference/sql/show-parameters"

func parameterEnums() []*g.Enum {
	return []*g.Enum{
		g.NewEnum("AccountParameter", "AccountParameters", withAdditionalParameterNames(
			parameterNamesForLevels(parameterdefs.ParameterLevelAccountExt),
			ssotdefs.ServiceCallerTokenValiditySecs,
		)...),
		g.NewEnum("SessionParameter", "SessionParameters", parameterNamesForLevels(parameterdefs.ParameterLevelSession)...),
		g.NewEnum("ObjectParameter", "ObjectParameters", parameterNamesForLevels(
			parameterdefs.ParameterLevelAccountExt,
			parameterdefs.ParameterLevelDatabase,
			parameterdefs.ParameterLevelSchema,
			parameterdefs.ParameterLevelTable,
			parameterdefs.ParameterLevelIcebergTable,
			parameterdefs.ParameterLevelHybridTable,
			parameterdefs.ParameterLevelTask,
			parameterdefs.ParameterLevelFunction,
			parameterdefs.ParameterLevelProcedure,
			parameterdefs.ParameterLevelProject,
			parameterdefs.ParameterLevelService,
			parameterdefs.ParameterLevelPipe,
			parameterdefs.ParameterLevelUser,
			parameterdefs.ParameterLevelOpenflowDeployment,
			parameterdefs.ParameterLevelWarehouse,
			parameterdefs.ParameterLevelWarehouseAdaptive,
			parameterdefs.ParameterLevelWarehouseInteractive,
		)...),
		g.NewEnum("UserParameter", "UserParameters", parameterNamesForLevels(parameterdefs.ParameterLevelUser)...),
		g.NewEnum("TaskParameter", "TaskParameters", withAdditionalParameterNames(
			parameterNamesForLevels(parameterdefs.ParameterLevelTask),
			ssotdefs.SearchPath,
		)...),
		g.NewEnum("WarehouseParameter", "WarehouseParameters", parameterNamesForLevels(
			parameterdefs.ParameterLevelWarehouse,
			parameterdefs.ParameterLevelWarehouseAdaptive,
			parameterdefs.ParameterLevelWarehouseInteractive,
		)...),
		g.NewEnum("DatabaseParameter", "DatabaseParameters", withAdditionalParameterNames(
			parameterNamesForLevels(parameterdefs.ParameterLevelDatabase),
			ssotdefs.MetricLevel,
			ssotdefs.ServiceCallerTokenValiditySecs,
		)...),
		g.NewEnum("SchemaParameter", "SchemaParameters", parameterNamesForLevels(parameterdefs.ParameterLevelSchema)...),
		g.NewEnum("FunctionParameter", "FunctionParameters", parameterNamesForLevels(parameterdefs.ParameterLevelFunction)...),
		g.NewEnum("ProcedureParameter", "ProcedureParameters", parameterNamesForLevels(parameterdefs.ParameterLevelProcedure)...),
		g.NewEnum("ServiceParameter", "ServiceParameters", parameterNamesForLevels(parameterdefs.ParameterLevelService)...),
		g.NewEnum("PipeParameter", "PipeParameters", parameterNamesForLevels(parameterdefs.ParameterLevelPipe)...),
		g.NewEnum("TableParameter", "TableParameters", parameterNamesForLevels(parameterdefs.ParameterLevelTable)...),
		g.NewEnum("IcebergTableParameter", "IcebergTableParameters", parameterNamesForLevels(parameterdefs.ParameterLevelIcebergTable)...),
		g.NewEnum("HybridTableParameter", "HybridTableParameters", parameterNamesForLevels(parameterdefs.ParameterLevelHybridTable)...),
		g.NewEnum("OpenflowDeploymentParameter", "OpenflowDeploymentParameters", parameterNamesForLevels(parameterdefs.ParameterLevelOpenflowDeployment)...),
		g.NewEnum("ActivePythonProfiler", "ActivePythonProfilers", "LINE", "MEMORY"),
		g.NewEnum("DefaultNullOrdering", "DefaultNullOrderings", "FIRST", "LAST"),
		g.NewEnum("GeographyOutputFormat", "GeographyOutputFormats", "GeoJSON", "WKT", "WKB", "EWKT", "EWKB"),
		g.NewEnum("GeometryOutputFormat", "GeometryOutputFormats", "GeoJSON", "WKT", "WKB", "EWKT", "EWKB"),
		g.NewEnum("BinaryInputFormat", "BinaryInputFormats", "HEX", "BASE64", "UTF8").WithAliases("UTF8", "UTF-8"),
		g.NewEnum("BinaryOutputFormat", "BinaryOutputFormats", "HEX", "BASE64"),
		g.NewEnum("ClientTimestampTypeMapping", "ClientTimestampTypeMappings", "TIMESTAMP_LTZ", "TIMESTAMP_NTZ"),
		g.NewEnum("TimestampTypeMapping", "TimestampTypeMappings", "TIMESTAMP_LTZ", "TIMESTAMP_NTZ", "TIMESTAMP_TZ"),
		g.NewEnum("TransactionDefaultIsolationLevel", "TransactionDefaultIsolationLevels", "READ COMMITTED"),
		g.NewEnum("UnsupportedDDLAction", "UnsupportedDDLActions", "IGNORE", "FAIL"),
		g.NewEnum("ParameterValueKind", "ParameterValueKinds", "STRING", "IDENTIFIER", "INTEGER", "FLOAT", "BOOLEAN"),
		g.NewEnum("ParameterType", "ParameterTypes",
			"SYSTEM", "ACCOUNT", "USER", "SESSION", "OBJECT", "WAREHOUSE",
			"DATABASE", "SCHEMA", "TABLE", "TASK", "FUNCTION", "PROCEDURE",
			"OPENFLOW_DEPLOYMENT", "SERVICE", "PIPE",
		).WithEmptyValue("SnowflakeDefault"),
	}
}

var parametersDef = g.NewInterface(
	"Parameters",
	"Parameter",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts(g.PartDefault, g.PartDto, g.PartDtoBuilders, g.PartImpl, g.PartValidations, g.PartEnums, g.PartUnitTests).
	CustomOperation(
		"SetAccountParameterRaw",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-account",
		g.NewQueryStruct("SetAccountParameterRaw").
			SQL("ALTER ACCOUNT").
			SQL("SET").
			Text("Parameter", g.KeywordOptions().NoQuotes().Required()).
			SQL("=").
			Text("Value", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidateValueSet, "Parameter").
			WithValidation(g.ValidateValueSet, "Value"),
	).
	CustomOperation(
		"UnsetAccountParameterRaw",
		"https://docs.snowflake.com/en/sql-reference/sql/alter-account",
		g.NewQueryStruct("UnsetAccountParameterRaw").
			SQL("ALTER ACCOUNT").
			SQL("UNSET").
			Text("Parameter", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidateValueSet, "Parameter"),
	).
	CustomOperation(
		"SetObjectParameterOnObjectRaw",
		"https://docs.snowflake.com/en/sql-reference/parameters",
		g.NewQueryStruct("SetObjectParameterOnObjectRaw").
			SQL("ALTER").
			PredefinedQueryStructField("ObjectType", "ObjectType", g.KeywordOptions().Required()).
			Identifier("ObjectName", "ObjectIdentifier", g.IdentifierOptions().Required()).
			SQL("SET").
			Text("Parameter", g.KeywordOptions().NoQuotes().Required()).
			SQL("=").
			Text("Value", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidateValueSet, "ObjectType").
			WithValidation(g.ValidIdentifier, "ObjectName").
			WithValidation(g.ValidateValueSet, "Parameter").
			WithValidation(g.ValidateValueSet, "Value"),
	).
	CustomOperation(
		"UnsetObjectParameterOnObjectRaw",
		"https://docs.snowflake.com/en/sql-reference/parameters",
		g.NewQueryStruct("UnsetObjectParameterOnObjectRaw").
			SQL("ALTER").
			PredefinedQueryStructField("ObjectType", "ObjectType", g.KeywordOptions().Required()).
			Identifier("ObjectName", "ObjectIdentifier", g.IdentifierOptions().Required()).
			SQL("UNSET").
			Text("Parameter", g.KeywordOptions().NoQuotes().Required()).
			WithValidation(g.ValidateValueSet, "ObjectType").
			WithValidation(g.ValidIdentifier, "ObjectName").
			WithValidation(g.ValidateValueSet, "Parameter"),
	).
	WithCustomInterfaceMethod(
		"SetObjectParameterOnAccountViaCatalog",
		"Catalog write path: the value is parsed and quoted by the parameter catalog.\nThe methods return false, and send nothing, when the catalog does not handle the parameter.",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"bool", "error",
	).
	WithCustomInterfaceMethod(
		"UnsetObjectParameterOnAccountViaCatalog",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "ObjectParameter")},
		"bool", "error",
	).
	WithCustomInterfaceMethod(
		"SetObjectParameterOnObjectViaCatalog",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("object", "Object"),
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"bool", "error",
	).
	WithCustomInterfaceMethod(
		"UnsetObjectParameterOnObjectViaCatalog",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("object", "Object"),
			g.NewMethodParameter("parameter", "ObjectParameter"),
		},
		"bool", "error",
	).
	WithCustomInterfaceMethod(
		"SetAccountParameterLegacy",
		"Legacy write path: the SQL of v2.22.x.",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "AccountParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetAccountParameterLegacy",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "AccountParameter")},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetObjectParameterOnAccountLegacy",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetObjectParameterOnAccountLegacy",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "ObjectParameter")},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetAccountParameterWithGenericFallback",
		"Hybrid write path: the catalog write path when it supports the parameter, otherwise the raw write path\n(the value is rendered by the kind, so it is quoted the same way as by the catalog).",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", g.KindString),
			g.NewMethodParameter("kind", "ParameterValueKind"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetAccountParameterWithGenericFallback",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetObjectParameterOnObjectWithGenericFallback",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("object", "Object"),
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("kind", "ParameterValueKind"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetObjectParameterOnObjectWithGenericFallback",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("object", "Object"),
			g.NewMethodParameter("parameter", "ObjectParameter"),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"ShowParameters",
		parametersDoc,
		[]*g.MethodParameter{g.NewMethodParameter("opts", "*ShowParametersOptions")},
		"[]*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"ShowAccountParameter",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "AccountParameter")},
		"*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"ShowSessionParameter",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "SessionParameter")},
		"*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"ShowUserParameter",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "UserParameter"),
			g.NewMethodParameter("user", g.KindOfT[sdkcommons.AccountObjectIdentifier]()),
		},
		"*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"ShowObjectParameter",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("object", "Object"),
		},
		"*Parameter", "error",
	).
	WithEnums(parameterEnums()...)
