package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceMonitorsSchema = map[string]*schema.Schema{
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the output with **case-insensitive** pattern, with support for SQL wildcard characters (`%` and `_`).",
	},
	"resource_monitors": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all resource monitor details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW RESOURCE MONITORS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowResourceMonitorSchema,
					},
				},
			},
		},
	},
}

func ResourceMonitors() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.ResourceMonitors, ReadResourceMonitors),
		Schema:      resourceMonitorsSchema,
		Description: "Data source used to get details of filtered resource monitors. Filtering is aligned with the current possibilities for [SHOW RESOURCE MONITORS](https://docs.snowflake.com/en/sql-reference/sql/show-resource-monitors) query (`like` is supported). The results of SHOW is encapsulated in show_output collection.",
	}
}

func ReadResourceMonitors(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	req := sdk.NewShowResourceMonitorRequest()

	if likePattern, ok := d.GetOk("like"); ok {
		req.WithLike(sdk.Like{
			Pattern: sdk.String(likePattern.(string)),
		})
	}

	resourceMonitors, err := client.ResourceMonitors.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("resource_monitors_read")

	flattenedResourceMonitors := make([]map[string]any, len(resourceMonitors))
	for i, resourceMonitor := range resourceMonitors {
		flattenedResourceMonitors[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.ResourceMonitorToSchema(&resourceMonitor)},
		}
	}

	err = d.Set("resource_monitors", flattenedResourceMonitors)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
