package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var trustCenterScannersSchema = map[string]*schema.Schema{
	"scanner_package_id": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the results to scanners belonging to the specified scanner package.",
	},
	"like": likeSchema,
	"scanners": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all scanner queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW TRUST CENTER SCANNERS.",
					Elem: &schema.Resource{
						Schema: schemas.ShowTrustCenterScannerSchema,
					},
				},
			},
		},
	},
}

func TrustCenterScanners() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.TrustCenterScanners, ReadTrustCenterScanners),
		Schema:     trustCenterScannersSchema,
		Description: "Data source for listing Trust Center scanners. For more information, see [Trust Center](https://docs.snowflake.com/en/user-guide/trust-center).",
	}
}

func ReadTrustCenterScanners(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	req := &sdk.ShowScannersRequest{}

	if v, ok := d.GetOk("scanner_package_id"); ok {
		scannerPackageId := v.(string)
		req.ScannerPackageId = &scannerPackageId
	}

	if v, ok := d.GetOk("like"); ok {
		like := v.(string)
		req.Like = &like
	}

	scannersList, err := client.TrustCenter.ShowScanners(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("trust_center_scanners")

	flattenedScanners := make([]map[string]any, len(scannersList))
	for i, s := range scannersList {
		s := s
		flattenedScanners[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.ScannerToSchema(&s)},
		}
	}

	if err := d.Set("scanners", flattenedScanners); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
