package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ApplicationClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewApplicationClient(context *TestClientContext, idsGenerator *IdsGenerator) *ApplicationClient {
	return &ApplicationClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ApplicationClient) client() sdk.Applications {
	return c.context.client.Applications
}

func (c *ApplicationClient) CreateApplication(t *testing.T, packageId sdk.AccountObjectIdentifier, version string) (*sdk.Application, func()) {
	t.Helper()
	ctx := context.Background()
	id := c.ids.RandomAccountObjectIdentifier()
	err := c.client().Create(ctx, sdk.NewCreateApplicationRequest(id, packageId).WithVersion(*sdk.NewApplicationVersionRequest().WithVersionAndPatch(*sdk.NewVersionAndPatchRequest(version, nil))))
	require.NoError(t, err)

	application, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return application, c.DropApplicationFunc(t, id)
}

func (c *ApplicationClient) DropApplicationFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropApplicationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
