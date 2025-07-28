//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_DatabaseRoles(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	assertDatabaseRole := func(t *testing.T, databaseRole *sdk.DatabaseRole, expectedId sdk.DatabaseObjectIdentifier, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, databaseRole.CreatedOn)
		assert.Equal(t, expectedId.Name(), databaseRole.Name)
		assert.Equal(t, "ACCOUNTADMIN", databaseRole.Owner)
		assert.Equal(t, expectedComment, databaseRole.Comment)
		assert.Equal(t, 0, databaseRole.GrantedToRoles)
		assert.Equal(t, 0, databaseRole.GrantedToDatabaseRoles)
		assert.Equal(t, 0, databaseRole.GrantedDatabaseRoles)
	}

	cleanupDatabaseRoleProvider := func(id sdk.DatabaseObjectIdentifier) func() {
		return func() {
			err := client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id))
			require.NoError(t, err)
		}
	}

	createDatabaseRole := func(t *testing.T) *sdk.DatabaseRole {
		t.Helper()
		id := testClientHelper().Ids.RandomDatabaseObjectIdentifier()

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		return databaseRole
	}

	t.Run("create database_role: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		comment := random.Comment()

		request := sdk.NewCreateDatabaseRoleRequest(id).WithComment(comment).WithIfNotExists(true)
		err := client.DatabaseRoles.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)

		require.NoError(t, err)
		assertDatabaseRole(t, databaseRole, id, comment)
	})

	t.Run("create database_role: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomDatabaseObjectIdentifier()

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assertDatabaseRole(t, databaseRole, id, "")
	})

	t.Run("drop database_role: existing", func(t *testing.T) {
		id := testClientHelper().Ids.RandomDatabaseObjectIdentifier()

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)

		err = client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id))
		require.NoError(t, err)

		_, err = client.DatabaseRoles.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop database_role: non-existing", func(t *testing.T) {
		id := NonExistingDatabaseObjectIdentifier

		err := client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("alter database_role: set value and unset value", func(t *testing.T) {
		id := testClientHelper().Ids.RandomDatabaseObjectIdentifier()

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		alterRequest := sdk.NewAlterDatabaseRoleRequest(id).WithSet(*sdk.NewDatabaseRoleSetRequest("new comment"))
		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredDatabaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredDatabaseRole.Comment)

		alterRequest = sdk.NewAlterDatabaseRoleRequest(id).WithUnset(*sdk.NewDatabaseRoleUnsetRequest())
		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredDatabaseRole, err = client.DatabaseRoles.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredDatabaseRole.Comment)
	})

	t.Run("alter database_role: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomDatabaseObjectIdentifier()

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)

		newId := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
		alterRequest := sdk.NewAlterDatabaseRoleRequest(id).WithRename(newId)

		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupDatabaseRoleProvider(id))
		} else {
			t.Cleanup(cleanupDatabaseRoleProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.DatabaseRoles.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		databaseRole, err := client.DatabaseRoles.ShowByID(ctx, newId)
		require.NoError(t, err)

		assertDatabaseRole(t, databaseRole, newId, "")
	})

	t.Run("alter database_role: rename to other database", func(t *testing.T) {
		secondDatabase, secondDatabaseCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(secondDatabaseCleanup)

		id := testClientHelper().Ids.RandomDatabaseObjectIdentifier()

		err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupDatabaseRoleProvider(id))

		newId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(secondDatabase.ID())
		alterRequest := sdk.NewAlterDatabaseRoleRequest(id).WithRename(newId)

		err = client.DatabaseRoles.Alter(ctx, alterRequest)
		assert.ErrorIs(t, err, sdk.ErrDifferentDatabase)
	})

	t.Run("show database_role: without like", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		role2 := createDatabaseRole(t)

		showRequest := sdk.NewShowDatabaseRoleRequest(testClientHelper().Ids.DatabaseId())
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Len(t, returnedDatabaseRoles, 2)
		assert.Contains(t, returnedDatabaseRoles, *role1)
		assert.Contains(t, returnedDatabaseRoles, *role2)
	})

	t.Run("show database_role: with like", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		role2 := createDatabaseRole(t)

		showRequest := sdk.NewShowDatabaseRoleRequest(testClientHelper().Ids.DatabaseId()).WithLike(sdk.Like{Pattern: &role1.Name})
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Len(t, returnedDatabaseRoles, 1)
		assert.Contains(t, returnedDatabaseRoles, *role1)
		assert.NotContains(t, returnedDatabaseRoles, *role2)
	})

	t.Run("show database_role: with like and limit", func(t *testing.T) {
		prefix := "SHOW_TEST_ROLE_"
		roleId1 := testClientHelper().Ids.AlphaWithPrefix(prefix + "1")
		roleId2 := testClientHelper().Ids.AlphaWithPrefix(prefix + "2")

		role1, cleanupRole1 := testClientHelper().DatabaseRole.CreateDatabaseRoleWithName(t, roleId1)
		t.Cleanup(cleanupRole1)

		role2, cleanupRole2 := testClientHelper().DatabaseRole.CreateDatabaseRoleWithName(t, roleId2)
		t.Cleanup(cleanupRole2)

		showRequest := sdk.NewShowDatabaseRoleRequest(testClientHelper().Ids.DatabaseId()).
			WithLike(sdk.Like{
				Pattern: sdk.Pointer(prefix + "%"),
			}).
			WithLimit(sdk.LimitFrom{
				Rows: sdk.Pointer(1),
				From: sdk.Pointer(roleId1),
			})
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Len(t, returnedDatabaseRoles, 1)
		assert.NotContains(t, returnedDatabaseRoles, *role1)
		assert.Contains(t, returnedDatabaseRoles, *role2)
	})

	t.Run("show database_role: no matches", func(t *testing.T) {
		showRequest := sdk.NewShowDatabaseRoleRequest(testClientHelper().Ids.DatabaseId()).WithLike(sdk.Like{Pattern: sdk.Pointer("non-existent")})
		returnedDatabaseRoles, err := client.DatabaseRoles.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Empty(t, returnedDatabaseRoles)
	})

	t.Run("grant and revoke database_role: to database role", func(t *testing.T) {
		role1 := createDatabaseRole(t)
		id1 := testClientHelper().Ids.NewDatabaseObjectIdentifier(role1.Name)
		role2 := createDatabaseRole(t)
		id2 := testClientHelper().Ids.NewDatabaseObjectIdentifier(role2.Name)

		grantRequest := sdk.NewGrantDatabaseRoleRequest(id1).WithDatabaseRole(id2)
		err := client.DatabaseRoles.Grant(ctx, grantRequest)
		require.NoError(t, err)

		extractedRole, err := client.DatabaseRoles.ShowByID(ctx, id1)
		require.NoError(t, err)
		assert.Equal(t, 0, extractedRole.GrantedToRoles)
		assert.Equal(t, 1, extractedRole.GrantedToDatabaseRoles)
		assert.Equal(t, 0, extractedRole.GrantedDatabaseRoles)

		extractedRole, err = client.DatabaseRoles.ShowByID(ctx, id2)
		require.NoError(t, err)
		assert.Equal(t, 0, extractedRole.GrantedToRoles)
		assert.Equal(t, 0, extractedRole.GrantedToDatabaseRoles)
		assert.Equal(t, 1, extractedRole.GrantedDatabaseRoles)

		revokeRequest := sdk.NewRevokeDatabaseRoleRequest(id1).WithDatabaseRole(id2)
		err = client.DatabaseRoles.Revoke(ctx, revokeRequest)
		require.NoError(t, err)
	})

	t.Run("grant and revoke database_role: to account role", func(t *testing.T) {
		role := createDatabaseRole(t)
		roleId := testClientHelper().Ids.NewDatabaseObjectIdentifier(role.Name)

		accountRole, accountRoleCleanup := testClientHelper().Role.CreateRole(t)
		t.Cleanup(accountRoleCleanup)

		grantRequest := sdk.NewGrantDatabaseRoleRequest(roleId).WithAccountRole(accountRole.ID())
		err := client.DatabaseRoles.Grant(ctx, grantRequest)
		require.NoError(t, err)

		extractedRole, err := client.DatabaseRoles.ShowByID(ctx, roleId)
		require.NoError(t, err)
		assert.Equal(t, 1, extractedRole.GrantedToRoles)
		assert.Equal(t, 0, extractedRole.GrantedToDatabaseRoles)
		assert.Equal(t, 0, extractedRole.GrantedDatabaseRoles)

		revokeRequest := sdk.NewRevokeDatabaseRoleRequest(roleId).WithAccountRole(accountRole.ID())
		err = client.DatabaseRoles.Revoke(ctx, revokeRequest)
		require.NoError(t, err)
	})

	t.Run("grant and revoke database_role: to share", func(t *testing.T) {
		role := createDatabaseRole(t)
		roleId := testClientHelper().Ids.NewDatabaseObjectIdentifier(role.Name)

		share, shareCleanup := testClientHelper().Share.CreateShare(t)
		t.Cleanup(shareCleanup)

		err := client.Grants.GrantPrivilegeToShare(ctx, []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}, &sdk.ShareGrantOn{Database: testClientHelper().Ids.DatabaseId()}, share.ID())
		require.NoError(t, err)

		grantRequest := sdk.NewGrantDatabaseRoleToShareRequest(roleId, share.ID())
		err = client.DatabaseRoles.GrantToShare(ctx, grantRequest)
		require.NoError(t, err)

		revokeRequest := sdk.NewRevokeDatabaseRoleFromShareRequest(roleId, share.ID())
		err = client.DatabaseRoles.RevokeFromShare(ctx, revokeRequest)
		require.NoError(t, err)
	})
}
