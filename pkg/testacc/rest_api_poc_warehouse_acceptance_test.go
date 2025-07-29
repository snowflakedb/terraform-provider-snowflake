package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_RestApiPoc_WarehouseInitialCheck(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	id := testClient().Ids.RandomAccountObjectIdentifier()

	userWithPat := testClient().SetUpTemporaryLegacyServiceUserWithPat(t)
	userWithPatConfig := testClient().TempTomlConfigForServiceUserWithPat(t, userWithPat)

	testClient().Grant.GrantGlobalPrivilegesOnAccount(t, userWithPat.RoleId, []sdk.GlobalPrivilege{sdk.GlobalPrivilegeCreateWarehouse})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactoriesWithPluginPoc,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					t.Setenv(snowflakeenvs.ConfigPath, userWithPatConfig.Path)
				},
				Config: config.FromModels(t, userWithPatConfig) + warehouseRestApiPocResourceConfig(id),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr("snowflake_warehouse_rest_api_poc.test", "id", id.Name())),
					assert.Check(resource.TestCheckResourceAttr("snowflake_warehouse_rest_api_poc.test", "fully_qualified_name", id.FullyQualifiedName())),
					objectassert.Warehouse(t, id).
						HasName(id.Name()).
						HasState(sdk.WarehouseStateStarted).
						HasType(sdk.WarehouseTypeStandard).
						HasSize(sdk.WarehouseSizeXSmall).
						HasMaxClusterCount(1).
						HasMinClusterCount(1).
						HasScalingPolicy(sdk.ScalingPolicyStandard).
						HasAutoSuspend(600).
						HasAutoResume(true).
						HasResourceMonitor(sdk.AccountObjectIdentifier{}).
						HasComment("").
						HasEnableQueryAcceleration(false).
						HasQueryAccelerationMaxScaleFactor(8),
					objectparametersassert.WarehouseParameters(t, id).
						HasAllDefaults(),
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
