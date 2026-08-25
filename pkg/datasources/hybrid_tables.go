package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var hybridTablesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC TABLE for each hybrid table returned by SHOW HYBRID TABLES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"with_parameters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs SHOW PARAMETERS FOR TABLE for each hybrid table returned by SHOW HYBRID TABLES. The output is saved to the parameters field. By default this value is set to true.",
	},
	"with_keys": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Runs SHOW PRIMARY KEYS, SHOW UNIQUE KEYS, and SHOW IMPORTED KEYS for each hybrid table returned by SHOW HYBRID TABLES. The merged constraints are saved to the show_keys_output field, grouped by constraint name and ordered by kind, then by column names. By default this value is set to false.",
	},
	"with_indexes": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Runs SHOW INDEXES for each hybrid table returned by SHOW HYBRID TABLES. The output is saved to the show_indexes field. By default this value is set to false.",
	},
	"like":        likeSchema,
	"in":          inSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"hybrid_tables": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all hybrid table details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW HYBRID TABLES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowHybridTableSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE TABLE.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeHybridTableSchema,
					},
				},
				resources.ParametersAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PARAMETERS FOR TABLE.",
					Elem: &schema.Resource{
						Schema: schemas.ShowHybridTableParametersSchema,
					},
				},
				resources.ShowKeysOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the result of `SHOW PRIMARY KEYS`, `SHOW UNIQUE KEYS`, and `SHOW IMPORTED KEYS` for the given hybrid table, merged and grouped by constraint name and ordered by kind, then by column names. The `referenced_table`, `referenced_columns`, `delete_rule`, and `update_rule` fields are populated for FOREIGN KEY constraints only.",
					Elem: &schema.Resource{
						Schema: schemas.ShowHybridTableConstraintSchema,
					},
				},
				resources.ShowIndexesAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW INDEXES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowHybridTableIndexSchema,
					},
				},
			},
		},
	},
}

func HybridTables() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.HybridTablesDatasource), TrackingReadWrapper(datasources.HybridTables, ReadHybridTables)),
		Schema:      hybridTablesSchema,
		Description: "Data source used to get details of filtered hybrid tables. Filtering is aligned with the current possibilities for [SHOW HYBRID TABLES](https://docs.snowflake.com/en/sql-reference/sql/show-hybrid-tables) query (`like`, `in`, `starts_with`, `limit`). The results of SHOW, DESCRIBE, SHOW PARAMETERS, SHOW PRIMARY KEYS, SHOW UNIQUE KEYS, SHOW IMPORTED KEYS, and SHOW INDEXES are encapsulated in one output collection `hybrid_tables`.",
	}
}

func ReadHybridTables(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	req := sdk.NewShowHybridTableRequest()

	handleLike(d, &req.Like)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	hybridTables, err := client.HybridTables.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("hybrid_tables_read")

	flattened := make([]map[string]any, len(hybridTables))
	for i := range hybridTables {
		table := hybridTables[i]
		id := table.ID()

		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.HybridTables.Describe(ctx, id)
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = schemas.HybridTableDetailsListToSchema(details)
		}

		var tableParameters []map[string]any
		if d.Get("with_parameters").(bool) {
			parameters, err := client.HybridTables.ShowParameters(ctx, id)
			if err != nil {
				return diag.FromErr(err)
			}
			tableParameters = []map[string]any{schemas.HybridTableParametersToSchema(parameters, providerCtx)}
		}

		var keysOutput []map[string]any
		if d.Get("with_keys").(bool) {
			constraints, err := client.HybridTables.GetConstraints(ctx, id)
			if err != nil {
				return diag.FromErr(err)
			}
			keysOutput = collections.Map(constraints, func(c sdk.HybridTableConstraint) map[string]any {
				return schemas.HybridTableConstraintToSchema(&c)
			})
		}

		var indexesOutput []map[string]any
		if d.Get("with_indexes").(bool) {
			indexes, err := client.HybridTables.ShowIndexes(ctx, sdk.NewShowIndexesHybridTableRequest().WithIn(sdk.TableIn{Table: id}))
			if err != nil {
				return diag.FromErr(err)
			}
			indexesOutput = collections.Map(indexes, func(idx sdk.HybridTableIndex) map[string]any {
				return schemas.HybridTableIndexToSchema(&idx)
			})
		}

		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.HybridTableToSchema(&table)},
			resources.DescribeOutputAttributeName: describeOutput,
			resources.ParametersAttributeName:     tableParameters,
			resources.ShowKeysOutputAttributeName: keysOutput,
			resources.ShowIndexesAttributeName:    indexesOutput,
		}
	}
	if err := d.Set("hybrid_tables", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
