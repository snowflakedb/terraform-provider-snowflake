package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type AgentClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewAgentClient(context *TestClientContext, idsGenerator *IdsGenerator) *AgentClient {
	return &AgentClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *AgentClient) Create(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	spec := `orchestration:
  budget:
    seconds: 30
    tokens: 16000
instructions:
  response: "Test agent for acceptance tests"
`
	createSQL := fmt.Sprintf("CREATE OR REPLACE AGENT %s FROM SPECIFICATION $$%s$$", id.FullyQualifiedName(), spec)
	_, err := c.context.client.ExecForTests(ctx, createSQL)
	require.NoError(t, err)

	cleanup := func() {
		dropSQL := fmt.Sprintf("DROP AGENT IF EXISTS %s", id.FullyQualifiedName())
		_, dropErr := c.context.client.ExecForTests(context.Background(), dropSQL)
		require.NoError(t, dropErr)
	}
	return cleanup
}
