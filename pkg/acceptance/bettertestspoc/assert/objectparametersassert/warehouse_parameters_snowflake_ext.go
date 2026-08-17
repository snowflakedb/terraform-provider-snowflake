package objectparametersassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasAllDefaultsForEnvironment is an environment-aware variant of HasAllDefaults.
// It replaces HasAllDefaults() in tests that run on preprod, where STATEMENT_TIMEOUT_IN_SECONDS
// is set at ACCOUNT level rather than the Snowflake default level.
func (w *WarehouseParametersAssert) HasAllDefaultsForEnvironment(t *testing.T, defaults *helpers.SnowflakeDefaultsClient) *WarehouseParametersAssert {
	t.Helper()
	return w.
		HasDefaultParameterValueOnLevel(sdk.WarehouseParameterMaxConcurrencyLevel, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds, sdk.ParameterTypeSnowflakeDefault).
		HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementTimeoutInSeconds, defaults.DefaultStatementTimeoutInSecondsLevel(t))
}
