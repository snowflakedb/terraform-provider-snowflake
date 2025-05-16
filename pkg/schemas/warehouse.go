package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var WarehouseDescribeSchema = map[string]*schema.Schema{
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

func WarehouseDescriptionToSchema(description sdk.WarehouseDetails) []map[string]any {
	return []map[string]any{
		{
			"created_on": description.CreatedOn.String(),
			"name":       description.Name,
			"kind":       description.Kind,
		},
	}
}
