package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type SessionPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSessionPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *SessionPolicyClient {
	return &SessionPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *SessionPolicyClient) client() sdk.SessionPolicies {
	return c.context.client.SessionPolicies
}

func (c *SessionPolicyClient) CreateSessionPolicy(t *testing.T) (*sdk.SessionPolicy, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	return c.CreateSessionPolicyWithOptions(t, id, sdk.NewCreateSessionPolicyRequest(id))
}

func (c *SessionPolicyClient) CreateSessionPolicyWithOptions(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateSessionPolicyRequest) (*sdk.SessionPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	sessionPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return sessionPolicy, c.DropSessionPolicyFunc(t, id)
}

func (c *SessionPolicyClient) DropSessionPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropSessionPolicyRequest(id).WithIfExists(sdk.Bool(true)))
		require.NoError(t, err)
	}
}
