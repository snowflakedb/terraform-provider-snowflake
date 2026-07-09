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

var apiIntegrationsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC API INTEGRATION for each integration returned by SHOW API INTEGRATIONS. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like": likeSchema,
	"api_integrations": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all API integration details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW API INTEGRATIONS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowApiIntegrationSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE API INTEGRATION.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeApiIntegrationAllDetailsSchema,
					},
				},
			},
		},
	},
}

func ApiIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ApiIntegrationsDatasource), TrackingReadWrapper(datasources.ApiIntegrations, ReadApiIntegrations)),
		Schema:      apiIntegrationsSchema,
		Description: "Data source used to get details of filtered API integrations. Filtering is aligned with the current possibilities for [SHOW API INTEGRATIONS](https://docs.snowflake.com/en/sql-reference/sql/show-integrations) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `api_integrations`.",
	}
}

func ReadApiIntegrations(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	showRequest := sdk.NewShowApiIntegrationRequest()

	handleLike(d, &showRequest.Like)

	apiIntegrations, err := client.ApiIntegrations.Show(ctx, showRequest)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("api_integrations_read")

	flattened := make([]map[string]any, len(apiIntegrations))
	for i := range apiIntegrations {
		ai := apiIntegrations[i]
		var describeOut []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.ApiIntegrations.DescribeAllDetails(ctx, ai.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOut = []map[string]any{schemas.ApiIntegrationAllDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.ApiIntegrationToSchema(&ai)},
			resources.DescribeOutputAttributeName: describeOut,
		}
	}
	if err := d.Set("api_integrations", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
