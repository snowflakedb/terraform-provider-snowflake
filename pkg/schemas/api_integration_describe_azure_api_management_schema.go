package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeAzureApiManagementApiIntegrationSchema = map[string]*schema.Schema{
	"enabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"api_key": {
		Type:      schema.TypeString,
		Computed:  true,
		Sensitive: true,
	},
	"api_provider": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"azure_tenant_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"azure_ad_application_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"azure_multi_tenant_app_name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"azure_consent_url": {
		Type:     schema.TypeString,
		Computed: true,
	},
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
	"comment": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func ApiIntegrationAzureApiManagementDetailsToSchema(d *sdk.ApiIntegrationAzureDetails) map[string]any {
	result := make(map[string]any)
	result["enabled"] = d.Enabled
	result["api_key"] = d.ApiKey
	result["api_provider"] = strings.ToLower(d.ApiProvider)
	result["azure_tenant_id"] = d.AzureTenantId
	result["azure_ad_application_id"] = d.AzureAdApplicationId
	result["azure_multi_tenant_app_name"] = d.AzureMultiTenantAppName
	result["azure_consent_url"] = d.AzureConsentUrl
	result["allowed_prefixes"] = d.AllowedPrefixes
	result["blocked_prefixes"] = d.BlockedPrefixes
	result["comment"] = d.Comment
	return result
}
