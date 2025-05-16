package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var SchemaDescribeSchema = map[string]*schema.Schema{
	"created_on": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"kind": {
		Type:     schema.TypeString,
		Computed: true,
	},
}

func SchemaDescriptionToSchema(description []sdk.SchemaDetails) []map[string]any {
	result := make([]map[string]any, len(description))
	for i, row := range description {
		result[i] = map[string]any{
			"created_on": row.CreatedOn.String(),
			"name":       row.Name,
			"kind":       row.Kind,
		}
	}
	return result
}
