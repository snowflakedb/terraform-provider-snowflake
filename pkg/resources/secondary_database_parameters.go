package resources

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/defs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func handleSecondaryDatabaseParametersCreate(d *schema.ResourceData, createOpts *sdk.CreateSecondaryDatabaseRequest) diag.Diagnostics {
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
