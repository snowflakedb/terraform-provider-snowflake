//go:build non_account_level_tests

package testint

import (
	"errors"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_AlertsShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	alertTest, alertCleanup := testClientHelper().Alert.CreateAlert(t)
	t.Cleanup(alertCleanup)

	alert2Test, alert2Cleanup := testClientHelper().Alert.CreateAlert(t)
	t.Cleanup(alert2Cleanup)

	t.Run("without show options", func(t *testing.T) {
		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest())
		require.NoError(t, err)
		assert.Len(t, alerts, 2)
	})

	t.Run("with show options", func(t *testing.T) {
		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Contains(t, alerts, *alertTest)
		assert.Contains(t, alerts, *alert2Test)
		assert.Len(t, alerts, 2)
	})

	t.Run("with show options and like", func(t *testing.T) {
		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(alertTest.Name)}).
			WithIn(sdk.In{Database: testClientHelper().Ids.DatabaseId()}))
		require.NoError(t, err)
		assert.Contains(t, alerts, *alertTest)
		assert.Len(t, alerts, 1)
	})

	t.Run("when searching a non-existent alert", func(t *testing.T) {
		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String("non-existent")}))
		require.NoError(t, err)
		assert.Empty(t, alerts)
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}).
			WithLimit(sdk.LimitFrom{Rows: new(1)}))
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
	})
}

func TestInt_AlertCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := "SELECT 1"
		comment := random.Comment()
		err := client.Alerts.Create(ctx, sdk.NewCreateAlertRequest(id, testClientHelper().Ids.WarehouseId(), schedule, sdk.NewAlertConditionFromString(condition), action).
			WithOrReplace(true).
			WithIfNotExists(false).
			WithComment(comment))
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testClientHelper().Ids.WarehouseId().Name(), alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, comment, *alertDetails.Comment)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(name)}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Len(t, alert, 1)
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, comment, *alert[0].Comment)
	})

	t.Run("test if_not_exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := "SELECT 1"
		comment := random.Comment()
		err := client.Alerts.Create(ctx, sdk.NewCreateAlertRequest(id, testClientHelper().Ids.WarehouseId(), schedule, sdk.NewAlertConditionFromString(condition), action).
			WithOrReplace(false).
			WithIfNotExists(true).
			WithComment(comment))
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testClientHelper().Ids.WarehouseId().Name(), alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, comment, *alertDetails.Comment)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(name)}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Len(t, alert, 1)
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, comment, *alert[0].Comment)
	})

	t.Run("test no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := "SELECT 1"
		err := client.Alerts.Create(ctx, sdk.NewCreateAlertRequest(id, testClientHelper().Ids.WarehouseId(), schedule, sdk.NewAlertConditionFromString(condition), action))
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testClientHelper().Ids.WarehouseId().Name(), alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(name)}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Len(t, alert, 1)
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, "", *alert[0].Comment)
	})

	t.Run("test multiline action", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		name := id.Name()
		schedule := "USING CRON * * * * TUE,THU UTC"
		condition := "SELECT 1"
		action := `
			select
				case
					when true then
						1
					else
						2
				end
		`
		err := client.Alerts.Create(ctx, sdk.NewCreateAlertRequest(id, testClientHelper().Ids.WarehouseId(), schedule, sdk.NewAlertConditionFromString(condition), action))
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testClientHelper().Ids.WarehouseId().Name(), alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, strings.TrimSpace(action), alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(name)}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Len(t, alert, 1)
		assert.Equal(t, name, alert[0].Name)
		assert.Equal(t, "", *alert[0].Comment)
	})
}

func TestInt_AlertDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	alert, alertCleanup := testClientHelper().Alert.CreateAlert(t)
	t.Cleanup(alertCleanup)

	t.Run("when alert exists", func(t *testing.T) {
		alertDetails, err := client.Alerts.Describe(ctx, alert.ID())
		require.NoError(t, err)
		assert.Equal(t, alert.Name, alertDetails.Name)
	})

	t.Run("when alert does not exist", func(t *testing.T) {
		_, err := client.Alerts.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_AlertAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when setting and unsetting a value", func(t *testing.T) {
		alert, alertCleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(alertCleanup)
		newSchedule := "USING CRON * * * * TUE,FRI GMT"

		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(alert.ID()).WithSet(*sdk.NewAlertSetRequest().WithSchedule(newSchedule)))
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(alert.Name)}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, newSchedule, alerts[0].Schedule)
	})

	t.Run("when modifying condition and action", func(t *testing.T) {
		alert, alertCleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(alertCleanup)
		newCondition := "select * from DUAL where false"

		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(alert.ID()).WithModifyCondition([]string{newCondition}))
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(alert.Name)}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, newCondition, alerts[0].Condition)

		newAction := "create table FOO(ID INT)"

		err = client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(alert.ID()).WithModifyAction(newAction))
		require.NoError(t, err)
		alerts, err = client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(alert.Name)}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, newAction, alerts[0].Action)
	})

	t.Run("resume and then suspend", func(t *testing.T) {
		alert, alertCleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(alertCleanup)

		err := client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(alert.ID()).WithAction(sdk.AlertActionResume))
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(alert.Name)}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, alerts[0].State, sdk.AlertStateStarted)

		err = client.Alerts.Alter(ctx, sdk.NewAlterAlertRequest(alert.ID()).WithAction(sdk.AlertActionSuspend))
		require.NoError(t, err)
		alerts, err = client.Alerts.Show(ctx, sdk.NewShowAlertRequest().
			WithLike(sdk.Like{Pattern: sdk.String(alert.Name)}).
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, alerts[0].State, sdk.AlertStateSuspended)
	})
}

func TestInt_AlertDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when alert exists", func(t *testing.T) {
		alert, _ := testClientHelper().Alert.CreateAlert(t)
		id := alert.ID()
		err := client.Alerts.Drop(ctx, sdk.NewDropAlertRequest(id))
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when alert does not exist", func(t *testing.T) {
		err := client.Alerts.Drop(ctx, sdk.NewDropAlertRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_AlertsShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	warehouseId := testClientHelper().Ids.WarehouseId()
	cleanupAlertHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.Alerts.Drop(ctx, sdk.NewDropAlertRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createAlertHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		schedule, condition, action := "USING CRON * * * * * UTC", "SELECT 1", "SELECT 1"
		err := client.Alerts.Create(ctx, sdk.NewCreateAlertRequest(id, warehouseId, schedule, sdk.NewAlertConditionFromString(condition), action))
		require.NoError(t, err)
		t.Cleanup(cleanupAlertHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createAlertHandle(t, id1)
		createAlertHandle(t, id2)

		e1, err := client.Alerts.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.Alerts.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})

	t.Run("show by id: check fields", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		createAlertHandle(t, id)

		alert, err := client.Alerts.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "ROLE", alert.OwnerRoleType)
	})
}
