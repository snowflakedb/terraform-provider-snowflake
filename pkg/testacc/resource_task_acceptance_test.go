//go:build !account_level_tests

package testacc

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	configvariable "github.com/hashicorp/terraform-plugin-testing/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO(SNOW-1822118): Create more complicated tests for task

func TestAcc_Task_Basic(t *testing.T) {
	currentRole := testClient().Context.CurrentRole(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	configModel := model.TaskWithId("test", id, false, statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfter().
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment("").
						HasWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasScheduleEmpty().
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")).
						HasLastCommittedOn("").
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			{
				ResourceName: configModel.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedTaskResource(t, helpers.EncodeResourceIdentifier(id)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanFalse).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfterEmpty().
						HasWhenString("").
						HasSqlStatementString(statement),
				),
			},
		},
	})
}

func TestAcc_Task_Complete(t *testing.T) {
	currentRole := testClient().Context.CurrentRole(t)

	errorNotificationIntegration, errorNotificationIntegrationCleanup := testClient().NotificationIntegration.CreateWithGcpPubSub(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	taskConfig := `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`
	comment := random.Comment()
	condition := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
	configModel := model.TaskWithId("test", id, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(10).
		WithConfigValue(configvariable.StringVariable(taskConfig)).
		WithAllowOverlappingExecution(r.BooleanTrue).
		WithErrorIntegration(errorNotificationIntegration.ID().Name()).
		WithComment(comment).
		WithWhen(condition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(testClient().Ids.WarehouseId().Name()).
						HasScheduleMinutes(10).
						HasConfigString(taskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfterEmpty().
						HasWhenString(condition).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment(comment).
						HasWarehouse(testClient().Ids.WarehouseId()).
						HasScheduleMinutes(10).
						HasPredecessors().
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID()).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig(taskConfig).
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			{
				ResourceName:    configModel.ResourceReference(),
				ImportState:     true,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedTaskResource(t, helpers.EncodeResourceIdentifier(id)).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(testClient().Ids.WarehouseId().Name()).
						HasScheduleMinutes(10).
						HasConfigString(taskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfterEmpty().
						HasWhenString(condition).
						HasSqlStatementString(statement),
				),
			},
		},
	})
}

func TestAcc_Task_Updates(t *testing.T) {
	currentRole := testClient().Context.CurrentRole(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	newStatement := "SELECT 123"
	basicConfigModel := model.TaskWithId("test", id, false, statement)

	// TODO(SNOW-1736173): New warehouse created, because the common one has lower-case letters that won't work
	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	errorNotificationIntegration, errorNotificationIntegrationCleanup := testClient().NotificationIntegration.CreateWithGcpPubSub(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	taskConfig := `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`
	comment := random.Comment()
	condition := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
	completeConfigModel := model.TaskWithId("test", id, true, newStatement).
		WithWarehouse(warehouse.ID().Name()).
		WithScheduleMinutes(5).
		WithConfigValue(configvariable.StringVariable(taskConfig)).
		WithAllowOverlappingExecution(r.BooleanTrue).
		WithErrorIntegration(errorNotificationIntegration.ID().Name()).
		WithComment(comment).
		WithWhen(condition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, basicConfigModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, basicConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfter().
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, basicConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment("").
						HasWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasScheduleEmpty().
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")).
						HasLastCommittedOn("").
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// Set
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, completeConfigModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(completeConfigModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.TaskResource(t, completeConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(warehouse.ID().Name()).
						HasScheduleMinutes(5).
						HasConfigString(taskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfter().
						HasWhenString(condition).
						HasSqlStatementString(newStatement),
					resourceshowoutputassert.TaskShowOutput(t, completeConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasWarehouse(warehouse.ID()).
						HasComment(comment).
						HasScheduleMinutes(5).
						HasPredecessors().
						HasState(sdk.TaskStateStarted).
						HasDefinition(newStatement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID()).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig(taskConfig).
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// Unset
			{
				Config: config.FromModels(t, basicConfigModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicConfigModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.TaskResource(t, basicConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfter().
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, basicConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment("").
						HasWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasScheduleEmpty().
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOnNotEmpty().
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
		},
	})
}

/*
DAG structure (the test proves child3 won't have any issues with updates in the following scenario):

		 child1
		/		\
	 root 		 child3
		\		/
		 child2
*/
func TestAcc_Task_UpdatesInComplexDAG(t *testing.T) {
	rootTask, rootTaskCleanup := testClient().Task.CreateWithSchedule(t)
	t.Cleanup(rootTaskCleanup)

	child1, child1Cleanup := testClient().Task.CreateWithAfter(t, rootTask.ID())
	t.Cleanup(child1Cleanup)

	child2, child2Cleanup := testClient().Task.CreateWithAfter(t, rootTask.ID())
	t.Cleanup(child2Cleanup)

	testClient().Task.Alter(t, sdk.NewAlterTaskRequest(child1.ID()).WithResume(true))
	testClient().Task.Alter(t, sdk.NewAlterTaskRequest(child2.ID()).WithResume(true))
	testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootTask.ID()).WithResume(true))
	t.Cleanup(func() { testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootTask.ID()).WithSuspend(true)) })

	child3Id := testClient().Ids.RandomSchemaObjectIdentifier()
	basicConfigModel := model.TaskWithId("test", child3Id, true, "SELECT 1").
		WithAfterValue(configvariable.SetVariable(
			configvariable.StringVariable(child1.ID().FullyQualifiedName()),
			configvariable.StringVariable(child2.ID().FullyQualifiedName()),
		))

	comment := random.Comment()
	basicConfigModelAfterUpdate := model.TaskWithId("test", child3Id, true, "SELECT 123").
		WithAfterValue(configvariable.SetVariable(
			configvariable.StringVariable(child1.ID().FullyQualifiedName()),
			configvariable.StringVariable(child2.ID().FullyQualifiedName()),
		)).
		WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, basicConfigModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, basicConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(child3Id.FullyQualifiedName()).
						HasDatabaseString(child3Id.DatabaseName()).
						HasSchemaString(child3Id.SchemaName()).
						HasNameString(child3Id.Name()).
						HasStartedString(r.BooleanTrue).
						HasAfter(child1.ID(), child2.ID()).
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, basicConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(child3Id.Name()).
						HasDatabaseName(child3Id.DatabaseName()).
						HasSchemaName(child3Id.SchemaName()).
						HasState(sdk.TaskStateStarted).
						HasDefinition("SELECT 1"),
				),
			},
			// Update some fields in child3
			{
				Config: config.FromModels(t, basicConfigModelAfterUpdate),
				Check: assertThat(t,
					resourceassert.TaskResource(t, basicConfigModelAfterUpdate.ResourceReference()).
						HasFullyQualifiedNameString(child3Id.FullyQualifiedName()).
						HasDatabaseString(child3Id.DatabaseName()).
						HasSchemaString(child3Id.SchemaName()).
						HasNameString(child3Id.Name()).
						HasStartedString(r.BooleanTrue).
						HasCommentString(comment).
						HasAfter(child1.ID(), child2.ID()).
						HasSqlStatementString("SELECT 123"),
					resourceshowoutputassert.TaskShowOutput(t, basicConfigModelAfterUpdate.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(child3Id.Name()).
						HasDatabaseName(child3Id.DatabaseName()).
						HasSchemaName(child3Id.SchemaName()).
						HasState(sdk.TaskStateStarted).
						HasComment(comment).
						HasDefinition("SELECT 123"),
				),
			},
		},
	})
}

func TestAcc_Task_StatementSpaces(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	when := "1 > 2"
	configModel := model.TaskWithId("test", id, false, statement).WithWhen(when)

	statementWithSpaces := "    SELECT    1    "
	whenWithSpaces := "     1      >       2      "
	configModelWithSpacesInStatements := model.TaskWithId("test", id, false, statementWithSpaces).WithWhen(whenWithSpaces)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasWhenString(when).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasName(id.Name()).
						HasCondition(when).
						HasDefinition(statement),
				),
			},
			{
				Config: config.FromModels(t, configModelWithSpacesInStatements),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasWhenString(when).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasName(id.Name()).
						HasCondition(when).
						HasDefinition(statement),
				),
			},
		},
	})
}

func TestAcc_Task_ExternalChanges(t *testing.T) {
	currentRole := testClient().Context.CurrentRole(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	basicConfigModel := model.TaskWithId("test", id, false, statement)

	// TODO(SNOW-1736173): New warehouse created, because the common one has lower-case letters that won't work
	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	errorNotificationIntegration, errorNotificationIntegrationCleanup := testClient().NotificationIntegration.CreateWithGcpPubSub(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	taskConfig := `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`
	comment := random.Comment()
	condition := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
	completeConfigModel := model.TaskWithId("test", id, true, statement).
		WithWarehouse(warehouse.ID().Name()).
		WithScheduleMinutes(5).
		WithConfigValue(configvariable.StringVariable(taskConfig)).
		WithAllowOverlappingExecution(r.BooleanTrue).
		WithErrorIntegration(errorNotificationIntegration.ID().Name()).
		WithComment(comment).
		WithWhen(condition)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			// Optionals set
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, completeConfigModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, completeConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(warehouse.ID().Name()).
						HasScheduleMinutes(5).
						HasConfigString(taskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfter().
						HasWhenString(condition).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, completeConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasWarehouse(warehouse.ID()).
						HasComment(comment).
						HasScheduleMinutes(5).
						HasPredecessors().
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID()).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig(taskConfig).
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// External change - unset all optional fields and expect no change
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithSuspend(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithUnset(*sdk.NewTaskUnsetRequest().
						WithWarehouse(true).
						WithConfig(true).
						WithAllowOverlappingExecution(true).
						WithErrorIntegration(true).
						WithComment(true).
						WithSchedule(true),
					))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithRemoveWhen(true))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicConfigModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, completeConfigModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, completeConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasWarehouseString(warehouse.ID().Name()).
						HasScheduleMinutes(5).
						HasConfigString(taskConfig).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()).
						HasCommentString(comment).
						HasFinalizeString("").
						HasAfter().
						HasWhenString(condition).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, completeConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasWarehouse(warehouse.ID()).
						HasComment(comment).
						HasScheduleMinutes(5).
						HasPredecessors().
						HasState(sdk.TaskStateStarted).
						HasDefinition(statement).
						HasCondition(condition).
						HasAllowOverlappingExecution(true).
						HasErrorIntegration(errorNotificationIntegration.ID()).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOnNotEmpty().
						HasOwnerRoleType("ROLE").
						HasConfig(taskConfig).
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// Unset optional values
			{
				Config: config.FromModels(t, basicConfigModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, basicConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfter().
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, basicConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment("").
						HasWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasScheduleEmpty().
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOnNotEmpty().
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// External change - set all optional fields and expect no change
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithSuspend(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithSet(*sdk.NewTaskSetRequest().
						WithWarehouse(warehouse.ID()).
						WithConfig(taskConfig).
						WithAllowOverlappingExecution(true).
						WithErrorIntegration(errorNotificationIntegration.ID()).
						WithComment(comment).
						WithSchedule("5 MINUTE"),
					))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithModifyWhen(condition))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithModifyAs("SELECT 123"))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(basicConfigModel.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: config.FromModels(t, basicConfigModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, basicConfigModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasWarehouseString("").
						HasNoScheduleSet().
						HasConfigString("").
						HasAllowOverlappingExecutionString(r.BooleanDefault).
						HasErrorIntegrationString("").
						HasCommentString("").
						HasFinalizeString("").
						HasAfter().
						HasWhenString("").
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, basicConfigModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasIdNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(currentRole.Name()).
						HasComment("").
						HasWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasScheduleEmpty().
						HasPredecessors().
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement).
						HasCondition("").
						HasAllowOverlappingExecution(false).
						HasErrorIntegration(sdk.NewAccountObjectIdentifier("")).
						HasLastCommittedOnNotEmpty().
						HasLastSuspendedOnNotEmpty().
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
		},
	})
}

func TestAcc_Task_CallingProcedure(t *testing.T) {
	procedure, cleanupProcedure := testClient().Procedure.Create(t, sdk.DataTypeNumber)
	t.Cleanup(cleanupProcedure)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := fmt.Sprintf("call %s(123)", procedure.Name)

	configModel := model.TaskWithId("test", id, false, statement).WithUserTaskManagedInitialWarehouseSizeEnum(sdk.WarehouseSizeXSmall)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasUserTaskManagedInitialWarehouseSizeEnum(sdk.WarehouseSizeXSmall).
						HasSqlStatementString(statement),
					resourceshowoutputassert.TaskShowOutput(t, configModel.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasName(id.Name()).
						HasState(sdk.TaskStateSuspended).
						HasDefinition(statement),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall),
				),
			},
		},
	})
}

func TestAcc_Task_CronAndMinutes(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	minutes := 5
	cron := "*/5 * * * * UTC"

	configModelWithoutSchedule := model.TaskWithId("test", id, false, "SELECT 1")
	configModelWithMinutes := model.TaskWithId("test", id, true, "SELECT 1").WithScheduleMinutes(minutes)
	configModelWithCron := model.TaskWithId("test", id, true, "SELECT 1").WithScheduleCron(cron)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			// create with minutes
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithMinutes),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithMinutes.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(minutes).
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithMinutes.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleMinutes(minutes),
				),
			},
			// Unset schedule (from minutes)
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithoutSchedule),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithoutSchedule.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasNoScheduleSet().
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithoutSchedule.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleEmpty(),
				),
			},
			// Create with cron
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithCron),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithCron.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasScheduleCron(cron).
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithCron.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleCron(cron),
				),
			},
			// Change to minutes
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithMinutes),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithMinutes.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(minutes).
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithMinutes.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleMinutes(minutes),
				),
			},
			// Change back to cron
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithCron),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithCron.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanTrue).
						HasScheduleCron(cron).
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithCron.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleCron(cron),
				),
			},
			// Unset schedule (from cron)
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithoutSchedule),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithoutSchedule.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasNoScheduleSet().
						HasSqlStatementString("SELECT 1"),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithoutSchedule.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleEmpty(),
				),
			},
		},
	})
}

func TestAcc_Task_CronAndMinutes_ExternalChanges(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	minutes := 5
	cron := "*/5 * * * * UTC"

	configModelWithoutSchedule := model.TaskWithId("test", id, false, "SELECT 1")
	configModelWithMinutes := model.TaskWithId("test", id, false, "SELECT 1").WithScheduleMinutes(minutes)
	configModelWithCron := model.TaskWithId("test", id, false, "SELECT 1").WithScheduleCron(cron)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			// Create without a schedule
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithoutSchedule),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithoutSchedule.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasNoScheduleSet(),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithoutSchedule.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleEmpty(),
				),
			},
			// External change - set minutes
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithSet(*sdk.NewTaskSetRequest().WithSchedule("5 MINUTES")))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithoutSchedule),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithoutSchedule.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasNoScheduleSet(),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithoutSchedule.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleEmpty(),
				),
			},
			// External change - set cron
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithSet(*sdk.NewTaskSetRequest().WithSchedule(fmt.Sprintf("USING CRON %s", cron))))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithoutSchedule),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithoutSchedule.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasNoScheduleSet(),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithoutSchedule.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleEmpty(),
				),
			},
			// Set minutes schedule
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithMinutes),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithMinutes.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasScheduleMinutes(minutes),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithMinutes.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleMinutes(minutes),
				),
			},
			// External change - unset schedule
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(id).WithUnset(*sdk.NewTaskUnsetRequest().WithSchedule(true)))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithMinutes),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithMinutes.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasScheduleMinutes(minutes),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithMinutes.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleMinutes(minutes),
				),
			},
			// Set cron schedule
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithCron),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithCron.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasScheduleCron(cron),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithCron.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleCron(cron),
				),
			},
			// External change - unset schedule
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithCron),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelWithCron.ResourceReference()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasScheduleCron(cron),
					resourceshowoutputassert.TaskShowOutput(t, configModelWithCron.ResourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasScheduleCron(cron),
				),
			},
		},
	})
}

func TestAcc_Task_ScheduleSchemaValidation(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				Config:      taskConfigInvalidScheduleSetMultipleOrEmpty(id, true),
				ExpectError: regexp.MustCompile("\"schedule.0.minutes\": only one of `schedule.0.minutes,schedule.0.using_cron`"),
			},
			{
				Config:      taskConfigInvalidScheduleSetMultipleOrEmpty(id, false),
				ExpectError: regexp.MustCompile("\"schedule.0.minutes\": one of `schedule.0.minutes,schedule.0.using_cron`"),
			},
		},
	})
}

func taskConfigInvalidScheduleSetMultipleOrEmpty(id sdk.SchemaObjectIdentifier, setMultiple bool) string {
	var scheduleString string
	scheduleBuffer := new(bytes.Buffer)
	scheduleBuffer.WriteString("schedule {\n")
	if setMultiple {
		scheduleBuffer.WriteString("minutes = 10\n")
		scheduleBuffer.WriteString("using_cron = \"*/5 * * * * UTC\"\n")
	}
	scheduleBuffer.WriteString("}\n")
	scheduleString = scheduleBuffer.String()

	return fmt.Sprintf(`
resource "snowflake_task" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	started = false
	sql_statement = "SELECT 1"

	%[4]s
}`, id.DatabaseName(), id.SchemaName(), id.Name(), scheduleString)
}

func TestAcc_Task_AllParameters(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	configModel := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(5)
	configModelWithAllParametersSet := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(5).
		WithSuspendTaskAfterNumFailures(15).
		WithTaskAutoRetryAttempts(15).
		WithUserTaskManagedInitialWarehouseSizeEnum(sdk.WarehouseSizeXSmall).
		WithUserTaskMinimumTriggerIntervalInSeconds(30).
		WithUserTaskTimeoutMs(1000).
		WithAbortDetachedQuery(true).
		WithAutocommit(false).
		WithBinaryInputFormatEnum(sdk.BinaryInputFormatUTF8).
		WithBinaryOutputFormatEnum(sdk.BinaryOutputFormatBase64).
		WithClientMemoryLimit(1024).
		WithClientMetadataRequestUseConnectionCtx(true).
		WithClientPrefetchThreads(2).
		WithClientResultChunkSize(48).
		WithClientResultColumnCaseInsensitive(true).
		WithClientSessionKeepAlive(true).
		WithClientSessionKeepAliveHeartbeatFrequency(2400).
		WithClientTimestampTypeMappingEnum(sdk.ClientTimestampTypeMappingNtz).
		WithDateInputFormat("YYYY-MM-DD").
		WithDateOutputFormat("YY-MM-DD").
		WithEnableUnloadPhysicalTypeOptimization(false).
		WithErrorOnNondeterministicMerge(false).
		WithErrorOnNondeterministicUpdate(true).
		WithGeographyOutputFormatEnum(sdk.GeographyOutputFormatWKB).
		WithGeometryOutputFormatEnum(sdk.GeometryOutputFormatWKB).
		WithJdbcUseSessionTimezone(false).
		WithJsonIndent(4).
		WithLockTimeout(21222).
		WithLogLevelEnum(sdk.LogLevelError).
		WithMultiStatementCount(0).
		WithNoorderSequenceAsDefault(false).
		WithOdbcTreatDecimalAsInt(true).
		WithQueryTag("some_tag").
		WithQuotedIdentifiersIgnoreCase(true).
		WithRowsPerResultset(2).
		WithS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
		WithSearchPath("$public, $current").
		WithStatementQueuedTimeoutInSeconds(10).
		WithStatementTimeoutInSeconds(10).
		WithStrictJsonOutput(true).
		WithTimestampDayIsAlways24h(true).
		WithTimestampInputFormat("YYYY-MM-DD").
		WithTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimestampTypeMappingEnum(sdk.TimestampTypeMappingLtz).
		WithTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
		WithTimezone("Europe/Warsaw").
		WithTimeInputFormat("HH24:MI").
		WithTimeOutputFormat("HH24:MI").
		WithTraceLevelEnum(sdk.TraceLevelPropagate).
		WithTransactionAbortOnError(true).
		WithTransactionDefaultIsolationLevelEnum(sdk.TransactionDefaultIsolationLevelReadCommitted).
		WithTwoDigitCenturyStart(1980).
		WithUnsupportedDdlActionEnum(sdk.UnsupportedDDLActionFail).
		WithUseCachedResult(false).
		WithWeekOfYearPolicy(1).
		WithWeekStart(1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { TestAccPreCheck(t) },
		CheckDestroy: CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			// create with default values for all the parameters
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				Check: assertThat(t,
					objectparametersassert.TaskParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
			// import when no parameter set
			{
				ResourceName:    configModel.ResourceReference(),
				ImportState:     true,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				ImportStateCheck: assertThatImport(t,
					resourceparametersassert.ImportedTaskResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasAllDefaults(),
				),
			},
			// set all parameters
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithAllParametersSet),
				Check: assertThat(t,
					objectparametersassert.TaskParameters(t, id).
						HasSuspendTaskAfterNumFailures(15).
						HasTaskAutoRetryAttempts(15).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall).
						HasUserTaskMinimumTriggerIntervalInSeconds(30).
						HasUserTaskTimeoutMs(1000).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelPropagate).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1),
					resourceparametersassert.TaskResourceParameters(t, configModelWithAllParametersSet.ResourceReference()).
						HasSuspendTaskAfterNumFailures(15).
						HasTaskAutoRetryAttempts(15).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall).
						HasUserTaskMinimumTriggerIntervalInSeconds(30).
						HasUserTaskTimeoutMs(1000).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelPropagate).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1),
				),
			},
			// import when all parameters set
			{
				ResourceName:    configModelWithAllParametersSet.ResourceReference(),
				ImportState:     true,
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelWithAllParametersSet),
				ImportStateCheck: assertThatImport(t,
					resourceparametersassert.ImportedTaskResourceParameters(t, helpers.EncodeResourceIdentifier(id)).
						HasSuspendTaskAfterNumFailures(15).
						HasTaskAutoRetryAttempts(15).
						HasUserTaskManagedInitialWarehouseSize(sdk.WarehouseSizeXSmall).
						HasUserTaskMinimumTriggerIntervalInSeconds(30).
						HasUserTaskTimeoutMs(1000).
						HasAbortDetachedQuery(true).
						HasAutocommit(false).
						HasBinaryInputFormat(sdk.BinaryInputFormatUTF8).
						HasBinaryOutputFormat(sdk.BinaryOutputFormatBase64).
						HasClientMemoryLimit(1024).
						HasClientMetadataRequestUseConnectionCtx(true).
						HasClientPrefetchThreads(2).
						HasClientResultChunkSize(48).
						HasClientResultColumnCaseInsensitive(true).
						HasClientSessionKeepAlive(true).
						HasClientSessionKeepAliveHeartbeatFrequency(2400).
						HasClientTimestampTypeMapping(sdk.ClientTimestampTypeMappingNtz).
						HasDateInputFormat("YYYY-MM-DD").
						HasDateOutputFormat("YY-MM-DD").
						HasEnableUnloadPhysicalTypeOptimization(false).
						HasErrorOnNondeterministicMerge(false).
						HasErrorOnNondeterministicUpdate(true).
						HasGeographyOutputFormat(sdk.GeographyOutputFormatWKB).
						HasGeometryOutputFormat(sdk.GeometryOutputFormatWKB).
						HasJdbcUseSessionTimezone(false).
						HasJsonIndent(4).
						HasLockTimeout(21222).
						HasLogLevel(sdk.LogLevelError).
						HasMultiStatementCount(0).
						HasNoorderSequenceAsDefault(false).
						HasOdbcTreatDecimalAsInt(true).
						HasQueryTag("some_tag").
						HasQuotedIdentifiersIgnoreCase(true).
						HasRowsPerResultset(2).
						HasS3StageVpceDnsName("vpce-id.s3.region.vpce.amazonaws.com").
						HasSearchPath("$public, $current").
						HasStatementQueuedTimeoutInSeconds(10).
						HasStatementTimeoutInSeconds(10).
						HasStrictJsonOutput(true).
						HasTimestampDayIsAlways24h(true).
						HasTimestampInputFormat("YYYY-MM-DD").
						HasTimestampLtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampNtzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimestampTypeMapping(sdk.TimestampTypeMappingLtz).
						HasTimestampTzOutputFormat("YYYY-MM-DD HH24:MI:SS").
						HasTimezone("Europe/Warsaw").
						HasTimeInputFormat("HH24:MI").
						HasTimeOutputFormat("HH24:MI").
						HasTraceLevel(sdk.TraceLevelPropagate).
						HasTransactionAbortOnError(true).
						HasTransactionDefaultIsolationLevel(sdk.TransactionDefaultIsolationLevelReadCommitted).
						HasTwoDigitCenturyStart(1980).
						HasUnsupportedDdlAction(sdk.UnsupportedDDLActionFail).
						HasUseCachedResult(false).
						HasWeekOfYearPolicy(1).
						HasWeekStart(1),
				),
			},
			// unset all the parameters
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModel),
				Check: assertThat(t,
					objectparametersassert.TaskParameters(t, id).
						HasAllDefaults().
						HasAllDefaultsExplicit(),
					resourceparametersassert.TaskResourceParameters(t, configModel.ResourceReference()).
						HasAllDefaults(),
				),
			},
		},
	})
}

func TestAcc_Task_Enabled(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	configModelEnabled := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(5)
	configModelDisabled := model.TaskWithId("test", id, false, statement).
		WithScheduleMinutes(5)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelDisabled),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelDisabled.ResourceReference()).
						HasStartedString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, configModelDisabled.ResourceReference()).
						HasState(sdk.TaskStateSuspended),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelEnabled),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelEnabled.ResourceReference()).
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, configModelEnabled.ResourceReference()).
						HasState(sdk.TaskStateStarted),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, configModelDisabled),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModelDisabled.ResourceReference()).
						HasStartedString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, configModelDisabled.ResourceReference()).
						HasState(sdk.TaskStateSuspended),
				),
			},
		},
	})
}

func TestAcc_Task_ConvertStandaloneTaskToSubtask(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	id2 := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"

	firstTaskStandaloneModel := model.TaskWithId("root", id, true, statement).
		WithScheduleMinutes(5).
		WithSuspendTaskAfterNumFailures(1)
	secondTaskStandaloneModel := model.TaskWithId("child", id2, true, statement).
		WithScheduleMinutes(5)

	rootTaskModel := model.TaskWithId("root", id, true, statement).
		WithScheduleMinutes(5).
		WithSuspendTaskAfterNumFailures(2)
	childTaskModel := model.TaskWithId("child", id2, true, statement).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(id.FullyQualifiedName())))
	childTaskModel.SetDependsOn(rootTaskModel.ResourceReference())

	firstTaskStandaloneModelDisabled := model.TaskWithId("root", id, false, statement).
		WithScheduleMinutes(5)
	secondTaskStandaloneModelDisabled := model.TaskWithId("child", id2, false, statement).
		WithScheduleMinutes(5)
	secondTaskStandaloneModelDisabled.SetDependsOn(firstTaskStandaloneModelDisabled.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", firstTaskStandaloneModel, secondTaskStandaloneModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, firstTaskStandaloneModel.ResourceReference()).
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceshowoutputassert.TaskShowOutput(t, firstTaskStandaloneModel.ResourceReference()).
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, secondTaskStandaloneModel.ResourceReference()).
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, secondTaskStandaloneModel.ResourceReference()).
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change the second task to run after the first one (creating a DAG)
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskModel, childTaskModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskModel.ResourceReference()).
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("2"),
					resourceshowoutputassert.TaskShowOutput(t, rootTaskModel.ResourceReference()).
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, childTaskModel.ResourceReference()).
						HasAfter(id).
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, childTaskModel.ResourceReference()).
						HasPredecessors(id).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change tasks in DAG to standalone tasks (disabled to check if resuming/suspending works correctly)
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", firstTaskStandaloneModelDisabled, secondTaskStandaloneModelDisabled),
				Check: assertThat(t,
					resourceassert.TaskResource(t, firstTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanFalse).
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceshowoutputassert.TaskShowOutput(t, firstTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateSuspended),
					resourceassert.TaskResource(t, secondTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleMinutes(5).
						HasStartedString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, secondTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleMinutes(5).
						HasState(sdk.TaskStateSuspended),
				),
			},
		},
	})
}

func TestAcc_Task_ConvertStandaloneTaskToFinalizer(t *testing.T) {
	rootTaskId := testClient().Ids.RandomSchemaObjectIdentifier()
	finalizerTaskId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	firstTaskStandaloneModel := model.TaskWithId("root", rootTaskId, true, statement).
		WithScheduleMinutes(schedule).
		WithSuspendTaskAfterNumFailures(1)
	secondTaskStandaloneModel := model.TaskWithId("child", finalizerTaskId, true, statement).
		WithScheduleMinutes(schedule)

	rootTaskModel := model.TaskWithId("root", rootTaskId, true, statement).
		WithScheduleMinutes(schedule).
		WithSuspendTaskAfterNumFailures(2)
	childTaskModel := model.TaskWithId("child", finalizerTaskId, true, statement).
		WithFinalize(rootTaskId.FullyQualifiedName())
	childTaskModel.SetDependsOn(rootTaskModel.ResourceReference())

	rootTaskStandaloneModelDisabled := model.TaskWithId("root", rootTaskId, false, statement).
		WithScheduleMinutes(schedule)
	childTaskStandaloneModelDisabled := model.TaskWithId("child", finalizerTaskId, false, statement).
		WithScheduleMinutes(schedule)
	childTaskStandaloneModelDisabled.SetDependsOn(rootTaskStandaloneModelDisabled.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", firstTaskStandaloneModel, secondTaskStandaloneModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, firstTaskStandaloneModel.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceshowoutputassert.TaskShowOutput(t, firstTaskStandaloneModel.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasTaskRelations(sdk.TaskRelations{}).
						HasState(sdk.TaskStateStarted),
					resourceassert.TaskResource(t, secondTaskStandaloneModel.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, secondTaskStandaloneModel.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasTaskRelations(sdk.TaskRelations{}).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change the second task to run after the first one (creating a DAG)
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskModel, childTaskModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskModel.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("2"),
					resourceshowoutputassert.TaskShowOutput(t, rootTaskModel.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasState(sdk.TaskStateStarted),
					// For task relations to be present, in show_output we would have to modify the root task in a way that would
					// trigger show_output recomputing by our custom diff.
					objectassert.Task(t, rootTaskId).HasTaskRelations(sdk.TaskRelations{FinalizerTask: &finalizerTaskId}),
					resourceassert.TaskResource(t, childTaskModel.ResourceReference()).
						HasStartedString(r.BooleanTrue),
					resourceshowoutputassert.TaskShowOutput(t, childTaskModel.ResourceReference()).
						HasTaskRelations(sdk.TaskRelations{FinalizedRootTask: &rootTaskId}).
						HasState(sdk.TaskStateStarted),
				),
			},
			// Change tasks in DAG to standalone tasks (disabled to check if resuming/suspending works correctly)
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskStandaloneModelDisabled, childTaskStandaloneModelDisabled),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanFalse).
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceshowoutputassert.TaskShowOutput(t, rootTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasTaskRelations(sdk.TaskRelations{}).
						HasState(sdk.TaskStateSuspended),
					resourceassert.TaskResource(t, childTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasStartedString(r.BooleanFalse),
					resourceshowoutputassert.TaskShowOutput(t, childTaskStandaloneModelDisabled.ResourceReference()).
						HasScheduleMinutes(schedule).
						HasTaskRelations(sdk.TaskRelations{}).
						HasState(sdk.TaskStateSuspended),
				),
			},
		},
	})
}

func TestAcc_Task_SwitchScheduledWithAfter(t *testing.T) {
	rootId := testClient().Ids.RandomSchemaObjectIdentifier()
	childId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5
	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithScheduleMinutes(schedule).
		WithSuspendTaskAfterNumFailures(1)
	childTaskConfigModel := model.TaskWithId("child", childId, true, statement).
		WithScheduleMinutes(schedule)

	rootTaskConfigModelAfterSuspendFailuresUpdate := model.TaskWithId("root", rootId, true, statement).
		WithScheduleMinutes(schedule).
		WithSuspendTaskAfterNumFailures(2)
	childTaskConfigModelWithAfter := model.TaskWithId("child", childId, true, statement).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName())))
	childTaskConfigModelWithAfter.SetDependsOn(rootTaskConfigModelAfterSuspendFailuresUpdate.ResourceReference())

	rootTaskConfigModelDisabled := model.TaskWithId("root", rootId, false, statement).
		WithScheduleMinutes(schedule)
	childTaskConfigModelDisabled := model.TaskWithId("child", childId, false, statement).
		WithScheduleMinutes(schedule)
	childTaskConfigModelDisabled.SetDependsOn(rootTaskConfigModelDisabled.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceassert.TaskResource(t, childTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasAfter().
						HasSuspendTaskAfterNumFailuresString("10"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModelAfterSuspendFailuresUpdate, childTaskConfigModelWithAfter),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModelAfterSuspendFailuresUpdate.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasSuspendTaskAfterNumFailuresString("2"),
					resourceassert.TaskResource(t, childTaskConfigModelWithAfter.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasNoScheduleSet().
						HasAfter(rootId).
						HasSuspendTaskAfterNumFailuresString("10"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasSuspendTaskAfterNumFailuresString("1"),
					resourceassert.TaskResource(t, childTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule).
						HasAfter().
						HasSuspendTaskAfterNumFailuresString("10"),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModelDisabled, childTaskConfigModelDisabled),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModelDisabled.ResourceReference()).
						HasStartedString(r.BooleanFalse).
						HasScheduleMinutes(schedule).
						HasSuspendTaskAfterNumFailuresString("10"),
					resourceassert.TaskResource(t, childTaskConfigModelDisabled.ResourceReference()).
						HasStartedString(r.BooleanFalse).
						HasScheduleMinutes(schedule).
						HasAfter().
						HasSuspendTaskAfterNumFailuresString("10"),
				),
			},
		},
	})
}

func TestAcc_Task_WithAfter(t *testing.T) {
	rootId := testClient().Ids.RandomSchemaObjectIdentifier()
	childId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithSqlStatement(statement)

	childTaskConfigModelWithoutAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithAfter),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithAfter.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasAfter(rootId),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutAfter),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithoutAfter.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasAfter(),
				),
			},
		},
	})
}

func TestAcc_Task_WithFinalizer(t *testing.T) {
	rootId := testClient().Ids.RandomSchemaObjectIdentifier()
	childId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithFinalize(rootId.FullyQualifiedName()).
		WithSqlStatement(statement)

	childTaskConfigModelWithoutFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithFinalizer),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithFinalizer.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(rootId.FullyQualifiedName()),
				),
			},
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutFinalizer),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithoutFinalizer.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(""),
				),
			},
		},
	})
}

func TestAcc_Task_UpdateFinalizerExternally(t *testing.T) {
	rootId := testClient().Ids.RandomSchemaObjectIdentifier()
	childId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithoutFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithComment("abc").
		WithSqlStatement(statement)

	childTaskConfigModelWithFinalizer := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithFinalize(rootId.FullyQualifiedName()).
		WithComment("abc").
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutFinalizer),
			},
			// Set finalizer externally
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithSuspend(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSuspend(true))

					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithUnset(*sdk.NewTaskUnsetRequest().WithSchedule(true)))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSetFinalize(rootId))

					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithResume(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithResume(true))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutFinalizer),
				Check: assertThat(t,
					resourceassert.TaskResource(t, childTaskConfigModelWithoutFinalizer.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(""),
					resourceshowoutputassert.TaskShowOutput(t, childTaskConfigModelWithoutFinalizer.ResourceReference()).
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// Set finalizer in config
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithFinalizer),
				Check: assertThat(t,
					resourceassert.TaskResource(t, childTaskConfigModelWithFinalizer.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(rootId.FullyQualifiedName()),
					resourceshowoutputassert.TaskShowOutput(t, childTaskConfigModelWithFinalizer.ResourceReference()).
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{FinalizedRootTask: &rootId}),
				),
			},
			// Unset finalizer externally
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithSuspend(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSuspend(true))

					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithUnsetFinalize(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSet(*sdk.NewTaskSetRequest().WithSchedule(fmt.Sprintf("%d minutes", schedule))))

					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithResume(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithResume(true))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithFinalizer),
				Check: assertThat(t,
					resourceassert.TaskResource(t, childTaskConfigModelWithFinalizer.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(rootId.FullyQualifiedName()),
					resourceshowoutputassert.TaskShowOutput(t, childTaskConfigModelWithFinalizer.ResourceReference()).
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{FinalizedRootTask: &rootId}),
				),
			},
			// Unset finalizer in config
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutFinalizer),
				Check: assertThat(t,
					resourceassert.TaskResource(t, childTaskConfigModelWithoutFinalizer.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasFinalizeString(""),
					resourceshowoutputassert.TaskShowOutput(t, childTaskConfigModelWithoutFinalizer.ResourceReference()).
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
		},
	})
}

func TestAcc_Task_UpdateAfterExternally(t *testing.T) {
	rootId := testClient().Ids.RandomSchemaObjectIdentifier()
	childId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModelWithoutAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithComment("abc").
		WithSqlStatement(statement)

	childTaskConfigModelWithAfter := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithComment("abc").
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutAfter),
			},
			// Set after externally
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithSuspend(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSuspend(true))

					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithUnset(*sdk.NewTaskUnsetRequest().WithSchedule(true)))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithAddAfter([]sdk.SchemaObjectIdentifier{rootId}))

					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithResume(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithResume(true))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutAfter),
				Check: assertThat(t,
					resourceassert.TaskResource(t, childTaskConfigModelWithoutAfter.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasAfter(),
					resourceshowoutputassert.TaskShowOutput(t, childTaskConfigModelWithoutAfter.ResourceReference()).
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
			// Set after in config
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithAfter),
				Check: assertThat(t,
					resourceassert.TaskResource(t, childTaskConfigModelWithAfter.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasAfter(rootId),
					resourceshowoutputassert.TaskShowOutput(t, childTaskConfigModelWithAfter.ResourceReference()).
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{Predecessors: []sdk.SchemaObjectIdentifier{rootId}}),
				),
			},
			// Unset after externally
			{
				PreConfig: func() {
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithSuspend(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSuspend(true))

					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithRemoveAfter([]sdk.SchemaObjectIdentifier{rootId}))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithSet(*sdk.NewTaskSetRequest().WithSchedule(fmt.Sprintf("%d MINUTES", schedule))))

					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(childId).WithResume(true))
					testClient().Task.Alter(t, sdk.NewAlterTaskRequest(rootId).WithResume(true))
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithAfter),
				Check: assertThat(t,
					resourceassert.TaskResource(t, childTaskConfigModelWithAfter.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasAfter(rootId),
					resourceshowoutputassert.TaskShowOutput(t, childTaskConfigModelWithAfter.ResourceReference()).
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{Predecessors: []sdk.SchemaObjectIdentifier{rootId}}),
				),
			},
			// Unset after in config
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithoutAfter),
				Check: assertThat(t,
					resourceassert.TaskResource(t, childTaskConfigModelWithoutAfter.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasAfter(),
					resourceshowoutputassert.TaskShowOutput(t, childTaskConfigModelWithoutAfter.ResourceReference()).
						HasState(sdk.TaskStateStarted).
						HasTaskRelations(sdk.TaskRelations{}),
				),
			},
		},
	})
}

func TestAcc_Task_issue2207(t *testing.T) {
	rootId := testClient().Ids.RandomSchemaObjectIdentifier()
	childId := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5

	rootTaskConfigModel := model.TaskWithId("root", rootId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	childTaskConfigModel := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithComment("abc").
		WithSqlStatement(statement)

	childTaskConfigModelWithDifferentComment := model.TaskWithId("child", childId, true, statement).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootId.FullyQualifiedName()))).
		WithComment("def").
		WithSqlStatement(statement)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, childTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasAfter(rootId).
						HasCommentString("abc"),
				),
			},
			// change comment
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(childTaskConfigModelWithDifferentComment.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/with_task_dependency"),
				ConfigVariables: config.ConfigVariablesFromModels(t, "tasks", rootTaskConfigModel, childTaskConfigModelWithDifferentComment),
				Check: assertThat(t,
					resourceassert.TaskResource(t, rootTaskConfigModel.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasScheduleMinutes(schedule),
					resourceassert.TaskResource(t, childTaskConfigModelWithDifferentComment.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasAfter(rootId).
						HasCommentString("def"),
				),
			},
		},
	})
}

func TestAcc_Task_issue2036(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5
	when := "TRUE"

	taskConfigModelWithoutWhen := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement)

	taskConfigModelWithWhen := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement).
		WithWhen(when)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			// create without when
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, taskConfigModelWithoutWhen),
				Check: assertThat(t,
					resourceassert.TaskResource(t, taskConfigModelWithoutWhen.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasWhenString(""),
				),
			},
			// add when
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, taskConfigModelWithWhen),
				Check: assertThat(t,
					resourceassert.TaskResource(t, taskConfigModelWithWhen.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasWhenString("TRUE"),
				),
			},
			// remove when
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables: config.ConfigVariablesFromModel(t, taskConfigModelWithoutWhen),
				Check: assertThat(t,
					resourceassert.TaskResource(t, taskConfigModelWithoutWhen.ResourceReference()).
						HasStartedString(r.BooleanTrue).
						HasWhenString(""),
				),
			},
		},
	})
}

func TestAcc_Task_issue3113(t *testing.T) {
	errorNotificationIntegration, errorNotificationIntegrationCleanup := testClient().NotificationIntegration.CreateWithGcpPubSub(t)
	t.Cleanup(errorNotificationIntegrationCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	schedule := 5
	configModel := model.TaskWithId("test", id, true, statement).
		WithScheduleMinutes(schedule).
		WithSqlStatement(statement).
		WithErrorIntegration(errorNotificationIntegration.ID().Name())

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetV097CompatibleConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.97.0"),
				Config:            taskConfigWithErrorIntegration(id, errorNotificationIntegration.ID()),
				ExpectError:       regexp.MustCompile("error_integration: '' expected type 'string', got unconvertible type 'sdk.AccountObjectIdentifier'"),
			},
			{
				PreConfig: func() {
					testClient().Task.DropFunc(t, id)()
					UnsetConfigPathEnv(t)
				},
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigDirectory:          ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables:          config.ConfigVariablesFromModel(t, configModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasErrorIntegrationString(errorNotificationIntegration.ID().Name()),
				),
			},
		},
	})
}

func TestAcc_Task_StateUpgrade_NoOptionalFields(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	configModel := model.TaskWithId("test", id, false, statement)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.98.0"),
				Config:            taskNoOptionalFieldsConfigV0980(id),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "allow_overlapping_execution", "false"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigDirectory:          ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables:          config.ConfigVariablesFromModel(t, configModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasAllowOverlappingExecutionString(r.BooleanDefault),
				),
			},
		},
	})
}

func TestAcc_Task_StateUpgrade(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	condition := "2 < 1"
	configModel := model.TaskWithId("test", id, false, statement).
		WithScheduleMinutes(5).
		WithAllowOverlappingExecution(r.BooleanTrue).
		WithSuspendTaskAfterNumFailures(10).
		WithWhen(condition).
		WithUserTaskManagedInitialWarehouseSizeEnum(sdk.WarehouseSizeXSmall)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.98.0"),
				Config:            taskBasicConfigV0980(id, condition),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "enabled", "false"),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "allow_overlapping_execution", "true"),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "schedule", "5 MINUTES"),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "suspend_task_after_num_failures", "10"),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "when", condition),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "user_task_managed_initial_warehouse_size", "XSMALL"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				ConfigDirectory:          ConfigurationDirectory("TestAcc_Task/basic"),
				ConfigVariables:          config.ConfigVariablesFromModel(t, configModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasScheduleMinutes(5).
						HasAllowOverlappingExecutionString(r.BooleanTrue).
						HasSuspendTaskAfterNumFailuresString("10").
						HasWhenString(condition).
						HasUserTaskManagedInitialWarehouseSizeEnum(sdk.WarehouseSizeXSmall),
				),
			},
		},
	})
}

func TestAcc_Task_StateUpgradeWithAfter(t *testing.T) {
	rootTask, rootTaskCleanup := testClient().Task.Create(t)
	t.Cleanup(rootTaskCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	statement := "SELECT 1"
	comment := random.Comment()
	configModel := model.TaskWithId("test", id, false, statement).
		WithUserTaskTimeoutMs(50).
		WithWarehouse(testClient().Ids.WarehouseId().Name()).
		WithAfterValue(configvariable.SetVariable(configvariable.StringVariable(rootTask.ID().FullyQualifiedName()))).
		WithComment(comment).
		WithLogLevelEnum(sdk.LogLevelInfo).
		WithAutocommit(false).
		WithJsonIndent(4)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Task),
		Steps: []resource.TestStep{
			{
				PreConfig:         func() { SetLegacyConfigPathEnv(t) },
				ExternalProviders: ExternalProviderWithExactVersion("0.98.0"),
				Config:            taskCompleteConfigV0980(id, rootTask.ID(), testClient().Ids.WarehouseId(), 50, comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "after.#", "1"),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "after.0", rootTask.ID().Name()),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "warehouse", testClient().Ids.WarehouseId().Name()),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "user_task_timeout_ms", "50"),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "comment", comment),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "session_parameters.LOG_LEVEL", "INFO"),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "session_parameters.AUTOCOMMIT", "false"),
					resource.TestCheckResourceAttr(configModel.ResourceReference(), "session_parameters.JSON_INDENT", "4"),
				),
			},
			{
				PreConfig:                func() { UnsetConfigPathEnv(t) },
				ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, configModel),
				Check: assertThat(t,
					resourceassert.TaskResource(t, configModel.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNameString(id.Name()).
						HasStartedString(r.BooleanFalse).
						HasSqlStatementString(statement).
						HasAfter(rootTask.ID()).
						HasWarehouseString(testClient().Ids.WarehouseId().Name()).
						HasUserTaskTimeoutMsString("50").
						HasLogLevelString(string(sdk.LogLevelInfo)).
						HasAutocommitString("false").
						HasJsonIndentString("4").
						HasCommentString(comment),
				),
			},
		},
	})
}

func taskNoOptionalFieldsConfigV0980(id sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_task" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	sql_statement = "SELECT 1"
}
`, id.DatabaseName(), id.SchemaName(), id.Name())
}

func taskConfigWithErrorIntegration(id sdk.SchemaObjectIdentifier, errorIntegrationId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_task" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	schedule = "5 MINUTES"
	sql_statement = "SELECT 1"
	enabled = true
	error_integration = "%[4]s"
}
`, id.DatabaseName(), id.SchemaName(), id.Name(), errorIntegrationId.Name())
}

func taskBasicConfigV0980(id sdk.SchemaObjectIdentifier, condition string) string {
	return fmt.Sprintf(`
resource "snowflake_task" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	enabled = false
	sql_statement = "SELECT 1"
	schedule = "5 MINUTES"
	allow_overlapping_execution = true
	suspend_task_after_num_failures = 10
	when = "%[4]s"
	user_task_managed_initial_warehouse_size = "XSMALL"
}
`, id.DatabaseName(), id.SchemaName(), id.Name(), condition)
}

func taskCompleteConfigV0980(
	id sdk.SchemaObjectIdentifier,
	rootTaskId sdk.SchemaObjectIdentifier,
	warehouseId sdk.AccountObjectIdentifier,
	userTaskTimeoutMs int,
	comment string,
) string {
	return fmt.Sprintf(`
resource "snowflake_task" "test" {
	database = "%[1]s"
	schema = "%[2]s"
	name = "%[3]s"
	enabled = false
	sql_statement = "SELECT 1"

	after = [%[4]s]
	warehouse = "%[5]s"
	user_task_timeout_ms = %[6]d
	comment = "%[7]s"
	session_parameters = {
		LOG_LEVEL = "INFO",
		AUTOCOMMIT = false,
		JSON_INDENT = 4,
	}
}
`, id.DatabaseName(), id.SchemaName(), id.Name(),
		strconv.Quote(rootTaskId.Name()),
		warehouseId.Name(),
		userTaskTimeoutMs,
		comment,
	)
}
