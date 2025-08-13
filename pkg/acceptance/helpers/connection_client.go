package helpers

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ConnectionClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewConnectionClient(context *TestClientContext, idsGenerator *IdsGenerator) *ConnectionClient {
	return &ConnectionClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ConnectionClient) client() sdk.Connections {
	return c.context.client.Connections
}

func (c *ConnectionClient) Create(t *testing.T) (*sdk.Connection, func()) {
	t.Helper()
	return c.CreateWithIdentifier(t, c.ids.RandomAccountObjectIdentifier())
}

func (c *ConnectionClient) CreateWithIdentifier(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Connection, func()) {
	t.Helper()
	ctx := context.Background()
	request := sdk.NewCreateConnectionRequest(id)
	err := c.client().Create(ctx, request)
	require.NoError(t, err)
	connection, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return connection, c.DropFunc(t, id)
}

func (c *ConnectionClient) CreateReplication(t *testing.T, id sdk.AccountObjectIdentifier, replicaOf sdk.ExternalObjectIdentifier) (*sdk.Connection, func()) {
	t.Helper()
	ctx := context.Background()
	request := sdk.NewCreateConnectionRequest(id).WithAsReplicaOf(replicaOf)
	err := c.client().Create(ctx, request)
	require.NoError(t, err)
	connection, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)
	return connection, c.DropFunc(t, id)
}

func (c *ConnectionClient) Alter(t *testing.T, req *sdk.AlterConnectionRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *ConnectionClient) Drop(t *testing.T, id sdk.AccountObjectIdentifier) error {
	t.Helper()
	ctx := context.Background()

	return c.client().Drop(ctx, sdk.NewDropConnectionRequest(id).WithIfExists(true))
}

func (c *ConnectionClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropConnectionRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *ConnectionClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Connection, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *ConnectionClient) GetConnectionUrl(organizationName, objectName string) string {
	return strings.ToLower(fmt.Sprintf("%s-%s.snowflakecomputing.com", organizationName, objectName))
}
