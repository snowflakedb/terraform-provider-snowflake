package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

func (c *PipeClient) CopyStatement(t *testing.T, table *sdk.Table, stage *sdk.Stage) string {
	t.Helper()
	require.NotNil(t, table, "table has to be created")
	require.NotNil(t, stage, "stage has to be created")
	return fmt.Sprintf("COPY INTO %s\nFROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())
}

func (c *PipeClient) CreatePipe(t *testing.T, copyStatement string) (*sdk.Pipe, func()) {
	t.Helper()
	return c.CreatePipeInSchema(t, c.ids.SchemaId(), copyStatement)
}

func (c *PipeClient) CreatePipeInSchema(t *testing.T, schemaId sdk.DatabaseObjectIdentifier, copyStatement string) (*sdk.Pipe, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	err := c.client().Create(ctx, sdk.NewCreatePipeRequest(id, copyStatement))
	require.NoError(t, err)

	pipe, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return pipe, c.DropPipeFunc(t, id)
}

func (c *PipeClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Pipe, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByID(ctx, id)
}

func (c *PipeClient) DropPipeFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropPipeRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *PipeClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.PipeDetails, error) {
	t.Helper()
	return c.client().Describe(context.Background(), id)
}
