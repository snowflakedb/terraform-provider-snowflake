package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var servicesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC SERVICE for each service returned by SHOW SERVICES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"jobs_only": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, only jobs will be returned. If false, normal services will be included in the output.",
	},
	"exclude_jobs": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If true, jobs will be excluded from the output. If false, jobs will be included in the output.",
	},
	"like":        likeSchema,
	"starts_with": startsWithSchema,
	"limit":       limitFromSchema,
	"in":          serviceInSchema,
	"services": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all services details queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW SERVICES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowServiceSchema,
					},
				},
				resources.DescribeOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of DESCRIBE SERVICE.",
					Elem: &schema.Resource{
						Schema: schemas.DescribeServiceSchema,
					},
				},
			},
		},
	},
}

func Services() *schema.Resource {
	return &schema.Resource{
		ReadContext: PreviewFeatureReadWrapper(string(previewfeatures.ServicesDatasource), TrackingReadWrapper(datasources.Services, ReadServices)),
		Schema:      servicesSchema,
		Description: "Data source used to get details of filtered services. Filtering is aligned with the current possibilities for [SHOW SERVICES](https://docs.snowflake.com/en/sql-reference/sql/show-services) query. The results of SHOW and DESCRIBE are encapsulated in one output collection `services`.",
	}
}

func ReadServices(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*provider.Context).Client
	req := sdk.ShowServiceRequest{}

	handleLike(d, &req.Like)
	handleStartsWith(d, &req.StartsWith)
	handleLimitFrom(d, &req.Limit)
	if err := handleServiceIn(d, &req.In); err != nil {
		return diag.FromErr(err)
	}
	if d.Get("jobs_only").(bool) {
		req.Job = sdk.Bool(true)
	}
	if d.Get("exclude_jobs").(bool) {
		req.ExcludeJobs = sdk.Bool(true)
	}

	services, err := client.Services.Show(ctx, &req)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("services_read")

	flattenedServices := make([]map[string]any, len(services))
	for i, service := range services {
		service := service
		var serviceDetails []map[string]any
		if d.Get("with_describe").(bool) {
			describeResult, err := client.Services.Describe(ctx, service.ID())
			if err != nil {
				return diag.FromErr(err)
			}
			serviceDetails = []map[string]any{schemas.ServiceDetailsToSchema(describeResult)}
		}
		flattenedServices[i] = map[string]any{
			resources.ShowOutputAttributeName:     []map[string]any{schemas.ServiceToSchema(&service)},
			resources.DescribeOutputAttributeName: serviceDetails,
		}
	}
	if err := d.Set("services", flattenedServices); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
