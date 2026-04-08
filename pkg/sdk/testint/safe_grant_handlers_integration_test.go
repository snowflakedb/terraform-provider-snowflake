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

	t.Run("privilege never granted", func(t *testing.T) {
		// Snowflake returns success (0 rows affected) when revoking a privilege that was never granted.
		// This validates that ErrObjectNotExistOrAuthorized is only produced for missing objects/roles,
		// not for already-absent grants.
		err := testClient(t).Grants.RevokePrivilegesFromAccountRole(ctx, tablePrivileges, tableOn(table.ID()), role.ID(), nil)
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
			privileges := &sdk.AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect},
			}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					SchemaObject: &sdk.Object{
						ObjectType: tt.ObjectType,
						Name:       NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema,
					},
				},
			}
			err := client.Grants.RevokePrivilegesFromAccountRoleSafely(ctx, privileges, on, role.ID(), nil)
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
		On         *sdk.AccountRoleGrantOn
	}{
		{ObjectType: sdk.ObjectTypeDatabase, On: &sdk.AccountRoleGrantOn{AccountObject: &sdk.GrantOnAccountObject{Database: sdk.Pointer(nonExistingId)}}},
		{ObjectType: sdk.ObjectTypeWarehouse, On: &sdk.AccountRoleGrantOn{AccountObject: &sdk.GrantOnAccountObject{Warehouse: sdk.Pointer(nonExistingId)}}},
		{ObjectType: sdk.ObjectTypeComputePool, On: &sdk.AccountRoleGrantOn{AccountObject: &sdk.GrantOnAccountObject{ComputePool: sdk.Pointer(nonExistingId)}}},
		{ObjectType: sdk.ObjectTypeExternalVolume, On: &sdk.AccountRoleGrantOn{AccountObject: &sdk.GrantOnAccountObject{ExternalVolume: sdk.Pointer(nonExistingId)}}},
		{ObjectType: sdk.ObjectTypeUser, On: &sdk.AccountRoleGrantOn{AccountObject: &sdk.GrantOnAccountObject{User: sdk.Pointer(nonExistingId)}}},
		{ObjectType: sdk.ObjectTypeResourceMonitor, On: &sdk.AccountRoleGrantOn{AccountObject: &sdk.GrantOnAccountObject{ResourceMonitor: sdk.Pointer(nonExistingId)}}},
		{ObjectType: sdk.ObjectTypeIntegration, On: &sdk.AccountRoleGrantOn{AccountObject: &sdk.GrantOnAccountObject{Integration: sdk.Pointer(nonExistingId)}}},
		{ObjectType: sdk.ObjectTypeFailoverGroup, On: &sdk.AccountRoleGrantOn{AccountObject: &sdk.GrantOnAccountObject{FailoverGroup: sdk.Pointer(nonExistingId)}}},
		{ObjectType: sdk.ObjectTypeReplicationGroup, On: &sdk.AccountRoleGrantOn{AccountObject: &sdk.GrantOnAccountObject{ReplicationGroup: sdk.Pointer(nonExistingId)}}},
	}

	for _, tt := range testCases {
		t.Run(tt.ObjectType.String(), func(t *testing.T) {
			privileges := &sdk.AccountRoleGrantPrivileges{
				AccountObjectPrivileges: []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage},
			}
			err := client.Grants.RevokePrivilegesFromAccountRoleSafely(ctx, privileges, tt.On, role.ID(), nil)
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
		Privileges *sdk.AccountRoleGrantPrivileges
		On         *sdk.AccountRoleGrantOn
	}{
		{
			Name:       "future tables in non-existing schema",
			Privileges: &sdk.AccountRoleGrantPrivileges{SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect}},
			On:         &sdk.AccountRoleGrantOn{SchemaObject: &sdk.GrantOnSchemaObject{Future: &sdk.GrantOnSchemaObjectIn{PluralObjectType: sdk.PluralObjectTypeTables, InSchema: sdk.Pointer(NonExistingDatabaseObjectIdentifier)}}},
		},
		{
			Name:       "future tables in non-existing database",
			Privileges: &sdk.AccountRoleGrantPrivileges{SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect}},
			On:         &sdk.AccountRoleGrantOn{SchemaObject: &sdk.GrantOnSchemaObject{Future: &sdk.GrantOnSchemaObjectIn{PluralObjectType: sdk.PluralObjectTypeTables, InDatabase: sdk.Pointer(NonExistingAccountObjectIdentifier)}}},
		},
		{
			Name:       "future schemas in non-existing database",
			Privileges: &sdk.AccountRoleGrantPrivileges{SchemaPrivileges: []sdk.SchemaPrivilege{sdk.SchemaPrivilegeUsage}},
			On:         &sdk.AccountRoleGrantOn{Schema: &sdk.GrantOnSchema{FutureSchemasInDatabase: sdk.Pointer(NonExistingAccountObjectIdentifier)}},
		},
		{
			Name:       "all tables in non-existing schema",
			Privileges: &sdk.AccountRoleGrantPrivileges{SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeSelect}},
			On:         &sdk.AccountRoleGrantOn{SchemaObject: &sdk.GrantOnSchemaObject{All: &sdk.GrantOnSchemaObjectIn{PluralObjectType: sdk.PluralObjectTypeTables, InSchema: sdk.Pointer(NonExistingDatabaseObjectIdentifier)}}},
		},
		{
			Name:       "all schemas in non-existing database",
			Privileges: &sdk.AccountRoleGrantPrivileges{SchemaPrivileges: []sdk.SchemaPrivilege{sdk.SchemaPrivilegeUsage}},
			On:         &sdk.AccountRoleGrantOn{Schema: &sdk.GrantOnSchema{AllSchemasInDatabase: sdk.Pointer(NonExistingAccountObjectIdentifier)}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := client.Grants.RevokePrivilegesFromAccountRoleSafely(ctx, tt.Privileges, tt.On, role.ID(), nil)
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
	err := client.Grants.GrantPrivilegesToAccountRole(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor},
		},
		&sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypePipes,
					InSchema:         sdk.Pointer(testClientHelper().Ids.SchemaId()),
				},
			},
		},
		role.ID(),
		&sdk.GrantPrivilegesToAccountRoleOptions{},
	)
	require.NoError(t, err)

	// Drop the role — pipes still exist, so Pipes.Show succeeds,
	// but each per-pipe REVOKE will fail with ErrObjectNotExistOrAuthorized.
	roleCleanup()

	// RevokePrivilegesFromAccountRoleSafely must suppress the per-pipe errors individually,
	// rather than swallowing a joined error that may also contain unexpected errors.
	err = client.Grants.RevokePrivilegesFromAccountRoleSafely(
		ctx,
		&sdk.AccountRoleGrantPrivileges{
			SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{sdk.SchemaObjectPrivilegeMonitor},
		},
		&sdk.AccountRoleGrantOn{
			SchemaObject: &sdk.GrantOnSchemaObject{
				All: &sdk.GrantOnSchemaObjectIn{
					PluralObjectType: sdk.PluralObjectTypePipes,
					InSchema:         sdk.Pointer(testClientHelper().Ids.SchemaId()),
				},
			},
		},
		role.ID(),
		nil,
	)
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
			_, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				On: &sdk.ShowGrantsOn{
					Object: &sdk.Object{
						ObjectType: tt.ObjectType,
						Name:       NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema,
					},
				},
			})
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
			_, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				On: &sdk.ShowGrantsOn{
					Object: tt.Object,
				},
			})
			assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		})
	}
}

func TestInt_ShowGrantsToNonExistingRole(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	_, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
		To: &sdk.ShowGrantsTo{
			Role: NonExistingAccountObjectIdentifier,
		},
	})
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}

func TestInt_ShowFutureGrantsInNonExistingContainer(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	testCases := []struct {
		Name string
		In   *sdk.ShowGrantsIn
	}{
		{
			Name: "in non-existing schema",
			In:   &sdk.ShowGrantsIn{Schema: sdk.Pointer(NonExistingDatabaseObjectIdentifier)},
		},
		{
			Name: "in non-existing database",
			In:   &sdk.ShowGrantsIn{Database: sdk.Pointer(NonExistingAccountObjectIdentifier)},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				Future: sdk.Bool(true),
				In:     tt.In,
			})
			assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
		})
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
		err := client.Grants.RevokePrivilegeFromShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, &sdk.ShareGrantOn{
			Table: &sdk.OnTable{
				AllInSchema: testClientHelper().Ids.SchemaId(),
			},
		}, share.ID())
		require.NoError(t, err)
	})

	t.Run("non-existing share", func(t *testing.T) {
		err := client.Grants.RevokePrivilegeFromShareSafely(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: testClientHelper().Ids.DatabaseId(),
		}, NonExistingAccountObjectIdentifier)
		require.NoError(t, err)
	})

	t.Run("non-existing database", func(t *testing.T) {
		err := client.Grants.RevokePrivilegeFromShareSafely(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{
			Database: NonExistingAccountObjectIdentifier,
		}, share.ID())
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
		On         *sdk.ShareGrantOn
	}{
		{Name: "non-existing schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, On: &sdk.ShareGrantOn{Schema: NonExistingDatabaseObjectIdentifier}},
		{Name: "non-existing schema in non-existing database", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, On: &sdk.ShareGrantOn{Schema: NonExistingDatabaseObjectIdentifierWithNonExistingDatabase}},
		{Name: "non-existing table", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: &sdk.ShareGrantOn{Table: &sdk.OnTable{Name: NonExistingSchemaObjectIdentifier}}},
		{Name: "non-existing table in non-existing database and schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: &sdk.ShareGrantOn{Table: &sdk.OnTable{Name: NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema}}},
		{Name: "all tables in non-existing schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: &sdk.ShareGrantOn{Table: &sdk.OnTable{AllInSchema: NonExistingDatabaseObjectIdentifier}}},
		{Name: "all tables in non-existing schema in non-existing database", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: &sdk.ShareGrantOn{Table: &sdk.OnTable{AllInSchema: NonExistingDatabaseObjectIdentifierWithNonExistingDatabase}}},
		{Name: "non-existing view", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: &sdk.ShareGrantOn{View: NonExistingSchemaObjectIdentifier}},
		{Name: "non-existing view in non-existing database and schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeSelect}, On: &sdk.ShareGrantOn{View: NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema}},
		{Name: "non-existing tag", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeRead}, On: &sdk.ShareGrantOn{Tag: NonExistingSchemaObjectIdentifier}},
		{Name: "non-existing tag in non-existing database and schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeRead}, On: &sdk.ShareGrantOn{Tag: NonExistingSchemaObjectIdentifierWithNonExistingDatabaseAndSchema}},
		{Name: "non-existing function", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, On: &sdk.ShareGrantOn{Function: NonExistingSchemaObjectIdentifierWithArguments}},
		{Name: "non-existing function in non-existing database and schema", Privileges: []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, On: &sdk.ShareGrantOn{Function: NonExistingSchemaObjectIdentifierWithArgumentsWithNonExistingDatabaseAndSchema}},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			err := client.Grants.RevokePrivilegeFromShareSafely(ctx, tt.Privileges, tt.On, share.ID())
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
		err := client.Roles.RevokeSafely(ctx, sdk.NewRevokeRoleRequest(NonExistingAccountObjectIdentifier, sdk.RevokeRole{Role: sdk.Pointer(parentRole.ID())}))
		assert.NoError(t, err)
	})

	t.Run("revoke from non-existing grantee role", func(t *testing.T) {
		err := client.Roles.RevokeSafely(ctx, sdk.NewRevokeRoleRequest(role.ID(), sdk.RevokeRole{Role: &NonExistingAccountObjectIdentifier}))
		assert.NoError(t, err)
	})

	t.Run("revoke role that was never granted", func(t *testing.T) {
		err := client.Roles.RevokeSafely(ctx, sdk.NewRevokeRoleRequest(NonExistingAccountObjectIdentifier, sdk.RevokeRole{Role: sdk.Pointer(NonExistingAccountObjectIdentifier)}))
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
