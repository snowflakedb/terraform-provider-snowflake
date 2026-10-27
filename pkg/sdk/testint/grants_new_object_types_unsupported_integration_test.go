package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// unsupportedFutureAndAllGrantTestCase is shared by the non-account and account-level
// builds. Snowflake Intelligence is an account-level object, so its case lives in the
// account_level_tests file; the other object types stay in the non-account file.
type unsupportedFutureAndAllGrantTestCase struct {
	objectTypePlural    sdk.PluralObjectType
	privilege           sdk.SchemaObjectPrivilege
	createObject        func(t *testing.T) (sdk.ObjectIdentifier, func())
	expectedFutureError string
	expectedAllError    string
}

// testGrantPrivilegesOnFutureAndAllUnsupportedObjectTypes verifies that granting on future/all
// for specific object types results in an error from Snowflake (these object types are not yet
// supported for bulk grants despite being listed in the docs).
func testGrantPrivilegesOnFutureAndAllUnsupportedObjectTypes(t *testing.T, testCases []unsupportedFutureAndAllGrantTestCase) {
	t.Helper()
	client := testClient(t)
	ctx := testContext(t)

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
