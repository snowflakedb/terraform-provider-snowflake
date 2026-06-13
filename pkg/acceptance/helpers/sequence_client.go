package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type SequenceClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSequenceClient(context *TestClientContext, idsGenerator *IdsGenerator) *SequenceClient {
	return &SequenceClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *SequenceClient) client() sdk.Sequences {
	return c.context.client.Sequences
}

func (c *SequenceClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropSequenceRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *SequenceClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	return c.CreateWithId(t, c.ids.RandomSchemaObjectIdentifier())
}

func (c *SequenceClient) CreateWithId(t *testing.T, id sdk.SchemaObjectIdentifier) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Create(ctx, sdk.NewCreateSequenceRequest(id))
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}
