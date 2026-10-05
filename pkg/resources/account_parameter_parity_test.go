package resources

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

var oldAccountParameterSupportedParameters = []sdk.AccountParameter{
	sdk.AccountParameterAllowClientMfaCaching,
	sdk.AccountParameterAllowIdToken,
	sdk.AccountParameterClientEncryptionKeySize,
	sdk.AccountParameterCortexCodeCliDailyEstCreditLimitPerUser,
	sdk.AccountParameterCortexCodeDesktopDailyEstCreditLimitPerUser,
	sdk.AccountParameterCortexCodeSnowsightDailyEstCreditLimitPerUser,
	sdk.AccountParameterCortexEnabledCrossRegion,
	sdk.AccountParameterCortexModelsAllowlist,
	sdk.AccountParameterDefaultStreamlitComputePool,
	sdk.AccountParameterDisableUserPrivilegeGrants,
	sdk.AccountParameterEnableIdentifierFirstLogin,
	sdk.AccountParameterEnableInternalStagesPrivatelink,
	sdk.AccountParameterEnablePerAccountAppServicePrivatelinkUrl,
	sdk.AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository,
	sdk.AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage,
	sdk.AccountParameterEnableUnhandledExceptionsReporting,
	sdk.AccountParameterEnforceNetworkRulesForInternalStages,
	sdk.AccountParameterEventTable,
	sdk.AccountParameterExternalOauthAddPrivilegedRolesToBlockedList,
	sdk.AccountParameterInitialReplicationSizeLimitInTb,
	sdk.AccountParameterMinDataRetentionTimeInDays,
	sdk.AccountParameterNetworkPolicy,
	sdk.AccountParameterOauthAddPrivilegedRolesToBlockedList,
	sdk.AccountParameterPeriodicDataRekeying,
	sdk.AccountParameterPreventLoadFromInlineURL,
	sdk.AccountParameterPreventUnloadToInlineUrl,
	sdk.AccountParameterRequireStorageIntegrationForStageCreation,
	sdk.AccountParameterRequireStorageIntegrationForStageOperation,
	sdk.AccountParameterSsoLoginPage,

	sdk.AccountParameterAbortDetachedQuery,
	sdk.AccountParameterActivePythonProfiler,
	sdk.AccountParameterAutocommit,
	sdk.AccountParameterBinaryInputFormat,
	sdk.AccountParameterBinaryOutputFormat,
	sdk.AccountParameterClientEnableLogInfoStatementParameters,
	sdk.AccountParameterClientMemoryLimit,
	sdk.AccountParameterClientMetadataRequestUseConnectionCtx,
	sdk.AccountParameterClientMetadataUseSessionDatabase,
	sdk.AccountParameterClientPrefetchThreads,
	sdk.AccountParameterClientResultChunkSize,
	sdk.AccountParameterClientSessionKeepAlive,
	sdk.AccountParameterClientSessionKeepAliveHeartbeatFrequency,
	sdk.AccountParameterClientTimestampTypeMapping,
	sdk.AccountParameterEnableUnloadPhysicalTypeOptimization,
	sdk.AccountParameterClientResultColumnCaseInsensitive,
	sdk.AccountParameterCsvTimestampFormat,
	sdk.AccountParameterDateInputFormat,
	sdk.AccountParameterDateOutputFormat,
	sdk.AccountParameterErrorOnNondeterministicMerge,
	sdk.AccountParameterErrorOnNondeterministicUpdate,
	sdk.AccountParameterGeographyOutputFormat,
	sdk.AccountParameterGeometryOutputFormat,
	sdk.AccountParameterHybridTableLockTimeout,
	sdk.AccountParameterJdbcTreatDecimalAsInt,
	sdk.AccountParameterJdbcTreatTimestampNtzAsUtc,
	sdk.AccountParameterJdbcUseSessionTimezone,
	sdk.AccountParameterJsonIndent,
	sdk.AccountParameterJsTreatIntegerAsBigint,
	sdk.AccountParameterLockTimeout,
	sdk.AccountParameterMultiStatementCount,
	sdk.AccountParameterNoorderSequenceAsDefault,
	sdk.AccountParameterOdbcTreatDecimalAsInt,
	sdk.AccountParameterPythonProfilerModules,
	sdk.AccountParameterPythonProfilerTargetStage,
	sdk.AccountParameterQueryTag,
	sdk.AccountParameterQuotedIdentifiersIgnoreCase,
	sdk.AccountParameterRowsPerResultset,
	sdk.AccountParameterS3StageVpceDnsName,
	sdk.AccountParameterSearchPath,
	sdk.AccountParameterSimulatedDataSharingConsumer,
	sdk.AccountParameterStatementTimeoutInSeconds,
	sdk.AccountParameterStrictJsonOutput,
	sdk.AccountParameterTimeInputFormat,
	sdk.AccountParameterTimeOutputFormat,
	sdk.AccountParameterTimestampDayIsAlways24h,
	sdk.AccountParameterTimestampInputFormat,
	sdk.AccountParameterTimestampLtzOutputFormat,
	sdk.AccountParameterTimestampNtzOutputFormat,
	sdk.AccountParameterTimestampOutputFormat,
	sdk.AccountParameterTimestampTypeMapping,
	sdk.AccountParameterTimestampTzOutputFormat,
	sdk.AccountParameterTimezone,
	sdk.AccountParameterTransactionAbortOnError,
	sdk.AccountParameterTransactionDefaultIsolationLevel,
	sdk.AccountParameterTwoDigitCenturyStart,
	sdk.AccountParameterUnsupportedDdlAction,
	sdk.AccountParameterUseCachedResult,
	sdk.AccountParameterWeekOfYearPolicy,
	sdk.AccountParameterWeekStart,

	sdk.AccountParameterCatalog,
	sdk.AccountParameterDataRetentionTimeInDays,
	sdk.AccountParameterDefaultDdlCollation,
	sdk.AccountParameterDefaultNotebookComputePoolCpu,
	sdk.AccountParameterDefaultNotebookComputePoolGpu,
	sdk.AccountParameterExternalVolume,
	sdk.AccountParameterLogLevel,
	sdk.AccountParameterLogEventLevel,
	sdk.AccountParameterMaxConcurrencyLevel,
	sdk.AccountParameterMaxDataExtensionTimeInDays,
	sdk.AccountParameterPipeExecutionPaused,
	sdk.AccountParameterPreventUnloadToInternalStages,
	sdk.AccountParameterReplaceInvalidCharacters,
	sdk.AccountParameterStatementQueuedTimeoutInSeconds,
	sdk.AccountParameterStorageSerializationPolicy,
	sdk.AccountParameterShareRestrictions,
	sdk.AccountParameterSuspendTaskAfterNumFailures,
	sdk.AccountParameterTraceLevel,
	sdk.AccountParameterUserTaskManagedInitialWarehouseSize,
	sdk.AccountParameterUserTaskTimeoutMs,
	sdk.AccountParameterTaskAutoRetryAttempts,
	sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds,
	sdk.AccountParameterMetricLevel,
	sdk.AccountParameterEnableConsoleOutput,
	sdk.AccountParameterEnableUnredactedQuerySyntaxError,
	sdk.AccountParameterEnablePersonalDatabase,

	sdk.AccountParameterAllowBindValuesAccess,
	sdk.AccountParameterAllowedSpcsWorkloadTypes,
	sdk.AccountParameterDataMetricSchedule,
	sdk.AccountParameterDefaultDbtVersion,
	sdk.AccountParameterDisallowedSpcsWorkloadTypes,
	sdk.AccountParameterEnableBudgetEventLogging,
	sdk.AccountParameterEnableCortexAnalyst,
	sdk.AccountParameterEnableDataCompaction,
	sdk.AccountParameterEnableGetDdlUseDataTypeAlias,
	sdk.AccountParameterEnableIcebergMergeOnRead,
	sdk.AccountParameterEnableNotebookCreationInPersonalDb,
	sdk.AccountParameterEnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement,
	sdk.AccountParameterEnableTagPropagationEventLogging,
	sdk.AccountParameterIcebergVersionDefault,
	sdk.AccountParameterReadConsistencyMode,
	sdk.AccountParameterRowTimestampDefault,
	sdk.AccountParameterSqlTraceQueryText,
	sdk.AccountParameterUseWorkspacesForSql,
}

// GH #4764 / SNOW-3562870 tracks AccountParameterDisableUiDownloadButton specifically; the rest were
// also missing from the old allowlist but untracked.
var parametersAddedByCatalogMigration = []sdk.AccountParameter{
	sdk.AccountParameterDisableUiDownloadButton,
	sdk.AccountParameterBaseLocationPrefix,
	sdk.AccountParameterCatalogSync,
	sdk.AccountParameterDefaultNullOrdering,
	sdk.AccountParameterDefaultStreamlitNotebookWarehouse,
	sdk.AccountParameterEnableAutomaticSensitiveDataClassificationLog,
	sdk.AccountParameterEnableEgressCostOptimizer,
	sdk.AccountParameterEnableUnredactedSecureObjectError,
	sdk.AccountParameterListingAutoFulfillmentReplicationRefreshSchedule,
	sdk.AccountParameterServerlessTaskMaxStatementSize,
	sdk.AccountParameterServerlessTaskMinStatementSize,
}

func TestAccountParameterSupportedParameters_CatalogParity(t *testing.T) {
	expected := make(map[sdk.AccountParameter]struct{}, len(oldAccountParameterSupportedParameters)+len(parametersAddedByCatalogMigration))
	for _, p := range oldAccountParameterSupportedParameters {
		expected[p] = struct{}{}
	}
	for _, p := range parametersAddedByCatalogMigration {
		expected[p] = struct{}{}
	}

	actual := make(map[sdk.AccountParameter]struct{}, len(accountParameterSupportedParameters))
	for _, p := range accountParameterSupportedParameters {
		actual[p] = struct{}{}
	}

	require.Equal(t, expected, actual)
}
