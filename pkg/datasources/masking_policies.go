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

var maskingPoliciesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC MASKING POLICY for each masking policy returned by SHOW MASKING POLICIES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"in": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "IN clause to filter the list of masking policies",
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
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Returns records for the current database in use or for a specified database.",
					ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
					ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
				},
				"schema": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Returns records for the current schema in use or a specified schema. Use fully qualified name.",
					ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
					ValidateDiagFunc: resources.IsValidIdentifier[sdk.DatabaseObjectIdentifier](),
				},
				"application": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Returns records for the specified application.",
					ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
					ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
				},
				"application_package": {
					Type:             schema.TypeString,
					Optional:         true,
					Description:      "Returns records for the specified application package.",
					ExactlyOneOf:     []string{"in.0.account", "in.0.database", "in.0.schema", "in.0.application", "in.0.application_package"},
					ValidateDiagFunc: resources.IsValidIdentifier[sdk.AccountObjectIdentifier](),
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
	"masking_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all views details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW MASKING POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowMaskingPolicySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE MASKING POLICY.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeMaskingPolicySchema,
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

func MaskingPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.MaskingPolicies, ReadMaskingPolicies),
		Schema:      maskingPoliciesSchema,
		Description: "Data source used to get details of filtered masking policies. Filtering is aligned with the current possibilities for [SHOW MASKING POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-masking-policies) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `masking_policies`.",
	}
}

func ReadMaskingPolicies(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowMaskingPolicyOptions{}

	handleLike(d, &req.Like)
	handleLimitFrom(d, &req.Limit)
	err := handleExtendedIn(d, &req.In)
	if err != nil {
		return diag.FromErr(err)
	}

	maskingPolicies, err := client.MaskingPolicies.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("masking_policies_read")

	flattenedMaskingPolicies := make([]map[string]any, len(maskingPolicies))
	for i, policy := range maskingPolicies {
		policy := policy
		var policyDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			describeOutput, err := client.MaskingPolicies.Describe(ctx, policy.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			policyDescriptions = []map[string]any{schemas.MaskingPolicyDescriptionToSchema(*describeOutput)}
		}

		flattenedMaskingPolicies[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.MaskingPolicyToSchema(&policy)},
			resources.DescribeOutputAttributeName: policyDescriptions,
		}
	}
	if err := d.Set("masking_policies", flattenedMaskingPolicies); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
