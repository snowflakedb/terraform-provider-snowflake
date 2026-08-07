package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO(SNOW-1325215): change raw sqls to proper client
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

// Create follows the other integration helpers: it creates through the typed client and returns the object.
func (c *ExternalAccessIntegrationClient) Create(t *testing.T, networkRuleId sdk.SchemaObjectIdentifier) (*sdk.ExternalAccessIntegration, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	err := c.client().ExternalAccessIntegrations.Create(ctx, sdk.NewCreateExternalAccessIntegrationRequest(id, []sdk.SchemaObjectIdentifier{networkRuleId}, true))
	require.NoError(t, err)

	externalAccessIntegration, err := c.client().ExternalAccessIntegrations.ShowByID(ctx, id)
	require.NoError(t, err)

	return externalAccessIntegration, c.DropExternalAccessIntegrationFunc(t, id)
}

func (c *ExternalAccessIntegrationClient) CreateExternalAccessIntegration(t *testing.T, networkRuleId sdk.SchemaObjectIdentifier) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ENABLED = TRUE`, id.FullyQualifiedName(), networkRuleId.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropExternalAccessIntegrationFunc(t, id)
}

func (c *ExternalAccessIntegrationClient) CreateExternalAccessIntegrationWithNetworkRuleAndSecret(t *testing.T, networkRuleId sdk.SchemaObjectIdentifier, secretId sdk.SchemaObjectIdentifier) (sdk.AccountObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`CREATE EXTERNAL ACCESS INTEGRATION %s ALLOWED_NETWORK_RULES = (%s) ALLOWED_AUTHENTICATION_SECRETS = (%s) ENABLED = TRUE`, id.FullyQualifiedName(), networkRuleId.FullyQualifiedName(), secretId.FullyQualifiedName()))
	require.NoError(t, err)
	return id, c.DropExternalAccessIntegrationFunc(t, id)
}

func (c *ExternalAccessIntegrationClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ExternalAccessIntegration, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ExternalAccessIntegrations.ShowByID(ctx, id)
}

func (c *ExternalAccessIntegrationClient) DescribeDetails(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ExternalAccessIntegrationDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ExternalAccessIntegrations.DescribeDetails(ctx, id)
}

func (c *ExternalAccessIntegrationClient) DropExternalAccessIntegrationFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP EXTERNAL ACCESS INTEGRATION IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
