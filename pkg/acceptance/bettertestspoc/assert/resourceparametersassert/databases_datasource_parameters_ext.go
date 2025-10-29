package resourceparametersassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// DatabasesDatasourceParameters is a temporary workaround to have better parameter assertions in data source acceptance tests.
func DatabasesDatasourceParameters(t *testing.T, datasourceReference string) *DatabaseResourceParametersAssert {
	t.Helper()

	d := DatabaseResourceParametersAssert{
		ResourceAssert: assert.NewDatasourceAssert(datasourceReference, "parameters", "databases.0."),
	}
	d.AddAssertion(assert.ValueSet("parameters.#", "1"))
	return &d
}

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
