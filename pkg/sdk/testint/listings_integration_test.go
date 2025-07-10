//go:build !account_level_tests

package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestInt_Listings(t *testing.T) {
	//stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	//t.Cleanup(stageCleanup)
	//
	//stageLocationPath := "dir/subdir"
	//stageLocation := sdk.NewStageLocation(stage.ID(), stageLocationPath)

	share, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, testClientHelper().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	applicationPackage, applicationPackageCleanup := testClientHelper().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(applicationPackageCleanup)
	testClientHelper().ApplicationPackage.SetDefaultReleaseDirective(t, applicationPackage.ID())

	client := testClient(t)
	ctx := testContext(t)

	basicManifest := testClientHelper().Listing.BasicManifest(t)
	manifestWithTarget := testClientHelper().Listing.ManifestBuilder(t).
		WithTitle("title").
		WithSubtitle("subtitle").
		WithDescription("description").
		WithListingTermsType(helpers.ListingTermsTypeOffline).
		WithTarget(testClientHelper().Context.CurrentAccountId(t)).
		Build()

	comment := random.Comment()

	t.Run("create from manifest: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithAs(basicManifest).
			WithReview(false).
			WithPublish(false))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		// TODO: Add more assertions
		assertThatObject(t,
			objectassert.Listing(t, id).
				HasName(id.Name()).
				HasTitle("title").
				HasSubtitle("subtitle").
				HasState(sdk.ListingStateDraft),
		)
	})

	t.Run("create from manifest: complete with share", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithAs(manifestWithTarget).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithIfNotExists(true).
			WithPublish(true).
			WithReview(true).
			WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		// TODO: Assert more
		assertThatObject(t,
			objectassert.Listing(t, id).
				HasGlobalNameNotEmpty().
				HasName(id.Name()).
				HasTitle("title").
				HasComment(comment).
				HasState(sdk.ListingStatePublished),
			// Profile                 string
			// CreatedOn               string
			// UpdatedOn               string
			// PublishedOn             *string
			// ReviewState             *string
			// Comment                 *string
			// Owner                   string
			// OwnerRoleType           string
			// Regions                 *string
			// TargetAccounts          string
			// IsMonetized             bool
			// IsApplication           bool
			// IsTargeted              bool
			// IsLimitedTrial          *bool
			// IsByRequest             *bool
			// Distribution            *string
			// IsMountlessQueryable    *bool
			// RejectedOn              *string
			// OrganizationProfileName *string
			// UniformListingLocator   *string
			// DetailedTargetAccounts  *string
		)
	})

	t.Run("create from manifest: complete with application packages", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithAs(basicManifest).
			WithWith(*sdk.NewListingWithRequest().WithApplicationPackage(applicationPackage.ID())).
			WithIfNotExists(true).
			WithPublish(true).
			WithReview(true).
			WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		// TODO: Assert more
		assertThatObject(t,
			objectassert.Listing(t, id).
				HasName(id.Name()).
				HasTitle("title").
				HasState(sdk.ListingStatePublished),
		)
	})

	//t.Run("create from stage: no optionals", func(t *testing.T) {
	//	id := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//	err := client.Listings.CreateFromStage(ctx, sdk.NewCreateFromStageListingRequest(id, stageLocation))
	//	assert.NoError(t, err)
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
	//	assert.NoError(t, err)
	//})
	//
	//t.Run("create from stage: complete with stage", func(t *testing.T) {
	//	id := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//	err := client.Listings.CreateFromStage(ctx, sdk.NewCreateFromStageListingRequest(id, stageLocation).
	//		WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
	//		WithIfNotExists(true).
	//		WithPublish(true).
	//		WithReview(true))
	//	assert.NoError(t, err)
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
	//	assert.NoError(t, err)
	//})
	//
	//t.Run("create from stage: complete with application packages", func(t *testing.T) {
	//	id := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//	err := client.Listings.CreateFromStage(ctx, sdk.NewCreateFromStageListingRequest(id, stageLocation).
	//		WithWith(*sdk.NewListingWithRequest().WithApplicationPackage(applicationPackage.ID())).
	//		WithIfNotExists(true).
	//		WithPublish(true).
	//		WithReview(true))
	//	assert.NoError(t, err)
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
	//	assert.NoError(t, err)
	//})
	//
	//t.Run("alter: change state", func(t *testing.T) {
	//	id := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//	err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest))
	//	assert.NoError(t, err)
	//	t.Cleanup(testClientHelper().Listing.DropFunc(t, id))
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
	//		WithIfExists(true).
	//		WithReview(true))
	//	assert.NoError(t, err)
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
	//		WithIfExists(true).
	//		WithPublish(true))
	//	assert.NoError(t, err)
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
	//		WithIfExists(true).
	//		WithUnpublish(true))
	//	assert.NoError(t, err)
	//})
	//
	//t.Run("alter: change manifest with optional values", func(t *testing.T) {
	//	id := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//	comment := random.Comment()
	//	err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest))
	//	assert.NoError(t, err)
	//	t.Cleanup(testClientHelper().Listing.DropFunc(t, id))
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
	//		WithAlterListingAs(*sdk.NewAlterListingAsRequest(manifest).
	//			WithPublish(true).
	//			WithReview(true).
	//			WithComment(comment),
	//		))
	//	assert.NoError(t, err)
	//
	//	// TODO: Assert
	//})
	//
	//t.Run("alter: add version", func(t *testing.T) {
	//	id := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//	comment := random.Comment()
	//	versionName := random.AlphaN(10)
	//
	//	err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest))
	//	assert.NoError(t, err)
	//	t.Cleanup(testClientHelper().Listing.DropFunc(t, id))
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
	//		WithAddVersion(*sdk.NewAddListingVersionRequest(versionName, stageLocation).
	//			WithIfNotExists(true).
	//			WithComment(comment)))
	//	assert.NoError(t, err)
	//
	//	// TODO: Assert
	//})
	//
	//t.Run("alter: rename", func(t *testing.T) {
	//	id := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//	newId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//
	//	err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest))
	//	assert.NoError(t, err)
	//	t.Cleanup(testClientHelper().Listing.DropFunc(t, id))
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithRenameTo(newId))
	//	assert.NoError(t, err)
	//	t.Cleanup(testClientHelper().Listing.DropFunc(t, newId))
	//
	//	// TODO: Assert
	//})
	//
	//t.Run("alter: set", func(t *testing.T) {
	//	id := testClientHelper().Ids.RandomAccountObjectIdentifier()
	//	comment := random.Comment()
	//	newComment := random.Comment()
	//
	//	err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest).WithComment(comment))
	//	assert.NoError(t, err)
	//	t.Cleanup(testClientHelper().Listing.DropFunc(t, id))
	//
	//	// TODO: Assert
	//
	//	err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithSet(*sdk.NewListingSetRequest().WithComment(newComment)))
	//	assert.NoError(t, err)
	//
	//	// TODO: Assert
	//})
	//
	//t.Run("drop: existing", func(t *testing.T) {
	//	listing, listingCleanup := testClientHelper().Listing.Create(t)
	//	t.Cleanup(listingCleanup)
	//
	//	err := client.Listings.Drop(ctx, sdk.NewDropListingRequest(listing.ID()))
	//	assert.NoError(t, err)
	//
	//	listingAfterDrop, err := client.Listings.ShowByID(ctx, listing.ID())
	//	assert.NoError(t, err)
	//	// TODO: Assert listingAfterDrop
	//	_ = listingAfterDrop
	//})
	//
	////t.Run("show: default", func(t *testing.T) {
	////	listing, listingCleanup := testClientHelper().Listing.Create(t)
	////	t.Cleanup(listingCleanup)
	////
	////	listings, err := client.Listings.Show(ctx, sdk.NewShowListingRequest())
	////	assert.NoError(t, err)
	////	assert.Greater(t, len(listings), 1)
	////
	////	listingFound, err := collections.FindFirst(listings, func(l sdk.Listing) bool { return l.ID() == listing.ID() })
	////	assert.NoError(t, err)
	////})
	//
	//t.Run("show: with options", func(t *testing.T) {
	//	listing, listingCleanup := testClientHelper().Listing.Create(t)
	//	t.Cleanup(listingCleanup)
	//
	//	listings, err := client.Listings.Show(ctx, sdk.NewShowListingRequest().
	//		WithLike(sdk.Like{Pattern: sdk.String(listing.ID().Name())}).
	//		WithStartsWith(listing.ID().Name()).
	//		WithLimit(sdk.LimitFrom{
	//			Rows: sdk.Int(1),
	//			From: sdk.String(listing.ID().Name()),
	//		}))
	//	assert.NoError(t, err)
	//	assert.Equal(t, len(listings), 1)
	//})
	//
	//t.Run("describe: default", func(t *testing.T) {
	//})
}
