package resources

import (
	"context"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var taskParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelTask) {
		s := parameterSchema(p)
		if p.SqlName == defs.UserTaskManagedInitialWarehouseSize.SqlName {
			s.ConflictsWith = []string{"warehouse"}
		}
		taskParametersSchema[p.FieldName()] = s
	}
	// SEARCH_PATH has no ParameterLevelTask membership in the catalog - it is not settable on a task at
	// all, so it is kept here only as a legacy field whose only behavior is rejecting any configured
	// value (see handleTaskParametersCreate/Changes below).
	taskParametersSchema[defs.SearchPath.FieldName()] = parameterSchema(defs.SearchPath)
}

func handleTaskParametersCreate(d *schema.ResourceData, createOpts *sdk.CreateTaskRequest) diag.Diagnostics {
	if v, ok := d.GetOk(defs.UserTaskManagedInitialWarehouseSize.FieldName()); ok {
		size, err := sdk.ToWarehouseSize(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		createOpts.WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(size))
	}
	return JoinDiags(
		// task parameters
		handleParameterCreate(d, defs.UserTaskTimeoutMs.FieldName(), &createOpts.UserTaskTimeoutMs),
		handleParameterCreate(d, defs.SuspendTaskAfterNumFailures.FieldName(), &createOpts.SuspendTaskAfterNumFailures),
		handleParameterCreate(d, defs.TaskAutoRetryAttempts.FieldName(), &createOpts.TaskAutoRetryAttempts),
		handleParameterCreate(d, defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), &createOpts.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterCreateWithMapping(d, defs.ServerlessTaskMinStatementSize.FieldName(), &createOpts.ServerlessTaskMinStatementSize, stringToStringEnumProvider(sdk.ToWarehouseSize)),
		handleParameterCreateWithMapping(d, defs.ServerlessTaskMaxStatementSize.FieldName(), &createOpts.ServerlessTaskMaxStatementSize, stringToStringEnumProvider(sdk.ToWarehouseSize)),
		// session parameters
		handleParameterCreate(d, defs.AbortDetachedQuery.FieldName(), &createOpts.AbortDetachedQuery),
		handleParameterCreateWithMapping(d, defs.BinaryInputFormat.FieldName(), &createOpts.BinaryInputFormat, stringToStringEnumProvider(sdk.ToBinaryInputFormat)),
		handleParameterCreateWithMapping(d, defs.BinaryOutputFormat.FieldName(), &createOpts.BinaryOutputFormat, stringToStringEnumProvider(sdk.ToBinaryOutputFormat)),
		handleParameterCreate(d, defs.ClientMemoryLimit.FieldName(), &createOpts.ClientMemoryLimit),
		handleParameterCreate(d, defs.ClientMetadataRequestUseConnectionCtx.FieldName(), &createOpts.ClientMetadataRequestUseConnectionCtx),
		handleParameterCreate(d, defs.ClientPrefetchThreads.FieldName(), &createOpts.ClientPrefetchThreads),
		handleParameterCreate(d, defs.ClientResultChunkSize.FieldName(), &createOpts.ClientResultChunkSize),
		handleParameterCreate(d, defs.ClientResultColumnCaseInsensitive.FieldName(), &createOpts.ClientResultColumnCaseInsensitive),
		handleParameterCreate(d, defs.ClientSessionKeepAlive.FieldName(), &createOpts.ClientSessionKeepAlive),
		handleParameterCreate(d, defs.ClientSessionKeepAliveHeartbeatFrequency.FieldName(), &createOpts.ClientSessionKeepAliveHeartbeatFrequency),
		handleParameterCreateWithMapping(d, defs.ClientTimestampTypeMapping.FieldName(), &createOpts.ClientTimestampTypeMapping, stringToStringEnumProvider(sdk.ToClientTimestampTypeMapping)),
		handleParameterCreate(d, defs.DateInputFormat.FieldName(), &createOpts.DateInputFormat),
		handleParameterCreate(d, defs.DateOutputFormat.FieldName(), &createOpts.DateOutputFormat),
		handleParameterCreate(d, defs.EnableUnloadPhysicalTypeOptimization.FieldName(), &createOpts.EnableUnloadPhysicalTypeOptimization),
		handleParameterCreate(d, defs.ErrorOnNondeterministicMerge.FieldName(), &createOpts.ErrorOnNondeterministicMerge),
		handleParameterCreate(d, defs.ErrorOnNondeterministicUpdate.FieldName(), &createOpts.ErrorOnNondeterministicUpdate),
		handleParameterCreateWithMapping(d, defs.GeographyOutputFormat.FieldName(), &createOpts.GeographyOutputFormat, stringToStringEnumProvider(sdk.ToGeographyOutputFormat)),
		handleParameterCreateWithMapping(d, defs.GeometryOutputFormat.FieldName(), &createOpts.GeometryOutputFormat, stringToStringEnumProvider(sdk.ToGeometryOutputFormat)),
		handleParameterCreate(d, defs.JdbcTreatTimestampNtzAsUtc.FieldName(), &createOpts.JdbcTreatTimestampNtzAsUtc),
		handleParameterCreate(d, defs.JdbcUseSessionTimezone.FieldName(), &createOpts.JdbcUseSessionTimezone),
		handleParameterCreate(d, defs.JsonIndent.FieldName(), &createOpts.JsonIndent),
		handleParameterCreate(d, defs.LockTimeout.FieldName(), &createOpts.LockTimeout),
		handleParameterCreateWithMapping(d, defs.LogLevel.FieldName(), &createOpts.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterCreateWithMapping(d, defs.LogEventLevel.FieldName(), &createOpts.LogEventLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterCreate(d, defs.MultiStatementCount.FieldName(), &createOpts.MultiStatementCount),
		handleParameterCreate(d, defs.NoorderSequenceAsDefault.FieldName(), &createOpts.NoorderSequenceAsDefault),
		handleParameterCreate(d, defs.OdbcTreatDecimalAsInt.FieldName(), &createOpts.OdbcTreatDecimalAsInt),
		handleParameterCreate(d, defs.QueryTag.FieldName(), &createOpts.QueryTag),
		handleParameterCreate(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), &createOpts.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, defs.RowsPerResultset.FieldName(), &createOpts.RowsPerResultset),
		handleParameterCreate(d, defs.S3StageVpceDnsName.FieldName(), &createOpts.S3StageVpceDnsName),
		handleParameterCreate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &createOpts.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, defs.StatementTimeoutInSeconds.FieldName(), &createOpts.StatementTimeoutInSeconds),
		handleParameterCreate(d, defs.StrictJsonOutput.FieldName(), &createOpts.StrictJsonOutput),
		handleParameterCreate(d, defs.TimestampDayIsAlways24h.FieldName(), &createOpts.TimestampDayIsAlways24H),
		handleParameterCreate(d, defs.TimestampInputFormat.FieldName(), &createOpts.TimestampInputFormat),
		handleParameterCreate(d, defs.TimestampLtzOutputFormat.FieldName(), &createOpts.TimestampLtzOutputFormat),
		handleParameterCreate(d, defs.TimestampNtzOutputFormat.FieldName(), &createOpts.TimestampNtzOutputFormat),
		handleParameterCreate(d, defs.TimestampOutputFormat.FieldName(), &createOpts.TimestampOutputFormat),
		handleParameterCreateWithMapping(d, defs.TimestampTypeMapping.FieldName(), &createOpts.TimestampTypeMapping, stringToStringEnumProvider(sdk.ToTimestampTypeMapping)),
		handleParameterCreate(d, defs.TimestampTzOutputFormat.FieldName(), &createOpts.TimestampTzOutputFormat),
		handleParameterCreate(d, defs.Timezone.FieldName(), &createOpts.Timezone),
		handleParameterCreate(d, defs.TimeInputFormat.FieldName(), &createOpts.TimeInputFormat),
		handleParameterCreate(d, defs.TimeOutputFormat.FieldName(), &createOpts.TimeOutputFormat),
		handleParameterCreateWithMapping(d, defs.TraceLevel.FieldName(), &createOpts.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
		handleParameterCreate(d, defs.TransactionAbortOnError.FieldName(), &createOpts.TransactionAbortOnError),
		handleParameterCreateWithMapping(d, defs.TransactionDefaultIsolationLevel.FieldName(), &createOpts.TransactionDefaultIsolationLevel, stringToStringEnumProvider(sdk.ToTransactionDefaultIsolationLevel)),
		handleParameterCreate(d, defs.TwoDigitCenturyStart.FieldName(), &createOpts.TwoDigitCenturyStart),
		handleParameterCreateWithMapping(d, defs.UnsupportedDdlAction.FieldName(), &createOpts.UnsupportedDdlAction, stringToStringEnumProvider(sdk.ToUnsupportedDDLAction)),
		handleParameterCreate(d, defs.UseCachedResult.FieldName(), &createOpts.UseCachedResult),
		handleParameterCreate(d, defs.WeekOfYearPolicy.FieldName(), &createOpts.WeekOfYearPolicy),
		handleParameterCreate(d, defs.WeekStart.FieldName(), &createOpts.WeekStart),
		rejectAutocommitFalseOnCreate(d, &createOpts.Autocommit),
		rejectSearchPathIfSet(d),
	)
}

func handleTaskParametersChanges(d *schema.ResourceData, set *sdk.TaskSetRequest, unset *sdk.TaskUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		// task parameters
		handleParameterUpdateWithMapping(d, defs.UserTaskManagedInitialWarehouseSize.FieldName(), &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, stringToStringEnumProvider(sdk.ToWarehouseSize)),
		handleParameterUpdate(d, defs.UserTaskTimeoutMs.FieldName(), &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
		handleParameterUpdate(d, defs.SuspendTaskAfterNumFailures.FieldName(), &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleParameterUpdate(d, defs.TaskAutoRetryAttempts.FieldName(), &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleParameterUpdate(d, defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterUpdateWithMapping(d, defs.ServerlessTaskMinStatementSize.FieldName(), &set.ServerlessTaskMinStatementSize, &unset.ServerlessTaskMinStatementSize, stringToStringEnumProvider(sdk.ToWarehouseSize)),
		handleParameterUpdateWithMapping(d, defs.ServerlessTaskMaxStatementSize.FieldName(), &set.ServerlessTaskMaxStatementSize, &unset.ServerlessTaskMaxStatementSize, stringToStringEnumProvider(sdk.ToWarehouseSize)),
		// session parameters
		handleParameterUpdate(d, defs.AbortDetachedQuery.FieldName(), &set.AbortDetachedQuery, &unset.AbortDetachedQuery),
		handleParameterUpdateWithMapping(d, defs.BinaryInputFormat.FieldName(), &set.BinaryInputFormat, &unset.BinaryInputFormat, stringToStringEnumProvider(sdk.ToBinaryInputFormat)),
		handleParameterUpdateWithMapping(d, defs.BinaryOutputFormat.FieldName(), &set.BinaryOutputFormat, &unset.BinaryOutputFormat, stringToStringEnumProvider(sdk.ToBinaryOutputFormat)),
		handleParameterUpdate(d, defs.ClientMemoryLimit.FieldName(), &set.ClientMemoryLimit, &unset.ClientMemoryLimit),
		handleParameterUpdate(d, defs.ClientMetadataRequestUseConnectionCtx.FieldName(), &set.ClientMetadataRequestUseConnectionCtx, &unset.ClientMetadataRequestUseConnectionCtx),
		handleParameterUpdate(d, defs.ClientPrefetchThreads.FieldName(), &set.ClientPrefetchThreads, &unset.ClientPrefetchThreads),
		handleParameterUpdate(d, defs.ClientResultChunkSize.FieldName(), &set.ClientResultChunkSize, &unset.ClientResultChunkSize),
		handleParameterUpdate(d, defs.ClientResultColumnCaseInsensitive.FieldName(), &set.ClientResultColumnCaseInsensitive, &unset.ClientResultColumnCaseInsensitive),
		handleParameterUpdate(d, defs.ClientSessionKeepAlive.FieldName(), &set.ClientSessionKeepAlive, &unset.ClientSessionKeepAlive),
		handleParameterUpdate(d, defs.ClientSessionKeepAliveHeartbeatFrequency.FieldName(), &set.ClientSessionKeepAliveHeartbeatFrequency, &unset.ClientSessionKeepAliveHeartbeatFrequency),
		handleParameterUpdateWithMapping(d, defs.ClientTimestampTypeMapping.FieldName(), &set.ClientTimestampTypeMapping, &unset.ClientTimestampTypeMapping, stringToStringEnumProvider(sdk.ToClientTimestampTypeMapping)),
		handleParameterUpdate(d, defs.DateInputFormat.FieldName(), &set.DateInputFormat, &unset.DateInputFormat),
		handleParameterUpdate(d, defs.DateOutputFormat.FieldName(), &set.DateOutputFormat, &unset.DateOutputFormat),
		handleParameterUpdate(d, defs.EnableUnloadPhysicalTypeOptimization.FieldName(), &set.EnableUnloadPhysicalTypeOptimization, &unset.EnableUnloadPhysicalTypeOptimization),
		handleParameterUpdate(d, defs.ErrorOnNondeterministicMerge.FieldName(), &set.ErrorOnNondeterministicMerge, &unset.ErrorOnNondeterministicMerge),
		handleParameterUpdate(d, defs.ErrorOnNondeterministicUpdate.FieldName(), &set.ErrorOnNondeterministicUpdate, &unset.ErrorOnNondeterministicUpdate),
		handleParameterUpdateWithMapping(d, defs.GeographyOutputFormat.FieldName(), &set.GeographyOutputFormat, &unset.GeographyOutputFormat, stringToStringEnumProvider(sdk.ToGeographyOutputFormat)),
		handleParameterUpdateWithMapping(d, defs.GeometryOutputFormat.FieldName(), &set.GeometryOutputFormat, &unset.GeometryOutputFormat, stringToStringEnumProvider(sdk.ToGeometryOutputFormat)),
		handleParameterUpdate(d, defs.JdbcTreatTimestampNtzAsUtc.FieldName(), &set.JdbcTreatTimestampNtzAsUtc, &unset.JdbcTreatTimestampNtzAsUtc),
		handleParameterUpdate(d, defs.JdbcUseSessionTimezone.FieldName(), &set.JdbcUseSessionTimezone, &unset.JdbcUseSessionTimezone),
		handleParameterUpdate(d, defs.JsonIndent.FieldName(), &set.JsonIndent, &unset.JsonIndent),
		handleParameterUpdate(d, defs.LockTimeout.FieldName(), &set.LockTimeout, &unset.LockTimeout),
		handleParameterUpdateWithMapping(d, defs.LogLevel.FieldName(), &set.LogLevel, &unset.LogLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterUpdateWithMapping(d, defs.LogEventLevel.FieldName(), &set.LogEventLevel, &unset.LogEventLevel, stringToStringEnumProvider(sdk.ToLogLevel)),
		handleParameterUpdate(d, defs.MultiStatementCount.FieldName(), &set.MultiStatementCount, &unset.MultiStatementCount),
		handleParameterUpdate(d, defs.NoorderSequenceAsDefault.FieldName(), &set.NoorderSequenceAsDefault, &unset.NoorderSequenceAsDefault),
		handleParameterUpdate(d, defs.OdbcTreatDecimalAsInt.FieldName(), &set.OdbcTreatDecimalAsInt, &unset.OdbcTreatDecimalAsInt),
		handleParameterUpdate(d, defs.QueryTag.FieldName(), &set.QueryTag, &unset.QueryTag),
		handleParameterUpdate(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, defs.RowsPerResultset.FieldName(), &set.RowsPerResultset, &unset.RowsPerResultset),
		handleParameterUpdate(d, defs.S3StageVpceDnsName.FieldName(), &set.S3StageVpceDnsName, &unset.S3StageVpceDnsName),
		handleParameterUpdate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, defs.StatementTimeoutInSeconds.FieldName(), &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
		handleParameterUpdate(d, defs.StrictJsonOutput.FieldName(), &set.StrictJsonOutput, &unset.StrictJsonOutput),
		handleParameterUpdate(d, defs.TimestampDayIsAlways24h.FieldName(), &set.TimestampDayIsAlways24H, &unset.TimestampDayIsAlways24H),
		handleParameterUpdate(d, defs.TimestampInputFormat.FieldName(), &set.TimestampInputFormat, &unset.TimestampInputFormat),
		handleParameterUpdate(d, defs.TimestampLtzOutputFormat.FieldName(), &set.TimestampLtzOutputFormat, &unset.TimestampLtzOutputFormat),
		handleParameterUpdate(d, defs.TimestampNtzOutputFormat.FieldName(), &set.TimestampNtzOutputFormat, &unset.TimestampNtzOutputFormat),
		handleParameterUpdate(d, defs.TimestampOutputFormat.FieldName(), &set.TimestampOutputFormat, &unset.TimestampOutputFormat),
		handleParameterUpdateWithMapping(d, defs.TimestampTypeMapping.FieldName(), &set.TimestampTypeMapping, &unset.TimestampTypeMapping, stringToStringEnumProvider(sdk.ToTimestampTypeMapping)),
		handleParameterUpdate(d, defs.TimestampTzOutputFormat.FieldName(), &set.TimestampTzOutputFormat, &unset.TimestampTzOutputFormat),
		handleParameterUpdate(d, defs.Timezone.FieldName(), &set.Timezone, &unset.Timezone),
		handleParameterUpdate(d, defs.TimeInputFormat.FieldName(), &set.TimeInputFormat, &unset.TimeInputFormat),
		handleParameterUpdate(d, defs.TimeOutputFormat.FieldName(), &set.TimeOutputFormat, &unset.TimeOutputFormat),
		handleParameterUpdateWithMapping(d, defs.TraceLevel.FieldName(), &set.TraceLevel, &unset.TraceLevel, stringToStringEnumProvider(sdk.ToTraceLevel)),
		handleParameterUpdate(d, defs.TransactionAbortOnError.FieldName(), &set.TransactionAbortOnError, &unset.TransactionAbortOnError),
		handleParameterUpdateWithMapping(d, defs.TransactionDefaultIsolationLevel.FieldName(), &set.TransactionDefaultIsolationLevel, &unset.TransactionDefaultIsolationLevel, stringToStringEnumProvider(sdk.ToTransactionDefaultIsolationLevel)),
		handleParameterUpdate(d, defs.TwoDigitCenturyStart.FieldName(), &set.TwoDigitCenturyStart, &unset.TwoDigitCenturyStart),
		handleParameterUpdateWithMapping(d, defs.UnsupportedDdlAction.FieldName(), &set.UnsupportedDdlAction, &unset.UnsupportedDdlAction, stringToStringEnumProvider(sdk.ToUnsupportedDDLAction)),
		handleParameterUpdate(d, defs.UseCachedResult.FieldName(), &set.UseCachedResult, &unset.UseCachedResult),
		handleParameterUpdate(d, defs.WeekOfYearPolicy.FieldName(), &set.WeekOfYearPolicy, &unset.WeekOfYearPolicy),
		handleParameterUpdate(d, defs.WeekStart.FieldName(), &set.WeekStart, &unset.WeekStart),
		rejectAutocommitFalseOnUpdate(d, &set.Autocommit, &unset.Autocommit),
		rejectSearchPathIfSet(d),
	)
}

// rejectAutocommitFalseOnCreate mirrors Snowflake's own validation: AUTOCOMMIT may only be TRUE on a
// task, never FALSE. Returning a diagnostic (rather than failing the SDK call) lets the config surface
// the error at plan time instead of relying on a mid-apply SQL failure.
func rejectAutocommitFalseOnCreate(d *schema.ResourceData, autocommit **bool) diag.Diagnostics {
	key := defs.Autocommit.FieldName()
	if v := GetConfigPropertyAsPointerAllowingZeroValue[bool](d, key); v != nil {
		if !*v {
			return diag.Diagnostics{{Severity: diag.Warning, Summary: "Invalid value for AUTOCOMMIT parameter: cannot be set to FALSE on a task"}}
		}
		*autocommit = v
	}
	return nil
}

func rejectAutocommitFalseOnUpdate(d *schema.ResourceData, set, unset **bool) diag.Diagnostics {
	key := defs.Autocommit.FieldName()
	if d.HasChange(key) || !d.GetRawPlan().AsValueMap()[key].IsKnown() {
		if !d.GetRawConfig().AsValueMap()[key].IsNull() {
			if !d.Get(key).(bool) {
				return diag.Diagnostics{{Severity: diag.Warning, Summary: "Invalid value for AUTOCOMMIT parameter: cannot be set to FALSE on a task"}}
			}
			*set = sdk.Bool(true)
		} else {
			*unset = sdk.Bool(true)
		}
	}
	return nil
}

// rejectSearchPathIfSet preserves the pre-catalog behavior: SEARCH_PATH has no ParameterLevelTask
// membership (it is not settable on a task at all), so the field stays in the schema only to surface
// this rejection instead of a less helpful "unknown attribute" error.
func rejectSearchPathIfSet(d *schema.ResourceData) diag.Diagnostics {
	key := defs.SearchPath.FieldName()
	if v := GetConfigPropertyAsPointerAllowingZeroValue[string](d, key); v != nil {
		return diag.Diagnostics{{Severity: diag.Warning, Summary: "Invalid value for SEARCH_PATH parameter: cannot be set on a task"}}
	}
	return nil
}

// handleTaskParameterRead sets every catalog-driven field from the typed parameters, plus SEARCH_PATH
// (not part of the ParameterLevelTask catalog set) read from the raw SHOW PARAMETERS result the caller
// already fetches for the parameters output block.
func handleTaskParameterRead(d *schema.ResourceData, parameters *sdk.TaskParametersDetails, rawParameters []*sdk.Parameter) diag.Diagnostics {
	return JoinDiags(
		// task parameters
		setResourceData(d, defs.UserTaskManagedInitialWarehouseSize.FieldName(), parameters.UserTaskManagedInitialWarehouseSize.Value),
		setResourceData(d, defs.UserTaskTimeoutMs.FieldName(), parameters.UserTaskTimeoutMs.Value),
		setResourceData(d, defs.SuspendTaskAfterNumFailures.FieldName(), parameters.SuspendTaskAfterNumFailures.Value),
		setResourceData(d, defs.TaskAutoRetryAttempts.FieldName(), parameters.TaskAutoRetryAttempts.Value),
		setResourceData(d, defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), parameters.UserTaskMinimumTriggerIntervalInSeconds.Value),
		setResourceData(d, defs.ServerlessTaskMinStatementSize.FieldName(), parameters.ServerlessTaskMinStatementSize.Value),
		setResourceData(d, defs.ServerlessTaskMaxStatementSize.FieldName(), parameters.ServerlessTaskMaxStatementSize.Value),
		// session parameters
		setResourceData(d, defs.AbortDetachedQuery.FieldName(), parameters.AbortDetachedQuery.Value),
		setResourceData(d, defs.Autocommit.FieldName(), parameters.Autocommit.Value),
		setResourceData(d, defs.BinaryInputFormat.FieldName(), parameters.BinaryInputFormat.Value),
		setResourceData(d, defs.BinaryOutputFormat.FieldName(), parameters.BinaryOutputFormat.Value),
		setResourceData(d, defs.ClientMemoryLimit.FieldName(), parameters.ClientMemoryLimit.Value),
		setResourceData(d, defs.ClientMetadataRequestUseConnectionCtx.FieldName(), parameters.ClientMetadataRequestUseConnectionCtx.Value),
		setResourceData(d, defs.ClientPrefetchThreads.FieldName(), parameters.ClientPrefetchThreads.Value),
		setResourceData(d, defs.ClientResultChunkSize.FieldName(), parameters.ClientResultChunkSize.Value),
		setResourceData(d, defs.ClientResultColumnCaseInsensitive.FieldName(), parameters.ClientResultColumnCaseInsensitive.Value),
		setResourceData(d, defs.ClientSessionKeepAlive.FieldName(), parameters.ClientSessionKeepAlive.Value),
		setResourceData(d, defs.ClientSessionKeepAliveHeartbeatFrequency.FieldName(), parameters.ClientSessionKeepAliveHeartbeatFrequency.Value),
		setResourceData(d, defs.ClientTimestampTypeMapping.FieldName(), parameters.ClientTimestampTypeMapping.Value),
		setResourceData(d, defs.DateInputFormat.FieldName(), parameters.DateInputFormat.Value),
		setResourceData(d, defs.DateOutputFormat.FieldName(), parameters.DateOutputFormat.Value),
		setResourceData(d, defs.EnableUnloadPhysicalTypeOptimization.FieldName(), parameters.EnableUnloadPhysicalTypeOptimization.Value),
		setResourceData(d, defs.ErrorOnNondeterministicMerge.FieldName(), parameters.ErrorOnNondeterministicMerge.Value),
		setResourceData(d, defs.ErrorOnNondeterministicUpdate.FieldName(), parameters.ErrorOnNondeterministicUpdate.Value),
		setResourceData(d, defs.GeographyOutputFormat.FieldName(), parameters.GeographyOutputFormat.Value),
		setResourceData(d, defs.GeometryOutputFormat.FieldName(), parameters.GeometryOutputFormat.Value),
		setResourceData(d, defs.JdbcTreatTimestampNtzAsUtc.FieldName(), parameters.JdbcTreatTimestampNtzAsUtc.Value),
		setResourceData(d, defs.JdbcUseSessionTimezone.FieldName(), parameters.JdbcUseSessionTimezone.Value),
		setResourceData(d, defs.JsonIndent.FieldName(), parameters.JsonIndent.Value),
		setResourceData(d, defs.LockTimeout.FieldName(), parameters.LockTimeout.Value),
		setResourceData(d, defs.LogLevel.FieldName(), parameters.LogLevel.Value),
		setResourceData(d, defs.LogEventLevel.FieldName(), parameters.LogEventLevel.Value),
		setResourceData(d, defs.MultiStatementCount.FieldName(), parameters.MultiStatementCount.Value),
		setResourceData(d, defs.NoorderSequenceAsDefault.FieldName(), parameters.NoorderSequenceAsDefault.Value),
		setResourceData(d, defs.OdbcTreatDecimalAsInt.FieldName(), parameters.OdbcTreatDecimalAsInt.Value),
		setResourceData(d, defs.QueryTag.FieldName(), parameters.QueryTag.Value),
		setResourceData(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), parameters.QuotedIdentifiersIgnoreCase.Value),
		setResourceData(d, defs.RowsPerResultset.FieldName(), parameters.RowsPerResultset.Value),
		setResourceData(d, defs.S3StageVpceDnsName.FieldName(), parameters.S3StageVpceDnsName.Value),
		setResourceData(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds.Value),
		setResourceData(d, defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds.Value),
		setResourceData(d, defs.StrictJsonOutput.FieldName(), parameters.StrictJsonOutput.Value),
		setResourceData(d, defs.TimestampDayIsAlways24h.FieldName(), parameters.TimestampDayIsAlways24H.Value),
		setResourceData(d, defs.TimestampInputFormat.FieldName(), parameters.TimestampInputFormat.Value),
		setResourceData(d, defs.TimestampLtzOutputFormat.FieldName(), parameters.TimestampLtzOutputFormat.Value),
		setResourceData(d, defs.TimestampNtzOutputFormat.FieldName(), parameters.TimestampNtzOutputFormat.Value),
		setResourceData(d, defs.TimestampOutputFormat.FieldName(), parameters.TimestampOutputFormat.Value),
		setResourceData(d, defs.TimestampTypeMapping.FieldName(), parameters.TimestampTypeMapping.Value),
		setResourceData(d, defs.TimestampTzOutputFormat.FieldName(), parameters.TimestampTzOutputFormat.Value),
		setResourceData(d, defs.Timezone.FieldName(), parameters.Timezone.Value),
		setResourceData(d, defs.TimeInputFormat.FieldName(), parameters.TimeInputFormat.Value),
		setResourceData(d, defs.TimeOutputFormat.FieldName(), parameters.TimeOutputFormat.Value),
		setResourceData(d, defs.TraceLevel.FieldName(), parameters.TraceLevel.Value),
		setResourceData(d, defs.TransactionAbortOnError.FieldName(), parameters.TransactionAbortOnError.Value),
		setResourceData(d, defs.TransactionDefaultIsolationLevel.FieldName(), parameters.TransactionDefaultIsolationLevel.Value),
		setResourceData(d, defs.TwoDigitCenturyStart.FieldName(), parameters.TwoDigitCenturyStart.Value),
		setResourceData(d, defs.UnsupportedDdlAction.FieldName(), parameters.UnsupportedDdlAction.Value),
		setResourceData(d, defs.UseCachedResult.FieldName(), parameters.UseCachedResult.Value),
		setResourceData(d, defs.WeekOfYearPolicy.FieldName(), parameters.WeekOfYearPolicy.Value),
		setResourceData(d, defs.WeekStart.FieldName(), parameters.WeekStart.Value),
		func() diag.Diagnostics {
			for _, p := range rawParameters {
				if strings.EqualFold(p.Key, defs.SearchPath.SqlName) {
					return setResourceData(d, defs.SearchPath.FieldName(), p.Value)
				}
			}
			return nil
		}(),
	)
}

var taskParametersCustomDiff = ParametersCustomDiffFromTypedParameters(taskParametersProvider, taskParameterDiffFunctions)

func taskParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.TaskParametersDetails, error) {
	id, err := sdk.ParseSchemaObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.Tasks.ShowParametersDetails(ctx, id)
}

func taskParameterDiffFunctions(parameters *sdk.TaskParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		StringTypedParameterValueComputedIf(defs.UserTaskManagedInitialWarehouseSize.FieldName(), parameters.UserTaskManagedInitialWarehouseSize, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.UserTaskTimeoutMs.FieldName(), parameters.UserTaskTimeoutMs, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.SuspendTaskAfterNumFailures.FieldName(), parameters.SuspendTaskAfterNumFailures, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.TaskAutoRetryAttempts.FieldName(), parameters.TaskAutoRetryAttempts, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), parameters.UserTaskMinimumTriggerIntervalInSeconds, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.ServerlessTaskMinStatementSize.FieldName(), parameters.ServerlessTaskMinStatementSize, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.ServerlessTaskMaxStatementSize.FieldName(), parameters.ServerlessTaskMaxStatementSize, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.AbortDetachedQuery.FieldName(), parameters.AbortDetachedQuery, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.Autocommit.FieldName(), parameters.Autocommit, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.BinaryInputFormat.FieldName(), parameters.BinaryInputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.BinaryOutputFormat.FieldName(), parameters.BinaryOutputFormat, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.ClientMemoryLimit.FieldName(), parameters.ClientMemoryLimit, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.ClientMetadataRequestUseConnectionCtx.FieldName(), parameters.ClientMetadataRequestUseConnectionCtx, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.ClientPrefetchThreads.FieldName(), parameters.ClientPrefetchThreads, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.ClientResultChunkSize.FieldName(), parameters.ClientResultChunkSize, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.ClientResultColumnCaseInsensitive.FieldName(), parameters.ClientResultColumnCaseInsensitive, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.ClientSessionKeepAlive.FieldName(), parameters.ClientSessionKeepAlive, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.ClientSessionKeepAliveHeartbeatFrequency.FieldName(), parameters.ClientSessionKeepAliveHeartbeatFrequency, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.ClientTimestampTypeMapping.FieldName(), parameters.ClientTimestampTypeMapping, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.DateInputFormat.FieldName(), parameters.DateInputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.DateOutputFormat.FieldName(), parameters.DateOutputFormat, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.EnableUnloadPhysicalTypeOptimization.FieldName(), parameters.EnableUnloadPhysicalTypeOptimization, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.ErrorOnNondeterministicMerge.FieldName(), parameters.ErrorOnNondeterministicMerge, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.ErrorOnNondeterministicUpdate.FieldName(), parameters.ErrorOnNondeterministicUpdate, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.GeographyOutputFormat.FieldName(), parameters.GeographyOutputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.GeometryOutputFormat.FieldName(), parameters.GeometryOutputFormat, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.JdbcTreatTimestampNtzAsUtc.FieldName(), parameters.JdbcTreatTimestampNtzAsUtc, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.JdbcUseSessionTimezone.FieldName(), parameters.JdbcUseSessionTimezone, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.JsonIndent.FieldName(), parameters.JsonIndent, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.LockTimeout.FieldName(), parameters.LockTimeout, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.LogLevel.FieldName(), parameters.LogLevel, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.LogEventLevel.FieldName(), parameters.LogEventLevel, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.MultiStatementCount.FieldName(), parameters.MultiStatementCount, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.NoorderSequenceAsDefault.FieldName(), parameters.NoorderSequenceAsDefault, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.OdbcTreatDecimalAsInt.FieldName(), parameters.OdbcTreatDecimalAsInt, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.QueryTag.FieldName(), parameters.QueryTag, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.QuotedIdentifiersIgnoreCase.FieldName(), parameters.QuotedIdentifiersIgnoreCase, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.RowsPerResultset.FieldName(), parameters.RowsPerResultset, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.S3StageVpceDnsName.FieldName(), parameters.S3StageVpceDnsName, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.StrictJsonOutput.FieldName(), parameters.StrictJsonOutput, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.TimestampDayIsAlways24h.FieldName(), parameters.TimestampDayIsAlways24H, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TimestampInputFormat.FieldName(), parameters.TimestampInputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TimestampLtzOutputFormat.FieldName(), parameters.TimestampLtzOutputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TimestampNtzOutputFormat.FieldName(), parameters.TimestampNtzOutputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TimestampOutputFormat.FieldName(), parameters.TimestampOutputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TimestampTypeMapping.FieldName(), parameters.TimestampTypeMapping, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TimestampTzOutputFormat.FieldName(), parameters.TimestampTzOutputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.Timezone.FieldName(), parameters.Timezone, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TimeInputFormat.FieldName(), parameters.TimeInputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TimeOutputFormat.FieldName(), parameters.TimeOutputFormat, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TraceLevel.FieldName(), parameters.TraceLevel, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.TransactionAbortOnError.FieldName(), parameters.TransactionAbortOnError, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.TransactionDefaultIsolationLevel.FieldName(), parameters.TransactionDefaultIsolationLevel, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.TwoDigitCenturyStart.FieldName(), parameters.TwoDigitCenturyStart, sdk.ParameterTypeTask),
		StringTypedParameterValueComputedIf(defs.UnsupportedDdlAction.FieldName(), parameters.UnsupportedDdlAction, sdk.ParameterTypeTask),
		BoolTypedParameterValueComputedIf(defs.UseCachedResult.FieldName(), parameters.UseCachedResult, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.WeekOfYearPolicy.FieldName(), parameters.WeekOfYearPolicy, sdk.ParameterTypeTask),
		IntTypedParameterValueComputedIf(defs.WeekStart.FieldName(), parameters.WeekStart, sdk.ParameterTypeTask),
	}
}
