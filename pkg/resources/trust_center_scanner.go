package resources

import (
	"context"
	"errors"
	"fmt"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var trustCenterScannerSchema = map[string]*schema.Schema{
	"scanner_package_id": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The scanner package identifier (e.g., 'SECURITY_ESSENTIALS').",
	},
	"scanner_id": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The scanner identifier (e.g., 'SECURITY_ESSENTIALS_MFA_REQUIRED_FOR_USERS_CHECK').",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Whether the scanner is enabled.",
	},
	"schedule": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "CRON expression for the scan schedule (e.g., 'USING CRON 0 2 * * * UTC').",
	},
	"notification": {
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Notification configuration for the scanner.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"notify_admins": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Whether to notify administrators.",
				},
				"severity_threshold": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Minimum severity level for notifications (e.g., 'High', 'Medium', 'Low', 'Critical').",
				},
				"users": {
					Type:        schema.TypeSet,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "List of users to notify.",
				},
			},
		},
	},
	ShowOutputAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW TRUST CENTER SCANNERS` for the given scanner.",
		Elem: &schema.Resource{
			Schema: schemas.ShowTrustCenterScannerSchema,
		},
	},
}

func TrustCenterScanner() *schema.Resource {
	return &schema.Resource{
		Schema: trustCenterScannerSchema,

		CreateContext: TrackingCreateWrapper(resources.TrustCenterScanner, CreateContextTrustCenterScanner),
		ReadContext:   TrackingReadWrapper(resources.TrustCenterScanner, ReadContextTrustCenterScanner),
		UpdateContext: TrackingUpdateWrapper(resources.TrustCenterScanner, UpdateContextTrustCenterScanner),
		DeleteContext: TrackingDeleteWrapper(resources.TrustCenterScanner, DeleteContextTrustCenterScanner),
		Description:   "Resource for managing individual Trust Center scanner configuration. For more information, see [Trust Center](https://docs.snowflake.com/en/user-guide/trust-center).",

		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.TrustCenterScanner, ImportTrustCenterScanner),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextTrustCenterScanner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	scannerPackageId := d.Get("scanner_package_id").(string)
	scannerId := d.Get("scanner_id").(string)
	enabled := d.Get("enabled").(bool)

	req := &sdk.SetScannerConfigurationRequest{
		ScannerPackageId: scannerPackageId,
		ScannerId:        scannerId,
		Enabled:          &enabled,
	}

	if v, ok := d.GetOk("schedule"); ok {
		schedule := v.(string)
		req.Schedule = &schedule
	}

	if v, ok := d.GetOk("notification"); ok {
		req.Notification = expandNotificationConfiguration(v)
	}

	err := client.TrustCenter.SetScannerConfiguration(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error configuring scanner %s/%s: %w", scannerPackageId, scannerId, err))
	}

	d.SetId(sdk.TrustCenterScannerId{
		Source:           "SNOWFLAKE",
		ScannerPackageId: scannerPackageId,
		ScannerId:        scannerId,
	}.String())

	return ReadContextTrustCenterScanner(ctx, d, meta)
}

func ReadContextTrustCenterScanner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	parsedId, err := sdk.ParseTrustCenterScannerId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	scanner, err := client.TrustCenter.ShowScannerByID(ctx, parsedId.ScannerPackageId, parsedId.ScannerId)
	if err != nil {
		if errors.Is(err, collections.ErrObjectNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading scanner %s/%s: %w", parsedId.ScannerPackageId, parsedId.ScannerId, err))
	}

	if err := d.Set("scanner_package_id", scanner.ScannerPackageId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("scanner_id", scanner.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", scanner.State == "TRUE"); err != nil {
		return diag.FromErr(err)
	}

	if scanner.Schedule != "" {
		if err := d.Set("schedule", scanner.Schedule); err != nil {
			return diag.FromErr(err)
		}
	}

	notificationList, err := flattenNotificationConfiguration(scanner.Notification)
	if err == nil && notificationList != nil {
		if err := d.Set("notification", notificationList); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set(ShowOutputAttributeName, []map[string]any{schemas.ScannerToSchema(scanner)}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContextTrustCenterScanner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	scannerPackageId := d.Get("scanner_package_id").(string)
	scannerId := d.Get("scanner_id").(string)

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		req := &sdk.SetScannerConfigurationRequest{
			ScannerPackageId: scannerPackageId,
			ScannerId:        scannerId,
			Enabled:          &enabled,
		}
		if err := client.TrustCenter.SetScannerConfiguration(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("schedule") {
		if v, ok := d.GetOk("schedule"); ok {
			schedule := v.(string)
			req := &sdk.SetScannerConfigurationRequest{
				ScannerPackageId: scannerPackageId,
				ScannerId:        scannerId,
				Schedule:         &schedule,
			}
			if err := client.TrustCenter.SetScannerConfiguration(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		} else {
			req := &sdk.UnsetScannerConfigurationRequest{
				ScannerPackageId: scannerPackageId,
				ScannerId:        scannerId,
				UnsetSchedule:    true,
			}
			if err := client.TrustCenter.UnsetScannerConfiguration(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("notification") {
		if v, ok := d.GetOk("notification"); ok {
			notification := expandNotificationConfiguration(v)
			if notification != nil {
				req := &sdk.SetScannerConfigurationRequest{
					ScannerPackageId: scannerPackageId,
					ScannerId:        scannerId,
					Notification:     notification,
				}
				if err := client.TrustCenter.SetScannerConfiguration(ctx, req); err != nil {
					return diag.FromErr(err)
				}
			}
		} else {
			req := &sdk.UnsetScannerConfigurationRequest{
				ScannerPackageId:  scannerPackageId,
				ScannerId:         scannerId,
				UnsetNotification: true,
			}
			if err := client.TrustCenter.UnsetScannerConfiguration(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadContextTrustCenterScanner(ctx, d, meta)
}

func DeleteContextTrustCenterScanner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	scannerPackageId := d.Get("scanner_package_id").(string)
	scannerId := d.Get("scanner_id").(string)

	enabled := false
	req := &sdk.SetScannerConfigurationRequest{
		ScannerPackageId: scannerPackageId,
		ScannerId:        scannerId,
		Enabled:          &enabled,
	}

	if err := client.TrustCenter.SetScannerConfiguration(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error disabling scanner %s/%s: %w", scannerPackageId, scannerId, err))
	}

	unsetReq := &sdk.UnsetScannerConfigurationRequest{
		ScannerPackageId:  scannerPackageId,
		ScannerId:         scannerId,
		UnsetSchedule:     true,
		UnsetNotification: true,
	}
	_ = client.TrustCenter.UnsetScannerConfiguration(ctx, unsetReq)

	d.SetId("")
	return nil
}

func ImportTrustCenterScanner(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parsedId, err := sdk.ParseTrustCenterScannerId(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("scanner_package_id", parsedId.ScannerPackageId); err != nil {
		return nil, err
	}
	if err := d.Set("scanner_id", parsedId.ScannerId); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
