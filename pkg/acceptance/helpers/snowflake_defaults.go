package helpers

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
)

type SnowflakeDefaultsClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSnowflakeDefaultsClient(context *TestClientContext) *SnowflakeDefaultsClient {
	return &SnowflakeDefaultsClient{
		context: context,
	}
}

func (c *SnowflakeDefaultsClient) WarehouseGenerationEmptyByDefault(t *testing.T) bool {
	t.Helper()
	if c.context.snowflakeEnvironment == testenvs.SnowflakePreProdGovEnvironment {
		return true
	}
	return false
}

func (c *SnowflakeDefaultsClient) WarehouseEnableQueryAcceleration(t *testing.T) bool {
	t.Helper()
	if c.context.snowflakeEnvironment == testenvs.SnowflakeNonProdEnvironment || c.context.snowflakeEnvironment == testenvs.SnowflakePreProdGovEnvironment {
		return true
	}
	return false
}

func (c *SnowflakeDefaultsClient) WarehouseQueryAccelerationMaxScaleFactor(t *testing.T) int {
	t.Helper()
	if c.context.snowflakeEnvironment == testenvs.SnowflakeNonProdEnvironment || c.context.snowflakeEnvironment == testenvs.SnowflakePreProdGovEnvironment {
		return 2
	}
	return 8
}
