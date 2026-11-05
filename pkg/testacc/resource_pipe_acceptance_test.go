//go:build non_account_level_tests

package testacc

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/previewfeatures"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_Pipe_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	statement := fmt.Sprintf("COPY INTO %s FROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())

	currentRole := testClient().Context.CurrentRole(t).Name()

	basic := model.PipeWithId("test", id, statement)

	complete := model.PipeWithId("test", id, statement).
		WithComment("initial").
		WithAutoIngest(false)

	updated := model.PipeWithId("test", id, statement).
		WithComment("updated")

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasNoErrorIntegration(),
		resourceshowoutputassert.PipeShowOutput(t, ref).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasDefinition(statement).
			HasOwner(currentRole).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
		resourceshowoutputassert.PipeDescribeOutput(t, ref).
			HasName(id.Name()).
			HasDefinition(statement).
			HasOwner(currentRole),
		resourceparametersassert.PipeResourceParameters(t, ref).
			HasLogEventLevel(sdk.LogLevelOff).
			HasPipeExecutionPaused(false).
			HasLogEventLevelLevel("").
			HasPipeExecutionPausedLevel(""),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
	}

	initialCommentAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("initial").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasNoErrorIntegration(),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("initial").
			HasErrorIntegration("").
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("initial").
			HasErrorIntegration("").
			HasNotificationChannel(""),
	}

	updatedCommentAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("updated").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasNoErrorIntegration(),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("updated").
			HasErrorIntegration("").
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("updated").
			HasErrorIntegration("").
			HasNotificationChannel(""),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Pipe),
		Steps: []resource.TestStep{
			// Create - no options
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Import - no options
			{
				Config:            config.FromModels(t, basic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// integration and error_integration are imported as empty strings when they are unset in Snowflake
				ImportStateVerifyIgnore: []string{"integration", "error_integration"},
			},
			// Update - set comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, initialCommentAssertions...),
			},
			// Import - with comment
			{
				Config:            config.FromModels(t, complete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// integration and error_integration are imported as empty strings when they are unset in Snowflake
				ImportStateVerifyIgnore: []string{"integration", "error_integration"},
			},
			// Update - change comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, updated),
				Check:  assertThat(t, updatedCommentAssertions...),
			},
			// Update - unset comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Update - external comment change
			{
				PreConfig: func() {
					require.NoError(t, atc.defaultTestEnv.client.Pipes.Alter(
						context.Background(),
						sdk.NewAlterPipeRequest(id).WithSet(
							*sdk.NewPipeSetRequest().
								WithComment("external"),
						),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// Create - recreate dropped pipe
			{
				PreConfig: func() {
					require.NoError(t, atc.defaultTestEnv.client.Pipes.Drop(
						context.Background(),
						sdk.NewDropPipeRequest(id),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// No change
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, basic),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider()),
			},
		},
	})
}

func TestAcc_Pipe_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	statement := fmt.Sprintf("COPY INTO %s FROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())

	integration1, integration1Cleanup := testClient().NotificationIntegration.CreateWithAmazonPush(t)
	t.Cleanup(integration1Cleanup)
	integration2, integration2Cleanup := testClient().NotificationIntegration.CreateWithAmazonPush(t)
	t.Cleanup(integration2Cleanup)

	currentRole := testClient().Context.CurrentRole(t).Name()

	complete := model.PipeWithId("test", id, statement).
		WithComment("complete").
		WithAutoIngest(false).
		WithErrorIntegration(integration1.Name)

	commentUpdated := model.PipeWithId("test", id, statement).
		WithComment("updated").
		WithErrorIntegration(integration1.Name)

	updated := model.PipeWithId("test", id, statement).
		WithComment("updated").
		WithErrorIntegration(integration2.Name)

	basic := model.PipeWithId("test", id, statement)

	setIntegration := model.PipeWithId("test", id, statement).
		WithErrorIntegration(integration1.Name)

	ref := basic.ResourceReference()

	completeAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("complete").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasErrorIntegrationString(integration1.Name),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("complete").
			HasErrorIntegration(integration1.Name).
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("complete").
			HasErrorIntegration(integration1.Name).
			HasNotificationChannel(""),
	}

	commentUpdatedAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("updated").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasErrorIntegrationString(integration1.Name),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("updated").
			HasErrorIntegration(integration1.Name).
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("updated").
			HasErrorIntegration(integration1.Name).
			HasNotificationChannel(""),
	}

	updatedAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("updated").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasErrorIntegrationString(integration2.Name),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("updated").
			HasErrorIntegration(integration2.Name).
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("updated").
			HasErrorIntegration(integration2.Name).
			HasNotificationChannel(""),
	}

	// Unsetting error_integration leaves an empty string in the state, while a pipe created without it has no value at all.
	basicAfterUnsetAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasErrorIntegrationString(""),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
	}

	setIntegrationAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasErrorIntegrationString(integration1.Name),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration(integration1.Name).
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration(integration1.Name).
			HasNotificationChannel(""),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Pipe),
		Steps: []resource.TestStep{
			// Create - with comment and error integration
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Import - with comment and error integration
			{
				Config:            config.FromModels(t, complete),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// integration and error_integration are imported as empty strings when they are unset in Snowflake
				ImportStateVerifyIgnore: []string{"integration", "error_integration"},
			},
			// Update - change comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, commentUpdated),
				Check:  assertThat(t, commentUpdatedAssertions...),
			},
			// Update - replace error integration
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, updated),
				Check:  assertThat(t, updatedAssertions...),
			},
			// Import - with replaced error integration
			{
				Config:            config.FromModels(t, updated),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// integration and error_integration are imported as empty strings when they are unset in Snowflake
				ImportStateVerifyIgnore: []string{"integration", "error_integration"},
			},
			// Update - unset comment and error integration
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAfterUnsetAssertions...),
			},
			// Update - set error integration
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, setIntegration),
				Check:  assertThat(t, setIntegrationAssertions...),
			},
			// Update - unset error integration
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAfterUnsetAssertions...),
			},
			// Update - external comment and error integration
			{
				PreConfig: func() {
					require.NoError(t, atc.defaultTestEnv.client.Pipes.Alter(
						context.Background(),
						sdk.NewAlterPipeRequest(id).WithSet(
							*sdk.NewPipeSetRequest().
								WithComment("external").
								WithErrorIntegration(integration2.Name),
						),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// Create - recreate dropped pipe
			{
				PreConfig: func() {
					require.NoError(t, atc.defaultTestEnv.client.Pipes.Drop(
						context.Background(),
						sdk.NewDropPipeRequest(id),
					))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, complete),
				Check:  assertThat(t, completeAssertions...),
			},
			// No change
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, complete),
			},
			// Destroy
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroy),
					},
				},
				Config: config.FromModels(t, providermodel.SnowflakeProvider()),
			},
		},
	})
}

func TestAcc_Pipe_Replacement(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	renamedId := testClient().Ids.RandomSchemaObjectIdentifier()

	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)
	stage, stageCleanup := testClient().Stage.CreateStageWithURL(t)
	t.Cleanup(stageCleanup)
	statement := fmt.Sprintf("COPY INTO %s FROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())
	changedStatement := statement + " ON_ERROR = CONTINUE"

	currentRole := testClient().Context.CurrentRole(t).Name()

	basic := model.PipeWithId("test", id, statement)

	renamed := model.PipeWithId("test", renamedId, statement)

	changed := model.PipeWithId("test", renamedId, changedStatement)

	automatic := model.PipeWithId("test", id, statement).
		WithAutoIngest(true).
		WithComment("auto ingest")

	replacedAutomatic := model.PipeWithId("test", renamedId, changedStatement).
		WithAutoIngest(true).
		WithComment("auto ingest")

	ref := basic.ResourceReference()

	renamedAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, renamedId, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasNoErrorIntegration(),
		newPipeObjectAssertions(t, renamedId).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, renamedId).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
	}

	changedAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, renamedId, currentRole).
			HasCopyStatementString(changedStatement).
			HasCommentString("").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasNoErrorIntegration(),
		newPipeObjectAssertions(t, renamedId).
			HasDefinition(changedStatement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, renamedId).
			HasDefinition(changedStatement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
	}

	automaticAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasAutoIngest(true).
			HasNotificationChannelNotEmpty().
			HasCommentString("auto ingest"),
		newPipeObjectAssertions(t, id).
			HasNotificationChannelNotEmpty().
			HasComment("auto ingest"),
	}

	manualAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasNoErrorIntegration(),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
	}

	replacedAutomaticAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, renamedId, currentRole).
			HasAutoIngest(true).
			HasNotificationChannelNotEmpty(),
		newPipeObjectAssertions(t, renamedId).
			HasNotificationChannelNotEmpty(),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Pipe),
		Steps: []resource.TestStep{
			// Create - with auto ingest
			{
				Config: config.FromModels(t, automatic),
				Check:  assertThat(t, automaticAssertions...),
			},
			// Import - with auto ingest
			{
				Config:            config.FromModels(t, automatic),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// integration and error_integration are imported as empty strings when they are unset in Snowflake
				ImportStateVerifyIgnore: []string{"integration", "error_integration"},
			},
			// Replace - unset auto ingest
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, manualAssertions...),
			},
			// Replace - rename
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, renamed),
				Check:  assertThat(t, renamedAssertions...),
			},
			// Replace - change copy statement
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, changed),
				Check:  assertThat(t, changedAssertions...),
			},
			// Replace - set auto ingest
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, replacedAutomatic),
				Check:  assertThat(t, replacedAutomaticAssertions...),
			},
		},
	})
}

func TestAcc_Pipe_CopyStatementFormattingAndEmptyComment(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	statement := fmt.Sprintf("COPY INTO %s FROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())

	currentRole := testClient().Context.CurrentRole(t).Name()

	basic := model.PipeWithId("test", id, statement)

	formatted := model.PipeWithId("test", id, statement).
		WithCopyStatementMultiline("\n" + statement + ";\n").
		WithComment("").
		WithAutoIngest(false)

	ref := basic.ResourceReference()

	basicAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasNoErrorIntegration(),
		newPipeObjectAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
		newPipeDetailsAssertions(t, id).
			HasDefinition(statement).
			HasComment("").
			HasErrorIntegration("").
			HasNotificationChannel(""),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Pipe),
		Steps: []resource.TestStep{
			// Create - no options
			{
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
			},
			// No change - formatted copy statement and empty comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, formatted),
			},
			// Import - formatted copy statement
			{
				Config:            config.FromModels(t, formatted),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// integration and error_integration are imported as empty strings when they are unset in Snowflake
				ImportStateVerifyIgnore: []string{"integration", "error_integration"},
			},
		},
	})
}

// Each call returns a new assert. Chain step-specific checks on the returned value.
func newPipeResourceAssertions(t *testing.T, ref string, id sdk.SchemaObjectIdentifier, currentRole string) *resourceassert.PipeResourceAssert {
	t.Helper()
	return resourceassert.PipeResource(t, ref).
		HasNameString(id.Name()).
		HasDatabaseString(id.DatabaseName()).
		HasSchemaString(id.SchemaName()).
		HasFullyQualifiedNameString(id.FullyQualifiedName()).
		HasOwnerString(currentRole)
}

func newPipeObjectAssertions(t *testing.T, id sdk.SchemaObjectIdentifier) *objectassert.PipeAssert {
	t.Helper()
	return objectassert.Pipe(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasIntegration("").
		HasOwnerRoleType("ROLE").
		HasInvalidReason("")
}

func newPipeDetailsAssertions(t *testing.T, id sdk.SchemaObjectIdentifier) *objectassert.PipeDetailsAssert {
	t.Helper()
	return objectassert.PipeDetails(t, id).
		HasCreatedOnNotEmpty().
		HasName(id.Name()).
		HasDatabaseName(id.DatabaseName()).
		HasSchemaName(id.SchemaName()).
		HasIntegration("").
		HasOwnerRoleType("ROLE").
		HasInvalidReason("")
}

func TestAcc_Pipe_Parameters(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	statement := fmt.Sprintf("COPY INTO %s FROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())

	basic := model.PipeWithId("test", id, statement)
	withParameters := model.PipeWithId("test", id, statement).
		WithLogEventLevelSdk(sdk.LogLevelInfo).
		WithPipeExecutionPaused(true)
	updatedParameters := model.PipeWithId("test", id, statement).
		WithLogEventLevelSdk(sdk.LogLevelWarn).
		WithPipeExecutionPaused(false)
	withPipeExecutionPausedFalse := model.PipeWithId("test", id, statement).
		WithPipeExecutionPaused(false)

	ref := basic.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Pipe),
		Steps: []resource.TestStep{
			// Create - with parameters
			{
				Config: config.FromModels(t, withParameters),
				Check: assertThat(t,
					resourceparametersassert.PipeResourceParameters(t, ref).
						HasLogEventLevel(sdk.LogLevelInfo).
						HasLogEventLevelLevel(sdk.ParameterTypePipe).
						HasPipeExecutionPaused(true).
						HasPipeExecutionPausedLevel(sdk.ParameterTypePipe),
					objectparametersassert.PipeParameters(t, id).
						HasLogEventLevel(sdk.LogLevelInfo).
						HasLogEventLevelLevel(sdk.ParameterTypePipe).
						HasPipeExecutionPaused(true).
						HasPipeExecutionPausedLevel(sdk.ParameterTypePipe),
				),
			},
			// Import
			{
				Config:            config.FromModels(t, withParameters),
				ResourceName:      ref,
				ImportState:       true,
				ImportStateVerify: true,
				// integration and error_integration are imported as empty strings when they are unset in Snowflake
				ImportStateVerifyIgnore: []string{"integration", "error_integration"},
			},
			// Update - change parameters
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, updatedParameters),
				Check: assertThat(t,
					resourceparametersassert.PipeResourceParameters(t, ref).
						HasLogEventLevel(sdk.LogLevelWarn).
						HasPipeExecutionPaused(false),
					objectparametersassert.PipeParameters(t, id).
						HasLogEventLevel(sdk.LogLevelWarn).
						HasPipeExecutionPaused(false),
				),
			},
			// Update - unset parameters
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check: assertThat(t,
					resourceparametersassert.PipeResourceParameters(t, ref).
						HasLogEventLevel(sdk.LogLevelOff).
						HasLogEventLevelLevel("").
						HasPipeExecutionPaused(false).
						HasPipeExecutionPausedLevel(""),
					objectparametersassert.PipeParameters(t, id).
						HasDefaultLogEventLevelValue().
						HasDefaultPipeExecutionPausedValue(),
				),
			},
			// Destroy
			{
				Destroy: true,
				Config:  config.FromModels(t, basic),
			},
			// Create - with pipe_execution_paused set to false, which is the zero value
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, withPipeExecutionPausedFalse),
				Check: assertThat(t,
					resourceparametersassert.PipeResourceParameters(t, ref).
						HasPipeExecutionPaused(false).
						HasPipeExecutionPausedLevel(sdk.ParameterTypePipe),
					objectparametersassert.PipeParameters(t, id).
						HasPipeExecutionPaused(false).
						HasPipeExecutionPausedLevel(sdk.ParameterTypePipe),
				),
			},
		},
	})
}

func TestAcc_Pipe_migrateFromV2_22_1_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	statement := fmt.Sprintf("COPY INTO %s FROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())

	pipeModel := model.PipeWithId("test", id, statement).WithComment("comment")
	ref := pipeModel.ResourceReference()
	providerModel := providermodel.SnowflakeProvider().WithPreviewFeaturesEnabled(string(previewfeatures.PipeResource))

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Pipe),
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.22.1"),
				Config:            config.FromModels(t, providerModel, pipeModel),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(ref, "id", helpers.EncodeSnowflakeID(id)),
				),
			},
			{
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, pipeModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(ref, "id", id.FullyQualifiedName())),
					resourceassert.PipeResource(t, ref).
						HasCommentString("comment"),
					resourceshowoutputassert.PipeShowOutput(t, ref).
						HasName(id.Name()),
				),
			},
		},
	})
}

func TestAcc_Pipe_ExternalRecreateWithDifferentErrorIntegration(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	table, tableCleanup := testClient().Table.Create(t)
	t.Cleanup(tableCleanup)
	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	statement := fmt.Sprintf("COPY INTO %s FROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())

	integration1, integration1Cleanup := testClient().NotificationIntegration.CreateWithAmazonPush(t)
	t.Cleanup(integration1Cleanup)
	integration2, integration2Cleanup := testClient().NotificationIntegration.CreateWithAmazonPush(t)
	t.Cleanup(integration2Cleanup)

	withIntegration1 := model.PipeWithId("test", id, statement).WithErrorIntegration(integration1.Name)
	withIntegration2 := model.PipeWithId("test", id, statement).WithErrorIntegration(integration2.Name)
	ref := withIntegration1.ResourceReference()

	recreateWithErrorIntegration := func(errorIntegration string) func() {
		return func() {
			require.NoError(t, atc.defaultTestEnv.client.Pipes.Drop(context.Background(), sdk.NewDropPipeRequest(id)))
			require.NoError(t, atc.defaultTestEnv.client.Pipes.Create(
				context.Background(),
				sdk.NewCreatePipeRequest(id, statement).WithErrorIntegration(errorIntegration),
			))
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             CheckDestroy(t, resources.Pipe),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, withIntegration1),
				Check:  assertThat(t, objectassert.Pipe(t, id).HasErrorIntegration(integration1.Name)),
			},
			// The pipe is recreated outside of Terraform with a different error integration: the configuration is restored in place.
			{
				PreConfig: recreateWithErrorIntegration(integration2.Name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, withIntegration1),
				Check:  assertThat(t, objectassert.Pipe(t, id).HasErrorIntegration(integration1.Name)),
			},
			// The pipe is recreated outside of Terraform again, and the configuration is adjusted to match: no changes.
			{
				PreConfig: recreateWithErrorIntegration(integration2.Name),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionNoop),
					},
				},
				Config: config.FromModels(t, withIntegration2),
				Check:  assertThat(t, objectassert.Pipe(t, id).HasErrorIntegration(integration2.Name)),
			},
		},
	})
}
