package schemas

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeExternalMcpOAuth2ApiIntegrationSchema = map[string]*schema.Schema{
	"enabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"api_provider": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"user_auth_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_grant": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_client_id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_client_auth_method": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_token_endpoint": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_authorization_endpoint": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_access_token_validity": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"oauth_refresh_token_validity": {
		Type:     schema.TypeInt,
		Computed: true,
	},
	"oauth_allowed_scopes": {
		Type:     schema.TypeList,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Computed: true,
	},
	"oauth_username": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_assertion_issuer": {
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

func ApiIntegrationExternalMcpOAuth2DetailsToSchema(d *sdk.ApiIntegrationExternalMcpDetails) map[string]any {
	result := make(map[string]any)
	result["enabled"] = d.Enabled
	result["api_provider"] = strings.ToLower(d.ApiProvider)
	result["user_auth_type"] = d.UserAuthType
	result["oauth_grant"] = d.OauthGrant
	result["oauth_client_id"] = d.OauthClientId
	result["oauth_client_auth_method"] = d.OauthClientAuthMethod
	result["oauth_token_endpoint"] = d.OauthTokenEndpoint
	result["oauth_authorization_endpoint"] = d.OauthAuthorizationEndpoint
	result["oauth_access_token_validity"] = d.OauthAccessTokenValidity
	result["oauth_refresh_token_validity"] = d.OauthRefreshTokenValidity
	result["oauth_allowed_scopes"] = d.OauthAllowedScopes
	result["oauth_username"] = d.OauthUsername
	result["oauth_assertion_issuer"] = d.OauthAssertionIssuer
	result["allowed_prefixes"] = d.AllowedPrefixes
	result["blocked_prefixes"] = d.BlockedPrefixes
	result["comment"] = d.Comment
	return result
}
