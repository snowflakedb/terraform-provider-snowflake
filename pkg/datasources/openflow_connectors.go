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

var openflowConnectorsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC OPENFLOW CONNECTOR for each connector returned by SHOW OPENFLOW CONNECTORS. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          inSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"openflow_connectors": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all Openflow connector details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW OPENFLOW CONNECTORS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowOpenflowConnectorSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE OPENFLOW CONNECTOR.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeOpenflowConnectorDetailsSchema,
					},
				},
			},
		},
	},
}

func OpenflowConnectors() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.OpenflowConnectorsDatasource), TrackingReadWrapper(datasources.OpenflowConnectors, ReadOpenflowConnectors)),
		Schema:      openflowConnectorsSchema,
		Description: "Data source used to get details of filtered Openflow connectors. Filtering is aligned with the current possibilities for [SHOW OPENFLOW CONNECTORS](https://docs.snowflake.com/en/sql-reference/sql/show-openflow-connectors). The results of SHOW and DESCRIBE are encapsulated in one output collection `openflow_connectors`.",
	}
}

func ReadOpenflowConnectors(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowOpenflowConnectorRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	connectors, err := client.OpenflowConnectors.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("openflow_connectors_read")

	flattenedConnectors := make([]map[string]any, len(connectors))
	for i, connector := range connectors {
		var connectorDetails []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.OpenflowConnectors.Describe(ctx, connector.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			connectorDetails = []map[string]any{schemas.OpenflowConnectorDetailsToSchema(describeResult)}
		}
		flattenedConnectors[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.OpenflowConnectorToSchema(&connector)},
			resources.DescribeOutputAttributeName: connectorDetails,
		}
	}
	if err := d.Set("openflow_connectors", flattenedConnectors); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
