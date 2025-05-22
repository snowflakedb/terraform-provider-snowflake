package helpers

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk"
	"github.com/stretchr/testify/require"
)

type ShareClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewShareClient(context *TestClientContext, idsGenerator *IdsGenerator) *ShareClient {
	return &ShareClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ShareClient) client() sdk.Shares {
	return c.context.client.Shares
}

func (c *ShareClient) CreateShare(t *testing.T) (*sdk.Share, func()) {
	t.Helper()
	return c.CreateShareWithIdentifier(t, c.ids.RandomAccountObjectIdentifier())
}

func (c *ShareClient) CreateShareWithIdentifier(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Share, func()) {
	t.Helper()
	return c.CreateShareWithOptions(t, id, &sdk.CreateShareOptions{})
}

func (c *ShareClient) CreateShareWithOptions(t *testing.T, id sdk.AccountObjectIdentifier, opts *sdk.CreateShareOptions) (*sdk.Share, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, id, opts)
	require.NoError(t, err)

	share, err := c.client().ShowByID(ctx, id)
	require.NoError(t, err)

	return share, c.DropShareFunc(t, id)
}

func (c *ShareClient) DropShareFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		err := c.client().Drop(ctx, id, &sdk.DropShareOptions{IfExists: sdk.Bool(true)})
		require.NoError(t, err)
	}
}

func (c *ShareClient) SetAccountOnShare(t *testing.T, accountId sdk.AccountIdentifier, shareId sdk.AccountObjectIdentifier) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, shareId, &sdk.AlterShareOptions{
		Set: &sdk.ShareSet{
			Accounts: []sdk.AccountIdentifier{accountId},
		},
	})
	require.NoError(t, err)
}
