package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var DescribeExternalVolumeSchema = map[string]*schema.Schema{
	"parent": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"type": {
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
}

var _ = ShowExternalVolumeSchema

func ExternalVolumeDescriptionToSchema(description []sdk.ExternalVolumeProperty) []map[string]any {
	result := make([]map[string]any, len(description))
	for i, row := range description {
		result[i] = map[string]any{
			"parent":  row.Parent,
			"name":    row.Name,
			"type":    row.Type,
			"value":   row.Value,
			"default": row.Default,
		}
	}
	return result
}

var _ = ExternalVolumeDescriptionToSchema
