package schemas

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func CatalogIntegrationIcebergRestToSchema(catalogIntegration *sdk.CatalogIntegration) map[string]any {
	return CatalogIntegrationToSchema(catalogIntegration)
}
