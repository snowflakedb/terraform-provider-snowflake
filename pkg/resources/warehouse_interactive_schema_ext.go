package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// warehouseInteractiveAttributesModificationExt applies schema details the attribute generator does not emit.
// warehouse_type is computed. auto_suspend keeps the IntDefault sentinel even though 0 is invalid (min 1).
// tables is a set of schema-object identifiers, applied incrementally in BeforeAlter.
func warehouseInteractiveAttributesModificationExt(s map[string]*schema.Schema) map[string]*schema.Schema {
	s["warehouse_type"].Optional = false
	s["warehouse_type"].Computed = true
	s["auto_suspend"].Default = IntDefault
	s["tables"].Elem.(*schema.Schema).ValidateDiagFunc = IsValidIdentifier[sdk.SchemaObjectIdentifier]()
	s["tables"].DiffSuppressFunc = NormalizeAndCompareIdentifiersInSet("tables")
	return s
}

var warehouseInteractiveShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSES` for the given interactive warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseInteractiveSchema,
		},
	},
}
