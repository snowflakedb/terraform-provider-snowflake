package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type NotebookClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewNotebookClient(context *TestClientContext, idsGenerator *IdsGenerator) *NotebookClient {
	return &NotebookClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *NotebookClient) client() sdk.Notebooks {
	return c.context.client.Notebooks
}

func (c *NotebookClient) Create(t *testing.T, name sdk.SchemaObjectIdentifier) (*sdk.Notebook, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateNotebookRequest(name))
}

func (c *NotebookClient) CreateWithRequest(t *testing.T, req *sdk.CreateNotebookRequest) (*sdk.Notebook, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	notebook, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return notebook, c.DropFunc(t, req.GetName())
}

func (c *NotebookClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropNotebookRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *NotebookClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Notebook, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
