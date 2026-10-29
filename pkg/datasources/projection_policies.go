package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var projectionPoliciesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC PROJECTION POLICY for each projection policy returned by SHOW PROJECTION POLICIES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          inSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"projection_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all projection policy details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PROJECTION POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowProjectionPolicySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE PROJECTION POLICY.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeProjectionPolicyDetailsSchema,
					},
				},
			},
		},
	},
}

func ProjectionPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ProjectionPoliciesDatasource), TrackingReadWrapper(datasources.ProjectionPolicies, ReadProjectionPolicies)),
		Schema:      projectionPoliciesSchema,
		Description: "Data source used to get details of filtered projection policies. Filtering is aligned with the current possibilities for [SHOW PROJECTION POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-projection-policies) query. " +
			"The results of SHOW and DESCRIBE are encapsulated in one output collection `projection_policies`.",
	}
}

func ReadProjectionPolicies(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowProjectionPolicyRequest{}

	handleLike(d, &req.Like)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	projectionPolicies, err := client.ProjectionPolicies.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("projection_policies_read")

	flattened := make([]map[string]any, len(projectionPolicies))
	for i := range projectionPolicies {
		policy := projectionPolicies[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.ProjectionPolicies.Describe(ctx, policy.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.ProjectionPolicyDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.ProjectionPolicyToSchema(&policy)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("projection_policies", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
