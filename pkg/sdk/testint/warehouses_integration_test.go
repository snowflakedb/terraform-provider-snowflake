//go:build account_level_tests

package testint

import (
	"fmt"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Warehouses(t *testing.T) {
	client := testClient(t)
	secondaryClient := testSecondaryClient(t)
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
		assert.Equal(t, sdk.WarehouseSizeXSmall, warehouses[0].Size)
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

	t.Run("create: with resource constraint & generation", func(t *testing.T) {
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
			ResourceConstraint:              sdk.Pointer(sdk.WarehouseResourceConstraintStandardGen1),
			Generation:                      sdk.Pointer(sdk.WarehouseGenerationStandardGen1),
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
			HasGeneration(sdk.WarehouseGenerationStandardGen1).
			HasNoResourceConstraint())

		warehouse, err := client.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), warehouse.Name)
		assert.Equal(t, sdk.WarehouseTypeStandard, warehouse.Type)
		assert.Equal(t, sdk.WarehouseSizeSmall, warehouse.Size)
		assert.Equal(t, 8, warehouse.MaxClusterCount)
		assert.Equal(t, 2, warehouse.MinClusterCount)
		assert.Equal(t, sdk.ScalingPolicyEconomy, warehouse.ScalingPolicy)
		assert.Equal(t, 1000, warehouse.AutoSuspend)
		assert.Equal(t, true, warehouse.AutoResume)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateResuming, sdk.WarehouseStateStarted}, warehouse.State)
		assert.Equal(t, resourceMonitor.ID().Name(), warehouse.ResourceMonitor.Name())
		assert.Equal(t, "comment", warehouse.Comment)
		assert.Equal(t, true, warehouse.EnableQueryAcceleration)
		assert.Equal(t, 90, warehouse.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouse.ResourceConstraint)
		assert.NotNil(t, warehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *warehouse.Generation)

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
			HasGeneration(sdk.WarehouseGenerationStandardGen1))

		tag1Value, err := client.SystemFunctions.GetTag(ctx, tag.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer("v1"), tag1Value)
		tag2Value, err := client.SystemFunctions.GetTag(ctx, tag2.ID(), warehouse.ID(), sdk.ObjectTypeWarehouse)
		require.NoError(t, err)
		assert.Equal(t, sdk.Pointer("v2"), tag2Value)
	})

	t.Run("create: with generation - BCR 2025_07", func(t *testing.T) {
		secondaryTestClientHelper().BcrBundles.EnableBcrBundle(t, "2025_07")
		id := secondaryTestClientHelper().Ids.RandomAccountObjectIdentifier()

		err := secondaryClient.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{
			Generation: sdk.Pointer(sdk.WarehouseGenerationStandardGen1),
		})
		require.NoError(t, err)

		result, err := secondaryClient.Warehouses.ShowByID(ctx, id)
		require.NoError(t, err)
		assertThatObject(t, objectassert.WarehouseFromObject(t, result).
			HasGeneration(sdk.WarehouseGenerationStandardGen1).
			HasNoResourceConstraint(),
		)
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
		assert.Equal(t, sdk.WarehouseSizeXSmall, result.Size)
		assert.Equal(t, 1, result.MaxClusterCount)
		assert.Equal(t, 1, result.MinClusterCount)
		assert.Equal(t, sdk.ScalingPolicyStandard, result.ScalingPolicy)
		assert.Equal(t, 600, result.AutoSuspend)
		assert.Equal(t, true, result.AutoResume)
		assert.Contains(t, []sdk.WarehouseState{sdk.WarehouseStateResuming, sdk.WarehouseStateStarted}, result.State)
		assert.Equal(t, "", result.Comment)
		assert.Equal(t, false, result.EnableQueryAcceleration)
		assert.Equal(t, 8, result.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, result.ResourceConstraint)
		assert.NotNil(t, result.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *result.Generation)
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

	t.Run("alter: set and unset", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		assert.Equal(t, sdk.WarehouseSizeXSmall, warehouse.Size)
		assert.Equal(t, 1, warehouse.MaxClusterCount)
		assert.Equal(t, 1, warehouse.MinClusterCount)
		assert.Equal(t, sdk.ScalingPolicyStandard, warehouse.ScalingPolicy)
		assert.Equal(t, 600, warehouse.AutoSuspend)
		assert.Equal(t, true, warehouse.AutoResume)
		assert.Equal(t, "", warehouse.ResourceMonitor.Name())
		assert.Equal(t, "", warehouse.Comment)
		assert.Equal(t, false, warehouse.EnableQueryAcceleration)
		assert.Equal(t, 8, warehouse.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouse.ResourceConstraint)
		assert.NotNil(t, warehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *warehouse.Generation)

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
		assert.Equal(t, sdk.WarehouseSizeMedium, warehouseAfterSet.Size)
		assert.Equal(t, 3, warehouseAfterSet.MaxClusterCount)
		assert.Equal(t, 2, warehouseAfterSet.MinClusterCount)
		assert.Equal(t, sdk.ScalingPolicyEconomy, warehouseAfterSet.ScalingPolicy)
		assert.Equal(t, 1234, warehouseAfterSet.AutoSuspend)
		assert.Equal(t, false, warehouseAfterSet.AutoResume)
		assert.Equal(t, resourceMonitor.ID().Name(), warehouseAfterSet.ResourceMonitor.Name())
		assert.Equal(t, "new comment", warehouseAfterSet.Comment)
		assert.Equal(t, true, warehouseAfterSet.EnableQueryAcceleration)
		assert.Equal(t, 2, warehouseAfterSet.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouseAfterSet.ResourceConstraint)
		assert.NotNil(t, warehouseAfterSet.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *warehouseAfterSet.Generation)

		alterOptions = &sdk.AlterWarehouseOptions{
			// WarehouseSize omitted on purpose - UNSET is not supported for warehouse size
			// WarehouseType, ScalingPolicy, AutoSuspend, and AutoResume omitted on purpose - UNSET do not work correctly
			// WaitForCompletion omitted on purpose - no unset
			Unset: &sdk.WarehouseUnset{
				MaxClusterCount:                 sdk.Bool(true),
				MinClusterCount:                 sdk.Bool(true),
				ResourceMonitor:                 sdk.Bool(true),
				Comment:                         sdk.Bool(true),
				EnableQueryAcceleration:         sdk.Bool(true),
				QueryAccelerationMaxScaleFactor: sdk.Bool(true),
			},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		warehouseAfterUnset, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, 1, warehouseAfterUnset.MaxClusterCount)
		assert.Equal(t, 1, warehouseAfterUnset.MinClusterCount)
		assert.Equal(t, "", warehouseAfterUnset.ResourceMonitor.Name())
		assert.Equal(t, "", warehouseAfterUnset.Comment)
		assert.Equal(t, false, warehouseAfterUnset.EnableQueryAcceleration)
		assert.Equal(t, 8, warehouseAfterUnset.QueryAccelerationMaxScaleFactor)
		assert.Nil(t, warehouseAfterUnset.ResourceConstraint)
		assert.NotNil(t, warehouseAfterUnset.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *warehouseAfterUnset.Generation)
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

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
		require.Eventually(t, func() bool { return sdk.WarehouseStateStarted == returnedWarehouse.State }, 5*time.Second, time.Second)

		alterOptions := &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{WarehouseType: sdk.Pointer(sdk.WarehouseTypeSnowparkOptimized)},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.Equal(t, sdk.WarehouseTypeSnowparkOptimized, returnedWarehouse.Type)
		assert.Contains(t, []any{sdk.WarehouseStateStarted, sdk.WarehouseStateResuming}, returnedWarehouse.State)

		// TODO [SNOW-1473453]: uncomment and test when UNSET starts working correctly (expecting to unset to default type STANDARD)
		// alterOptions = &sdk.AlterWarehouseOptions{
		//	Unset: &sdk.WarehouseUnset{WarehouseType: sdk.Bool(true)},
		// }
		// err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		// require.NoError(t, err)
		//
		// returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		// require.NoError(t, err)
		// assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
		// assert.Equal(t, sdk.WarehouseStateStarted, returnedWarehouse.State)
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

		// TODO [SNOW-1473453]: uncomment and test when UNSET starts working correctly (expecting to unset to default type STANDARD)
		// alterOptions = &sdk.AlterWarehouseOptions{
		//	Unset: &sdk.WarehouseUnset{WarehouseType: sdk.Bool(true)},
		// }
		// err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		// require.NoError(t, err)
		//
		// returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		// require.NoError(t, err)
		// assert.Equal(t, sdk.WarehouseTypeStandard, returnedWarehouse.Type)
		// assert.Equal(t, sdk.WarehouseStateStarted, returnedWarehouse.State)
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

	setAndUnsetGenerationOld := func(t *testing.T, client *sdk.Client, testClientHelper func() *helpers.TestClient) {
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
			Set: &sdk.WarehouseSet{ResourceConstraint: sdk.Pointer(sdk.WarehouseResourceConstraintStandardGen2)},
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
			Unset: &sdk.WarehouseUnset{ResourceConstraint: sdk.Bool(true)},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.NotNil(t, returnedWarehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *returnedWarehouse.Generation)
	}

	t.Run("alter: set and unset generation (old resource constraint syntax)", func(t *testing.T) {
		setAndUnsetGenerationOld(t, client, testClientHelper)
	})

	t.Run("alter: set and unset generation (old resource constraint syntax) - BCR 2025_07", func(t *testing.T) {
		secondaryTestClientHelper().BcrBundles.EnableBcrBundle(t, "2025_07")
		setAndUnsetGenerationOld(t, secondaryClient, secondaryTestClientHelper)
	})

	setAndUnsetGeneration := func(t *testing.T, client *sdk.Client, testClientHelper func() *helpers.TestClient) {
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
	}

	t.Run("alter: set and unset generation (new generation syntax)", func(t *testing.T) {
		setAndUnsetGeneration(t, client, testClientHelper)
	})

	t.Run("alter: set and unset generation (new generation syntax) - BCR 2025_07", func(t *testing.T) {
		testClientHelper().BcrBundles.EnableBcrBundle(t, "2025_07")
		setAndUnsetGeneration(t, secondaryClient, secondaryTestClientHelper)
	})

	t.Run("alter: set and unset generation (both at the same time)", func(t *testing.T) {
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
			Set: &sdk.WarehouseSet{
				Generation:         sdk.Pointer(sdk.WarehouseGenerationStandardGen2),
				ResourceConstraint: sdk.Pointer(sdk.WarehouseResourceConstraintStandardGen2),
			},
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
			Unset: &sdk.WarehouseUnset{
				Generation:         sdk.Bool(true),
				ResourceConstraint: sdk.Bool(true),
			},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err = client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		assert.NotNil(t, returnedWarehouse.Generation)
		assert.Equal(t, sdk.WarehouseGenerationStandardGen1, *returnedWarehouse.Generation)

		// setting incompatible values should fail
		alterOptions = &sdk.AlterWarehouseOptions{
			Set: &sdk.WarehouseSet{
				Generation:         sdk.Pointer(sdk.WarehouseGenerationStandardGen2),
				ResourceConstraint: sdk.Pointer(sdk.WarehouseResourceConstraintStandardGen1),
			},
		}
		err = client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.ErrorContains(t, err, "invalid property combination 'RESOURCE_CONSTRAINT'='STANDARD_GEN_1' and 'GENERATION'='2'")
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
		// assert.Equal(t, "600", returnedWarehouse.AutoSuspend)
		assert.Equal(t, 0, returnedWarehouse.AutoSuspend)
	})

	t.Run("alter: prove problems with unset warehouse type", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{WarehouseType: sdk.Bool(true)},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		// TODO [SNOW-1473453]: change when UNSET starts working correctly (expecting to unset to default type STANDARD)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid type of property 'null' for 'WAREHOUSE_TYPE'")
	})

	t.Run("alter: prove problems with unset scaling policy", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{ScalingPolicy: sdk.Bool(true)},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		// TODO [SNOW-1473453]: change when UNSET starts working correctly (expecting to unset to default scaling policy STANDARD)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid type of property 'null' for 'SCALING_POLICY'")
	})

	t.Run("alter: prove problems with unset auto resume", func(t *testing.T) {
		// new warehouse created on purpose
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)

		alterOptions := &sdk.AlterWarehouseOptions{
			Unset: &sdk.WarehouseUnset{AutoResume: sdk.Bool(true)},
		}
		err := client.Warehouses.Alter(ctx, warehouse.ID(), alterOptions)
		require.NoError(t, err)

		returnedWarehouse, err := client.Warehouses.ShowByID(ctx, warehouse.ID())
		require.NoError(t, err)
		// TODO [SNOW-1473453]: change when UNSET starts working correctly (expecting to unset to default auto resume TRUE)
		// assert.Equal(t, true, returnedWarehouse.AutoResume)
		assert.Equal(t, false, returnedWarehouse.AutoResume)
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
		assert.Equal(t, 1, result.Running)
		assert.Equal(t, 0, result.Queued)

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
		assert.Equal(t, 0, result.Running)
		assert.Equal(t, 0, result.Queued)
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
		assert.Equal(t, 1, result.Running)
		assert.Equal(t, 0, result.Queued)

		warehouseSuspended := func() bool {
			r, e := client.Warehouses.ShowByID(ctx, warehouse.ID())
			return e == nil && r.State == sdk.WarehouseStateSuspended
		}
		assert.Eventually(t, warehouseSuspended, 10*time.Second, time.Second)
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
		assert.Equal(t, sdk.WarehouseSizeMedium, result.Size)
		assert.Equal(t, 1, result.Running)
		assert.Equal(t, 0, result.Queued)
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
