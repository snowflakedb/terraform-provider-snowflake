package main

import (
	"fmt"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// Default warehouse values
const (
	defaultAutoSuspend                     = 600
	defaultMinClusterCount                 = 1
	defaultMaxClusterCount                 = 1
	defaultQueryAccelerationMaxScaleFactor = 8
)

func HandleWarehouses(config *Config, csvInput [][]string) (string, error) {
	return HandleResources[WarehouseCsvRow, WarehouseRepresentation](config, csvInput, MapWarehouseToModel)
}

func MapWarehouseToModel(warehouse WarehouseRepresentation) (accconfig.ResourceModel, *ImportModel, error) {
	warehouseId := sdk.NewAccountObjectIdentifier(warehouse.Name)
	resourceId := NormalizeResourceId(fmt.Sprintf("warehouse_%s", warehouseId.FullyQualifiedName()))
	resourceModel := model.Warehouse(resourceId, warehouse.Name)

	handleIfNotEmpty(warehouse.Comment, resourceModel.WithComment)

	if warehouse.Type != sdk.WarehouseTypeStandard {
		resourceModel.WithWarehouseTypeEnum(warehouse.Type)
	}

	if warehouse.Size != sdk.WarehouseSizeXSmall {
		resourceModel.WithWarehouseSizeEnum(warehouse.Size)
	}

	if warehouse.AutoSuspend != defaultAutoSuspend {
		resourceModel.WithAutoSuspend(warehouse.AutoSuspend)
	}

	if !warehouse.AutoResume {
		resourceModel.WithAutoResume(r.BooleanFalse)
	}

	if warehouse.MinClusterCount != defaultMinClusterCount {
		resourceModel.WithMinClusterCount(warehouse.MinClusterCount)
	}
	if warehouse.MaxClusterCount != defaultMaxClusterCount {
		resourceModel.WithMaxClusterCount(warehouse.MaxClusterCount)
	}

	if warehouse.ScalingPolicy != sdk.ScalingPolicyStandard {
		resourceModel.WithScalingPolicyEnum(warehouse.ScalingPolicy)
	}

	if warehouse.EnableQueryAcceleration {
		resourceModel.WithEnableQueryAcceleration(r.BooleanTrue)
	}

	if warehouse.QueryAccelerationMaxScaleFactor != defaultQueryAccelerationMaxScaleFactor {
		resourceModel.WithQueryAccelerationMaxScaleFactor(warehouse.QueryAccelerationMaxScaleFactor)
	}

	if warehouse.ResourceMonitor.Name() != "" {
		resourceModel.WithResourceMonitor(warehouse.ResourceMonitor.Name())
	}

	handleOptionalFieldWithBuilder(warehouse.MaxConcurrencyLevel, resourceModel.WithMaxConcurrencyLevel)
	handleOptionalFieldWithBuilder(warehouse.StatementQueuedTimeoutInSeconds, resourceModel.WithStatementQueuedTimeoutInSeconds)
	handleOptionalFieldWithBuilder(warehouse.StatementTimeoutInSeconds, resourceModel.WithStatementTimeoutInSeconds)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		warehouseId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}
