package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (hybridTableDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (hybridTableDetailsToSchemaMapper) additionalToSchema(src *sdk.HybridTableDetails, dst map[string]any) {
	dst["type"] = src.DataTypeRaw
}

func HybridTableDetailsListToSchema(details []sdk.HybridTableDetails) []map[string]any {
	result := make([]map[string]any, len(details))
	for i := range details {
		result[i] = HybridTableDetailsToSchema(&details[i])
	}
	return result
}
