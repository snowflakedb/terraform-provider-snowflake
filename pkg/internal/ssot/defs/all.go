package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

func All() []ssot.Object {
	objs := []ssot.Object{
		ComputePool,
		ExternalVolumeObject,
		Schema,
		AuthenticationPolicy,
		CatalogIntegrationOpenCatalog,
		Warehouse,
	}
	for _, obj := range objs {
		ssot.ValidateObject(obj)
	}
	return objs
}
