//go:build account_level_tests

package testint

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Warehouses(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	prefix := random.StringN(6)
	precreatedWarehouseId := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	precreatedWarehouseId2 := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	// new warehouses created on purpose
	_, precreatedWarehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(precreatedWarehouseId))
	t.Cleanup(precreatedWarehouseCleanup)
	_, precreatedWarehouse2Cleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(precreatedWarehouseId2))
	t.Cleanup(precreatedWarehouse2Cleanup)

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	tag2, tag2Cleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tag2Cleanup)

	resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
	t.Cleanup(resourceMonitorCleanup)

	startLongRunningQuery := func() {
		go client.ExecForTests(ctx, "CALL SYSTEM$WAIT(15);") //nolint:errcheck // we don't care if this eventually errors, as long as it runs for a little while
		time.Sleep(3 * time.Second)
	}

	warehouseCondition := func(id sdk.AccountObjectIdentifier, pred func(r *sdk.Warehouse) bool) func() bool {
		return func() bool {
			r, e := client.Warehouses.ShowByID(ctx, id)
			return e == nil && pred(r)
		}
	}

	t.Run("show: without options", func(t *testing.T) {
		warehouses, err := client.Warehouses.Show(ctx, sdk.NewShowWarehouseRequest())
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(warehouses))
	})

	t.Run("show: like", func(t *testing.T) {
		warehouses, err := client.Warehouses.Show(ctx, sdk.NewShowWarehouseRequest().
			WithLike(sdk.Like{Pattern: sdk.Pointer(prefix + "%")}))
		require.NoError(t, err)
		assert.Len(t, warehouses, 2)
	})

	t.Run("show: with options", func(t *testing.T) {
		warehouses, err := client.Warehouses.Show(ctx, sdk.NewShowWarehouseRequest().
			WithLike(sdk.Like{Pattern: sdk.Pointer(precreatedWarehouseId.Name())}))
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		assert.Equal(t, precreatedWarehouseId.Name(), warehouses[0].Name)
		assert.Equal(t, sdk.Pointer(sdk.WarehouseSizeXSmall), warehouses[0].Size)
		assert.Equal(t, "ROLE", warehouses[0].OwnerRoleType)
	})

	t.Run("show: when searching a non-existent warehouse", func(t *testing.T) {
		warehouses, err := client.Warehouses.Show(ctx, sdk.NewShowWarehouseRequest().
			WithLike(sdk.Like{Pattern: sdk.String("non-existent")}))
		require.NoError(t, err)
		assert.Len(t, warehouses, 0)
	})

	t.Run("show: with starts with, and limit", func(t *testing.T) {
		warehouses, err := client.Warehouses.Show(ctx, sdk.NewShowWarehouseRequest().
			WithStartsWith(precreatedWarehouseId.Name()).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)

		require.Len(t, warehouses, 1)
		require.Equal(t, precreatedWarehouseId.Name(), warehouses[0].Name)
	})

	t.Run("create: with resource constraint", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, sdk.NewCreateWarehouseRequest(id).
			WithResourceConstraint(sdk.WarehouseResourceConstraintMemory1X).
			WithWarehouseType(sdk.WarehouseTypeSnowparkOptimized).
			WithWarehouseSize(sdk.WarehouseSizeMedium))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		result, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(
			t, objectassert.WarehouseFromObject(t, result).
				HasResourceConstraint(sdk.WarehouseResourceConstraintMemory1X).
				HasNoGeneration().
				HasType(sdk.WarehouseTypeSnowparkOptimized).
				HasSize(sdk.WarehouseSizeMedium),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, sdk.NewCreateWarehouseRequest(id).
			WithOrReplace(true).
			WithWarehouseType(sdk.WarehouseTypeStandard).
			WithWarehouseSize(sdk.WarehouseSizeSmall).
			WithMaxClusterCount(8).
			WithMinClusterCount(2).
			WithScalingPolicy(sdk.ScalingPolicyEconomy).
			WithAutoSuspend(1000).
			WithAutoResume(true).
			WithInitiallySuspended(false).
			WithResourceMonitor(resourceMonitor.ID()).
			WithComment("comment").
			WithEnableQueryAcceleration(true).
			WithQueryAccelerationMaxScaleFactor(90).
			WithMaxConcurrencyLevel(10).
			WithStatementQueuedTimeoutInSeconds(2000).
			WithStatementTimeoutInSeconds(3000).
			WithGeneration(sdk.WarehouseGenerationStandardGen2).
			WithTag([]sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: "v1",
				},
				{
					Name:  tag2.ID(),
					Value: "v2",
				},
			}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		// we can use the same assertion builder in the SDK tests
		assertThatObject(t, objectassert.Warehouse(t, id).
			HasName(id.Name()).
			HasType(sdk.WarehouseTypeStandard).
			HasSize(sdk.WarehouseSizeSmall).
			HasMaxClusterCount(8).
			HasMinClusterCount(2).
			HasScalingPolicy(sdk.ScalingPolicyEconomy).
			HasAutoSuspend(1000).
			HasAutoResume(true).
			HasStateOneOf(sdk.WarehouseStateResuming, sdk.WarehouseStateStarted).
			HasResourceMonitor(resourceMonitor.ID()).
			HasComment("comment").
			HasEnableQueryAcceleration(true).
			HasQueryAccelerationMaxScaleFactor(90).
			HasGeneration(sdk.WarehouseGenerationStandardGen2).
			HasNoResourceConstraint().
			HasNoMaxQueryPerformanceLevel().
			HasNoQueryThroughputMultiplier())

		warehouse, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), warehouse.Name)
		assert.Equal(t, sdk.WarehouseTypeStandard, warehouse.Type)
		assert.Equal(t, sdk.Pointer(sdk.WarehouseSizeSmall), warehouse.Size)
		assert.Equal(t, sdk.Pointer(8), warehouse.MaxClusterCount)
		assert.Equal(t, sdk.Pointer(2), warehouse.MinClusterCount)
		assert.Equal(t, sdk.Pointer(sdk.ScalingPolicyEconomy), warehouse.ScalingPolicy)
		assert.Equal(t, sdk.Pointer(1000), warehouse.AutoSuspend)
		assert.True(t, warehouse.AutoResume)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateResuming, sdk.WarehouseStateStarted}, warehouse.State)
		assert.Equal(t, resourceMonitor.ID().Name(), warehouse.ResourceMonitor.Name())
		assert.Equal(t, "comment", warehouse.Comment)
		assert.Equal(t, sdk.Pointer(true), warehouse.EnableQueryAcceleration)
		assert.Equal(t, sdk.Pointer(90), warehouse.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouse.ResourceConstraint)
		assert.NotNil(t, warehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen2, *warehouse.Generation)
		assert.Nil(t, warehouse.MaxQueryPerformanceLevel)
		assert.Nil(t, warehouse.QueryThroughputMultiplier)

		// we can also use the read object to initialize:
		assertThatObject(t, objectassert.WarehouseFromObject(t, warehouse).
			HasName(id.Name()).
			HasType(sdk.WarehouseTypeStandard).
			HasSize(sdk.WarehouseSizeSmall).
			HasMaxClusterCount(8).
			HasMinClusterCount(2).
			HasScalingPolicy(sdk.ScalingPolicyEconomy).
			HasAutoSuspend(1000).
			HasAutoResume(true).
			HasStateOneOf(sdk.WarehouseStateResuming, sdk.WarehouseStateStarted).
			HasResourceMonitor(resourceMonitor.ID()).
			HasComment("comment").
			HasEnableQueryAcceleration(true).
			HasQueryAccelerationMaxScaleFactor(90).
			HasNoResourceConstraint().
			HasGeneration(sdk.WarehouseGenerationStandardGen2).
			HasNoMaxQueryPerformanceLevel().
			HasNoQueryThroughputMultiplier())

		tag1Value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer("v1"), tag1Value)
		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer("v2"), tag2Value)
	})

	t.Run("create: no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, sdk.NewCreateWarehouseRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		result, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), result.Name)
		assert.Equal(t, sdk.WarehouseTypeStandard, result.Type)
		assert.Equal(t, sdk.Pointer(sdk.WarehouseSizeXSmall), result.Size)
		assert.Equal(t, sdk.Pointer(1), result.MaxClusterCount)
		assert.Equal(t, sdk.Pointer(1), result.MinClusterCount)
		assert.Equal(t, sdk.Pointer(sdk.ScalingPolicyStandard), result.ScalingPolicy)
		assert.Equal(t, sdk.Pointer(600), result.AutoSuspend)
		assert.True(t, result.AutoResume)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateResuming, sdk.WarehouseStateStarted}, result.State)
		assert.Equal(t, "", result.Comment)
		assert.Equal(t, sdk.Pointer(true), result.EnableQueryAcceleration)
		assert.Equal(t, sdk.Pointer(2), result.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, result.ResourceConstraint)
		assert.NotNil(t, result.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen2, *result.Generation)
		assert.Nil(t, result.MaxQueryPerformanceLevel)
		assert.Nil(t, result.QueryThroughputMultiplier)
	})

	t.Run("create: empty comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, sdk.NewCreateWarehouseRequest(id).WithComment(""))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		result, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", result.Comment)
	})

	t.Run("create adaptive: minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.CreateAdaptive(ctx, sdk.NewCreateAdaptiveWarehouseRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		assertThatObject(
			t, objectassert.Warehouse(t, id).
				HasName(id.Name()).
				HasType(sdk.WarehouseTypeAdaptive).
				HasComment("").
				HasNoSize().
				HasNoGeneration().
				HasNoResourceConstraint().
				HasNoMaxClusterCount().
				HasNoMinClusterCount().
				HasNoScalingPolicy().
				HasNoAutoSuspend().
				HasAutoResume(true).
				HasNoEnableQueryAcceleration().
				HasNoQueryAccelerationMaxScaleFactor().
				HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelXLarge).
				HasQueryThroughputMultiplier(2),
		)
		assertThatObject(
			t, objectparametersassert.WarehouseParameters(t, id).
				HasStatementQueuedTimeoutInSeconds(0).
				HasStatementTimeoutInSeconds(172800),
		)
	})

	t.Run("create adaptive: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.CreateAdaptive(ctx, sdk.NewCreateAdaptiveWarehouseRequest(id).
			WithComment("test adaptive warehouse").
			WithMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelMedium).
			WithQueryThroughputMultiplier(22).
			WithStatementQueuedTimeoutInSeconds(30).
			WithStatementTimeoutInSeconds(60))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		assertThatObject(
			t, objectassert.Warehouse(t, id).
				HasName(id.Name()).
				HasType(sdk.WarehouseTypeAdaptive).
				HasComment("test adaptive warehouse").
				HasNoSize().
				HasNoGeneration().
				HasNoResourceConstraint().
				HasNoMaxClusterCount().
				HasNoMinClusterCount().
				HasNoScalingPolicy().
				HasNoAutoSuspend().
				HasAutoResume(true).
				HasNoEnableQueryAcceleration().
				HasNoQueryAccelerationMaxScaleFactor().
				HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelMedium).
				HasQueryThroughputMultiplier(22),
		)
		assertThatObject(
			t, objectparametersassert.WarehouseParameters(t, id).
				HasStatementQueuedTimeoutInSeconds(30).
				HasStatementTimeoutInSeconds(60),
		)
	})

	t.Run("alter: set and unset", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		assert.Equal(t, sdk.Pointer(sdk.WarehouseSizeXSmall), warehouse.Size)
		assert.Equal(t, sdk.Pointer(1), warehouse.MaxClusterCount)
		assert.Equal(t, sdk.Pointer(1), warehouse.MinClusterCount)
		assert.Equal(t, sdk.Pointer(sdk.ScalingPolicyStandard), warehouse.ScalingPolicy)
		assert.Equal(t, sdk.Pointer(600), warehouse.AutoSuspend)
		assert.True(t, warehouse.AutoResume)
		assert.Equal(t, "", warehouse.ResourceMonitor.Name())
		assert.Equal(t, "", warehouse.Comment)
		assert.Equal(t, sdk.Pointer(true), warehouse.EnableQueryAcceleration)
		assert.Equal(t, sdk.Pointer(2), warehouse.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouse.ResourceConstraint)
		assert.NotNil(t, warehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen2, *warehouse.Generation)
		assert.Nil(t, warehouse.MaxQueryPerformanceLevel)
		assert.Nil(t, warehouse.QueryThroughputMultiplier)

		// WarehouseType omitted on purpose - it requires suspending the warehouse (separate test cases)
		// ResourceConstraint omitted on purpose - it requires setting WarehouseType to SNOWPARK_OPTIMIZED (separate test cases)
		err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().
				WithWarehouseSize(sdk.WarehouseSizeMedium).
				WithWaitForCompletion(true).
				WithMaxClusterCount(3).
				WithMinClusterCount(2).
				WithScalingPolicy(sdk.ScalingPolicyEconomy).
				WithAutoSuspend(1234).
				WithAutoResume(false).
				WithResourceMonitor(resourceMonitor.ID()).
				WithComment("new comment").
				WithEnableQueryAcceleration(true).
				WithQueryAccelerationMaxScaleFactor(2)))
		require.NoError(t, err)

		warehouseAfterSet, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer(sdk.WarehouseSizeMedium), warehouseAfterSet.Size)
		assert.Equal(t, sdk.Pointer(3), warehouseAfterSet.MaxClusterCount)
		assert.Equal(t, sdk.Pointer(2), warehouseAfterSet.MinClusterCount)
		assert.Equal(t, sdk.Pointer(sdk.ScalingPolicyEconomy), warehouseAfterSet.ScalingPolicy)
		assert.Equal(t, sdk.Pointer(1234), warehouseAfterSet.AutoSuspend)
		assert.False(t, warehouseAfterSet.AutoResume)
		assert.Equal(t, resourceMonitor.ID().Name(), warehouseAfterSet.ResourceMonitor.Name())
		assert.Equal(t, "new comment", warehouseAfterSet.Comment)
		assert.Equal(t, sdk.Pointer(true), warehouseAfterSet.EnableQueryAcceleration)
		assert.Equal(t, sdk.Pointer(2), warehouseAfterSet.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouseAfterSet.ResourceConstraint)
		assert.NotNil(t, warehouseAfterSet.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen2, *warehouseAfterSet.Generation)
		assert.Nil(t, warehouseAfterSet.MaxQueryPerformanceLevel)
		assert.Nil(t, warehouseAfterSet.QueryThroughputMultiplier)

		// WarehouseSize omitted on purpose - UNSET is not supported for warehouse size
		// AutoSuspend omitted on purpose - UNSET works incorrectly (returns 0 instead of default 600)
		// WaitForCompletion omitted on purpose - no unset
		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithUnset(*sdk.NewWarehouseUnsetRequest().
				WithMaxClusterCount(true).
				WithMinClusterCount(true).
				WithResourceMonitor(true).
				WithComment(true).
				WithEnableQueryAcceleration(true).
				WithQueryAccelerationMaxScaleFactor(true).
				WithWarehouseType(true).
				WithScalingPolicy(true).
				WithAutoResume(true)))
		require.NoError(t, err)

		warehouseAfterUnset, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer(1), warehouseAfterUnset.MaxClusterCount)
		assert.Equal(t, sdk.Pointer(1), warehouseAfterUnset.MinClusterCount)
		assert.Equal(t, "", warehouseAfterUnset.ResourceMonitor.Name())
		assert.Equal(t, "", warehouseAfterUnset.Comment)
		assert.Equal(t, sdk.Pointer(true), warehouseAfterUnset.EnableQueryAcceleration)
		assert.Equal(t, sdk.Pointer(2), warehouseAfterUnset.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouseAfterUnset.ResourceConstraint)
		assert.NotNil(t, warehouseAfterUnset.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen2, *warehouseAfterUnset.Generation)
		assert.Equal(t, sdk.WarehouseTypeStandard, warehouseAfterUnset.Type)
		assert.Equal(t, sdk.Pointer(sdk.ScalingPolicyStandard), warehouseAfterUnset.ScalingPolicy)
		assert.True(t, warehouseAfterUnset.AutoResume)
		assert.Nil(t, warehouseAfterUnset.MaxQueryPerformanceLevel)
		assert.Nil(t, warehouseAfterUnset.QueryThroughputMultiplier)
	})

	t.Run("alter adaptive: change warehouse type", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(
			t, sdk.NewCreateWarehouseRequest(id).
				WithWarehouseSize(sdk.WarehouseSizeMedium),
		)
		t.Cleanup(warehouseCleanup)

		// Wait for the warehouse to be started and confirm it's standard type
		condition := warehouseCondition(warehouse.ID(), func(r *sdk.Warehouse) bool {
			return r.State == sdk.WarehouseStateStarted && r.Type == sdk.WarehouseTypeStandard
		})
		require.Eventually(t, condition, 5*time.Second, time.Second)

		// Change warehouse type from standard to adaptive
		err := client.Warehouses.AlterWithSuspend(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().WithWarehouseType(sdk.WarehouseTypeAdaptive)))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Warehouse(t, warehouse.ID()).
				HasType(sdk.WarehouseTypeAdaptive).
				HasNoSize().
				HasNoGeneration().
				HasNoResourceConstraint().
				HasNoMaxClusterCount().
				HasNoMinClusterCount().
				HasNoScalingPolicy().
				HasNoAutoSuspend().
				HasAutoResume(true).
				HasNoEnableQueryAcceleration().
				HasNoQueryAccelerationMaxScaleFactor().
				HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelLarge).
				HasQueryThroughputMultiplier(2),
		)

		// Change warehouse type back from adaptive to standard
		err = client.Warehouses.AlterWithSuspend(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().
				WithWarehouseType(sdk.WarehouseTypeStandard).
				WithWarehouseSize(sdk.WarehouseSizeMedium)))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Warehouse(t, warehouse.ID()).
				HasType(sdk.WarehouseTypeStandard).
				HasSize(sdk.WarehouseSizeMedium).
				HasGeneration(sdk.WarehouseGenerationStandardGen2).
				HasNoResourceConstraint().
				HasMaxClusterCount(1).
				HasMinClusterCount(1).
				HasScalingPolicy(sdk.ScalingPolicyStandard).
				HasAutoSuspend(600).
				HasAutoResume(true).
				HasEnableQueryAcceleration(true).
				HasQueryAccelerationMaxScaleFactor(2).
				HasNoMaxQueryPerformanceLevel().
				HasNoQueryThroughputMultiplier(),
		)
	})

	t.Run("alter adaptive: set and unset all adaptive params", func(t *testing.T) {
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateAdaptive(t)
		t.Cleanup(warehouseCleanup)

		assertThatObject(
			t, objectassert.Warehouse(t, warehouse.ID()).
				HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelXLarge).
				HasQueryThroughputMultiplier(2),
		)
		assertThatObject(
			t, objectparametersassert.WarehouseParameters(t, warehouse.ID()).
				HasStatementQueuedTimeoutInSeconds(0).
				HasStatementTimeoutInSeconds(172800),
		)

		err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().
				WithMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelXSmall).
				WithQueryThroughputMultiplier(5).
				WithStatementQueuedTimeoutInSeconds(100).
				WithStatementTimeoutInSeconds(200)))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Warehouse(t, warehouse.ID()).
				HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelXSmall).
				HasQueryThroughputMultiplier(5),
		)
		assertThatObject(
			t, objectparametersassert.WarehouseParameters(t, warehouse.ID()).
				HasStatementQueuedTimeoutInSeconds(100).
				HasStatementTimeoutInSeconds(200),
		)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithUnset(*sdk.NewWarehouseUnsetRequest().
				WithMaxQueryPerformanceLevel(true).
				WithQueryThroughputMultiplier(true).
				WithStatementQueuedTimeoutInSeconds(true).
				WithStatementTimeoutInSeconds(true)))
		require.NoError(t, err)

		assertThatObject(
			t, objectassert.Warehouse(t, warehouse.ID()).
				HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelXLarge).
				HasQueryThroughputMultiplier(2),
		)
		assertThatObject(
			t, objectparametersassert.WarehouseParameters(t, warehouse.ID()).
				HasStatementQueuedTimeoutInSeconds(0).
				HasStatementTimeoutInSeconds(172800),
		)
	})

	t.Run("alter: set and unset parameters", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		parameters, err := client.Warehouses.ShowParameters(ctx, warehouse.ID())
		require.NoError(t, err)

		assert.Equal(t, "8", helpers.FindParameter(t, parameters, sdk.AccountParameterMaxConcurrencyLevel).Value)
		assert.Equal(t, "0", helpers.FindParameter(t, parameters, sdk.AccountParameterStatementQueuedTimeoutInSeconds).Value)
		assert.Equal(t, "172800", helpers.FindParameter(t, parameters, sdk.AccountParameterStatementTimeoutInSeconds).Value)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().
				WithMaxConcurrencyLevel(4).
				WithStatementQueuedTimeoutInSeconds(2).
				WithStatementTimeoutInSeconds(86400)))
		require.NoError(t, err)

		parametersAfterSet, err := client.Warehouses.ShowParameters(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, "4", helpers.FindParameter(t, parametersAfterSet, sdk.AccountParameterMaxConcurrencyLevel).Value)
		assert.Equal(t, "2", helpers.FindParameter(t, parametersAfterSet, sdk.AccountParameterStatementQueuedTimeoutInSeconds).Value)
		assert.Equal(t, "86400", helpers.FindParameter(t, parametersAfterSet, sdk.AccountParameterStatementTimeoutInSeconds).Value)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithUnset(*sdk.NewWarehouseUnsetRequest().
				WithMaxConcurrencyLevel(true).
				WithStatementQueuedTimeoutInSeconds(true).
				WithStatementTimeoutInSeconds(true)))
		require.NoError(t, err)

		parametersAfterUnset, err := client.Warehouses.ShowParameters(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, "8", helpers.FindParameter(t, parametersAfterUnset, sdk.AccountParameterMaxConcurrencyLevel).Value)
		assert.Equal(t, "0", helpers.FindParameter(t, parametersAfterUnset, sdk.AccountParameterStatementQueuedTimeoutInSeconds).Value)
		assert.Equal(t, "172800", helpers.FindParameter(t, parametersAfterUnset, sdk.AccountParameterStatementTimeoutInSeconds).Value)
	})

	t.Run("alter: set and unset warehouse type with started warehouse", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// new warehouse created on purpose - we need medium to be able to use snowpark-optimized type
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(
			t, sdk.NewCreateWarehouseRequest(id).
				WithWarehouseSize(sdk.WarehouseSizeMedium),
		)
		t.Cleanup(warehouseCleanup)

		condition := warehouseCondition(warehouse.ID(), func(r *sdk.Warehouse) bool {
			return r.State == sdk.WarehouseStateStarted && r.Type == sdk.WarehouseTypeStandard
		})
		require.Eventually(t, condition, 5*time.Second, time.Second)

		err := client.Warehouses.AlterWithSuspend(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().WithWarehouseType(sdk.WarehouseTypeSnowparkOptimized)))
		require.NoError(t, err)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeSnowparkOptimized, returnedWarehouse.Type)
		assert.Contains(t, []any{sdk.WarehouseStateStarted, sdk.WarehouseStateResuming}, returnedWarehouse.State)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithUnset(*sdk.NewWarehouseUnsetRequest().WithWarehouseType(true)))
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
	})

	t.Run("alter: set and unset warehouse type with suspended warehouse", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// new warehouse created on purpose - we need medium to be able to use snowpark-optimized type
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(
			t, sdk.NewCreateWarehouseRequest(id).
				WithWarehouseSize(sdk.WarehouseSizeMedium).
				WithInitiallySuspended(true),
		)
		t.Cleanup(warehouseCleanup)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
		assert.Contains(t, []any{sdk.WarehouseStateSuspended, sdk.WarehouseStateSuspending}, returnedWarehouse.State)

		err = client.Warehouses.AlterWithSuspend(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().WithWarehouseType(sdk.WarehouseTypeSnowparkOptimized)))
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeSnowparkOptimized, returnedWarehouse.Type)
		assert.Equal(t, sdk.WarehouseStateSuspended, returnedWarehouse.State)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithUnset(*sdk.NewWarehouseUnsetRequest().WithWarehouseType(true)))
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
	})

	t.Run("alter: set and unset resource constraint", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// new warehouse created on purpose - we need medium to be able to use snowpark-optimized type
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(
			t, sdk.NewCreateWarehouseRequest(id).
				WithWarehouseType(sdk.WarehouseTypeSnowparkOptimized).
				WithWarehouseSize(sdk.WarehouseSizeMedium),
		)
		t.Cleanup(warehouseCleanup)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeSnowparkOptimized, returnedWarehouse.Type)
		assert.Nil(t, returnedWarehouse.Generation)
		require.NotNil(t, returnedWarehouse.ResourceConstraint)
		assert.Equal(t, sdk.WarehouseResourceConstraintMemory16X, *returnedWarehouse.ResourceConstraint)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().WithResourceConstraint(sdk.WarehouseResourceConstraintMemory1X)))
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assertThatObject(
			t, objectassert.WarehouseFromObject(t, returnedWarehouse).
				HasResourceConstraint(sdk.WarehouseResourceConstraintMemory1X).
				HasNoGeneration().
				HasType(sdk.WarehouseTypeSnowparkOptimized).
				HasSize(sdk.WarehouseSizeMedium),
		)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithUnset(*sdk.NewWarehouseUnsetRequest().WithResourceConstraint(true)))
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		require.NotNil(t, returnedWarehouse.ResourceConstraint)
		assert.Equal(t, sdk.WarehouseResourceConstraintMemory16X, *returnedWarehouse.ResourceConstraint)
	})

	t.Run("alter: set and unset generation", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(id))
		t.Cleanup(warehouseCleanup)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
		assert.Nil(t, returnedWarehouse.ResourceConstraint)
		assert.NotNil(t, returnedWarehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen2, *returnedWarehouse.Generation)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().WithGeneration(sdk.WarehouseGenerationStandardGen1)))
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assertThatObject(
			t, objectassert.WarehouseFromObject(t, returnedWarehouse).
				HasGeneration(sdk.WarehouseGenerationStandardGen1).
				HasNoResourceConstraint().
				HasType(sdk.WarehouseTypeStandard),
		)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithUnset(*sdk.NewWarehouseUnsetRequest().WithGeneration(true)))
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.NotNil(t, returnedWarehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen2, *returnedWarehouse.Generation)
	})

	t.Run("alter: prove problems with unset auto suspend", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithUnset(*sdk.NewWarehouseUnsetRequest().WithAutoSuspend(true)))
		require.NoError(t, err)
		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		// TODO [SNOW-1473453]: change when UNSET starts working correctly (expecting to unset to default 600)
		// assert.Equal(t, sdk.Pointer(600), returnedWarehouse.AutoSuspend)
		assert.Nil(t, returnedWarehouse.AutoSuspend)
	})

	t.Run("alter: rename", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		newID := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithRenameTo(newID))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, newID))

		result, err := client.Warehouses.Describe(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), result.Name)
	})

	// This proves that we don't have to handle empty comment inside the resource.
	t.Run("alter: set empty comment versus unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, sdk.NewCreateWarehouseRequest(id).WithComment("abc"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		// can't use normal way, because of our SDK validation
		_, err = client.ExecForTests(ctx, fmt.Sprintf("ALTER WAREHOUSE %s SET COMMENT = ''", id.FullyQualifiedName()))
		require.NoError(t, err)

		warehouse, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", warehouse.Comment)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).
			WithSet(*sdk.NewWarehouseSetRequest().WithComment("abc")))
		require.NoError(t, err)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).
			WithUnset(*sdk.NewWarehouseUnsetRequest().WithComment(true)))
		require.NoError(t, err)

		warehouse, err = client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", warehouse.Comment)
	})

	t.Run("alter: suspend and resume", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSuspend(true))
		require.NoError(t, err)

		result, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateSuspended, sdk.WarehouseStateSuspending}, result.State)

		// check what happens if we suspend the already suspended warehouse
		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSuspend(true))
		require.ErrorContains(t, err, "090064 (22000): Invalid state.")

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithResume(true))
		require.NoError(t, err)
		result, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateStarted, sdk.WarehouseStateResuming}, result.State)

		// check what happens if we resume the already started warehouse
		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithResume(true))
		require.ErrorContains(t, err, "090063 (22000): Invalid state.")
	})

	t.Run("alter: resume without suspending", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithResume(true).
			WithIfSuspended(true))
		require.NoError(t, err)

		result, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateStarted, sdk.WarehouseStateResuming}, result.State)
	})

	t.Run("alter: abort all queries", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		startLongRunningQuery()

		// Check that query is running
		result, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer(1), result.Running)
		assert.Equal(t, sdk.Pointer(0), result.Queued)

		// Abort all queries
		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithAbortAllQueries(true))
		require.NoError(t, err)

		// Wait for abort to be effective
		time.Sleep(2 * time.Second)

		// Check no query is running
		result, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer(0), result.Running)
		assert.Equal(t, sdk.Pointer(0), result.Queued)
	})

	t.Run("alter: suspend with a long running-query", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		startLongRunningQuery()

		// Suspend the warehouse
		err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSuspend(true))
		require.NoError(t, err)

		// check the state - it seems that the warehouse is suspended despite having a running query on it
		result, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer(1), result.Running)
		assert.Equal(t, sdk.Pointer(0), result.Queued)

		condition := warehouseCondition(warehouse.ID(), func(r *sdk.Warehouse) bool { return r.State == sdk.WarehouseStateSuspended })
		require.Eventually(t, condition, 10*time.Second, time.Second)
	})

	t.Run("alter: resize with a long running-query", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		startLongRunningQuery()

		// Resize the warehouse
		err := client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(warehouse.ID()).
			WithSet(*sdk.NewWarehouseSetRequest().WithWarehouseSize(sdk.WarehouseSizeMedium)))
		require.NoError(t, err)

		// check the state - it seems it's resized despite query being run on it
		result, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateResizing, sdk.WarehouseStateStarted}, result.State)
		assert.Equal(t, sdk.Pointer(sdk.WarehouseSizeMedium), result.Size)
		assert.Equal(t, sdk.Pointer(1), result.Running)
		assert.Equal(t, sdk.Pointer(0), result.Queued)
	})

	t.Run("describe: when warehouse exists", func(t *testing.T) {
		result, err := client.Warehouses.Describe(ctx, precreatedWarehouseId)
		require.NoError(t, err)
		assert.Equal(t, precreatedWarehouseId.Name(), result.Name)
		assert.Equal(t, "WAREHOUSE", result.Kind)
		assert.NotEmpty(t, result.CreatedOn)
	})

	t.Run("describe: when warehouse does not exist", func(t *testing.T) {
		id := NonExistingAccountObjectIdentifier
		_, err := client.Warehouses.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when warehouse exists", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		err := client.Warehouses.Drop(ctx, sdk.NewDropWarehouseRequest(warehouse.ID()))
		require.NoError(t, err)
		_, err = client.Warehouses.Describe(ctx, warehouse.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when warehouse does not exist", func(t *testing.T) {
		err := client.Warehouses.Drop(ctx, sdk.NewDropWarehouseRequest(NonExistingAccountObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_Warehouses_Experimental(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	prefix := random.StringN(6) + "_"
	warehouseId1 := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	warehouseId2 := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	warehouseId3 := testClientHelper().Ids.RandomAccountObjectIdentifier()
	_, warehouse1Cleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(warehouseId1))
	t.Cleanup(warehouse1Cleanup)
	_, warehouse2Cleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(warehouseId2))
	t.Cleanup(warehouse2Cleanup)
	_, warehouse3Cleanup := testClientHelper().Warehouse.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(warehouseId3))
	t.Cleanup(warehouse3Cleanup)

	t.Run("show experimental", func(t *testing.T) {
		wh, err := client.Warehouses.ShowByIDExperimental(ctx, warehouseId1)
		require.NoError(t, err)
		assert.Equal(t, warehouseId1.Name(), wh.Name)

		wh, err = client.Warehouses.ShowByIDExperimental(ctx, warehouseId2)
		require.NoError(t, err)
		assert.Equal(t, warehouseId2.Name(), wh.Name)

		wh, err = client.Warehouses.ShowByIDExperimental(ctx, warehouseId3)
		require.NoError(t, err)
		assert.Equal(t, warehouseId3.Name(), wh.Name)
	})

	t.Run("show experimental safely", func(t *testing.T) {
		wh, err := client.Warehouses.ShowByIDExperimentalSafely(ctx, warehouseId1)
		require.NoError(t, err)
		assert.Equal(t, warehouseId1.Name(), wh.Name)

		wh, err = client.Warehouses.ShowByIDExperimentalSafely(ctx, warehouseId2)
		require.NoError(t, err)
		assert.Equal(t, warehouseId2.Name(), wh.Name)

		wh, err = client.Warehouses.ShowByIDExperimentalSafely(ctx, warehouseId3)
		require.NoError(t, err)
		assert.Equal(t, warehouseId3.Name(), wh.Name)
	})

	t.Run("show using starts with prefix", func(t *testing.T) {
		warehouses, err := client.Warehouses.Show(ctx, sdk.NewShowWarehouseRequest().
			WithLike(sdk.Like{Pattern: sdk.String(warehouseId2.Name())}).
			WithStartsWith(prefix).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, warehouses, 1)
		assert.Equal(t, warehouseId2.Name(), warehouses[0].Name)
	})
}

func TestInt_Warehouses_Interactive(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("create interactive: minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.CreateInteractive(ctx, sdk.NewCreateInteractiveWarehouseRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		assertThatObject(
			t, objectassert.Warehouse(t, id).
				HasName(id.Name()).
				HasType(sdk.WarehouseTypeInteractive).
				HasComment("").
				HasSize(sdk.WarehouseSizeXSmall).
				HasMaxClusterCount(1).
				HasMinClusterCount(1).
				HasAutoSuspend(86400).
				HasAutoResume(true).
				HasScalingPolicy(sdk.ScalingPolicyStandard).
				HasEnableQueryAcceleration(false).
				HasQueryAccelerationMaxScaleFactor(8).
				HasNoResourceConstraint().
				HasNoGeneration().
				HasNoMaxQueryPerformanceLevel().
				HasNoQueryThroughputMultiplier().
				HasTables(),
		)
		assertThatObject(
			t, objectparametersassert.WarehouseInteractiveParameters(t, id).
				HasMaxConcurrencyLevel(8).
				HasStatementQueuedTimeoutInSeconds(0).
				HasStatementTimeoutInSeconds(172800).
				HasFallbackWarehouse(""),
		)
	})

	t.Run("create interactive: with tables", func(t *testing.T) {
		table1, table1Cleanup := testClientHelper().Table.CreateInteractiveTable(t)
		t.Cleanup(table1Cleanup)
		table2, table2Cleanup := testClientHelper().Table.CreateInteractiveTable(t)
		t.Cleanup(table2Cleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.CreateInteractive(ctx, sdk.NewCreateInteractiveWarehouseRequest(id).
			WithTables([]sdk.SchemaObjectIdentifier{table1, table2}).
			WithWarehouseSize(sdk.WarehouseSizeXSmall).
			WithComment("interactive warehouse"))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		assertThatObject(
			t, objectassert.Warehouse(t, id).
				HasName(id.Name()).
				HasType(sdk.WarehouseTypeInteractive).
				HasComment("interactive warehouse").
				HasSize(sdk.WarehouseSizeXSmall).
				HasMaxClusterCount(1).
				HasMinClusterCount(1).
				HasAutoSuspend(86400).
				HasAutoResume(true).
				HasScalingPolicy(sdk.ScalingPolicyStandard).
				HasEnableQueryAcceleration(false).
				HasQueryAccelerationMaxScaleFactor(8).
				HasNoResourceConstraint().
				HasNoGeneration().
				HasNoMaxQueryPerformanceLevel().
				HasNoQueryThroughputMultiplier().
				HasTables(table1, table2),
		)
		assertThatObject(
			t, objectparametersassert.WarehouseInteractiveParameters(t, id).
				HasMaxConcurrencyLevel(8).
				HasStatementQueuedTimeoutInSeconds(0).
				HasStatementTimeoutInSeconds(172800).
				HasFallbackWarehouse(""),
		)
	})

	t.Run("alter: add and drop tables", func(t *testing.T) {
		table1, table1Cleanup := testClientHelper().Table.CreateInteractiveTable(t)
		t.Cleanup(table1Cleanup)
		table2, table2Cleanup := testClientHelper().Table.CreateInteractiveTable(t)
		t.Cleanup(table2Cleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.CreateInteractive(ctx, sdk.NewCreateInteractiveWarehouseRequest(id).
			WithTables([]sdk.SchemaObjectIdentifier{table1}))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithAddTables([]sdk.SchemaObjectIdentifier{table2}))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Warehouse(t, id).HasTables(table1, table2))

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithDropTables([]sdk.SchemaObjectIdentifier{table1}))
		require.NoError(t, err)
		assertThatObject(t, objectassert.Warehouse(t, id).HasTables(table2))
	})

	t.Run("alter: set and unset fallback warehouse", func(t *testing.T) {
		fallback, fallbackCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(fallbackCleanup)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.CreateInteractive(ctx, sdk.NewCreateInteractiveWarehouseRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).
			WithSet(*sdk.NewWarehouseSetRequest().WithFallbackWarehouse(fallback.ID())))
		require.NoError(t, err)

		assertThatObject(
			t, objectparametersassert.WarehouseInteractiveParameters(t, id).
				HasFallbackWarehouse(fallback.ID().Name()),
		)

		err = client.Warehouses.Alter(ctx, sdk.NewAlterWarehouseRequest(id).
			WithUnset(*sdk.NewWarehouseUnsetRequest().WithFallbackWarehouse(true)))
		require.NoError(t, err)

		assertThatObject(
			t, objectparametersassert.WarehouseInteractiveParameters(t, id).
				HasFallbackWarehouse(""),
		)
	})
}
