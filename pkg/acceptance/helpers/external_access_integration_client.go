package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ExternalAccessIntegrationClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewExternalAccessIntegrationClient(context *TestClientContext, idsGenerator *IdsGenerator) *ExternalAccessIntegrationClient {
	return &ExternalAccessIntegrationClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ExternalAccessIntegrationClient) client() *sdk.Client {
	return c.context.client
}

func (c *ExternalAccessIntegrationClient) CreateExternalAccessIntegration(t *testing.T, networkRuleId sdk.SchemaObjectIdentifier) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	req := sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRuleId}, true)
	err := c.client().ExternalAccessIntegrations.Create(ctx, req)
	require.NoError(t, err)
	return id, c.DropExternalAccessIntegrationFunc(t, id)
}

func (c *ExternalAccessIntegrationClient) CreateExternalAccessIntegrationWithNetworkRuleAndSecret(t *testing.T, networkRuleId sdk.SchemaObjectIdentifier, secretId sdk.SchemaObjectIdentifier) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	req := sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRuleId}, true).
		WithAllowedAuthenticationSecrets([]sdk.SchemaObjectIdentifier{secretId})
	err := c.client().ExternalAccessIntegrations.Create(ctx, req)
	require.NoError(t, err)
	return id, c.DropExternalAccessIntegrationFunc(t, id)
}

func (c *ExternalAccessIntegrationClient) DropExternalAccessIntegrationFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().ExternalAccessIntegrations.DropSafely(ctx, id)
		require.NoError(t, err)
	}
}
