package resourceparametersassert

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

func (s *SchemaResourceParametersAssert) HasAllDefaultParameters() *SchemaResourceParametersAssert {
	return s.
		HasCatalog("").
		HasDataRetentionTimeInDays(1).
		HasDefaultDdlCollation("").
		HasDefaultNotebookComputePoolCpu("SYSTEM_COMPUTE_POOL_CPU").
		HasDefaultNotebookComputePoolGpu("SYSTEM_COMPUTE_POOL_GPU").
		HasEnableConsoleOutput(false).
		HasExternalVolume("").
		HasLogEventLevel(sdk.LogLevelOff).
		HasLogLevel(sdk.LogLevelOff).
		HasMaxDataExtensionTimeInDays(14).
		HasPipeExecutionPaused(false).
		HasQuotedIdentifiersIgnoreCase(false).
		HasReplaceInvalidCharacters(false).
		HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
		HasSuspendTaskAfterNumFailures(10).
		HasTaskAutoRetryAttempts(0).
		HasTraceLevel(sdk.TraceLevelOff).
		HasUserTaskManagedInitialWarehouseSize("Medium").
		HasUserTaskMinimumTriggerIntervalInSeconds(30).
		HasUserTaskTimeoutMs(3600000)
}
