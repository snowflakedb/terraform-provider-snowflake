package resourceparametersassert

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/testenvidentifiers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (u *UserResourceParametersAssert) HasAllDefaults() *UserResourceParametersAssert {
	return u.
		HasEnableUnredactedQuerySyntaxError(false).
		HasNetworkPolicy(testenvidentifiers.NetworkPolicy.Name()).
		HasPreventUnloadToInternalStages(false).
		HasAbortDetachedQuery(false).
		HasAutocommit(true).
		HasBinaryInputFormat(sdk.BinaryInputFormatHex).
		HasBinaryOutputFormat(sdk.BinaryOutputFormatHex).
		HasClientMemoryLimit(1536).
		HasClientMetadataRequestUseConnectionCtx(false).
		HasClientPrefetchThreads(4).
		HasClientResultChunkSize(160).
		HasClientResultColumnCaseInsensitive(false).
		HasClientSessionKeepAlive(false).
		HasClientSessionKeepAliveHeartbeatFrequency(3600).
		HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingLtz).
		HasDateInputFormat("AUTO").
		HasDateOutputFormat("YYYY-MM-DD").
		HasEnableUnloadPhysicalTypeOptimization(true).
		HasErrorOnNondeterministicMerge(true).
		HasErrorOnNondeterministicUpdate(false).
		HasGeographyOutputFormat(sdk.GeographyOutputFormatGeoJSON).
		HasGeometryOutputFormat(sdk.GeometryOutputFormatGeoJSON).
		HasJdbcTreatDecimalAsInt(true).
		HasJdbcTreatTimestampNtzAsUtc(false).
		HasJdbcUseSessionTimezone(true).
		HasJsonIndent(2).
		HasLockTimeout(43200).
		HasLogLevel(sdk.LogLevelOff).
		HasMultiStatementCount(1).
		HasNoorderSequenceAsDefault(true).
		HasOdbcTreatDecimalAsInt(false).
		HasQueryTag("").
		HasQuotedIdentifiersIgnoreCase(false).
		HasRowsPerResultset(0).
		HasS3StageVpceDnsName("").
		HasSearchPath("$current, $public").
		HasSimulatedDataSharingConsumer("").
		HasStatementQueuedTimeoutInSeconds(0).
		HasStatementTimeoutInSeconds(172800).
		HasStrictJsonOutput(false).
		HasTimestampDayIsAlways24h(false).
		HasTimestampInputFormat("AUTO").
		HasTimestampLtzOutputFormat("").
		HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS.FF3").
		HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS.FF3 TZHTZM").
		HasTimestampTypeMapping(sdk.TimestampTypeMappingNtz).
		HasTimestampTzOutputFormat("").
		HasTimezone("America/Los_Angeles").
		HasTimeInputFormat("AUTO").
		HasTimeOutputFormat("HH24:MI:SS").
		HasTraceLevel(sdk.TraceLevelOff).
		HasTransactionAbortOnError(false).
		HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
		HasTwoDigitCenturyStart(1970).
		HasUnsupportedDdlAction(sdk.UnsupportedDDLAction(strings.ToLower(string(sdk.UnsupportedDDLActionIgnore)))).
		HasUseCachedResult(true).
		HasWeekOfYearPolicy(0).
		HasWeekStart(0)
}
