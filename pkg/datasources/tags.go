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

var tagsSchema = map[string]*schema.Schema{
	"like": likeSchema,
	"in":   extendedInSchema,
	"tags": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all tags details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW TAGS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowTagSchema,
					},
				},
			},
		},
	},
}

func Tags() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Tags, ReadTags),
		Schema:      tagsSchema,
		Description: "Data source used to get details of filtered tags. Filtering is aligned with the current possibilities for [SHOW TAGS](https://docs.snowflake.com/en/sql-reference/sql/show-tags) query. The results of SHOW are encapsulated in one output collection `tags`.",
	}
}

func ReadTags(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowTagRequest{}

	handleLike(d, &req.Like)
	err := handleExtendedIn(d, &req.In)
	if err != nil {
		return diag.FromErr(err)
	}

	tags, err := client.Tags.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("tags_read")

	flattenedTags := make([]map[string]any, len(tags))
	for i, tag := range tags {
		tag := tag
		flattenedTags[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.TagToSchema(&tag)},
		}
	}
	if err := d.Set("tags", flattenedTags); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
