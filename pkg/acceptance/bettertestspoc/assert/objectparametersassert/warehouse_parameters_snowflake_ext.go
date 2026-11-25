package objectparametersassert

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

// HasAllDefaultsForEnvironment is an environment-aware variant of HasAllDefaults.
// Non-prod accounts report STATEMENT_TIMEOUT_IN_SECONDS at ACCOUNT level. The value
// is still the catalog default. SYSTEM is accepted in place of the Snowflake default level.
func (w *WarehouseParametersAssert) HasAllDefaultsForEnvironment(t *testing.T, defaults *helpers.SnowflakeDefaultsClient) *WarehouseParametersAssert {
	t.Helper()
	return w.
		HasDefaultParameterValueOnDefaultLevel(sdk.WarehouseParameterMaxConcurrencyLevel).
		HasDefaultParameterValueOnDefaultLevel(sdk.WarehouseParameterStatementQueuedTimeoutInSeconds).
		HasDefaultParameterValueOnLevel(sdk.WarehouseParameterStatementTimeoutInSeconds, defaults.DefaultStatementTimeoutInSecondsLevel(t))
}
