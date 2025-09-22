//go:build account_level_tests

package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// EnableBcrBundle enables the specified behavior change bundle in the account.
// If the bundle is not enabled by default, it registers a cleanup to disable the bundle after the test.
func (c *BcrBundlesClient) EnableBcrBundle(t *testing.T, name string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().EnableBehaviorChangeBundle(ctx, name)
	require.NoError(t, err)

	bundle := c.getBcrInfo(t, name)
	if !bundle.IsDefault {
		t.Cleanup(c.DisableBcrBundleFunc(t, name))
	}
}

// DisableBcrBundle disables the specified behavior change bundle in the account.
// If the bundle is enabled by default, it registers a cleanup to enable the bundle after the test.
func (c *BcrBundlesClient) DisableBcrBundle(t *testing.T, name string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().DisableBehaviorChangeBundle(ctx, name)
	require.NoError(t, err)

	bundle := c.getBcrInfo(t, name)
	if bundle.IsDefault {
		t.Cleanup(c.EnableBcrBundleFunc(t, name))
	}
}

func (c *BcrBundlesClient) DisableBcrBundleFunc(t *testing.T, name string) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().DisableBehaviorChangeBundle(ctx, name)
		require.NoError(t, err)
	}
}

func (c *BcrBundlesClient) EnableBcrBundleFunc(t *testing.T, name string) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().EnableBehaviorChangeBundle(ctx, name)
		require.NoError(t, err)
	}
}

func (c *BcrBundlesClient) getBcrInfo(t *testing.T, name string) sdk.BehaviorChangeBundleInfo {
	t.Helper()
	ctx := context.Background()

	bundles, err := c.client().ShowActiveBehaviorChangeBundles(ctx)
	require.NoError(t, err)

	info, err := collections.FindFirst(bundles, func(bundle sdk.BehaviorChangeBundleInfo) bool {
		return bundle.Name == name
	})
	require.NoError(t, err)
	require.NotNil(t, info)

	return *info
}
