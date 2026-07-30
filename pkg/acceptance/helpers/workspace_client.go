package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO [SNOW-3825231]: change raw sqls to proper client
type WorkspaceClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewWorkspaceClient(context *TestClientContext, idsGenerator *IdsGenerator) *WorkspaceClient {
	return &WorkspaceClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *WorkspaceClient) client() *sdk.Client {
	return c.context.client
}

func (c *WorkspaceClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE WORKSPACE %s`, id.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *WorkspaceClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP WORKSPACE IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
