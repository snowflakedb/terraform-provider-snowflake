package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (resourceMonitorToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Generated snake_case is suspend_immediately_at; keep the existing public key.
		"suspend_immediate_at": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func (resourceMonitorToSchemaMapper) additionalToSchema(src *sdk.ResourceMonitor, dst map[string]any) {
	if src.SuspendImmediatelyAt != nil {
		dst["suspend_immediate_at"] = (*src.SuspendImmediatelyAt)
	}
}
