//go:build !account_level_tests

package testint

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestInt_Listings(t *testing.T) {
	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	share, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, testClientHelper().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	applicationPackage, applicationPackageCleanup := testClientHelper().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(applicationPackageCleanup)
	//testClientHelper().ApplicationPackage.SetQaReleaseChannel(t, applicationPackage.ID(), stage.ID())
	//testClientHelper().ApplicationPackage.SetDefaultReleaseDirective(t, applicationPackage.ID())

	client := testClient(t)
	ctx := testContext(t)

	accountId := testClientHelper().Context.CurrentAccountId(t)
	basicManifest := `
title: title
subtitle: subtitle
description: description
listing_terms:
  type: OFFLINE
`
	testClientHelper().Stage.PutOnStageDirectoryWithContent(t, stage.ID(), "manifest.yml", "basic", basicManifest)
	basicManifestStageLocation := sdk.NewStageLocation(stage.ID(), "basic/")

	basicWithDifferentSubtitleManifest := `
title: title
subtitle: different_subtitle
description: description
listing_terms:
  type: OFFLINE
`
	testClientHelper().Stage.PutOnStageDirectoryWithContent(t, stage.ID(), "manifest.yml", "basic_different_subtitle", basicWithDifferentSubtitleManifest)
	basicManifestWithDifferentSubtitleStageLocation := sdk.NewStageLocation(stage.ID(), "basic_different_subtitle/")

	basicManifestWithTarget := fmt.Sprintf(`
title: title
subtitle: subtitle
description: description
listing_terms:
  type: OFFLINE
targets:
  accounts: [%s.%s]
`, accountId.OrganizationName(), accountId.AccountName())
	testClientHelper().Stage.PutOnStageDirectoryWithContent(t, stage.ID(), "manifest.yml", "with_target", basicManifestWithTarget)
	basicManifestWithTargetStageLocation := sdk.NewStageLocation(stage.ID(), "with_target/")

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
			WithAs(basicManifestWithTarget).
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

	t.Run("create from stage: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).WithFrom(basicManifestStageLocation).
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

	t.Run("create from stage: complete with stage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).WithFrom(basicManifestWithTargetStageLocation).
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
		)
	})

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

	t.Run("alter: change state", func(t *testing.T) {
		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		assertThatObject(t,
			objectassert.ListingFromObject(t, listing).
				HasState(sdk.ListingStateDraft).
				HasReviewState("UNSENT"),
		)

		err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(listing.ID()).WithReview(true))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, listing.ID()).
				HasState(sdk.ListingStateDraft).
				HasReviewState("UNSENT"),
		)

		// TODO: Too much to fulfill to check
		//err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithPublish(true))
		//assert.NoError(t, err)
		//
		//assertThatObject(t,
		//	objectassert.Listing(t, id).
		//		HasName(id.Name()).
		//		HasState(sdk.ListingStatePublished).
		//		HasReviewState(""),
		//)
		//
		//err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithUnpublish(true))
		//assert.NoError(t, err)
		//
		//assertThatObject(t,
		//	objectassert.Listing(t, id).
		//		HasName(id.Name()).
		//		HasState(sdk.ListingStateUnpublished).
		//		HasReviewState(""),
		//)
	})

	t.Run("alter: change manifest with optional values", func(t *testing.T) {
		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		assertThatObject(t,
			objectassert.ListingFromObject(t, listing).
				HasSubtitle("subtitle").
				HasNoComment(),
		)

		basicManifestWithDifferentSubtitle := `
title: title
subtitle: different_subtitle
description: description
listing_terms:
  type: OFFLINE
`

		err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(listing.ID()).
			WithAlterListingAs(*sdk.NewAlterListingAsRequest(basicManifestWithDifferentSubtitle).
				WithPublish(false).
				WithReview(false).
				WithComment(comment),
			))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, listing.ID()).
				HasSubtitle("different_subtitle").
				HasNoComment(),
			// TODO: Should be HasComment(comment), but it seems the comment is not set on alter or this comment is set somewhere else
		)
	})

	t.Run("alter: add version", func(t *testing.T) {
		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		assertThatObject(t,
			objectassert.ListingFromObject(t, listing).
				HasSubtitle("subtitle").
				HasNoComment(),
		)

		err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(listing.ID()).
			WithAddVersion(*sdk.NewAddListingVersionRequest("v2", basicManifestWithDifferentSubtitleStageLocation).
				WithIfNotExists(true).
				WithComment(comment)))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.ListingFromObject(t, listing).
				HasSubtitle("different_subtitle").
				HasNoComment(),
			// TODO: Should be HasComment(comment), but it seems the comment is not set on alter or this comment is set somewhere else
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		newId := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(listing.ID()).WithRenameTo(newId))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, newId))

		_, err = client.Listings.ShowByID(ctx, listing.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

		assertThatObject(t, objectassert.Listing(t, newId).HasName(newId.Name()))
	})

	t.Run("alter: set", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		newComment := random.Comment()

		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithAs(basicManifest).
			WithReview(false).
			WithPublish(false).
			WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		assertThatObject(t, objectassert.Listing(t, id).
			HasName(id.Name()).
			HasComment(comment),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithSet(*sdk.NewListingSetRequest().WithComment(newComment)))
		assert.NoError(t, err)

		assertThatObject(t, objectassert.Listing(t, id).
			HasName(id.Name()).
			HasComment(newComment),
		)
	})

	t.Run("drop", func(t *testing.T) {
		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		err := client.Listings.Drop(ctx, sdk.NewDropListingRequest(listing.ID()))
		assert.NoError(t, err)

		_, err = client.Listings.ShowByID(ctx, listing.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

		err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(listing.ID()))
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

		err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(listing.ID()).WithIfExists(true))
		assert.NoError(t, err)
	})

	t.Run("drop safely", func(t *testing.T) {
		// TODO: Show how it behaves for listings with draft/published statuses
	})

	t.Run("show: with options", func(t *testing.T) {
		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		listings, err := client.Listings.Show(ctx, sdk.NewShowListingRequest().
			WithLike(sdk.Like{Pattern: sdk.String(listing.ID().Name())}).
			WithStartsWith(listing.ID().Name()),
		)
		//WithLimit(sdk.LimitFrom{
		//	Rows: sdk.Int(1),
		//	From: sdk.String(listing.ID().Name()),
		//}))

		assert.NoError(t, err)
		assert.Len(t, listings, 1)
		assert.Equal(t, listing.ID().Name(), listings[0].Name)
	})

	t.Run("describe: default", func(t *testing.T) {
		// TODO: Revision draft and published
	})
}
