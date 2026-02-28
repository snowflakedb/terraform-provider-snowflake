//go:build non_account_level_tests

package testacc

import (
	"regexp"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/invokeactionassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
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
		*sdk.NewTableColumnRequest(`"a1"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a2"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a3"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a4"`, sdk.DataTypeNumber),
	})
	t.Cleanup(table1Cleanup)
	table2, table2Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest(`"a1"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a2"`, sdk.DataTypeNumber),
	})
	t.Cleanup(table2Cleanup)

	logicalTable1 := sdk.LogicalTable{}
	logicalTable1.WithLogicalTableAlias("lt1").
		WithTableName(table1.ID()).
		WithPrimaryKeys([]sdk.SemanticViewColumn{{Name: "a1"}}).
		WithUniqueKeys([][]sdk.SemanticViewColumn{{{Name: "a2"}}, {{Name: "a3"}, {Name: "a4"}}}).
		WithSynonyms([]sdk.Synonym{{"orders"}, {"sales"}}).
		WithComment("logical table 1")

	logicalTable2 := sdk.LogicalTable{}
	logicalTable2.WithLogicalTableAlias("lt2").
		WithTableName(table2.ID()).
		WithPrimaryKeys([]sdk.SemanticViewColumn{{Name: "a1"}})

	semExp1 := sdk.SemanticExpression{}
	semExp1.WithQualifiedExpressionName(`"lt1"."m1"`).
		WithSqlExpression(`SUM("lt1"."a1")`).
		WithSynonyms([]sdk.Synonym{{Synonym: "sem1"}, {Synonym: "baseSem"}}).
		WithComment("semantic expression 1")

	metric1 := sdk.MetricDefinition{}
	metric1.WithSemanticExpression(&semExp1).
		WithIsPrivate(true)

	relTableAlias := sdk.RelationshipTableAlias{}
	relTableAlias.WithRelationshipTableAlias("lt1").
		WithRelationshipTableName(table1.ID())

	relTableColumns := []sdk.SemanticViewColumn{
		{
			Name: "a1",
		},
		{
			Name: "a2",
		},
	}

	refTableAlias := sdk.RelationshipTableAlias{}
	refTableAlias.WithRelationshipTableAlias("lt2").
		WithRelationshipTableName(table2.ID())
	refRelTableColumns := []sdk.SemanticViewColumn{
		{
			Name: "a1",
		},
		{
			Name: "a2",
		},
	}

	rel1 := sdk.SemanticViewRelationship{}
	rel1.WithRelationshipAlias("r1").
		WithTableNameOrAlias(relTableAlias).
		WithRelationshipColumnsNames(relTableColumns).
		WithRefTableNameOrAlias(refTableAlias).
		WithRelationshipRefColumnsNames(refRelTableColumns)

	factSemExp1 := sdk.SemanticExpression{}
	factSemExp1.WithQualifiedExpressionName(`"lt1"."f1"`).
		WithSqlExpression(`"lt1"."a2"`).
		WithSynonyms([]sdk.Synonym{{Synonym: "fact1"}}).
		WithComment("fact 1")

	fact1 := sdk.FactDefinition{}
	fact1.WithSemanticExpression(&factSemExp1).
		WithIsPrivate(false)

	dimensionSemExp1 := sdk.SemanticExpression{}
	dimensionSemExp1.WithQualifiedExpressionName(`"lt1"."d1"`).
		WithSqlExpression(`"lt1"."a1"`).
		WithSynonyms([]sdk.Synonym{{Synonym: "dim1"}}).
		WithComment("dimension 1")

	dimension1 := sdk.DimensionDefinition{}
	dimension1.WithSemanticExpression(&dimensionSemExp1)

	relTableAlias2 := sdk.RelationshipTableAlias{}
	relTableAlias2.WithRelationshipTableAlias("lt2").
		WithRelationshipTableName(table1.ID())

	relTableColumns2 := []sdk.SemanticViewColumn{
		{
			Name: "a1",
		},
		{
			Name: "a2",
		},
	}

	refTableAlias2 := sdk.RelationshipTableAlias{}
	refTableAlias2.WithRelationshipTableAlias("lt1").
		WithRelationshipTableName(table2.ID())

	refRelTableColumns2 := []sdk.SemanticViewColumn{
		{
			Name: "a1",
		},
		{
			Name: "a2",
		},
	}

	rel2 := sdk.SemanticViewRelationship{}
	rel2.WithRelationshipAlias("r2").
		WithTableNameOrAlias(relTableAlias2).
		WithRelationshipColumnsNames(relTableColumns2).
		WithRefTableNameOrAlias(refTableAlias2).
		WithRelationshipRefColumnsNames(refRelTableColumns2)

	factSemExp2 := sdk.SemanticExpression{}
	factSemExp2.WithQualifiedExpressionName(`"lt1"."f2"`).
		WithSqlExpression(`"lt1"."a1"`).
		WithSynonyms([]sdk.Synonym{{Synonym: "fact2"}}).
		WithComment("fact 2")

	fact2 := sdk.FactDefinition{}
	fact2.WithSemanticExpression(&factSemExp2).
		WithIsPrivate(true)

	dimensionSemExp2 := sdk.SemanticExpression{}
	dimensionSemExp2.WithQualifiedExpressionName(`"lt1"."d2"`).
		WithSqlExpression(`"lt1"."a2"`).
		WithSynonyms([]sdk.Synonym{{Synonym: "dim2"}}).
		WithComment("dimension 2")

	dimension2 := sdk.DimensionDefinition{}
	dimension2.WithSemanticExpression(&dimensionSemExp2)

	windowFunc1 := sdk.WindowFunctionMetricDefinition{}
	windowFunc1.WithQualifiedExpressionName(`"lt1"."wf1"`).
		WithSqlExpression(`SUM("lt1"."m1")`).
		WithOverClause(sdk.WindowFunctionOverClause{PartitionBy: sdk.Pointer(`"lt1"."d2"`)})

	metric2 := sdk.MetricDefinition{}
	metric2.WithWindowFunctionMetricDefinition(&windowFunc1).
		WithIsPrivate(false)

	lt1Request := sdk.NewLogicalTableRequest(table1.ID()).WithLogicalTableAlias(sdk.LogicalTableAliasRequest{LogicalTableAlias: "lt1"})
	lt2Request := sdk.NewLogicalTableRequest(table2.ID()).WithLogicalTableAlias(sdk.LogicalTableAliasRequest{LogicalTableAlias: "lt2"})
	seRequest := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: `"lt1"."m2"`}, &sdk.SemanticSqlExpressionRequest{SqlExpression: `SUM("lt1"."a1")`})
	wfRequest := sdk.NewWindowFunctionMetricDefinitionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: `"lt1"."wf2"`}, &sdk.SemanticSqlExpressionRequest{SqlExpression: `SUM("lt1"."m2")`}).WithOverClause(*sdk.NewWindowFunctionOverClauseRequest().WithPartitionBy(`"lt1"."d1"`))
	m1Request := sdk.NewMetricDefinitionRequest().WithSemanticExpression(*seRequest)
	m2Request := sdk.NewMetricDefinitionRequest().WithWindowFunctionMetricDefinition(*wfRequest)
	dseRequest := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: `"lt1"."d1"`}, &sdk.SemanticSqlExpressionRequest{SqlExpression: `"lt1"."a2"`})
	d1Request := sdk.NewDimensionDefinitionRequest().WithSemanticExpression(*dseRequest)

	modelBasic := model.SemanticViewWithMetrics(
		"test",
		id,
		[]sdk.LogicalTable{logicalTable1},
		[]sdk.MetricDefinition{metric1},
	)

	modelComplete := model.SemanticViewWithMetrics(
		"test",
		id,
		[]sdk.LogicalTable{logicalTable1, logicalTable2},
		[]sdk.MetricDefinition{metric1},
	).WithComment(comment).
		WithRelationships([]sdk.SemanticViewRelationship{rel1}).
		WithFacts([]sdk.FactDefinition{fact1}).
		WithDimensions([]sdk.DimensionDefinition{dimension1})

	modelCompleteWithDifferentValues := model.SemanticViewWithMetrics(
		"test",
		id,
		[]sdk.LogicalTable{logicalTable1, logicalTable2},
		[]sdk.MetricDefinition{metric1, metric2},
	).WithComment(changedComment).
		WithRelationships([]sdk.SemanticViewRelationship{rel2}).
		WithFacts([]sdk.FactDefinition{fact2}).
		WithDimensions([]sdk.DimensionDefinition{dimension2})

	t1Alias, t2Alias, dimensionName, factName, privateFactName, metricName, relationshipName := "lt1", "lt2", "d1", "f1", "f2", "m1", "r1"

	// semantic view related details
	commentDetails := objectassert.NewSemanticViewDetails(nil, nil, nil, "COMMENT", comment)

	// logical table 1 related details
	table1DatabaseName := objectassert.NewSemanticViewDetailsTable(t1Alias, "BASE_TABLE_DATABASE_NAME", table1.ID().DatabaseName())
	table1SchemaName := objectassert.NewSemanticViewDetailsTable(t1Alias, "BASE_TABLE_SCHEMA_NAME", table1.ID().SchemaName())
	table1Name := objectassert.NewSemanticViewDetailsTable(t1Alias, "BASE_TABLE_NAME", table1.ID().Name())
	table1Synonyms := objectassert.NewSemanticViewDetailsTable(t1Alias, "SYNONYMS", `["sales","orders"]`)
	table1PrimaryKey := objectassert.NewSemanticViewDetailsTable(t1Alias, "PRIMARY_KEY", `["a1"]`)
	table1UniqueKey := objectassert.NewSemanticViewDetailsTable(t1Alias, "UNIQUE_KEY", `[["a2"],["a3","a4"]]`)
	table1Comment := objectassert.NewSemanticViewDetailsTable(t1Alias, "COMMENT", `logical table 1`)

	// dimension related details
	dimensionTable := objectassert.NewSemanticViewDetailsDimension(dimensionName, t1Alias, "TABLE", t1Alias)
	dimensionExpression := objectassert.NewSemanticViewDetailsDimension(dimensionName, t1Alias, "EXPRESSION", `"lt1"."a1"`)
	dimensionDataType := objectassert.NewSemanticViewDetailsDimension(dimensionName, t1Alias, "DATA_TYPE", "NUMBER(38,0)")
	dimensionSynonyms := objectassert.NewSemanticViewDetailsDimension(dimensionName, t1Alias, "SYNONYMS", `["dim1"]`)
	dimensionComment := objectassert.NewSemanticViewDetailsDimension(dimensionName, t1Alias, "COMMENT", "dimension 1")
	dimensionAccessModifier := objectassert.NewSemanticViewDetailsDimension(dimensionName, t1Alias, "ACCESS_MODIFIER", "PUBLIC")

	// fact related details
	factTable := objectassert.NewSemanticViewDetailsFact(factName, t1Alias, "TABLE", t1Alias)
	factExpression := objectassert.NewSemanticViewDetailsFact(factName, t1Alias, "EXPRESSION", `"lt1"."a2"`)
	factDataType := objectassert.NewSemanticViewDetailsFact(factName, t1Alias, "DATA_TYPE", "NUMBER(38,0)")
	factSynonyms := objectassert.NewSemanticViewDetailsFact(factName, t1Alias, "SYNONYMS", `["fact1"]`)
	factComment := objectassert.NewSemanticViewDetailsFact(factName, t1Alias, "COMMENT", "fact 1")
	factAccessModifier := objectassert.NewSemanticViewDetailsFact(factName, t1Alias, "ACCESS_MODIFIER", "PUBLIC")

	privateFactTable := objectassert.NewSemanticViewDetailsFact(privateFactName, t1Alias, "TABLE", t1Alias)
	privateFactExpression := objectassert.NewSemanticViewDetailsFact(privateFactName, t1Alias, "EXPRESSION", `"lt1"."a1"`)
	privateFactDataType := objectassert.NewSemanticViewDetailsFact(privateFactName, t1Alias, "DATA_TYPE", "NUMBER(38,0)")
	privateFactSynonyms := objectassert.NewSemanticViewDetailsFact(privateFactName, t1Alias, "SYNONYMS", `["fact2"]`)
	privateFactComment := objectassert.NewSemanticViewDetailsFact(privateFactName, t1Alias, "COMMENT", "fact 2")
	privateFactAccessModifier := objectassert.NewSemanticViewDetailsFact(privateFactName, t1Alias, "ACCESS_MODIFIER", "PRIVATE")

	// metric related details
	metricTable := objectassert.NewSemanticViewDetailsMetric(metricName, t1Alias, "TABLE", t1Alias)
	metricExpression := objectassert.NewSemanticViewDetailsMetric(metricName, t1Alias, "EXPRESSION", `SUM("lt1"."a1")`)
	metricDataType := objectassert.NewSemanticViewDetailsMetric(metricName, t1Alias, "DATA_TYPE", "NUMBER(38,0)")
	metricAccessModifier := objectassert.NewSemanticViewDetailsMetric(metricName, t1Alias, "ACCESS_MODIFIER", "PRIVATE")

	// logical table 2 related details
	table2DatabaseName := objectassert.NewSemanticViewDetailsTable(t2Alias, "BASE_TABLE_DATABASE_NAME", table2.ID().DatabaseName())
	table2SchemaName := objectassert.NewSemanticViewDetailsTable(t2Alias, "BASE_TABLE_SCHEMA_NAME", table2.ID().SchemaName())
	table2Name := objectassert.NewSemanticViewDetailsTable(t2Alias, "BASE_TABLE_NAME", table2.ID().Name())

	// relationship related details
	relationshipTable := objectassert.NewSemanticViewDetailsRelationship(relationshipName, t1Alias, "TABLE", t1Alias)
	relationshipRefTable := objectassert.NewSemanticViewDetailsRelationship(relationshipName, t1Alias, "REF_TABLE", t2Alias)
	relationshipForeignKey := objectassert.NewSemanticViewDetailsRelationship(relationshipName, t1Alias, "FOREIGN_KEY", `["a1","a2"]`)
	relationshipRefKey := objectassert.NewSemanticViewDetailsRelationship(relationshipName, t1Alias, "REF_KEY", `["a1","a2"]`)

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
						HasDimensionsEmpty().
						HasFactsEmpty().
						HasRelationshipsEmpty().
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
						HasDimensionsEmpty().
						HasFactsEmpty().
						HasMetricsEmpty().
						HasRelationshipsEmpty().
						HasNoTables().
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
			// add optional attributes - recreate
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelComplete.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: accconfig.FromModels(t, modelComplete),
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, modelComplete.ResourceReference()).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					objectassert.SemanticViewDetails(t, id).
						HasDetailsCount(34).
						ContainsDetail(commentDetails).
						ContainsDetail(table1DatabaseName).
						ContainsDetail(table1SchemaName).
						ContainsDetail(table1Name).
						ContainsDetail(table1Synonyms).
						ContainsDetail(table1PrimaryKey).
						ContainsDetail(table1UniqueKey).
						ContainsDetail(table1Comment).
						ContainsDetail(dimensionTable).
						ContainsDetail(dimensionExpression).
						ContainsDetail(dimensionDataType).
						ContainsDetail(dimensionSynonyms).
						ContainsDetail(dimensionComment).
						ContainsDetail(dimensionAccessModifier).
						ContainsDetail(factTable).
						ContainsDetail(factExpression).
						ContainsDetail(factDataType).
						ContainsDetail(factSynonyms).
						ContainsDetail(factComment).
						ContainsDetail(factAccessModifier).
						ContainsDetail(metricTable).
						ContainsDetail(metricExpression).
						ContainsDetail(metricDataType).
						ContainsDetail(metricAccessModifier).
						ContainsDetail(table2DatabaseName).
						ContainsDetail(table2SchemaName).
						ContainsDetail(table2Name).
						ContainsDetail(relationshipTable).
						ContainsDetail(relationshipRefTable).
						ContainsDetail(relationshipForeignKey).
						ContainsDetail(relationshipRefKey),
				),
			},
			// import complete
			{
				Config:       accconfig.FromModels(t, modelComplete),
				ResourceName: modelComplete.ResourceReference(),
				ImportState:  true,
				ImportStateCheck: assertThatImport(t,
					resourceassert.ImportedSemanticViewResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasCommentString(comment).
						HasDimensionsEmpty().
						HasFactsEmpty().
						HasMetricsEmpty().
						HasRelationshipsEmpty().
						HasNoTables().
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.ImportedSemanticViewShowOutput(t, helpers.EncodeResourceIdentifier(id)).
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
			// change values in config - recreate
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
						HasComment(changedComment),
				),
			},
			// change externally - alter
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
						HasComment(changedComment),
					objectassert.SemanticViewDetails(t, id).
						ContainsDetail(privateFactTable).
						ContainsDetail(privateFactExpression).
						ContainsDetail(privateFactDataType).
						ContainsDetail(privateFactSynonyms).
						ContainsDetail(privateFactComment).
						ContainsDetail(privateFactAccessModifier),
				),
			},
			// change externally - no recreation yet
			// TODO [SNOW-2852837]: Handle external changes
			{
				PreConfig: func() {

					_, semanticViewCleanup := testClient().SemanticView.CreateWithRequest(t, sdk.NewCreateSemanticViewRequest(id, []sdk.LogicalTableRequest{*lt1Request, *lt2Request}).
						WithSemanticViewMetrics([]sdk.MetricDefinitionRequest{*m1Request, *m2Request}).
						WithSemanticViewDimensions([]sdk.DimensionDefinitionRequest{*d1Request}).
						WithComment(changedComment).WithOrReplace(true))
					t.Cleanup(semanticViewCleanup)
				},
				Config: accconfig.FromModels(t, modelCompleteWithDifferentValues),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelCompleteWithDifferentValues.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, modelCompleteWithDifferentValues.ResourceReference()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// recreate to basic
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
						HasDimensionsEmpty().
						HasFactsEmpty().
						HasRelationshipsEmpty().
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_SemanticView_Rename(t *testing.T) {
	secondSchema, secondSchemaCleanup := testClient().Schema.CreateSchema(t)
	t.Cleanup(secondSchemaCleanup)

	secondDatabase, secondDatabaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(secondDatabaseCleanup)

	schemaInSecondDatabase, schemaInSecondDatabaseCleanup := testClient().Schema.CreateSchemaInDatabase(t, secondDatabase.ID())
	t.Cleanup(schemaInSecondDatabaseCleanup)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	newId := testClient().Ids.RandomSchemaObjectIdentifier()
	newIdInDifferentSchema := testClient().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())
	newIdInDifferentDatabase := testClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaInSecondDatabase.ID())

	table1, table1Cleanup := testClient().Table.CreateWithColumns(t, []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest(`"a1"`, sdk.DataTypeNumber),
		*sdk.NewTableColumnRequest(`"a2"`, sdk.DataTypeNumber),
	})
	t.Cleanup(table1Cleanup)

	logicalTable1 := sdk.LogicalTable{}
	logicalTable1.WithLogicalTableAlias("lt1").
		WithTableName(table1.ID()).
		WithPrimaryKeys([]sdk.SemanticViewColumn{{Name: "a1"}}).
		WithUniqueKeys([][]sdk.SemanticViewColumn{{{Name: "a2"}}})

	semExp1 := sdk.SemanticExpression{}
	semExp1.WithQualifiedExpressionName(`"lt1"."se1"`).
		WithSqlExpression(`SUM("lt1"."a1")`)

	metric1 := sdk.MetricDefinition{}
	metric1.WithSemanticExpression(&semExp1).
		WithIsPrivate(false)

	modelBasic := model.SemanticViewWithMetrics(
		"test",
		id,
		[]sdk.LogicalTable{logicalTable1},
		[]sdk.MetricDefinition{metric1},
	).WithComment("old comment")

	renamedAndChanged := model.SemanticViewWithMetrics(
		"test",
		newId,
		[]sdk.LogicalTable{logicalTable1},
		[]sdk.MetricDefinition{metric1},
	).WithComment("new comment")

	renamedDifferentSchema := model.SemanticViewWithMetrics(
		"test",
		newIdInDifferentSchema,
		[]sdk.LogicalTable{logicalTable1},
		[]sdk.MetricDefinition{metric1},
	).WithComment("new comment")

	renamedDifferentDatabase := model.SemanticViewWithMetrics(
		"test",
		newIdInDifferentDatabase,
		[]sdk.LogicalTable{logicalTable1},
		[]sdk.MetricDefinition{metric1},
	).WithComment("new comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: CheckDestroy(t, resources.SemanticView),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, modelBasic),
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, modelBasic.ResourceReference()).
						HasNameString(id.Name()).
						HasCommentString("old comment").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
				),
			},
			// rename with one param changed
			{
				Config: accconfig.FromModels(t, renamedAndChanged),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(renamedAndChanged.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, renamedAndChanged.ResourceReference()).
						HasNameString(newId.Name()).
						HasCommentString("new comment").
						HasFullyQualifiedNameString(newId.FullyQualifiedName()),
					invokeactionassert.SemanticViewDoesNotExist(t, id),
				),
			},
			// rename - different schema
			{
				Config: accconfig.FromModels(t, renamedDifferentSchema),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(renamedDifferentSchema.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, renamedDifferentSchema.ResourceReference()).
						HasNameString(newIdInDifferentSchema.Name()).
						HasCommentString("new comment").
						HasFullyQualifiedNameString(newIdInDifferentSchema.FullyQualifiedName()),
					invokeactionassert.SemanticViewDoesNotExist(t, newId),
				),
			},
			// rename - different database
			{
				Config: accconfig.FromModels(t, renamedDifferentDatabase),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(renamedDifferentDatabase.ResourceReference(), plancheck.ResourceActionUpdate),
					},
				},
				Check: assertThat(t,
					resourceassert.SemanticViewResource(t, renamedDifferentDatabase.ResourceReference()).
						HasNameString(newIdInDifferentDatabase.Name()).
						HasCommentString("new comment").
						HasFullyQualifiedNameString(newIdInDifferentDatabase.FullyQualifiedName()),
					invokeactionassert.SemanticViewDoesNotExist(t, newIdInDifferentSchema),
				),
			},
		},
	})
}

func TestAcc_SemanticView_Validations(t *testing.T) {
	id := testClient().Ids.RandomSchemaObjectIdentifier()
	tableId := testClient().Ids.RandomSchemaObjectIdentifier()
	logicalTable1 := sdk.LogicalTable{}
	logicalTable1.WithLogicalTableAlias("lt1").
		WithTableName(tableId).
		WithPrimaryKeys([]sdk.SemanticViewColumn{{Name: "a1"}}).
		WithUniqueKeys([][]sdk.SemanticViewColumn{{{Name: "a2"}}, {{Name: "a3"}, {Name: "a4"}}}).
		WithComment("logical table 1")

	modelWithoutMetricNorDimension := model.SemanticView(
		"test",
		id.DatabaseName(),
		id.SchemaName(),
		id.Name(),
		[]sdk.LogicalTable{logicalTable1},
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:      accconfig.FromModels(t, modelWithoutMetricNorDimension),
				ExpectError: regexp.MustCompile("one of `dimensions,metrics` must be specified"),
			},
		},
	})
}
