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

var gitRepositoriesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC GIT REPOSITORY for each git repository returned by SHOW GIT REPOSITORIES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":  likeSchema,
	"in":    inSchema,
	"limit": limitFromSchema,
	"git_repositories": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all git repositories details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW GIT REPOSITORIES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowGitRepositorySchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE GIT REPOSITORY.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeGitRepositorySchema,
					},
				},
			},
		},
	},
}

func GitRepositories() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.GitRepositoriesDatasource), TrackingReadWrapper(datasources.GitRepositories, ReadGitRepositories)),
		Schema:      gitRepositoriesSchema,
		Description: "Data source used to get details of filtered git repositories. Filtering is aligned with the current possibilities for [SHOW GIT REPOSITORIES](https://docs.snowflake.com/en/sql-reference/sql/show-git-repositories) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `git_repositories`.",
	}
}

func ReadGitRepositories(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowGitRepositoryRequest{}

	handleLike(d, &req.Like)
	err := handleIn(d, &req.In)
	if err != nil {
		return diag.FromErr(err)
	}
	handleLimitFrom(d, &req.Limit)

	gitRepositories, err := client.GitRepositories.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("git_repositories_read")

	flattenedGitRepositories := make([]map[string]any, len(gitRepositories))
	for i, gitRepository := range gitRepositories {
		gitRepository := gitRepository
		var gitRepositoryDetails []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.GitRepositories.Describe(ctx, gitRepository.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			gitRepositoryDetails = []map[string]any{schemas.GitRepositoryDetailsToSchema(*describeResult)}
		}
		flattenedGitRepositories[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.GitRepositoryToSchema(&gitRepository)},
			resources.DescribeOutputAttributeName: gitRepositoryDetails,
		}
	}
	if err := d.Set("git_repositories", flattenedGitRepositories); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
