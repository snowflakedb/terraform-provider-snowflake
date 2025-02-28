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
	"in": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IN clause to filter the list of tables",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"account": {
					Type:         schema.TypeBool,
					Optional:     true,
					Description:  "Returns records for the entire account.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema"},
				},
				"database": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the current database in use or for a specified database.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema"},
				},
				"schema": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the current schema in use or a specified schema. Use fully qualified name.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema"},
				},
			},
		},
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"starts_with": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-sensitive** characters indicating the beginning of the object name.",
	},
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
		Description: "Datasource used to get details of filtered tables. Filtering is aligned with the current possibilities for [SHOW VIEWS](https://docs.snowflake.com/en/sql-reference/sql/show-tables) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `tables`.",
	}
}

func ReadTables(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowTableRequest()

	if v, ok := d.GetOk("in"); ok {
		in := v.([]any)[0].(map[string]any)
		if v, ok := in["account"]; ok && v.(bool) {
			req.WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Bool(true)}})
		}
		if v, ok := in["database"]; ok {
			database := v.(string)
			if database != "" {
				req.WithIn(sdk.ExtendedIn{In: sdk.In{Database: sdk.NewAccountObjectIdentifier(database)}})
			}
		}
		if v, ok := in["schema"]; ok {
			schema := v.(string)
			if schema != "" {
				schemaId, err := sdk.ParseDatabaseObjectIdentifier(schema)
				if err != nil {
					return diag.FromErr(err)
				}
				req.WithIn(sdk.ExtendedIn{In: sdk.In{Schema: schemaId}})
			}
		}
	}

	if likePattern, ok := d.GetOk("like"); ok {
		req.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}

	if v, ok := d.GetOk("starts_with"); ok {
		req.WithStartsWith(v.(string))
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
