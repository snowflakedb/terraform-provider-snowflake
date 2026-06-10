package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var semanticViewPairs = g.StructPair("semanticViewDBRow", "SemanticView").
	Time("created_on").
	Text("name").
	Text("database_name").
	Text("schema_name").
	OptionalText("comment").
	Text("owner").
	Text("owner_role_type").
	OptionalText("extension")

var semanticViewDetailsPairs = g.StructPair("semanticViewDetailsRow", "SemanticViewDetails").
	OptionalText("object_kind").
	OptionalText("object_name").
	OptionalText("parent_entity").
	Text("property").
	Text("property_value")

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
			ListQueryStructField("LogicalTables", logicalTable, g.ParameterOptions().Required().Parentheses().NoEquals().SQL("TABLES")).
			ListQueryStructField("SemanticViewRelationships", semanticViewRelationship, g.ParameterOptions().Parentheses().NoEquals().SQL("RELATIONSHIPS")).
			ListQueryStructField("SemanticViewFacts", factDefinition, g.ParameterOptions().Parentheses().NoEquals().SQL("FACTS")).
			ListQueryStructField("SemanticViewDimensions", dimensionDefinition, g.ParameterOptions().Parentheses().NoEquals().SQL("DIMENSIONS")).
			ListQueryStructField("SemanticViewMetrics", metricDefinition, g.ParameterOptions().Parentheses().NoEquals().SQL("METRICS")).
			OptionalComment().
			OptionalCopyGrants().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace").
			WithAdditionalValidations(),
		synonym,
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
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSlice,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-semantic-view",
		semanticViewDetailsPairs,
		g.NewQueryStruct("DescribeSemanticView").
			Describe().
			SQL("SEMANTIC VIEW").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views",
		semanticViewPairs,
		g.NewQueryStruct("ShowSemanticViews").
			Show().
			Terse().
			SQL("SEMANTIC VIEWS").
			OptionalLike().
			OptionalIn().
			OptionalStartsWith().
			OptionalLimitFrom(),
		g.ShowByIDInFiltering,
		g.ShowByIDLikeFiltering,
	).
	WithCustomInterfaceMethod(
		"DescribeSemanticViewDetails",
		"DescribeSemanticViewDetails returns converted describe output for semantic views.",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"*SemanticViewDescribeDetails", "error",
	)

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
	OptionalQueryStructField("LogicalTableAlias", logicalTableAlias, g.KeywordOptions()).
	Identifier("TableName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
	OptionalQueryStructField("PrimaryKeys", primaryKey, g.ParameterOptions().NoEquals()).
	ListQueryStructField("UniqueKeys", uniqueKey, g.ListOptions().NoEquals().NoComma()).
	OptionalQueryStructField("Synonyms", synonyms, g.ParameterOptions().NoEquals()).
	OptionalComment()

var relationshipAlias = g.NewQueryStruct("RelationshipAlias").
	Text("RelationshipAlias", g.KeywordOptions().DoubleQuotes().Required()).
	SQL("AS")

var relationshipTableNameOrAlias = g.NewQueryStruct("RelationshipTableAlias").
	OptionalIdentifier("RelationshipTableName", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions()).
	OptionalText("RelationshipTableAlias", g.KeywordOptions().DoubleQuotes()).
	WithValidation(g.ExactlyOneValueSet, "RelationshipTableName", "RelationshipTableAlias")

var semanticViewRelationship = g.NewQueryStruct("SemanticViewRelationship").
	OptionalQueryStructField("RelationshipAlias", relationshipAlias, g.KeywordOptions()).
	OptionalQueryStructField("TableNameOrAlias", relationshipTableNameOrAlias, g.KeywordOptions().Required()).
	ListQueryStructField("RelationshipColumnNames", semanticViewColumn, g.ListOptions().NoEquals().Parentheses().Required()).
	SQL("REFERENCES").
	OptionalQueryStructField("RefTableNameOrAlias", relationshipTableNameOrAlias, g.KeywordOptions().Required()).
	ListQueryStructField("RelationshipRefColumnNames", semanticViewColumn, g.ListOptions().NoEquals().Parentheses())

func qualifiedExpressionName() *g.QueryStruct {
	return g.NewQueryStruct("QualifiedExpressionName").
		Text("QualifiedExpressionName", g.KeywordOptions().Required())
}

var semanticSqlExpression = g.NewQueryStruct("SemanticSqlExpression").
	Text("SqlExpression", g.KeywordOptions().NoQuotes().Required())

// TODO [SNOW-2398097]: replace qualifiedExpressionName with table_alias and fact_or_metric fields
func semanticExpression() *g.QueryStruct {
	return g.NewQueryStruct("SemanticExpression").
		OptionalQueryStructField("QualifiedExpressionName", qualifiedExpressionName(), g.KeywordOptions().Required()).
		SQL("AS").
		OptionalQueryStructField("SqlExpression", semanticSqlExpression, g.KeywordOptions().Required()).
		OptionalQueryStructField("Synonyms", synonyms, g.ParameterOptions().NoEquals()).
		OptionalComment()
}

var windowFunctionOverClause = g.NewQueryStruct("WindowFunctionOverClause").
	OptionalTextAssignment("PARTITION BY", g.ParameterOptions().NoEquals()).
	OptionalTextAssignment("ORDER BY", g.ParameterOptions().NoEquals()).
	OptionalText("WindowFrameClause", g.KeywordOptions())

// TODO [SNOW-2398097]: sqlExpression could be replaced with <window_function>(<metric>)
// TODO [SNOW-2398097]: windowFunctionMetricDefinition could be merged with semanticExpression to have syntax for metrics definition (different than for facts and dimensions)
var windowFunctionMetricDefinition = g.NewQueryStruct("WindowFunctionMetricDefinition").
	OptionalQueryStructField("QualifiedExpressionName", qualifiedExpressionName(), g.KeywordOptions().Required()).
	SQL("AS").
	OptionalQueryStructField("SqlExpression", semanticSqlExpression, g.KeywordOptions().Required()).
	OptionalQueryStructField("OverClause", windowFunctionOverClause, g.ListOptions().Parentheses().NoComma().SQL("OVER"))

var metricDefinition = g.NewQueryStruct("MetricDefinition").
	OptionalSQLWithCustomFieldName("IsPrivate", "PRIVATE").
	OptionalQueryStructField("SemanticExpression", semanticExpression(), g.KeywordOptions()).
	OptionalQueryStructField("WindowFunctionMetricDefinition", windowFunctionMetricDefinition, g.KeywordOptions()).
	WithValidation(g.ExactlyOneValueSet, "SemanticExpression", "WindowFunctionMetricDefinition")

var factDefinition = g.NewQueryStruct("FactDefinition").
	OptionalSQLWithCustomFieldName("IsPrivate", "PRIVATE").
	OptionalQueryStructField("SemanticExpression", semanticExpression(), g.KeywordOptions())

var dimensionDefinition = g.NewQueryStruct("DimensionDefinition").
	OptionalQueryStructField("SemanticExpression", semanticExpression(), g.KeywordOptions())
