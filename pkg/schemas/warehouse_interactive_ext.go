package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (warehouseInteractiveToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tables": warehouseTablesSchema(),
	}
}

func (warehouseInteractiveToSchemaMapper) additionalToSchema(src *sdk.WarehouseInteractive, dst map[string]any) {
	dst["tables"] = warehouseTablesToSchema(src.Tables)
}
