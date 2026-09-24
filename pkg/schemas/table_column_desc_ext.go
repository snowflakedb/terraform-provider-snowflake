package schemas

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (tableColumnDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"check": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (tableColumnDetailsToSchemaMapper) additionalToSchema(src *sdk.TableColumnDetails, dst map[string]any) {
	if src.Check != nil {
		dst["check"] = strconv.FormatBool(*src.Check)
	}
}

func TableColumnDetailsListToSchema(details []sdk.TableColumnDetails) []map[string]any {
	result := make([]map[string]any, len(details))
	for i := range details {
		result[i] = TableColumnDetailsToSchema(&details[i])
	}
	return result
}
