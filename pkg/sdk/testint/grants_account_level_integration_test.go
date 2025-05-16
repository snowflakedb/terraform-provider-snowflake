//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ShowGrants_To_Users(t *testing.T) {
	secondaryTestClientHelper().BcrBundles.EnableBcrBundle(t, "2025_02")

	t.Run("handles granteeName for account role granted to user", func(t *testing.T) {
		user, userCleanup := secondaryTestClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		accountRole, accountRoleCleanup := secondaryTestClientHelper().Role.CreateRole(t)
		t.Cleanup(accountRoleCleanup)

		secondaryTestClientHelper().Role.GrantRoleToUser(t, accountRole.ID(), user.ID())
		grants, err := secondaryTestClientHelper().Grant.ShowGrantsOfAccountRole(t, accountRole.ID())
		require.NoError(t, err)

		assert.Len(t, grants, 1)
		assert.Equal(t, sdk.ObjectTypeUser, grants[0].GrantedTo)
		assert.Equal(t, grants[0].GranteeName.FullyQualifiedName(), user.ID().FullyQualifiedName())
	})

	t.Run("handles granteeName for database role granted to user", func(t *testing.T) {
		user, userCleanup := secondaryTestClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		databaseRole, databaseRoleCleanup := secondaryTestClientHelper().DatabaseRole.CreateDatabaseRole(t)
		t.Cleanup(databaseRoleCleanup)

		secondaryTestClientHelper().Grant.GrantDatabaseRoleToUser(t, databaseRole.ID(), user.ID())
		grants, err := secondaryTestClientHelper().Grant.ShowGrantsOfDatabaseRole(t, databaseRole.ID())
		require.NoError(t, err)

		assert.Len(t, grants, 1)
		assert.Equal(t, sdk.ObjectTypeUser, grants[0].GrantedTo)
		assert.Equal(t, grants[0].GranteeName.FullyQualifiedName(), user.ID().FullyQualifiedName())
	})
}
