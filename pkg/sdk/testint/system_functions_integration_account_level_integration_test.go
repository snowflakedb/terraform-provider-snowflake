//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Our assumptions for Snowflake behavior are:
// 1. We have at least 2 bundles.
// 2. We don't assume which bundle is active/inactive by default.
// 3. In this test we don't want to assert the bundle names as they may change. So, we simply assert that they are not empty.
// 4. After each test, we clean up the bundle state by conditionally reverting the operation, i.e. by enabling the bundle if isDefault=true or disabling the bundle if isDefault=false.
func TestInt_BcrBundles_AccountLevel(t *testing.T) {
	client := testSecondaryClient(t)
	ctx := testContext(t)

	bundles := secondaryTestClientHelper().BcrBundles.ShowActiveBundles(t)
	require.GreaterOrEqual(t, len(bundles), 2)

	t.Run("show active bundles", func(t *testing.T) {
		bundles, err := client.SystemFunctions.ShowActiveBehaviorChangeBundles(ctx)
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(bundles), 2)
		for _, bundle := range bundles {
			assert.NotEmpty(t, bundle.Name)
		}
	})

	t.Run("enable a valid bundle", func(t *testing.T) {
		bundle := bundles[1]
		err := client.SystemFunctions.EnableBehaviorChangeBundle(ctx, bundle.Name)
		require.NoError(t, err)
		if !bundle.IsDefault {
			t.Cleanup(secondaryTestClientHelper().BcrBundles.DisableBcrBundleCleanupFunc(t, bundle.Name))
		}
		status := secondaryTestClientHelper().BcrBundles.BehaviorChangeBundleStatus(t, bundle.Name)
		require.Equal(t, sdk.BehaviorChangeBundleStatusEnabled, status)
	})

	t.Run("disable a valid bundle", func(t *testing.T) {
		bundle := bundles[0]
		err := client.SystemFunctions.DisableBehaviorChangeBundle(ctx, bundle.Name)
		require.NoError(t, err)
		if bundle.IsDefault {
			t.Cleanup(secondaryTestClientHelper().BcrBundles.EnableBcrBundleCleanupFunc(t, bundle.Name))
		}
		status := secondaryTestClientHelper().BcrBundles.BehaviorChangeBundleStatus(t, bundle.Name)
		require.Equal(t, sdk.BehaviorChangeBundleStatusDisabled, status)
	})
}
