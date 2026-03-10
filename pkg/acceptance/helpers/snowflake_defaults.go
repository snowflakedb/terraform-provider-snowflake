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
