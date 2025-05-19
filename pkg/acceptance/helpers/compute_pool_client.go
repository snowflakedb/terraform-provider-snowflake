package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ComputePoolClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewComputePoolClient(context *TestClientContext, idsGenerator *IdsGenerator) *ComputePoolClient {
	return &ComputePoolClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ComputePoolClient) client() sdk.ComputePools {
	return c.context.client.ComputePools
}

func (c *ComputePoolClient) Create(t *testing.T) (*sdk.ComputePool, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	err := c.client().Create(ctx, sdk.NewCreateComputePoolRequest(id, 1, 1, sdk.ComputePoolInstanceFamilyCpuX64XS))
	require.NoError(t, err)
	computePool, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return computePool, c.DropFunc(t, id)
}

func (c *ComputePoolClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropComputePoolRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
