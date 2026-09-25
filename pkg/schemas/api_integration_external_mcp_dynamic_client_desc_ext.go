package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (apiIntegrationExternalMcpDynamicClientToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return apiIntegrationMergeAdditional(
		apiIntegrationApiProviderSchema(),
		apiIntegrationAllowedBlockedPrefixesSchema(),
	)
}

func (apiIntegrationExternalMcpDynamicClientToSchemaMapper) additionalToSchema(src *sdk.ApiIntegrationExternalMcpDynamicClient, dst map[string]any) {
	apiIntegrationApiProviderToSchema(src.ApiProvider, dst)
	apiIntegrationAllowedBlockedPrefixesToSchema(src.AllowedPrefixes, src.BlockedPrefixes, dst)
}
