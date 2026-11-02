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

var openflowConnectorDefinitionsSchema = map[string]*schema.Schema{
	// SHOW OPENFLOW CONNECTOR DEFINITIONS accepts STARTS WITH but returns every definition regardless
	// (SNOW-4024255), so only the filters that narrow the result are exposed. LIKE takes patterns, which
	// covers prefix matching.
	"like":  likeSchema,
	"limit": limitFromSchema,
	"openflow_connector_definitions": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all Openflow connector definition queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW OPENFLOW CONNECTOR DEFINITIONS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowOpenflowConnectorDefinitionSchema,
					},
				},
			},
		},
	},
}

func OpenflowConnectorDefinitions() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.OpenflowConnectorDefinitionsDatasource), TrackingReadWrapper(datasources.OpenflowConnectorDefinitions, ReadOpenflowConnectorDefinitions)),
		Schema:      openflowConnectorDefinitionsSchema,
		Description: docs.WithPreviewFeatureMark() + "Data source used to get details of filtered Openflow connector definitions, the Snowflake-managed templates a connector can be created from. Filtering is aligned with the current possibilities for [SHOW OPENFLOW CONNECTOR DEFINITIONS](https://docs.snowflake.com/en/sql-reference/sql/show-openflow-connector-definitions). Definitions are read-only, so there is no describe output and no matching resource.",
	}
}

func ReadOpenflowConnectorDefinitions(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowOpenflowConnectorDefinitionRequest{}

	handleLike(d, &req.Like)
	handleLimitFrom(d, &req.Limit)

	definitions, err := client.OpenflowConnectorDefinitions.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("openflow_connector_definitions_read")

	flattenedDefinitions := make([]map[string]any, len(definitions))
	for i, definition := range definitions {
		flattenedDefinitions[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.OpenflowConnectorDefinitionToSchema(&definition)},
		}
	}
	if err := d.Set("openflow_connector_definitions", flattenedDefinitions); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
