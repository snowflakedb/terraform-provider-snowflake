package resources

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/schemas"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// =============================================================================
// Schema maps
// =============================================================================

var schemaParametersAttributesSchemaExt = func() map[string]*schema.Schema {
	m := make(map[string]*schema.Schema)
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelSchema) {
		m[p.FieldName()] = parameterSchema(p)
	}
	return m
}()

var schemaParametersOutputSchemaExt = map[string]*schema.Schema{
	ParametersAttributeName: {
		Type:        schema.TypeList,
		Computed:    true,
		Description: "Outputs the result of `SHOW PARAMETERS IN SCHEMA` for the given object.",
		Elem: &schema.Resource{
			Schema: schemas.ShowSchemaParametersSchema,
		},
	},
}

var schemaParametersFieldNamesExt = collections.Map(
	defs.ParameterDefsForLevel(parameterdefs.ParameterLevelSchema),
	func(p parameterdefs.ParameterDef) string { return p.FieldName() },
)

var schemaParametersCustomDiffExt = ParametersCustomDiffFromTypedParameters(
	schemaParametersProviderExt,
	schemaParameterDiffFunctionsExt,
)

func schemaParametersProviderExt(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.SchemaParametersDetails, error) {
	id, err := schemaParseIdExt(d.Id())
	if err != nil {
		return nil, err
	}
	return schemaShowParametersDetailsInSdkExt(ctx, meta, id)
}

func schemaParameterDiffFunctionsExt(parameters *sdk.SchemaParametersDetails) []schema.CustomizeDiffFunc {
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

// =============================================================================
// Create / Read / Update / CustomizeDiff handlers
// =============================================================================

func schemaApplyParametersCreateExt(d *schema.ResourceData, req *sdk.CreateSchemaRequest) diag.Diagnostics {
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

func schemaApplyParametersChangesExt(d *schema.ResourceData, reqs *schemaAlterRequestsExt) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdateWithMapping(d, defs.Catalog.FieldName(), &reqs.set.Catalog, &reqs.unset.Catalog, stringToAccountObjectIdentifier),
		handleParameterUpdate(d, defs.DataRetentionTimeInDays.FieldName(), &reqs.set.DataRetentionTimeInDays, &reqs.unset.DataRetentionTimeInDays),
		handleParameterUpdateWithMapping(d, defs.DefaultDdlCollation.FieldName(), &reqs.set.DefaultDdlCollation, &reqs.unset.DefaultDdlCollation, func(value string) (sdk.StringAllowEmpty, error) { return sdk.StringAllowEmpty{Value: value}, nil }),
		handleParameterUpdate(d, defs.DefaultNotebookComputePoolCpu.FieldName(), &reqs.set.DefaultNotebookComputePoolCpu, &reqs.unset.DefaultNotebookComputePoolCpu),
		handleParameterUpdate(d, defs.DefaultNotebookComputePoolGpu.FieldName(), &reqs.set.DefaultNotebookComputePoolGpu, &reqs.unset.DefaultNotebookComputePoolGpu),
		handleParameterUpdate(d, defs.EnableConsoleOutput.FieldName(), &reqs.set.EnableConsoleOutput, &reqs.unset.EnableConsoleOutput),
		handleParameterUpdateWithMapping(d, defs.ExternalVolume.FieldName(), &reqs.set.ExternalVolume, &reqs.unset.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterUpdateWithMapping(d, defs.LogEventLevel.FieldName(), &reqs.set.LogEventLevel, &reqs.unset.LogEventLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, defs.LogLevel.FieldName(), &reqs.set.LogLevel, &reqs.unset.LogLevel, sdk.ToLogLevel),
		handleParameterUpdate(d, defs.MaxDataExtensionTimeInDays.FieldName(), &reqs.set.MaxDataExtensionTimeInDays, &reqs.unset.MaxDataExtensionTimeInDays),
		handleParameterUpdate(d, defs.PipeExecutionPaused.FieldName(), &reqs.set.PipeExecutionPaused, &reqs.unset.PipeExecutionPaused),
		handleParameterUpdate(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), &reqs.set.QuotedIdentifiersIgnoreCase, &reqs.unset.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, defs.ReplaceInvalidCharacters.FieldName(), &reqs.set.ReplaceInvalidCharacters, &reqs.unset.ReplaceInvalidCharacters),
		handleParameterUpdateWithMapping(d, defs.StorageSerializationPolicy.FieldName(), &reqs.set.StorageSerializationPolicy, &reqs.unset.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterUpdate(d, defs.SuspendTaskAfterNumFailures.FieldName(), &reqs.set.SuspendTaskAfterNumFailures, &reqs.unset.SuspendTaskAfterNumFailures),
		handleParameterUpdate(d, defs.TaskAutoRetryAttempts.FieldName(), &reqs.set.TaskAutoRetryAttempts, &reqs.unset.TaskAutoRetryAttempts),
		handleParameterUpdateWithMapping(d, defs.TraceLevel.FieldName(), &reqs.set.TraceLevel, &reqs.unset.TraceLevel, sdk.ToTraceLevel),
		handleParameterUpdateWithMapping(d, defs.UserTaskManagedInitialWarehouseSize.FieldName(), &reqs.set.UserTaskManagedInitialWarehouseSize, &reqs.unset.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterUpdate(d, defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), &reqs.set.UserTaskMinimumTriggerIntervalInSeconds, &reqs.unset.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterUpdate(d, defs.UserTaskTimeoutMs.FieldName(), &reqs.set.UserTaskTimeoutMs, &reqs.unset.UserTaskTimeoutMs),
	)
}

func schemaSetParametersFieldsExt(d *schema.ResourceData, parameters *sdk.SchemaParametersDetails) diag.Diagnostics {
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

// =============================================================================
// SDK incision points
// =============================================================================

func schemaShowParametersInSdkExt(ctx context.Context, meta any, id sdk.DatabaseObjectIdentifier) ([]*sdk.Parameter, error) {
	client := meta.(*provider.Context).Client
	return client.Schemas.ShowParameters(ctx, id)
}

func schemaShowParametersDetailsInSdkExt(ctx context.Context, meta any, id sdk.DatabaseObjectIdentifier) (*sdk.SchemaParametersDetails, error) {
	client := meta.(*provider.Context).Client
	return client.Schemas.ShowParametersDetails(ctx, id)
}
