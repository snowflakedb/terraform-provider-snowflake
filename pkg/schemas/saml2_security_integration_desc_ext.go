package schemas

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var Saml2PropertiesNames = []string{
	"COMMENT",
	"SAML2_ISSUER",
	"SAML2_SSO_URL",
	"SAML2_PROVIDER",
	"SAML2_SP_INITIATED_LOGIN_PAGE_LABEL",
	"SAML2_REQUESTED_NAMEID_FORMAT",
	"SAML2_POST_LOGOUT_REDIRECT_URL",
	"SAML2_SNOWFLAKE_ISSUER_URL",
	"SAML2_SNOWFLAKE_ACS_URL",
	"SAML2_SNOWFLAKE_METADATA",
	"SAML2_DIGEST_METHODS_USED",
	"SAML2_SIGNATURE_METHODS_USED",
	"SAML2_ENABLE_SP_INITIATED",
	"SAML2_SIGN_REQUEST",
	"SAML2_FORCE_AUTHN",
	"ALLOWED_USER_DOMAINS",
	"ALLOWED_EMAIL_PATTERNS",
}

// TODO [v3]: flatten nested DescribePropertyListSchema to typed scalars + DescribeSaml2Details
// (describe_output.0.<key>.0.value → describe_output.0.<key>).
func (saml2SecurityIntegrationDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"saml2_issuer":                        DescribePropertyListSchema,
		"saml2_sso_url":                       DescribePropertyListSchema,
		"saml2_provider":                      DescribePropertyListSchema,
		"saml2_sp_initiated_login_page_label": DescribePropertyListSchema,
		"saml2_enable_sp_initiated":           DescribePropertyListSchema,
		"saml2_sign_request":                  DescribePropertyListSchema,
		"saml2_requested_nameid_format":       DescribePropertyListSchema,
		"saml2_post_logout_redirect_url":      DescribePropertyListSchema,
		"saml2_force_authn":                   DescribePropertyListSchema,
		"saml2_snowflake_issuer_url":          DescribePropertyListSchema,
		"saml2_snowflake_acs_url":             DescribePropertyListSchema,
		"saml2_snowflake_metadata":            DescribePropertyListSchema,
		"saml2_digest_methods_used":           DescribePropertyListSchema,
		"saml2_signature_methods_used":        DescribePropertyListSchema,
		"allowed_user_domains":                DescribePropertyListSchema,
		"allowed_email_patterns":              DescribePropertyListSchema,
		"comment":                             DescribePropertyListSchema,
	}
}

func (saml2SecurityIntegrationDetailsToSchemaMapper) additionalToSchema(src *sdk.Saml2SecurityIntegrationDetails, dst map[string]any) {
	if src == nil {
		return
	}
	mapSecurityIntegrationProperty(dst, "saml2_issuer", src.Saml2Issuer)
	mapSecurityIntegrationProperty(dst, "saml2_sso_url", src.Saml2SsoUrl)
	mapSecurityIntegrationProperty(dst, "saml2_provider", src.Saml2Provider)
	mapSecurityIntegrationProperty(dst, "saml2_sp_initiated_login_page_label", src.Saml2SpInitiatedLoginPageLabel)
	mapSecurityIntegrationProperty(dst, "saml2_enable_sp_initiated", src.Saml2EnableSpInitiated)
	mapSecurityIntegrationProperty(dst, "saml2_sign_request", src.Saml2SignRequest)
	mapSecurityIntegrationProperty(dst, "saml2_requested_nameid_format", src.Saml2RequestedNameidFormat)
	mapSecurityIntegrationProperty(dst, "saml2_post_logout_redirect_url", src.Saml2PostLogoutRedirectUrl)
	mapSecurityIntegrationProperty(dst, "saml2_force_authn", src.Saml2ForceAuthn)
	mapSecurityIntegrationProperty(dst, "saml2_snowflake_issuer_url", src.Saml2SnowflakeIssuerUrl)
	mapSecurityIntegrationProperty(dst, "saml2_snowflake_acs_url", src.Saml2SnowflakeAcsUrl)
	mapSecurityIntegrationProperty(dst, "saml2_snowflake_metadata", src.Saml2SnowflakeMetadata)
	mapSecurityIntegrationProperty(dst, "saml2_digest_methods_used", src.Saml2DigestMethodsUsed)
	mapSecurityIntegrationProperty(dst, "saml2_signature_methods_used", src.Saml2SignatureMethodsUsed)
	mapSecurityIntegrationProperty(dst, "allowed_user_domains", src.AllowedUserDomains)
	mapSecurityIntegrationProperty(dst, "allowed_email_patterns", src.AllowedEmailPatterns)
	mapSecurityIntegrationProperty(dst, "comment", src.Comment)
}

func DescribeSaml2IntegrationToSchema(props []sdk.SecurityIntegrationProperty) map[string]any {
	for _, property := range props {
		if !slices.Contains(Saml2PropertiesNames, property.Name) {
			log.Printf("[WARN] unexpected property %v in saml2 security integration returned from Snowflake", property.Name)
		}
	}
	return Saml2SecurityIntegrationDetailsToSchema(sdk.AsSaml2(props))
}
