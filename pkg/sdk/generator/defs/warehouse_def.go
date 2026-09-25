package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var warehouseTypeEnum = g.NewEnum(
	"WarehouseType", "WarehouseTypes",
	"STANDARD", "SNOWPARK-OPTIMIZED", "ADAPTIVE", "INTERACTIVE",
)

var warehouseSizeEnum = g.NewEnum(
	"WarehouseSize", "WarehouseSizes",
	"XSMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "XXLARGE", "XXXLARGE", "X4LARGE", "X5LARGE", "X6LARGE",
).WithAliases("XSMALL", "X-SMALL").
	WithAliases("XLARGE", "X-LARGE").
	WithAliases("XXLARGE", "X2LARGE", "2X-LARGE").
	WithAliases("XXXLARGE", "X3LARGE", "3X-LARGE").
	WithAliases("X4LARGE", "4X-LARGE").
	WithAliases("X5LARGE", "5X-LARGE").
	WithAliases("X6LARGE", "6X-LARGE")

var scalingPolicyEnum = g.NewEnum(
	"ScalingPolicy", "ScalingPolicies",
	"STANDARD", "ECONOMY",
)

var warehouseResourceConstraintEnum = g.NewEnum(
	"WarehouseResourceConstraint", "WarehouseResourceConstraints",
	"MEMORY_1X", "MEMORY_1X_x86", "MEMORY_16X", "MEMORY_16X_x86", "MEMORY_64X", "MEMORY_64X_x86",
)

var maxQueryPerformanceLevelEnum = g.NewEnum(
	"MaxQueryPerformanceLevel", "MaxQueryPerformanceLevels",
	"XSMALL", "SMALL", "MEDIUM", "LARGE", "XLARGE", "XXLARGE", "XXXLARGE", "X4LARGE",
).WithAliases("XSMALL", "X-SMALL").
	WithAliases("XLARGE", "X-LARGE").
	WithAliases("XXLARGE", "2X-LARGE").
	WithAliases("XXXLARGE", "3X-LARGE").
	WithAliases("X4LARGE", "4X-LARGE")

// ENABLED is reported for adaptive warehouses, which have no suspend/resume lifecycle.
var warehouseStateEnum = g.NewEnum(
	"WarehouseState", "WarehouseStates",
	"SUSPENDED", "SUSPENDING", "STARTED", "RESIZING", "RESUMING", "ENABLED",
)

var warehousePairs = g.StructPair("warehouseDBRow", "Warehouse").
	Text("name").
	Field("state", "string", "WarehouseState", g.WithManualConvert()).
	Field("type", "string", "WarehouseType", g.WithManualConvert()).
	OptionalEnum("size", warehouseSizeEnum).
	OptionalNumber("min_cluster_count").
	OptionalNumber("max_cluster_count").
	OptionalNumber("started_clusters").
	OptionalNumber("running").
	OptionalNumber("queued").
	BoolFromText("is_default").
	BoolFromText("is_current").
	OptionalNumber("auto_suspend").
	Bool("auto_resume").
	Field("available", "string", "float64", g.WithManualConvert()).
	Field("provisioning", "string", "float64", g.WithManualConvert()).
	Field("quiescing", "string", "float64", g.WithManualConvert()).
	Field("other", "string", "float64", g.WithManualConvert()).
	Time("created_on").
	Time("resumed_on").
	Time("updated_on").
	Text("owner").
	Text("comment").
	OptionalBool("enable_query_acceleration").
	OptionalNumber("query_acceleration_max_scale_factor").
	Field("resource_monitor", "string", "AccountObjectIdentifier", g.WithCustomParser("ParseAccountObjectIdentifierExcludingExplicitNullString")).
	Text("actives").
	Text("pendings").
	Text("failed").
	Text("suspended").
	Text("uuid").
	OptionalEnum("scaling_policy", scalingPolicyEnum).
	OptionalText("owner_role_type", g.WithRequiredInPlain()).
	OptionalEnum("resource_constraint", warehouseResourceConstraintEnum, g.WithManualConvert()).
	Field("generation", "sql.NullString", "*WarehouseGeneration", g.WithManualConvert()).
	OptionalEnum("max_query_performance_level", maxQueryPerformanceLevelEnum).
	OptionalNumber("query_throughput_multiplier").
	Field("tables", "sql.NullString", "[]SchemaObjectIdentifier", g.WithCustomParser("ParseCommaSeparatedSchemaObjectIdentifierArray"))

var warehouseDetailsPairs = g.StructPair("warehouseDetailsRow", "WarehouseDetails").
	Time("created_on").
	Text("name").
	Text("kind")

// SHOW projections for type-specific Terraform schemas. SHOW WAREHOUSES is one row type;
// Show() still returns []Warehouse. These structs are in-memory subsets (storage-integration
// AllDetails analog). Regular covers STANDARD and SNOWPARK-OPTIMIZED SHOW columns.
var warehouseRegularDef = g.PlainStruct("WarehouseRegular").
	Text("Name").
	Enum("State", warehouseStateEnum).
	Enum("Type", warehouseTypeEnum).
	OptionalField("Size", warehouseSizeEnum.Kind()).
	OptionalNumber("MinClusterCount").
	OptionalNumber("MaxClusterCount").
	OptionalNumber("StartedClusters").
	OptionalNumber("Running").
	OptionalNumber("Queued").
	Bool("IsDefault").
	Bool("IsCurrent").
	OptionalNumber("AutoSuspend").
	Bool("AutoResume").
	Field("Available", "float64").
	Field("Provisioning", "float64").
	Field("Quiescing", "float64").
	Field("Other", "float64").
	Time("CreatedOn").
	Time("ResumedOn").
	Time("UpdatedOn").
	Text("Owner").
	Text("Comment").
	OptionalBool("EnableQueryAcceleration").
	OptionalNumber("QueryAccelerationMaxScaleFactor").
	Field("ResourceMonitor", "AccountObjectIdentifier").
	OptionalField("ScalingPolicy", scalingPolicyEnum.Kind()).
	Text("OwnerRoleType").
	OptionalField("ResourceConstraint", warehouseResourceConstraintEnum.Kind()).
	OptionalField("Generation", "WarehouseGeneration")

var warehouseAdaptiveDef = g.PlainStruct("WarehouseAdaptive").
	Text("Name").
	Enum("State", warehouseStateEnum).
	Enum("Type", warehouseTypeEnum).
	OptionalNumber("Running").
	OptionalNumber("Queued").
	Bool("IsDefault").
	Bool("IsCurrent").
	Bool("AutoResume").
	Field("Available", "float64").
	Field("Provisioning", "float64").
	Field("Quiescing", "float64").
	Field("Other", "float64").
	Time("CreatedOn").
	Time("ResumedOn").
	Time("UpdatedOn").
	Text("Owner").
	Text("Comment").
	Field("ResourceMonitor", "AccountObjectIdentifier").
	Text("OwnerRoleType").
	OptionalField("MaxQueryPerformanceLevel", maxQueryPerformanceLevelEnum.Kind()).
	OptionalNumber("QueryThroughputMultiplier")

var warehouseInteractiveDef = g.PlainStruct("WarehouseInteractive").
	Text("Name").
	Enum("State", warehouseStateEnum).
	Enum("Type", warehouseTypeEnum).
	OptionalField("Size", warehouseSizeEnum.Kind()).
	OptionalNumber("MinClusterCount").
	OptionalNumber("MaxClusterCount").
	OptionalNumber("StartedClusters").
	OptionalNumber("Running").
	OptionalNumber("Queued").
	Bool("IsDefault").
	Bool("IsCurrent").
	OptionalNumber("AutoSuspend").
	Bool("AutoResume").
	Field("Available", "float64").
	Field("Provisioning", "float64").
	Field("Quiescing", "float64").
	Field("Other", "float64").
	Time("CreatedOn").
	Time("ResumedOn").
	Time("UpdatedOn").
	Text("Owner").
	Text("Comment").
	Field("ResourceMonitor", "AccountObjectIdentifier").
	Text("OwnerRoleType").
	Field("Tables", "[]SchemaObjectIdentifier")

var warehouseSetStruct = g.NewQueryStruct("WarehouseSet").
	OptionalEnumAssignment("WAREHOUSE_TYPE", warehouseTypeEnum, g.ParameterOptions().SingleQuotes()).
	OptionalEnumAssignment("WAREHOUSE_SIZE", warehouseSizeEnum, g.ParameterOptions().SingleQuotes()).
	OptionalBooleanAssignment("WAIT_FOR_COMPLETION", nil).
	OptionalNumberAssignment("MAX_CLUSTER_COUNT", g.ParameterOptions()).
	OptionalNumberAssignment("MIN_CLUSTER_COUNT", g.ParameterOptions()).
	OptionalEnumAssignment("SCALING_POLICY", scalingPolicyEnum, g.ParameterOptions().SingleQuotes()).
	OptionalNumberAssignment("AUTO_SUSPEND", g.ParameterOptions()).
	OptionalBooleanAssignment("AUTO_RESUME", nil).
	OptionalIdentifier("ResourceMonitor", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RESOURCE_MONITOR").Equals()).
	OptionalComment().
	OptionalBooleanAssignment("ENABLE_QUERY_ACCELERATION", nil).
	OptionalNumberAssignment("QUERY_ACCELERATION_MAX_SCALE_FACTOR", g.ParameterOptions()).
	OptionalEnumAssignment("RESOURCE_CONSTRAINT", warehouseResourceConstraintEnum, g.ParameterOptions().SingleQuotes()).
	WithField(g.OptionalEnumLegacy[sdkcommons.WarehouseGeneration]("Generation", g.ParameterOptions().SingleQuotes().SQL("GENERATION"))).
	OptionalNumberAssignment("QUERY_THROUGHPUT_MULTIPLIER", g.ParameterOptions()).
	OptionalEnumAssignment("MAX_QUERY_PERFORMANCE_LEVEL", maxQueryPerformanceLevelEnum, g.ParameterOptions().SingleQuotes()).
	OptionalNumberAssignment("MAX_CONCURRENCY_LEVEL", g.ParameterOptions()).
	OptionalNumberAssignment("STATEMENT_QUEUED_TIMEOUT_IN_SECONDS", g.ParameterOptions()).
	OptionalNumberAssignment("STATEMENT_TIMEOUT_IN_SECONDS", g.ParameterOptions()).
	OptionalIdentifier("FallbackWarehouse", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FALLBACK_WAREHOUSE").Equals()).
	WithValidation(g.AtLeastOneValueSet, "WarehouseType", "WarehouseSize", "WaitForCompletion", "MaxClusterCount", "MinClusterCount", "ScalingPolicy", "AutoSuspend", "AutoResume", "ResourceMonitor", "Comment", "EnableQueryAcceleration", "QueryAccelerationMaxScaleFactor", "ResourceConstraint", "Generation", "QueryThroughputMultiplier", "MaxQueryPerformanceLevel", "MaxConcurrencyLevel", "StatementQueuedTimeoutInSeconds", "StatementTimeoutInSeconds", "FallbackWarehouse").
	WithAdditionalValidations()

var warehouseUnsetStruct = g.NewQueryStruct("WarehouseUnset").
	OptionalSQL("WAREHOUSE_TYPE").
	OptionalSQL("WAIT_FOR_COMPLETION").
	OptionalSQL("MAX_CLUSTER_COUNT").
	OptionalSQL("MIN_CLUSTER_COUNT").
	OptionalSQL("SCALING_POLICY").
	OptionalSQL("AUTO_SUSPEND").
	OptionalSQL("AUTO_RESUME").
	OptionalSQL("RESOURCE_MONITOR").
	OptionalSQL("COMMENT").
	OptionalSQL("ENABLE_QUERY_ACCELERATION").
	OptionalSQL("QUERY_ACCELERATION_MAX_SCALE_FACTOR").
	OptionalSQL("RESOURCE_CONSTRAINT").
	OptionalSQL("GENERATION").
	OptionalSQL("MAX_CONCURRENCY_LEVEL").
	OptionalSQL("STATEMENT_QUEUED_TIMEOUT_IN_SECONDS").
	OptionalSQL("STATEMENT_TIMEOUT_IN_SECONDS").
	OptionalSQL("QUERY_THROUGHPUT_MULTIPLIER").
	OptionalSQL("MAX_QUERY_PERFORMANCE_LEVEL").
	OptionalSQL("FALLBACK_WAREHOUSE").
	WithValidation(g.AtLeastOneValueSet, "WarehouseType", "WaitForCompletion", "MaxClusterCount", "MinClusterCount", "ScalingPolicy", "AutoSuspend", "AutoResume", "ResourceMonitor", "Comment", "EnableQueryAcceleration", "QueryAccelerationMaxScaleFactor", "ResourceConstraint", "Generation", "MaxConcurrencyLevel", "StatementQueuedTimeoutInSeconds", "StatementTimeoutInSeconds", "QueryThroughputMultiplier", "MaxQueryPerformanceLevel", "FallbackWarehouse")

var warehousesDef = g.NewInterface(
	"Warehouses",
	"Warehouse",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-warehouse",
	g.NewQueryStruct("CreateWarehouse").
		Create().
		OrReplace().
		SQL("WAREHOUSE").
		IfNotExists().
		Name().
		OptionalEnumAssignment("WAREHOUSE_TYPE", warehouseTypeEnum, g.ParameterOptions().SingleQuotes()).
		OptionalEnumAssignment("WAREHOUSE_SIZE", warehouseSizeEnum, g.ParameterOptions().SingleQuotes()).
		OptionalNumberAssignment("MAX_CLUSTER_COUNT", g.ParameterOptions()).
		OptionalNumberAssignment("MIN_CLUSTER_COUNT", g.ParameterOptions()).
		OptionalEnumAssignment("SCALING_POLICY", scalingPolicyEnum, g.ParameterOptions().SingleQuotes()).
		OptionalNumberAssignment("AUTO_SUSPEND", g.ParameterOptions()).
		OptionalBooleanAssignment("AUTO_RESUME", nil).
		OptionalBooleanAssignment("INITIALLY_SUSPENDED", nil).
		OptionalIdentifier("ResourceMonitor", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RESOURCE_MONITOR").Equals()).
		OptionalComment().
		OptionalBooleanAssignment("ENABLE_QUERY_ACCELERATION", nil).
		OptionalNumberAssignment("QUERY_ACCELERATION_MAX_SCALE_FACTOR", g.ParameterOptions()).
		OptionalEnumAssignment("RESOURCE_CONSTRAINT", warehouseResourceConstraintEnum, g.ParameterOptions().SingleQuotes()).
		WithField(g.OptionalEnumLegacy[sdkcommons.WarehouseGeneration]("Generation", g.ParameterOptions().SingleQuotes().SQL("GENERATION"))).
		OptionalNumberAssignment("MAX_CONCURRENCY_LEVEL", g.ParameterOptions()).
		OptionalNumberAssignment("STATEMENT_QUEUED_TIMEOUT_IN_SECONDS", g.ParameterOptions()).
		OptionalNumberAssignment("STATEMENT_TIMEOUT_IN_SECONDS", g.ParameterOptions()).
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).CustomOperation(
	"CreateAdaptive",
	"https://docs.snowflake.com/en/user-guide/warehouses-adaptive",
	g.NewQueryStruct("CreateAdaptiveWarehouse").
		Create().
		OrReplace().
		SQL("WAREHOUSE").
		IfNotExists().
		Name().
		SQLWithCustomFieldName("warehouseType", "WAREHOUSE_TYPE = 'ADAPTIVE'").
		OptionalComment().
		OptionalEnumAssignment("MAX_QUERY_PERFORMANCE_LEVEL", maxQueryPerformanceLevelEnum, g.ParameterOptions().SingleQuotes()).
		OptionalNumberAssignment("QUERY_THROUGHPUT_MULTIPLIER", g.ParameterOptions()).
		OptionalIdentifier("ResourceMonitor", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RESOURCE_MONITOR").Equals()).
		OptionalTags().
		OptionalNumberAssignment("STATEMENT_QUEUED_TIMEOUT_IN_SECONDS", g.ParameterOptions()).
		OptionalNumberAssignment("STATEMENT_TIMEOUT_IN_SECONDS", g.ParameterOptions()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).CustomOperation(
	"CreateInteractive",
	"https://docs.snowflake.com/en/user-guide/warehouses-interactive",
	g.NewQueryStruct("CreateInteractiveWarehouse").
		Create().
		OrReplace().
		SQL("INTERACTIVE WAREHOUSE").
		IfNotExists().
		Name().
		ListAssignment("TABLES", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().Parentheses().NoEquals()).
		OptionalEnumAssignment("WAREHOUSE_SIZE", warehouseSizeEnum, g.ParameterOptions().SingleQuotes()).
		OptionalNumberAssignment("MAX_CLUSTER_COUNT", g.ParameterOptions()).
		OptionalNumberAssignment("MIN_CLUSTER_COUNT", g.ParameterOptions()).
		OptionalNumberAssignment("AUTO_SUSPEND", g.ParameterOptions()).
		OptionalBooleanAssignment("AUTO_RESUME", nil).
		OptionalBooleanAssignment("INITIALLY_SUSPENDED", nil).
		OptionalIdentifier("ResourceMonitor", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RESOURCE_MONITOR").Equals()).
		OptionalComment().
		OptionalNumberAssignment("MAX_CONCURRENCY_LEVEL", g.ParameterOptions()).
		OptionalNumberAssignment("STATEMENT_QUEUED_TIMEOUT_IN_SECONDS", g.ParameterOptions()).
		OptionalNumberAssignment("STATEMENT_TIMEOUT_IN_SECONDS", g.ParameterOptions()).
		OptionalIdentifier("FallbackWarehouse", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FALLBACK_WAREHOUSE").Equals()).
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-warehouse",
	g.NewQueryStruct("AlterWarehouse").
		Alter().
		SQL("WAREHOUSE").
		IfExists().
		Name().
		OptionalSQL("SUSPEND").
		OptionalSQL("RESUME").
		OptionalSQL("IF SUSPENDED").
		OptionalSQL("ABORT ALL QUERIES").
		RenameTo().
		OptionalQueryStructField("Set", warehouseSetStruct, g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("Unset", warehouseUnsetStruct, g.ListOptions().NoParentheses().SQL("UNSET")).
		ListAssignment("ADD TABLES", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().Parentheses().NoEquals()).
		ListAssignment("DROP TABLES", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().Parentheses().NoEquals()).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Suspend", "Resume", "AbortAllQueries", "RenameTo", "Set", "Unset", "AddTables", "DropTables", "SetTags", "UnsetTags").
		WithAdditionalValidations(),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-warehouse",
	g.NewQueryStruct("DropWarehouse").
		Drop().
		SQL("WAREHOUSE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-warehouses",
	warehousePairs,
	g.NewQueryStruct("ShowWarehouses").
		Show().
		SQL("WAREHOUSES").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimitFrom(),
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-warehouse",
	warehouseDetailsPairs,
	g.NewQueryStruct("DescribeWarehouse").
		Describe().
		SQL("WAREHOUSE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
	warehouseRegularDef,
	warehouseAdaptiveDef,
	warehouseInteractiveDef,
).ShowParameters("AccountObjectIdentifier").
	WithCustomInterfaceMethod(
		"ShowByIDExperimental", "ShowByIDExperimental is a show by id function with improved performance (using starts with and limit)",
		[]*g.MethodParameter{g.NewMethodParameter("id", "AccountObjectIdentifier")},
		"*Warehouse", "error",
	).WithCustomInterfaceMethod(
	"ShowByIDExperimentalSafely", "",
	[]*g.MethodParameter{g.NewMethodParameter("id", "AccountObjectIdentifier")},
	"*Warehouse", "error",
).WithCustomInterfaceMethod(
	"AlterWithSuspend", "AlterWithSuspend wraps Alter with automatic suspend/resume when changing warehouse type, or resizing an interactive warehouse",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*AlterWarehouseRequest")},
	"error",
).WithCustomInterfaceMethod(
	"CreateInteractivePreservingSession", "CreateInteractivePreservingSession wraps CreateInteractive to work around a known Snowflake behavior: creating a warehouse implicitly switches the session to it, which for interactive warehouses (5 second statement timeout) can cause subsequent statements in the same session to time out. It restores whichever warehouse (or lack of one) was active in the session before creation",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*CreateInteractiveWarehouseRequest")},
	"error",
).WithEnums(
	warehouseTypeEnum,
	warehouseSizeEnum,
	scalingPolicyEnum,
	warehouseResourceConstraintEnum,
	maxQueryPerformanceLevelEnum,
	warehouseStateEnum,
)
