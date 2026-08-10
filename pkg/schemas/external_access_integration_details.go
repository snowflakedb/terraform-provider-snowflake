package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeExternalAccessIntegrationDetailsSchema represents output of DESCRIBE query for the single ExternalAccessIntegration.
var DescribeExternalAccessIntegrationDetailsSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
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
	"enabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ExternalAccessIntegrationDetailsToSchema(details sdk.ExternalAccessIntegrationDetails) map[string]any {
	networkRules := make([]string, len(details.AllowedNetworkRules))
	for i, v := range details.AllowedNetworkRules {
		networkRules[i] = v.FullyQualifiedName()
	}
	return map[string]any{
		"id":                    details.Id.Name(),
		"allowed_network_rules": networkRules,
		"allowed_api_authentication_integrations": details.AllowedApiAuthenticationIntegrations,
		"allowed_authentication_secrets":          details.AllowedAuthenticationSecrets,
		"enabled":                                 details.Enabled,
		"comment":                                 details.Comment,
	}
}
