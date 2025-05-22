package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type StreamClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewStreamClient(context *TestClientContext, idsGenerator *IdsGenerator) *StreamClient {
	return &StreamClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *StreamClient) client() sdk.Streams {
	return c.context.client.Streams
}

func (c *StreamClient) CreateOnTable(t *testing.T, tableId sdk.SchemaObjectIdentifier) (*sdk.Stream, func()) {
	t.Helper()

	return c.CreateOnTableWithRequest(t, sdk.NewCreateOnTableStreamRequest(c.ids.RandomSchemaObjectIdentifier(), tableId))
}

func (c *StreamClient) CreateOnTableWithRequest(t *testing.T, req *sdk.CreateOnTableStreamRequest) (*sdk.Stream, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateOnTable(ctx, req)
	require.NoError(t, err)

	stream, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return stream, c.DropFunc(t, req.GetName())
}

func (c *StreamClient) CreateOnViewWithRequest(t *testing.T, req *sdk.CreateOnViewStreamRequest) (*sdk.Stream, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().CreateOnView(ctx, req)
	require.NoError(t, err)

	stream, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return stream, c.DropFunc(t, req.GetName())
}

func (c *StreamClient) Update(t *testing.T, request *sdk.AlterStreamRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}

func (c *StreamClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropStreamRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *StreamClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Stream, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *StreamClient) Alter(t *testing.T, req *sdk.AlterStreamRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}
