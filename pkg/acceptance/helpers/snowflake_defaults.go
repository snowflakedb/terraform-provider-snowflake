package helpers

import (
	"fmt"
	"slices"
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

func (c *SnowflakeDefaultsClient) EnabledForSnowflakeOauthSecurityIntegration(t *testing.T) bool {
	t.Helper()
	if slices.Contains([]testenvs.SnowflakeEnvironment{testenvs.SnowflakeNonProdEnvironment, testenvs.SnowflakePreProdGovEnvironment}, c.context.snowflakeEnvironment) {
		return true
	}
	return false
}

func (c *SnowflakeDefaultsClient) StageIdentifierOutputFormatForStreamOnDirectoryTable(t *testing.T, id sdk.SchemaObjectIdentifier) string {
	t.Helper()
	if slices.Contains([]testenvs.SnowflakeEnvironment{testenvs.SnowflakeNonProdEnvironment, testenvs.SnowflakePreProdGovEnvironment}, c.context.snowflakeEnvironment) {
		return fmt.Sprintf(`"%s"."%s".%s`, id.DatabaseName(), id.SchemaName(), id.Name())
	}
	return id.Name()
}

func (c *SnowflakeDefaultsClient) WarehouseGenerationEmptyByDefault(t *testing.T) bool {
	t.Helper()
	if c.context.snowflakeEnvironment == testenvs.SnowflakePreProdGovEnvironment {
		return true
	}
	return false
}
