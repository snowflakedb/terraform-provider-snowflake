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

var semanticViewsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC SEMANTIC VIEW for each semantic view returned by SHOW SEMANTIC VIEWS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          inSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"semantic_views": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all semantic view details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW SEMANTIC VIEWS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowSemanticViewSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE SEMANTIC VIEW.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeSemanticViewSchema,
					},
				},
			},
		},
	},
}

func SemanticViews() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SemanticViewsDatasource), TrackingReadWrapper(datasources.SemanticViews, ReadSemanticViews)),
		Schema:      semanticViewsSchema,
		Description: "Data source used to get details of filtered semantic views. Filtering is aligned with the current possibilities for [SHOW SEMANTIC VIEWS](https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views) query." +
			" The results of SHOW and DESCRIBE are encapsulated in one output collection `semantic_views`.",
	}
}

func ReadSemanticViews(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowSemanticViewRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	err := handleIn(d, &req.In)
	if err != nil {
		return diag.FromErr(err)
	}
	handleLimitFrom(d, &req.Limit)

	semanticViews, err := client.SemanticViews.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("semantic_views_read")

	flattenedSemanticViews := make([]map[string]any, len(semanticViews))
	for i, semanticView := range semanticViews {
		semanticView := semanticView
		var semanticViewDetails [][]map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.SemanticViews.Describe(ctx, semanticView.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			semanticViewDetails = [][]map[string]any{schemas.SemanticViewDetailsToSchema(describeResult)}
		}
		flattenedSemanticViews[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.SemanticViewToSchema(&semanticView)},
			resources.DescribeOutputAttributeName: semanticViewDetails,
		}
	}
	if err := d.Set("semantic_views", flattenedSemanticViews); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
