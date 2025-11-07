//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SystemGetSnowflakePlatformInfo(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: snowflakePlatformInfo(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_snowflake_platform_info.p", "aws_vpc_ids.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_snowflake_platform_info.p", "azure_vnet_subnet_ids.#"),
				),
			},
		},
	})
}

func snowflakePlatformInfo() string {
	s := `
	data snowflake_system_get_snowflake_platform_info "p" {}
	`
	return s
}
