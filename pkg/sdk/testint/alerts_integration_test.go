//go:build !account_level_tests

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
		alerts, err := client.Alerts.Show(ctx, nil)
		require.NoError(t, err)
		assert.Len(t, alerts, 2)
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &sdk.ShowAlertOptions{
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, alerts, *alertTest)
		assert.Contains(t, alerts, *alert2Test)
		assert.Len(t, alerts, 2)
	})

	t.Run("with show options and like", func(t *testing.T) {
		showOptions := &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alertTest.Name),
			},
			In: &sdk.In{
				Database: testClientHelper().Ids.DatabaseId(),
			},
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, alerts, *alertTest)
		assert.Len(t, alerts, 1)
	})

	t.Run("when searching a non-existent alert", func(t *testing.T) {
		showOptions := &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Empty(t, alerts)
	})

	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &sdk.ShowAlertOptions{
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
			Limit: sdk.Int(1),
		}
		alerts, err := client.Alerts.Show(ctx, showOptions)
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
		err := client.Alerts.Create(ctx, id, testClientHelper().Ids.WarehouseId(), schedule, condition, action, &sdk.CreateAlertOptions{
			OrReplace:   sdk.Bool(true),
			IfNotExists: sdk.Bool(false),
			Comment:     sdk.String(comment),
		})
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testClientHelper().Ids.WarehouseId().Name(), alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, comment, *alertDetails.Comment)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
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
		err := client.Alerts.Create(ctx, id, testClientHelper().Ids.WarehouseId(), schedule, condition, action, &sdk.CreateAlertOptions{
			OrReplace:   sdk.Bool(false),
			IfNotExists: sdk.Bool(true),
			Comment:     sdk.String(comment),
		})
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testClientHelper().Ids.WarehouseId().Name(), alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, comment, *alertDetails.Comment)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
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
		err := client.Alerts.Create(ctx, id, testClientHelper().Ids.WarehouseId(), schedule, condition, action, nil)
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testClientHelper().Ids.WarehouseId().Name(), alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, action, alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
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
		err := client.Alerts.Create(ctx, id, testClientHelper().Ids.WarehouseId(), schedule, condition, action, nil)
		require.NoError(t, err)
		alertDetails, err := client.Alerts.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, name, alertDetails.Name)
		assert.Equal(t, testClientHelper().Ids.WarehouseId().Name(), alertDetails.Warehouse)
		assert.Equal(t, schedule, alertDetails.Schedule)
		assert.Equal(t, condition, alertDetails.Condition)
		assert.Equal(t, strings.TrimSpace(action), alertDetails.Action)

		alert, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(name),
			},
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
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

		alterOptions := &sdk.AlterAlertOptions{
			Set: &sdk.AlertSet{
				Schedule: &newSchedule,
			},
		}

		err := client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, newSchedule, alerts[0].Schedule)
	})

	t.Run("when modifying condition and action", func(t *testing.T) {
		alert, alertCleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(alertCleanup)
		newCondition := "select * from DUAL where false"

		alterOptions := &sdk.AlterAlertOptions{
			ModifyCondition: &[]string{newCondition},
		}

		err := client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, newCondition, alerts[0].Condition)

		newAction := "create table FOO(ID INT)"

		alterOptions = &sdk.AlterAlertOptions{
			ModifyAction: &newAction,
		}

		err = client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err = client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, newAction, alerts[0].Action)
	})

	t.Run("resume and then suspend", func(t *testing.T) {
		alert, alertCleanup := testClientHelper().Alert.CreateAlert(t)
		t.Cleanup(alertCleanup)

		alterOptions := &sdk.AlterAlertOptions{
			Action: &sdk.AlertActionResume,
		}

		err := client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err := client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
		require.NoError(t, err)
		assert.Len(t, alerts, 1)
		assert.Equal(t, alerts[0].State, sdk.AlertStateStarted)

		alterOptions = &sdk.AlterAlertOptions{
			Action: &sdk.AlertActionSuspend,
		}

		err = client.Alerts.Alter(ctx, alert.ID(), alterOptions)
		require.NoError(t, err)
		alerts, err = client.Alerts.Show(ctx, &sdk.ShowAlertOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(alert.Name),
			},
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		})
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
		err := client.Alerts.Drop(ctx, id, &sdk.DropAlertOptions{})
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when alert does not exist", func(t *testing.T) {
		err := client.Alerts.Drop(ctx, NonExistingSchemaObjectIdentifier, &sdk.DropAlertOptions{})
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
			err := client.Alerts.Drop(ctx, id, &sdk.DropAlertOptions{})
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createAlertHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		schedule, condition, action := "USING CRON * * * * * UTC", "SELECT 1", "SELECT 1"
		err := client.Alerts.Create(ctx, id, warehouseId, schedule, condition, action, &sdk.CreateAlertOptions{})
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
