package sdk

func init() {
	sessionsTests.Alter.
		withModifyAndExpectedErr(
			case_Sessions_validation_Alter_Set_SessionParameters_ValidateValue,
			func(opts *AlterSessionOptions) {
				opts.Set = &SessionSet{SessionParameters: &SessionParameters{ClientPrefetchThreads: new(-1)}}
			},
			errIntValue("SessionParameters", "ClientPrefetchThreads", IntErrGreaterOrEqual, 0),
		).
		withModifyAndExpectedErr(
			case_Sessions_validation_Alter_Unset_SessionParametersUnset_ValidateValue,
			func(opts *AlterSessionOptions) {
				opts.Unset = &SessionUnset{SessionParametersUnset: &SessionParametersUnset{}}
			},
			errAtLeastOneOf("SessionParametersUnset", "AbortDetachedQuery", "ActivePythonProfiler", "Autocommit", "BinaryInputFormat", "BinaryOutputFormat", "ClientEnableLogInfoStatementParameters", "ClientMemoryLimit", "ClientMetadataRequestUseConnectionCtx", "ClientPrefetchThreads", "ClientResultChunkSize", "ClientResultColumnCaseInsensitive", "ClientMetadataUseSessionDatabase", "ClientSessionKeepAlive", "ClientSessionKeepAliveHeartbeatFrequency", "ClientTimestampTypeMapping", "CsvTimestampFormat", "DateInputFormat", "DateOutputFormat", "EnableCortexAnalyst", "EnableGetDdlUseDataTypeAlias", "EnableUnloadPhysicalTypeOptimization", "ErrorOnNondeterministicMerge", "ErrorOnNondeterministicUpdate", "GeographyOutputFormat", "GeometryOutputFormat", "HybridTableLockTimeout", "JdbcTreatDecimalAsInt", "JdbcTreatTimestampNtzAsUtc", "JdbcUseSessionTimezone", "JsonIndent", "JsTreatIntegerAsBigInt", "LockTimeout", "LogLevel", "LogEventLevel", "MultiStatementCount", "NoorderSequenceAsDefault", "OdbcTreatDecimalAsInt", "PythonProfilerModules", "PythonProfilerTargetStage", "QueryTag", "QuotedIdentifiersIgnoreCase", "RowsPerResultset", "S3StageVpceDnsName", "SearchPath", "SimulatedDataSharingConsumer", "StatementQueuedTimeoutInSeconds", "StatementTimeoutInSeconds", "StrictJsonOutput", "TimestampDayIsAlways24h", "TimestampInputFormat", "TimestampLTZOutputFormat", "TimestampNTZOutputFormat", "TimestampOutputFormat", "TimestampTypeMapping", "TimestampTZOutputFormat", "Timezone", "TimeInputFormat", "TimeOutputFormat", "TraceLevel", "TransactionAbortOnError", "TransactionDefaultIsolationLevel", "TwoDigitCenturyStart", "UnsupportedDDLAction", "UseCachedResult", "WeekOfYearPolicy", "WeekStart"),
		).
		withModifyAndExpectedSqlf(
			case_Sessions_sql_Alter_Set,
			func(opts *AlterSessionOptions) {
				opts.Set = &SessionSet{SessionParameters: &SessionParameters{
					AbortDetachedQuery:    new(true),
					ClientPrefetchThreads: new(5),
				}}
			},
			`ALTER SESSION SET ABORT_DETACHED_QUERY = true, CLIENT_PREFETCH_THREADS = 5`,
		).
		withModifyAndExpectedSqlf(
			case_Sessions_sql_Alter_Unset,
			func(opts *AlterSessionOptions) {
				opts.Unset = &SessionUnset{SessionParametersUnset: &SessionParametersUnset{
					AbortDetachedQuery:    new(true),
					ClientPrefetchThreads: new(true),
				}}
			},
			`ALTER SESSION UNSET ABORT_DETACHED_QUERY, CLIENT_PREFETCH_THREADS`,
		)

	warehouseId := sessionsTestIdAccountObjectIdentifier
	databaseId := sessionsTestIdAccountObjectIdentifier
	schemaId := sessionsTestIdDatabaseObjectIdentifier
	roleId := sessionsTestIdAccountObjectIdentifier

	sessionsTests.UseWarehouse.
		withExpectedSqlf(case_Sessions_sql_UseWarehouse_basic, `USE WAREHOUSE %s`, warehouseId.FullyQualifiedName())
	sessionsTests.UseDatabase.
		withExpectedSqlf(case_Sessions_sql_UseDatabase_basic, `USE DATABASE %s`, databaseId.FullyQualifiedName())
	sessionsTests.UseSchema.
		withExpectedSqlf(case_Sessions_sql_UseSchema_basic, `USE SCHEMA %s`, schemaId.FullyQualifiedName())
	sessionsTests.UseRole.
		withExpectedSqlf(case_Sessions_sql_UseRole_basic, `USE ROLE %s`, roleId.FullyQualifiedName())
	sessionsTests.UseSecondaryRoles.
		withModifyAndExpectedSqlf(
			case_Sessions_sql_UseSecondaryRoles_basic,
			func(opts *UseSecondaryRolesSessionOptions) { opts.Option = SecondaryRoleOptionAll },
			`USE SECONDARY ROLES ALL`,
		)
}
