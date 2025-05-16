package schemas

import (
	"slices"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ShowTaskParametersSchema = make(map[string]*schema.Schema)
	taskParameters           = []sdk.TaskParameter{
		// task parameters
		sdk.TaskParameterSuspendTaskAfterNumFailures,
		sdk.TaskParameterTaskAutoRetryAttempts,
		sdk.TaskParameterUserTaskManagedInitialWarehouseSize,
		sdk.TaskParameterUserTaskMinimumTriggerIntervalInSeconds,
		sdk.TaskParameterUserTaskTimeoutMs,
		// session parameters
		sdk.TaskParameterAbortDetachedQuery,
		sdk.TaskParameterAutocommit,
		sdk.TaskParameterBinaryInputFormat,
		sdk.TaskParameterBinaryOutputFormat,
		sdk.TaskParameterClientMemoryLimit,
		sdk.TaskParameterClientMetadataRequestUseConnectionCtx,
		sdk.TaskParameterClientPrefetchThreads,
		sdk.TaskParameterClientResultChunkSize,
		sdk.TaskParameterClientResultColumnCaseInsensitive,
		sdk.TaskParameterClientSessionKeepAlive,
		sdk.TaskParameterClientSessionKeepAliveHeartbeatFrequency,
		sdk.TaskParameterClientTimestampTypeMapping,
		sdk.TaskParameterDateInputFormat,
		sdk.TaskParameterDateOutputFormat,
		sdk.TaskParameterEnableUnloadPhysicalTypeOptimization,
		sdk.TaskParameterErrorOnNondeterministicMerge,
		sdk.TaskParameterErrorOnNondeterministicUpdate,
		sdk.TaskParameterGeographyOutputFormat,
		sdk.TaskParameterGeometryOutputFormat,
		sdk.TaskParameterJdbcTreatTimestampNtzAsUtc,
		sdk.TaskParameterJdbcUseSessionTimezone,
		sdk.TaskParameterJsonIndent,
		sdk.TaskParameterLockTimeout,
		sdk.TaskParameterLogLevel,
		sdk.TaskParameterMultiStatementCount,
		sdk.TaskParameterNoorderSequenceAsDefault,
		sdk.TaskParameterOdbcTreatDecimalAsInt,
		sdk.TaskParameterQueryTag,
		sdk.TaskParameterQuotedIdentifiersIgnoreCase,
		sdk.TaskParameterRowsPerResultset,
		sdk.TaskParameterS3StageVpceDnsName,
		sdk.TaskParameterSearchPath,
		sdk.TaskParameterStatementQueuedTimeoutInSeconds,
		sdk.TaskParameterStatementTimeoutInSeconds,
		sdk.TaskParameterStrictJsonOutput,
		sdk.TaskParameterTimestampDayIsAlways24h,
		sdk.TaskParameterTimestampInputFormat,
		sdk.TaskParameterTimestampLtzOutputFormat,
		sdk.TaskParameterTimestampNtzOutputFormat,
		sdk.TaskParameterTimestampOutputFormat,
		sdk.TaskParameterTimestampTypeMapping,
		sdk.TaskParameterTimestampTzOutputFormat,
		sdk.TaskParameterTimezone,
		sdk.TaskParameterTimeInputFormat,
		sdk.TaskParameterTimeOutputFormat,
		sdk.TaskParameterTraceLevel,
		sdk.TaskParameterTransactionAbortOnError,
		sdk.TaskParameterTransactionDefaultIsolationLevel,
		sdk.TaskParameterTwoDigitCenturyStart,
		sdk.TaskParameterUnsupportedDdlAction,
		sdk.TaskParameterUseCachedResult,
		sdk.TaskParameterWeekOfYearPolicy,
		sdk.TaskParameterWeekStart,
	}
)

func init() {
	for _, param := range taskParameters {
		ShowTaskParametersSchema[strings.ToLower(string(param))] = ParameterListSchema
	}
}

func TaskParametersToSchema(parameters []*sdk.Parameter) map[string]any {
	taskParametersValue := make(map[string]any)
	for _, param := range parameters {
		if slices.Contains(taskParameters, sdk.TaskParameter(param.Key)) {
			taskParametersValue[strings.ToLower(param.Key)] = []map[string]any{ParameterToSchema(param)}
		}
	}
	return taskParametersValue
}
