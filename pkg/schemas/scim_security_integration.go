package schemas

import (
	"log"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DescribeScimSecurityIntegrationSchema represents output of DESCRIBE query for the single SecurityIntegration.
var DescribeScimSecurityIntegrationSchema = map[string]*schema.Schema{
	"enabled":        DescribePropertyListSchema,
	"network_policy": DescribePropertyListSchema,
	"run_as_role":    DescribePropertyListSchema,
	"sync_password":  DescribePropertyListSchema,
	"comment":        DescribePropertyListSchema,
}

var (
	_                   = DescribeScimSecurityIntegrationSchema
	ScimPropertiesNames = []string{
		"ENABLED",
		"NETWORK_POLICY",
		"RUN_AS_ROLE",
		"SYNC_PASSWORD",
		"COMMENT",
	}
)

func ScimSecurityIntegrationPropertiesToSchema(securityIntegrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	securityIntegrationSchema := make(map[string]any)
	for _, property := range securityIntegrationProperties {
		property := property
		if slices.Contains(ScimPropertiesNames, property.Name) {
			securityIntegrationSchema[strings.ToLower(property.Name)] = []map[string]any{SecurityIntegrationPropertyToSchema(&property)}
		} else {
			log.Printf("[WARN] unexpected property %v in scim security integration returned from Snowflake", property.Name)
		}
	}
	return securityIntegrationSchema
}

var _ = ScimSecurityIntegrationPropertiesToSchema
