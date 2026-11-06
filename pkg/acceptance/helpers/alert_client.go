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
		sdk.NewAlertConditionFromString("SELECT 1"),
		"SELECT 1",
	).
		WithWarehouse(c.ids.WarehouseId()).
		WithSchedule("USING CRON * * * * * UTC"))
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

func (c *AlertClient) Alter(t *testing.T, req *sdk.AlterAlertRequest) {
	t.Helper()
	err := c.client().Alter(context.Background(), req)
	require.NoError(t, err)
}

func (c *AlertClient) DropAlertFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropAlertRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *AlertClient) Describe(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.AlertDetails, error) {
	t.Helper()
	return c.client().Describe(context.Background(), id)
}

func (c *AlertClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Alert, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}
