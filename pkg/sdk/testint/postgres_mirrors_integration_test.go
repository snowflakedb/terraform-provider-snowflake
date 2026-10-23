//go:build account_level_tests

package testint

import (
	"database/sql"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_PostgresMirrors(t *testing.T) {
	instance := testClientHelper().PostgresMirror.FixtureInstance(t)
	client := testClient(t)
	ctx := testContext(t)

	newCreateArguments := func() (*sdk.CreatePostgresMirrorArgumentsRequest, sdk.AccountObjectIdentifier) {
		targetDatabase := testClientHelper().Ids.RandomAccountObjectIdentifier()
		arguments := sdk.NewCreatePostgresMirrorArgumentsRequest(
			testClientHelper().Ids.RandomPostgresMirrorName(),
			instance,
			helpers.PostgresMirrorFixtureDatabase,
			targetDatabase.Name(),
		)
		return arguments, targetDatabase
	}

	describe := func(t *testing.T, name string) *sdk.PostgresMirrorDetails {
		t.Helper()
		details, err := testClientHelper().PostgresMirror.Describe(t, name)
		require.NoError(t, err)
		return details
	}

	t.Run("create: tables", func(t *testing.T) {
		arguments, targetDatabase := newCreateArguments()
		request := sdk.NewCreatePostgresMirrorRequest(*arguments.
			WithPostgresTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable1, helpers.PostgresMirrorFixtureMixedCaseTable)))

		details := testClientHelper().PostgresMirror.CreateWithRequest(t, request)

		assertThatObject(t, objectassert.PostgresMirrorDetailsFromObject(t, details).
			HasName(request.GetName()).
			HasPostgresInstance(instance).
			HasPostgresDatabase(helpers.PostgresMirrorFixtureDatabase).
			HasTargetDatabase(targetDatabase.Name()).
			HasRefreshInterval(sdk.PostgresMirrorDefaultRefreshInterval).
			HasTargetTableType(sdk.PostgresMirrorTargetTableTypeSnowflake).
			HasNoExternalVolume().
			HasNoWarehouse().
			// Returned sorted, in Postgres quoting.
			HasPostgresTables(helpers.PostgresMirrorFixtureMixedCaseTable, helpers.PostgresMirrorFixtureTable1).
			HasNoPostgresSchemas(),
		)
		assert.True(t, details.IsMembershipKnown())
		assert.Contains(t, []sdk.PostgresMirrorStatus{sdk.PostgresMirrorStatusStarting, sdk.PostgresMirrorStatusActive}, details.Status)

		_, err := testClientHelper().Database.Show(t, targetDatabase)
		require.NoError(t, err)
	})

	t.Run("create: tables, schemas, refresh interval, target type, and warehouse", func(t *testing.T) {
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)
		t.Cleanup(testClientHelper().Grant.GrantUsageOnWarehouseToSnowflakeApplication(t, warehouse.ID()))

		arguments, targetDatabase := newCreateArguments()
		request := sdk.NewCreatePostgresMirrorRequest(*arguments.
			WithPostgresTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable2)).
			WithPostgresSchemas(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureSchemaA)).
			WithRefreshInterval("1 minute").
			WithTargetTableType(sdk.PostgresMirrorTargetTableTypeIceberg).
			WithWarehouse(warehouse.ID().Name()))

		details := testClientHelper().PostgresMirror.CreateWithRequest(t, request)

		assertThatObject(t, objectassert.PostgresMirrorDetailsFromObject(t, details).
			HasName(request.GetName()).
			HasTargetDatabase(targetDatabase.Name()).
			HasRefreshInterval("1 minute").
			HasTargetTableType(sdk.PostgresMirrorTargetTableTypeIceberg).
			HasNoExternalVolume().
			HasWarehouse(warehouse.ID().Name()).
			HasPostgresTables(helpers.PostgresMirrorFixtureTable2).
			HasPostgresSchemas(helpers.PostgresMirrorFixtureSchemaA),
		)
	})

	t.Run("create: mirror name is case-insensitive", func(t *testing.T) {
		name := testClientHelper().Ids.RandomPostgresMirrorName()
		request := sdk.NewCreatePostgresMirrorRequest(*sdk.NewCreatePostgresMirrorArgumentsRequest(strings.ToUpper(name), instance, helpers.PostgresMirrorFixtureDatabase, testClientHelper().Ids.RandomAccountObjectIdentifier().Name()).
			WithPostgresTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable1)))

		err := client.PostgresMirrors.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresMirror.DropFunc(t, name))

		assertThatObject(t, objectassert.PostgresMirrorDetailsFromObject(t, describe(t, strings.ToUpper(name))).
			HasName(name),
		)
	})

	t.Run("create: target database already exists", func(t *testing.T) {
		database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(databaseCleanup)

		request := sdk.NewCreatePostgresMirrorRequest(*sdk.NewCreatePostgresMirrorArgumentsRequest(testClientHelper().Ids.RandomPostgresMirrorName(), instance, helpers.PostgresMirrorFixtureDatabase, database.ID().Name()).
			WithPostgresTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable1)))
		t.Cleanup(testClientHelper().PostgresMirror.DropFunc(t, request.GetName()))

		// CREATE_MIRROR only sees databases the SNOWFLAKE application has privileges on, so it misses a database owned
		// by the caller. The error then comes from its CREATE DATABASE, not from its [ALREADY_EXISTS] check.
		err := client.PostgresMirrors.Create(ctx, request)
		require.ErrorContains(t, err, "already exists")
	})

	t.Run("alter: tables, schemas, and refresh interval", func(t *testing.T) {
		details := testClientHelper().PostgresMirror.CreateWithTables(t, instance, helpers.PostgresMirrorFixtureTable1, helpers.PostgresMirrorFixtureTable2)

		testClientHelper().PostgresMirror.Alter(t, sdk.NewAlterPostgresMirrorRequest(*sdk.NewAlterPostgresMirrorArgumentsRequest(details.Name).
			WithAddTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable3)).
			WithRemoveTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable1)).
			WithAddSchemas(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureSchemaB)).
			WithRefreshInterval("30 seconds"),
		))
		updated := describe(t, details.Name)
		assertThatObject(t, objectassert.PostgresMirrorDetailsFromObject(t, updated).
			HasRefreshInterval("30 seconds").
			HasPostgresTables(helpers.PostgresMirrorFixtureTable2, helpers.PostgresMirrorFixtureTable3).
			HasPostgresSchemas(helpers.PostgresMirrorFixtureSchemaB),
		)
		assert.True(t, updated.IsMembershipKnown())

		testClientHelper().PostgresMirror.Alter(t, sdk.NewAlterPostgresMirrorRequest(*sdk.NewAlterPostgresMirrorArgumentsRequest(details.Name).
			WithRemoveSchemas(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureSchemaB)),
		))
		assertThatObject(t, objectassert.PostgresMirrorDetailsFromObject(t, describe(t, details.Name)).
			HasRefreshInterval("30 seconds").
			HasPostgresTables(helpers.PostgresMirrorFixtureTable2, helpers.PostgresMirrorFixtureTable3).
			HasNoPostgresSchemas(),
		)
	})

	t.Run("alter: set and unset warehouse", func(t *testing.T) {
		warehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)
		t.Cleanup(testClientHelper().Grant.GrantUsageOnWarehouseToSnowflakeApplication(t, warehouse.ID()))

		details := testClientHelper().PostgresMirror.CreateWithTables(t, instance, helpers.PostgresMirrorFixtureTable1)

		testClientHelper().PostgresMirror.Alter(t, sdk.NewAlterPostgresMirrorRequest(*sdk.NewAlterPostgresMirrorArgumentsRequest(details.Name).WithWarehouse(warehouse.ID().Name())))
		assertThatObject(t, objectassert.PostgresMirrorDetailsFromObject(t, describe(t, details.Name)).
			HasWarehouse(warehouse.ID().Name()),
		)

		testClientHelper().PostgresMirror.Alter(t, sdk.NewAlterPostgresMirrorRequest(*sdk.NewAlterPostgresMirrorArgumentsRequest(details.Name).WithUnsetWarehouse(true)))
		assertThatObject(t, objectassert.PostgresMirrorDetailsFromObject(t, describe(t, details.Name)).
			HasNoWarehouse(),
		)
	})

	t.Run("list", func(t *testing.T) {
		details := testClientHelper().PostgresMirror.CreateWithTables(t, instance, helpers.PostgresMirrorFixtureTable1)

		mirrors, err := client.PostgresMirrors.List(ctx, sdk.NewListPostgresMirrorRequest(*sdk.NewListPostgresMirrorArgumentsRequest().WithPostgresInstance(instance)))
		require.NoError(t, err)
		mirror, err := collections.FindFirst(mirrors, func(mirror sdk.PostgresMirror) bool { return mirror.Name == details.Name })
		require.NoError(t, err)
		assert.Equal(t, instance, mirror.PostgresInstance)
		assert.Equal(t, helpers.PostgresMirrorFixtureDatabase, mirror.PostgresDatabase)
		assert.Equal(t, details.TargetDatabase, mirror.TargetDatabase)
		assert.Equal(t, sdk.PostgresMirrorDefaultRefreshInterval, mirror.RefreshInterval)
		assert.NotNil(t, mirror.TableCount)
		assert.NotNil(t, mirror.InitialSnapshotTablesRemaining)

		allMirrors, err := client.PostgresMirrors.List(ctx, sdk.NewListPostgresMirrorRequest(*sdk.NewListPostgresMirrorArgumentsRequest()))
		require.NoError(t, err)
		assert.True(t, slices.ContainsFunc(allMirrors, func(mirror sdk.PostgresMirror) bool { return mirror.Name == details.Name }))
	})

	t.Run("wait for the initial snapshot", func(t *testing.T) {
		details := testClientHelper().PostgresMirror.CreateWithTables(t, instance, helpers.PostgresMirrorFixtureTable1, helpers.PostgresMirrorFixtureTable2)

		testClientHelper().PostgresMirror.WaitForInitialSnapshot(t, instance, details.Name, 2, 15*time.Minute)

		assertThatObject(t, objectassert.PostgresMirrorDetailsFromObject(t, describe(t, details.Name)).
			HasStatus(sdk.PostgresMirrorStatusActive),
		)
	})

	t.Run("describe: missing mirror", func(t *testing.T) {
		_, err := client.PostgresMirrors.Describe(ctx, sdk.NewDescribePostgresMirrorRequest(*sdk.NewDescribePostgresMirrorArgumentsRequest(testClientHelper().Ids.RandomPostgresMirrorName())))
		require.ErrorIs(t, err, sql.ErrNoRows)
	})

	t.Run("drop: keep the target database", func(t *testing.T) {
		arguments, targetDatabase := newCreateArguments()
		request := sdk.NewCreatePostgresMirrorRequest(*arguments.WithPostgresTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable1)))
		details := testClientHelper().PostgresMirror.CreateWithRequest(t, request)

		err := client.PostgresMirrors.Drop(ctx, sdk.NewDropPostgresMirrorRequest(*sdk.NewDropPostgresMirrorArgumentsRequest(details.Name).WithDropTargetDatabase(false)))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().PostgresMirror.DropKeptTargetDatabaseFunc(t, targetDatabase))

		_, err = testClientHelper().PostgresMirror.Describe(t, details.Name)
		require.ErrorIs(t, err, sql.ErrNoRows)
		_, err = testClientHelper().Database.Show(t, targetDatabase)
		require.NoError(t, err)

		// A kept target database is marked as previously mirrored and can never be a target again.
		recreateRequest := sdk.NewCreatePostgresMirrorRequest(*sdk.NewCreatePostgresMirrorArgumentsRequest(testClientHelper().Ids.RandomPostgresMirrorName(), instance, helpers.PostgresMirrorFixtureDatabase, targetDatabase.Name()).
			WithPostgresTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable1)))
		t.Cleanup(testClientHelper().PostgresMirror.DropFunc(t, recreateRequest.GetName()))
		err = client.PostgresMirrors.Create(ctx, recreateRequest)
		require.ErrorContains(t, err, "[ALREADY_EXISTS]")
	})

	t.Run("drop: target database is dropped by default and its name can be reused", func(t *testing.T) {
		arguments, targetDatabase := newCreateArguments()
		request := sdk.NewCreatePostgresMirrorRequest(*arguments.WithPostgresTables(helpers.PostgresMirrorNames(helpers.PostgresMirrorFixtureTable1)))
		details := testClientHelper().PostgresMirror.CreateWithRequest(t, request)

		err := client.PostgresMirrors.Drop(ctx, sdk.NewDropPostgresMirrorRequest(*sdk.NewDropPostgresMirrorArgumentsRequest(details.Name)))
		require.NoError(t, err)

		_, err = testClientHelper().Database.Show(t, targetDatabase)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)

		testClientHelper().PostgresMirror.CreateWithRequest(t, request)
	})

	t.Run("drop: missing mirror", func(t *testing.T) {
		request := sdk.NewDropPostgresMirrorRequest(*sdk.NewDropPostgresMirrorArgumentsRequest(testClientHelper().Ids.RandomPostgresMirrorName()))

		err := client.PostgresMirrors.Drop(ctx, request)
		require.ErrorContains(t, err, "[NOT_FOUND] Mirror does not exist")

		err = client.PostgresMirrors.DropSafely(ctx, request)
		require.NoError(t, err)
	})
}
