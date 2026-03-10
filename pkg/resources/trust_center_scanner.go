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
	// Computed fields
	"state": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The current state of the scanner ('TRUE' or 'FALSE').",
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
	"last_scan_timestamp": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Timestamp of the last scan run.",
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

	err := client.TrustCenter.SetScannerConfiguration(ctx, req)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error creating Trust Center scanner configuration",
				Detail:   fmt.Sprintf("error configuring scanner %s/%s: %v", scannerPackageId, scannerId, err),
			},
		}
	}

	// Set the resource ID using source/package_id/scanner_id format
	d.SetId(helpers.EncodeResourceIdentifier(sdk.TrustCenterScannerId{
		Source:           "SNOWFLAKE",
		ScannerPackageId: scannerPackageId,
		ScannerId:        scannerId,
	}.String()))

	return ReadContextTrustCenterScanner(ctx, d, meta)
}

func ReadContextTrustCenterScanner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	// Parse the ID to get the scanner package ID and scanner ID
	id := d.Id()
	parts := strings.Split(id, "/")
	if len(parts) != 3 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid resource ID format",
				Detail:   fmt.Sprintf("expected format 'SOURCE/SCANNER_PACKAGE_ID/SCANNER_ID', got: %s", id),
			},
		}
	}
	scannerPackageId := parts[1]
	scannerId := parts[2]

	scanner, err := client.TrustCenter.ShowScannerByID(ctx, scannerPackageId, scannerId)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Scanner not found",
					Detail:   fmt.Sprintf("Scanner %s/%s not found, removing from state", scannerPackageId, scannerId),
				},
			}
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error reading Trust Center scanner",
				Detail:   fmt.Sprintf("error reading scanner %s/%s: %v", scannerPackageId, scannerId, err),
			},
		}
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
	if err := d.Set("state", scanner.State); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", scanner.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("short_description", scanner.ShortDescription); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", scanner.Description); err != nil {
		return diag.FromErr(err)
	}

	// Set schedule if available
	if scanner.Schedule != "" {
		if err := d.Set("schedule", scanner.Schedule); err != nil {
			return diag.FromErr(err)
		}
	}

	// Set last_scan_timestamp if available
	if scanner.LastScanTimestamp != nil {
		if err := d.Set("last_scan_timestamp", *scanner.LastScanTimestamp); err != nil {
			return diag.FromErr(err)
		}
	}

	// Parse and set notification configuration
	if scanner.Notification != "" {
		notification, err := sdk.ParseNotificationConfiguration(scanner.Notification)
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

func UpdateContextTrustCenterScanner(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*provider.Context).Client

	scannerPackageId := d.Get("scanner_package_id").(string)
	scannerId := d.Get("scanner_id").(string)

	// Handle enabled changes
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

	// Handle schedule changes
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
			// Unset schedule
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
			// Unset notification
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

	// Disable the scanner (can't delete first-party scanners, just disable)
	enabled := false
	req := &sdk.SetScannerConfigurationRequest{
		ScannerPackageId: scannerPackageId,
		ScannerId:        scannerId,
		Enabled:          &enabled,
	}

	if err := client.TrustCenter.SetScannerConfiguration(ctx, req); err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Error disabling Trust Center scanner",
				Detail:   fmt.Sprintf("error disabling scanner %s/%s: %v", scannerPackageId, scannerId, err),
			},
		}
	}

	// Also unset schedule and notification to return to defaults
	unsetReq := &sdk.UnsetScannerConfigurationRequest{
		ScannerPackageId:  scannerPackageId,
		ScannerId:         scannerId,
		UnsetSchedule:     true,
		UnsetNotification: true,
	}
	// Ignore errors on unset - scanner is already disabled
	_ = client.TrustCenter.UnsetScannerConfiguration(ctx, unsetReq)

	d.SetId("")
	return nil
}

func ImportTrustCenterScanner(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// Import ID format: SNOWFLAKE/SECURITY_ESSENTIALS/MFA_CHECK
	id := d.Id()
	parts := strings.Split(id, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid import ID format, expected 'SOURCE/SCANNER_PACKAGE_ID/SCANNER_ID', got: %s", id)
	}

	if err := d.Set("scanner_package_id", parts[1]); err != nil {
		return nil, err
	}
	if err := d.Set("scanner_id", parts[2]); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
