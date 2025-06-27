package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TerraformPluginFrameworkPoc_InitialSetup(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

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
				Config: someResourceConfig(id),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_some.test", "id", id.FullyQualifiedName())),
				),
			},
		},
	})
}

func someResourceConfig(id sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_some" "test" {
  name = "%s"
}
`, id.Name())
}
