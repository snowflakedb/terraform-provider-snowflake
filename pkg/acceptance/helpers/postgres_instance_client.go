package helpers

import (
	"context"
	"testing"
	"time"

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

func (c *PostgresInstanceClient) WaitForReady(t *testing.T, id sdk.AccountObjectIdentifier, timeout time.Duration) *sdk.PostgresInstance {
	t.Helper()
	ctx := context.Background()

	var instance *sdk.PostgresInstance
	require.Eventually(t, func() bool {
		var err error
		instance, err = c.client().ShowByID(ctx, id)
		require.NoError(t, err)
		return instance.State == sdk.PostgresInstanceStateReady
	}, timeout, 3*time.Second)
	return instance
}

func (c *PostgresInstanceClient) CreateAndWaitForReady(t *testing.T) (*sdk.PostgresInstance, func()) {
	t.Helper()
	instance, cleanup := c.Create(t)
	instance = c.WaitForReady(t, instance.ID(), 5*time.Minute)
	return instance, cleanup
}

func (c *PostgresInstanceClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.PostgresInstanceDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeDetails(ctx, id)
}

func (c *PostgresInstanceClient) Alter(t *testing.T, req *sdk.AlterPostgresInstanceRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}
