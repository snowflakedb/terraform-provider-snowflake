//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Listings_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "1")
	idTwo := testClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix + "2")
	idThree := testClient().Ids.RandomAccountObjectIdentifier()

	manifest, _ := testClient().Listing.BasicManifestWithUnquotedValues(t)

	// Listings without share/app package cannot be published
	listingModel1 := model.ListingWithInlineManifest("listing_1", idOne.Name(), manifest).WithPublish("false")
	listingModel2 := model.ListingWithInlineManifest("listing_2", idTwo.Name(), manifest).WithPublish("false")
	listingModel3 := model.ListingWithInlineManifest("listing_3", idThree.Name(), manifest).WithPublish("false")

	listingsLike := datasourcemodel.Listings("like").
		WithLike(idOne.Name()).
		WithDependsOn(listingModel1.ResourceReference(), listingModel2.ResourceReference(), listingModel3.ResourceReference())

	listingsStartsWith := datasourcemodel.Listings("test").
		WithStartsWith(prefix).
		WithDependsOn(listingModel1.ResourceReference(), listingModel2.ResourceReference(), listingModel3.ResourceReference())

	listingsLimit := datasourcemodel.Listings("test").
		WithRowsAndFrom(1, prefix).
		WithDependsOn(listingModel1.ResourceReference(), listingModel2.ResourceReference(), listingModel3.ResourceReference())

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.ListingsDatasource))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, listingModel1, listingModel2, listingModel3, listingsLike),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(listingsLike.DatasourceReference(), "listings.#", "1"),
					resource.TestCheckResourceAttr(listingsLike.DatasourceReference(), "listings.0.show_output.#", "1"),
					resource.TestCheckResourceAttr(listingsLike.DatasourceReference(), "listings.0.show_output.0.name", idOne.Name()),
				),
			},
			{
				Config: config.FromModels(t, providerModel, listingModel1, listingModel2, listingModel3, listingsStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(listingsStartsWith.DatasourceReference(), "listings.#", "2"),
				),
			},
			{
				Config: config.FromModels(t, providerModel, listingModel1, listingModel2, listingModel3, listingsLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(listingsLimit.DatasourceReference(), "listings.#", "1"),
				),
			},
		},
	})
}

func TestAcc_Listings_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	manifest, title := testClient().Listing.BasicManifestWithUnquotedValuesAndTargetAccounts(t)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	listingModel := model.ListingWithInlineManifest("listing_complete", id.Name(), manifest).
		WithShare(share.ID().Name()).
		WithPublish("true").
		WithComment(comment)

	listingsModelWithoutAdditional := datasourcemodel.Listings("without_additional").
		WithLike(id.Name()).
		WithStartsWith(id.Name()).
		WithLimit(1).
		WithWithDescribe(false).
		WithDependsOn(listingModel.ResourceReference())

	listingsModel := datasourcemodel.Listings("with_additional").
		WithLike(id.Name()).
		WithStartsWith(id.Name()).
		WithLimit(1).
		WithDependsOn(listingModel.ResourceReference())

	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.ListingsDatasource))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Listing),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, providerModel, listingModel, listingsModelWithoutAdditional),
				Check: assertThat(t,
					resourceshowoutputassert.ListingsDatasourceShowOutput(t, listingsModelWithoutAdditional.DatasourceReference()).
						HasName(id.Name()).
						HasTitle(title).
						HasState(sdk.ListingStatePublished).
						HasComment(comment),

					assert.Check(resource.TestCheckResourceAttr(listingsModelWithoutAdditional.DatasourceReference(), "listings.0.describe_output.#", "0")),
				),
			},
			{
				Config: config.FromModels(t, providerModel, listingModel, listingsModel),
				Check: assertThat(t,
					resourceshowoutputassert.ListingsDatasourceShowOutput(t, listingsModel.DatasourceReference()).
						HasName(id.Name()).
						HasTitle(title).
						HasState(sdk.ListingStatePublished).
						HasComment(comment),

					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.title", title)),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.description", "description")),
					assert.Check(resource.TestCheckResourceAttrSet(listingsModel.DatasourceReference(), "listings.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttrSet(listingsModel.DatasourceReference(), "listings.0.describe_output.0.updated_on")),
					assert.Check(resource.TestCheckResourceAttrSet(listingsModel.DatasourceReference(), "listings.0.describe_output.0.owner")),
					assert.Check(resource.TestCheckResourceAttrSet(listingsModel.DatasourceReference(), "listings.0.describe_output.0.owner_role_type")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.state", "PUBLISHED")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.share", share.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.application_package", "")),
					assert.Check(resource.TestCheckResourceAttrSet(listingsModel.DatasourceReference(), "listings.0.describe_output.0.manifest_yaml")),
					assert.Check(resource.TestCheckResourceAttrSet(listingsModel.DatasourceReference(), "listings.0.describe_output.0.listing_terms")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.business_needs", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.usage_examples", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.data_attributes", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.categories", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.resources", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.profile", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.customized_contact_info", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.data_dictionary", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.data_preview", "")),
					assert.Check(resource.TestCheckResourceAttrSet(listingsModel.DatasourceReference(), "listings.0.describe_output.0.published_on")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.revisions", "PUBLISHED")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.target_accounts", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.regions", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.refresh_schedule", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.refresh_type", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.review_state", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.rejection_reason", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.unpublished_by_admin_reasons", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.is_monetized", "false")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.is_application", "false")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.is_targeted", "true")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.is_limited_trial", "false")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.is_by_request", "false")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.limited_trial_plan", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.retried_on", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.scheduled_drop_time", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.distribution", "EXTERNAL")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.is_mountless_queryable", "false")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.organization_profile_name", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.uniform_listing_locator", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.approver_contact", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.support_contact", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.published_version_name", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.published_version_alias", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.is_share", "true")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.request_approval_type", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.monetization_display_order", "")),
					assert.Check(resource.TestCheckResourceAttr(listingsModel.DatasourceReference(), "listings.0.describe_output.0.legacy_uniform_listing_locators", "")),
				),
			},
		},
	})
}
