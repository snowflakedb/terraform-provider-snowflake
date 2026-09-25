package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (apiIntegrationGitRepositoryOauth2ToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return apiIntegrationMergeAdditional(
		apiIntegrationAllowedBlockedPrefixesSchema(),
		apiIntegrationOauthAllowedScopesSchema(),
	)
}

func (apiIntegrationGitRepositoryOauth2ToSchemaMapper) additionalToSchema(src *sdk.ApiIntegrationGitRepositoryOauth2, dst map[string]any) {
	apiIntegrationAllowedBlockedPrefixesToSchema(src.AllowedPrefixes, src.BlockedPrefixes, dst)
	dst["oauth_allowed_scopes"] = src.OauthAllowedScopes
}
