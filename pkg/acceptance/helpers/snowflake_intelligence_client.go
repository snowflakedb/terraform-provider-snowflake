package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type SnowflakeIntelligenceClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSnowflakeIntelligenceClient(context *TestClientContext, idsGenerator *IdsGenerator) *SnowflakeIntelligenceClient {
	return &SnowflakeIntelligenceClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *SnowflakeIntelligenceClient) client() *sdk.Client {
	return c.context.client
}

func (c *SnowflakeIntelligenceClient) Create(t *testing.T) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	err := c.client().SnowflakeIntelligences.Create(ctx, sdk.NewCreateSnowflakeIntelligenceRequest(id))
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *SnowflakeIntelligenceClient) ShowAgents(t *testing.T, id sdk.AccountObjectIdentifier) ([]sdk.SnowflakeIntelligenceAgent, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().SnowflakeIntelligences.ShowAgents(ctx, sdk.NewShowAgentsSnowflakeIntelligenceRequest(id))
}

// DropAgent detaches the agent when it is currently attached. It is a no-op when the agent is already absent.
func (c *SnowflakeIntelligenceClient) DropAgent(t *testing.T, id sdk.AccountObjectIdentifier, agentId sdk.SchemaObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().SnowflakeIntelligences.DropAgentSafely(ctx, id, agentId)
	require.NoError(t, err)
}

func (c *SnowflakeIntelligenceClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.SnowflakeIntelligence, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().SnowflakeIntelligences.ShowByID(ctx, id)
}

func (c *SnowflakeIntelligenceClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.SnowflakeIntelligenceDetails, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().SnowflakeIntelligences.Describe(ctx, id)
}

func (c *SnowflakeIntelligenceClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().SnowflakeIntelligences.Drop(ctx, sdk.NewDropSnowflakeIntelligenceRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
