package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ListingClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewListingClient(context *TestClientContext, idsGenerator *IdsGenerator) *ListingClient {
	return &ListingClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *ListingClient) client() sdk.Listings {
	return c.context.client.Listings
}

func (c *ListingClient) Create(t *testing.T) (*sdk.Listing, func()) {
	t.Helper()
	return c.CreateWithId(t, c.ids.RandomAccountObjectIdentifier())
}

func (c *ListingClient) CreateWithId(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Listing, func()) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Create(ctx, sdk.NewCreateListingRequest(id).
		WithAs(c.BasicManifest(t)).
		WithReview(false).
		WithPublish(false),
	)
	assert.NoError(t, err)

	listing, err := c.client().ShowByID(ctx, id)
	assert.NoError(t, err)

	return listing, c.DropFunc(t, id)
}

func (c *ListingClient) Alter(t *testing.T, req *sdk.AlterListingRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	require.NoError(t, err)
}

func (c *ListingClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		if err := c.client().DropSafely(ctx, id); !errors.Is(err, sdk.ErrObjectNotFound) {
			assert.NoError(t, err)
		}
	}
}

func (c *ListingClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Listing, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *ListingClient) BasicManifest(t *testing.T) string {
	t.Helper()
	return `title: "title"
subtitle: "subtitle"
description: "description"
listing_terms:
  type: "OFFLINE"
`
}
