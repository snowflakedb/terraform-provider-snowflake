//go:build account_level_tests

package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_OrganizationAccount_SelfAlter(t *testing.T) {
	testClientHelper().EnsureValidNonProdOrganizationAccountIsUsed(t)

	client := testClient(t)
	ctx := testContext(t)

	t.Cleanup(testClientHelper().Role.UseRole(t, snowflakeroles.GlobalOrgAdmin))

	require.NoError(t, client.Sessions.UseRole(ctx, snowflakeroles.GlobalOrgAdmin))
	t.Cleanup(func() { require.NoError(t, client.Sessions.UseRole(ctx, snowflakeroles.Accountadmin)) })

	t.Run("set / unset comment", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier(testClientHelper().OrganizationAccount.ShowCurrent(t).AccountName)
		comment := random.Comment()

		assertThatObject(t, objectassert.OrganizationAccount(t, id).HasNoComment())

		err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithComment(comment)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.OrganizationAccount(t, id).HasComment(comment))

		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithComment(true)))
		require.NoError(t, err)

		assertThatObject(t, objectassert.OrganizationAccount(t, id).HasNoComment())
	})

	t.Run("set / unset resource monitor", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		resourceMonitor2, resourceMonitor2Cleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitor2Cleanup)

		err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithResourceMonitor(resourceMonitor.ID())))
		require.NoError(t, err)

		// Set another resource monitor without unsetting the previous one
		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithResourceMonitor(resourceMonitor2.ID())))
		require.NoError(t, err)

		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithResourceMonitor(true)))
		require.NoError(t, err)

		// TODO(SNOW-2184799): Currently, there's no way to query resource monitor to verify to was unset properly.
	})

	t.Run("set / unset password policy", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)

		passwordPolicy2, passwordPolicy2Cleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicy2Cleanup)

		err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithPasswordPolicy(passwordPolicy.ID())))
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, sdk.PolicyKindPasswordPolicy, passwordPolicy.ID())

		// Set another password policy without unsetting the previous one
		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithPasswordPolicy(passwordPolicy2.ID())))
		assert.ErrorContains(t, err, fmt.Sprintf("Only one %s is allowed at a time", sdk.PolicyKindPasswordPolicy))

		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithPasswordPolicy(true)))
		require.NoError(t, err)
		assertThatNoPolicyIsSetOnAccount(t)

		// Unset password policy when there's no password policy attached
		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithPasswordPolicy(true)))
		require.NoError(t, err)
	})

	t.Run("set / unset session policy", func(t *testing.T) {
		sessionPolicy, sessionPolicyCleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicyCleanup)

		sessionPolicy2, sessionPolicy2Cleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicy2Cleanup)

		err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithSessionPolicy(sessionPolicy.ID())))
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, sdk.PolicyKindSessionPolicy, sessionPolicy.ID())

		// Set another session policy without unsetting the previous one
		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithSessionPolicy(sessionPolicy2.ID())))
		assert.ErrorContains(t, err, "already has a SESSION POLICY on ACCOUNT. The existing policy must be removed")

		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithSessionPolicy(true)))
		require.NoError(t, err)
		assertThatNoPolicyIsSetOnAccount(t)

		// Unset session policy when there's no password policy attached
		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithSessionPolicy(true)))
		require.NoError(t, err)
	})

	t.Run("set / unset policies safely", func(t *testing.T) {
		sessionPolicy, sessionPolicyCleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicyCleanup)

		sessionPolicy2, sessionPolicy2Cleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicy2Cleanup)

		err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithSessionPolicy(sessionPolicy.ID())))
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, sdk.PolicyKindSessionPolicy, sessionPolicy.ID())

		// Try to set policy when there's already one set
		err = client.OrganizationAccounts.SetPolicySafely(ctx, sdk.PolicyKindSessionPolicy, sessionPolicy2.ID())
		require.NoError(t, err)

		// Try to set unsupported policy kind
		err = client.OrganizationAccounts.SetPolicySafely(ctx, sdk.PolicyKindAuthenticationPolicy, sessionPolicy2.ID())
		assert.ErrorContains(t, err, fmt.Sprintf("policy kind %s is not supported for organization account policies", sdk.PolicyKindAuthenticationPolicy))

		// Try to unset unsupported policy kind
		err = client.OrganizationAccounts.UnsetPolicySafely(ctx, sdk.PolicyKindAuthenticationPolicy)
		assert.ErrorContains(t, err, fmt.Sprintf("policy kind %s is not supported for organization account policies", sdk.PolicyKindAuthenticationPolicy))

		err = client.OrganizationAccounts.UnsetPolicySafely(ctx, sdk.PolicyKindSessionPolicy)
		require.NoError(t, err)
		assertThatNoPolicyIsSetOnAccount(t)

		// Try to unset policy when there's no policy set
		err = client.OrganizationAccounts.UnsetPolicySafely(ctx, sdk.PolicyKindSessionPolicy)
		require.NoError(t, err)
	})

	t.Run("set / unset parameters",
		setAndUnsetAccountParametersTest(
			func(ctx context.Context, parameters sdk.AccountParameters) error {
				return client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithParameters(parameters)))
			},
			client.OrganizationAccounts.UnsetAllParameters,
			client.OrganizationAccounts.ShowParameters,
		),
	)
}
