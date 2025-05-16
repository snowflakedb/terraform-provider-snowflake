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

var networkPoliciesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC NETWORK POLICY for each network policy returned by SHOW NETWORK POLICIES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"network_policies": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all network policies details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW NETWORK POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowNetworkPolicySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE NETWORK POLICIES.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeNetworkPolicySchema,
					},
				},
			},
		},
	},
}

func NetworkPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.NetworkPolicies, ReadNetworkPolicies),
		Schema:      networkPoliciesSchema,
		Description: "Data source used to get details of filtered network policies. Filtering is aligned with the current possibilities for [SHOW NETWORK POLICIES](https://docs.snowflake.com/en/sql-reference/sql/show-network-policies) query (`like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection.",
	}
}

func ReadNetworkPolicies(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowNetworkPolicyRequest()

	if likePattern, ok := d.GetOk("like"); ok {
		req.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}

	networkPolicies, err := client.NetworkPolicies.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("network_policies_read")

	flattenedNetworkPolicies := make([]map[string]any, len(networkPolicies))
	for i, networkPolicy := range networkPolicies {
		networkPolicy := networkPolicy

		var networkPolicyDescribeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.NetworkPolicies.Describe(ctx, networkPolicy.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			networkPolicyDescribeOutput = []map[string]any{schemas.NetworkPolicyPropertiesToSchema(describeResult)}
		}

		flattenedNetworkPolicies[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.NetworkPolicyToSchema(&networkPolicy)},
			resources.DescribeOutputAttributeName: networkPolicyDescribeOutput,
		}
	}

	if err = d.Set("network_policies", flattenedNetworkPolicies); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
