package resourceassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func (c *CurrentAccountResourceAssert) HasAllDefaultParameters() *CurrentAccountResourceAssert {
	return c.
		HasAbortDetachedQueryString("false").
		HasAllowBindValuesAccessString("true").
		HasAllowClientMfaCachingString("false").
		HasAllowIdTokenString("false").
		HasAllowedSpcsWorkloadTypesString("ALL").
		HasAutocommitString("true").
		HasBaseLocationPrefixEmpty().
		HasBinaryInputFormatString(string(sdk.BinaryInputFormatHex)).
		HasBinaryOutputFormatString(string(sdk.BinaryOutputFormatHex)).
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
		HasClientTimestampTypeMappingString(string(sdk.ClientTimestampTypeMappingLtz)).
		HasCortexEnabledCrossRegionString("DISABLED").
		HasCortexModelsAllowlistString("ALL").
		HasCsvTimestampFormatEmpty().
		HasDataMetricScheduleString("60 MINUTES").
		HasDataRetentionTimeInDaysString("1").
		HasDateInputFormatString("AUTO").
		HasDateOutputFormatString("YYYY-MM-DD").
		HasDefaultDdlCollationEmpty().
		HasDefaultDbtVersionString("1.9.4").
		HasDefaultNotebookComputePoolCpuString("SYSTEM_COMPUTE_POOL_CPU").
		HasDefaultNotebookComputePoolGpuString("SYSTEM_COMPUTE_POOL_GPU").
		HasDefaultNullOrderingString(string(sdk.DefaultNullOrderingLast)).
		// TODO [SNOW-3797718]: changed temporarily - what is the REGRESS warehouse?
		HasDefaultStreamlitNotebookWarehouseString(sdk.NewAccountObjectIdentifier("SYSTEM$STREAMLIT_NOTEBOOK_WH").FullyQualifiedName()).
		HasDisableUiDownloadButtonString("false").
		HasDisableUserPrivilegeGrantsString("false").
		HasDisallowedSpcsWorkloadTypesEmpty().
		HasEnableAutomaticSensitiveDataClassificationLogString("true").
		HasEnableBudgetEventLoggingString("true").
		HasEnableCortexAnalystString("false").
		HasEnableDataCompactionString("true").
		HasEnableEgressCostOptimizerString("true").
		HasEnableGetDdlUseDataTypeAliasString("false").
		HasEnableIcebergMergeOnReadString("true").
		HasEnableNotebookCreationInPersonalDbString("false").
		HasEnableSpcsBlockStorageSnowflakeFullEncryptionEnforcementString("false").
		HasEnableTagPropagationEventLoggingString("false").
		// TODO [SNOW-3797718]: commented out as some other test rewrites this value
		// HasEnableIdentifierFirstLoginString("true").
		// TODO [SNOW-3797718]: commented out as some other test rewrites this value
		// HasEnableTriSecretAndRekeyOptOutForImageRepositoryString("false").
		HasEnableTriSecretAndRekeyOptOutForSpcsBlockStorageString("false").
		HasEnableUnhandledExceptionsReportingString("true").
		HasEnableUnloadPhysicalTypeOptimizationString("true").
		HasEnableUnredactedQuerySyntaxErrorString("false").
		HasEnableUnredactedSecureObjectErrorString("false").
		HasEnforceNetworkRulesForInternalStagesString("false").
		HasErrorOnNondeterministicMergeString("true").
		HasErrorOnNondeterministicUpdateString("false").
		HasEventTableString(sdk.NewSchemaObjectIdentifier("snowflake", "telemetry", "events").FullyQualifiedName()).
		HasExternalOauthAddPrivilegedRolesToBlockedListString("true").
		HasExternalVolumeEmpty().
		HasGeographyOutputFormatString(string(sdk.GeographyOutputFormatGeoJSON)).
		HasGeometryOutputFormatString(string(sdk.GeometryOutputFormatGeoJSON)).
		HasHybridTableLockTimeoutString("3600").
		HasIcebergVersionDefaultString("2").
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
		HasReadConsistencyModeString("SESSION").
		HasReplaceInvalidCharactersString("false").
		HasRequireStorageIntegrationForStageCreationString("false").
		HasRequireStorageIntegrationForStageOperationString("false").
		HasRowTimestampDefaultString("false").
		HasRowsPerResultsetString("0").
		HasSearchPathString("$current, $public").
		HasServerlessTaskMaxStatementSizeString(string(sdk.WarehouseSizeXXLarge)).
		HasServerlessTaskMinStatementSizeString(string(sdk.WarehouseSizeXSmall)).
		HasSsoLoginPageString("false").
		HasSqlTraceQueryTextString("OFF").
		HasStatementQueuedTimeoutInSecondsString("0").
		HasStatementTimeoutInSecondsString("172800").
		HasStorageSerializationPolicyString(string(sdk.StorageSerializationPolicyOptimized)).
		HasStrictJsonOutputString("false").
		HasSuspendTaskAfterNumFailuresString("10").
		HasTaskAutoRetryAttemptsString("0").
		HasTimestampDayIsAlways24hString("false").
		HasTimestampInputFormatString("AUTO").
		HasTimestampLtzOutputFormatEmpty().
		HasTimestampNtzOutputFormatString("YYYY-MM-DD HH24:MI:SS.FF3").
		HasTimestampOutputFormatString("YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM").
		HasTimestampTypeMappingString(string(sdk.TimestampTypeMappingNtz)).
		HasTimestampTzOutputFormatEmpty().
		HasTimezoneString("America/Los_Angeles").
		HasTimeInputFormatString("AUTO").
		HasTimeOutputFormatString("HH24:MI:SS").
		HasTraceLevelString("OFF").
		HasTransactionAbortOnErrorString("false").
		HasTransactionDefaultIsolationLevelString(string(sdk.TransactionDefaultIsolationLevelReadCommitted)).
		HasTwoDigitCenturyStartString("1970").
		HasUnsupportedDdlActionString(string(sdk.UnsupportedDDLActionIgnore)).
		HasUserTaskManagedInitialWarehouseSizeString(string(sdk.WarehouseSizeMedium)).
		HasUserTaskMinimumTriggerIntervalInSecondsString("30").
		HasUserTaskTimeoutMsString("3600000").
		HasUseCachedResultString("true").
		HasUseWorkspacesForSqlEmpty().
		HasWeekOfYearPolicyString("0").
		HasWeekStartString("0")
}
