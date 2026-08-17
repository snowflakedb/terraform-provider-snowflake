package objectparametersassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasAllDefaultsForEnvironment is an environment-aware variant of HasAllDefaults.
// It replaces HasAllDefaults() in tests that run on preprod, where QUOTED_IDENTIFIERS_IGNORE_CASE
// and STATEMENT_TIMEOUT_IN_SECONDS are set at ACCOUNT level rather than the Snowflake default level.
func (u *UserParametersAssert) HasAllDefaultsForEnvironment(t *testing.T, defaults *helpers.SnowflakeDefaultsClient) *UserParametersAssert {
	t.Helper()
	return u.
		HasDefaultParameterValueOnLevel(sdk.UserParameterEnableUnredactedQuerySyntaxError, sdk.ParameterTypeSnowflakeDefault).
		HasNetworkPolicy("RESTRICTED_ACCESS").
		HasNetworkPolicyLevel(sdk.ParameterTypeAccount).
		HasDefaultParameterValueOnLevel(sdk.UserParameterPreventUnloadToInternalStages, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterAbortDetachedQuery, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterAutocommit, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterBinaryInputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterBinaryOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientMemoryLimit, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientMetadataRequestUseConnectionCtx, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientPrefetchThreads, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientResultChunkSize, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientResultColumnCaseInsensitive, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientSessionKeepAlive, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterClientTimestampTypeMapping, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterDateInputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterDateOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterEnableUnloadPhysicalTypeOptimization, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterErrorOnNondeterministicMerge, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterErrorOnNondeterministicUpdate, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterGeographyOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterGeometryOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterJdbcTreatDecimalAsInt, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterJdbcTreatTimestampNtzAsUtc, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterJdbcUseSessionTimezone, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterJsonIndent, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterLockTimeout, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterLogLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterLogEventLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterMultiStatementCount, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterNoorderSequenceAsDefault, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterOdbcTreatDecimalAsInt, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterQueryTag, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterQuotedIdentifiersIgnoreCase, defaults.DefaultQuotedIdentifiersIgnoreCaseLevel(t)).
		HasDefaultParameterValueOnLevel(sdk.UserParameterRowsPerResultset, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterS3StageVpceDnsName, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterSearchPath, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterSimulatedDataSharingConsumer, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterStatementQueuedTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterStatementTimeoutInSeconds, defaults.DefaultStatementTimeoutInSecondsLevel(t)).
		HasDefaultParameterValueOnLevel(sdk.UserParameterStrictJsonOutput, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampDayIsAlways24h, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampInputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampLtzOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampNtzOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampTypeMapping, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimestampTzOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimezone, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimeInputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTimeOutputFormat, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTraceLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTransactionAbortOnError, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTransactionDefaultIsolationLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterTwoDigitCenturyStart, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterUnsupportedDdlAction, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterUseCachedResult, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterWeekOfYearPolicy, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.UserParameterWeekStart, sdk.ParameterTypeSnowflakeDefault)
}
