package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

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
			OptionalEnumAssignment("MAX_QUERY_PERFORMANCE_LEVEL", ssotdefs.MaxQueryPerformanceLevel, g.ParameterOptions().SingleQuotes()).
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
		OptionalEnumAssignment("WAREHOUSE_SIZE", ssotdefs.WarehouseSize, g.ParameterOptions().SingleQuotes()).
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
	DescribeOperationFromObject(ssotdefs.Warehouse).
	ShowParametersFromObject(ssotdefs.Warehouse).
	ShowParametersDetailsFromObject(ssotdefs.Warehouse).
	CustomInterfaceMethodsFromObject(ssotdefs.Warehouse).
	WithCustomInterfaceMethod(
		"CreateInteractivePreservingSession", "CreateInteractivePreservingSession wraps CreateInteractive to work around a known Snowflake behavior: creating a warehouse implicitly switches the session to it, which for interactive warehouses (5 second statement timeout) can cause subsequent statements in the same session to time out. It restores whichever warehouse (or lack of one) was active in the session before creation",
		[]*g.MethodParameter{g.NewMethodParameter("request", "*CreateInteractiveWarehouseRequest")},
		"error",
	).WithEnums(
	ssotdefs.WarehouseType,
	ssotdefs.WarehouseSize,
	ssotdefs.ScalingPolicy,
	ssotdefs.WarehouseResourceConstraint,
	ssotdefs.MaxQueryPerformanceLevel,
	ssotdefs.WarehouseState,
)
