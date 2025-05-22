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

var streamlitsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC STREAMLIT for each streamlit returned by SHOW STREAMLITS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"in": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IN clause to filter the list of streamlits",
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
					Description:  "Returns records for the current database in use or for a specified database (db_name).",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema"},
				},
				"schema": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the current schema in use or a specified schema (schema_name).",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema"},
				},
			},
		},
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
	"streamlits": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all streamlits details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW STREAMLITS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowStreamlitSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE STREAMLIT.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeStreamlitSchema,
					},
				},
			},
		},
	},
}

func Streamlits() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Streamlits, ReadStreamlits),
		Schema:      streamlitsSchema,
		Description: "Data source used to get details of filtered streamlits. Filtering is aligned with the current possibilities for [SHOW STREAMLITS](https://docs.snowflake.com/en/sql-reference/sql/show-streamlits) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `streamlits`.",
	}
}

func ReadStreamlits(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowStreamlitRequest()

	if likePattern, ok := d.GetOk("like"); ok {
		req.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}
	if v, ok := d.GetOk("in"); ok {
		in := v.([]interface{})[0].(map[string]interface{})
		if v, ok := in["account"]; ok {
			account := v.(bool)
			if account {
				req.WithIn(sdk.In{Account: sdk.Bool(account)})
			}
		}
		if v, ok := in["database"]; ok {
			database := v.(string)
			if database != "" {
				req.WithIn(sdk.In{Database: sdk.NewAccountObjectIdentifier(database)})
			}
		}
		if v, ok := in["schema"]; ok {
			schema := v.(string)
			if schema != "" {
				req.WithIn(sdk.In{Schema: sdk.NewDatabaseObjectIdentifierFromFullyQualifiedName(schema)})
			}
		}
	}
	if v, ok := d.GetOk("limit"); ok {
		l := v.([]interface{})[0].(map[string]interface{})
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
	streamlits, err := client.Streamlits.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("streamlits_read")

	flattenedStreamlits := make([]map[string]any, len(streamlits))

	for i, streamlit := range streamlits {
		streamlit := streamlit
		var streamlitDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			descriptions, err := client.Streamlits.Describe(ctx, streamlit.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			streamlitDescriptions = make([]map[string]any, 1)
			streamlitDescriptions[0], err = schemas.StreamlitPropertiesToSchema(*descriptions)
			if err != nil {
				return diag.FromErr(err)
			}
		}

		flattenedStreamlits[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.StreamlitToSchema(&streamlit)},
			resources.DescribeOutputAttributeName: streamlitDescriptions,
		}
	}

	err = d.Set("streamlits", flattenedStreamlits)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
