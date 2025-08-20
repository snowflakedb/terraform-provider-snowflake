package sdk

import "testing"

func TestSemanticViews_Create(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	logicalTableId := randomSchemaObjectIdentifier()
	// Minimal valid CreateSemanticViewOptions
	defaultOpts := func() *CreateSemanticViewOptions {
		return &CreateSemanticViewOptions{
			name: id,
			logicalTables: []LogicalTable{
				{
					TableName: logicalTableId,
				},
			},
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *CreateSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: conflicting fields for [opts.IfNotExists opts.OrReplace]", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfNotExists = Bool(true)
		opts.OrReplace = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errOneOf("CreateSemanticViewOptions", "IfNotExists", "OrReplace"))
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		opts.OrReplace = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "CREATE OR REPLACE SEMANTIC VIEW %s TABLES (%s)", id.FullyQualifiedName(), logicalTableId.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		logicalTableId1 := randomSchemaObjectIdentifier()
		logicalTableId2 := randomSchemaObjectIdentifier()
		tableAlias1 := "table1"
		tableAlias2 := "table2"
		relationshipAlias1 := "rel1"
		logicalTableComment1 := String("logical table comment 1")
		logicalTableComment2 := String("logical table comment 2")
		factExpression := "fact_sql_expression"
		factName := "fact_name"
		dimensionExpression := "dimension_sql_expression"
		dimensionName := "dimension_name"
		metricExpression := "metric_sql_expression"
		metricName := "metric_name"
		tablesObj := []LogicalTable{
			{
				logicalTableAlias: &LogicalTableAlias{LogicalTableAlias: tableAlias1},
				TableName:         logicalTableId1,
				primaryKeys: &PrimaryKeys{PrimaryKey: []SemanticViewColumn{
					{
						Name: "pk1.1",
					},
					{
						Name: "pk1.2",
					},
				}},
				uniqueKeys: []UniqueKeys{
					{
						Unique: []SemanticViewColumn{
							{
								Name: "uk1.3",
							},
						},
					},
					{
						Unique: []SemanticViewColumn{
							{
								Name: "uk1.4",
							},
						},
					},
				},
				synonyms: &Synonyms{WithSynonyms: []string{"'test1'", "'test2'"}},
				Comment:  logicalTableComment1,
			},
			{
				logicalTableAlias: &LogicalTableAlias{LogicalTableAlias: tableAlias2},
				TableName:         logicalTableId2,
				primaryKeys: &PrimaryKeys{PrimaryKey: []SemanticViewColumn{
					{
						Name: "pk2.1",
					},
					{
						Name: "pk2.2",
					},
				}},
				synonyms: &Synonyms{WithSynonyms: []string{"'test3'", "'test4'"}},
				Comment:  logicalTableComment2,
			},
		}
		relationshipsObj := []SemanticViewRelationship{
			{
				relationshipAlias: &RelationshipAlias{RelationshipAlias: relationshipAlias1},
				tableName:         &RelationshipTableAlias{RelationshipTableAlias: tableAlias1},
				relationshipColumnNames: []SemanticViewColumn{
					{
						Name: "pk1.1",
					},
					{
						Name: "pk1.2",
					},
				},
				refTableName: &RelationshipTableAlias{RelationshipTableAlias: tableAlias2},
				relationshipRefColumnNames: []SemanticViewColumn{
					{
						Name: "pk2.1",
					},
					{
						Name: "pk2.2",
					},
				},
			},
		}
		factsObj := []SemanticExpression{
			{
				qualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: factName},
				sqlExpression:           &SemanticSqlExpression{SqlExpression: factExpression},
				synonyms:                &Synonyms{WithSynonyms: []string{"'test1'", "'test2'"}},
				Comment:                 String("fact_comment"),
			},
		}
		dimensionsObj := []SemanticExpression{
			{
				qualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: dimensionName},
				sqlExpression:           &SemanticSqlExpression{SqlExpression: dimensionExpression},
				synonyms:                &Synonyms{WithSynonyms: []string{"'test3'", "'test4'"}},
				Comment:                 String("dimension_comment"),
			},
		}
		metricsObj := []MetricDefinition{
			{
				semanticExpression: &SemanticExpression{
					qualifiedExpressionName: &QualifiedExpressionName{QualifiedExpressionName: metricName},
					sqlExpression:           &SemanticSqlExpression{SqlExpression: metricExpression},
					synonyms:                &Synonyms{WithSynonyms: []string{"'test5'", "'test6'"}},
					Comment:                 String("metric_comment"),
				},
			},
			{
				windowFunctionMetricDefinition: &WindowFunctionMetricDefinition{
					WindowFunction: "metric1",
					as:             true,
					Metric:         "SUM(table_1.metric_1)",
					OverClause: &WindowFunctionOverClause{
						PartitionBy:       Bool(true),
						PartitionByClause: String("table_1.dimension_2, table_1.dimension_3"),
						OrderBy:           Bool(true),
						OrderByClause:     String("table_1.dimension_2"),
					},
				},
			},
		}

		opts := &CreateSemanticViewOptions{
			name:                      id,
			Comment:                   String("comment"),
			IfNotExists:               Bool(true),
			logicalTables:             tablesObj,
			Relationships:             Bool(true),
			semanticViewRelationships: relationshipsObj,
			Facts:                     Bool(true),
			semanticViewFacts:         factsObj,
			Dimensions:                Bool(true),
			semanticViewDimensions:    dimensionsObj,
			Metrics:                   Bool(true),
			semanticViewMetrics:       metricsObj,
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE SEMANTIC VIEW IF NOT EXISTS %s TABLES (%s AS %s PRIMARY KEY (pk1.1, pk1.2) UNIQUE (uk1.3) UNIQUE (uk1.4) WITH SYNONYMS ('test1', 'test2') COMMENT = '%s', %s AS %s PRIMARY KEY (pk2.1, pk2.2) WITH SYNONYMS ('test3', 'test4') COMMENT = '%s') RELATIONSHIPS (%s AS %s (pk1.1, pk1.2) REFERENCES %s (pk2.1, pk2.2)) FACTS (%s AS %s WITH SYNONYMS ('test1', 'test2') COMMENT = '%s') DIMENSIONS (%s AS %s WITH SYNONYMS ('test3', 'test4') COMMENT = '%s') METRICS (%s AS %s WITH SYNONYMS ('test5', 'test6') COMMENT = '%s', %s AS %s OVER (PARTITION BY %s ORDER BY %s)) COMMENT = '%s'`,
			id.FullyQualifiedName(), tableAlias1, logicalTableId1.FullyQualifiedName(), *logicalTableComment1, tableAlias2, logicalTableId2.FullyQualifiedName(), *logicalTableComment2, relationshipAlias1, tableAlias1, tableAlias2, factName, factExpression, *factsObj[0].Comment, dimensionName, dimensionExpression, *dimensionsObj[0].Comment, metricName, metricExpression, *metricsObj[0].semanticExpression.Comment, metricsObj[1].windowFunctionMetricDefinition.WindowFunction, metricsObj[1].windowFunctionMetricDefinition.Metric, *metricsObj[1].windowFunctionMetricDefinition.OverClause.PartitionByClause, *metricsObj[1].windowFunctionMetricDefinition.OverClause.OrderByClause, "comment")
	})
}

func TestSemanticViews_Drop(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid DropSemanticViewOptions
	defaultOpts := func() *DropSemanticViewOptions {
		return &DropSemanticViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DropSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DROP SEMANTIC VIEW %s", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "DROP SEMANTIC VIEW IF EXISTS %s", id.FullyQualifiedName())
	})
}

func TestSemanticViews_Describe(t *testing.T) {
	id := randomSchemaObjectIdentifier()
	// Minimal valid DescribeSemanticViewOptions
	defaultOpts := func() *DescribeSemanticViewOptions {
		return &DescribeSemanticViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *DescribeSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "DESCRIBE SEMANTIC VIEW %s", id.FullyQualifiedName())
	})
}

func TestSemanticViews_Show(t *testing.T) {
	// Minimal valid ShowSemanticViewOptions
	defaultOpts := func() *ShowSemanticViewOptions {
		return &ShowSemanticViewOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *ShowSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("basic", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, "SHOW SEMANTIC VIEWS")
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.Terse = Bool(true)
		opts.Like = &Like{
			Pattern: String("my_account"),
		}
		opts.In = &In{
			Account: Bool(true),
		}
		opts.StartsWith = String("sem")
		opts.Limit = &LimitFrom{Rows: Int(10)}
		assertOptsValidAndSQLEquals(t, opts, "SHOW TERSE SEMANTIC VIEWS LIKE 'my_account' IN ACCOUNT STARTS WITH 'sem' LIMIT 10")
	})
}

func TestSemanticViews_Alter(t *testing.T) {

	id := randomSchemaObjectIdentifier()
	// Minimal valid AlterSemanticViewOptions
	defaultOpts := func() *AlterSemanticViewOptions {
		return &AlterSemanticViewOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *AlterSemanticViewOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, ErrNilOptions)
	})

	t.Run("validation: invalid identifier for [opts.name]", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: invalid options", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = emptySchemaObjectIdentifier
		opts.SetComment = String("comment")
		opts.UnsetComment = Bool(true)
		assertOptsInvalidJoinedErrors(t, opts, errExactlyOneOf("AlterSemanticViewOptions", "SetComment", "UnsetComment"))
	})

	t.Run("set comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("comment")
		assertOptsValidAndSQLEquals(t, opts, "ALTER SEMANTIC VIEW %s SET COMMENT = 'comment'", id.FullyQualifiedName())
	})

	t.Run("unset comment", func(t *testing.T) {
		opts := defaultOpts()
		opts.UnsetComment = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER SEMANTIC VIEW %s UNSET COMMENT", id.FullyQualifiedName())
	})

	t.Run("all options", func(t *testing.T) {
		opts := defaultOpts()
		opts.SetComment = String("comment")
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, "ALTER SEMANTIC VIEW IF EXISTS %s SET COMMENT = 'comment'", id.FullyQualifiedName())
	})
}
