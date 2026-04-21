//go:build non_account_level_tests

package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func TestInt_SemanticView(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// create 2 tables and add them to the cleanup queue
	columns1 := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest(`"first_a"`, sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		*sdk.NewTableColumnRequest(`"first_b"`, sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		*sdk.NewTableColumnRequest(`"first_c"`, sdk.DataTypeVARCHAR).WithInlineConstraint(sdk.NewColumnInlineConstraintRequest("pkey", sdk.ColumnConstraintTypePrimaryKey)),
	}
	table1Id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("lowercase")
	_, table1Cleanup := testClientHelper().Table.CreateWithRequest(t, sdk.NewCreateTableRequest(table1Id, columns1))
	t.Cleanup(table1Cleanup)

	columns2 := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest(`"second_a"`, sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		*sdk.NewTableColumnRequest(`"second_b"`, sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		*sdk.NewTableColumnRequest(`"second_c"`, sdk.DataTypeVARCHAR).WithInlineConstraint(sdk.NewColumnInlineConstraintRequest("pkey", sdk.ColumnConstraintTypePrimaryKey)),
	}
	table2Id := testClientHelper().Ids.RandomSchemaObjectIdentifierWithPrefix("lowercase")
	_, table2Cleanup := testClientHelper().Table.CreateWithRequest(t, sdk.NewCreateTableRequest(table2Id, columns2))
	t.Cleanup(table2Cleanup)

	// create logical table entities using the 2 tables created above
	alias1 := sdk.NewLogicalTableAliasRequest().WithLogicalTableAlias("table1")
	pk1 := sdk.NewPrimaryKeysRequest().WithPrimaryKey([]sdk.SemanticViewColumn{
		{
			Name: "first_c",
		},
	})
	logicalTable1 := sdk.NewLogicalTableRequest(table1Id).WithLogicalTableAlias(*alias1).WithPrimaryKeys(*pk1)
	alias2 := sdk.NewLogicalTableAliasRequest().WithLogicalTableAlias("table2")
	logicalTable2 := sdk.NewLogicalTableRequest(table2Id).WithLogicalTableAlias(*alias2)

	logicalTables := []sdk.LogicalTableRequest{
		*logicalTable1,
		*logicalTable2,
	}

	// create a simple metric to be used in the semantic view definition
	metricSemanticExpression := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: `"table1"."metric1"`}, &sdk.SemanticSqlExpressionRequest{SqlExpression: `SUM("table1"."first_a")`})
	metric := sdk.NewMetricDefinitionRequest().WithSemanticExpression(*metricSemanticExpression)
	metrics := []sdk.MetricDefinitionRequest{
		*metric,
	}

	t.Run("create: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics).WithComment("comment")

		// create the semantic view with logical tables and a metric
		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		// check that the semantic view was created
		semanticView, err := client.SemanticViews.ShowByID(ctx, id)
		require.NoError(t, err)

		// check that the semantic view's properties match our settings
		assertThatObject(t, objectassert.SemanticViewFromObject(t, semanticView).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasName(id.Name()).
			HasCreatedOnNotEmpty().
			HasOwner("ACCOUNTADMIN").
			HasOwnerRoleType("ROLE").
			HasComment("comment"),
		)
	})

	t.Run("create: table in a different schema", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics)

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		semanticView, err := client.SemanticViews.ShowByID(ctx, id)
		require.NoError(t, err)

		// check that the semantic view's properties match our settings
		assertThatObject(t, objectassert.SemanticViewFromObject(t, semanticView).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasName(id.Name()),
		)
	})

	t.Run("create: without queryable expression", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		logicalTableNoAlias := sdk.NewLogicalTableRequest(table1Id)

		request := sdk.NewCreateSemanticViewRequest(id, []sdk.LogicalTableRequest{*logicalTableNoAlias})

		err := client.SemanticViews.Create(ctx, request)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))
		require.ErrorContains(t, err, "No queryable expression is defined in the semantic view")
	})

	// TODO [SNOW-2852837]: Clarify if creation without table alias is possible
	t.Run("create: without table alias", func(t *testing.T) {
		t.Skip("SNOW-2852837: Skipped as by current docs the table alias is optional but we can't figure out the syntax for metrics without it.")

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		logicalTableNoAlias := sdk.NewLogicalTableRequest(table1Id)
		metricOnNoAliasTable := sdk.NewMetricDefinitionRequest().WithSemanticExpression(*sdk.NewSemanticExpressionRequest(
			&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: fmt.Sprintf(`%s."metric1"`, table1Id.FullyQualifiedName())},
			&sdk.SemanticSqlExpressionRequest{SqlExpression: fmt.Sprintf(`SUM(%s."first_a")`, table1Id.FullyQualifiedName())},
		))

		request := sdk.NewCreateSemanticViewRequest(id, []sdk.LogicalTableRequest{*logicalTableNoAlias}).WithSemanticViewMetrics([]sdk.MetricDefinitionRequest{*metricOnNoAliasTable})

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		semanticView, err := client.SemanticViews.ShowByID(ctx, id)
		require.NoError(t, err)

		// check that the semantic view's properties match our settings
		assertThatObject(t, objectassert.SemanticViewFromObject(t, semanticView).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasName(id.Name()).
			HasCreatedOnNotEmpty().
			HasOwner("ACCOUNTADMIN").
			HasOwnerRoleType("ROLE").
			HasComment("comment"),
		)
	})

	t.Run("create: all fields", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		// relationships
		tableAlias := sdk.NewRelationshipTableAliasRequest().WithRelationshipTableAlias("table2")
		relCol := sdk.NewSemanticViewColumnRequest("second_c")
		relColumnNames := []sdk.SemanticViewColumnRequest{
			*relCol,
		}
		refTableAlias := sdk.NewRelationshipTableAliasRequest().WithRelationshipTableAlias("table1")
		relAliasRequest := sdk.NewRelationshipAliasRequest().WithRelationshipAlias("rel1")
		relRefCol := sdk.NewSemanticViewColumnRequest("first_c")

		relationships := sdk.NewSemanticViewRelationshipRequest(
			tableAlias,
			relColumnNames,
			refTableAlias,
		).WithRelationshipAlias(*relAliasRequest).WithRelationshipRefColumnNames([]sdk.SemanticViewColumnRequest{*relRefCol})

		// facts
		factSynonymRequest := sdk.NewSynonymsRequest().WithWithSynonyms([]sdk.Synonym{{Synonym: "F1"}, {Synonym: "FA"}})
		factSemanticExpression := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: `"table1"."fact1"`}, &sdk.SemanticSqlExpressionRequest{SqlExpression: `"first_c"`}).
			WithSynonyms(*factSynonymRequest).
			WithComment("fact comment")
		fact := sdk.NewFactDefinitionRequest().WithSemanticExpression(*factSemanticExpression)

		factSynonymRequest2 := sdk.NewSynonymsRequest().WithWithSynonyms([]sdk.Synonym{{Synonym: "F2"}})
		factSemanticExpression2 := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: `"table1"."fact2"`}, &sdk.SemanticSqlExpressionRequest{SqlExpression: `"first_b"`}).
			WithSynonyms(*factSynonymRequest2)
		fact2 := sdk.NewFactDefinitionRequest().WithSemanticExpression(*factSemanticExpression2).WithIsPrivate(true)

		// dimensions
		dimensionExpressionNameRaw := `"table1"."d1"`
		dimensionExpressionRaw := `"table1"."first_c"`
		dimensionSynonymRequest := sdk.NewSynonymsRequest().WithWithSynonyms([]sdk.Synonym{{Synonym: "D1"}})
		dimensionSemanticExpression := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: dimensionExpressionNameRaw}, &sdk.SemanticSqlExpressionRequest{SqlExpression: dimensionExpressionRaw}).
			WithSynonyms(*dimensionSynonymRequest).
			WithComment("dimension comment")
		dimension := sdk.NewDimensionDefinitionRequest().WithSemanticExpression(*dimensionSemanticExpression)

		windowFunctionExpression := sdk.NewWindowFunctionMetricDefinitionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: `"table1"."metric2"`}, &sdk.SemanticSqlExpressionRequest{SqlExpression: `SUM("table1"."metric1")`}).WithOverClause(*sdk.NewWindowFunctionOverClauseRequest().WithPartitionBy(`"table1"."d1"`))
		windowFunctionMetric := sdk.NewMetricDefinitionRequest().
			WithWindowFunctionMetricDefinition(*windowFunctionExpression).
			WithIsPrivate(true)

		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).
			WithSemanticViewMetrics([]sdk.MetricDefinitionRequest{*metric, *windowFunctionMetric}).
			WithComment("comment").
			WithSemanticViewRelationships([]sdk.SemanticViewRelationshipRequest{*relationships}).
			WithSemanticViewFacts([]sdk.FactDefinitionRequest{*fact, *fact2}).
			WithSemanticViewDimensions([]sdk.DimensionDefinitionRequest{*dimension})

		// create the semantic view with logical tables and a metric
		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		// check that the semantic view was created
		semanticView, err := client.SemanticViews.ShowByID(ctx, id)
		require.NoError(t, err)

		// check that the semantic view's properties match our settings
		assertThatObject(t, objectassert.SemanticViewFromObject(t, semanticView).
			HasDatabaseName(id.DatabaseName()).
			HasSchemaName(id.SchemaName()).
			HasName(id.Name()).
			HasCreatedOnNotEmpty().
			HasOwner("ACCOUNTADMIN").
			HasOwnerRoleType("ROLE").
			HasComment("comment"),
		)

		t1Alias, t2Alias, dimensionName, factName, factName2, metricName, metric2Name, relationshipName := "table1", "table2", "d1", "fact1", "fact2", "metric1", "metric2", "rel1"

		expectedTable1 := sdk.SemanticViewTableDetails{
			TableNameOrAlias: t1Alias,
			BaseTable:        table1Id,
			PrimaryKeys:      []string{"first_c"},
		}

		expectedTable2 := sdk.SemanticViewTableDetails{
			TableNameOrAlias: t2Alias,
			BaseTable:        table2Id,
		}

		expectedDimension := sdk.SemanticViewDimensionDetails{
			DimensionAlias: dimensionName,
			Properties: sdk.CommonProperties{
				TableNameOrAlias: t1Alias,
				Expression:       dimensionExpressionRaw,
				// TODO [SNOW-2852837]: there is a currently open BCR changing the VARCHAR default size (VARCHAR(16777216) vs VARCHAR(134217728)), update when generally available
				DataType:       "VARCHAR(16777216)",
				AccessModifier: "PUBLIC",
			},
			Synonyms:     []string{"D1"},
			Comment:      "dimension comment",
			ParentEntity: t1Alias,
		}

		expectedFact1 := sdk.SemanticViewFactDetails{
			FactAlias: factName,
			Properties: sdk.CommonProperties{
				TableNameOrAlias: t1Alias,
				Expression:       `"first_c"`,
				// TODO [SNOW-2852837]: there is a currently open BCR changing the VARCHAR default size (VARCHAR(16777216) vs VARCHAR(134217728)), update when generally available
				DataType:       "VARCHAR(16777216)",
				AccessModifier: "PUBLIC",
			},
			Synonyms:     []string{"F1", "FA"},
			Comment:      "fact comment",
			ParentEntity: t1Alias,
		}

		expectedFact2 := sdk.SemanticViewFactDetails{
			FactAlias: factName2,
			Properties: sdk.CommonProperties{
				TableNameOrAlias: t1Alias,
				Expression:       `"first_b"`,
				DataType:         "NUMBER(38,0)",
				AccessModifier:   "PRIVATE",
			},
			Synonyms:     []string{"F2"},
			ParentEntity: t1Alias,
		}

		expectedMetric1 := sdk.SemanticViewMetricDetails{
			MetricAlias: metricName,
			Properties: sdk.CommonProperties{
				TableNameOrAlias: t1Alias,
				Expression:       `SUM("table1"."first_a")`,
				DataType:         "NUMBER(38,0)",
				AccessModifier:   "PUBLIC",
			},
			ParentEntity: t1Alias,
		}

		expectedMetric2 := sdk.SemanticViewMetricDetails{
			MetricAlias: metric2Name,
			Properties: sdk.CommonProperties{
				TableNameOrAlias: t1Alias,
				Expression:       `SUM("table1"."metric1") OVER (PARTITION BY "table1"."d1")`,
				DataType:         "NUMBER(38,0)",
				AccessModifier:   "PRIVATE",
			},
			ParentEntity: t1Alias,
		}

		expectedRelationship := sdk.SemanticViewRelationshipDetails{
			RelationshipAlias:   relationshipName,
			TableNameOrAlias:    t2Alias,
			ForeignKeys:         []string{"second_c"},
			RefTableNameOrAlias: t1Alias,
			RefKeys:             []string{"first_c"},
			ParentEntity:        t2Alias,
		}

		assertThatObject(t, objectassert.SemanticViewDetails(t, id).
			HasDetailsCount(37).
			HasComment("comment").
			ContainsTable(expectedTable1).
			ContainsTable(expectedTable2).
			ContainsDimension(expectedDimension).
			ContainsFact(expectedFact1).
			ContainsFact(expectedFact2).
			ContainsMetric(expectedMetric1).
			ContainsMetric(expectedMetric2).
			ContainsRelationship(expectedRelationship),
		)
	})

	t.Run("create: non-existing table", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateSemanticViewRequest(id, []sdk.LogicalTableRequest{*sdk.NewLogicalTableRequest(NonExistingSchemaObjectIdentifier)}).WithSemanticViewMetrics(metrics)

		err := client.SemanticViews.Create(ctx, request)
		require.ErrorContains(t, err, "object does not exist or not authorized")
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))
	})

	t.Run("describe semantic view", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics)

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		t1Alias, t2Alias, metricName := "table1", "table2", "metric1"

		expectedTable1 := sdk.SemanticViewTableDetails{
			TableNameOrAlias: t1Alias,
			BaseTable:        table1Id,
			PrimaryKeys:      []string{"first_c"},
		}

		expectedMetric := sdk.SemanticViewMetricDetails{
			MetricAlias: metricName,
			Properties: sdk.CommonProperties{
				TableNameOrAlias: t1Alias,
				Expression:       `SUM("table1"."first_a")`,
				DataType:         "NUMBER(38,0)",
				AccessModifier:   "PUBLIC",
			},
			ParentEntity: t1Alias,
		}

		expectedTable2 := sdk.SemanticViewTableDetails{
			TableNameOrAlias: t2Alias,
			BaseTable:        table2Id,
		}

		// confirm the semantic view details are correct
		assertThatObject(t, objectassert.SemanticViewDetails(t, id).
			HasDetailsCount(11).
			ContainsTable(expectedTable1).
			ContainsTable(expectedTable2).
			ContainsMetric(expectedMetric),
		)
	})

	t.Run("alter semantic view", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics).WithComment("comment")

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		semanticView, err := client.SemanticViews.ShowByID(ctx, id)
		require.NoError(t, err)

		// check that the semantic view was created with a comment
		assertThatObject(t, objectassert.SemanticViewFromObject(t, semanticView).
			HasComment("comment"),
		)

		// alter the semantic view to unset the comment
		alterRequest := sdk.NewAlterSemanticViewRequest(id).WithIfExists(true).WithUnsetComment(true)
		err = client.SemanticViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// semantic view should still exist
		semanticView, err = client.SemanticViews.ShowByID(ctx, id)
		require.NoError(t, err)

		// check that the semantic view no longer has a comment
		assertThatObject(t, objectassert.SemanticViewFromObject(t, semanticView).
			HasNoComment(),
		)

		// add a new comment to the semantic view
		alterRequest = sdk.NewAlterSemanticViewRequest(id).WithSetComment("updated comment")
		err = client.SemanticViews.Alter(ctx, alterRequest)
		require.NoError(t, err)

		// semantic view should still exist
		semanticView, err = client.SemanticViews.ShowByID(ctx, id)
		require.NoError(t, err)

		// semantic view should now have the new comment
		assertThatObject(t, objectassert.SemanticViewFromObject(t, semanticView).
			HasComment("updated comment"),
		)
	})

	t.Run("alter: rename", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics).WithComment("comment")

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		alterRequest := sdk.NewAlterSemanticViewRequest(id).WithRenameTo(newId)

		err = client.SemanticViews.Alter(ctx, alterRequest)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, newId))

		result, err := client.SemanticViews.ShowByID(ctx, newId)
		require.NoError(t, err)
		require.Equal(t, newId.Name(), result.Name)

		_, err = client.SemanticViews.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("alter: rename to different schema", func(t *testing.T) {
		secondSchema, secondSchemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(secondSchemaCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics).WithComment("comment")

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())
		alterRequest := sdk.NewAlterSemanticViewRequest(id).WithRenameTo(newId)

		err = client.SemanticViews.Alter(ctx, alterRequest)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, newId))

		result, err := client.SemanticViews.ShowByID(ctx, newId)
		require.NoError(t, err)
		require.Equal(t, newId.Name(), result.Name)
		require.Equal(t, newId.SchemaName(), result.SchemaName)

		_, err = client.SemanticViews.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("alter: rename to different database", func(t *testing.T) {
		secondDatabase, secondDatabaseCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
		t.Cleanup(secondDatabaseCleanup)

		secondSchema, secondSchemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, secondDatabase.ID())
		t.Cleanup(secondSchemaCleanup)

		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics)

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		newId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(secondSchema.ID())
		alterRequest := sdk.NewAlterSemanticViewRequest(id).WithRenameTo(newId)

		err = client.SemanticViews.Alter(ctx, alterRequest)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, newId))

		result, err := client.SemanticViews.ShowByID(ctx, newId)
		require.NoError(t, err)
		require.Equal(t, newId.Name(), result.Name)
		require.Equal(t, newId.SchemaName(), result.SchemaName)
		require.Equal(t, newId.DatabaseName(), result.DatabaseName)

		_, err = client.SemanticViews.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotFound)
	})

	t.Run("drop: basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics)

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		dropRequest := sdk.NewDropSemanticViewRequest(id).WithIfExists(true)
		err = client.SemanticViews.Drop(ctx, dropRequest)
		require.NoError(t, err)

		// confirm the semantic view no longer exists
		_, err = client.SemanticViews.ShowByID(ctx, id)
		require.Error(t, err)

		// with if exists set to true, calling Drop again should not return an error
		err = client.SemanticViews.Drop(ctx, dropRequest)
		require.NoError(t, err)
	})

	t.Run("drop: not existing", func(t *testing.T) {
		dropRequest := sdk.NewDropSemanticViewRequest(NonExistingSchemaObjectIdentifier)
		err := client.SemanticViews.Drop(ctx, dropRequest)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}
