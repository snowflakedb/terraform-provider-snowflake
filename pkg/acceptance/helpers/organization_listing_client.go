package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

// OrganizationListingClient manages organization listings in tests.
// Organization listings only have a dedicated CREATE command; all the other operations
// (ALTER, DROP, SHOW, DESCRIBE) are shared with regular listings, hence the sdk.Listings interface usage.
type OrganizationListingClient struct {
	context *TestClientContext
	ids     *IdsGenerator
}

func NewOrganizationListingClient(context *TestClientContext, idsGenerator *IdsGenerator) *OrganizationListingClient {
	return &OrganizationListingClient{
		context: context,
		ids:     idsGenerator,
	}
}

func (c *OrganizationListingClient) client() sdk.Listings {
	return c.context.client.Listings
}

func (c *OrganizationListingClient) Create(t *testing.T) (*sdk.Listing, func()) {
	t.Helper()
	return c.CreateWithId(t, c.ids.RandomAccountObjectIdentifier())
}

func (c *OrganizationListingClient) CreateWithId(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Listing, func()) {
	t.Helper()
	ctx := context.Background()

	manifest, _ := c.BasicManifest(t)
	err := c.client().CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
		WithAs(manifest).
		WithPublish(false),
	)
	assert.NoError(t, err)

	listing, err := c.client().ShowByID(ctx, id)
	assert.NoError(t, err)

	return listing, c.DropFunc(t, id)
}

func (c *OrganizationListingClient) Alter(t *testing.T, req *sdk.AlterListingRequest) {
	t.Helper()
	ctx := context.Background()

	err := c.client().Alter(ctx, req)
	assert.NoError(t, err)
}

func (c *OrganizationListingClient) DropFunc(t *testing.T, id sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()

	return func() {
		assert.NoError(t, c.client().DropSafely(ctx, id))
	}
}

func (c *OrganizationListingClient) Show(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Listing, error) {
	t.Helper()
	return c.client().ShowByID(context.Background(), id)
}

func (c *OrganizationListingClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ListingDetails, error) {
	t.Helper()
	return c.client().Describe(context.Background(), sdk.NewDescribeListingRequest(id))
}

func (c *OrganizationListingClient) BasicManifest(t *testing.T) (string, string) {
	t.Helper()
	return c.basicManifest(t, "basic_", "subtitle")
}

func (c *OrganizationListingClient) BasicManifestWithDifferentSubtitle(t *testing.T) (string, string) {
	t.Helper()
	return c.basicManifest(t, "basic_with_diff_subtitle_", "different_subtitle")
}

func (c *OrganizationListingClient) BasicManifestWithTargetAccounts(t *testing.T, targetAccounts ...sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicManifestWithTargetAccounts(t, "with_target_accounts_", "subtitle", targetAccounts...)
}

// basicManifest returns a minimal valid organization listing manifest; see
// https://docs.snowflake.com/en/user-guide/collaboration/listings/organizational/org-listing-manifest-reference.
func (c *OrganizationListingClient) basicManifest(t *testing.T, titleSuffix string, subtitle string) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	return fmt.Sprintf(`title: "%s"
subtitle: "%s"
description: "description"
organization_targets:
  access:
  - all_internal_accounts: true
locations:
  access_regions:
  - name: "ALL"
`, title, subtitle), title
}

func (c *OrganizationListingClient) basicManifestWithTargetAccounts(t *testing.T, titleSuffix string, subtitle string, targetAccounts ...sdk.AccountIdentifier) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	accessEntries := collections.Map(targetAccounts, func(id sdk.AccountIdentifier) string {
		return fmt.Sprintf("  - account: \"%s\"", id.AccountName())
	})
	return fmt.Sprintf(`title: "%s"
subtitle: "%s"
description: "description"
organization_targets:
  access:
%s
locations:
  access_regions:
  - name: "ALL"
`, title, subtitle, collections.JoinStrings(accessEntries, "\n")), title
}
