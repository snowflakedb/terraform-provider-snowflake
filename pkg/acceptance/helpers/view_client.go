package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ViewClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewViewClient(context *TestClientContext, idsGenerator *IdsGenerator) *ViewClient {
	return &ViewClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ViewClient) client() sdk.Views {
	return c.context.client.Views
}

func (c *ViewClient) CreateView(t *testing.T, query string) (*sdk.View, func()) {
	t.Helper()
	return c.CreateViewInSchema(t, query, c.ids.SchemaId())
}

func (c *ViewClient) CreateViewInSchema(t *testing.T, query string, schemaId sdk.DatabaseObjectIdentifier) (*sdk.View, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifierInSchema(schemaId)

	err := c.client().Create(ctx, sdk.NewCreateViewRequest(id, query))
	require.NoError(t, err)

	view, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return view, c.DropViewFunc(t, id)
}

func (c *ViewClient) RecreateView(t *testing.T, id sdk.SchemaObjectIdentifier, query string) *sdk.View {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, sdk.NewCreateViewRequest(id, query).WithOrReplace(true))
	require.NoError(t, err)

	view, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return view
}

func (c *ViewClient) Alter(t *testing.T, req *sdk.AlterViewRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *ViewClient) DropViewFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropViewRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *ViewClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.View, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
