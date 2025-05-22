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

var streamsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC STREAM for each user returned by SHOW STREAMS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          extendedInSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"streams": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all streams details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW STREAMS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowStreamSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE STREAM.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeStreamSchema,
					},
				},
			},
		},
	},
}

func Streams() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Streams, ReadStreams),
		Schema:      streamsSchema,
		Description: "Data source used to get details of filtered streams. Filtering is aligned with the current possibilities for [SHOW STREAMS](https://docs.snowflake.com/en/sql-reference/sql/show-streams) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `streams`.",
	}
}

func ReadStreams(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowStreamRequest{}

	handleLike(d, &req.Like)
	handleLimitFrom(d, &req.Limit)
	handleStartsWith(d, &req.StartsWith)
	err := handleExtendedIn(d, &req.In)
	if err != nil {
		return diag.FromErr(err)
	}

	streams, err := client.Streams.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("streams_read")

	flattenedStreams := make([]map[string]any, len(streams))
	for i, stream := range streams {
		stream := stream
		var streamDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			describeOutput, err := client.Streams.Describe(ctx, stream.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			streamDescriptions = []map[string]any{schemas.StreamDescriptionToSchema(*describeOutput)}
		}

		flattenedStreams[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.StreamToSchema(&stream)},
			resources.DescribeOutputAttributeName: streamDescriptions,
		}
	}
	if err := d.Set("streams", flattenedStreams); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
