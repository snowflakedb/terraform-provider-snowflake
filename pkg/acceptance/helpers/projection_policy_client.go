package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

const projectionPolicyBody = "PROJECTION_CONSTRAINT(ALLOW => false)"

type ProjectionPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewProjectionPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *ProjectionPolicyClient {
	return &ProjectionPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ProjectionPolicyClient) client() sdk.ProjectionPolicies {
	return c.context.client.ProjectionPolicies
}

func (c *ProjectionPolicyClient) CreateProjectionPolicy(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	policy, cleanup := c.CreateWithRequest(t, sdk.NewCreateProjectionPolicyRequest(id, projectionPolicyBody))
	return policy.ID(), cleanup
}

func (c *ProjectionPolicyClient) CreateWithRequest(t *testing.T, request *sdk.CreateProjectionPolicyRequest) (*sdk.ProjectionPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	id := request.GetName()
	policy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return policy, c.DropProjectionPolicyFunc(t, id)
}

func (c *ProjectionPolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.ProjectionPolicy, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *ProjectionPolicyClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.ProjectionPolicyDetails, error) {
	t.Helper()
	return c.client().DescribeDetails(context.Background(), id)
}

func (c *ProjectionPolicyClient) DropProjectionPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropProjectionPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
