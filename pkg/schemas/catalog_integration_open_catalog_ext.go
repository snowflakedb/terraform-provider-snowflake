package schemas

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func CatalogIntegrationOpenCatalogToSchema(catalogIntegration *sdk.CatalogIntegration) map[string]any {
	return CatalogIntegrationToSchema(catalogIntegration)
}

// Show rows are the shared sdk.CatalogIntegration type. The generated resource schema
// looks up that block as ShowCatalogIntegrationOpenCatalogSchema.
var ShowCatalogIntegrationOpenCatalogSchema = ShowCatalogIntegrationSchema
