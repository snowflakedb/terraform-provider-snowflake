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

var alertsSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC ALERT for each alert returned by SHOW ALERTS. The output of describe is saved to the describe_output field. By default this value is set to true.",
	},
	"like":        likeSchema,
	"in":          extendedInSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"alerts": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all alert details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW ALERTS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowAlertSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE ALERT.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeAlertDetailsSchema,
					},
				},
			},
		},
	},
}

func Alerts() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.AlertsDatasource), TrackingReadWrapper(datasources.Alerts, ReadAlerts)),
		Schema:      alertsSchema,
		Description: docs.WithPreviewFeatureMark() + "Data source used to get details of filtered alerts. Filtering is aligned with the current possibilities for [SHOW ALERTS](https://docs.snowflake.com/en/sql-reference/sql/show-alerts) query. " +
			"The results of SHOW and DESCRIBE are encapsulated in one output collection `alerts`.",
	}
}

func ReadAlerts(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowAlertRequest{}

	handleLike(d, &req.Like)
	if err := handleExtendedIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)

	alerts, err := client.Alerts.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("alerts_read")

	flattened := make([]map[string]any, len(alerts))
	for i := range alerts {
		alert := alerts[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.Alerts.Describe(ctx, alert.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.AlertDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.AlertToSchema(&alert)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("alerts", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
