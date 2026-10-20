package schemas

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

// SHOW is one sdk.CatalogIntegration row for every catalog source. The resource schema
// names the show block ShowCatalogIntegrationAwsGlueSchema.
var ShowCatalogIntegrationAwsGlueSchema = ShowCatalogIntegrationSchema

func CatalogIntegrationAwsGlueToSchema(catalogIntegration *sdk.CatalogIntegration) map[string]any {
	return CatalogIntegrationToSchema(catalogIntegration)
}
