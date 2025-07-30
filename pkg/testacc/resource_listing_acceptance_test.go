//go:build !account_level_tests

package testacc

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"testing"
)

func TestAcc_Listing_Basic_Inlined(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	basicManifest, listingTitle := testClient().Listing.BasicManifest(t)
	listingTitleEscaped := fmt.Sprintf(`\"%s\"`, listingTitle)

	manifestWithTargetAccounts, listingTitleWithTargetAccounts := testClient().Listing.BasicManifestWithTargetAccount(t, testClient().Context.CurrentAccountId(t))
	listingTitleWithTargetAccountsEscaped := fmt.Sprintf(`\"%s\"`, listingTitleWithTargetAccounts)

	comment, newComment := random.Comment(), random.Comment()

	modelBasic := model.ListingWithInlineManifest("test", id.Name(), basicManifest).
		// Has to be set when a listing is not associated with a share or an application package
		WithPublish(r.BooleanFalse)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	modelComplete := model.ListingWithInlineManifest("test", id.Name(), manifestWithTargetAccounts).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue).
		WithComment(comment)

	modelCompleteWithDifferentValues := model.ListingWithInlineManifest("test", id.Name(), manifestWithTargetAccounts).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanFalse).
		WithComment(newComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Listing),
		Steps: []resource.TestStep{
			// create without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareEmpty().
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitleEscaped).
						HasSubtitle(`\"subtitle\"`).
						HasState(sdk.ListingStateDraft).
						HasComment(""),
				),
			},
			// import without optionals
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedListingResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasNoManifest().
						HasShareEmpty().
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasTitle(listingTitleEscaped).
						HasSubtitle(`\"subtitle\"`).
						HasState(sdk.ListingStateDraft).
						HasComment(""),
				),
			},
			// set optionals (expect re-creation as share is set)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ListingShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitleWithTargetAccountsEscaped).
						HasSubtitle(`\"subtitle\"`).
						HasState(sdk.ListingStatePublished).
						HasComment(comment),
				),
			},
			// import complete object
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedListingResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasNoManifest().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasTitle(listingTitleWithTargetAccountsEscaped).
						HasSubtitle(`\"subtitle\"`).
						HasState(sdk.ListingStatePublished).
						HasComment(comment),
				),
			},
			// alter
			{
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentString(newComment),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitleWithTargetAccountsEscaped).
						HasSubtitle(`\"subtitle\"`).
						HasState(sdk.ListingStateUnpublished).
						HasComment(newComment),
				),
			},
			// change externally
			{
				PreConfig: func() {
					testClient().Listing.Alter(t, sdk.NewAlterListingRequest(id).WithPublish(true))
					testClient().Listing.Alter(t, sdk.NewAlterListingRequest(id).WithSet(*sdk.NewListingSetRequest().WithComment(comment)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(modelCompleteWithDifferentValues.ResourceReference(), "publish", sdk.String(r.BooleanFalse), sdk.String(r.BooleanTrue)),
						planchecks.ExpectDrift(modelCompleteWithDifferentValues.ResourceReference(), "comment", sdk.String(newComment), sdk.String(comment)),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentString(newComment),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitleWithTargetAccountsEscaped).
						HasSubtitle(`\"subtitle\"`).
						HasState(sdk.ListingStateUnpublished).
						HasComment(newComment),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareEmpty().
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitleEscaped).
						HasSubtitle(`\"subtitle\"`).
						HasState(sdk.ListingStateDraft).
						HasComment(""),
				),
			},
		},
	})
}

func TestAcc_Listing_Basic_FromStage(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	basicManifest, listingTitle := testClient().Listing.BasicManifest(t)
	manifestWithTargetAccounts, listingTitleWithTargetAccounts := testClient().Listing.BasicManifestWithTargetAccount(t, testClient().Context.CurrentAccountId(t))

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	_ = testClient().Stage.PutInLocationWithContent(t, stage.Location()+"/basic", "manifest.yml", basicManifest)
	_ = testClient().Stage.PutInLocationWithContent(t, stage.Location()+"/with_targets", "manifest.yml", manifestWithTargetAccounts)

	comment, newComment := random.Comment(), random.Comment()

	modelBasic := model.ListingWithStagedManifestWithLocation("test", id.Name(), stage.ID(), "basic").
		// Has to be set when a listing is not associated with a share or an application package
		WithPublish(r.BooleanFalse)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)

	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	modelComplete := model.ListingWithStagedManifestWithLocation("test", id.Name(), stage.ID(), "with_targets").
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue).
		WithComment(comment)

	modelCompleteWithDifferentValues := model.ListingWithStagedManifestWithLocation("test", id.Name(), stage.ID(), "with_targets").
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanFalse).
		WithComment(newComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Listing),
		Steps: []resource.TestStep{
			// create without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageNotEmpty().
						HasShareEmpty().
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitle).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStateDraft).
						HasComment(""),
				),
			},
			// import without optionals
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedListingResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasNoManifest().
						HasShareEmpty().
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasTitle(listingTitle).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStateDraft).
						HasComment(""),
				),
			},
			// set optionals (expect re-creation as share is set)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ListingShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitleWithTargetAccounts).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStatePublished).
						HasComment(comment),
				),
			},
			// import complete object
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedListingResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasNoManifest().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasTitle(listingTitleWithTargetAccounts).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStatePublished).
						HasComment(comment),
				),
			},
			// alter
			{
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentString(newComment),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitleWithTargetAccounts).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStateUnpublished).
						HasComment(newComment),
				),
			},
			// change externally
			{
				PreConfig: func() {
					testClient().Listing.Alter(t, sdk.NewAlterListingRequest(id).WithPublish(true))
					testClient().Listing.Alter(t, sdk.NewAlterListingRequest(id).WithSet(*sdk.NewListingSetRequest().WithComment(comment)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(modelCompleteWithDifferentValues.ResourceReference(), "publish", sdk.String(r.BooleanFalse), sdk.String(r.BooleanTrue)),
						planchecks.ExpectDrift(modelCompleteWithDifferentValues.ResourceReference(), "comment", sdk.String(newComment), sdk.String(comment)),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentString(newComment),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitleWithTargetAccounts).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStateUnpublished).
						HasComment(newComment),
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageNotEmpty().
						HasShareEmpty().
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitle).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStateDraft).
						HasComment(""),
				),
			},
		},
	})
}

func TestAcc_Listing_Complete_Inlined(t *testing.T)   {}
func TestAcc_Listing_Complete_FromStage(t *testing.T) {}

func TestAcc_Listing_NewVersions_Inlined(t *testing.T)   {}
func TestAcc_Listing_NewVersions_FromStage(t *testing.T) {}

func TestAcc_Listing_Updates_Inlined(t *testing.T)   {}
func TestAcc_Listing_Updates_FromStage(t *testing.T) {}

func TestAcc_Listing_UpdateManifestSource(t *testing.T) {}
func TestAcc_Listing_Validations(t *testing.T)          {}
