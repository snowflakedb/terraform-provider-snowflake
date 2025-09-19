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
	opts := c.Database.TestParametersSet()
	opts.IfNotExists = sdk.Bool(ifNotExists)
	err := c.context.client.Databases.Create(ctx, id, opts)
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
	err := c.context.client.Schemas.Create(ctx, id, &sdk.CreateSchemaOptions{IfNotExists: sdk.Bool(ifNotExists)})
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
	err := c.context.client.Warehouses.Create(ctx, id, &sdk.CreateWarehouseOptions{IfNotExists: sdk.Bool(ifNotExists)})
	if err != nil {
		return nil, cleanup, err
	}
	warehouse, err := c.context.client.Warehouses.ShowByID(ctx, id)
	return warehouse, cleanup, err
}
