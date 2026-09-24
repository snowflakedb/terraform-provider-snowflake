package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (externalAccessIntegrationDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allowed_network_rules": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"allowed_api_authentication_integrations": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"allowed_authentication_secrets": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (externalAccessIntegrationDetailsToSchemaMapper) additionalToSchema(src *sdk.ExternalAccessIntegrationDetails, dst map[string]any) {
	networkRules := make([]string, len(src.AllowedNetworkRules))
	for i, v := range src.AllowedNetworkRules {
		networkRules[i] = v.FullyQualifiedName()
	}
	dst["allowed_network_rules"] = networkRules
	dst["allowed_api_authentication_integrations"] = src.AllowedApiAuthenticationIntegrations
	dst["allowed_authentication_secrets"] = src.AllowedAuthenticationSecrets
}
