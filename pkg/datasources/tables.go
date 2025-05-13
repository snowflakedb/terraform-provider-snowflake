package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var tablesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC TABLE for each table returned by SHOW TABLES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"in":          extendedInSchema,
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"tables": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all tables details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW TABLES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowTableSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE TABLES.",
					Elem: &schema.Resource{
						Schema: schemas.TableDescribeSchema,
					},
				},
			},
		},
	},
}

func Tables() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.TablesDatasource), TrackingReadWrapper(datasources.Tables, ReadTables)),
		Schema:      tablesSchema,
		Description: "Datasource used to get details of filtered tables. Filtering is aligned with the current possibilities for [SHOW TABLES](https://docs.snowflake.com/en/sql-reference/sql/show-tables) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `tables`.",
	}
}

func ReadTables(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowTableRequest()

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)
	err := handleExtendedIn(d, &req.In)
	if err != nil {
		return diag.FromErr(err)
	}

	tables, err := client.Tables.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("tables_read")

	flattenedTables := make([]map[string]any, len(tables))
	for i, table := range tables {
		table := table
		var tableDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			describeOutput, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(table.ID()))
			if err != nil {
				return diag.FromErr(err)
			}
			tableDescriptions = schemas.TableDescriptionToSchema(describeOutput)
		}

		flattenedTables[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.TableToSchema(&table)},
			resources.DescribeOutputAttributeName: tableDescriptions,
		}
	}

	if err := d.Set("tables", flattenedTables); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
