//go:build non_account_level_tests

package testint

import (
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_PasswordPoliciesShow(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	passwordPolicyTest, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(passwordPolicyCleanup)

	passwordPolicy2Test, passwordPolicy2Cleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(passwordPolicy2Cleanup)

	t.Run("without show options", func(t *testing.T) {
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, sdk.NewShowPasswordPolicyRequest())
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(passwordPolicies))
	})

	t.Run("with show options", func(t *testing.T) {
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, sdk.NewShowPasswordPolicyRequest().
			WithIn(sdk.In{Schema: testClientHelper().Ids.SchemaId()}))
		require.NoError(t, err)
		assert.Contains(t, passwordPolicies, *passwordPolicyTest)
		assert.Contains(t, passwordPolicies, *passwordPolicy2Test)
		assert.Len(t, passwordPolicies, 2)
	})

	t.Run("with show options and like", func(t *testing.T) {
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, sdk.NewShowPasswordPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String(passwordPolicyTest.Name)}).
			WithIn(sdk.In{Database: testClientHelper().Ids.DatabaseId()}))
		require.NoError(t, err)
		assert.Contains(t, passwordPolicies, *passwordPolicyTest)
		assert.Len(t, passwordPolicies, 1)
	})

	t.Run("when searching a non-existent password policy", func(t *testing.T) {
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, sdk.NewShowPasswordPolicyRequest().
			WithLike(sdk.Like{Pattern: sdk.String("non-existent")}))
		require.NoError(t, err)
		assert.Empty(t, passwordPolicies)
	})

	/* there appears to be a bug in the Snowflake API. LIMIT is not actually limiting the number of results
	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &ShowPasswordPolicyOptions{
			In: &In{
				Schema: String(schemaTest.FullyQualifiedName()),
			},
			Limit: &LimitFrom{Rows: Int(1)},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Equal(t, 1, len(passwordPolicies))
	})*/
}

func TestInt_PasswordPolicyCreate(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("test complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.PasswordPolicies.Create(ctx, sdk.NewCreatePasswordPolicyRequest(id).
			WithOrReplace(true).
			WithPasswordMinLength(10).
			WithPasswordMaxLength(20).
			WithPasswordMinUpperCaseChars(1).
			WithPasswordMinLowerCaseChars(1).
			WithPasswordMinNumericChars(1).
			WithPasswordMinSpecialChars(1).
			WithPasswordMinAgeDays(25).
			WithPasswordMaxAgeDays(30).
			WithPasswordMaxRetries(5).
			WithPasswordLockoutTimeMins(30).
			WithPasswordHistory(15).
			WithComment("test comment"))
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.DescribeDetails(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), passwordPolicyDetails.Name)
		assert.Equal(t, 10, passwordPolicyDetails.PasswordMinLength)
		assert.Equal(t, 20, passwordPolicyDetails.PasswordMaxLength)
		assert.Equal(t, 1, passwordPolicyDetails.PasswordMinUpperCaseChars)
		assert.Equal(t, 1, passwordPolicyDetails.PasswordMinLowerCaseChars)
		assert.Equal(t, 1, passwordPolicyDetails.PasswordMinNumericChars)
		assert.Equal(t, 1, passwordPolicyDetails.PasswordMinSpecialChars)
		assert.Equal(t, 25, passwordPolicyDetails.PasswordMinAgeDays)
		assert.Equal(t, 30, passwordPolicyDetails.PasswordMaxAgeDays)
		assert.Equal(t, 5, passwordPolicyDetails.PasswordMaxRetries)
		assert.Equal(t, 30, passwordPolicyDetails.PasswordLockoutTimeMins)
		assert.Equal(t, 15, passwordPolicyDetails.PasswordHistory)
		assert.Equal(t, "test comment", passwordPolicyDetails.Comment)
	})

	t.Run("test if_not_exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.PasswordPolicies.Create(ctx, sdk.NewCreatePasswordPolicyRequest(id).
			WithIfNotExists(true).
			WithPasswordMinLength(10).
			WithPasswordMaxLength(20).
			WithPasswordMinUpperCaseChars(5).
			WithComment("test comment"))
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.DescribeDetails(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), passwordPolicyDetails.Name)
		assert.Equal(t, "test comment", passwordPolicyDetails.Comment)
		assert.Equal(t, 10, passwordPolicyDetails.PasswordMinLength)
		assert.Equal(t, 20, passwordPolicyDetails.PasswordMaxLength)
		assert.Equal(t, 5, passwordPolicyDetails.PasswordMinUpperCaseChars)
	})

	t.Run("test no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.PasswordPolicies.Create(ctx, sdk.NewCreatePasswordPolicyRequest(id))
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.DescribeDetails(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), passwordPolicyDetails.Name)
		assert.Equal(t, "", passwordPolicyDetails.Comment)
		// All parameters use Snowflake defaults when not specified
		assert.Greater(t, passwordPolicyDetails.PasswordMinLength, 0)
		assert.Greater(t, passwordPolicyDetails.PasswordMaxLength, 0)
	})
}

func TestInt_PasswordPolicyDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(passwordPolicyCleanup)

	t.Run("when password policy exists", func(t *testing.T) {
		passwordPolicyDetails, err := client.PasswordPolicies.DescribeDetails(ctx, passwordPolicy.ID())
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name)
		assert.Equal(t, passwordPolicy.Comment, passwordPolicyDetails.Comment)
	})

	t.Run("when password policy does not exist", func(t *testing.T) {
		_, err := client.PasswordPolicies.Describe(ctx, NonExistingSchemaObjectIdentifier)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_PasswordPolicyAlter(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when setting new values", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)
		err := client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(passwordPolicy.ID()).
			WithSet(*sdk.NewPasswordPolicySetRequest().
				WithPasswordMinLength(10).
				WithPasswordMaxLength(20).
				WithComment("new comment")))
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.DescribeDetails(ctx, passwordPolicy.ID())
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name)
		assert.Equal(t, 10, passwordPolicyDetails.PasswordMinLength)
		assert.Equal(t, 20, passwordPolicyDetails.PasswordMaxLength)
		assert.Equal(t, "new comment", passwordPolicyDetails.Comment)
	})

	t.Run("when renaming", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		oldID := passwordPolicy.ID()
		t.Cleanup(passwordPolicyCleanup)
		newID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(oldID).WithNewName(newID))
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.DescribeDetails(ctx, newID)
		require.NoError(t, err)
		// rename back to original name, so it can be cleaned up
		assert.Equal(t, newID.Name(), passwordPolicyDetails.Name)
		err = client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(newID).WithNewName(oldID))
		require.NoError(t, err)
	})

	t.Run("when unsetting values", func(t *testing.T) {
		createReq := sdk.NewCreatePasswordPolicyRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()).
			WithPasswordMaxAgeDays(20).
			WithPasswordMaxRetries(10).
			WithComment("test comment")
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicyWithOptions(t, createReq)
		id := passwordPolicy.ID()
		t.Cleanup(passwordPolicyCleanup)
		err := client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(id).
			WithUnset(*sdk.NewPasswordPolicyUnsetRequest().WithPasswordMaxRetries(true)))
		require.NoError(t, err)
		err = client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(id).
			WithUnset(*sdk.NewPasswordPolicyUnsetRequest().
				WithPasswordMaxAgeDays(true).
				WithComment(true)))
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.DescribeDetails(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name)
		assert.Equal(t, "", passwordPolicyDetails.Comment)
		// After unsetting PASSWORD_MAX_RETRIES, it should revert to Snowflake's default
		assert.Greater(t, passwordPolicyDetails.PasswordMaxRetries, 0)
	})

	t.Run("when unsetting multiple values at same time", func(t *testing.T) {
		createReq := sdk.NewCreatePasswordPolicyRequest(testClientHelper().Ids.RandomSchemaObjectIdentifier()).
			WithPasswordMaxAgeDays(20).
			WithPasswordMaxRetries(10).
			WithComment("test comment")
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicyWithOptions(t, createReq)
		id := passwordPolicy.ID()
		t.Cleanup(passwordPolicyCleanup)
		err := client.PasswordPolicies.Alter(ctx, sdk.NewAlterPasswordPolicyRequest(id).
			WithUnset(*sdk.NewPasswordPolicyUnsetRequest().
				WithPasswordMaxAgeDays(true).
				WithPasswordMaxRetries(true).
				WithComment(true)))
		require.NoError(t, err)
	})
}

func TestInt_PasswordPolicyDrop(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	t.Run("when password policy exists", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)
		id := passwordPolicy.ID()
		err := client.PasswordPolicies.Drop(ctx, sdk.NewDropPasswordPolicyRequest(id))
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when password policy does not exist", func(t *testing.T) {
		err := client.PasswordPolicies.Drop(ctx, sdk.NewDropPasswordPolicyRequest(NonExistingSchemaObjectIdentifier))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when password policy exists and if exists is true", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)
		id := passwordPolicy.ID()
		err := client.PasswordPolicies.Drop(ctx, sdk.NewDropPasswordPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_PasswordPoliciesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	cleanupPasswordPolicyHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.PasswordPolicies.Drop(ctx, sdk.NewDropPasswordPolicyRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createPasswordPolicyHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		err := client.PasswordPolicies.Create(ctx, sdk.NewCreatePasswordPolicyRequest(id))
		require.NoError(t, err)
		t.Cleanup(cleanupPasswordPolicyHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		id2 := testClientHelper().Ids.NewSchemaObjectIdentifierInSchema(id1.Name(), schema.ID())

		createPasswordPolicyHandle(t, id1)
		createPasswordPolicyHandle(t, id2)

		e1, err := client.PasswordPolicies.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.PasswordPolicies.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})

	t.Run("show by id: check fields", func(t *testing.T) {
		id1 := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		createPasswordPolicyHandle(t, id1)

		sl, err := client.PasswordPolicies.ShowByID(ctx, id1)
		require.NoError(t, err)
		assert.Equal(t, "ROLE", sl.OwnerRoleType)
	})
}
