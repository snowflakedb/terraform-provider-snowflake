package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeGitRepositoryOauth2ApiIntegrationSchema = map[string]*schema.Schema{
	"enabled": {
		Type:     schema.TypeBool,
		Computed: true,
	},
	"user_auth_type": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"oauth_client_id": {
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

func ApiIntegrationGitRepositoryOauth2DetailsToSchema(d *sdk.ApiIntegrationGitHttpsApiDetails) map[string]any {
	result := make(map[string]any)
	result["enabled"] = d.Enabled
	result["user_auth_type"] = d.UserAuthType
	result["oauth_client_id"] = d.OauthClientId
	result["oauth_token_endpoint"] = d.OauthTokenEndpoint
	result["oauth_authorization_endpoint"] = d.OauthAuthorizationEndpoint
	result["oauth_access_token_validity"] = d.OauthAccessTokenValidity
	result["oauth_refresh_token_validity"] = d.OauthRefreshTokenValidity
	result["oauth_allowed_scopes"] = d.OauthAllowedScopes
	result["oauth_username"] = d.OauthUsername
	result["allowed_prefixes"] = d.AllowedPrefixes
	result["blocked_prefixes"] = d.BlockedPrefixes
	result["comment"] = d.Comment
	return result
}
