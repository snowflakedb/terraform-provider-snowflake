package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type NetworkPolicyClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewNetworkPolicyClient(context *TestClientContext, idsGenerator *IdsGenerator) *NetworkPolicyClient {
	return &NetworkPolicyClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *NetworkPolicyClient) client() sdk.NetworkPolicies {
	return c.context.client.NetworkPolicies
}

func (c *NetworkPolicyClient) CreateNetworkPolicy(t *testing.T) (*sdk.NetworkPolicy, func()) {
	t.Helper()
	return c.CreateNetworkPolicyWithRequest(t, sdk.NewCreateNetworkPolicyRequest(c.ids.RandomAccountObjectIdentifier()))
}

// CreateNetworkPolicyNotEmpty tackles SF Error: 098519 (22023): Empty network policy [id] cannot be active.
func (c *NetworkPolicyClient) CreateNetworkPolicyNotEmpty(t *testing.T) (*sdk.NetworkPolicy, func()) {
	t.Helper()
	return c.CreateNetworkPolicyWithRequest(t, sdk.NewCreateNetworkPolicyRequest(c.ids.RandomAccountObjectIdentifier()).
		WithBlockedIpList([]sdk.IPRequest{*sdk.NewIPRequest("1.1.1.1")}),
	)
}

func (c *NetworkPolicyClient) CreateNetworkPolicyWithRequest(t *testing.T, request *sdk.CreateNetworkPolicyRequest) (*sdk.NetworkPolicy, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	networkPolicy, err := c.client().ShowByID(ctx, request.GetName())
	require.NoError(t, err)

	return networkPolicy, c.DropNetworkPolicyFunc(t, request.GetName())
}

func (c *NetworkPolicyClient) Update(t *testing.T, request *sdk.AlterNetworkPolicyRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, request)
	require.NoError(t, err)
}

func (c *NetworkPolicyClient) DropNetworkPolicyFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropNetworkPolicyRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}

func (c *NetworkPolicyClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.NetworkPolicy, error) {
	t.Helper()
	ctx := context.Background()

	return c.client().ShowByID(ctx, id)
}
