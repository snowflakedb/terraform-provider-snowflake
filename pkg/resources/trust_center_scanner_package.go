package resources

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var trustCenterScannerPackageSchema = map[string]*schema.Schema{
	"scanner_package_id": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "The scanner package identifier (e.g., 'SECURITY_ESSENTIALS').",
	},
	"enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Whether the scanner package is enabled.",
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
		Description: "Notification configuration for the scanner package.",
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
		Description: "Outputs the result of `SHOW TRUST CENTER SCANNER PACKAGES` for the given scanner package.",
		Elem: &schema.Resource{
			Schema: schemas.ShowTrustCenterScannerPackageSchema,
		},
	},
}

func TrustCenterScannerPackage() *schema.Resource {
	return &schema.Resource{
		Schema: trustCenterScannerPackageSchema,

		CreateContext: TrackingCreateWrapper(resources.TrustCenterScannerPackage, CreateContextTrustCenterScannerPackage),
		ReadContext:   TrackingReadWrapper(resources.TrustCenterScannerPackage, ReadContextTrustCenterScannerPackage),
		UpdateContext: TrackingUpdateWrapper(resources.TrustCenterScannerPackage, UpdateContextTrustCenterScannerPackage),
		DeleteContext: TrackingDeleteWrapper(resources.TrustCenterScannerPackage, DeleteContextTrustCenterScannerPackage),
		Description:   "Resource for managing Trust Center scanner package configuration. For more information, see [Trust Center](https://docs.snowflake.com/en/user-guide/trust-center).",

		Importer: &schema.ResourceImporter{
			StateContext: TrackingImportWrapper(resources.TrustCenterScannerPackage, ImportTrustCenterScannerPackage),
		},
		Timeouts: defaultTimeouts,
	}
}

func CreateContextTrustCenterScannerPackage(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	scannerPackageId := d.Get("scanner_package_id").(string)
	enabled := d.Get("enabled").(bool)

	req := &sdk.SetPackageConfigurationRequest{
		ScannerPackageId: scannerPackageId,
		Enabled:          &enabled,
	}

	if v, ok := d.GetOk("schedule"); ok {
		schedule := v.(string)
		req.Schedule = &schedule
	}

	if v, ok := d.GetOk("notification"); ok {
		req.Notification = expandNotificationConfiguration(v)
	}

	err := client.TrustCenter.SetPackageConfiguration(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error configuring scanner package %s: %w", scannerPackageId, err))
	}

	d.SetId(sdk.TrustCenterScannerPackageId{
		Source:           "SNOWFLAKE",
		ScannerPackageId: scannerPackageId,
	}.String())

	return ReadContextTrustCenterScannerPackage(ctx, d, meta)
}

func ReadContextTrustCenterScannerPackage(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	parsedId, err := sdk.ParseTrustCenterScannerPackageId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	pkg, err := client.TrustCenter.ShowScannerPackageByID(ctx, parsedId.ScannerPackageId)
	if err != nil {
		if errors.Is(err, collections.ErrObjectNotFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading scanner package %s: %w", parsedId.ScannerPackageId, err))
	}

	if err := d.Set("scanner_package_id", pkg.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", pkg.State == "TRUE"); err != nil {
		return diag.FromErr(err)
	}

	if pkg.Schedule != "" {
		if err := d.Set("schedule", pkg.Schedule); err != nil {
			return diag.FromErr(err)
		}
	}

	notificationList, err := flattenNotificationConfiguration(pkg.Notification)
	if err == nil && notificationList != nil {
		if err := d.Set("notification", notificationList); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set(ShowOutputAttributeName, []map[string]any{schemas.ScannerPackageToSchema(pkg)}); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func UpdateContextTrustCenterScannerPackage(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	scannerPackageId := d.Get("scanner_package_id").(string)

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		req := &sdk.SetPackageConfigurationRequest{
			ScannerPackageId: scannerPackageId,
			Enabled:          &enabled,
		}
		if err := client.TrustCenter.SetPackageConfiguration(ctx, req); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("schedule") {
		if v, ok := d.GetOk("schedule"); ok {
			schedule := v.(string)
			req := &sdk.SetPackageConfigurationRequest{
				ScannerPackageId: scannerPackageId,
				Schedule:         &schedule,
			}
			if err := client.TrustCenter.SetPackageConfiguration(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		} else {
			req := &sdk.UnsetPackageConfigurationRequest{
				ScannerPackageId: scannerPackageId,
				UnsetSchedule:    true,
			}
			if err := client.TrustCenter.UnsetPackageConfiguration(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("notification") {
		if v, ok := d.GetOk("notification"); ok {
			notification := expandNotificationConfiguration(v)
			if notification != nil {
				req := &sdk.SetPackageConfigurationRequest{
					ScannerPackageId: scannerPackageId,
					Notification:     notification,
				}
				if err := client.TrustCenter.SetPackageConfiguration(ctx, req); err != nil {
					return diag.FromErr(err)
				}
			}
		} else {
			req := &sdk.UnsetPackageConfigurationRequest{
				ScannerPackageId:  scannerPackageId,
				UnsetNotification: true,
			}
			if err := client.TrustCenter.UnsetPackageConfiguration(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return ReadContextTrustCenterScannerPackage(ctx, d, meta)
}

func DeleteContextTrustCenterScannerPackage(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	scannerPackageId := d.Get("scanner_package_id").(string)

	enabled := false
	req := &sdk.SetPackageConfigurationRequest{
		ScannerPackageId: scannerPackageId,
		Enabled:          &enabled,
	}

	if err := client.TrustCenter.SetPackageConfiguration(ctx, req); err != nil {
		return diag.FromErr(fmt.Errorf("error disabling scanner package %s: %w", scannerPackageId, err))
	}

	unsetReq := &sdk.UnsetPackageConfigurationRequest{
		ScannerPackageId:  scannerPackageId,
		UnsetSchedule:     true,
		UnsetNotification: true,
	}
	if err := client.TrustCenter.UnsetPackageConfiguration(ctx, unsetReq); err != nil {
		log.Printf("[DEBUG] failed to unset schedule/notification for scanner package %s during delete: %s", scannerPackageId, err)
	}

	d.SetId("")
	return nil
}

func ImportTrustCenterScannerPackage(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parsedId, err := sdk.ParseTrustCenterScannerPackageId(d.Id())
	if err != nil {
		return nil, err
	}

	if err := d.Set("scanner_package_id", parsedId.ScannerPackageId); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
