package gen

import (
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvidentifiers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
)

type SnowflakeObjectParameters struct {
	Name                    string
	IdType                  string
	Level                   sdk.ParameterType
	Parameters              []SnowflakeParameter
	ParameterConstantPrefix string
	ObjectTypeName          string
	IgnoreIdInProvider      bool
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
		// Adjust defaults according to our testing environments: override NETWORK_POLICY underneath.
		Parameters: append(
			slices.DeleteFunc(snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelUser)), func(p SnowflakeParameter) bool {
				return p.ParameterName == string(sdk.UserParameterNetworkPolicy)
			}),
			SnowflakeParameter{ParameterName: string(sdk.UserParameterNetworkPolicy), ParameterType: "string", DefaultValue: testenvidentifiers.NetworkPolicy.Name(), DefaultLevel: "sdk.ParameterTypeAccount"},
		),
	},
	{
		Name:       "Warehouse",
		IdType:     "sdk.AccountObjectIdentifier",
		Level:      sdk.ParameterTypeWarehouse,
		Parameters: snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouse)),
	},
	{
		Name:                    "WarehouseAdaptive",
		IdType:                  "sdk.AccountObjectIdentifier",
		Level:                   sdk.ParameterTypeWarehouse,
		ParameterConstantPrefix: "Warehouse",
		ObjectTypeName:          "Warehouse",
		Parameters:              snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouseAdaptive)),
	},
	{
		Name:                    "WarehouseInteractive",
		IdType:                  "sdk.AccountObjectIdentifier",
		Level:                   sdk.ParameterTypeWarehouse,
		ParameterConstantPrefix: "Warehouse",
		ObjectTypeName:          "Warehouse",
		Parameters:              snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelWarehouseInteractive)),
	},
	{
		Name:       "Database",
		IdType:     "sdk.AccountObjectIdentifier",
		Level:      sdk.ParameterTypeDatabase,
		Parameters: snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelDatabase)),
	},
	{
		Name:   "IcebergTable",
		IdType: "sdk.SchemaObjectIdentifier",
		Level:  sdk.ParameterTypeObject,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.IcebergTableParameterAllowRowTimestamp), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterCatalog), ParameterType: "string", DefaultValue: "SNOWFLAKE", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterCatalogSync), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterDataMetricSchedule), ParameterType: "string", DefaultValue: "60 MINUTES", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterDataRetentionTimeInDays), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeDatabase"},
			{ParameterName: string(sdk.IcebergTableParameterDefaultDdlCollation), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterEnableDataCompaction), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterEnableIcebergMergeOnRead), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterExternalVolume), ParameterType: "string", DefaultValue: "SNOWFLAKE_MANAGED", DefaultLevel: "sdk.ParameterTypeTable"},
			{ParameterName: string(sdk.IcebergTableParameterIcebergMergeOnReadBehavior), ParameterType: "string", DefaultValue: "auto", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterLogEventLevel), ParameterType: "string", DefaultValue: "OFF", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterMaxDataExtensionTimeInDays), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeDatabase"},
			{ParameterName: string(sdk.IcebergTableParameterOptimizeDataLayout), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterQuotedIdentifiersIgnoreCase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterReplaceInvalidCharacters), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.IcebergTableParameterStorageSerializationPolicy), ParameterType: "sdk.StorageSerializationPolicy", DefaultValue: "sdk.StorageSerializationPolicyOptimized", DefaultLevel: "sdk.ParameterTypeTable"},
			{ParameterName: string(sdk.IcebergTableParameterTargetFileSize), ParameterType: "sdk.IcebergTableTargetFileSize", DefaultValue: "sdk.IcebergTableTargetFileSizeAuto", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:   "Task",
		IdType: "sdk.SchemaObjectIdentifier",
		Level:  sdk.ParameterTypeTask,
		Parameters: append(snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelTask)), SnowflakeParameter{
			ParameterName: string(sdk.TaskParameterSearchPath), ParameterType: "string", DefaultValue: "$current, $public", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault",
		}),
	},
	{
		Name:               "Account",
		IdType:             "sdk.AccountIdentifier",
		Level:              sdk.ParameterTypeAccount,
		IgnoreIdInProvider: true,
		Parameters: []SnowflakeParameter{
			// Bool parameters
			{ParameterName: string(sdk.AccountParameterAbortDetachedQuery), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterAllowBindValuesAccess), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterAllowClientMfaCaching), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterAllowIdToken), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterAutocommit), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientEnableLogInfoStatementParameters), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientMetadataRequestUseConnectionCtx), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientMetadataUseSessionDatabase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientResultColumnCaseInsensitive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientSessionKeepAlive), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDisallowedSpcsWorkloadTypes), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDisableUiDownloadButton), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDisableUserPrivilegeGrants), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableAutomaticSensitiveDataClassificationLog), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableBudgetEventLogging), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableDataCompaction), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableEgressCostOptimizer), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableGetDdlUseDataTypeAlias), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableIcebergMergeOnRead), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableNotebookCreationInPersonalDb), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableSpcsBlockStorageSnowflakeFullEncryptionEnforcement), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableTagPropagationEventLogging), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableIdentifierFirstLogin), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableInternalStagesPrivatelink), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableTriSecretAndRekeyOptOutForImageRepository), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableTriSecretAndRekeyOptOutForSpcsBlockStorage), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableUnhandledExceptionsReporting), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableUnloadPhysicalTypeOptimization), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableUnredactedQuerySyntaxError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnableUnredactedSecureObjectError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEnforceNetworkRulesForInternalStages), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterErrorOnNondeterministicMerge), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterErrorOnNondeterministicUpdate), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterExternalOauthAddPrivilegedRolesToBlockedList), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterJdbcTreatDecimalAsInt), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterJdbcTreatTimestampNtzAsUtc), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterJdbcUseSessionTimezone), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterJsTreatIntegerAsBigint), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterNoorderSequenceAsDefault), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterOauthAddPrivilegedRolesToBlockedList), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterOdbcTreatDecimalAsInt), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterPeriodicDataRekeying), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterPipeExecutionPaused), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterPreventUnloadToInlineUrl), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterPreventUnloadToInternalStages), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterQuotedIdentifiersIgnoreCase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterRowTimestampDefault), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterReplaceInvalidCharacters), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterRequireStorageIntegrationForStageCreation), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterRequireStorageIntegrationForStageOperation), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterSsoLoginPage), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterStrictJsonOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimestampDayIsAlways24h), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTransactionAbortOnError), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterUseCachedResult), ParameterType: "bool", DefaultValue: "true", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},

			// Int parameters
			{ParameterName: string(sdk.AccountParameterClientEncryptionKeySize), ParameterType: "int", DefaultValue: "128", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientMemoryLimit), ParameterType: "int", DefaultValue: "1536", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientPrefetchThreads), ParameterType: "int", DefaultValue: "4", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientResultChunkSize), ParameterType: "int", DefaultValue: "160", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientSessionKeepAliveHeartbeatFrequency), ParameterType: "int", DefaultValue: "3600", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDataRetentionTimeInDays), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterHybridTableLockTimeout), ParameterType: "int", DefaultValue: "3600", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterIcebergVersionDefault), ParameterType: "int", DefaultValue: "2", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterJsonIndent), ParameterType: "int", DefaultValue: "2", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterLockTimeout), ParameterType: "int", DefaultValue: "43200", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterMaxConcurrencyLevel), ParameterType: "int", DefaultValue: "14", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterMaxDataExtensionTimeInDays), ParameterType: "int", DefaultValue: "14", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterMinDataRetentionTimeInDays), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterMultiStatementCount), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterRowsPerResultset), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterStatementQueuedTimeoutInSeconds), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterStatementTimeoutInSeconds), ParameterType: "int", DefaultValue: "172800", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterSuspendTaskAfterNumFailures), ParameterType: "int", DefaultValue: "10", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTaskAutoRetryAttempts), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTwoDigitCenturyStart), ParameterType: "int", DefaultValue: "1970", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterUserTaskMinimumTriggerIntervalInSeconds), ParameterType: "int", DefaultValue: "30", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterUserTaskTimeoutMs), ParameterType: "int", DefaultValue: "3600000", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterWeekOfYearPolicy), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterWeekStart), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},

			// String parameters
			{ParameterName: string(sdk.AccountParameterActivePythonProfiler), ParameterType: "sdk.ActivePythonProfiler", DefaultValue: "sdk.ActivePythonProfilerLine", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterAllowedSpcsWorkloadTypes), ParameterType: "string", DefaultValue: "ALL", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterBaseLocationPrefix), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterBinaryInputFormat), ParameterType: "sdk.BinaryInputFormat", DefaultValue: "sdk.BinaryInputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterBinaryOutputFormat), ParameterType: "sdk.BinaryOutputFormat", DefaultValue: "sdk.BinaryOutputFormatHex", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterCatalog), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterCatalogSync), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterClientTimestampTypeMapping), ParameterType: "sdk.ClientTimestampTypeMapping", DefaultValue: "sdk.ClientTimestampTypeMappingLtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterCortexCodeCliDailyEstCreditLimitPerUser), ParameterType: "int", DefaultValue: "-1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterCortexCodeDesktopDailyEstCreditLimitPerUser), ParameterType: "int", DefaultValue: "-1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterCortexCodeSnowsightDailyEstCreditLimitPerUser), ParameterType: "int", DefaultValue: "-1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterCortexEnabledCrossRegion), ParameterType: "string", DefaultValue: "DISABLED", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterCortexModelsAllowlist), ParameterType: "string", DefaultValue: "ALL", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterCsvTimestampFormat), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDataMetricSchedule), ParameterType: "string", DefaultValue: "60 MINUTES", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDateInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDateOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDefaultDbtVersion), ParameterType: "string", DefaultValue: "1.9.4", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDefaultDdlCollation), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDefaultNotebookComputePoolCpu), ParameterType: "string", DefaultValue: "SYSTEM_COMPUTE_POOL_CPU", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDefaultNotebookComputePoolGpu), ParameterType: "string", DefaultValue: "SYSTEM_COMPUTE_POOL_GPU", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDefaultNullOrdering), ParameterType: "sdk.DefaultNullOrdering", DefaultValue: "sdk.DefaultNullOrderingLast", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDefaultStreamlitComputePool), ParameterType: "string", DefaultValue: "SYSTEM_COMPUTE_POOL_CPU", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterDefaultStreamlitNotebookWarehouse), ParameterType: "string", DefaultValue: "REGRESS", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterEventTable), ParameterType: "string", DefaultValue: "snowflake.telemetry.events", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterExternalVolume), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterGeographyOutputFormat), ParameterType: "sdk.GeographyOutputFormat", DefaultValue: "sdk.GeographyOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterGeometryOutputFormat), ParameterType: "sdk.GeometryOutputFormat", DefaultValue: "sdk.GeometryOutputFormatGeoJSON", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterInitialReplicationSizeLimitInTb), ParameterType: "string", DefaultValue: "10", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterListingAutoFulfillmentReplicationRefreshSchedule), ParameterType: "string", DefaultValue: "1440 MINUTE", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterLogEventLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterMetricLevel), ParameterType: "sdk.MetricLevel", DefaultValue: "sdk.MetricLevelNone", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterNetworkPolicy), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterPythonProfilerModules), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterPythonProfilerTargetStage), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterQueryTag), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterReadConsistencyMode), ParameterType: "string", DefaultValue: "SESSION", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterS3StageVpceDnsName), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterSearchPath), ParameterType: "string", DefaultValue: "$current, $public", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterServerlessTaskMaxStatementSize), ParameterType: "sdk.WarehouseSize", DefaultValue: `sdk.WarehouseSize("X2Large")`, DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterServerlessTaskMinStatementSize), ParameterType: "sdk.WarehouseSize", DefaultValue: "sdk.WarehouseSizeXSmall", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterSimulatedDataSharingConsumer), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterSqlTraceQueryText), ParameterType: "string", DefaultValue: "OFF", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterStorageSerializationPolicy), ParameterType: "sdk.StorageSerializationPolicy", DefaultValue: "sdk.StorageSerializationPolicyOptimized", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimestampInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimestampLtzOutputFormat), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimestampNtzOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimestampOutputFormat), ParameterType: "string", DefaultValue: "YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimestampTypeMapping), ParameterType: "sdk.TimestampTypeMapping", DefaultValue: "sdk.TimestampTypeMappingNtz", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimestampTzOutputFormat), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimezone), ParameterType: "string", DefaultValue: "America/Los_Angeles", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimeInputFormat), ParameterType: "string", DefaultValue: "AUTO", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTimeOutputFormat), ParameterType: "string", DefaultValue: "HH24:MI:SS", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterTransactionDefaultIsolationLevel), ParameterType: "string", DefaultValue: "READ COMMITTED", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterUnsupportedDdlAction), ParameterType: "string", DefaultValue: "ignore", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterUserTaskManagedInitialWarehouseSize), ParameterType: "sdk.WarehouseSize", DefaultValue: "sdk.WarehouseSizeMedium", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.AccountParameterUseWorkspacesForSql), ParameterType: "string", DefaultValue: "unset", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:   "Function",
		IdType: "sdk.SchemaObjectIdentifierWithArguments",
		Level:  sdk.ParameterTypeFunction,
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.FunctionParameterEnableConsoleOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.FunctionParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.FunctionParameterLogEventLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.FunctionParameterMetricLevel), ParameterType: "sdk.MetricLevel", DefaultValue: "sdk.MetricLevelNone", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.FunctionParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:       "Procedure",
		IdType:     "sdk.SchemaObjectIdentifierWithArguments",
		Level:      sdk.ParameterTypeProcedure,
		Parameters: snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelProcedure)),
	},
	{
		Name:       "Service",
		IdType:     "sdk.SchemaObjectIdentifier",
		Level:      sdk.ParameterTypeService,
		Parameters: snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelService)),
	},
	{
		Name:                    "HybridTable",
		IdType:                  "sdk.SchemaObjectIdentifier",
		Level:                   sdk.ParameterTypeObject,
		ParameterConstantPrefix: "Object",
		Parameters:              snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelHybridTable)),
	},
	{
		Name:                    "Schema",
		IdType:                  "sdk.DatabaseObjectIdentifier",
		Level:                   sdk.ParameterTypeObject,
		ParameterConstantPrefix: "Object",
		Parameters: []SnowflakeParameter{
			{ParameterName: string(sdk.ObjectParameterCatalog), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterDataRetentionTimeInDays), ParameterType: "int", DefaultValue: "1", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterDefaultDdlCollation), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterDefaultNotebookComputePoolCpu), ParameterType: "string", DefaultValue: "SYSTEM_COMPUTE_POOL_CPU", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterDefaultNotebookComputePoolGpu), ParameterType: "string", DefaultValue: "SYSTEM_COMPUTE_POOL_GPU", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterEnableConsoleOutput), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterExternalVolume), ParameterType: "string", DefaultValue: "", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterLogEventLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterLogLevel), ParameterType: "sdk.LogLevel", DefaultValue: "sdk.LogLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterMaxDataExtensionTimeInDays), ParameterType: "int", DefaultValue: "14", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterPipeExecutionPaused), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterQuotedIdentifiersIgnoreCase), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterReplaceInvalidCharacters), ParameterType: "bool", DefaultValue: "false", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterStorageSerializationPolicy), ParameterType: "sdk.StorageSerializationPolicy", DefaultValue: "sdk.StorageSerializationPolicyOptimized", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterSuspendTaskAfterNumFailures), ParameterType: "int", DefaultValue: "10", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterTaskAutoRetryAttempts), ParameterType: "int", DefaultValue: "0", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterTraceLevel), ParameterType: "sdk.TraceLevel", DefaultValue: "sdk.TraceLevelOff", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterUserTaskManagedInitialWarehouseSize), ParameterType: "sdk.WarehouseSize", DefaultValue: `sdk.WarehouseSize("Medium")`, DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds), ParameterType: "int", DefaultValue: "30", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
			{ParameterName: string(sdk.ObjectParameterUserTaskTimeoutMs), ParameterType: "int", DefaultValue: "3600000", DefaultLevel: "sdk.ParameterTypeSnowflakeDefault"},
		},
	},
	{
		Name:   "OpenflowDeployment",
		IdType: "sdk.AccountObjectIdentifier",
		Level:  sdk.ParameterTypeOpenflowDeployment,
		// EVENT_TABLE has no fixed default. An account that sets its own EVENT_TABLE hands that value
		// down, so what a deployment reports before anything is set differs between accounts and only
		// the level tells the two apart. Assert HasEventTable and HasEventTableLevel rather than the
		// generated defaults assertions.
		Parameters: snowflakeParameters(defs.ParameterDefsForLevel(parameterdefs.ParameterLevelOpenflowDeployment)),
	},
}

func snowflakeParameters(parameters []parameterdefs.ParameterDef) []SnowflakeParameter {
	out := make([]SnowflakeParameter, 0, len(parameters))
	for _, parameter := range parameters {
		out = append(out, SnowflakeParameter{
			ParameterName: parameter.SqlName,
			ParameterType: defs.AssertionParameterType(parameter.Kind),
			DefaultValue:  parameter.DefaultValue,
			DefaultLevel:  parameter.DefaultLevel,
		})
	}
	return out
}
