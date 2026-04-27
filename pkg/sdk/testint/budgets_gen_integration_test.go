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

	t.Run("SetEmailNotifications and GetNotificationIntegrations", func(t *testing.T) {
		t.Skip("TODO: implement once notification integration setup is available")
		budgetId, budgetCleanup := testClientHelper().Budget.Create(t)
		t.Cleanup(budgetCleanup)
		_ = budgetId
	})

	t.Run("SetCycleStartAction and GetCycleStartAction", func(t *testing.T) {
		t.Skip("TODO: implement once scalar result scanning and procedure reference setup is available")
		budgetId, budgetCleanup := testClientHelper().Budget.Create(t)
		t.Cleanup(budgetCleanup)
		_ = budgetId
	})
}
