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
	_, precreatedWarehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, precreatedWarehouseId, nil)
	t.Cleanup(precreatedWarehouseCleanup)
	_, precreatedWarehouse2Cleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, precreatedWarehouseId2, nil)
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
		warehouses, err := client.Warehouses.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(warehouses))
	})

	t.Run("show: like", func(t *testing.T) {
		showOptions := &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.Pointer(prefix + "%"),
			},
		}
		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Len(t, warehouses, 2)
	})

	t.Run("show: with options", func(t *testing.T) {
		showOptions := &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.Pointer(precreatedWarehouseId.Name()),
			},
		}
		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(warehouses))
		assert.Equal(t, precreatedWarehouseId.Name(), warehouses[0].Name)
		assert.Equal(t, sdk.Pointer(sdk.WarehouseSizeXSmall), warehouses[0].Size)
		assert.Equal(t, "ROLE", warehouses[0].OwnerRoleType)
	})

	t.Run("show: when searching a non-existent warehouse", func(t *testing.T) {
		showOptions := &sdk.ShowWarehouseOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Len(t, warehouses, 0)
	})

	t.Run("show: with starts with, and limit", func(t *testing.T) {
		t.Skip("TODO(SNOW-2683898): Unskip after the regression is fixed in Snowflake")
		showOptions := &sdk.ShowWarehouseOptions{
			StartsWith: sdk.String(precreatedWarehouseId.Name()),
			LimitFrom: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		}

		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)

		require.Len(t, warehouses, 1)
		require.Equal(t, precreatedWarehouseId.Name(), warehouses[0].Name)
	})

	t.Run("create: with resource constraint", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{
			ResourceConstraint: sdk.Pointer(sdk.WarehouseResourceConstraintMemory1X),
			WarehouseType:      sdk.Pointer(sdk.WarehouseTypeSnowparkOptimized),
			WarehouseSize:      sdk.Pointer(sdk.WarehouseSizeMedium),
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		result, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.WarehouseFromObject(t, result).
			HasResourceConstraint(sdk.WarehouseResourceConstraintMemory1X).
			HasNoGeneration().
			HasType(sdk.WarehouseTypeSnowparkOptimized).
			HasSize(sdk.WarehouseSizeMedium),
		)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{
			OrReplace:                       sdk.Bool(true),
			WarehouseType:                   sdk.Pointer(sdk.WarehouseTypeStandard),
			WarehouseSize:                   sdk.Pointer(sdk.WarehouseSizeSmall),
			MaxClusterCount:                 sdk.Int(8),
			MinClusterCount:                 sdk.Int(2),
			ScalingPolicy:                   sdk.Pointer(sdk.ScalingPolicyEconomy),
			AutoSuspend:                     sdk.Int(1000),
			AutoResume:                      sdk.Bool(true),
			InitiallySuspended:              sdk.Bool(false),
			ResourceMonitor:                 sdk.Pointer(resourceMonitor.ID()),
			Comment:                         sdk.String("comment"),
			EnableQueryAcceleration:         sdk.Bool(true),
			QueryAccelerationMaxScaleFactor: sdk.Int(90),
			MaxConcurrencyLevel:             sdk.Int(10),
			StatementQueuedTimeoutInSeconds: sdk.Int(2000),
			StatementTimeoutInSeconds:       sdk.Int(3000),
			Generation:                      sdk.Pointer(sdk.WarehouseGenerationStandardGen2),
			Tag: []sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: "v1",
				},
				{
					Name:  tag2.ID(),
					Value: "v2",
				},
			},
		})
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
		err := client.Warehouses.Create(ctx, id, nil)
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
		assert.Equal(t, sdk.Pointer(false), result.EnableQueryAcceleration)
		assert.Equal(t, sdk.Pointer(8), result.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, result.ResourceConstraint)
		assert.NotNil(t, result.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *result.Generation)
		assert.Nil(t, result.MaxQueryPerformanceLevel)
		assert.Nil(t, result.QueryThroughputMultiplier)
	})

	t.Run("create: empty comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{Comment: sdk.String("")})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		result, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", result.Comment)
	})

	t.Run("create adaptive: minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.CreateAdaptive(ctx, id, &sdk.CreateAdaptiveWarehouseOptions{})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		assertThatObject(t, objectassert.Warehouse(t, id).
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
			HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelLarge).
			HasQueryThroughputMultiplier(0),
		)
		assertThatObject(t, objectparametersassert.WarehouseParameters(t, id).
			HasStatementQueuedTimeoutInSeconds(0).
			HasStatementTimeoutInSeconds(172800),
		)
	})

	t.Run("create adaptive: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.CreateAdaptive(ctx, id, &sdk.CreateAdaptiveWarehouseOptions{
			Comment:                         sdk.String("test adaptive warehouse"),
			MaxQueryPerformanceLevel:        sdk.Pointer(sdk.MaxQueryPerformanceLevelMedium),
			QueryThroughputMultiplier:       sdk.Int(22),
			StatementQueuedTimeoutInSeconds: sdk.Int(30),
			StatementTimeoutInSeconds:       sdk.Int(60),
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		assertThatObject(t, objectassert.Warehouse(t, id).
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
		assertThatObject(t, objectparametersassert.WarehouseParameters(t, id).
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
		assert.Equal(t, sdk.Pointer(false), warehouse.EnableQueryAcceleration)
		assert.Equal(t, sdk.Pointer(8), warehouse.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouse.ResourceConstraint)
		assert.NotNil(t, warehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *warehouse.Generation)
		assert.Nil(t, warehouse.MaxQueryPerformanceLevel)
		assert.Nil(t, warehouse.QueryThroughputMultiplier)

		alterOptions := &sdk.AlterWarehouseOptions{
			// WarehouseType omitted on purpose - it requires suspending the warehouse (separate test cases)
			// ResourceConstraint omitted on purpose - it requires setting WarehouseType to SNOWPARK_OPTIMIZED (separate test cases)
			Set: &sdk.WarehouseSet{
				WarehouseSize:                   sdk.Pointer(sdk.WarehouseSizeMedium),
				WaitForCompletion:               sdk.Bool(true),
				MaxClusterCount:                 sdk.Int(3),
				MinClusterCount:                 sdk.Int(2),
				ScalingPolicy:                   sdk.Pointer(sdk.ScalingPolicyEconomy),
				AutoSuspend:                     sdk.Int(1234),
				AutoResume:                      sdk.Bool(false),
				ResourceMonitor:                 resourceMonitor.ID(),
				Comment:                         sdk.String("new comment"),
				EnableQueryAcceleration:         sdk.Bool(true),
				QueryAccelerationMaxScaleFactor: sdk.Int(2),
			},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
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
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *warehouseAfterSet.Generation)
		assert.Nil(t, warehouseAfterSet.MaxQueryPerformanceLevel)
		assert.Nil(t, warehouseAfterSet.QueryThroughputMultiplier)

		alterOptions = &sdk.AlterWarehouseOptions{
			// WarehouseSize omitted on purpose - UNSET is not supported for warehouse size
			// AutoSuspend omitted on purpose - UNSET works incorrectly (returns 0 instead of default 600)
			// WaitForCompletion omitted on purpose - no unset
			Unset: &sdk.WarehouseUnset{
				MaxClusterCount:                 sdk.Bool(true),
				MinClusterCount:                 sdk.Bool(true),
				ResourceMonitor:                 sdk.Bool(true),
				Comment:                         sdk.Bool(true),
				EnableQueryAcceleration:         sdk.Bool(true),
				QueryAccelerationMaxScaleFactor: sdk.Bool(true),
				WarehouseType:                   sdk.Bool(true),
				ScalingPolicy:                   sdk.Bool(true),
				AutoResume:                      sdk.Bool(true),
			},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		warehouseAfterUnset, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer(1), warehouseAfterUnset.MaxClusterCount)
		assert.Equal(t, sdk.Pointer(1), warehouseAfterUnset.MinClusterCount)
		assert.Equal(t, "", warehouseAfterUnset.ResourceMonitor.Name())
		assert.Equal(t, "", warehouseAfterUnset.Comment)
		assert.Equal(t, sdk.Pointer(false), warehouseAfterUnset.EnableQueryAcceleration)
		assert.Equal(t, sdk.Pointer(8), warehouseAfterUnset.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouseAfterUnset.ResourceConstraint)
		assert.NotNil(t, warehouseAfterUnset.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *warehouseAfterUnset.Generation)
		assert.Equal(t, sdk.WarehouseTypeStandard, warehouseAfterUnset.Type)
		assert.Equal(t, sdk.Pointer(sdk.ScalingPolicyStandard), warehouseAfterUnset.ScalingPolicy)
		assert.True(t, warehouseAfterUnset.AutoResume)
		assert.Nil(t, warehouseAfterUnset.MaxQueryPerformanceLevel)
		assert.Nil(t, warehouseAfterUnset.QueryThroughputMultiplier)
	})

	t.Run("alter adaptive: change warehouse type", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, id, &sdk.CreateWarehouseOptions{
			WarehouseSize: sdk.Pointer(sdk.WarehouseSizeMedium),
		})
		t.Cleanup(warehouseCleanup)

		// Wait for the warehouse to be started and confirm it's standard type
		condition := warehouseCondition(warehouse.ID(), func(r *sdk.Warehouse) bool {
			return r.State == sdk.WarehouseStateStarted && r.Type == sdk.WarehouseTypeStandard
		})
		require.Eventually(t, condition, 5*time.Second, time.Second)

		// Change warehouse type from standard to adaptive
		err := client.Warehouses.Alter(ctx, warehouse.ID(), &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{WarehouseType: sdk.Pointer(sdk.WarehouseTypeAdaptive)},
		})
		require.NoError(t, err)

		assertThatObject(t, objectassert.Warehouse(t, warehouse.ID()).
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
			HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelMedium).
			HasQueryThroughputMultiplier(0),
		)

		// Change warehouse type back from adaptive to standard
		err = client.Warehouses.Alter(ctx, warehouse.ID(), &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{
				WarehouseType: sdk.Pointer(sdk.WarehouseTypeStandard),
				WarehouseSize: sdk.Pointer(sdk.WarehouseSizeMedium),
			},
		})
		require.NoError(t, err)

		assertThatObject(t, objectassert.Warehouse(t, warehouse.ID()).
			HasType(sdk.WarehouseTypeStandard).
			HasSize(sdk.WarehouseSizeMedium).
			HasGeneration(sdk.WarehouseGenerationStandardGen1).
			HasNoResourceConstraint().
			HasMaxClusterCount(1).
			HasMinClusterCount(1).
			HasScalingPolicy(sdk.ScalingPolicyStandard).
			HasAutoSuspend(600).
			HasAutoResume(true).
			HasEnableQueryAcceleration(false).
			HasQueryAccelerationMaxScaleFactor(8).
			HasNoMaxQueryPerformanceLevel().
			HasNoQueryThroughputMultiplier(),
		)
	})

	t.Run("alter adaptive: set and unset all adaptive params", func(t *testing.T) {
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateAdaptive(t)
		t.Cleanup(warehouseCleanup)

		assertThatObject(t, objectassert.Warehouse(t, warehouse.ID()).
			HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelLarge).
			HasQueryThroughputMultiplier(0),
		)
		assertThatObject(t, objectparametersassert.WarehouseParameters(t, warehouse.ID()).
			HasStatementQueuedTimeoutInSeconds(0).
			HasStatementTimeoutInSeconds(172800),
		)

		err := client.Warehouses.Alter(ctx, warehouse.ID(), &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{
				MaxQueryPerformanceLevel:        sdk.Pointer(sdk.MaxQueryPerformanceLevelXSmall),
				QueryThroughputMultiplier:       sdk.Int(5),
				StatementQueuedTimeoutInSeconds: sdk.Int(100),
				StatementTimeoutInSeconds:       sdk.Int(200),
			},
		})
		require.NoError(t, err)

		assertThatObject(t, objectassert.Warehouse(t, warehouse.ID()).
			HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelXSmall).
			HasQueryThroughputMultiplier(5),
		)
		assertThatObject(t, objectparametersassert.WarehouseParameters(t, warehouse.ID()).
			HasStatementQueuedTimeoutInSeconds(100).
			HasStatementTimeoutInSeconds(200),
		)

		err = client.Warehouses.Alter(ctx, warehouse.ID(), &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{
				MaxQueryPerformanceLevel:        sdk.Bool(true),
				QueryThroughputMultiplier:       sdk.Bool(true),
				StatementQueuedTimeoutInSeconds: sdk.Bool(true),
				StatementTimeoutInSeconds:       sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		assertThatObject(t, objectassert.Warehouse(t, warehouse.ID()).
			HasMaxQueryPerformanceLevel(sdk.MaxQueryPerformanceLevelLarge).
			HasQueryThroughputMultiplier(2),
		)
		assertThatObject(t, objectparametersassert.WarehouseParameters(t, warehouse.ID()).
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

		alterOptions := &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{
				MaxConcurrencyLevel:             sdk.Int(4),
				StatementQueuedTimeoutInSeconds: sdk.Int(2),
				StatementTimeoutInSeconds:       sdk.Int(86400),
			},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		parametersAfterSet, err := client.Warehouses.ShowParameters(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, "4", helpers.FindParameter(t, parametersAfterSet, sdk.AccountParameterMaxConcurrencyLevel).Value)
		assert.Equal(t, "2", helpers.FindParameter(t, parametersAfterSet, sdk.AccountParameterStatementQueuedTimeoutInSeconds).Value)
		assert.Equal(t, "86400", helpers.FindParameter(t, parametersAfterSet, sdk.AccountParameterStatementTimeoutInSeconds).Value)

		alterOptions = &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{
				MaxConcurrencyLevel:             sdk.Bool(true),
				StatementQueuedTimeoutInSeconds: sdk.Bool(true),
				StatementTimeoutInSeconds:       sdk.Bool(true),
			},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
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
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, id, &sdk.CreateWarehouseOptions{
			WarehouseSize: sdk.Pointer(sdk.WarehouseSizeMedium),
		})
		t.Cleanup(warehouseCleanup)

		condition := warehouseCondition(warehouse.ID(), func(r *sdk.Warehouse) bool {
			return r.State == sdk.WarehouseStateStarted && r.Type == sdk.WarehouseTypeStandard
		})
		require.Eventually(t, condition, 5*time.Second, time.Second)

		alterOptions := &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{WarehouseType: sdk.Pointer(sdk.WarehouseTypeSnowparkOptimized)},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeSnowparkOptimized, returnedWarehouse.Type)
		assert.Contains(t, []any{sdk.WarehouseStateStarted, sdk.WarehouseStateResuming}, returnedWarehouse.State)

		alterOptions = &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{WarehouseType: sdk.Bool(true)},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
	})

	t.Run("alter: set and unset warehouse type with suspended warehouse", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// new warehouse created on purpose - we need medium to be able to use snowpark-optimized type
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, id, &sdk.CreateWarehouseOptions{
			WarehouseSize:      sdk.Pointer(sdk.WarehouseSizeMedium),
			InitiallySuspended: sdk.Bool(true),
		})
		t.Cleanup(warehouseCleanup)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
		assert.Contains(t, []any{sdk.WarehouseStateSuspended, sdk.WarehouseStateSuspending}, returnedWarehouse.State)

		alterOptions := &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{WarehouseType: sdk.Pointer(sdk.WarehouseTypeSnowparkOptimized)},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeSnowparkOptimized, returnedWarehouse.Type)
		assert.Equal(t, sdk.WarehouseStateSuspended, returnedWarehouse.State)

		alterOptions = &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{WarehouseType: sdk.Bool(true)},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
	})

	t.Run("alter: set and unset resource constraint", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// new warehouse created on purpose - we need medium to be able to use snowpark-optimized type
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, id, &sdk.CreateWarehouseOptions{
			WarehouseType: sdk.Pointer(sdk.WarehouseTypeSnowparkOptimized),
			WarehouseSize: sdk.Pointer(sdk.WarehouseSizeMedium),
		})
		t.Cleanup(warehouseCleanup)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeSnowparkOptimized, returnedWarehouse.Type)
		assert.Nil(t, returnedWarehouse.Generation)
		require.NotNil(t, returnedWarehouse.ResourceConstraint)
		assert.Equal(t, sdk.WarehouseResourceConstraintMemory16X, *returnedWarehouse.ResourceConstraint)

		alterOptions := &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{ResourceConstraint: sdk.Pointer(sdk.WarehouseResourceConstraintMemory1X)},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assertThatObject(t, objectassert.WarehouseFromObject(t, returnedWarehouse).
			HasResourceConstraint(sdk.WarehouseResourceConstraintMemory1X).
			HasNoGeneration().
			HasType(sdk.WarehouseTypeSnowparkOptimized).
			HasSize(sdk.WarehouseSizeMedium),
		)

		alterOptions = &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{ResourceConstraint: sdk.Bool(true)},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		require.NotNil(t, returnedWarehouse.ResourceConstraint)
		assert.Equal(t, sdk.WarehouseResourceConstraintMemory16X, *returnedWarehouse.ResourceConstraint)
	})

	t.Run("alter: set and unset generation", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, id, &sdk.CreateWarehouseOptions{})
		t.Cleanup(warehouseCleanup)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
		assert.Nil(t, returnedWarehouse.ResourceConstraint)
		assert.NotNil(t, returnedWarehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *returnedWarehouse.Generation)

		alterOptions := &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{Generation: sdk.Pointer(sdk.WarehouseGenerationStandardGen2)},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assertThatObject(t, objectassert.WarehouseFromObject(t, returnedWarehouse).
			HasGeneration(sdk.WarehouseGenerationStandardGen2).
			HasNoResourceConstraint().
			HasType(sdk.WarehouseTypeStandard),
		)

		alterOptions = &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{Generation: sdk.Bool(true)},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.NotNil(t, returnedWarehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *returnedWarehouse.Generation)
	})

	t.Run("alter: prove problems with unset auto suspend", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)
		alterOptions := &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{AutoSuspend: sdk.Bool(true)},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
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
		alterOptions := &sdk.AlterWarehouseOptions{
			NewName: &newID,
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, newID))

		result, err := client.Warehouses.Describe(ctx, newID)
		require.NoError(t, err)
		assert.Equal(t, newID.Name(), result.Name)
	})

	// This proves that we don't have to handle empty comment inside the resource.
	t.Run("alter: set empty comment versus unset", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{Comment: sdk.String("abc")})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Warehouse.DropWarehouseFunc(t, id))

		// can't use normal way, because of our SDK validation
		_, err = client.ExecForTests(ctx, fmt.Sprintf("ALTER WAREHOUSE %s SET COMMENT = ''", id.FullyQualifiedName()))
		require.NoError(t, err)

		warehouse, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", warehouse.Comment)

		alterOptions := &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{
				Comment: sdk.String("abc"),
			},
		}
		err = client.Warehouses.Alter(ctx, id, alterOptions)
		require.NoError(t, err)

		alterOptions = &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{
				Comment: sdk.Bool(true),
			},
		}
		err = client.Warehouses.Alter(ctx, id, alterOptions)
		require.NoError(t, err)

		warehouse, err = client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "", warehouse.Comment)
	})

	t.Run("alter: suspend and resume", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Suspend: sdk.Bool(true),
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		result, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateSuspended, sdk.WarehouseStateSuspending}, result.State)

		// check what happens if we suspend the already suspended warehouse
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.ErrorContains(t, err, "090064 (22000): Invalid state.")

		alterOptions = &sdk.AlterWarehouseOptions{
			Resume: sdk.Bool(true),
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)
		result, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateStarted, sdk.WarehouseStateResuming}, result.State)

		// check what happens if we resume the already started warehouse
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.ErrorContains(t, err, "090063 (22000): Invalid state.")
	})

	t.Run("alter: resume without suspending", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Resume:      sdk.Bool(true),
			IfSuspended: sdk.Bool(true),
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
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
		alterOptions := &sdk.AlterWarehouseOptions{
			AbortAllQueries: sdk.Bool(true),
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
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
		alterOptions := &sdk.AlterWarehouseOptions{
			Suspend: sdk.Bool(true),
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
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
		err := client.Warehouses.Alter(ctx, warehouse.ID(), &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{WarehouseSize: sdk.Pointer(sdk.WarehouseSizeMedium)},
		})
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

		err := client.Warehouses.Drop(ctx, warehouse.ID(), nil)
		require.NoError(t, err)
		_, err = client.Warehouses.Describe(ctx, warehouse.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: when warehouse does not exist", func(t *testing.T) {
		err := client.Warehouses.Drop(ctx, NonExistingAccountObjectIdentifier, nil)
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
	_, warehouse1Cleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, warehouseId1, nil)
	t.Cleanup(warehouse1Cleanup)
	_, warehouse2Cleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, warehouseId2, nil)
	t.Cleanup(warehouse2Cleanup)
	_, warehouse3Cleanup := testClientHelper().Warehouse.CreateWarehouseWithOptions(t, warehouseId3, nil)
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
		showOptions := &sdk.ShowWarehouseOptions{
			Like:       &sdk.Like{Pattern: sdk.String(warehouseId2.Name())},
			StartsWith: sdk.String(prefix),
			LimitFrom: &sdk.LimitFrom{
				Rows: sdk.Int(1),
			},
		}

		warehouses, err := client.Warehouses.Show(ctx, showOptions)
		require.NoError(t, err)
		require.Len(t, warehouses, 1)
		assert.Equal(t, warehouseId2.Name(), warehouses[0].Name)
	})
}
