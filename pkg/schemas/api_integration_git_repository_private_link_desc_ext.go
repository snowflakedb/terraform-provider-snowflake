package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (apiIntegrationGitRepositoryPrivateLinkToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return collections.MergeMaps(
		apiIntegrationApiProviderSchema(),
		apiIntegrationAllowedBlockedPrefixesSchema(),
		apiIntegrationTlsTrustedCertificatesSchema(),
	)
}

func (apiIntegrationGitRepositoryPrivateLinkToSchemaMapper) additionalToSchema(src *sdk.ApiIntegrationGitRepositoryPrivateLink, dst map[string]any) {
	apiIntegrationApiProviderToSchema(src.ApiProvider, dst)
	apiIntegrationAllowedBlockedPrefixesToSchema(src.AllowedPrefixes, src.BlockedPrefixes, dst)
	dst["tls_trusted_certificates"] = src.TlsTrustedCertificates
}
