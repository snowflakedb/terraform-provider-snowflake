//go:build acceptance

package testacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TrustCenterScannerPackages_Basic(t *testing.T) {
	datasourceName := "data.snowflake_trust_center_scanner_packages.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: trustCenterScannerPackagesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Verify that at least one scanner package is returned
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.#"),
				),
			},
		},
	})
}

func TestAcc_TrustCenterScannerPackages_WithLikeFilter(t *testing.T) {
	datasourceName := "data.snowflake_trust_center_scanner_packages.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: trustCenterScannerPackagesWithLikeConfig("SECURITY%"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "scanner_packages.#"),
				),
			},
		},
	})
}

func trustCenterScannerPackagesDataSourceConfig() string {
	return `
data "snowflake_trust_center_scanner_packages" "test" {
}
`
}

func trustCenterScannerPackagesWithLikeConfig(like string) string {
	return `
data "snowflake_trust_center_scanner_packages" "test" {
  like = "` + like + `"
}
`
}
