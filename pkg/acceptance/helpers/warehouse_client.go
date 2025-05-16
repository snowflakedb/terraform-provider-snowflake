package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type WarehouseClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewWarehouseClient(context *TestClientContext, idsGenerator *IdsGenerator) *WarehouseClient {
	return &WarehouseClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *WarehouseClient) client() sdk.Warehouses {
	return c.context.client.Warehouses
}

func (c *WarehouseClient) UseWarehouse(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	err := c.context.client.Sessions.UseWarehouse(ctx, id)
	require.NoError(t, err)
	return func() {
		err = c.context.client.Sessions.UseWarehouse(ctx, c.ids.WarehouseId())
		require.NoError(t, err)
	}
}

func (c *WarehouseClient) CreateWarehouse(t *testing.T) (*sdk.Warehouse, func()) {
	t.Helper()
	return c.CreateWarehouseWithOptions(t, c.ids.RandomAccountObjectIdentifier(), &sdk.CreateWarehouseOptions{})
}

// CreateTestWarehouseIfNotExists should be used to create the main warehouse used throughout the acceptance tests.
// It's created only if it does not exist already.
func (c *WarehouseClient) CreateTestWarehouseIfNotExists(t *testing.T) (*sdk.Warehouse, func()) {
	t.Helper()
	return c.CreateWarehouseWithOptions(t, c.ids.WarehouseId(), &sdk.CreateWarehouseOptions{IfNotExists: sdk.Bool(true)})
}

func (c *WarehouseClient) CreateWarehouseWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateWarehouseOptions) (*sdk.Warehouse, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)

	warehouse, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return warehouse, c.DropWarehouseFunc(t, id)
}

func (c *WarehouseClient) DropWarehouseFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropWarehouseOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
		err = c.context.client.Sessions.UseWarehouse(ctx, c.ids.WarehouseId())
		require.NoError(t, err)
	}
}

func (c *WarehouseClient) UpdateMaxConcurrencyLevel(t *testing.T, id sdk.AccountObjectIdentifier, level int) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Set: &sdk.WarehouseSet{MaxConcurrencyLevel: sdk.Int(level)}})
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateWarehouseSize(t *testing.T, id sdk.AccountObjectIdentifier, newSize sdk.WarehouseSize) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Set: &sdk.WarehouseSet{WarehouseSize: sdk.Pointer(newSize)}})
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateWarehouseType(t *testing.T, id sdk.AccountObjectIdentifier, newType sdk.WarehouseType) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Set: &sdk.WarehouseSet{WarehouseType: sdk.Pointer(newType)}})
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateStatementTimeoutInSeconds(t *testing.T, id sdk.AccountObjectIdentifier, newValue int) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Set: &sdk.WarehouseSet{StatementTimeoutInSeconds: sdk.Int(newValue)}})
	require.NoError(t, err)
}

func (c *WarehouseClient) UnsetStatementTimeoutInSeconds(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Unset: &sdk.WarehouseUnset{StatementTimeoutInSeconds: sdk.Bool(true)}})
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateAutoResume(t *testing.T, id sdk.AccountObjectIdentifier, newAutoResume bool) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Set: &sdk.WarehouseSet{AutoResume: sdk.Pointer(newAutoResume)}})
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateAutoSuspend(t *testing.T, id sdk.AccountObjectIdentifier, newAutoSuspend int) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Set: &sdk.WarehouseSet{AutoSuspend: sdk.Int(newAutoSuspend)}})
	require.NoError(t, err)
}

func (c *WarehouseClient) Suspend(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, id, &sdk.AlterWarehouseOptions{Suspend: sdk.Bool(true)})
	require.NoError(t, err)
}

func (c *WarehouseClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Warehouse, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
