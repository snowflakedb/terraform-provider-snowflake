package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var pipeExecutionStateEnumDef = g.NewEnum(
	"PipeExecutionState", "PipeExecutionStates",
	"FAILING_OVER", "PAUSED", "READ_ONLY", "RUNNING", "STOPPED_BY_SNOWFLAKE_ADMIN", "STOPPED_CLONED",
	"STOPPED_FEATURE_DISABLED", "STOPPED_STAGE_ALTERED", "STOPPED_STAGE_DROPPED", "STOPPED_FILE_FORMAT_DROPPED",
	"STOPPED_NOTIFICATION_INTEGRATION_DROPPED", "STOPPED_MISSING_PIPE", "STOPPED_MISSING_TABLE",
	"STALLED_COMPILATION_ERROR", "STALLED_INITIALIZATION_ERROR", "STALLED_EXECUTION_ERROR",
	"STALLED_INTERNAL_ERROR", "STALLED_STAGE_PERMISSION_ERROR",
)

var forceResumePipeOptionEnumDef = g.NewEnum(
	"ForceResumePipeOption", "ForceResumePipeOptions",
	"STALENESS_CHECK_OVERRIDE", "OWNERSHIP_TRANSFER_CHECK_OVERRIDE",
)

var behaviorChangeBundleStatusEnumDef = g.NewEnum(
	"BehaviorChangeBundleStatus", "BehaviorChangeBundleStatuses",
	"ENABLED", "DISABLED", "RELEASED",
)

func systemFunctionCall(name, functionSQL, column string, args *g.QueryStruct) *g.QueryStruct {
	qs := g.NewQueryStruct(name).SQLWithCustomFieldName("selectExpr", functionSQL)
	if args != nil {
		qs = qs.QueryStructField("Arguments", args, g.ListOptions().MustParentheses().Required())
	}
	if column != "" {
		qs = qs.SQLWithCustomFieldName("resultAlias", `AS \"`+column+`\"`)
	}
	return qs
}

func behaviorChangeBundleArguments() *g.QueryStruct {
	return g.NewQueryStruct("BehaviorChangeBundleArguments").
		Text("Bundle", g.KeywordOptions().SingleQuotes().Required())
}

// NameSingular is empty on purpose. A non-empty singular ("SystemFunction") would produce
// GetTagSystemFunctionOptions etc.
// Identifier kind is a required placeholder — system functions are not a named Snowflake object.
var systemFunctionsDef = g.NewInterface(
	"SystemFunctions",
	"",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CustomShowOperationWithPairedStructsAndOpts(
	"GetTag",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/system_get_tag",
	g.StructPair("getTagRow", "GetTag").
		OptionalText("TAG"),
	systemFunctionCall("GetTag", "SELECT SYSTEM$GET_TAG", "TAG", g.NewQueryStruct("GetTagArguments").
		Identifier("TagId", "ObjectIdentifier", g.IdentifierOptions().SingleQuotes().Required()).
		Identifier("ObjectId", "ObjectIdentifier", g.IdentifierOptions().SingleQuotes().Required()).
		PredefinedQueryStructField("ObjectType", "ObjectType", g.KeywordOptions().SingleQuotes().Required()).
		WithValidation(g.ValidIdentifier, "TagId").
		WithValidation(g.ValidIdentifier, "ObjectId")).
		WithAdditionalValidations(),
	[]g.CustomOperationOption{g.WithRequestAdjust()},
).CustomShowOperationWithPairedStructs(
	"PipeStatus",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/system_pipe_status",
	g.StructPair("pipeStatusRow", "PipeStatus").
		Field("PIPE_STATUS", "string", "PipeExecutionState", g.WithPlainFieldName("ExecutionState"), g.WithCustomParser("parsePipeExecutionState")),
	systemFunctionCall("PipeStatus", "SELECT SYSTEM$PIPE_STATUS", "PIPE_STATUS", g.NewQueryStruct("PipeStatusArguments").
		Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SingleQuotes().Required()).
		WithValidation(g.ValidIdentifier, "Name")),
).CustomOperation(
	"PipeForceResume",
	"https://docs.snowflake.com/en/sql-reference/functions/system_pipe_force_resume",
	systemFunctionCall("PipeForceResume", "SELECT SYSTEM$PIPE_FORCE_RESUME", "", g.NewQueryStruct("PipeForceResumeArguments").
		Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SingleQuotes().Required()).
		ListAssignmentWithFieldName("", "ForceResumePipeOption", g.ParameterOptions().NoEquals().SingleQuotes().NoParentheses(), "Options").
		WithValidation(g.ValidIdentifier, "Name")),
).CustomOperation(
	"EnableBehaviorChangeBundle",
	"https://docs.snowflake.com/en/sql-reference/functions/system_enable_behavior_change_bundle",
	systemFunctionCall("EnableBehaviorChangeBundle", "SELECT SYSTEM$ENABLE_BEHAVIOR_CHANGE_BUNDLE", "", behaviorChangeBundleArguments()),
).CustomOperation(
	"DisableBehaviorChangeBundle",
	"https://docs.snowflake.com/en/sql-reference/functions/system_disable_behavior_change_bundle",
	systemFunctionCall("DisableBehaviorChangeBundle", "SELECT SYSTEM$DISABLE_BEHAVIOR_CHANGE_BUNDLE", "", behaviorChangeBundleArguments()),
).CustomShowOperationWithPairedStructsAndOpts(
	"ShowActiveBehaviorChangeBundles",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/system_show_active_behavior_change_bundles",
	g.StructPair("activeBehaviorChangeBundlesRow", "ActiveBehaviorChangeBundles").
		JsonField("BUNDLES", "[]BehaviorChangeBundleInfo"),
	g.NewQueryStruct("ShowActiveBehaviorChangeBundles").
		SQLWithCustomFieldName("selectExpr", `SELECT SYSTEM$SHOW_ACTIVE_BEHAVIOR_CHANGE_BUNDLES() AS \"BUNDLES\"`),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructs(
	"BehaviorChangeBundleStatus",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/system_behavior_change_bundle_status",
	g.StructPair("behaviorChangeBundleStatusRow", "BehaviorChangeBundleStatusResult").
		Enum("STATUS", behaviorChangeBundleStatusEnumDef),
	systemFunctionCall("BehaviorChangeBundleStatus", "SELECT SYSTEM$BEHAVIOR_CHANGE_BUNDLE_STATUS", "STATUS", behaviorChangeBundleArguments()),
).CustomShowOperationWithPairedStructs(
	"GetIcebergTableInformation",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/system_get_iceberg_table_information",
	g.StructPair("icebergTableInformationRow", "IcebergTableInformation").
		Text("ICEBERG_TABLE_INFORMATION", g.WithPlainFieldName("MetadataLocation"), g.WithCustomParser("parseIcebergTableMetadataLocation")),
	systemFunctionCall("GetIcebergTableInformation", "SELECT SYSTEM$GET_ICEBERG_TABLE_INFORMATION", "ICEBERG_TABLE_INFORMATION", g.NewQueryStruct("GetIcebergTableInformationArguments").
		Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SingleQuotes().Required()).
		WithValidation(g.ValidIdentifier, "Name")),
).CustomShowOperationWithPairedStructs(
	"GetClusteringInformation",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/system_clustering_information",
	g.StructPair("clusteringInformationRow", "ClusteringInformation").
		Text("CLUSTERING_INFORMATION", g.WithPlainFieldName("raw"), g.WithManualConvert()).
		PlainOnlyField("Version", "string").
		PlainOnlyField("ClusterByKeys", "string").
		PlainOnlyField("TotalPartitionCount", "int").
		PlainOnlyField("TotalConstantPartitionCount", "int").
		PlainOnlyField("AverageOverlaps", "float64").
		PlainOnlyField("AverageDepth", "float64").
		PlainOnlyField("PartitionDepthHistogram", "map[string]int").
		PlainOnlyField("ClusteringErrors", "[]ClusteringError"),
	systemFunctionCall(
		"GetClusteringInformation",
		"SELECT SYSTEM$CLUSTERING_INFORMATION",
		"CLUSTERING_INFORMATION",
		g.NewQueryStruct("GetClusteringInformationArguments").
			Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SingleQuotes().Required()).
			ListQueryStructField("Columns", g.NewQueryStruct("ClusteringInformationColumn").
				Text("Name", g.KeywordOptions().DoubleQuotes()), g.ParameterOptions().NoEquals().SingleQuotes().Parentheses()).
			WithValidation(g.ValidIdentifier, "Name"),
	),
	g.PlainStruct("ClusteringError").
		Text("Timestamp").
		Text("Error"),
).CustomShowOperationWithPairedStructs(
	"GetCatalogLinkedDatabaseConfig",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/system_get_catalog_linked_database_config",
	g.StructPair("catalogLinkedDatabaseConfigRow", "CatalogLinkedDatabaseConfig").
		Text("CONFIG", g.WithPlainFieldName("raw"), g.WithManualConvert()).
		PlainOnlyField("CatalogIntegration", "AccountObjectIdentifier").
		PlainOnlyField("CatalogName", "*string").
		PlainOnlyField("ExternalVolume", "*AccountObjectIdentifier").
		PlainOnlyField("SyncIntervalSeconds", "*int").
		PlainOnlyField("NamespaceMode", "*CatalogLinkedDatabaseNamespaceMode").
		PlainOnlyField("NamespaceFlattenDelimiter", "*string").
		PlainOnlyField("AllowedWriteOperations", "*CatalogLinkedDatabaseAllowedWriteOperations").
		PlainOnlyField("CatalogCaseSensitivity", "*DatabaseCatalogCaseSensitivity").
		PlainOnlyField("IsSuspended", "*bool").
		PlainOnlyField("AllowedNamespaces", "[]string").
		PlainOnlyField("BlockedNamespaces", "[]string"),
	systemFunctionCall("GetCatalogLinkedDatabaseConfig", "SELECT SYSTEM$GET_CATALOG_LINKED_DATABASE_CONFIG", "CONFIG", g.NewQueryStruct("GetCatalogLinkedDatabaseConfigArguments").
		Identifier("Name", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SingleQuotes().Required()).
		WithValidation(g.ValidIdentifier, "Name")),
).CustomShowOperationWithPairedStructs(
	"GetCatalogLinkStatus",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/system_catalog_link_status",
	g.StructPair("catalogLinkStatusRow", "CatalogLinkStatus").
		Text("STATUS", g.WithPlainFieldName("raw"), g.WithManualConvert()).
		PlainOnlyField("ExecutionState", "string").
		PlainOnlyField("FailedExecutionStateReason", "*string").
		PlainOnlyField("FailedExecutionStateErrorCode", "*string").
		PlainOnlyField("LastLinkAttemptStartTime", "*time.Time").
		PlainOnlyField("FailureDetails", "[]CatalogLinkFailureDetail"),
	systemFunctionCall("GetCatalogLinkStatus", "SELECT SYSTEM$CATALOG_LINK_STATUS", "STATUS", g.NewQueryStruct("GetCatalogLinkStatusArguments").
		Identifier("Name", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SingleQuotes().Required()).
		WithValidation(g.ValidIdentifier, "Name")),
	g.PlainStruct("CatalogLinkFailureDetail").
		Text("QualifiedEntityName").
		Text("EntityDomain").
		Text("Operation").
		Text("ErrorCode").
		Text("ErrorMessage"),
).WithEnums(
	pipeExecutionStateEnumDef,
	forceResumePipeOptionEnumDef,
	behaviorChangeBundleStatusEnumDef,
)
