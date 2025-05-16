package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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

func (c *CatalogIntegrationClient) exec(sql string) error {
	ctx := context.Background()
	_, err := c.context.client.ExecForTests(ctx, sql)
	return err
}

// TODO(SNOW-999142): Use SDK implementation for Catalog once it's available
func (c *CatalogIntegrationClient) Create(t *testing.T) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	id := c.ids.RandomAccountObjectIdentifier()
	err := c.exec(fmt.Sprintf(`
create catalog integration %s
  catalog_source=object_store
  table_format=iceberg
  enabled=true
`, id.FullyQualifiedName()))
	require.NoError(t, err)

	return id, c.DropFunc(t, id)
}

func (c *CatalogIntegrationClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()

	return func() {
		err := c.exec(fmt.Sprintf(`drop catalog integration if exists %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
