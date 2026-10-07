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

func (c *SnowflakeIntelligenceClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().SnowflakeIntelligences.Drop(ctx, sdk.NewDropSnowflakeIntelligenceRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
