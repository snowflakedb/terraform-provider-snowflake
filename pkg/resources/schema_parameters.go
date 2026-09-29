package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var schemaParametersSchema = make(map[string]*schema.Schema)

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelSchema) {
		schemaParametersSchema[p.FieldName()] = parameterSchema(p)
	}
}

func handleSchemaParametersCreate(d *schema.ResourceData, req *sdk.CreateSchemaRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreateWithMapping(d, defs.Catalog.FieldName(), &req.Catalog, stringToAccountObjectIdentifier),
		handleParameterCreate(d, defs.DataRetentionTimeInDays.FieldName(), &req.DataRetentionTimeInDays),
		handleParameterCreateWithMapping(d, defs.DefaultDdlCollation.FieldName(), &req.DefaultDdlCollation, func(value string) (sdk.StringAllowEmpty, error) { return sdk.StringAllowEmpty{Value: value}, nil }),
		handleParameterCreate(d, defs.DefaultNotebookComputePoolCpu.FieldName(), &req.DefaultNotebookComputePoolCpu),
		handleParameterCreate(d, defs.DefaultNotebookComputePoolGpu.FieldName(), &req.DefaultNotebookComputePoolGpu),
		handleParameterCreate(d, defs.EnableConsoleOutput.FieldName(), &req.EnableConsoleOutput),
		handleParameterCreateWithMapping(d, defs.ExternalVolume.FieldName(), &req.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, defs.LogEventLevel.FieldName(), &req.LogEventLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, defs.LogLevel.FieldName(), &req.LogLevel, sdk.ToLogLevel),
		handleParameterCreate(d, defs.MaxDataExtensionTimeInDays.FieldName(), &req.MaxDataExtensionTimeInDays),
		handleParameterCreate(d, defs.PipeExecutionPaused.FieldName(), &req.PipeExecutionPaused),
		handleParameterCreate(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), &req.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, defs.ReplaceInvalidCharacters.FieldName(), &req.ReplaceInvalidCharacters),
		handleParameterCreateWithMapping(d, defs.StorageSerializationPolicy.FieldName(), &req.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterCreate(d, defs.SuspendTaskAfterNumFailures.FieldName(), &req.SuspendTaskAfterNumFailures),
		handleParameterCreate(d, defs.TaskAutoRetryAttempts.FieldName(), &req.TaskAutoRetryAttempts),
		handleParameterCreateWithMapping(d, defs.TraceLevel.FieldName(), &req.TraceLevel, sdk.ToTraceLevel),
		handleParameterCreateWithMapping(d, defs.UserTaskManagedInitialWarehouseSize.FieldName(), &req.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterCreate(d, defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), &req.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterCreate(d, defs.UserTaskTimeoutMs.FieldName(), &req.UserTaskTimeoutMs),
	)
}

func handleSchemaParametersChanges(d *schema.ResourceData, set *sdk.SchemaSetRequest, unset *sdk.SchemaUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdateWithMapping(d, defs.Catalog.FieldName(), &set.Catalog, &unset.Catalog, stringToAccountObjectIdentifier),
		handleParameterUpdate(d, defs.DataRetentionTimeInDays.FieldName(), &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleParameterUpdateWithMapping(d, defs.DefaultDdlCollation.FieldName(), &set.DefaultDdlCollation, &unset.DefaultDdlCollation, func(value string) (sdk.StringAllowEmpty, error) { return sdk.StringAllowEmpty{Value: value}, nil }),
		handleParameterUpdate(d, defs.DefaultNotebookComputePoolCpu.FieldName(), &set.DefaultNotebookComputePoolCpu, &unset.DefaultNotebookComputePoolCpu),
		handleParameterUpdate(d, defs.DefaultNotebookComputePoolGpu.FieldName(), &set.DefaultNotebookComputePoolGpu, &unset.DefaultNotebookComputePoolGpu),
		handleParameterUpdate(d, defs.EnableConsoleOutput.FieldName(), &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
		handleParameterUpdateWithMapping(d, defs.ExternalVolume.FieldName(), &set.ExternalVolume, &unset.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterUpdateWithMapping(d, defs.LogEventLevel.FieldName(), &set.LogEventLevel, &unset.LogEventLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, defs.LogLevel.FieldName(), &set.LogLevel, &unset.LogLevel, sdk.ToLogLevel),
		handleParameterUpdate(d, defs.MaxDataExtensionTimeInDays.FieldName(), &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleParameterUpdate(d, defs.PipeExecutionPaused.FieldName(), &set.PipeExecutionPaused, &unset.PipeExecutionPaused),
		handleParameterUpdate(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, defs.ReplaceInvalidCharacters.FieldName(), &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleParameterUpdateWithMapping(d, defs.StorageSerializationPolicy.FieldName(), &set.StorageSerializationPolicy, &unset.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterUpdate(d, defs.SuspendTaskAfterNumFailures.FieldName(), &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleParameterUpdate(d, defs.TaskAutoRetryAttempts.FieldName(), &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleParameterUpdateWithMapping(d, defs.TraceLevel.FieldName(), &set.TraceLevel, &unset.TraceLevel, sdk.ToTraceLevel),
		handleParameterUpdateWithMapping(d, defs.UserTaskManagedInitialWarehouseSize.FieldName(), &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterUpdate(d, defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterUpdate(d, defs.UserTaskTimeoutMs.FieldName(), &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
	)
}

func handleSchemaParameterRead(d *schema.ResourceData, parameters *sdk.SchemaParametersDetails) diag.Diagnostics {
	set := func(key string, value any) diag.Diagnostics {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	return JoinDiags(
		set(defs.Catalog.FieldName(), parameters.Catalog.Value.FullyQualifiedName()),
		set(defs.DataRetentionTimeInDays.FieldName(), parameters.DataRetentionTimeInDays.Value),
		set(defs.DefaultDdlCollation.FieldName(), parameters.DefaultDdlCollation.Value),
		set(defs.DefaultNotebookComputePoolCpu.FieldName(), parameters.DefaultNotebookComputePoolCpu.Value),
		set(defs.DefaultNotebookComputePoolGpu.FieldName(), parameters.DefaultNotebookComputePoolGpu.Value),
		set(defs.EnableConsoleOutput.FieldName(), parameters.EnableConsoleOutput.Value),
		set(defs.ExternalVolume.FieldName(), parameters.ExternalVolume.Value.FullyQualifiedName()),
		set(defs.LogEventLevel.FieldName(), parameters.LogEventLevel.Value),
		set(defs.LogLevel.FieldName(), parameters.LogLevel.Value),
		set(defs.MaxDataExtensionTimeInDays.FieldName(), parameters.MaxDataExtensionTimeInDays.Value),
		set(defs.PipeExecutionPaused.FieldName(), parameters.PipeExecutionPaused.Value),
		set(defs.QuotedIdentifiersIgnoreCase.FieldName(), parameters.QuotedIdentifiersIgnoreCase.Value),
		set(defs.ReplaceInvalidCharacters.FieldName(), parameters.ReplaceInvalidCharacters.Value),
		set(defs.StorageSerializationPolicy.FieldName(), parameters.StorageSerializationPolicy.Value),
		set(defs.SuspendTaskAfterNumFailures.FieldName(), parameters.SuspendTaskAfterNumFailures.Value),
		set(defs.TaskAutoRetryAttempts.FieldName(), parameters.TaskAutoRetryAttempts.Value),
		set(defs.TraceLevel.FieldName(), parameters.TraceLevel.Value),
		set(defs.UserTaskManagedInitialWarehouseSize.FieldName(), parameters.UserTaskManagedInitialWarehouseSize.Value),
		set(defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), parameters.UserTaskMinimumTriggerIntervalInSeconds.Value),
		set(defs.UserTaskTimeoutMs.FieldName(), parameters.UserTaskTimeoutMs.Value),
	)
}

var schemaParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	schemaParametersProvider,
	schemaParameterDiffFunctions,
)

func schemaParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.SchemaParametersDetails, error) {
	id, err := sdk.ParseDatabaseObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.Schemas.ShowParametersDetails(ctx, id)
}

func schemaParameterDiffFunctions(parameters *sdk.SchemaParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IdentifierTypedParameterValueComputedIf(defs.Catalog.FieldName(), parameters.Catalog, sdk.ParameterTypeSchema),
		IntTypedParameterValueComputedIf(defs.DataRetentionTimeInDays.FieldName(), parameters.DataRetentionTimeInDays, sdk.ParameterTypeSchema),
		StringTypedParameterValueComputedIf(defs.DefaultDdlCollation.FieldName(), parameters.DefaultDdlCollation, sdk.ParameterTypeSchema),
		StringTypedParameterValueComputedIf(defs.DefaultNotebookComputePoolCpu.FieldName(), parameters.DefaultNotebookComputePoolCpu, sdk.ParameterTypeSchema),
		StringTypedParameterValueComputedIf(defs.DefaultNotebookComputePoolGpu.FieldName(), parameters.DefaultNotebookComputePoolGpu, sdk.ParameterTypeSchema),
		BoolTypedParameterValueComputedIf(defs.EnableConsoleOutput.FieldName(), parameters.EnableConsoleOutput, sdk.ParameterTypeSchema),
		IdentifierTypedParameterValueComputedIf(defs.ExternalVolume.FieldName(), parameters.ExternalVolume, sdk.ParameterTypeSchema),
		StringTypedParameterValueComputedIf(defs.LogEventLevel.FieldName(), parameters.LogEventLevel, sdk.ParameterTypeSchema),
		StringTypedParameterValueComputedIf(defs.LogLevel.FieldName(), parameters.LogLevel, sdk.ParameterTypeSchema),
		IntTypedParameterValueComputedIf(defs.MaxDataExtensionTimeInDays.FieldName(), parameters.MaxDataExtensionTimeInDays, sdk.ParameterTypeSchema),
		BoolTypedParameterValueComputedIf(defs.PipeExecutionPaused.FieldName(), parameters.PipeExecutionPaused, sdk.ParameterTypeSchema),
		BoolTypedParameterValueComputedIf(defs.QuotedIdentifiersIgnoreCase.FieldName(), parameters.QuotedIdentifiersIgnoreCase, sdk.ParameterTypeSchema),
		BoolTypedParameterValueComputedIf(defs.ReplaceInvalidCharacters.FieldName(), parameters.ReplaceInvalidCharacters, sdk.ParameterTypeSchema),
		StringTypedParameterValueComputedIf(defs.StorageSerializationPolicy.FieldName(), parameters.StorageSerializationPolicy, sdk.ParameterTypeSchema),
		IntTypedParameterValueComputedIf(defs.SuspendTaskAfterNumFailures.FieldName(), parameters.SuspendTaskAfterNumFailures, sdk.ParameterTypeSchema),
		IntTypedParameterValueComputedIf(defs.TaskAutoRetryAttempts.FieldName(), parameters.TaskAutoRetryAttempts, sdk.ParameterTypeSchema),
		StringTypedParameterValueComputedIf(defs.TraceLevel.FieldName(), parameters.TraceLevel, sdk.ParameterTypeSchema),
		StringTypedParameterValueComputedIf(defs.UserTaskManagedInitialWarehouseSize.FieldName(), parameters.UserTaskManagedInitialWarehouseSize, sdk.ParameterTypeSchema),
		IntTypedParameterValueComputedIf(defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), parameters.UserTaskMinimumTriggerIntervalInSeconds, sdk.ParameterTypeSchema),
		IntTypedParameterValueComputedIf(defs.UserTaskTimeoutMs.FieldName(), parameters.UserTaskTimeoutMs, sdk.ParameterTypeSchema),
	}
}
