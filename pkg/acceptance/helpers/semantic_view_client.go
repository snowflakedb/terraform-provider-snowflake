package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type SemanticViewClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewSemanticViewClient(context *TestClientContext, idsGenerator *IdsGenerator) *SemanticViewClient {
	return &SemanticViewClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *SemanticViewClient) client() sdk.SemanticViews {
	return c.context.client.SemanticViews
}

func (c *SemanticViewClient) Create(t *testing.T) (*sdk.SemanticView, func()) {
	t.Helper()
	return c.CreateWithId(t, c.ids.RandomSchemaObjectIdentifier())
}

func (c *SemanticViewClient) CreateWithId(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.SemanticView, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateSemanticViewRequest(id, []sdk.LogicalTableRequest{}))
}

func (c *SemanticViewClient) CreateWithRequest(t *testing.T, req *sdk.CreateSemanticViewRequest) (*sdk.SemanticView, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)
	semanticView, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)
	return semanticView, c.DropFunc(t, req.GetName())
}

func (c *SemanticViewClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropSemanticViewRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *SemanticViewClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.SemanticView, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().ShowByID(ctx, id)
}

func (c *SemanticViewClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) ([]sdk.SemanticViewDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.client().Describe(ctx, id)
}
