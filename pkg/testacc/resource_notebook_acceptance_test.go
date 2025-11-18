package testacc

import (
	"strconv"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Notebook_basic(t *testing.T) {
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

	path := testClient().Stage.PutOnStageWithPath(t, stage.ID(), "example.ipynb")
	changedPath := testClient().Stage.PutOnStageWithPath(t, changedStage.ID(), "example.ipynb")

	idleAutoShutdownTimeSeconds, changedIdleAutoShutdownTimeSeconds := 3600, 2400

	modelBasic := model.NotebookFromId("test", id)

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
			// create with only required parameters
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
				),
			},
			// import minimal state
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedNotebookResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasNoFromString().
						HasNoMainFile().
						HasNoQueryWarehouse().
						HasNoIdleAutoShutdownTimeSeconds().
						HasWarehouseString("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"").
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImportedNotebookShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasCodeWarehouse(sdk.NewAccountObjectIdentifier("\"SYSTEM$STREAMLIT_NOTEBOOK_WH\"")).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(""),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.name", id.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.main_file", "notebook_app.ipynb")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.query_warehouse", "")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.idle_auto_shutdown_time_seconds", "1800")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.code_warehouse", "SYSTEM$STREAMLIT_NOTEBOOK_WH")),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.CheckImport(importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "describe_output.0.comment", "")),
				),
			},
			// add optional attributes
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
				),
			},
			// import complete state
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"from", "main_file", "idle_auto_shutdown_time_seconds"},
			},
			// alter
			{
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
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
				),
			},
			// change externally alter
			{
				PreConfig: func() {
					testClient().Notebook.Alter(t, sdk.NewAlterNotebookRequest(id).WithSet(
						*sdk.NewNotebookSetRequest().
							WithComment(comment)))
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
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
				),
			},
			// change externally create
			{
				PreConfig: func() {
					newComment := random.Comment()

					newWarehouse, newWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
					t.Cleanup(newWarehouseCleanup)

					newQueryWarehouse, newQueryWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
					t.Cleanup(newQueryWarehouseCleanup)

					newStage, newStageCleanup := testClient().Stage.CreateStage(t)
					t.Cleanup(newStageCleanup)

					testClient().Stage.PutOnStage(t, newStage.ID(), "example.ipynb")
					newLocation := sdk.NewStageLocation(newStage.ID(), "")

					newIdleAutoShutdownTimeSeconds := 4800

					testClient().Notebook.CreateWithRequest(t, sdk.NewCreateNotebookRequest(id).
						WithComment(newComment).
						WithFrom(newLocation).
						WithMainFile("example.ipynb").
						WithQueryWarehouse(newQueryWarehouse.ID()).
						WithIdleAutoShutdownTimeSeconds(newIdleAutoShutdownTimeSeconds).
						WithWarehouse(newWarehouse.ID()).
						WithOrReplace(true),
					)
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
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
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
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
						HasNoWarehouse().
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.NotebookShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(sdk.NewAccountObjectIdentifier("")).
						HasCodeWarehouse(sdk.NewAccountObjectIdentifier("SYSTEM$STREAMLIT_NOTEBOOK_WH")).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(""),
				),
			},
		},
	})
}

func TestAcc_Notebook_complete(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()

	warehouse, warehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(warehouseCleanup)

	queryWarehouse, queryWarehouseCleanup := testClient().Warehouse.CreateWarehouse(t)
	t.Cleanup(queryWarehouseCleanup)

	stage, stageCleanup := testClient().Stage.CreateStage(t)
	t.Cleanup(stageCleanup)

	path := testClient().Stage.PutOnStageWithPath(t, stage.ID(), "example.ipynb")

	idleAutoShutdownTimeSeconds := 3600

	modelComplete := model.NotebookFromId("test", id).
		WithComment(comment).
		WithFrom(path, stage.ID()).
		WithMainFile("example.ipynb").
		WithQueryWarehouse(queryWarehouse.ID().FullyQualifiedName()).
		WithIdleAutoShutdownTimeSeconds(idleAutoShutdownTimeSeconds).
		WithWarehouse(warehouse.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.Notebook),
		Steps: []resource.TestStep{
			// create with all attributes set
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
			{
				Config:                  accconfig.FromModels(t, modelComplete),
				ResourceName:            modelComplete.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"from", "main_file", "idle_auto_shutdown_time_seconds"},
			},
		},
	})
}
