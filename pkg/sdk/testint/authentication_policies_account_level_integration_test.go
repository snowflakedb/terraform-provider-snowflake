package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_AuthenticationPolicies_AccountLevel(t *testing.T) {
	client := testSecondaryClient(t)
	ctx := testSecondaryContext(t)
	secondaryTestClientHelper().BcrBundles.EnableBcrBundle(t, "2025_06")

	t.Run("Create - with options deprecated in 2025_06", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		err := client.AuthenticationPolicies.Create(ctx, sdk.NewCreateAuthenticationPolicyRequest(id).
			WithMfaAuthenticationMethods([]sdk.MfaAuthenticationMethods{
				{Method: sdk.MfaAuthenticationMethodsPassword},
				{Method: sdk.MfaAuthenticationMethodsSaml},
			}))
		require.ErrorContains(t, err, "001420 (22023): SQL compilation error:\ninvalid property 'MFA_AUTHENTICATION_METHODS' for 'AUTHENTICATION_POLICY'")
	})

	t.Run("Alter - set and unset options deprecated in 2025_06", func(t *testing.T) {
		authenticationPolicy, cleanupAuthPolicy := secondaryTestClientHelper().AuthenticationPolicy.Create(t)
		t.Cleanup(cleanupAuthPolicy)

		err := client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(authenticationPolicy.ID()).
			WithSet(*sdk.NewAuthenticationPolicySetRequest().
				WithMfaAuthenticationMethods([]sdk.MfaAuthenticationMethods{
					{Method: sdk.MfaAuthenticationMethodsPassword},
					{Method: sdk.MfaAuthenticationMethodsSaml},
				})))
		require.ErrorContains(t, err, "003639 (01P01): SQL Compilation Error: MFA_AUTHENTICATION_METHODS is deprecated, please use MFA_POLICY=(ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=ALL | NONE) instead.")

		err = client.AuthenticationPolicies.Alter(ctx, sdk.NewAlterAuthenticationPolicyRequest(authenticationPolicy.ID()).
			WithUnset(*sdk.NewAuthenticationPolicyUnsetRequest().
				WithMfaAuthenticationMethods(true)))
		require.ErrorContains(t, err, "003639 (01P01): SQL Compilation Error: MFA_AUTHENTICATION_METHODS is deprecated, please use MFA_POLICY=(ENFORCE_MFA_ON_EXTERNAL_AUTHENTICATION=ALL | NONE) instead.")
	})
}
