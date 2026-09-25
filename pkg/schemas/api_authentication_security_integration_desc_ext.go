package schemas

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ApiAuthenticationPropertiesNames = []string{
	"ENABLED",
	"OAUTH_ACCESS_TOKEN_VALIDITY",
	"OAUTH_REFRESH_TOKEN_VALIDITY",
	"OAUTH_CLIENT_AUTH_METHOD",
	"OAUTH_AUTHORIZATION_ENDPOINT",
	"OAUTH_TOKEN_ENDPOINT",
	"OAUTH_ALLOWED_SCOPES",
	"OAUTH_GRANT",
	"PARENT_INTEGRATION",
	"AUTH_TYPE",
	"COMMENT",
}

// TODO [v3]: flatten nested DescribePropertyListSchema to typed scalars + DescribeApiAuthDetails
// (describe_output.0.<key>.0.value → describe_output.0.<key>).
func (apiAuthenticationSecurityIntegrationDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled":                      DescribePropertyListSchema,
		"oauth_access_token_validity":  DescribePropertyListSchema,
		"oauth_refresh_token_validity": DescribePropertyListSchema,
		"oauth_client_auth_method":     DescribePropertyListSchema,
		"oauth_authorization_endpoint": DescribePropertyListSchema,
		"oauth_token_endpoint":         DescribePropertyListSchema,
		"oauth_allowed_scopes":         DescribePropertyListSchema,
		"oauth_grant":                  DescribePropertyListSchema,
		"parent_integration":           DescribePropertyListSchema,
		"auth_type":                    DescribePropertyListSchema,
		"comment":                      DescribePropertyListSchema,
	}
}

func (apiAuthenticationSecurityIntegrationDetailsToSchemaMapper) additionalToSchema(src *sdk.ApiAuthenticationSecurityIntegrationDetails, dst map[string]any) {
	if src == nil {
		return
	}
	mapSecurityIntegrationProperty(dst, "enabled", src.Enabled)
	mapSecurityIntegrationProperty(dst, "oauth_access_token_validity", src.OauthAccessTokenValidity)
	mapSecurityIntegrationProperty(dst, "oauth_refresh_token_validity", src.OauthRefreshTokenValidity)
	mapSecurityIntegrationProperty(dst, "oauth_client_auth_method", src.OauthClientAuthMethod)
	mapSecurityIntegrationProperty(dst, "oauth_authorization_endpoint", src.OauthAuthorizationEndpoint)
	mapSecurityIntegrationProperty(dst, "oauth_token_endpoint", src.OauthTokenEndpoint)
	mapSecurityIntegrationProperty(dst, "oauth_allowed_scopes", src.OauthAllowedScopes)
	mapSecurityIntegrationProperty(dst, "oauth_grant", src.OauthGrant)
	mapSecurityIntegrationProperty(dst, "parent_integration", src.ParentIntegration)
	mapSecurityIntegrationProperty(dst, "auth_type", src.AuthType)
	mapSecurityIntegrationProperty(dst, "comment", src.Comment)
}

func ApiAuthSecurityIntegrationPropertiesToSchema(securityIntegrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	for _, property := range securityIntegrationProperties {
		if !slices.Contains(ApiAuthenticationPropertiesNames, property.Name) {
			log.Printf("[WARN] unexpected property %v in api auth security integration returned from Snowflake", property.Name)
		}
	}
	return ApiAuthenticationSecurityIntegrationDetailsToSchema(sdk.AsApiAuth(securityIntegrationProperties))
}
