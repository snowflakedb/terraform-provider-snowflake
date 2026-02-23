package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type HybridTableClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewHybridTableClient(context *TestClientContext, idsGenerator *IdsGenerator) *HybridTableClient {
	return &HybridTableClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *HybridTableClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	id := c.ids.RandomSchemaObjectIdentifier()
	err := c.context.client.HybridTables.Create(ctx, sdk.NewCreateHybridTableRequest(
		id,
		sdk.HybridTableColumnsConstraintsAndIndexes{
			Columns: []sdk.HybridTableColumn{
				{
					Name: "id",
					Type: sdk.DataType("INT"),
					InlineConstraint: &sdk.HybridTableColumnInlineConstraint{
						Type: sdk.ColumnConstraintTypePrimaryKey,
					},
				},
			},
		},
	))
	require.NoError(t, err)

	return id, c.DropFunc(t, id)
}

func (c *HybridTableClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()

	return func() {
		ctx := context.Background()
		err := c.context.client.HybridTables.Drop(ctx, sdk.NewDropHybridTableRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *HybridTableClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.HybridTable, error) {
	t.Helper()
	ctx := context.Background()
	return c.context.client.HybridTables.ShowByID(ctx, id)
}

func (c *HybridTableClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) ([]sdk.HybridTableDetails, error) {
	t.Helper()
	ctx := context.Background()
	return c.context.client.HybridTables.Describe(ctx, id)
}

func (c *HybridTableClient) DescribeDetails(t *testing.T, id sdk.SchemaObjectIdentifier) (*[]sdk.HybridTableDetails, error) {
	t.Helper()
	details, err := c.Describe(t, id)
	if err != nil {
		return nil, err
	}
	return &details, nil
}
