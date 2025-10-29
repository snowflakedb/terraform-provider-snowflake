//go:build account_level_tests

package testacc

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Experimental_Warehouse_ShowImprovedPerformance(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	warehouseId := testClient().Ids.RandomAccountObjectIdentifier()

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary).
		WithExperimentalFeaturesEnabled(string(experimentalfeatures.WarehouseShowImprovedPerformance))
	warehouseModel := model.Warehouse("test", warehouseId.Name())

	expectedWarehouseQuery := fmt.Sprintf("SHOW WAREHOUSES LIKE '%[1]s' STARTS WITH '%[1]s' LIMIT 1", warehouseId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, warehouseModel),
				Check: assertThat(t,
					resourceassert.WarehouseResource(t, warehouseModel.ResourceReference()).
						HasNameString(warehouseId.Name()),
					invokeactionassert.QueryHistoryEntry(t, secondaryTestClient(), expectedWarehouseQuery, tracking.CreateOperation, 100),
				),
			},
			{
				Config:       config.FromModels(t, providerModel, warehouseModel),
				ResourceName: warehouseModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedWarehouseResource(t, helpers.EncodeResourceIdentifier(warehouseId)).
						HasNameString(warehouseId.Name()),
					invokeactionassert.QueryHistoryEntryInImport(t, secondaryTestClient(), expectedWarehouseQuery, tracking.ImportOperation, 100),
				),
			},
		},
	})
}
