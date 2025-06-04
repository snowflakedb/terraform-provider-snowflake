//go:build !account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
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
	//t.Skip("TODO(SNOW-1920881): Adjust the test so that self alters will be done on newly created account - not the main test one")
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

		err := client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Set: &sdk.AccountSet{
				Parameters: &sdk.AccountParameters{
					AbortDetachedQuery: sdk.Bool(true),
					// TODO: ActivePythonProfiler:                     sdk.Pointer(sdk.ActivePythonProfilerMemory),
					AllowClientMFACaching:                            sdk.Bool(true),
					AllowIDToken:                                     sdk.Bool(true),
					Autocommit:                                       sdk.Bool(false),
					BaseLocationPrefix:                               sdk.String("STORAGE_BASE_URL/"),
					BinaryInputFormat:                                sdk.Pointer(sdk.BinaryInputFormatBase64),
					BinaryOutputFormat:                               sdk.Pointer(sdk.BinaryOutputFormatBase64),
					Catalog:                                          sdk.String("SNOWFLAKE"),
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
					// TODO: PythonProfilerTargetStage:                  sdk.Pointer(stage.ID()),
					QueryTag:                                   sdk.String("test-query-tag"),
					QuotedIdentifiersIgnoreCase:                sdk.Bool(true),
					ReplaceInvalidCharacters:                   sdk.Bool(true),
					RequireStorageIntegrationForStageCreation:  sdk.Bool(true),
					RequireStorageIntegrationForStageOperation: sdk.Bool(true),
					RowsPerResultset:                           sdk.Int(1000),
					SearchPath:                                 sdk.String("$current, $public"),
					ServerlessTaskMaxStatementSize:             sdk.Pointer(sdk.WarehouseSizeXLarge),
					ServerlessTaskMinStatementSize:             sdk.Pointer(sdk.WarehouseSizeSmall),
					SsoLoginPage:                               sdk.Bool(true),
					StatementQueuedTimeoutInSeconds:            sdk.Int(1),
					StatementTimeoutInSeconds:                  sdk.Int(1),
					StorageSerializationPolicy:                 sdk.Pointer(sdk.StorageSerializationPolicyOptimized),
					StrictJsonOutput:                           sdk.Bool(true),
					SuspendTaskAfterNumFailures:                sdk.Int(3),
					TaskAutoRetryAttempts:                      sdk.Int(3),
					TimestampDayIsAlways24h:                    sdk.Bool(true),
					TimestampInputFormat:                       sdk.String("YYYY-MM-DD"),
					TimestampLtzOutputFormat:                   sdk.String("YYYY-MM-DD"),
					TimestampNtzOutputFormat:                   sdk.String("YYYY-MM-DD"),
					TimestampOutputFormat:                      sdk.String("YYYY-MM-DD"),
					TimestampTypeMapping:                       sdk.Pointer(sdk.TimestampTypeMappingLtz),
					TimestampTzOutputFormat:                    sdk.String("YYYY-MM-DD"),
					Timezone:                                   sdk.String("Europe/London"),
					TimeInputFormat:                            sdk.String("YYYY-MM-DD"),
					TimeOutputFormat:                           sdk.String("YYYY-MM-DD"),
					TraceLevel:                                 sdk.Pointer(sdk.TraceLevelPropagate),
					TransactionAbortOnError:                    sdk.Bool(true),
					TransactionDefaultIsolationLevel:           sdk.Pointer(sdk.TransactionDefaultIsolationLevelReadCommitted),
					TwoDigitCenturyStart:                       sdk.Int(1971),
					UnsupportedDdlAction:                       sdk.Pointer(sdk.UnsupportedDDLActionFail),
					UserTaskManagedInitialWarehouseSize:        sdk.Pointer(sdk.WarehouseSizeSmall),
					UserTaskMinimumTriggerIntervalInSeconds:    sdk.Int(10),
					UserTaskTimeoutMs:                          sdk.Int(10),
					UseCachedResult:                            sdk.Bool(false),
					WeekOfYearPolicy:                           sdk.Int(1),
					WeekStart:                                  sdk.Int(1),
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
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterCatalog), "SNOWFLAKE")
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
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterServerlessTaskMaxStatementSize), string(sdk.WarehouseSizeXLarge))
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
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUserTaskManagedInitialWarehouseSize), string(sdk.WarehouseSizeSmall))
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds), "10")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUserTaskTimeoutMs), "10")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterUseCachedResult), "false")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterWeekOfYearPolicy), "1")
		assertParameterValueSetOnAccount(t, parameters, string(sdk.AccountParameterWeekStart), "1")

		err = client.Accounts.Alter(ctx, &sdk.AlterAccountOptions{
			Unset: &sdk.AccountUnset{
				Parameters: &sdk.AccountParametersUnset{
					AbortDetachedQuery:                               sdk.Bool(true),
					ActivePythonProfiler:                             sdk.Bool(true),
					AllowClientMFACaching:                            sdk.Bool(true),
					AllowIDToken:                                     sdk.Bool(true),
					Autocommit:                                       sdk.Bool(true),
					BaseLocationPrefix:                               sdk.Bool(true),
					BinaryInputFormat:                                sdk.Bool(true),
					BinaryOutputFormat:                               sdk.Bool(true),
					Catalog:                                          sdk.Bool(true),
					CatalogSync:                                      sdk.Bool(true),
					ClientEnableLogInfoStatementParameters:           sdk.Bool(true),
					ClientEncryptionKeySize:                          sdk.Bool(true),
					ClientMemoryLimit:                                sdk.Bool(true),
					ClientMetadataRequestUseConnectionCtx:            sdk.Bool(true),
					ClientMetadataUseSessionDatabase:                 sdk.Bool(true),
					ClientPrefetchThreads:                            sdk.Bool(true),
					ClientResultChunkSize:                            sdk.Bool(true),
					ClientResultColumnCaseInsensitive:                sdk.Bool(true),
					ClientSessionKeepAlive:                           sdk.Bool(true),
					ClientSessionKeepAliveHeartbeatFrequency:         sdk.Bool(true),
					ClientTimestampTypeMapping:                       sdk.Bool(true),
					CortexEnabledCrossRegion:                         sdk.Bool(true),
					CortexModelsAllowlist:                            sdk.Bool(true),
					CsvTimestampFormat:                               sdk.Bool(true),
					DataRetentionTimeInDays:                          sdk.Bool(true),
					DateInputFormat:                                  sdk.Bool(true),
					DateOutputFormat:                                 sdk.Bool(true),
					DefaultDDLCollation:                              sdk.Bool(true),
					DefaultNotebookComputePoolCpu:                    sdk.Bool(true),
					DefaultNotebookComputePoolGpu:                    sdk.Bool(true),
					DefaultNullOrdering:                              sdk.Bool(true),
					DefaultStreamlitNotebookWarehouse:                sdk.Bool(true),
					DisableUiDownloadButton:                          sdk.Bool(true),
					DisableUserPrivilegeGrants:                       sdk.Bool(true),
					EnableAutomaticSensitiveDataClassificationLog:    sdk.Bool(true),
					EnableEgressCostOptimizer:                        sdk.Bool(true),
					EnableIdentifierFirstLogin:                       sdk.Bool(true),
					EnableInternalStagesPrivatelink:                  sdk.Bool(true),
					EnableTriSecretAndRekeyOptOutForImageRepository:  sdk.Bool(true),
					EnableTriSecretAndRekeyOptOutForSpcsBlockStorage: sdk.Bool(true),
					EnableUnhandledExceptionsReporting:               sdk.Bool(true),
					EnableUnloadPhysicalTypeOptimization:             sdk.Bool(true),
					EnableUnredactedQuerySyntaxError:                 sdk.Bool(true),
					EnableUnredactedSecureObjectError:                sdk.Bool(true),
					EnforceNetworkRulesForInternalStages:             sdk.Bool(true),
					ErrorOnNondeterministicMerge:                     sdk.Bool(true),
					ErrorOnNondeterministicUpdate:                    sdk.Bool(true),
					EventTable:                                       sdk.Bool(true),
					ExternalOAuthAddPrivilegedRolesToBlockedList:     sdk.Bool(true),
					ExternalVolume:                                   sdk.Bool(true),
					GeographyOutputFormat:                            sdk.Bool(true),
					GeometryOutputFormat:                             sdk.Bool(true),
					HybridTableLockTimeout:                           sdk.Bool(true),
					InitialReplicationSizeLimitInTB:                  sdk.Bool(true),
					JdbcTreatDecimalAsInt:                            sdk.Bool(true),
					JdbcTreatTimestampNtzAsUtc:                       sdk.Bool(true),
					JdbcUseSessionTimezone:                           sdk.Bool(true),
					JsonIndent:                                       sdk.Bool(true),
					JsTreatIntegerAsBigInt:                           sdk.Bool(true),
					ListingAutoFulfillmentReplicationRefreshSchedule: sdk.Bool(true),
					LockTimeout:                                      sdk.Bool(true),
					LogLevel:                                         sdk.Bool(true),
					MaxConcurrencyLevel:                              sdk.Bool(true),
					MaxDataExtensionTimeInDays:                       sdk.Bool(true),
					MetricLevel:                                      sdk.Bool(true),
					MinDataRetentionTimeInDays:                       sdk.Bool(true),
					MultiStatementCount:                              sdk.Bool(true),
					NetworkPolicy:                                    sdk.Bool(true),
					NoorderSequenceAsDefault:                         sdk.Bool(true),
					OAuthAddPrivilegedRolesToBlockedList:             sdk.Bool(true),
					OdbcTreatDecimalAsInt:                            sdk.Bool(true),
					PeriodicDataRekeying:                             sdk.Bool(true),
					PipeExecutionPaused:                              sdk.Bool(true),
					PreventLoadFromInlineURL:                         sdk.Bool(true),
					PreventUnloadToInlineURL:                         sdk.Bool(true),
					PreventUnloadToInternalStages:                    sdk.Bool(true),
					PythonProfilerModules:                            sdk.Bool(true),
					PythonProfilerTargetStage:                        sdk.Bool(true),
					QueryTag:                                         sdk.Bool(true),
					QuotedIdentifiersIgnoreCase:                      sdk.Bool(true),
					ReplaceInvalidCharacters:                         sdk.Bool(true),
					RequireStorageIntegrationForStageCreation:        sdk.Bool(true),
					RequireStorageIntegrationForStageOperation:       sdk.Bool(true),
					RowsPerResultset:                                 sdk.Bool(true),
					S3StageVpceDnsName:                               sdk.Bool(true),
					SamlIdentityProvider:                             sdk.Bool(true),
					SearchPath:                                       sdk.Bool(true),
					ServerlessTaskMaxStatementSize:                   sdk.Bool(true),
					ServerlessTaskMinStatementSize:                   sdk.Bool(true),
					SimulatedDataSharingConsumer:                     sdk.Bool(true),
					SsoLoginPage:                                     sdk.Bool(true),
					StatementQueuedTimeoutInSeconds:                  sdk.Bool(true),
					StatementTimeoutInSeconds:                        sdk.Bool(true),
					StorageSerializationPolicy:                       sdk.Bool(true),
					StrictJsonOutput:                                 sdk.Bool(true),
					SuspendTaskAfterNumFailures:                      sdk.Bool(true),
					TaskAutoRetryAttempts:                            sdk.Bool(true),
					TimestampDayIsAlways24h:                          sdk.Bool(true),
					TimestampInputFormat:                             sdk.Bool(true),
					TimestampLtzOutputFormat:                         sdk.Bool(true),
					TimestampNtzOutputFormat:                         sdk.Bool(true),
					TimestampOutputFormat:                            sdk.Bool(true),
					TimestampTypeMapping:                             sdk.Bool(true),
					TimestampTzOutputFormat:                          sdk.Bool(true),
					Timezone:                                         sdk.Bool(true),
					TimeInputFormat:                                  sdk.Bool(true),
					TimeOutputFormat:                                 sdk.Bool(true),
					TraceLevel:                                       sdk.Bool(true),
					TransactionAbortOnError:                          sdk.Bool(true),
					TransactionDefaultIsolationLevel:                 sdk.Bool(true),
					TwoDigitCenturyStart:                             sdk.Bool(true),
					UnsupportedDdlAction:                             sdk.Bool(true),
					UserTaskManagedInitialWarehouseSize:              sdk.Bool(true),
					UserTaskMinimumTriggerIntervalInSeconds:          sdk.Bool(true),
					UserTaskTimeoutMs:                                sdk.Bool(true),
					UseCachedResult:                                  sdk.Bool(true),
					WeekOfYearPolicy:                                 sdk.Bool(true),
					WeekStart:                                        sdk.Bool(true),
				},
			},
		})
		require.NoError(t, err)

		parameters, err = client.Accounts.ShowParameters(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, parameters)

		for _, parameter := range sdk.AllAccountParameters {
			assertParameterIsDefault(t, parameters, string(parameter))
		}
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
