//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	resourcehelpers "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

// BYOC rather than Snowflake-managed: it settles at INACTIVE in about ninety seconds against six or seven
// minutes to reach ACTIVE, which makes it the cheaper of the two to exercise across nine steps.
func TestAcc_OpenflowDeploymentByoc_BasicUseCase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	currentRole := testClient().Context.CurrentRole(t)
	// event_table is a parameter, so with none set on the deployment it reports what the account has.
	accountEventTable := helpers.FindParameter(t, testClient().Parameter.ShowAccountParameters(t), sdk.AccountParameterEventTable).Value

	eventTable, eventTableCleanup := testClient().EventTable.Create(t)
	t.Cleanup(eventTableCleanup)
	comment := random.Comment()
	displayName := random.AlphaN(12)
	externallyChangedComment := random.Comment()

	// vpc_type is required: Snowflake refuses a BYOC deployment without it.
	basic := model.OpenflowDeploymentByoc("t", id.Name(), string(sdk.OpenflowVpcTypeManaged))
	withOptionals := model.OpenflowDeploymentByoc("t", id.Name(), string(sdk.OpenflowVpcTypeManaged)).
		WithDisplayName(displayName).
		WithComment(comment).
		WithEventTable(eventTable.ID().FullyQualifiedName())
	// vpc_type, custom_ingress_hostname and the private-link flags are ForceNew, so reaching them means
	// creating the deployment afresh.
	complete := model.OpenflowDeploymentByoc("t", id.Name(), string(sdk.OpenflowVpcTypeProvided)).
		WithCustomIngressHostname("openflow.example.com").
		WithUsePrivateLink(r.BooleanTrue).
		WithUseUserAuthOverPrivatelink(r.BooleanTrue).
		WithDisplayName(displayName).
		WithComment(comment).
		WithEventTable(eventTable.ID().FullyQualifiedName())
	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentByocResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeByoc)).
			HasVpcTypeString(string(sdk.OpenflowVpcTypeManaged)).
			HasNoDisplayName().
			HasNoComment().
			HasEventTableString(accountEventTable).
			HasUsePrivateLinkString(r.BooleanDefault).
			HasUseUserAuthOverPrivatelinkString(r.BooleanDefault).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowDeploymentShowOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeByoc).
			HasVpcType(sdk.OpenflowVpcTypeManaged).
			HasUsePrivateLink(false).
			HasUseUserAuthOverPrivateLink(false).
			HasDisplayName("").
			HasComment("").
			HasOwner(currentRole.Name()).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		resourceshowoutputassert.OpenflowDeploymentDescribeOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeByoc).
			HasVpcType(sdk.OpenflowVpcTypeManaged).
			HasOwner(currentRole.Name()),
		resourceparametersassert.OpenflowDeploymentResourceParameters(t, ref).
			HasEventTable(accountEventTable).
			HasEventTableLevel(sdk.ParameterTypeAccount),
	}

	assertOptionalsUnset := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentByocResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeByoc)).
			HasVpcTypeString(string(sdk.OpenflowVpcTypeManaged)).
			HasDisplayNameString("").
			HasCommentString("").
			HasEventTableString(accountEventTable).
			HasUsePrivateLinkString(r.BooleanDefault).
			HasUseUserAuthOverPrivatelinkString(r.BooleanDefault).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowDeploymentShowOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeByoc).
			HasVpcType(sdk.OpenflowVpcTypeManaged).
			HasUsePrivateLink(false).
			HasUseUserAuthOverPrivateLink(false).
			HasDisplayName("").
			HasComment("").
			HasOwner(currentRole.Name()).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		resourceparametersassert.OpenflowDeploymentResourceParameters(t, ref).
			HasEventTable(accountEventTable).
			HasEventTableLevel(sdk.ParameterTypeAccount),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentByocResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeByoc)).
			HasVpcTypeString(string(sdk.OpenflowVpcTypeManaged)).
			HasDisplayNameString(displayName).
			HasCommentString(comment).
			HasEventTableString(eventTable.ID().FullyQualifiedName()).
			HasUsePrivateLinkString(r.BooleanDefault).
			HasUseUserAuthOverPrivatelinkString(r.BooleanDefault).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowDeploymentShowOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeByoc).
			HasVpcType(sdk.OpenflowVpcTypeManaged).
			HasUsePrivateLink(false).
			HasUseUserAuthOverPrivateLink(false).
			HasDisplayName(displayName).
			HasComment(comment).
			HasOwner(currentRole.Name()).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		resourceparametersassert.OpenflowDeploymentResourceParameters(t, ref).
			HasEventTable(eventTable.ID().FullyQualifiedName()).
			HasEventTableLevel(sdk.ParameterTypeOpenflowDeployment),
	}

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowDeploymentByocResource(t, ref).
			HasNameString(id.Name()).
			HasTypeString(string(sdk.OpenflowDeploymentTypeByoc)).
			HasVpcTypeString(string(sdk.OpenflowVpcTypeProvided)).
			HasCustomIngressHostnameString("openflow.example.com").
			HasUsePrivateLinkString(r.BooleanTrue).
			HasUseUserAuthOverPrivatelinkString(r.BooleanTrue).
			HasDisplayNameString(displayName).
			HasCommentString(comment).
			HasEventTableString(eventTable.ID().FullyQualifiedName()).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowDeploymentShowOutput(t, ref).
			HasName(id.Name()).
			HasType(sdk.OpenflowDeploymentTypeByoc).
			HasVpcType(sdk.OpenflowVpcTypeProvided).
			HasCustomIngressHostname("openflow.example.com").
			HasUsePrivateLink(true).
			HasUseUserAuthOverPrivateLink(true).
			HasDisplayName(displayName).
			HasComment(comment).
			HasOwner(currentRole.Name()).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		resourceparametersassert.OpenflowDeploymentResourceParameters(t, ref).
			HasEventTable(eventTable.ID().FullyQualifiedName()).
			HasEventTableLevel(sdk.ParameterTypeOpenflowDeployment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowDeploymentByoc),
		Steps: []resource.TestStep{
			// 1. Create - without optionals
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// 2. Import - without optionals
			{
				Config:       config.FromModels(t, basic),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedOpenflowDeploymentByocResource(t, resourcehelpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasTypeString(string(sdk.OpenflowDeploymentTypeByoc)).
						HasVpcTypeString(string(sdk.OpenflowVpcTypeManaged)).
						HasNoDisplayName().
						HasNoComment().
						HasUsePrivateLinkString(r.BooleanDefault).
						HasUseUserAuthOverPrivatelinkString(r.BooleanDefault).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImportedOpenflowDeploymentShowOutput(t, resourcehelpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasType(sdk.OpenflowDeploymentTypeByoc).
						HasVpcType(sdk.OpenflowVpcTypeManaged).
						HasOwner(currentRole.Name()),
				),
			},
			// 3. Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 4. Import - with optionals
			{
				Config:       config.FromModels(t, withOptionals),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedOpenflowDeploymentByocResource(t, resourcehelpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasTypeString(string(sdk.OpenflowDeploymentTypeByoc)).
						HasVpcTypeString(string(sdk.OpenflowVpcTypeManaged)).
						HasDisplayNameString(displayName).
						HasCommentString(comment).
						HasUsePrivateLinkString(r.BooleanDefault).
						HasUseUserAuthOverPrivatelinkString(r.BooleanDefault).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImportedOpenflowDeploymentShowOutput(t, resourcehelpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasDisplayName(displayName).
						HasComment(comment),
				),
			},
			// 5. Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
						planchecks.ExpectChange(ref, "display_name", tfjson.ActionUpdate, sdk.String(displayName), nil),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertOptionalsUnset...),
			},
			// 6. Destroy
			{
				Config:  config.FromModels(t, basic),
				Destroy: true,
			},
			// 7. Create - with all fields, including the ForceNew networking options
			{
				PreConfig: func() {
					_, err := testClient().OpenflowDeployment.Show(t, id)
					require.ErrorIs(t, err, sdk.ErrObjectNotFound)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
			// 8. Import - with all fields
			{
				Config:                  config.FromModels(t, complete),
				ResourceName:            ref,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"use_private_link", "use_user_auth_over_privatelink"},
			},
			// 9. Update - external changes
			{
				PreConfig: func() {
					testClient().OpenflowDeployment.Alter(t, sdk.NewAlterOpenflowDeploymentRequest(id).WithSet(
						*sdk.NewOpenflowDeploymentSetRequest().WithComment(externallyChangedComment),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(complete.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externallyChangedComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externallyChangedComment), sdk.String(comment)),
						planchecks.ExpectNoChangeOnField(ref, "vpc_type"),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_OpenflowDeploymentByoc_Rename(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	id := testClient().Ids.RandomAccountObjectIdentifier()
	newId := testClient().Ids.RandomAccountObjectIdentifier()
	renamedComment := random.Comment()

	basic := model.OpenflowDeploymentByoc("t", id.Name(), string(sdk.OpenflowVpcTypeManaged))
	// The comment changes in the same step, proving the rename composes with a SET.
	renamed := model.OpenflowDeploymentByoc("t", newId.Name(), string(sdk.OpenflowVpcTypeManaged)).
		WithComment(renamedComment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowDeploymentByoc),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, basic),
				Check: assertThat(
					t,
					resourceassert.OpenflowDeploymentByocResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Update, not replace: a rename must not destroy the deployment.
						plancheck.ExpectResourceAction(renamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, renamed),
				Check: assertThat(
					t,
					resourceassert.OpenflowDeploymentByocResource(t, renamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasFullyQualifiedNameString(newId.FullyQualifiedName()).
						HasCommentString(renamedComment),
					// Asserted against Snowflake as well as state, so a rename that updated only Terraform
					// would fail here.
					objectassert.OpenflowDeployment(t, newId).HasName(newId.Name()),
				),
			},
		},
	})
}

// Plan-only, so nothing is created and this runs in seconds.
func TestAcc_OpenflowDeploymentByoc_Validations(t *testing.T) {
	id := testClient().Ids.RandomAccountObjectIdentifier()

	invalidVpcType := model.OpenflowDeploymentByoc("t", id.Name(), "INVALID")
	invalidUsePrivateLink := model.OpenflowDeploymentByoc("t", id.Name(), string(sdk.OpenflowVpcTypeManaged)).
		WithUsePrivateLink("invalid")
	invalidUseUserAuthOverPrivatelink := model.OpenflowDeploymentByoc("t", id.Name(), string(sdk.OpenflowVpcTypeManaged)).
		WithUseUserAuthOverPrivatelink("invalid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowDeploymentByoc),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidVpcType),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid openflow vpc type: INVALID`),
			},
			{
				Config:      config.FromModels(t, invalidUsePrivateLink),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected.*use_private_link.*to be one of \["true" "false"\], got invalid`),
			},
			{
				Config:      config.FromModels(t, invalidUseUserAuthOverPrivatelink),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected.*use_user_auth_over_privatelink.*to be one of \["true" "false"\], got invalid`),
			},
		},
	})
}
