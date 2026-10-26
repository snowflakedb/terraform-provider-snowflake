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
	AlterVariantOperationsFromObject(ssotdefs.Warehouse).
	DropOperationFromObject(ssotdefs.Warehouse).
	ShowOperationFromObject(ssotdefs.Warehouse).
	DescribeOperationFromObject(ssotdefs.Warehouse).
	ShowParametersFromObject(ssotdefs.Warehouse).
	ShowParametersDetailsFromObject(ssotdefs.Warehouse).
	CustomInterfaceMethodsFromObject(ssotdefs.Warehouse).
	WithEnums(
		ssotdefs.WarehouseType,
		ssotdefs.WarehouseSize,
		ssotdefs.ScalingPolicy,
		ssotdefs.WarehouseResourceConstraint,
		ssotdefs.MaxQueryPerformanceLevel,
		ssotdefs.WarehouseState,
	)
