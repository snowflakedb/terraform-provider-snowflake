package helpers

import (
	"context"
	"fmt"
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

func (c *McpServerClient) Create(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	spec := `tools:
  - title: "SQL Execution Tool"
    name: "sql_exec_tool"
    type: "SYSTEM_EXECUTE_SQL"
    description: "For acceptance tests."
`
	createSQL := fmt.Sprintf("CREATE OR REPLACE MCP SERVER %s FROM SPECIFICATION $$%s$$", id.FullyQualifiedName(), spec)
	_, err := c.context.client.ExecForTests(ctx, createSQL)
	require.NoError(t, err)

	cleanup := func() {
		dropSQL := fmt.Sprintf("DROP MCP SERVER IF EXISTS %s", id.FullyQualifiedName())
		_, dropErr := c.context.client.ExecForTests(context.Background(), dropSQL)
		require.NoError(t, dropErr)
	}
	return cleanup
}
