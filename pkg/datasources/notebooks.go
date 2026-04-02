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

var notebooksSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC NOTEBOOK for each notebook returned by SHOW NOTEBOOKS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"notebooks": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all notebooks details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW NOTEBOOKS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowNotebookSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE NOTEBOOK",
					Elem: &schema.Resource{
						Schema: schemas.DescribeNotebookSchema,
					},
				},
			},
		},
	},
}

func Notebooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.NotebooksDatasource), TrackingReadWrapper(datasources.Notebooks, ReadNotebooks)),
		Schema:      notebooksSchema,
		Description: "Data source used to get details of filtered notebooks. Filtering is aligned with the current possibilities for [SHOW NOTEBOOKS](https://docs.snowflake.com/en/sql-reference/sql/show-notebooks) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `notebooks`.",
	}
}

func ReadNotebooks(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowNotebookRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	notebooks, err := client.Notebooks.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("notebooks_read")

	flattenedNotebooks := make([]map[string]any, len(notebooks))
	for i, notebook := range notebooks {
		notebook := notebook
		var notebookDetails []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.Notebooks.Describe(ctx, notebook.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			notebookDetails = []map[string]any{schemas.NotebookDetailsToSchema(describeResult)}
		}
		flattenedNotebooks[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.NotebookToSchema(&notebook)},
			resources.DescribeOutputAttributeName: notebookDetails,
		}
	}
	if err := d.Set("notebooks", flattenedNotebooks); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
