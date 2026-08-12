package sdk

import "fmt"

func init() {
	id := semanticViewsTestIdSchemaObjectIdentifier
	logicalTableId := randomSchemaObjectIdentifier()

	semanticViewsTests.Create.
		withDefaultOpts(func() *CreateSemanticViewOptions {
			return &CreateSemanticViewOptions{
				name: id,
				LogicalTables: []LogicalTable{
					{TableName: logicalTableId},
				},
			}
		}).
		withModify(
			case_SemanticViews_validation_Create_opts_SemanticViewMetrics_ExactlyOneValueSet_MoreThanOneSet,
			func(opts *CreateSemanticViewOptions) {
				opts.SemanticViewMetrics = []MetricDefinition{
					{
						SemanticExpression:             &SemanticExpression{},
						WindowFunctionMetricDefinition: &WindowFunctionMetricDefinition{},
					},
				}
			},
		).
		withModify(
			case_SemanticViews_validation_Create_opts_SemanticViewMetrics_ExactlyOneValueSet_OneValidOneInvalid,
			func(opts *CreateSemanticViewOptions) {
				opts.SemanticViewMetrics = []MetricDefinition{
					{SemanticExpression: &SemanticExpression{}},
					{},
				}
			},
		).
		withExpectedSqlf(
			case_SemanticViews_sql_Create_basic,
			`CREATE SEMANTIC VIEW %s TABLES (%s)`,
			id.FullyQualifiedName(), logicalTableId.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Create_orReplace",
			func(opts *CreateSemanticViewOptions) { opts.OrReplace = new(true) },
			`CREATE OR REPLACE SEMANTIC VIEW %s TABLES (%s)`,
			id.FullyQualifiedName(), logicalTableId.FullyQualifiedName(),
		)

	logicalTableId1 := randomSchemaObjectIdentifier()
	logicalTableId2 := randomSchemaObjectIdentifier()
	tableAlias1 := "table1"
	tableAlias2 := "table2"
	relationshipAlias1 := "rel1"
	logicalTableComment1 := "logical table comment 1"
	logicalTableComment2 := "logical table comment 2"
	factExpression := "fact_sql_expression"
	factName := "fact_name"
	dimensionExpression := "dimension_sql_expression"
	dimensionName := "dimension_name"
	metricExpression := "metric_sql_expression"
	metricName := `"table1"."metric_name"`

	tablesObj := []LogicalTable{
		{
			LogicalTableAlias: &LogicalTableAlias{LogicalTableAlias: tableAlias1},
			TableName:         logicalTableId1,
			PrimaryKeys: &PrimaryKeys{PrimaryKey: []SemanticViewColumn{
				{Name: "pk1.1"},
				{Name: "pk1.2"},
			}},
			UniqueKeys: []UniqueKeys{
				{Unique: []SemanticViewColumn{{Name: "uk1.3"}}},
				{Unique: []SemanticViewColumn{{Name: "uk1.4"}}},
			},
			Synonyms: &Synonyms{WithSynonyms: []Synonym{{Synonym: "test1"}, {Synonym: "test2"}}},
			Comment:  new(logicalTableComment1),
		},
		{
			LogicalTableAlias: &LogicalTableAlias{LogicalTableAlias: tableAlias2},
			TableName:         logicalTableId2,
			PrimaryKeys: &PrimaryKeys{PrimaryKey: []SemanticViewColumn{
				{Name: "pk2.1"},
				{Name: "pk2.2"},
			}},
			Synonyms: &Synonyms{WithSynonyms: []Synonym{{Synonym: "test3"}, {Synonym: "test4"}}},
			Comment:  new(logicalTableComment2),
		},
	}
	relationshipsObj := []SemanticViewRelationship{
		{
			RelationshipAlias: &RelationshipAlias{RelationshipAlias: relationshipAlias1},
			TableNameOrAlias:  &RelationshipTableAlias{RelationshipTableAlias: new(tableAlias1)},
			RelationshipColumnNames: []SemanticViewColumn{
				{Name: "pk1.1"},
				{Name: "pk1.2"},
			},
			RefTableNameOrAlias: &RelationshipTableAlias{RelationshipTableAlias: new(tableAlias2)},
			RelationshipRefColumnNames: []SemanticViewColumn{
				{Name: "pk2.1"},
				{Name: "pk2.2"},
			},
		},
	}
	factsObj := []FactDefinition{
		{
			IsPrivate: new(true),
			SemanticExpression: &SemanticExpression{
				QualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: fmt.Sprintf(`"%s"`, factName)},
				SqlExpression:           &SemanticSqlExpression{SqlExpression: factExpression},
				Synonyms:                &Synonyms{WithSynonyms: []Synonym{{Synonym: "test1"}, {Synonym: "test2"}}},
				Comment:                 new("fact_comment"),
			},
		},
	}
	dimensionsObj := []DimensionDefinition{
		{
			SemanticExpression: &SemanticExpression{
				QualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: fmt.Sprintf(`"%s"`, dimensionName)},
				SqlExpression:           &SemanticSqlExpression{SqlExpression: dimensionExpression},
				Synonyms:                &Synonyms{WithSynonyms: []Synonym{{Synonym: "test3"}, {Synonym: "test4"}}},
				Comment:                 new("dimension_comment"),
			},
		},
	}
	metricsObj := []MetricDefinition{
		{
			IsPrivate: new(true),
			SemanticExpression: &SemanticExpression{
				QualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: metricName},
				SqlExpression:           &SemanticSqlExpression{SqlExpression: metricExpression},
				Synonyms:                &Synonyms{WithSynonyms: []Synonym{{Synonym: "test5"}, {Synonym: "test6"}}},
				Comment:                 new("metric_comment"),
			},
		},
		{
			WindowFunctionMetricDefinition: &WindowFunctionMetricDefinition{
				QualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: `"table1"."metric1"`},
				SqlExpression:           &SemanticSqlExpression{SqlExpression: fmt.Sprintf(`SUM(%s)`, metricName)},
				OverClause: &WindowFunctionOverClause{
					PartitionBy: new("table_1.dimension_2, table_1.dimension_3"),
					OrderBy:     new("table_1.dimension_2"),
				},
			},
		},
	}

	semanticViewsTests.Create.
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Create_all,
			func(opts *CreateSemanticViewOptions) {
				opts.Comment = new("comment")
				opts.IfNotExists = new(true)
				opts.LogicalTables = tablesObj
				opts.SemanticViewRelationships = relationshipsObj
				opts.SemanticViewFacts = factsObj
				opts.SemanticViewDimensions = dimensionsObj
				opts.SemanticViewMetrics = metricsObj
			},
			`CREATE SEMANTIC VIEW IF NOT EXISTS %s TABLES ("%s" AS %s PRIMARY KEY ("pk1.1", "pk1.2") UNIQUE ("uk1.3") UNIQUE ("uk1.4") WITH SYNONYMS ('test1', 'test2') COMMENT = '%s', "%s" AS %s PRIMARY KEY ("pk2.1", "pk2.2") WITH SYNONYMS ('test3', 'test4') COMMENT = '%s') RELATIONSHIPS ("%s" AS "%s" ("pk1.1", "pk1.2") REFERENCES "%s" ("pk2.1", "pk2.2")) FACTS (PRIVATE "%s" AS %s WITH SYNONYMS ('test1', 'test2') COMMENT = '%s') DIMENSIONS ("%s" AS %s WITH SYNONYMS ('test3', 'test4') COMMENT = '%s') METRICS (PRIVATE %s AS %s WITH SYNONYMS ('test5', 'test6') COMMENT = '%s', %s AS %s OVER (PARTITION BY %s ORDER BY %s)) COMMENT = '%s'`,
			id.FullyQualifiedName(),
			tableAlias1,
			logicalTableId1.FullyQualifiedName(),
			logicalTableComment1,
			tableAlias2,
			logicalTableId2.FullyQualifiedName(),
			logicalTableComment2,
			relationshipAlias1,
			tableAlias1,
			tableAlias2,
			factName,
			factExpression,
			"fact_comment",
			dimensionName,
			dimensionExpression,
			"dimension_comment",
			metricName,
			metricExpression,
			"metric_comment",
			metricsObj[1].WindowFunctionMetricDefinition.QualifiedExpressionName.QualifiedExpressionName,
			metricsObj[1].WindowFunctionMetricDefinition.SqlExpression.SqlExpression,
			*metricsObj[1].WindowFunctionMetricDefinition.OverClause.PartitionBy,
			*metricsObj[1].WindowFunctionMetricDefinition.OverClause.OrderBy,
			"comment",
		)

	renameToId := randomSchemaObjectIdentifier()

	semanticViewsTests.Alter.
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Alter_SetComment,
			func(opts *AlterSemanticViewOptions) { opts.SetComment = new("comment") },
			`ALTER SEMANTIC VIEW %s SET COMMENT = 'comment'`, id.FullyQualifiedName(),
		).
		withAdditionalSqlCasef(
			"sql_Alter_SetComment_ifExists",
			func(opts *AlterSemanticViewOptions) {
				opts.IfExists = new(true)
				opts.SetComment = new("comment")
			},
			`ALTER SEMANTIC VIEW IF EXISTS %s SET COMMENT = 'comment'`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Alter_UnsetComment,
			func(opts *AlterSemanticViewOptions) { opts.UnsetComment = new(true) },
			`ALTER SEMANTIC VIEW %s UNSET COMMENT`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Alter_RenameTo,
			func(opts *AlterSemanticViewOptions) { opts.RenameTo = new(renameToId) },
			`ALTER SEMANTIC VIEW %s RENAME TO %s`, id.FullyQualifiedName(), renameToId.FullyQualifiedName(),
		)

	semanticViewsTests.Drop.
		withExpectedSqlf(
			case_SemanticViews_sql_Drop_basic,
			`DROP SEMANTIC VIEW %s`, id.FullyQualifiedName(),
		).
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Drop_all,
			func(opts *DropSemanticViewOptions) { opts.IfExists = new(true) },
			`DROP SEMANTIC VIEW IF EXISTS %s`, id.FullyQualifiedName(),
		)

	semanticViewsTests.Describe.
		withExpectedSqlf(
			case_SemanticViews_sql_Describe_basic,
			`DESCRIBE SEMANTIC VIEW %s`, id.FullyQualifiedName(),
		)

	semanticViewsTests.Show.
		withExpectedSqlf(
			case_SemanticViews_sql_Show_basic,
			`SHOW SEMANTIC VIEWS`,
		).
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Show_all,
			func(opts *ShowSemanticViewOptions) {
				opts.Terse = new(true)
				opts.Like = &Like{Pattern: new("my_account")}
				opts.In = &In{Account: new(true)}
				opts.StartsWith = new("sem")
				opts.Limit = &LimitFrom{Rows: new(10)}
			},
			`SHOW TERSE SEMANTIC VIEWS LIKE 'my_account' IN ACCOUNT STARTS WITH 'sem' LIMIT 10`,
		).
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Show_Like,
			func(opts *ShowSemanticViewOptions) { opts.Like = &Like{Pattern: new("my_account")} },
			`SHOW SEMANTIC VIEWS LIKE 'my_account'`,
		).
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Show_In,
			func(opts *ShowSemanticViewOptions) { opts.In = &In{Account: new(true)} },
			`SHOW SEMANTIC VIEWS IN ACCOUNT`,
		).
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Show_StartsWith,
			func(opts *ShowSemanticViewOptions) { opts.StartsWith = new("sem") },
			`SHOW SEMANTIC VIEWS STARTS WITH 'sem'`,
		).
		withModifyAndExpectedSqlf(
			case_SemanticViews_sql_Show_Limit,
			func(opts *ShowSemanticViewOptions) { opts.Limit = &LimitFrom{Rows: new(10)} },
			`SHOW SEMANTIC VIEWS LIMIT 10`,
		)
}
