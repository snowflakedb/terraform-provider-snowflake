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

var usersSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC USER for each user returned by SHOW USERS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"with_parameters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs SHOW PARAMETERS FOR USER for each user returned by SHOW USERS. The output of describe is saved to the parameters field as a map. By default this value is set to true.",
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
	"users": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all user details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW USERS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowUserSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE USER.",
					Elem: &schema.Resource{
						Schema: schemas.UserDescribeSchema,
					},
				},
				resources.ParametersAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PARAMETERS FOR USER.",
					Elem: &schema.Resource{
						Schema: schemas.ShowUserParametersSchema,
					},
				},
			},
		},
	},
}

func Users() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Users, ReadUsers),
		Schema:      usersSchema,
		Description: "Data source used to get details of filtered users. Filtering is aligned with the current possibilities for [SHOW USERS](https://docs.snowflake.com/en/sql-reference/sql/show-users) query. The results of SHOW, DESCRIBE, and SHOW PARAMETERS IN are encapsulated in one output collection. Important note is that when querying users you don't have permissions to, the querying options are limited. You won't get almost any field in `show_output` (only empty or default values), the DESCRIBE command will return error when called, so you have to set `with_describe = false`; the SHOW PARAMETERS command will return error when called too, so you have to set `with_parameters = false`.",
	}
}

func ReadUsers(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	var opts sdk.ShowUserOptions

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
		opts.Limit = &rows

		if from, ok := limitMap["from"].(string); ok {
			opts.From = &from
		}
	}

	users, err := client.Users.Show(ctx, &opts)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("users_read")

	flattenedUsers := make([]map[string]any, len(users))

	for i, user := range users {
		user := user
		var userDescription []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.Users.Describe(ctx, user.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			userDescription = schemas.UserDescriptionToSchema(*describeResult)
		}

		var userParameters []map[string]any
		if d.Get("with_parameters").(bool) {
			parameters, err := client.Users.ShowParameters(ctx, user.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			userParameters = []map[string]any{schemas.UserParametersToSchema(parameters)}
		}

		flattenedUsers[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.UserToSchema(&user)},
			resources.DescribeOutputAttributeName: userDescription,
			resources.ParametersAttributeName:     userParameters,
		}
	}

	err = d.Set("users", flattenedUsers)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
