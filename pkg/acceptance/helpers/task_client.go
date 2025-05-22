package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type TaskClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewTaskClient(context *TestClientContext, idsGenerator *IdsGenerator) *TaskClient {
	return &TaskClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *TaskClient) client() sdk.Tasks {
	return c.context.client.Tasks
}

func (c *TaskClient) defaultCreateTaskRequest(t *testing.T) *sdk.CreateTaskRequest {
	t.Helper()
	id := c.ids.RandomSchemaObjectIdentifier()
	warehouseReq := sdk.NewCreateTaskWarehouseRequest().WithWarehouse(c.ids.WarehouseId())
	return sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP").WithWarehouse(*warehouseReq)
}

func (c *TaskClient) Create(t *testing.T) (*sdk.Task, func()) {
	t.Helper()
	return c.CreateWithRequest(t, c.defaultCreateTaskRequest(t))
}

func (c *TaskClient) CreateWithSchedule(t *testing.T) (*sdk.Task, func()) {
	t.Helper()
	return c.CreateWithRequest(t, c.defaultCreateTaskRequest(t).WithSchedule("60 MINUTES"))
}

func (c *TaskClient) CreateWithAfter(t *testing.T, after ...sdk.SchemaObjectIdentifier) (*sdk.Task, func()) {
	t.Helper()
	return c.CreateWithRequest(t, c.defaultCreateTaskRequest(t).WithAfter(after))
}

func (c *TaskClient) CreateWithRequest(t *testing.T, request *sdk.CreateTaskRequest) (*sdk.Task, func()) {
	t.Helper()
	ctx := context.Background()

	id := request.GetName()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	task, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return task, c.DropFunc(t, id)
}

func (c *TaskClient) Alter(t *testing.T, req *sdk.AlterTaskRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *TaskClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropTaskRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *TaskClient) Show(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.Task, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
