package schemas

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

var AuthenticationPolicyNames = []string{
	"NAME",
	"OWNER",
	"COMMENT",
	"AUTHENTICATION_METHODS",
	"CLIENT_TYPES",
	"CLIENT_POLICY",
	"SECURITY_INTEGRATIONS",
	"MFA_ENROLLMENT",
	"MFA_POLICY",
	"PAT_POLICY",
	"WORKLOAD_IDENTITY_POLICY",
}

// AuthenticationPolicyDetailsToSchema maps the original DESC rows. Formatting the parsed fields
// back into strings drops properties the parser does not model.
func AuthenticationPolicyDetailsToSchema(details *sdk.AuthenticationPolicyDetails) map[string]any {
	if details == nil {
		return map[string]any{}
	}
	return AuthenticationPolicyDescriptionsToSchema(details.Descriptions())
}

func AuthenticationPolicyDescriptionsToSchema(authenticationPolicyDescription []sdk.AuthenticationPolicyDescription) map[string]any {
	for _, property := range authenticationPolicyDescription {
		if !slices.Contains(AuthenticationPolicyNames, property.Property) {
			log.Printf("[WARN] unexpected property %v in authentication policy returned from Snowflake", property.Property)
		}
	}
	details := sdk.AsAuthenticationPolicyDescribe(authenticationPolicyDescription)
	if details == nil {
		return map[string]any{}
	}
	return AuthenticationPolicyDescribeDetailsToSchema(details)
}
