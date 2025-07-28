//go:build !account_level_tests

package testacc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CurrentRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: currentRole(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_current_role.p", "name"),
				),
			},
		},
	})
}

func currentRole() string {
	s := `
	data snowflake_current_role p {}
	`
	return s
}
