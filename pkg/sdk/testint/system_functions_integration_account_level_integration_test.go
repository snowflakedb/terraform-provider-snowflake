//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Our assumptions for Snowflake behavior are:
// 1. We have 2 bundles.
// 2. The first bundle is always active by default.
// 3. The second bundle is always not active by default.
// 4. In this test we don't want to assert the bundle names as they may change. So, we simply assert that they are not empty.
// 5. After each test, we clean up the bundle state by reverting the operation, i.e. by enabling the first bundle or disabling the second bundle.
func TestInt_BcrBundles_AccountLevel(t *testing.T) {
	client := testSecondaryClient(t)
	ctx := testContext(t)

	bundles := secondaryTestClientHelper().BcrBundles.ShowActiveBundles(t)
	require.Len(t, bundles, 2)

	t.Run("show active bundles", func(t *testing.T) {
		bundles, err := client.SystemFunctions.ShowActiveBehaviorChangeBundles(ctx)
		require.NoError(t, err)
		assert.Len(t, bundles, 2)
		assert.NotEmpty(t, bundles[0].Name)
		assert.NotEmpty(t, bundles[1].Name)
		assert.True(t, bundles[0].IsDefault)
		assert.False(t, bundles[1].IsDefault)
		assert.True(t, bundles[0].IsEnabled)
		assert.False(t, bundles[1].IsEnabled)
	})

	t.Run("enable a valid bundle", func(t *testing.T) {
		err := client.SystemFunctions.EnableBehaviorChangeBundle(ctx, bundles[1].Name)
		require.NoError(t, err)
		t.Cleanup(secondaryTestClientHelper().BcrBundles.DisableBcrBundleFunc(t, bundles[1].Name))
		status := secondaryTestClientHelper().BcrBundles.BehaviorChangeBundleStatus(t, bundles[1].Name)
		require.Equal(t, sdk.BehaviorChangeBundleStatusEnabled, status)
	})

	t.Run("disable a valid bundle", func(t *testing.T) {
		err := client.SystemFunctions.DisableBehaviorChangeBundle(ctx, bundles[0].Name)
		require.NoError(t, err)
		t.Cleanup(secondaryTestClientHelper().BcrBundles.EnableBcrBundleFunc(t, bundles[0].Name))
		status := secondaryTestClientHelper().BcrBundles.BehaviorChangeBundleStatus(t, bundles[0].Name)
		require.Equal(t, sdk.BehaviorChangeBundleStatusDisabled, status)
	})
}
