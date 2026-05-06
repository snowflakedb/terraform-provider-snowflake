package helpers

import (
	"context"
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

func (c *AgentClient) client() sdk.CortexAgents {
	return c.context.client.CortexAgents
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
	err := c.client().Create(ctx, sdk.NewCreateCortexAgentRequest(id, spec))
	require.NoError(t, err)

	return c.DropFunc(t, id)
}

func (c *AgentClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropCortexAgentRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *AgentClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.CortexAgent, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *AgentClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.CortexAgentDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Describe(ctx, id)
}
