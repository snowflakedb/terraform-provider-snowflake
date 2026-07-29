package datasources

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider/validators"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ShowServicesType string

const (
	ShowServicesTypeAll          ShowServicesType = "ALL"
	ShowServicesTypeJobsOnly     ShowServicesType = "JOBS_ONLY"
	ShowServicesTypeServicesOnly ShowServicesType = "SERVICES_ONLY"
)

var allShowServicesTypes = []ShowServicesType{
	ShowServicesTypeAll,
	ShowServicesTypeJobsOnly,
	ShowServicesTypeServicesOnly,
}

func toShowServicesType(s string) (ShowServicesType, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(allShowServicesTypes, ShowServicesType(s)) {
		return "", fmt.Errorf("invalid show services type: %s", s)
	}
	return ShowServicesType(s), nil
}

var servicesSchema = map[string]*schema.Schema{
	"with_describe": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     true,
		Description: "Runs DESC SERVICE for each service returned by SHOW SERVICES. The output of describe is saved to the description field. By default this value is set to true.",
	},
	"service_type": {
		Type:             schema.TypeString,
		Optional:         true,
		Default:          string(ShowServicesTypeAll),
		Description:      fmt.Sprintf("The type filtering of `SHOW SERVICES` results. `ALL` returns both services and job services. `JOBS_ONLY` returns only job services (`JOB` option in SQL). `SERVICES_ONLY` returns only services (`EXCLUDE_JOBS` option in SQL)."),
		ValidateDiagFunc: validators.NormalizeValidation(toShowServicesType),
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
		ReadContext: TrackingReadWrapper(datasources.Services, ReadServices),
		Schema:      servicesSchema,
		Description: "Data source used to get details of filtered services. Filtering is aligned with the current possibilities for [SHOW SERVICES](https://docs.snowflake.com/en/sql-reference/sql/show-services) query." +
			" The results of SHOW and DESCRIBE are encapsulated in one output collection `services`. By default, the results includes both services and job services. If you want to filter only services or job service, set `service_type` with a relevant option.",
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
	if v := d.Get("service_type").(string); v != "" {
		mode, err := toShowServicesType(v)
		if err != nil {
			return diag.FromErr(err)
		}
		switch mode {
		case ShowServicesTypeJobsOnly:
			req.Job = sdk.Bool(true)
		case ShowServicesTypeServicesOnly:
			req.ExcludeJobs = sdk.Bool(true)
		case ShowServicesTypeAll:
			// no-op, filtering is not applied
		}
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
