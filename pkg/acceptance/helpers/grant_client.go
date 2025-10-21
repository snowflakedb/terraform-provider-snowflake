package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type GrantClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewGrantClient(context *TestClientContext, idsGenerator *IdsGenerator) *GrantClient {
	return &GrantClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *GrantClient) client() sdk.Grants {
	return c.context.client.Grants
}

func (c *GrantClient) GrantGlobalPrivilegesOnAccount(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	privileges []sdk.GlobalPrivilege,
) {
	t.Helper()
	ctx := context.Background()

	accountRoleGrantPrivileges := &sdk.AccountRoleGrantPrivileges{
		GlobalPrivileges: privileges,
	}
	on := &sdk.AccountRoleGrantOn{
		Account: sdk.Bool(true),
	}
	opts := &sdk.GrantPrivilegesToAccountRoleOptions{
		WithGrantOption: sdk.Bool(false),
	}
	err := c.client().GrantPrivilegesToAccountRole(ctx, accountRoleGrantPrivileges, on, accountRoleId, opts)
	require.NoError(t, err)
}

func (c *GrantClient) RevokeGlobalPrivilegesFromAccount(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	privileges []sdk.GlobalPrivilege,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().RevokePrivilegesFromAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			GlobalPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			Account: sdk.Bool(true),
		},
		accountRoleId,
		&sdk.RevokePrivilegesFromAccountRoleOptions{},
	)
	require.NoError(t, err)
}

func (c *GrantClient) GrantOnSchemaToAccountRole(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, accountRoleId sdk.AccountObjectIdentifier, privileges ...sdk.SchemaPrivilege) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantPrivilegesToAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			Schema: &sdk.GrantOnSchema{
				Schema: &schemaId,
			},
		},
		accountRoleId,
		new(sdk.GrantPrivilegesToAccountRoleOptions),
	)
	require.NoError(t, err)
}

func (c *GrantClient) RevokePrivilegesOnSchemaFromAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	schemaId sdk.DatabaseObjectIdentifier,
	privileges []sdk.SchemaPrivilege,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().RevokePrivilegesFromAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			Schema: &sdk.GrantOnSchema{
				Schema: &schemaId,
			},
		},
		accountRoleId,
		new(sdk.RevokePrivilegesFromAccountRoleOptions),
	)

	require.NoError(t, err)
}

func (c *GrantClient) RevokePrivilegesOnSchemaObjectFromAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	objectType sdk.ObjectType,
	schemaObjectIdentifier sdk.SchemaObjectIdentifier,
	privileges []sdk.SchemaObjectPrivilege,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().RevokePrivilegesFromAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				SchemaObject: &sdk.Object{
					ObjectType: objectType,
					Name:       schemaObjectIdentifier,
				},
			},
		},
		accountRoleId,
		new(sdk.RevokePrivilegesFromAccountRoleOptions),
	)

	require.NoError(t, err)
}

func (c *GrantClient) GrantPrivilegesOnSchemaObjectToAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	objectType sdk.ObjectType,
	schemaObjectIdentifier sdk.SchemaObjectIdentifier,
	privileges []sdk.SchemaObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantPrivilegesToAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				SchemaObject: &sdk.Object{
					ObjectType: objectType,
					Name:       schemaObjectIdentifier,
				},
			},
		},
		accountRoleId,
		&sdk.GrantPrivilegesToAccountRoleOptions{
			WithGrantOption: sdk.Bool(withGrantOption),
		},
	)
	require.NoError(t, err)
}

func (c *GrantClient) RevokePrivilegesOnDatabaseFromDatabaseRole(
	t *testing.T,
	databaseRoleId sdk.DatabaseObjectIdentifier,
	databaseId sdk.AccountObjectIdentifier,
	privileges []sdk.AccountObjectPrivilege,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().RevokePrivilegesFromDatabaseRole(
		ctx,
		&sdk.DatabaseRoleGrantPrivileges{
			DatabasePrivileges: privileges,
		},
		&sdk.DatabaseRoleGrantOn{
			Database: sdk.Pointer(databaseId),
		},
		databaseRoleId,
		new(sdk.RevokePrivilegesFromDatabaseRoleOptions),
	)
	require.NoError(t, err)
}

func (c *GrantClient) RevokePrivilegesOnDatabaseFromAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	databaseId sdk.AccountObjectIdentifier,
	privileges []sdk.AccountObjectPrivilege,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().RevokePrivilegesFromAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			AccountObjectPrivileges: privileges,
		},
		&sdk.AccountRoleGrantOn{
			AccountObject: &sdk.GrantOnAccountObject{
				Database: sdk.Pointer(databaseId),
			},
		},
		accountRoleId,
		new(sdk.RevokePrivilegesFromAccountRoleOptions),
	)
	require.NoError(t, err)
}

func (c *GrantClient) GrantPrivilegesOnDatabaseToAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	databaseId sdk.AccountObjectIdentifier,
	privileges []sdk.AccountObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	c.grantPrivilegesOnAccountLevelObjectToAccountRole(
		t,
		accountRoleId,
		&sdk.AccountRoleGrantOn{
			AccountObject: &sdk.GrantOnAccountObject{
				Database: sdk.Pointer(databaseId),
			},
		},
		privileges,
		withGrantOption,
	)
}

func (c *GrantClient) GrantPrivilegesOnWarehouseToAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	warehouseId sdk.AccountObjectIdentifier,
	privileges []sdk.AccountObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	c.grantPrivilegesOnAccountLevelObjectToAccountRole(
		t,
		accountRoleId,
		&sdk.AccountRoleGrantOn{
			AccountObject: &sdk.GrantOnAccountObject{
				Warehouse: sdk.Pointer(warehouseId),
			},
		},
		privileges,
		withGrantOption,
	)
}

func (c *GrantClient) grantPrivilegesOnAccountLevelObjectToAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	accountObjectGrantOn *sdk.AccountRoleGrantOn,
	privileges []sdk.AccountObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantPrivilegesToAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			AccountObjectPrivileges: privileges,
		},
		accountObjectGrantOn,
		accountRoleId,
		&sdk.GrantPrivilegesToAccountRoleOptions{
			WithGrantOption: sdk.Bool(withGrantOption),
		},
	)
	require.NoError(t, err)
}

func (c *GrantClient) GrantPrivilegesOnDatabaseToDatabaseRole(
	t *testing.T,
	databaseRoleId sdk.DatabaseObjectIdentifier,
	databaseId sdk.AccountObjectIdentifier,
	privileges []sdk.AccountObjectPrivilege,
	withGrantOption bool,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantPrivilegesToDatabaseRole(
		ctx,
		&sdk.DatabaseRoleGrantPrivileges{
			DatabasePrivileges: privileges,
		},
		&sdk.DatabaseRoleGrantOn{
			Database: sdk.Pointer(databaseId),
		},
		databaseRoleId,
		&sdk.GrantPrivilegesToDatabaseRoleOptions{
			WithGrantOption: sdk.Bool(withGrantOption),
		},
	)
	require.NoError(t, err)
}

func (c *GrantClient) GrantOwnershipToAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	objectType sdk.ObjectType,
	objectName sdk.ObjectIdentifier,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantOwnership(
		ctx,
		sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: objectType,
				Name:       objectName,
			},
		},
		sdk.OwnershipGrantTo{
			AccountRoleName: &accountRoleId,
		},
		new(sdk.GrantOwnershipOptions),
	)
	require.NoError(t, err)
}

func (c *GrantClient) GrantOwnershipToAccountRoleWithOwnershipOptions(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	objectType sdk.ObjectType,
	objectName sdk.ObjectIdentifier,
	options sdk.GrantOwnershipOptions,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantOwnership(
		ctx,
		sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: objectType,
				Name:       objectName,
			},
		},
		sdk.OwnershipGrantTo{
			AccountRoleName: &accountRoleId,
		},
		&options,
	)
	require.NoError(t, err)
}

func (c *GrantClient) GrantOwnershipOnSchemaObjectToAccountRole(
	t *testing.T,
	accountRoleId sdk.AccountObjectIdentifier,
	objectType sdk.ObjectType,
	objectId sdk.SchemaObjectIdentifier,
	outboundPrivileges sdk.OwnershipCurrentGrantsOutboundPrivileges,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantOwnership(
		ctx,
		sdk.OwnershipGrantOn{
			Object: &sdk.Object{
				ObjectType: objectType,
				Name:       objectId,
			},
		},
		sdk.OwnershipGrantTo{
			AccountRoleName: sdk.Pointer(accountRoleId),
		},
		&sdk.GrantOwnershipOptions{
			CurrentGrants: &sdk.OwnershipCurrentGrants{
				OutboundPrivileges: outboundPrivileges,
			},
		},
	)
	require.NoError(t, err)
}

func (c *GrantClient) GrantPrivilegeOnDatabaseToShare(
	t *testing.T,
	databaseId sdk.AccountObjectIdentifier,
	shareId sdk.AccountObjectIdentifier,
	privileges []sdk.ObjectPrivilege,
) func() {
	t.Helper()
	ctx := context.Background()

	err := c.client().GrantPrivilegeToShare(ctx, privileges, &sdk.ShareGrantOn{Database: databaseId}, shareId)
	require.NoError(t, err)

	return func() {
		c.RevokePrivilegeOnDatabaseFromShare(t, databaseId, shareId, privileges)
	}
}

func (c *GrantClient) RevokePrivilegeOnDatabaseFromShare(
	t *testing.T,
	databaseId sdk.AccountObjectIdentifier,
	shareId sdk.AccountObjectIdentifier,
	privileges []sdk.ObjectPrivilege,
) {
	t.Helper()
	ctx := context.Background()

	err := c.client().RevokePrivilegeFromShare(ctx, privileges, &sdk.ShareGrantOn{Database: databaseId}, shareId)
	require.NoError(t, err)
}

func (c *GrantClient) ShowGrantsToShare(t *testing.T, shareId sdk.AccountObjectIdentifier) ([]sdk.Grant, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Show(ctx, &sdk.ShowGrantOptions{
		To: &sdk.ShowGrantsTo{
			Share: &sdk.ShowGrantsToShare{
				Name: shareId,
			},
		},
	})
}

func (c *GrantClient) ShowGrantsOfAccountRole(t *testing.T, accountRoleId sdk.AccountObjectIdentifier) ([]sdk.Grant, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Show(ctx, &sdk.ShowGrantOptions{
		Of: &sdk.ShowGrantsOf{
			Role: accountRoleId,
		},
	})
}

func (c *GrantClient) ShowGrantsToAccountRole(t *testing.T, accountRoleId sdk.AccountObjectIdentifier) ([]sdk.Grant, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Show(ctx, &sdk.ShowGrantOptions{
		To: &sdk.ShowGrantsTo{
			Role: accountRoleId,
		},
	})
}

func (c *GrantClient) ShowGrantsOfDatabaseRole(t *testing.T, databaseRoleId sdk.DatabaseObjectIdentifier) ([]sdk.Grant, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Show(ctx, &sdk.ShowGrantOptions{
		Of: &sdk.ShowGrantsOf{
			DatabaseRole: databaseRoleId,
		},
	})
}

func (c *GrantClient) ShowGrantsToDatabaseRole(t *testing.T, databaseRoleId sdk.DatabaseObjectIdentifier) ([]sdk.Grant, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().Show(ctx, &sdk.ShowGrantOptions{
		To: &sdk.ShowGrantsTo{
			DatabaseRole: databaseRoleId,
		},
	})
}

func (c *GrantClient) GrantDatabaseRoleToUser(t *testing.T, databaseRoleId sdk.DatabaseObjectIdentifier, userId sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	// TODO(SNOW-2095669): Update when the client is updated to support this
	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf("GRANT DATABASE ROLE %s TO USER %s", databaseRoleId.FullyQualifiedName(), userId.FullyQualifiedName()))
	require.NoError(t, err)
}

func (c *GrantClient) GrantPrivilegesOnDatabaseToUser(t *testing.T, databaseId sdk.AccountObjectIdentifier, userId sdk.AccountObjectIdentifier, privileges ...sdk.AccountObjectPrivilege) {
	t.Helper()
	ctx := context.Background()

	// TODO(SNOW-2095669): Update when the client is updated to support this
	_, err := c.context.client.ExecForTests(ctx, fmt.Sprintf("GRANT %s ON DATABASE %s TO USER %s", collections.JoinStrings(privileges, ","), databaseId.FullyQualifiedName(), userId.FullyQualifiedName()))
	require.NoError(t, err)
}
