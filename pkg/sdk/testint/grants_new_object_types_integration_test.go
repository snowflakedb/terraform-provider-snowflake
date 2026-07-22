//go:build non_account_level_tests

package testint

import (
	"regexp"
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
	}

	for _, tc := range testCases {
		// --- FUTURE GRANTS ---

		t.Run("account role - future "+tc.objectTypePlural.String(), func(t *testing.T) {
			database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)

			privileges := &sdk.AccountRoleGrantPrivileges{
				SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{tc.privilege},
			}
			on := &sdk.AccountRoleGrantOn{
				SchemaObject: &sdk.GrantOnSchemaObject{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectTypePlural,
						InDatabase:       sdk.Pointer(database.ID()),
					},
				},
			}
			err := client.Grants.GrantPrivilegesToAccountRole(ctx, privileges, on, role.ID(), nil)
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				Future: sdk.Bool(true),
				To:     &sdk.ShowGrantsTo{Role: role.ID()},
			})
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

			err := client.Grants.GrantPrivilegesToDatabaseRole(
				ctx,
				&sdk.DatabaseRoleGrantPrivileges{
					SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{tc.privilege},
				},
				&sdk.DatabaseRoleGrantOn{
					SchemaObject: &sdk.GrantOnSchemaObject{
						Future: &sdk.GrantOnSchemaObjectIn{
							PluralObjectType: tc.objectTypePlural,
							InDatabase:       sdk.Pointer(database.ID()),
						},
					},
				},
				databaseRole.ID(), nil,
			)
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				Future: sdk.Bool(true),
				To:     &sdk.ShowGrantsTo{DatabaseRole: databaseRole.ID()},
			})
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

			err := client.Grants.GrantOwnership(
				ctx,
				sdk.OwnershipGrantOn{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectTypePlural,
						InDatabase:       sdk.Pointer(database.ID()),
					},
				},
				sdk.OwnershipGrantTo{AccountRoleName: &roleId},
				nil,
			)
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				Future: sdk.Bool(true),
				To:     &sdk.ShowGrantsTo{Role: roleId},
			})
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

			err := client.Grants.GrantPrivilegesToAccountRole(
				ctx,
				&sdk.AccountRoleGrantPrivileges{
					SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{tc.privilege},
				},
				&sdk.AccountRoleGrantOn{
					SchemaObject: &sdk.GrantOnSchemaObject{
						All: &sdk.GrantOnSchemaObjectIn{
							PluralObjectType: tc.objectTypePlural,
							InSchema:         sdk.Pointer(testClientHelper().Ids.SchemaId()),
						},
					},
				},
				role.ID(), nil,
			)
			require.NoError(t, err)

			grants, err := client.Grants.Show(ctx, &sdk.ShowGrantOptions{
				On: &sdk.ShowGrantsOn{Object: &sdk.Object{
					ObjectType: tc.expectedGrantOn,
					Name:       objectId,
				}},
			})
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
			_, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			databaseRole, databaseRoleCleanup := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
			t.Cleanup(databaseRoleCleanup)

			err := client.Grants.GrantPrivilegesToDatabaseRole(
				ctx,
				&sdk.DatabaseRoleGrantPrivileges{
					SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{tc.privilege},
				},
				&sdk.DatabaseRoleGrantOn{
					SchemaObject: &sdk.GrantOnSchemaObject{
						All: &sdk.GrantOnSchemaObjectIn{
							PluralObjectType: tc.objectTypePlural,
							InSchema:         sdk.Pointer(testClientHelper().Ids.SchemaId()),
						},
					},
				},
				databaseRole.ID(), nil,
			)
			require.NoError(t, err)
		})

		t.Run("ownership - all "+tc.objectTypePlural.String(), func(t *testing.T) {
			_, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)
			roleId := role.ID()

			err := client.Grants.GrantOwnership(
				ctx,
				sdk.OwnershipGrantOn{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectTypePlural,
						InSchema:         sdk.Pointer(testClientHelper().Ids.SchemaId()),
					},
				},
				sdk.OwnershipGrantTo{AccountRoleName: &roleId},
				nil,
			)
			require.NoError(t, err)
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
		createObject        func(t *testing.T) (sdk.SchemaObjectIdentifier, func())
		expectedFutureError *regexp.Regexp
		expectedAllError    *regexp.Regexp
	}

	testCases := []testCase{
		{
			objectTypePlural: sdk.PluralObjectTypeExperiments,
			privilege:        sdk.SchemaObjectPrivilegeUsage,
			createObject: func(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
				t.Helper()
				return testClientHelper().Experiment.Create(t)
			},
			expectedFutureError: regexp.MustCompile(`Unsupported feature 'EXPERIMENT'`),
			expectedAllError:    regexp.MustCompile(`Unsupported feature 'GRANT on all objects of type EXPERIMENT'`),
		},
		{
			objectTypePlural: sdk.PluralObjectTypeGateways,
			privilege:        sdk.SchemaObjectPrivilegeUsage,
			createObject: func(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
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
			expectedFutureError: regexp.MustCompile(`syntax error line 0 at position 0 unexpected 'TOK_GATEWAY'`),
			expectedAllError:    regexp.MustCompile(`syntax error line 0 at position 0 unexpected 'TOK_GATEWAY'`),
		},
	}

	for _, tc := range testCases {
		t.Run("account role - future "+tc.objectTypePlural.String(), func(t *testing.T) {
			database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)

			err := client.Grants.GrantPrivilegesToAccountRole(
				ctx,
				&sdk.AccountRoleGrantPrivileges{
					SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{tc.privilege},
				},
				&sdk.AccountRoleGrantOn{
					SchemaObject: &sdk.GrantOnSchemaObject{
						Future: &sdk.GrantOnSchemaObjectIn{
							PluralObjectType: tc.objectTypePlural,
							InDatabase:       sdk.Pointer(database.ID()),
						},
					},
				},
				role.ID(), nil,
			)
			require.Error(t, err)
			assert.Regexp(t, tc.expectedFutureError, err.Error())
		})

		t.Run("account role - all "+tc.objectTypePlural.String(), func(t *testing.T) {
			_, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)

			err := client.Grants.GrantPrivilegesToAccountRole(
				ctx,
				&sdk.AccountRoleGrantPrivileges{
					SchemaObjectPrivileges: []sdk.SchemaObjectPrivilege{tc.privilege},
				},
				&sdk.AccountRoleGrantOn{
					SchemaObject: &sdk.GrantOnSchemaObject{
						All: &sdk.GrantOnSchemaObjectIn{
							PluralObjectType: tc.objectTypePlural,
							InSchema:         sdk.Pointer(testClientHelper().Ids.SchemaId()),
						},
					},
				},
				role.ID(), nil,
			)
			require.Error(t, err)
			assert.Regexp(t, tc.expectedAllError, err.Error())
		})

		t.Run("ownership - future "+tc.objectTypePlural.String(), func(t *testing.T) {
			database, databaseCleanup := testClientHelper().Database.CreateDatabase(t)
			t.Cleanup(databaseCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)
			roleId := role.ID()

			err := client.Grants.GrantOwnership(
				ctx,
				sdk.OwnershipGrantOn{
					Future: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectTypePlural,
						InDatabase:       sdk.Pointer(database.ID()),
					},
				},
				sdk.OwnershipGrantTo{AccountRoleName: &roleId},
				nil,
			)
			require.Error(t, err)
			assert.Regexp(t, tc.expectedFutureError, err.Error())
		})

		t.Run("ownership - all "+tc.objectTypePlural.String(), func(t *testing.T) {
			_, objectCleanup := tc.createObject(t)
			t.Cleanup(objectCleanup)

			role, roleCleanup := testClientHelper().Role.CreateRole(t)
			t.Cleanup(roleCleanup)
			roleId := role.ID()

			err := client.Grants.GrantOwnership(
				ctx,
				sdk.OwnershipGrantOn{
					All: &sdk.GrantOnSchemaObjectIn{
						PluralObjectType: tc.objectTypePlural,
						InSchema:         sdk.Pointer(testClientHelper().Ids.SchemaId()),
					},
				},
				sdk.OwnershipGrantTo{AccountRoleName: &roleId},
				nil,
			)
			require.Error(t, err)
			assert.Regexp(t, tc.expectedAllError, err.Error())
		})
	}
}
