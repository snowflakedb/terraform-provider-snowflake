//go:build !account_level_tests

package testacc

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"testing"
)

func TestAcc_CurrentAccount_Minimal(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	resourceMonitor, resourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "snowflake_current_account" "current" {
 resource_monitor = "%s"
}
`, resourceMonitor.ID().Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_current_account.current", "resource_monitor", resourceMonitor.ID().Name()),
				),
			},
		},
	})
}
