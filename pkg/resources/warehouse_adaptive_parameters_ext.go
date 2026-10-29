package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Schema maps
// =============================================================================

var warehouseAdaptiveParametersOutputSchemaExt = map[string]*schema.Schema{
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN WAREHOUSE` for the given adaptive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseAdaptiveParametersSchema,
		},
	},
}
