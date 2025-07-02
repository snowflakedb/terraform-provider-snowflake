//go:build !account_level_tests

package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_OrganizationAccount_SelfAlter(t *testing.T) {
	testClientHelper().EnsureValidNonProdOrganizationAccountIsUsed(t)

	client := testClient(t)
	ctx := testContext(t)

	t.Cleanup(testClientHelper().Role.UseRoleWithClient(t, client, snowflakeroles.GlobalOrgAdmin))
	t.Cleanup(testClientHelper().Role.UseRole(t, snowflakeroles.GlobalOrgAdmin))

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
		assert.ErrorContains(t, err, fmt.Sprintf("Only one %s is allowed at a time", sdk.PolicyKindSessionPolicy))

		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithSessionPolicy(true)))
		require.NoError(t, err)
		assertThatNoPolicyIsSetOnAccount(t)

		// Unset session policy when there's no password policy attached
		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithSessionPolicy(true)))
		assert.ErrorContains(t, err, fmt.Sprintf("Any policy of kind %s is not attached to ACCOUNT", sdk.PolicyKindSessionPolicy))
	})

	t.Run("set / unset parameters",
		SetAndUnsetAccountParametersTest(
			func(parameters sdk.AccountParameters) error {
				return client.OrganizationAccounts.Alter(context.Background(), sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithParameters(parameters)))
			},
			client.OrganizationAccounts.UnsetAllParameters,
			client.OrganizationAccounts.ShowParameters,
		),
	)
}
