package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var userParametersSchema = make(map[string]*schema.Schema)

var userParameterFieldNames = collections.Map(
	defs.ParameterDefsForLevel(parameterdefs.ParameterLevelUser),
	func(p parameterdefs.ParameterDef) string { return p.FieldName() },
)

type parameterDef[T ~string] struct {
	Name          T
	Type          schema.ValueType
	Description   string
	DiffSuppress  schema.SchemaDiffSuppressFunc
	ValidateDiag  schema.SchemaValidateDiagFunc
	ConflictsWith []string
}

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelUser) {
		// User parameters previously had no numeric bounds. Keep that behavior to avoid
		// introducing breaking validation as part of the catalog migration.
		userParametersSchema[p.FieldName()] = parameterSchemaWithoutModifiers(p)
	}
}

func handleUserParametersCreate(d *schema.ResourceData, create *sdk.CreateUserRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.AbortDetachedQuery.FieldName(), &create.AbortDetachedQuery),
		handleParameterCreate(d, defs.Autocommit.FieldName(), &create.Autocommit),
		handleParameterCreateWithMapping(d, defs.BinaryInputFormat.FieldName(), &create.BinaryInputFormat, sdk.ToBinaryInputFormat),
		handleParameterCreateWithMapping(d, defs.BinaryOutputFormat.FieldName(), &create.BinaryOutputFormat, sdk.ToBinaryOutputFormat),
		handleParameterCreate(d, defs.ClientMemoryLimit.FieldName(), &create.ClientMemoryLimit),
		handleParameterCreate(d, defs.ClientMetadataRequestUseConnectionCtx.FieldName(), &create.ClientMetadataRequestUseConnectionCtx),
		handleParameterCreate(d, defs.ClientPrefetchThreads.FieldName(), &create.ClientPrefetchThreads),
		handleParameterCreate(d, defs.ClientResultChunkSize.FieldName(), &create.ClientResultChunkSize),
		handleParameterCreate(d, defs.ClientResultColumnCaseInsensitive.FieldName(), &create.ClientResultColumnCaseInsensitive),
		handleParameterCreate(d, defs.ClientSessionKeepAlive.FieldName(), &create.ClientSessionKeepAlive),
		handleParameterCreate(d, defs.ClientSessionKeepAliveHeartbeatFrequency.FieldName(), &create.ClientSessionKeepAliveHeartbeatFrequency),
		handleParameterCreateWithMapping(d, defs.ClientTimestampTypeMapping.FieldName(), &create.ClientTimestampTypeMapping, sdk.ToClientTimestampTypeMapping),
		handleParameterCreate(d, defs.DateInputFormat.FieldName(), &create.DateInputFormat),
		handleParameterCreate(d, defs.DateOutputFormat.FieldName(), &create.DateOutputFormat),
		handleParameterCreate(d, defs.EnableUnloadPhysicalTypeOptimization.FieldName(), &create.EnableUnloadPhysicalTypeOptimization),
		handleParameterCreate(d, defs.ErrorOnNondeterministicMerge.FieldName(), &create.ErrorOnNondeterministicMerge),
		handleParameterCreate(d, defs.ErrorOnNondeterministicUpdate.FieldName(), &create.ErrorOnNondeterministicUpdate),
		handleParameterCreateWithMapping(d, defs.GeographyOutputFormat.FieldName(), &create.GeographyOutputFormat, sdk.ToGeographyOutputFormat),
		handleParameterCreateWithMapping(d, defs.GeometryOutputFormat.FieldName(), &create.GeometryOutputFormat, sdk.ToGeometryOutputFormat),
		handleParameterCreate(d, defs.JdbcTreatDecimalAsInt.FieldName(), &create.JdbcTreatDecimalAsInt),
		handleParameterCreate(d, defs.JdbcTreatTimestampNtzAsUtc.FieldName(), &create.JdbcTreatTimestampNtzAsUtc),
		handleParameterCreate(d, defs.JdbcUseSessionTimezone.FieldName(), &create.JdbcUseSessionTimezone),
		handleParameterCreate(d, defs.JsonIndent.FieldName(), &create.JsonIndent),
		handleParameterCreate(d, defs.LockTimeout.FieldName(), &create.LockTimeout),
		handleParameterCreateWithMapping(d, defs.LogLevel.FieldName(), &create.LogLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, defs.LogEventLevel.FieldName(), &create.LogEventLevel, sdk.ToLogLevel),
		handleParameterCreate(d, defs.MultiStatementCount.FieldName(), &create.MultiStatementCount),
		handleParameterCreate(d, defs.NoorderSequenceAsDefault.FieldName(), &create.NoorderSequenceAsDefault),
		handleParameterCreate(d, defs.OdbcTreatDecimalAsInt.FieldName(), &create.OdbcTreatDecimalAsInt),
		handleParameterCreate(d, defs.QueryTag.FieldName(), &create.QueryTag),
		handleParameterCreate(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), &create.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, defs.RowsPerResultset.FieldName(), &create.RowsPerResultset),
		handleParameterCreate(d, defs.S3StageVpceDnsName.FieldName(), &create.S3StageVpceDnsName),
		handleParameterCreate(d, defs.SearchPath.FieldName(), &create.SearchPath),
		handleParameterCreate(d, defs.SimulatedDataSharingConsumer.FieldName(), &create.SimulatedDataSharingConsumer),
		handleParameterCreate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &create.StatementQueuedTimeoutInSeconds),
		handleParameterCreate(d, defs.StatementTimeoutInSeconds.FieldName(), &create.StatementTimeoutInSeconds),
		handleParameterCreate(d, defs.StrictJsonOutput.FieldName(), &create.StrictJsonOutput),
		handleParameterCreate(d, defs.TimestampDayIsAlways24h.FieldName(), &create.TimestampDayIsAlways24H),
		handleParameterCreate(d, defs.TimestampInputFormat.FieldName(), &create.TimestampInputFormat),
		handleParameterCreate(d, defs.TimestampLtzOutputFormat.FieldName(), &create.TimestampLtzOutputFormat),
		handleParameterCreate(d, defs.TimestampNtzOutputFormat.FieldName(), &create.TimestampNtzOutputFormat),
		handleParameterCreate(d, defs.TimestampOutputFormat.FieldName(), &create.TimestampOutputFormat),
		handleParameterCreateWithMapping(d, defs.TimestampTypeMapping.FieldName(), &create.TimestampTypeMapping, sdk.ToTimestampTypeMapping),
		handleParameterCreate(d, defs.TimestampTzOutputFormat.FieldName(), &create.TimestampTzOutputFormat),
		handleParameterCreate(d, defs.Timezone.FieldName(), &create.Timezone),
		handleParameterCreate(d, defs.TimeInputFormat.FieldName(), &create.TimeInputFormat),
		handleParameterCreate(d, defs.TimeOutputFormat.FieldName(), &create.TimeOutputFormat),
		handleParameterCreateWithMapping(d, defs.TraceLevel.FieldName(), &create.TraceLevel, sdk.ToTraceLevel),
		handleParameterCreate(d, defs.TransactionAbortOnError.FieldName(), &create.TransactionAbortOnError),
		handleParameterCreateWithMapping(d, defs.TransactionDefaultIsolationLevel.FieldName(), &create.TransactionDefaultIsolationLevel, sdk.ToTransactionDefaultIsolationLevel),
		handleParameterCreate(d, defs.TwoDigitCenturyStart.FieldName(), &create.TwoDigitCenturyStart),
		handleParameterCreateWithMapping(d, defs.UnsupportedDdlAction.FieldName(), &create.UnsupportedDdlAction, sdk.ToUnsupportedDDLAction),
		handleParameterCreate(d, defs.UseCachedResult.FieldName(), &create.UseCachedResult),
		handleParameterCreate(d, defs.WeekOfYearPolicy.FieldName(), &create.WeekOfYearPolicy),
		handleParameterCreate(d, defs.WeekStart.FieldName(), &create.WeekStart),
		handleParameterCreate(d, defs.EnableUnredactedQuerySyntaxError.FieldName(), &create.EnableUnredactedQuerySyntaxError),
		handleParameterCreateWithMapping(d, defs.NetworkPolicy.FieldName(), &create.NetworkPolicy, stringToAccountObjectIdentifier),
		handleParameterCreate(d, defs.PreventUnloadToInternalStages.FieldName(), &create.PreventUnloadToInternalStages),
	)
}

func handleUserParametersUpdate(d *schema.ResourceData, set *sdk.UserSetRequest, unset *sdk.UserUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.AbortDetachedQuery.FieldName(), &set.AbortDetachedQuery, &unset.AbortDetachedQuery),
		handleParameterUpdate(d, defs.Autocommit.FieldName(), &set.Autocommit, &unset.Autocommit),
		handleParameterUpdateWithMapping(d, defs.BinaryInputFormat.FieldName(), &set.BinaryInputFormat, &unset.BinaryInputFormat, sdk.ToBinaryInputFormat),
		handleParameterUpdateWithMapping(d, defs.BinaryOutputFormat.FieldName(), &set.BinaryOutputFormat, &unset.BinaryOutputFormat, sdk.ToBinaryOutputFormat),
		handleParameterUpdate(d, defs.ClientMemoryLimit.FieldName(), &set.ClientMemoryLimit, &unset.ClientMemoryLimit),
		handleParameterUpdate(d, defs.ClientMetadataRequestUseConnectionCtx.FieldName(), &set.ClientMetadataRequestUseConnectionCtx, &unset.ClientMetadataRequestUseConnectionCtx),
		handleParameterUpdate(d, defs.ClientPrefetchThreads.FieldName(), &set.ClientPrefetchThreads, &unset.ClientPrefetchThreads),
		handleParameterUpdate(d, defs.ClientResultChunkSize.FieldName(), &set.ClientResultChunkSize, &unset.ClientResultChunkSize),
		handleParameterUpdate(d, defs.ClientResultColumnCaseInsensitive.FieldName(), &set.ClientResultColumnCaseInsensitive, &unset.ClientResultColumnCaseInsensitive),
		handleParameterUpdate(d, defs.ClientSessionKeepAlive.FieldName(), &set.ClientSessionKeepAlive, &unset.ClientSessionKeepAlive),
		handleParameterUpdate(d, defs.ClientSessionKeepAliveHeartbeatFrequency.FieldName(), &set.ClientSessionKeepAliveHeartbeatFrequency, &unset.ClientSessionKeepAliveHeartbeatFrequency),
		handleParameterUpdateWithMapping(d, defs.ClientTimestampTypeMapping.FieldName(), &set.ClientTimestampTypeMapping, &unset.ClientTimestampTypeMapping, sdk.ToClientTimestampTypeMapping),
		handleParameterUpdate(d, defs.DateInputFormat.FieldName(), &set.DateInputFormat, &unset.DateInputFormat),
		handleParameterUpdate(d, defs.DateOutputFormat.FieldName(), &set.DateOutputFormat, &unset.DateOutputFormat),
		handleParameterUpdate(d, defs.EnableUnloadPhysicalTypeOptimization.FieldName(), &set.EnableUnloadPhysicalTypeOptimization, &unset.EnableUnloadPhysicalTypeOptimization),
		handleParameterUpdate(d, defs.ErrorOnNondeterministicMerge.FieldName(), &set.ErrorOnNondeterministicMerge, &unset.ErrorOnNondeterministicMerge),
		handleParameterUpdate(d, defs.ErrorOnNondeterministicUpdate.FieldName(), &set.ErrorOnNondeterministicUpdate, &unset.ErrorOnNondeterministicUpdate),
		handleParameterUpdateWithMapping(d, defs.GeographyOutputFormat.FieldName(), &set.GeographyOutputFormat, &unset.GeographyOutputFormat, sdk.ToGeographyOutputFormat),
		handleParameterUpdateWithMapping(d, defs.GeometryOutputFormat.FieldName(), &set.GeometryOutputFormat, &unset.GeometryOutputFormat, sdk.ToGeometryOutputFormat),
		handleParameterUpdate(d, defs.JdbcTreatDecimalAsInt.FieldName(), &set.JdbcTreatDecimalAsInt, &unset.JdbcTreatDecimalAsInt),
		handleParameterUpdate(d, defs.JdbcTreatTimestampNtzAsUtc.FieldName(), &set.JdbcTreatTimestampNtzAsUtc, &unset.JdbcTreatTimestampNtzAsUtc),
		handleParameterUpdate(d, defs.JdbcUseSessionTimezone.FieldName(), &set.JdbcUseSessionTimezone, &unset.JdbcUseSessionTimezone),
		handleParameterUpdate(d, defs.JsonIndent.FieldName(), &set.JsonIndent, &unset.JsonIndent),
		handleParameterUpdate(d, defs.LockTimeout.FieldName(), &set.LockTimeout, &unset.LockTimeout),
		handleParameterUpdateWithMapping(d, defs.LogLevel.FieldName(), &set.LogLevel, &unset.LogLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, defs.LogEventLevel.FieldName(), &set.LogEventLevel, &unset.LogEventLevel, sdk.ToLogLevel),
		handleParameterUpdate(d, defs.MultiStatementCount.FieldName(), &set.MultiStatementCount, &unset.MultiStatementCount),
		handleParameterUpdate(d, defs.NoorderSequenceAsDefault.FieldName(), &set.NoorderSequenceAsDefault, &unset.NoorderSequenceAsDefault),
		handleParameterUpdate(d, defs.OdbcTreatDecimalAsInt.FieldName(), &set.OdbcTreatDecimalAsInt, &unset.OdbcTreatDecimalAsInt),
		handleParameterUpdate(d, defs.QueryTag.FieldName(), &set.QueryTag, &unset.QueryTag),
		handleParameterUpdate(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, defs.RowsPerResultset.FieldName(), &set.RowsPerResultset, &unset.RowsPerResultset),
		handleParameterUpdate(d, defs.S3StageVpceDnsName.FieldName(), &set.S3StageVpceDnsName, &unset.S3StageVpceDnsName),
		handleParameterUpdate(d, defs.SearchPath.FieldName(), &set.SearchPath, &unset.SearchPath),
		handleParameterUpdate(d, defs.SimulatedDataSharingConsumer.FieldName(), &set.SimulatedDataSharingConsumer, &unset.SimulatedDataSharingConsumer),
		handleParameterUpdate(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), &set.StatementQueuedTimeoutInSeconds, &unset.StatementQueuedTimeoutInSeconds),
		handleParameterUpdate(d, defs.StatementTimeoutInSeconds.FieldName(), &set.StatementTimeoutInSeconds, &unset.StatementTimeoutInSeconds),
		handleParameterUpdate(d, defs.StrictJsonOutput.FieldName(), &set.StrictJsonOutput, &unset.StrictJsonOutput),
		handleParameterUpdate(d, defs.TimestampDayIsAlways24h.FieldName(), &set.TimestampDayIsAlways24H, &unset.TimestampDayIsAlways24H),
		handleParameterUpdate(d, defs.TimestampInputFormat.FieldName(), &set.TimestampInputFormat, &unset.TimestampInputFormat),
		handleParameterUpdate(d, defs.TimestampLtzOutputFormat.FieldName(), &set.TimestampLtzOutputFormat, &unset.TimestampLtzOutputFormat),
		handleParameterUpdate(d, defs.TimestampNtzOutputFormat.FieldName(), &set.TimestampNtzOutputFormat, &unset.TimestampNtzOutputFormat),
		handleParameterUpdate(d, defs.TimestampOutputFormat.FieldName(), &set.TimestampOutputFormat, &unset.TimestampOutputFormat),
		handleParameterUpdateWithMapping(d, defs.TimestampTypeMapping.FieldName(), &set.TimestampTypeMapping, &unset.TimestampTypeMapping, sdk.ToTimestampTypeMapping),
		handleParameterUpdate(d, defs.TimestampTzOutputFormat.FieldName(), &set.TimestampTzOutputFormat, &unset.TimestampTzOutputFormat),
		handleParameterUpdate(d, defs.Timezone.FieldName(), &set.Timezone, &unset.Timezone),
		handleParameterUpdate(d, defs.TimeInputFormat.FieldName(), &set.TimeInputFormat, &unset.TimeInputFormat),
		handleParameterUpdate(d, defs.TimeOutputFormat.FieldName(), &set.TimeOutputFormat, &unset.TimeOutputFormat),
		handleParameterUpdateWithMapping(d, defs.TraceLevel.FieldName(), &set.TraceLevel, &unset.TraceLevel, sdk.ToTraceLevel),
		handleParameterUpdate(d, defs.TransactionAbortOnError.FieldName(), &set.TransactionAbortOnError, &unset.TransactionAbortOnError),
		handleParameterUpdateWithMapping(d, defs.TransactionDefaultIsolationLevel.FieldName(), &set.TransactionDefaultIsolationLevel, &unset.TransactionDefaultIsolationLevel, sdk.ToTransactionDefaultIsolationLevel),
		handleParameterUpdate(d, defs.TwoDigitCenturyStart.FieldName(), &set.TwoDigitCenturyStart, &unset.TwoDigitCenturyStart),
		handleParameterUpdateWithMapping(d, defs.UnsupportedDdlAction.FieldName(), &set.UnsupportedDdlAction, &unset.UnsupportedDdlAction, sdk.ToUnsupportedDDLAction),
		handleParameterUpdate(d, defs.UseCachedResult.FieldName(), &set.UseCachedResult, &unset.UseCachedResult),
		handleParameterUpdate(d, defs.WeekOfYearPolicy.FieldName(), &set.WeekOfYearPolicy, &unset.WeekOfYearPolicy),
		handleParameterUpdate(d, defs.WeekStart.FieldName(), &set.WeekStart, &unset.WeekStart),
		handleParameterUpdate(d, defs.EnableUnredactedQuerySyntaxError.FieldName(), &set.EnableUnredactedQuerySyntaxError, &unset.EnableUnredactedQuerySyntaxError),
		handleParameterUpdateWithMapping(d, defs.NetworkPolicy.FieldName(), &set.NetworkPolicy, &unset.NetworkPolicy, stringToAccountObjectIdentifier),
		handleParameterUpdate(d, defs.PreventUnloadToInternalStages.FieldName(), &set.PreventUnloadToInternalStages, &unset.PreventUnloadToInternalStages),
	)
}

func handleUserParameterRead(d *schema.ResourceData, parameters *sdk.UserParametersDetails) diag.Diagnostics {
	return JoinDiags(
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
		setResourceData(d, defs.EnableUnredactedQuerySyntaxError.FieldName(), parameters.EnableUnredactedQuerySyntaxError.Value),
		setResourceData(d, defs.ErrorOnNondeterministicMerge.FieldName(), parameters.ErrorOnNondeterministicMerge.Value),
		setResourceData(d, defs.ErrorOnNondeterministicUpdate.FieldName(), parameters.ErrorOnNondeterministicUpdate.Value),
		setResourceData(d, defs.GeographyOutputFormat.FieldName(), parameters.GeographyOutputFormat.Value),
		setResourceData(d, defs.GeometryOutputFormat.FieldName(), parameters.GeometryOutputFormat.Value),
		setResourceData(d, defs.JdbcTreatDecimalAsInt.FieldName(), parameters.JdbcTreatDecimalAsInt.Value),
		setResourceData(d, defs.JdbcTreatTimestampNtzAsUtc.FieldName(), parameters.JdbcTreatTimestampNtzAsUtc.Value),
		setResourceData(d, defs.JdbcUseSessionTimezone.FieldName(), parameters.JdbcUseSessionTimezone.Value),
		setResourceData(d, defs.JsonIndent.FieldName(), parameters.JsonIndent.Value),
		setResourceData(d, defs.LockTimeout.FieldName(), parameters.LockTimeout.Value),
		setResourceData(d, defs.LogEventLevel.FieldName(), parameters.LogEventLevel.Value),
		setResourceData(d, defs.LogLevel.FieldName(), parameters.LogLevel.Value),
		setResourceData(d, defs.MultiStatementCount.FieldName(), parameters.MultiStatementCount.Value),
		setResourceData(d, defs.NetworkPolicy.FieldName(), parameters.NetworkPolicy.Value.FullyQualifiedName()),
		setResourceData(d, defs.NoorderSequenceAsDefault.FieldName(), parameters.NoorderSequenceAsDefault.Value),
		setResourceData(d, defs.OdbcTreatDecimalAsInt.FieldName(), parameters.OdbcTreatDecimalAsInt.Value),
		setResourceData(d, defs.PreventUnloadToInternalStages.FieldName(), parameters.PreventUnloadToInternalStages.Value),
		setResourceData(d, defs.QueryTag.FieldName(), parameters.QueryTag.Value),
		setResourceData(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), parameters.QuotedIdentifiersIgnoreCase.Value),
		setResourceData(d, defs.RowsPerResultset.FieldName(), parameters.RowsPerResultset.Value),
		setResourceData(d, defs.S3StageVpceDnsName.FieldName(), parameters.S3StageVpceDnsName.Value),
		setResourceData(d, defs.SearchPath.FieldName(), parameters.SearchPath.Value),
		setResourceData(d, defs.SimulatedDataSharingConsumer.FieldName(), parameters.SimulatedDataSharingConsumer.Value),
		setResourceData(d, defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds.Value),
		setResourceData(d, defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds.Value),
		setResourceData(d, defs.StrictJsonOutput.FieldName(), parameters.StrictJsonOutput.Value),
		setResourceData(d, defs.TimeInputFormat.FieldName(), parameters.TimeInputFormat.Value),
		setResourceData(d, defs.TimeOutputFormat.FieldName(), parameters.TimeOutputFormat.Value),
		setResourceData(d, defs.TimestampDayIsAlways24h.FieldName(), parameters.TimestampDayIsAlways24H.Value),
		setResourceData(d, defs.TimestampInputFormat.FieldName(), parameters.TimestampInputFormat.Value),
		setResourceData(d, defs.TimestampLtzOutputFormat.FieldName(), parameters.TimestampLtzOutputFormat.Value),
		setResourceData(d, defs.TimestampNtzOutputFormat.FieldName(), parameters.TimestampNtzOutputFormat.Value),
		setResourceData(d, defs.TimestampOutputFormat.FieldName(), parameters.TimestampOutputFormat.Value),
		setResourceData(d, defs.TimestampTypeMapping.FieldName(), parameters.TimestampTypeMapping.Value),
		setResourceData(d, defs.TimestampTzOutputFormat.FieldName(), parameters.TimestampTzOutputFormat.Value),
		setResourceData(d, defs.Timezone.FieldName(), parameters.Timezone.Value),
		setResourceData(d, defs.TraceLevel.FieldName(), parameters.TraceLevel.Value),
		setResourceData(d, defs.TransactionAbortOnError.FieldName(), parameters.TransactionAbortOnError.Value),
		setResourceData(d, defs.TransactionDefaultIsolationLevel.FieldName(), parameters.TransactionDefaultIsolationLevel.Value),
		setResourceData(d, defs.TwoDigitCenturyStart.FieldName(), parameters.TwoDigitCenturyStart.Value),
		setResourceData(d, defs.UnsupportedDdlAction.FieldName(), parameters.UnsupportedDdlAction.Value),
		setResourceData(d, defs.UseCachedResult.FieldName(), parameters.UseCachedResult.Value),
		setResourceData(d, defs.WeekOfYearPolicy.FieldName(), parameters.WeekOfYearPolicy.Value),
		setResourceData(d, defs.WeekStart.FieldName(), parameters.WeekStart.Value),
	)
}

var userParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	userParametersProvider,
	userParameterDiffFunctions,
)

func userParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.UserParametersDetails, error) {
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.Users.ShowParametersDetails(ctx, id)
}

func userParameterDiffFunctions(parameters *sdk.UserParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		BoolTypedParameterValueComputedIf(defs.AbortDetachedQuery.FieldName(), parameters.AbortDetachedQuery, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.Autocommit.FieldName(), parameters.Autocommit, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.BinaryInputFormat.FieldName(), parameters.BinaryInputFormat, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.BinaryOutputFormat.FieldName(), parameters.BinaryOutputFormat, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.ClientMemoryLimit.FieldName(), parameters.ClientMemoryLimit, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.ClientMetadataRequestUseConnectionCtx.FieldName(), parameters.ClientMetadataRequestUseConnectionCtx, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.ClientPrefetchThreads.FieldName(), parameters.ClientPrefetchThreads, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.ClientResultChunkSize.FieldName(), parameters.ClientResultChunkSize, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.ClientResultColumnCaseInsensitive.FieldName(), parameters.ClientResultColumnCaseInsensitive, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.ClientSessionKeepAlive.FieldName(), parameters.ClientSessionKeepAlive, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.ClientSessionKeepAliveHeartbeatFrequency.FieldName(), parameters.ClientSessionKeepAliveHeartbeatFrequency, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.ClientTimestampTypeMapping.FieldName(), parameters.ClientTimestampTypeMapping, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.DateInputFormat.FieldName(), parameters.DateInputFormat, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.DateOutputFormat.FieldName(), parameters.DateOutputFormat, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.EnableUnloadPhysicalTypeOptimization.FieldName(), parameters.EnableUnloadPhysicalTypeOptimization, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.EnableUnredactedQuerySyntaxError.FieldName(), parameters.EnableUnredactedQuerySyntaxError, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.ErrorOnNondeterministicMerge.FieldName(), parameters.ErrorOnNondeterministicMerge, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.ErrorOnNondeterministicUpdate.FieldName(), parameters.ErrorOnNondeterministicUpdate, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.GeographyOutputFormat.FieldName(), parameters.GeographyOutputFormat, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.GeometryOutputFormat.FieldName(), parameters.GeometryOutputFormat, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.JdbcTreatDecimalAsInt.FieldName(), parameters.JdbcTreatDecimalAsInt, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.JdbcTreatTimestampNtzAsUtc.FieldName(), parameters.JdbcTreatTimestampNtzAsUtc, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.JdbcUseSessionTimezone.FieldName(), parameters.JdbcUseSessionTimezone, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.JsonIndent.FieldName(), parameters.JsonIndent, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.LockTimeout.FieldName(), parameters.LockTimeout, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.LogEventLevel.FieldName(), parameters.LogEventLevel, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.LogLevel.FieldName(), parameters.LogLevel, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.MultiStatementCount.FieldName(), parameters.MultiStatementCount, sdk.ParameterTypeUser),
		IdentifierTypedParameterValueComputedIf(defs.NetworkPolicy.FieldName(), parameters.NetworkPolicy, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.NoorderSequenceAsDefault.FieldName(), parameters.NoorderSequenceAsDefault, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.OdbcTreatDecimalAsInt.FieldName(), parameters.OdbcTreatDecimalAsInt, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.PreventUnloadToInternalStages.FieldName(), parameters.PreventUnloadToInternalStages, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.QueryTag.FieldName(), parameters.QueryTag, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.QuotedIdentifiersIgnoreCase.FieldName(), parameters.QuotedIdentifiersIgnoreCase, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.RowsPerResultset.FieldName(), parameters.RowsPerResultset, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.S3StageVpceDnsName.FieldName(), parameters.S3StageVpceDnsName, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.SearchPath.FieldName(), parameters.SearchPath, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.SimulatedDataSharingConsumer.FieldName(), parameters.SimulatedDataSharingConsumer, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.StatementQueuedTimeoutInSeconds.FieldName(), parameters.StatementQueuedTimeoutInSeconds, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.StatementTimeoutInSeconds.FieldName(), parameters.StatementTimeoutInSeconds, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.StrictJsonOutput.FieldName(), parameters.StrictJsonOutput, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TimeInputFormat.FieldName(), parameters.TimeInputFormat, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TimeOutputFormat.FieldName(), parameters.TimeOutputFormat, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.TimestampDayIsAlways24h.FieldName(), parameters.TimestampDayIsAlways24H, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TimestampInputFormat.FieldName(), parameters.TimestampInputFormat, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TimestampLtzOutputFormat.FieldName(), parameters.TimestampLtzOutputFormat, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TimestampNtzOutputFormat.FieldName(), parameters.TimestampNtzOutputFormat, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TimestampOutputFormat.FieldName(), parameters.TimestampOutputFormat, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TimestampTypeMapping.FieldName(), parameters.TimestampTypeMapping, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TimestampTzOutputFormat.FieldName(), parameters.TimestampTzOutputFormat, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.Timezone.FieldName(), parameters.Timezone, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TraceLevel.FieldName(), parameters.TraceLevel, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.TransactionAbortOnError.FieldName(), parameters.TransactionAbortOnError, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.TransactionDefaultIsolationLevel.FieldName(), parameters.TransactionDefaultIsolationLevel, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.TwoDigitCenturyStart.FieldName(), parameters.TwoDigitCenturyStart, sdk.ParameterTypeUser),
		StringTypedParameterValueComputedIf(defs.UnsupportedDdlAction.FieldName(), parameters.UnsupportedDdlAction, sdk.ParameterTypeUser),
		BoolTypedParameterValueComputedIf(defs.UseCachedResult.FieldName(), parameters.UseCachedResult, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.WeekOfYearPolicy.FieldName(), parameters.WeekOfYearPolicy, sdk.ParameterTypeUser),
		IntTypedParameterValueComputedIf(defs.WeekStart.FieldName(), parameters.WeekStart, sdk.ParameterTypeUser),
	}
}
