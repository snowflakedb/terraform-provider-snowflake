package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var semanticViewDbRow = g.DbStruct("semanticViewDBRow").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	OptionalText("comment").
	Text("owner").
	Text("owner_role_type").
	OptionalText("extension")

var semanticView = g.PlainStruct("SemanticView").
	Time("CreatedOn").
	Text("Name").
	Text("DatabaseName").
	Text("SchemaName").
	OptionalText("Comment").
	Text("Owner").
	Text("OwnerRoleType").
	OptionalText("Extension")

var semanticViewDetailsDbRow = g.DbStruct("semanticViewDetailsRow").
	OptionalText("object_kind").
	OptionalText("object_name").
	OptionalText("parent_entity").
	Text("property").
	Text("property_value")

var semanticViewDetails = g.PlainStruct("SemanticViewDetails").
	OptionalText("ObjectKind").
	OptionalText("ObjectName").
	OptionalText("ParentEntity").
	Text("Property").
	Text("PropertyValue")

var semanticViewsDef = g.NewInterface(
	"SemanticViews",
	"SemanticView",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-semantic-view",
		g.NewQueryStruct("CreateSemanticView").
			Create().
			OrReplace().
			SQL("SEMANTIC VIEW").
			IfNotExists().
			Name().
			PredefinedQueryStructField("logicalTables", "[]LogicalTable", g.ParameterOptions().Required().Parentheses().NoEquals().SQL("TABLES")).
			PredefinedQueryStructField("semanticViewRelationships", "[]SemanticViewRelationship", g.ParameterOptions().Parentheses().NoEquals().SQL("RELATIONSHIPS")).
			PredefinedQueryStructField("semanticViewFacts", "[]FactDefinition", g.ParameterOptions().Parentheses().NoEquals().SQL("FACTS")).
			PredefinedQueryStructField("semanticViewDimensions", "[]DimensionDefinition", g.ParameterOptions().Parentheses().NoEquals().SQL("DIMENSIONS")).
			PredefinedQueryStructField("semanticViewMetrics", "[]MetricDefinition", g.ParameterOptions().Parentheses().NoEquals().SQL("METRICS")).
			OptionalComment().
			OptionalCopyGrants().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
		logicalTable,
		synonym,
		semanticViewRelationship,
		semanticExpression,
		factDefinition,
		dimensionDefinition,
		metricDefinition,
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-semantic-view",
		g.NewQueryStruct("AlterSemanticView").
			Alter().
			SQL("SEMANTIC VIEW").
			IfExists().
			Name().
			OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalSQL("UNSET COMMENT").
			OptionalIdentifier("RenameTo", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "SetComment", "UnsetComment", "RenameTo"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-semantic-view",
		g.NewQueryStruct("DropSemanticView").
			Drop().
			SQL("SEMANTIC VIEW").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	DescribeOperation(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-semantic-view",
		semanticViewDetailsDbRow,
		semanticViewDetails,
		g.NewQueryStruct("DescribeSemanticView").
			Describe().
			SQL("SEMANTIC VIEW").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views",
		semanticViewDbRow,
		semanticView,
		g.NewQueryStruct("ShowSemanticViews").
			Show().
			Terse().
			SQL("SEMANTIC VIEWS").
			OptionalLike().
			OptionalIn().
			OptionalStartsWith().
			OptionalLimitFrom(),
	).
	ShowByIdOperationWithFiltering(g.ShowByIDInFiltering, g.ShowByIDLikeFiltering)

var primaryKey = g.NewQueryStruct("PrimaryKeys").
	ListAssignment("PRIMARY KEY", "SemanticViewColumn", g.ParameterOptions().Parentheses().NoEquals().Required())

var uniqueKey = g.NewQueryStruct("UniqueKeys").
	ListAssignment("UNIQUE", "SemanticViewColumn", g.ParameterOptions().Parentheses().NoEquals().Required())

var synonym = g.NewQueryStruct("Synonym").
	Text("Synonym", g.KeywordOptions().SingleQuotes().Required())

var synonyms = g.NewQueryStruct("Synonyms").
	ListAssignment("WITH SYNONYMS", "Synonym", g.ParameterOptions().NoEquals().Parentheses().Required())

var logicalTableAlias = g.NewQueryStruct("LogicalTableAlias").
	Text("LogicalTableAlias", g.KeywordOptions().DoubleQuotes().Required()).
	SQL("AS")

var semanticViewColumn = g.NewQueryStruct("SemanticViewColumn").
	Text("Name", g.KeywordOptions().DoubleQuotes().Required())

var logicalTable = g.NewQueryStruct("LogicalTable").
	OptionalQueryStructField("logicalTableAlias", logicalTableAlias, g.KeywordOptions()).
	Identifier("TableName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
	OptionalQueryStructField("primaryKeys", primaryKey, g.ParameterOptions().NoEquals()).
	ListQueryStructField("uniqueKeys", uniqueKey, g.ListOptions().NoEquals().NoComma()).
	OptionalQueryStructField("synonyms", synonyms, g.ParameterOptions().NoEquals()).
	OptionalComment()

var relationshipAlias = g.NewQueryStruct("RelationshipAlias").
	Text("RelationshipAlias", g.KeywordOptions().DoubleQuotes().Required()).
	SQL("AS")

var relationshipTableNameOrAlias = g.NewQueryStruct("RelationshipTableAlias").
	OptionalIdentifier("RelationshipTableName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions()).
	OptionalText("RelationshipTableAlias", g.KeywordOptions().DoubleQuotes()).
	WithValidation(g.ExactlyOneValueSet, "RelationshipTableName", "RelationshipTableAlias")

var semanticViewRelationship = g.NewQueryStruct("SemanticViewRelationship").
	OptionalQueryStructField("relationshipAlias", relationshipAlias, g.KeywordOptions()).
	OptionalQueryStructField("tableNameOrAlias", relationshipTableNameOrAlias, g.KeywordOptions().Required()).
	ListQueryStructField("relationshipColumnNames", semanticViewColumn, g.ListOptions().NoEquals().Parentheses().Required()).
	SQL("REFERENCES").
	OptionalQueryStructField("refTableNameOrAlias", relationshipTableNameOrAlias, g.KeywordOptions().Required()).
	ListQueryStructField("relationshipRefColumnNames", semanticViewColumn, g.ListOptions().NoEquals().Parentheses())

var qualifiedExpressionName = g.NewQueryStruct("QualifiedExpressionName").
	Text("QualifiedExpressionName", g.KeywordOptions().Required())

var semanticSqlExpression = g.NewQueryStruct("SemanticSqlExpression").
	Text("SqlExpression", g.KeywordOptions().NoQuotes().Required())

// TODO [SNOW-2398097]: replace qualifiedExpressionName with table_alias and fact_or_metric fields
var semanticExpression = g.NewQueryStruct("SemanticExpression").
	OptionalQueryStructField("qualifiedExpressionName", qualifiedExpressionName, g.KeywordOptions().Required()).
	SQL("AS").
	OptionalQueryStructField("sqlExpression", semanticSqlExpression, g.KeywordOptions().Required()).
	OptionalQueryStructField("synonyms", synonyms, g.ParameterOptions().NoEquals()).
	OptionalComment()

var windowFunctionOverClause = g.NewQueryStruct("WindowFunctionOverClause").
	OptionalTextAssignment("PARTITION BY", g.ParameterOptions().NoEquals()).
	OptionalTextAssignment("ORDER BY", g.ParameterOptions().NoEquals()).
	OptionalText("WindowFrameClause", g.KeywordOptions())

// TODO [SNOW-2398097]: sqlExpression could be replaced with <window_function>(<metric>)
// TODO [SNOW-2398097]: windowFunctionMetricDefinition could be merged with semanticExpression to have syntax for metrics definition (different than for facts and dimensions)
var windowFunctionMetricDefinition = g.NewQueryStruct("WindowFunctionMetricDefinition").
	OptionalQueryStructField("qualifiedExpressionName", qualifiedExpressionName, g.KeywordOptions().Required()).
	SQL("AS").
	OptionalQueryStructField("sqlExpression", semanticSqlExpression, g.KeywordOptions().Required()).
	OptionalQueryStructField("OverClause", windowFunctionOverClause, g.ListOptions().Parentheses().NoComma().SQL("OVER"))

var metricDefinition = g.NewQueryStruct("MetricDefinition").
	OptionalText("IsPrivate", g.KeywordOptions()).
	OptionalQueryStructField("semanticExpression", semanticExpression, g.KeywordOptions()).
	OptionalQueryStructField("windowFunctionMetricDefinition", windowFunctionMetricDefinition, g.KeywordOptions()).
	WithValidation(g.ExactlyOneValueSet, "semanticExpression", "windowFunctionMetricDefinition")

var factDefinition = g.NewQueryStruct("FactDefinition").
	OptionalText("IsPrivate", g.KeywordOptions()).
	OptionalQueryStructField("semanticExpression", semanticExpression, g.KeywordOptions())

var dimensionDefinition = g.NewQueryStruct("DimensionDefinition").
	OptionalQueryStructField("semanticExpression", semanticExpression, g.KeywordOptions())
