//go:build non_account_level_tests

package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SafeDropOnAccountObjectIdentifier(t *testing.T) {
	networkPolicy, cleanupNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(cleanupNetworkPolicy)

	ctx := context.Background()
	networkPolicyDrop := func() error {
		return testClient(t).NetworkPolicies.Drop(ctx, sdk.NewDropNetworkPolicyRequest(networkPolicy.ID()).WithIfExists(true))
	}

	err := sdk.SafeDrop(testClient(t), networkPolicyDrop, ctx, networkPolicy.ID())
	assert.NoError(t, err)

	err = sdk.SafeDrop(testClient(t), networkPolicyDrop, ctx, networkPolicy.ID())
	assert.NoError(t, err)
}

func TestInt_SafeDropOnDatabaseObjectIdentifier(t *testing.T) {
	databaseRole, cleanupDatabaseRole := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(cleanupDatabaseRole)

	ctx := context.Background()
	databaseRoleDrop := func(id sdk.DatabaseObjectIdentifier) func() error {
		return func() error {
			return testClient(t).DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id).WithIfExists(true))
		}
	}

	err := sdk.SafeDrop(testClient(t), databaseRoleDrop(databaseRole.ID()), testContext(t), databaseRole.ID())
	assert.NoError(t, err)

	invalidDatabaseRoleId := NonExistingDatabaseObjectIdentifierWithNonExistingDatabase
	err = sdk.SafeDrop(testClient(t), databaseRoleDrop(invalidDatabaseRoleId), testContext(t), invalidDatabaseRoleId)
	assert.NoError(t, err)
}

func TestInt_SafeDropOnSchemaObjectIdentifier(t *testing.T) {
	table, cleanupTable := testClientHelper().Table.Create(t)
	t.Cleanup(cleanupTable)

	ctx := context.Background()
	tableDrop := func(id sdk.SchemaObjectIdentifier) func() error {
		return func() error {
			return testClient(t).Tables.Drop(ctx, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))
		}
	}

	err := sdk.SafeDrop(testClient(t), tableDrop(table.ID()), ctx, table.ID())
	assert.NoError(t, err)

	invalidTableIdOnValidDatabase := NonExistingSchemaObjectIdentifierWithNonExistingSchema
	err = sdk.SafeDrop(testClient(t), tableDrop(invalidTableIdOnValidDatabase), ctx, invalidTableIdOnValidDatabase)
	assert.NoError(t, err)

	invalidTableId := NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema
	err = sdk.SafeDrop(testClient(t), tableDrop(invalidTableId), ctx, invalidTableId)
	assert.NoError(t, err)
}

func TestInt_SafeDropOnSchemaObjectIdentifierWithArguments(t *testing.T) {
	procedure, procedureCleanup := testClientHelper().Procedure.Create(t, sdk.DataTypeInt)
	t.Cleanup(procedureCleanup)

	ctx := context.Background()
	procedureDrop := func(id sdk.SchemaObjectIdentifierWithArguments) func() error {
		return func() error {
			return testClient(t).Procedures.Drop(ctx, sdk.NewDropProcedureRequest(id).WithIfExists(true))
		}
	}

	err := sdk.SafeDrop(testClient(t), procedureDrop(procedure.ID()), ctx, procedure.ID())
	assert.NoError(t, err)

	invalidProcedureIdOnValidDatabase := NonExistingSchemaObjectIdentifierWithArgumentsWithNonExistingSchema
	err = sdk.SafeDrop(testClient(t), procedureDrop(invalidProcedureIdOnValidDatabase), ctx, invalidProcedureIdOnValidDatabase)
	assert.NoError(t, err)

	invalidProcedureId := NonExistingSchemaObjectIdentifierWithArgumentsWithNonExistingDatabaseAndSchema
	err = sdk.SafeDrop(testClient(t), procedureDrop(invalidProcedureId), ctx, invalidProcedureId)
	assert.NoError(t, err)
}

func TestInt_SafeRevokePrivilegesFromAccountRole(t *testing.T) {
	client := testClient(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	ctx := context.Background()

	tablePrivileges := &sdk.AccountRoleGrantPrivileges{
		SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect},
	}
	tableOn := func(id sdk.SchemaObjectIdentifier) *sdk.AccountRoleGrantOn {
		return &sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				SchemaObject: &sdk.Object{
					ObjectType: sdk.ObjectTypeTable,
					Name:       id,
				},
			},
		}
	}

	revoke := func(privileges *sdk.AccountRoleGrantPrivileges, on *sdk.AccountRoleGrantOn, r sdk.AccountObjectIdentifier) error {
		return client.Grants.RevokePrivilegesFromAccountRoleSafely(ctx, privileges, on, r, nil)
	}

	t.Run("existing object", func(t *testing.T) {
		err := testClient(t).Grants.GrantPrivilegesToAccountRole(ctx, tablePrivileges, tableOn(table.ID()), role.ID(), nil)
		require.NoError(t, err)

		err = revoke(tablePrivileges, tableOn(table.ID()), role.ID())
		assert.NoError(t, err)
	})

	t.Run("privilege never granted", func(t *testing.T) {
		// Snowflake returns success (0 rows affected) when revoking a privilege that was never granted.
		// This validates that ErrObjectNotExistOrAuthorized is only produced for missing objects/roles,
		// not for already-absent grants.
		err := testClient(t).Grants.RevokePrivilegesFromAccountRole(ctx, tablePrivileges, tableOn(table.ID()), role.ID(), nil)
		assert.NoError(t, err)
	})

	t.Run("missing schema object", func(t *testing.T) {
		err := revoke(tablePrivileges, tableOn(NonExistingSchemaObjectIdentifier), role.ID())
		assert.NoError(t, err)

		err = revoke(tablePrivileges, tableOn(NonExistingSchemaObjectIdentifierWithNonExistingSchema), role.ID())
		assert.NoError(t, err)

		err = revoke(tablePrivileges, tableOn(NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema), role.ID())
		assert.NoError(t, err)
	})

	t.Run("missing schema", func(t *testing.T) {
		schemaPrivileges := &sdk.AccountRoleGrantPrivileges{
			SchemaPrivileges: []sdk.SchemaPrivilege{sdk.SchemaPrivilegeUsage},
		}
		on := &sdk.AccountRoleGrantOn{
			Schema: &sdk.GrantOnSchema{
				Schema: sdk.Pointer(NonExistingDatabaseObjectIdentifier),
			},
		}
		err := revoke(schemaPrivileges, on, role.ID())
		assert.NoError(t, err)

		on = &sdk.AccountRoleGrantOn{
			Schema: &sdk.GrantOnSchema{
				Schema: sdk.Pointer(NonExistingDatabaseObjectIdentifierWithNonExistingDatabase),
			},
		}
		err = revoke(schemaPrivileges, on, role.ID())
		assert.NoError(t, err)
	})

	t.Run("missing database", func(t *testing.T) {
		dbPrivileges := &sdk.AccountRoleGrantPrivileges{
			AccountObjectPrivileges: []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage},
		}
		on := &sdk.AccountRoleGrantOn{
			AccountObject: &sdk.GrantOnAccountObject{
				Database: sdk.Pointer(NonExistingAccountObjectIdentifier),
			},
		}
		err := revoke(dbPrivileges, on, role.ID())
		assert.NoError(t, err)
	})

	t.Run("missing account object", func(t *testing.T) {
		whPrivileges := &sdk.AccountRoleGrantPrivileges{
			AccountObjectPrivileges: []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage},
		}
		on := &sdk.AccountRoleGrantOn{
			AccountObject: &sdk.GrantOnAccountObject{
				Warehouse: sdk.Pointer(NonExistingAccountObjectIdentifier),
			},
		}
		err := revoke(whPrivileges, on, role.ID())
		assert.NoError(t, err)
	})

	t.Run("missing role", func(t *testing.T) {
		err := revoke(tablePrivileges, tableOn(table.ID()), NonExistingAccountObjectIdentifier)
		assert.NoError(t, err)
	})

	t.Run("insufficient privileges", func(t *testing.T) {
		// Create a role with no grants and switch to it.
		limitedRole, limitedRoleCleanup := testClientHelper().Role.CreateRoleGrantedToCurrentRole(t)
		t.Cleanup(limitedRoleCleanup)
		useRoleCleanup := testClientHelper().Role.UseRole(t, limitedRole.ID())
		t.Cleanup(useRoleCleanup)

		// The raw revoke should fail with an authorization error because the limited role
		// cannot see the table (ErrObjectNotExistOrAuthorized).
		rawErr := testClient(t).Grants.RevokePrivilegesFromAccountRole(ctx, tablePrivileges, tableOn(table.ID()), role.ID(), nil)
		assert.ErrorIs(t, rawErr, sdk.ErrObjectNotExistOrAuthorized)

		// RevokePrivilegesFromAccountRoleSafely treats this the same as a missing object and returns nil.
		err := revoke(tablePrivileges, tableOn(table.ID()), role.ID())
		assert.NoError(t, err)
	})
}

func TestInt_SafeRemoveProgrammaticAccessToken(t *testing.T) {
	user, cleanupUser := testClientHelper().User.CreateUser(t)
	t.Cleanup(cleanupUser)

	token, cleanupToken := testClientHelper().User.AddProgrammaticAccessToken(t, user.ID())
	t.Cleanup(cleanupToken)

	ctx := context.Background()
	err := sdk.SafeRemoveProgrammaticAccessToken(testClient(t), ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(user.ID(), token.ID()))
	assert.NoError(t, err)

	err = sdk.SafeRemoveProgrammaticAccessToken(testClient(t), ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(user.ID(), token.ID()))
	assert.NoError(t, err)

	invalidUserId := NonExistingAccountObjectIdentifier
	err = sdk.SafeRemoveProgrammaticAccessToken(testClient(t), ctx, sdk.NewRemoveUserProgrammaticAccessTokenRequest(invalidUserId, token.ID()))
	assert.NoError(t, err)
}
