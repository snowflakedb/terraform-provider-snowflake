package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-1564954): change raw sqls to proper client
type AggregationPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewAggregationPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *AggregationPolicyClient {
	return &AggregationPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *AggregationPolicyClient) client() *sdk.Client {
	return c.context.client
}

func (c *AggregationPolicyClient) CreateAggregationPolicy(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE AGGREGATION POLICY %s AS () RETURNS AGGREGATION_CONSTRAINT -> AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 5)`, id.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropAggregationPolicyFunc(t, id)
}

func (c *AggregationPolicyClient) DropAggregationPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP AGGREGATION POLICY IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
