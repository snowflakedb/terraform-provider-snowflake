package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/datasources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var connectionsSchema = map[string]*schema.Schema{
	"like": likeSchema,
	"connections": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all connections details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW CONNECTIONS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowConnectionSchema,
					},
				},
			},
		},
	},
}

func Connections() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.Connections, ReadConnections),
		Schema:      connectionsSchema,
		Description: "Data source used to get details of filtered connections. Filtering is aligned with the current possibilities for [SHOW CONNECTIONS](https://docs.snowflake.com/en/sql-reference/sql/show-connections) query. The results of SHOW is encapsulated in one output collection `connections`.",
	}
}

func ReadConnections(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowConnectionRequest{}

	handleLike(d, &req.Like)

	connections, err := client.Connections.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("connections_read")

	flattenedConnections := make([]map[string]any, len(connections))
	for i, connection := range connections {
		connection := connection
		flattenedConnections[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.ConnectionToSchema(&connection)},
		}
	}
	if err := d.Set("connections", flattenedConnections); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
