package helpers

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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

// CreateOrganization creates an organization listing. Organization listings only have a
// dedicated CREATE command (CREATE ORGANIZATION LISTING); all the other operations (ALTER,
// DROP, SHOW, DESCRIBE) are shared with regular listings, hence the rest of ListingClient is reused.
func (c *ListingClient) CreateOrganization(t *testing.T) (*sdk.Listing, func()) {
	t.Helper()
	return c.CreateOrganizationWithId(t, c.ids.RandomAccountObjectIdentifier())
}

func (c *ListingClient) CreateOrganizationWithId(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.Listing, func()) {
	t.Helper()
	ctx := context.Background()

	manifest, _ := c.OrganizationBasicManifest(t)
	err := c.client().CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
		WithAs(manifest).
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

func (c *ListingClient) Describe(t *testing.T, id sdk.AccountObjectIdentifier) (*sdk.ListingDetails, error) {
	t.Helper()
	return c.client().Describe(context.Background(), sdk.NewDescribeListingRequest(id))
}

func (c *ListingClient) ShowVersions(t *testing.T, id sdk.AccountObjectIdentifier) ([]sdk.ListingVersion, error) {
	t.Helper()
	return c.client().ShowVersions(context.Background(), sdk.NewShowVersionsListingRequest(id))
}

// RequestListingAndWaitForSuccess calls SYSTEM$REQUEST_LISTING_AND_WAIT for the given listing
// global name with the given in-call wait and fails the test unless the procedure reports
// success.
//
// See: https://docs.snowflake.com/en/sql-reference/stored-procedures/system_request_listing_and_wait
func (c *ListingClient) RequestListingAndWaitForSuccess(t *testing.T, globalName string, waitInMinutes int) {
	t.Helper()

	rows, err := c.context.client.QueryUnsafe(context.Background(), fmt.Sprintf("CALL SYSTEM$REQUEST_LISTING_AND_WAIT('%s', %d)", globalName, waitInMinutes))
	require.NoError(t, err)
	require.Len(t, rows, 1)

	statusPtr := rows[0]["SYSTEM$REQUEST_LISTING_AND_WAIT"]
	require.NotNil(t, statusPtr)
	require.NotNil(t, *statusPtr)

	status, ok := (*statusPtr).(string)
	require.True(t, ok)
	require.Truef(t, strings.HasPrefix(status, "Success:"), "SYSTEM$REQUEST_LISTING_AND_WAIT did not succeed: %s", status)
}

// AcceptLegalTermsWithRetry calls SYSTEM$ACCEPT_LEGAL_TERMS for the given listing global name
// and retries until it succeeds or the timeout elapses. This doubles as a visibility probe: a
// freshly-published cross-account listing's global name is not immediately resolvable on the
// consumer account, and SYSTEM$ACCEPT_LEGAL_TERMS is the first consumer-side call that returns
// a clear SQL error ("does not exist or not authorized") in that window.
//
// See: https://docs.snowflake.com/en/sql-reference/stored-procedures/system_accept_legal_terms
func (c *ListingClient) AcceptLegalTermsWithRetry(t *testing.T, globalName string, timeout time.Duration, tick time.Duration) {
	t.Helper()

	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		_, err := c.context.client.QueryUnsafe(context.Background(), fmt.Sprintf("CALL SYSTEM$ACCEPT_LEGAL_TERMS('DATA_EXCHANGE_LISTING', '%s')", globalName))
		assert.NoError(collect, err)
	}, timeout, tick)
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

func (c *ListingClient) BasicManifestWithTargetAccounts(t *testing.T, targetAccounts ...sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicManifestWithTargetAccount(t, "with_target_accounts_", "subtitle", targetAccounts...)
}

func (c *ListingClient) BasicManifestWithTargetAccountsAndDifferentSubtitle(t *testing.T, targetAccounts ...sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicManifestWithTargetAccount(t, "with_target_accounts_and_different_subtitle_", "different_subtitle", targetAccounts...)
}

func (c *ListingClient) BasicManifestWithUnquotedValuesAndTargetAccounts(t *testing.T, targetAccounts ...sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicManifestWithUnquotedValuesAndTargetAccount(t, "with_target_accounts_", "subtitle", targetAccounts...)
}

func (c *ListingClient) BasicManifestWithUnquotedValuesAndTargetAccountsAndDifferentSubtitle(t *testing.T, targetAccounts ...sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicManifestWithUnquotedValuesAndTargetAccount(t, "with_target_accounts_and_different_subtitle_", "different_subtitle", targetAccounts...)
}

// OrganizationBasicManifest returns a minimal valid organization listing manifest; see
// https://docs.snowflake.com/en/user-guide/collaboration/listings/organizational/org-listing-manifest-reference.
func (c *ListingClient) OrganizationBasicManifest(t *testing.T) (string, string) {
	t.Helper()
	return c.basicOrganizationManifest(t, "basic_", "subtitle")
}

func (c *ListingClient) OrganizationBasicManifestWithDifferentSubtitle(t *testing.T) (string, string) {
	t.Helper()
	return c.basicOrganizationManifest(t, "basic_with_diff_subtitle_", "different_subtitle")
}

func (c *ListingClient) OrganizationBasicManifestWithTargetAccounts(t *testing.T, targetAccounts ...sdk.AccountIdentifier) (string, string) {
	t.Helper()
	return c.basicOrganizationManifestWithTargetAccounts(t, "with_target_accounts_", "subtitle", targetAccounts...)
}

func (c *ListingClient) basicOrganizationManifest(t *testing.T, titleSuffix string, subtitle string) (string, string) {
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

func (c *ListingClient) basicOrganizationManifestWithTargetAccounts(t *testing.T, titleSuffix string, subtitle string, targetAccounts ...sdk.AccountIdentifier) (string, string) {
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

func (c *ListingClient) basicManifest(t *testing.T, titleSuffix string, subtitle string) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	return fmt.Sprintf(`title: "%s"
subtitle: "%s"
description: "description"
listing_terms:
  type: "STANDARD"
`, title, subtitle), title
}

func (c *ListingClient) basicManifestWithUnquotedValues(t *testing.T, titleSuffix string, subtitle string) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	return fmt.Sprintf(`title: %s
subtitle: %s
description: description
listing_terms:
  type: STANDARD
`, title, subtitle), title
}

func (c *ListingClient) basicManifestWithTargetAccount(t *testing.T, titleSuffix string, subtitle string, targetAccounts ...sdk.AccountIdentifier) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	mappedTargetAccounts := collections.Map(targetAccounts, func(id sdk.AccountIdentifier) string {
		return fmt.Sprintf("%s.%s", id.OrganizationName(), id.AccountName())
	})
	return fmt.Sprintf(`title: "%s"
subtitle: "%s"
description: "description"
listing_terms:
  type: "STANDARD"
targets:
  accounts: [%s]
`, title, subtitle, collections.JoinStrings(mappedTargetAccounts, ", ")), title
}

func (c *ListingClient) basicManifestWithUnquotedValuesAndTargetAccount(t *testing.T, titleSuffix string, subtitle string, targetAccounts ...sdk.AccountIdentifier) (string, string) {
	t.Helper()
	title := c.ids.WithTestObjectSuffix(titleSuffix)
	mappedTargetAccounts := collections.Map(targetAccounts, func(id sdk.AccountIdentifier) string {
		return fmt.Sprintf("%s.%s", id.OrganizationName(), id.AccountName())
	})
	return fmt.Sprintf(`title: %s
subtitle: %s
description: description
listing_terms:
  type: STANDARD
targets:
  accounts: [%s]
`, title, subtitle, collections.JoinStrings(mappedTargetAccounts, ", ")), title
}
