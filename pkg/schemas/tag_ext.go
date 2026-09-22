package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (tagToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// allowed_values is []string in the SDK; public show_output is a TypeSet.
		"allowed_values": {
			Type:     schema.TypeSet,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Computed: true,
		},
	}
}

func (tagToSchemaMapper) additionalToSchema(src *sdk.Tag, dst map[string]any) {
	dst["allowed_values"] = src.AllowedValues
}
