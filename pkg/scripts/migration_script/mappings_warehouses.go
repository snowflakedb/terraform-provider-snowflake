package main

import (
	"errors"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func HandleWarehouses(config *Config, csvInput [][]string) (string, error) {
	return HandleResources[WarehouseCsvRow, WarehouseRepresentation](config, csvInput, MapWarehouseToModel)
}

func MapWarehouseToModel(warehouse WarehouseRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	warehouseId := sdk.NewAccountObjectIdentifier(warehouse.Name)
	resourceId := ResourceId(resources.Warehouse, warehouseId.FullyQualifiedName())
	resourceModel := model.Warehouse(resourceId, warehouse.Name)

	// always include fields with default values
	handleBooleanString(warehouse.AutoResume, resourceModel.WithAutoResume)
	if warehouse.Type == sdk.WarehouseTypeAdaptive {
		return nil, nil, errors.New("adaptive warehouses are not supported")
	}
	resourceModel.WithWarehouseTypeEnum(warehouse.Type)
	handleIfNotNil(warehouse.Size, resourceModel.WithWarehouseSizeEnum)
	handleIfNotNil(warehouse.ScalingPolicy, resourceModel.WithScalingPolicyEnum)
	handleIfNotNil(warehouse.AutoSuspend, resourceModel.WithAutoSuspend)
	handleIfNotNil(warehouse.MinClusterCount, resourceModel.WithMinClusterCount)
	handleIfNotNil(warehouse.MaxClusterCount, resourceModel.WithMaxClusterCount)
	handleIfNotNil(warehouse.QueryAccelerationMaxScaleFactor, resourceModel.WithQueryAccelerationMaxScaleFactor)
	handleIfNotEmpty(warehouse.Comment, resourceModel.WithComment)
	if warehouse.EnableQueryAcceleration != nil {
		handleIf(*warehouse.EnableQueryAcceleration, resourceModel.WithEnableQueryAcceleration)
	}
	handleIfNotEmpty(warehouse.ResourceMonitor.Name(), resourceModel.WithResourceMonitor)
	handleOptionalFieldWithBuilder(warehouse.Generation, resourceModel.WithGenerationEnum)
	handleOptionalFieldWithBuilder(warehouse.ResourceConstraint, resourceModel.WithResourceConstraintEnum)

	handleOptionalFieldWithBuilder(warehouse.MaxConcurrencyLevel, resourceModel.WithMaxConcurrencyLevel)
	handleOptionalFieldWithBuilder(warehouse.StatementQueuedTimeoutInSeconds, resourceModel.WithStatementQueuedTimeoutInSeconds)
	handleOptionalFieldWithBuilder(warehouse.StatementTimeoutInSeconds, resourceModel.WithStatementTimeoutInSeconds)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		warehouseId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}
