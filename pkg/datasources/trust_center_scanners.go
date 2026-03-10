package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
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
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the results to scanners with names matching the specified pattern.",
	},
	"scanners": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of scanners.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The scanner identifier.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The display name of the scanner.",
				},
				"short_description": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Short description of the scanner.",
				},
				"description": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Full description of the scanner.",
				},
				"scanner_package_id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The scanner package this scanner belongs to.",
				},
				"state": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The current state of the scanner ('TRUE' or 'FALSE').",
				},
				"schedule": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The current schedule for the scanner.",
				},
				"notification": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The notification configuration as JSON.",
				},
				"last_scan_timestamp": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Timestamp of the last scan run.",
				},
			},
		},
	},
}

func TrustCenterScanners() *schema.Resource {
	return &schema.Resource{
		ReadContext: ReadTrustCenterScanners,
		Schema:      trustCenterScannersSchema,
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

	scanners := make([]map[string]interface{}, len(scannersList))
	for i, s := range scannersList {
		scanner := map[string]interface{}{
			"id":                 s.Id,
			"name":               s.Name,
			"short_description":  s.ShortDescription,
			"description":        s.Description,
			"scanner_package_id": s.ScannerPackageId,
			"state":              s.State,
			"schedule":           s.Schedule,
			"notification":       s.Notification,
		}

		if s.LastScanTimestamp != nil {
			scanner["last_scan_timestamp"] = *s.LastScanTimestamp
		}

		scanners[i] = scanner
	}

	if err := d.Set("scanners", scanners); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("trust_center_scanners")

	return nil
}
