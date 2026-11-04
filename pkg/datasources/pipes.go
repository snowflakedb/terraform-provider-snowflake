package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/docs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var pipesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC PIPE for each pipe returned by SHOW PIPES. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"with_parameters": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs SHOW PARAMETERS FOR PIPE for each pipe returned by SHOW PIPES. The output of describe is saved to the parameters field as a map. By default this value is set to true.",
	},
	"like": likeSchema,
	"in":   extendedInSchema,
	// TODO(SNOW-4410188): add starts_with and limit when Snowflake applies STARTS WITH and LIMIT in SHOW PIPES.
	"pipes": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all pipe details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PIPES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowPipeSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE PIPE.",
					Elem: &schema.Resource{
						Schema: schemas.DescribePipeDetailsSchema,
					},
				},
				resources.ParametersAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW PARAMETERS FOR PIPE.",
					Elem: &schema.Resource{
						Schema: schemas.ShowPipeParametersSchema,
					},
				},
			},
		},
	},
}

func Pipes() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.PipesDatasource), TrackingReadWrapper(datasources.Pipes, ReadPipes)),
		Schema:      pipesSchema,
		Description: docs.WithPreviewFeatureMark() + "Data source used to get details of filtered pipes. Filtering is aligned with the current possibilities for [SHOW PIPES](https://docs.snowflake.com/en/sql-reference/sql/show-pipes) query (`like` and `in`, including `application` and `application_package`). " +
			"The results of SHOW, DESCRIBE, and SHOW PARAMETERS are encapsulated in one output collection `pipes`.",
	}
}

func ReadPipes(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	providerCtx := meta.(*provider.Context)
	client := providerCtx.Client
	req := sdk.ShowPipeRequest{}

	handleLike(d, &req.Like)
	if err := handleExtendedIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	pipes, err := client.Pipes.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("pipes_read")

	flattened := make([]map[string]any, len(pipes))
	for i := range pipes {
		pipe := pipes[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.Pipes.Describe(ctx, pipe.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.PipeDetailsToSchema(details)}
		}
		var parameters []map[string]any
		if d.Get("with_parameters").(bool) {
			pipeParameters, err := client.Pipes.ShowParameters(ctx, pipe.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			parameters = []map[string]any{schemas.PipeParametersToSchema(pipeParameters, providerCtx)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.PipeToSchema(&pipe)},
			resources.DescribeOutputAttributeName: describeOutput,
			resources.ParametersAttributeName:     parameters,
		}
	}
	if err := d.Set("pipes", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
