//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_Users_BCR_2025_05(t *testing.T) {
	ctx := testContext(t)
	client := testSecondaryClient(t)
	secondaryTestClientHelper().BcrBundles.EnableBcrBundle(t, "2025_05")
	t.Cleanup(func() {
		secondaryTestClientHelper().BcrBundles.DisableBcrBundle(t, "2025_05")
	})

	t.Run("alter: unset type", func(t *testing.T) {
		user, userCleanup := secondaryTestClientHelper().User.CreateServiceUser(t)
		t.Cleanup(userCleanup)

		assertThatObjectSecondary(t, objectassert.UserFromObject(t, user).
			HasType(string(sdk.UserTypeService)),
		)

		alterOpts := &sdk.AlterUserOptions{Unset: &sdk.UserUnset{
			ObjectProperties: &sdk.UserObjectPropertiesUnset{
				Type: sdk.Bool(true),
			},
		}}

		err := client.Users.Alter(ctx, user.ID(), alterOpts)
		require.NoError(t, err)

		assertThatObjectSecondary(t, objectassert.User(t, user.ID()).
			HasType(string(sdk.UserTypePerson)),
		)
	})
}
