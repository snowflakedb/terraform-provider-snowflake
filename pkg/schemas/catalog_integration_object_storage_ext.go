package schemas

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func CatalogIntegrationObjectStorageToSchema(catalogIntegration *sdk.CatalogIntegration) map[string]any {
	return CatalogIntegrationToSchema(catalogIntegration)
}
