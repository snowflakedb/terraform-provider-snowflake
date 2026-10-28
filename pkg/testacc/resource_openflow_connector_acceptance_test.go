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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_OpenflowConnector_BasicUseCase(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	currentRole := testClient().Context.CurrentRole(t)
	runtimeId := testClient().OpenflowRuntime.ActiveRuntime(t)
	definition := testClient().OpenflowConnectorDefinition.ForTesting(t)

	// A connector is created in its runtime's schema, so its identifier is derived from the runtime's.
	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())
	comment := random.Comment()
	displayName := random.AlphaN(12)
	externallyChangedComment := random.Comment()
	externallyChangedDisplayName := random.AlphaN(12)

	basic := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), definition.Name, runtimeId.FullyQualifiedName())
	withOptionals := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), definition.Name, runtimeId.FullyQualifiedName()).
		WithDisplayName(displayName).
		WithComment(comment)
	ref := basic.ResourceReference()

	// A connector from a definition settles on STOPPED: it only runs once a version is committed, which is an
	// operational step this resource does not own. default_version defaults to LAST, and every other version
	// column stays empty until a commit happens, which is why they are asserted as empty rather than skipped.
	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowConnectorResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasRuntimeString(runtimeId.FullyQualifiedName()).
			HasFromDefinition(definition.Name).
			HasNoDisplayName().
			HasNoComment().
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowConnectorShowOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowConnectorStatusStopped).
			HasRuntime(runtimeId.Name()).
			HasConnectorDefinition(definition.Name).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasDisplayName("").
			HasComment("").
			HasOwner(currentRole.Name()).
			HasDefaultVersion("LAST").
			HasDefaultVersionName("").
			HasDefaultVersionAlias("").
			HasDefaultVersionLocationUri("").
			HasDefaultVersionSourceLocationUri("").
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasConnectorUrlNotEmpty().
			HasLiveVersionLocationUriNotEmpty(),
		resourceshowoutputassert.OpenflowConnectorDescribeOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowConnectorStatusStopped).
			HasRuntime(runtimeId.Name()).
			HasConnectorDefinition(definition.Name).
			HasDisplayName("").
			HasComment("").
			HasOwner(currentRole.Name()).
			HasDefaultVersion("LAST").
			HasDefaultVersionName("").
			HasDefaultVersionAlias("").
			HasDefaultVersionLocationUri("").
			HasDefaultVersionSourceLocationUri("").
			HasDefaultVersionGitCommitHash("").
			HasLastVersionName("").
			HasLastVersionAlias("").
			HasLastVersionLocationUri("").
			HasLastVersionSourceLocationUri("").
			HasLastVersionGitCommitHash("").
			HasConnectorUrlNotEmpty().
			HasLiveVersionLocationUriNotEmpty(),
	}

	assertWithOptionals := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowConnectorResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasRuntimeString(runtimeId.FullyQualifiedName()).
			HasFromDefinition(definition.Name).
			HasDisplayNameString(displayName).
			HasCommentString(comment).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowConnectorShowOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowConnectorStatusStopped).
			HasRuntime(runtimeId.Name()).
			HasConnectorDefinition(definition.Name).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasDisplayName(displayName).
			HasComment(comment).
			HasOwner(currentRole.Name()).
			HasDefaultVersion("LAST").
			HasDefaultVersionName("").
			HasDefaultVersionAlias("").
			HasDefaultVersionLocationUri("").
			HasDefaultVersionSourceLocationUri("").
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasConnectorUrlNotEmpty().
			HasLiveVersionLocationUriNotEmpty(),
		resourceshowoutputassert.OpenflowConnectorDescribeOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowConnectorStatusStopped).
			HasRuntime(runtimeId.Name()).
			HasConnectorDefinition(definition.Name).
			HasDisplayName(displayName).
			HasComment(comment).
			HasOwner(currentRole.Name()).
			HasDefaultVersion("LAST").
			HasDefaultVersionName("").
			HasDefaultVersionAlias("").
			HasDefaultVersionLocationUri("").
			HasDefaultVersionSourceLocationUri("").
			HasDefaultVersionGitCommitHash("").
			HasLastVersionName("").
			HasLastVersionAlias("").
			HasLastVersionLocationUri("").
			HasLastVersionSourceLocationUri("").
			HasLastVersionGitCommitHash("").
			HasConnectorUrlNotEmpty().
			HasLiveVersionLocationUriNotEmpty(),
	}

	assertOptionalsUnset := []assert.TestCheckFuncProvider{
		resourceassert.OpenflowConnectorResource(t, ref).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasRuntimeString(runtimeId.FullyQualifiedName()).
			HasFromDefinition(definition.Name).
			HasDisplayNameString("").
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.OpenflowConnectorShowOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowConnectorStatusStopped).
			HasRuntime(runtimeId.Name()).
			HasConnectorDefinition(definition.Name).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasDisplayName("").
			HasComment("").
			HasOwner(currentRole.Name()).
			HasDefaultVersion("LAST").
			HasDefaultVersionName("").
			HasDefaultVersionAlias("").
			HasDefaultVersionLocationUri("").
			HasDefaultVersionSourceLocationUri("").
			HasCreatedOnNotEmpty().
			HasUpdatedOnNotEmpty().
			HasConnectorUrlNotEmpty().
			HasLiveVersionLocationUriNotEmpty(),
		resourceshowoutputassert.OpenflowConnectorDescribeOutput(t, ref).
			HasName(id.Name()).
			HasStatus(sdk.OpenflowConnectorStatusStopped).
			HasRuntime(runtimeId.Name()).
			HasConnectorDefinition(definition.Name).
			HasDisplayName("").
			HasComment("").
			HasOwner(currentRole.Name()).
			HasDefaultVersion("LAST").
			HasDefaultVersionName("").
			HasDefaultVersionAlias("").
			HasDefaultVersionLocationUri("").
			HasDefaultVersionSourceLocationUri("").
			HasDefaultVersionGitCommitHash("").
			HasLastVersionName("").
			HasLastVersionAlias("").
			HasLastVersionLocationUri("").
			HasLastVersionSourceLocationUri("").
			HasLastVersionGitCommitHash("").
			HasConnectorUrlNotEmpty().
			HasLiveVersionLocationUriNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowConnector),
		Steps: []resource.TestStep{
			// 1. Create - required fields only
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, assertBasic...),
			},
			// 2. Import - required fields only
			{
				Config:       config.FromModels(t, basic),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedOpenflowConnectorResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasRuntimeString(runtimeId.FullyQualifiedName()).
						// Not imported, because SHOW reports a definition whichever source created the
						// connector, so an imported one could be wrong.
						HasFromEmpty().
						HasNoDisplayName().
						HasNoComment(),
				),
			},
			// 3. Update - set the optionals in place
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
				Config:       config.FromModels(t, withOptionals),
				ResourceName: ref,
				ImportState:  true,
				ImportStateCheck: assertThatImport(
					t,
					resourceassert.ImportedOpenflowConnectorResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasRuntimeString(runtimeId.FullyQualifiedName()).
						HasFromEmpty().
						HasDisplayNameString(displayName).
						HasCommentString(comment),
				),
			},
			// 5. External change - both configured fields are changed outside Terraform and must be reconciled
			{
				PreConfig: func() {
					testClient().OpenflowConnector.Alter(t, sdk.NewAlterOpenflowConnectorRequest(id).
						WithSet(*sdk.NewOpenflowConnectorSetRequest().
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

func TestAcc_OpenflowConnector_Rename(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	runtimeId := testClient().OpenflowRuntime.ActiveRuntime(t)
	definition := testClient().OpenflowConnectorDefinition.ForTesting(t)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())
	newId := testClient().Ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())

	basic := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), definition.Name, runtimeId.FullyQualifiedName())
	renamed := model.OpenflowConnector("t", newId.DatabaseName(), newId.SchemaName(), newId.Name(), definition.Name, runtimeId.FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowConnector),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, basic),
				Check: assertThat(
					t,
					resourceassert.OpenflowConnectorResource(t, basic.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Update, not replace: a rename must not destroy the connector.
						plancheck.ExpectResourceAction(renamed.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, renamed),
				Check: assertThat(
					t,
					resourceassert.OpenflowConnectorResource(t, renamed.ResourceReference()).
						HasNameString(newId.Name()).
						HasFullyQualifiedNameString(newId.FullyQualifiedName()),
					objectassert.OpenflowConnector(t, newId).HasName(newId.Name()),
				),
			},
		},
	})
}

// Covers `from`, which is the only way to create a connector that already carries a configuration. It applies
// the configuration rather than only planning it: the framework fails on a non-empty plan after apply, which is
// what catches a field read into state that the configuration cannot hold. Both a bare stage and a stage with a
// path are created from, because Snowflake resolves the two differently.
func TestAcc_OpenflowConnector_FromStage(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	const bundlePath = "orders"

	runtimeId := testClient().OpenflowRuntime.ActiveRuntime(t)
	definition := testClient().OpenflowConnectorDefinition.ForTesting(t)

	stage, stageCleanup := testClient().Stage.CreateStageInSchema(t, runtimeId.SchemaId())
	t.Cleanup(stageCleanup)

	// Snowflake reads config.json from the location and checks that the definition it names matches the
	// connector's. Nothing else in the bundle is needed to create one. A copy goes at the stage root and
	// another under a subdirectory, one for each step.
	bundle := fmt.Sprintf(`{"configFormatVersion":1,"connectorDefinitionId":%q,"configuration":[]}`, definition.Name)
	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "config.json", bundle)
	testClient().Stage.PutInLocationWithContent(t, fmt.Sprintf("@%s/%s", stage.ID().FullyQualifiedName(), bundlePath), "config.json", bundle)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())
	withPath := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), "", runtimeId.FullyQualifiedName()).
		WithFromStage(stage.ID(), bundlePath)
	atRoot := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), "", runtimeId.FullyQualifiedName()).
		WithFromStage(stage.ID(), "")
	ref := withPath.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowConnector),
		Steps: []resource.TestStep{
			// 1. Create from a subdirectory of the stage. Snowflake only finds the bundle when the path it is
			// given ends in a slash, so this fails unless the provider adds one.
			{
				Config: config.FromModels(t, withPath),
				Check: assertThat(
					t,
					resourceassert.OpenflowConnectorResource(t, ref).
						HasNameString(id.Name()).
						HasRuntimeString(runtimeId.FullyQualifiedName()).
						// The path is held as written, without the slash the create added, and the definition
						// stays empty even though Snowflake resolved one from the bundle.
						HasFromStage(stage.ID(), bundlePath),
					resourceshowoutputassert.OpenflowConnectorShowOutput(t, ref).
						HasName(id.Name()).
						HasConnectorDefinition(definition.Name).
						HasStatus(sdk.OpenflowConnectorStatusStopped),
				),
			},
			// 2. Re-planning the same configuration must be a no-op. It is what catches the definition being
			// read into state, which the configuration cannot hold next to `stage`, and it also catches the
			// added slash leaking into state.
			{
				Config:   config.FromModels(t, withPath),
				PlanOnly: true,
			},
			// 3. Create from the stage root instead. `from` is create-only, so dropping the path replaces the
			// connector rather than altering it.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, atRoot),
				Check: assertThat(
					t,
					resourceassert.OpenflowConnectorResource(t, ref).
						HasNameString(id.Name()).
						HasFromStage(stage.ID(), ""),
					resourceshowoutputassert.OpenflowConnectorShowOutput(t, ref).
						HasConnectorDefinition(definition.Name).
						HasStatus(sdk.OpenflowConnectorStatusStopped),
				),
			},
		},
	})
}

func TestAcc_OpenflowConnector_VersionTrigger(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	runtimeId := testClient().OpenflowRuntime.ActiveRuntime(t)
	definition := testClient().OpenflowConnectorDefinition.ForTesting(t)

	stage, stageCleanup := testClient().Stage.CreateStageInSchema(t, runtimeId.SchemaId())
	t.Cleanup(stageCleanup)

	// The two bundles differ only in whitespace: Snowflake checks that config.json names the connector's
	// definition, and nothing else in it matters for the version to be added.
	bundleV1 := fmt.Sprintf(`{"configFormatVersion":1,"connectorDefinitionId":%q,"configuration":[]}`, definition.Name)
	bundleV2 := fmt.Sprintf(`{"configFormatVersion": 1, "connectorDefinitionId": %q, "configuration": []}`, definition.Name)
	testClient().Stage.PutOnStageWithContent(t, stage.ID(), "config.json", bundleV1)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())
	withV1 := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), "", runtimeId.FullyQualifiedName()).
		WithFromStage(stage.ID(), "").
		WithVersionTrigger("v1")
	withV2 := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), "", runtimeId.FullyQualifiedName()).
		WithFromStage(stage.ID(), "").
		WithVersionTrigger("v2")
	withoutTrigger := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), "", runtimeId.FullyQualifiedName()).
		WithFromStage(stage.ID(), "")
	ref := withV1.ResourceReference()

	// The version a step leaves behind, so the next step can tell whether one was added.
	var lastVersion string
	recordLastVersion := resource.TestCheckResourceAttrWith(ref, "describe_output.0.last_version_name", func(value string) error {
		lastVersion = value
		return nil
	})
	lastVersionChanged := resource.TestCheckResourceAttrWith(ref, "describe_output.0.last_version_name", func(value string) error {
		if value == "" || value == lastVersion {
			return fmt.Errorf("expected a new last version, got %q after %q", value, lastVersion)
		}
		lastVersion = value
		return nil
	})
	lastVersionUnchanged := resource.TestCheckResourceAttrWith(ref, "describe_output.0.last_version_name", func(value string) error {
		if value != lastVersion {
			return fmt.Errorf("expected the last version to stay %q, got %q", lastVersion, value)
		}
		return nil
	})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowConnector),
		Steps: []resource.TestStep{
			// 1. Create without a trigger. The bundle is read by CREATE.
			{
				Config: config.FromModels(t, withoutTrigger),
				Check: assertThat(
					t,
					resourceassert.OpenflowConnectorResource(t, ref).
						HasNoVersionTrigger(),
					resourceshowoutputassert.OpenflowConnectorShowOutput(t, ref).
						HasStatus(sdk.OpenflowConnectorStatusStopped),
					assert.Check(recordLastVersion),
				),
			},
			// 2. Adding the trigger updates state only and adds no version.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withV1),
				Check: assertThat(
					t,
					resourceassert.OpenflowConnectorResource(t, ref).
						HasVersionTriggerString("v1"),
					assert.Check(lastVersionUnchanged),
				),
			},
			// 3. Change the bundle and the trigger. The same connector takes the bundle as a new version, in
			// place, and a stopped connector stays stopped.
			{
				PreConfig: func() {
					testClient().Stage.PutOnStageWithContent(t, stage.ID(), "config.json", bundleV2)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue(ref, tfjsonpath.New("describe_output")),
					},
				},
				Config: config.FromModels(t, withV2),
				Check: assertThat(
					t,
					resourceassert.OpenflowConnectorResource(t, ref).
						HasNameString(id.Name()).
						HasVersionTriggerString("v2"),
					resourceshowoutputassert.OpenflowConnectorShowOutput(t, ref).
						HasStatus(sdk.OpenflowConnectorStatusStopped),
					assert.Check(lastVersionChanged),
				),
			},
			// 4. Removing the trigger updates state only and adds no version.
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withoutTrigger),
				Check: assertThat(
					t,
					resourceassert.OpenflowConnectorResource(t, ref).
						HasVersionTriggerEmpty(),
					assert.Check(lastVersionUnchanged),
				),
			},
		},
	})
}

// Plan-only, so nothing is created and this runs in seconds.
func TestAcc_OpenflowConnector_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	runtimeId := testClient().Ids.RandomSchemaObjectIdentifier()
	stageId := testClient().Ids.RandomSchemaObjectIdentifier()

	invalidRuntime := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), "OPENFLOW_POSTGRES_CDC", "not-a-valid-identifier")

	// A connector comes either from a definition or from a bundle on a stage, never both and never neither.
	bothSources := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), "", runtimeId.FullyQualifiedName()).
		WithFromValue(tfconfig.ListVariable(tfconfig.MapVariable(map[string]tfconfig.Variable{
			"definition": tfconfig.StringVariable("OPENFLOW_POSTGRES_CDC"),
			"stage":      tfconfig.StringVariable(stageId.FullyQualifiedName()),
		})))
	noFrom := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), "", runtimeId.FullyQualifiedName()).
		WithFromValue(nil)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.OpenflowConnector),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, invalidRuntime),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Invalid identifier type`),
			},
			{
				Config:      config.FromModels(t, bothSources),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`only one of .from.0.definition,from.0.stage.`),
			},
			{
				// Only caught because the block is required. ExactlyOneOf on the two sources does not fire
				// when the block is absent, so this would otherwise plan a create and fail against Snowflake.
				Config:      config.FromModels(t, noFrom),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`Insufficient from blocks`),
			},
		},
	})
}
