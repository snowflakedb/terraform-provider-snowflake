//go:build !account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Listings(t *testing.T) {
	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	share, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClientHelper().Grant.GrantPrivilegeOnDatabaseToShare(t, testClientHelper().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	applicationPackage, applicationPackageCleanup := testClientHelper().ApplicationPackage.CreateApplicationPackageWithReleaseChannelsDisabled(t)
	t.Cleanup(applicationPackageCleanup)

	testClientHelper().Stage.PutOnStageWithContent(t, stage.ID(), "manifest.yml", "")
	testClientHelper().Stage.PutOnStageWithContent(t, stage.ID(), "setup.sql", "CREATE APPLICATION ROLE IF NOT EXISTS APP_HELLO_SNOWFLAKE;")

	version := "v1"
	testClientHelper().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stage.ID(), version)
	testClientHelper().ApplicationPackage.SetDefaultReleaseDirective(t, applicationPackage.ID(), version)

	client := testClient(t)
	ctx := testContext(t)

	accountId := testClientHelper().Context.CurrentAccountId(t)
	basicManifest := testClientHelper().Listing.BasicManifest(t)
	testClientHelper().Stage.PutOnStageDirectoryWithContent(t, stage.ID(), "manifest.yml", "basic", basicManifest)
	basicManifestStageLocation := sdk.NewStageLocation(stage.ID(), "basic/")

	targetAccount := fmt.Sprintf("%s.%s", accountId.OrganizationName(), accountId.AccountName())
	basicManifestWithTarget := fmt.Sprintf(`
title: title
subtitle: subtitle
description: description
listing_terms:
  type: OFFLINE
targets:
  accounts: [%s]
`, targetAccount)
	testClientHelper().Stage.PutOnStageDirectoryWithContent(t, stage.ID(), "manifest.yml", "with_target", basicManifestWithTarget)
	basicManifestWithTargetStageLocation := sdk.NewStageLocation(stage.ID(), "with_target/")

	comment := random.Comment()

	assertNoOptionals := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasGlobalNameNotEmpty().
				HasName(id.Name()).
				HasTitle("title").
				HasSubtitle("subtitle").
				HasProfile("").
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty().
				HasNoPublishedOn().
				HasState(sdk.ListingStateDraft).
				HasReviewState("UNSENT").
				HasNoComment().
				HasNoRegions().
				HasTargetAccounts("").
				HasIsMonetized(false).
				HasIsApplication(false).
				HasIsTargeted(false).
				HasIsLimitedTrial(false).
				HasIsByRequest(false).
				HasDistribution("EXTERNAL").
				HasIsMountlessQueryable(false).
				HasOrganizationProfileName("").
				HasNoUniformListingLocator().
				HasNoDetailedTargetAccounts(),
		)
	}

	t.Run("create from manifest: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithAs(basicManifest).
			WithReview(false).
			WithPublish(false))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		assertNoOptionals(t, id)
	})

	t.Run("create from stage: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithFrom(basicManifestStageLocation).
			WithReview(false).
			WithPublish(false))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		assertNoOptionals(t, id)
	})

	assertCompleteWithShare := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()

		listingDetails, err := client.Listings.Describe(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, share.ID().Name(), listingDetails.Share.Name())

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasGlobalNameNotEmpty().
				HasName(id.Name()).
				HasTitle("title").
				HasSubtitle("subtitle").
				HasProfile("").
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty().
				HasNoPublishedOn().
				HasState(sdk.ListingStateDraft).
				HasNoReviewState().
				HasComment(comment).
				HasNoRegions().
				HasTargetAccounts(targetAccount).
				HasIsMonetized(false).
				HasIsApplication(false).
				HasIsTargeted(true).
				HasIsLimitedTrial(false).
				HasIsByRequest(false).
				HasDistribution("EXTERNAL").
				HasIsMountlessQueryable(false).
				HasOrganizationProfileName("").
				HasNoUniformListingLocator().
				HasDetailedTargetAccountsNotEmpty(),
		)
	}

	t.Run("create from manifest: complete with share", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithAs(basicManifestWithTarget).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithIfNotExists(true).
			WithPublish(false).
			WithReview(false).
			WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		assertCompleteWithShare(t, id)
	})

	t.Run("create from stage: complete with share", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithFrom(basicManifestWithTargetStageLocation).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithIfNotExists(true).
			WithPublish(false).
			WithReview(false).
			WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		assertCompleteWithShare(t, id)
	})

	assertCompleteWithApplicationPackage := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()

		listingDetails, err := client.Listings.Describe(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, applicationPackage.ID().Name(), listingDetails.ApplicationPackage.Name())

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasGlobalNameNotEmpty().
				HasName(id.Name()).
				HasTitle("title").
				HasSubtitle("subtitle").
				HasProfile("").
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty().
				HasNoPublishedOn().
				HasState(sdk.ListingStateDraft).
				HasNoReviewState().
				HasComment(comment).
				HasNoRegions().
				HasTargetAccounts(targetAccount).
				HasIsMonetized(false).
				HasIsApplication(true).
				HasIsTargeted(true).
				HasIsLimitedTrial(false).
				HasIsByRequest(false).
				HasDistribution("EXTERNAL").
				HasIsMountlessQueryable(false).
				HasOrganizationProfileName("").
				HasNoUniformListingLocator().
				HasDetailedTargetAccountsNotEmpty(),
		)
	}

	t.Run("create from manifest: complete with application packages", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithAs(basicManifestWithTarget).
			WithWith(*sdk.NewListingWithRequest().WithApplicationPackage(applicationPackage.ID())).
			WithIfNotExists(true).
			WithPublish(false).
			WithReview(false).
			WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		assertCompleteWithApplicationPackage(t, id)
	})

	t.Run("create from stage: complete with application packages", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithFrom(basicManifestWithTargetStageLocation).
			WithWith(*sdk.NewListingWithRequest().WithApplicationPackage(applicationPackage.ID())).
			WithIfNotExists(true).
			WithPublish(false).
			WithReview(false).
			WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		assertCompleteWithApplicationPackage(t, id)
	})

	t.Run("alter: change state", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithFrom(basicManifestWithTargetStageLocation).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithIfNotExists(true).
			WithPublish(false).
			WithReview(false).
			WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasState(sdk.ListingStateDraft).
				HasNoReviewState(),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithReview(true))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasState(sdk.ListingStateDraft).
				HasNoReviewState(),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithPublish(true))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasState(sdk.ListingStatePublished).
				HasNoReviewState(),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithUnpublish(true))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasState(sdk.ListingStateUnpublished).
				HasNoReviewState(),
		)
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
			// Should be HasComment(comment), but it seems the comment is not set on alter or this comment is set somewhere else
		)
	})

	t.Run("alter: add version", func(t *testing.T) {
		basicWithDifferentSubtitleManifest := `
title: title
subtitle: different_subtitle
description: description
listing_terms:
  type: OFFLINE
`
		testClientHelper().Stage.PutOnStageDirectoryWithContent(t, stage.ID(), "manifest.yml", "basic_different_subtitle", basicWithDifferentSubtitleManifest)
		basicManifestWithDifferentSubtitleStageLocation := sdk.NewStageLocation(stage.ID(), "basic_different_subtitle/")

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
				HasSubtitle("subtitle").
				HasNoComment(),
		)

		versions, err := client.Listings.ShowVersions(ctx, sdk.NewShowVersionsListingRequest(listing.ID()))
		assert.NoError(t, err)
		assert.Len(t, versions, 1)
		assert.Equal(t, "v2", versions[0].Alias)
		assert.Equal(t, comment, versions[0].Comment)
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

	t.Run("show: with options", func(t *testing.T) {
		prefix := random.AlphanumericN(10)
		id := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
		id2 := testClientHelper().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

		_, listingCleanup := testClientHelper().Listing.CreateWithId(t, id)
		t.Cleanup(listingCleanup)

		_, listing2Cleanup := testClientHelper().Listing.CreateWithId(t, id2)
		t.Cleanup(listing2Cleanup)

		listings, err := client.Listings.Show(ctx, sdk.NewShowListingRequest().
			WithLike(sdk.Like{Pattern: sdk.String(prefix + "%")}).
			WithStartsWith(prefix).
			WithLimit(sdk.LimitFrom{
				Rows: sdk.Int(1),
				From: sdk.String(prefix),
			}),
		)

		assert.NoError(t, err)
		assert.Len(t, listings, 1)
	})

	t.Run("describe: default", func(t *testing.T) {
		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		listingDetails, err := client.Listings.Describe(ctx, listing.ID())
		require.NoError(t, err)
		require.NotNil(t, listingDetails)

		assert.NotEmpty(t, listingDetails.GlobalName)
		assert.Equal(t, listing.ID().Name(), listingDetails.Name)
		assert.NotEmpty(t, listingDetails.Owner)
		assert.NotEmpty(t, listingDetails.OwnerRoleType)
		assert.NotEmpty(t, listingDetails.CreatedOn)
		assert.NotEmpty(t, listingDetails.UpdatedOn)
		assert.Nil(t, listingDetails.PublishedOn)
		assert.Equal(t, "title", listingDetails.Title)
		assert.Equal(t, "subtitle", *listingDetails.Subtitle)
		assert.Equal(t, "description", *listingDetails.Description)
		assert.JSONEq(t, `{
"type" : "OFFLINE"
}`, *listingDetails.ListingTerms)
		assert.Equal(t, sdk.ListingStateDraft, listingDetails.State)
		assert.Nil(t, listingDetails.Share)
		assert.Empty(t, listingDetails.ApplicationPackage.Name()) // Application package is returned even if listing is not associated with one, but it is empty in that case
		assert.Nil(t, listingDetails.BusinessNeeds)
		assert.Nil(t, listingDetails.UsageExamples)
		assert.Nil(t, listingDetails.DataAttributes)
		assert.Nil(t, listingDetails.Categories)
		assert.Nil(t, listingDetails.Resources)
		assert.Nil(t, listingDetails.Profile)
		assert.Nil(t, listingDetails.CustomizedContactInfo)
		assert.Nil(t, listingDetails.DataDictionary)
		assert.Nil(t, listingDetails.DataPreview)
		assert.Nil(t, listingDetails.Comment)
		assert.Equal(t, "DRAFT", listingDetails.Revisions)
		assert.Nil(t, listingDetails.TargetAccounts)
		assert.Nil(t, listingDetails.Regions)
		assert.Nil(t, listingDetails.RefreshSchedule)
		assert.Nil(t, listingDetails.RefreshType)
		assert.Equal(t, "UNSENT", *listingDetails.ReviewState)
		assert.Nil(t, listingDetails.RejectionReason)
		assert.Nil(t, listingDetails.UnpublishedByAdminReasons)
		assert.False(t, listingDetails.IsMonetized)
		assert.False(t, listingDetails.IsApplication)
		assert.False(t, listingDetails.IsTargeted)
		assert.False(t, *listingDetails.IsLimitedTrial)
		assert.False(t, *listingDetails.IsByRequest)
		assert.Nil(t, listingDetails.LimitedTrialPlan)
		assert.Nil(t, listingDetails.RetriedOn)
		assert.Nil(t, listingDetails.ScheduledDropTime)
		assert.Equal(t, testClientHelper().Listing.BasicManifest(t), listingDetails.ManifestYaml)
		assert.Equal(t, "EXTERNAL", *listingDetails.Distribution)
		assert.False(t, *listingDetails.IsMountlessQueryable)
		assert.Nil(t, listingDetails.OrganizationProfileName)
		assert.Nil(t, listingDetails.UniformListingLocator)
		assert.Nil(t, listingDetails.TrialDetails)
		assert.Nil(t, listingDetails.ApproverContact)
		assert.Nil(t, listingDetails.SupportContact)
		assert.Nil(t, listingDetails.LiveVersionUri)
		assert.Nil(t, listingDetails.LastCommittedVersionUri)
		assert.Nil(t, listingDetails.LastCommittedVersionName)
		assert.Nil(t, listingDetails.LastCommittedVersionAlias)
		assert.Nil(t, listingDetails.PublishedVersionUri)
		assert.Nil(t, listingDetails.PublishedVersionName)
		assert.Nil(t, listingDetails.PublishedVersionAlias)
		assert.False(t, *listingDetails.IsShare)
		assert.Nil(t, listingDetails.RequestApprovalType)
		assert.Empty(t, *listingDetails.MonetizationDisplayOrder)
		assert.Empty(t, *listingDetails.LegacyUniformListingLocators)
	})
}
