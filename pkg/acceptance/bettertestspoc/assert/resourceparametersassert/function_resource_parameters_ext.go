package resourceparametersassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
)

func (f *FunctionResourceParametersAssert) HasAllDefaults() *FunctionResourceParametersAssert {
	return f.
		HasEnableConsoleOutput(false).
		HasLogLevel(sdk.LogLevelOff).
		HasMetricLevel(sdk.MetricLevelNone).
		HasTraceLevel(sdk.TraceLevelOff)
}
