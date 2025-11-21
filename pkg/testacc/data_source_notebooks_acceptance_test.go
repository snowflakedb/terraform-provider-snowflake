package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Notebooks(t *testing.T) {
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

	completeModel := model.NotebookFromId("test", id).
		WithComment(comment).
		WithFrom(path, stage.ID()).
		WithMainFile("example.ipynb").
		WithQueryWarehouse(queryWarehouse.ID().FullyQualifiedName()).
		WithIdleAutoShutdownTimeSeconds(idleAutoShutdownTimeSeconds).
		WithWarehouse(warehouse.ID().FullyQualifiedName())

	notebooksModel := datasourcemodel.Notebooks("test").
		WithLike(id.Name()).
		WithDependsOn(completeModel.ResourceReference())

	notebooksModelWithoutOptionals := datasourcemodel.Notebooks("test").
		WithWithDescribe(false).
		WithLike(id.Name()).
		WithDependsOn(completeModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, completeModel, notebooksModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(notebooksModel.DatasourceReference(), "notebooks.#", "1")),

					resourceshowoutputassert.NotebooksDatasourceShowOutput(t, "snowflake_notebooks.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(queryWarehouse.ID()).
						HasCodeWarehouse(warehouse.ID()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttr(notebooksModel.DatasourceReference(), "notebooks.0.describe_output.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(notebooksModel.DatasourceReference(), "notebooks.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(notebooksModel.DatasourceReference(), "notebooks.0.describe_output.0.main_file", "example.ipynb")),
					assert.Check(resource.TestCheckResourceAttr(notebooksModel.DatasourceReference(), "notebooks.0.describe_output.0.query_warehouse", queryWarehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(notebooksModel.DatasourceReference(), "notebooks.0.describe_output.0.idle_auto_shutdown_time_seconds", "3600")),
					assert.Check(resource.TestCheckResourceAttr(notebooksModel.DatasourceReference(), "notebooks.0.describe_output.0.code_warehouse", warehouse.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(notebooksModel.DatasourceReference(), "notebooks.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(notebooksModel.DatasourceReference(), "notebooks.0.describe_output.0.comment", comment)),
				),
			},
			{
				Config: accconfig.FromModels(t, completeModel, notebooksModelWithoutOptionals),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(notebooksModelWithoutOptionals.DatasourceReference(), "notebooks.#", "1")),
					resourceshowoutputassert.NotebooksDatasourceShowOutput(t, "snowflake_notebooks.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasQueryWarehouse(queryWarehouse.ID()).
						HasCodeWarehouse(warehouse.ID()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment),
					assert.Check(resource.TestCheckResourceAttr(notebooksModelWithoutOptionals.DatasourceReference(), "notebooks.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Notebooks_Filtering(t *testing.T) {
	prefix := random.AlphaN(4)
	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.NotebookFromId("test1", id1)
	model2 := model.NotebookFromId("test2", id2)
	model3 := model.NotebookFromId("test3", id3)
	notebooksModelLikeFirstOne := datasourcemodel.Notebooks("test").
		WithLike(id1.Name()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	notebooksModelLikePrefix := datasourcemodel.Notebooks("test").
		WithLike(prefix+"%").
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	notebooksWithStartsWithModel := datasourcemodel.Notebooks("test").
		WithStartsWith(prefix).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())
	notebooksWithLimitModel := datasourcemodel.Notebooks("test").
		WithRowsAndFrom(1, prefix).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, notebooksModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(notebooksModelLikeFirstOne.DatasourceReference(), "notebooks.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, notebooksModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(notebooksModelLikePrefix.DatasourceReference(), "notebooks.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, notebooksWithStartsWithModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(notebooksWithStartsWithModel.DatasourceReference(), "notebooks.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, notebooksWithLimitModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(notebooksWithLimitModel.DatasourceReference(), "notebooks.#", "1"),
				),
			},
		},
	})
}
