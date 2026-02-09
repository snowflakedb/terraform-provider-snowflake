package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
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
	// Computed fields
	"state": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The current state of the scanner package ('TRUE' or 'FALSE').",
	},
	"provider_type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The provider of the scanner package (e.g., 'Snowflake').",
	},
	"last_enabled_timestamp": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Timestamp when the scanner package was last enabled.",
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

	// Set schedule if provided
	if v, ok := d.GetOk("schedule"); ok {
		schedule := v.(string)
		req.Schedule = &schedule
	}

	// Set notification if provided
	if v, ok := d.GetOk("notification"); ok {
		notificationList := v.([]interface{})
		if len(notificationList) > 0 {
			notificationMap := notificationList[0].(map[string]interface{})
			notification := &sdk.NotificationConfiguration{}

			if notifyAdmins, ok := notificationMap["notify_admins"].(bool); ok {
				notification.NotifyAdmins = &notifyAdmins
			}
			if severityThreshold, ok := notificationMap["severity_threshold"].(string); ok && severityThreshold != "" {
				notification.SeverityThreshold = &severityThreshold
			}
			if usersSet, ok := notificationMap["users"].(*schema.Set); ok && usersSet.Len() > 0 {
				users := make([]string, 0, usersSet.Len())
				for _, u := range usersSet.List() {
					users = append(users, u.(string))
				}
				notification.Users = users
			}
			req.Notification = notification
		}
	}

	err := client.TrustCenter.SetPackageConfiguration(ctx, req)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error creating Trust Center scanner package configuration",
				Detail:   fmt.Sprintf("error configuring scanner package %s: %v", scannerPackageId, err),
			},
		}
	}

	// Set the resource ID using source/package_id format
	d.SetId(helpers.EncodeResourceIdentifier(sdk.TrustCenterScannerPackageId{
		Source:           "SNOWFLAKE",
		ScannerPackageId: scannerPackageId,
	}.String()))

	return ReadContextTrustCenterScannerPackage(ctx, d, meta)
}

func ReadContextTrustCenterScannerPackage(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	// Parse the ID to get the scanner package ID
	id := d.Id()
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid resource ID format",
				Detail:   fmt.Sprintf("expected format 'SOURCE/SCANNER_PACKAGE_ID', got: %s", id),
			},
		}
	}
	scannerPackageId := parts[1]

	pkg, err := client.TrustCenter.ShowScannerPackageByID(ctx, scannerPackageId)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Scanner package not found",
					Detail:   fmt.Sprintf("Scanner package %s not found, removing from state", scannerPackageId),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Trust Center scanner package",
				Detail:   fmt.Sprintf("error reading scanner package %s: %v", scannerPackageId, err),
			},
		}
	}

	if err := d.Set("scanner_package_id", pkg.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", pkg.State == "TRUE"); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", pkg.State); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("provider_type", pkg.Provider); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", pkg.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("default_schedule", pkg.DefaultSchedule); err != nil {
		return diag.FromErr(err)
	}

	// Set schedule if it differs from default
	if pkg.Schedule != "" && pkg.Schedule != pkg.DefaultSchedule {
		if err := d.Set("schedule", pkg.Schedule); err != nil {
			return diag.FromErr(err)
		}
	}

	// Set last_enabled_timestamp if available
	if pkg.LastEnabledTimestamp != nil {
		if err := d.Set("last_enabled_timestamp", *pkg.LastEnabledTimestamp); err != nil {
			return diag.FromErr(err)
		}
	}

	// Parse and set notification configuration
	if pkg.Notification != "" {
		notification, err := sdk.ParseNotificationConfiguration(pkg.Notification)
		if err == nil && notification != nil {
			notificationList := []interface{}{
				map[string]interface{}{
					"notify_admins":      notification.NotifyAdmins != nil && *notification.NotifyAdmins,
					"severity_threshold": "",
					"users":              notification.Users,
				},
			}
			if notification.SeverityThreshold != nil {
				notificationList[0].(map[string]interface{})["severity_threshold"] = *notification.SeverityThreshold
			}
			if err := d.Set("notification", notificationList); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func UpdateContextTrustCenterScannerPackage(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	scannerPackageId := d.Get("scanner_package_id").(string)

	// Handle enabled changes
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

	// Handle schedule changes
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
			// Unset schedule
			req := &sdk.UnsetPackageConfigurationRequest{
				ScannerPackageId: scannerPackageId,
				UnsetSchedule:    true,
			}
			if err := client.TrustCenter.UnsetPackageConfiguration(ctx, req); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// Handle notification changes
	if d.HasChange("notification") {
		if v, ok := d.GetOk("notification"); ok {
			notificationList := v.([]interface{})
			if len(notificationList) > 0 {
				notificationMap := notificationList[0].(map[string]interface{})
				notification := &sdk.NotificationConfiguration{}

				if notifyAdmins, ok := notificationMap["notify_admins"].(bool); ok {
					notification.NotifyAdmins = &notifyAdmins
				}
				if severityThreshold, ok := notificationMap["severity_threshold"].(string); ok && severityThreshold != "" {
					notification.SeverityThreshold = &severityThreshold
				}
				if usersSet, ok := notificationMap["users"].(*schema.Set); ok && usersSet.Len() > 0 {
					users := make([]string, 0, usersSet.Len())
					for _, u := range usersSet.List() {
						users = append(users, u.(string))
					}
					notification.Users = users
				}

				req := &sdk.SetPackageConfigurationRequest{
					ScannerPackageId: scannerPackageId,
					Notification:     notification,
				}
				if err := client.TrustCenter.SetPackageConfiguration(ctx, req); err != nil {
					return diag.FromErr(err)
				}
			}
		} else {
			// Unset notification
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

	// Disable the scanner package (can't delete first-party packages, just disable)
	enabled := false
	req := &sdk.SetPackageConfigurationRequest{
		ScannerPackageId: scannerPackageId,
		Enabled:          &enabled,
	}

	if err := client.TrustCenter.SetPackageConfiguration(ctx, req); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error disabling Trust Center scanner package",
				Detail:   fmt.Sprintf("error disabling scanner package %s: %v", scannerPackageId, err),
			},
		}
	}

	// Also unset schedule and notification to return to defaults
	unsetReq := &sdk.UnsetPackageConfigurationRequest{
		ScannerPackageId:  scannerPackageId,
		UnsetSchedule:     true,
		UnsetNotification: true,
	}
	// Ignore errors on unset - package is already disabled
	_ = client.TrustCenter.UnsetPackageConfiguration(ctx, unsetReq)

	d.SetId("")
	return nil
}

func ImportTrustCenterScannerPackage(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Import ID format: SNOWFLAKE/SECURITY_ESSENTIALS
	id := d.Id()
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import ID format, expected 'SOURCE/SCANNER_PACKAGE_ID', got: %s", id)
	}

	if err := d.Set("scanner_package_id", parts[1]); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
