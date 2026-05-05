package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type PostgresInstanceClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewPostgresInstanceClient(context *TestClientContext, idsGenerator *IdsGenerator) *PostgresInstanceClient {
	return &PostgresInstanceClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *PostgresInstanceClient) client() sdk.PostgresInstances {
	return c.context.client.PostgresInstances
}

func (c *PostgresInstanceClient) Create(t *testing.T) (*sdk.PostgresInstance, func()) {
	t.Helper()

	id := c.ids.RandomAccountObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreatePostgresInstanceRequest(id, "STANDARD_M", 10, sdk.PostgresInstanceAuthenticationAuthorityPostgres))
}

func (c *PostgresInstanceClient) CreateWithRequest(t *testing.T, req *sdk.CreatePostgresInstanceRequest) (*sdk.PostgresInstance, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	id := req.GetName()
	postgresInstance, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return postgresInstance, c.DropFunc(t, id)
}

func (c *PostgresInstanceClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropPostgresInstanceRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *PostgresInstanceClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.PostgresInstance, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByID(ctx, id)
}
