//go:build account_level_tests

package helpers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// EnableBcrBundle enables the specified behavior change bundle in the account.
// If the bundle is not enabled by default, it registers a cleanup to disable the bundle after the test.
func (c *BcrBundlesClient) EnableBcrBundle(t *testing.T, name string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().EnableBehaviorChangeBundle(ctx, name)
	require.NoError(t, err)

	bundle := c.GetBcrInfo(t, name)
	if !bundle.IsDefault {
		t.Cleanup(c.DisableBcrBundleCleanupFunc(t, name))
	}
}

// DisableBcrBundle disables the specified behavior change bundle in the account.
// If the bundle is enabled by default, it registers a cleanup to enable the bundle after the test.
func (c *BcrBundlesClient) DisableBcrBundle(t *testing.T, name string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().DisableBehaviorChangeBundle(ctx, name)
	require.NoError(t, err)

	bundle := c.GetBcrInfo(t, name)
	if bundle.IsDefault {
		t.Cleanup(c.EnableBcrBundleCleanupFunc(t, name))
	}
}

func (c *BcrBundlesClient) DisableBcrBundleCleanupFunc(t *testing.T, name string) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().DisableBehaviorChangeBundle(ctx, name)
		require.NoError(t, err)
	}
}

func (c *BcrBundlesClient) EnableBcrBundleCleanupFunc(t *testing.T, name string) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().EnableBehaviorChangeBundle(ctx, name)
		require.NoError(t, err)
	}
}
