package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type DatabaseRoleClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewDatabaseRoleClient(context *TestClientContext, idsGenerator *IdsGenerator) *DatabaseRoleClient {
	return &DatabaseRoleClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *DatabaseRoleClient) client() sdk.DatabaseRoles {
	return c.context.client.DatabaseRoles
}

func (c *DatabaseRoleClient) CreateDatabaseRole(t *testing.T) (*sdk.DatabaseRole, func()) {
	t.Helper()
	return c.CreateDatabaseRoleInDatabase(t, c.ids.DatabaseId())
}

func (c *DatabaseRoleClient) CreateDatabaseRoleInDatabase(t *testing.T, databaseId sdk.AccountObjectIdentifier) (*sdk.DatabaseRole, func()) {
	t.Helper()
	return c.CreateDatabaseRoleInDatabaseWithName(t, databaseId, c.ids.Alpha())
}

func (c *DatabaseRoleClient) CreateDatabaseRoleWithName(t *testing.T, name string) (*sdk.DatabaseRole, func()) {
	t.Helper()
	return c.CreateDatabaseRoleInDatabaseWithName(t, c.ids.DatabaseId(), name)
}

func (c *DatabaseRoleClient) CreateDatabaseRoleInDatabaseWithName(t *testing.T, databaseId sdk.AccountObjectIdentifier, name string) (*sdk.DatabaseRole, func()) {
	t.Helper()
	ctx := context.Background()

	id := sdk.NewDatabaseObjectIdentifier(databaseId.Name(), name)

	err := c.client().Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
	require.NoError(t, err)

	databaseRole, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return databaseRole, c.CleanupDatabaseRoleFunc(t, id)
}

func (c *DatabaseRoleClient) CleanupDatabaseRoleFunc(t *testing.T, id sdk.DatabaseObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		// to prevent error when db was removed before the role
		_, err := c.context.client.Databases.ShowByID(ctx, id.DatabaseId())
		if errors.Is(err, sdk.ErrObjectNotFound) {
			return
		}

		err = c.client().Drop(ctx, sdk.NewDropDatabaseRoleRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *DatabaseRoleClient) Show(t *testing.T, id sdk.DatabaseObjectIdentifier) (*sdk.DatabaseRole, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}
