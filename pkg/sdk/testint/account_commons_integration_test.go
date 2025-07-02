package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// SetAndUnsetAccountParametersTest is a common test used for different account kinds.
func SetAndUnsetAccountParametersTest(
	setParameters func(sdk.AccountParameters) error,
	unsetAllParameters func(ctx context.Context) error,
	showParameters func(ctx context.Context) ([]*sdk.Parameter, error),
) func(t *testing.T) {
	return func(t *testing.T) {
		id := testClientHelper().Context.CurrentAccountId(t)

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
		err := setParameters(sdk.AccountParameters{
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
		})
		require.NoError(t, err)

		parameters, err := showParameters(context.Background())
		require.NoError(t, err)

		objectparametersassert.AccountParametersPrefetched(t, id, parameters).
			HasAbortDetachedQuery(true).
			HasAllowClientMfaCaching(true).
			HasAllowIdToken(true).
			HasAutocommit(false).
			HasBaseLocationPrefix("STORAGE_BASE_URL/").
			HasBinaryInputFormat(sdk.BinaryInputFormatBase64).
			HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
			HasCatalog(helpers.TestDatabaseCatalog.Name()).
			HasClientEnableLogInfoStatementParameters(true).
			HasClientEncryptionKeySize(256).
			HasClientMemoryLimit(1540).
			HasClientMetadataRequestUseConnectionCtx(true).
			HasClientMetadataUseSessionDatabase(true).
			HasClientPrefetchThreads(5).
			HasClientResultChunkSize(159).
			HasClientResultColumnCaseInsensitive(true).
			HasClientSessionKeepAlive(true).
			HasClientSessionKeepAliveHeartbeatFrequency(3599).
			HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
			HasCortexEnabledCrossRegion("ANY_REGION").
			HasCortexModelsAllowlist("All").
			HasCsvTimestampFormat("YYYY-MM-DD").
			HasDataRetentionTimeInDays(2).
			HasDateInputFormat("YYYY-MM-DD").
			HasDateOutputFormat("YYYY-MM-DD").
			HasDefaultDdlCollation("en-cs").
			HasDefaultNotebookComputePoolCpu("CPU_X64_S").
			HasDefaultNotebookComputePoolGpu("GPU_NV_S").
			HasDefaultNullOrdering(sdk.DefaultNullOrderingFirst).
			HasDefaultStreamlitNotebookWarehouse(warehouseId.Name()).
			HasDisableUiDownloadButton(true).
			HasDisableUserPrivilegeGrants(true).
			HasEnableAutomaticSensitiveDataClassificationLog(false).
			HasEnableEgressCostOptimizer(false).
			HasEnableIdentifierFirstLogin(false).
			HasEnableTriSecretAndRekeyOptOutForImageRepository(true).
			HasEnableTriSecretAndRekeyOptOutForSpcsBlockStorage(true).
			HasEnableUnhandledExceptionsReporting(false).
			HasEnableUnloadPhysicalTypeOptimization(false).
			HasEnableUnredactedQuerySyntaxError(true).
			HasEnableUnredactedSecureObjectError(true).
			HasEnforceNetworkRulesForInternalStages(true).
			HasErrorOnNondeterministicMerge(false).
			HasErrorOnNondeterministicUpdate(true).
			HasEventTable(eventTable.ID().FullyQualifiedName()).
			HasExternalOauthAddPrivilegedRolesToBlockedList(false).
			HasExternalVolume(externalVolumeId.Name()).
			HasGeographyOutputFormat(sdk.GeographyOutputFormatWKT).
			HasGeometryOutputFormat(sdk.GeometryOutputFormatWKT).
			HasHybridTableLockTimeout(3599).
			HasInitialReplicationSizeLimitInTb("9.9").
			HasJdbcTreatDecimalAsInt(false).
			HasJdbcTreatTimestampNtzAsUtc(true).
			HasJdbcUseSessionTimezone(false).
			HasJsonIndent(4).
			HasJsTreatIntegerAsBigint(true).
			HasListingAutoFulfillmentReplicationRefreshSchedule("2 minutes").
			HasLockTimeout(43201).
			HasLogLevel(sdk.LogLevelInfo).
			HasMaxConcurrencyLevel(7).
			HasMaxDataExtensionTimeInDays(13).
			HasMetricLevel(sdk.MetricLevelAll).
			HasMinDataRetentionTimeInDays(1).
			HasMultiStatementCount(0).
			HasNetworkPolicy(networkPolicy.ID().Name()).
			HasNoorderSequenceAsDefault(false).
			HasOauthAddPrivilegedRolesToBlockedList(false).
			HasOdbcTreatDecimalAsInt(true).
			HasPeriodicDataRekeying(false).
			HasPipeExecutionPaused(true).
			HasPreventUnloadToInlineUrl(true).
			HasPreventUnloadToInternalStages(true).
			HasQueryTag("test-query-tag").
			HasQuotedIdentifiersIgnoreCase(true).
			HasReplaceInvalidCharacters(true).
			HasRequireStorageIntegrationForStageCreation(true).
			HasRequireStorageIntegrationForStageOperation(true).
			HasRowsPerResultset(1000).
			HasSearchPath("$current, $public").
			HasServerlessTaskMaxStatementSize("6X-LARGE").
			HasServerlessTaskMinStatementSize(sdk.WarehouseSizeSmall).
			HasSsoLoginPage(true).
			HasStatementQueuedTimeoutInSeconds(1).
			HasStatementTimeoutInSeconds(1).
			HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
			HasStrictJsonOutput(true).
			HasSuspendTaskAfterNumFailures(3).
			HasTaskAutoRetryAttempts(3).
			HasTimestampDayIsAlways24h(true).
			HasTimestampInputFormat("YYYY-MM-DD").
			HasTimestampLtzOutputFormat("YYYY-MM-DD").
			HasTimestampNtzOutputFormat("YYYY-MM-DD").
			HasTimestampOutputFormat("YYYY-MM-DD").
			HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
			HasTimestampTzOutputFormat("YYYY-MM-DD").
			HasTimezone("Europe/London").
			HasTimeInputFormat("YYYY-MM-DD").
			HasTimeOutputFormat("YYYY-MM-DD").
			HasTraceLevel(sdk.TraceLevelPropagate).
			HasTransactionAbortOnError(true).
			HasTransactionDefaultIsolationLevel(string(sdk.TransactionDefaultIsolationLevelReadCommitted)).
			HasTwoDigitCenturyStart(1971).
			HasUnsupportedDdlAction(string(sdk.UnsupportedDDLActionFail)).
			HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeX6Large).
			HasUserTaskMinimumTriggerIntervalInSeconds(10).
			HasUserTaskTimeoutMs(10).
			HasUseCachedResult(false).
			HasWeekOfYearPolicy(1).
			HasWeekStart(1)

		err = unsetAllParameters(context.Background())
		require.NoError(t, err)

		parameters, err = showParameters(context.Background())
		require.NoError(t, err)

		objectparametersassert.AccountParametersPrefetched(t, id, parameters).HasAllDefaults()
	}
}

func assertThatPolicyIsSetOnAccount(t *testing.T, kind sdk.PolicyKind, id sdk.SchemaObjectIdentifier) {
	t.Helper()

	policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(testClient(t).GetAccountLocator()), sdk.PolicyEntityDomainAccount)
	require.NoError(t, err)
	_, err = collections.FindFirst(policies, func(reference sdk.PolicyReference) bool {
		return reference.PolicyName == id.Name() && reference.PolicyKind == kind
	})
	require.NoError(t, err)
}

func assertThatNoPolicyIsSetOnAccount(t *testing.T) {
	t.Helper()

	policies, err := testClientHelper().PolicyReferences.GetPolicyReferences(t, sdk.NewAccountObjectIdentifier(testClient(t).GetAccountLocator()), sdk.PolicyEntityDomainAccount)
	require.Empty(t, policies)
	require.NoError(t, err)
}
