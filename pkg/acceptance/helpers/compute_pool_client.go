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

	id := c.ids.RandomAccountObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreateComputePoolRequest(id, 1, 1, sdk.ComputePoolInstanceFamilyCpuX64XS))
}

func (c *ComputePoolClient) CreateWithRequest(t *testing.T, req *sdk.CreateComputePoolRequest) (*sdk.ComputePool, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	id := req.GetName()
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

func (c *ComputePoolClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ComputePool, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByID(ctx, id)
}

func (c *ComputePoolClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ComputePoolDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Describe(ctx, id)
}

func (c *ComputePoolClient) Alter(t *testing.T, req *sdk.AlterComputePoolRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}
