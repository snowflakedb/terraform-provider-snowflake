package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type PasswordPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewPasswordPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *PasswordPolicyClient {
	return &PasswordPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *PasswordPolicyClient) client() sdk.PasswordPolicies {
	return c.context.client.PasswordPolicies
}

func (c *PasswordPolicyClient) CreatePasswordPolicy(t *testing.T) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	return c.CreatePasswordPolicyInSchema(t, c.ids.SchemaId())
}

func (c *PasswordPolicyClient) CreatePasswordPolicyInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	return c.CreatePasswordPolicyWithOptions(t, sdk.NewCreatePasswordPolicyRequest(c.ids.RandomSchemaObjectIdentifierInSchema(schemaId)))
}

func (c *PasswordPolicyClient) CreatePasswordPolicyWithOptions(t *testing.T, req *sdk.CreatePasswordPolicyRequest) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	passwordPolicy, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return passwordPolicy, c.DropPasswordPolicyFunc(t, req.GetName())
}

func (c *PasswordPolicyClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.PasswordPolicy, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByID(ctx, id)
}

func (c *PasswordPolicyClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.PasswordPolicyDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().DescribeDetails(ctx, id)
}

func (c *PasswordPolicyClient) Alter(t *testing.T, request *sdk.AlterPasswordPolicyRequest) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}

func (c *PasswordPolicyClient) DropPasswordPolicyFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropPasswordPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
