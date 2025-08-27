//go:build account_level_tests

package testint

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func TestInt_SemanticView(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	db, dbCleanup := testClientHelper().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(dbCleanup)

	schema, schemaCleanup := testClientHelper().Schema.CreateSchemaInDatabase(t, db.ID())
	t.Cleanup(schemaCleanup)

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
	alias1 := sdk.NewLogicalTableAliasRequest().WithLogicalTableAlias("table1")
	logicalTable1 := sdk.NewLogicalTableRequest(table1.ID()).WithLogicalTableAlias(*alias1)
	alias2 := sdk.NewLogicalTableAliasRequest().WithLogicalTableAlias("table2")
	logicalTable2 := sdk.NewLogicalTableRequest(table2.ID()).WithLogicalTableAlias(*alias2)

	logicalTables := []sdk.LogicalTableRequest{
		*logicalTable1,
		*logicalTable2,
	}
	semanticExpression := sdk.NewSemanticExpressionRequest(&sdk.QualifiedExpressionNameRequest{QualifiedExpressionName: "table1.metric1"}, &sdk.SemanticSqlExpressionRequest{SqlExpression: "SUM(table1.FIRST_A)"})
	metric := sdk.NewMetricDefinitionRequest().WithSemanticExpression(*semanticExpression)
	metrics := []sdk.MetricDefinitionRequest{
		*metric,
	}

	t.Run("create - basic", func(t *testing.T) {
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics).WithComment("comment")

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		semanticView, err := client.SemanticViews.ShowByID(ctx, id)
		require.NoError(t, err)

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
		id := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
		request := sdk.NewCreateSemanticViewRequest(id, logicalTables).WithSemanticViewMetrics(metrics)

		err := client.SemanticViews.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(testClientHelper().SemanticView.DropFunc(t, id))

		semanticViewDetails, err := client.SemanticViews.Describe(ctx, id)
		require.NoError(t, err)

		parentEntity := "TABLE1"
		expectedDescription := []sdk.SemanticViewDetails{
			{
				ObjectKind:    "TABLE",
				ObjectName:    "TABLE1",
				ParentEntity:  nil,
				Property:      "BASE_TABLE_DATABASE_NAME",
				PropertyValue: table1.DatabaseName,
			},
			{
				ObjectKind:    "TABLE",
				ObjectName:    "TABLE1",
				ParentEntity:  nil,
				Property:      "BASE_TABLE_SCHEMA_NAME",
				PropertyValue: table1.SchemaName,
			},
			{
				ObjectKind:    "TABLE",
				ObjectName:    "TABLE1",
				ParentEntity:  nil,
				Property:      "BASE_TABLE_NAME",
				PropertyValue: table1.Name,
			},
			{
				ObjectKind:    "METRIC",
				ObjectName:    "METRIC1",
				ParentEntity:  &parentEntity,
				Property:      "TABLE",
				PropertyValue: "TABLE1",
			},
			{
				ObjectKind:    "METRIC",
				ObjectName:    "METRIC1",
				ParentEntity:  &parentEntity,
				Property:      "EXPRESSION",
				PropertyValue: "SUM(table1.FIRST_A)",
			},
			{
				ObjectKind:    "METRIC",
				ObjectName:    "METRIC1",
				ParentEntity:  &parentEntity,
				Property:      "DATA_TYPE",
				PropertyValue: "NUMBER(38,0)",
			},
			{
				ObjectKind:    "METRIC",
				ObjectName:    "METRIC1",
				ParentEntity:  &parentEntity,
				Property:      "ACCESS_MODIFIER",
				PropertyValue: "PUBLIC",
			}, {
				ObjectKind:    "TABLE",
				ObjectName:    "TABLE2",
				ParentEntity:  nil,
				Property:      "BASE_TABLE_DATABASE_NAME",
				PropertyValue: table2.DatabaseName,
			},
			{
				ObjectKind:    "TABLE",
				ObjectName:    "TABLE2",
				ParentEntity:  nil,
				Property:      "BASE_TABLE_SCHEMA_NAME",
				PropertyValue: table2.SchemaName,
			},
			{
				ObjectKind:    "TABLE",
				ObjectName:    "TABLE2",
				ParentEntity:  nil,
				Property:      "BASE_TABLE_NAME",
				PropertyValue: table2.Name,
			},
		}

		_ = expectedDescription
		//assert.Equal(t, "TABLE", semanticViewDetails[0].ObjectKind)
		//assert.Equal(t, "TABLE1", semanticViewDetails[0].ObjectName)
		//assert.Equal(t, nil, semanticViewDetails[0].ParentEntity)
		//assert.Equal(t, "BASE_TABLE_DATABASE_NAME", semanticViewDetails[0].Property)
		//assert.Equal(t, table1.DatabaseName, semanticViewDetails[0].PropertyValue)

		assert.Equal(t, expectedDescription, semanticViewDetails)
	})
}
