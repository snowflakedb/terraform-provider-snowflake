//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-1920887): Some of the account features cannot be currently tested as they require two Snowflake organizations
// TODO(SNOW-1342761): Adjust the tests, so they can be run in their own pipeline
// For now, those tests should be run manually. The account/admin user running those tests is required to:
// - Be privileged with ORGADMIN and ACCOUNTADMIN roles.
// - Shouldn't be any of the "main" accounts/admin users, because those tests alter the current account.

func TestInt_Account(t *testing.T) {
	testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	client := testClient(t)
	ctx := testContext(t)
	currentAccountName := testClientHelper().Context.CurrentAccountName(t)

	assertAccountQueriedByOrgAdmin := func(t *testing.T, account sdk.Account, accountName string) {
		t.Helper()
		assert.NotEmpty(t, account.OrganizationName)
		assert.Equal(t, accountName, account.AccountName)
		assert.Nil(t, account.RegionGroup)
		assert.NotEmpty(t, account.SnowflakeRegion)
		assert.Equal(t, sdk.EditionEnterprise, *account.Edition)
		assert.NotEmpty(t, *account.AccountURL)
		assert.NotEmpty(t, *account.CreatedOn)
		assert.Equal(t, "SNOWFLAKE", *account.Comment)
		assert.NotEmpty(t, account.AccountLocator)
		assert.NotEmpty(t, *account.AccountLocatorUrl)
		assert.Zero(t, *account.ManagedAccounts)
		assert.NotEmpty(t, *account.ConsumptionBillingEntityName)
		assert.Nil(t, account.MarketplaceConsumerBillingEntityName)
		assert.NotNil(t, account.MarketplaceProviderBillingEntityName)
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

	assertCreateResponse := func(t *testing.T, response *sdk.AccountCreateResponse, account sdk.Account) {
		t.Helper()
		require.NotNil(t, response)
		assert.Equal(t, account.AccountLocator, response.AccountLocator)
		assert.Equal(t, *account.AccountLocatorUrl, response.AccountLocatorUrl)
		assert.Equal(t, account.AccountName, response.AccountName)
		assert.Equal(t, *account.AccountURL, response.Url)
		assert.Equal(t, account.OrganizationName, response.OrganizationName)
		assert.Equal(t, *account.Edition, response.Edition)
		assert.NotEmpty(t, response.RegionGroup)
		assert.NotEmpty(t, response.Cloud)
		assert.NotEmpty(t, response.Region)
	}

	t.Run("create: minimal", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSensitiveAccountObjectIdentifier()
		name := random.AdminName()
		password := random.Password()
		email := random.Email()

		createResponse, err := client.Accounts.Create(ctx, id, &sdk.CreateAccountOptions{
			AdminName:     name,
			AdminPassword: sdk.String(password),
			Email:         email,
			Edition:       sdk.EditionStandard,
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
		assertCreateResponse(t, createResponse, *acc)
	})

	t.Run("create: user type service", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSensitiveAccountObjectIdentifier()
		name := random.AdminName()
		key, _ := random.GenerateRSAPublicKey(t)
		email := random.Email()

		createResponse, err := client.Accounts.Create(ctx, id, &sdk.CreateAccountOptions{
			AdminName:         name,
			AdminRSAPublicKey: sdk.String(key),
			AdminUserType:     sdk.Pointer(sdk.UserTypeService),
			Email:             email,
			Edition:           sdk.EditionStandard,
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
		assertCreateResponse(t, createResponse, *acc)
	})

	t.Run("create: user type legacy service", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSensitiveAccountObjectIdentifier()
		name := random.AdminName()
		password := random.Password()
		email := random.Email()

		createResponse, err := client.Accounts.Create(ctx, id, &sdk.CreateAccountOptions{
			AdminName:     name,
			AdminPassword: sdk.String(password),
			AdminUserType: sdk.Pointer(sdk.UserTypeLegacyService),
			Email:         email,
			Edition:       sdk.EditionStandard,
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
		assertCreateResponse(t, createResponse, *acc)
	})

	t.Run("create: complete", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSensitiveAccountObjectIdentifier()
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

		createResponse, err := client.Accounts.Create(ctx, id, &sdk.CreateAccountOptions{
			AdminName:          name,
			AdminPassword:      sdk.String(password),
			FirstName:          sdk.String("firstName"),
			LastName:           sdk.String("lastName"),
			Email:              email,
			MustChangePassword: sdk.Bool(true),
			Edition:            sdk.EditionStandard,
			RegionGroup:        sdk.String("PUBLIC"),
			Region:             sdk.String(currentRegion.SnowflakeRegion),
			Comment:            sdk.String(comment),
			// TODO(SNOW-1895880): with polaris Snowflake returns an error saying: "invalid property polaris for account"
			// Polaris: sdk.Bool(true),
		})
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Account.DropFunc(t, id))

		acc, err := client.Accounts.ShowByID(ctx, id)
		require.NoError(t, err)
		require.Equal(t, id, acc.ID())
		assertCreateResponse(t, createResponse, *acc)
	})

	t.Run("alter: set / unset is org admin", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		require.Equal(t, false, *account.IsOrgAdmin)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
				Name:     account.ID(),
				OrgAdmin: true,
			},
		})
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, true, *acc.IsOrgAdmin)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			SetIsOrgAdmin: &sdk.AccountSetIsOrgAdmin{
				Name:     account.ID(),
				OrgAdmin: false,
			},
		})
		require.NoError(t, err)

		acc, err = client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, false, *acc.IsOrgAdmin)
	})

	t.Run("alter: rename", func(t *testing.T) {
		oldAccount, oldAccountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(oldAccountCleanup)

		newName := testClientHelper().Ids.RandomSensitiveAccountObjectIdentifier()
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:    oldAccount.ID(),
				NewName: newName,
			},
		})
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

		newName := testClientHelper().Ids.RandomSensitiveAccountObjectIdentifier()
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:       account.ID(),
				NewName:    newName,
				SaveOldURL: sdk.Bool(false),
			},
		})
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)

		acc, err := client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.NotEqual(t, account.AccountURL, acc.AccountURL)
		require.Empty(t, acc.OldAccountURL)
	})

	t.Run("alter: drop url when there's no old url", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Drop: &sdk.AccountDrop{
				Name:   account.ID(),
				OldUrl: sdk.Bool(true),
			},
		})
		require.ErrorContains(t, err, "The account has no old url")
	})

	t.Run("alter: drop url after rename", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		newName := testClientHelper().Ids.RandomSensitiveAccountObjectIdentifier()
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Rename: &sdk.AccountRename{
				Name:    account.ID(),
				NewName: newName,
			},
		})
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.NotEmpty(t, acc.OldAccountURL)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Drop: &sdk.AccountDrop{
				Name:   newName,
				OldUrl: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		acc, err = client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.Empty(t, acc.OldAccountURL)
	})

	t.Run("drop: without options", func(t *testing.T) {
		err := client.Accounts.Drop(ctx, NonExistingAccountObjectIdentifier, 3, &sdk.DropAccountOptions{})
		require.Error(t, err)

		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		err = client.Accounts.Drop(ctx, account.ID(), 3, &sdk.DropAccountOptions{})
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("drop: with if exists", func(t *testing.T) {
		err := client.Accounts.Drop(ctx, NonExistingAccountObjectIdentifier, 3, &sdk.DropAccountOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)

		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		err = client.Accounts.Drop(ctx, account.ID(), 3, &sdk.DropAccountOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)

		_, err = client.Accounts.ShowByID(ctx, account.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("undrop", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		require.NoError(t, testClientHelper().Account.Drop(t, account.ID()))

		err := client.Accounts.Undrop(ctx, account.ID())
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, account.ID(), acc.ID())
	})

	t.Run("show: with like", func(t *testing.T) {
		currentAccount := testClientHelper().Context.CurrentAccount(t)
		accounts, err := client.Accounts.Show(ctx, &sdk.ShowAccountOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(currentAccount),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(accounts))
		assertAccountQueriedByOrgAdmin(t, accounts[0], currentAccountName)
	})

	t.Run("show: with history", func(t *testing.T) {
		currentAccount := testClientHelper().Context.CurrentAccount(t)
		accounts, err := client.Accounts.Show(ctx, &sdk.ShowAccountOptions{
			History: sdk.Bool(true),
			Like: &sdk.Like{
				Pattern: sdk.String(currentAccount),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(accounts))
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
		accounts, err := client.Accounts.Show(ctx, &sdk.ShowAccountOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(currentAccount),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(accounts))
		assertAccountQueriedByAccountAdmin(t, accounts[0], currentAccountName)
	})
}

func TestInt_Account_SelfAlter(t *testing.T) {
	t.Skip("TODO(SNOW-1920881): Adjust the test so that self alters will be done on newly created account - not the main test one")
	testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	// This client should be operating on a different account than the "main" one (because it will be altered here).
	// Cannot use a newly created account because ORGADMIN role is necessary,
	// and it is propagated only after some time (e.g., 1 hour) making it hard to automate.
	client := testClient(t)
	ctx := testContext(t)
	t.Cleanup(testClientHelper().Role.UseRole(t, snowflakeroles.Accountadmin))

	assertParameterIsDefault := func(t *testing.T, parameters []*sdk.Parameter, parameterKey string) {
		t.Helper()
		param, err := collections.FindFirst(parameters, func(parameter *sdk.Parameter) bool { return parameter.Key == parameterKey })
		require.NoError(t, err)
		require.NotNil(t, param)
		require.Equal(t, (*param).Default, (*param).Value)
		require.Equal(t, sdk.ParameterType(""), (*param).Level)
	}

	assertParameterValueSetOnAccount := func(t *testing.T, parameters []*sdk.Parameter, parameterKey string, parameterValue string) {
		t.Helper()
		param, err := collections.FindFirst(parameters, func(parameter *sdk.Parameter) bool { return parameter.Key == parameterKey })
		require.NoError(t, err)
		require.NotNil(t, param)
		require.Equal(t, parameterValue, (*param).Value)
		require.Equal(t, sdk.ParameterTypeAccount, (*param).Level)
	}

	t.Run("set / unset parameters", func(t *testing.T) {
		parameters, err := client.Accounts.ShowParameters(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, parameters)

		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterMinDataRetentionTimeInDays))
		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterJsonIndent))
		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterUserTaskTimeoutMs))
		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterEnableUnredactedQuerySyntaxError))

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				Parameters: &sdk.AccountLevelParameters{
					AccountParameters: &sdk.AccountParameters{
						MinDataRetentionTimeInDays: sdk.Int(15), // default is 0
					},
					SessionParameters: &sdk.SessionParameters{
						JsonIndent: sdk.Int(8), // default is 2
					},
					ObjectParameters: &sdk.ObjectParameters{
						UserTaskTimeoutMs: sdk.Int(100), // default is 3600000
					},
					UserParameters: &sdk.UserParameters{
						EnableUnredactedQuerySyntaxError: sdk.Bool(true), // default is false
					},
				},
			},
		})
		require.NoError(t, err)

		parameters, err = client.Accounts.ShowParameters(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, parameters)

		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterMinDataRetentionTimeInDays), "15")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterJsonIndent), "8")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUserTaskTimeoutMs), "100")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableUnredactedQuerySyntaxError), "true")

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				Parameters: &sdk.AccountLevelParametersUnset{
					AccountParameters: &sdk.AccountParametersUnset{
						MinDataRetentionTimeInDays: sdk.Bool(true),
					},
					SessionParameters: &sdk.SessionParametersUnset{
						JsonIndent: sdk.Bool(true),
					},
					ObjectParameters: &sdk.ObjectParametersUnset{
						UserTaskTimeoutMs: sdk.Bool(true),
					},
					UserParameters: &sdk.UserParametersUnset{
						EnableUnredactedQuerySyntaxError: sdk.Bool(true),
					},
				},
			},
		})
		require.NoError(t, err)

		parameters, err = client.Accounts.ShowParameters(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, parameters)

		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterMinDataRetentionTimeInDays))
		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterJsonIndent))
		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterUserTaskTimeoutMs))
		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterEnableUnredactedQuerySyntaxError))
	})

	assertPolicySet := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(client.GetAccountLocator()), sdk.PolicyEntityDomainAccount)
		require.NoError(t, err)
		_, err = collections.FindFirst(policies, func(reference sdk.PolicyReference) bool {
			return reference.PolicyName == id.Name()
		})
		require.NoError(t, err)
	}

	assertPolicyNotSet := func(t *testing.T) {
		t.Helper()

		policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(client.GetAccountLocator()), sdk.PolicyEntityDomainAccount)
		require.Len(t, policies, 0)
		require.NoError(t, err)
	}

	t.Run("set / unset resource monitor", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		require.Nil(t, resourceMonitor.Level)
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				ResourceMonitor: resourceMonitor.ID(),
			},
		})
		require.NoError(t, err)

		resourceMonitor, err = testClientHelper().ResourceMonitor.Show(t, resourceMonitor.ID())
		require.NoError(t, err)
		require.NotNil(t, resourceMonitor.Level)
		require.Equal(t, sdk.ResourceMonitorLevelAccount, *resourceMonitor.Level)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				ResourceMonitor: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		resourceMonitor, err = testClientHelper().ResourceMonitor.Show(t, resourceMonitor.ID())
		require.NoError(t, err)
		require.Nil(t, resourceMonitor.Level)
	})

	t.Run("set / unset policies", func(t *testing.T) {
		authPolicy, authPolicyCleanup := testClientHelper().AuthenticationPolicy.Create(t)
		t.Cleanup(authPolicyCleanup)

		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)

		sessionPolicy, sessionPolicyCleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicyCleanup)

		packagesPolicyId, packagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(packagesPolicyCleanup)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				PackagesPolicy: packagesPolicyId,
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
				Unset: &sdk.AccountUnset{
					PackagesPolicy: sdk.Bool(true),
				},
			}))
		})

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				PasswordPolicy: passwordPolicy.ID(),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
				Unset: &sdk.AccountUnset{
					PasswordPolicy: sdk.Bool(true),
				},
			}))
		})

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				SessionPolicy: sessionPolicy.ID(),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
				Unset: &sdk.AccountUnset{
					SessionPolicy: sdk.Bool(true),
				},
			}))
		})

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				AuthenticationPolicy: authPolicy.ID(),
			},
		})
		require.NoError(t, err)
		t.Cleanup(func() {
			assert.NoError(t, client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
				Unset: &sdk.AccountUnset{
					AuthenticationPolicy: sdk.Bool(true),
				},
			}))
		})

		assertPolicySet(t, authPolicy.ID())
		assertPolicySet(t, passwordPolicy.ID())
		assertPolicySet(t, sessionPolicy.ID())
		assertPolicySet(t, packagesPolicyId)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				PackagesPolicy: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				PasswordPolicy: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				SessionPolicy: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				AuthenticationPolicy: sdk.Bool(true),
			},
		})
		require.NoError(t, err)

		assertPolicyNotSet(t)
	})

	t.Run("force new packages policy", func(t *testing.T) {
		packagesPolicyId, packagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(packagesPolicyCleanup)

		newPackagesPolicyId, newPackagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(newPackagesPolicyCleanup)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				PackagesPolicy: packagesPolicyId,
			},
		})
		require.NoError(t, err)
		assertPolicySet(t, packagesPolicyId)
		t.Cleanup(func() {
			err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
				Unset: &sdk.AccountUnset{
					PackagesPolicy: sdk.Bool(true),
				},
			})
			require.NoError(t, err)
			assertPolicyNotSet(t)
		})

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				PackagesPolicy: newPackagesPolicyId,
			},
		})
		require.Error(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				PackagesPolicy: newPackagesPolicyId,
				Force:          sdk.Bool(true),
			},
		})
		require.NoError(t, err)
		assertPolicySet(t, newPackagesPolicyId)
	})
}
