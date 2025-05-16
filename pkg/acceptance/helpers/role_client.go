package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type RoleClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewRoleClient(context *TestClientContext, idsGenerator *IdsGenerator) *RoleClient {
	return &RoleClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *RoleClient) client() sdk.Roles {
	return c.context.client.Roles
}

func (c *RoleClient) UseRole(t *testing.T, roleId sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	currentRole, err := c.context.client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	err = c.context.client.Sessions.UseRole(ctx, roleId)
	require.NoError(t, err)

	return func() {
		err = c.context.client.Sessions.UseRole(ctx, currentRole)
		require.NoError(t, err)
	}
}

func (c *RoleClient) CreateRole(t *testing.T) (*sdk.Role, func()) {
	t.Helper()
	return c.CreateRoleWithRequest(t, sdk.NewCreateRoleRequest(c.ids.RandomAccountObjectIdentifier()))
}

func (c *RoleClient) CreateRoleWithIdentifier(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Role, func()) {
	t.Helper()
	return c.CreateRoleWithRequest(t, sdk.NewCreateRoleRequest(id))
}

func (c *RoleClient) CreateRoleGrantedToCurrentUser(t *testing.T) (*sdk.Role, func()) {
	t.Helper()
	ctx := context.Background()

	role, roleCleanup := c.CreateRole(t)

	currentUser, err := c.context.client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)

	c.GrantRoleToUser(t, role.ID(), currentUser)
	return role, roleCleanup
}

func (c *RoleClient) CreateRoleWithRequest(t *testing.T, req *sdk.CreateRoleRequest) (*sdk.Role, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	role, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)
	return role, c.DropRoleFunc(t, req.GetName())
}

func (c *RoleClient) DropRoleFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropRoleRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *RoleClient) GrantRoleToUser(t *testing.T, id sdk.AccountObjectIdentifier, userId sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Grant(ctx, sdk.NewGrantRoleRequest(id, sdk.GrantRole{
		User: sdk.Pointer(userId),
	}))
	require.NoError(t, err)
}

func (c *RoleClient) GrantRoleToCurrentRole(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	currentRole, err := c.context.client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	err = c.client().Grant(ctx, sdk.NewGrantRoleRequest(id, sdk.GrantRole{
		Role: sdk.Pointer(currentRole),
	}))
	require.NoError(t, err)
}
