package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type HybridTableClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewHybridTableClient(context *TestClientContext, idsGenerator *IdsGenerator) *HybridTableClient {
	return &HybridTableClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *HybridTableClient) exec(sql string) error {
	ctx := context.Background()
	_, err := c.context.client.ExecForTests(ctx, sql)
	return err
}

// TODO(SNOW-999142): Use SDK implementation for Hybrid Table once it's available
func (c *HybridTableClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	err := c.exec(fmt.Sprintf(`
create hybrid table %s (
  id INT AUTOINCREMENT PRIMARY KEY
)
`, id.FullyQualifiedName()))
	require.NoError(t, err)

	return id, c.DropFunc(t, id)
}

func (c *HybridTableClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()

	return func() {
		err := c.exec(fmt.Sprintf(`drop table if exists %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
