package main

import (
	"fmt"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func HandleWarehouses(config *Config, csvInput [][]string) (string, error) {
	return HandleResources[WarehouseCsvRow, WarehouseRepresentation](config, csvInput, MapWarehouseToModel)
}

func MapWarehouseToModel(warehouse WarehouseRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	warehouseId := sdk.NewAccountObjectIdentifier(warehouse.Name)
	resourceId := NormalizeResourceId(fmt.Sprintf("warehouse_%s", warehouseId.FullyQualifiedName()))
	resourceModel := model.Warehouse(resourceId, warehouse.Name)

	handleOptionalFieldWithStringBuilder(warehouse.Comment, resourceModel.WithComment)
	handleOptionalFieldWithStringBuilder(warehouse.AutoResume, resourceModel.WithAutoResume)
	handleOptionalFieldWithBuilder(warehouse.AutoSuspend, resourceModel.WithAutoSuspend)
	handleOptionalFieldWithStringBuilder(warehouse.EnableQueryAcceleration, resourceModel.WithEnableQueryAcceleration)
	handleOptionalFieldWithBuilder(warehouse.Generation, resourceModel.WithGeneration)
	handleOptionalFieldWithBuilder(warehouse.InitiallySuspended, resourceModel.WithInitiallySuspended)
	handleOptionalFieldWithBuilder(warehouse.MaxClusterCount, resourceModel.WithMaxClusterCount)
	handleOptionalFieldWithBuilder(warehouse.MinClusterCount, resourceModel.WithMinClusterCount)
	handleOptionalFieldWithBuilder(warehouse.QueryAccelerationMaxScaleFactor, resourceModel.WithQueryAccelerationMaxScaleFactor)
	handleOptionalFieldWithBuilder(warehouse.ResourceConstraint, resourceModel.WithResourceConstraint)
	handleOptionalFieldWithBuilder(warehouse.ResourceMonitor, resourceModel.WithResourceMonitor)
	handleOptionalFieldWithBuilder(warehouse.ScalingPolicy, resourceModel.WithScalingPolicy)
	handleOptionalFieldWithBuilder(warehouse.WarehouseSize, resourceModel.WithWarehouseSize)
	handleOptionalFieldWithBuilder(warehouse.WarehouseType, resourceModel.WithWarehouseType)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		warehouseId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}
