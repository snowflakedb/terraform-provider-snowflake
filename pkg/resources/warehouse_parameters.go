package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var warehouseParametersSchema = func() map[string]*schema.Schema {
	m := make(map[string]*schema.Schema)
	for _, p := range defs.ParameterDefsForLevel(ssot.ParameterLevelWarehouse) {
		m[p.FieldName()] = parameterSchema(p)
	}
	return m
}()

var warehouseParametersOutputSchema = map[string]*schema.Schema{
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN WAREHOUSE` for the given warehouse.",
		Elem: &schema.Resource{
			Schema: schemas.ShowWarehouseParametersSchema,
		},
	},
}

var warehouseParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	warehouseParametersProvider,
	warehouseParameterDiffFunctionsExt,
)
