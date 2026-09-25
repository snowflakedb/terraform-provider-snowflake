package resources

import (
	"context"
	"slices"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/parameterdefs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	databaseParametersSchema       = make(map[string]*schema.Schema)
	sharedDatabaseParametersSchema = make(map[string]*schema.Schema)
)

var databaseParameterIntRanges = map[string]schema.SchemaValidateDiagFunc{
	// Choosing a higher range than the documented one (for the standard edition or transient databases, the maximum number is 1).
	defs.DataRetentionTimeInDays.SqlName:     validation.ToDiagFunc(validation.IntBetween(0, 90)),
	defs.MaxDataExtensionTimeInDays.SqlName:  validation.ToDiagFunc(validation.IntBetween(0, 90)),
	defs.UserTaskTimeoutMs.SqlName:           validation.ToDiagFunc(validation.IntBetween(0, 86400000)),
	defs.SuspendTaskAfterNumFailures.SqlName: validation.ToDiagFunc(validation.IntAtLeast(0)),
	defs.TaskAutoRetryAttempts.SqlName:       validation.ToDiagFunc(validation.IntAtLeast(0)),
}

var databaseValidAccountObjectIdentifierParameters = []parameterdefs.ParameterDef{
	defs.DefaultNotebookComputePoolCpu,
	defs.DefaultNotebookComputePoolGpu,
}

var sharedDatabaseNotApplicableParameters = []string{
	defs.DataRetentionTimeInDays.SqlName,
	defs.MaxDataExtensionTimeInDays.SqlName,
}

func databaseParameterSchemaFor(p parameterdefs.ParameterDef) *schema.Schema {
	s := parameterSchema(p)

	if validate, ok := databaseParameterIntRanges[p.SqlName]; ok {
		s.ValidateDiagFunc = validate
	}
	if slices.ContainsFunc(databaseValidAccountObjectIdentifierParameters, func(validIdentifierParameter parameterdefs.ParameterDef) bool {
		return validIdentifierParameter.SqlName == p.SqlName
	}) {
		s.ValidateDiagFunc = IsValidIdentifier[sdk.AccountObjectIdentifier]()
		s.DiffSuppressFunc = suppressIdentifierQuoting
	}

	return s
}

func init() {
	for _, p := range defs.ParameterDefsForLevel(parameterdefs.ParameterLevelDatabase) {
		databaseParametersSchema[p.FieldName()] = databaseParameterSchemaFor(p)

		if !slices.Contains(sharedDatabaseNotApplicableParameters, p.SqlName) {
			sharedSchema := databaseParameterSchemaFor(p)
			sharedSchema.ForceNew = true
			sharedDatabaseParametersSchema[p.FieldName()] = sharedSchema
		}
	}
}

func handleDatabaseParametersCreate(d *schema.ResourceData, createOpts *sdk.CreateDatabaseRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, defs.DataRetentionTimeInDays.FieldName(), &createOpts.DataRetentionTimeInDays),
		handleParameterCreate(d, defs.MaxDataExtensionTimeInDays.FieldName(), &createOpts.MaxDataExtensionTimeInDays),
		handleParameterCreateWithMapping(d, defs.ExternalVolume.FieldName(), &createOpts.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, defs.Catalog.FieldName(), &createOpts.Catalog, stringToAccountObjectIdentifier),
		handleParameterCreate(d, defs.ReplaceInvalidCharacters.FieldName(), &createOpts.ReplaceInvalidCharacters),
		handleParameterCreateWithMapping(d, defs.DefaultDdlCollation.FieldName(), &createOpts.DefaultDdlCollation, func(value string) (sdk.StringAllowEmpty, error) { return sdk.StringAllowEmpty{Value: value}, nil }),
		handleParameterCreate(d, defs.DefaultNotebookComputePoolCpu.FieldName(), &createOpts.DefaultNotebookComputePoolCpu),
		handleParameterCreate(d, defs.DefaultNotebookComputePoolGpu.FieldName(), &createOpts.DefaultNotebookComputePoolGpu),
		handleParameterCreateWithMapping(d, defs.StorageSerializationPolicy.FieldName(), &createOpts.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterCreateWithMapping(d, defs.LogLevel.FieldName(), &createOpts.LogLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, defs.LogEventLevel.FieldName(), &createOpts.LogEventLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, defs.TraceLevel.FieldName(), &createOpts.TraceLevel, sdk.ToTraceLevel),
		handleParameterCreate(d, defs.SuspendTaskAfterNumFailures.FieldName(), &createOpts.SuspendTaskAfterNumFailures),
		handleParameterCreate(d, defs.TaskAutoRetryAttempts.FieldName(), &createOpts.TaskAutoRetryAttempts),
		handleParameterCreateWithMapping(d, defs.UserTaskManagedInitialWarehouseSize.FieldName(), &createOpts.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterCreate(d, defs.UserTaskTimeoutMs.FieldName(), &createOpts.UserTaskTimeoutMs),
		handleParameterCreate(d, defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), &createOpts.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterCreate(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), &createOpts.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, defs.EnableConsoleOutput.FieldName(), &createOpts.EnableConsoleOutput),
	)
}

func handleDatabaseParametersUpdate(d *schema.ResourceData, set *sdk.DatabaseSetRequest, unset *sdk.DatabaseUnsetRequest) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, defs.DataRetentionTimeInDays.FieldName(), &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleParameterUpdate(d, defs.MaxDataExtensionTimeInDays.FieldName(), &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleParameterUpdateWithMapping(d, defs.ExternalVolume.FieldName(), &set.ExternalVolume, &unset.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterUpdateWithMapping(d, defs.Catalog.FieldName(), &set.Catalog, &unset.Catalog, stringToAccountObjectIdentifier),
		handleParameterUpdate(d, defs.ReplaceInvalidCharacters.FieldName(), &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleParameterUpdateWithMapping(d, defs.DefaultDdlCollation.FieldName(), &set.DefaultDdlCollation, &unset.DefaultDdlCollation, func(value string) (sdk.StringAllowEmpty, error) { return sdk.StringAllowEmpty{Value: value}, nil }),
		handleParameterUpdate(d, defs.DefaultNotebookComputePoolCpu.FieldName(), &set.DefaultNotebookComputePoolCpu, &unset.DefaultNotebookComputePoolCpu),
		handleParameterUpdate(d, defs.DefaultNotebookComputePoolGpu.FieldName(), &set.DefaultNotebookComputePoolGpu, &unset.DefaultNotebookComputePoolGpu),
		handleParameterUpdateWithMapping(d, defs.StorageSerializationPolicy.FieldName(), &set.StorageSerializationPolicy, &unset.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterUpdateWithMapping(d, defs.LogLevel.FieldName(), &set.LogLevel, &unset.LogLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, defs.LogEventLevel.FieldName(), &set.LogEventLevel, &unset.LogEventLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, defs.TraceLevel.FieldName(), &set.TraceLevel, &unset.TraceLevel, sdk.ToTraceLevel),
		handleParameterUpdate(d, defs.SuspendTaskAfterNumFailures.FieldName(), &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleParameterUpdate(d, defs.TaskAutoRetryAttempts.FieldName(), &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleParameterUpdateWithMapping(d, defs.UserTaskManagedInitialWarehouseSize.FieldName(), &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterUpdate(d, defs.UserTaskTimeoutMs.FieldName(), &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
		handleParameterUpdate(d, defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterUpdate(d, defs.QuotedIdentifiersIgnoreCase.FieldName(), &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, defs.EnableConsoleOutput.FieldName(), &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
	)
}

func handleDatabaseParameterRead(d *schema.ResourceData, databaseParameters *sdk.DatabaseParametersDetails) diag.Diagnostics {
	set := func(key string, value any) diag.Diagnostics {
		if err := d.Set(key, value); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	return JoinDiags(
		set(defs.Catalog.FieldName(), databaseParameters.Catalog.Value.FullyQualifiedName()),
		set(defs.DataRetentionTimeInDays.FieldName(), databaseParameters.DataRetentionTimeInDays.Value),
		set(defs.DefaultDdlCollation.FieldName(), databaseParameters.DefaultDdlCollation.Value),
		set(defs.DefaultNotebookComputePoolCpu.FieldName(), databaseParameters.DefaultNotebookComputePoolCpu.Value),
		set(defs.DefaultNotebookComputePoolGpu.FieldName(), databaseParameters.DefaultNotebookComputePoolGpu.Value),
		set(defs.EnableConsoleOutput.FieldName(), databaseParameters.EnableConsoleOutput.Value),
		set(defs.ExternalVolume.FieldName(), databaseParameters.ExternalVolume.Value.FullyQualifiedName()),
		set(defs.LogEventLevel.FieldName(), databaseParameters.LogEventLevel.Value),
		set(defs.LogLevel.FieldName(), databaseParameters.LogLevel.Value),
		set(defs.MaxDataExtensionTimeInDays.FieldName(), databaseParameters.MaxDataExtensionTimeInDays.Value),
		set(defs.QuotedIdentifiersIgnoreCase.FieldName(), databaseParameters.QuotedIdentifiersIgnoreCase.Value),
		set(defs.ReplaceInvalidCharacters.FieldName(), databaseParameters.ReplaceInvalidCharacters.Value),
		set(defs.StorageSerializationPolicy.FieldName(), databaseParameters.StorageSerializationPolicy.Value),
		set(defs.SuspendTaskAfterNumFailures.FieldName(), databaseParameters.SuspendTaskAfterNumFailures.Value),
		set(defs.TaskAutoRetryAttempts.FieldName(), databaseParameters.TaskAutoRetryAttempts.Value),
		set(defs.TraceLevel.FieldName(), databaseParameters.TraceLevel.Value),
		set(defs.UserTaskManagedInitialWarehouseSize.FieldName(), databaseParameters.UserTaskManagedInitialWarehouseSize.Value),
		set(defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), databaseParameters.UserTaskMinimumTriggerIntervalInSeconds.Value),
		set(defs.UserTaskTimeoutMs.FieldName(), databaseParameters.UserTaskTimeoutMs.Value),
	)
}

var databaseParametersCustomDiff = ParametersCustomDiffFromTypedParameters(
	databaseParametersProvider,
	databaseParameterDiffFunctions,
)

func databaseParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) (*sdk.DatabaseParametersDetails, error) {
	id, err := sdk.ParseAccountObjectIdentifier(d.Id())
	if err != nil {
		return nil, err
	}
	return meta.(*provider.Context).Client.Databases.ShowParametersDetails(ctx, id)
}

func databaseParameterDiffFunctions(parameters *sdk.DatabaseParametersDetails) []schema.CustomizeDiffFunc {
	return []schema.CustomizeDiffFunc{
		IdentifierTypedParameterValueComputedIf(defs.Catalog.FieldName(), parameters.Catalog, sdk.ParameterTypeDatabase),
		IntTypedParameterValueComputedIf(defs.DataRetentionTimeInDays.FieldName(), parameters.DataRetentionTimeInDays, sdk.ParameterTypeDatabase),
		StringTypedParameterValueComputedIf(defs.DefaultDdlCollation.FieldName(), parameters.DefaultDdlCollation, sdk.ParameterTypeDatabase),
		StringTypedParameterValueComputedIf(defs.DefaultNotebookComputePoolCpu.FieldName(), parameters.DefaultNotebookComputePoolCpu, sdk.ParameterTypeDatabase),
		StringTypedParameterValueComputedIf(defs.DefaultNotebookComputePoolGpu.FieldName(), parameters.DefaultNotebookComputePoolGpu, sdk.ParameterTypeDatabase),
		BoolTypedParameterValueComputedIf(defs.EnableConsoleOutput.FieldName(), parameters.EnableConsoleOutput, sdk.ParameterTypeDatabase),
		IdentifierTypedParameterValueComputedIf(defs.ExternalVolume.FieldName(), parameters.ExternalVolume, sdk.ParameterTypeDatabase),
		StringTypedParameterValueComputedIf(defs.LogEventLevel.FieldName(), parameters.LogEventLevel, sdk.ParameterTypeDatabase),
		StringTypedParameterValueComputedIf(defs.LogLevel.FieldName(), parameters.LogLevel, sdk.ParameterTypeDatabase),
		IntTypedParameterValueComputedIf(defs.MaxDataExtensionTimeInDays.FieldName(), parameters.MaxDataExtensionTimeInDays, sdk.ParameterTypeDatabase),
		BoolTypedParameterValueComputedIf(defs.QuotedIdentifiersIgnoreCase.FieldName(), parameters.QuotedIdentifiersIgnoreCase, sdk.ParameterTypeDatabase),
		BoolTypedParameterValueComputedIf(defs.ReplaceInvalidCharacters.FieldName(), parameters.ReplaceInvalidCharacters, sdk.ParameterTypeDatabase),
		StringTypedParameterValueComputedIf(defs.StorageSerializationPolicy.FieldName(), parameters.StorageSerializationPolicy, sdk.ParameterTypeDatabase),
		IntTypedParameterValueComputedIf(defs.SuspendTaskAfterNumFailures.FieldName(), parameters.SuspendTaskAfterNumFailures, sdk.ParameterTypeDatabase),
		IntTypedParameterValueComputedIf(defs.TaskAutoRetryAttempts.FieldName(), parameters.TaskAutoRetryAttempts, sdk.ParameterTypeDatabase),
		StringTypedParameterValueComputedIf(defs.TraceLevel.FieldName(), parameters.TraceLevel, sdk.ParameterTypeDatabase),
		StringTypedParameterValueComputedIf(defs.UserTaskManagedInitialWarehouseSize.FieldName(), parameters.UserTaskManagedInitialWarehouseSize, sdk.ParameterTypeDatabase),
		IntTypedParameterValueComputedIf(defs.UserTaskMinimumTriggerIntervalInSeconds.FieldName(), parameters.UserTaskMinimumTriggerIntervalInSeconds, sdk.ParameterTypeDatabase),
		IntTypedParameterValueComputedIf(defs.UserTaskTimeoutMs.FieldName(), parameters.UserTaskTimeoutMs, sdk.ParameterTypeDatabase),
	}
}
