package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	return c.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(c.ids.RandomAccountObjectIdentifier()))
}

// CreateTestWarehouseIfNotExists should be used to create the main warehouse used throughout the acceptance tests.
// It's created only if it does not exist already.
func (c *WarehouseClient) CreateTestWarehouseIfNotExists(t *testing.T) (*sdk.Warehouse, func()) {
	t.Helper()
	return c.CreateWarehouseWithRequest(t, sdk.NewCreateWarehouseRequest(c.ids.WarehouseId()).WithIfNotExists(true))
}

func (c *WarehouseClient) CreateWarehouseWithRequest(t *testing.T, request *sdk.CreateWarehouseRequest) (*sdk.Warehouse, func()) {
	t.Helper()
	ctx := context.Background()

	id := request.ID()
	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	warehouse, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return warehouse, c.DropWarehouseFunc(t, id)
}

func (c *WarehouseClient) DropWarehouseFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropWarehouseRequest(id).WithIfExists(true))
		require.NoError(t, err)
		err = c.context.client.Sessions.UseWarehouse(ctx, c.ids.WarehouseId())
		require.NoError(t, err)
	}
}

func (c *WarehouseClient) UpdateMaxConcurrencyLevel(t *testing.T, id sdk.AccountObjectIdentifier, level int) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{MaxConcurrencyLevel: sdk.Int(level)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateWarehouseSize(t *testing.T, id sdk.AccountObjectIdentifier, newSize sdk.WarehouseSize) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{WarehouseSize: sdk.Pointer(newSize)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateWarehouseType(t *testing.T, id sdk.AccountObjectIdentifier, newType sdk.WarehouseType) {
	t.Helper()

	ctx := context.Background()

	err := c.client().AlterWithSuspend(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{WarehouseType: sdk.Pointer(newType)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateResourceConstraint(t *testing.T, id sdk.AccountObjectIdentifier, newResourceConstraint sdk.WarehouseResourceConstraint) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{ResourceConstraint: sdk.Pointer(newResourceConstraint)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateGeneration(t *testing.T, id sdk.AccountObjectIdentifier, newGeneration sdk.WarehouseGeneration) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{Generation: sdk.Pointer(newGeneration)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateWarehouseTypeAndResourceConstraint(t *testing.T, id sdk.AccountObjectIdentifier, newType sdk.WarehouseType, newResourceConstraint sdk.WarehouseResourceConstraint) {
	t.Helper()

	ctx := context.Background()

	err := c.client().AlterWithSuspend(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{WarehouseType: sdk.Pointer(newType), ResourceConstraint: sdk.Pointer(newResourceConstraint)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateWarehouseTypeAndGeneration(t *testing.T, id sdk.AccountObjectIdentifier, newType sdk.WarehouseType, newGeneration sdk.WarehouseGeneration) {
	t.Helper()

	ctx := context.Background()

	err := c.client().AlterWithSuspend(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{WarehouseType: sdk.Pointer(newType), Generation: sdk.Pointer(newGeneration)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateStatementTimeoutInSeconds(t *testing.T, id sdk.AccountObjectIdentifier, newValue int) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{StatementTimeoutInSeconds: sdk.Int(newValue)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UnsetStatementTimeoutInSeconds(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithUnset(sdk.WarehouseUnsetRequest{StatementTimeoutInSeconds: sdk.Bool(true)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateAutoResume(t *testing.T, id sdk.AccountObjectIdentifier, newAutoResume bool) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{AutoResume: sdk.Pointer(newAutoResume)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) UpdateAutoSuspend(t *testing.T, id sdk.AccountObjectIdentifier, newAutoSuspend int) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSet(sdk.WarehouseSetRequest{AutoSuspend: sdk.Int(newAutoSuspend)}))
	require.NoError(t, err)
}

func (c *WarehouseClient) Suspend(t *testing.T, id sdk.AccountObjectIdentifier) {
	t.Helper()

	ctx := context.Background()

	err := c.client().Alter(ctx, sdk.NewAlterWarehouseRequest(id).WithSuspend(true))
	require.NoError(t, err)
}

func (c *WarehouseClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Warehouse, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}

func (c *WarehouseClient) CreateAdaptive(t *testing.T) (*sdk.Warehouse, func()) {
	t.Helper()
	return c.CreateAdaptiveWithRequest(t, sdk.NewCreateAdaptiveWarehouseRequest(c.ids.RandomAccountObjectIdentifier()))
}

func (c *WarehouseClient) CreateAdaptiveWithRequest(t *testing.T, request *sdk.CreateAdaptiveWarehouseRequest) (*sdk.Warehouse, func()) {
	t.Helper()
	ctx := context.Background()

	id := request.ID()
	err := c.client().CreateAdaptive(ctx, request)
	require.NoError(t, err)

	warehouse, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return warehouse, c.DropWarehouseFunc(t, id)
}
