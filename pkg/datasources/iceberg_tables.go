package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var icebergTablesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC ICEBERG TABLE for each iceberg table returned by SHOW ICEBERG TABLES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"with_parameters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs SHOW PARAMETERS FOR ICEBERG TABLE for each iceberg table returned by SHOW ICEBERG TABLES. The output is saved to the parameters field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          inSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"iceberg_tables": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all iceberg table details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW ICEBERG TABLES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowIcebergTableSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE ICEBERG TABLE.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeIcebergTableSchema,
					},
				},
				resources.ParametersAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PARAMETERS FOR ICEBERG TABLE.",
					Elem: &schema.Resource{
						Schema: schemas.ShowIcebergTableParametersSchema,
					},
				},
			},
		},
	},
}

func IcebergTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.IcebergTablesDatasource), TrackingReadWrapper(datasources.IcebergTables, ReadIcebergTables)),
		Schema:      icebergTablesSchema,
		Description: "Data source used to get details of filtered iceberg tables. Filtering is aligned with the current possibilities for [SHOW ICEBERG TABLES](https://docs.snowflake.com/en/sql-reference/sql/show-iceberg-tables) query (`like`, `in`, `starts_with`, `limit`). The results of SHOW, DESCRIBE, and SHOW PARAMETERS are encapsulated in one output collection `iceberg_tables`.",
	}
}

func ReadIcebergTables(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	req := sdk.NewShowIcebergTableRequest()

	handleLike(d, &req.Like)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	icebergTables, err := client.IcebergTables.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("iceberg_tables_read")

	flattened := make([]map[string]any, len(icebergTables))
	for i := range icebergTables {
		table := icebergTables[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.IcebergTables.Describe(ctx, table.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = schemas.IcebergTableDetailsToSchema(details)
		}

		var tableParameters []map[string]any
		if d.Get("with_parameters").(bool) {
			parameters, err := client.IcebergTables.ShowParameters(ctx, table.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			tableParameters = []map[string]any{schemas.IcebergTableParametersToSchema(parameters, providerCtx)}
		}

		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.IcebergTableToSchema(&table)},
			resources.DescribeOutputAttributeName: describeOutput,
			resources.ParametersAttributeName:     tableParameters,
		}
	}
	if err := d.Set("iceberg_tables", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
