//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_PolicyReferences(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	passwordPolicyId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
	err := client.PasswordPolicies.Create(ctx, passwordPolicyId, &sdk.CreatePasswordPolicyOptions{})
	require.NoError(t, err)

	t.Cleanup(func() {
		err := client.PasswordPolicies.Drop(ctx, passwordPolicyId, &sdk.DropPasswordPolicyOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	})

	t.Run("user domain", func(t *testing.T) {
		user, userCleanup := testClientHelper().User.CreateUser(t)
		t.Cleanup(userCleanup)

		err = client.Users.Alter(ctx, user.ID(), &sdk.AlterUserOptions{
			Set: &sdk.UserSet{
				PasswordPolicy: &passwordPolicyId,
			},
		})
		require.NoError(t, err)

		policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(user.ID(), sdk.PolicyEntityDomainUser))
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(policyReferences), 1)

		passwordPolicy, err := collections.FindFirst(policyReferences, func(reference sdk.PolicyReference) bool {
			return reference.PolicyKind == sdk.PolicyKindPasswordPolicy && reference.PolicyName == passwordPolicyId.Name()
		})
		require.NoError(t, err)
		require.Equal(t, passwordPolicyId.Name(), passwordPolicy.PolicyName)
		require.Equal(t, sdk.PolicyKindPasswordPolicy, passwordPolicy.PolicyKind)
	})

	t.Run("tag domain", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := testClientHelper().MaskingPolicy.CreateMaskingPolicy(t)
		t.Cleanup(maskingPolicyCleanup)

		tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
		t.Cleanup(tagCleanup)

		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(tag.ID()).WithSet(
			sdk.NewTagSetRequest().WithMaskingPolicies([]sdk.SchemaObjectIdentifier{maskingPolicy.ID()}),
		))
		require.NoError(t, err)

		policyReferences, err := client.PolicyReferences.GetForEntity(ctx, sdk.NewGetForEntityPolicyReferenceRequest(tag.ID(), sdk.PolicyEntityDomainTag))
		require.NoError(t, err)
		require.Len(t, policyReferences, 1)
		require.Equal(t, maskingPolicy.ID().Name(), policyReferences[0].PolicyName)
		require.Equal(t, sdk.PolicyKindMaskingPolicy, policyReferences[0].PolicyKind)

		err = client.Tags.Alter(ctx, sdk.NewAlterTagRequest(tag.ID()).WithUnset(
			sdk.NewTagUnsetRequest().WithMaskingPolicies([]sdk.SchemaObjectIdentifier{maskingPolicy.ID()}),
		))
		require.NoError(t, err)
	})
}
