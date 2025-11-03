//go:build account_level_tests

// These tests are temporarily moved to account level tests due to flakiness caused by changes in the higher-level parameters.

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Warehouses_BaseUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := testClient().Ids.RandomAccountObjectIdentifier()

	warehouseModel1 := model.Warehouse("test1", idOne.Name())
	warehouseModel2 := model.Warehouse("test2", idTwo.Name())
	warehouseModel3 := model.Warehouse("test3", idThree.Name())

	warehousesModelLikeFirstOne := datasourcemodel.Warehouses("test").
		WithLike(idOne.Name()).
		WithDependsOn(warehouseModel1.ResourceReference(), warehouseModel2.ResourceReference(), warehouseModel3.ResourceReference())

	warehousesModelLikePrefix := datasourcemodel.Warehouses("test").
		WithLike(prefix+"%").
		WithDependsOn(warehouseModel1.ResourceReference(), warehouseModel2.ResourceReference(), warehouseModel3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, warehouseModel1, warehouseModel2, warehouseModel3, warehousesModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehousesModelLikeFirstOne.DatasourceReference(), "warehouses.#", "1"),
				),
			},
			{
				Config: config.FromModels(t, warehouseModel1, warehouseModel2, warehouseModel3, warehousesModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(warehousesModelLikePrefix.DatasourceReference(), "warehouses.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Warehouses_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	warehouseModel := model.Warehouse("test", id.Name()).
		WithComment(comment)

	warehousesModel := datasourcemodel.Warehouses("test").
		WithLike(id.Name()).
		WithDependsOn(warehouseModel.ResourceReference())

	warehousesModelWithoutOptionals := datasourcemodel.Warehouses("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(warehouseModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, warehouseModel, warehousesModel),
				Check: assertThat(t,
					resourceshowoutputassert.WarehousesDatasourceShowOutput(t, warehousesModel.DatasourceReference()).
						HasName(id.Name()).
						HasStateNotEmpty().
						HasType(sdk.WarehouseTypeStandard).
						HasSize(sdk.WarehouseSizeXSmall).
						HasMinClusterCount(1).
						HasMaxClusterCount(1).
						HasStartedClustersNotEmpty().
						HasRunningNotEmpty().
						HasQueuedNotEmpty().
						HasIsDefault(false).
						HasAutoSuspend(600).
						HasAutoResume(true).
						HasAvailableNotEmpty().
						HasProvisioningNotEmpty().
						HasQuiescingNotEmpty().
						HasOtherNotEmpty().
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(comment).
						HasEnableQueryAcceleration(false).
						HasQueryAccelerationMaxScaleFactor(8).
						HasResourceMonitorEmpty().
						HasScalingPolicy(sdk.ScalingPolicyStandard).
						HasOwnerRoleTypeNotEmpty().
						HasResourceConstraintEmpty().
						HasGenerationNotEmpty(),

					assert.Check(resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.name")),
					assert.Check(resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.kind", "WAREHOUSE")),

					resourceparametersassert.WarehousesDatasourceParameters(t, warehousesModel.DatasourceReference()).
						HasDefaultMaxConcurrencyLevel().
						HasDefaultStatementQueuedTimeoutInSeconds().
						HasDefaultStatementTimeoutInSeconds(),
				),
			},
			{
				Config: config.FromModels(t, warehouseModel, warehousesModelWithoutOptionals),
				Check: assertThat(t,
					resourceshowoutputassert.WarehousesDatasourceShowOutput(t, "test").
						HasName(id.Name()).
						HasStateNotEmpty().
						HasType(sdk.WarehouseTypeStandard).
						HasSize(sdk.WarehouseSizeXSmall).
						HasMinClusterCount(1).
						HasMaxClusterCount(1).
						HasStartedClustersNotEmpty().
						HasRunningNotEmpty().
						HasQueuedNotEmpty().
						HasIsDefault(false).
						HasAutoSuspend(600).
						HasAutoResume(true).
						HasAvailableNotEmpty().
						HasProvisioningNotEmpty().
						HasQuiescingNotEmpty().
						HasOtherNotEmpty().
						HasCreatedOnNotEmpty().
						HasResumedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasOwnerNotEmpty().
						HasComment(comment).
						HasEnableQueryAcceleration(false).
						HasQueryAccelerationMaxScaleFactor(8).
						HasResourceMonitorEmpty().
						HasScalingPolicy(sdk.ScalingPolicyStandard).
						HasOwnerRoleTypeNotEmpty().
						HasResourceConstraintEmpty().
						HasGenerationNotEmpty(),

					assert.Check(resource.TestCheckResourceAttr(warehousesModelWithoutOptionals.DatasourceReference(), "warehouses.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(warehousesModelWithoutOptionals.DatasourceReference(), "warehouses.0.parameters.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Warehouses_WarehouseNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Warehouses/without_warehouse"),
				ExpectError:     regexp.MustCompile("there should be at least one warehouse"),
			},
		},
	})
}
