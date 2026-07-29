package schemas

import (
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeSaml2IntegrationSchema = map[string]*schema.Schema{
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

func DescribeSaml2IntegrationToSchema(props []sdk.SecurityIntegrationProperty) map[string]any {
	propsSchema := make(map[string]any)
	for _, property := range props {
		property := property
		if slices.Contains(Saml2PropertiesNames, property.Name) {
			propsSchema[strings.ToLower(property.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&property)}
		} else {
			log.Printf("[WARN] unexpected property %v in saml2 security integration returned from Snowflake", property.Name)
		}
	}
	return propsSchema
}
