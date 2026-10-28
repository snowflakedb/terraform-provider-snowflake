//go:build non_account_level_tests

package testacc

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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

	initialCommentAssertions := []assert.TestCheckFuncProvider{
		newPipeResourceAssertions(t, ref, id, currentRole).
			HasCopyStatementString(statement).
			HasCommentString("initial").
			HasAutoIngest(false).
			HasNotificationChannelEmpty().
			HasErrorIntegrationString(""),
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
			HasErrorIntegrationString(""),
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

	basicAssertions := []assert.TestCheckFuncProvider{
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
			},
			// Update - unset comment and error integration
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(ref, plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basic),
				Check:  assertThat(t, basicAssertions...),
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
				Check:  assertThat(t, basicAssertions...),
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
			HasErrorIntegrationString(""),
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
			HasErrorIntegrationString(""),
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
