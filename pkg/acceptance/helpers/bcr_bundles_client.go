package helpers

import (
	"context"
	"testing"

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

func (c *BcrBundlesClient) EnableBcrBundle(t *testing.T, name string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().EnableBehaviorChangeBundle(ctx, name)
	require.NoError(t, err)

	t.Cleanup(c.DisableBcrBundleFunc(t, name))
}

func (c *BcrBundlesClient) DisableBcrBundle(t *testing.T, name string) {
	t.Helper()
	ctx := context.Background()

	err := c.client().DisableBehaviorChangeBundle(ctx, name)
	require.NoError(t, err)

	t.Cleanup(c.EnableBcrBundleFunc(t, name))
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
