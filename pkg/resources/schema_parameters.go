package resources

import (
	"context"
	"strconv"
	"strings"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	schemaParametersSchema     = make(map[string]*schema.Schema)
	schemaParametersCustomDiff = ParametersCustomDiff(
		schemaParametersProvider,
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterDataRetentionTimeInDays, valueTypeInt, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterMaxDataExtensionTimeInDays, valueTypeInt, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterExternalVolume, valueTypeString, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterCatalog, valueTypeString, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterReplaceInvalidCharacters, valueTypeBool, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterDefaultDDLCollation, valueTypeString, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterStorageSerializationPolicy, valueTypeString, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterLogLevel, valueTypeString, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterTraceLevel, valueTypeString, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterSuspendTaskAfterNumFailures, valueTypeInt, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterTaskAutoRetryAttempts, valueTypeInt, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, valueTypeString, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterUserTaskTimeoutMs, valueTypeInt, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, valueTypeInt, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterQuotedIdentifiersIgnoreCase, valueTypeBool, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterEnableConsoleOutput, valueTypeBool, sdk.ParameterTypeSchema},
		parameter[sdk.ObjectParameter]{sdk.ObjectParameterPipeExecutionPaused, valueTypeBool, sdk.ParameterTypeSchema},
	)
)

func init() {
	additionalSchemaParameterFields := []parameterDef[sdk.ObjectParameter]{
		{Name: sdk.ObjectParameterPipeExecutionPaused, Type: schema.TypeBool, Description: "Specifies whether to pause a running pipe, primarily in preparation for transferring ownership of the pipe to a different role."},
	}

	additionalSchemaParameters := make(map[string]*schema.Schema)
	for _, field := range additionalSchemaParameterFields {
		fieldName := strings.ToLower(string(field.Name))

		additionalSchemaParameters[fieldName] = &schema.Schema{
			Type:        field.Type,
			Description: enrichWithReferenceToParameterDocs(field.Name, field.Description),
			Computed:    true,
			Optional:    true,
		}
	}
	schemaParametersSchema = collections.MergeMaps(databaseParametersSchema, additionalSchemaParameters)
}

func schemaParametersProvider(ctx context.Context, d ResourceIdProvider, meta any) ([]*sdk.Parameter, error) {
	return parametersProvider(ctx, d, meta.(*provider.Context), schemaParametersProviderFunc, sdk.ParseDatabaseObjectIdentifier)
}

func schemaParametersProviderFunc(c *sdk.Client) showParametersFunc[sdk.DatabaseObjectIdentifier] {
	return c.Schemas.ShowParameters
}

func handleSchemaParameterRead(d *schema.ResourceData, databaseParameters []*sdk.Parameter) diag.Diagnostics {
	for _, parameter := range databaseParameters {
		switch parameter.Key {
		case
			string(sdk.ObjectParameterDataRetentionTimeInDays),
			string(sdk.ObjectParameterMaxDataExtensionTimeInDays),
			string(sdk.ObjectParameterSuspendTaskAfterNumFailures),
			string(sdk.ObjectParameterTaskAutoRetryAttempts),
			string(sdk.ObjectParameterUserTaskTimeoutMs),
			string(sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds):
			value, err := strconv.Atoi(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(strings.ToLower(parameter.Key), value); err != nil {
				return diag.FromErr(err)
			}
		case
			string(sdk.ObjectParameterExternalVolume),
			string(sdk.ObjectParameterCatalog),
			string(sdk.ObjectParameterDefaultDDLCollation),
			string(sdk.ObjectParameterStorageSerializationPolicy),
			string(sdk.ObjectParameterLogLevel),
			string(sdk.ObjectParameterTraceLevel),
			string(sdk.ObjectParameterUserTaskManagedInitialWarehouseSize):
			if err := d.Set(strings.ToLower(parameter.Key), parameter.Value); err != nil {
				return diag.FromErr(err)
			}
		case
			string(sdk.ObjectParameterPipeExecutionPaused),
			string(sdk.ObjectParameterReplaceInvalidCharacters),
			string(sdk.ObjectParameterQuotedIdentifiersIgnoreCase),
			string(sdk.ObjectParameterEnableConsoleOutput):
			value, err := strconv.ParseBool(parameter.Value)
			if err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set(strings.ToLower(parameter.Key), value); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func handleSchemaParametersCreate(d *schema.ResourceData, createOpts *sdk.CreateSchemaOptions) diag.Diagnostics {
	return JoinDiags(
		handleParameterCreate(d, sdk.ObjectParameterDataRetentionTimeInDays, &createOpts.DataRetentionTimeInDays),
		handleParameterCreate(d, sdk.ObjectParameterMaxDataExtensionTimeInDays, &createOpts.MaxDataExtensionTimeInDays),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterExternalVolume, &createOpts.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterCatalog, &createOpts.Catalog, stringToAccountObjectIdentifier),
		handleParameterCreate(d, sdk.ObjectParameterPipeExecutionPaused, &createOpts.PipeExecutionPaused),
		handleParameterCreate(d, sdk.ObjectParameterReplaceInvalidCharacters, &createOpts.ReplaceInvalidCharacters),
		handleParameterCreate(d, sdk.ObjectParameterDefaultDDLCollation, &createOpts.DefaultDDLCollation),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterStorageSerializationPolicy, &createOpts.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterLogLevel, &createOpts.LogLevel, sdk.ToLogLevel),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterTraceLevel, &createOpts.TraceLevel, sdk.ToTraceLevel),
		handleParameterCreate(d, sdk.ObjectParameterSuspendTaskAfterNumFailures, &createOpts.SuspendTaskAfterNumFailures),
		handleParameterCreate(d, sdk.ObjectParameterTaskAutoRetryAttempts, &createOpts.TaskAutoRetryAttempts),
		handleParameterCreateWithMapping(d, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, &createOpts.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskTimeoutMs, &createOpts.UserTaskTimeoutMs),
		handleParameterCreate(d, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, &createOpts.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterCreate(d, sdk.ObjectParameterQuotedIdentifiersIgnoreCase, &createOpts.QuotedIdentifiersIgnoreCase),
		handleParameterCreate(d, sdk.ObjectParameterEnableConsoleOutput, &createOpts.EnableConsoleOutput),
	)
}

func handleSchemaParametersChanges(d *schema.ResourceData, set *sdk.SchemaSet, unset *sdk.SchemaUnset) diag.Diagnostics {
	return JoinDiags(
		handleParameterUpdate(d, sdk.ObjectParameterDataRetentionTimeInDays, &set.DataRetentionTimeInDays, &unset.DataRetentionTimeInDays),
		handleParameterUpdate(d, sdk.ObjectParameterMaxDataExtensionTimeInDays, &set.MaxDataExtensionTimeInDays, &unset.MaxDataExtensionTimeInDays),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterExternalVolume, &set.ExternalVolume, &unset.ExternalVolume, stringToAccountObjectIdentifier),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterCatalog, &set.Catalog, &unset.Catalog, stringToAccountObjectIdentifier),
		handleParameterUpdate(d, sdk.ObjectParameterPipeExecutionPaused, &set.PipeExecutionPaused, &unset.PipeExecutionPaused),
		handleParameterUpdate(d, sdk.ObjectParameterReplaceInvalidCharacters, &set.ReplaceInvalidCharacters, &unset.ReplaceInvalidCharacters),
		handleParameterUpdate(d, sdk.ObjectParameterDefaultDDLCollation, &set.DefaultDDLCollation, &unset.DefaultDDLCollation),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterStorageSerializationPolicy, &set.StorageSerializationPolicy, &unset.StorageSerializationPolicy, sdk.ToStorageSerializationPolicy),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterLogLevel, &set.LogLevel, &unset.LogLevel, sdk.ToLogLevel),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterTraceLevel, &set.TraceLevel, &unset.TraceLevel, sdk.ToTraceLevel),
		handleParameterUpdate(d, sdk.ObjectParameterSuspendTaskAfterNumFailures, &set.SuspendTaskAfterNumFailures, &unset.SuspendTaskAfterNumFailures),
		handleParameterUpdate(d, sdk.ObjectParameterTaskAutoRetryAttempts, &set.TaskAutoRetryAttempts, &unset.TaskAutoRetryAttempts),
		handleParameterUpdateWithMapping(d, sdk.ObjectParameterUserTaskManagedInitialWarehouseSize, &set.UserTaskManagedInitialWarehouseSize, &unset.UserTaskManagedInitialWarehouseSize, sdk.ToWarehouseSize),
		handleParameterUpdate(d, sdk.ObjectParameterUserTaskTimeoutMs, &set.UserTaskTimeoutMs, &unset.UserTaskTimeoutMs),
		handleParameterUpdate(d, sdk.ObjectParameterUserTaskMinimumTriggerIntervalInSeconds, &set.UserTaskMinimumTriggerIntervalInSeconds, &unset.UserTaskMinimumTriggerIntervalInSeconds),
		handleParameterUpdate(d, sdk.ObjectParameterQuotedIdentifiersIgnoreCase, &set.QuotedIdentifiersIgnoreCase, &unset.QuotedIdentifiersIgnoreCase),
		handleParameterUpdate(d, sdk.ObjectParameterEnableConsoleOutput, &set.EnableConsoleOutput, &unset.EnableConsoleOutput),
	)
}
