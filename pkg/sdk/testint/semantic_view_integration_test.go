//go:build non_account_level_tests

package testint

import (
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
		factSynonymRequest := sdk.NewSynonymsRequest().WithWithSynonyms([]sdk.Synonym{{Synonym: "F1"}})
		factSemanticExpression := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: `"table1"."fact1"`}, &sdk.SemanticSqlExpressionRequest{SqlExpression: `"first_c"`}).
			WithSynonyms(*factSynonymRequest).
			WithComment("fact comment")

		// dimensions
		dimensionSynonymRequest := sdk.NewSynonymsRequest().WithWithSynonyms([]sdk.Synonym{{Synonym: "D1"}})
		dimensionSemanticExpression := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: `"table1"."first_c"`}, &sdk.SemanticSqlExpressionRequest{SqlExpression: `"table1"."first_c"`}).
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

		tableKind, metricKind := "TABLE", "METRIC"
		t1Name, t2Name, metricName := "table1", "table2", "metric1"

		tableDatabaseName1 := objectassert.NewSemanticViewDetails(&tableKind, &t1Name, nil, "BASE_TABLE_DATABASE_NAME", table1Id.DatabaseName())
		tableSchemaName1 := objectassert.NewSemanticViewDetails(&tableKind, &t1Name, nil, "BASE_TABLE_SCHEMA_NAME", table1Id.SchemaName())
		tableName1 := objectassert.NewSemanticViewDetails(&tableKind, &t1Name, nil, "BASE_TABLE_NAME", table1Id.Name())
		pk := objectassert.NewSemanticViewDetails(&tableKind, &t1Name, nil, "PRIMARY_KEY", "[\"first_c\"]")
		metricTable := objectassert.NewSemanticViewDetails(&metricKind, &metricName, &t1Name, "TABLE", "table1")
		metricExpression := objectassert.NewSemanticViewDetails(&metricKind, &metricName, &t1Name, "EXPRESSION", `SUM("table1"."first_a")`) // alias
		metricDataType := objectassert.NewSemanticViewDetails(&metricKind, &metricName, &t1Name, "DATA_TYPE", "NUMBER(38,0)")
		metricAccessModifier := objectassert.NewSemanticViewDetails(&metricKind, &metricName, &t1Name, "ACCESS_MODIFIER", "PUBLIC")
		tableDatabaseName2 := objectassert.NewSemanticViewDetails(&tableKind, &t2Name, nil, "BASE_TABLE_DATABASE_NAME", table2Id.DatabaseName())
		tableSchemaName2 := objectassert.NewSemanticViewDetails(&tableKind, &t2Name, nil, "BASE_TABLE_SCHEMA_NAME", table2Id.SchemaName())
		tableName2 := objectassert.NewSemanticViewDetails(&tableKind, &t2Name, nil, "BASE_TABLE_NAME", table2Id.Name())

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

		_, err := client.SemanticViews.ShowByID(ctx, id)
		require.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
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
