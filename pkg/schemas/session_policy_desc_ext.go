package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (sessionPolicyDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allowed_secondary_roles": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
		"blocked_secondary_roles": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (sessionPolicyDetailsToSchemaMapper) additionalToSchema(src *sdk.SessionPolicyDetails, dst map[string]any) {
	dst["allowed_secondary_roles"] = src.AllowedSecondaryRoles
	dst["blocked_secondary_roles"] = src.BlockedSecondaryRoles
}
