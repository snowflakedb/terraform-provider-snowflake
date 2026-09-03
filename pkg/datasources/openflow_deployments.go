package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var openflowDeploymentsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC OPENFLOW DEPLOYMENT for each deployment returned by SHOW OPENFLOW DEPLOYMENTS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"openflow_deployments": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all Openflow deployment details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW OPENFLOW DEPLOYMENTS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowOpenflowDeploymentSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE OPENFLOW DEPLOYMENT.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeOpenflowDeploymentSchema,
					},
				},
			},
		},
	},
}

func OpenflowDeployments() *schema.Resource {
	return &schema.Resource{
		// TODO(SNOW-4039167): Add PreviewFeatureReadWrapper when this data source is moved to the production
		// provider. It is registered only in the acceptance test provider for now, so there is no preview
		// feature to gate on yet.
		ReadContext: TrackingReadWrapper(datasources.OpenflowDeployments, ReadOpenflowDeployments),
		Schema:      openflowDeploymentsSchema,
		Description: "Data source used to get details of filtered Openflow deployments. Both Snowflake-managed and BYOC deployments are returned; the `type` field in `show_output` distinguishes them. The results of SHOW and DESCRIBE are encapsulated in one output collection `openflow_deployments`.",
	}
}

func ReadOpenflowDeployments(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowOpenflowDeploymentRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	deployments, err := client.OpenflowDeployments.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("openflow_deployments_read")

	flattenedDeployments := make([]map[string]any, len(deployments))
	for i, deployment := range deployments {
		var deploymentDetails []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.OpenflowDeployments.Describe(ctx, deployment.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			deploymentDetails = []map[string]any{schemas.OpenflowDeploymentDetailsToSchema(*describeResult)}
		}
		flattenedDeployments[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.OpenflowDeploymentToSchema(&deployment)},
			resources.DescribeOutputAttributeName: deploymentDetails,
		}
	}
	if err := d.Set("openflow_deployments", flattenedDeployments); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
