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

var networkRulesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC NETWORK RULE for each network rule returned by SHOW NETWORK RULES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"in":          inSchema,
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"network_rules": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all network rules details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW NETWORK RULES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowNetworkRuleSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE NETWORK RULE.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeNetworkRuleSchema,
					},
				},
			},
		},
	},
}

func NetworkRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.NetworkRulesDatasource), TrackingReadWrapper(datasources.NetworkRules, ReadNetworkRules)),
		Schema:      networkRulesSchema,
		Description: "Data source used to get details of filtered network rules. Filtering is aligned with the current possibilities for [SHOW NETWORK RULES](https://docs.snowflake.com/en/sql-reference/sql/show-network-rules) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `network_rules`.",
	}
}

func ReadNetworkRules(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowNetworkRuleRequest()

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	networkRules, err := client.NetworkRules.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("network_rules_read")

	flattenedNetworkRules := make([]map[string]any, len(networkRules))
	for i, networkRule := range networkRules {
		var networkRuleDescribeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.NetworkRules.Describe(ctx, networkRule.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			networkRuleDescribeOutput = []map[string]any{schemas.NetworkRuleDetailsToSchema(describeResult)}
		}

		flattenedNetworkRules[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.NetworkRuleToSchema(&networkRule)},
			resources.DescribeOutputAttributeName: networkRuleDescribeOutput,
		}
	}

	if err := d.Set("network_rules", flattenedNetworkRules); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
