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
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, nil)
		require.NoError(t, err)
		assert.LessOrEqual(t, 2, len(passwordPolicies))
	})

	t.Run("with show options", func(t *testing.T) {
		showOptions := &sdk.ShowPasswordPolicyOptions{
			In: &sdk.In{
				Schema: testClientHelper().Ids.SchemaId(),
			},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, passwordPolicies, *passwordPolicyTest)
		assert.Contains(t, passwordPolicies, *passwordPolicy2Test)
		assert.Len(t, passwordPolicies, 2)
	})

	t.Run("with show options and like", func(t *testing.T) {
		showOptions := &sdk.ShowPasswordPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(passwordPolicyTest.Name),
			},
			In: &sdk.In{
				Database: testClientHelper().Ids.DatabaseId(),
			},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Contains(t, passwordPolicies, *passwordPolicyTest)
		assert.Len(t, passwordPolicies, 1)
	})

	t.Run("when searching a non-existent password policy", func(t *testing.T) {
		showOptions := &sdk.ShowPasswordPolicyOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		passwordPolicies, err := client.PasswordPolicies.Show(ctx, showOptions)
		require.NoError(t, err)
		assert.Empty(t, passwordPolicies)
	})

	/* there appears to be a bug in the Snowflake API. LIMIT is not actually limiting the number of results
	t.Run("when limiting the number of results", func(t *testing.T) {
		showOptions := &ShowPasswordPolicyOptions{
			In: &In{
				Schema: String(schemaTest.FullyQualifiedName()),
			},
			Limit: Int(1),
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
		err := client.PasswordPolicies.Create(ctx, id, &sdk.CreatePasswordPolicyOptions{
			OrReplace:                 sdk.Bool(true),
			PasswordMinLength:         sdk.Int(10),
			PasswordMaxLength:         sdk.Int(20),
			PasswordMinUpperCaseChars: sdk.Int(1),
			PasswordMinLowerCaseChars: sdk.Int(1),
			PasswordMinNumericChars:   sdk.Int(1),
			PasswordMinSpecialChars:   sdk.Int(1),
			PasswordMinAgeDays:        sdk.Int(25),
			PasswordMaxAgeDays:        sdk.Int(30),
			PasswordMaxRetries:        sdk.Int(5),
			PasswordLockoutTimeMins:   sdk.Int(30),
			PasswordHistory:           sdk.Int(15),
			Comment:                   sdk.String("test comment"),
		})
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), passwordPolicyDetails.Name.Value)
		assert.Equal(t, 10, *passwordPolicyDetails.PasswordMinLength.Value)
		assert.Equal(t, 20, *passwordPolicyDetails.PasswordMaxLength.Value)
		assert.Equal(t, 1, *passwordPolicyDetails.PasswordMinUpperCaseChars.Value)
		assert.Equal(t, 1, *passwordPolicyDetails.PasswordMinLowerCaseChars.Value)
		assert.Equal(t, 1, *passwordPolicyDetails.PasswordMinNumericChars.Value)
		assert.Equal(t, 1, *passwordPolicyDetails.PasswordMinSpecialChars.Value)
		assert.Equal(t, 25, *passwordPolicyDetails.PasswordMinAgeDays.Value)
		assert.Equal(t, 30, *passwordPolicyDetails.PasswordMaxAgeDays.Value)
		assert.Equal(t, 5, *passwordPolicyDetails.PasswordMaxRetries.Value)
		assert.Equal(t, 30, *passwordPolicyDetails.PasswordLockoutTimeMins.Value)
		assert.Equal(t, 15, *passwordPolicyDetails.PasswordHistory.Value)
		assert.Equal(t, "test comment", passwordPolicyDetails.Comment.Value)
	})

	t.Run("test if_not_exists", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.PasswordPolicies.Create(ctx, id, &sdk.CreatePasswordPolicyOptions{
			OrReplace:                 sdk.Bool(false),
			IfNotExists:               sdk.Bool(true),
			PasswordMinLength:         sdk.Int(10),
			PasswordMaxLength:         sdk.Int(20),
			PasswordMinUpperCaseChars: sdk.Int(5),
			Comment:                   sdk.String("test comment"),
		})
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), passwordPolicyDetails.Name.Value)
		assert.Equal(t, "test comment", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, 10, *passwordPolicyDetails.PasswordMinLength.Value)
		assert.Equal(t, 20, *passwordPolicyDetails.PasswordMaxLength.Value)
		assert.Equal(t, 5, *passwordPolicyDetails.PasswordMinUpperCaseChars.Value)
	})

	t.Run("test no options", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		err := client.PasswordPolicies.Create(ctx, id, nil)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), passwordPolicyDetails.Name.Value)
		assert.Equal(t, "", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinLength.Value, *passwordPolicyDetails.PasswordMinLength.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMaxLength.Value, *passwordPolicyDetails.PasswordMaxLength.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinUpperCaseChars.Value, *passwordPolicyDetails.PasswordMinUpperCaseChars.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinLowerCaseChars.Value, *passwordPolicyDetails.PasswordMinLowerCaseChars.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinNumericChars.Value, *passwordPolicyDetails.PasswordMinNumericChars.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMinSpecialChars.Value, *passwordPolicyDetails.PasswordMinSpecialChars.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMaxAgeDays.Value, *passwordPolicyDetails.PasswordMaxAgeDays.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordMaxRetries.Value, *passwordPolicyDetails.PasswordMaxRetries.DefaultValue)
		assert.Equal(t, *passwordPolicyDetails.PasswordLockoutTimeMins.Value, *passwordPolicyDetails.PasswordLockoutTimeMins.DefaultValue)
	})
}

func TestInt_PasswordPolicyDescribe(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
	t.Cleanup(passwordPolicyCleanup)

	t.Run("when password policy exists", func(t *testing.T) {
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, passwordPolicy.ID())
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, passwordPolicy.Comment, passwordPolicyDetails.Comment.Value)
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
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Set: &sdk.PasswordPolicySet{
				PasswordMinLength: sdk.Int(10),
				PasswordMaxLength: sdk.Int(20),
				Comment:           sdk.String("new comment"),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, passwordPolicy.ID(), alterOptions)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, passwordPolicy.ID())
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, 10, *passwordPolicyDetails.PasswordMinLength.Value)
		assert.Equal(t, 20, *passwordPolicyDetails.PasswordMaxLength.Value)
		assert.Equal(t, "new comment", passwordPolicyDetails.Comment.Value)
	})

	t.Run("when renaming", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		oldID := passwordPolicy.ID()
		t.Cleanup(passwordPolicyCleanup)
		newID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			NewName: &newID,
		}
		err := client.PasswordPolicies.Alter(ctx, oldID, alterOptions)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, newID)
		require.NoError(t, err)
		// rename back to original name, so it can be cleaned up
		assert.Equal(t, newID.Name(), passwordPolicyDetails.Name.Value)
		alterOptions = &sdk.AlterPasswordPolicyOptions{
			NewName: &oldID,
		}
		err = client.PasswordPolicies.Alter(ctx, newID, alterOptions)
		require.NoError(t, err)
	})

	t.Run("when unsetting values", func(t *testing.T) {
		createOptions := &sdk.CreatePasswordPolicyOptions{
			PasswordMaxAgeDays: sdk.Int(20),
			PasswordMaxRetries: sdk.Int(10),
			Comment:            sdk.String("test comment"),
		}
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicyWithOptions(t, createOptions)
		id := passwordPolicy.ID()
		t.Cleanup(passwordPolicyCleanup)
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Unset: &sdk.PasswordPolicyUnset{
				PasswordMaxRetries: sdk.Bool(true),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		alterOptions = &sdk.AlterPasswordPolicyOptions{
			Unset: &sdk.PasswordPolicyUnset{
				PasswordMaxAgeDays: sdk.Bool(true),
				Comment:            sdk.Bool(true),
			},
		}
		err = client.PasswordPolicies.Alter(ctx, id, alterOptions)
		require.NoError(t, err)
		passwordPolicyDetails, err := client.PasswordPolicies.Describe(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, passwordPolicy.Name, passwordPolicyDetails.Name.Value)
		assert.Equal(t, "", passwordPolicyDetails.Comment.Value)
		assert.Equal(t, *passwordPolicyDetails.PasswordMaxRetries.Value, *passwordPolicyDetails.PasswordMaxRetries.DefaultValue)
	})

	t.Run("when unsetting multiple values at same time", func(t *testing.T) {
		createOptions := &sdk.CreatePasswordPolicyOptions{
			PasswordMaxAgeDays: sdk.Int(20),
			PasswordMaxRetries: sdk.Int(10),
			Comment:            sdk.String("test comment"),
		}
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicyWithOptions(t, createOptions)
		id := passwordPolicy.ID()
		t.Cleanup(passwordPolicyCleanup)
		alterOptions := &sdk.AlterPasswordPolicyOptions{
			Unset: &sdk.PasswordPolicyUnset{
				PasswordMaxAgeDays: sdk.Bool(true),
				PasswordMaxRetries: sdk.Bool(true),
				Comment:            sdk.Bool(true),
			},
		}
		err := client.PasswordPolicies.Alter(ctx, id, alterOptions)
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
		err := client.PasswordPolicies.Drop(ctx, id, nil)
		require.NoError(t, err)
		_, err = client.PasswordPolicies.Describe(ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when password policy does not exist", func(t *testing.T) {
		err := client.PasswordPolicies.Drop(ctx, NonExistingSchemaObjectIdentifier, nil)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("when password policy exists and if exists is true", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)
		id := passwordPolicy.ID()
		dropOptions := &sdk.DropPasswordPolicyOptions{IfExists: sdk.Bool(true)}
		err := client.PasswordPolicies.Drop(ctx, id, dropOptions)
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
			err := client.PasswordPolicies.Drop(ctx, id, nil)
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createPasswordPolicyHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		err := client.PasswordPolicies.Create(ctx, id, nil)
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
