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
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
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
		WithWithDescribe(false).
		WithWithParameters(false).
		WithLike(idOne.Name()).
		WithDependsOn(warehouseModel1.ResourceReference(), warehouseModel2.ResourceReference(), warehouseModel3.ResourceReference())

	warehousesModelLikePrefix := datasourcemodel.Warehouses("test").
		WithWithDescribe(false).
		WithWithParameters(false).
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

	warehousesModelWithoutOptionals := datasourcemodel.Warehouses("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(warehouseModel.ResourceReference())

	warehousesModel := datasourcemodel.Warehouses("test").
		WithLike(id.Name()).
		WithDependsOn(warehouseModel.ResourceReference())

	commonShowAssert := func(t *testing.T, datasourceReference string) *resourceshowoutputassert.WarehouseShowOutputAssert {
		t.Helper()
		assert := resourceshowoutputassert.WarehousesDatasourceShowOutput(t, datasourceReference).
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
			// TODO(SNOW-2852741): Different auto suspend default on different test environments
			// HasAutoSuspend(600).
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
			HasEnableQueryAcceleration(true).
			HasQueryAccelerationMaxScaleFactor(2).
			HasResourceMonitorEmpty().
			HasScalingPolicy(sdk.ScalingPolicyStandard).
			HasOwnerRoleTypeNotEmpty().
			HasResourceConstraintEmpty().
			HasMaxQueryPerformanceLevelEmpty().
			HasQueryThroughputMultiplier(0).
			HasNoTables()
		if testClient().SnowflakeDefaults.WarehouseGenerationEmptyByDefault(t) {
			assert = assert.HasGenerationEmpty()
		} else {
			assert = assert.HasGenerationNotEmpty()
		}

		return assert
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Warehouse),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, warehouseModel, warehousesModelWithoutOptionals),
				Check: assertThat(
					t,
					commonShowAssert(t, warehousesModelWithoutOptionals.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(warehousesModelWithoutOptionals.DatasourceReference(), "warehouses.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(warehousesModelWithoutOptionals.DatasourceReference(), "warehouses.0.parameters.#", "0")),
				),
			},
			{
				Config: config.FromModels(t, warehouseModel, warehousesModel),
				Check: assertThat(
					t,
					commonShowAssert(t, warehousesModel.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.name")),
					assert.Check(resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.kind", "WAREHOUSE")),

					resourceparametersassert.WarehousesDatasourceParameters(t, warehousesModel.DatasourceReference()).
						HasDefaultMaxConcurrencyLevel().
						HasDefaultStatementQueuedTimeoutInSeconds().
						HasDefaultStatementTimeoutInSeconds().
						HasDefaultFallbackWarehouse(),
				),
			},
		},
	})
}

func TestAcc_Warehouses_CompleteUseCase_AdaptiveWarehouse(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	warehouseModel := model.WarehouseAdaptive("test", id.Name()).
		WithComment(comment).
		WithMaxQueryPerformanceLevel(string(sdk.MaxQueryPerformanceLevelMedium)).
		WithQueryThroughputMultiplier(2).
		WithStatementQueuedTimeoutInSeconds(300).
		WithStatementTimeoutInSeconds(86400)
	warehousesModelWithoutOptionals := datasourcemodel.Warehouses("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(warehouseModel.ResourceReference())

	warehousesModel := datasourcemodel.Warehouses("test").
		WithLike(id.Name()).
		WithDependsOn(warehouseModel.ResourceReference())

	commonShowAssert := func(t *testing.T, datasourceReference string) *resourceshowoutputassert.WarehouseShowOutputAssert {
		t.Helper()
		return resourceshowoutputassert.WarehousesDatasourceShowOutput(t, datasourceReference).
			HasName(id.Name()).
			HasStateNotEmpty().
			HasType(sdk.WarehouseTypeAdaptive).
			HasSizeEmpty().
			HasMinClusterCount(0).
			HasMaxClusterCount(0).
			HasIsDefault(false).
			HasAutoSuspend(0).
			HasAutoResume(true).
			HasCreatedOnNotEmpty().
			HasResumedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasOwnerNotEmpty().
			HasComment(comment).
			HasEnableQueryAcceleration(false).
			HasQueryAccelerationMaxScaleFactor(0).
			HasResourceMonitorEmpty().
			HasScalingPolicyEmpty().
			HasOwnerRoleTypeNotEmpty().
			HasResourceConstraintEmpty().
			HasGenerationEmpty().
			HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelMedium).
			HasQueryThroughputMultiplier(2).
			HasNoTables()
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseAdaptive),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, warehouseModel, warehousesModelWithoutOptionals),
				Check: assertThat(
					t,
					commonShowAssert(t, warehousesModelWithoutOptionals.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(warehousesModelWithoutOptionals.DatasourceReference(), "warehouses.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(warehousesModelWithoutOptionals.DatasourceReference(), "warehouses.0.parameters.#", "0")),
				),
			},
			{
				Config: config.FromModels(t, warehouseModel, warehousesModel),
				Check: assertThat(
					t,
					commonShowAssert(t, warehousesModel.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.name")),
					assert.Check(resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.kind", "WAREHOUSE")),

					resourceparametersassert.WarehousesDatasourceParameters(t, warehousesModel.DatasourceReference()).
						HasStatementQueuedTimeoutInSeconds(300).
						HasStatementQueuedTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse).
						HasStatementTimeoutInSeconds(86400).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse).
						HasDefaultFallbackWarehouse(),
				),
			},
		},
	})
}

func TestAcc_Warehouses_CompleteUseCase_InteractiveWarehouse(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	resourceMonitor, resourceMonitorCleanup := testClient().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	fallbackWarehouse, fallbackWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(fallbackWarehouseCleanup)

	interactiveTable, interactiveTableCleanup := testClient().Table.CreateInteractiveTable(t)
	t.Cleanup(interactiveTableCleanup)

	interactiveWarehouseModel := model.WarehouseInteractiveWithId(id).
		WithInitiallySuspended(true).
		WithWarehouseSize(string(sdk.WarehouseSizeSmall)).
		WithMinClusterCount(2).
		WithMaxClusterCount(3).
		WithAutoSuspend(98000).
		WithAutoResume(r.BooleanFalse).
		WithResourceMonitor(resourceMonitor.ID().Name()).
		WithFallbackWarehouse(fallbackWarehouse.ID().Name()).
		WithComment(comment).
		WithMaxConcurrencyLevel(4).
		WithStatementQueuedTimeoutInSeconds(30).
		WithStatementTimeoutInSeconds(45).
		WithTables(interactiveTable.FullyQualifiedName())

	warehousesModelWithoutOptionals := datasourcemodel.Warehouses("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithWithParameters(false).
		WithDependsOn(interactiveWarehouseModel.ResourceReference())

	warehousesModel := datasourcemodel.Warehouses("test").
		WithLike(id.Name()).
		WithDependsOn(interactiveWarehouseModel.ResourceReference())

	commonShowAssert := func(t *testing.T, datasourceReference string) *resourceshowoutputassert.WarehouseShowOutputAssert {
		t.Helper()
		return resourceshowoutputassert.WarehousesDatasourceShowOutput(t, datasourceReference).
			HasName(id.Name()).
			HasStateNotEmpty().
			HasType(sdk.WarehouseTypeInteractive).
			HasSize(sdk.WarehouseSizeSmall).
			HasMinClusterCount(2).
			HasMaxClusterCount(3).
			HasIsDefault(false).
			HasAutoSuspend(98000).
			HasAutoResume(false).
			HasCreatedOnNotEmpty().
			HasResumedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasOwnerNotEmpty().
			HasComment(comment).
			HasEnableQueryAcceleration(false).
			HasQueryAccelerationMaxScaleFactor(0).
			HasResourceMonitor(resourceMonitor.ID()).
			HasScalingPolicyEmpty().
			HasOwnerRoleTypeNotEmpty().
			HasResourceConstraintEmpty().
			HasGenerationEmpty().
			HasMaxQueryPerformanceLevelEmpty().
			HasQueryThroughputMultiplier(0).
			HasTables(interactiveTable.FullyQualifiedName())
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: interactiveWarehouseProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.WarehouseInteractive),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, interactiveWarehouseModel, warehousesModelWithoutOptionals),
				Check: assertThat(
					t,
					commonShowAssert(t, warehousesModelWithoutOptionals.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(warehousesModelWithoutOptionals.DatasourceReference(), "warehouses.0.describe_output.#", "0")),
					assert.Check(resource.TestCheckResourceAttr(warehousesModelWithoutOptionals.DatasourceReference(), "warehouses.0.parameters.#", "0")),
				),
			},
			{
				Config: config.FromModels(t, interactiveWarehouseModel, warehousesModel),
				Check: assertThat(
					t,
					commonShowAssert(t, warehousesModel.DatasourceReference()),

					assert.Check(resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.name")),
					assert.Check(resource.TestCheckResourceAttr(warehousesModel.DatasourceReference(), "warehouses.0.describe_output.0.kind", "WAREHOUSE")),

					resourceparametersassert.WarehousesDatasourceParameters(t, warehousesModel.DatasourceReference()).
						HasMaxConcurrencyLevel(4).
						HasMaxConcurrencyLevelLevel(sdk.ParameterTypeWarehouse).
						HasStatementQueuedTimeoutInSeconds(30).
						HasStatementQueuedTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse).
						HasStatementTimeoutInSeconds(45).
						HasStatementTimeoutInSecondsLevel(sdk.ParameterTypeWarehouse).
						HasFallbackWarehouse(fallbackWarehouse.ID().Name()).
						HasFallbackWarehouseLevel(sdk.ParameterTypeWarehouse),
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
