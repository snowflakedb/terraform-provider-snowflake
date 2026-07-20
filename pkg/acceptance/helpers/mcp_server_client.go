package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type McpServerClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewMcpServerClient(context *TestClientContext, idsGenerator *IdsGenerator) *McpServerClient {
	return &McpServerClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *McpServerClient) client() sdk.McpServers {
	return c.context.client.McpServers
}

func (c *McpServerClient) Create(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	return c.CreateWithSpec(t, id, c.DefaultSpec())
}

func (c *McpServerClient) CreateWithSpec(t *testing.T, id sdk.SchemaObjectIdentifier, spec string) func() {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateMcpServerRequest(id, spec))
}

func (c *McpServerClient) CreateWithRequest(t *testing.T, req *sdk.CreateMcpServerRequest) func() {
	t.Helper()
	ctx := context.Background()
	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	return c.DropFunc(t, req.GetName())
}

func (c *McpServerClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	return func() {
		err := c.client().Drop(ctx, sdk.NewDropMcpServerRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *McpServerClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.McpServer, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByID(ctx, id)
}

func (c *McpServerClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.McpServerDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().Describe(ctx, id)
}

func (c *McpServerClient) DefaultSpec() string {
	return `tools:
  - title: "SQL Execution Tool"
    name: "sql_exec_tool"
    type: "SYSTEM_EXECUTE_SQL"
    description: "For acceptance tests."
`
}
