package schemas

import (
	"strconv"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (viewDetailsToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"check": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func (viewDetailsToSchemaMapper) additionalToSchema(src *sdk.ViewDetails, dst map[string]any) {
	if src.Check != nil {
		dst["check"] = strconv.FormatBool(*src.Check)
	}
}

func ViewDetailsListToSchema(details []sdk.ViewDetails) []map[string]any {
	result := make([]map[string]any, len(details))
	for i := range details {
		result[i] = ViewDetailsToSchema(&details[i])
	}
	return result
}
