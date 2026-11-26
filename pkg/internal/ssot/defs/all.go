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
		BackupSet,
		AggregationPolicy,
		ProjectionPolicy,
		JoinPolicy,
		PostgresMirror,
		Share,
		ReplicationGroup,
		ReplicationGroupSecondary,
		FailoverGroup,
		FailoverGroupSecondary,
	}
	for _, obj := range objs {
		ssot.ValidateObject(obj)
	}
	return objs
}

// VariantsOf returns objects whose Variant.Of name is base.Name, in All order.
func VariantsOf(base ssot.Object) []ssot.Object {
	var out []ssot.Object
	for _, obj := range All() {
		if baseOf := obj.VariantBase(); baseOf != nil && baseOf.Name == base.Name {
			out = append(out, obj)
		}
	}
	return out
}
