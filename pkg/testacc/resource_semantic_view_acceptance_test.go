//go:build account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SemanticView_basic(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment, changedComment := "comment 1", "comment 2"
	table1, table1Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("a1", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("a2", sdk.DataTypeNumber),
	})
	t.Cleanup(table1Cleanup)
	table2, table2Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("a1", sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest("a2", sdk.DataTypeNumber),
	})
	t.Cleanup(table2Cleanup)
	logicalTable1 := model.LogicalTableWithProps("lt1", table1.ID(), []sdk.SemanticViewColumn{{Name: "a1"}}, nil, nil, "logical table 1")
	logicalTable2 := model.LogicalTableWithProps("lt2", table2.ID(), []sdk.SemanticViewColumn{{Name: "a1"}}, nil, nil, "")
	semExp1 := model.SemanticExpressionWithProps("lt1.se1", "SUM(lt1.a1)", []sdk.Synonym{{Synonym: "sem1"}, {Synonym: "baseSem"}}, "semantic expression 1")
	partitionBy := "lt1.d2"
	windowFunc1 := model.WindowFunctionMetricDefinitionWithProps("lt1.wf1", "sum(lt1.se1)", sdk.WindowFunctionOverClause{PartitionBy: &partitionBy})
	metric1 := model.MetricDefinitionWithProps(semExp1, nil)
	metric2 := model.MetricDefinitionWithProps(nil, windowFunc1)
	relTableAlias := model.RelationshipTableAliasWithProps("lt1", table1.ID())
	relTableColumns := []sdk.SemanticViewColumn{
		{
			Name: "a1",
		},
		{
			Name: "a2",
		},
	}
	refTableAlias := model.RelationshipTableAliasWithProps("lt2", table2.ID())
	refRelTableColumns := []sdk.SemanticViewColumn{
		{
			Name: "a1",
		},
		{
			Name: "a2",
		},
	}

	rel1 := model.RelationshipWithProps("r1", *relTableAlias, relTableColumns, *refTableAlias, refRelTableColumns)
	fact1 := model.SemanticExpressionWithProps("lt1.f1", "lt1.a2", []sdk.Synonym{{Synonym: "fact1"}}, "fact 1")
	dimension1 := model.SemanticExpressionWithProps("lt1.d1", "lt1.a1", []sdk.Synonym{{Synonym: "dim1"}}, "dimension 1")

	relTableAlias2 := model.RelationshipTableAliasWithProps("lt2", table1.ID())
	relTableColumns2 := []sdk.SemanticViewColumn{
		{
			Name: "a1",
		},
		{
			Name: "a2",
		},
	}
	refTableAlias2 := model.RelationshipTableAliasWithProps("lt1", table2.ID())
	refRelTableColumns2 := []sdk.SemanticViewColumn{
		{
			Name: "a1",
		},
		{
			Name: "a2",
		},
	}
	rel2 := model.RelationshipWithProps("r2", *relTableAlias2, relTableColumns2, *refTableAlias2, refRelTableColumns2)
	fact2 := model.SemanticExpressionWithProps("lt1.f2", "lt1.a1", []sdk.Synonym{{Synonym: "fact2"}}, "fact 2")
	dimension2 := model.SemanticExpressionWithProps("lt1.d2", "lt1.a2", []sdk.Synonym{{Synonym: "dim2"}}, "dimension 2")

	lt1Request := sdk.NewLogicalTableRequest(table1.ID()).WithLogicalTableAlias(sdk.LogicalTableAliasRequest{LogicalTableAlias: "lt1"})
	lt2Request := sdk.NewLogicalTableRequest(table2.ID()).WithLogicalTableAlias(sdk.LogicalTableAliasRequest{LogicalTableAlias: "lt2"})
	seRequest := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: "lt1.se1"}, &sdk.SemanticSqlExpressionRequest{SqlExpression: "SUM(lt1.a1)"})
	wfRequest := sdk.NewWindowFunctionMetricDefinitionRequest("lt2.wf1", "sum(lt2.a1)")
	m1Request := sdk.NewMetricDefinitionRequest().WithSemanticExpression(*seRequest)
	m2Request := sdk.NewMetricDefinitionRequest().WithWindowFunctionMetricDefinition(*wfRequest)

	modelBasic := model.SemanticView(
		"test",
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		[]sdk.LogicalTable{*logicalTable1},
		[]sdk.MetricDefinition{*metric1},
	)

	modelComplete := model.SemanticView(
		"test",
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		[]sdk.LogicalTable{*logicalTable1, *logicalTable2},
		[]sdk.MetricDefinition{*metric1},
	).WithComment(comment).
		WithRelationships([]sdk.SemanticViewRelationship{*rel1}).
		WithFacts([]sdk.SemanticExpression{*fact1}).
		WithDimensions([]sdk.SemanticExpression{*dimension1})

	modelCompleteWithDifferentValues := model.SemanticView(
		"test",
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		[]sdk.LogicalTable{*logicalTable1, *logicalTable2},
		[]sdk.MetricDefinition{*metric1, *metric2},
	).WithComment(changedComment).
		WithRelationships([]sdk.SemanticViewRelationship{*rel2}).
		WithFacts([]sdk.SemanticExpression{*fact2}).
		WithDimensions([]sdk.SemanticExpression{*dimension2})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SemanticView),

		Steps: []resource.TestStep{
			// create with only required attributes
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.SemanticViewShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment("").
						HasExtension(""),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.#", "11")),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.property", "BASE_TABLE_DATABASE_NAME")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.0.property_value", table1.DatabaseName)),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.1.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.1.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.1.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.1.property", "BASE_TABLE_SCHEMA_NAME")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.1.property_value", table1.SchemaName)),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.2.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.2.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.2.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.2.property", "BASE_TABLE_NAME")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.2.property_value", table1.Name)),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.3.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.3.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.3.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.3.property", "PRIMARY_KEY")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.3.property_value", "[\"A1\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.4.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.4.object_name", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.4.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.4.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.4.property_value", "logical table 1")),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.5.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.5.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.5.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.5.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.5.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.6.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.6.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.6.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.6.property", "EXPRESSION")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.6.property_value", "SUM(lt1.a1)")),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.7.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.7.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.7.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.7.property", "DATA_TYPE")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.7.property_value", "NUMBER(38,0)")),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.8.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.8.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.8.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.8.property", "SYNONYMS")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.8.property_value", "[\"sem1\",\"baseSem\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.9.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.9.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.9.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.9.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.9.property_value", "semantic expression 1")),

					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.10.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.10.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.10.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.10.property", "ACCESS_MODIFIER")),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.10.property_value", "PUBLIC")),
				),
			},
			//// import minimal state
			{
				Config:       accconfig.FromModels(t, modelBasic),
				ResourceName: modelBasic.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedSemanticViewResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImportedSemanticViewShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment("").
						HasExtension(""),
				),
			},
			// add optional attributes
			{
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.SemanticViewShowOutput(t, modelComplete.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(comment).
						HasExtension(""),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.#", "32")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.object_kind", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.object_name", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.property_value", comment)),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.6.object_kind", "RELATIONSHIP")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.6.object_name", "R1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.6.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.6.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.6.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.7.object_kind", "RELATIONSHIP")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.7.object_name", "R1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.7.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.7.property", "REF_TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.7.property_value", "LT2")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.8.object_kind", "RELATIONSHIP")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.8.object_name", "R1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.8.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.8.property", "FOREIGN_KEY")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.8.property_value", "[\"A1\",\"A2\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.9.object_kind", "RELATIONSHIP")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.9.object_name", "R1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.9.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.9.property", "REF_KEY")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.9.property_value", "[\"A1\",\"A2\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.10.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.10.object_name", "D1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.10.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.10.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.10.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.11.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.11.object_name", "D1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.11.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.11.property", "EXPRESSION")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.11.property_value", "lt1.a1")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.12.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.12.object_name", "D1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.12.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.12.property", "DATA_TYPE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.12.property_value", "NUMBER(38,0)")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.13.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.13.object_name", "D1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.13.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.13.property", "SYNONYMS")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.13.property_value", "[\"dim1\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.14.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.14.object_name", "D1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.14.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.14.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.14.property_value", "dimension 1")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.15.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.15.object_name", "D1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.15.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.15.property", "ACCESS_MODIFIER")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.15.property_value", "PUBLIC")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.16.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.16.object_name", "F1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.16.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.16.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.16.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.17.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.17.object_name", "F1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.17.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.17.property", "EXPRESSION")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.17.property_value", "lt1.a2")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.18.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.18.object_name", "F1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.18.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.18.property", "DATA_TYPE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.18.property_value", "NUMBER(38,0)")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.19.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.19.object_name", "F1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.19.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.19.property", "SYNONYMS")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.19.property_value", "[\"fact1\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.20.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.20.object_name", "F1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.20.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.20.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.20.property_value", "fact 1")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.21.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.21.object_name", "F1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.21.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.21.property", "ACCESS_MODIFIER")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.21.property_value", "PUBLIC")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.22.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.22.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.22.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.22.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.22.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.23.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.23.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.23.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.23.property", "EXPRESSION")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.23.property_value", "SUM(lt1.a1)")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.24.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.24.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.24.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.24.property", "DATA_TYPE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.24.property_value", "NUMBER(38,0)")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.25.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.25.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.25.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.25.property", "SYNONYMS")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.25.property_value", "[\"sem1\",\"baseSem\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.26.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.26.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.26.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.26.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.26.property_value", "semantic expression 1")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.27.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.27.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.27.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.27.property", "ACCESS_MODIFIER")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.27.property_value", "PUBLIC")),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.28.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.28.object_name", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.28.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.28.property", "BASE_TABLE_DATABASE_NAME")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.28.property_value", table2.DatabaseName)),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.29.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.29.object_name", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.29.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.29.property", "BASE_TABLE_SCHEMA_NAME")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.29.property_value", table2.SchemaName)),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.30.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.30.object_name", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.30.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.30.property", "BASE_TABLE_NAME")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.30.property_value", table2.Name)),

					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.31.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.31.object_name", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.31.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.31.property", "PRIMARY_KEY")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.31.property_value", "[\"A1\"]")),
				),
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
					resourceassert.SemanticViewResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(changedComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.SemanticViewShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment).
						HasExtension(""),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.#", "36")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.object_kind", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.object_name", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.0.property_value", changedComment)),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.6.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.6.object_name", "D2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.6.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.6.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.6.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.7.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.7.object_name", "D2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.7.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.7.property", "EXPRESSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.7.property_value", "lt1.a2")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.8.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.8.object_name", "D2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.8.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.8.property", "DATA_TYPE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.8.property_value", "NUMBER(38,0)")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.9.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.9.object_name", "D2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.9.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.9.property", "SYNONYMS")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.9.property_value", "[\"dim2\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.10.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.10.object_name", "D2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.10.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.10.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.10.property_value", "dimension 2")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.11.object_kind", "DIMENSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.11.object_name", "D2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.11.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.11.property", "ACCESS_MODIFIER")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.11.property_value", "PUBLIC")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.12.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.12.object_name", "F2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.12.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.12.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.12.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.13.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.13.object_name", "F2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.13.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.13.property", "EXPRESSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.13.property_value", "lt1.a1")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.14.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.14.object_name", "F2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.14.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.14.property", "DATA_TYPE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.14.property_value", "NUMBER(38,0)")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.15.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.15.object_name", "F2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.15.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.15.property", "SYNONYMS")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.15.property_value", "[\"fact2\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.16.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.16.object_name", "F2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.16.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.16.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.16.property_value", "fact 2")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.17.object_kind", "FACT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.17.object_name", "F2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.17.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.17.property", "ACCESS_MODIFIER")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.17.property_value", "PUBLIC")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.18.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.18.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.18.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.18.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.18.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.19.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.19.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.19.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.19.property", "EXPRESSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.19.property_value", "SUM(lt1.a1)")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.20.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.20.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.20.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.20.property", "DATA_TYPE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.20.property_value", "NUMBER(38,0)")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.21.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.21.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.21.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.21.property", "SYNONYMS")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.21.property_value", "[\"sem1\",\"baseSem\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.22.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.22.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.22.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.22.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.22.property_value", "semantic expression 1")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.23.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.23.object_name", "SE1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.23.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.23.property", "ACCESS_MODIFIER")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.23.property_value", "PUBLIC")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.24.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.24.object_name", "WF1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.24.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.24.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.24.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.25.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.25.object_name", "WF1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.25.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.25.property", "EXPRESSION")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.25.property_value", "sum(lt1.se1) OVER (PARTITION BY lt1.d2)")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.26.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.26.object_name", "WF1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.26.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.26.property", "DATA_TYPE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.26.property_value", "NUMBER(38,0)")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.27.object_kind", "METRIC")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.27.object_name", "WF1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.27.parent_entity", "LT1")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.27.property", "ACCESS_MODIFIER")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.27.property_value", "PUBLIC")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.28.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.28.object_name", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.28.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.28.property", "BASE_TABLE_DATABASE_NAME")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.28.property_value", table2.DatabaseName)),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.29.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.29.object_name", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.29.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.29.property", "BASE_TABLE_SCHEMA_NAME")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.29.property_value", table2.SchemaName)),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.30.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.30.object_name", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.30.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.30.property", "BASE_TABLE_NAME")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.30.property_value", table2.Name)),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.31.object_kind", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.31.object_name", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.31.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.31.property", "PRIMARY_KEY")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.31.property_value", "[\"A1\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.32.object_kind", "RELATIONSHIP")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.32.object_name", "R2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.32.parent_entity", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.32.property", "TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.32.property_value", "LT2")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.33.object_kind", "RELATIONSHIP")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.33.object_name", "R2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.33.parent_entity", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.33.property", "REF_TABLE")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.33.property_value", "LT1")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.34.object_kind", "RELATIONSHIP")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.34.object_name", "R2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.34.parent_entity", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.34.property", "FOREIGN_KEY")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.34.property_value", "[\"A1\",\"A2\"]")),

					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.35.object_kind", "RELATIONSHIP")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.35.object_name", "R2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.35.parent_entity", "LT2")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.35.property", "REF_KEY")),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.35.property_value", "[\"A1\",\"A2\"]")),
				),
			},
			// change externally alter
			{
				PreConfig: func() {
					testClient().SemanticView.Alter(t, sdk.NewAlterSemanticViewRequest(id).WithSetComment(comment))
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(changedComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.SemanticViewShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment).
						HasExtension(""),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.#", "36")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.object_kind", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.object_name", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.parent_entity", "")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.property", "COMMENT")),
					assert.Check(resource.TestCheckResourceAttr(modelComplete.ResourceReference(), "describe_output.0.property_value", changedComment)),
				),
			},
			// change externally create
			{
				PreConfig: func() {
					_, semanticViewCleanup := testClient().SemanticView.CreateWithRequest(t, sdk.NewCreateSemanticViewRequest(id, []sdk.LogicalTableRequest{*lt1Request, *lt2Request}).
						WithSemanticViewMetrics([]sdk.MetricDefinitionRequest{*m1Request, *m2Request}).
						WithComment(comment).WithOrReplace(true))
					t.Cleanup(semanticViewCleanup)
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(changedComment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.SemanticViewShowOutput(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment(changedComment).
						HasExtension(""),
					assert.Check(resource.TestCheckResourceAttr(modelCompleteWithDifferentValues.ResourceReference(), "describe_output.#", "15")),
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
					resourceassert.SemanticViewResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.SemanticViewShowOutput(t, modelBasic.ResourceReference()).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasOwnerRoleType("ROLE").
						HasComment("").
						HasExtension(""),
					assert.Check(resource.TestCheckResourceAttr(modelBasic.ResourceReference(), "describe_output.#", "11")),
				),
			},
		},
	})
}
