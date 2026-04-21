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

var externalVolumesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC EXTERNAL VOLUME for each external volume returned by SHOW EXTERNAL VOLUMES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like": likeSchema,
	"external_volumes": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all external volume details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW EXTERNAL VOLUMES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowExternalVolumeSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE EXTERNAL VOLUME.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeExternalVolumeSchema,
					},
				},
			},
		},
	},
}

func ExternalVolumes() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ExternalVolumesDatasource), TrackingReadWrapper(datasources.ExternalVolumes, ReadExternalVolumes)),
		Schema:      externalVolumesSchema,
		Description: "Data source used to get details of filtered external volumes. Filtering is aligned with the current possibilities for [SHOW EXTERNAL VOLUMES](https://docs.snowflake.com/en/sql-reference/sql/show-external-volumes) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `external_volumes`.",
	}
}

func ReadExternalVolumes(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowExternalVolumeRequest{}

	handleLike(d, &req.Like)

	externalVolumes, err := client.ExternalVolumes.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("external_volumes_read")

	flattenedExternalVolumes := make([]map[string]any, len(externalVolumes))
	for i, ev := range externalVolumes {
		var evDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			props, err := client.ExternalVolumes.Describe(ctx, ev.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			details, err := sdk.ParseExternalVolumeDescribed(props)
			if err != nil {
				return diag.FromErr(err)
			}
			evDescriptions = []map[string]any{schemas.ExternalVolumeDetailsToSchema(details)}
		}
		flattenedExternalVolumes[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.ExternalVolumeToSchema(&ev)},
			resources.DescribeOutputAttributeName: evDescriptions,
		}
	}
	if err := d.Set("external_volumes", flattenedExternalVolumes); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
