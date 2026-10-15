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
		g.NewEnum("ParameterType", "ParameterTypes",
			"SYSTEM", "ACCOUNT", "USER", "SESSION", "OBJECT", "WAREHOUSE",
			"DATABASE", "SCHEMA", "TABLE", "TASK", "FUNCTION", "PROCEDURE",
			"OPENFLOW_DEPLOYMENT", "SERVICE",
		).WithEmptyValue("SnowflakeDefault"),
	}
}

var parametersDef = g.NewInterface(
	"Parameters",
	"Parameter",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	WithAllowedGenerationParts(g.PartDefault, g.PartImpl, g.PartEnums, g.PartUnitTests).
	WithCustomInterfaceMethod(
		"SetAccountParameter",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "AccountParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetAccountParameter",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "AccountParameter")},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetObjectParameterOnAccount",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetObjectParameterOnObject",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("object", "Object"),
			g.NewMethodParameter("parameter", "ObjectParameter"),
			g.NewMethodParameter("value", g.KindString),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetObjectParameterOnAccount",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("parameter", "ObjectParameter")},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetObjectParameterOnObject",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("object", "Object"),
			g.NewMethodParameter("parameter", "ObjectParameter"),
		},
		"error",
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
		"SetObjectParameterOnAccountViaCatalog",
		"",
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
