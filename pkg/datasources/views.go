package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var viewsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC VIEW for each view returned by SHOW VIEWS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"in": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IN clause to filter the list of views",
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
	"limit": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Limits the number of rows returned. If the `limit.from` is set, then the limit will start from the first element matched by the expression. The expression is only used to match with the first element, later on the elements are not matched by the prefix, but you can enforce a certain pattern with `starts_with` or `like`.",
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
	"views": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all views details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW VIEWS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowViewSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE VIEW.",
					Elem: &schema.Resource{
						Schema: schemas.ViewDescribeSchema,
					},
				},
			},
		},
	},
}

func Views() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Views, ReadViews),
		Schema:      viewsSchema,
		Description: "Data source used to get details of filtered views. Filtering is aligned with the current possibilities for [SHOW VIEWS](https://docs.snowflake.com/en/sql-reference/sql/show-views) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `views`.",
	}
}

func ReadViews(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowViewRequest()

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

	if v, ok := d.GetOk("limit"); ok {
		l := v.([]any)[0].(map[string]any)
		limit := sdk.LimitFrom{}
		if v, ok := l["rows"]; ok {
			rows := v.(int)
			limit.Rows = sdk.Int(rows)
		}
		if v, ok := l["from"]; ok {
			from := v.(string)
			limit.From = sdk.String(from)
		}
		req.WithLimit(limit)
	}

	views, err := client.Views.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("views_read")

	flattenedViews := make([]map[string]any, len(views))
	for i, view := range views {
		view := view
		var viewDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			describeOutput, err := client.Views.Describe(ctx, view.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			viewDescriptions = schemas.ViewDescriptionToSchema(describeOutput)
		}

		flattenedViews[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.ViewToSchema(&view)},
			resources.DescribeOutputAttributeName: viewDescriptions,
		}
	}

	if err := d.Set("views", flattenedViews); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
