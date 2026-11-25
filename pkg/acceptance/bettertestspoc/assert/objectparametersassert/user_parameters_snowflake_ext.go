package objectparametersassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvidentifiers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasAllDefaultsForEnvironment is an environment-aware variant of HasAllDefaults.
// It replaces HasAllDefaults() in tests that run on preprod, where QUOTED_IDENTIFIERS_IGNORE_CASE
// and STATEMENT_TIMEOUT_IN_SECONDS are set at ACCOUNT level rather than the Snowflake default level.
// SYSTEM is accepted in place of the Snowflake default level.
func (u *UserParametersAssert) HasAllDefaultsForEnvironment(t *testing.T, defaults *helpers.SnowflakeDefaultsClient) *UserParametersAssert {
	t.Helper()
	return u.
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterEnableUnredactedQuerySyntaxError).
		HasNetworkPolicy(testenvidentifiers.NetworkPolicy.Name()).
		HasNetworkPolicyLevel(sdk.ParameterTypeAccount).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterPreventUnloadToInternalStages).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterAbortDetachedQuery).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterAutocommit).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterBinaryInputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterBinaryOutputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterClientMemoryLimit).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterClientMetadataRequestUseConnectionCtx).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterClientPrefetchThreads).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterClientResultChunkSize).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterClientResultColumnCaseInsensitive).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterClientSessionKeepAlive).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterClientTimestampTypeMapping).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterDateInputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterDateOutputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterEnableUnloadPhysicalTypeOptimization).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterErrorOnNondeterministicMerge).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterErrorOnNondeterministicUpdate).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterGeographyOutputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterGeometryOutputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterJdbcTreatDecimalAsInt).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterJdbcTreatTimestampNtzAsUtc).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterJdbcUseSessionTimezone).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterJsonIndent).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterLockTimeout).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterLogLevel).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterLogEventLevel).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterMultiStatementCount).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterNoorderSequenceAsDefault).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterOdbcTreatDecimalAsInt).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterQueryTag).
		HasDefaultParameterValueOnLevel(sdk.UserParameterQuotedIdentifiersIgnoreCase, defaults.DefaultQuotedIdentifiersIgnoreCaseLevel(t)).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterRowsPerResultset).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterS3StageVpceDnsName).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterSearchPath).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterSimulatedDataSharingConsumer).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterStatementQueuedTimeoutInSeconds).
		HasDefaultParameterValueOnLevel(sdk.UserParameterStatementTimeoutInSeconds, defaults.DefaultStatementTimeoutInSecondsLevel(t)).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterStrictJsonOutput).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimestampDayIsAlways24h).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimestampInputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimestampLtzOutputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimestampNtzOutputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimestampOutputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimestampTypeMapping).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimestampTzOutputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimezone).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimeInputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTimeOutputFormat).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTraceLevel).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTransactionAbortOnError).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTransactionDefaultIsolationLevel).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterTwoDigitCenturyStart).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterUnsupportedDdlAction).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterUseCachedResult).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterWeekOfYearPolicy).
		HasDefaultParameterValueOnDefaultLevel(sdk.UserParameterWeekStart)
}
