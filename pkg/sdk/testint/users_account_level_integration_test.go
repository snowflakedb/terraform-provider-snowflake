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

	t.Run("alter: unset type with 2025_05 bundle", func(t *testing.T) {
		secondaryTestClientHelper().BcrBundles.EnableBcrBundle(t, "2025_05")

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

	t.Run("alter: unset type without 2025_05 bundle", func(t *testing.T) {
		secondaryTestClientHelper().BcrBundles.DisableBcrBundle(t, "2025_05")

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
			HasType(""),
		)
	})
}
