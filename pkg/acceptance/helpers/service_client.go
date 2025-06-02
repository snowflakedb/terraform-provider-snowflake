package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ServiceClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewServiceClient(context *TestClientContext, idsGenerator *IdsGenerator) *ServiceClient {
	return &ServiceClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ServiceClient) client() sdk.Services {
	return c.context.client.Services
}

func (c *ServiceClient) Create(t *testing.T, computePoolId sdk.AccountObjectIdentifier) (*sdk.Service, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	err := c.client().Create(ctx, sdk.NewCreateServiceRequest(id, computePoolId))
	require.NoError(t, err)
	service, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return service, c.DropFunc(t, id)
}

func (c *ServiceClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropServiceRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
