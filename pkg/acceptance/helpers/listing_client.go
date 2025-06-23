package helpers

import (
	"context"
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

func (c *ListingClient) SampleListingManifest(t *testing.T) string {
	t.Helper()
	return `
title: "MyListing"
subtitle: "Subtitle for MyListing"
description: "Description for MyListing"
listing_terms:
 type: "OFFLINE"
`
}

func (c *ListingClient) Create(t *testing.T) (*sdk.Listing, func()) {
	t.Helper()
	ctx := context.Background()

	id := c.ids.RandomAccountObjectIdentifier()
	err := c.client().Create(ctx, sdk.NewCreateListingRequest(id, c.SampleListingManifest(t)))
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
		err := c.client().Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	}
}
