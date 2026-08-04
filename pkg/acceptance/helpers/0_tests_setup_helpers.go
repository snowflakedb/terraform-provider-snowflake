package helpers

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *TestClient) CreateTestDatabase(ctx context.Context, ifNotExists bool) (*sdk.Database, func(), error) {
	id := c.Ids.DatabaseId()
	cleanup := func() {
		_ = c.context.client.Databases.DropSafely(ctx, id)
	}
	req := c.Database.TestParametersSet(id).WithIfNotExists(ifNotExists)
	err := c.context.client.Databases.Create(ctx, req)
	if err != nil {
		return nil, cleanup, err
	}
	database, err := c.context.client.Databases.ShowByID(ctx, id)
	return database, cleanup, err
}

func (c *TestClient) CreateTestSchema(ctx context.Context, ifNotExists bool) (*sdk.Schema, func(), error) {
	id := c.Ids.SchemaId()
	cleanup := func() {
		_ = c.context.client.Schemas.DropSafely(ctx, id)
	}
	err := c.context.client.Schemas.Create(ctx, sdk.NewCreateSchemaRequest(id).WithIfNotExists(ifNotExists))
	if err != nil {
		return nil, cleanup, err
	}
	schema, err := c.context.client.Schemas.ShowByID(ctx, id)
	return schema, cleanup, err
}

func (c *TestClient) CreateTestWarehouse(ctx context.Context, ifNotExists bool) (*sdk.Warehouse, func(), error) {
	id := c.Ids.WarehouseId()
	cleanup := func() {
		_ = c.context.client.Warehouses.DropSafely(ctx, id)
	}
	err := c.context.client.Warehouses.Create(ctx, sdk.NewCreateWarehouseRequest(id).WithIfNotExists(ifNotExists))
	if err != nil {
		return nil, cleanup, err
	}
	warehouse, err := c.context.client.Warehouses.ShowByID(ctx, id)
	return warehouse, cleanup, err
}

// CreateTestGcpPubSubNotificationIntegration creates a single notification integration shared by all the tests
// needing one (e.g. as a task's error integration). It's created once per test run on purpose: the account-wide
// limit of GCP Pubsub queue integrations is 10, so creating one per test exhausts it quickly.
func (c *TestClient) CreateTestGcpPubSubNotificationIntegration(ctx context.Context) (*sdk.NotificationIntegration, func(), error) {
	id := c.Ids.RandomAccountObjectIdentifier()
	cleanup := func() {
		_ = c.context.client.NotificationIntegrations.DropSafely(ctx, id)
	}
	req := sdk.NewCreateNotificationIntegrationRequest(id, true).
		WithAutomatedDataLoadsParams(
			*sdk.NewAutomatedDataLoadsParamsRequest().
				WithGoogleAutoParams(*sdk.NewGoogleAutoParamsRequest(gcpPubsubSubscriptionName)),
		)
	err := c.context.client.NotificationIntegrations.Create(ctx, req)
	if err != nil {
		return nil, cleanup, err
	}
	integration, err := c.context.client.NotificationIntegrations.ShowByID(ctx, id)
	return integration, cleanup, err
}
