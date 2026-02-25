//go:build acceptance

package testacc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TrustCenterScannerPackage_Basic(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner_package.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create with enabled=true
			{
				Config: trustCenterScannerPackageBasicConfig(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", "SECURITY_ESSENTIALS"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "provider_type"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "SNOWFLAKE/SECURITY_ESSENTIALS",
				ImportStateVerify: true,
			},
			// Update to disabled
			{
				Config: trustCenterScannerPackageBasicConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScannerPackage_WithSchedule(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner_package.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create with schedule
			{
				Config: trustCenterScannerPackageWithScheduleConfig("USING CRON 0 2 * * * UTC"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", "SECURITY_ESSENTIALS"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 2 * * * UTC"),
				),
			},
			// Update schedule
			{
				Config: trustCenterScannerPackageWithScheduleConfig("USING CRON 0 4 * * * UTC"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 4 * * * UTC"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScannerPackage_WithNotification(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner_package.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create with notification
			{
				Config: trustCenterScannerPackageWithNotificationConfig("High"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", "SECURITY_ESSENTIALS"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.notify_admins", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.severity_threshold", "High"),
				),
			},
			// Update notification severity
			{
				Config: trustCenterScannerPackageWithNotificationConfig("Critical"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "notification.0.severity_threshold", "Critical"),
				),
			},
		},
	})
}

func trustCenterScannerPackageBasicConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "snowflake_trust_center_scanner_package" "test" {
  scanner_package_id = "SECURITY_ESSENTIALS"
  enabled            = %t
}
`, enabled)
}

func trustCenterScannerPackageWithScheduleConfig(schedule string) string {
	return fmt.Sprintf(`
resource "snowflake_trust_center_scanner_package" "test" {
  scanner_package_id = "SECURITY_ESSENTIALS"
  enabled            = true
  schedule           = "%s"
}
`, schedule)
}

func trustCenterScannerPackageWithNotificationConfig(severity string) string {
	return fmt.Sprintf(`
resource "snowflake_trust_center_scanner_package" "test" {
  scanner_package_id = "SECURITY_ESSENTIALS"
  enabled            = true

  notification {
    notify_admins      = true
    severity_threshold = "%s"
  }
}
`, severity)
}
