package helpers

import (
	"context"
	"fmt"
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

	manifest, _ := c.BasicManifest(t)
	err := c.client().Create(ctx, sdk.NewCreateListingRequest(id).
		WithAs(manifest).
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
		assert.NoError(t, c.client().DropSafely(ctx, id))
	}
}

func (c *ListingClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Listing, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *ListingClient) ShowVersions(t *testing.T, id sdk.AccountObjectIdentifier) ([]sdk.ListingVersion, error) {
	t.Helper()
	return c.client().ShowVersions(context.Background(), sdk.NewShowVersionsListingRequest(id))
}

func (c *ListingClient) BasicManifest(t *testing.T) (string, string) {
	t.Helper()
	return c.basicManifest(t, "basic_", "subtitle")
}

func (c *ListingClient) BasicManifestWithDifferentSubtitle(t *testing.T) (string, string) {
	t.Helper()
	return c.basicManifest(t, "basic_with_diff_subtitle_", "different_subtitle")
}

func (c *ListingClient) BasicManifestWithUnquotedValues(t *testing.T) (string, string) {
	t.Helper()
	return c.basicManifestWithUnquotedValues(t, "basic_", "subtitle")
}

func (c *ListingClient) BasicManifestWithUnquotedValuesAndDifferentSubtitle(t *testing.T) (string, string) {
	t.Helper()
	return c.basicManifestWithUnquotedValues(t, "basic_with_diff_subtitle_", "different_subtitle")
}

func (c *ListingClient) BasicManifestWithTargetAccount(t *testing.T, targetAccount sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicManifestWithTargetAccount(t, "with_target_accounts_", "subtitle", targetAccount)
}

func (c *ListingClient) BasicManifestWithTargetAccountAndDifferentSubtitle(t *testing.T, targetAccount sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicManifestWithTargetAccount(t, "with_target_accounts_and_different_subtitle_", "different_subtitle", targetAccount)
}

func (c *ListingClient) BasicManifestWithUnquotedValuesAndTargetAccount(t *testing.T, targetAccount sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicManifestWithUnquotedValuesAndTargetAccount(t, "with_target_accounts_", "subtitle", targetAccount)
}

func (c *ListingClient) BasicManifestWithUnquotedValuesAndTargetAccountAndDifferentSubtitle(t *testing.T, targetAccount sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicManifestWithUnquotedValuesAndTargetAccount(t, "with_target_accounts_and_different_subtitle_", "different_subtitle", targetAccount)
}

func (c *ListingClient) basicManifest(t *testing.T, titleSuffix string, subtitle string) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	return fmt.Sprintf(`title: "%s"
subtitle: "%s"
description: "description"
listing_terms:
  type: "OFFLINE"
`, title, subtitle), title
}

func (c *ListingClient) basicManifestWithUnquotedValues(t *testing.T, titleSuffix string, subtitle string) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	return fmt.Sprintf(`title: %s
subtitle: %s
description: description
listing_terms:
  type: OFFLINE
`, title, subtitle), title
}

func (c *ListingClient) basicManifestWithTargetAccount(t *testing.T, titleSuffix string, subtitle string, targetAccount sdk.AccountIdentifier) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	return fmt.Sprintf(`title: "%s"
subtitle: "%s"
description: "description"
listing_terms:
  type: "OFFLINE"
targets:
  accounts: [%s.%s]
`, title, subtitle, targetAccount.OrganizationName(), targetAccount.AccountName()), title
}

func (c *ListingClient) basicManifestWithUnquotedValuesAndTargetAccount(t *testing.T, titleSuffix string, subtitle string, targetAccount sdk.AccountIdentifier) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	return fmt.Sprintf(`title: %s
subtitle: %s
description: description
listing_terms:
  type: OFFLINE
targets:
  accounts: [%s.%s]
`, title, subtitle, targetAccount.OrganizationName(), targetAccount.AccountName()), title
}
