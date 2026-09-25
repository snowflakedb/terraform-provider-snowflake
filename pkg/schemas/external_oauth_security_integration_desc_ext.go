package schemas

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

// TODO [v3]: flatten nested DescribePropertyListSchema to typed scalars + DescribeExternalOauthDetails
// (describe_output.0.<key>.0.value → describe_output.0.<key>).
func (externalOauthSecurityIntegrationDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
}

func (externalOauthSecurityIntegrationDetailsToSchemaMapper) additionalToSchema(src *sdk.ExternalOauthSecurityIntegrationDetails, dst map[string]any) {
	if src == nil {
		return
	}
	mapSecurityIntegrationProperty(dst, "enabled", src.Enabled)
	mapSecurityIntegrationProperty(dst, "external_oauth_issuer", src.ExternalOauthIssuer)
	mapSecurityIntegrationProperty(dst, "external_oauth_jws_keys_url", src.ExternalOauthJwsKeysUrl)
	mapSecurityIntegrationProperty(dst, "external_oauth_any_role_mode", src.ExternalOauthAnyRoleMode)
	mapSecurityIntegrationProperty(dst, "external_oauth_rsa_public_key", src.ExternalOauthRsaPublicKey)
	mapSecurityIntegrationProperty(dst, "external_oauth_rsa_public_key_2", src.ExternalOauthRsaPublicKey2)
	mapSecurityIntegrationProperty(dst, "external_oauth_blocked_roles_list", src.ExternalOauthBlockedRolesList)
	mapSecurityIntegrationProperty(dst, "external_oauth_allowed_roles_list", src.ExternalOauthAllowedRolesList)
	mapSecurityIntegrationProperty(dst, "external_oauth_audience_list", src.ExternalOauthAudienceList)
	mapSecurityIntegrationProperty(dst, "external_oauth_token_user_mapping_claim", src.ExternalOauthTokenUserMappingClaim)
	mapSecurityIntegrationProperty(dst, "external_oauth_snowflake_user_mapping_attribute", src.ExternalOauthSnowflakeUserMappingAttribute)
	mapSecurityIntegrationProperty(dst, "external_oauth_scope_delimiter", src.ExternalOauthScopeDelimiter)
	mapSecurityIntegrationProperty(dst, "comment", src.Comment)
}

func ExternalOauthSecurityIntegrationPropertiesToSchema(securityIntegrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	for _, property := range securityIntegrationProperties {
		if !slices.Contains(ExternalOauthPropertiesNames, property.Name) {
			log.Printf("[WARN] unexpected property %v in external oauth security integration returned from Snowflake", property.Name)
		}
	}
	return ExternalOauthSecurityIntegrationDetailsToSchema(sdk.AsExternalOauth(securityIntegrationProperties))
}
