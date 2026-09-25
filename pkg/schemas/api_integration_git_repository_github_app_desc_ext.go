package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (apiIntegrationGitRepositoryGithubAppToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return apiIntegrationMergeAdditional(
		apiIntegrationApiProviderSchema(),
		apiIntegrationAllowedBlockedPrefixesSchema(),
	)
}

func (apiIntegrationGitRepositoryGithubAppToSchemaMapper) additionalToSchema(src *sdk.ApiIntegrationGitRepositoryGithubApp, dst map[string]any) {
	apiIntegrationApiProviderToSchema(src.ApiProvider, dst)
	apiIntegrationAllowedBlockedPrefixesToSchema(src.AllowedPrefixes, src.BlockedPrefixes, dst)
}
