package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (materializedViewDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (materializedViewDetailsToSchemaMapper) additionalToSchema(src *sdk.MaterializedViewDetails, dst map[string]any) {
	dst["type"] = src.TypeString()
}
