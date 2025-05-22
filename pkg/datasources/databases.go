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

var databasesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC DATABASE for each database returned by SHOW DATABASES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"with_parameters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs SHOW PARAMETERS FOR DATABASE for each database returned by SHOW DATABASES. The output of describe is saved to the parameters field as a map. By default this value is set to true.",
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
	"limit": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Limits the number of rows returned. If the `limit.from` is set, then the limit wll start from the first element matched by the expression. The expression is only used to match with the first element, later on the elements are not matched by the prefix, but you can enforce a certain pattern with `starts_with` or `like`.",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rows": {
					Type:        schema.TypeInt,
					Required:    true,
					Description: "The maximum number of rows to return.",
				},
				"from": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Specifies a **case-sensitive** pattern that is used to match object name. After the first match, the limit on the number of rows will be applied.",
				},
			},
		},
	},
	"databases": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all database details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW DATABASES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowDatabaseSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE DATABASE.",
					Elem: &schema.Resource{
						Schema: schemas.DatabaseDescribeSchema,
					},
				},
				resources.ParametersAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PARAMETERS FOR DATABASE.",
					Elem: &schema.Resource{
						Schema: schemas.ShowDatabaseParametersSchema,
					},
				},
			},
		},
	},
}

func Databases() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Databases, ReadDatabases),
		Schema:      databasesSchema,
		Description: "Data source used to get details of filtered databases. Filtering is aligned with the current possibilities for [SHOW DATABASES](https://docs.snowflake.com/en/sql-reference/sql/show-databases) query (`like`, `starts_with`, and `limit` are all supported). The results of SHOW, DESCRIBE, and SHOW PARAMETERS IN are encapsulated in one output collection.",
	}
}

// ReadDatabases read the current snowflake account information.
func ReadDatabases(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	var opts sdk.ShowDatabasesOptions

	if likePattern, ok := d.GetOk("like"); ok {
		opts.Like = &sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		}
	}

	if startsWith, ok := d.GetOk("starts_with"); ok {
		opts.StartsWith = sdk.String(startsWith.(string))
	}

	if limit, ok := d.GetOk("limit"); ok && len(limit.([]any)) == 1 {
		limitMap := limit.([]any)[0].(map[string]any)

		rows := limitMap["rows"].(int)
		opts.LimitFrom = &sdk.LimitFrom{
			Rows: &rows,
		}

		if from, ok := limitMap["from"].(string); ok {
			opts.LimitFrom.From = &from
		}
	}

	databases, err := client.Databases.Show(ctx, &opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("databases_read")

	flattenedDatabases := make([]map[string]any, len(databases))

	for i, database := range databases {
		database := database
		var databaseDescription []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.Databases.Describe(ctx, database.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			databaseDescription = schemas.DatabaseDescriptionToSchema(*describeResult)
		}

		var databaseParameters []map[string]any
		if d.Get("with_parameters").(bool) {
			parameters, err := client.Databases.ShowParameters(ctx, database.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			databaseParameters = []map[string]any{schemas.DatabaseParametersToSchema(parameters)}
		}

		flattenedDatabases[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.DatabaseToSchema(&database)},
			resources.DescribeOutputAttributeName: databaseDescription,
			resources.ParametersAttributeName:     databaseParameters,
		}
	}

	err = d.Set("databases", flattenedDatabases)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
