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

	t.Run("Create and Drop", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.Budgets.Create(ctx, sdk.NewCreateBudgetRequest(id))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Budget.DropFunc(t, id))
	})

	t.Run("SetSpendingLimit and GetSpendingLimit", func(t *testing.T) {
		id, cleanup := testClientHelper().Budget.Create(t)
		t.Cleanup(cleanup)
		t.Skip("TODO: implement once scalar result scanning is in place")
		_ = id
	})

	t.Run("SetEmailNotifications and GetNotificationIntegrations", func(t *testing.T) {
		id, cleanup := testClientHelper().Budget.Create(t)
		t.Cleanup(cleanup)
		t.Skip("TODO: implement once notification integration setup is available")
		_ = id
	})

	t.Run("SetCycleStartAction and GetCycleStartAction", func(t *testing.T) {
		id, cleanup := testClientHelper().Budget.Create(t)
		t.Cleanup(cleanup)
		t.Skip("TODO: implement once scalar result scanning and procedure reference setup is available")
		_ = id
	})
}
