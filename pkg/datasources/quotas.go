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

var quotasSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Reads the configuration of each quota returned by SHOW SNOWFLAKE.CORE.QUOTA with its GET_* methods (scope, config, notifications, and actions). The output is saved to the describe_output field. By default this value is set to true.",
	},
	"like":  likeSchema,
	"in":    inSchema,
	"limit": limitFromSchema,
	"quotas": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all quota details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW SNOWFLAKE.CORE.QUOTA.",
					Elem: &schema.Resource{
						Schema: schemas.ShowQuotaSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the quota configuration read with its GET_* methods.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeQuotaDetailsSchema,
					},
				},
			},
		},
	},
}

func Quotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.QuotasDatasource), TrackingReadWrapper(datasources.Quotas, ReadQuotas)),
		Schema:      quotasSchema,
		Description: "Data source used to get details of filtered per-user quotas (`SNOWFLAKE.CORE.QUOTA`). Filtering is aligned with the current possibilities for SHOW SNOWFLAKE.CORE.QUOTA." +
			" The results of SHOW and of the quota's GET_* methods are encapsulated in one output collection `quotas`.",
	}
}

func ReadQuotas(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.NewShowQuotaRequest()

	handleLike(d, &req.Like)
	handleLimitFrom(d, &req.Limit)
	if err := handleIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}

	quotas, err := client.Quotas.Show(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("quotas_read")

	flattened := make([]map[string]any, len(quotas))
	for i := range quotas {
		quota := quotas[i]
		var describeOutput []map[string]any
		if d.Get("with_describe").(bool) {
			details, err := client.Quotas.DescribeDetails(ctx, quota.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			describeOutput = []map[string]any{schemas.QuotaDetailsToSchema(details)}
		}
		flattened[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.QuotaToSchema(&quota)},
			resources.DescribeOutputAttributeName: describeOutput,
		}
	}
	if err := d.Set("quotas", flattened); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
