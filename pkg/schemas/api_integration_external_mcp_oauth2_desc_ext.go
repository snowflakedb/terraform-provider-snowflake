package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (apiIntegrationExternalMcpOauth2ToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return apiIntegrationMergeAdditional(
		apiIntegrationApiProviderSchema(),
		apiIntegrationAllowedBlockedPrefixesSchema(),
		apiIntegrationOauthAllowedScopesSchema(),
	)
}

func (apiIntegrationExternalMcpOauth2ToSchemaMapper) additionalToSchema(src *sdk.ApiIntegrationExternalMcpOauth2, dst map[string]any) {
	apiIntegrationApiProviderToSchema(src.ApiProvider, dst)
	apiIntegrationAllowedBlockedPrefixesToSchema(src.AllowedPrefixes, src.BlockedPrefixes, dst)
	dst["oauth_allowed_scopes"] = src.OauthAllowedScopes
}
