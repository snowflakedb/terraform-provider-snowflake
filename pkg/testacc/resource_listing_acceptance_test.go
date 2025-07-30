//go:build !account_level_tests

package testacc

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"slices"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
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
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
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

func TestAcc_Listing_Complete_Inlined(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	manifest, title := testClient().Listing.BasicManifestWithTargetAccount(t, testClient().Context.CurrentAccountId(t))
	titleEscaped := fmt.Sprintf(`\"%s\"`, title)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	applicationPackage, applicationPackageCleanup := testClient().ApplicationPackage.CreateApplicationPackageWithReleaseChannelsDisabled(t)
	t.Cleanup(applicationPackageCleanup)

	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "manifest.yml", "")
	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "setup.sql", "CREATE APPLICATION ROLE IF NOT EXISTS APP_HELLO_SNOWFLAKE;")

	version := "v1"
	testClient().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stage.ID(), version)
	testClient().ApplicationPackage.SetDefaultReleaseDirective(t, applicationPackage.ID(), version)

	modelComplete := model.ListingWithInlineManifest("test", id.Name(), manifest).
		WithApplicationPackage(applicationPackage.ID().Name()).
		WithPublish(r.BooleanTrue).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Listing),
		Steps: []resource.TestStep{
			// create complete with all optionals
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareEmpty().
						HasApplicationPackageString(applicationPackage.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ListingShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasTitle(titleEscaped).
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
						HasShareEmpty().
						HasApplicationPackageString(applicationPackage.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasTitle(titleEscaped).
						HasSubtitle(`\"subtitle\"`).
						HasState(sdk.ListingStatePublished).
						HasComment(comment),
				),
			},
		},
	})
}

func TestAcc_Listing_Complete_FromStage(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	manifest, title := testClient().Listing.BasicManifestWithTargetAccount(t, testClient().Context.CurrentAccountId(t))

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	_ = testClient().Stage.PutInLocationWithContent(t, stage.Location()+"/listing", "manifest.yml", manifest)

	applicationPackage, applicationPackageCleanup := testClient().ApplicationPackage.CreateApplicationPackageWithReleaseChannelsDisabled(t)
	t.Cleanup(applicationPackageCleanup)

	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "manifest.yml", "")
	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "setup.sql", "CREATE APPLICATION ROLE IF NOT EXISTS APP_HELLO_SNOWFLAKE;")

	version := "v1"
	testClient().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stage.ID(), version)
	testClient().ApplicationPackage.SetDefaultReleaseDirective(t, applicationPackage.ID(), version)

	modelComplete := model.ListingWithStagedManifestWithOptionals("test", id.Name(), stage.ID(), "v0", "", "listing").
		WithApplicationPackage(applicationPackage.ID().Name()).
		WithPublish(r.BooleanTrue).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Listing),
		Steps: []resource.TestStep{
			// create complete with all optionals
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageNotEmpty().
						HasShareEmpty().
						HasApplicationPackageString(applicationPackage.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ListingShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title).
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
						HasShareEmpty().
						HasApplicationPackageString(applicationPackage.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasTitle(title).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStatePublished).
						HasComment(comment),
				),
			},
		},
	})
}

func TestAcc_Listing_NewVersions_Inlined(t *testing.T) {}

func TestAcc_Listing_NewVersions_FromStage(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	manifest1, title1 := testClient().Listing.BasicManifestWithTargetAccount(t, testClient().Context.CurrentAccountId(t))
	manifest2, title2 := testClient().Listing.BasicManifestWithTargetAccountAndDifferentSubtitle(t, testClient().Context.CurrentAccountId(t))

	stage1, stage1Cleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stage1Cleanup)
	_ = testClient().Stage.PutInLocationWithContent(t, stage1.Location()+"/v1", "manifest.yml", manifest1)
	_ = testClient().Stage.PutInLocationWithContent(t, stage1.Location()+"/v2", "manifest.yml", manifest2)

	stage2, stage2Cleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stage2Cleanup)
	_ = testClient().Stage.PutInLocationWithContent(t, stage2.Location()+"/v2", "manifest.yml", manifest2)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	modelInitial := model.ListingWithStagedManifestWithLocation("test", id.Name(), stage1.ID(), "v1").
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue)

	modelWithNewManifestLocation := model.ListingWithStagedManifestWithLocation("test", id.Name(), stage1.ID(), "v2").
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue)

	versionComment := random.Comment()
	modelWithVersionNameAndComment := model.ListingWithStagedManifestWithOptionals("test", id.Name(), stage1.ID(), "version_name", versionComment, "v2").
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue)

	modelWithNewVersionName := model.ListingWithStagedManifestWithOptionals("test", id.Name(), stage1.ID(), "other_version_name", versionComment, "v2").
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue)

	modelWithNewStage := model.ListingWithStagedManifestWithLocation("test", id.Name(), stage2.ID(), "v2").
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Listing),
		Steps: []resource.TestStep{
			// create initial listing with staged manifest
			{
				Config: accconfig.FromModels(t, modelInitial),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelInitial.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageStageId(stage1.ID()).
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue).
						HasCommentEmpty(),
					resourceshowoutputassert.ListingShowOutput(t, modelInitial.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title1).
						HasState(sdk.ListingStatePublished),
					// CREATE LISTING does not support specifying version name, so it's always "VERSION$1" with no alias
					assert.Check(assertContainsListingVersion(t, id, "VERSION$1", "null")),
				),
			},
			// Change manifest location (points to a different manifest, but it shouldn't matter) - should create a new version
			{
				Config: accconfig.FromModels(t, modelWithNewManifestLocation),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelWithNewManifestLocation.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageStageId(stage1.ID()).
						HasManifestFromStageVersionName("").
						HasManifestFromStageVersionComment("").
						HasManifestFromStageLocation("v2"),
					resourceshowoutputassert.ListingShowOutput(t, modelWithNewManifestLocation.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title2).
						HasState(sdk.ListingStatePublished),
					assert.Check(assertContainsListingVersion(t, id, "VERSION$2", "")),
				),
			},
			// add optional values (version name and version comment) - should create a new version
			{
				Config: accconfig.FromModels(t, modelWithVersionNameAndComment),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelWithVersionNameAndComment.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageStageId(stage1.ID()).
						HasManifestFromStageVersionName("version_name").
						HasManifestFromStageVersionComment(versionComment).
						HasManifestFromStageLocation("v2"),
					resourceshowoutputassert.ListingShowOutput(t, modelWithVersionNameAndComment.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title2).
						HasState(sdk.ListingStatePublished),
					assert.Check(assertContainsListingVersion(t, id, "VERSION$3", "version_name")),
				),
			},
			// change version_name - should create a new version
			{
				Config: accconfig.FromModels(t, modelWithNewVersionName),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelWithNewVersionName.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageStageId(stage1.ID()).
						HasManifestFromStageVersionName("other_version_name").
						HasManifestFromStageVersionComment(versionComment).
						HasManifestFromStageLocation("v2"),
					resourceshowoutputassert.ListingShowOutput(t, modelWithNewVersionName.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title2).
						HasState(sdk.ListingStatePublished),
					assert.Check(assertContainsListingVersion(t, id, "VERSION$4", "other_version_name")),
				),
			},
			// change stage and location - should create a new version
			{
				Config: accconfig.FromModels(t, modelWithNewStage),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelWithNewStage.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageStageId(stage2.ID()).
						HasManifestFromStageVersionName("").
						HasManifestFromStageVersionComment("").
						HasManifestFromStageLocation("v2"),
					resourceshowoutputassert.ListingShowOutput(t, modelWithNewStage.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title2).
						HasState(sdk.ListingStatePublished),
					assert.Check(assertContainsListingVersion(t, id, "VERSION$5", "")),
				),
			},
			// manifest changed externally, but the stage, version name, and location are the same - should create a new version, but no planned changes should be produced
			{
				PreConfig: func() {
					_ = testClient().Stage.PutInLocationWithContent(t, stage2.Location()+"/v2", "manifest.yml", manifest1)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithNewStage.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Config: accconfig.FromModels(t, modelWithNewStage),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelWithNewStage.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageStageId(stage2.ID()).
						HasManifestFromStageVersionName("").
						HasManifestFromStageVersionComment("").
						HasManifestFromStageLocation("v2"),
					resourceshowoutputassert.ListingShowOutput(t, modelWithNewStage.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title2).
						HasState(sdk.ListingStatePublished),
					assert.Check(assertContainsListingVersion(t, id, "VERSION$6", "")),
				),
			},
		},
	})
}

// TODO: Updates are tested above, maybe is there some situation that should be showed like external changes (e.g. to the manifest or changes to the publish state)?
func TestAcc_Listing_Updates_Inlined(t *testing.T)   {}
func TestAcc_Listing_Updates_FromStage(t *testing.T) {}

func TestAcc_Listing_UpdateManifestSource(t *testing.T) {}

func TestAcc_Listing_Validations(t *testing.T) {
	// validations
	// cannot set from_string and from_stage at the same time
	// cannot set share and application_package at the same time
}

func assertContainsListingVersion(t *testing.T, id sdk.AccountObjectIdentifier, expectedName string, expectedAlias string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		versions := testClient().Listing.ShowVersions(t, id)

		versionNamesAndAliases := collections.Map(versions, func(v sdk.ListingVersion) string {
			alias := "null"
			if v.Alias != nil {
				alias = *v.Alias
			}
			return fmt.Sprintf("%s_%s", v.Name, alias)
		})
		expectedNameWithAlias := fmt.Sprintf("%s_%s", expectedName, expectedAlias)

		if !slices.Contains(versionNamesAndAliases, expectedNameWithAlias) {
			return fmt.Errorf("expected version name '%s' with alias '%s' to be present, but was not found", expectedName, expectedAlias)
		}
		return nil
	}
}

func assertDoesNotContainListingVersion(t *testing.T, id sdk.AccountObjectIdentifier, expectedName string, expectedAlias string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := assertContainsListingVersion(t, id, expectedName, expectedAlias)(s); err == nil {
			return fmt.Errorf("expected version name '%s' with alias '%s' to not be present, but was found", expectedName, expectedAlias)
		}
		return nil
	}
}
