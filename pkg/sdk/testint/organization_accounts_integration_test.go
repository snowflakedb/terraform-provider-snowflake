//go:build !account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_OrganizationAccount_SelfAlter(t *testing.T) {
	testClientHelper().EnsureValidNonProdAccountIsUsed(t)

	client := testClient(t)
	ctx := testContext(t)

	err := client.Sessions.UseRole(ctx, snowflakeroles.GlobalOrgAdmin)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := client.Sessions.UseRole(ctx, snowflakeroles.Accountadmin)
		require.NoError(t, err)
	})
	t.Cleanup(testClientHelper().Role.UseRole(t, snowflakeroles.GlobalOrgAdmin))

	assertParameterValueSetOnAccount := func(t *testing.T, parameters []*sdk.Parameter, parameterKey string, parameterValue string) {
		t.Helper()
		param, err := collections.FindFirst(parameters, func(parameter *sdk.Parameter) bool { return parameter.Key == parameterKey })
		require.NoError(t, err)
		require.NotNil(t, param)
		require.Equal(t, parameterValue, (*param).Value)
		require.Equal(t, sdk.ParameterTypeAccount, (*param).Level)
	}

	assertParameterIsDefault := func(t *testing.T, parameters []*sdk.Parameter, parameterKey string) {
		t.Helper()
		param, err := collections.FindFirst(parameters, func(parameter *sdk.Parameter) bool { return parameter.Key == parameterKey })
		require.NoError(t, err, "parameter %v not found", parameterKey)
		require.NotNil(t, param)
		require.Equal(t, (*param).Default, (*param).Value)
		require.Equal(t, sdk.ParameterType(""), (*param).Level)
	}

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

		resourceMonitor2, resourceMonitor2Cleanup := testClientHelper().ResourceMonitor.CreateResourceMonitor(t)
		t.Cleanup(resourceMonitor2Cleanup)

		err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithResourceMonitor(resourceMonitor.ID())))
		require.NoError(t, err)

		// Set another resource monitor without unsetting the previous one
		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithResourceMonitor(resourceMonitor2.ID())))
		require.NoError(t, err)

		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithResourceMonitor(true)))
		require.NoError(t, err)

		// TODO(ticket number): Currently, there's no way to query resource monitor to verify to was unset properly.
	})

	t.Run("set / unset password policy", func(t *testing.T) {
		passwordPolicy, passwordPolicyCleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicyCleanup)

		passwordPolicy2, passwordPolicy2Cleanup := testClientHelper().PasswordPolicy.CreatePasswordPolicy(t)
		t.Cleanup(passwordPolicy2Cleanup)

		err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithPasswordPolicy(passwordPolicy.ID())))
		require.NoError(t, err)
		assertThatPolicyIsSetOnAccount(t, passwordPolicy.ID())

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
		assertThatPolicyIsSetOnAccount(t, sessionPolicy.ID())

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
		err := client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithSet(*sdk.NewOrganizationAccountSetRequest().WithParameters(sdk.AccountParameters{
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
		})))
		require.NoError(t, err)

		// TODO: Make for organization accounts, but passthrough the account interface
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

		err = client.OrganizationAccounts.Alter(ctx, sdk.NewAlterOrganizationAccountRequest().WithUnset(*sdk.NewOrganizationAccountUnsetRequest().WithParameters(sdk.AccountParametersUnset{
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
		})))
		require.NoError(t, err)

		parameters, err = client.Accounts.ShowParameters(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, parameters)

		for _, parameter := range sdk.AllAccountParameters {
			assertParameterIsDefault(t, parameters, string(parameter))
		}
	})
}
