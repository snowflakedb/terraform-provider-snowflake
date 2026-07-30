//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_PasswordPolicies(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("create password_policy: complete case", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		comment := random.Comment()

		request := sdk.NewCreatePasswordPolicyRequest(id).
			WithOrReplace(true).
			WithPasswordMinLength(10).
			WithPasswordMaxLength(20).
			WithPasswordMinUpperCaseChars(2).
			WithPasswordMinLowerCaseChars(3).
			WithPasswordMinNumericChars(4).
			WithPasswordMinSpecialChars(1).
			WithPasswordMinAgeDays(25).
			WithPasswordMaxAgeDays(30).
			WithPasswordMaxRetries(3).
			WithPasswordLockoutTimeMins(30).
			WithPasswordHistory(15).
			WithComment(comment)

		passwordPolicy, cleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicyWithOptions(t, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.PasswordPolicyFromObject(t, passwordPolicy).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("PASSWORD_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasOwnerRoleType("ROLE").
			HasOptions(""))

		assertThatObject(t, objectassert.PasswordPolicyDetails(t, id).
			HasName(id.Name()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment(comment).
			HasPasswordMinLength(10).
			HasPasswordMaxLength(20).
			HasPasswordMinUpperCaseChars(2).
			HasPasswordMinLowerCaseChars(3).
			HasPasswordMinNumericChars(4).
			HasPasswordMinSpecialChars(1).
			HasPasswordMinAgeDays(25).
			HasPasswordMaxAgeDays(30).
			HasPasswordMaxRetries(3).
			HasPasswordLockoutTimeMins(30).
			HasPasswordHistory(15))
	})

	t.Run("create password_policy: no optionals", func(t *testing.T) {
		passwordPolicy, cleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(cleanup)
		id := passwordPolicy.ID()

		assertThatObject(t, objectassert.PasswordPolicyFromObject(t, passwordPolicy).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasKind("PASSWORD_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(""))

		assertThatObject(t, objectassert.PasswordPolicyDetails(t, id).
			HasName(passwordPolicy.Name).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasPasswordMinLength(14).
			HasPasswordMaxLength(256).
			HasPasswordMinUpperCaseChars(1).
			HasPasswordMinLowerCaseChars(1).
			HasPasswordMinNumericChars(1).
			HasPasswordMinSpecialChars(0).
			HasPasswordMinAgeDays(0).
			HasPasswordMaxAgeDays(90).
			HasPasswordMaxRetries(5).
			HasPasswordLockoutTimeMins(15).
			HasPasswordHistory(5))
	})

	t.Run("create password_policy: if not exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		request := sdk.NewCreatePasswordPolicyRequest(id).
			WithIfNotExists(true).
			WithPasswordMinLength(10).
			WithPasswordMaxLength(20).
			WithPasswordMinUpperCaseChars(5).
			WithComment("test comment")

		_, cleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicyWithOptions(t, request)
		t.Cleanup(cleanup)

		assertThatObject(t, objectassert.PasswordPolicyDetails(t, id).
			HasName(id.Name()).
			HasComment("test comment").
			HasPasswordMinLength(10).
			HasPasswordMaxLength(20).
			HasPasswordMinUpperCaseChars(5))

		// Creating again with IF NOT EXISTS should succeed without error
		err := client.PasswordPolicies.Create(ctx, sdk.NewCreatePasswordPolicyRequest(id).
			WithIfNotExists(true).
			WithPasswordMinLength(99))
		require.NoError(t, err)

		// Original values should remain unchanged
		assertThatObject(t, objectassert.PasswordPolicyDetails(t, id).
			HasPasswordMinLength(10))
	})

	t.Run("drop password_policy: existing", func(t *testing.T) {
		passwordPolicy, _ := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		id := passwordPolicy.ID()

		err := client.PasswordPolicies.Drop(ctx, sdk.NewDropPasswordPolicyRequest(id))
		require.NoError(t, err)

		_, err = client.PasswordPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop password_policy: non-existing", func(t *testing.T) {
		err := client.PasswordPolicies.Drop(ctx, sdk.NewDropPasswordPolicyRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("drop password_policy: existing with if exists", func(t *testing.T) {
		passwordPolicy, _ := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		id := passwordPolicy.ID()

		err := client.PasswordPolicies.Drop(ctx, sdk.NewDropPasswordPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)

		_, err = client.PasswordPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("alter password_policy: set value and unset value", func(t *testing.T) {
		passwordPolicy, cleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(cleanup)
		id := passwordPolicy.ID()

		alterRequest := sdk.NewAlterPasswordPolicyRequest(id).WithSet(*sdk.NewPasswordPolicySetRequest().
			WithPasswordMinLength(10).
			WithPasswordMaxLength(20).
			WithPasswordMinUpperCaseChars(2).
			WithPasswordMinLowerCaseChars(3).
			WithPasswordMinNumericChars(4).
			WithPasswordMinSpecialChars(1).
			WithPasswordMinAgeDays(1).
			WithPasswordMaxAgeDays(30).
			WithPasswordMaxRetries(10).
			WithPasswordLockoutTimeMins(30).
			WithPasswordHistory(5).
			WithComment("new comment"))
		err := client.PasswordPolicies.Alter(ctx, alterRequest)
		require.NoError(t, err)

		assertThatObject(t, objectassert.PasswordPolicyDetails(t, id).
			HasComment("new comment").
			HasPasswordMinLength(10).
			HasPasswordMaxLength(20).
			HasPasswordMinUpperCaseChars(2).
			HasPasswordMinLowerCaseChars(3).
			HasPasswordMinNumericChars(4).
			HasPasswordMinSpecialChars(1).
			HasPasswordMinAgeDays(1).
			HasPasswordMaxAgeDays(30).
			HasPasswordMaxRetries(10).
			HasPasswordLockoutTimeMins(30).
			HasPasswordHistory(5))

		err = client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(id).
			WithUnset(*sdk.NewPasswordPolicyUnsetRequest().
				WithPasswordMinLength(true).
				WithPasswordMaxLength(true).
				WithPasswordMinUpperCaseChars(true).
				WithPasswordMinLowerCaseChars(true).
				WithPasswordMinNumericChars(true).
				WithPasswordMinSpecialChars(true).
				WithPasswordMinAgeDays(true).
				WithPasswordMaxAgeDays(true).
				WithPasswordMaxRetries(true).
				WithPasswordLockoutTimeMins(true).
				WithPasswordHistory(true).
				WithComment(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.PasswordPolicyDetails(t, id).
			HasComment("").
			HasPasswordMinLength(14).
			HasPasswordMaxLength(256).
			HasPasswordMinUpperCaseChars(1).
			HasPasswordMinLowerCaseChars(1).
			HasPasswordMinNumericChars(1).
			HasPasswordMinSpecialChars(0).
			HasPasswordMinAgeDays(0).
			HasPasswordMaxAgeDays(90).
			HasPasswordMaxRetries(5).
			HasPasswordLockoutTimeMins(15).
			HasPasswordHistory(5))
	})

	t.Run("alter password_policy: rename", func(t *testing.T) {
		passwordPolicy, _ := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		id := passwordPolicy.ID()

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(id).WithRenameTo(newId))
		if err != nil {
			t.Cleanup(testClientHelper().PasswordPolicy.DropPasswordPolicyFunc(t, id))
		} else {
			t.Cleanup(testClientHelper().PasswordPolicy.DropPasswordPolicyFunc(t, newId))
		}
		require.NoError(t, err)

		_, err = client.PasswordPolicies.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		assertThatObject(t, objectassert.PasswordPolicy(t, newId).
			HasName(newId.Name()).
			HasDatabaseName(newId.DatabaseName()).
			HasSchemaName(newId.SchemaName()).
			HasKind("PASSWORD_POLICY").
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasComment("").
			HasOwnerRoleType("ROLE").
			HasOptions(""))
	})

	t.Run("show password policies", func(t *testing.T) {
		db, dbCleanup := testClientHelper().Database.CreateDatabase(t)
		t.Cleanup(dbCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_password_policy_1_")
		id2 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_password_policy_2_")
		id3 := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("test_password_policy_3_")
		id4 := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(sdk.NewDatabaseObjectIdentifier(db.Name, "PUBLIC"))
		ids := []sdk.SchemaObjectIdentifier{id1, id2, id3, id4}
		for _, id := range ids {
			_, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicyWithOptions(t, sdk.NewCreatePasswordPolicyRequest(id))
			t.Cleanup(passwordPolicyCleanup)
		}

		t.Run("like", func(t *testing.T) {
			passwordPolicies, err := client.PasswordPolicies.Show(ctx, sdk.NewShowPasswordPolicyRequest().
				WithLike(sdk.Like{Pattern: sdk.String("test_password_policy_2_%")}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, passwordPolicies, 1)
		})

		t.Run("in_account", func(t *testing.T) {
			passwordPolicies, err := client.PasswordPolicies.Show(ctx, sdk.NewShowPasswordPolicyRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Account: sdk.Bool(true)}}))
			require.NoError(t, err)
			assert.GreaterOrEqual(t, len(passwordPolicies), 4)
		})

		t.Run("in_database", func(t *testing.T) {
			passwordPolicies, err := client.PasswordPolicies.Show(ctx, sdk.NewShowPasswordPolicyRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Database: id1.DatabaseId()}}))
			require.NoError(t, err)
			assert.Len(t, passwordPolicies, 3)
		})

		t.Run("in_schema", func(t *testing.T) {
			passwordPolicies, err := client.PasswordPolicies.Show(ctx, sdk.NewShowPasswordPolicyRequest().
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, passwordPolicies, 3)
		})

		t.Run("limit", func(t *testing.T) {
			passwordPolicies, err := client.PasswordPolicies.Show(ctx, sdk.NewShowPasswordPolicyRequest().
				WithLimit(sdk.LimitFrom{Rows: sdk.Int(1)}).
				WithIn(sdk.ExtendedIn{In: sdk.In{Schema: id1.SchemaId()}}))
			require.NoError(t, err)
			assert.Len(t, passwordPolicies, 1)
		})
	})

	t.Run("describe password_policy: non-existing", func(t *testing.T) {
		_, err := client.PasswordPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		_, cleanup1 := testClientHelper().PasswordPolicy.CreatePasswordPolicyWithOptions(t, sdk.NewCreatePasswordPolicyRequest(id1))
		t.Cleanup(cleanup1)
		_, cleanup2 := testClientHelper().PasswordPolicy.CreatePasswordPolicyWithOptions(t, sdk.NewCreatePasswordPolicyRequest(id2))
		t.Cleanup(cleanup2)

		e1, err := client.PasswordPolicies.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.PasswordPolicies.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
