//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_Budgets(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("create and drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Budgets.Create(ctx, sdk.NewCreateBudgetRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Budget.DropFunc(t, id))

		err = client.Budgets.Drop(ctx, sdk.NewDropBudgetRequest(id))
		require.NoError(t, err)
	})

	t.Run("SetSpendingLimit and GetSpendingLimit", func(t *testing.T) {
		budgetId, budgetCleanup := testClientHelper().Budget.Create(t)
		t.Cleanup(budgetCleanup)

		result, err := client.Budgets.SetSpendingLimit(ctx, sdk.NewSetSpendingLimitBudgetRequest(budgetId, *sdk.NewBudgetSetSpendingLimitArgsRequest(500)))
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Contains(t, *result, "The spending limit has been updated to 500 credits.")

		spendingLimit, err := client.Budgets.GetSpendingLimit(ctx, sdk.NewGetSpendingLimitBudgetRequest(budgetId))
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, 500, *spendingLimit)
	})

	t.Run("SetEmailNotifications, GetNotificationEmail, and GetNotificationIntegrationName", func(t *testing.T) {
		budgetId, budgetCleanup := testClientHelper().Budget.Create(t)
		t.Cleanup(budgetCleanup)

		integration, integrationCleanup := testClientHelper().NotificationIntegration.Create(t)
		t.Cleanup(integrationCleanup)

		revokePrivilege := testClientHelper().Grant.GrantUsageOnIntegrationToSnowflakeApplication(t, integration.ID())
		t.Cleanup(revokePrivilege)

		result, err := client.Budgets.SetEmailNotifications(ctx, sdk.NewSetEmailNotificationsBudgetRequest(
			budgetId,
			*sdk.NewBudgetSetEmailNotificationsArgsRequestFromEmails("artur.sawicki@snowflake.com").
				WithNotificationIntegration(integration.ID()),
		))
		require.NoError(t, err)
		require.NotNil(t, result)

		email, err := client.Budgets.GetNotificationEmail(ctx, sdk.NewGetNotificationEmailBudgetRequest(budgetId))
		require.NoError(t, err)
		require.NotNil(t, email)
		require.Equal(t, *email, "artur.sawicki@snowflake.com")

		integrationName, err := client.Budgets.GetNotificationIntegrationName(ctx, sdk.NewGetNotificationIntegrationNameBudgetRequest(budgetId))
		require.NoError(t, err)
		require.NotNil(t, integrationName)
		require.Equal(t, integration.ID().Name(), *integrationName)
	})

	t.Run("SetCycleStartAction and GetCycleStartAction", func(t *testing.T) {
		budgetId, budgetCleanup := testClientHelper().Budget.Create(t)
		t.Cleanup(budgetCleanup)

		procedureId, procCleanup := testClientHelper().Procedure.Create(t)
		t.Cleanup(procCleanup)

		result, err := client.Budgets.SetCycleStartAction(ctx, sdk.NewSetCycleStartActionBudgetRequest(
			budgetId, *sdk.NewBudgetSetCycleStartActionArgsRequest(procedureId.ID().SchemaObjectId(), []string{}),
		))
		require.NoError(t, err)
		require.NotNil(t, result)

		action, err := client.Budgets.GetCycleStartAction(ctx, sdk.NewGetCycleStartActionBudgetRequest(budgetId))
		require.NoError(t, err)
		require.NotNil(t, action)
		require.Equal(t, procedureId.ID().SchemaObjectId().FullyQualifiedName(), action.ProcedureFqn)
	})
}
