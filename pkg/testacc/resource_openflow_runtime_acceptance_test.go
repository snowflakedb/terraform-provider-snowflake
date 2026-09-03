//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OpenflowRuntime_BasicUseCase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	currentRole := testClient().Context.CurrentRole(t)
	deploymentId := testClient().OpenflowDeployment.ActiveDeploymentForRuntimes(t)

	networkRule, networkRuleCleanup := testClient().NetworkRule.Create(t)
	t.Cleanup(networkRuleCleanup)
	externalAccessIntegrationId, externalAccessIntegrationCleanup := testClient().ExternalAccessIntegration.CreateExternalAccessIntegration(t, networkRule.ID())
	t.Cleanup(externalAccessIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	displayName := random.AlphaN(12)
	externallyChangedComment := random.Comment()
	externallyChangedDisplayName := random.AlphaN(12)

	basic := model.OpenflowRuntime("t", id.DatabaseName(), id.SchemaName(), id.Name(),
		deploymentId.Name(), currentRole.Name(), 1, 1, string(sdk.OpenflowRuntimeNodeTypeSmall))
	withOptionals := model.OpenflowRuntime("t", id.DatabaseName(), id.SchemaName(), id.Name(),
		deploymentId.Name(), currentRole.Name(), 2, 1, string(sdk.OpenflowRuntimeNodeTypeSmall)).
		WithDisplayName(displayName).
		WithComment(comment).
		WithExternalAccessIntegrations(externalAccessIntegrationId)
	ref := basic.ResourceReference()

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowRuntimeResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasDeploymentString(deploymentId.Name()).
			HasNodeTypeString(string(sdk.OpenflowRuntimeNodeTypeSmall)).
			HasMinNodesString("1").
			HasMaxNodesString("1").
			HasExecuteAsRoleString(currentRole.Name()).
			HasNoDisplayName().
			HasNoComment().
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowRuntimeShowOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowRuntimeStatusActive).
			HasDeployment(deploymentId.Name()).
			HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
			HasMinNodes(1).
			HasMaxNodes(1).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasExecuteAsRole(currentRole.Name()).
			HasDisplayName("").
			HasComment("").
			HasInitiallySuspended(false).
			HasKeyNotEmpty().
			HasNoExternalAccessIntegrations().
			HasOwner(currentRole.Name()).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		resourceshowoutputassert.OpenflowRuntimeDescribeOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowRuntimeStatusActive).
			HasDeployment(deploymentId.Name()).
			HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
			HasMinNodes(1).
			HasMaxNodes(1).
			HasExecuteAsRole(currentRole.Name()).
			HasDisplayName("").
			HasComment("").
			HasInitiallySuspended(false).
			HasKeyNotEmpty().
			HasServerUrlNotEmpty().
			HasNodeTypeTierNotEmpty().
			HasNoExternalAccessIntegrations().
			HasOwner(currentRole.Name()),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowRuntimeResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasDeploymentString(deploymentId.Name()).
			HasNodeTypeString(string(sdk.OpenflowRuntimeNodeTypeSmall)).
			HasMinNodesString("1").
			HasMaxNodesString("2").
			HasExecuteAsRoleString(currentRole.Name()).
			HasDisplayNameString(displayName).
			HasCommentString(comment).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowRuntimeShowOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowRuntimeStatusActive).
			HasDeployment(deploymentId.Name()).
			HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
			HasMinNodes(1).
			HasMaxNodes(2).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasExecuteAsRole(currentRole.Name()).
			HasDisplayName(displayName).
			HasComment(comment).
			HasInitiallySuspended(false).
			HasKeyNotEmpty().
			HasExternalAccessIntegrations(externalAccessIntegrationId).
			HasOwner(currentRole.Name()).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		resourceshowoutputassert.OpenflowRuntimeDescribeOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowRuntimeStatusActive).
			HasDeployment(deploymentId.Name()).
			HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
			HasMinNodes(1).
			HasMaxNodes(2).
			HasExecuteAsRole(currentRole.Name()).
			HasDisplayName(displayName).
			HasComment(comment).
			HasInitiallySuspended(false).
			HasKeyNotEmpty().
			HasServerUrlNotEmpty().
			HasNodeTypeTierNotEmpty().
			HasExternalAccessIntegrations(externalAccessIntegrationId).
			HasOwner(currentRole.Name()),
	}

	assertOptionalsUnset := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowRuntimeResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasDeploymentString(deploymentId.Name()).
			HasNodeTypeString(string(sdk.OpenflowRuntimeNodeTypeSmall)).
			HasMinNodesString("1").
			HasMaxNodesString("1").
			HasExecuteAsRoleString(currentRole.Name()).
			HasDisplayNameString("").
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowRuntimeShowOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowRuntimeStatusActive).
			HasDeployment(deploymentId.Name()).
			HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
			HasMinNodes(1).
			HasMaxNodes(1).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasExecuteAsRole(currentRole.Name()).
			HasDisplayName("").
			HasComment("").
			HasInitiallySuspended(false).
			HasKeyNotEmpty().
			HasNoExternalAccessIntegrations().
			HasOwner(currentRole.Name()).
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty(),
		resourceshowoutputassert.OpenflowRuntimeDescribeOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowRuntimeStatusActive).
			HasDeployment(deploymentId.Name()).
			HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
			HasMinNodes(1).
			HasMaxNodes(1).
			HasExecuteAsRole(currentRole.Name()).
			HasDisplayName("").
			HasComment("").
			HasInitiallySuspended(false).
			HasKeyNotEmpty().
			HasServerUrlNotEmpty().
			HasNodeTypeTierNotEmpty().
			HasNoExternalAccessIntegrations().
			HasOwner(currentRole.Name()),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowRuntime),
		Steps: []resource.TestStep{
			// 1. Create - required fields only
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// 2. Import - required fields only
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// 3. Update - set the optionals and scale max_nodes up, all in place
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 4. Import - with the optionals set
			{
				Config:            config.FromModels(t, withOptionals),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// 5. External change - both configured fields are changed outside Terraform and must be reconciled
			{
				PreConfig: func() {
					testClient().OpenflowRuntime.Alter(t, sdk.NewAlterOpenflowRuntimeRequest(id).
						WithSet(*sdk.NewOpenflowRuntimeSetRequest().
							WithComment(externallyChangedComment).
							WithDisplayName(externallyChangedDisplayName).
							WithMinNodes(2).
							WithMaxNodes(3)))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(ref, "comment", sdk.String(comment), sdk.String(externallyChangedComment)),
						planchecks.ExpectChange(ref, "comment", tfjson.ActionUpdate, sdk.String(externallyChangedComment), sdk.String(comment)),
						planchecks.ExpectDrift(ref, "display_name", sdk.String(displayName), sdk.String(externallyChangedDisplayName)),
						planchecks.ExpectChange(ref, "display_name", tfjson.ActionUpdate, sdk.String(externallyChangedDisplayName), sdk.String(displayName)),
						planchecks.ExpectDrift(ref, "min_nodes", sdk.String("1"), sdk.String("2")),
						planchecks.ExpectChange(ref, "min_nodes", tfjson.ActionUpdate, sdk.String("2"), sdk.String("1")),
						planchecks.ExpectDrift(ref, "max_nodes", sdk.String("2"), sdk.String("3")),
						planchecks.ExpectChange(ref, "max_nodes", tfjson.ActionUpdate, sdk.String("3"), sdk.String("2")),
					},
				},
				Config: config.FromModels(t, withOptionals),
				Check:  assertThat(t, assertWithOptionals...),
			},
			// 6. Update - unset the optionals again
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
		},
	})
}

func TestAcc_OpenflowRuntime_Rename(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	currentRole := testClient().Context.CurrentRole(t)
	deploymentId := testClient().OpenflowDeployment.ActiveDeploymentForRuntimes(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(id.SchemaId())

	basic := model.OpenflowRuntime("t", id.DatabaseName(), id.SchemaName(), id.Name(),
		deploymentId.Name(), currentRole.Name(), 1, 1, string(sdk.OpenflowRuntimeNodeTypeSmall))
	renamed := model.OpenflowRuntime("t", newId.DatabaseName(), newId.SchemaName(), newId.Name(),
		deploymentId.Name(), currentRole.Name(), 1, 1, string(sdk.OpenflowRuntimeNodeTypeSmall))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowRuntime),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, basic),
				Check: assertThat(
					t,
					resourceassert.OpenflowRuntimeResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Update, not replace: a rename must not destroy the runtime.
						plancheck.ExpectResourceAction(renamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, renamed),
				Check: assertThat(
					t,
					resourceassert.OpenflowRuntimeResource(t, renamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasFullyQualifiedNameString(newId.FullyQualifiedName()),
					objectassert.OpenflowRuntime(t, newId).HasName(newId.Name()),
				),
			},
		},
	})
}

// Plan-only, so nothing is created and this runs in seconds.
func TestAcc_OpenflowRuntime_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	deployment := testClient().Ids.RandomAccountObjectIdentifier()
	role := testClient().Ids.RandomAccountObjectIdentifier()

	invalidNodeType := model.OpenflowRuntime("t", id.DatabaseName(), id.SchemaName(), id.Name(),
		deployment.Name(), role.Name(), 1, 1, "INVALID")
	invalidMinNodes := model.OpenflowRuntime("t", id.DatabaseName(), id.SchemaName(), id.Name(),
		deployment.Name(), role.Name(), 1, 0, string(sdk.OpenflowRuntimeNodeTypeSmall))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowRuntime),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidNodeType),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`invalid openflow runtime node type: INVALID`),
			},
			{
				Config:      config.FromModels(t, invalidMinNodes),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected min_nodes to be at least \(1\), got 0`),
			},
		},
	})
}
