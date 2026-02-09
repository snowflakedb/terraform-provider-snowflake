//go:build acceptance

package testacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TrustCenterScanners_Basic(t *testing.T) {
	datasourceName := "data.snowflake_trust_center_scanners.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: trustCenterScannersDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Verify that at least one scanner is returned
					resource.TestCheckResourceAttrSet(datasourceName, "scanners.#"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScanners_WithPackageFilter(t *testing.T) {
	datasourceName := "data.snowflake_trust_center_scanners.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: trustCenterScannersWithPackageFilterConfig("SECURITY_ESSENTIALS"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "scanners.#"),
				),
			},
		},
	})
}

func trustCenterScannersDataSourceConfig() string {
	return `
data "snowflake_trust_center_scanners" "test" {
}
`
}

func trustCenterScannersWithPackageFilterConfig(packageId string) string {
	return `
data "snowflake_trust_center_scanners" "test" {
  scanner_package_id = "` + packageId + `"
}
`
}
