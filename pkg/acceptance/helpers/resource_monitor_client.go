package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ResourceMonitorClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewResourceMonitorClient(context *TestClientContext, idsGenerator *IdsGenerator) *ResourceMonitorClient {
	return &ResourceMonitorClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ResourceMonitorClient) client() sdk.ResourceMonitors {
	return c.context.client.ResourceMonitors
}

func (c *ResourceMonitorClient) CreateResourceMonitor(t *testing.T) (*sdk.ResourceMonitor, func()) {
	t.Helper()
	id := c.ids.RandomAccountObjectIdentifier()
	req := sdk.NewCreateResourceMonitorRequest(id).
		WithWith(sdk.ResourceMonitorWithRequest{
			CreditQuota: sdk.Pointer(100),
			Triggers: []sdk.TriggerDefinitionRequest{
				{
					Threshold:     100,
					TriggerAction: sdk.TriggerActionSuspend,
				},
				{
					Threshold:     70,
					TriggerAction: sdk.TriggerActionSuspendImmediate,
				},
				{
					Threshold:     90,
					TriggerAction: sdk.TriggerActionNotify,
				},
			},
		})
	return c.CreateResourceMonitorWithRequest(t, req)
}

func (c *ResourceMonitorClient) CreateResourceMonitorWithRequest(t *testing.T, req *sdk.CreateResourceMonitorRequest) (*sdk.ResourceMonitor, func()) {
	t.Helper()
	ctx := context.Background()

	id := req.ID()
	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	resourceMonitor, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return resourceMonitor, c.DropResourceMonitorFunc(t, id)
}

func (c *ResourceMonitorClient) Alter(t *testing.T, req *sdk.AlterResourceMonitorRequest) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *ResourceMonitorClient) DropResourceMonitorFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropResourceMonitorRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *ResourceMonitorClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ResourceMonitor, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}
