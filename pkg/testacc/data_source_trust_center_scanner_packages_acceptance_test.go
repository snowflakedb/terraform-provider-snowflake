//go:build account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TrustCenterScannerPackages_Basic(t *testing.T) {
	datasourceName := "data.snowflake_trust_center_scanner_packages.test"

	dsModel := datasourcemodel.TrustCenterScannerPackages("test")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, dsModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.0.show_output.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.0.show_output.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.0.show_output.0.state"),
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.0.show_output.0.provider_name"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScannerPackages_WithLikeFilter(t *testing.T) {
	datasourceName := "data.snowflake_trust_center_scanner_packages.test"

	dsModel := datasourcemodel.TrustCenterScannerPackages("test").WithLike("SECURITY%")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, dsModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.#"),
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.0.show_output.0.name"),
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.0.show_output.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.0.show_output.0.description"),
				),
			},
		},
	})
}
