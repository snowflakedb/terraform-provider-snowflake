package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-3648593): change raw sqls to proper client
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
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE SNOWFLAKE INTELLIGENCE %s`, id.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *SnowflakeIntelligenceClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP SNOWFLAKE INTELLIGENCE IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
