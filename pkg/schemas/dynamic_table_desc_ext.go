package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (dynamicTableDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (dynamicTableDetailsToSchemaMapper) additionalToSchema(src *sdk.DynamicTableDetails, dst map[string]any) {
	dst["type"] = src.TypeString()
}
