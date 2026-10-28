package schemas

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var authenticationPolicyDescribeProperties = map[string]string{
	"NAME":                     "name",
	"OWNER":                    "owner",
	"AUTHENTICATION_METHODS":   "authentication_methods",
	"MFA_ENROLLMENT":           "mfa_enrollment",
	"CLIENT_TYPES":             "client_types",
	"SECURITY_INTEGRATIONS":    "security_integrations",
	"COMMENT":                  "comment",
	"CLIENT_POLICY":            "client_policy",
	"MFA_POLICY":               "mfa_policy",
	"PAT_POLICY":               "pat_policy",
	"WORKLOAD_IDENTITY_POLICY": "workload_identity_policy",
}

// authenticationPolicyManualDescribeKeys are the describe_output keys that are not mapped from the typed details.
// They keep the raw DESCRIBE string, which the typed values cannot rebuild without dropping properties the details do not model.
var authenticationPolicyManualDescribeKeys = []string{
	"authentication_methods",
	"client_types",
	"security_integrations",
	"comment",
	"client_policy",
	"mfa_policy",
	"pat_policy",
	"workload_identity_policy",
}

func (authenticationPolicyDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	manual := make(map[string]*schema.Schema, len(authenticationPolicyManualDescribeKeys))
	for _, key := range authenticationPolicyManualDescribeKeys {
		manual[key] = &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		}
	}
	return manual
}

func (authenticationPolicyDetailsToSchemaMapper) additionalToSchema(details *sdk.AuthenticationPolicyDetails, detailsSchema map[string]any) {
	for property, key := range authenticationPolicyDescribeProperties {
		if slices.Contains(authenticationPolicyManualDescribeKeys, key) {
			detailsSchema[key] = details.Descriptions().Raw(property)
		}
	}
}

// AuthenticationPolicyDescriptionsToSchema maps the DESCRIBE rows to the describe_output values without parsing them.
func AuthenticationPolicyDescriptionsToSchema(properties []sdk.AuthenticationPolicyProperty) map[string]any {
	if properties == nil {
		return map[string]any{}
	}
	result := make(map[string]any, len(authenticationPolicyDescribeProperties))
	for _, key := range authenticationPolicyDescribeProperties {
		result[key] = ""
	}
	for _, property := range properties {
		key, ok := authenticationPolicyDescribeProperties[property.Property]
		if !ok {
			log.Printf("[WARN] unexpected property %v in authentication policy returned from Snowflake", property.Property)
			continue
		}
		result[key] = property.Value
	}
	return result
}
