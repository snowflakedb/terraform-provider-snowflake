package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (cortexSearchServiceDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"attribute_columns": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"columns": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"primary_key_columns": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func (cortexSearchServiceDetailsToSchemaMapper) additionalToSchema(src *sdk.CortexSearchServiceDetails, dst map[string]any) {
	dst["attribute_columns"] = src.AttributeColumns
	dst["columns"] = src.Columns
	dst["primary_key_columns"] = src.PrimaryKeyColumns
}
