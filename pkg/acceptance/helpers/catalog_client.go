package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type CatalogIntegrationClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewCatalogIntegrationClient(context *TestClientContext, idsGenerator *IdsGenerator) *CatalogIntegrationClient {
	return &CatalogIntegrationClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *CatalogIntegrationClient) client() sdk.CatalogIntegrations {
	return c.context.client.CatalogIntegrations
}

func (c *CatalogIntegrationClient) Create(t *testing.T) (*sdk.CatalogIntegration, func()) {
	t.Helper()
	return c.CreateObjectStore(t)
}

func (c *CatalogIntegrationClient) CreateObjectStore(t *testing.T) (*sdk.CatalogIntegration, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateCatalogIntegrationRequest(c.ids.RandomAccountObjectIdentifier(), true).
		WithObjectStoreParams(*sdk.NewObjectStoreParamsRequest(sdk.TableFormatIceberg)),
	)
}

func (c *CatalogIntegrationClient) CreateGlue(t *testing.T, glueAwsRoleArn string, glueCatalogId string) (*sdk.CatalogIntegration, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateCatalogIntegrationRequest(c.ids.RandomAccountObjectIdentifier(), true).
		WithGlueParams(*sdk.NewGlueParamsRequest(sdk.TableFormatIceberg, glueAwsRoleArn, glueCatalogId)),
	)
}

func (c *CatalogIntegrationClient) CreateWithRequest(t *testing.T, request *sdk.CreateCatalogIntegrationRequest) (*sdk.CatalogIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	integration, err := c.client().ShowByID(ctx, request.GetName())
	require.NoError(t, err)

	return integration, c.DropFunc(t, request.GetName())
}

func (c *CatalogIntegrationClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropCatalogIntegrationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
