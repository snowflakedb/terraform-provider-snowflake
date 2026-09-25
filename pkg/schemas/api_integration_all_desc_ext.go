package schemas

import (
	"maps"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO [next PRs]: first move — []string prefixes/scopes/certs stay here until MapToSchemaField maps slices natively.
// api_key is Sensitive and api_provider needs strings.ToLower; neither can be generated.

func apiIntegrationMergeAdditional(parts ...map[string]*schema.Schema) map[string]*schema.Schema {
	out := make(map[string]*schema.Schema)
	for _, part := range parts {
		maps.Copy(out, part)
	}
	return out
}

func apiIntegrationApiKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_key": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func apiIntegrationApiProviderSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_provider": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func apiIntegrationAllowedBlockedPrefixesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allowed_prefixes": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"blocked_prefixes": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func apiIntegrationOauthAllowedScopesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"oauth_allowed_scopes": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func apiIntegrationTlsTrustedCertificatesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tls_trusted_certificates": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func apiIntegrationApiKeyToSchema(apiKey string, dst map[string]any) {
	dst["api_key"] = apiKey
}

func apiIntegrationApiProviderToSchema(apiProvider string, dst map[string]any) {
	dst["api_provider"] = strings.ToLower(apiProvider)
}

func apiIntegrationAllowedBlockedPrefixesToSchema(allowed, blocked []string, dst map[string]any) {
	dst["allowed_prefixes"] = allowed
	dst["blocked_prefixes"] = blocked
}

func (apiIntegrationAllDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return apiIntegrationMergeAdditional(
		apiIntegrationApiKeySchema(),
		apiIntegrationApiProviderSchema(),
		apiIntegrationAllowedBlockedPrefixesSchema(),
		apiIntegrationOauthAllowedScopesSchema(),
		apiIntegrationTlsTrustedCertificatesSchema(),
	)
}

func (apiIntegrationAllDetailsToSchemaMapper) additionalToSchema(src *sdk.ApiIntegrationAllDetails, dst map[string]any) {
	apiIntegrationApiKeyToSchema(src.ApiKey, dst)
	apiIntegrationApiProviderToSchema(src.ApiProvider, dst)
	apiIntegrationAllowedBlockedPrefixesToSchema(src.AllowedPrefixes, src.BlockedPrefixes, dst)
	dst["oauth_allowed_scopes"] = src.OauthAllowedScopes
	dst["tls_trusted_certificates"] = src.TlsTrustedCertificates
}
