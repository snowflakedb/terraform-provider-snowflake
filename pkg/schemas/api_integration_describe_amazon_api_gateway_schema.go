package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeAmazonApiGatewayApiIntegrationSchema = map[string]*schema.Schema{
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
	"api_aws_role_arn": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"api_aws_iam_user_arn": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"api_aws_external_id": {
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

func ApiIntegrationAmazonApiGatewayDetailsToSchema(d *sdk.ApiIntegrationAwsDetails) map[string]any {
	result := make(map[string]any)
	result["enabled"] = d.Enabled
	result["api_key"] = d.ApiKey
	result["api_provider"] = strings.ToLower(d.ApiProvider)
	result["api_aws_role_arn"] = d.ApiAwsRoleArn
	result["api_aws_iam_user_arn"] = d.ApiAwsIamUserArn
	result["api_aws_external_id"] = d.ApiAwsExternalId
	result["allowed_prefixes"] = d.AllowedPrefixes
	result["blocked_prefixes"] = d.BlockedPrefixes
	result["comment"] = d.Comment
	return result
}
