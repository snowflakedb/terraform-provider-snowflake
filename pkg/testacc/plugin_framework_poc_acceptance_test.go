package testacc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TerraformPluginFrameworkPoc_InitialSetup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// TODO [mux-PR]: 1.6?
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { TestAccPreCheck(t) },
		// TODO [mux-PR]: fill check destroy
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: someResourceConfig("abc"),
				Check:  nil,
			},
		},
	})
}

func someResourceConfig(value string) string {
	return fmt.Sprintf(`
resource "snowflake_some" "test" {
  todo = "%s"
}
`, value)
}
