//go:build !account_level_tests

package testacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SystemGetPrivateLinkConfig_aws(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: privateLinkConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "account_name"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "account_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "ocsp_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "aws_vpce_id"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "regionless_account_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "regionless_snowsight_url"),
					resource.TestCheckResourceAttrSet("data.snowflake_system_get_privatelink_config.p", "snowsight_url"),
				),
			},
		},
	})
}

func privateLinkConfig() string {
	s := `
	data snowflake_system_get_privatelink_config p {}
	`
	return s
}
