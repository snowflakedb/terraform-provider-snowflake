package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var StageDescribeSchema = map[string]*schema.Schema{
	"directory_table": {
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"enable": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"auto_refresh": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
		Computed: true,
	},
}

func StageDescribeToSchema(properties sdk.StageDetails) (map[string]any, error) {
	schema := make(map[string]any)

	if properties.DirectoryTable != nil {
		schema["directory_table"] = []map[string]any{
			{
				"enable":       properties.DirectoryTable.Enable,
				"auto_refresh": properties.DirectoryTable.AutoRefresh,
			},
		}
	}
	return schema, nil
}
