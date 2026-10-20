package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var warehouseIdentitySchemaExt = map[string]*schema.Schema{
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      blocklistedCharactersFieldDescription("Identifier for the virtual warehouse; must be unique for your account."),
		DiffSuppressFunc: suppressIdentifierQuoting,
	},
}

var warehouseShowOutputSchemaExt = map[string]*schema.Schema{
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW WAREHOUSES` for the given warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseRegularSchema,
		},
	},
}
