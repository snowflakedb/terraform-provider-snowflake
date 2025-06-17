//go:build !account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CurrentAccount_Parameters(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	warehouseId := testClient().Ids.WarehouseId()

	eventTable, eventTableCleanup := testClient().EventTable.Create(t)
	t.Cleanup(eventTableCleanup)

	externalVolumeId, externalVolumeCleanup := testClient().ExternalVolume.Create(t)
	t.Cleanup(externalVolumeCleanup)

	createNetworkPolicyRequest := sdk.NewCreateNetworkPolicyRequest(testClient().Ids.RandomAccountObjectIdentifier()).WithAllowedIpList([]sdk.IPRequest{*sdk.NewIPRequest("0.0.0.0/0")})
	networkPolicy, networkPolicyCleanup := testClient().NetworkPolicy.CreateNetworkPolicyWithRequest(t, createNetworkPolicyRequest)
	t.Cleanup(networkPolicyCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	unsetParametersModel := model.CurrentAccount("test")

	setParametersModel := model.CurrentAccount("test").
		WithAbortDetachedQuery(true).
		WithAllowClientMfaCaching(true).
		WithAllowIdToken(true).
		WithAutocommit(false).
		WithBaseLocationPrefix("STORAGE_BASE_URL/").
		WithBinaryInputFormat(string(sdk.BinaryInputFormatBase64)).
		WithBinaryOutputFormat(string(sdk.BinaryOutputFormatBase64)).
		WithCatalog(helpers.TestDatabaseCatalog.Name()).
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
		WithClientTimestampTypeMapping(string(sdk.ClientTimestampTypeMappingNtz)).
		WithCortexEnabledCrossRegion("ANY_REGION").
		WithCortexModelsAllowlist("All").
		WithCsvTimestampFormat("YYYY-MM-DD").
		WithDataRetentionTimeInDays(2).
		WithDateInputFormat("YYYY-MM-DD").
		WithDateOutputFormat("YYYY-MM-DD").
		WithDefaultDdlCollation("en-cs").
		WithDefaultNotebookComputePoolCpu("CPU_X64_S").
		WithDefaultNotebookComputePoolGpu("GPU_NV_S").
		WithDefaultNullOrdering(string(sdk.DefaultNullOrderingFirst)).
		WithDefaultStreamlitNotebookWarehouse(warehouseId.Name()).
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
		WithEventTable(eventTable.ID().FullyQualifiedName()).
		WithExternalOauthAddPrivilegedRolesToBlockedList(false).
		WithExternalVolume(externalVolumeId.Name()).
		WithGeographyOutputFormat(string(sdk.GeographyOutputFormatWKT)).
		WithGeometryOutputFormat(string(sdk.GeometryOutputFormatWKT)).
		WithHybridTableLockTimeout(3599).
		WithInitialReplicationSizeLimitInTb("9.9").
		WithJdbcTreatDecimalAsInt(false).
		WithJdbcTreatTimestampNtzAsUtc(true).
		WithJdbcUseSessionTimezone(false).
		WithJsonIndent(4).
		WithJsTreatIntegerAsBigint(true).
		WithListingAutoFulfillmentReplicationRefreshSchedule("2 minutes").
		WithLockTimeout(43201).
		WithLogLevel(string(sdk.LogLevelInfo)).
		WithMaxConcurrencyLevel(7).
		WithMaxDataExtensionTimeInDays(13).
		WithMetricLevel(string(sdk.MetricLevelAll)).
		WithMinDataRetentionTimeInDays(1).
		WithMultiStatementCount(0).
		WithNetworkPolicy(networkPolicy.ID().Name()).
		WithNoorderSequenceAsDefault(false).
		WithOauthAddPrivilegedRolesToBlockedList(false).
		WithOdbcTreatDecimalAsInt(true).
		WithPeriodicDataRekeying(false).
		WithPipeExecutionPaused(true).
		WithPreventUnloadToInlineUrl(true).
		WithPreventUnloadToInternalStages(true).
		WithPythonProfilerTargetStage(stage.ID().FullyQualifiedName()).
		WithQueryTag("test-query-tag").
		WithQuotedIdentifiersIgnoreCase(true).
		WithReplaceInvalidCharacters(true).
		WithRequireStorageIntegrationForStageCreation(true).
		WithRequireStorageIntegrationForStageOperation(true).
		WithRowsPerResultset(1000).
		WithSearchPath("$current, $public").
		WithServerlessTaskMaxStatementSize(string(sdk.WarehouseSizeXLarge)).
		WithServerlessTaskMinStatementSize(string(sdk.WarehouseSizeSmall)).
		WithSsoLoginPage(true).
		WithStatementQueuedTimeoutInSeconds(1).
		WithStatementTimeoutInSeconds(1).
		WithStorageSerializationPolicy(string(sdk.StorageSerializationPolicyOptimized)).
		WithStrictJsonOutput(true).
		WithSuspendTaskAfterNumFailures(3).
		WithTaskAutoRetryAttempts(3).
		WithTimestampDayIsAlways24h(true).
		WithTimestampInputFormat("YYYY-MM-DD").
		WithTimestampLtzOutputFormat("YYYY-MM-DD").
		WithTimestampNtzOutputFormat("YYYY-MM-DD").
		WithTimestampOutputFormat("YYYY-MM-DD").
		WithTimestampTypeMapping(string(sdk.TimestampTypeMappingLtz)).
		WithTimestampTzOutputFormat("YYYY-MM-DD").
		WithTimezone("Europe/London").
		WithTimeInputFormat("YYYY-MM-DD").
		WithTimeOutputFormat("YYYY-MM-DD").
		WithTraceLevel(string(sdk.TraceLevelPropagate)).
		WithTransactionAbortOnError(true).
		WithTransactionDefaultIsolationLevel(string(sdk.TransactionDefaultIsolationLevelReadCommitted)).
		WithTwoDigitCenturyStart(1971).
		WithUnsupportedDdlAction(string(sdk.UnsupportedDDLActionFail)).
		WithUserTaskManagedInitialWarehouseSize(string(sdk.WarehouseSizeSmall)).
		WithUserTaskMinimumTriggerIntervalInSeconds(10).
		WithUserTaskTimeoutMs(10).
		WithUseCachedResult(false).
		WithWeekOfYearPolicy(1).
		WithWeekStart(1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// resource with unset parameters
			{
				Config: config.FromModels(t, unsetParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentAccountResource(t, unsetParametersModel.ResourceReference()).
						HasAbortDetachedQueryString("false").
						HasAllowClientMfaCachingString("false").
						HasAllowIdTokenString("false").
						HasAutocommitString("true").
						HasBaseLocationPrefixEmpty().
						HasBinaryInputFormatString("HEX").
						HasBinaryOutputFormatString("HEX").
						HasCatalogEmpty().
						HasCatalogSyncEmpty().
						HasClientEnableLogInfoStatementParametersString("false").
						HasClientEncryptionKeySizeString("128").
						HasClientMemoryLimitString("1536").
						HasClientMetadataRequestUseConnectionCtxString("false").
						HasClientMetadataUseSessionDatabaseString("false").
						HasClientPrefetchThreadsString("4").
						HasClientResultChunkSizeString("160").
						HasClientResultColumnCaseInsensitiveString("false").
						HasClientSessionKeepAliveString("false").
						HasClientSessionKeepAliveHeartbeatFrequencyString("3600").
						HasClientTimestampTypeMappingString("TIMESTAMP_LTZ").
						HasCortexEnabledCrossRegionString("DISABLED").
						HasCortexModelsAllowlistString("ALL").
						HasCsvTimestampFormatEmpty().
						HasDataRetentionTimeInDaysString("1").
						HasDateInputFormatString("AUTO").
						HasDateOutputFormatString("YYYY-MM-DD").
						HasDefaultDdlCollationEmpty().
						HasDefaultNotebookComputePoolCpuString("SYSTEM_COMPUTE_POOL_CPU").
						HasDefaultNotebookComputePoolGpuString("SYSTEM_COMPUTE_POOL_GPU").
						HasDefaultNullOrderingString("LAST").
						HasDefaultStreamlitNotebookWarehouseString("REGRESS").
						HasDisableUiDownloadButtonString("false").
						HasDisableUserPrivilegeGrantsString("false").
						HasEnableAutomaticSensitiveDataClassificationLogString("true").
						HasEnableEgressCostOptimizerString("true").
						HasEnableIdentifierFirstLoginString("true").
						HasEnableTriSecretAndRekeyOptOutForImageRepositoryString("false").
						HasEnableTriSecretAndRekeyOptOutForSpcsBlockStorageString("false").
						HasEnableUnhandledExceptionsReportingString("true").
						HasEnableUnloadPhysicalTypeOptimizationString("true").
						HasEnableUnredactedQuerySyntaxErrorString("false").
						HasEnableUnredactedSecureObjectErrorString("false").
						HasEnforceNetworkRulesForInternalStagesString("false").
						HasErrorOnNondeterministicMergeString("true").
						HasErrorOnNondeterministicUpdateString("false").
						HasEventTableString("snowflake.telemetry.events").
						HasExternalOauthAddPrivilegedRolesToBlockedListString("true").
						HasExternalVolumeEmpty().
						HasGeographyOutputFormatString("GeoJSON").
						HasGeometryOutputFormatString("GeoJSON").
						HasHybridTableLockTimeoutString("3600").
						HasInitialReplicationSizeLimitInTbString("10.0").
						HasJdbcTreatDecimalAsIntString("true").
						HasJdbcTreatTimestampNtzAsUtcString("false").
						HasJdbcUseSessionTimezoneString("true").
						HasJsonIndentString("2").
						HasJsTreatIntegerAsBigintString("false").
						HasListingAutoFulfillmentReplicationRefreshScheduleString("1440 MINUTE").
						HasLockTimeoutString("43200").
						HasLogLevelString("OFF").
						HasMaxConcurrencyLevelString("8").
						HasMaxDataExtensionTimeInDaysString("14").
						HasMetricLevelString("NONE").
						HasMinDataRetentionTimeInDaysString("0").
						HasMultiStatementCountString("1").
						HasNetworkPolicyEmpty().
						HasNoorderSequenceAsDefaultString("true").
						HasOauthAddPrivilegedRolesToBlockedListString("true").
						HasOdbcTreatDecimalAsIntString("false").
						HasPeriodicDataRekeyingString("false").
						HasPipeExecutionPausedString("false").
						HasPreventUnloadToInlineUrlString("false").
						HasPreventUnloadToInternalStagesString("false").
						HasPythonProfilerTargetStageEmpty().
						HasQueryTagEmpty().
						HasQuotedIdentifiersIgnoreCaseString("false").
						HasReplaceInvalidCharactersString("false").
						HasRequireStorageIntegrationForStageCreationString("false").
						HasRequireStorageIntegrationForStageOperationString("false").
						HasRowsPerResultsetString("0").
						HasSearchPathString("$current, $public").
						HasServerlessTaskMaxStatementSizeString("X2Large").
						HasServerlessTaskMinStatementSizeString("XSMALL").
						HasSsoLoginPageString("false").
						HasStatementQueuedTimeoutInSecondsString("0").
						HasStatementTimeoutInSecondsString("172800").
						HasStorageSerializationPolicyString("OPTIMIZED").
						HasStrictJsonOutputString("false").
						HasSuspendTaskAfterNumFailuresString("10").
						HasTaskAutoRetryAttemptsString("0").
						HasTimestampDayIsAlways24hString("false").
						HasTimestampInputFormatString("AUTO").
						HasTimestampLtzOutputFormatEmpty().
						HasTimestampNtzOutputFormatString("YYYY-MM-DD HH24:MI:SS.FF3").
						HasTimestampOutputFormatString("YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM").
						HasTimestampTypeMappingString("TIMESTAMP_NTZ").
						HasTimestampTzOutputFormatEmpty().
						HasTimezoneString("America/Los_Angeles").
						HasTimeInputFormatString("AUTO").
						HasTimeOutputFormatString("HH24:MI:SS").
						HasTraceLevelString("OFF").
						HasTransactionAbortOnErrorString("false").
						HasTransactionDefaultIsolationLevelString("READ COMMITTED").
						HasTwoDigitCenturyStartString("1970").
						HasUnsupportedDdlActionString("ignore").
						HasUserTaskManagedInitialWarehouseSizeString("Medium").
						HasUserTaskMinimumTriggerIntervalInSecondsString("30").
						HasUserTaskTimeoutMsString("3600000").
						HasUseCachedResultString("true").
						HasWeekOfYearPolicyString("0").
						HasWeekStartString("0"),
				),
			},
			// import with unset parameters
			{
				Config:       config.FromModels(t, unsetParametersModel),
				ResourceName: unsetParametersModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentAccountResource(t, "current_account").
						HasAbortDetachedQueryString("false").
						HasAllowClientMfaCachingString("false").
						HasAllowIdTokenString("false").
						HasAutocommitString("true").
						HasBaseLocationPrefixEmpty().
						HasBinaryInputFormatString("HEX").
						HasBinaryOutputFormatString("HEX").
						HasCatalogEmpty().
						HasCatalogSyncEmpty().
						HasClientEnableLogInfoStatementParametersString("false").
						HasClientEncryptionKeySizeString("128").
						HasClientMemoryLimitString("1536").
						HasClientMetadataRequestUseConnectionCtxString("false").
						HasClientMetadataUseSessionDatabaseString("false").
						HasClientPrefetchThreadsString("4").
						HasClientResultChunkSizeString("160").
						HasClientResultColumnCaseInsensitiveString("false").
						HasClientSessionKeepAliveString("false").
						HasClientSessionKeepAliveHeartbeatFrequencyString("3600").
						HasClientTimestampTypeMappingString("TIMESTAMP_LTZ").
						HasCortexEnabledCrossRegionString("DISABLED").
						HasCortexModelsAllowlistString("ALL").
						HasCsvTimestampFormatEmpty().
						HasDataRetentionTimeInDaysString("1").
						HasDateInputFormatString("AUTO").
						HasDateOutputFormatString("YYYY-MM-DD").
						HasDefaultDdlCollationEmpty().
						HasDefaultNotebookComputePoolCpuString("SYSTEM_COMPUTE_POOL_CPU").
						HasDefaultNotebookComputePoolGpuString("SYSTEM_COMPUTE_POOL_GPU").
						HasDefaultNullOrderingString("LAST").
						HasDefaultStreamlitNotebookWarehouseString("REGRESS").
						HasDisableUiDownloadButtonString("false").
						HasDisableUserPrivilegeGrantsString("false").
						HasEnableAutomaticSensitiveDataClassificationLogString("true").
						HasEnableEgressCostOptimizerString("true").
						HasEnableIdentifierFirstLoginString("true").
						HasEnableTriSecretAndRekeyOptOutForImageRepositoryString("false").
						HasEnableTriSecretAndRekeyOptOutForSpcsBlockStorageString("false").
						HasEnableUnhandledExceptionsReportingString("true").
						HasEnableUnloadPhysicalTypeOptimizationString("true").
						HasEnableUnredactedQuerySyntaxErrorString("false").
						HasEnableUnredactedSecureObjectErrorString("false").
						HasEnforceNetworkRulesForInternalStagesString("false").
						HasErrorOnNondeterministicMergeString("true").
						HasErrorOnNondeterministicUpdateString("false").
						HasEventTableString("snowflake.telemetry.events").
						HasExternalOauthAddPrivilegedRolesToBlockedListString("true").
						HasExternalVolumeEmpty().
						HasGeographyOutputFormatString("GeoJSON").
						HasGeometryOutputFormatString("GeoJSON").
						HasHybridTableLockTimeoutString("3600").
						HasInitialReplicationSizeLimitInTbString("10.0").
						HasJdbcTreatDecimalAsIntString("true").
						HasJdbcTreatTimestampNtzAsUtcString("false").
						HasJdbcUseSessionTimezoneString("true").
						HasJsonIndentString("2").
						HasJsTreatIntegerAsBigintString("false").
						HasListingAutoFulfillmentReplicationRefreshScheduleString("1440 MINUTE").
						HasLockTimeoutString("43200").
						HasLogLevelString("OFF").
						HasMaxConcurrencyLevelString("8").
						HasMaxDataExtensionTimeInDaysString("14").
						HasMetricLevelString("NONE").
						HasMinDataRetentionTimeInDaysString("0").
						HasMultiStatementCountString("1").
						HasNetworkPolicyEmpty().
						HasNoorderSequenceAsDefaultString("true").
						HasOauthAddPrivilegedRolesToBlockedListString("true").
						HasOdbcTreatDecimalAsIntString("false").
						HasPeriodicDataRekeyingString("false").
						HasPipeExecutionPausedString("false").
						HasPreventUnloadToInlineUrlString("false").
						HasPreventUnloadToInternalStagesString("false").
						HasPythonProfilerTargetStageEmpty().
						HasQueryTagEmpty().
						HasQuotedIdentifiersIgnoreCaseString("false").
						HasReplaceInvalidCharactersString("false").
						HasRequireStorageIntegrationForStageCreationString("false").
						HasRequireStorageIntegrationForStageOperationString("false").
						HasRowsPerResultsetString("0").
						HasSearchPathString("$current, $public").
						HasServerlessTaskMaxStatementSizeString("X2Large").
						HasServerlessTaskMinStatementSizeString("XSMALL").
						HasSsoLoginPageString("false").
						HasStatementQueuedTimeoutInSecondsString("0").
						HasStatementTimeoutInSecondsString("172800").
						HasStorageSerializationPolicyString("OPTIMIZED").
						HasStrictJsonOutputString("false").
						HasSuspendTaskAfterNumFailuresString("10").
						HasTaskAutoRetryAttemptsString("0").
						HasTimestampDayIsAlways24hString("false").
						HasTimestampInputFormatString("AUTO").
						HasTimestampLtzOutputFormatEmpty().
						HasTimestampNtzOutputFormatString("YYYY-MM-DD HH24:MI:SS.FF3").
						HasTimestampOutputFormatString("YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM").
						HasTimestampTypeMappingString("TIMESTAMP_NTZ").
						HasTimestampTzOutputFormatEmpty().
						HasTimezoneString("America/Los_Angeles").
						HasTimeInputFormatString("AUTO").
						HasTimeOutputFormatString("HH24:MI:SS").
						HasTraceLevelString("OFF").
						HasTransactionAbortOnErrorString("false").
						HasTransactionDefaultIsolationLevelString("READ COMMITTED").
						HasTwoDigitCenturyStartString("1970").
						HasUnsupportedDdlActionString("ignore").
						HasUserTaskManagedInitialWarehouseSizeString("Medium").
						HasUserTaskMinimumTriggerIntervalInSecondsString("30").
						HasUserTaskTimeoutMsString("3600000").
						HasUseCachedResultString("true").
						HasWeekOfYearPolicyString("0").
						HasWeekStartString("0"),
				),
			},
			// set all parameters
			{
				Config: config.FromModels(t, setParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentAccountResource(t, setParametersModel.ResourceReference()).
						HasAbortDetachedQueryString("true").
						HasAllowClientMfaCachingString("true").
						HasAllowIdTokenString("true").
						HasAutocommitString("false").
						HasBaseLocationPrefixString("STORAGE_BASE_URL/").
						HasBinaryInputFormatString(string(sdk.BinaryInputFormatBase64)).
						HasBinaryOutputFormatString(string(sdk.BinaryOutputFormatBase64)).
						HasCatalogString(helpers.TestDatabaseCatalog.Name()).
						HasClientEnableLogInfoStatementParametersString("true").
						HasClientEncryptionKeySizeString("256").
						HasClientMemoryLimitString("1540").
						HasClientMetadataRequestUseConnectionCtxString("true").
						HasClientMetadataUseSessionDatabaseString("true").
						HasClientPrefetchThreadsString("5").
						HasClientResultChunkSizeString("159").
						HasClientResultColumnCaseInsensitiveString("true").
						HasClientSessionKeepAliveString("true").
						HasClientSessionKeepAliveHeartbeatFrequencyString("3599").
						HasClientTimestampTypeMappingString(string(sdk.ClientTimestampTypeMappingNtz)).
						HasCortexEnabledCrossRegionString("ANY_REGION").
						HasCortexModelsAllowlistString("All").
						HasCsvTimestampFormatString("YYYY-MM-DD").
						HasDataRetentionTimeInDaysString("2").
						HasDateInputFormatString("YYYY-MM-DD").
						HasDateOutputFormatString("YYYY-MM-DD").
						HasDefaultDdlCollationString("en-cs").
						HasDefaultNotebookComputePoolCpuString("CPU_X64_S").
						HasDefaultNotebookComputePoolGpuString("GPU_NV_S").
						HasDefaultNullOrderingString(string(sdk.DefaultNullOrderingFirst)).
						HasDefaultStreamlitNotebookWarehouseString(warehouseId.Name()).
						HasDisableUiDownloadButtonString("true").
						HasDisableUserPrivilegeGrantsString("true").
						HasEnableAutomaticSensitiveDataClassificationLogString("false").
						HasEnableEgressCostOptimizerString("false").
						HasEnableIdentifierFirstLoginString("false").
						HasEnableTriSecretAndRekeyOptOutForImageRepositoryString("true").
						HasEnableTriSecretAndRekeyOptOutForSpcsBlockStorageString("true").
						HasEnableUnhandledExceptionsReportingString("false").
						HasEnableUnloadPhysicalTypeOptimizationString("false").
						HasEnableUnredactedQuerySyntaxErrorString("true").
						HasEnableUnredactedSecureObjectErrorString("true").
						HasEnforceNetworkRulesForInternalStagesString("true").
						HasErrorOnNondeterministicMergeString("false").
						HasErrorOnNondeterministicUpdateString("true").
						HasEventTableString(eventTable.ID().FullyQualifiedName()).
						HasExternalOauthAddPrivilegedRolesToBlockedListString("false").
						HasExternalVolumeString(externalVolumeId.Name()).
						HasGeographyOutputFormatString(string(sdk.GeographyOutputFormatWKT)).
						HasGeometryOutputFormatString(string(sdk.GeometryOutputFormatWKT)).
						HasHybridTableLockTimeoutString("3599").
						HasInitialReplicationSizeLimitInTbString("9.9").
						HasJdbcTreatDecimalAsIntString("false").
						HasJdbcTreatTimestampNtzAsUtcString("true").
						HasJdbcUseSessionTimezoneString("false").
						HasJsonIndentString("4").
						HasJsTreatIntegerAsBigintString("true").
						HasListingAutoFulfillmentReplicationRefreshScheduleString("2 minutes").
						HasLockTimeoutString("43201").
						HasLogLevelString(string(sdk.LogLevelInfo)).
						HasMaxConcurrencyLevelString("7").
						HasMaxDataExtensionTimeInDaysString("13").
						HasMetricLevelString(string(sdk.MetricLevelAll)).
						HasMinDataRetentionTimeInDaysString("1").
						HasMultiStatementCountString("0").
						HasNetworkPolicyString(networkPolicy.ID().Name()).
						HasNoorderSequenceAsDefaultString("false").
						HasOauthAddPrivilegedRolesToBlockedListString("false").
						HasOdbcTreatDecimalAsIntString("true").
						HasPeriodicDataRekeyingString("false").
						HasPipeExecutionPausedString("true").
						HasPreventUnloadToInlineUrlString("true").
						HasPreventUnloadToInternalStagesString("true").
						HasPythonProfilerTargetStageString(stage.ID().FullyQualifiedName()).
						HasQueryTagString("test-query-tag").
						HasQuotedIdentifiersIgnoreCaseString("true").
						HasReplaceInvalidCharactersString("true").
						HasRequireStorageIntegrationForStageCreationString("true").
						HasRequireStorageIntegrationForStageOperationString("true").
						HasRowsPerResultsetString("1000").
						HasSearchPathString("$current, $public").
						HasServerlessTaskMaxStatementSizeString(string(sdk.WarehouseSizeXLarge)).
						HasServerlessTaskMinStatementSizeString(string(sdk.WarehouseSizeSmall)).
						HasSsoLoginPageString("true").
						HasStatementQueuedTimeoutInSecondsString("1").
						HasStatementTimeoutInSecondsString("1").
						HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyOptimized)).
						HasStrictJsonOutputString("true").
						HasSuspendTaskAfterNumFailuresString("3").
						HasTaskAutoRetryAttemptsString("3").
						HasTimestampDayIsAlways24hString("true").
						HasTimestampInputFormatString("YYYY-MM-DD").
						HasTimestampLtzOutputFormatString("YYYY-MM-DD").
						HasTimestampNtzOutputFormatString("YYYY-MM-DD").
						HasTimestampOutputFormatString("YYYY-MM-DD").
						HasTimestampTypeMappingString(string(sdk.TimestampTypeMappingLtz)).
						HasTimestampTzOutputFormatString("YYYY-MM-DD").
						HasTimezoneString("Europe/London").
						HasTimeInputFormatString("YYYY-MM-DD").
						HasTimeOutputFormatString("YYYY-MM-DD").
						HasTraceLevelString(string(sdk.TraceLevelPropagate)).
						HasTransactionAbortOnErrorString("true").
						HasTransactionDefaultIsolationLevelString(string(sdk.TransactionDefaultIsolationLevelReadCommitted)).
						HasTwoDigitCenturyStartString("1971").
						HasUnsupportedDdlActionString(string(sdk.UnsupportedDDLActionFail)).
						HasUserTaskManagedInitialWarehouseSizeString(string(sdk.WarehouseSizeSmall)).
						HasUserTaskMinimumTriggerIntervalInSecondsString("10").
						HasUserTaskTimeoutMsString("10").
						HasUseCachedResultString("false").
						HasWeekOfYearPolicyString("1").
						HasWeekStartString("1"),
				),
			},
			// import with all parameters set
			{
				Config:       config.FromModels(t, setParametersModel),
				ResourceName: unsetParametersModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedCurrentAccountResource(t, "current_account").
						HasAbortDetachedQueryString("true").
						HasAllowClientMfaCachingString("true").
						HasAllowIdTokenString("true").
						HasAutocommitString("false").
						HasBaseLocationPrefixString("STORAGE_BASE_URL/").
						HasBinaryInputFormatString(string(sdk.BinaryInputFormatBase64)).
						HasBinaryOutputFormatString(string(sdk.BinaryOutputFormatBase64)).
						HasCatalogString(helpers.TestDatabaseCatalog.Name()).
						HasClientEnableLogInfoStatementParametersString("true").
						HasClientEncryptionKeySizeString("256").
						HasClientMemoryLimitString("1540").
						HasClientMetadataRequestUseConnectionCtxString("true").
						HasClientMetadataUseSessionDatabaseString("true").
						HasClientPrefetchThreadsString("5").
						HasClientResultChunkSizeString("159").
						HasClientResultColumnCaseInsensitiveString("true").
						HasClientSessionKeepAliveString("true").
						HasClientSessionKeepAliveHeartbeatFrequencyString("3599").
						HasClientTimestampTypeMappingString(string(sdk.ClientTimestampTypeMappingNtz)).
						HasCortexEnabledCrossRegionString("ANY_REGION").
						HasCortexModelsAllowlistString("All").
						HasCsvTimestampFormatString("YYYY-MM-DD").
						HasDataRetentionTimeInDaysString("2").
						HasDateInputFormatString("YYYY-MM-DD").
						HasDateOutputFormatString("YYYY-MM-DD").
						HasDefaultDdlCollationString("en-cs").
						HasDefaultNotebookComputePoolCpuString("CPU_X64_S").
						HasDefaultNotebookComputePoolGpuString("GPU_NV_S").
						HasDefaultNullOrderingString(string(sdk.DefaultNullOrderingFirst)).
						HasDefaultStreamlitNotebookWarehouseString(warehouseId.Name()).
						HasDisableUiDownloadButtonString("true").
						HasDisableUserPrivilegeGrantsString("true").
						HasEnableAutomaticSensitiveDataClassificationLogString("false").
						HasEnableEgressCostOptimizerString("false").
						HasEnableIdentifierFirstLoginString("false").
						HasEnableTriSecretAndRekeyOptOutForImageRepositoryString("true").
						HasEnableTriSecretAndRekeyOptOutForSpcsBlockStorageString("true").
						HasEnableUnhandledExceptionsReportingString("false").
						HasEnableUnloadPhysicalTypeOptimizationString("false").
						HasEnableUnredactedQuerySyntaxErrorString("true").
						HasEnableUnredactedSecureObjectErrorString("true").
						HasEnforceNetworkRulesForInternalStagesString("true").
						HasErrorOnNondeterministicMergeString("false").
						HasErrorOnNondeterministicUpdateString("true").
						HasEventTableString(eventTable.ID().FullyQualifiedName()).
						HasExternalOauthAddPrivilegedRolesToBlockedListString("false").
						HasExternalVolumeString(externalVolumeId.Name()).
						HasGeographyOutputFormatString(string(sdk.GeographyOutputFormatWKT)).
						HasGeometryOutputFormatString(string(sdk.GeometryOutputFormatWKT)).
						HasHybridTableLockTimeoutString("3599").
						HasInitialReplicationSizeLimitInTbString("9.9").
						HasJdbcTreatDecimalAsIntString("false").
						HasJdbcTreatTimestampNtzAsUtcString("true").
						HasJdbcUseSessionTimezoneString("false").
						HasJsonIndentString("4").
						HasJsTreatIntegerAsBigintString("true").
						HasListingAutoFulfillmentReplicationRefreshScheduleString("2 minutes").
						HasLockTimeoutString("43201").
						HasLogLevelString(string(sdk.LogLevelInfo)).
						HasMaxConcurrencyLevelString("7").
						HasMaxDataExtensionTimeInDaysString("13").
						HasMetricLevelString(string(sdk.MetricLevelAll)).
						HasMinDataRetentionTimeInDaysString("1").
						HasMultiStatementCountString("0").
						HasNetworkPolicyString(networkPolicy.ID().Name()).
						HasNoorderSequenceAsDefaultString("false").
						HasOauthAddPrivilegedRolesToBlockedListString("false").
						HasOdbcTreatDecimalAsIntString("true").
						HasPeriodicDataRekeyingString("false").
						HasPipeExecutionPausedString("true").
						HasPreventUnloadToInlineUrlString("true").
						HasPreventUnloadToInternalStagesString("true").
						HasPythonProfilerTargetStageString(stage.ID().FullyQualifiedName()).
						HasQueryTagString("test-query-tag").
						HasQuotedIdentifiersIgnoreCaseString("true").
						HasReplaceInvalidCharactersString("true").
						HasRequireStorageIntegrationForStageCreationString("true").
						HasRequireStorageIntegrationForStageOperationString("true").
						HasRowsPerResultsetString("1000").
						HasSearchPathString("$current, $public").
						HasServerlessTaskMaxStatementSizeString(string(sdk.WarehouseSizeXLarge)).
						HasServerlessTaskMinStatementSizeString(string(sdk.WarehouseSizeSmall)).
						HasSsoLoginPageString("true").
						HasStatementQueuedTimeoutInSecondsString("1").
						HasStatementTimeoutInSecondsString("1").
						HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyOptimized)).
						HasStrictJsonOutputString("true").
						HasSuspendTaskAfterNumFailuresString("3").
						HasTaskAutoRetryAttemptsString("3").
						HasTimestampDayIsAlways24hString("true").
						HasTimestampInputFormatString("YYYY-MM-DD").
						HasTimestampLtzOutputFormatString("YYYY-MM-DD").
						HasTimestampNtzOutputFormatString("YYYY-MM-DD").
						HasTimestampOutputFormatString("YYYY-MM-DD").
						HasTimestampTypeMappingString(string(sdk.TimestampTypeMappingLtz)).
						HasTimestampTzOutputFormatString("YYYY-MM-DD").
						HasTimezoneString("Europe/London").
						HasTimeInputFormatString("YYYY-MM-DD").
						HasTimeOutputFormatString("YYYY-MM-DD").
						HasTraceLevelString(string(sdk.TraceLevelPropagate)).
						HasTransactionAbortOnErrorString("true").
						HasTransactionDefaultIsolationLevelString(string(sdk.TransactionDefaultIsolationLevelReadCommitted)).
						HasTwoDigitCenturyStartString("1971").
						HasUnsupportedDdlActionString(string(sdk.UnsupportedDDLActionFail)).
						HasUserTaskManagedInitialWarehouseSizeString(string(sdk.WarehouseSizeSmall)).
						HasUserTaskMinimumTriggerIntervalInSecondsString("10").
						HasUserTaskTimeoutMsString("10").
						HasUseCachedResultString("false").
						HasWeekOfYearPolicyString("1").
						HasWeekStartString("1"),
				),
			},
			// unset parameters
			{
				Config: config.FromModels(t, unsetParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentAccountResource(t, setParametersModel.ResourceReference()).
						HasAbortDetachedQueryString("true").
						HasAllowClientMfaCachingString("true").
						HasAllowIdTokenString("true").
						HasAutocommitString("false").
						HasBaseLocationPrefixString("STORAGE_BASE_URL/").
						HasBinaryInputFormatString(string(sdk.BinaryInputFormatBase64)).
						HasBinaryOutputFormatString(string(sdk.BinaryOutputFormatBase64)).
						HasCatalogString(helpers.TestDatabaseCatalog.Name()).
						HasClientEnableLogInfoStatementParametersString("true").
						HasClientEncryptionKeySizeString("256").
						HasClientMemoryLimitString("1540").
						HasClientMetadataRequestUseConnectionCtxString("true").
						HasClientMetadataUseSessionDatabaseString("true").
						HasClientPrefetchThreadsString("5").
						HasClientResultChunkSizeString("159").
						HasClientResultColumnCaseInsensitiveString("true").
						HasClientSessionKeepAliveString("true").
						HasClientSessionKeepAliveHeartbeatFrequencyString("3599").
						HasClientTimestampTypeMappingString(string(sdk.ClientTimestampTypeMappingNtz)).
						HasCortexEnabledCrossRegionString("ANY_REGION").
						HasCortexModelsAllowlistString("All").
						HasCsvTimestampFormatString("YYYY-MM-DD").
						HasDataRetentionTimeInDaysString("2").
						HasDateInputFormatString("YYYY-MM-DD").
						HasDateOutputFormatString("YYYY-MM-DD").
						HasDefaultDdlCollationString("en-cs").
						HasDefaultNotebookComputePoolCpuString("CPU_X64_S").
						HasDefaultNotebookComputePoolGpuString("GPU_NV_S").
						HasDefaultNullOrderingString(string(sdk.DefaultNullOrderingFirst)).
						HasDefaultStreamlitNotebookWarehouseString(warehouseId.Name()).
						HasDisableUiDownloadButtonString("true").
						HasDisableUserPrivilegeGrantsString("true").
						HasEnableAutomaticSensitiveDataClassificationLogString("false").
						HasEnableEgressCostOptimizerString("false").
						HasEnableIdentifierFirstLoginString("false").
						HasEnableTriSecretAndRekeyOptOutForImageRepositoryString("true").
						HasEnableTriSecretAndRekeyOptOutForSpcsBlockStorageString("true").
						HasEnableUnhandledExceptionsReportingString("false").
						HasEnableUnloadPhysicalTypeOptimizationString("false").
						HasEnableUnredactedQuerySyntaxErrorString("true").
						HasEnableUnredactedSecureObjectErrorString("true").
						HasEnforceNetworkRulesForInternalStagesString("true").
						HasErrorOnNondeterministicMergeString("false").
						HasErrorOnNondeterministicUpdateString("true").
						HasEventTableString(eventTable.ID().FullyQualifiedName()).
						HasExternalOauthAddPrivilegedRolesToBlockedListString("false").
						HasExternalVolumeString(externalVolumeId.Name()).
						HasGeographyOutputFormatString(string(sdk.GeographyOutputFormatWKT)).
						HasGeometryOutputFormatString(string(sdk.GeometryOutputFormatWKT)).
						HasHybridTableLockTimeoutString("3599").
						HasInitialReplicationSizeLimitInTbString("9.9").
						HasJdbcTreatDecimalAsIntString("false").
						HasJdbcTreatTimestampNtzAsUtcString("true").
						HasJdbcUseSessionTimezoneString("false").
						HasJsonIndentString("4").
						HasJsTreatIntegerAsBigintString("true").
						HasListingAutoFulfillmentReplicationRefreshScheduleString("2 minutes").
						HasLockTimeoutString("43201").
						HasLogLevelString(string(sdk.LogLevelInfo)).
						HasMaxConcurrencyLevelString("7").
						HasMaxDataExtensionTimeInDaysString("13").
						HasMetricLevelString(string(sdk.MetricLevelAll)).
						HasMinDataRetentionTimeInDaysString("1").
						HasMultiStatementCountString("0").
						HasNetworkPolicyString(networkPolicy.ID().Name()).
						HasNoorderSequenceAsDefaultString("false").
						HasOauthAddPrivilegedRolesToBlockedListString("false").
						HasOdbcTreatDecimalAsIntString("true").
						HasPeriodicDataRekeyingString("false").
						HasPipeExecutionPausedString("true").
						HasPreventUnloadToInlineUrlString("true").
						HasPreventUnloadToInternalStagesString("true").
						HasPythonProfilerTargetStageString(stage.ID().FullyQualifiedName()).
						HasQueryTagString("test-query-tag").
						HasQuotedIdentifiersIgnoreCaseString("true").
						HasReplaceInvalidCharactersString("true").
						HasRequireStorageIntegrationForStageCreationString("true").
						HasRequireStorageIntegrationForStageOperationString("true").
						HasRowsPerResultsetString("1000").
						HasSearchPathString("$current, $public").
						HasServerlessTaskMaxStatementSizeString(string(sdk.WarehouseSizeXLarge)).
						HasServerlessTaskMinStatementSizeString(string(sdk.WarehouseSizeSmall)).
						HasSsoLoginPageString("true").
						HasStatementQueuedTimeoutInSecondsString("1").
						HasStatementTimeoutInSecondsString("1").
						HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyOptimized)).
						HasStrictJsonOutputString("true").
						HasSuspendTaskAfterNumFailuresString("3").
						HasTaskAutoRetryAttemptsString("3").
						HasTimestampDayIsAlways24hString("true").
						HasTimestampInputFormatString("YYYY-MM-DD").
						HasTimestampLtzOutputFormatString("YYYY-MM-DD").
						HasTimestampNtzOutputFormatString("YYYY-MM-DD").
						HasTimestampOutputFormatString("YYYY-MM-DD").
						HasTimestampTypeMappingString(string(sdk.TimestampTypeMappingLtz)).
						HasTimestampTzOutputFormatString("YYYY-MM-DD").
						HasTimezoneString("Europe/London").
						HasTimeInputFormatString("YYYY-MM-DD").
						HasTimeOutputFormatString("YYYY-MM-DD").
						HasTraceLevelString(string(sdk.TraceLevelPropagate)).
						HasTransactionAbortOnErrorString("true").
						HasTransactionDefaultIsolationLevelString(string(sdk.TransactionDefaultIsolationLevelReadCommitted)).
						HasTwoDigitCenturyStartString("1971").
						HasUnsupportedDdlActionString(string(sdk.UnsupportedDDLActionFail)).
						HasUserTaskManagedInitialWarehouseSizeString(string(sdk.WarehouseSizeSmall)).
						HasUserTaskMinimumTriggerIntervalInSecondsString("10").
						HasUserTaskTimeoutMsString("10").
						HasUseCachedResultString("false").
						HasWeekOfYearPolicyString("1").
						HasWeekStartString("1"),
				),
			},
			// Test for external changes
			{
				PreConfig: func() {
					testClient().Account.Alter(t, &sdk.AlterAccountOptions{Set: &sdk.AccountSet{Parameters: &sdk.AccountParameters{AbortDetachedQuery: sdk.Bool(false)}}})
				},
				Config: config.FromModels(t, unsetParametersModel),
				Check: assertThat(t,
					resourceassert.CurrentAccountResource(t, setParametersModel.ResourceReference()).
						HasAbortDetachedQueryString("true"),
				),
			},
		},
	})
}
