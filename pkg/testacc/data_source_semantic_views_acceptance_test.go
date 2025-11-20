//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TODO [SNOW-2108211]: show output assertions
func TestAcc_SemanticViews_Basic(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	table1, table1Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest(`"a1"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a2"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a3"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a4"`, sdk.DataTypeNumber),
	})
	t.Cleanup(table1Cleanup)

	logicalTable1 := model.LogicalTableWithProps("lt1", table1.ID(), []sdk.SemanticViewColumn{{Name: "a1"}}, [][]sdk.SemanticViewColumn{{{Name: "a2"}}, {{Name: "a3"}, {Name: "a4"}}}, nil, "logical table 1")
	semExp1 := model.SemanticExpressionWithProps(`"lt1"."se1"`, `SUM("lt1"."a1")`, []sdk.Synonym{{Synonym: "sem1"}, {Synonym: "baseSem"}}, "semantic expression 1")

	metric1 := model.MetricDefinitionWithProps(semExp1, nil)

	semanticViewModel := model.SemanticViewWithMetrics(
		"test",
		id,
		[]sdk.LogicalTable{*logicalTable1},
		[]sdk.MetricDefinition{*metric1},
	).WithComment(comment)

	dataSourceModel := datasourcemodel.SemanticViews("test").
		WithLike(id.Name()).
		WithInDatabase(id.DatabaseId()).
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
						HasExtension("").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(dataSourceModel.DatasourceReference(), "semantic_views.0.show_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.show_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.show_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.show_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.show_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.show_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.show_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.show_output.0.extension", "")),
				),
			},
		},
	})
}

func TestAcc_SemanticViews_Filtering(t *testing.T) {
	table1, table1Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest(`"a1"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a2"`, sdk.DataTypeNumber),
	})
	t.Cleanup(table1Cleanup)

	table2, table2Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest(`"a1"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a2"`, sdk.DataTypeNumber),
	})
	t.Cleanup(table2Cleanup)

	table3, table3Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest(`"a1"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a2"`, sdk.DataTypeNumber),
	})
	t.Cleanup(table3Cleanup)

	logicalTable1 := model.LogicalTableWithProps("lt1", table1.ID(), []sdk.SemanticViewColumn{{Name: "a1"}}, nil, nil, "logical table 1")
	logicalTable2 := model.LogicalTableWithProps("lt2", table2.ID(), []sdk.SemanticViewColumn{{Name: "a1"}}, nil, nil, "logical table 2")
	logicalTable3 := model.LogicalTableWithProps("lt3", table3.ID(), []sdk.SemanticViewColumn{{Name: "a1"}}, nil, nil, "logical table 3")
	semExp1 := model.SemanticExpressionWithProps(`"lt1"."se1"`, `SUM("lt1"."a1")`, []sdk.Synonym{{Synonym: "sem1"}, {Synonym: "baseSem"}}, "semantic expression 1")
	semExp2 := model.SemanticExpressionWithProps(`"lt2"."se1"`, `SUM("lt2"."a1")`, []sdk.Synonym{{Synonym: "sem1"}, {Synonym: "baseSem"}}, "semantic expression 1")
	semExp3 := model.SemanticExpressionWithProps(`"lt3"."se1"`, `SUM("lt3"."a1")`, []sdk.Synonym{{Synonym: "sem1"}, {Synonym: "baseSem"}}, "semantic expression 1")

	metric1 := model.MetricDefinitionWithProps(semExp1, nil)
	metric2 := model.MetricDefinitionWithProps(semExp2, nil)
	metric3 := model.MetricDefinitionWithProps(semExp3, nil)

	prefix := random.AlphaUpperN(4)

	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.SemanticViewWithMetrics("test1", id1, []sdk.LogicalTable{*logicalTable1}, []sdk.MetricDefinition{*metric1})
	model2 := model.SemanticViewWithMetrics("test2", id2, []sdk.LogicalTable{*logicalTable2}, []sdk.MetricDefinition{*metric2})
	model3 := model.SemanticViewWithMetrics("test3", id3, []sdk.LogicalTable{*logicalTable3}, []sdk.MetricDefinition{*metric3})

	dataSourceModelLikeFirstOne := datasourcemodel.SemanticViews("test").
		WithLike(id1.Name()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	dataSourceModelLikePrefix := datasourcemodel.SemanticViews("test").
		WithLike(prefix+"%").
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	dataSourceModelWithIn := datasourcemodel.SemanticViews("test").
		WithInDatabase(id1.DatabaseId()).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	dataSourceModelWithLimit := datasourcemodel.SemanticViews("test").
		WithRowsAndFrom(2, "").
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, model1, model2, model3, dataSourceModelLikeFirstOne),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelLikeFirstOne.DatasourceReference(), "semantic_views.#", "1"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, dataSourceModelLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelLikePrefix.DatasourceReference(), "semantic_views.#", "2"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, dataSourceModelWithIn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelWithIn.DatasourceReference(), "semantic_views.#", "3"),
				),
			},
			{
				Config: accconfig.FromModels(t, model1, model2, model3, dataSourceModelWithLimit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceModelWithLimit.DatasourceReference(), "semantic_views.#", "2"),
				),
			},
		},
	})
}

func TestAcc_SemanticViews_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, datasourcemodel.SemanticViews("test").WithEmptyIn()),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func TestAcc_SemanticViews_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: ConfigurationDirectory("TestAcc_SemanticViews/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one semantic view"),
			},
		},
	})
}
