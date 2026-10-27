package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// warehouseAdaptiveAttributesModificationExt marks warehouse_type computed.
// The attribute generator emits Optional and has no Computed flag.
func warehouseAdaptiveAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["warehouse_type"].Optional = false
	s["warehouse_type"].Computed = true
	return s
}

var warehouseAdaptiveShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSES` for the given adaptive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseAdaptiveSchema,
		},
	},
}
