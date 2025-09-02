//go:build !account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tasks_Like_RootTask(t *testing.T) {
	// Created to show LIKE is working
	_, standaloneTaskCleanup := testClient().Task.Create(t)
	t.Cleanup(standaloneTaskCleanup)

	createRootReq := sdk.NewCreateTaskRequest(testClient().Ids.RandomSchemaObjectIdentifier(), "SELECT 1").
		WithSchedule("1 MINUTE").
		WithComment("some comment").
		WithAllowOverlappingExecution(true).
		WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(testClient().Ids.WarehouseId()))
	rootTask, rootTaskCleanup := testClient().Task.CreateWithRequest(t, createRootReq)
	t.Cleanup(rootTaskCleanup)

	childTask, childTaskCleanup := testClient().Task.CreateWithAfter(t, rootTask.ID())
	t.Cleanup(childTaskCleanup)

	tasksModel := datasourcemodel.Tasks("test").
		WithLike(rootTask.ID().Name()).
		WithInDatabase(rootTask.ID().DatabaseId()).
		WithRootOnly(true)
	tasksModelLikeChildRootOnly := datasourcemodel.Tasks("test").
		WithLike(childTask.ID().Name()).
		WithInDatabase(rootTask.ID().DatabaseId()).
		WithRootOnly(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tasksModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(tasksModel.DatasourceReference(), "tasks.#", "1")),
					resourceshowoutputassert.TaskDatasourceShowOutput(t, "snowflake_tasks.test").
						HasName(rootTask.Name).
						HasSchemaName(rootTask.SchemaName).
						HasDatabaseName(rootTask.DatabaseName).
						HasCreatedOnNotEmpty().
						HasIdNotEmpty().
						HasOwnerNotEmpty().
						HasComment("some comment").
						HasWarehouse(testClient().Ids.WarehouseId()).
						HasSchedule("1 MINUTE").
						HasPredecessors().
						HasDefinition("SELECT 1").
						HasCondition("").
						HasAllowOverlappingExecution(true).
						HasErrorIntegrationEmpty().
						HasLastCommittedOn("").
						HasLastSuspendedOn("").
						HasOwnerRoleType("ROLE").
						HasConfig("").
						HasBudget("").
						HasTaskRelations(sdk.TaskRelations{}).
						HasLastSuspendedReason(""),
					resourceparametersassert.TaskDatasourceParameters(t, "snowflake_tasks.test").
						HasAllDefaults(),
				),
			},
			{
				Config: accconfig.FromModels(t, tasksModelLikeChildRootOnly),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tasksModelLikeChildRootOnly.DatasourceReference(), "tasks.#", "0"),
				),
			},
		},
	})
}

func TestAcc_Tasks_In_StartsWith(t *testing.T) {
	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	prefix := testClient().Ids.AlphaN(4)
	taskId1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	taskId2 := testClient().Ids.RandomSchemaObjectIdentifierInSchemaWithPrefix(prefix, schema.ID())
	taskId3 := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	_, standaloneTaskCleanup := testClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(taskId1, "SELECT 1"))
	t.Cleanup(standaloneTaskCleanup)

	_, standaloneTask2Cleanup := testClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(taskId2, "SELECT 1"))
	t.Cleanup(standaloneTask2Cleanup)

	_, standaloneTask3Cleanup := testClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(taskId3, "SELECT 1"))
	t.Cleanup(standaloneTask3Cleanup)

	tasksModelInAccountStartsWith := datasourcemodel.Tasks("test").
		WithStartsWith(prefix).
		WithInAccount()
	tasksModelInDatabaseStartsWith := datasourcemodel.Tasks("test").
		WithStartsWith(prefix).
		WithInDatabase(taskId1.DatabaseId())
	tasksModelInSchemaStartsWith := datasourcemodel.Tasks("test").
		WithStartsWith(prefix).
		WithInSchema(schema.ID())
	tasksModelInSchema := datasourcemodel.Tasks("test").
		WithInSchema(schema.ID())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tasksModelInAccountStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tasksModelInAccountStartsWith.DatasourceReference(), "tasks.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, tasksModelInDatabaseStartsWith),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tasksModelInDatabaseStartsWith.DatasourceReference(), "tasks.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, tasksModelInSchemaStartsWith),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(tasksModelInSchemaStartsWith.DatasourceReference(), "tasks.#", "1")),
					resourceshowoutputassert.TaskDatasourceShowOutput(t, "snowflake_tasks.test").
						HasName(taskId2.Name()).
						HasSchemaName(taskId2.SchemaName()).
						HasDatabaseName(taskId2.DatabaseName()),
				),
			},
			{
				Config: accconfig.FromModels(t, tasksModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tasksModelInSchema.DatasourceReference(), "tasks.#", "2"),
				),
			},
		},
	})
}

func TestAcc_Tasks_Limit(t *testing.T) {
	prefix := testClient().Ids.AlphaN(4)
	taskId1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	taskId2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)

	_, standaloneTaskCleanup := testClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(taskId1, "SELECT 1"))
	t.Cleanup(standaloneTaskCleanup)

	_, standaloneTask2Cleanup := testClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(taskId2, "SELECT 1"))
	t.Cleanup(standaloneTask2Cleanup)

	tasksModelLimitWithPrefix := datasourcemodel.Tasks("test").
		WithLimitRowsAndFrom(2, prefix).
		WithInDatabase(taskId1.DatabaseId())
	tasksModelLimit := datasourcemodel.Tasks("test").
		WithLimitRows(1).
		WithInDatabase(taskId1.DatabaseId())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, tasksModelLimitWithPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tasksModelLimitWithPrefix.DatasourceReference(), "tasks.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, tasksModelLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tasksModelLimit.DatasourceReference(), "tasks.#", "1"),
				),
			},
		},
	})
}
