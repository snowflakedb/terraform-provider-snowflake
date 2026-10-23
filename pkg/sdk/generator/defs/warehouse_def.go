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
	CreateVariantOperationsFromObject(ssotdefs.Warehouse).
	AlterOperationFromObject(ssotdefs.Warehouse).
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
