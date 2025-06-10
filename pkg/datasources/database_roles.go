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

var databaseRolesSchema = map[string]*schema.Schema{
	"in_database": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The database from which to return the database roles from.",
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
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
	"database_roles": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all database role details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW DATABASE ROLES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowDatabaseRoleSchema,
					},
				},
			},
		},
	},
}

func DatabaseRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.DatabaseRoles, ReadDatabaseRoles),
		Schema:      databaseRolesSchema,
		Description: "Data source used to get details of filtered database roles. Filtering is aligned with the current possibilities for [SHOW DATABASE ROLES](https://docs.snowflake.com/en/sql-reference/sql/show-database-roles) query (`like` and `limit` are supported). The results of SHOW is encapsulated in show_output collection.",
	}
}

func ReadDatabaseRoles(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowDatabaseRoleRequest(sdk.NewAccountObjectIdentifier(d.Get("in_database").(string)))

	if likePattern, ok := d.GetOk("like"); ok {
		req.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}

	if limit, ok := d.GetOk("limit"); ok && len(limit.([]any)) == 1 {
		limitMap := limit.([]any)[0].(map[string]any)

		rows := limitMap["rows"].(int)
		limitFrom := sdk.LimitFrom{
			Rows: &rows,
		}

		if from, ok := limitMap["from"].(string); ok {
			limitFrom.From = &from
		}

		req.WithLimit(limitFrom)
	}

	databaseRoles, err := client.DatabaseRoles.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("database_roles_read")

	flattenedDatabaseRoles := make([]map[string]any, len(databaseRoles))
	for i, databaseRole := range databaseRoles {
		databaseRole := databaseRole
		flattenedDatabaseRoles[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.DatabaseRoleToSchema(&databaseRole)},
		}
	}

	err = d.Set("database_roles", flattenedDatabaseRoles)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
