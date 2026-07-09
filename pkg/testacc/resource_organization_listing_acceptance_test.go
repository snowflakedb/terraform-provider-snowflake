//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OrganizationListing_Basic_Inlined(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	basicManifest, _ := testClient().Listing.OrganizationBasicManifest(t)

	comment, newComment := random.Comment(), random.Comment()

	modelBasic := model.OrganizationListingWithInlineManifest("test", id.Name(), basicManifest).
		WithPublish(r.BooleanFalse)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	modelComplete := model.OrganizationListingWithInlineManifest("test", id.Name(), basicManifest).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanFalse).
		WithComment(comment)

	modelCompleteUpdated := model.OrganizationListingWithInlineManifest("test", id.Name(), basicManifest).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanFalse).
		WithComment(newComment)

	modelUnset := model.OrganizationListingWithInlineManifest("test", id.Name(), basicManifest).
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
						HasState(sdk.ListingStateDraft).
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
				),
			},
			// set optionals (expect re-creation because share is ForceNew)
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
						HasShareString(share.ID().Name()).
						HasPublishString(r.BooleanFalse).
						HasCommentString(comment),
				),
			},
			// update comment
			{
				Config: accconfig.FromModels(t, modelCompleteUpdated),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelCompleteUpdated.ResourceReference()).
						HasCommentString(newComment),
				),
			},
			// unset comment
			{
				Config: accconfig.FromModels(t, modelUnset),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelUnset.ResourceReference()).
						HasCommentEmpty(),
				),
			},
		},
	})
}

func TestAcc_OrganizationListing_Rename(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()

	manifest, _ := testClient().Listing.OrganizationBasicManifest(t)

	modelOriginal := model.OrganizationListingWithInlineManifest("test", id.Name(), manifest).
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
				Config: accconfig.FromModels(t, modelOriginal),
				Check: assertThat(
					t,
					resourceassert.OrganizationListingResource(t, modelOriginal.ResourceReference()).
						HasNameString(id.Name()),
				),
			},
			// rename in-place (ALTER LISTING ... RENAME TO)
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
						HasNameString(newId.Name()),
				),
			},
		},
	})
}

func TestAcc_OrganizationListing_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()
	manifest, _ := testClient().Listing.OrganizationBasicManifest(t)

	organizationListingModelWithoutManifest := func(resourceName string, name string) *model.OrganizationListingModel {
		l := &model.OrganizationListingModel{ResourceModelMeta: accconfig.Meta(resourceName, resources.OrganizationListing)}
		l.WithName(name)
		return l
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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// share and application_package are mutually exclusive
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
			// invalid stage identifier
			{
				Config:      accconfig.FromModels(t, modelWithInvalidStageId),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Expected SchemaObjectIdentifier identifier type`),
			},
		},
	})
}
