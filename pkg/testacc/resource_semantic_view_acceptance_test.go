//go:build !account_level_tests

package testacc

import (
	"testing"

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
	logicalTable1 := model.LogicalTableWithProps("lt1", table1.ID(), nil, nil, nil, "logical table 1")
	logicalTable2 := model.LogicalTableWithProps("lt2", table2.ID(), nil, nil, nil, "")
	semExp1 := model.SemanticExpressionWithProps("lt1.se1", "SUM(lt1.a1)", []sdk.Synonym{{Synonym: "sem1"}, {Synonym: "baseSem"}}, "semantic expression 1")
	partitionBy := "a1"
	windowFunc1 := model.WindowFunctionMetricDefinitionWithProps("lt2.wf1", "sum(lt2.a1)", sdk.WindowFunctionOverClause{PartitionBy: &partitionBy})
	metric1 := model.MetricDefinitionWithProps(semExp1, nil)
	metric2 := model.MetricDefinitionWithProps(nil, windowFunc1)

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
		[]sdk.LogicalTable{*logicalTable1},
		[]sdk.MetricDefinition{*metric1},
	).WithComment(comment)

	modelCompleteWithDifferentValues := model.SemanticView(
		"test",
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		[]sdk.LogicalTable{*logicalTable2},
		[]sdk.MetricDefinition{*metric2},
	).WithComment(changedComment)

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
				),
			},
			// import minimal state
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
				),
			},
			// alter
			{
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
				),
			},
			// unset
			{
				Config: accconfig.FromModels(t, modelBasic),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelBasic.ResourceReference(), plancheck.ResourceActionUpdate),
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
				),
			},
		},
	})
}
