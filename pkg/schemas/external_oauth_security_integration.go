package schemas

import (
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeExternalOauthSecurityIntegrationSchema represents output of DESCRIBE query for the single SecurityIntegration.
var DescribeExternalOauthSecurityIntegrationSchema = map[string]*schema.Schema{
	"enabled":                                         DescribePropertyListSchema,
	"external_oauth_issuer":                           DescribePropertyListSchema,
	"external_oauth_jws_keys_url":                     DescribePropertyListSchema,
	"external_oauth_any_role_mode":                    DescribePropertyListSchema,
	"external_oauth_rsa_public_key":                   DescribePropertyListSchema,
	"external_oauth_rsa_public_key_2":                 DescribePropertyListSchema,
	"external_oauth_blocked_roles_list":               DescribePropertyListSchema,
	"external_oauth_allowed_roles_list":               DescribePropertyListSchema,
	"external_oauth_audience_list":                    DescribePropertyListSchema,
	"external_oauth_token_user_mapping_claim":         DescribePropertyListSchema,
	"external_oauth_snowflake_user_mapping_attribute": DescribePropertyListSchema,
	"external_oauth_scope_delimiter":                  DescribePropertyListSchema,
	"comment":                                         DescribePropertyListSchema,
}

var ExternalOauthPropertiesNames = []string{
	"ENABLED",
	"EXTERNAL_OAUTH_ISSUER",
	"EXTERNAL_OAUTH_JWS_KEYS_URL",
	"EXTERNAL_OAUTH_ANY_ROLE_MODE",
	"EXTERNAL_OAUTH_RSA_PUBLIC_KEY",
	"EXTERNAL_OAUTH_RSA_PUBLIC_KEY_2",
	"EXTERNAL_OAUTH_BLOCKED_ROLES_LIST",
	"EXTERNAL_OAUTH_ALLOWED_ROLES_LIST",
	"EXTERNAL_OAUTH_AUDIENCE_LIST",
	"EXTERNAL_OAUTH_TOKEN_USER_MAPPING_CLAIM",
	"EXTERNAL_OAUTH_SNOWFLAKE_USER_MAPPING_ATTRIBUTE",
	"EXTERNAL_OAUTH_SCOPE_DELIMITER",
	"COMMENT",
}

func ExternalOauthSecurityIntegrationPropertiesToSchema(securityIntegrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationSchema := make(map[string]any)
	for _, securityIntegrationProperty := range securityIntegrationProperties {
		securityIntegrationProperty := securityIntegrationProperty
		if slices.Contains(ExternalOauthPropertiesNames, securityIntegrationProperty.Name) {
			securityIntegrationSchema[strings.ToLower(securityIntegrationProperty.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&securityIntegrationProperty)}
		} else {
			log.Printf("[WARN] unexpected property %v in external oauth security integration returned from Snowflake", securityIntegrationProperty.Name)
		}
	}
	return securityIntegrationSchema
}

var _ = ExternalOauthSecurityIntegrationPropertiesToSchema

var ShowExternalOauthParametersSchema = map[string]*schema.Schema{
	strings.ToLower(string(sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList)): ParameterListSchema,
}

func ExternalOauthParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	schemaMap := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains([]sdk.AccountParameter{sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList}, sdk.AccountParameter(param.Key)) {
			schemaMap[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return schemaMap
}
