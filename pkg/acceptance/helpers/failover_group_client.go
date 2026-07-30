package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type FailoverGroupClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewFailoverGroupClient(context *TestClientContext, idsGenerator *IdsGenerator) *FailoverGroupClient {
	return &FailoverGroupClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *FailoverGroupClient) client() sdk.FailoverGroups {
	return c.context.client.FailoverGroups
}

func (c *FailoverGroupClient) Create(t *testing.T) (*sdk.FailoverGroup, func()) {
	t.Helper()
	return c.CreateWithRequest(t, sdk.NewCreateFailoverGroupRequest(
		c.ids.RandomAccountObjectIdentifier(),
		[]sdk.PluralObjectType{sdk.PluralObjectTypeRoles},
		[]sdk.AccountIdentifier{c.ids.AccountIdentifierWithLocator()},
	))
}

func (c *FailoverGroupClient) CreateWithRequest(t *testing.T, req *sdk.CreateFailoverGroupRequest) (*sdk.FailoverGroup, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, req)
	require.NoError(t, err)

	failoverGroup, err := c.client().ShowByID(ctx, req.GetName())
	require.NoError(t, err)

	return failoverGroup, c.DropFunc(t, req.GetName())
}

func (c *FailoverGroupClient) RemoveAllowedAccounts(t *testing.T, id sdk.AccountObjectIdentifier, accounts ...sdk.AccountIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().AlterSource(ctx, sdk.NewAlterSourceFailoverGroupRequest(id).WithRemove(*sdk.NewFailoverGroupRemoveRequest().WithAllowedAccounts(accounts)))
	require.NoError(t, err)
}

func (c *FailoverGroupClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, sdk.NewDropFailoverGroupRequest(id).WithIfExists(true))
		require.NoError(t, err)
	}
}
