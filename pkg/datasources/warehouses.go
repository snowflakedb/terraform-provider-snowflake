package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var warehousesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC WAREHOUSE for each warehouse returned by SHOW WAREHOUSES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"with_parameters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs SHOW PARAMETERS FOR WAREHOUSE for each warehouse returned by SHOW WAREHOUSES. The output of describe is saved to the parameters field as a map. By default this value is set to true.",
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"warehouses": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all warehouse details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW WAREHOUSES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowWarehouseSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE WAREHOUSE.",
					Elem: &schema.Resource{
						Schema: schemas.WarehouseDescribeSchema,
					},
				},
				resources.ParametersAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PARAMETERS FOR WAREHOUSE.",
					Elem: &schema.Resource{
						Schema: schemas.ShowWarehouseParametersSchema,
					},
				},
			},
		},
	},
}

func Warehouses() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Warehouses, ReadWarehouses),
		Schema:      warehousesSchema,
		Description: "Data source used to get details of filtered warehouses. Filtering is aligned with the current possibilities for [SHOW WAREHOUSES](https://docs.snowflake.com/en/sql-reference/sql/show-warehouses) query (only `like` is supported). The results of SHOW, DESCRIBE, and SHOW PARAMETERS IN are encapsulated in one output collection.",
	}
}

func ReadWarehouses(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	var opts sdk.ShowWarehouseOptions

	if likePattern, ok := d.GetOk("like"); ok {
		opts.Like = &sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		}
	}

	warehouses, err := client.Warehouses.Show(ctx, &opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("warehouses_read")

	flattenedWarehouses := make([]map[string]any, len(warehouses))

	for i, warehouse := range warehouses {
		warehouse := warehouse
		var warehouseDescription []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.Warehouses.Describe(ctx, warehouse.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			warehouseDescription = schemas.WarehouseDescriptionToSchema(*describeResult)
		}

		var warehouseParameters []map[string]any
		if d.Get("with_parameters").(bool) {
			parameters, err := client.Warehouses.ShowParameters(ctx, warehouse.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			warehouseParameters = []map[string]any{schemas.WarehouseParametersToSchema(parameters)}
		}

		flattenedWarehouses[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.WarehouseToSchema(&warehouse)},
			resources.DescribeOutputAttributeName: warehouseDescription,
			resources.ParametersAttributeName:     warehouseParameters,
		}
	}

	err = d.Set("warehouses", flattenedWarehouses)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
