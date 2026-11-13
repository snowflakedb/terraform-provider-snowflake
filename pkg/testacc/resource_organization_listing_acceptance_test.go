//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
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

	manifest, title := testClient().Listing.OrganizationBasicManifestWithUnquotedValues(t)
	manifestWithDifferentSubtitle, titleWithDifferentSubtitle := testClient().Listing.OrganizationBasicManifestWithUnquotedValuesAndDifferentSubtitle(t)

	share, shareCleanup := testClient().Share.Create(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	modelBasic := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
		WithPublish(r.BooleanFalse)

	modelComplete := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanTrue)

	modelCompleteWithDifferentValues := model.OrganizationListingWithInlineManifest("test", id.Name(), manifestWithDifferentSubtitle).
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
						HasPublishString(r.BooleanFalse),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStateDraft).
						HasDistribution("ORGANIZATION"),
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
						HasPublishString(r.BooleanFalse),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasTitle(title).
						HasState(sdk.ListingStateDraft),
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
						HasPublishString(r.BooleanTrue),
					resourceshowoutputassert.ListingShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStatePublished).
						HasOrganizationProfileName("INTERNAL"),
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
						HasPublishString(r.BooleanTrue),
					resourceshowoutputassert.ImportedListingShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasState(sdk.ListingStatePublished),
				),
			},
			// alter manifest and unpublish
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
						HasManifestFromStringNotEmpty().
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanFalse),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasName(id.Name()).
						HasTitle(titleWithDifferentSubtitle).
						HasSubtitle("different_subtitle").
						HasState(sdk.ListingStateUnpublished),
				),
			},
			// change externally
			{
				PreConfig: func() {
					testClient().Listing.Alter(t, sdk.NewAlterListingRequest(id).WithPublish(true))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(modelCompleteWithDifferentValues.ResourceReference(), "publish", sdk.String(r.BooleanFalse), sdk.String(r.BooleanTrue)),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasPublishString(r.BooleanFalse),
					resourceshowoutputassert.ListingShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasState(sdk.ListingStateUnpublished),
				),
			},
		},
	})
}

func TestAcc_OrganizationListing_PublishNotSet(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	manifest, _ := testClient().Listing.OrganizationBasicManifestWithUnquotedValues(t)
	manifestWithDifferentSubtitle, titleWithDifferentSubtitle := testClient().Listing.OrganizationBasicManifestWithUnquotedValuesAndDifferentSubtitle(t)

	share, shareCleanup := testClient().Share.Create(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	modelWithoutPublish := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
		WithShare(share.ID().Name())

	modelWithoutPublishWithDifferentManifest := model.OrganizationListingWithInlineManifest("test", id.Name(), manifestWithDifferentSubtitle).
		WithShare(share.ID().Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OrganizationListing),
		Steps: []resource.TestStep{
			// create without publish (Snowflake publishes the listing by default)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: accconfig.FromModels(t, modelWithoutPublish),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelWithoutPublish.ResourceReference()).
						HasNameString(id.Name()).
						HasShareString(share.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanDefault),
					resourceshowoutputassert.ListingShowOutput(t, modelWithoutPublish.ResourceReference()).
						HasState(sdk.ListingStatePublished),
				),
			},
			// unpublish externally
			{
				PreConfig: func() {
					testClient().Listing.Alter(t, sdk.NewAlterListingRequest(id).WithUnpublish(true))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithoutPublish.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(modelWithoutPublish.ResourceReference(), "publish", sdk.String(r.BooleanDefault), sdk.String(r.BooleanFalse)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: accconfig.FromModels(t, modelWithoutPublish),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelWithoutPublish.ResourceReference()).
						HasPublishString(r.BooleanDefault),
					resourceshowoutputassert.ListingShowOutput(t, modelWithoutPublish.ResourceReference()).
						HasState(sdk.ListingStateUnpublished),
				),
			},
			// alter manifest without publish (the listing is not published again)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithoutPublishWithDifferentManifest.ResourceReference(), plancheck.ResourceActionUpdate),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Config: accconfig.FromModels(t, modelWithoutPublishWithDifferentManifest),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelWithoutPublishWithDifferentManifest.ResourceReference()).
						HasPublishString(r.BooleanDefault),
					resourceshowoutputassert.ListingShowOutput(t, modelWithoutPublishWithDifferentManifest.ResourceReference()).
						HasTitle(titleWithDifferentSubtitle).
						HasSubtitle("different_subtitle").
						HasState(sdk.ListingStateUnpublished),
				),
			},
		},
	})
}

func TestAcc_OrganizationListing_BasicUseCase_FromStage(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	manifest, title := testClient().Listing.OrganizationBasicManifest(t)
	manifestWithDifferentSubtitle, titleWithDifferentSubtitle := testClient().Listing.OrganizationBasicManifestWithDifferentSubtitle(t)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	_ = testClient().Stage.PutInLocationWithContent(t, stage.Location()+"/v1", "manifest.yml", manifest)
	_ = testClient().Stage.PutInLocationWithContent(t, stage.Location()+"/v2", "manifest.yml", manifestWithDifferentSubtitle)

	versionComment := random.Comment()

	modelBasic := model.OrganizationListingWithStagedManifestWithLocation("test", id.Name(), stage.ID(), "v1").
		WithPublish(r.BooleanFalse)

	modelWithNewVersion := model.OrganizationListingWithStagedManifestWithOptionals("test", id.Name(), stage.ID(), "version_name", versionComment, "v2").
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
						HasManifestFromStageStageId(stage.ID()).
						HasManifestFromStageLocation("v1").
						HasShareEmpty().
						HasApplicationPackageEmpty().
						HasPublishString(r.BooleanFalse),
					resourceshowoutputassert.ListingShowOutput(t, modelBasic.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title).
						HasSubtitle("subtitle").
						HasState(sdk.ListingStateDraft).
						HasDistribution("ORGANIZATION"),
					assert.Check(assertContainsListingVersion(t, id, "VERSION$1", "null")),
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
						HasPublishString(r.BooleanFalse),
				),
			},
			// add a new version with version name and comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithNewVersion.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelWithNewVersion),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelWithNewVersion.ResourceReference()).
						HasManifestFromStageStageId(stage.ID()).
						HasManifestFromStageLocation("v2").
						HasManifestFromStageVersionName("version_name").
						HasManifestFromStageVersionComment(versionComment).
						HasPublishString(r.BooleanFalse),
					resourceshowoutputassert.ListingShowOutput(t, modelWithNewVersion.ResourceReference()).
						HasTitle(titleWithDifferentSubtitle).
						HasSubtitle("different_subtitle").
						HasState(sdk.ListingStateDraft),
					assert.Check(assertContainsListingVersion(t, id, "VERSION$2", "version_name")),
				),
			},
		},
	})
}

func TestAcc_OrganizationListing_ApplicationPackage(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	manifest, title := testClient().Listing.OrganizationBasicManifestWithUnquotedValues(t)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	applicationPackage, applicationPackageCleanup := testClient().ApplicationPackage.CreateApplicationPackageWithReleaseChannelsDisabled(t)
	t.Cleanup(applicationPackageCleanup)

	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "manifest.yml", "")
	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "setup.sql", "CREATE APPLICATION ROLE IF NOT EXISTS APP_HELLO_SNOWFLAKE;")

	version := "v1"
	testClient().ApplicationPackage.AddApplicationPackageVersion(t, applicationPackage.ID(), stage.ID(), version)
	testClient().ApplicationPackage.SetDefaultReleaseDirective(t, applicationPackage.ID(), version)

	modelComplete := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
		WithApplicationPackage(applicationPackage.ID().Name()).
		WithPublish(r.BooleanTrue)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OrganizationListing),
		Steps: []resource.TestStep{
			// create complete with all optionals
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasManifestFromStringNotEmpty().
						HasShareEmpty().
						HasApplicationPackageString(applicationPackage.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue),
					resourceshowoutputassert.ListingShowOutput(t, modelComplete.ResourceReference()).
						HasName(id.Name()).
						HasTitle(title).
						HasState(sdk.ListingStatePublished).
						HasIsApplication(true),
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
						HasShareEmpty().
						HasApplicationPackageString(applicationPackage.ID().FullyQualifiedName()).
						HasPublishString(r.BooleanTrue),
				),
			},
		},
	})
}

func TestAcc_OrganizationListing_Rename(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	manifest, _ := testClient().Listing.OrganizationBasicManifestWithUnquotedValues(t)

	modelBasic := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
		WithPublish(r.BooleanFalse)
	modelRenamed := model.OrganizationListingWithInlineManifest("test", newId.Name(), manifest).
		WithPublish(r.BooleanFalse)

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
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// rename in place (allowed only in the DRAFT state)
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelRenamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelRenamed),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelRenamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasFullyQualifiedNameString(newId.FullyQualifiedName()),
					resourceshowoutputassert.ListingShowOutput(t, modelRenamed.ResourceReference()).
						HasName(newId.Name()).
						HasState(sdk.ListingStateDraft),
				),
			},
		},
	})
}

func TestAcc_OrganizationListing_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	manifest, _ := testClient().Listing.OrganizationBasicManifestWithUnquotedValues(t)

	modelWithBothShareAndApplicationPackage := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
		WithShare("test_share").
		WithApplicationPackage("test_app_package")

	modelWithInvalidStageId := (&model.OrganizationListingModel{ResourceModelMeta: accconfig.Meta("test", resources.OrganizationListing)}).
		WithName(id.Name()).
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
				Config:      accconfig.FromModels(t, modelWithBothShareAndApplicationPackage),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`"share": conflicts with application_package`),
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
