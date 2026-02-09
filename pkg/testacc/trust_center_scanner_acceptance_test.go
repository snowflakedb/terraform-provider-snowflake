//go:build acceptance

package testacc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TrustCenterScanner_Basic(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create with enabled=true
			{
				Config: trustCenterScannerBasicConfig(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", "SECURITY_ESSENTIALS"),
					resource.TestCheckResourceAttr(resourceName, "scanner_id", "SECURITY_ESSENTIALS_MFA_REQUIRED_FOR_USERS_CHECK"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "SNOWFLAKE/SECURITY_ESSENTIALS/SECURITY_ESSENTIALS_MFA_REQUIRED_FOR_USERS_CHECK",
				ImportStateVerify: true,
			},
			// Update to disabled
			{
				Config: trustCenterScannerBasicConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScanner_WithSchedule(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// Create with schedule
			{
				Config: trustCenterScannerWithScheduleConfig("USING CRON 0 0 * * * UTC"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 0 * * * UTC"),
				),
			},
			// Update schedule
			{
				Config: trustCenterScannerWithScheduleConfig("USING CRON 0 6 * * * UTC"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 6 * * * UTC"),
				),
			},
		},
	})
}

func trustCenterScannerBasicConfig(enabled bool) string {
	return fmt.Sprintf(`
resource "snowflake_trust_center_scanner" "test" {
  scanner_package_id = "SECURITY_ESSENTIALS"
  scanner_id         = "SECURITY_ESSENTIALS_MFA_REQUIRED_FOR_USERS_CHECK"
  enabled            = %t
}
`, enabled)
}

func trustCenterScannerWithScheduleConfig(schedule string) string {
	return fmt.Sprintf(`
resource "snowflake_trust_center_scanner" "test" {
  scanner_package_id = "SECURITY_ESSENTIALS"
  scanner_id         = "SECURITY_ESSENTIALS_MFA_REQUIRED_FOR_USERS_CHECK"
  enabled            = true
  schedule           = "%s"
}
`, schedule)
}
