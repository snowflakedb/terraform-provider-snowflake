package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [this PR]: set up this test (new user with PAT, role assigned with only CREATE WAREHOUSE privilege, use proper config, turn of configure client once)
func TestAcc_RestApiPoc_WarehouseInitialCheck(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: warehouseRestApiPocResourceConfig(id),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_warehouse_rest_api_poc.test", "id", id.Name())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_warehouse_rest_api_poc.test", "fully_qualified_name", id.FullyQualifiedName())),
				),
			},
		},
	})
}

func warehouseRestApiPocResourceConfig(id sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_warehouse_rest_api_poc" "test" {
  name = "%s"
}
`, id.Name())
}
