package testacc

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Notebook_BasicUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	queryWarehouse, queryWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(queryWarehouseCleanup)

	idleAutoShutdownTimeSeconds := 3600

	modelBasic := model.NotebookFromId("test", id)

	assertBasic := []assert.TestCheckFuncProvider{
		resourceassert.NotebookResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNoFromString().
			HasNoMainFile().
			HasNoQueryWarehouse().
			HasNoIdleAutoShutdownTimeSeconds().
			HasNoWarehouse().
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.NotebookShowOutput(t, modelBasic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
			HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.main_file", "notebook_app.ipynb")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.query_warehouse", "")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", "1800")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.code_warehouse", "SYSTEM$STREAMLIT_NOTEBOOK_WH")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
	}

	assertAfterUnset := []assert.TestCheckFuncProvider{
		resourceassert.NotebookResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNoFromString().
			HasNoMainFile().
			HasQueryWarehouseString("").
			HasIdleAutoShutdownTimeSecondsString("0").
			HasWarehouseString("").
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.NotebookShowOutput(t, modelBasic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
			HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.main_file", "notebook_app.ipynb")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.query_warehouse", "")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", "1800")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.code_warehouse", "SYSTEM$STREAMLIT_NOTEBOOK_WH")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
	}

	modelComplete := model.NotebookFromId("test", id).
		WithComment(comment).
		WithQueryWarehouse(queryWarehouse.ID().FullyQualifiedName()).
		WithIdleAutoShutdownTimeSeconds(idleAutoShutdownTimeSeconds).
		WithWarehouse(warehouse.ID().FullyQualifiedName())

	assertComplete := []assert.TestCheckFuncProvider{
		resourceassert.NotebookResource(t, modelComplete.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNoFromString().
			HasNoMainFile().
			HasQueryWarehouseString(queryWarehouse.ID().FullyQualifiedName()).
			HasIdleAutoShutdownTimeSecondsString(strconv.Itoa(idleAutoShutdownTimeSeconds)).
			HasWarehouseString(warehouse.ID().FullyQualifiedName()).
			HasCommentString(comment).
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.NotebookShowOutput(t, modelComplete.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasQueryWarehouse(queryWarehouse.ID()).
			HasCodeWarehouse(warehouse.ID()).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(comment),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Notebook),

		Steps: []resource.TestStep{
			// Create - without optionals
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
			// Import - without optionals
			{
				Config:                  accconfig.FromModels(t, modelBasic),
				ResourceName:            modelBasic.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"warehouse"},
			},
			// Update - set optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check:  assertThat(t, assertComplete...),
			},
			// Import - with optionals
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"from", "main_file", "idle_auto_shutdown_time_seconds"},
			},
			// Update - unset optionals
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertAfterUnset...),
			},
			// Update - detect external changes
			{
				PreConfig: func() {
					testClient().Notebook.Alter(t, sdk.NewAlterNotebookRequest(id).WithSet(
						*sdk.NewNotebookSetRequest().
							WithComment(comment)))
				},
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t, assertAfterUnset...),
			},
			// Destroy - ensure notebook is destroyed before the next step
			{
				Destroy: true,
				Config:  accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					invokeactionassert.NotebookDoesNotExist(t, id),
				),
			},
			// Create - with optionals
			{
				Config: accconfig.FromModels(t, modelComplete),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionCreate),
					},
				},
				Check: assertThat(t, assertComplete...),
			},
		},
	})
}

func TestAcc_Notebook_CompleteUseCase(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := random.Comment(), random.Comment()

	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)
	changedWarehouse, changedWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(changedWarehouseCleanup)

	queryWarehouse, queryWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(queryWarehouseCleanup)
	changedQueryWarehouse, changedQueryWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(changedQueryWarehouseCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)
	changedStage, changedStageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(changedStageCleanup)

	path := testClient().Stage.PutOnStageWithPath(t, stage.ID(), "testdata", "example.ipynb")
	changedPath := testClient().Stage.PutOnStageWithPath(t, changedStage.ID(), "testdata", "example.ipynb")

	idleAutoShutdownTimeSeconds, changedIdleAutoShutdownTimeSeconds := 3600, 2400

	modelComplete := model.NotebookFromId("test", id).
		WithComment(comment).
		WithFrom(path, stage.ID()).
		WithMainFile("example.ipynb").
		WithQueryWarehouse(queryWarehouse.ID().FullyQualifiedName()).
		WithIdleAutoShutdownTimeSeconds(idleAutoShutdownTimeSeconds).
		WithWarehouse(warehouse.ID().FullyQualifiedName())

	modelCompleteWithDifferentValues := model.NotebookFromId("test", id).
		WithComment(changedComment).
		WithFrom(changedPath, changedStage.ID()).
		WithMainFile("example.ipynb").
		WithQueryWarehouse(changedQueryWarehouse.ID().FullyQualifiedName()).
		WithIdleAutoShutdownTimeSeconds(changedIdleAutoShutdownTimeSeconds).
		WithWarehouse(changedWarehouse.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Notebook),
		Steps: []resource.TestStep{
			// Create - with all attributes
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.NotebookResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFromString("testdata/example.ipynb", stage.ID().FullyQualifiedName()).
						HasMainFileString("example.ipynb").
						HasQueryWarehouseString(queryWarehouse.ID().FullyQualifiedName()).
						HasIdleAutoShutdownTimeSecondsString(strconv.Itoa(idleAutoShutdownTimeSeconds)).
						HasWarehouseString(warehouse.ID().FullyQualifiedName()).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.NotebookShowOutput(t, modelComplete.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(queryWarehouse.ID()).
						HasCodeWarehouse(warehouse.ID()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.main_file", "example.ipynb")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.query_warehouse", queryWarehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", fmt.Sprint(idleAutoShutdownTimeSeconds))),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.code_warehouse", warehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", comment)),
				),
			},
			// Import - with all attributes
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"from", "main_file", "idle_auto_shutdown_time_seconds"},
			},
			// Update - forcenew attribues
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				Check: assertThat(t,
					resourceassert.NotebookResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFromString("testdata/example.ipynb", changedStage.ID().FullyQualifiedName()).
						HasMainFileString("example.ipynb").
						HasQueryWarehouseString(changedQueryWarehouse.ID().FullyQualifiedName()).
						HasIdleAutoShutdownTimeSecondsString(strconv.Itoa(changedIdleAutoShutdownTimeSeconds)).
						HasWarehouseString(changedWarehouse.ID().FullyQualifiedName()).
						HasCommentString(changedComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.NotebookShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(changedQueryWarehouse.ID()).
						HasCodeWarehouse(changedWarehouse.ID()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.main_file", "example.ipynb")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.query_warehouse", changedQueryWarehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", fmt.Sprint(changedIdleAutoShutdownTimeSeconds))),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.code_warehouse", changedWarehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.comment", changedComment)),
				),
			},
		},
	})
}

func TestAcc_Notebook_SimultaneousWarehousesChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	queryWarehouse, queryWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(queryWarehouseCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	path := testClient().Stage.PutOnStageWithPath(t, stage.ID(), "testdata", "example.ipynb")

	idleAutoShutdownTimeSeconds := 3600

	modelComplete := model.NotebookFromId("test", id).
		WithComment(comment).
		WithFrom(path, stage.ID()).
		WithMainFile("example.ipynb").
		WithQueryWarehouse(queryWarehouse.ID().FullyQualifiedName()).
		WithIdleAutoShutdownTimeSeconds(idleAutoShutdownTimeSeconds).
		WithWarehouse(warehouse.ID().FullyQualifiedName())

	modelCompleteWithChangedWarehouses := model.NotebookFromId("test", id).
		WithComment(comment).
		WithFrom(path, stage.ID()).
		WithMainFile("example.ipynb").
		WithIdleAutoShutdownTimeSeconds(idleAutoShutdownTimeSeconds)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Notebook),
		Steps: []resource.TestStep{
			// create object
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.NotebookResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFromString("testdata/example.ipynb", stage.ID().FullyQualifiedName()).
						HasMainFileString("example.ipynb").
						HasQueryWarehouseString(queryWarehouse.ID().FullyQualifiedName()).
						HasIdleAutoShutdownTimeSecondsString(strconv.Itoa(idleAutoShutdownTimeSeconds)).
						HasWarehouseString(warehouse.ID().FullyQualifiedName()).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.NotebookShowOutput(t, modelComplete.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(queryWarehouse.ID()).
						HasCodeWarehouse(warehouse.ID()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.main_file", "example.ipynb")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.query_warehouse", queryWarehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", "3600")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.code_warehouse", warehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.comment", comment)),
				),
			},
			// unset warehouse and query_warehouse at the same time
			{
				Config: accconfig.FromModels(t, modelCompleteWithChangedWarehouses),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithChangedWarehouses.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.NotebookResource(t, modelCompleteWithChangedWarehouses.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFromString("testdata/example.ipynb", stage.ID().FullyQualifiedName()).
						HasMainFileString("example.ipynb").
						HasIdleAutoShutdownTimeSecondsString(strconv.Itoa(idleAutoShutdownTimeSeconds)).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.NotebookShowOutput(t, modelCompleteWithChangedWarehouses.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithChangedWarehouses.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithChangedWarehouses.ResourceReference(), "describe_output.0.main_file", "example.ipynb")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithChangedWarehouses.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithChangedWarehouses.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", "3600")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithChangedWarehouses.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithChangedWarehouses.ResourceReference(), "describe_output.0.code_warehouse", "SYSTEM$STREAMLIT_NOTEBOOK_WH")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithChangedWarehouses.ResourceReference(), "describe_output.0.comment", comment)),
				),
			},
		},
	})
}

func TestAcc_Notebook_Rename(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelBasic := model.Notebook("test", id.DatabaseName(), id.SchemaName(), "name")
	modelWithChangedName := model.Notebook("test", id.DatabaseName(), id.SchemaName(), "new_name")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Notebook),
		Steps: []resource.TestStep{
			// create object
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.NotebookResource(t, modelBasic.ResourceReference()).
						HasNameString("name").
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNoFromString().
						HasNoMainFile().
						HasNoQueryWarehouse().
						HasNoIdleAutoShutdownTimeSeconds().
						HasNoWarehouse().
						HasCommentString(""),
					resourceshowoutputassert.NotebookShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName("name").
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(""),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", "name")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.main_file", "notebook_app.ipynb")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", "1800")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.code_warehouse", "SYSTEM$STREAMLIT_NOTEBOOK_WH")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
				),
			},
			// rename object
			{
				Config: accconfig.FromModels(t, modelWithChangedName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithChangedName.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.NotebookResource(t, modelWithChangedName.ResourceReference()).
						HasNameString("new_name").
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNoFromString().
						HasNoMainFile().
						HasNoQueryWarehouse().
						HasNoIdleAutoShutdownTimeSeconds().
						HasNoWarehouse().
						HasCommentString(""),
					resourceshowoutputassert.NotebookShowOutput(t, modelWithChangedName.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName("new_name").
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(""),
					assert.Check(resource.TestCheckResourceAttr(modelWithChangedName.ResourceReference(), "describe_output.0.name", "new_name")),
					assert.Check(resource.TestCheckResourceAttr(modelWithChangedName.ResourceReference(), "describe_output.0.main_file", "notebook_app.ipynb")),
					assert.Check(resource.TestCheckResourceAttr(modelWithChangedName.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(modelWithChangedName.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", "1800")),
					assert.Check(resource.TestCheckResourceAttr(modelWithChangedName.ResourceReference(), "describe_output.0.code_warehouse", "SYSTEM$STREAMLIT_NOTEBOOK_WH")),
					assert.Check(resource.TestCheckResourceAttr(modelWithChangedName.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelWithChangedName.ResourceReference(), "describe_output.0.comment", "")),
				),
			},
		},
	})
}

func TestAcc_Notebook_ExternalWarehouseChange(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	modelBasic := model.NotebookFromId("test", id)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Notebook),
		Steps: []resource.TestStep{
			// create object
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.NotebookResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNoFromString().
						HasNoMainFile().
						HasNoQueryWarehouse().
						HasNoIdleAutoShutdownTimeSeconds().
						HasNoWarehouse().
						HasCommentString(""),
					resourceshowoutputassert.NotebookShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(""),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.main_file", "notebook_app.ipynb")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", "1800")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.code_warehouse", "SYSTEM$STREAMLIT_NOTEBOOK_WH")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
				),
			},
			// change warehouse externally
			{
				PreConfig: func() {
					testClient().Notebook.Alter(t, sdk.NewAlterNotebookRequest(id).WithSet(
						*sdk.NewNotebookSetRequest().
							WithWarehouse(warehouse.ID())))
				},
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift("snowflake_notebook.test", "warehouse", nil, sdk.String(warehouse.ID().Name())),
						planchecks.ExpectChange("snowflake_notebook.test", "warehouse", tfjson.ActionUpdate, sdk.String(warehouse.ID().Name()), nil),
					},
				},
				Check: assertThat(t,
					resourceassert.NotebookResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNoFromString().
						HasNoMainFile().
						HasNoQueryWarehouse().
						HasNoIdleAutoShutdownTimeSeconds().
						HasWarehouseString("").
						HasCommentString(""),
					resourceshowoutputassert.NotebookShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(""),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.main_file", "notebook_app.ipynb")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.query_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", "1800")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.code_warehouse", "SYSTEM$STREAMLIT_NOTEBOOK_WH")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
				),
			},
		},
	})
}

func TestAcc_notebook_WarehouseSchemaLevelChange(t *testing.T) {
	schema, schemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	modelBasic := model.NotebookFromId("test", id)

	assertBasic := []assert.TestCheckFuncProvider{
		objectassert.Notebook(t, id).
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasNoComment().
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")),
		resourceassert.NotebookResource(t, modelBasic.ResourceReference()).
			HasNameString(id.Name()).
			HasDatabaseString(id.DatabaseName()).
			HasSchemaString(id.SchemaName()).
			HasNoFromString().
			HasNoMainFile().
			HasNoQueryWarehouse().
			HasNoIdleAutoShutdownTimeSeconds().
			HasNoWarehouse().
			HasCommentString("").
			HasFullyQualifiedNameString(id.FullyQualifiedName()),
		resourceshowoutputassert.NotebookShowOutput(t, modelBasic.ResourceReference()).
			HasCreatedOnNotEmpty().
			HasName(id.Name()).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
			HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")).
			HasOwner(snowflakeroles.Accountadmin.Name()).
			HasOwnerRoleType("ROLE").
			HasComment(""),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.name", id.Name())),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.main_file", "notebook_app.ipynb")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.query_warehouse", "")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.idle_auto_shutdown_time_seconds", "1800")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.code_warehouse", "SYSTEM$STREAMLIT_NOTEBOOK_WH")),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
		assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.comment", "")),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Notebook),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
			{
				PreConfig: func() {
					testClient().Schema.AlterDefaultStreamlitNotebookWarehouse(t, sdk.NewDatabaseObjectIdentifier(id.DatabaseName(), id.SchemaName()), warehouse.ID())
				},
				Config: accconfig.FromModels(t, modelBasic),
				Check:  assertThat(t, assertBasic...),
			},
		},
	})
}

func TestAcc_Notebook_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()

	modelInvalidIdleAutoShutdownTimeSeconds := model.NotebookFromId("test", id).WithIdleAutoShutdownTimeSeconds(0)

	stage := testClient().Ids.RandomSchemaObjectIdentifier()
	path := "test/path"

	modelFromWithoutMainFile := model.NotebookFromId("test", id).WithFrom(path, stage)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Notebook),
		Steps: []resource.TestStep{
			{
				Config:      config.FromModels(t, modelInvalidIdleAutoShutdownTimeSeconds),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile(`expected idle_auto_shutdown_time_seconds to be at least \(1\), got 0`),
			},
			{
				Config:      config.FromModels(t, modelFromWithoutMainFile),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("\"from\": all of `from,main_file` must be specified"),
			},
		},
	})
}
