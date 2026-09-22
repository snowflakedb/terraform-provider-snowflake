package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (cortexAgentDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"profile": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"display_name": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"avatar": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"color": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (cortexAgentDetailsToSchemaMapper) additionalToSchema(src *sdk.CortexAgentDetails, dst map[string]any) {
	dst["profile"] = []map[string]any{
		{
			"display_name": src.Profile.DisplayName,
			"avatar":       src.Profile.Avatar,
			"color":        src.Profile.Color,
		},
	}
}
