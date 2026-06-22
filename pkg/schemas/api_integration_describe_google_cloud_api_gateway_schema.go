package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeGoogleCloudApiGatewayApiIntegrationSchema = map[string]*schema.Schema{
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
	"google_audience": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"google_api_service_account": {
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

func ApiIntegrationGoogleCloudApiGatewayDetailsToSchema(d *sdk.ApiIntegrationGoogleDetails) map[string]any {
	result := make(map[string]any)
	result["enabled"] = d.Enabled
	result["api_key"] = d.ApiKey
	result["api_provider"] = strings.ToLower(d.ApiProvider)
	result["google_audience"] = d.GoogleAudience
	result["google_api_service_account"] = d.GoogleApiServiceAccount
	result["allowed_prefixes"] = d.AllowedPrefixes
	result["blocked_prefixes"] = d.BlockedPrefixes
	result["comment"] = d.Comment
	return result
}
