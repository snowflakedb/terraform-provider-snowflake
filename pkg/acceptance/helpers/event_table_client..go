package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type EventTableClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewEventTableClient(context *TestClientContext, idsGenerator *IdsGenerator) *EventTableClient {
	return &EventTableClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *EventTableClient) client() sdk.EventTables {
	return c.context.client.EventTables
}

func (c *EventTableClient) Create(t *testing.T) (*sdk.EventTable, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()
	err := c.client().Create(ctx, sdk.NewCreateEventTableRequest(id))
	require.NoError(t, err)

	integration, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return integration, c.DropFunc(t, id)
}

func (c *EventTableClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropEventTableRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
