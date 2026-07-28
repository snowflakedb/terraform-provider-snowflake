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

var fileFormatsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC FILE FORMAT for each file format returned by SHOW FILE FORMATS. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like": likeSchema,
	"in":   inSchema,
	"file_formats": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all file formats details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW FILE FORMATS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowFileFormatSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE FILE FORMAT. Because every file format type returns a different set of properties, the type-specific properties are nested under the field named after the file format type (only one of them is filled at a time).",
					Elem: &schema.Resource{
						Schema: schemas.FileFormatsDatasourceDescribeSchema(),
					},
				},
			},
		},
	},
}

func FileFormats() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.FileFormatsDatasource), TrackingReadWrapper(datasources.FileFormats, ReadFileFormats)),
		Schema:      fileFormatsSchema,
		Description: "Data source used to get details of filtered file formats. Filtering is aligned with the current possibilities for [SHOW FILE FORMATS](https://docs.snowflake.com/en/sql-reference/sql/show-file-formats) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `file_formats`.",
	}
}

func ReadFileFormats(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowFileFormatRequest{}

	handleLike(d, &req.Like)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	fileFormats, err := client.FileFormats.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("file_formats_read")

	flattenedFileFormats := make([]map[string]any, len(fileFormats))
	for i, fileFormat := range fileFormats {
		var fileFormatDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.FileFormats.DescribeAllDetails(ctx, fileFormat.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			fileFormatDescriptions = []map[string]any{schemas.FileFormatAllDetailsToSchema(*details)}
		}
		flattenedFileFormats[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.FileFormatToSchema(&fileFormat)},
			resources.DescribeOutputAttributeName: fileFormatDescriptions,
		}
	}
	return diag.FromErr(d.Set("file_formats", flattenedFileFormats))
}
