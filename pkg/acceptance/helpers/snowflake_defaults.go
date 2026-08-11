package helpers

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

// DefaultQueryAccelerationMaxScaleFactor returns the Snowflake default for
// QUERY_ACCELERATION_MAX_SCALE_FACTOR. Prod accounts default to 2; preprod still defaults to 8.
func (c *SnowflakeDefaultsClient) DefaultQueryAccelerationMaxScaleFactor(t *testing.T) int {
	t.Helper()
	if c.context.snowflakeEnvironment != testenvs.SnowflakeProdEnvironment {
		return 8
	}
	return 2
}

// DefaultStatementTimeoutInSecondsLevel returns the expected parameter level for
// STATEMENT_TIMEOUT_IN_SECONDS on a newly created warehouse. Prod inherits the
// Snowflake default (empty level); non-prod accounts currently surface ACCOUNT because they need to be overridden.
func (c *SnowflakeDefaultsClient) DefaultStatementTimeoutInSecondsLevel(t *testing.T) sdk.ParameterType {
	t.Helper()
	if c.context.snowflakeEnvironment != testenvs.SnowflakeProdEnvironment {
		return sdk.ParameterTypeAccount
	}
	return sdk.ParameterTypeSnowflakeDefault
}
