package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (apiIntegrationAwsDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return apiIntegrationMergeAdditional(
		apiIntegrationApiKeySchema(),
		apiIntegrationApiProviderSchema(),
		apiIntegrationAllowedBlockedPrefixesSchema(),
	)
}

func (apiIntegrationAwsDetailsToSchemaMapper) additionalToSchema(src *sdk.ApiIntegrationAwsDetails, dst map[string]any) {
	apiIntegrationApiKeyToSchema(src.ApiKey, dst)
	apiIntegrationApiProviderToSchema(src.ApiProvider, dst)
	apiIntegrationAllowedBlockedPrefixesToSchema(src.AllowedPrefixes, src.BlockedPrefixes, dst)
}
