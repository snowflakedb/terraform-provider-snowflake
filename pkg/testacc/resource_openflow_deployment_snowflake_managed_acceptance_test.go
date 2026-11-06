//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/experimentalfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

// A Snowflake-managed deployment takes six to seven minutes to reach ACTIVE, with a comparable terminate, so
// this covers the lifecycle in six steps rather than the nine the cheaper BYOC resource uses.
//
// use_private_link is ForceNew, so setting it replaces the deployment. display_name and comment still
// change in place. Rename and the other networking options are covered on the BYOC resource.
func TestAcc_OpenflowDeploymentSnowflakeManaged_BasicUseCase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	currentRole := testClient().Context.CurrentRole(t)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()
	displayName := random.AlphaN(12)
	externallyChangedComment := random.Comment()
	externallyChangedDisplayName := random.AlphaN(12)

	basic := model.OpenflowDeploymentSnowflakeManaged("t", id.Name())
	withOptionals := model.OpenflowDeploymentSnowflakeManaged("t", id.Name()).
		WithDisplayName(displayName).
		WithComment(comment).
		WithUsePrivateLink(r.BooleanTrue)
	ref := basic.ResourceReference()
	providerModelImportPermissiveMode := providermodel.SnowflakeProvider().
		WithExperimentalFeaturesEnabled(experimentalfeatures.ImportPermissiveMode)

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentSnowflakeManagedResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeSnowflake)).
			HasNoDisplayName().
			HasNoComment().
			HasUsePrivateLinkString(r.BooleanDefault).
			HasEventTableString(accountEventTable).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowDeploymentShowOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeSnowflake).
			HasStatus(sdk.OpenflowDeploymentStatusActive).
			HasUsePrivateLink(false).
			HasOwner(currentRole.Name()),
		resourceshowoutputassert.OpenflowDeploymentDescribeOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeSnowflake).
			HasUsePrivateLink(false).
			HasOwner(currentRole.Name()),
		objectassert.OpenflowDeployment(t, id).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeSnowflake).
			// Unlike BYOC, a Snowflake-managed deployment reaches ACTIVE on its own, so create is only
			// complete once it gets there.
			HasStatus(sdk.OpenflowDeploymentStatusActive).
			HasOwner(currentRole.Name()).
			HasNoDisplayName().
			HasNoComment().
			HasUsePrivateLink(false).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
	}

	assertOptionalsUnset := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentSnowflakeManagedResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeSnowflake)).
			HasDisplayNameString("").
			HasCommentString("").
			HasUsePrivateLinkString(r.BooleanDefault).
			HasEventTableString(accountEventTable).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		objectassert.OpenflowDeployment(t, id).
			HasName(id.Name()).
			HasNoDisplayName().
			HasNoComment().
			HasUsePrivateLink(false),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentSnowflakeManagedResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeSnowflake)).
			HasDisplayNameString(displayName).
			HasCommentString(comment).
			HasUsePrivateLinkString(r.BooleanTrue),
		resourceshowoutputassert.OpenflowDeploymentShowOutput(t, ref).
			HasName(id.Name()).
			HasDisplayName(displayName).
			HasComment(comment).
			HasUsePrivateLink(true),
		objectassert.OpenflowDeployment(t, id).
			HasDisplayName(displayName).
			HasComment(comment).
			HasUsePrivateLink(true).
			HasStatus(sdk.OpenflowDeploymentStatusActive),
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowDeploymentSnowflakeManaged),
		Steps: []resource.TestStep{
			// 1. Create - required fields only
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, basic),
				Check:                    assertThat(t, assertBasic...),
			},
			// 2. Import - permissive mode. An omitted use_private_link stays "default" and does not replace.
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   config.FromModels(t, providerModelImportPermissiveMode, basic),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			// 3. Set the optionals. display_name and comment can change in place; use_private_link cannot,
			// so the deployment is replaced.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "use_private_link", tfjson.ActionDelete, new(r.BooleanDefault), new(r.BooleanTrue)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 4. Import - permissive mode, with the optionals set
			{
				ProtoV6ProviderFactories: importPermissiveModeProviderFactory,
				Config:                   config.FromModels(t, providerModelImportPermissiveMode, withOptionals),
				ResourceName:             ref,
				ImportState:              true,
				ImportStateKind:          resource.ImportBlockWithID,
			},
			// 5. External change - both configured fields are changed outside Terraform and must be reconciled
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				PreConfig: func() {
					testClient().OpenflowDeployment.Alter(t, sdk.NewAlterOpenflowDeploymentRequest(id).
						WithSet(*sdk.NewOpenflowDeploymentSetRequest().
							WithComment(externallyChangedComment).
							WithDisplayName(externallyChangedDisplayName)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externallyChangedComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externallyChangedComment), sdk.String(comment)),
						planchecks.ExpectDrift(ref, "display_name", sdk.String(displayName), sdk.String(externallyChangedDisplayName)),
						planchecks.ExpectChange(ref, "display_name", tfjson.ActionUpdate, sdk.String(externallyChangedDisplayName), sdk.String(displayName)),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 6. Unset the optionals. Dropping use_private_link replaces the deployment again.
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "use_private_link", tfjson.ActionDelete, new(r.BooleanTrue), new(r.BooleanDefault)),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertOptionalsUnset...),
			},
		},
	})
}

// Enforcing import writes the SHOW value of use_private_link. The attribute is ForceNew, so applying a
// config that omits it replaces the deployment.
func TestAcc_OpenflowDeploymentSnowflakeManaged_BasicUseCase_Import(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	basic := model.OpenflowDeploymentSnowflakeManaged("t", id.Name())
	ref := basic.ResourceReference()
	resourceId := resourcehelpers.EncodeResourceIdentifier(id)

	_, cleanup := testClient().OpenflowDeployment.CreateWithRequest(t, sdk.NewCreateOpenflowDeploymentRequest(id, sdk.OpenflowDeploymentTypeSnowflake))
	t.Cleanup(cleanup)

	deployment, err := testClient().OpenflowDeployment.Show(t, id)
	require.NoError(t, err)
	usePrivateLink := fmt.Sprintf("%t", deployment.UsePrivateLink)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowDeploymentSnowflakeManaged),
		Steps: []resource.TestStep{
			{
				Config:             config.FromModels(t, basic),
				ResourceName:       ref,
				ImportState:        true,
				ImportStateId:      resourceId,
				ImportStatePersist: true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedOpenflowDeploymentSnowflakeManagedResource(t, resourceId).
						HasName(id.Name()).
						HasFullyQualifiedName(id.FullyQualifiedName()).
						HasType(string(sdk.OpenflowDeploymentTypeSnowflake)).
						HasUsePrivateLink(usePrivateLink).
						HasEventTable(accountEventTable).
						HasNoDisplayName().
						HasNoComment(),
					resourceshowoutputassert.ImportedOpenflowDeploymentShowOutput(t, resourceId).
						HasName(id.Name()).
						HasType(sdk.OpenflowDeploymentTypeSnowflake).
						HasUsePrivateLink(deployment.UsePrivateLink),
				),
			},
			{
				Config: config.FromModels(t, basic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
						planchecks.ExpectChange(ref, "use_private_link", tfjson.ActionDelete, new(usePrivateLink), new(r.BooleanDefault)),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assertThat(
					t,
					resourceassert.OpenflowDeploymentSnowflakeManagedResource(t, ref).
						HasName(id.Name()).
						HasFullyQualifiedName(id.FullyQualifiedName()).
						HasType(string(sdk.OpenflowDeploymentTypeSnowflake)).
						HasUsePrivateLink(r.BooleanDefault).
						HasEventTable(accountEventTable).
						HasNoDisplayName().
						HasNoComment(),
				),
			},
		},
	})
}

// v2.22.1 state has no use_private_link. The attribute is ForceNew, so leaving the upgraded value empty
// would plan a replacement. The upgrader writes "default", which matches an omitted configuration.
func TestAcc_OpenflowDeploymentSnowflakeManaged_migrateFromV2_22_1(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	basic := model.OpenflowDeploymentSnowflakeManaged("t", id.Name())
	ref := basic.ResourceReference()
	providerModel := providermodel.SnowflakeProvider().
		WithPreviewFeaturesEnabled(string(previewfeatures.OpenflowDeploymentSnowflakeManagedResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowDeploymentSnowflakeManaged),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.22.1"),
				Config:            config.FromModels(t, providerModel, basic),
				Check: assertThat(t,
					resourceassert.OpenflowDeploymentSnowflakeManagedResource(t, ref).
						HasNameString(id.Name()).
						HasNoUsePrivateLink(),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, providerModel, basic),
				Check: assertThat(t,
					resourceassert.OpenflowDeploymentSnowflakeManagedResource(t, ref).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasUsePrivateLinkString(r.BooleanDefault),
					resourceshowoutputassert.OpenflowDeploymentShowOutput(t, ref).
						HasName(id.Name()).
						HasUsePrivateLink(false),
				),
			},
		},
	})
}

// Plan-only, so nothing is created and this runs in seconds.
func TestAcc_OpenflowDeploymentSnowflakeManaged_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	invalidUsePrivateLink := model.OpenflowDeploymentSnowflakeManaged("t", id.Name()).
		WithUsePrivateLink("invalid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowDeploymentSnowflakeManaged),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidUsePrivateLink),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected.*use_private_link.*to be one of \["true" "false"\], got invalid`),
			},
		},
	})
}
