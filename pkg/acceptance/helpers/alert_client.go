package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type AlertClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewAlertClient(context *TestClientContext, idsGenerator *IdsGenerator) *AlertClient {
	return &AlertClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *AlertClient) client() sdk.Alerts {
	return c.context.client.Alerts
}

func (c *AlertClient) CreateAlert(t *testing.T) (*sdk.Alert, func()) {
	t.Helper()
	return c.CreateAlertWithRequest(t, sdk.NewCreateAlertRequest(
		c.ids.RandomSchemaObjectIdentifier(),
		c.ids.WarehouseId(),
		"USING CRON * * * * * UTC",
		sdk.NewAlertConditionFromString("SELECT 1"),
		"SELECT 1",
	))
}

func (c *AlertClient) CreateAlertWithRequest(t *testing.T, req *sdk.CreateAlertRequest) (*sdk.Alert, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	id := req.GetName()
	alert, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return alert, c.DropAlertFunc(t, id)
}

func (c *AlertClient) DropAlertFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropAlertRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
