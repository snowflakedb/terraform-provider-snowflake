package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type CortexSearchServiceClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewCortexSearchServiceClient(context *TestClientContext, idsGenerator *IdsGenerator) *CortexSearchServiceClient {
	return &CortexSearchServiceClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *CortexSearchServiceClient) client() sdk.CortexSearchServices {
	return c.context.client.CortexSearchServices
}

func (c *CortexSearchServiceClient) CreateCortexSearchService(t *testing.T, tableId sdk.SchemaObjectIdentifier) (*sdk.CortexSearchService, func()) {
	t.Helper()
	return c.CreateCortexSearchServiceWithOptions(t, c.ids.RandomSchemaObjectIdentifier(), c.ids.WarehouseId(), tableId, "some_text_column")
}

func (c *CortexSearchServiceClient) CreateCortexSearchServiceWithOptions(t *testing.T, id sdk.SchemaObjectIdentifier, warehouseId sdk.AccountObjectIdentifier, tableId sdk.SchemaObjectIdentifier, column string) (*sdk.CortexSearchService, func()) {
	t.Helper()
	targetLag := "2 minutes"
	query := fmt.Sprintf(`select %s from %s`, column, tableId.FullyQualifiedName())
	comment := random.Comment()
	ctx := context.Background()

	err := c.client().Create(ctx, sdk.NewCreateCortexSearchServiceRequest(id, column, warehouseId, targetLag, query).WithComment(comment))
	require.NoError(t, err)

	contextSearchService, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return contextSearchService, c.DropCortexSearchServiceFunc(t, id)
}

func (c *CortexSearchServiceClient) DropCortexSearchServiceFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropCortexSearchServiceRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
