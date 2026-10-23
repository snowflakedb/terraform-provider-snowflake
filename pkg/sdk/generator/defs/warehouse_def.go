package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
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

func warehouseShowIdentity(ps *g.Plain) *g.Plain {
	return ps.
		Text("Name").
		Enum("State", warehouseStateEnum).
		Enum("Type", warehouseTypeEnum)
}

func warehouseShowLoad(ps *g.Plain) *g.Plain {
	return ps.
		OptionalNumber("Running").
		OptionalNumber("Queued").
		Bool("IsDefault").
		Bool("IsCurrent")
}

func warehouseShowAvailability(ps *g.Plain) *g.Plain {
	return ps.
		Bool("AutoResume").
		Field("Available", "float64").
		Field("Provisioning", "float64").
		Field("Quiescing", "float64").
		Field("Other", "float64").
		Time("CreatedOn").
		Time("ResumedOn").
		Time("UpdatedOn").
		Text("Owner").
		Text("Comment")
}

// SHOW projections for type-specific Terraform schemas. SHOW WAREHOUSES is one row type;
// Show() still returns []Warehouse. These structs are in-memory subsets (storage-integration
// AllDetails analog). Regular covers STANDARD and SNOWPARK-OPTIMIZED SHOW columns.
var warehouseRegularDef = warehouseShowAvailability(
	warehouseShowLoad(
		warehouseShowIdentity(g.PlainStruct("WarehouseRegular")).
			OptionalField("Size", warehouseSizeEnum.Kind()).
			OptionalNumber("MinClusterCount").
			OptionalNumber("MaxClusterCount").
			OptionalNumber("StartedClusters"),
	).
		OptionalNumber("AutoSuspend"),
).
	OptionalBool("EnableQueryAcceleration").
	OptionalNumber("QueryAccelerationMaxScaleFactor").
	Field("ResourceMonitor", "AccountObjectIdentifier").
	OptionalField("ScalingPolicy", scalingPolicyEnum.Kind()).
	Text("OwnerRoleType").
	OptionalField("ResourceConstraint", warehouseResourceConstraintEnum.Kind()).
	OptionalField("Generation", "WarehouseGeneration")

var warehouseAdaptiveDef = warehouseShowAvailability(
	warehouseShowLoad(
		warehouseShowIdentity(g.PlainStruct("WarehouseAdaptive")),
	),
).
	Field("ResourceMonitor", "AccountObjectIdentifier").
	Text("OwnerRoleType").
	OptionalField("MaxQueryPerformanceLevel", maxQueryPerformanceLevelEnum.Kind()).
	OptionalNumber("QueryThroughputMultiplier")

var warehouseInteractiveDef = warehouseShowAvailability(
	warehouseShowLoad(
		warehouseShowIdentity(g.PlainStruct("WarehouseInteractive")).
			OptionalField("Size", warehouseSizeEnum.Kind()).
			OptionalNumber("MinClusterCount").
			OptionalNumber("MaxClusterCount").
			OptionalNumber("StartedClusters"),
	).
		OptionalNumber("AutoSuspend"),
).
	Field("ResourceMonitor", "AccountObjectIdentifier").
	Text("OwnerRoleType").
	Field("Tables", "[]SchemaObjectIdentifier")

var warehousesDef = g.NewInterface(
	"Warehouses",
	"Warehouse",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperationFromObject(ssotdefs.Warehouse).
	CustomOperation(
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
		WithValidation(g.ValidIdentifier, "Tables").
		WithAdditionalValidations(),
).AlterOperationFromObject(ssotdefs.Warehouse).
	DropOperationFromObject(ssotdefs.Warehouse).
	ShowOperationFromObject(ssotdefs.Warehouse).
	DescribeOperationFromObject(ssotdefs.Warehouse, warehouseRegularDef, warehouseAdaptiveDef, warehouseInteractiveDef).
	ShowParametersFromObject(ssotdefs.Warehouse).
	ShowParametersDetailsFromObject(ssotdefs.Warehouse).
	CustomInterfaceMethodsFromObject(ssotdefs.Warehouse).
	WithCustomInterfaceMethod(
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
