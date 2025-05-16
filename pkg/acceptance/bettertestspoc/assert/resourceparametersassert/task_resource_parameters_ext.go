package resourceparametersassert

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

// TaskDatasourceParameters is a temporary workaround to have better parameter assertions in data source acceptance tests.
func TaskDatasourceParameters(t *testing.T, name string) *TaskResourceParametersAssert {
	t.Helper()

	taskAssert := TaskResourceParametersAssert{
		ResourceAssert: assert.NewDatasourceAssert("data."+name, "parameters", "tasks.0."),
	}
	taskAssert.AddAssertion(assert.ValueSet("parameters.#", "1"))
	return &taskAssert
}

func (u *TaskResourceParametersAssert) HasAllDefaults() *TaskResourceParametersAssert {
	return u.
		HasSuspendTaskAfterNumFailures(10).
		HasTaskAutoRetryAttempts(0).
		HasUserTaskManagedInitialWarehouseSize("Medium").
		HasUserTaskMinimumTriggerIntervalInSeconds(30).
		HasUserTaskTimeoutMs(3600000).
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
