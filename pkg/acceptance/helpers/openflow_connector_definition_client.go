package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// OpenflowConnectorDefinitionClient wraps the read-only SHOW OPENFLOW CONNECTOR DEFINITIONS. Definitions
// are global, immutable and managed by Snowflake, so there is nothing to create or clean up.
type OpenflowConnectorDefinitionClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewOpenflowConnectorDefinitionClient(context *TestClientContext, idsGenerator *IdsGenerator) *OpenflowConnectorDefinitionClient {
	return &OpenflowConnectorDefinitionClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *OpenflowConnectorDefinitionClient) client() sdk.OpenflowConnectorDefinitions {
	return c.context.client.OpenflowConnectorDefinitions
}

func (c *OpenflowConnectorDefinitionClient) Show(t *testing.T) []sdk.OpenflowConnectorDefinition {
	t.Helper()
	ctx := context.Background()

	definitions, err := c.client().Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest())
	require.NoError(t, err)
	return definitions
}

// ShowWithLike returns the definitions matching the given LIKE pattern.
func (c *OpenflowConnectorDefinitionClient) ShowWithLike(t *testing.T, pattern string) []sdk.OpenflowConnectorDefinition {
	t.Helper()
	ctx := context.Background()

	definitions, err := c.client().Show(ctx, sdk.NewShowOpenflowConnectorDefinitionRequest().WithLike(sdk.Like{Pattern: sdk.String(pattern)}))
	require.NoError(t, err)
	return definitions
}

// PostgresCdcDefinitionName is the definition tests instantiate connectors from. COMMIT succeeds without
// any configuration, but reaching RUNNING needs a reachable source database, so START settles on
// START_FAILED instead.
const PostgresCdcDefinitionName = "OPENFLOW_POSTGRES_CDC"

// ForTesting returns the definition tests instantiate connectors from.
func (c *OpenflowConnectorDefinitionClient) ForTesting(t *testing.T) sdk.OpenflowConnectorDefinition {
	t.Helper()

	definitions := c.ShowWithLike(t, PostgresCdcDefinitionName)
	require.NotEmptyf(t, definitions, "the %s connector definition is not available on this account", PostgresCdcDefinitionName)
	return definitions[0]
}
