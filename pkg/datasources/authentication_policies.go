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

var authenticationPoliciesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC AUTHENTICATION POLICY for each service returned by SHOW AUTHENTICATION POLICIES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"in":          extendedInSchema,
	"on":          onSchema,
	"authentication_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all authentication policies details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW AUTHENTICATION POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowAuthenticationPolicySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE AUTHENTICATION POLICY.",
					Elem: &schema.Resource{
						Schema: schemas.AuthenticationPolicyDescribeSchema,
					},
				},
			},
		},
	},
}

func AuthenticationPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.AuthenticationPoliciesDatasource), TrackingReadWrapper(datasources.AuthenticationPolicies, ReadAuthenticationPolicies)),
		Schema:      authenticationPoliciesSchema,
		Description: "Data source used to get details of filtered authentication policies. Filtering is aligned with the current possibilities for [SHOW AUTHENTICATION POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-authentication-policies) query." +
			" The results of SHOW and DESCRIBE are encapsulated in one output collection `authentication_policies`.",
	}
}

func ReadAuthenticationPolicies(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowAuthenticationPolicyRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)
	if err := errors.Join(
		handleExtendedIn(d, &req.In),
		handleOn(d, &req.On),
	); err != nil {
		return diag.FromErr(err)
	}

	authenticationPolicies, err := client.AuthenticationPolicies.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("authentication_policies_read")

	flattenedAuthenticationPolicies := make([]map[string]any, len(authenticationPolicies))
	for i, authenticationPolicy := range authenticationPolicies {
		authenticationPolicy := authenticationPolicy
		var authenticationPolicyDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.AuthenticationPolicies.Describe(ctx, authenticationPolicy.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			authenticationPolicyDescriptions = []map[string]any{schemas.AuthenticationPolicyDescriptionToSchema(describeResult)}
		}
		flattenedAuthenticationPolicies[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.AuthenticationPolicyToSchema(&authenticationPolicy)},
			resources.DescribeOutputAttributeName: authenticationPolicyDescriptions,
		}
	}
	if err := d.Set("authentication_policies", flattenedAuthenticationPolicies); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
