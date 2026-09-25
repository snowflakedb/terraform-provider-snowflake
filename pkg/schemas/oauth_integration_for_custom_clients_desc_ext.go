package schemas

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var OauthIntegrationForCustomClientsPropertiesNames = []string{
	"OAUTH_CLIENT_TYPE",
	"ENABLED",
	"OAUTH_ALLOW_NON_TLS_REDIRECT_URI",
	"OAUTH_ENFORCE_PKCE",
	"OAUTH_USE_SECONDARY_ROLES",
	"PRE_AUTHORIZED_ROLES_LIST",
	"ALLOWED_ROLES_LIST",
	"BLOCKED_ROLES_LIST",
	"OAUTH_ISSUE_REFRESH_TOKENS",
	"OAUTH_REFRESH_TOKEN_VALIDITY",
	"NETWORK_POLICY",
	"OAUTH_CLIENT_RSA_PUBLIC_KEY_FP",
	"OAUTH_CLIENT_RSA_PUBLIC_KEY_2_FP",
	"COMMENT",
	"OAUTH_AUTHORIZATION_ENDPOINT",
	"OAUTH_TOKEN_ENDPOINT",
	"OAUTH_ALLOWED_AUTHORIZATION_ENDPOINTS",
	"OAUTH_ALLOWED_TOKEN_ENDPOINTS",
}

// TODO [v3]: flatten nested DescribePropertyListSchema to typed scalars + DescribeOauthCustomDetails
// (describe_output.0.<key>.0.value → describe_output.0.<key>).
func (oauthIntegrationForCustomClientsDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"oauth_client_type":                     DescribePropertyListSchema,
		"enabled":                               DescribePropertyListSchema,
		"oauth_allow_non_tls_redirect_uri":      DescribePropertyListSchema,
		"oauth_enforce_pkce":                    DescribePropertyListSchema,
		"oauth_use_secondary_roles":             DescribePropertyListSchema,
		"pre_authorized_roles_list":             DescribePropertyListSchema,
		"allowed_roles_list":                    DescribePropertyListSchema,
		"blocked_roles_list":                    DescribePropertyListSchema,
		"oauth_issue_refresh_tokens":            DescribePropertyListSchema,
		"oauth_refresh_token_validity":          DescribePropertyListSchema,
		"network_policy":                        DescribePropertyListSchema,
		"oauth_client_rsa_public_key_fp":        DescribePropertyListSchema,
		"oauth_client_rsa_public_key_2_fp":      DescribePropertyListSchema,
		"comment":                               DescribePropertyListSchema,
		"oauth_authorization_endpoint":          DescribePropertyListSchema,
		"oauth_token_endpoint":                  DescribePropertyListSchema,
		"oauth_allowed_authorization_endpoints": DescribePropertyListSchema,
		"oauth_allowed_token_endpoints":         DescribePropertyListSchema,
	}
}

func (oauthIntegrationForCustomClientsDetailsToSchemaMapper) additionalToSchema(src *sdk.OauthIntegrationForCustomClientsDetails, dst map[string]any) {
	if src == nil {
		return
	}
	mapSecurityIntegrationProperty(dst, "oauth_client_type", src.OauthClientType)
	mapSecurityIntegrationProperty(dst, "enabled", src.Enabled)
	mapSecurityIntegrationProperty(dst, "oauth_allow_non_tls_redirect_uri", src.OauthAllowNonTlsRedirectUri)
	mapSecurityIntegrationProperty(dst, "oauth_enforce_pkce", src.OauthEnforcePkce)
	mapSecurityIntegrationProperty(dst, "oauth_use_secondary_roles", src.OauthUseSecondaryRoles)
	mapSecurityIntegrationProperty(dst, "pre_authorized_roles_list", src.PreAuthorizedRolesList)
	mapSecurityIntegrationProperty(dst, "allowed_roles_list", src.AllowedRolesList)
	mapSecurityIntegrationProperty(dst, "blocked_roles_list", src.BlockedRolesList)
	mapSecurityIntegrationProperty(dst, "oauth_issue_refresh_tokens", src.OauthIssueRefreshTokens)
	mapSecurityIntegrationProperty(dst, "oauth_refresh_token_validity", src.OauthRefreshTokenValidity)
	mapSecurityIntegrationProperty(dst, "network_policy", src.NetworkPolicy)
	mapSecurityIntegrationProperty(dst, "oauth_client_rsa_public_key_fp", src.OauthClientRsaPublicKeyFp)
	mapSecurityIntegrationProperty(dst, "oauth_client_rsa_public_key_2_fp", src.OauthClientRsaPublicKey2Fp)
	mapSecurityIntegrationProperty(dst, "comment", src.Comment)
	mapSecurityIntegrationProperty(dst, "oauth_authorization_endpoint", src.OauthAuthorizationEndpoint)
	mapSecurityIntegrationProperty(dst, "oauth_token_endpoint", src.OauthTokenEndpoint)
	mapSecurityIntegrationProperty(dst, "oauth_allowed_authorization_endpoints", src.OauthAllowedAuthorizationEndpoints)
	mapSecurityIntegrationProperty(dst, "oauth_allowed_token_endpoints", src.OauthAllowedTokenEndpoints)
}

func DescribeOauthIntegrationForCustomClientsToSchema(integrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	for _, property := range integrationProperties {
		if !slices.Contains(OauthIntegrationForCustomClientsPropertiesNames, property.Name) {
			log.Printf("[WARN] unexpected property %v in oauth security integration for custom clients returned from Snowflake", property.Name)
		}
	}
	return OauthIntegrationForCustomClientsDetailsToSchema(sdk.AsOauthCustom(integrationProperties))
}
