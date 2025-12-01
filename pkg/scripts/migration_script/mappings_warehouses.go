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

	handleIfNotEmpty(warehouse.Comment, resourceModel.WithComment)
	handleOptionalFieldWithBuilder(warehouse.MaxConcurrencyLevel, resourceModel.WithMaxConcurrencyLevel)
	handleOptionalFieldWithBuilder(warehouse.StatementQueuedTimeoutInSeconds, resourceModel.WithStatementQueuedTimeoutInSeconds)
	handleOptionalFieldWithBuilder(warehouse.StatementTimeoutInSeconds, resourceModel.WithStatementTimeoutInSeconds)

	importModel := NewImportModel(
		resourceModel.ResourceReference(),
		warehouseId.FullyQualifiedName(),
	)

	return resourceModel, importModel, nil
}
