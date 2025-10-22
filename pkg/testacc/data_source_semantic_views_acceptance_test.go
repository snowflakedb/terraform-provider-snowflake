//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SemanticViews_Basic(t *testing.T) {
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

	semanticViewModel := model.SemanticViewWithMetrics(
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

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.object_kind", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.object_name", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.0.property_value", comment)),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.1.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.1.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.1.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.1.property", "BASE_TABLE_DATABASE_NAME")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.1.property_value", table1.DatabaseName)),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.2.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.2.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.2.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.2.property", "BASE_TABLE_SCHEMA_NAME")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.2.property_value", table1.SchemaName)),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.3.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.3.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.3.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.3.property", "BASE_TABLE_NAME")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.3.property_value", table1.Name)),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.4.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.4.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.4.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.4.property", "PRIMARY_KEY")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.4.property_value", "[\"A1\"]")),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.5.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.5.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.5.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.5.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.5.property_value", "logical table 1")),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.6.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.6.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.6.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.6.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.6.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.7.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.7.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.7.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.7.property", "EXPRESSION")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.7.property_value", "SUM(lt1.a1)")),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.8.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.8.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.8.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.8.property", "DATA_TYPE")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.8.property_value", "NUMBER(38,0)")),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.9.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.9.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.9.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.9.property", "SYNONYMS")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.9.property_value", "[\"sem1\",\"baseSem\"]")),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.10.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.10.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.10.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.10.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.10.property_value", "semantic expression 1")),

					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.11.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.11.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.11.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.11.property", "ACCESS_MODIFIER")),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModel.DatasourceReference(), "semantic_views.0.describe_output.11.property_value", "PUBLIC")),
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
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(dataSourceModelWithoutOptionals.DatasourceReference(), "semantic_views.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_SemanticViews_Filtering(t *testing.T) {
	table1, table1Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("a1", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("a2", sdk.DataTypeNumber),
	})
	t.Cleanup(table1Cleanup)

	logicalTable1 := model.LogicalTableWithProps("lt1", table1.ID(), []sdk.SemanticViewColumn{{Name: "a1"}}, nil, nil, "logical table 1")
	semExp1 := model.SemanticExpressionWithProps("lt1.se1", "SUM(lt1.a1)", []sdk.Synonym{{Synonym: "sem1"}, {Synonym: "baseSem"}}, "semantic expression 1")

	metric1 := model.MetricDefinitionWithProps(semExp1, nil)

	prefix := random.AlphaUpperN(4)

	id1 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := testClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := testClient().Ids.RandomSchemaObjectIdentifier()

	model1 := model.SemanticViewWithMetrics("test1", id1.DatabaseName(), id1.SchemaName(), id1.Name(), []sdk.LogicalTable{*logicalTable1}, []sdk.MetricDefinition{*metric1})
	model2 := model.SemanticViewWithMetrics("test2", id2.DatabaseName(), id2.SchemaName(), id2.Name(), []sdk.LogicalTable{*logicalTable1}, []sdk.MetricDefinition{*metric1})
	model3 := model.SemanticViewWithMetrics("test3", id3.DatabaseName(), id3.SchemaName(), id3.Name(), []sdk.LogicalTable{*logicalTable1}, []sdk.MetricDefinition{*metric1})

	dataSourceModelLikeFirstOne := datasourcemodel.SemanticViews("test").
		WithLike(id1.Name()).
		WithWithDescribe(false).
		WithDependsOn(model1.ResourceReference(), model2.ResourceReference(), model3.ResourceReference())

	dataSourceModelLikePrefix := datasourcemodel.SemanticViews("test").
		WithLike(prefix+"%").
		WithWithDescribe(false).
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
