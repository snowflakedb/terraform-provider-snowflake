//go:build !account_level_tests

package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInt_Listings(t *testing.T) {
	stage, stageCleanup := testClientHelper().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	stageLocationPath := "dir/subdir"
	stageLocation := sdk.NewStageLocation(stage.ID(), stageLocationPath)

	share, shareCleanup := testClientHelper().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	applicationPackage, applicationPackageCleanup := testClientHelper().ApplicationPackage.CreateApplicationPackage(t)
	t.Cleanup(applicationPackageCleanup)

	client := testClient(t)
	ctx := testContext(t)

	manifest := testClientHelper().Listing.SampleListingManifest(t)

	t.Run("create from manifest: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest))
		assert.NoError(t, err)

		// TODO: Assert

		err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	})

	t.Run("create from manifest: complete with stage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithIfNotExists(true).
			WithPublish(true).
			WithReview(true).
			WithComment(comment))
		assert.NoError(t, err)

		// TODO: Assert

		err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	})

	t.Run("create from manifest: complete with application packages", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest).
			WithWith(*sdk.NewListingWithRequest().WithApplicationPackage(applicationPackage.ID())).
			WithIfNotExists(true).
			WithPublish(true).
			WithReview(true).
			WithComment(comment))
		assert.NoError(t, err)

		// TODO: Assert

		err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	})

	t.Run("create from stage: no optionals", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.CreateFromStage(ctx, sdk.NewCreateFromStageListingRequest(id, stageLocation))
		assert.NoError(t, err)

		// TODO: Assert

		err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	})

	t.Run("create from stage: complete with stage", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.CreateFromStage(ctx, sdk.NewCreateFromStageListingRequest(id, stageLocation).
			WithWith(*sdk.NewListingWithRequest().WithShare(share.ID())).
			WithIfNotExists(true).
			WithPublish(true).
			WithReview(true))
		assert.NoError(t, err)

		// TODO: Assert

		err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	})

	t.Run("create from stage: complete with application packages", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.CreateFromStage(ctx, sdk.NewCreateFromStageListingRequest(id, stageLocation).
			WithWith(*sdk.NewListingWithRequest().WithApplicationPackage(applicationPackage.ID())).
			WithIfNotExists(true).
			WithPublish(true).
			WithReview(true))
		assert.NoError(t, err)

		// TODO: Assert

		err = client.Listings.Drop(ctx, sdk.NewDropListingRequest(id).WithIfExists(true))
		assert.NoError(t, err)
	})

	t.Run("alter: change state", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(id))

		// TODO: Assert

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
			WithIfExists(true).
			WithReview(true))
		assert.NoError(t, err)

		// TODO: Assert

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
			WithIfExists(true).
			WithPublish(true))
		assert.NoError(t, err)

		// TODO: Assert

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
			WithIfExists(true).
			WithUnpublish(true))
		assert.NoError(t, err)
	})

	t.Run("alter: change manifest with optional values", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(id))

		// TODO: Assert

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
			WithAlterListingAs(*sdk.NewAlterListingAsRequest(manifest).
				WithPublish(true).
				WithReview(true).
				WithComment(comment),
			))
		assert.NoError(t, err)

		// TODO: Assert
	})

	t.Run("alter: add version", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		versionName := random.AlphaN(10)

		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(id))

		// TODO: Assert

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).
			WithAddVersion(*sdk.NewAddListingVersionRequest(versionName, stageLocation).
				WithIfNotExists(true).
				WithComment(comment)))
		assert.NoError(t, err)

		// TODO: Assert
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		newId := testClientHelper().Ids.RandomAccountObjectIdentifier()

		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(id))

		// TODO: Assert

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithRenameTo(newId))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(newId))

		// TODO: Assert
	})

	t.Run("alter: set", func(t *testing.T) {
		id := testClientHelper().Ids.RandomAccountObjectIdentifier()
		comment := random.Comment()
		newComment := random.Comment()

		err := client.Listings.Create(ctx, sdk.NewCreateListingRequest(id, manifest).WithComment(comment))
		assert.NoError(t, err)
		t.Cleanup(testClientHelper().Listing.DropFunc(id))

		// TODO: Assert

		err = client.Listings.Alter(ctx, sdk.NewAlterListingRequest(id).WithSet(*sdk.NewListingSetRequest().WithComment(newComment)))
		assert.NoError(t, err)

		// TODO: Assert
	})

	t.Run("drop: existing", func(t *testing.T) {

	})

	t.Run("show: default", func(t *testing.T) {

	})

	t.Run("show: with options", func(t *testing.T) {

	})

	t.Run("describe: default", func(t *testing.T) {

	})
}
