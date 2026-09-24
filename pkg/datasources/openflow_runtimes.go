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

var openflowRuntimesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC OPENFLOW RUNTIME for each runtime returned by SHOW OPENFLOW RUNTIMES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          inSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"openflow_runtimes": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all Openflow runtime details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW OPENFLOW RUNTIMES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowOpenflowRuntimeSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE OPENFLOW RUNTIME.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeOpenflowRuntimeDetailsSchema,
					},
				},
			},
		},
	},
}

func OpenflowRuntimes() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.OpenflowRuntimesDatasource), TrackingReadWrapper(datasources.OpenflowRuntimes, ReadOpenflowRuntimes)),
		Schema:      openflowRuntimesSchema,
		Description: "Data source used to get details of filtered Openflow runtimes. The results of SHOW and DESCRIBE are encapsulated in one output collection `openflow_runtimes`.",
	}
}

func ReadOpenflowRuntimes(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowOpenflowRuntimeRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	runtimes, err := client.OpenflowRuntimes.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("openflow_runtimes_read")

	flattenedRuntimes := make([]map[string]any, len(runtimes))
	for i, runtime := range runtimes {
		var runtimeDetails []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.OpenflowRuntimes.Describe(ctx, runtime.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			runtimeDetails = []map[string]any{schemas.OpenflowRuntimeDetailsToSchema(describeResult)}
		}
		flattenedRuntimes[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.OpenflowRuntimeToSchema(&runtime)},
			resources.DescribeOutputAttributeName: runtimeDetails,
		}
	}
	if err := d.Set("openflow_runtimes", flattenedRuntimes); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
