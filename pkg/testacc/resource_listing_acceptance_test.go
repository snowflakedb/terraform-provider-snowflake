//go:build !account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	testifyassert "github.com/stretchr/testify/assert"
)

func TestAcc_Listing_Basic_Inlined(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	basicManifest, listingTitle := testClient().Listing.BasicManifestWithUnquotedValues(t)
	manifestWithTargetAccounts, listingTitleWithTargetAccounts := testClient().Listing.BasicManifestWithUnquotedValuesAndTargetAccount(t, testClient().Context.CurrentAccountId(t))

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

	modelUnset := model.ListingWithInlineManifest("test", id.Name(), manifestWithTargetAccounts).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanFalse)

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
						HasManifestFromStringNotEmpty().
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
						HasManifestFromStringNotEmpty().
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
						HasManifestFromStringNotEmpty().
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
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelUnset.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelUnset),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelUnset.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ListingShowOutput(t, modelUnset.ResourceReference()).
						HasName(id.Name()).
						HasTitle(listingTitleWithTargetAccounts).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStateUnpublished).
						HasComment(""),
				),
			},
		},
	})
}

func TestAcc_Listing_Basic_FromStage(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	basicManifest, listingTitle := testClient().Listing.BasicManifest(t)
	manifestWithTargetAccounts, listingTitleWithTargetAccounts := testClient().Listing.BasicManifestWithUnquotedValuesAndTargetAccount(t, testClient().Context.CurrentAccountId(t))

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

	manifest, title := testClient().Listing.BasicManifestWithUnquotedValuesAndTargetAccount(t, testClient().Context.CurrentAccountId(t))

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

func TestAcc_Listing_Complete_FromStage(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	manifest, title := testClient().Listing.BasicManifestWithUnquotedValuesAndTargetAccount(t, testClient().Context.CurrentAccountId(t))

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

func TestAcc_Listing_NewVersions_Inlined(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	manifest1, title1 := testClient().Listing.BasicManifestWithUnquotedValuesAndTargetAccount(t, testClient().Context.CurrentAccountId(t))
	manifest2, title2 := testClient().Listing.BasicManifestWithUnquotedValuesAndTargetAccountAndDifferentSubtitle(t, testClient().Context.CurrentAccountId(t))

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	_ = testClient().Stage.PutInLocationWithContent(t, stage.Location()+"/manifest", "manifest.yml", manifest2)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	modelInitialInlined := model.ListingWithInlineManifest("test", id.Name(), manifest1).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue)

	modelModifiedInlined := model.ListingWithInlineManifest("test", id.Name(), manifest2).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue)

	modelStaged := model.ListingWithStagedManifestWithLocation("test", id.Name(), stage.ID(), "manifest").
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
			// Create listing with inlined manifest (inline manifests don't track versions)
			{
				Config: accconfig.FromModels(t, modelInitialInlined),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelInitialInlined.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue),
					resourceshowoutputassert.ListingShowOutput(t, modelInitialInlined.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title1).
						HasState(sdk.ListingStatePublished),
					assert.Check(assertListingIsNotVersioned(t, id)),
				),
			},
			// Modify the manifest and show that no version was produced
			{
				Config: accconfig.FromModels(t, modelModifiedInlined),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelModifiedInlined.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue),
					resourceshowoutputassert.ListingShowOutput(t, modelModifiedInlined.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title2).
						HasState(sdk.ListingStatePublished),
					assert.Check(assertListingIsNotVersioned(t, id)),
				),
			},
			// Change the manifest source from inlined to staged (the manifest is the same; a new version is produced)
			{
				Config: accconfig.FromModels(t, modelStaged),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelStaged.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageStageId(stage.ID()).
						HasManifestFromStageLocation("manifest").
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue),
					resourceshowoutputassert.ListingShowOutput(t, modelStaged.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title2).
						HasState(sdk.ListingStatePublished),
					assert.Check(assertContainsListingVersion(t, id, "VERSION$1", "null")),
				),
			},
			// Switch back to the inlined manifest (the manifest is the same; a new version is produced).
			// After listing being sourced from stage once, it will be recording every version from now on.
			// Now, it doesn't matter if it's sourced from stage or back to inlined.
			{
				Config: accconfig.FromModels(t, modelModifiedInlined),
				Check: assertThat(t,
					resourceassert.ListingResource(t, modelModifiedInlined.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue),
					resourceshowoutputassert.ListingShowOutput(t, modelModifiedInlined.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title2).
						HasState(sdk.ListingStatePublished),
					assert.Check(assertContainsListingVersion(t, id, "VERSION$2", "")),
				),
			},
		},
	})
}

func TestAcc_Listing_NewVersions_FromStage(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	manifest1, title1 := testClient().Listing.BasicManifestWithUnquotedValuesAndTargetAccount(t, testClient().Context.CurrentAccountId(t))
	manifest2, title2 := testClient().Listing.BasicManifestWithUnquotedValuesAndTargetAccountAndDifferentSubtitle(t, testClient().Context.CurrentAccountId(t))

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
					assert.Check(assertDoesNotContainListingVersion(t, id, "VERSION$6", "")),
				),
			},
		},
	})
}

func TestAcc_Listing_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	manifest, _ := testClient().Listing.BasicManifestWithUnquotedValues(t)

	listingModelWithoutManifest := func(resourceName string, name string) *model.ListingModel {
		l := &model.ListingModel{ResourceModelMeta: accconfig.Meta(resourceName, resources.Listing)}
		l.WithName(name)
		return l
	}

	modelWithBothShareAndApplicationPackage := model.ListingWithInlineManifest("test", id.Name(), manifest).
		WithShare("test_share").
		WithApplicationPackage("test_app_package")

	modelWithInvalidStageId := listingModelWithoutManifest("test", id.Name()).
		WithManifestValue(tfconfig.ListVariable(
			tfconfig.MapVariable(map[string]tfconfig.Variable{
				"from_stage": tfconfig.ListVariable(
					tfconfig.MapVariable(map[string]tfconfig.Variable{
						"stage": tfconfig.StringVariable("invalid.stage.identifier.name"),
					}),
				),
			}),
		))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Listing),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelWithBothShareAndApplicationPackage),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"application_package": conflicts with share`),
			},
			{
				Config:      accconfig.FromModels(t, modelWithBothShareAndApplicationPackage),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"share": conflicts with application_package`),
			},
			{
				Config:      accconfig.FromModels(t, modelWithInvalidStageId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Expected SchemaObjectIdentifier identifier type`),
			},
		},
	})
}

func assertContainsListingVersion(t *testing.T, id sdk.AccountObjectIdentifier, expectedName string, expectedAlias string) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		versions, err := testClient().Listing.ShowVersions(t, id)
		testifyassert.NoError(t, err)

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
	t.Helper()
	return func(s *terraform.State) error {
		if err := assertContainsListingVersion(t, id, expectedName, expectedAlias)(s); err == nil {
			return fmt.Errorf("expected version name '%s' with alias '%s' to not be present, but was found", expectedName, expectedAlias)
		}
		return nil
	}
}

func assertListingIsNotVersioned(t *testing.T, id sdk.AccountObjectIdentifier) resource.TestCheckFunc {
	t.Helper()
	return func(s *terraform.State) error {
		_, err := testClient().Listing.ShowVersions(t, id)
		if !strings.Contains(err.Error(), "Attached stage not exists") {
			return err
		}
		return nil
	}
}
