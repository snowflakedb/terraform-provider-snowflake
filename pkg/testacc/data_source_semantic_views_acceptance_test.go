//go:build non_account_level_tests

package testacc

import (
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SemanticViews_basic(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	table1, table1Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("a1", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("a2", sdk.DataTypeNumber),
	})
	t.Cleanup(table1Cleanup)

	logicalTable1 := model.LogicalTableWithProps("lt1", table1.ID(), []sdk.SemanticViewColumn{{Name: "a1"}}, nil, nil, "logical table 1")
	semExp1 := model.SemanticExpressionWithProps("lt1.se1", "SUM(lt1.a1)", []sdk.Synonym{{Synonym: "sem1"}, {Synonym: "baseSem"}}, "semantic expression 1")

	metric1 := model.MetricDefinitionWithProps(semExp1, nil)

	semanticViewModel := model.SemanticView(
		"test",
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		[]sdk.LogicalTable{*logicalTable1},
		[]sdk.MetricDefinition{*metric1},
	).WithComment(comment)

	dataSourceModel := datasourcemodel.SemanticViews("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(semanticViewModel.ResourceReference())

	dataSourceModelWithoutOptionals := datasourcemodel.SemanticViews("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(semanticViewModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, semanticViewModel, dataSourceModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.#", "1")),

					resourceshowoutputassert.SemanticViewsDatasourceShowOutput(t, dataSourceModel.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty().
						HasComment(comment).
						HasOwnerRoleType("ROLE").
						HasNoExtension(),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				Config: accconfig.FromModels(t, semanticViewModel, dataSourceModelWithoutOptionals),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dataSourceModelWithoutOptionals.DatasourceReference(), "semantic_views.#", "1")),

					resourceshowoutputassert.SemanticViewsDatasourceShowOutput(t, dataSourceModelWithoutOptionals.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasCreatedOnNotEmpty().
						HasComment(comment).
						HasNoExtension().
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModelWithoutOptionals.DatasourceReference(), "semantic_views.0.describe_output.#", "0")),
				),
			},
		},
	})
}
