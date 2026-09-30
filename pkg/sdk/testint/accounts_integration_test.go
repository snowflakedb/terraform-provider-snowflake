//go:build account_level_tests

package testint

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-1920887): Some of the account features cannot be currently tested as they require two Snowflake organizations
// TODO(SNOW-1342761): Adjust the tests, so they can be run in their own pipeline
// For now, those tests should be run manually. The account/admin user running those tests is required to:
// - Be privileged with ORGADMIN and ACCOUNTADMIN roles.
// - Shouldn't be any of the "main" accounts/admin users, because those tests alter the current account.

func TestInt_AccountsShowParametersDetails(t *testing.T) {
	testClientHelper().EnsureValidNonProdAccountIsUsed(t)

	client := testClient(t)
	ctx := testContext(t)
	id := testClientHelper().Context.CurrentAccountId(t)

	details, err := client.Accounts.ShowParametersDetails(ctx)
	require.NoError(t, err)

	assertThatObject(t, objectparametersassert.AccountParameters(t, id).
		HasAutocommit(details.Autocommit.Value).
		HasJsonIndent(details.JsonIndent.Value).
		HasLogLevel(details.LogLevel.Value).
		HasTimezone(details.Timezone.Value))

	assert.Equal(t, string(sdk.AccountParameterAutocommit), details.Autocommit.Key)
	assert.Equal(t, sdk.ParameterTypeAccount, details.Autocommit.Level)
	assert.NotEmpty(t, details.Autocommit.Description)

	// Organization accounts expose the same parameters, so both accessors have to agree.
	organizationAccountDetails, err := client.OrganizationAccounts.ShowParametersDetails(ctx)
	require.NoError(t, err)
	assert.Equal(t, details.Autocommit, organizationAccountDetails.Autocommit)
	assert.Equal(t, details.JsonIndent, organizationAccountDetails.JsonIndent)
	assert.Equal(t, details.LogLevel, organizationAccountDetails.LogLevel)
	assert.Equal(t, details.Timezone, organizationAccountDetails.Timezone)
}

func TestInt_Account(t *testing.T) {
	testClientHelper().EnsureValidNonProdAccountIsUsed(t)

	client := testClient(t)
	ctx := testContext(t)
	currentAccountId := testClientHelper().Context.CurrentAccountId(t)
	currentAccountName := currentAccountId.AccountName()
	defaultConsumptionBillingEntity := testClientHelper().Context.DefaultConsumptionBillingEntity(t).Name()

	assertAccountQueriedByOrgAdmin := func(t *testing.T, account sdk.Account, accountName string) {
		t.Helper()
		assert.NotEmpty(t, account.OrganizationName)
		assert.Equal(t, accountName, account.AccountName)
		assert.Nil(t, account.RegionGroup)
		assert.NotEmpty(t, account.SnowflakeRegion)
		// TODO [SNOW-3797718]: changed to business critical temporarily
		assert.Equal(t, sdk.AccountEditionBusinessCritical, *account.Edition)
		assert.NotEmpty(t, *account.AccountURL)
		assert.NotEmpty(t, *account.CreatedOn)
		// TODO [SNOW-3797718]: changed to nil temporarily - having direct dereference panics here; use object assertions
		assert.Nil(t, account.Comment)
		assert.NotEmpty(t, account.AccountLocator)
		assert.NotEmpty(t, *account.AccountLocatorUrl)
		assert.Zero(t, *account.ManagedAccounts)
		assert.NotEmpty(t, *account.ConsumptionBillingEntityName)
		assert.Nil(t, account.MarketplaceConsumerBillingEntityName)
		// TODO [SNOW-3797718]: changed to nil temporarily
		assert.Nil(t, account.MarketplaceProviderBillingEntityName)
		assert.Empty(t, *account.OldAccountURL)
		assert.True(t, *account.IsOrgAdmin)
		assert.Nil(t, account.AccountOldUrlSavedOn)
		assert.Nil(t, account.AccountOldUrlLastUsed)
		assert.Empty(t, *account.OrganizationOldUrl)
		assert.Nil(t, account.OrganizationOldUrlSavedOn)
		assert.Nil(t, account.OrganizationOldUrlLastUsed)
		assert.False(t, *account.IsEventsAccount)
		assert.False(t, account.IsOrganizationAccount)
	}

	assertAccountQueriedByAccountAdmin := func(t *testing.T, account sdk.Account, accountName string) {
		t.Helper()
		assert.NotEmpty(t, account.OrganizationName)
		assert.Equal(t, accountName, account.AccountName)
		assert.NotEmpty(t, account.SnowflakeRegion)
		assert.NotEmpty(t, account.AccountLocator)
		assert.False(t, account.IsOrganizationAccount)
		assert.Nil(t, account.RegionGroup)
		assert.Nil(t, account.Edition)
		assert.Nil(t, account.AccountURL)
		assert.Nil(t, account.CreatedOn)
		assert.Nil(t, account.Comment)
		assert.Nil(t, account.AccountLocatorUrl)
		assert.Nil(t, account.ManagedAccounts)
		assert.Nil(t, account.ConsumptionBillingEntityName)
		assert.Nil(t, account.MarketplaceConsumerBillingEntityName)
		assert.Nil(t, account.MarketplaceProviderBillingEntityName)
		assert.Nil(t, account.OldAccountURL)
		assert.Nil(t, account.IsOrgAdmin)
		assert.Nil(t, account.IsOrgAdmin)
		assert.Nil(t, account.AccountOldUrlSavedOn)
		assert.Nil(t, account.AccountOldUrlLastUsed)
		assert.Nil(t, account.OrganizationOldUrl)
		assert.Nil(t, account.OrganizationOldUrlSavedOn)
		assert.Nil(t, account.OrganizationOldUrlLastUsed)
		assert.Nil(t, account.IsEventsAccount)
	}

	assertHistoryAccount := func(t *testing.T, account sdk.Account, accountName string) {
		t.Helper()
		assertAccountQueriedByOrgAdmin(t, account, currentAccountName)
		assert.Nil(t, account.DroppedOn)
		assert.Nil(t, account.ScheduledDeletionTime)
		assert.Nil(t, account.RestoredOn)
		assert.Empty(t, account.MovedToOrganization)
		assert.Nil(t, account.MovedOn)
		assert.Nil(t, account.OrganizationUrlExpirationOn)
	}

	t.Run("create: minimal", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier(random.AccountName())
		name := random.AdminName()
		password := random.Password()
		email := random.Email()

		err := client.Accounts.Create(ctx, sdk.NewCreateAccountRequest(id, name, email).
			WithAdminPassword(password).
			WithEdition(sdk.AccountEditionStandard))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
	})

	t.Run("create: user type service", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier(random.AccountName())
		name := random.AdminName()
		key, _ := random.GenerateRSAPublicKey(t)
		email := random.Email()

		err := client.Accounts.Create(ctx, sdk.NewCreateAccountRequest(id, name, email).
			WithAdminRsaPublicKey(key).
			WithAdminUserType(sdk.UserTypeService).
			WithEdition(sdk.AccountEditionStandard))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
	})

	t.Run("create: user type legacy service", func(t *testing.T) {
		t.Skip("SNOW-3797718 - can we drop this test already?")
		id := sdk.NewAccountObjectIdentifier(random.AccountName())
		name := random.AdminName()
		password := random.Password()
		email := random.Email()

		err := client.Accounts.Create(ctx, sdk.NewCreateAccountRequest(id, name, email).
			WithAdminPassword(password).
			WithAdminUserType(sdk.UserTypeLegacyService).
			WithEdition(sdk.AccountEditionStandard))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
	})

	// TODO [SNOW-3797718] - billing entity was removed; use contract number
	t.Run("create: complete", func(t *testing.T) {
		id := sdk.NewAccountObjectIdentifier(random.AccountName())
		name := random.AdminName()
		password := random.Password()
		email := random.Email()
		region := testClientHelper().Context.CurrentRegion(t)
		regions := testClientHelper().Account.ShowRegions(t)
		currentRegion, err := collections.FindFirst(regions, func(r helpers.Region) bool {
			return r.SnowflakeRegion == region
		})
		require.NoError(t, err)
		comment := random.Comment()

		err = client.Accounts.Create(ctx, sdk.NewCreateAccountRequest(id, name, email).
			WithAdminPassword(password).
			WithFirstName("firstName").
			WithLastName("lastName").
			WithMustChangePassword(true).
			WithEdition(sdk.AccountEditionStandard).
			WithRegionGroup("PUBLIC").
			WithRegion(currentRegion.SnowflakeRegion).
			WithComment(comment))
		// TODO(SNOW-1895880): with polaris Snowflake returns an error saying: "invalid property polaris for account"
		// .WithPolaris(true)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
	})

	t.Run("alter: set / unset is org admin", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		require.False(t, *account.IsOrgAdmin)

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithName(account.ID()).
			WithSet(*sdk.NewAccountSetRequest().WithOrgAdmin(true)))
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.True(t, *acc.IsOrgAdmin)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithName(account.ID()).
			WithSet(*sdk.NewAccountSetRequest().WithOrgAdmin(false)))
		require.NoError(t, err)

		acc, err = client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.False(t, *acc.IsOrgAdmin)
	})

	t.Run("alter: rename", func(t *testing.T) {
		oldAccount, oldAccountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(oldAccountCleanup)

		newName := sdk.NewAccountObjectIdentifier(random.AccountName())
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithName(oldAccount.ID()).
			WithRenameTo(newName))
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, oldAccount.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		newAccount, err := client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.NotNil(t, newAccount)
		require.NotEqual(t, oldAccount.AccountURL, newAccount.AccountURL)
		require.Equal(t, oldAccount.AccountURL, newAccount.OldAccountURL)
	})

	t.Run("alter: rename with new url", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		newName := sdk.NewAccountObjectIdentifier(random.AccountName())
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithName(account.ID()).
			WithRenameTo(newName).WithSaveOldURL(false))
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		acc, err := client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.NotEqual(t, account.AccountURL, acc.AccountURL)
		require.Empty(t, acc.OldAccountURL)
	})

	t.Run("alter: set / unset consumption billing entity", func(t *testing.T) {
		t.Skip("SNOW-3797718 - address with billing entity removal")
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		require.Equal(t, defaultConsumptionBillingEntity, *account.ConsumptionBillingEntityName)

		// We are not able to create consumption billing entities, because of that, we use the default one.
		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithName(account.ID()).
			WithSet(*sdk.NewAccountSetRequest().WithConsumptionBillingEntity(defaultConsumptionBillingEntity)))
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, defaultConsumptionBillingEntity, *acc.ConsumptionBillingEntityName)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithName(account.ID()).
			WithUnset(*sdk.NewAccountUnsetRequest().WithConsumptionBillingEntity(true)))
		require.NoError(t, err)

		acc, err = client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, defaultConsumptionBillingEntity, *acc.ConsumptionBillingEntityName)
	})

	t.Run("alter: drop url when there's no old url", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithName(account.ID()).
			WithDrop(*sdk.NewAccountDropRequest().WithOldUrl(true)))
		require.ErrorContains(t, err, "The account has no old url")
	})

	t.Run("alter: drop url after rename", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		newName := sdk.NewAccountObjectIdentifier(random.AccountName())
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithName(account.ID()).
			WithRenameTo(newName))
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.NotEmpty(t, acc.OldAccountURL)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithName(newName).
			WithDrop(*sdk.NewAccountDropRequest().WithOldUrl(true)))
		require.NoError(t, err)

		acc, err = client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.Empty(t, acc.OldAccountURL)
	})

	t.Run("drop: without options", func(t *testing.T) {
		err := client.Accounts.Drop(ctx, sdk.NewDropAccountRequest(NonExistingAccountObjectIdentifier).WithGracePeriodInDays(3))
		require.Error(t, err)

		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		err = client.Accounts.Drop(ctx, sdk.NewDropAccountRequest(account.ID()).WithGracePeriodInDays(3))
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: with if exists", func(t *testing.T) {
		err := client.Accounts.Drop(ctx, sdk.NewDropAccountRequest(NonExistingAccountObjectIdentifier).WithGracePeriodInDays(3).WithIfExists(true))
		require.NoError(t, err)

		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		err = client.Accounts.Drop(ctx, sdk.NewDropAccountRequest(account.ID()).WithGracePeriodInDays(3).WithIfExists(true))
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("undrop", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		require.NoError(t, testClientHelper().Account.Drop(t, account.ID()))

		err := client.Accounts.Undrop(ctx, sdk.NewUndropAccountRequest(account.ID()))
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, account.ID(), acc.ID())
	})

	t.Run("show: with like", func(t *testing.T) {
		currentAccount := testClientHelper().Context.CurrentAccount(t)
		accounts, err := client.Accounts.Show(ctx, sdk.NewShowAccountRequest().
			WithLike(sdk.Like{Pattern: new(currentAccount)}))
		require.NoError(t, err)
		assert.Len(t, accounts, 1)
		assertAccountQueriedByOrgAdmin(t, accounts[0], currentAccountName)
	})

	t.Run("show: with history", func(t *testing.T) {
		currentAccount := testClientHelper().Context.CurrentAccount(t)
		accounts, err := client.Accounts.Show(ctx, sdk.NewShowAccountRequest().
			WithHistory(true).
			WithLike(sdk.Like{Pattern: new(currentAccount)}))
		require.NoError(t, err)
		assert.Len(t, accounts, 1)
		assertHistoryAccount(t, accounts[0], currentAccountName)
	})

	t.Run("show: with accountadmin role", func(t *testing.T) {
		err := client.Roles.Use(ctx, sdk.NewUseRoleRequest(snowflakeroles.Accountadmin))
		require.NoError(t, err)
		t.Cleanup(func() {
			err = client.Roles.Use(ctx, sdk.NewUseRoleRequest(snowflakeroles.Orgadmin))
			require.NoError(t, err)
		})

		currentAccount := testClientHelper().Context.CurrentAccount(t)
		accounts, err := client.Accounts.Show(ctx, sdk.NewShowAccountRequest().
			WithLike(sdk.Like{Pattern: new(currentAccount)}))
		require.NoError(t, err)
		assert.Len(t, accounts, 1)
		assertAccountQueriedByAccountAdmin(t, accounts[0], currentAccountName)
	})
}

func TestInt_Account_SelfAlter(t *testing.T) {
	// TODO(SNOW-1920881): Adjust the test so that self alters will be done on newly created account - not the main test one
	testClientHelper().EnsureValidNonProdAccountIsUsed(t)

	// This client should be operating on a different account than the "main" one (because it will be altered here).
	// Cannot use a newly created account because ORGADMIN role is necessary,
	// and it is propagated only after some time (e.g., 1 hour) making it hard to automate.
	client := testClient(t)
	ctx := testContext(t)
	t.Cleanup(testClientHelper().Role.UseRole(t, snowflakeroles.Accountadmin))

	id := testClientHelper().Context.CurrentAccountId(t)

	t.Run("set / unset legacy parameters", func(t *testing.T) {
		objectparametersassert.AccountParameters(t, id).
			HasDefaultMinDataRetentionTimeInDaysValue().
			HasDefaultJsonIndentValue().
			HasDefaultUserTaskTimeoutMsValue().
			HasDefaultEnableUnredactedQuerySyntaxErrorValue()

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().WithLegacyParameters(
				*sdk.NewAccountLevelParametersRequest().
					WithAccountParameters(sdk.LegacyAccountParameters{
						MinDataRetentionTimeInDays: new(15), // default is 0
					}).
					WithSessionParameters(sdk.SessionParameters{
						JsonIndent: new(8), // default is 2
					}).
					WithObjectParameters(sdk.ObjectParameters{
						UserTaskTimeoutMs: new(100), // default is 3600000
					}).
					WithUserParameters(sdk.UserParameters{
						EnableUnredactedQuerySyntaxError: new(true), // default is false
					}),
			)))
		require.NoError(t, err)

		objectparametersassert.AccountParameters(t, id).
			HasMinDataRetentionTimeInDays(15).
			HasJsonIndent(8).
			HasUserTaskTimeoutMs(100).
			HasEnableUnredactedQuerySyntaxError(true)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithUnset(*sdk.NewAccountUnsetRequest().WithLegacyParameters(
				*sdk.NewAccountLevelParametersUnsetRequest().
					WithAccountParameters(sdk.LegacyAccountParametersUnset{
						MinDataRetentionTimeInDays: new(true),
					}).
					WithSessionParameters(sdk.SessionParametersUnset{
						JsonIndent: new(true),
					}).
					WithObjectParameters(sdk.ObjectParametersUnset{
						UserTaskTimeoutMs: new(true),
					}).
					WithUserParameters(sdk.UserParametersUnset{
						EnableUnredactedQuerySyntaxError: new(true),
					}),
			)))
		require.NoError(t, err)

		objectparametersassert.AccountParameters(t, id).
			HasDefaultMinDataRetentionTimeInDaysValue().
			HasDefaultJsonIndentValue().
			HasDefaultUserTaskTimeoutMsValue().
			HasDefaultEnableUnredactedQuerySyntaxErrorValue()
	})

	t.Run(
		"set / unset parameters",
		setAndUnsetAccountParametersTest(
			func(ctx context.Context, parameters sdk.AccountParametersRequest) error {
				return client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
					WithSet(*sdk.NewAccountSetRequest().WithParameters(parameters)))
			},
			client.Accounts.UnsetAllParameters,
			client.Accounts.ShowParameters,
		),
	)

	t.Run("set / unset resource monitor", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		require.Nil(t, resourceMonitor.Level)
		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().WithResourceMonitor(resourceMonitor.ID())))
		require.NoError(t, err)

		resourceMonitor, err = testClientHelper().ResourceMonitor.Show(t, resourceMonitor.ID())
		require.NoError(t, err)
		require.NotNil(t, resourceMonitor.Level)
		require.Equal(t, sdk.ResourceMonitorLevelAccount, *resourceMonitor.Level)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithUnset(*sdk.NewAccountUnsetRequest().WithResourceMonitor(true)))
		require.NoError(t, err)

		resourceMonitor, err = testClientHelper().ResourceMonitor.Show(t, resourceMonitor.ID())
		require.NoError(t, err)
		require.Nil(t, resourceMonitor.Level)
	})

	t.Run("set / unset policies", func(t *testing.T) {
		authPolicyForAccountAndPersonUsers, authPolicyForAccountAndPersonUsersCleanup := testClientHelper().AuthenticationPolicy.Create(t)
		t.Cleanup(authPolicyForAccountAndPersonUsersCleanup)

		authPolicyForServiceUsers, authPolicyForServiceUsersCleanup := testClientHelper().AuthenticationPolicy.Create(t)
		t.Cleanup(authPolicyForServiceUsersCleanup)

		sessionPolicyForPersonAndServiceUsers, sessionPolicyForPersonAndServiceUsersCleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicyForPersonAndServiceUsersCleanup)

		sessionPolicyForAccount, sessionPolicyForAccountCleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicyForAccountCleanup)

		featurePolicyId, featurePolicyCleanup := testClientHelper().FeaturePolicy.Create(t)
		t.Cleanup(featurePolicyCleanup)

		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)

		packagesPolicyId, packagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(packagesPolicyCleanup)

		t.Cleanup(func() {
			assert.NoError(t, client.Accounts.UnsetAllPoliciesSafely(ctx))
			assertThatNoPolicyIsSetOnAccount(t)
		})

		// Authentication policies: one policy attached account-wide and FOR ALL PERSON USERS, another attached
		// FOR ALL SERVICE USERS.
		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().
				WithAuthenticationPolicySet(*sdk.NewAccountAuthenticationPolicySetRequest().
					WithAuthenticationPolicy(authPolicyForAccountAndPersonUsers.ID()))))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().
				WithAuthenticationPolicySet(*sdk.NewAccountAuthenticationPolicySetRequest().
					WithAuthenticationPolicy(authPolicyForAccountAndPersonUsers.ID()).
					WithForAllPersonUsers(true))))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().
				WithAuthenticationPolicySet(*sdk.NewAccountAuthenticationPolicySetRequest().
					WithAuthenticationPolicy(authPolicyForServiceUsers.ID()).
					WithForAllServiceUsers(true))))
		require.NoError(t, err)

		// Session policies: one policy attached FOR ALL PERSON USERS and FOR ALL SERVICE USERS, another attached
		// account-wide.
		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().
				WithSessionPolicySet(*sdk.NewAccountSessionPolicySetRequest().
					WithSessionPolicy(sessionPolicyForPersonAndServiceUsers.ID()).
					WithForAllPersonUsers(true))))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().
				WithSessionPolicySet(*sdk.NewAccountSessionPolicySetRequest().
					WithSessionPolicy(sessionPolicyForPersonAndServiceUsers.ID()).
					WithForAllServiceUsers(true))))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().
				WithSessionPolicySet(*sdk.NewAccountSessionPolicySetRequest().
					WithSessionPolicy(sessionPolicyForAccount.ID()))))
		require.NoError(t, err)

		// Remaining policy kinds.
		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithFeaturePolicySet(*sdk.NewAccountFeaturePolicySetRequest().WithFeaturePolicy(featurePolicyId))))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithPackagesPolicy(packagesPolicyId)))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithPasswordPolicy(passwordPolicy.ID())))
		require.NoError(t, err)

		assertThatPolicyIsSetOnAccount(
			t,
			expectedAccountPolicy{sdk.PolicyKindAuthenticationPolicy, authPolicyForAccountAndPersonUsers.ID()},
			expectedAccountPolicy{sdk.PolicyKindAuthenticationPolicy, authPolicyForServiceUsers.ID()},
			expectedAccountPolicy{sdk.PolicyKindSessionPolicy, sessionPolicyForPersonAndServiceUsers.ID()},
			expectedAccountPolicy{sdk.PolicyKindSessionPolicy, sessionPolicyForAccount.ID()},
			expectedAccountPolicy{sdk.PolicyKindFeaturePolicy, featurePolicyId},
			expectedAccountPolicy{sdk.PolicyKindPasswordPolicy, passwordPolicy.ID()},
			expectedAccountPolicy{sdk.PolicyKindPackagesPolicy, packagesPolicyId},
		)

		assertThatAuthenticationPolicyIsSetOnAccountWithScopes(t, authPolicyForAccountAndPersonUsers.ID(), sdk.AuthenticationPolicyTargetScopeAccount, sdk.AuthenticationPolicyTargetScopePersonUsers)
		assertThatAuthenticationPolicyIsSetOnAccountWithScopes(t, authPolicyForServiceUsers.ID(), sdk.AuthenticationPolicyTargetScopeServiceUsers)
		assertThatSessionPolicyIsSetOnAccountWithScopes(t, sessionPolicyForPersonAndServiceUsers.ID(), sdk.SessionPolicyTargetScopePersonUsers, sdk.SessionPolicyTargetScopeServiceUsers)
		assertThatSessionPolicyIsSetOnAccountWithScopes(t, sessionPolicyForAccount.ID(), sdk.SessionPolicyTargetScopeAccount)
	})

	t.Run("force new packages policy", func(t *testing.T) {
		packagesPolicyId, packagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(packagesPolicyCleanup)

		newPackagesPolicyId, newPackagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(newPackagesPolicyCleanup)

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithPackagesPolicy(packagesPolicyId)))
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, expectedAccountPolicy{sdk.PolicyKindPackagesPolicy, packagesPolicyId})
		t.Cleanup(func() {
			err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithUnset(*sdk.NewAccountUnsetRequest().WithPackagesPolicy(true)))
			require.NoError(t, err)
			assertThatNoPolicyIsSetOnAccount(t)
		})

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithPackagesPolicy(newPackagesPolicyId)))
		require.Error(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithPackagesPolicy(newPackagesPolicyId).WithForce(true)))
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, expectedAccountPolicy{sdk.PolicyKindPackagesPolicy, newPackagesPolicyId})
	})

	t.Run("force new feature policy", func(t *testing.T) {
		featurePolicyId, featurePolicyCleanup := testClientHelper().FeaturePolicy.Create(t)
		t.Cleanup(featurePolicyCleanup)

		newFeaturePolicyId, newFeaturePolicyCleanup := testClientHelper().FeaturePolicy.Create(t)
		t.Cleanup(newFeaturePolicyCleanup)

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithFeaturePolicySet(*sdk.NewAccountFeaturePolicySetRequest().WithFeaturePolicy(featurePolicyId))))
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, expectedAccountPolicy{sdk.PolicyKindFeaturePolicy, featurePolicyId})
		t.Cleanup(func() {
			err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithUnset(*sdk.NewAccountUnsetRequest().WithFeaturePolicyUnset(*sdk.NewAccountFeaturePolicyUnsetRequest().WithFeaturePolicy(true))))
			require.NoError(t, err)
			assertThatNoPolicyIsSetOnAccount(t)
		})

		// Here we expect to get an error as there is another feature policy set on the account.
		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithFeaturePolicySet(*sdk.NewAccountFeaturePolicySetRequest().WithFeaturePolicy(newFeaturePolicyId))))
		require.Error(t, err)

		// To set a new feature policy on the account without firstly unsetting it, we can use FORCE parameter.
		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithFeaturePolicySet(*sdk.NewAccountFeaturePolicySetRequest().WithFeaturePolicy(newFeaturePolicyId)).WithForce(true)))
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, expectedAccountPolicy{sdk.PolicyKindFeaturePolicy, newFeaturePolicyId})
	})

	t.Run("unset policy safely", func(t *testing.T) {
		t.Skip("SNOW-3797718 - No such error is currently thrown")
		authenticationPolicy, authenticationPolicyCleanup := testClientHelper().AuthenticationPolicy.Create(t)
		t.Cleanup(authenticationPolicyCleanup)

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithUnset(*sdk.NewAccountUnsetRequest().WithAuthenticationPolicyUnset(*sdk.NewAccountAuthenticationPolicyUnsetRequest().WithAuthenticationPolicy(true))))
		assert.ErrorContains(t, err, fmt.Sprintf("Any policy of kind %s is not attached to ACCOUNT", sdk.PolicyKindAuthenticationPolicy))

		err = client.Accounts.UnsetPolicySafely(ctx, sdk.PolicyKindAuthenticationPolicy)
		assert.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithAuthenticationPolicySet(*sdk.NewAccountAuthenticationPolicySetRequest().WithAuthenticationPolicy(authenticationPolicy.ID()))))
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, expectedAccountPolicy{sdk.PolicyKindAuthenticationPolicy, authenticationPolicy.ID()})

		err = client.Accounts.UnsetPolicySafely(ctx, sdk.PolicyKindAuthenticationPolicy)
		assert.NoError(t, err)
		assertThatNoPolicyIsSetOnAccount(t)
	})

	t.Run("unset all", func(t *testing.T) {
		authPolicy, authPolicyCleanup := testClientHelper().AuthenticationPolicy.Create(t)
		t.Cleanup(authPolicyCleanup)

		featurePolicyId, featurePolicyCleanup := testClientHelper().FeaturePolicy.Create(t)
		t.Cleanup(featurePolicyCleanup)

		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)

		sessionPolicy, sessionPolicyCleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicyCleanup)

		packagesPolicyId, packagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(packagesPolicyCleanup)

		warehouseId := testClientHelper().Ids.WarehouseId()

		eventTable, eventTableCleanup := testClientHelper().EventTable.Create(t)
		t.Cleanup(eventTableCleanup)

		externalVolumeId, externalVolumeCleanup := testClientHelper().ExternalVolume.Create(t)
		t.Cleanup(externalVolumeCleanup)

		createNetworkPolicyRequest := sdk.NewCreateNetworkPolicyRequest(testClientHelper().Ids.RandomAccountObjectIdentifier()).WithAllowedIpList([]sdk.IPRequest{*sdk.NewIPRequest("0.0.0.0/0")})
		networkPolicy, networkPolicyCleanup := testClientHelper().NetworkPolicy.CreateNetworkPolicyWithRequest(t, createNetworkPolicyRequest)
		t.Cleanup(networkPolicyCleanup)

		stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(stageCleanup)

		t.Cleanup(func() {
			err := client.Accounts.UnsetAllPoliciesSafely(ctx)
			assert.NoError(t, err)
			err = client.Accounts.UnsetAllParameters(ctx)
			require.NoError(t, err)
		})

		err := client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithAuthenticationPolicySet(*sdk.NewAccountAuthenticationPolicySetRequest().WithAuthenticationPolicy(authPolicy.ID()))))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithFeaturePolicySet(*sdk.NewAccountFeaturePolicySetRequest().WithFeaturePolicy(featurePolicyId))))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithPackagesPolicy(packagesPolicyId)))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithPasswordPolicy(passwordPolicy.ID())))
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().WithSet(*sdk.NewAccountSetRequest().WithSessionPolicySet(*sdk.NewAccountSessionPolicySetRequest().WithSessionPolicy(sessionPolicy.ID()))))
		require.NoError(t, err)

		// TODO(SNOW-2138715): Test all parameters, the following parameters were not tested due to more complex setup:
		// - ActivePythonProfiler
		// - CatalogSync
		// - EnableInternalStagesPrivatelink
		// - PythonProfilerModules
		// - S3StageVpceDnsName
		// - SimulatedDataSharingConsumer
		err = client.Accounts.Alter(ctx, sdk.NewAlterAccountRequest().
			WithSet(*sdk.NewAccountSetRequest().WithParameters(
				*sdk.NewAccountParametersRequest().
					WithAbortDetachedQuery(true).
					WithAllowClientMfaCaching(true).
					WithAllowIdToken(true).
					WithAutocommit(false).
					WithBaseLocationPrefix("STORAGE_BASE_URL/").
					WithBinaryInputFormat(sdk.BinaryInputFormatBase64).
					WithBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
					WithCatalog(helpers.TestDatabaseCatalog).
					WithClientEnableLogInfoStatementParameters(true).
					WithClientEncryptionKeySize(256).
					WithClientMemoryLimit(1540).
					WithClientMetadataRequestUseConnectionCtx(true).
					WithClientMetadataUseSessionDatabase(true).
					WithClientPrefetchThreads(5).
					WithClientResultChunkSize(159).
					WithClientResultColumnCaseInsensitive(true).
					WithClientSessionKeepAlive(true).
					WithClientSessionKeepAliveHeartbeatFrequency(3599).
					WithClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
					WithCortexEnabledCrossRegion("ANY_REGION").
					WithCortexModelsAllowlist("All").
					WithCsvTimestampFormat("YYYY-MM-DD").
					WithDataRetentionTimeInDays(2).
					WithDateInputFormat("YYYY-MM-DD").
					WithDateOutputFormat("YYYY-MM-DD").
					WithDefaultDdlCollation(sdk.StringAllowEmpty{Value: "en-cs"}).
					WithDefaultNotebookComputePoolCpu("CPU_X64_S").
					WithDefaultNotebookComputePoolGpu("GPU_NV_S").
					WithDefaultNullOrdering(sdk.DefaultNullOrderingFirst).
					WithDefaultStreamlitComputePool("SYSTEM_COMPUTE_POOL_GPU").
					WithDefaultStreamlitNotebookWarehouse(warehouseId).
					WithDisableUiDownloadButton(true).
					WithDisableUserPrivilegeGrants(true).
					WithEnableAutomaticSensitiveDataClassificationLog(false).
					WithEnableEgressCostOptimizer(false).
					WithEnableIdentifierFirstLogin(false).
					WithEnableTriSecretAndRekeyOptOutForImageRepository(true).
					WithEnableTriSecretAndRekeyOptOutForSpcsBlockStorage(true).
					WithEnableUnhandledExceptionsReporting(false).
					WithEnableUnloadPhysicalTypeOptimization(false).
					WithEnableUnredactedQuerySyntaxError(true).
					WithEnableUnredactedSecureObjectError(true).
					WithEnforceNetworkRulesForInternalStages(true).
					WithErrorOnNondeterministicMerge(false).
					WithErrorOnNondeterministicUpdate(true).
					WithEventTable(eventTable.ID()).
					WithExternalOauthAddPrivilegedRolesToBlockedList(false).
					WithExternalVolume(externalVolumeId).
					WithGeographyOutputFormat(sdk.GeographyOutputFormatWKT).
					WithGeometryOutputFormat(sdk.GeometryOutputFormatWKT).
					WithHybridTableLockTimeout(3599).
					WithInitialReplicationSizeLimitInTb("9.9").
					WithJdbcTreatDecimalAsInt(false).
					WithJdbcTreatTimestampNtzAsUtc(true).
					WithJdbcUseSessionTimezone(false).
					WithJsonIndent(4).
					WithJsTreatIntegerAsBigint(true).
					WithListingAutoFulfillmentReplicationRefreshSchedule("2 minutes").
					WithLockTimeout(43201).
					WithLogLevel(sdk.LogLevelInfo).
					WithMaxConcurrencyLevel(7).
					WithMaxDataExtensionTimeInDays(13).
					WithMetricLevel(sdk.MetricLevelAll).
					WithMinDataRetentionTimeInDays(1).
					WithMultiStatementCount(0).
					WithNetworkPolicy(networkPolicy.ID()).
					WithNoorderSequenceAsDefault(false).
					WithOauthAddPrivilegedRolesToBlockedList(false).
					WithOdbcTreatDecimalAsInt(true).
					WithPeriodicDataRekeying(false).
					WithPipeExecutionPaused(true).
					WithPreventUnloadToInlineUrl(true).
					WithPreventUnloadToInternalStages(true).
					WithPythonProfilerTargetStage(stage.ID()).
					WithQueryTag("test-query-tag").
					WithQuotedIdentifiersIgnoreCase(true).
					WithReplaceInvalidCharacters(true).
					WithRequireStorageIntegrationForStageCreation(true).
					WithRequireStorageIntegrationForStageOperation(true).
					WithRowsPerResultset(1000).
					WithSearchPath("$current, $public").
					WithServerlessTaskMaxStatementSize(sdk.WarehouseSize("6X-LARGE")).
					WithServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
					WithSsoLoginPage(true).
					WithStatementQueuedTimeoutInSeconds(1).
					WithStatementTimeoutInSeconds(1).
					WithStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
					WithStrictJsonOutput(true).
					WithSuspendTaskAfterNumFailures(3).
					WithTaskAutoRetryAttempts(3).
					WithTimestampDayIsAlways24H(true).
					WithTimestampInputFormat("YYYY-MM-DD").
					WithTimestampLtzOutputFormat("YYYY-MM-DD").
					WithTimestampNtzOutputFormat("YYYY-MM-DD").
					WithTimestampOutputFormat("YYYY-MM-DD").
					WithTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
					WithTimestampTzOutputFormat("YYYY-MM-DD").
					WithTimezone("Europe/London").
					WithTimeInputFormat("YYYY-MM-DD").
					WithTimeOutputFormat("YYYY-MM-DD").
					WithTraceLevel(sdk.TraceLevelPropagate).
					WithTransactionAbortOnError(true).
					WithTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
					WithTwoDigitCenturyStart(1971).
					WithUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
					WithUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeX6Large).
					WithUserTaskMinimumTriggerIntervalInSeconds(10).
					WithUserTaskTimeoutMs(10).
					WithUseCachedResult(false).
					WithWeekOfYearPolicy(1).
					WithWeekStart(1),
			)))
		require.NoError(t, err)

		err = client.Accounts.UnsetAll(ctx)
		require.NoError(t, err)

		objectparametersassert.AccountParameters(t, id).HasAllDefaults()

		assertThatNoPolicyIsSetOnAccount(t)
	})
}
