//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_OrganizationListings(t *testing.T) {
	// Organization listings only have a dedicated CREATE command (CREATE ORGANIZATION LISTING);
	// all the other commands (ALTER, DROP, SHOW, DESCRIBE) are shared with regular listings.
	// A role used to create organization listings must have at least the CREATE ORGANIZATION LISTING
	// (or CREATE LISTING) privilege on the account.

	share, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, testClientHelper().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	client := testClient(t)
	ctx := testContext(t)

	comment := random.Comment()

	t.Run("create from manifest inlined: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		manifest, _ := testClientHelper().Listing.OrganizationBasicManifest(t)

		err := client.Listings.CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
			WithAs(manifest).
			WithPublish(false))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasName(id.Name()).
				HasGlobalNameNotEmpty().
				HasState(sdk.ListingStateDraft).
				HasNoComment().
				HasNoPublishedOn(),
		)
	})

	t.Run("create from manifest inlined: with share", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		manifest, _ := testClientHelper().Listing.OrganizationBasicManifest(t)

		err := client.Listings.CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithAs(manifest).
			WithPublish(false))
		require.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		listing, err := client.Listings.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id.Name(), listing.Name)
	})

	t.Run("alter: publish and unpublish", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		manifest, _ := testClientHelper().Listing.OrganizationBasicManifest(t)

		require.NoError(t, client.Listings.CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithAs(manifest).
			WithPublish(false)))
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		require.NoError(t, client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithPublish(true)))
		listing, err := client.Listings.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.ListingStatePublished, listing.State)

		require.NoError(t, client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithUnpublish(true)))
		listing, err = client.Listings.ShowByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, sdk.ListingStateUnpublished, listing.State)
	})

	t.Run("alter: update manifest via AS", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		manifest, _ := testClientHelper().Listing.OrganizationBasicManifest(t)
		updatedManifest, _ := testClientHelper().Listing.OrganizationBasicManifestWithDifferentSubtitle(t)

		require.NoError(t, client.Listings.CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
			WithAs(manifest).
			WithPublish(false)))
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		require.NoError(t, client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
			WithAlterListingAs(*sdk.NewAlterListingAsRequest(updatedManifest))))
	})

	t.Run("alter: set and unset comment", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		manifest, _ := testClientHelper().Listing.OrganizationBasicManifest(t)

		require.NoError(t, client.Listings.CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
			WithAs(manifest).
			WithPublish(false)))
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		require.NoError(t, client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
			WithSet(*sdk.NewListingSetRequest().WithComment(comment))))

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasComment(comment),
		)

		require.NoError(t, client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
			WithUnset(*sdk.NewListingUnsetRequest().WithComment(true))))

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasNoComment(),
		)
	})

	t.Run("alter: rename to", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		newId := testClientHelper().Ids.RandomAccountObjectIdentifier()
		manifest, _ := testClientHelper().Listing.OrganizationBasicManifest(t)

		require.NoError(t, client.Listings.CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
			WithAs(manifest).
			WithPublish(false)))
		t.Cleanup(testClientHelper().Listing.DropFunc(t, newId))

		require.NoError(t, client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithRenameTo(newId)))

		_, err := client.Listings.ShowByID(ctx, newId)
		require.NoError(t, err)
	})

	t.Run("show: with like filter", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		manifest, _ := testClientHelper().Listing.OrganizationBasicManifest(t)

		require.NoError(t, client.Listings.CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
			WithAs(manifest).
			WithPublish(false)))
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		listings, err := client.Listings.Show(ctx, sdk.NewShowListingRequest().
			WithLike(sdk.Like{Pattern: sdk.String(id.Name())}))
		require.NoError(t, err)
		assert.Len(t, listings, 1)
		assert.Equal(t, id.Name(), listings[0].Name)
	})

	t.Run("describe", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		manifest, _ := testClientHelper().Listing.OrganizationBasicManifest(t)

		require.NoError(t, client.Listings.CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
			WithAs(manifest).
			WithPublish(false)))
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		details, err := client.Listings.Describe(ctx, sdk.NewDescribeListingRequest(id))
		require.NoError(t, err)
		assert.Equal(t, id.Name(), details.Name)
		assert.NotEmpty(t, details.ManifestYaml)
	})

	t.Run("drop safely: published listing is unpublished first", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		manifest, _ := testClientHelper().Listing.OrganizationBasicManifest(t)

		require.NoError(t, client.Listings.CreateOrganization(ctx, sdk.NewCreateOrganizationListingRequest(id).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithAs(manifest).
			WithPublish(true)))

		require.NoError(t, client.Listings.DropSafely(ctx, id))

		_, err := client.Listings.ShowByID(ctx, id)
		assert.Error(t, err)
	})
}
