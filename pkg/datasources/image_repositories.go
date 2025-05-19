package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var imageRepositoriesSchema = map[string]*schema.Schema{
	"like": likeSchema,
	"in":   inSchema,
	"image_repositories": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all image repositories details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW IMAGE REPOSITORIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowImageRepositorySchema,
					},
				},
			},
		},
	},
}

func ImageRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ImageRepositoriesDatasource), TrackingReadWrapper(datasources.ImageRepositories, ReadImageRepositories)),
		Schema:      imageRepositoriesSchema,
		Description: "Data source used to get details of filtered image repositories. Filtering is aligned with the current possibilities for [SHOW IMAGE REPOSITORIES](https://docs.snowflake.com/en/sql-reference/sql/show-image-repositories) query. The results of SHOW are encapsulated in one output collection `image_repositories`.",
	}
}

func ReadImageRepositories(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowImageRepositoryRequest{}

	handleLike(d, &req.Like)
	err := handleIn(d, &req.In)
	if err != nil {
		return diag.FromErr(err)
	}

	imageRepositories, err := client.ImageRepositories.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("image_repositories_read")

	flattenedImageRepositories := make([]map[string]any, len(imageRepositories))
	for i, imageRepository := range imageRepositories {
		imageRepository := imageRepository
		flattenedImageRepositories[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.ImageRepositoryToSchema(&imageRepository)},
		}
	}
	if err := d.Set("image_repositories", flattenedImageRepositories); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
