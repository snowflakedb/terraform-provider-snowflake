package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-1790174]: change raw sqls to proper client
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

func (c *ComputePoolClient) client() *sdk.Client {
	return c.context.client
}

func (c *ComputePoolClient) CreateComputePool(t *testing.T) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE COMPUTE POOL %s MIN_NODES = 1 MAX_NODES = 1 INSTANCE_FAMILY = CPU_X64_XS`, id.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropComputePoolFunc(t, id)
}

func (c *ComputePoolClient) DropComputePoolFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP COMPUTE POOL IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
