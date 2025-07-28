package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TerraformPluginFrameworkPoc_WarehouseInitialCheck(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { TestAccPreCheck(t) },
		// TODO [mux-PR]: fill check destroy; the protected enum is used now and this test is for a PoC resource
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: warehousePocResourceConfig(id),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_warehouse_poc.test", "id", id.Name())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_warehouse_poc.test", "fully_qualified_name", id.FullyQualifiedName())),
				),
			},
		},
	})
}

func warehousePocResourceConfig(id sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse_poc" "test" {
  name = "%s"
}
`, id.Name())
}
