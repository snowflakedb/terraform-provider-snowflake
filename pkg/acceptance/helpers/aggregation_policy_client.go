package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const aggregationPolicyBody = "AGGREGATION_CONSTRAINT(MIN_GROUP_SIZE => 5)"

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

func (c *AggregationPolicyClient) client() sdk.AggregationPolicies {
	return c.context.client.AggregationPolicies
}

func (c *AggregationPolicyClient) CreateAggregationPolicy(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	policy, cleanup := c.CreateWithRequest(t, sdk.NewCreateAggregationPolicyRequest(id, aggregationPolicyBody))
	return policy.ID(), cleanup
}

func (c *AggregationPolicyClient) CreateWithRequest(t *testing.T, request *sdk.CreateAggregationPolicyRequest) (*sdk.AggregationPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	id := request.GetName()
	policy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return policy, c.DropAggregationPolicyFunc(t, id)
}

func (c *AggregationPolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.AggregationPolicy, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *AggregationPolicyClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.AggregationPolicyDetails, error) {
	t.Helper()
	return c.client().DescribeDetails(context.Background(), id)
}

func (c *AggregationPolicyClient) DropAggregationPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropAggregationPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
