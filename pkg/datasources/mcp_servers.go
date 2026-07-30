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

var mcpServersSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC MCP SERVER for each MCP server returned by SHOW MCP SERVERS. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like": likeSchema,
	"in":   inSchema,
	"mcp_servers": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all MCP server details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW MCP SERVERS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowMcpServerSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE MCP SERVER.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeMcpServerDetailsSchema,
					},
				},
			},
		},
	},
}

func McpServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.McpServersDatasource), TrackingReadWrapper(datasources.McpServers, ReadMcpServers)),
		Schema:      mcpServersSchema,
		Description: "Data source used to get details of filtered MCP servers. Filtering is aligned with the current possibilities for [SHOW MCP SERVERS](https://docs.snowflake.com/en/sql-reference/sql/show-mcp-servers) query (`like`, `in`). The results of SHOW and DESCRIBE are encapsulated in one output collection `mcp_servers`.",
	}
}

func ReadMcpServers(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowMcpServerRequest()

	handleLike(d, &req.Like)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	mcpServers, err := client.McpServers.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("mcp_servers_read")

	flattened := make([]map[string]any, len(mcpServers))
	for i := range mcpServers {
		server := mcpServers[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.McpServers.Describe(ctx, server.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.McpServerDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.McpServerToSchema(&server)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("mcp_servers", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
