package schemas

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func warehouseTablesToSchema(tables []sdk.SchemaObjectIdentifier) []string {
	return collections.Map(tables, sdk.SchemaObjectIdentifier.FullyQualifiedName)
}

func warehouseTablesSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func (warehouseToSchemaMapper) additionalSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"tables": warehouseTablesSchema(),
	}
}

func (warehouseToSchemaMapper) additionalToSchema(src *sdk.Warehouse, dst map[string]any) {
	dst["tables"] = warehouseTablesToSchema(src.Tables)
}
