//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_GrantPrivileges_OnFutureAndAll_NewObjectTypes(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	type testCase struct {
		objectTypePlural sdk.PluralObjectType
		expectedGrantOn  sdk.ObjectType
		privilege        sdk.SchemaObjectPrivilege
		createObject     func(t *testing.T) (sdk.SchemaObjectIdentifier, func())
	}

	testCases := []testCase{
		{
			objectTypePlural: sdk.PluralObjectTypeWorkspaces,
			expectedGrantOn:  sdk.ObjectTypeWorkspace,
			privilege:        sdk.SchemaObjectPrivilegeRead,
			createObject: func(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
				t.Helper()
				return testClientHelper().Workspace.Create(t)
			},
		},
		{
			objectTypePlural: sdk.PluralObjectTypeInteractiveTables,
			expectedGrantOn:  sdk.ObjectTypeInteractiveTable,
			privilege:        sdk.SchemaObjectPrivilegeSelect,
			createObject: func(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
				t.Helper()
				it, cleanup := testClientHelper().Table.CreateInteractiveTable(t)
				return it.ID(), cleanup
			},
		},
	}

	for _, tc := range testCases {
		// --- FUTURE GRANTS ---

		t.Run("account role - future "+tc.objectTypePlural.String(), func(t *testing.T) {
			database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)

			privileges := sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{tc.privilege})
			on := sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInDatabase(database.ID())))
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, sdk.NewGrantPrivilegesToAccountRoleRequest(role.ID()).WithPrivileges(*privileges).WithOn(*on))
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithFuture(true).WithTo(*sdk.NewShowGrantsToRequest().WithRole(role.ID())))
			require.NoError(t, err)
			require.Len(t, grants, 1)
			assert.Equal(t, tc.privilege.String(), grants[0].Privilege)
			assert.Equal(t, tc.expectedGrantOn, grants[0].GrantOn)
			assert.Equal(t, sdk.ObjectTypeRole, grants[0].GrantTo)
			assert.Equal(t, role.ID().Name(), grants[0].GranteeName.Name())
		})

		t.Run("database role - future "+tc.objectTypePlural.String(), func(t *testing.T) {
			database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRoleInDatabase(t, database.ID())
			t.Cleanup(databaseRoleCleanup)

			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, sdk.NewGrantPrivilegesToDatabaseRoleRequest(databaseRole.ID()).WithPrivileges(*sdk.NewDatabaseRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{tc.privilege})).WithOn(*sdk.NewDatabaseRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInDatabase(database.ID())))))
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithFuture(true).WithTo(*sdk.NewShowGrantsToRequest().WithDatabaseRole(databaseRole.ID())))
			require.NoError(t, err)
			require.Len(t, grants, 1)
			assert.Equal(t, tc.privilege.String(), grants[0].Privilege)
			assert.Equal(t, tc.expectedGrantOn, grants[0].GrantOn)
			assert.Equal(t, sdk.ObjectTypeDatabaseRole, grants[0].GrantTo)
		})

		t.Run("ownership - future "+tc.objectTypePlural.String(), func(t *testing.T) {
			database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)
			roleId := role.ID()

			err := client.Grants.GrantOwnership(ctx, sdk.NewGrantOwnershipRequest(*sdk.NewOwnershipGrantOnRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInDatabase(database.ID())), *sdk.NewOwnershipGrantToRequest().WithAccountRoleName(roleId)))
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithFuture(true).WithTo(*sdk.NewShowGrantsToRequest().WithRole(roleId)))
			require.NoError(t, err)
			require.Len(t, grants, 1)
			assert.Equal(t, sdk.SchemaObjectOwnership.String(), grants[0].Privilege)
			assert.Equal(t, tc.expectedGrantOn, grants[0].GrantOn)
			assert.Equal(t, sdk.ObjectTypeRole, grants[0].GrantTo)
			assert.Equal(t, roleId, grants[0].GranteeName)
		})

		// --- ALL GRANTS ---

		t.Run("account role - all "+tc.objectTypePlural.String(), func(t *testing.T) {
			objectId, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)

			err := client.Grants.GrantPrivilegesToAccountRole(ctx, sdk.NewGrantPrivilegesToAccountRoleRequest(role.ID()).WithPrivileges(*sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{tc.privilege})).WithOn(*sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInSchema(testClientHelper().Ids.SchemaId())))))
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
				ObjectType: tc.expectedGrantOn,
				Name:       objectId,
			})))
			require.NoError(t, err)
			found := false
			for _, g := range grants {
				if g.Privilege == tc.privilege.String() && g.GranteeName.Name() == role.ID().Name() {
					found = true
					assert.Equal(t, tc.expectedGrantOn, g.GrantedOn)
					break
				}
			}
			assert.True(t, found, "expected privilege %s granted to %s on %s", tc.privilege, role.ID().Name(), objectId.FullyQualifiedName())
		})

		t.Run("database role - all "+tc.objectTypePlural.String(), func(t *testing.T) {
			objectId, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
			t.Cleanup(databaseRoleCleanup)

			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, sdk.NewGrantPrivilegesToDatabaseRoleRequest(databaseRole.ID()).WithPrivileges(*sdk.NewDatabaseRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{tc.privilege})).WithOn(*sdk.NewDatabaseRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInSchema(testClientHelper().Ids.SchemaId())))))
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
				ObjectType: tc.expectedGrantOn,
				Name:       objectId,
			})))
			require.NoError(t, err)
			found := false
			for _, g := range grants {
				if g.Privilege == tc.privilege.String() && g.GranteeName.Name() == databaseRole.ID().Name() {
					found = true
					assert.Equal(t, tc.expectedGrantOn, g.GrantedOn)
					break
				}
			}
			assert.True(t, found, "expected privilege %s granted to %s on %s", tc.privilege, databaseRole.ID().Name(), objectId.FullyQualifiedName())
		})

		t.Run("ownership - all "+tc.objectTypePlural.String(), func(t *testing.T) {
			objectId, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)
			roleId := role.ID()

			err := client.Grants.GrantOwnership(ctx, sdk.NewGrantOwnershipRequest(*sdk.NewOwnershipGrantOnRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInSchema(testClientHelper().Ids.SchemaId())), *sdk.NewOwnershipGrantToRequest().WithAccountRoleName(roleId)))
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, sdk.NewShowGrantsRequest().WithOn(*sdk.NewShowGrantsOnRequest().WithObject(sdk.Object{
				ObjectType: tc.expectedGrantOn,
				Name:       objectId,
			})))
			require.NoError(t, err)
			found := false
			for _, g := range grants {
				if g.Privilege == sdk.SchemaObjectOwnership.String() && g.GranteeName.Name() == roleId.Name() {
					found = true
					assert.Equal(t, tc.expectedGrantOn, g.GrantedOn)
					break
				}
			}
			assert.True(t, found, "expected privilege OWNERSHIP granted to %s on %s", roleId.Name(), objectId.FullyQualifiedName())
		})
	}
}

// TestInt_GrantPrivileges_OnFutureAndAll_UnsupportedObjectTypes verifies that granting on future/all
// for specific object types results in an error from Snowflake (these object types are not yet
// supported for bulk grants despite being listed in the docs).
func TestInt_GrantPrivileges_OnFutureAndAll_UnsupportedObjectTypes(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	type testCase struct {
		objectTypePlural    sdk.PluralObjectType
		privilege           sdk.SchemaObjectPrivilege
		createObject        func(t *testing.T) (sdk.ObjectIdentifier, func())
		expectedFutureError string
		expectedAllError    string
	}

	testCases := []testCase{
		{
			objectTypePlural: sdk.PluralObjectTypeExperiments,
			privilege:        sdk.SchemaObjectPrivilegeUsage,
			createObject: func(t *testing.T) (sdk.ObjectIdentifier, func()) {
				t.Helper()
				return testClientHelper().Experiment.Create(t)
			},
			expectedFutureError: "Unsupported feature 'EXPERIMENT'",
		},
		{
			objectTypePlural: sdk.PluralObjectTypeGateways,
			privilege:        sdk.SchemaObjectPrivilegeUsage,
			createObject: func(t *testing.T) (sdk.ObjectIdentifier, func()) {
				t.Helper()

				computePool, computePoolCleanup := testClientHelper().ComputePool.Create(t)
				t.Cleanup(computePoolCleanup)

				endpointName := "endpoint"
				spec := testClientHelper().Service.SampleSpecWithEndpoint(t, endpointName)
				serviceId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
				_, serviceCleanup := testClientHelper().Service.CreateWithRequest(
					t,
					sdk.NewCreateServiceRequest(serviceId, computePool.ID()).
						WithFromSpecification(*sdk.NewServiceFromSpecificationRequest().WithSpecification(spec)),
				)
				t.Cleanup(serviceCleanup)

				// TODO [SNOW-3825229]: uncomment when the below is fixed
				// testClientHelper().Service.WaitForStatus(t, serviceId, sdk.ServiceStatusRunning, 90*time.Second)
				testClientHelper().Gateway.GrantUsageOfAllServiceEndpointsToRole(t, serviceId, snowflakeroles.Accountadmin)

				// TODO [SNOW-3825229]: fix it
				// err: 398529 (02000): Service specified in gateway does not exist or not authorized: <FQN>
				// return testClientHelper().Gateway.Create(t, serviceId, endpointName)
				return sdk.SchemaObjectIdentifier{}, func() {}
			},
			expectedFutureError: "syntax error line 0 at position 0 unexpected 'TOK_GATEWAY'",
			expectedAllError:    "syntax error line 0 at position 0 unexpected 'TOK_GATEWAY'",
		},
		{
			objectTypePlural: sdk.PluralObjectTypeSnowflakeIntelligences,
			privilege:        sdk.SchemaObjectPrivilegeUsage,
			createObject: func(t *testing.T) (sdk.ObjectIdentifier, func()) {
				t.Helper()
				return testClientHelper().SnowflakeIntelligence.Create(t)
			},
			expectedFutureError: "unexpected 'SNOWFLAKE'.",
			expectedAllError:    "unexpected 'INTELLIGENCES'.",
		},
	}

	for _, tc := range testCases {
		t.Run("account role - future "+tc.objectTypePlural.String(), func(t *testing.T) {
			database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)

			err := client.Grants.GrantPrivilegesToAccountRole(ctx, sdk.NewGrantPrivilegesToAccountRoleRequest(role.ID()).WithPrivileges(*sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{tc.privilege})).WithOn(*sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInDatabase(database.ID())))))
			require.ErrorContains(t, err, tc.expectedFutureError)
		})

		t.Run("database role - future "+tc.objectTypePlural.String(), func(t *testing.T) {
			database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRoleInDatabase(t, database.ID())
			t.Cleanup(databaseRoleCleanup)

			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, sdk.NewGrantPrivilegesToDatabaseRoleRequest(databaseRole.ID()).WithPrivileges(*sdk.NewDatabaseRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{tc.privilege})).WithOn(*sdk.NewDatabaseRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInDatabase(database.ID())))))
			require.ErrorContains(t, err, tc.expectedFutureError)
		})

		t.Run("account role - all "+tc.objectTypePlural.String(), func(t *testing.T) {
			_, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)

			err := client.Grants.GrantPrivilegesToAccountRole(ctx, sdk.NewGrantPrivilegesToAccountRoleRequest(role.ID()).WithPrivileges(*sdk.NewAccountRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{tc.privilege})).WithOn(*sdk.NewAccountRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInSchema(testClientHelper().Ids.SchemaId())))))

			if tc.expectedAllError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedAllError)
			}
		})

		t.Run("database role - all "+tc.objectTypePlural.String(), func(t *testing.T) {
			_, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
			t.Cleanup(databaseRoleCleanup)

			err := client.Grants.GrantPrivilegesToDatabaseRole(ctx, sdk.NewGrantPrivilegesToDatabaseRoleRequest(databaseRole.ID()).WithPrivileges(*sdk.NewDatabaseRoleGrantPrivilegesRequest().WithSchemaObjectPrivileges([]sdk.SchemaObjectPrivilege{tc.privilege})).WithOn(*sdk.NewDatabaseRoleGrantOnRequest().WithSchemaObject(*sdk.NewGrantOnSchemaObjectRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInSchema(testClientHelper().Ids.SchemaId())))))

			if tc.expectedAllError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedAllError)
			}
		})

		t.Run("ownership - future "+tc.objectTypePlural.String(), func(t *testing.T) {
			database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)
			roleId := role.ID()

			err := client.Grants.GrantOwnership(ctx, sdk.NewGrantOwnershipRequest(*sdk.NewOwnershipGrantOnRequest().WithFuture(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInDatabase(database.ID())), *sdk.NewOwnershipGrantToRequest().WithAccountRoleName(roleId)))
			require.ErrorContains(t, err, tc.expectedFutureError)
		})

		t.Run("ownership - all "+tc.objectTypePlural.String(), func(t *testing.T) {
			_, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)
			roleId := role.ID()

			err := client.Grants.GrantOwnership(ctx, sdk.NewGrantOwnershipRequest(*sdk.NewOwnershipGrantOnRequest().WithAll(*sdk.NewGrantOnSchemaObjectInRequest(tc.objectTypePlural).WithInSchema(testClientHelper().Ids.SchemaId())), *sdk.NewOwnershipGrantToRequest().WithAccountRoleName(roleId)))

			if tc.expectedAllError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.expectedAllError)
			}
		})
	}
}
