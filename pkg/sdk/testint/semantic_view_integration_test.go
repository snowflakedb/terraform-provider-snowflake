//go:build non_account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func TestInt_SemanticView(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// create 2 tables and add them to the cleanup queue
	columns1 := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("FIRST_A", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		*sdk.NewTableColumnRequest("FIRST_B", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		*sdk.NewTableColumnRequest("FIRST_C", sdk.DataTypeVARCHAR).WithInlineConstraint(sdk.NewColumnInlineConstraintRequest("pkey", sdk.ColumnConstraintTypePrimaryKey)),
	}
	table1ID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
	table1, table1Cleanup := testClientHelper().Table.CreateWithRequest(t, sdk.NewCreateTableRequest(table1ID, columns1))
	t.Cleanup(table1Cleanup)

	columns2 := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("SECOND_A", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		*sdk.NewTableColumnRequest("SECOND_B", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		*sdk.NewTableColumnRequest("SECOND_C", sdk.DataTypeVARCHAR).WithInlineConstraint(sdk.NewColumnInlineConstraintRequest("pkey", sdk.ColumnConstraintTypePrimaryKey)),
	}
	table2ID := testClientHelper().Ids.RandomSchemaObjectIdentifier()
	table2, table2Cleanup := testClientHelper().Table.CreateWithRequest(t, sdk.NewCreateTableRequest(table2ID, columns2))
	t.Cleanup(table2Cleanup)

	// create logical table entities using the 2 tables created above
	alias1 := sdk.NewLogicalTableAliasRequest().WithLogicalTableAlias("table1")
	pk1 := sdk.NewPrimaryKeysRequest().WithPrimaryKey([]sdk.SemanticViewColumn{
		{
			Name: "FIRST_C",
		},
	})
	logicalTable1 := sdk.NewLogicalTableRequest(table1.ID()).WithLogicalTableAlias(*alias1).WithPrimaryKeys(*pk1)
	alias2 := sdk.NewLogicalTableAliasRequest().WithLogicalTableAlias("table2")
	logicalTable2 := sdk.NewLogicalTableRequest(table2.ID()).WithLogicalTableAlias(*alias2)

	logicalTables := []sdk.LogicalTableRequest{
		*logicalTable1,
		*logicalTable2,
	}

	// create a simple metric to be used in the semantic view definition
	metricSemanticExpression := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: "table1.metric1"}, &sdk.SemanticSqlExpressionRequest{SqlExpression: "SUM(table1.FIRST_A)"})
	metric := sdk.NewMetricDefinitionRequest().WithSemanticExpression(*metricSemanticExpression)
	metrics := []sdk.MetricDefinitionRequest{
		*metric,
	}

	t.Run("create and show", func(t *testing.T) {
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

	t.Run("create - all fields", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()

		// relationships
		tableAlias := sdk.NewRelationshipTableAliasRequest().WithRelationshipTableAlias("table2")
		relCol := sdk.NewSemanticViewColumnRequest("SECOND_C")
		relColumnNames := []sdk.SemanticViewColumnRequest{
			*relCol,
		}
		refTableAlias := sdk.NewRelationshipTableAliasRequest().WithRelationshipTableAlias("table1")
		relAliasRequest := sdk.NewRelationshipAliasRequest().WithRelationshipAlias("rel1")
		relRefCol := sdk.NewSemanticViewColumnRequest("FIRST_C")

		relationships := sdk.NewSemanticViewRelationshipRequest(
			tableAlias,
			relColumnNames,
			refTableAlias,
		).WithRelationshipAlias(*relAliasRequest).WithRelationshipRefColumnNames([]sdk.SemanticViewColumnRequest{*relRefCol})

		// facts
		factSynonymRequest := sdk.NewSynonymsRequest().WithWithSynonyms([]sdk.Synonym{{Synonym: "F1"}})
		factSemanticExpression := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: "table1.fact1"}, &sdk.SemanticSqlExpressionRequest{SqlExpression: "FIRST_C"}).
			WithSynonyms(*factSynonymRequest).
			WithComment("fact comment")

		// dimensions
		dimensionSynonymRequest := sdk.NewSynonymsRequest().WithWithSynonyms([]sdk.Synonym{{Synonym: "D1"}})
		dimensionSemanticExpression := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: "table1.FIRST_C"}, &sdk.SemanticSqlExpressionRequest{SqlExpression: "table1.FIRST_C"}).
			WithSynonyms(*dimensionSynonymRequest).
			WithComment("dimension comment")

		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).
			WithSemanticViewMetrics(metrics).
			WithComment("comment").
			WithSemanticViewRelationships([]sdk.SemanticViewRelationshipRequest{*relationships}).
			WithSemanticViewFacts([]sdk.SemanticExpressionRequest{*factSemanticExpression}).
			WithSemanticViewDimensions([]sdk.SemanticExpressionRequest{*dimensionSemanticExpression})

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

	t.Run("describe semantic view", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifier()
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics)

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		tableKind, metricKind := "TABLE", "METRIC"
		t1Name, t2Name, metricName := "TABLE1", "TABLE2", "METRIC1"

		tableDatabaseName1 := objectassert.NewSemanticViewDetails(&tableKind, &t1Name, nil, "BASE_TABLE_DATABASE_NAME", table1.DatabaseName)
		tableSchemaName1 := objectassert.NewSemanticViewDetails(&tableKind, &t1Name, nil, "BASE_TABLE_SCHEMA_NAME", table1.SchemaName)
		tableName1 := objectassert.NewSemanticViewDetails(&tableKind, &t1Name, nil, "BASE_TABLE_NAME", table1.Name)
		pk := objectassert.NewSemanticViewDetails(&tableKind, &t1Name, nil, "PRIMARY_KEY", "[\"FIRST_C\"]")
		metricTable := objectassert.NewSemanticViewDetails(&metricKind, &metricName, &t1Name, "TABLE", "TABLE1")
		metricExpression := objectassert.NewSemanticViewDetails(&metricKind, &metricName, &t1Name, "EXPRESSION", "SUM(table1.FIRST_A)")
		metricDataType := objectassert.NewSemanticViewDetails(&metricKind, &metricName, &t1Name, "DATA_TYPE", "NUMBER(38,0)")
		metricAccessModifier := objectassert.NewSemanticViewDetails(&metricKind, &metricName, &t1Name, "ACCESS_MODIFIER", "PUBLIC")
		tableDatabaseName2 := objectassert.NewSemanticViewDetails(&tableKind, &t2Name, nil, "BASE_TABLE_DATABASE_NAME", table2.DatabaseName)
		tableSchemaName2 := objectassert.NewSemanticViewDetails(&tableKind, &t2Name, nil, "BASE_TABLE_SCHEMA_NAME", table2.SchemaName)
		tableName2 := objectassert.NewSemanticViewDetails(&tableKind, &t2Name, nil, "BASE_TABLE_NAME", table2.Name)

		// confirm the semantic view details are correct
		assertThatObject(t, objectassert.SemanticViewDetails(t, id).
			HasDetailsCount(11).
			ContainsDetail(tableDatabaseName1).
			ContainsDetail(tableSchemaName1).
			ContainsDetail(tableName1).
			ContainsDetail(pk).
			ContainsDetail(metricTable).
			ContainsDetail(metricExpression).
			ContainsDetail(metricDataType).
			ContainsDetail(metricAccessModifier).
			ContainsDetail(tableDatabaseName2).
			ContainsDetail(tableSchemaName2).
			ContainsDetail(tableName2),
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

	t.Run("drop semantic view", func(t *testing.T) {
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
}
