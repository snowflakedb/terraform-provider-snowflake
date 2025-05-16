//go:build !account_level_tests

package testint

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/provider/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DynamicTableCreateAndDrop(t *testing.T) {
	client := testClient(t)

	tableTest, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	ctx := context.Background()
	t.Run("test complete", func(t *testing.T) {
		name := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		targetLag := sdk.TargetLag{
			MaximumDuration: sdk.String("2 minutes"),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := random.Comment()
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(name, testClientHelper().Ids.WarehouseId(), targetLag, query).WithOrReplace(true).WithComment(&comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(name))
			require.NoError(t, err)
		})
		entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(name.Name())}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, name.Name(), entity.Name)
		require.Equal(t, testClientHelper().Ids.WarehouseId().Name(), entity.Warehouse)
		require.Equal(t, *targetLag.MaximumDuration, entity.TargetLag)

		dynamicTableById, err := client.DynamicTables.ShowByID(ctx, name)
		require.NoError(t, err)
		require.NotNil(t, dynamicTableById)
		require.Equal(t, name.Name(), dynamicTableById.Name)
		require.Equal(t, testClientHelper().Ids.WarehouseId().Name(), dynamicTableById.Warehouse)
		require.Equal(t, *targetLag.MaximumDuration, dynamicTableById.TargetLag)
		assert.Equal(t, "ROLE", dynamicTableById.OwnerRoleType)
	})

	t.Run("create with usage tracking comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		plainQuery := fmt.Sprintf("SELECT id FROM %s", tableTest.ID().FullyQualifiedName())
		query, err := tracking.AppendMetadata(plainQuery, tracking.NewVersionedResourceMetadata(resources.DynamicTable, tracking.CreateOperation))
		require.NoError(t, err)

		err = client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(id, testClientHelper().Ids.WarehouseId(), sdk.TargetLag{
			MaximumDuration: sdk.String("2 minutes"),
		}, query))
		require.NoError(t, err)

		dynamicTable, err := client.DynamicTables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, fmt.Sprintf("CREATE DYNAMIC TABLE %s lag = '2 minutes' refresh_mode = 'AUTO' initialize = 'ON_CREATE' warehouse = %s AS %s", id.FullyQualifiedName(), testClientHelper().Ids.WarehouseId().FullyQualifiedName(), plainQuery), dynamicTable.Text)
	})

	t.Run("test complete with target lag", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		targetLag := sdk.TargetLag{
			Downstream: sdk.Bool(true),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := random.Comment()
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(id, testClientHelper().Ids.WarehouseId(), targetLag, query).WithOrReplace(true).WithComment(&comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(id))
			require.NoError(t, err)
		})
		entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(id.Name())}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, id.Name(), entity.Name)
		require.Equal(t, testClientHelper().Ids.WarehouseId().Name(), entity.Warehouse)
		require.Equal(t, "DOWNSTREAM", entity.TargetLag)
		require.Equal(t, sdk.DynamicTableRefreshModeIncremental, entity.RefreshMode)
		require.Contains(t, entity.Text, "initialize = 'ON_CREATE'")
		require.Contains(t, entity.Text, "refresh_mode = 'AUTO'")
	})

	t.Run("test complete with refresh mode and initialize", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		targetLag := sdk.TargetLag{
			MaximumDuration: sdk.String("2 minutes"),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		comment := random.Comment()
		refreshMode := sdk.DynamicTableRefreshModeFull
		initialize := sdk.DynamicTableInitializeOnSchedule
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(id, testClientHelper().Ids.WarehouseId(), targetLag, query).WithOrReplace(true).WithInitialize(initialize).WithRefreshMode(refreshMode).WithComment(&comment))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(id))
			require.NoError(t, err)
		})
		entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(id.Name())}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))

		entity := entities[0]
		require.Equal(t, id.Name(), entity.Name)
		require.Equal(t, testClientHelper().Ids.WarehouseId().Name(), entity.Warehouse)
		require.Equal(t, *targetLag.MaximumDuration, entity.TargetLag)
		require.Equal(t, sdk.DynamicTableRefreshModeFull, entity.RefreshMode)
		require.Contains(t, entity.Text, "initialize = 'ON_SCHEDULE'")
		require.Contains(t, entity.Text, "refresh_mode = 'FULL'")
	})
}

func TestInt_DynamicTableDescribe(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	dynamicTable, dynamicTableCleanup := testClientHelper().DynamicTable.CreateDynamicTable(t, table.ID())
	t.Cleanup(dynamicTableCleanup)

	t.Run("when dynamic table exists", func(t *testing.T) {
		_, err := client.DynamicTables.Describe(ctx, sdk.NewDescribeDynamicTableRequest(dynamicTable.ID()))
		require.NoError(t, err)
	})

	t.Run("when dynamic table does not exist", func(t *testing.T) {
		_, err := client.DynamicTables.Describe(ctx, sdk.NewDescribeDynamicTableRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_DynamicTableAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	t.Run("alter with suspend or resume", func(t *testing.T) {
		dynamicTable, dynamicTableCleanup := testClientHelper().DynamicTable.CreateDynamicTable(t, table.ID())
		t.Cleanup(dynamicTableCleanup)

		entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))
		require.Equal(t, sdk.DynamicTableSchedulingStateActive, entities[0].SchedulingState)

		err = client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithSuspend(sdk.Bool(true)))
		require.NoError(t, err)

		entities, err = client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))
		require.Equal(t, sdk.DynamicTableSchedulingStateSuspended, entities[0].SchedulingState)

		err = client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithResume(sdk.Bool(true)))
		require.NoError(t, err)

		entities, err = client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))
		require.Equal(t, sdk.DynamicTableSchedulingStateActive, entities[0].SchedulingState)
	})

	t.Run("alter with refresh", func(t *testing.T) {
		dynamicTable, dynamicTableCleanup := testClientHelper().DynamicTable.CreateDynamicTable(t, table.ID())
		t.Cleanup(dynamicTableCleanup)

		err := client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithRefresh(sdk.Bool(true)))
		require.NoError(t, err)

		entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
		require.NoError(t, err)
		require.Equal(t, 1, len(entities))
	})

	t.Run("alter with suspend and resume", func(t *testing.T) {
		dynamicTable, dynamicTableCleanup := testClientHelper().DynamicTable.CreateDynamicTable(t, table.ID())
		t.Cleanup(dynamicTableCleanup)

		err := client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithSuspend(sdk.Bool(true)).WithResume(sdk.Bool(true)))
		require.Error(t, err)
		sdk.ErrorsEqual(t, sdk.JoinErrors(sdk.ErrExactlyOneOf("alterDynamicTableOptions", "Suspend", "Resume", "Refresh", "Set")), err)
	})

	t.Run("alter with set", func(t *testing.T) {
		dynamicTable, dynamicTableCleanup := testClientHelper().DynamicTable.CreateDynamicTable(t, table.ID())
		t.Cleanup(dynamicTableCleanup)

		targetLagCases := []string{"10 minutes", "DOWNSTREAM"}
		for _, value := range targetLagCases {
			err := client.DynamicTables.Alter(ctx, sdk.NewAlterDynamicTableRequest(dynamicTable.ID()).WithSet(sdk.NewDynamicTableSetRequest().WithTargetLag(sdk.TargetLag{MaximumDuration: sdk.String(value)})))
			require.NoError(t, err)
			entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(dynamicTable.Name)}))
			require.NoError(t, err)
			require.Equal(t, 1, len(entities))
			require.Equal(t, value, entities[0].TargetLag)
		}
	})
}

func TestInt_DynamicTablesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	warehouseId := testClientHelper().Ids.WarehouseId()

	cleanupDynamicTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createDynamicTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		tableTest, tableCleanup := testClientHelper().Table.Create(t)
		t.Cleanup(tableCleanup)
		targetLag := sdk.TargetLag{
			MaximumDuration: sdk.String("2 minutes"),
		}
		query := "select id from " + tableTest.ID().FullyQualifiedName()
		err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(id, warehouseId, targetLag, query).WithOrReplace(true))
		require.NoError(t, err)
		t.Cleanup(cleanupDynamicTableHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createDynamicTableHandle(t, id1)
		createDynamicTableHandle(t, id2)

		e1, err := client.DynamicTables.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.DynamicTables.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
