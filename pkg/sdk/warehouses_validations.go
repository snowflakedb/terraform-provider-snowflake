package sdk

// ValidWarehouseSizesString is based on https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties
var ValidWarehouseSizesString = []string{
	string(WarehouseSizeXSmall),
	"X-SMALL",
	string(WarehouseSizeSmall),
	string(WarehouseSizeMedium),
	string(WarehouseSizeLarge),
	string(WarehouseSizeXLarge),
	"X-LARGE",
	string(WarehouseSizeXXLarge),
	"X2LARGE",
	"2X-LARGE",
	string(WarehouseSizeXXXLarge),
	"X3LARGE",
	"3X-LARGE",
	string(WarehouseSizeX4Large),
	"4X-LARGE",
	string(WarehouseSizeX5Large),
	"5X-LARGE",
	string(WarehouseSizeX6Large),
	"6X-LARGE",
}

// ValidWarehouseScalingPoliciesString is based on https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties
var ValidWarehouseScalingPoliciesString = []string{
	string(ScalingPolicyStandard),
	string(ScalingPolicyEconomy),
}

// ValidWarehouseTypesString is based on https://docs.snowflake.com/en/sql-reference/sql/create-warehouse#optional-properties-objectproperties
var ValidWarehouseTypesString = []string{
	string(WarehouseTypeStandard),
	string(WarehouseTypeSnowparkOptimized),
}

// ValidResourceConstraintsString is based on https://docs.snowflake.com/en/user-guide/warehouses-gen2
var ValidResourceConstraintsString = []string{
	string(ResourceConstraintStandardGen1),
	string(ResourceConstraintStandardGen2),
}

// WarehouseParameters is based on https://docs.snowflake.com/en/sql-reference/parameters#object-parameters
var WarehouseParameters = []ObjectParameter{
	ObjectParameterMaxConcurrencyLevel,
	ObjectParameterStatementQueuedTimeoutInSeconds,
	ObjectParameterStatementTimeoutInSeconds,
}
