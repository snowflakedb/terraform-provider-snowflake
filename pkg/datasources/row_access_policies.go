package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var rowAccessPoliciesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC ROW ACCESS POLICY for each row access policy returned by SHOW ROW ACCESS POLICIES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"in": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IN clause to filter the list of row access policies",
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"account": {
					Type:         schema.TypeBool,
					Optional:     true,
					Description:  "Returns records for the entire account.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				},
				"database": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the current database in use or for a specified database.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				},
				"schema": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the current schema in use or a specified schema. Use fully qualified name.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				},
				"application": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the specified application.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				},
				"application_package": {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "Returns records for the specified application package.",
					ExactlyOneOf: []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
				},
			},
		},
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
	"row_access_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all views details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW ROW ACCESS POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowRowAccessPolicySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE ROW ACCESS POLICY.",
					Elem: &schema.Resource{
						Schema: schemas.RowAccessPolicyDescribeSchema,
					},
				},
			},
		},
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
}

func RowAccessPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.RowAccessPolicies, ReadRowAccessPolicies),
		Schema:      rowAccessPoliciesSchema,
		Description: "Data source used to get details of filtered row access policies. Filtering is aligned with the current possibilities for [SHOW ROW ACCESS POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-row-access-policies) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `row_access_policies`.",
	}
}

func ReadRowAccessPolicies(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowRowAccessPolicyRequest()

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
		if v, ok := in["application"]; ok {
			if application := v.(string); application != "" {
				req.In = &sdk.ExtendedIn{Application: sdk.NewAccountObjectIdentifier(application)}
			}
		}
		if v, ok := in["application_package"]; ok {
			if applicationPackage := v.(string); applicationPackage != "" {
				req.In = &sdk.ExtendedIn{ApplicationPackage: sdk.NewAccountObjectIdentifier(applicationPackage)}
			}
		}
	}

	if likePattern, ok := d.GetOk("like"); ok {
		req.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}

	if v, ok := d.GetOk("limit"); ok {
		l := v.([]any)[0].(map[string]any)
		limit := &sdk.LimitFrom{}
		if v, ok := l["rows"]; ok {
			rows := v.(int)
			limit.Rows = sdk.Int(rows)
		}
		if v, ok := l["from"]; ok {
			from := v.(string)
			limit.From = sdk.String(from)
		}
		req.WithLimit(*limit)
	}

	rowAccessPolicies, err := client.RowAccessPolicies.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("row_access_policies_read")

	flattenedRowAccessPolicies := make([]map[string]any, len(rowAccessPolicies))
	for i, policy := range rowAccessPolicies {
		policy := policy
		var policyDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			describeOutput, err := client.RowAccessPolicies.Describe(ctx, policy.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			policyDescriptions = []map[string]any{schemas.RowAccessPolicyDescriptionToSchema(*describeOutput)}
		}

		flattenedRowAccessPolicies[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.RowAccessPolicyToSchema(&policy)},
			resources.DescribeOutputAttributeName: policyDescriptions,
		}
	}
	if err := d.Set("row_access_policies", flattenedRowAccessPolicies); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
