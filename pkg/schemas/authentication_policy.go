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
	"name":                       {Type: schema.TypeString, Computed: true},
	"owner":                      {Type: schema.TypeString, Computed: true},
	"authentication_methods":     {Type: schema.TypeString, Computed: true},
	"mfa_authentication_methods": {Type: schema.TypeString, Computed: true},
	"mfa_enrollment":             {Type: schema.TypeString, Computed: true},
	"client_types":               {Type: schema.TypeString, Computed: true},
	"security_integrations":      {Type: schema.TypeString, Computed: true},
	"comment":                    {Type: schema.TypeString, Computed: true},
	"mfa_policy":                 {Type: schema.TypeString, Computed: true},
	"pat_policy":                 {Type: schema.TypeString, Computed: true},
	"workload_identity_policy":   {Type: schema.TypeString, Computed: true},
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
	"MFA_POLICY",
	"PAT_POLICY",
	"WORKLOAD_IDENTITY_POLICY",
}

func AuthenticationPolicyDescriptionToSchema(authenticationPolicyDescription []sdk.AuthenticationPolicyDescription) map[string]any {
	authenticationPolicySchema := make(map[string]any)
	for _, property := range authenticationPolicyDescription {
		if slices.Contains(AuthenticationPolicyNames, property.Property) {
			authenticationPolicySchema[strings.ToLower(property.Property)] = property.Value
		} else {
			log.Printf("[WARN] unexpected property %v in authentication policy returned from Snowflake", property.Property)
		}
	}
	return authenticationPolicySchema
}
