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

var externalAccessIntegrationsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC EXTERNAL ACCESS INTEGRATION for each integration returned by SHOW EXTERNAL ACCESS INTEGRATIONS. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like": likeSchema,
	"external_access_integrations": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all external access integration details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW EXTERNAL ACCESS INTEGRATIONS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowExternalAccessIntegrationSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE EXTERNAL ACCESS INTEGRATION.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeExternalAccessIntegrationDetailsSchema,
					},
				},
			},
		},
	},
}

func ExternalAccessIntegrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ExternalAccessIntegrationsDatasource), TrackingReadWrapper(datasources.ExternalAccessIntegrations, ReadExternalAccessIntegrations)),
		Schema:      externalAccessIntegrationsSchema,
		Description: "Data source used to get details of filtered external access integrations. Filtering is aligned with the current possibilities for [SHOW EXTERNAL ACCESS INTEGRATIONS](https://docs.snowflake.com/en/sql-reference/sql/show-integrations) query (only `like` is supported). The results of SHOW and DESCRIBE are encapsulated in one output collection `external_access_integrations`.",
	}
}

func ReadExternalAccessIntegrations(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowExternalAccessIntegrationRequest()

	handleLike(d, &req.Like)

	items, err := client.ExternalAccessIntegrations.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("external_access_integrations_read")

	flattened := make([]map[string]any, len(items))
	for i := range items {
		item := items[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.ExternalAccessIntegrations.DescribeDetails(ctx, item.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.ExternalAccessIntegrationDetailsToSchema(*details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.ExternalAccessIntegrationToSchema(&item)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("external_access_integrations", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
