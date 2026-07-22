package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

// TODO [fill]: change raw sqls to proper client
type GatewayClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewGatewayClient(context *TestClientContext, idsGenerator *IdsGenerator) *GatewayClient {
	return &GatewayClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *GatewayClient) client() *sdk.Client {
	return c.context.client
}

func (c *GatewayClient) SampleSpec(serviceId sdk.SchemaObjectIdentifier, endpointName string) string {
	return fmt.Sprintf(`$$
spec:
  type: traffic_split
  split_type: custom
  targets:
  - type: endpoint
    value: %s.%s.%s!%s
    weight: 100
$$`, serviceId.DatabaseName(), serviceId.SchemaName(), serviceId.Name(), endpointName)
}

func (c *GatewayClient) Create(t *testing.T, serviceId sdk.SchemaObjectIdentifier, endpointName string) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	query := fmt.Sprintf(`CREATE GATEWAY %s FROM SPECIFICATION %s`, id.FullyQualifiedName(), c.SampleSpec(serviceId, endpointName))
	_, err := c.client().ExecForTests(ctx, query)
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *GatewayClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		_, err := c.client().ExecForTests(ctx, fmt.Sprintf(`DROP GATEWAY IF EXISTS %s`, id.FullyQualifiedName()))
		require.NoError(t, err)
	}
}
