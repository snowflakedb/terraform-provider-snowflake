//go:build account_level_tests

// These tests are temporarily moved to account level tests due to flakiness caused by changes in the higher-level parameters.
// Some tests might be also affected due to STATEMENT_TIMEOUT_IN_SECONDS being set on warehouse level and messing with the results.

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Tasks_BasicUseCase_DifferentFiltering(t *testing.T) {
	prefix := testClient().Ids.AlphaN(4)
	rootTaskId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	childTaskId := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	standaloneTaskId := testClient().Ids.RandomSchemaObjectIdentifier()

	_, rootTaskCleanup := testClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(rootTaskId, "SELECT 1"))
	t.Cleanup(rootTaskCleanup)

	_, childTaskCleanup := testClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(childTaskId, "SELECT 1").WithAfter([]sdk.SchemaObjectIdentifier{rootTaskId}))
	t.Cleanup(childTaskCleanup)

	_, standaloneTaskCleanup := testClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(standaloneTaskId, "SELECT 1"))
	t.Cleanup(standaloneTaskCleanup)

	datasourceModel1 := datasourcemodel.Tasks("test").
		WithLike(rootTaskId.Name()).
		WithInDatabase(rootTaskId.DatabaseId())

	datasourceModel2 := datasourcemodel.Tasks("test").
		WithLike(prefix + "%").
		WithInDatabase(rootTaskId.DatabaseId())

	datasourceModel3 := datasourcemodel.Tasks("test").
		WithStartsWith(prefix).
		WithInDatabase(rootTaskId.DatabaseId())

	datasourceModel4 := datasourcemodel.Tasks("test").
		WithLimitRows(1).
		WithInDatabase(rootTaskId.DatabaseId())

	datasourceModel5 := datasourcemodel.Tasks("test").
		WithLike(prefix + "%").
		WithInDatabase(rootTaskId.DatabaseId()).
		WithRootOnly(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, datasourceModel1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModel1.DatasourceReference(), "tasks.#", "1"),
					resource.TestCheckResourceAttr(datasourceModel1.DatasourceReference(), "tasks.0.show_output.0.name", rootTaskId.Name()),
				),
			},
			{
				Config: accconfig.FromModels(t, datasourceModel2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModel2.DatasourceReference(), "tasks.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, datasourceModel3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModel3.DatasourceReference(), "tasks.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, datasourceModel4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModel4.DatasourceReference(), "tasks.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, datasourceModel5),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceModel5.DatasourceReference(), "tasks.#", "1"),
					resource.TestCheckResourceAttr(datasourceModel5.DatasourceReference(), "tasks.0.show_output.0.name", rootTaskId.Name()),
				),
			},
		},
	})
}

func TestAcc_Tasks_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	_, rootTaskCleanup := testClient().Task.CreateWithRequest(t, sdk.NewCreateTaskRequest(id, "SELECT 1").
		WithSchedule("1 MINUTE").
		WithComment(comment).
		WithAllowOverlappingExecution(true).
		WithWarehouse(*sdk.NewCreateTaskWarehouseRequest().WithWarehouse(testClient().Ids.WarehouseId())))
	t.Cleanup(rootTaskCleanup)

	datasourceModelWithoutParameters := datasourcemodel.Tasks("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithWithParameters(false)

	datasourceModelWithParameters := datasourcemodel.Tasks("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithWithParameters(true)

	commonShowOutputAsserts := resourceshowoutputassert.TaskDatasourceShowOutput(t, "test").
		HasName(id.Name()).
		HasSchemaName(id.SchemaName()).
		HasDatabaseName(id.DatabaseName()).
		HasCreatedOnNotEmpty().
		HasIdNotEmpty().
		HasOwnerNotEmpty().
		HasComment(comment).
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
		HasLastSuspendedReason("")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, datasourceModelWithoutParameters),
				Check: assertThat(t,
					commonShowOutputAsserts,
					assert.Check(resource.TestCheckResourceAttr(datasourceModelWithoutParameters.DatasourceReference(), "tasks.0.parameters.#", "0")),
				),
			},
			{
				Config: accconfig.FromModels(t, datasourceModelWithParameters),
				Check: assertThat(t,
					commonShowOutputAsserts,
					resourceparametersassert.TaskDatasourceParameters(t, "snowflake_tasks.test").
						HasAllDefaults(),
				),
			},
		},
	})
}
