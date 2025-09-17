package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionParameters_setParam(t *testing.T) {
	tests := []struct {
		parameter     SessionParameter
		value         string
		expectedValue any
		accessor      func(*SessionParameters) any
	}{
		{parameter: SessionParameterAbortDetachedQuery, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.AbortDetachedQuery }},
		{parameter: SessionParameterActivePythonProfiler, value: "LINE", expectedValue: ActivePythonProfilerLine, accessor: func(sp *SessionParameters) any { return *sp.ActivePythonProfiler }},
		{parameter: SessionParameterAutocommit, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.Autocommit }},
		{parameter: SessionParameterBinaryInputFormat, value: "some", expectedValue: BinaryInputFormat("some"), accessor: func(sp *SessionParameters) any { return *sp.BinaryInputFormat }},
		{parameter: SessionParameterBinaryOutputFormat, value: "some", expectedValue: BinaryOutputFormat("some"), accessor: func(sp *SessionParameters) any { return *sp.BinaryOutputFormat }},
		{parameter: SessionParameterClientEnableLogInfoStatementParameters, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ClientEnableLogInfoStatementParameters }},
		{parameter: SessionParameterClientMemoryLimit, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.ClientMemoryLimit }},
		{parameter: SessionParameterClientMetadataRequestUseConnectionCtx, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ClientMetadataRequestUseConnectionCtx }},
		{parameter: SessionParameterClientPrefetchThreads, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.ClientPrefetchThreads }},
		{parameter: SessionParameterClientResultChunkSize, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.ClientResultChunkSize }},
		{parameter: SessionParameterClientResultColumnCaseInsensitive, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ClientResultColumnCaseInsensitive }},
		{parameter: SessionParameterClientMetadataUseSessionDatabase, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ClientMetadataUseSessionDatabase }},
		{parameter: SessionParameterClientSessionKeepAlive, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ClientSessionKeepAlive }},
		{parameter: SessionParameterClientSessionKeepAliveHeartbeatFrequency, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.ClientSessionKeepAliveHeartbeatFrequency }},
		{parameter: SessionParameterClientTimestampTypeMapping, value: "some", expectedValue: ClientTimestampTypeMapping("some"), accessor: func(sp *SessionParameters) any { return *sp.ClientTimestampTypeMapping }},
		{parameter: SessionParameterCsvTimestampFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.CsvTimestampFormat }},
		{parameter: SessionParameterDateInputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.DateInputFormat }},
		{parameter: SessionParameterDateOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.DateOutputFormat }},
		{parameter: SessionParameterEnableUnloadPhysicalTypeOptimization, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.EnableUnloadPhysicalTypeOptimization }},
		{parameter: SessionParameterErrorOnNondeterministicMerge, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ErrorOnNondeterministicMerge }},
		{parameter: SessionParameterErrorOnNondeterministicUpdate, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.ErrorOnNondeterministicUpdate }},
		{parameter: SessionParameterGeographyOutputFormat, value: "some", expectedValue: GeographyOutputFormat("some"), accessor: func(sp *SessionParameters) any { return *sp.GeographyOutputFormat }},
		{parameter: SessionParameterGeometryOutputFormat, value: "some", expectedValue: GeometryOutputFormat("some"), accessor: func(sp *SessionParameters) any { return *sp.GeometryOutputFormat }},
		{parameter: SessionParameterHybridTableLockTimeout, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.HybridTableLockTimeout }},
		{parameter: SessionParameterJdbcTreatDecimalAsInt, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.JdbcTreatDecimalAsInt }},
		{parameter: SessionParameterJdbcTreatTimestampNtzAsUtc, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.JdbcTreatTimestampNtzAsUtc }},
		{parameter: SessionParameterJdbcUseSessionTimezone, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.JdbcUseSessionTimezone }},
		{parameter: SessionParameterJsonIndent, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.JsonIndent }},
		{parameter: SessionParameterJsTreatIntegerAsBigInt, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.JsTreatIntegerAsBigInt }},
		{parameter: SessionParameterLockTimeout, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.LockTimeout }},
		{parameter: SessionParameterLogLevel, value: "some", expectedValue: LogLevel("some"), accessor: func(sp *SessionParameters) any { return *sp.LogLevel }},
		{parameter: SessionParameterMultiStatementCount, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.MultiStatementCount }},
		{parameter: SessionParameterNoorderSequenceAsDefault, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.NoorderSequenceAsDefault }},
		{parameter: SessionParameterOdbcTreatDecimalAsInt, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.OdbcTreatDecimalAsInt }},
		{parameter: SessionParameterPythonProfilerModules, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.PythonProfilerModules }},
		{parameter: SessionParameterPythonProfilerTargetStage, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.PythonProfilerTargetStage }},
		{parameter: SessionParameterQueryTag, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.QueryTag }},
		{parameter: SessionParameterQuotedIdentifiersIgnoreCase, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.QuotedIdentifiersIgnoreCase }},
		{parameter: SessionParameterRowsPerResultset, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.RowsPerResultset }},
		{parameter: SessionParameterS3StageVpceDnsName, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.S3StageVpceDnsName }},
		{parameter: SessionParameterSearchPath, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.SearchPath }},
		{parameter: SessionParameterSimulatedDataSharingConsumer, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.SimulatedDataSharingConsumer }},
		{parameter: SessionParameterStatementQueuedTimeoutInSeconds, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.StatementQueuedTimeoutInSeconds }},
		{parameter: SessionParameterStatementTimeoutInSeconds, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.StatementTimeoutInSeconds }},
		{parameter: SessionParameterStrictJsonOutput, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.StrictJsonOutput }},
		{parameter: SessionParameterTimestampDayIsAlways24h, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.TimestampDayIsAlways24h }},
		{parameter: SessionParameterTimestampInputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampInputFormat }},
		{parameter: SessionParameterTimestampLTZOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampLTZOutputFormat }},
		{parameter: SessionParameterTimestampNTZOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampNTZOutputFormat }},
		{parameter: SessionParameterTimestampOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampOutputFormat }},
		{parameter: SessionParameterTimestampTypeMapping, value: "some", expectedValue: TimestampTypeMapping("some"), accessor: func(sp *SessionParameters) any { return *sp.TimestampTypeMapping }},
		{parameter: SessionParameterTimestampTZOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimestampTZOutputFormat }},
		{parameter: SessionParameterTimezone, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.Timezone }},
		{parameter: SessionParameterTimeInputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimeInputFormat }},
		{parameter: SessionParameterTimeOutputFormat, value: "some", expectedValue: "some", accessor: func(sp *SessionParameters) any { return *sp.TimeOutputFormat }},
		{parameter: SessionParameterTraceLevel, value: "some", expectedValue: TraceLevel("some"), accessor: func(sp *SessionParameters) any { return *sp.TraceLevel }},
		{parameter: SessionParameterTransactionAbortOnError, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.TransactionAbortOnError }},
		{parameter: SessionParameterTransactionDefaultIsolationLevel, value: "some", expectedValue: TransactionDefaultIsolationLevel("some"), accessor: func(sp *SessionParameters) any { return *sp.TransactionDefaultIsolationLevel }},
		{parameter: SessionParameterTwoDigitCenturyStart, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.TwoDigitCenturyStart }},
		{parameter: SessionParameterUnsupportedDDLAction, value: "some", expectedValue: UnsupportedDDLAction("some"), accessor: func(sp *SessionParameters) any { return *sp.UnsupportedDDLAction }},
		{parameter: SessionParameterUseCachedResult, value: "true", expectedValue: true, accessor: func(sp *SessionParameters) any { return *sp.UseCachedResult }},
		{parameter: SessionParameterWeekOfYearPolicy, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.WeekOfYearPolicy }},
		{parameter: SessionParameterWeekStart, value: "1", expectedValue: 1, accessor: func(sp *SessionParameters) any { return *sp.WeekStart }},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("test valid value '%s' for parameter %s", tt.value, tt.parameter), func(t *testing.T) {
			sessionParameters := &SessionParameters{}

			err := sessionParameters.setParam(tt.parameter, tt.value)

			require.NoError(t, err)
			require.Equal(t, tt.expectedValue, tt.accessor(sessionParameters))
		})
	}

	// TODO [SNOW-2316589]: add validation to string type
	invalidCases := []struct {
		parameter SessionParameter
		value     string
	}{
		{parameter: SessionParameterAbortDetachedQuery, value: "true123"},
		{parameter: SessionParameterAutocommit, value: "true123"},
		// {parameter: SessionParameterBinaryInputFormat, value: "some"}, // add validation
		// {parameter: SessionParameterBinaryOutputFormat, value: "some"}, // add validation
		{parameter: SessionParameterClientMetadataRequestUseConnectionCtx, value: "true123"},
		{parameter: SessionParameterClientMetadataUseSessionDatabase, value: "true123"},
		{parameter: SessionParameterClientResultColumnCaseInsensitive, value: "true123"},
		// {parameter: SessionParameterDateInputFormat, value: "some"}, // add validation
		// {parameter: SessionParameterGeographyOutputFormat, value: "some"}, // add validation
		// {parameter: SessionParameterDateOutputFormat, value: "some"}, // add validation
		{parameter: SessionParameterErrorOnNondeterministicMerge, value: "true123"},
		{parameter: SessionParameterErrorOnNondeterministicUpdate, value: "true123"},
		{parameter: SessionParameterJsonIndent, value: "aaa"},
		{parameter: SessionParameterLockTimeout, value: "aaa"},
		{parameter: SessionParameterMultiStatementCount, value: "aaa"},
		// {parameter: SessionParameterQueryTag, value: "some"}, // add validation
		{parameter: SessionParameterQuotedIdentifiersIgnoreCase, value: "true123"},
		{parameter: SessionParameterRowsPerResultset, value: "aaa"},
		// {parameter: SessionParameterSimulatedDataSharingConsumer, value: "some"}, // add validation
		{parameter: SessionParameterStatementTimeoutInSeconds, value: "aaa"},
		{parameter: SessionParameterStrictJsonOutput, value: "true123"},
		// {parameter: SessionParameterTimeInputFormat, value: "some"}, // add validation
		// {parameter: SessionParameterTimeOutputFormat, value: "some"}, // add validation
		{parameter: SessionParameterTimestampDayIsAlways24h, value: "true123"},
		// {parameter: SessionParameterTimestampInputFormat, value: "some"}, // add validation
		// {parameter: SessionParameterTimestampLTZOutputFormat, value: "some"}, // add validation
		// {parameter: SessionParameterTimestampNTZOutputFormat, value: "some"}, // add validation
		// {parameter: SessionParameterTimestampOutputFormat, value: "some"}, // add validation
		// {parameter: SessionParameterTimestampTypeMapping, value: "some"}, // add validation
		// {parameter: SessionParameterTimestampTZOutputFormat, value: "some"}, // add validation
		// {parameter: SessionParameterTimezone, value: "some"}, // add validation
		// {parameter: SessionParameterTransactionDefaultIsolationLevel, value: "some"}, // add validation
		{parameter: SessionParameterTwoDigitCenturyStart, value: "aaa"},
		// {parameter: SessionParameterUnsupportedDDLAction, value: "some"}, // add validation
		{parameter: SessionParameterUseCachedResult, value: "true123"},
		{parameter: SessionParameterWeekOfYearPolicy, value: "aaa"},
		{parameter: SessionParameterWeekStart, value: "aaa"},
	}
	for _, tt := range invalidCases {
		t.Run(fmt.Sprintf("test invalid value '%s' for parameter %s", tt.value, tt.parameter), func(t *testing.T) {
			sessionParameters := &SessionParameters{}

			err := sessionParameters.setParam(tt.parameter, tt.value)

			require.Error(t, err)
		})
	}
}

func Test_LegacyAccountParameters_setParam(t *testing.T) {
	tests := []struct {
		parameter     AccountParameter
		value         string
		expectedValue any
		accessor      func(accountParameters *LegacyAccountParameters) any
	}{
		{parameter: AccountParameterAllowClientMFACaching, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.AllowClientMFACaching }},
		{parameter: AccountParameterAllowIDToken, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.AllowIDToken }},
		{parameter: AccountParameterClientEncryptionKeySize, value: "1", expectedValue: 1, accessor: func(lap *LegacyAccountParameters) any { return *lap.ClientEncryptionKeySize }},
		{parameter: AccountParameterCortexEnabledCrossRegion, value: "some", expectedValue: "some", accessor: func(lap *LegacyAccountParameters) any { return *lap.CortexEnabledCrossRegion }},
		{parameter: AccountParameterDisableUserPrivilegeGrants, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.DisableUserPrivilegeGrants }},
		{parameter: AccountParameterEnableIdentifierFirstLogin, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.EnableIdentifierFirstLogin }},
		{parameter: AccountParameterEnableInternalStagesPrivatelink, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.EnableInternalStagesPrivatelink }},
		{parameter: AccountParameterEnablePersonalDatabase, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.EnablePersonalDatabase }},
		{parameter: AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.EnableTriSecretAndRekeyOptOutForImageRepository }},
		{parameter: AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.EnableTriSecretAndRekeyOptOutForSpcsBlockStorage }},
		{parameter: AccountParameterEnableUnhandledExceptionsReporting, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.EnableUnhandledExceptionsReporting }},
		{parameter: AccountParameterEnableUnredactedQuerySyntaxError, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.EnableUnredactedQuerySyntaxError }},
		{parameter: AccountParameterEnforceNetworkRulesForInternalStages, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.EnforceNetworkRulesForInternalStages }},
		{parameter: AccountParameterEventTable, value: "some", expectedValue: "some", accessor: func(lap *LegacyAccountParameters) any { return *lap.EventTable }},
		{parameter: AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.ExternalOAuthAddPrivilegedRolesToBlockedList }},
		{parameter: AccountParameterInitialReplicationSizeLimitInTB, value: "some", expectedValue: "some", accessor: func(lap *LegacyAccountParameters) any { return *lap.InitialReplicationSizeLimitInTB }},
		{parameter: AccountParameterMetricLevel, value: string(MetricLevelAll), expectedValue: MetricLevelAll, accessor: func(lap *LegacyAccountParameters) any { return *lap.MetricLevel }},
		{parameter: AccountParameterMinDataRetentionTimeInDays, value: "1", expectedValue: 1, accessor: func(lap *LegacyAccountParameters) any { return *lap.MinDataRetentionTimeInDays }},
		{parameter: AccountParameterNetworkPolicy, value: "some", expectedValue: "some", accessor: func(lap *LegacyAccountParameters) any { return *lap.NetworkPolicy }},
		{parameter: AccountParameterOAuthAddPrivilegedRolesToBlockedList, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.OAuthAddPrivilegedRolesToBlockedList }},
		{parameter: AccountParameterPeriodicDataRekeying, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.PeriodicDataRekeying }},
		{parameter: AccountParameterPreventLoadFromInlineURL, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.PreventLoadFromInlineURL }},
		{parameter: AccountParameterPreventUnloadToInlineURL, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.PreventUnloadToInlineURL }},
		{parameter: AccountParameterPreventUnloadToInternalStages, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.PreventUnloadToInternalStages }},
		{parameter: AccountParameterRequireStorageIntegrationForStageCreation, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.RequireStorageIntegrationForStageCreation }},
		{parameter: AccountParameterRequireStorageIntegrationForStageOperation, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.RequireStorageIntegrationForStageOperation }},
		{parameter: AccountParameterSsoLoginPage, value: "true", expectedValue: true, accessor: func(lap *LegacyAccountParameters) any { return *lap.SSOLoginPage }},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("test valid value '%s' for account parameter %s", tt.value, tt.parameter), func(t *testing.T) {
			legacyAccountParameters := &LegacyAccountParameters{}

			matched, err := legacyAccountParameters.setParam(tt.parameter, tt.value)

			require.NoError(t, err)
			require.True(t, matched)
			require.Equal(t, tt.expectedValue, tt.accessor(legacyAccountParameters))
		})
	}

	// TODO [SNOW-2316589]: add validation to string type
	invalidCases := []struct {
		parameter AccountParameter
		value     string
	}{
		{parameter: AccountParameterAllowClientMFACaching, value: "true123"},
		{parameter: AccountParameterAllowIDToken, value: "true123"},
		{parameter: AccountParameterClientEncryptionKeySize, value: "aaa"},
		// {parameter: AccountParameterCortexEnabledCrossRegion, value: "some"}, // add validation
		{parameter: AccountParameterDisableUserPrivilegeGrants, value: "true123"},
		{parameter: AccountParameterEnableIdentifierFirstLogin, value: "true123"},
		{parameter: AccountParameterEnableInternalStagesPrivatelink, value: "true123"},
		{parameter: AccountParameterEnablePersonalDatabase, value: "true123"},
		{parameter: AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository, value: "true123"},
		{parameter: AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage, value: "true123"},
		{parameter: AccountParameterEnableUnhandledExceptionsReporting, value: "true123"},
		{parameter: AccountParameterEnableUnredactedQuerySyntaxError, value: "true123"},
		{parameter: AccountParameterEnforceNetworkRulesForInternalStages, value: "true123"},
		// {parameter: AccountParameterEventTable, value: "some"}, // add validation
		{parameter: AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList, value: "true123"},
		// {parameter: AccountParameterInitialReplicationSizeLimitInTB, value: "some"}, // add validation
		// {parameter: AccountParameterMetricLevel, value: "some"}, // add validation
		{parameter: AccountParameterMinDataRetentionTimeInDays, value: "aaa"},
		// {parameter: AccountParameterNetworkPolicy, value: "some"}, // add validation
		{parameter: AccountParameterOAuthAddPrivilegedRolesToBlockedList, value: "true123"},
		{parameter: AccountParameterPeriodicDataRekeying, value: "true123"},
		{parameter: AccountParameterPreventLoadFromInlineURL, value: "true123"},
		{parameter: AccountParameterPreventUnloadToInlineURL, value: "true123"},
		{parameter: AccountParameterPreventUnloadToInternalStages, value: "true123"},
		{parameter: AccountParameterRequireStorageIntegrationForStageCreation, value: "true123"},
		{parameter: AccountParameterRequireStorageIntegrationForStageOperation, value: "true123"},
		{parameter: AccountParameterSsoLoginPage, value: "true123"},
	}
	for _, tt := range invalidCases {
		t.Run(fmt.Sprintf("test invalid value '%s' for account parameter %s", tt.value, tt.parameter), func(t *testing.T) {
			legacyAccountParameters := &LegacyAccountParameters{}

			matched, err := legacyAccountParameters.setParam(tt.parameter, tt.value)

			require.Error(t, err)
			require.True(t, matched)
		})
	}

	nonAccountLevelParametersCases := []struct {
		parameter string
		value     string
	}{
		{parameter: string(SessionParameterWeekStart), value: "1"},
		{parameter: "non_existent_param", value: "some_value"},
	}
	for _, tt := range nonAccountLevelParametersCases {
		t.Run(fmt.Sprintf("test non-account level parameter %s", tt.parameter), func(t *testing.T) {
			legacyAccountParameters := &LegacyAccountParameters{}

			matched, err := legacyAccountParameters.setParam(AccountParameter(tt.parameter), tt.value)

			require.NoError(t, err)
			require.False(t, matched)
			require.Equal(t, LegacyAccountParameters{}, *legacyAccountParameters)
		})
	}
}
