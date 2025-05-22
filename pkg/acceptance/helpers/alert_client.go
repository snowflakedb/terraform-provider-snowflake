package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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
	schedule := "USING CRON * * * * * UTC"
	condition := "SELECT 1"
	action := "SELECT 1"
	return c.CreateAlertWithOptions(t, schedule, condition, action, &sdk.CreateAlertOptions{})
}

func (c *AlertClient) CreateAlertWithOptions(t *testing.T, schedule string, condition string, action string, opts *sdk.CreateAlertOptions) (*sdk.Alert, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomSchemaObjectIdentifier()

	err := c.client().Create(ctx, id, c.ids.WarehouseId(), schedule, condition, action, opts)
	require.NoError(t, err)

	alert, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return alert, c.DropAlertFunc(t, id)
}

func (c *AlertClient) DropAlertFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropAlertOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}
