package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type MaterializedViewClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewMaterializedViewClient(context *TestClientContext, idsGenerator *IdsGenerator) *MaterializedViewClient {
	return &MaterializedViewClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *MaterializedViewClient) client() sdk.MaterializedViews {
	return c.context.client.MaterializedViews
}

func (c *MaterializedViewClient) CreateMaterializedView(t *testing.T, query string, orReplace bool) (*sdk.MaterializedView, func()) {
	t.Helper()
	return c.CreateMaterializedViewWithName(t, c.ids.RandomSchemaObjectIdentifier(), query, orReplace)
}

func (c *MaterializedViewClient) CreateMaterializedViewWithName(t *testing.T, id sdk.SchemaObjectIdentifier, query string, orReplace bool) (*sdk.MaterializedView, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, sdk.NewCreateMaterializedViewRequest(id, query).WithOrReplace(orReplace))
	require.NoError(t, err)

	view, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return view, c.DropMaterializedViewFunc(t, id)
}

func (c *MaterializedViewClient) DropMaterializedViewFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropMaterializedViewRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
