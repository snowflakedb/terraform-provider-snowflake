package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var computePoolsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC COMPUTE POOL for each compute pool returned by SHOW COMPUTE POOLS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"compute_pools": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all compute pools details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW COMPUTE POOLS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowComputePoolSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE COMPUTE POOL.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeComputePoolSchema,
					},
				},
			},
		},
	},
}

func ComputePools() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.ComputePools, ReadComputePools),
		Schema:      computePoolsSchema,
		Description: "Data source used to get details of filtered compute pools. Filtering is aligned with the current possibilities for [SHOW COMPUTE POOLS](https://docs.snowflake.com/en/sql-reference/sql/show-compute-pools) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `compute_pools`.",
	}
}

func ReadComputePools(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowComputePoolRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	computePools, err := client.ComputePools.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("compute_pools_read")

	flattenedComputePools := make([]map[string]any, len(computePools))
	for i, computePool := range computePools {
		computePool := computePool
		var computePoolDetails []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.ComputePools.Describe(ctx, computePool.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			computePoolDetails = []map[string]any{schemas.ComputePoolDetailsToSchema(*describeResult)}
		}
		flattenedComputePools[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.ComputePoolToSchema(&computePool)},
			resources.DescribeOutputAttributeName: computePoolDetails,
		}
	}
	if err := d.Set("compute_pools", flattenedComputePools); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
