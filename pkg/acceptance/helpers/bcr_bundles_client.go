package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type BcrBundlesClient struct {
	context *TestClientContext
}

func NewBcrBundlesClient(context *TestClientContext) *BcrBundlesClient {
	return &BcrBundlesClient{
		context: context,
	}
}

func (c *BcrBundlesClient) client() sdk.SystemFunctions {
	return c.context.client.SystemFunctions
}

func (c *BcrBundlesClient) ShowActiveBundles(t *testing.T) []sdk.BehaviorChangeBundleInfo {
	t.Helper()
	ctx := context.Background()

	bundles, err := c.client().ShowActiveBehaviorChangeBundles(ctx)
	require.NoError(t, err)

	return bundles
}

func (c *BcrBundlesClient) BehaviorChangeBundleStatus(t *testing.T, bundle string) sdk.BehaviorChangeBundleStatus {
	t.Helper()
	ctx := context.Background()

	status, err := c.client().BehaviorChangeBundleStatus(ctx, bundle)
	require.NoError(t, err)

	return status
}

func (c *BcrBundlesClient) GetBcrInfo(t *testing.T, name string) sdk.BehaviorChangeBundleInfo {
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
