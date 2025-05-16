package gen

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvidentifiers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

type SnowflakeObjectParameters struct {
	Name              string
	IdType            string
	Level             sdk.ParameterType
	AdditionalImports []string
	Parameters        []SnowflakeParameter
}

func (p SnowflakeObjectParameters) ObjectName() string {
	return p.Name
}

type SnowflakeParameter struct {
	ParameterName string
	ParameterType string
	DefaultValue  string
	DefaultLevel  string
}

func GetAllSnowflakeObjectParameters() []SnowflakeObjectParameters {
	return allObjectsParameters
}

// TODO [SNOW-1501905]: use SDK definition after parameters rework (+ preprocessing here)
var allObjectsParameters = []SnowflakeObjectParameters{
	{
		Name:   "User",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeUser,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.UserParameterEnableUnredactedQuerySyntaxError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterNetworkPolicy), ParameterType: "string", DefaultValue: testenvidentifiers.NetworkPolicy.Name(), DefaultLevel: "sdk.ParameterTypeAccount"},
			{ParameterName: string(sdk.UserParameterPreventUnloadToInternalStages), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterAbortDetachedQuery), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterAutocommit), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterBinaryInputFormat), ParameterType: "sdk.BinaryInputFormat", DefaultValue: "sdk.BinaryInputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterBinaryOutputFormat), ParameterType: "sdk.BinaryOutputFormat", DefaultValue: "sdk.BinaryOutputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientMemoryLimit), ParameterType: "int", DefaultValue: "1536", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientMetadataRequestUseConnectionCtx), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientPrefetchThreads), ParameterType: "int", DefaultValue: "4", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientResultChunkSize), ParameterType: "int", DefaultValue: "160", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientResultColumnCaseInsensitive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientSessionKeepAlive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency), ParameterType: "int", DefaultValue: "3600", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterClientTimestampTypeMapping), ParameterType: "sdk.ClientTimestampTypeMapping", DefaultValue: "sdk.ClientTimestampTypeMappingLtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterDateInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterDateOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterEnableUnloadPhysicalTypeOptimization), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterErrorOnNondeterministicMerge), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterErrorOnNondeterministicUpdate), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterGeographyOutputFormat), ParameterType: "sdk.GeographyOutputFormat", DefaultValue: "sdk.GeographyOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterGeometryOutputFormat), ParameterType: "sdk.GeometryOutputFormat", DefaultValue: "sdk.GeometryOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterJdbcTreatDecimalAsInt), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterJdbcTreatTimestampNtzAsUtc), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterJdbcUseSessionTimezone), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterJsonIndent), ParameterType: "int", DefaultValue: "2", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterLockTimeout), ParameterType: "int", DefaultValue: "43200", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterMultiStatementCount), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterNoorderSequenceAsDefault), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterOdbcTreatDecimalAsInt), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterQueryTag), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterQuotedIdentifiersIgnoreCase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterRowsPerResultset), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterS3StageVpceDnsName), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterSearchPath), ParameterType: "string", DefaultValue: "$current, $public", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterSimulatedDataSharingConsumer), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterStatementQueuedTimeoutInSeconds), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterStatementTimeoutInSeconds), ParameterType: "int", DefaultValue: "172800", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterStrictJsonOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampDayIsAlways24h), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampLtzOutputFormat), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampNtzOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampTypeMapping), ParameterType: "sdk.TimestampTypeMapping", DefaultValue: "sdk.TimestampTypeMappingNtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimestampTzOutputFormat), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimezone), ParameterType: "string", DefaultValue: "America/Los_Angeles", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimeInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTimeOutputFormat), ParameterType: "string", DefaultValue: "HH24:MI:SS", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTransactionAbortOnError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTransactionDefaultIsolationLevel), ParameterType: "sdk.TransactionDefaultIsolationLevel", DefaultValue: "sdk.TransactionDefaultIsolationLevelReadCommitted", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterTwoDigitCenturyStart), ParameterType: "int", DefaultValue: "1970", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// TODO [SNOW-1501905]: quick workaround for now: lowercase for ignore in snowflake by default but uppercase for FAIL
			{ParameterName: string(sdk.UserParameterUnsupportedDdlAction), ParameterType: "sdk.UnsupportedDDLAction", DefaultValue: "sdk.UnsupportedDDLAction(strings.ToLower(string(sdk.UnsupportedDDLActionIgnore)))", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterUseCachedResult), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterWeekOfYearPolicy), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.UserParameterWeekStart), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
		// for the quickfix above
		AdditionalImports: []string{"strings"},
	},
	{
		Name:   "Warehouse",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeWarehouse,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.WarehouseParameterMaxConcurrencyLevel), ParameterType: "int", DefaultValue: "8", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.WarehouseParameterStatementTimeoutInSeconds), ParameterType: "int", DefaultValue: "172800", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:   "Database",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeDatabase,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.DatabaseParameterDataRetentionTimeInDays), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterMaxDataExtensionTimeInDays), ParameterType: "int", DefaultValue: "14", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterExternalVolume), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterCatalog), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterReplaceInvalidCharacters), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterDefaultDdlCollation), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterStorageSerializationPolicy), ParameterType: "sdk.StorageSerializationPolicy", DefaultValue: "sdk.StorageSerializationPolicyOptimized", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterSuspendTaskAfterNumFailures), ParameterType: "int", DefaultValue: "10", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterTaskAutoRetryAttempts), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterUserTaskManagedInitialWarehouseSize), ParameterType: "sdk.WarehouseSize", DefaultValue: "sdk.WarehouseSizeMedium", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterUserTaskTimeoutMs), ParameterType: "int", DefaultValue: "3600000", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterUserTaskMinimumTriggerIntervalInSeconds), ParameterType: "int", DefaultValue: "30", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterQuotedIdentifiersIgnoreCase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.DatabaseParameterEnableConsoleOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:   "Task",
		IdType: "sdk.SchemaObjectIdentifier",
		Level:  sdk.ParameterTypeTask,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.TaskParameterSuspendTaskAfterNumFailures), ParameterType: "int", DefaultValue: "10", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTaskAutoRetryAttempts), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterUserTaskManagedInitialWarehouseSize), ParameterType: "sdk.WarehouseSize", DefaultValue: `sdk.WarehouseSize("Medium")`, DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterUserTaskMinimumTriggerIntervalInSeconds), ParameterType: "int", DefaultValue: "30", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterUserTaskTimeoutMs), ParameterType: "int", DefaultValue: "3600000", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterAbortDetachedQuery), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterAutocommit), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterBinaryInputFormat), ParameterType: "sdk.BinaryInputFormat", DefaultValue: "sdk.BinaryInputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterBinaryOutputFormat), ParameterType: "sdk.BinaryOutputFormat", DefaultValue: "sdk.BinaryOutputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientMemoryLimit), ParameterType: "int", DefaultValue: "1536", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientMetadataRequestUseConnectionCtx), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientPrefetchThreads), ParameterType: "int", DefaultValue: "4", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientResultChunkSize), ParameterType: "int", DefaultValue: "160", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientResultColumnCaseInsensitive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientSessionKeepAlive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientSessionKeepAliveHeartbeatFrequency), ParameterType: "int", DefaultValue: "3600", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterClientTimestampTypeMapping), ParameterType: "sdk.ClientTimestampTypeMapping", DefaultValue: "sdk.ClientTimestampTypeMappingLtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterDateInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterDateOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterEnableUnloadPhysicalTypeOptimization), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterErrorOnNondeterministicMerge), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterErrorOnNondeterministicUpdate), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterGeographyOutputFormat), ParameterType: "sdk.GeographyOutputFormat", DefaultValue: "sdk.GeographyOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterGeometryOutputFormat), ParameterType: "sdk.GeometryOutputFormat", DefaultValue: "sdk.GeometryOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterJdbcTreatTimestampNtzAsUtc), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterJdbcUseSessionTimezone), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterJsonIndent), ParameterType: "int", DefaultValue: "2", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterLockTimeout), ParameterType: "int", DefaultValue: "43200", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterMultiStatementCount), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterNoorderSequenceAsDefault), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterOdbcTreatDecimalAsInt), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterQueryTag), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterQuotedIdentifiersIgnoreCase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterRowsPerResultset), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterS3StageVpceDnsName), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterSearchPath), ParameterType: "string", DefaultValue: "$current, $public", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterStatementQueuedTimeoutInSeconds), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterStatementTimeoutInSeconds), ParameterType: "int", DefaultValue: "172800", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterStrictJsonOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampDayIsAlways24h), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampLtzOutputFormat), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampNtzOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampTypeMapping), ParameterType: "sdk.TimestampTypeMapping", DefaultValue: "sdk.TimestampTypeMappingNtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimestampTzOutputFormat), ParameterType: "string", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimezone), ParameterType: "string", DefaultValue: "America/Los_Angeles", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimeInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTimeOutputFormat), ParameterType: "string", DefaultValue: "HH24:MI:SS", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTransactionAbortOnError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTransactionDefaultIsolationLevel), ParameterType: "sdk.TransactionDefaultIsolationLevel", DefaultValue: "sdk.TransactionDefaultIsolationLevelReadCommitted", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterTwoDigitCenturyStart), ParameterType: "int", DefaultValue: "1970", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// TODO [SNOW-1501905]: quick workaround for now: lowercase for ignore in snowflake by default but uppercase for FAIL
			{ParameterName: string(sdk.TaskParameterUnsupportedDdlAction), ParameterType: "sdk.UnsupportedDDLAction", DefaultValue: "sdk.UnsupportedDDLAction(strings.ToLower(string(sdk.UnsupportedDDLActionIgnore)))", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterUseCachedResult), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterWeekOfYearPolicy), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.TaskParameterWeekStart), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
		// for the quickfix above
		AdditionalImports: []string{"strings"},
	},
	{
		Name:   "Account",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeAccount,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.UserParameterEnableUnredactedQuerySyntaxError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			// TODO(SNOW-1348092 - next prs): Add parameters
		},
	},
	{
		Name:   "Function",
		IdType: "sdk.SchemaObjectIdentifierWithArguments",
		Level:  sdk.ParameterTypeFunction,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.FunctionParameterEnableConsoleOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.FunctionParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.FunctionParameterMetricLevel), ParameterType: "sdk.MetricLevel", DefaultValue: "sdk.MetricLevelNone", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.FunctionParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:   "Procedure",
		IdType: "sdk.SchemaObjectIdentifierWithArguments",
		Level:  sdk.ParameterTypeProcedure,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.ProcedureParameterAutoEventLogging), ParameterType: "sdk.AutoEventLogging", DefaultValue: "sdk.AutoEventLoggingOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ProcedureParameterEnableConsoleOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ProcedureParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ProcedureParameterMetricLevel), ParameterType: "sdk.MetricLevel", DefaultValue: "sdk.MetricLevelNone", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ProcedureParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
}
