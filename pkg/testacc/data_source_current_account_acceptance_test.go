//go:build !account_level_tests

package testacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CurrentAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: currentAccount(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_current_account.p", "account"),
					resource.TestCheckResourceAttrSet("data.snowflake_current_account.p", "region"),
					resource.TestCheckResourceAttrSet("data.snowflake_current_account.p", "url"),
				),
			},
		},
	})
}

func currentAccount() string {
	s := `
	data snowflake_current_account p {}
	`
	return s
}
