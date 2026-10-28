package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

func All() []ssot.Object {
	objs := []ssot.Object{
		ComputePool,
		ExternalVolumeObject,
		Schema,
		AuthenticationPolicy,
		CatalogIntegration,
		CatalogIntegrationAwsGlue,
		CatalogIntegrationIcebergRest,
		CatalogIntegrationObjectStorage,
		CatalogIntegrationOpenCatalog,
		Warehouse,
		WarehouseAdaptive,
		WarehouseInteractive,
		Alert,
		Table,
		Pipe,
		CodeBundle,
		BackupPolicy,
		ProjectionPolicy,
	}
	for _, obj := range objs {
		ssot.ValidateObject(obj)
	}
	return objs
}

// VariantsOf returns objects whose VariantOf name is base.Name, in All order.
func VariantsOf(base ssot.Object) []ssot.Object {
	var out []ssot.Object
	for _, obj := range All() {
		if obj.VariantOf != nil && obj.VariantOf.Name == base.Name {
			out = append(out, obj)
		}
	}
	return out
}
