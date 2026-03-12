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

func (c *CatalogIntegrationClient) exec(sql string) error {
	ctx := context.Background()
	_, err := c.context.client.ExecForTests(ctx, sql)
	return err
}

func (c *CatalogIntegrationClient) Create(t *testing.T) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	id := c.ids.RandomAccountObjectIdentifier()

	err := c.client().Create(ctx, sdk.NewCreateCatalogIntegrationRequest(id, true).
		WithObjectStorageCatalogSourceParams(*sdk.NewObjectStorageParamsRequest(sdk.CatalogIntegrationTableFormatIceberg)))
	require.NoError(t, err)

	return id, c.DropFunc(t, id)
}

func (c *CatalogIntegrationClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropCatalogIntegrationRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *CatalogIntegrationClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegration, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *CatalogIntegrationClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) ([]sdk.CatalogIntegrationProperty, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().Describe(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeAwsGlue(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationAwsGlueDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeAwsGlueDetails(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeObjectStorage(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationObjectStorageDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeObjectStorageDetails(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeOpenCatalog(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationOpenCatalogDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeOpenCatalogDetails(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeIcebergRest(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationIcebergRestDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeIcebergRestDetails(ctx, id)
}

func (c *CatalogIntegrationClient) DescribeSapBdc(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.CatalogIntegrationSapBdcDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeSapBdcDetails(ctx, id)
}
