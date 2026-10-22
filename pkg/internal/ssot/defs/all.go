package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

func All() []ssot.Object {
	objs := []ssot.Object{
		ComputePool,
		ExternalVolumeObject,
		Schema,
		AuthenticationPolicy,
		CatalogIntegrationAwsGlue,
		CatalogIntegrationIcebergRest,
		CatalogIntegrationObjectStorage,
		CatalogIntegrationOpenCatalog,
		Warehouse,
		WarehouseAdaptive,
		WarehouseInteractive,
		Alert,
		Pipe,
	}
	for _, obj := range objs {
		ssot.ValidateObject(obj)
	}
	return objs
}
