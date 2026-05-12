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

func (c *AgentClient) client() sdk.CortexAgents {
	return c.context.client.CortexAgents
}

func (c *AgentClient) CreateWithId(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateCortexAgentRequest(id, c.SampleSpecWithResponse(t, "Sample response")))
}

func (c *AgentClient) CreateWithRequest(t *testing.T, req *sdk.CreateCortexAgentRequest) func() {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	return c.DropFunc(t, req.GetName())
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

func (c *AgentClient) SampleSpecWithResponse(t *testing.T, response string) string {
	t.Helper()
	return fmt.Sprintf(`orchestration:
  budget:
    seconds: 30
    tokens: 16000
instructions:
  response: "%s"
`, response)
}
