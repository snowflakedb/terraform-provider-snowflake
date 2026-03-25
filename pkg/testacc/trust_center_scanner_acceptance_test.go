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

func checkTrustCenterScannerDisabled(scannerPackageId, scannerId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		scanner, err := client.TrustCenter.ShowScannerByID(ctx, scannerPackageId, scannerId)
		if err != nil {
			return fmt.Errorf("error checking scanner %s/%s: %w", scannerPackageId, scannerId, err)
		}
		if scanner.State != "FALSE" {
			return fmt.Errorf("scanner %s/%s expected to be disabled (FALSE), got: %s", scannerPackageId, scannerId, scanner.State)
		}
		return nil
	}
}

func TestAcc_TrustCenterScanner_Basic(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner.test"
	scannerPackageId := "SECURITY_ESSENTIALS"
	scannerId := "SECURITY_ESSENTIALS_MFA_REQUIRED_FOR_USERS_CHECK"

	modelEnabled := model.TrustCenterScanner("test", scannerPackageId, scannerId).WithEnabled(true)
	modelDisabled := model.TrustCenterScanner("test", scannerPackageId, scannerId).WithEnabled(false)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: checkTrustCenterScannerDisabled(scannerPackageId, scannerId),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelEnabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scanner_package_id", scannerPackageId),
					resource.TestCheckResourceAttr(resourceName, "scanner_id", scannerId),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "show_output.0.state"),
					resource.TestCheckResourceAttrSet(resourceName, "show_output.0.name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     "SNOWFLAKE|SECURITY_ESSENTIALS|SECURITY_ESSENTIALS_MFA_REQUIRED_FOR_USERS_CHECK",
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

func TestAcc_TrustCenterScanner_WithSchedule(t *testing.T) {
	resourceName := "snowflake_trust_center_scanner.test"
	scannerPackageId := "SECURITY_ESSENTIALS"
	scannerId := "SECURITY_ESSENTIALS_MFA_REQUIRED_FOR_USERS_CHECK"

	modelSchedule1 := model.TrustCenterScanner("test", scannerPackageId, scannerId).
		WithEnabled(true).
		WithSchedule("USING CRON 0 0 * * * UTC")
	modelSchedule2 := model.TrustCenterScanner("test", scannerPackageId, scannerId).
		WithEnabled(true).
		WithSchedule("USING CRON 0 6 * * * UTC")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: checkTrustCenterScannerDisabled(scannerPackageId, scannerId),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelSchedule1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 0 * * * UTC"),
				),
			},
			{
				Config: accconfig.FromModels(t, modelSchedule2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "schedule", "USING CRON 0 6 * * * UTC"),
				),
			},
		},
	})
}
