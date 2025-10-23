package schemas

import (
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AuthenticationPolicyDescribeSchema represents output of DESCRIBE query for the single AuthenticationPolicy.
var AuthenticationPolicyDescribeSchema = map[string]*schema.Schema{
	"name":                       DescribeAuthenticationPolicyPropertyListSchema,
	"owner":                      DescribeAuthenticationPolicyPropertyListSchema,
	"authentication_methods":     DescribeAuthenticationPolicyPropertyListSchema,
	"mfa_authentication_methods": DescribeAuthenticationPolicyPropertyListSchema,
	"mfa_enrollment":             DescribeAuthenticationPolicyPropertyListSchema,
	"client_types":               DescribeAuthenticationPolicyPropertyListSchema,
	"security_integrations":      DescribeAuthenticationPolicyPropertyListSchema,
	"comment":                    DescribeAuthenticationPolicyPropertyListSchema,
}

// DescribeAuthenticationPolicyPropertyListSchema represents Snowflake property object returned by DESCRIBE query.
var DescribeAuthenticationPolicyPropertyListSchema = &schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: DescribeAuthenticationPolicyPropertySchema,
	},
}

var _ = AuthenticationPolicyDescribeSchema

var AuthenticationPolicyNames = []string{
	"NAME",
	"OWNER",
	"COMMENT",
	"AUTHENTICATION_METHODS",
	"CLIENT_TYPES",
	"SECURITY_INTEGRATIONS",
	"MFA_ENROLLMENT",
	"MFA_AUTHENTICATION_METHODS",
}

func AuthenticationPolicyDescriptionToSchema(authenticationPolicyDescription []sdk.AuthenticationPolicyDescription) map[string]any {
	authenticationPolicySchema := make(map[string]any)
	for _, property := range authenticationPolicyDescription {
		if slices.Contains(AuthenticationPolicyNames, property.Property) {
			authenticationPolicySchema[strings.ToLower(property.Property)] = []map[string]any{AuthenticationPolicyDescribePropertyToSchema(&property)}
		} else {
			log.Printf("[WARN] unexpected property %v in authentication policy returned from Snowflake", property.Property)
		}
	}
	return authenticationPolicySchema
}

func AuthenticationPolicyDescribePropertyToSchema(property *sdk.AuthenticationPolicyDescription) map[string]any {
	authenticationPolicyPropertySchema := make(map[string]any)
	authenticationPolicyPropertySchema["property"] = property.Property
	authenticationPolicyPropertySchema["value"] = property.Value
	authenticationPolicyPropertySchema["default"] = property.Default
	authenticationPolicyPropertySchema["description"] = property.Description
	return authenticationPolicyPropertySchema
}

// DescribeAuthenticationPolicyPropertySchema represents output of SHOW query for the single SecurityIntegrationProperty.
var DescribeAuthenticationPolicyPropertySchema = map[string]*schema.Schema{
	"property": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"value": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"default": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"description": {
		Type:     schema.TypeString,
		Computed: true,
	},
}
