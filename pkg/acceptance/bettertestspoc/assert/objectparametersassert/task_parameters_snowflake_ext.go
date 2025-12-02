package objectparametersassert

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (t *TaskParametersAssert) HasServerlessTaskMinStatementSizeEnum(expected sdk.WarehouseSize) *TaskParametersAssert {
	t.AddAssertion(assert.SnowflakeParameterValueSet(sdk.TaskParameterServerlessTaskMinStatementSize, string(expected)))
	return t
}

func (t *TaskParametersAssert) HasServerlessTaskMaxStatementSizeEnum(expected sdk.WarehouseSize) *TaskParametersAssert {
	t.AddAssertion(assert.SnowflakeParameterValueSet(sdk.TaskParameterServerlessTaskMaxStatementSize, string(expected)))
	return t
}
