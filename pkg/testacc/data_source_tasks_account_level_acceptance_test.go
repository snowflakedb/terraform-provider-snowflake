//go:build account_level_tests

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

// All tests in this file are temporarily moved to account level tests due to STATEMENT_TIMEOUT_IN_SECONDS being set on warehouse level and messing with the results.

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
