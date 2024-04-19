package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ContextClient struct {
	context *TestClientContext
}

func NewContextClient(context *TestClientContext) *ContextClient {
	return &ContextClient{
		context: context,
	}
}

func (c *ContextClient) client() sdk.ContextFunctions {
	return c.context.client.ContextFunctions
}

func (c *ContextClient) CurrentAccount(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	currentAccount, err := c.client().CurrentAccount(ctx)
	require.NoError(t, err)

	return currentAccount
}

func (c *ContextClient) CurrentRole(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	currentRole, err := c.client().CurrentRole(ctx)
	require.NoError(t, err)

	return currentRole
}

func (c *ContextClient) CurrentRegion(t *testing.T) string {
	t.Helper()
	ctx := context.Background()

	currentRegion, err := c.client().CurrentRegion(ctx)
	require.NoError(t, err)

	return currentRegion
}
