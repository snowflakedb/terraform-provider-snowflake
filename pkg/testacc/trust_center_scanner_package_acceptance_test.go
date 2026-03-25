//go:build acceptance

package testacc

import (
	"context"
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func checkTrustCenterScannerPackageDisabled(scannerPackageId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		pkg, err := client.TrustCenter.ShowScannerPackageByID(ctx, scannerPackageId)
		if err != nil {
			return fmt.Errorf("error checking scanner package %s: %w", scannerPackageId, err)
		}
		if pkg.State != "FALSE" {
			return fmt.Errorf("scanner package %s expected to be disabled (FALSE), got: %s", scannerPackageId, pkg.State)
		}
		return nil
	}
}

func TestAcc_TrustCenterScannerPackage_Basic(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner_package.test"
	scannerPackageId := "SECURITY_ESSENTIALS"

	modelEnabled := model.TrustCenterScannerPackage("test", scannerPackageId).WithEnabled(true)
	modelDisabled := model.TrustCenterScannerPackage("test", scannerPackageId).WithEnabled(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: checkTrustCenterScannerPackageDisabled(scannerPackageId),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", scannerPackageId),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "show_output.0.state"),
					resource.TestCheckResourceAttrSet(resourceName, "show_output.0.provider_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "SNOWFLAKE|SECURITY_ESSENTIALS",
				ImportStateVerify: true,
			},
			{
				Config: accconfig.FromModels(t, modelDisabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScannerPackage_WithSchedule(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner_package.test"
	scannerPackageId := "SECURITY_ESSENTIALS"

	modelSchedule1 := model.TrustCenterScannerPackage("test", scannerPackageId).
		WithEnabled(true).
		WithSchedule("USING CRON 0 2 * * * UTC")
	modelSchedule2 := model.TrustCenterScannerPackage("test", scannerPackageId).
		WithEnabled(true).
		WithSchedule("USING CRON 0 4 * * * UTC")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: checkTrustCenterScannerPackageDisabled(scannerPackageId),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelSchedule1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", scannerPackageId),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 2 * * * UTC"),
				),
			},
			{
				Config: accconfig.FromModels(t, modelSchedule2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 4 * * * UTC"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScannerPackage_WithNotification(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner_package.test"
	scannerPackageId := "SECURITY_ESSENTIALS"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: checkTrustCenterScannerPackageDisabled(scannerPackageId),
		Steps: []resource.TestStep{
			{
				Config: trustCenterScannerPackageWithNotificationConfig("High"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", scannerPackageId),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.notify_admins", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.severity_threshold", "High"),
				),
			},
			{
				Config: trustCenterScannerPackageWithNotificationConfig("Critical"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "notification.0.severity_threshold", "Critical"),
				),
			},
		},
	})
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
