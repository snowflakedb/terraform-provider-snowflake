//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OrganizationListing_BasicUseCase_Inlined(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	targetAccount := testClient().Context.CurrentAccountId(t)
	targetAccounts := fmt.Sprintf("%s.%s", targetAccount.OrganizationName(), targetAccount.AccountName())
	uniformListingLocator := fmt.Sprintf("ORGDATACLOUD$INTERNAL$%s", id.Name())

	manifest, title := testClient().Listing.BasicOrganizationManifestWithUnquotedValues(t, targetAccount)
	differentSubtitleManifest, differentSubtitleTitle := testClient().Listing.BasicOrganizationManifestWithUnquotedValuesAndDifferentSubtitle(t, targetAccount)

	comment, newComment := random.Comment(), random.Comment()

	share, shareCleanup := testClient().Share.Create(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	modelBasic := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
		// Has to be set when a listing is not associated with a share or an application package
		WithPublish(r.BooleanFalse)

	modelComplete := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue).
		WithComment(comment)

	modelCompleteWithDifferentValues := model.OrganizationListingWithInlineManifest("test", id.Name(), differentSubtitleManifest).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanFalse).
		WithComment(newComment)

	modelUnset := model.OrganizationListingWithInlineManifest("test", id.Name(), differentSubtitleManifest).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanFalse)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OrganizationListing),
		Steps: []resource.TestStep{
			// create without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareEmpty().
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStateDraft).
						HasDistribution("ORGANIZATION").
						HasOrganizationProfileName("INTERNAL").
						HasUniformListingLocator(uniformListingLocator).
						HasIsTargeted(true).
						HasTargetAccounts(targetAccounts).
						HasComment(""),
				),
			},
			// import without optionals
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedOrganizationListingResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasNoManifest().
						HasShareEmpty().
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasTitle(title).
						HasState(sdk.ListingStateDraft).
						HasDistribution("ORGANIZATION"),
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
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ListingShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStatePublished).
						HasDistribution("ORGANIZATION").
						HasUniformListingLocator(uniformListingLocator).
						HasComment(comment),
				),
			},
			// import complete object
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedOrganizationListingResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasNoManifest().
						HasShareString(share.ID().FullyQualifiedName()).
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanTrue).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasState(sdk.ListingStatePublished).
						HasComment(comment),
				),
			},
			// alter manifest, publish, and comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromString(differentSubtitleManifest).
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanFalse).
						HasCommentString(newComment),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasTitle(differentSubtitleTitle).
						HasSubtitle("different_subtitle").
						HasState(sdk.ListingStateUnpublished).
						HasComment(newComment),
				),
			},
			// change externally (organization listings only accept comment changes through ALTER LISTING ... AS)
			{
				PreConfig: func() {
					testClient().Listing.Alter(t, sdk.NewAlterListingRequest(id).WithAlterListingAs(
						*sdk.NewAlterListingAsRequest(differentSubtitleManifest).
							WithPublish(true).
							WithReview(true).
							WithComment(comment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(modelCompleteWithDifferentValues.ResourceReference(), "publish", sdk.String(r.BooleanFalse), sdk.String(r.BooleanTrue)),
						planchecks.ExpectDrift(modelCompleteWithDifferentValues.ResourceReference(), "comment", sdk.String(newComment), sdk.String(comment)),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasPublishString(r.BooleanFalse).
						HasCommentString(newComment),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
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
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelUnset.ResourceReference()).
						HasNameString(id.Name()).
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanFalse).
						HasCommentEmpty(),
					resourceshowoutputassert.ListingShowOutput(t, modelUnset.ResourceReference()).
						HasState(sdk.ListingStateUnpublished).
						HasComment(""),
				),
			},
		},
	})
}

func TestAcc_OrganizationListing_BasicUseCase_FromStage(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	targetAccount := testClient().Context.CurrentAccountId(t)

	manifest, title := testClient().Listing.BasicOrganizationManifest(t, targetAccount)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	_ = testClient().Stage.PutInLocationWithContent(t, stage.Location()+"/organization", "manifest.yml", manifest)

	share, shareCleanup := testClient().Share.Create(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	modelBasic := model.OrganizationListingWithStagedManifestWithLocation("test", id.Name(), stage.ID(), "organization").
		// Has to be set when a listing is not associated with a share or an application package
		WithPublish(r.BooleanFalse)

	modelComplete := model.OrganizationListingWithStagedManifestWithLocation("test", id.Name(), stage.ID(), "organization").
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue)

	comment := random.Comment()
	modelCompleteWithComment := model.OrganizationListingWithStagedManifestWithLocation("test", id.Name(), stage.ID(), "organization").
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OrganizationListing),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageStageId(stage.ID()).
						HasManifestFromStageLocation("organization").
						HasShareEmpty().
						HasPublishString(r.BooleanFalse),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title).
						HasState(sdk.ListingStateDraft).
						HasDistribution("ORGANIZATION"),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStageStageId(stage.ID()).
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue),
					resourceshowoutputassert.ListingShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasState(sdk.ListingStatePublished).
						HasDistribution("ORGANIZATION"),
				),
			},
			// a comment change cannot be sent without an inline manifest, so it recreates the listing
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithComment.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithComment),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelCompleteWithComment.ResourceReference()).
						HasCommentString(comment),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithComment.ResourceReference()).
						HasState(sdk.ListingStatePublished).
						HasComment(comment),
				),
			},
		},
	})
}

func TestAcc_OrganizationListing_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	manifest, _ := testClient().Listing.BasicOrganizationManifestWithUnquotedValues(t, testClient().Context.CurrentAccountId(t))

	organizationListingModelWithoutManifest := func(resourceName string, name string) *model.OrganizationListingModel {
		o := &model.OrganizationListingModel{ResourceModelMeta: accconfig.Meta(resourceName, resources.OrganizationListing)}
		o.WithName(name)
		return o
	}

	modelWithBothShareAndApplicationPackage := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
		WithShare("test_share").
		WithApplicationPackage("test_app_package")

	modelWithInvalidStageId := organizationListingModelWithoutManifest("test", id.Name()).
		WithManifestValue(tfconfig.ListVariable(
			tfconfig.MapVariable(map[string]tfconfig.Variable{
				"from_stage": tfconfig.ListVariable(
					tfconfig.MapVariable(map[string]tfconfig.Variable{
						"stage": tfconfig.StringVariable("invalid.stage.identifier.name"),
					}),
				),
			}),
		))

	modelWithInvalidName := model.OrganizationListingWithInlineManifest("test", "_invalid_name", manifest)
	modelWithDoubleDollarInManifest := model.OrganizationListingWithInlineManifest("test", id.Name(), "manifest: contains $$ sequence")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OrganizationListing),
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelWithBothShareAndApplicationPackage),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"application_package": conflicts with share`),
			},
			{
				Config:      accconfig.FromModels(t, modelWithInvalidStageId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Expected SchemaObjectIdentifier identifier type`),
			},
			{
				Config:      accconfig.FromModels(t, modelWithInvalidName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Listing name must start with an alphabetic character and cannot contain spaces or special characters except for underscores and hyphens`),
			},
			{
				Config:      accconfig.FromModels(t, modelWithDoubleDollarInManifest),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`cannot contain the \$\$ sequence`),
			},
		},
	})
}
