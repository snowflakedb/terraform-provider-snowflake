package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type NetworkRuleClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewNetworkRuleClient(context *TestClientContext, idsGenerator *IdsGenerator) *NetworkRuleClient {
	return &NetworkRuleClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *NetworkRuleClient) client() sdk.NetworkRules {
	return c.context.client.NetworkRules
}

func (c *NetworkRuleClient) Create(t *testing.T) (*sdk.NetworkRule, func()) {
	t.Helper()
	return c.CreateEgressWithIdentifier(t, c.ids.RandomSchemaObjectIdentifier())
}

func (c *NetworkRuleClient) CreateIngress(t *testing.T) (*sdk.NetworkRule, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateNetworkRuleRequest(
		c.ids.RandomSchemaObjectIdentifier(),
		sdk.NetworkRuleTypeIpv4,
		[]sdk.NetworkRuleValue{},
		sdk.NetworkRuleModeIngress,
	))
}

func (c *NetworkRuleClient) CreateIP(t *testing.T) (*sdk.NetworkRule, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateNetworkRuleRequest(
		c.ids.RandomSchemaObjectIdentifier(),
		sdk.NetworkRuleTypeIpv4,
		[]sdk.NetworkRuleValue{{Value: "1.2.3.4"}},
		sdk.NetworkRuleModeIngress,
	))
}

func (c *NetworkRuleClient) CreateEgressWithIdentifier(t *testing.T, id sdk.SchemaObjectIdentifier) (*sdk.NetworkRule, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateNetworkRuleRequest(
		id,
		sdk.NetworkRuleTypeHostPort,
		[]sdk.NetworkRuleValue{},
		sdk.NetworkRuleModeEgress,
	))
}

func (c *NetworkRuleClient) CreateWithRequest(t *testing.T, request *sdk.CreateNetworkRuleRequest) (*sdk.NetworkRule, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, request)
	require.NoError(t, err)

	networkRule, err := c.client().ShowByID(ctx, request.GetName())
	require.NoError(t, err)

	return networkRule, c.DropFunc(t, request.GetName())
}

func (c *NetworkRuleClient) DropFunc(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropNetworkRuleRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
