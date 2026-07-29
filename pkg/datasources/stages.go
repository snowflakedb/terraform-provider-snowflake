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

var stagesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC STAGE for each stage returned by SHOW STAGES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like": likeSchema,
	"in":   extendedInSchema,
	"stages": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all stages details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW STAGES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowStageSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE STAGE.",
					Elem: &schema.Resource{
						Schema: schemas.StageDatasourceDescribeSchema(),
					},
				},
			},
		},
	},
}

func Stages() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.StagesDatasource), TrackingReadWrapper(datasources.Stages, ReadStages)),
		Schema:      stagesSchema,
		Description: "Data source used to get details of filtered stages. Filtering is aligned with the current possibilities for [SHOW STAGES](https://docs.snowflake.com/en/sql-reference/sql/show-stages) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `stages`.",
	}
}

func ReadStages(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowStageRequest{}

	handleLike(d, &req.Like)
	err := handleExtendedIn(d, &req.In)
	if err != nil {
		return diag.FromErr(err)
	}

	stages, err := client.Stages.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("stages_read")

	flattenedStages := make([]map[string]any, len(stages))
	for i, stage := range stages {
		var stageDescriptions []map[string]any
		if d.Get("with_describe").(bool) {
			properties, err := client.Stages.Describe(ctx, stage.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			details, err := sdk.ParseStageDetails(properties)
			if err != nil {
				return diag.FromErr(err)
			}
			describeSchema, err := schemas.StageDatasourceToDatasourceSchema(*details)
			if err != nil {
				return diag.FromErr(err)
			}
			stageDescriptions = []map[string]any{describeSchema}
		}
		flattenedStages[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.StageToSchema(&stage)},
			resources.DescribeOutputAttributeName: stageDescriptions,
		}
	}
	if err := d.Set("stages", flattenedStages); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
