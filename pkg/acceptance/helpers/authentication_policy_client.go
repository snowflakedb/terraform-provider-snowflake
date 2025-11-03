package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type AuthenticationPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewAuthenticationPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *AuthenticationPolicyClient {
	return &AuthenticationPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *AuthenticationPolicyClient) client() sdk.AuthenticationPolicies {
	return c.context.client.AuthenticationPolicies
}

func (c *AuthenticationPolicyClient) Create(t *testing.T) (*sdk.AuthenticationPolicy, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	return c.CreateWithOptions(t, id, sdk.NewCreateAuthenticationPolicyRequest(id))
}

func (c *AuthenticationPolicyClient) CreateWithOptions(t *testing.T, id sdk.SchemaObjectIdentifier, request *sdk.CreateAuthenticationPolicyRequest) (*sdk.AuthenticationPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	authenticationPolicy, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return authenticationPolicy, c.DropFunc(t, id)
}

func (c *AuthenticationPolicyClient) Alter(t *testing.T, request *sdk.AlterAuthenticationPolicyRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}

func (c *AuthenticationPolicyClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropAuthenticationPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *AuthenticationPolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.AuthenticationPolicy, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByID(ctx, id)
}
