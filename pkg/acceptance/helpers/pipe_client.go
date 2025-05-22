package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type PipeClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewPipeClient(context *TestClientContext, idsGenerator *IdsGenerator) *PipeClient {
	return &PipeClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *PipeClient) client() sdk.Pipes {
	return c.context.client.Pipes
}

func (c *PipeClient) CreatePipe(t *testing.T, copyStatement string) (*sdk.Pipe, func()) {
	t.Helper()
	return c.CreatePipeInSchema(t, c.ids.SchemaId(), copyStatement)
}

func (c *PipeClient) CreatePipeInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, copyStatement string) (*sdk.Pipe, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	err := c.client().Create(ctx, id, copyStatement, &sdk.CreatePipeOptions{})
	require.NoError(t, err)

	pipe, err := c.client().Describe(ctx, id)
	require.NoError(t, err)

	return pipe, c.DropPipeFunc(t, id)
}

func (c *PipeClient) DropPipeFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropPipeOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}
