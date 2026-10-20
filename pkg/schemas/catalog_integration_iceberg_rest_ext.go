package schemas

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func CatalogIntegrationIcebergRestToSchema(catalogIntegration *sdk.CatalogIntegration) map[string]any {
	return CatalogIntegrationToSchema(catalogIntegration)
}

// Show rows are the shared sdk.CatalogIntegration type. The generated resource schema
// looks up that block as ShowCatalogIntegrationIcebergRestSchema.
var ShowCatalogIntegrationIcebergRestSchema = ShowCatalogIntegrationSchema
