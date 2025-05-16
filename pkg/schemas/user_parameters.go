package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowUserParametersSchema = make(map[string]*schema.Schema)
	userParameters           = []sdk.UserParameter{
		sdk.UserParameterEnableUnredactedQuerySyntaxError,
		sdk.UserParameterNetworkPolicy,
		sdk.UserParameterPreventUnloadToInternalStages,
		sdk.UserParameterAbortDetachedQuery,
		sdk.UserParameterAutocommit,
		sdk.UserParameterBinaryInputFormat,
		sdk.UserParameterBinaryOutputFormat,
		sdk.UserParameterClientMemoryLimit,
		sdk.UserParameterClientMetadataRequestUseConnectionCtx,
		sdk.UserParameterClientPrefetchThreads,
		sdk.UserParameterClientResultChunkSize,
		sdk.UserParameterClientResultColumnCaseInsensitive,
		sdk.UserParameterClientSessionKeepAlive,
		sdk.UserParameterClientSessionKeepAliveHeartbeatFrequency,
		sdk.UserParameterClientTimestampTypeMapping,
		sdk.UserParameterDateInputFormat,
		sdk.UserParameterDateOutputFormat,
		sdk.UserParameterEnableUnloadPhysicalTypeOptimization,
		sdk.UserParameterErrorOnNondeterministicMerge,
		sdk.UserParameterErrorOnNondeterministicUpdate,
		sdk.UserParameterGeographyOutputFormat,
		sdk.UserParameterGeometryOutputFormat,
		sdk.UserParameterJdbcTreatDecimalAsInt,
		sdk.UserParameterJdbcTreatTimestampNtzAsUtc,
		sdk.UserParameterJdbcUseSessionTimezone,
		sdk.UserParameterJsonIndent,
		sdk.UserParameterLockTimeout,
		sdk.UserParameterLogLevel,
		sdk.UserParameterMultiStatementCount,
		sdk.UserParameterNoorderSequenceAsDefault,
		sdk.UserParameterOdbcTreatDecimalAsInt,
		sdk.UserParameterQueryTag,
		sdk.UserParameterQuotedIdentifiersIgnoreCase,
		sdk.UserParameterRowsPerResultset,
		sdk.UserParameterS3StageVpceDnsName,
		sdk.UserParameterSearchPath,
		sdk.UserParameterSimulatedDataSharingConsumer,
		sdk.UserParameterStatementQueuedTimeoutInSeconds,
		sdk.UserParameterStatementTimeoutInSeconds,
		sdk.UserParameterStrictJsonOutput,
		sdk.UserParameterTimestampDayIsAlways24h,
		sdk.UserParameterTimestampInputFormat,
		sdk.UserParameterTimestampLtzOutputFormat,
		sdk.UserParameterTimestampNtzOutputFormat,
		sdk.UserParameterTimestampOutputFormat,
		sdk.UserParameterTimestampTypeMapping,
		sdk.UserParameterTimestampTzOutputFormat,
		sdk.UserParameterTimezone,
		sdk.UserParameterTimeInputFormat,
		sdk.UserParameterTimeOutputFormat,
		sdk.UserParameterTraceLevel,
		sdk.UserParameterTransactionAbortOnError,
		sdk.UserParameterTransactionDefaultIsolationLevel,
		sdk.UserParameterTwoDigitCenturyStart,
		sdk.UserParameterUnsupportedDdlAction,
		sdk.UserParameterUseCachedResult,
		sdk.UserParameterWeekOfYearPolicy,
		sdk.UserParameterWeekStart,
	}
)

func init() {
	for _, param := range userParameters {
		ShowUserParametersSchema[strings.ToLower(string(param))] = ParameterListSchema
	}
}

func UserParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	userParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(userParameters, sdk.UserParameter(param.Key)) {
			userParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return userParametersValue
}
