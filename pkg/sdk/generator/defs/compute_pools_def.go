package defs

import (
	ssotdefs "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var ComputePoolBackupInstanceFamilyListItemDef = g.NewQueryStruct("ComputePoolBackupInstanceFamilyListItem").
	Enum("Value", ssotdefs.ComputePoolInstanceFamily, g.KeywordOptions().SingleQuotes().Required())

var computePoolsDef = g.NewInterface(
	"ComputePools",
	"ComputePool",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperationFromObject(ssotdefs.ComputePool, ComputePoolBackupInstanceFamilyListItemDef).
	AlterOperationFromObject(ssotdefs.ComputePool).
	DropOperationFromObject(ssotdefs.ComputePool).
	ShowOperationFromObject(ssotdefs.ComputePool).
	DescribeOperationFromObject(ssotdefs.ComputePool).
	WithEnums(
		ssotdefs.ComputePoolState,
		ssotdefs.ComputePoolInstanceFamily,
	)
