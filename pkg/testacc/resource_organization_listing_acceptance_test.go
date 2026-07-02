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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// orgListingManifest returns a minimal valid organization listing manifest for acceptance tests; see
// https://docs.snowflake.com/en/user-guide/collaboration/listings/organizational/org-listing-manifest-reference.
func orgListingManifest(title string) string {
	return `title: "` + title + `"
subtitle: "subtitle"
description: "description"
organization_targets:
  access:
  - all_internal_accounts: true
locations:
  access_regions:
  - name: "ALL"
`
}

func TestAcc_OrganizationListing_Basic_Inlined(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	basicManifest := orgListingManifest(id.Name())
	listingTitle := id.Name()

	comment, newComment := random.Comment(), random.Comment()

	modelBasic := model.OrganizationListingWithInlineManifest("test", id.Name(), basicManifest).
		WithPublish(r.BooleanFalse)

	share, shareCleanup := testClient().Share.CreateShare(t)
	t.Cleanup(shareCleanup)
	t.Cleanup(testClient().Grant.GrantPrivilegeOnDatabaseToShare(t, testClient().Ids.DatabaseId(), share.ID(), []sdk.ObjectPrivilege{sdk.ObjectPrivilegeUsage}))

	manifestWithShare := orgListingManifest(listingTitle)
	modelComplete := model.OrganizationListingWithInlineManifest("test", id.Name(), manifestWithShare).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanFalse).
		WithComment(comment)

	modelCompleteUpdated := model.OrganizationListingWithInlineManifest("test", id.Name(), manifestWithShare).
		WithShare(share.ID().Name()).
		WithPublish(r.BooleanFalse).
		WithComment(newComment)

	modelUnset := model.OrganizationListingWithInlineManifest("test", id.Name(), manifestWithShare).
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

	manifest := orgListingManifest(id.Name())

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
	manifest := orgListingManifest(id.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// share and application_package are mutually exclusive
			{
				Config:      organizationListingConfigWithShareAndAppPackage(id.Name(), manifest),
				ExpectError: regexp.MustCompile(`conflicts with`),
			},
			// invalid stage identifier
			{
				Config:      organizationListingConfigWithInvalidStage(id.Name()),
				ExpectError: regexp.MustCompile(`is not a valid identifier`),
			},
		},
	})
}

func organizationListingConfigWithShareAndAppPackage(name, manifest string) string {
	return `
resource "snowflake_organization_listing" "test" {
  name = "` + name + `"
  manifest {
    from_string = <<-EOT
` + manifest + `
EOT
  }
  share               = "my_share"
  application_package = "my_app_package"
}
`
}

func organizationListingConfigWithInvalidStage(name string) string {
	return `
resource "snowflake_organization_listing" "test" {
  name = "` + name + `"
  manifest {
    from_stage {
      stage = "not..a..valid..stage"
    }
  }
}
`
}
