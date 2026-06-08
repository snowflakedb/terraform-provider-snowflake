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

var cortexAgentsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC AGENT for each object returned by SHOW AGENTS. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"in":          extendedInSchema,
	"cortex_agents": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all cortex agent details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW AGENTS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowCortexAgentSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE AGENT.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeCortexAgentDetailsSchema,
					},
				},
			},
		},
	},
}

func CortexAgents() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.CortexAgentsDatasource), TrackingReadWrapper(datasources.CortexAgents, ReadCortexAgents)),
		Schema:      cortexAgentsSchema,
		Description: "Data source used to get details of filtered Cortex agents. Filtering is aligned with the current possibilities for [SHOW AGENTS](https://docs.snowflake.com/en/sql-reference/sql/show-agents) query." +
			" The results of SHOW and DESCRIBE are encapsulated in one output collection `cortex_agents`.",
	}
}

func ReadCortexAgents(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowCortexAgentRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)
	if err := handleExtendedIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	cortexAgents, err := client.CortexAgents.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("cortex_agents_read")

	flattened := make([]map[string]any, len(cortexAgents))
	for i := range cortexAgents {
		ca := cortexAgents[i]
		var describeOut []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.CortexAgents.Describe(ctx, ca.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOut = []map[string]any{schemas.CortexAgentDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.CortexAgentToSchema(&ca)},
			resources.DescribeOutputAttributeName: describeOut,
		}
	}
	if err := d.Set("cortex_agents", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
