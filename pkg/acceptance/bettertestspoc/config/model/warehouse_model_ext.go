package model

import (
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func BasicWarehouseModel(
	id sdk.AccountObjectIdentifier,
	comment string,
) *WarehouseModel {
	return WarehouseWithDefaultMeta(id.Name()).WithComment(comment)
}

func WarehouseSnowflakeDefaultWithoutParameters(
	id sdk.AccountObjectIdentifier,
	comment string,
) *WarehouseModel {
	return BasicWarehouseModel(id, comment).
		WithWarehouseTypeEnum(sdk.WarehouseTypeStandard).
		WithWarehouseSizeEnum(sdk.WarehouseSizeXSmall).
		WithMinClusterCount(1).
		WithMaxClusterCount(1).
		WithScalingPolicyEnum(sdk.ScalingPolicyStandard).
		WithAutoSuspend(600).
		WithAutoResume(r.BooleanTrue).
		WithInitiallySuspended(false).
		WithEnableQueryAcceleration(r.BooleanFalse).
		WithQueryAccelerationMaxScaleFactor(8)
}

// TODO [SNOW-1501905]: currently config builder are generated from the resource schema, so there is no direct connection to the source enum (like sdk.WarehouseSize)
// For now, we can just add extension methods manually.
// Later, we could provide type overrides map or even SDK object to automatically match by name.
func (w *WarehouseModel) WithWarehouseSizeEnum(warehouseSize sdk.WarehouseSize) *WarehouseModel {
	return w.WithWarehouseSize(string(warehouseSize))
}

func (w *WarehouseModel) WithWarehouseTypeEnum(warehouseType sdk.WarehouseType) *WarehouseModel {
	return w.WithWarehouseType(string(warehouseType))
}

func (w *WarehouseModel) WithScalingPolicyEnum(scalingPolicy sdk.ScalingPolicy) *WarehouseModel {
	return w.WithScalingPolicy(string(scalingPolicy))
}

func (w *WarehouseModel) WithResourceConstraintEnum(resourceConstraint sdk.WarehouseResourceConstraint) *WarehouseModel {
	return w.WithResourceConstraint(string(resourceConstraint))
}
