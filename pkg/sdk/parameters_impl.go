package sdk

import (
	"fmt"
	"strconv"
)

func GetSessionParametersFrom(params map[string]any) (*SessionParameters, error) {
	sessionParameters := &SessionParameters{}
	for k, v := range params {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expecting string value for parameter %s (current value: %v)", k, v)
		}
		err := sessionParameters.setParam(SessionParameter(k), s)
		if err != nil {
			return nil, err
		}
	}
	return sessionParameters, nil
}

// TODO [SNOW-1348330]: get type based on the tag in SessionParameters struct and handle in a generic way
// TODO [SNOW-1348330]: use sdk.ToX for the enums
func (sessionParameters *SessionParameters) setParam(parameter SessionParameter, value string) error {
	var err error
	switch parameter {
	case SessionParameterAbortDetachedQuery:
		err = setBooleanValue(parameter, value, &sessionParameters.AbortDetachedQuery)
	case SessionParameterActivePythonProfiler:
		sessionParameters.ActivePythonProfiler = Pointer(ActivePythonProfiler(value))
	case SessionParameterAutocommit:
		err = setBooleanValue(parameter, value, &sessionParameters.Autocommit)
	case SessionParameterBinaryInputFormat:
		sessionParameters.BinaryInputFormat = Pointer(BinaryInputFormat(value))
	case SessionParameterBinaryOutputFormat:
		sessionParameters.BinaryOutputFormat = Pointer(BinaryOutputFormat(value))
	case SessionParameterClientEnableLogInfoStatementParameters:
		err = setBooleanValue(parameter, value, &sessionParameters.ClientEnableLogInfoStatementParameters)
	case SessionParameterClientMemoryLimit:
		err = setIntegerValue(parameter, value, &sessionParameters.ClientMemoryLimit)
	case SessionParameterClientMetadataRequestUseConnectionCtx:
		err = setBooleanValue(parameter, value, &sessionParameters.ClientMetadataRequestUseConnectionCtx)
	case SessionParameterClientPrefetchThreads:
		err = setIntegerValue(parameter, value, &sessionParameters.ClientPrefetchThreads)
	case SessionParameterClientResultChunkSize:
		err = setIntegerValue(parameter, value, &sessionParameters.ClientResultChunkSize)
	case SessionParameterClientResultColumnCaseInsensitive:
		err = setBooleanValue(parameter, value, &sessionParameters.ClientResultColumnCaseInsensitive)
	case SessionParameterClientMetadataUseSessionDatabase:
		err = setBooleanValue(parameter, value, &sessionParameters.ClientMetadataUseSessionDatabase)
	case SessionParameterClientSessionKeepAlive:
		err = setBooleanValue(parameter, value, &sessionParameters.ClientSessionKeepAlive)
	case SessionParameterClientSessionKeepAliveHeartbeatFrequency:
		err = setIntegerValue(parameter, value, &sessionParameters.ClientSessionKeepAliveHeartbeatFrequency)
	case SessionParameterClientTimestampTypeMapping:
		sessionParameters.ClientTimestampTypeMapping = Pointer(ClientTimestampTypeMapping(value))
	case SessionParameterCsvTimestampFormat:
		sessionParameters.CsvTimestampFormat = &value
	case SessionParameterDateInputFormat:
		sessionParameters.DateInputFormat = &value
	case SessionParameterDateOutputFormat:
		sessionParameters.DateOutputFormat = &value
	case SessionParameterEnableUnloadPhysicalTypeOptimization:
		err = setBooleanValue(parameter, value, &sessionParameters.EnableUnloadPhysicalTypeOptimization)
	case SessionParameterErrorOnNondeterministicMerge:
		err = setBooleanValue(parameter, value, &sessionParameters.ErrorOnNondeterministicMerge)
	case SessionParameterErrorOnNondeterministicUpdate:
		err = setBooleanValue(parameter, value, &sessionParameters.ErrorOnNondeterministicUpdate)
	case SessionParameterGeographyOutputFormat:
		sessionParameters.GeographyOutputFormat = Pointer(GeographyOutputFormat(value))
	case SessionParameterGeometryOutputFormat:
		sessionParameters.GeometryOutputFormat = Pointer(GeometryOutputFormat(value))
	case SessionParameterHybridTableLockTimeout:
		err = setIntegerValue(parameter, value, &sessionParameters.HybridTableLockTimeout)
	case SessionParameterJdbcTreatDecimalAsInt:
		err = setBooleanValue(parameter, value, &sessionParameters.JdbcTreatDecimalAsInt)
	case SessionParameterJdbcTreatTimestampNtzAsUtc:
		err = setBooleanValue(parameter, value, &sessionParameters.JdbcTreatTimestampNtzAsUtc)
	case SessionParameterJdbcUseSessionTimezone:
		err = setBooleanValue(parameter, value, &sessionParameters.JdbcUseSessionTimezone)
	case SessionParameterJsonIndent:
		err = setIntegerValue(parameter, value, &sessionParameters.JsonIndent)
	case SessionParameterJsTreatIntegerAsBigInt:
		err = setBooleanValue(parameter, value, &sessionParameters.JsTreatIntegerAsBigInt)
	case SessionParameterLockTimeout:
		err = setIntegerValue(parameter, value, &sessionParameters.LockTimeout)
	case SessionParameterLogLevel:
		sessionParameters.LogLevel = Pointer(LogLevel(value))
	case SessionParameterMultiStatementCount:
		err = setIntegerValue(parameter, value, &sessionParameters.MultiStatementCount)
	case SessionParameterNoorderSequenceAsDefault:
		err = setBooleanValue(parameter, value, &sessionParameters.NoorderSequenceAsDefault)
	case SessionParameterOdbcTreatDecimalAsInt:
		err = setBooleanValue(parameter, value, &sessionParameters.OdbcTreatDecimalAsInt)
	case SessionParameterPythonProfilerModules:
		sessionParameters.PythonProfilerModules = &value
	case SessionParameterPythonProfilerTargetStage:
		sessionParameters.PythonProfilerTargetStage = &value
	case SessionParameterQueryTag:
		sessionParameters.QueryTag = &value
	case SessionParameterQuotedIdentifiersIgnoreCase:
		err = setBooleanValue(parameter, value, &sessionParameters.QuotedIdentifiersIgnoreCase)
	case SessionParameterRowsPerResultset:
		err = setIntegerValue(parameter, value, &sessionParameters.RowsPerResultset)
	case SessionParameterS3StageVpceDnsName:
		sessionParameters.S3StageVpceDnsName = &value
	case SessionParameterSearchPath:
		sessionParameters.SearchPath = &value
	case SessionParameterSimulatedDataSharingConsumer:
		sessionParameters.SimulatedDataSharingConsumer = &value
	case SessionParameterStatementQueuedTimeoutInSeconds:
		err = setIntegerValue(parameter, value, &sessionParameters.StatementQueuedTimeoutInSeconds)
	case SessionParameterStatementTimeoutInSeconds:
		err = setIntegerValue(parameter, value, &sessionParameters.StatementTimeoutInSeconds)
	case SessionParameterStrictJsonOutput:
		err = setBooleanValue(parameter, value, &sessionParameters.StrictJsonOutput)
	case SessionParameterTimestampDayIsAlways24h:
		err = setBooleanValue(parameter, value, &sessionParameters.TimestampDayIsAlways24h)
	case SessionParameterTimestampInputFormat:
		sessionParameters.TimestampInputFormat = &value
	case SessionParameterTimestampLTZOutputFormat:
		sessionParameters.TimestampLTZOutputFormat = &value
	case SessionParameterTimestampNTZOutputFormat:
		sessionParameters.TimestampNTZOutputFormat = &value
	case SessionParameterTimestampOutputFormat:
		sessionParameters.TimestampOutputFormat = &value
	case SessionParameterTimestampTypeMapping:
		sessionParameters.TimestampTypeMapping = Pointer(TimestampTypeMapping(value))
	case SessionParameterTimestampTZOutputFormat:
		sessionParameters.TimestampTZOutputFormat = &value
	case SessionParameterTimezone:
		sessionParameters.Timezone = &value
	case SessionParameterTimeInputFormat:
		sessionParameters.TimeInputFormat = &value
	case SessionParameterTimeOutputFormat:
		sessionParameters.TimeOutputFormat = &value
	case SessionParameterTraceLevel:
		sessionParameters.TraceLevel = Pointer(TraceLevel(value))
	case SessionParameterTransactionAbortOnError:
		err = setBooleanValue(parameter, value, &sessionParameters.TransactionAbortOnError)
	case SessionParameterTransactionDefaultIsolationLevel:
		sessionParameters.TransactionDefaultIsolationLevel = Pointer(TransactionDefaultIsolationLevel(value))
	case SessionParameterTwoDigitCenturyStart:
		err = setIntegerValue(parameter, value, &sessionParameters.TwoDigitCenturyStart)
	case SessionParameterUnsupportedDDLAction:
		sessionParameters.UnsupportedDDLAction = Pointer(UnsupportedDDLAction(value))
	case SessionParameterUseCachedResult:
		err = setBooleanValue(parameter, value, &sessionParameters.UseCachedResult)
	case SessionParameterWeekOfYearPolicy:
		err = setIntegerValue(parameter, value, &sessionParameters.WeekOfYearPolicy)
	case SessionParameterWeekStart:
		err = setIntegerValue(parameter, value, &sessionParameters.WeekStart)
	default:
		err = fmt.Errorf("%s session parameter is not supported", string(parameter))
	}
	return err
}

func GetSessionParametersUnsetFrom(params map[string]any) (*SessionParametersUnset, error) {
	sessionParametersUnset := &SessionParametersUnset{}
	for k := range params {
		err := sessionParametersUnset.setParam(SessionParameter(k))
		if err != nil {
			return nil, err
		}
	}
	return sessionParametersUnset, nil
}

func (sessionParametersUnset *SessionParametersUnset) setParam(parameter SessionParameter) error {
	var unsetField **bool
	switch parameter {
	case SessionParameterAbortDetachedQuery:
		unsetField = &sessionParametersUnset.AbortDetachedQuery
	case SessionParameterActivePythonProfiler:
		unsetField = &sessionParametersUnset.ActivePythonProfiler
	case SessionParameterAutocommit:
		unsetField = &sessionParametersUnset.Autocommit
	case SessionParameterBinaryInputFormat:
		unsetField = &sessionParametersUnset.BinaryInputFormat
	case SessionParameterBinaryOutputFormat:
		unsetField = &sessionParametersUnset.BinaryOutputFormat
	case SessionParameterClientEnableLogInfoStatementParameters:
		unsetField = &sessionParametersUnset.ClientEnableLogInfoStatementParameters
	case SessionParameterClientMemoryLimit:
		unsetField = &sessionParametersUnset.ClientMemoryLimit
	case SessionParameterClientMetadataRequestUseConnectionCtx:
		unsetField = &sessionParametersUnset.ClientMetadataRequestUseConnectionCtx
	case SessionParameterClientMetadataUseSessionDatabase:
		unsetField = &sessionParametersUnset.ClientMetadataUseSessionDatabase
	case SessionParameterClientResultColumnCaseInsensitive:
		unsetField = &sessionParametersUnset.ClientResultColumnCaseInsensitive
	case SessionParameterClientPrefetchThreads:
		unsetField = &sessionParametersUnset.ClientPrefetchThreads
	case SessionParameterClientResultChunkSize:
		unsetField = &sessionParametersUnset.ClientResultChunkSize
	case SessionParameterClientSessionKeepAlive:
		unsetField = &sessionParametersUnset.ClientSessionKeepAlive
	case SessionParameterClientSessionKeepAliveHeartbeatFrequency:
		unsetField = &sessionParametersUnset.ClientSessionKeepAliveHeartbeatFrequency
	case SessionParameterClientTimestampTypeMapping:
		unsetField = &sessionParametersUnset.ClientTimestampTypeMapping
	case SessionParameterCsvTimestampFormat:
		unsetField = &sessionParametersUnset.CsvTimestampFormat
	case SessionParameterDateInputFormat:
		unsetField = &sessionParametersUnset.DateInputFormat
	case SessionParameterDateOutputFormat:
		unsetField = &sessionParametersUnset.DateOutputFormat
	case SessionParameterEnableUnloadPhysicalTypeOptimization:
		unsetField = &sessionParametersUnset.EnableUnloadPhysicalTypeOptimization
	case SessionParameterErrorOnNondeterministicMerge:
		unsetField = &sessionParametersUnset.ErrorOnNondeterministicMerge
	case SessionParameterErrorOnNondeterministicUpdate:
		unsetField = &sessionParametersUnset.ErrorOnNondeterministicUpdate
	case SessionParameterGeographyOutputFormat:
		unsetField = &sessionParametersUnset.GeographyOutputFormat
	case SessionParameterGeometryOutputFormat:
		unsetField = &sessionParametersUnset.GeometryOutputFormat
	case SessionParameterHybridTableLockTimeout:
		unsetField = &sessionParametersUnset.HybridTableLockTimeout
	case SessionParameterJdbcTreatDecimalAsInt:
		unsetField = &sessionParametersUnset.JdbcTreatDecimalAsInt
	case SessionParameterJdbcTreatTimestampNtzAsUtc:
		unsetField = &sessionParametersUnset.JdbcTreatTimestampNtzAsUtc
	case SessionParameterJdbcUseSessionTimezone:
		unsetField = &sessionParametersUnset.JdbcUseSessionTimezone
	case SessionParameterJsonIndent:
		unsetField = &sessionParametersUnset.JsonIndent
	case SessionParameterJsTreatIntegerAsBigInt:
		unsetField = &sessionParametersUnset.JsTreatIntegerAsBigInt
	case SessionParameterLockTimeout:
		unsetField = &sessionParametersUnset.LockTimeout
	case SessionParameterLogLevel:
		unsetField = &sessionParametersUnset.LogLevel
	case SessionParameterMultiStatementCount:
		unsetField = &sessionParametersUnset.MultiStatementCount
	case SessionParameterNoorderSequenceAsDefault:
		unsetField = &sessionParametersUnset.NoorderSequenceAsDefault
	case SessionParameterOdbcTreatDecimalAsInt:
		unsetField = &sessionParametersUnset.OdbcTreatDecimalAsInt
	case SessionParameterPythonProfilerModules:
		unsetField = &sessionParametersUnset.PythonProfilerModules
	case SessionParameterPythonProfilerTargetStage:
		unsetField = &sessionParametersUnset.PythonProfilerTargetStage
	case SessionParameterQueryTag:
		unsetField = &sessionParametersUnset.QueryTag
	case SessionParameterQuotedIdentifiersIgnoreCase:
		unsetField = &sessionParametersUnset.QuotedIdentifiersIgnoreCase
	case SessionParameterRowsPerResultset:
		unsetField = &sessionParametersUnset.RowsPerResultset
	case SessionParameterS3StageVpceDnsName:
		unsetField = &sessionParametersUnset.S3StageVpceDnsName
	case SessionParameterSearchPath:
		unsetField = &sessionParametersUnset.SearchPath
	case SessionParameterSimulatedDataSharingConsumer:
		unsetField = &sessionParametersUnset.SimulatedDataSharingConsumer
	case SessionParameterStatementQueuedTimeoutInSeconds:
		unsetField = &sessionParametersUnset.StatementQueuedTimeoutInSeconds
	case SessionParameterStatementTimeoutInSeconds:
		unsetField = &sessionParametersUnset.StatementTimeoutInSeconds
	case SessionParameterStrictJsonOutput:
		unsetField = &sessionParametersUnset.StrictJsonOutput
	case SessionParameterTimestampDayIsAlways24h:
		unsetField = &sessionParametersUnset.TimestampDayIsAlways24h
	case SessionParameterTimestampInputFormat:
		unsetField = &sessionParametersUnset.TimestampInputFormat
	case SessionParameterTimestampLTZOutputFormat:
		unsetField = &sessionParametersUnset.TimestampLTZOutputFormat
	case SessionParameterTimestampNTZOutputFormat:
		unsetField = &sessionParametersUnset.TimestampNTZOutputFormat
	case SessionParameterTimestampOutputFormat:
		unsetField = &sessionParametersUnset.TimestampOutputFormat
	case SessionParameterTimestampTypeMapping:
		unsetField = &sessionParametersUnset.TimestampTypeMapping
	case SessionParameterTimestampTZOutputFormat:
		unsetField = &sessionParametersUnset.TimestampTZOutputFormat
	case SessionParameterTimezone:
		unsetField = &sessionParametersUnset.Timezone
	case SessionParameterTimeInputFormat:
		unsetField = &sessionParametersUnset.TimeInputFormat
	case SessionParameterTimeOutputFormat:
		unsetField = &sessionParametersUnset.TimeOutputFormat
	case SessionParameterTraceLevel:
		unsetField = &sessionParametersUnset.TraceLevel
	case SessionParameterTransactionAbortOnError:
		unsetField = &sessionParametersUnset.TransactionAbortOnError
	case SessionParameterTransactionDefaultIsolationLevel:
		unsetField = &sessionParametersUnset.TransactionDefaultIsolationLevel
	case SessionParameterTwoDigitCenturyStart:
		unsetField = &sessionParametersUnset.TwoDigitCenturyStart
	case SessionParameterUnsupportedDDLAction:
		unsetField = &sessionParametersUnset.UnsupportedDDLAction
	case SessionParameterUseCachedResult:
		unsetField = &sessionParametersUnset.UseCachedResult
	case SessionParameterWeekOfYearPolicy:
		unsetField = &sessionParametersUnset.WeekOfYearPolicy
	case SessionParameterWeekStart:
		unsetField = &sessionParametersUnset.WeekStart
	default:
		return fmt.Errorf("%s session parameter is not supported", string(parameter))
	}
	*unsetField = Bool(true)
	return nil
}

func (legacyAccountParameters *LegacyAccountParameters) setParam(parameter AccountParameter, value string) (matched bool, err error) {
	matched = true
	switch parameter {
	case AccountParameterAllowClientMFACaching:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.AllowClientMFACaching = b
	case AccountParameterAllowIDToken:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.AllowIDToken = b
	case AccountParameterClientEncryptionKeySize:
		v, err := strconv.Atoi(value)
		if err != nil {
			return matched, fmt.Errorf("CLIENT_ENCRYPTION_KEY_SIZE session parameter is an integer, got %v", value)
		}
		legacyAccountParameters.ClientEncryptionKeySize = Pointer(v)
	case AccountParameterCortexEnabledCrossRegion:
		legacyAccountParameters.CortexEnabledCrossRegion = &value
	case AccountParameterDisableUserPrivilegeGrants:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.DisableUserPrivilegeGrants = b
	case AccountParameterEnableIdentifierFirstLogin:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.EnableIdentifierFirstLogin = b
	case AccountParameterEnableInternalStagesPrivatelink:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.EnableInternalStagesPrivatelink = b
	case AccountParameterEnablePersonalDatabase:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.EnablePersonalDatabase = b
	case AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.EnableTriSecretAndRekeyOptOutForImageRepository = b
	case AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.EnableTriSecretAndRekeyOptOutForSpcsBlockStorage = b
	case AccountParameterEnableUnhandledExceptionsReporting:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.EnableUnhandledExceptionsReporting = b
	case AccountParameterEnableUnredactedQuerySyntaxError:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.EnableUnredactedQuerySyntaxError = b
	case AccountParameterEnforceNetworkRulesForInternalStages:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.EnforceNetworkRulesForInternalStages = b
	case AccountParameterEventTable:
		legacyAccountParameters.EventTable = &value
	case AccountParameterExternalOAuthAddPrivilegedRolesToBlockedList:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.ExternalOAuthAddPrivilegedRolesToBlockedList = b
	case AccountParameterInitialReplicationSizeLimitInTB:
		legacyAccountParameters.InitialReplicationSizeLimitInTB = &value
	case AccountParameterMetricLevel:
		legacyAccountParameters.MetricLevel = Pointer(MetricLevel(value))
	case AccountParameterMinDataRetentionTimeInDays:
		v, err := strconv.Atoi(value)
		if err != nil {
			return matched, fmt.Errorf("MIN_DATA_RETENTION_TIME_IN_DAYS session parameter is an integer, got %v", value)
		}
		legacyAccountParameters.MinDataRetentionTimeInDays = Pointer(v)
	case AccountParameterNetworkPolicy:
		legacyAccountParameters.NetworkPolicy = &value
	case AccountParameterOAuthAddPrivilegedRolesToBlockedList:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.OAuthAddPrivilegedRolesToBlockedList = b
	case AccountParameterPeriodicDataRekeying:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.PeriodicDataRekeying = b
	case AccountParameterPreventLoadFromInlineURL:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.PreventLoadFromInlineURL = b
	case AccountParameterPreventUnloadToInlineURL:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.PreventUnloadToInlineURL = b
	case AccountParameterPreventUnloadToInternalStages:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.PreventUnloadToInternalStages = b
	case AccountParameterRequireStorageIntegrationForStageCreation:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.RequireStorageIntegrationForStageCreation = b
	case AccountParameterRequireStorageIntegrationForStageOperation:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.RequireStorageIntegrationForStageOperation = b
	case AccountParameterSsoLoginPage:
		var b *bool
		b, err = parseBooleanParameter(string(parameter), value)
		if err != nil {
			return
		}
		legacyAccountParameters.SSOLoginPage = b
	default:
		matched = false
	}
	return matched, nil
}
