package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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
	return c.CreateResourceMonitorWithOptions(t, &sdk.CreateResourceMonitorOptions{
		With: &sdk.ResourceMonitorWith{
			CreditQuota: sdk.Pointer(100),
			Triggers: []sdk.TriggerDefinition{
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
		},
	})
}

func (c *ResourceMonitorClient) CreateResourceMonitorWithOptions(t *testing.T, opts *sdk.CreateResourceMonitorOptions) (*sdk.ResourceMonitor, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()

	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)

	resourceMonitor, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return resourceMonitor, c.DropResourceMonitorFunc(t, id)
}

func (c *ResourceMonitorClient) Alter(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.AlterResourceMonitorOptions) {
	t.Helper()
	ctx := context.Background()
	err := c.client().Alter(ctx, id, opts)
	require.NoError(t, err)
}

func (c *ResourceMonitorClient) DropResourceMonitorFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropResourceMonitorOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}

func (c *ResourceMonitorClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ResourceMonitor, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}
