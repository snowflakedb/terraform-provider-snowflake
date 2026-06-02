package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type IcebergTableClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewIcebergTableClient(context *TestClientContext, idsGenerator *IdsGenerator) *IcebergTableClient {
	return &IcebergTableClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *IcebergTableClient) Create(t *testing.T) (*sdk.IcebergTable, func()) {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	return c.CreateWithRequest(t, sdk.NewCreateIcebergTableRequest(id, sdk.IcebergTableColumnsAndConstraintsRequest{Columns: []sdk.IcebergTableColumnRequest{{Name: "ID", ColumnType: testdatatypes.DataTypeNumber}}}))
}

func (c *IcebergTableClient) CreateWithRequest(t *testing.T, request *sdk.CreateIcebergTableRequest) (*sdk.IcebergTable, func()) {
	t.Helper()
	ctx := context.Background()
	err := c.context.client.IcebergTables.Create(ctx, request)
	require.NoError(t, err)
	id := request.GetName()
	obj, err := c.context.client.IcebergTables.ShowByID(ctx, id)
	require.NoError(t, err)
	return obj, c.DropFunc(t, id)
}

func (c *IcebergTableClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()

	return func() {
		ctx := context.Background()
		err := c.context.client.IcebergTables.Drop(ctx, sdk.NewDropIcebergTableRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *IcebergTableClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.IcebergTable, error) {
	t.Helper()
	ctx := context.Background()
	return c.context.client.IcebergTables.ShowByID(ctx, id)
}

func (c *IcebergTableClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) ([]sdk.IcebergTableDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.context.client.IcebergTables.Describe(ctx, id)
}
