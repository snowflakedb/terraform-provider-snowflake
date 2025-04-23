package schemas

import (
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeOauthIntegrationForPartnerApplications = map[string]*schema.Schema{
	"oauth_client_type":                     DescribePropertyListSchema,
	"enabled":                               DescribePropertyListSchema,
	"oauth_allow_non_tls_redirect_uri":      DescribePropertyListSchema,
	"oauth_enforce_pkce":                    DescribePropertyListSchema,
	"oauth_use_secondary_roles":             DescribePropertyListSchema,
	"pre_authorized_roles_list":             DescribePropertyListSchema,
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

var OauthIntegrationForPartnerApplicationsPropertiesNames = []string{
	"OAUTH_CLIENT_TYPE",
	"ENABLED",
	"OAUTH_ALLOW_NON_TLS_REDIRECT_URI",
	"OAUTH_ENFORCE_PKCE",
	"OAUTH_USE_SECONDARY_ROLES",
	"PRE_AUTHORIZED_ROLES_LIST",
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

func DescribeOauthIntegrationForPartnerApplicationsToSchema(integrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationProperties := make(map[string]any)
	for _, property := range integrationProperties {
		property := property
		if slices.Contains(OauthIntegrationForPartnerApplicationsPropertiesNames, property.Name) {
			securityIntegrationProperties[strings.ToLower(property.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&property)}
		} else {
			log.Printf("[WARN] unexpected property %v in oauth security integration for partner applications returned from Snowflake", property.Name)
		}
	}
	return securityIntegrationProperties
}

var ShowOauthForPartnerApplicationsParametersSchema = map[string]*schema.Schema{
	strings.ToLower(string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList)): ParameterListSchema,
}

func OauthForPartnerApplicationsParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	schemaMap := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains([]sdk.AccountParameter{sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList}, sdk.AccountParameter(param.Key)) {
			schemaMap[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return schemaMap
}
