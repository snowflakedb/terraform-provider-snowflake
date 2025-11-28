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

// TODO [SNOW-2852837]: add describe output (handle correctly each subsection); add with_describe attribute
var semanticViewsSchema = map[string]*schema.Schema{
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
			},
		},
	},
}

func SemanticViews() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.SemanticViewDatasource), TrackingReadWrapper(datasources.SemanticViews, ReadSemanticViews)),
		Schema:      semanticViewsSchema,
		Description: "Data source used to get details of filtered semantic views. Filtering is aligned with the current possibilities for [SHOW SEMANTIC VIEWS](https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views) query." +
			" The results are encapsulated in one output collection `semantic_views`. DESCRIBE is not currently supported and will be added before promoting the resource to stable.",
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
		flattenedSemanticViews[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.SemanticViewToSchema(&semanticView)},
		}
	}
	if err := d.Set("semantic_views", flattenedSemanticViews); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
