package datasources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var trustCenterScannerPackagesSchema = map[string]*schema.Schema{
	"like": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Filters the results to scanner packages with names matching the specified pattern.",
	},
	"scanner_packages": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "List of scanner packages.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The scanner package identifier.",
				},
				"name": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The display name of the scanner package.",
				},
				"description": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Description of the scanner package.",
				},
				"default_schedule": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The default schedule for the scanner package.",
				},
				"state": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The current state of the scanner package ('TRUE' or 'FALSE').",
				},
				"schedule": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The current schedule for the scanner package.",
				},
				"notification": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The notification configuration as JSON.",
				},
				"provider": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "The provider of the scanner package (e.g., 'Snowflake').",
				},
				"last_enabled_timestamp": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Timestamp when the scanner package was last enabled.",
				},
				"last_disabled_timestamp": {
					Type:        schema.TypeString,
					Computed:    true,
					Description: "Timestamp when the scanner package was last disabled.",
				},
			},
		},
	},
}

func TrustCenterScannerPackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: ReadTrustCenterScannerPackages,
		Schema:      trustCenterScannerPackagesSchema,
		Description: "Data source for listing Trust Center scanner packages. For more information, see [Trust Center](https://docs.snowflake.com/en/user-guide/trust-center).",
	}
}

func ReadTrustCenterScannerPackages(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	req := &sdk.ShowScannerPackagesRequest{}

	if v, ok := d.GetOk("like"); ok {
		like := v.(string)
		req.Like = &like
	}

	packages, err := client.TrustCenter.ShowScannerPackages(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	scannerPackages := make([]map[string]interface{}, len(packages))
	for i, pkg := range packages {
		scannerPackage := map[string]interface{}{
			"id":               pkg.Id,
			"name":             pkg.Name,
			"description":      pkg.Description,
			"default_schedule": pkg.DefaultSchedule,
			"state":            pkg.State,
			"schedule":         pkg.Schedule,
			"notification":     pkg.Notification,
			"provider":         pkg.Provider,
		}

		if pkg.LastEnabledTimestamp != nil {
			scannerPackage["last_enabled_timestamp"] = *pkg.LastEnabledTimestamp
		}
		if pkg.LastDisabledTimestamp != nil {
			scannerPackage["last_disabled_timestamp"] = *pkg.LastDisabledTimestamp
		}

		scannerPackages[i] = scannerPackage
	}

	if err := d.Set("scanner_packages", scannerPackages); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("trust_center_scanner_packages")

	return nil
}
