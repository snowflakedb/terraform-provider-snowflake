package schemas

import (
	"log"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var ScimPropertiesNames = []string{
	"ENABLED",
	"NETWORK_POLICY",
	"RUN_AS_ROLE",
	"SYNC_PASSWORD",
	"COMMENT",
}

// TODO [v3]: flatten nested DescribePropertyListSchema to typed scalars + DescribeScimDetails
// (describe_output.0.enabled.0.value → describe_output.0.enabled).
func (scimSecurityIntegrationDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled":        DescribePropertyListSchema,
		"network_policy": DescribePropertyListSchema,
		"run_as_role":    DescribePropertyListSchema,
		"sync_password":  DescribePropertyListSchema,
		"comment":        DescribePropertyListSchema,
	}
}

func (scimSecurityIntegrationDetailsToSchemaMapper) additionalToSchema(src *sdk.ScimSecurityIntegrationDetails, dst map[string]any) {
	if src == nil {
		return
	}
	mapSecurityIntegrationProperty(dst, "enabled", src.Enabled)
	mapSecurityIntegrationProperty(dst, "network_policy", src.NetworkPolicy)
	mapSecurityIntegrationProperty(dst, "run_as_role", src.RunAsRole)
	mapSecurityIntegrationProperty(dst, "sync_password", src.SyncPassword)
	mapSecurityIntegrationProperty(dst, "comment", src.Comment)
}

func mapSecurityIntegrationProperty(dst map[string]any, key string, property *sdk.SecurityIntegrationProperty) {
	if property == nil {
		return
	}
	dst[key] = []map[string]any{SecurityIntegrationPropertyToSchema(property)}
}

func ScimSecurityIntegrationPropertiesToSchema(securityIntegrationProperties []sdk.SecurityIntegrationProperty) map[string]any {
	for _, property := range securityIntegrationProperties {
		if !slices.Contains(ScimPropertiesNames, property.Name) {
			log.Printf("[WARN] unexpected property %v in scim security integration returned from Snowflake", property.Name)
		}
	}
	return ScimSecurityIntegrationDetailsToSchema(sdk.AsScim(securityIntegrationProperties))
}
