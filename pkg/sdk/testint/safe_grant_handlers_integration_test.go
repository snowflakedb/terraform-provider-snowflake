//go:build non_account_level_tests

package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testvars"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_SafeRevokePrivilegesFromAccountRole(t *testing.T) {
	client := testClient(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	ctx := context.Background()

	tablePrivileges := sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect})
	tableOn := func(id sdk.SchemaObjectIdentifier) *sdk.AccountRoleGrantOnRequest {
		return sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithSchemaObject(sdk.Object{
			ObjectType: sdk.ObjectTypeTable,
			Name:       id,
		}))
	}

	revoke := func(privileges *sdk.AccountRoleGrantPrivilegesRequest, on *sdk.AccountRoleGrantOnRequest, r sdk.AccountObjectIdentifier) error {
		return client.Grants.RevokePrivilegesFromAccountRoleSafely(ctx, sdk.NewRevokePrivilegesFromAccountRoleRequest(r).WithPrivileges(*privileges).WithOn(*on))
	}

	t.Run("privilege never granted", func(t *testing.T) {
		// Snowflake returns success (0 rows affected) when revoking a privilege that was never granted.
		// This validates that ErrObjectNotExistOrAuthorized is only produced for missing objects/roles,
		// not for already-absent grants.
		err := testClient(t).Grants.RevokePrivilegesFromAccountRole(ctx, sdk.NewRevokePrivilegesFromAccountRoleRequest(role.ID()).WithPrivileges(*tablePrivileges).WithOn(*tableOn(table.ID())))
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
		rawErr := testClient(t).Grants.RevokePrivilegesFromAccountRole(ctx, sdk.NewRevokePrivilegesFromAccountRoleRequest(role.ID()).WithPrivileges(*tablePrivileges).WithOn(*tableOn(table.ID())))
		assert.ErrorIs(t, rawErr, sdk.ErrObjectNotExistOrAuthorized)

		// RevokePrivilegesFromAccountRoleSafely treats this the same as a missing object and returns nil.
		err := revoke(tablePrivileges, tableOn(table.ID()), role.ID())
		assert.NoError(t, err)
	})
}

func TestInt_SafeRevokeOnNonExistingSchemaObject(t *testing.T) {
	client := testClient(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	ctx := context.Background()

	testCases := []struct {
		ObjectType sdk.ObjectType
	}{
		{ObjectType: sdk.ObjectTypeTable},
		{ObjectType: sdk.ObjectTypeDynamicTable},
		{ObjectType: sdk.ObjectTypeCortexSearchService},
		{ObjectType: sdk.ObjectTypeExternalTable},
		{ObjectType: sdk.ObjectTypeEventTable},
		{ObjectType: sdk.ObjectTypeView},
		{ObjectType: sdk.ObjectTypeMaterializedView},
		{ObjectType: sdk.ObjectTypeSequence},
		{ObjectType: sdk.ObjectTypeStream},
		{ObjectType: sdk.ObjectTypeTask},
		{ObjectType: sdk.ObjectTypeMaskingPolicy},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy},
		{ObjectType: sdk.ObjectTypeTag},
		{ObjectType: sdk.ObjectTypeSecret},
		{ObjectType: sdk.ObjectTypeStage},
		{ObjectType: sdk.ObjectTypeFileFormat},
		{ObjectType: sdk.ObjectTypePipe},
		{ObjectType: sdk.ObjectTypeAlert},
		{ObjectType: sdk.ObjectTypeStreamlit},
		{ObjectType: sdk.ObjectTypeNetworkRule},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy},
		{ObjectType: sdk.ObjectTypeImageRepository},
		{ObjectType: sdk.ObjectTypeService},
		{ObjectType: sdk.ObjectTypeGitRepository},
		{ObjectType: sdk.ObjectTypeNotebook},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			privileges := sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect})
			on := sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithSchemaObject(sdk.Object{
				ObjectType: tt.ObjectType,
				Name:       NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema,
			}))
			err := client.Grants.RevokePrivilegesFromAccountRoleSafely(ctx, sdk.NewRevokePrivilegesFromAccountRoleRequest(role.ID()).WithPrivileges(*privileges).WithOn(*on))
			assert.NoError(t, err)
		})
	}
}

func TestInt_SafeRevokeOnNonExistingAccountObject(t *testing.T) {
	client := testClient(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	ctx := context.Background()
	nonExistingId := NonExistingAccountObjectIdentifier

	testCases := []struct {
		ObjectType sdk.ObjectType
		On         *sdk.AccountRoleGrantOnRequest
	}{
		{ObjectType: sdk.ObjectTypeDatabase, On: sdk.NewAccountRoleGrantOnRequest().WithAccountObject(*sdk.NewGrantOnAccountObjectRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: nonExistingId}))},
		{ObjectType: sdk.ObjectTypeWarehouse, On: sdk.NewAccountRoleGrantOnRequest().WithAccountObject(*sdk.NewGrantOnAccountObjectRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeWarehouse, Name: nonExistingId}))},
		{ObjectType: sdk.ObjectTypeComputePool, On: sdk.NewAccountRoleGrantOnRequest().WithAccountObject(*sdk.NewGrantOnAccountObjectRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeComputePool, Name: nonExistingId}))},
		{ObjectType: sdk.ObjectTypeExternalVolume, On: sdk.NewAccountRoleGrantOnRequest().WithAccountObject(*sdk.NewGrantOnAccountObjectRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeExternalVolume, Name: nonExistingId}))},
		{ObjectType: sdk.ObjectTypeUser, On: sdk.NewAccountRoleGrantOnRequest().WithAccountObject(*sdk.NewGrantOnAccountObjectRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeUser, Name: nonExistingId}))},
		{ObjectType: sdk.ObjectTypeResourceMonitor, On: sdk.NewAccountRoleGrantOnRequest().WithAccountObject(*sdk.NewGrantOnAccountObjectRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeResourceMonitor, Name: nonExistingId}))},
		{ObjectType: sdk.ObjectTypeIntegration, On: sdk.NewAccountRoleGrantOnRequest().WithAccountObject(*sdk.NewGrantOnAccountObjectRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeIntegration, Name: nonExistingId}))},
		{ObjectType: sdk.ObjectTypeFailoverGroup, On: sdk.NewAccountRoleGrantOnRequest().WithAccountObject(*sdk.NewGrantOnAccountObjectRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeFailoverGroup, Name: nonExistingId}))},
		{ObjectType: sdk.ObjectTypeReplicationGroup, On: sdk.NewAccountRoleGrantOnRequest().WithAccountObject(*sdk.NewGrantOnAccountObjectRequest().WithObject(sdk.Object{ObjectType: sdk.ObjectTypeReplicationGroup, Name: nonExistingId}))},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			privileges := sdk.NewAccountRoleGrantPrivilegesRequest().WithAccountObjectPrivileges([]sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage})
			err := client.Grants.RevokePrivilegesFromAccountRoleSafely(ctx, sdk.NewRevokePrivilegesFromAccountRoleRequest(role.ID()).WithPrivileges(*privileges).WithOn(*tt.On))
			assert.NoError(t, err)
		})
	}
}

func TestInt_SafeRevokeOnFutureGrantsInNonExistingObjectInHierarchy(t *testing.T) {
	client := testClient(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	ctx := context.Background()

	testCases := []struct {
		Name       string
		Privileges *sdk.AccountRoleGrantPrivilegesRequest
		On         *sdk.AccountRoleGrantOnRequest
	}{
		{
			Name:       "future tables in non-existing schema",
			Privileges: sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect}),
			On:         sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypeTables).WithInSchema(NonExistingDatabaseObjectIdentifier))),
		},
		{
			Name:       "future tables in non-existing database",
			Privileges: sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect}),
			On:         sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypeTables).WithInDatabase(NonExistingAccountObjectIdentifier))),
		},
		{
			Name:       "future schemas in non-existing database",
			Privileges: sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaPrivileges([]sdk.SchemaPrivilege{sdk.SchemaPrivilegeUsage}),
			On:         sdk.NewAccountRoleGrantOnRequest().WithSchema(*sdk.NewGrantOnSchemaRequest().WithFutureSchemasInDatabase(NonExistingAccountObjectIdentifier)),
		},
		{
			Name:       "all tables in non-existing schema",
			Privileges: sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect}),
			On:         sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypeTables).WithInSchema(NonExistingDatabaseObjectIdentifier))),
		},
		{
			Name:       "all schemas in non-existing database",
			Privileges: sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaPrivileges([]sdk.SchemaPrivilege{sdk.SchemaPrivilegeUsage}),
			On:         sdk.NewAccountRoleGrantOnRequest().WithSchema(*sdk.NewGrantOnSchemaRequest().WithAllSchemasInDatabase(NonExistingAccountObjectIdentifier)),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := client.Grants.RevokePrivilegesFromAccountRoleSafely(ctx, sdk.NewRevokePrivilegesFromAccountRoleRequest(role.ID()).WithPrivileges(*tt.Privileges).WithOn(*tt.On))
			assert.NoError(t, err)
		})
	}
}

func TestInt_SafeRevokeOnAllPipesWithMissingRole(t *testing.T) {
	client := testClient(t)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	copyStatement := createPipeCopyStatement(t, table, stage)

	_, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
	t.Cleanup(pipeCleanup)

	_, secondPipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
	t.Cleanup(secondPipeCleanup)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	ctx := context.Background()

	// Grant MONITOR on all pipes in schema to the role.
	err := client.Grants.GrantPrivilegesToAccountRole(ctx, sdk.NewGrantPrivilegesToAccountRoleRequest(role.ID()).WithPrivileges(*sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor})).WithOn(*sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypePipes).WithInSchema(testClientHelper().Ids.SchemaId())))))
	require.NoError(t, err)

	// Drop the role — pipes still exist, so Pipes.Show succeeds,
	// but each per-pipe REVOKE will fail with ErrObjectNotExistOrAuthorized.
	roleCleanup()

	// RevokePrivilegesFromAccountRoleSafely must suppress the per-pipe errors individually,
	// rather than swallowing a joined error that may also contain unexpected errors.
	err = client.Grants.RevokePrivilegesFromAccountRoleSafely(ctx, sdk.NewRevokePrivilegesFromAccountRoleRequest(role.ID()).WithPrivileges(*sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor})).WithOn(*sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypePipes).WithInSchema(testClientHelper().Ids.SchemaId())))))
	assert.NoError(t, err)
}

func TestInt_ShowGrantsOnNonExistingSchemaObject(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	testCases := []struct {
		ObjectType sdk.ObjectType
	}{
		{ObjectType: sdk.ObjectTypeTable},
		{ObjectType: sdk.ObjectTypeDynamicTable},
		{ObjectType: sdk.ObjectTypeCortexSearchService},
		{ObjectType: sdk.ObjectTypeExternalTable},
		{ObjectType: sdk.ObjectTypeEventTable},
		{ObjectType: sdk.ObjectTypeView},
		{ObjectType: sdk.ObjectTypeMaterializedView},
		{ObjectType: sdk.ObjectTypeSequence},
		{ObjectType: sdk.ObjectTypeStream},
		{ObjectType: sdk.ObjectTypeTask},
		{ObjectType: sdk.ObjectTypeMaskingPolicy},
		{ObjectType: sdk.ObjectTypeRowAccessPolicy},
		{ObjectType: sdk.ObjectTypeTag},
		{ObjectType: sdk.ObjectTypeSecret},
		{ObjectType: sdk.ObjectTypeStage},
		{ObjectType: sdk.ObjectTypeFileFormat},
		{ObjectType: sdk.ObjectTypePipe},
		{ObjectType: sdk.ObjectTypeAlert},
		{ObjectType: sdk.ObjectTypeStreamlit},
		{ObjectType: sdk.ObjectTypeNetworkRule},
		{ObjectType: sdk.ObjectTypeAuthenticationPolicy},
		{ObjectType: sdk.ObjectTypeImageRepository},
		{ObjectType: sdk.ObjectTypeService},
		{ObjectType: sdk.ObjectTypeGitRepository},
		{ObjectType: sdk.ObjectTypeNotebook},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			_, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
				ObjectType: tt.ObjectType,
				Name:       NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema,
			})))
			assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		})
	}
}

func TestInt_ShowGrantsOnNonExistingAccountObject(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()
	nonExistingId := NonExistingAccountObjectIdentifier

	testCases := []struct {
		ObjectType sdk.ObjectType
		Object     *sdk.Object
	}{
		{ObjectType: sdk.ObjectTypeDatabase, Object: &sdk.Object{ObjectType: sdk.ObjectTypeDatabase, Name: nonExistingId}},
		{ObjectType: sdk.ObjectTypeWarehouse, Object: &sdk.Object{ObjectType: sdk.ObjectTypeWarehouse, Name: nonExistingId}},
		{ObjectType: sdk.ObjectTypeComputePool, Object: &sdk.Object{ObjectType: sdk.ObjectTypeComputePool, Name: nonExistingId}},
		{ObjectType: sdk.ObjectTypeExternalVolume, Object: &sdk.Object{ObjectType: sdk.ObjectTypeExternalVolume, Name: nonExistingId}},
		{ObjectType: sdk.ObjectTypeUser, Object: &sdk.Object{ObjectType: sdk.ObjectTypeUser, Name: nonExistingId}},
		{ObjectType: sdk.ObjectTypeResourceMonitor, Object: &sdk.Object{ObjectType: sdk.ObjectTypeResourceMonitor, Name: nonExistingId}},
		{ObjectType: sdk.ObjectTypeIntegration, Object: &sdk.Object{ObjectType: sdk.ObjectTypeIntegration, Name: nonExistingId}},
		{ObjectType: sdk.ObjectTypeFailoverGroup, Object: &sdk.Object{ObjectType: sdk.ObjectTypeFailoverGroup, Name: nonExistingId}},
		{ObjectType: sdk.ObjectTypeReplicationGroup, Object: &sdk.Object{ObjectType: sdk.ObjectTypeReplicationGroup, Name: nonExistingId}},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			_, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(*tt.Object)))
			assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		})
	}
}

func TestInt_ShowGrantsToNonExistingRole(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	_, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithTo(*sdk.NewShowGrantsToRequest().WithRole(NonExistingAccountObjectIdentifier)))
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}

func TestInt_ShowFutureAndInheritedGrantsInNonExistingContainer(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	containers := []struct {
		Name string
		In   *sdk.ShowGrantsInRequest
	}{
		{
			Name: "in non-existing schema",
			In:   sdk.NewShowGrantsInRequest().WithSchema(NonExistingDatabaseObjectIdentifier),
		},
		{
			Name: "in non-existing database",
			In:   sdk.NewShowGrantsInRequest().WithDatabase(NonExistingAccountObjectIdentifier),
		},
	}

	grantTypes := []struct {
		Name      string
		Future    *bool
		Inherited *bool
	}{
		{Name: "future", Future: new(true)},
		{Name: "inherited", Inherited: new(true)},
	}

	for _, grantType := range grantTypes {
		for _, container := range containers {
			t.Run(grantType.Name+" "+container.Name, func(t *testing.T) {
				req := sdk.NewShowGrantsRequest().WithIn(*container.In)
				if grantType.Future != nil {
					req.WithFuture(*grantType.Future)
				}
				if grantType.Inherited != nil {
					req.WithInherited(*grantType.Inherited)
				}
				_, err := client.Grants.Show(ctx, req)
				assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
			})
		}
	}
}

func TestInt_SafeRevokePrivilegeFromShare(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	share, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	// Grant USAGE on the test database to the share (prerequisite for schema-level grants).
	revokeGrant := testClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, testClientHelper().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})
	t.Cleanup(revokeGrant)

	t.Run("privilege never granted", func(t *testing.T) {
		// Snowflake returns success (0 rows affected) when revoking a privilege that was never granted.
		err := client.Grants.RevokePrivilegeFromShare(ctx, sdk.NewRevokePrivilegeFromShareRequest(share.ID()).WithPrivileges([]sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}).WithOn(*sdk.NewShareGrantOnRequest().WithTable(*sdk.NewOnTableRequest().WithAllInSchema(testClientHelper().Ids.SchemaId()))))
		require.NoError(t, err)
	})

	t.Run("non-existing share", func(t *testing.T) {
		err := client.Grants.RevokePrivilegeFromShareSafely(ctx, sdk.NewRevokePrivilegeFromShareRequest(NonExistingAccountObjectIdentifier).WithPrivileges([]sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}).WithOn(*sdk.NewShareGrantOnRequest().WithDatabase(testClientHelper().Ids.DatabaseId())))
		require.NoError(t, err)
	})

	t.Run("non-existing database", func(t *testing.T) {
		err := client.Grants.RevokePrivilegeFromShareSafely(ctx, sdk.NewRevokePrivilegeFromShareRequest(share.ID()).WithPrivileges([]sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}).WithOn(*sdk.NewShareGrantOnRequest().WithDatabase(NonExistingAccountObjectIdentifier)))
		require.NoError(t, err)
	})
}

func TestInt_SafeRevokeFromShareOnNonExistingSchemaLevelObjects(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	share, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	// Grant USAGE on the test database to the share (prerequisite for schema-level grants).
	revokeGrant := testClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, testClientHelper().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage})
	t.Cleanup(revokeGrant)

	testCases := []struct {
		Name       string
		Privileges []sdk.ObjectPrivilege
		On         *sdk.ShareGrantOnRequest
	}{
		{Name: "non-existing schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, On: sdk.NewShareGrantOnRequest().WithSchema(NonExistingDatabaseObjectIdentifier)},
		{Name: "non-existing schema in non-existing database", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, On: sdk.NewShareGrantOnRequest().WithSchema(NonExistingDatabaseObjectIdentifierWithNonExistingDatabase)},
		{Name: "non-existing table", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: sdk.NewShareGrantOnRequest().WithTable(*sdk.NewOnTableRequest().WithName(NonExistingSchemaObjectIdentifier))},
		{Name: "non-existing table in non-existing database and schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: sdk.NewShareGrantOnRequest().WithTable(*sdk.NewOnTableRequest().WithName(NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema))},
		{Name: "all tables in non-existing schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: sdk.NewShareGrantOnRequest().WithTable(*sdk.NewOnTableRequest().WithAllInSchema(NonExistingDatabaseObjectIdentifier))},
		{Name: "all tables in non-existing schema in non-existing database", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: sdk.NewShareGrantOnRequest().WithTable(*sdk.NewOnTableRequest().WithAllInSchema(NonExistingDatabaseObjectIdentifierWithNonExistingDatabase))},
		{Name: "non-existing view", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: sdk.NewShareGrantOnRequest().WithView(NonExistingSchemaObjectIdentifier)},
		{Name: "non-existing view in non-existing database and schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: sdk.NewShareGrantOnRequest().WithView(NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema)},
		{Name: "non-existing tag", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeRead}, On: sdk.NewShareGrantOnRequest().WithTag(NonExistingSchemaObjectIdentifier)},
		{Name: "non-existing tag in non-existing database and schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeRead}, On: sdk.NewShareGrantOnRequest().WithTag(NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema)},
		{Name: "non-existing function", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, On: sdk.NewShareGrantOnRequest().WithFunction(NonExistingSchemaObjectIdentifierWithArguments)},
		{Name: "non-existing function in non-existing database and schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, On: sdk.NewShareGrantOnRequest().WithFunction(NonExistingSchemaObjectIdentifierWithArgumentsWithNonExistingDatabaseAndSchema)},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := client.Grants.RevokePrivilegeFromShareSafely(ctx, sdk.NewRevokePrivilegeFromShareRequest(share.ID()).WithPrivileges(tt.Privileges).WithOn(*tt.On))
			assert.NoError(t, err)
		})
	}
}

func TestInt_SafeRevokeAccountRole(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	parentRole, parentRoleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(parentRoleCleanup)

	t.Run("revoke non-existing role", func(t *testing.T) {
		err := client.Roles.RevokeSafely(ctx, sdk.NewRevokeRoleRequest(NonExistingAccountObjectIdentifier, *sdk.NewRevokeRoleFromRequest().WithRole(parentRole.ID())))
		assert.NoError(t, err)
	})

	t.Run("revoke from non-existing grantee role", func(t *testing.T) {
		err := client.Roles.RevokeSafely(ctx, sdk.NewRevokeRoleRequest(role.ID(), *sdk.NewRevokeRoleFromRequest().WithRole(NonExistingAccountObjectIdentifier)))
		assert.NoError(t, err)
	})

	t.Run("revoke role that was never granted", func(t *testing.T) {
		err := client.Roles.RevokeSafely(ctx, sdk.NewRevokeRoleRequest(NonExistingAccountObjectIdentifier, *sdk.NewRevokeRoleFromRequest().WithRole(NonExistingAccountObjectIdentifier)))
		assert.NoError(t, err)
	})
}

func TestInt_SafeRevokeDatabaseRole(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	accountRole, accountRoleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(accountRoleCleanup)

	t.Run("revoke non-existing database role from account role", func(t *testing.T) {
		err := client.DatabaseRoles.RevokeSafely(ctx, sdk.NewRevokeDatabaseRoleRequest(NonExistingDatabaseObjectIdentifier).WithAccountRole(accountRole.ID()))
		assert.NoError(t, err)
	})

	t.Run("revoke database role from non-existing account role", func(t *testing.T) {
		err := client.DatabaseRoles.RevokeSafely(ctx, sdk.NewRevokeDatabaseRoleRequest(databaseRole.ID()).WithAccountRole(NonExistingAccountObjectIdentifier))
		assert.NoError(t, err)
	})

	t.Run("revoke database role that was never granted", func(t *testing.T) {
		err := client.DatabaseRoles.RevokeSafely(ctx, sdk.NewRevokeDatabaseRoleRequest(NonExistingDatabaseObjectIdentifier).WithAccountRole(NonExistingAccountObjectIdentifier))
		assert.NoError(t, err)
	})
}

func TestInt_SafeRevokeDatabaseRoleFromShare(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(databaseRoleCleanup)

	share, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	t.Run("revoke non-existing database role from share", func(t *testing.T) {
		err := client.DatabaseRoles.RevokeFromShareSafely(ctx, sdk.NewRevokeFromShareDatabaseRoleRequest(NonExistingDatabaseObjectIdentifier, share.ID()))
		assert.NoError(t, err)
	})

	t.Run("revoke database role from non-existing share", func(t *testing.T) {
		err := client.DatabaseRoles.RevokeFromShareSafely(ctx, sdk.NewRevokeFromShareDatabaseRoleRequest(databaseRole.ID(), NonExistingAccountObjectIdentifier))
		assert.NoError(t, err)
	})

	t.Run("revoke non-existing database role from non-existing share", func(t *testing.T) {
		err := client.DatabaseRoles.RevokeFromShareSafely(ctx, sdk.NewRevokeFromShareDatabaseRoleRequest(NonExistingDatabaseObjectIdentifier, NonExistingAccountObjectIdentifier))
		assert.NoError(t, err)
	})
}

func TestInt_SafeRevokeApplicationRole(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	app := createApp(t)
	applicationRoleName := testvars.ApplicationRole1
	applicationRoleId := sdk.NewDatabaseObjectIdentifier(app.Name, applicationRoleName)

	accountRole, accountRoleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(accountRoleCleanup)

	t.Run("revoke non-existing application role from account role", func(t *testing.T) {
		err := client.ApplicationRoles.RevokeSafely(ctx, sdk.NewRevokeApplicationRoleRequest(NonExistingDatabaseObjectIdentifier).WithFrom(*sdk.NewKindOfRoleRequest().WithRoleName(accountRole.ID())))
		assert.NoError(t, err)
	})

	t.Run("revoke application role from non-existing account role", func(t *testing.T) {
		err := client.ApplicationRoles.RevokeSafely(ctx, sdk.NewRevokeApplicationRoleRequest(applicationRoleId).WithFrom(*sdk.NewKindOfRoleRequest().WithRoleName(accountRole.ID())))
		assert.NoError(t, err)
	})

	t.Run("revoke application role that was never granted", func(t *testing.T) {
		err := client.ApplicationRoles.RevokeSafely(ctx, sdk.NewRevokeApplicationRoleRequest(NonExistingDatabaseObjectIdentifier).WithFrom(*sdk.NewKindOfRoleRequest().WithRoleName(NonExistingAccountObjectIdentifier)))
		assert.NoError(t, err)
	})
}

func TestInt_SafeRevokePrivilegesFromDatabaseRole(t *testing.T) {
	client := testClient(t)

	dbRole, dbRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	ctx := context.Background()

	tablePrivileges := sdk.NewDatabaseRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect})
	tableOn := func(id sdk.SchemaObjectIdentifier) *sdk.DatabaseRoleGrantOnRequest {
		return sdk.NewDatabaseRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithSchemaObject(sdk.Object{
			ObjectType: sdk.ObjectTypeTable,
			Name:       id,
		}))
	}

	testCases := []struct {
		Name string
		On   *sdk.DatabaseRoleGrantOnRequest
		Role sdk.DatabaseObjectIdentifier
	}{
		{Name: "missing database role", On: tableOn(table.ID()), Role: NonExistingDatabaseObjectIdentifier},
		{Name: "missing database", On: tableOn(NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema), Role: dbRole.ID()},
		{Name: "missing schema", On: tableOn(NonExistingSchemaObjectIdentifierWithNonExistingSchema), Role: dbRole.ID()},
		{Name: "missing schema object", On: tableOn(NonExistingSchemaObjectIdentifier), Role: dbRole.ID()},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := client.Grants.RevokePrivilegesFromDatabaseRoleSafely(ctx, sdk.NewRevokePrivilegesFromDatabaseRoleRequest(tt.Role).WithPrivileges(*tablePrivileges).WithOn(*tt.On))
			assert.NoError(t, err)
		})
	}
}

func TestInt_SafeRevokePrivilegesFromDatabaseRole_AllPipesWithMissingRole(t *testing.T) {
	client := testClient(t)

	table, tableCleanup := testClientHelper().Table.Create(t)
	t.Cleanup(tableCleanup)

	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	copyStatement := createPipeCopyStatement(t, table, stage)

	_, pipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
	t.Cleanup(pipeCleanup)

	_, secondPipeCleanup := testClientHelper().Pipe.CreatePipe(t, copyStatement)
	t.Cleanup(secondPipeCleanup)

	dbRole, dbRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	ctx := context.Background()

	// Grant MONITOR on all pipes in schema to the database role.
	err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, sdk.NewGrantPrivilegesToDatabaseRoleRequest(dbRole.ID()).WithPrivileges(*sdk.NewDatabaseRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor})).WithOn(*sdk.NewDatabaseRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypePipes).WithInSchema(testClientHelper().Ids.SchemaId())))))
	require.NoError(t, err)

	// Drop the database role — pipes still exist, so Pipes.Show succeeds,
	// but each per-pipe REVOKE will fail with ErrObjectNotExistOrAuthorized.
	dbRoleCleanup()

	// RevokePrivilegesFromDatabaseRoleSafely must suppress the per-pipe errors individually.
	err = client.Grants.RevokePrivilegesFromDatabaseRoleSafely(ctx, sdk.NewRevokePrivilegesFromDatabaseRoleRequest(dbRole.ID()).WithPrivileges(*sdk.NewDatabaseRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor})).WithOn(*sdk.NewDatabaseRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(sdk.PluralObjectTypePipes).WithInSchema(testClientHelper().Ids.SchemaId())))))
	assert.NoError(t, err)
}

func TestInt_SafeRevokeInheritedPrivilegesFromAccountRole(t *testing.T) {
	client := testClient(t)

	role, roleCleanup := testClientHelper().Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	ctx := context.Background()

	databaseId := testClientHelper().Ids.DatabaseId()

	privileges := sdk.NewInheritedAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect})

	testCases := []struct {
		Name string
		In   sdk.InheritedAccountRoleGrantInRequest
		Role sdk.AccountObjectIdentifier
	}{
		{Name: "missing account role", In: *sdk.NewInheritedAccountRoleGrantInRequest().WithDatabase(databaseId), Role: NonExistingAccountObjectIdentifier},
		{Name: "missing database", In: *sdk.NewInheritedAccountRoleGrantInRequest().WithDatabase(NonExistingAccountObjectIdentifier), Role: role.ID()},
		{Name: "missing schema", In: *sdk.NewInheritedAccountRoleGrantInRequest().WithSchema(NonExistingDatabaseObjectIdentifier), Role: role.ID()},
		{Name: "missing schema in missing database", In: *sdk.NewInheritedAccountRoleGrantInRequest().WithSchema(NonExistingDatabaseObjectIdentifierWithNonExistingDatabase), Role: role.ID()},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := client.Grants.RevokeInheritedPrivilegesFromAccountRoleSafely(ctx, sdk.NewRevokeInheritedPrivilegesFromAccountRoleRequest(sdk.PluralObjectTypeTables, tt.In, tt.Role).WithPrivileges(*privileges))
			assert.NoError(t, err)
		})
	}
}

func TestInt_SafeRevokeInheritedPrivilegesFromDatabaseRole(t *testing.T) {
	client := testClient(t)

	dbRole, dbRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	ctx := context.Background()

	databaseId := testClientHelper().Ids.DatabaseId()

	privileges := sdk.NewInheritedDatabaseRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect})

	testCases := []struct {
		Name string
		In   sdk.InheritedDatabaseRoleGrantInRequest
		Role sdk.DatabaseObjectIdentifier
	}{
		{Name: "missing database role", In: *sdk.NewInheritedDatabaseRoleGrantInRequest().WithDatabase(databaseId), Role: NonExistingDatabaseObjectIdentifier},
		{Name: "missing database", In: *sdk.NewInheritedDatabaseRoleGrantInRequest().WithDatabase(NonExistingAccountObjectIdentifier), Role: dbRole.ID()},
		{Name: "missing schema", In: *sdk.NewInheritedDatabaseRoleGrantInRequest().WithSchema(NonExistingDatabaseObjectIdentifier), Role: dbRole.ID()},
		{Name: "missing schema in missing database", In: *sdk.NewInheritedDatabaseRoleGrantInRequest().WithSchema(NonExistingDatabaseObjectIdentifierWithNonExistingDatabase), Role: dbRole.ID()},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := client.Grants.RevokeInheritedPrivilegesFromDatabaseRoleSafely(ctx, sdk.NewRevokeInheritedPrivilegesFromDatabaseRoleRequest(sdk.PluralObjectTypeTables, tt.In, tt.Role).WithPrivileges(*privileges))
			assert.NoError(t, err)
		})
	}
}
