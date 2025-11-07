//go:build non_account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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

	basicManifest, basicManifestTitle := testClientHelper().Listing.BasicManifest(t)
	_ = testClientHelper().Stage.PutInLocationWithContent(t, stage.Location()+"/basic", "manifest.yml", basicManifest)
	basicManifestStageLocation := sdk.NewStageLocation(stage.ID(), "basic")

	basicManifestWithTarget, basicManifestWithTargetTitle := testClientHelper().Listing.BasicManifestWithTargetAccounts(t)
	testClientHelper().Stage.PutInLocationWithContent(t, stage.Location()+"/with_target", "manifest.yml", basicManifestWithTarget)
	basicManifestWithTargetStageLocation := sdk.NewStageLocation(stage.ID(), "with_target")

	comment := random.Comment()

	assertNoOptionals := func(t *testing.T, id sdk.AccountObjectIdentifier) {
		t.Helper()

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasGlobalNameNotEmpty().
				HasName(id.Name()).
				HasTitle(basicManifestTitle).
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

		listingDetails, err := client.Listings.Describe(ctx, sdk.NewDescribeListingRequest(id))
		assert.NoError(t, err)
		assert.Equal(t, share.ID().Name(), listingDetails.Share.Name())

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasGlobalNameNotEmpty().
				HasName(id.Name()).
				HasTitle(basicManifestWithTargetTitle).
				HasSubtitle("subtitle").
				HasProfile("").
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty().
				HasNoPublishedOn().
				HasState(sdk.ListingStateDraft).
				HasNoReviewState().
				HasComment(comment).
				HasNoRegions().
				HasTargetAccounts("").
				HasIsMonetized(false).
				HasIsApplication(false).
				HasIsTargeted(true).
				HasIsLimitedTrial(false).
				HasIsByRequest(false).
				HasDistribution("EXTERNAL").
				HasIsMountlessQueryable(false).
				HasOrganizationProfileName("").
				HasNoUniformListingLocator().
				HasNoDetailedTargetAccounts(),
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

		listingDetails, err := client.Listings.Describe(ctx, sdk.NewDescribeListingRequest(id))
		assert.NoError(t, err)
		assert.Equal(t, applicationPackage.ID().Name(), listingDetails.ApplicationPackage.Name())

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasGlobalNameNotEmpty().
				HasName(id.Name()).
				HasTitle(basicManifestWithTargetTitle).
				HasSubtitle("subtitle").
				HasProfile("").
				HasCreatedOnNotEmpty().
				HasUpdatedOnNotEmpty().
				HasNoPublishedOn().
				HasState(sdk.ListingStateDraft).
				HasNoReviewState().
				HasComment(comment).
				HasNoRegions().
				HasTargetAccounts("").
				HasIsMonetized(false).
				HasIsApplication(true).
				HasIsTargeted(true).
				HasIsLimitedTrial(false).
				HasIsByRequest(false).
				HasDistribution("EXTERNAL").
				HasIsMountlessQueryable(false).
				HasOrganizationProfileName("").
				HasNoUniformListingLocator().
				HasNoDetailedTargetAccounts(),
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
				HasTitle(basicManifestWithTargetTitle).
				HasState(sdk.ListingStateDraft).
				HasNoReviewState(),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithReview(true))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasTitle(basicManifestWithTargetTitle).
				HasState(sdk.ListingStateDraft).
				HasNoReviewState(),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithPublish(true))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasTitle(basicManifestWithTargetTitle).
				HasState(sdk.ListingStatePublished).
				HasNoReviewState(),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithUnpublish(true))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasTitle(basicManifestWithTargetTitle).
				HasState(sdk.ListingStateUnpublished).
				HasNoReviewState(),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithReview(true))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasTitle(basicManifestWithTargetTitle).
				HasState(sdk.ListingStateUnpublished).
				HasNoReviewState(),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithPublish(true))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, id).
				HasTitle(basicManifestWithTargetTitle).
				HasState(sdk.ListingStatePublished).
				HasNoReviewState(),
		)
	})

	t.Run("alter: change manifest with optional values", func(t *testing.T) {
		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		assertThatObject(t,
			objectassert.ListingFromObject(t, listing).
				HasTitle(basicManifestTitle).
				HasSubtitle("subtitle").
				HasNoComment(),
		)

		basicManifestWithDifferentSubtitle, title := testClientHelper().Listing.BasicManifestWithDifferentSubtitle(t)
		err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(listing.ID()).
			WithAlterListingAs(*sdk.NewAlterListingAsRequest(basicManifestWithDifferentSubtitle).
				WithPublish(false).
				WithReview(false).
				WithComment(comment),
			))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.Listing(t, listing.ID()).
				HasTitle(title).
				HasSubtitle("different_subtitle").
				HasComment(comment),
		)
	})

	t.Run("alter: add version", func(t *testing.T) {
		basicWithDifferentSubtitleManifest, title := testClientHelper().Listing.BasicManifest(t)
		testClientHelper().Stage.PutInLocationWithContent(t, stage.Location()+"/basic_different_subtitle", "manifest.yml", basicWithDifferentSubtitleManifest)
		basicManifestWithDifferentSubtitleStageLocation := sdk.NewStageLocation(stage.ID(), "basic_different_subtitle")

		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		assertThatObject(t,
			objectassert.ListingFromObject(t, listing).
				HasTitle(title).
				HasSubtitle("subtitle").
				HasNoComment(),
		)

		err := client.Listings.Alter(ctx, sdk.NewAlterListingRequest(listing.ID()).
			WithAddVersion(*sdk.NewAddListingVersionRequest(basicManifestWithDifferentSubtitleStageLocation).
				WithVersionName("v2").
				WithIfNotExists(true).
				WithComment(comment)))
		assert.NoError(t, err)

		assertThatObject(t,
			objectassert.ListingFromObject(t, listing).
				HasTitle(title).
				HasSubtitle("subtitle").
				HasNoComment(),
		)

		versions, err := client.Listings.ShowVersions(ctx, sdk.NewShowVersionsListingRequest(listing.ID()))
		assert.NoError(t, err)
		assert.Len(t, versions, 1)
		assert.NotEmpty(t, versions[0].CreatedOn)
		assert.NotEmpty(t, versions[0].Name)
		assert.Equal(t, "v2", *versions[0].Alias)
		assert.Empty(t, versions[0].LocationUrl)
		assert.True(t, versions[0].IsDefault)
		assert.False(t, versions[0].IsLive)
		assert.True(t, versions[0].IsFirst)
		assert.True(t, versions[0].IsLast)
		assert.Equal(t, comment, *versions[0].Comment)
		assert.Nil(t, versions[0].GitCommitHash)
	})

	t.Run("alter: add version with inlined manifest", func(t *testing.T) {
		versionStage, versionStageCleanup := testClientHelper().Stage.CreateStage(t)
		t.Cleanup(versionStageCleanup)

		_ = testClientHelper().Stage.PutInLocationWithContent(t, versionStage.Location(), "manifest.yml", basicManifest)
		manifestLocation := sdk.NewStageLocation(versionStage.ID(), "")

		listing, listingCleanup := testClientHelper().Listing.Create(t)
		t.Cleanup(listingCleanup)

		// Versions can be only added whenever listing was using manifest sourced from stage at any point
		_, err := client.Listings.ShowVersions(ctx, sdk.NewShowVersionsListingRequest(listing.ID()))
		assert.ErrorContains(t, err, "Attached stage not exists")

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(listing.ID()).WithAddVersion(*sdk.NewAddListingVersionRequest(manifestLocation).WithVersionName("v1")))
		assert.NoError(t, err)

		versions, err := client.Listings.ShowVersions(ctx, sdk.NewShowVersionsListingRequest(listing.ID()))
		assert.NoError(t, err)
		require.Len(t, versions, 1)
		require.NotNil(t, versions[0].Alias)
		assert.Equal(t, "v1", *versions[0].Alias)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(listing.ID()).WithAlterListingAs(*sdk.NewAlterListingAsRequest(basicManifest).WithReview(false).WithPublish(false)))
		assert.NoError(t, err)

		versions, err = client.Listings.ShowVersions(ctx, sdk.NewShowVersionsListingRequest(listing.ID()))
		assert.NoError(t, err)
		require.Len(t, versions, 2)

		inlineVersion, err := collections.FindFirst(versions, func(v sdk.ListingVersion) bool { return v.Name == "VERSION$2" })
		assert.NoError(t, err)
		require.Nil(t, inlineVersion.Alias)
		require.NotNil(t, inlineVersion.Comment)
		assert.Equal(t, "Inline update", *inlineVersion.Comment) // This is the default comment for inline updates

		// Removing referenced stage doesn't seem to break the listing's versioning
		versionStageCleanup()

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(listing.ID()).WithAlterListingAs(*sdk.NewAlterListingAsRequest(basicManifest).WithReview(false).WithPublish(false)))
		assert.NoError(t, err)

		versions, err = client.Listings.ShowVersions(ctx, sdk.NewShowVersionsListingRequest(listing.ID()))
		assert.NoError(t, err)
		require.Len(t, versions, 3)
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
			HasTitle(basicManifestTitle).
			HasComment(comment),
		)

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithSet(*sdk.NewListingSetRequest().WithComment(newComment)))
		assert.NoError(t, err)

		assertThatObject(t, objectassert.Listing(t, id).
			HasName(id.Name()).
			HasTitle(basicManifestTitle).
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
		accountId := testClientHelper().Context.CurrentAccountId(t)
		manifest, title := testClientHelper().Listing.BasicManifestWithTargetAccounts(t, accountId)

		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithAs(manifest).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithIfNotExists(true).
			WithPublish(false).
			WithReview(false))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		listingDetails, err := client.Listings.Describe(ctx, sdk.NewDescribeListingRequest(id))
		require.NoError(t, err)
		require.NotNil(t, listingDetails)

		assert.NotEmpty(t, listingDetails.GlobalName)
		assert.Equal(t, id.Name(), listingDetails.Name)
		assert.NotEmpty(t, listingDetails.Owner)
		assert.NotEmpty(t, listingDetails.OwnerRoleType)
		assert.NotEmpty(t, listingDetails.CreatedOn)
		assert.NotEmpty(t, listingDetails.UpdatedOn)
		assert.Nil(t, listingDetails.PublishedOn)
		assert.Equal(t, title, listingDetails.Title)
		assert.Equal(t, "subtitle", *listingDetails.Subtitle)
		assert.Equal(t, "description", *listingDetails.Description)
		assert.JSONEq(t, `{
"type" : "OFFLINE"
}`, *listingDetails.ListingTerms)
		assert.Equal(t, sdk.ListingStateDraft, listingDetails.State)
		assert.Equal(t, share.ID().Name(), listingDetails.Share.Name())
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
		assert.Equal(t, fmt.Sprintf("%s.%s", accountId.OrganizationName(), accountId.AccountName()), *listingDetails.TargetAccounts)
		assert.Nil(t, listingDetails.Regions)
		assert.Nil(t, listingDetails.RefreshSchedule)
		assert.Nil(t, listingDetails.RefreshType)
		assert.Nil(t, listingDetails.ReviewState)
		assert.Nil(t, listingDetails.RejectionReason)
		assert.Nil(t, listingDetails.UnpublishedByAdminReasons)
		assert.False(t, listingDetails.IsMonetized)
		assert.False(t, listingDetails.IsApplication)
		assert.True(t, listingDetails.IsTargeted)
		assert.False(t, *listingDetails.IsLimitedTrial)
		assert.False(t, *listingDetails.IsByRequest)
		assert.Nil(t, listingDetails.LimitedTrialPlan)
		assert.Nil(t, listingDetails.RetriedOn)
		assert.Nil(t, listingDetails.ScheduledDropTime)
		assert.NotEmpty(t, listingDetails.ManifestYaml)
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
		assert.True(t, *listingDetails.IsShare)
		assert.Nil(t, listingDetails.RequestApprovalType)
		assert.Empty(t, *listingDetails.MonetizationDisplayOrder)
		assert.NotNil(t, listingDetails.LegacyUniformListingLocators)
		assert.Equal(t, "", *listingDetails.LegacyUniformListingLocators)
	})

	// TODO(SNOW-2220593): Test describe with revisions
	// t.Run("describe: revisions", func(t *testing.T) {
	// })

	t.Run("drop safely", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id).
			WithFrom(basicManifestWithTargetStageLocation).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithIfNotExists(true).
			WithPublish(true).
			WithReview(true).
			WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(t, id))

		err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id))
		assert.Error(t, err)

		err = client.Listings.DropSafely(ctx, id)
		assert.NoError(t, err)

		// Show that it can be called for already dropped listings
		err = client.Listings.DropSafely(ctx, id)
		assert.NoError(t, err)
	})
}
