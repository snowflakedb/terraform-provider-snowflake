package resourceparametersassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (d *DatabaseResourceParametersAssert) HasAllDefaultParameters() *DatabaseResourceParametersAssert {
	return d.
		HasDataRetentionTimeInDays(1).
		HasMaxDataExtensionTimeInDays(14).
		HasExternalVolume("").
		HasCatalog("").
		HasReplaceInvalidCharacters(false).
		HasDefaultDdlCollation("").
		HasStorageSerializationPolicy(sdk.StorageSerializationPolicyOptimized).
		HasLogLevel(sdk.LogLevelOff).
		HasTraceLevel(sdk.TraceLevelOff).
		HasSuspendTaskAfterNumFailures(10).
		HasTaskAutoRetryAttempts(0).
		HasUserTaskManagedInitialWarehouseSize("Medium").
		HasUserTaskTimeoutMs(3600000).
		HasUserTaskMinimumTriggerIntervalInSeconds(30).
		HasQuotedIdentifiersIgnoreCase(false).
		HasEnableConsoleOutput(false)
}
