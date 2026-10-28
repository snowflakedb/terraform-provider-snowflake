//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Alerts(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	warehouse := testClientHelper().Ids.WarehouseId()
	currentRole := testClientHelper().Context.CurrentRole(t)

	assertAlert := func(t *testing.T, id sdk.SchemaObjectIdentifier, warehouseName, schedule, comment, condition, action string, state sdk.AlertState) {
		t.Helper()

		assertThatObject(
			t, objectassert.Alert(t, id).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner(currentRole.Name()).
				HasOwnerRoleType("ROLE").
				HasWarehouse(warehouseName).
				HasSchedule(schedule).
				HasComment(comment).
				HasCondition(condition).
				HasAction(action).
				HasState(state),
		)
		assertThatObject(
			t, objectassert.AlertDetails(t, id).
				HasCreatedOnNotEmpty().
				HasName(id.Name()).
				HasDatabaseName(id.DatabaseName()).
				HasSchemaName(id.SchemaName()).
				HasOwner(currentRole.Name()).
				HasWarehouse(warehouseName).
				HasSchedule(schedule).
				HasComment(comment).
				HasCondition(condition).
				HasAction(action).
				HasState(string(state)),
		)
	}

	t.Run("create: no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.Alerts.Create(ctx, sdk.NewCreateAlertRequest(
			id,
			warehouse,
			"5 MINUTE",
			sdk.NewAlertConditionFromString("SELECT 1"),
			"SELECT 1",
		))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Alert.DropAlertFunc(t, id))

		assertAlert(t, id, warehouse.Name(), "5 MINUTE", "", "SELECT 1", "SELECT 1", sdk.AlertStateSuspended)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()
		action := "SELECT\n  CASE WHEN true THEN 1 ELSE 2 END"

		err := client.Alerts.Create(ctx, sdk.NewCreateAlertRequest(
			id,
			warehouse,
			"USING CRON * * * * TUE,THU UTC",
			sdk.NewAlertConditionFromString("SELECT 1"),
			action,
		).
			WithIfNotExists(true).
			WithComment(comment))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Alert.DropAlertFunc(t, id))

		assertAlert(t, id, warehouse.Name(), "USING CRON * * * * TUE,THU UTC", comment, "SELECT 1", action, sdk.AlertStateSuspended)
	})

	t.Run("alter: set and unset", func(t *testing.T) {
		alert, cleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(cleanup)

		otherWarehouse, warehouseCleanup := testClientHelper().Warehouse.CreateWarehouse(t)
		t.Cleanup(warehouseCleanup)
		id := alert.ID()

		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).
			WithSet(
				*sdk.NewAlertSetRequest().
					WithWarehouse(otherWarehouse.ID()).
					WithSchedule("10 MINUTE").
					WithComment("changed"),
			))
		require.NoError(t, err)

		assertAlert(t, id, otherWarehouse.Name, "10 MINUTE", "changed", "SELECT 1", "SELECT 1", sdk.AlertStateSuspended)

		// Snowflake rejects UNSET WAREHOUSE COMMENT as a syntax error, so each property is unset on its own.
		// Schedule is omitted: unsetting it would convert the alert into a streaming alert.
		err = client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).
			WithUnset(
				*sdk.NewAlertUnsetRequest().
					WithWarehouse(true),
			))
		require.NoError(t, err)

		err = client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(id).
			WithUnset(
				*sdk.NewAlertUnsetRequest().
					WithComment(true),
			))
		require.NoError(t, err)

		assertAlert(t, id, "", "10 MINUTE", "", "SELECT 1", "SELECT 1", sdk.AlertStateSuspended)
	})

	t.Run("alter: modify condition and action", func(t *testing.T) {
		alert, cleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(cleanup)

		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(alert.ID()).
			WithModifyCondition([]string{"SELECT 2 WHERE false"}))
		require.NoError(t, err)

		err = client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(alert.ID()).
			WithModifyAction("SELECT 3"))
		require.NoError(t, err)

		assertAlert(t, alert.ID(), alert.Warehouse, alert.Schedule, "", "SELECT 2 WHERE false", "SELECT 3", sdk.AlertStateSuspended)
	})

	t.Run("alter: resume and suspend", func(t *testing.T) {
		alert, cleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(cleanup)

		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(alert.ID()).
			WithAction(sdk.AlertActionResume))
		require.NoError(t, err)

		assertAlert(t, alert.ID(), alert.Warehouse, alert.Schedule, "", alert.Condition, alert.Action, sdk.AlertStateStarted)

		err = client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(alert.ID()).
			WithAction(sdk.AlertActionSuspend))
		require.NoError(t, err)

		assertAlert(t, alert.ID(), alert.Warehouse, alert.Schedule, "", alert.Condition, alert.Action, sdk.AlertStateSuspended)
	})

	t.Run("describe", func(t *testing.T) {
		alert, cleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(cleanup)

		assertAlert(t, alert.ID(), alert.Warehouse, alert.Schedule, "", alert.Condition, alert.Action, sdk.AlertStateSuspended)
	})

	t.Run("show: with like", func(t *testing.T) {
		alert, cleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(cleanup)

		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(alert.Name)}))
		require.NoError(t, err)
		require.Len(t, alerts, 1)

		assertThatObject(
			t, objectassert.AlertFromObject(t, &alerts[0]).
				HasName(alert.Name).
				HasDatabaseName(alert.DatabaseName).
				HasSchemaName(alert.SchemaName).
				HasSchedule(alert.Schedule).
				HasState(sdk.AlertStateSuspended).
				HasCondition(alert.Condition).
				HasAction(alert.Action).
				HasOwnerRoleType("ROLE"),
		)
	})

	t.Run("show: in account, database and schema", func(t *testing.T) {
		alert, cleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(cleanup)
		id := alert.ID()

		names := func(alerts []sdk.Alert) []string {
			return collections.Map(alerts, func(alert sdk.Alert) string { return alert.Name })
		}

		inAccount, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithIn(sdk.In{Account: sdk.Bool(true)}))
		require.NoError(t, err)
		assert.Contains(t, names(inAccount), alert.Name)

		inDatabase, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithIn(sdk.In{Database: id.DatabaseId()}))
		require.NoError(t, err)
		assert.Contains(t, names(inDatabase), alert.Name)

		inSchema, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithIn(sdk.In{Schema: id.SchemaId()}))
		require.NoError(t, err)
		assert.Contains(t, names(inSchema), alert.Name)
	})

	t.Run("show: with starts with and limit", func(t *testing.T) {
		alert, cleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(cleanup)

		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithStartsWith(alert.Name).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		require.Len(t, alerts, 1)
		assert.Equal(t, alert.Name, alerts[0].Name)

		all, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithIn(sdk.In{Schema: alert.ID().SchemaId()}))
		require.NoError(t, err)
		require.NotEmpty(t, all)

		limited, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithIn(sdk.In{Schema: alert.ID().SchemaId()}).
			WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}))
		require.NoError(t, err)
		assert.Len(t, limited, 1)
		assert.LessOrEqual(t, len(limited), len(all))
	})

	t.Run("show: terse", func(t *testing.T) {
		alert, cleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(cleanup)

		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithTerse(true).
			WithLike(sdk.Like{Pattern: sdk.String(alert.Name)}))
		require.NoError(t, err)
		require.Len(t, alerts, 1)

		assertThatObject(
			t, objectassert.AlertFromObject(t, &alerts[0]).
				HasCreatedOnNotEmpty().
				HasName(alert.Name).
				HasDatabaseName(alert.DatabaseName).
				HasSchemaName(alert.SchemaName),
		)
	})

	t.Run("show by id: same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		for _, id := range []sdk.SchemaObjectIdentifier{id1, id2} {
			created, cleanup := testClientHelper().Alert.CreateAlertWithRequest(t, sdk.NewCreateAlertRequest(
				id,
				warehouse,
				"5 MINUTE",
				sdk.NewAlertConditionFromString("SELECT 1"),
				"SELECT 1",
			))
			t.Cleanup(cleanup)

			assert.Equal(t, id, created.ID())
			assertAlert(t, id, created.Warehouse, "5 MINUTE", "", "SELECT 1", "SELECT 1", sdk.AlertStateSuspended)
		}
	})

	// LIKE wildcards must not select a different alert. "prefix1" contains "prefix_" once "_" is a wildcard.
	t.Run("show by id: like wildcards do not select another alert", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)
		suffix := testClientHelper().Ids.Alpha()

		for _, name := range []string{"prefix1" + suffix, "prefix_" + suffix} {
			id := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(name, schema.ID())
			_, cleanup := testClientHelper().Alert.CreateAlertWithRequest(t, sdk.NewCreateAlertRequest(
				id,
				warehouse,
				"5 MINUTE",
				sdk.NewAlertConditionFromString("SELECT 1"),
				"SELECT 1",
			))
			t.Cleanup(cleanup)

			found, err := client.Alerts.ShowByID(ctx, id)
			require.NoError(t, err)
			require.Equal(t, id, found.ID())
		}
	})

	t.Run("drop: existing", func(t *testing.T) {
		alert, cleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(cleanup)

		err := client.Alerts.Drop(ctx, sdk.NewDropAlertRequest(alert.ID()))
		require.NoError(t, err)

		_, err = client.Alerts.ShowByID(ctx, alert.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: when an alert does not exist", func(t *testing.T) {
		err := client.Alerts.Drop(ctx, sdk.NewDropAlertRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop: if exists on a non-existing alert", func(t *testing.T) {
		err := client.Alerts.Drop(ctx, sdk.NewDropAlertRequest(NonExistingSchemaObjectIdentifier).
			WithIfExists(true))
		require.NoError(t, err)
	})
}
