package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const joinPolicyBody = "JOIN_CONSTRAINT(JOIN_REQUIRED => TRUE)"

type JoinPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewJoinPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *JoinPolicyClient {
	return &JoinPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *JoinPolicyClient) client() sdk.JoinPolicies {
	return c.context.client.JoinPolicies
}

func (c *JoinPolicyClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	policy, cleanup := c.CreateWithRequest(t, sdk.NewCreateJoinPolicyRequest(id, joinPolicyBody))
	return policy.ID(), cleanup
}

func (c *JoinPolicyClient) CreateWithRequest(t *testing.T, request *sdk.CreateJoinPolicyRequest) (*sdk.JoinPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	id := request.GetName()
	policy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return policy, c.DropFunc(t, id)
}

func (c *JoinPolicyClient) Alter(t *testing.T, request *sdk.AlterJoinPolicyRequest) {
	t.Helper()
	err := c.client().Alter(context.Background(), request)
	require.NoError(t, err)
}

func (c *JoinPolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.JoinPolicy, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *JoinPolicyClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.JoinPolicyDetails, error) {
	t.Helper()
	return c.client().DescribeDetails(context.Background(), id)
}

func (c *JoinPolicyClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropJoinPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
