package resourceparametersassert

import (
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvidentifiers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

// TODO [SNOW-1501905]: generate defaults for each parameter
func (u *UserResourceParametersAssert) HasEnableUnredactedQuerySyntaxErrorValueDefault() *UserResourceParametersAssert {
	return u.HasEnableUnredactedQuerySyntaxError(false)
}

func (u *UserResourceParametersAssert) HasEnableUnredactedQuerySyntaxErrorKey() *UserResourceParametersAssert {
	u.AddAssertion(assert.ResourceParameterKeySet(sdk.UserParameterEnableUnredactedQuerySyntaxError, string(sdk.UserParameterEnableUnredactedQuerySyntaxError)))
	return u
}

func (u *UserResourceParametersAssert) HasEnableUnredactedQuerySyntaxErrorDefault() *UserResourceParametersAssert {
	u.AddAssertion(assert.ResourceParameterDefaultSet(sdk.UserParameterEnableUnredactedQuerySyntaxError, "false"))
	return u
}

func (u *UserResourceParametersAssert) HasEnableUnredactedQuerySyntaxErrorDescriptionNotEmpty() *UserResourceParametersAssert {
	u.AddAssertion(assert.ResourceParameterDescriptionPresent(sdk.UserParameterEnableUnredactedQuerySyntaxError))
	return u
}

func (u *UserResourceParametersAssert) HasEnableUnredactedQuerySyntaxErrorKeyEmpty() *UserResourceParametersAssert {
	u.AddAssertion(assert.ResourceParameterKeySet(sdk.UserParameterEnableUnredactedQuerySyntaxError, ""))
	return u
}

func (u *UserResourceParametersAssert) HasEnableUnredactedQuerySyntaxErrorDefaultEmpty() *UserResourceParametersAssert {
	u.AddAssertion(assert.ResourceParameterDefaultSet(sdk.UserParameterEnableUnredactedQuerySyntaxError, ""))
	return u
}

func (u *UserResourceParametersAssert) HasEnableUnredactedQuerySyntaxErrorDescriptionEmpty() *UserResourceParametersAssert {
	u.AddAssertion(assert.ResourceParameterDescriptionSet(sdk.UserParameterEnableUnredactedQuerySyntaxError, ""))
	return u
}
