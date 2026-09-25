package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (apiIntegrationGoogleDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return apiIntegrationMergeAdditional(
		apiIntegrationApiKeySchema(),
		apiIntegrationApiProviderSchema(),
		apiIntegrationAllowedBlockedPrefixesSchema(),
	)
}

func (apiIntegrationGoogleDetailsToSchemaMapper) additionalToSchema(src *sdk.ApiIntegrationGoogleDetails, dst map[string]any) {
	apiIntegrationApiKeyToSchema(src.ApiKey, dst)
	apiIntegrationApiProviderToSchema(src.ApiProvider, dst)
	apiIntegrationAllowedBlockedPrefixesToSchema(src.AllowedPrefixes, src.BlockedPrefixes, dst)
}
