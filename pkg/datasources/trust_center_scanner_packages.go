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

var trustCenterScannerPackagesSchema = map[string]*schema.Schema{
	"like": likeSchema,
	"scanner_packages": {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Holds the aggregated output of all scanner package queries.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				resources.ShowOutputAttributeName: {
					Type:        schema.TypeList,
					Computed:    true,
					Description: "Holds the output of SHOW TRUST CENTER SCANNER PACKAGES.",
					Elem: &schema.Resource{
						Schema: schemas.ShowTrustCenterScannerPackageSchema,
					},
				},
			},
		},
	},
}

func TrustCenterScannerPackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: TrackingReadWrapper(datasources.TrustCenterScannerPackages, ReadTrustCenterScannerPackages),
		Schema:     trustCenterScannerPackagesSchema,
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

	d.SetId("trust_center_scanner_packages")

	flattenedPackages := make([]map[string]any, len(packages))
	for i, pkg := range packages {
		pkg := pkg
		flattenedPackages[i] = map[string]any{
			resources.ShowOutputAttributeName: []map[string]any{schemas.ScannerPackageToSchema(&pkg)},
		}
	}

	if err := d.Set("scanner_packages", flattenedPackages); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
