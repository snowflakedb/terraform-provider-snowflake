package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type BudgetClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewBudgetClient(context *TestClientContext, idsGenerator *IdsGenerator) *BudgetClient {
	return &BudgetClient{context: context, ids: idsGenerator}
}

func (c *BudgetClient) client() sdk.Budgets {
	return c.context.client.Budgets
}

func (c *BudgetClient) Create(t *testing.T) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	id := c.ids.RandomSchemaObjectIdentifier()
	err := c.client().Create(ctx, sdk.NewCreateBudgetRequest(id))
	require.NoError(t, err)
	return id, c.DropFunc(t, id)
}

func (c *BudgetClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	return func() {
		err := c.client().Drop(ctx, sdk.NewDropBudgetRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
