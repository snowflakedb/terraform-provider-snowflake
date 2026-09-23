package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (icebergTableDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"type": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (icebergTableDetailsToSchemaMapper) additionalToSchema(src *sdk.IcebergTableDetails, dst map[string]any) {
	// type is the public string form via TypeString(); data_type_raw is omitted.
	dst["type"] = src.TypeString()
}

func IcebergTableDetailsListToSchema(details []sdk.IcebergTableDetails) []map[string]any {
	result := make([]map[string]any, len(details))
	for i := range details {
		result[i] = IcebergTableDetailsToSchema(&details[i])
	}
	return result
}
