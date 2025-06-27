//go:build !account_level_tests

package testint

import (
	"fmt"
	"testing"

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
		id := sdk.NewAccountObjectIdentifier(random.AccountName())
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
		id := sdk.NewAccountObjectIdentifier(random.AccountName())
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
		id := sdk.NewAccountObjectIdentifier(random.AccountName())
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

		createResponse, err := client.Accounts.Create(ctx, id, &sdk.CreateAccountOptions{
			AdminName:                name,
			AdminPassword:            sdk.String(password),
			FirstName:                sdk.String("firstName"),
			LastName:                 sdk.String("lastName"),
			Email:                    email,
			MustChangePassword:       sdk.Bool(true),
			Edition:                  sdk.EditionStandard,
			RegionGroup:              sdk.String("PUBLIC"),
			Region:                   sdk.String(currentRegion.SnowflakeRegion),
			Comment:                  sdk.String(comment),
			ConsumptionBillingEntity: sdk.String(defaultConsumptionBillingEntity),
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

		require.False(t, *account.IsOrgAdmin)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Name: sdk.Pointer(account.ID()),
			Set:  &sdk.AccountSet{OrgAdmin: sdk.Bool(true)},
		})
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.True(t, *acc.IsOrgAdmin)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Name: sdk.Pointer(account.ID()),
			Set:  &sdk.AccountSet{OrgAdmin: sdk.Bool(false)},
		})
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

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Name: sdk.Pointer(oldAccount.ID()),
			Rename: &sdk.AccountRename{
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

		newName := sdk.NewAccountObjectIdentifier(random.AccountName())
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Name: sdk.Pointer(account.ID()),
			Rename: &sdk.AccountRename{
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

	t.Run("alter: set / unset consumption billing entity", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		require.Equal(t, defaultConsumptionBillingEntity, *account.ConsumptionBillingEntityName)

		// We are not able to create consumption billing entities, because of that, we use the default one.
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Name: sdk.Pointer(account.ID()),
			Set:  &sdk.AccountSet{ConsumptionBillingEntity: sdk.String(defaultConsumptionBillingEntity)},
		})
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, defaultConsumptionBillingEntity, *acc.ConsumptionBillingEntityName)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Name:  sdk.Pointer(account.ID()),
			Unset: &sdk.AccountUnset{ConsumptionBillingEntity: sdk.Bool(true)},
		})
		require.NoError(t, err)

		acc, err = client.Accounts.ShowByID(ctx, account.ID())
		require.NoError(t, err)
		require.Equal(t, defaultConsumptionBillingEntity, *acc.ConsumptionBillingEntityName)
	})

	t.Run("alter: drop url when there's no old url", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Name: sdk.Pointer(account.ID()),
			Drop: &sdk.AccountDrop{
				OldUrl: sdk.Bool(true),
			},
		})
		require.ErrorContains(t, err, "The account has no old url")
	})

	t.Run("alter: drop url after rename", func(t *testing.T) {
		account, accountCleanup := testClientHelper().Account.Create(t)
		t.Cleanup(accountCleanup)

		newName := sdk.NewAccountObjectIdentifier(random.AccountName())
		t.Cleanup(testClientHelper().Account.DropFunc(t, newName))

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Name: sdk.Pointer(account.ID()),
			Rename: &sdk.AccountRename{
				NewName: newName,
			},
		})
		require.NoError(t, err)

		acc, err := client.Accounts.ShowByID(ctx, newName)
		require.NoError(t, err)
		require.NotEmpty(t, acc.OldAccountURL)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Name: sdk.Pointer(newName),
			Drop: &sdk.AccountDrop{
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
		assert.Len(t, accounts, 1)
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
		accounts, err := client.Accounts.Show(ctx, &sdk.ShowAccountOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(currentAccount),
			},
		})
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

	assertParameterIsDefault := func(t *testing.T, parameters []*sdk.Parameter, parameterKey string) {
		t.Helper()
		param, err := collections.FindFirst(parameters, func(parameter *sdk.Parameter) bool { return parameter.Key == parameterKey })
		// TODO(SNOW-1325308): Improve collections.FindFirst error message to include more detail about missing item
		require.NoError(t, err, "parameter %v not found", parameterKey)
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

	t.Run("set / unset legacy parameters", func(t *testing.T) {
		parameters, err := client.Accounts.ShowParameters(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, parameters)

		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterMinDataRetentionTimeInDays))
		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterJsonIndent))
		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterUserTaskTimeoutMs))
		assertParameterIsDefault(t, parameters, string(sdk.AccountParameterEnableUnredactedQuerySyntaxError))

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				LegacyParameters: &sdk.AccountLevelParameters{
					AccountParameters: &sdk.LegacyAccountParameters{
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
				LegacyParameters: &sdk.AccountLevelParametersUnset{
					AccountParameters: &sdk.LegacyAccountParametersUnset{
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

	t.Run("set / unset parameters", func(t *testing.T) {
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

		// TODO(SNOW-2138715): Test all parameters, the following parameters were not tested due to more complex setup:
		// - ActivePythonProfiler
		// - CatalogSync
		// - EnableInternalStagesPrivatelink
		// - PythonProfilerModules
		// - S3StageVpceDnsName
		// - SamlIdentityProvider
		// - SimulatedDataSharingConsumer
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				Parameters: &sdk.AccountParameters{
					AbortDetachedQuery:                               sdk.Bool(true),
					AllowClientMFACaching:                            sdk.Bool(true),
					AllowIDToken:                                     sdk.Bool(true),
					Autocommit:                                       sdk.Bool(false),
					BaseLocationPrefix:                               sdk.String("STORAGE_BASE_URL/"),
					BinaryInputFormat:                                sdk.Pointer(sdk.BinaryInputFormatBase64),
					BinaryOutputFormat:                               sdk.Pointer(sdk.BinaryOutputFormatBase64),
					Catalog:                                          sdk.String(helpers.TestDatabaseCatalog.Name()),
					ClientEnableLogInfoStatementParameters:           sdk.Bool(true),
					ClientEncryptionKeySize:                          sdk.Int(256),
					ClientMemoryLimit:                                sdk.Int(1540),
					ClientMetadataRequestUseConnectionCtx:            sdk.Bool(true),
					ClientMetadataUseSessionDatabase:                 sdk.Bool(true),
					ClientPrefetchThreads:                            sdk.Int(5),
					ClientResultChunkSize:                            sdk.Int(159),
					ClientResultColumnCaseInsensitive:                sdk.Bool(true),
					ClientSessionKeepAlive:                           sdk.Bool(true),
					ClientSessionKeepAliveHeartbeatFrequency:         sdk.Int(3599),
					ClientTimestampTypeMapping:                       sdk.Pointer(sdk.ClientTimestampTypeMappingNtz),
					CortexEnabledCrossRegion:                         sdk.String("ANY_REGION"),
					CortexModelsAllowlist:                            sdk.String("All"),
					CsvTimestampFormat:                               sdk.String("YYYY-MM-DD"),
					DataRetentionTimeInDays:                          sdk.Int(2),
					DateInputFormat:                                  sdk.String("YYYY-MM-DD"),
					DateOutputFormat:                                 sdk.String("YYYY-MM-DD"),
					DefaultDDLCollation:                              sdk.String("en-cs"),
					DefaultNotebookComputePoolCpu:                    sdk.String("CPU_X64_S"),
					DefaultNotebookComputePoolGpu:                    sdk.String("GPU_NV_S"),
					DefaultNullOrdering:                              sdk.Pointer(sdk.DefaultNullOrderingFirst),
					DefaultStreamlitNotebookWarehouse:                sdk.Pointer(warehouseId),
					DisableUiDownloadButton:                          sdk.Bool(true),
					DisableUserPrivilegeGrants:                       sdk.Bool(true),
					EnableAutomaticSensitiveDataClassificationLog:    sdk.Bool(false),
					EnableEgressCostOptimizer:                        sdk.Bool(false),
					EnableIdentifierFirstLogin:                       sdk.Bool(false),
					EnableTriSecretAndRekeyOptOutForImageRepository:  sdk.Bool(true),
					EnableTriSecretAndRekeyOptOutForSpcsBlockStorage: sdk.Bool(true),
					EnableUnhandledExceptionsReporting:               sdk.Bool(false),
					EnableUnloadPhysicalTypeOptimization:             sdk.Bool(false),
					EnableUnredactedQuerySyntaxError:                 sdk.Bool(true),
					EnableUnredactedSecureObjectError:                sdk.Bool(true),
					EnforceNetworkRulesForInternalStages:             sdk.Bool(true),
					ErrorOnNondeterministicMerge:                     sdk.Bool(false),
					ErrorOnNondeterministicUpdate:                    sdk.Bool(true),
					EventTable:                                       sdk.Pointer(eventTable.ID()),
					ExternalOAuthAddPrivilegedRolesToBlockedList:     sdk.Bool(false),
					ExternalVolume:                                   sdk.Pointer(externalVolumeId),
					GeographyOutputFormat:                            sdk.Pointer(sdk.GeographyOutputFormatWKT),
					GeometryOutputFormat:                             sdk.Pointer(sdk.GeometryOutputFormatWKT),
					HybridTableLockTimeout:                           sdk.Int(3599),
					InitialReplicationSizeLimitInTB:                  sdk.String("9.9"),
					JdbcTreatDecimalAsInt:                            sdk.Bool(false),
					JdbcTreatTimestampNtzAsUtc:                       sdk.Bool(true),
					JdbcUseSessionTimezone:                           sdk.Bool(false),
					JsonIndent:                                       sdk.Int(4),
					JsTreatIntegerAsBigInt:                           sdk.Bool(true),
					ListingAutoFulfillmentReplicationRefreshSchedule: sdk.String("2 minutes"),
					LockTimeout:                                      sdk.Int(43201),
					LogLevel:                                         sdk.Pointer(sdk.LogLevelInfo),
					MaxConcurrencyLevel:                              sdk.Int(7),
					MaxDataExtensionTimeInDays:                       sdk.Int(13),
					MetricLevel:                                      sdk.Pointer(sdk.MetricLevelAll),
					MinDataRetentionTimeInDays:                       sdk.Int(1),
					MultiStatementCount:                              sdk.Int(0),
					NetworkPolicy:                                    sdk.Pointer(networkPolicy.ID()),
					NoorderSequenceAsDefault:                         sdk.Bool(false),
					OAuthAddPrivilegedRolesToBlockedList:             sdk.Bool(false),
					OdbcTreatDecimalAsInt:                            sdk.Bool(true),
					PeriodicDataRekeying:                             sdk.Bool(false),
					PipeExecutionPaused:                              sdk.Bool(true),
					PreventUnloadToInlineURL:                         sdk.Bool(true),
					PreventUnloadToInternalStages:                    sdk.Bool(true),
					PythonProfilerTargetStage:                        sdk.Pointer(stage.ID()),
					QueryTag:                                         sdk.String("test-query-tag"),
					QuotedIdentifiersIgnoreCase:                      sdk.Bool(true),
					ReplaceInvalidCharacters:                         sdk.Bool(true),
					RequireStorageIntegrationForStageCreation:        sdk.Bool(true),
					RequireStorageIntegrationForStageOperation:       sdk.Bool(true),
					RowsPerResultset:                                 sdk.Int(1000),
					SearchPath:                                       sdk.String("$current, $public"),
					ServerlessTaskMaxStatementSize:                   sdk.Pointer(sdk.WarehouseSize("6X-LARGE")),
					ServerlessTaskMinStatementSize:                   sdk.Pointer(sdk.WarehouseSizeSmall),
					SsoLoginPage:                                     sdk.Bool(true),
					StatementQueuedTimeoutInSeconds:                  sdk.Int(1),
					StatementTimeoutInSeconds:                        sdk.Int(1),
					StorageSerializationPolicy:                       sdk.Pointer(sdk.StorageSerializationPolicyOptimized),
					StrictJsonOutput:                                 sdk.Bool(true),
					SuspendTaskAfterNumFailures:                      sdk.Int(3),
					TaskAutoRetryAttempts:                            sdk.Int(3),
					TimestampDayIsAlways24h:                          sdk.Bool(true),
					TimestampInputFormat:                             sdk.String("YYYY-MM-DD"),
					TimestampLtzOutputFormat:                         sdk.String("YYYY-MM-DD"),
					TimestampNtzOutputFormat:                         sdk.String("YYYY-MM-DD"),
					TimestampOutputFormat:                            sdk.String("YYYY-MM-DD"),
					TimestampTypeMapping:                             sdk.Pointer(sdk.TimestampTypeMappingLtz),
					TimestampTzOutputFormat:                          sdk.String("YYYY-MM-DD"),
					Timezone:                                         sdk.String("Europe/London"),
					TimeInputFormat:                                  sdk.String("YYYY-MM-DD"),
					TimeOutputFormat:                                 sdk.String("YYYY-MM-DD"),
					TraceLevel:                                       sdk.Pointer(sdk.TraceLevelPropagate),
					TransactionAbortOnError:                          sdk.Bool(true),
					TransactionDefaultIsolationLevel:                 sdk.Pointer(sdk.TransactionDefaultIsolationLevelReadCommitted),
					TwoDigitCenturyStart:                             sdk.Int(1971),
					UnsupportedDdlAction:                             sdk.Pointer(sdk.UnsupportedDDLActionFail),
					UserTaskManagedInitialWarehouseSize:              sdk.Pointer(sdk.WarehouseSizeX6Large),
					UserTaskMinimumTriggerIntervalInSeconds:          sdk.Int(10),
					UserTaskTimeoutMs:                                sdk.Int(10),
					UseCachedResult:                                  sdk.Bool(false),
					WeekOfYearPolicy:                                 sdk.Int(1),
					WeekStart:                                        sdk.Int(1),
				},
			},
		})
		require.NoError(t, err)

		parameters, err := client.Accounts.ShowParameters(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, parameters)

		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterAbortDetachedQuery), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterAllowClientMFACaching), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterAllowIDToken), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterAutocommit), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterBaseLocationPrefix), "STORAGE_BASE_URL/")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterBinaryInputFormat), string(sdk.BinaryInputFormatBase64))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterBinaryOutputFormat), string(sdk.BinaryOutputFormatBase64))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterCatalog), helpers.TestDatabaseCatalog.Name())
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientEnableLogInfoStatementParameters), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientEncryptionKeySize), "256")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientMemoryLimit), "1540")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientMetadataRequestUseConnectionCtx), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientMetadataUseSessionDatabase), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientPrefetchThreads), "5")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientResultChunkSize), "159")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientResultColumnCaseInsensitive), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientSessionKeepAlive), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientSessionKeepAliveHeartbeatFrequency), "3599")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterClientTimestampTypeMapping), string(sdk.ClientTimestampTypeMappingNtz))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterCortexEnabledCrossRegion), "ANY_REGION")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterCortexModelsAllowlist), "All")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterCsvTimestampFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDataRetentionTimeInDays), "2")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDateInputFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDateOutputFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDefaultDDLCollation), "en-cs")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDefaultNotebookComputePoolCpu), "CPU_X64_S")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDefaultNotebookComputePoolGpu), "GPU_NV_S")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDefaultNullOrdering), string(sdk.DefaultNullOrderingFirst))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDefaultStreamlitNotebookWarehouse), warehouseId.Name())
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDisableUiDownloadButton), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterDisableUserPrivilegeGrants), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableAutomaticSensitiveDataClassificationLog), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableEgressCostOptimizer), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableIdentifierFirstLogin), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableUnhandledExceptionsReporting), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableUnloadPhysicalTypeOptimization), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableUnredactedQuerySyntaxError), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnableUnredactedSecureObjectError), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEnforceNetworkRulesForInternalStages), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterErrorOnNondeterministicMerge), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterErrorOnNondeterministicUpdate), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterEventTable), eventTable.ID().FullyQualifiedName())
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterExternalVolume), externalVolumeId.Name())
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterGeographyOutputFormat), string(sdk.GeographyOutputFormatWKT))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterGeometryOutputFormat), string(sdk.GeometryOutputFormatWKT))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterHybridTableLockTimeout), "3599")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterInitialReplicationSizeLimitInTB), "9.9")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterJdbcTreatDecimalAsInt), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterJdbcTreatTimestampNtzAsUtc), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterJdbcUseSessionTimezone), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterJsonIndent), "4")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterJsTreatIntegerAsBigInt), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterListingAutoFulfillmentReplicationRefreshSchedule), "2 minutes")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterLockTimeout), "43201")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterLogLevel), string(sdk.LogLevelInfo))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterMaxConcurrencyLevel), "7")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterMaxDataExtensionTimeInDays), "13")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterMetricLevel), string(sdk.MetricLevelAll))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterMinDataRetentionTimeInDays), "1")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterMultiStatementCount), "0")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterNetworkPolicy), networkPolicy.ID().Name())
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterNoorderSequenceAsDefault), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterOdbcTreatDecimalAsInt), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterPeriodicDataRekeying), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterPipeExecutionPaused), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterPreventUnloadToInlineURL), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterPreventUnloadToInternalStages), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterQueryTag), "test-query-tag")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterQuotedIdentifiersIgnoreCase), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterReplaceInvalidCharacters), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterRequireStorageIntegrationForStageCreation), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterRequireStorageIntegrationForStageOperation), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterRowsPerResultset), "1000")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterSearchPath), "$current, $public")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterServerlessTaskMaxStatementSize), "6X-LARGE")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterServerlessTaskMinStatementSize), string(sdk.WarehouseSizeSmall))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterSsoLoginPage), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterStatementQueuedTimeoutInSeconds), "1")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterStatementTimeoutInSeconds), "1")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterStorageSerializationPolicy), string(sdk.StorageSerializationPolicyOptimized))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterStrictJsonOutput), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterSuspendTaskAfterNumFailures), "3")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTaskAutoRetryAttempts), "3")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimestampDayIsAlways24h), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimestampInputFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimestampLtzOutputFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimestampNtzOutputFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimestampOutputFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimestampTypeMapping), string(sdk.TimestampTypeMappingLtz))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimestampTzOutputFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimezone), "Europe/London")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimeInputFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTimeOutputFormat), "YYYY-MM-DD")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTraceLevel), string(sdk.TraceLevelPropagate))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTransactionAbortOnError), "true")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTransactionDefaultIsolationLevel), string(sdk.TransactionDefaultIsolationLevelReadCommitted))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterTwoDigitCenturyStart), "1971")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUnsupportedDdlAction), string(sdk.UnsupportedDDLActionFail))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUserTaskManagedInitialWarehouseSize), string(sdk.WarehouseSizeX6Large))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds), "10")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUserTaskTimeoutMs), "10")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUseCachedResult), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterWeekOfYearPolicy), "1")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterWeekStart), "1")

		err = client.Accounts.UnsetAllParameters(ctx)
		require.NoError(t, err)

		parameters, err = client.Accounts.ShowParameters(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, parameters)

		for _, parameter := range sdk.AllAccountParameters {
			assertParameterIsDefault(t, parameters, string(parameter))
		}
	})

	assertThatPolicyIsSetOnAccount := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(client.GetAccountLocator()), sdk.PolicyEntityDomainAccount)
		require.NoError(t, err)
		_, err = collections.FindFirst(policies, func(reference sdk.PolicyReference) bool {
			return reference.PolicyName == id.Name()
		})
		require.NoError(t, err)
	}

	assertThatNoPolicyIsSetOnAccount := func(t *testing.T) {
		t.Helper()

		policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(client.GetAccountLocator()), sdk.PolicyEntityDomainAccount)
		require.Empty(t, policies)
		require.NoError(t, err)
	}

	t.Run("set / unset resource monitor", func(t *testing.T) {
		resourceMonitor, resourceMonitorCleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitorCleanup)

		require.Nil(t, resourceMonitor.Level)
		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				ResourceMonitor: sdk.Pointer(resourceMonitor.ID()),
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

		featurePolicyId, featurePolicyCleanup := testClientHelper().FeaturePolicy.Create(t)
		t.Cleanup(featurePolicyCleanup)

		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)

		sessionPolicy, sessionPolicyCleanup := testClientHelper().SessionPolicy.CreateSessionPolicy(t)
		t.Cleanup(sessionPolicyCleanup)

		packagesPolicyId, packagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(packagesPolicyCleanup)

		t.Cleanup(func() {
			err := client.Accounts.UnsetAllPoliciesSafely(ctx)
			assert.NoError(t, err)
			assertThatNoPolicyIsSetOnAccount(t)
		})

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{AuthenticationPolicy: sdk.Pointer(authPolicy.ID())}})
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{FeaturePolicySet: &sdk.AccountFeaturePolicySet{FeaturePolicy: &featurePolicyId}}})
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{PackagesPolicy: &packagesPolicyId}})
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{PasswordPolicy: sdk.Pointer(passwordPolicy.ID())}})
		require.NoError(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{SessionPolicy: sdk.Pointer(sessionPolicy.ID())}})
		require.NoError(t, err)

		assertThatPolicyIsSetOnAccount(t, featurePolicyId)
		assertThatPolicyIsSetOnAccount(t, authPolicy.ID())
		assertThatPolicyIsSetOnAccount(t, passwordPolicy.ID())
		assertThatPolicyIsSetOnAccount(t, sessionPolicy.ID())
		assertThatPolicyIsSetOnAccount(t, packagesPolicyId)
	})

	t.Run("force new packages policy", func(t *testing.T) {
		packagesPolicyId, packagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(packagesPolicyCleanup)

		newPackagesPolicyId, newPackagesPolicyCleanup := testClientHelper().PackagesPolicy.Create(t)
		t.Cleanup(newPackagesPolicyCleanup)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{PackagesPolicy: &packagesPolicyId}})
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, packagesPolicyId)
		t.Cleanup(func() {
			err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: &sdk.AccountUnset{PackagesPolicy: sdk.Bool(true)}})
			require.NoError(t, err)
			assertThatNoPolicyIsSetOnAccount(t)
		})

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{PackagesPolicy: &newPackagesPolicyId}})
		require.Error(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{PackagesPolicy: &newPackagesPolicyId, Force: sdk.Bool(true)}})
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, newPackagesPolicyId)
	})

	t.Run("force new feature policy", func(t *testing.T) {
		featurePolicyId, featurePolicyCleanup := testClientHelper().FeaturePolicy.Create(t)
		t.Cleanup(featurePolicyCleanup)

		newFeaturePolicyId, newFeaturePolicyCleanup := testClientHelper().FeaturePolicy.Create(t)
		t.Cleanup(newFeaturePolicyCleanup)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{FeaturePolicySet: &sdk.AccountFeaturePolicySet{FeaturePolicy: &featurePolicyId}}})
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, featurePolicyId)
		t.Cleanup(func() {
			err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: &sdk.AccountUnset{FeaturePolicyUnset: &sdk.AccountFeaturePolicyUnset{FeaturePolicy: sdk.Bool(true)}}})
			require.NoError(t, err)
			assertThatNoPolicyIsSetOnAccount(t)
		})

		// Here we expect to get an error as there is another feature policy set on the account.
		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{FeaturePolicySet: &sdk.AccountFeaturePolicySet{FeaturePolicy: &newFeaturePolicyId}}})
		require.Error(t, err)

		// To set a new feature policy on the account without firstly unsetting it, we can use FORCE parameter.
		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{FeaturePolicySet: &sdk.AccountFeaturePolicySet{FeaturePolicy: &newFeaturePolicyId}, Force: sdk.Bool(true)}})
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, newFeaturePolicyId)
	})

	t.Run("unset policy safely", func(t *testing.T) {
		featurePolicyId, featurePolicyCleanup := testClientHelper().FeaturePolicy.Create(t)
		t.Cleanup(featurePolicyCleanup)

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Unset: &sdk.AccountUnset{FeaturePolicyUnset: &sdk.AccountFeaturePolicyUnset{FeaturePolicy: sdk.Bool(true)}}})
		assert.ErrorContains(t, err, fmt.Sprintf("Any policy of kind %s is not attached to ACCOUNT", sdk.PolicyKindFeaturePolicy))

		err = client.Accounts.UnsetPolicySafely(ctx, sdk.PolicyKindFeaturePolicy)
		assert.NoError(t, err)

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{FeaturePolicySet: &sdk.AccountFeaturePolicySet{FeaturePolicy: &featurePolicyId}}})
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, featurePolicyId)

		err = client.Accounts.UnsetPolicySafely(ctx, sdk.PolicyKindFeaturePolicy)
		assert.NoError(t, err)
		assertThatNoPolicyIsSetOnAccount(t)
	})
}
