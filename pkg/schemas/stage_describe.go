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

func StageDescribeToSchema(properties []sdk.StageProperty) (map[string]any, error) {
	details, err := sdk.ParseStageDetails(properties)
	if err != nil {
		return nil, err
	}

	schema := make(map[string]any)

	if details.DirectoryTable != nil {
		schema["directory_table"] = []map[string]any{
			{
				"enable":       details.DirectoryTable.Enable,
				"auto_refresh": details.DirectoryTable.AutoRefresh,
			},
		}
	}
	return schema, nil
}
