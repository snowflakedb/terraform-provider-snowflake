package datasources

import (
	"context"
	"errors"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var sessionPoliciesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC SESSION POLICY for each object returned by SHOW SESSION POLICIES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"in":          extendedInSchema,
	"on":          onSchema,
	"session_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all session policy details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW SESSION POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowSessionPolicySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE SESSION POLICY.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeSessionPolicyDetailsSchema,
					},
				},
			},
		},
	},
}

func SessionPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SessionPoliciesDatasource), TrackingReadWrapper(datasources.SessionPolicies, ReadSessionPolicies)),
		Schema:      sessionPoliciesSchema,
		Description: "Data source used to get details of filtered session policies. Filtering is aligned with the current possibilities for [SHOW SESSION POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-session-policies) query." +
			" The results of SHOW and DESCRIBE are encapsulated in one output collection `session_policies`.",
	}
}

func ReadSessionPolicies(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowSessionPolicyRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)
	if err := errors.Join(
		handleExtendedIn(d, &req.In),
		handleOn(d, &req.On),
	); err != nil {
		return diag.FromErr(err)
	}

	sessionPolicies, err := client.SessionPolicies.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("session_policies_read")

	flattened := make([]map[string]any, len(sessionPolicies))
	for i := range sessionPolicies {
		sp := sessionPolicies[i]
		var describeOut []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.SessionPolicies.DescribeDetails(ctx, sp.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOut = []map[string]any{schemas.SessionPolicyDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.SessionPolicyToSchema(&sp)},
			resources.DescribeOutputAttributeName: describeOut,
		}
	}
	if err := d.Set("session_policies", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
