package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
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

func (c *FailoverGroupClient) CreateFailoverGroup(t *testing.T) (*sdk.FailoverGroup, func()) {
	t.Helper()
	objectTypes := []sdk.PluralObjectType{sdk.PluralObjectTypeRoles}
	accountID := c.ids.AccountIdentifierWithLocator()
	allowedAccounts := []sdk.AccountIdentifier{accountID}
	return c.CreateFailoverGroupWithOptions(t, objectTypes, allowedAccounts, nil)
}

func (c *FailoverGroupClient) CreateFailoverGroupWithOptions(t *testing.T, objectTypes []sdk.PluralObjectType, allowedAccounts []sdk.AccountIdentifier, opts *sdk.CreateFailoverGroupOptions) (*sdk.FailoverGroup, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()

	err := c.client().Create(ctx, id, objectTypes, allowedAccounts, opts)
	require.NoError(t, err)

	failoverGroup, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return failoverGroup, c.DropFailoverGroupFunc(t, id)
}

func (c *FailoverGroupClient) DropFailoverGroupFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropFailoverGroupOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}
