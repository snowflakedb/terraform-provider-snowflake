package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var computePoolsDef = g.NewInterface(
	"ComputePools",
	"ComputePool",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperationFromObject(ssotdefs.ComputePool).
	AlterOperationFromObject(ssotdefs.ComputePool).
	DropOperationFromObject(ssotdefs.ComputePool).
	ShowOperationFromObject(ssotdefs.ComputePool).
	DescribeOperationFromObject(ssotdefs.ComputePool).
	WithEnums(
		ssotdefs.ComputePoolState,
		ssotdefs.ComputePoolInstanceFamily,
	)
