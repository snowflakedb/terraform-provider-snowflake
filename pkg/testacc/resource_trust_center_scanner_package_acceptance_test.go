//go:build account_level_tests

package testacc

import (
	"context"
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/config"
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

func trustCenterScannerPackageNotificationVariable(notifyAdmins bool, severity string) config.Variable {
	return config.ListVariable(
		config.ObjectVariable(map[string]config.Variable{
			"notify_admins":      config.BoolVariable(notifyAdmins),
			"severity_threshold": config.StringVariable(severity),
		}),
	)
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

func TestAcc_TrustCenterScannerPackage_Complete(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner_package.test"
	scannerPackageId := "SECURITY_ESSENTIALS"

	completeModel := model.TrustCenterScannerPackage("test", scannerPackageId).
		WithEnabled(true).
		WithSchedule("USING CRON 0 2 * * * UTC").
		WithNotificationValue(trustCenterScannerPackageNotificationVariable(true, "High"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: checkTrustCenterScannerPackageDisabled(scannerPackageId),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, completeModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", scannerPackageId),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 2 * * * UTC"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.notify_admins", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.severity_threshold", "High"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "SNOWFLAKE|SECURITY_ESSENTIALS",
				ImportStateVerify: true,
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

	modelNotification1 := model.TrustCenterScannerPackage("test", scannerPackageId).
		WithEnabled(true).
		WithNotificationValue(trustCenterScannerPackageNotificationVariable(true, "High"))
	modelNotification2 := model.TrustCenterScannerPackage("test", scannerPackageId).
		WithEnabled(true).
		WithNotificationValue(trustCenterScannerPackageNotificationVariable(true, "Critical"))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: checkTrustCenterScannerPackageDisabled(scannerPackageId),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelNotification1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", scannerPackageId),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.notify_admins", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.severity_threshold", "High"),
				),
			},
			{
				Config: accconfig.FromModels(t, modelNotification2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "notification.0.severity_threshold", "Critical"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScannerPackage_ScheduleRemoval(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner_package.test"
	scannerPackageId := "SECURITY_ESSENTIALS"

	modelWithSchedule := model.TrustCenterScannerPackage("test", scannerPackageId).
		WithEnabled(true).
		WithSchedule("USING CRON 0 2 * * * UTC")
	modelWithoutSchedule := model.TrustCenterScannerPackage("test", scannerPackageId).
		WithEnabled(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: checkTrustCenterScannerPackageDisabled(scannerPackageId),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWithSchedule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 2 * * * UTC"),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithoutSchedule),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "schedule"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScannerPackage_NotificationRemoval(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner_package.test"
	scannerPackageId := "SECURITY_ESSENTIALS"

	modelWithNotification := model.TrustCenterScannerPackage("test", scannerPackageId).
		WithEnabled(true).
		WithNotificationValue(trustCenterScannerPackageNotificationVariable(true, "High"))
	modelWithoutNotification := model.TrustCenterScannerPackage("test", scannerPackageId).
		WithEnabled(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: checkTrustCenterScannerPackageDisabled(scannerPackageId),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelWithNotification),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "notification.0.notify_admins", "true"),
					resource.TestCheckResourceAttr(resourceName, "notification.0.severity_threshold", "High"),
				),
			},
			{
				Config: accconfig.FromModels(t, modelWithoutNotification),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "notification.0.notify_admins"),
				),
			},
		},
	})
}
