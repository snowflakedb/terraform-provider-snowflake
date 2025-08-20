package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

//go:generate go run ./poc/main.go

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
	Text("object_kind").
	Text("object_name").
	Text("parent_entity").
	Text("property").
	Text("property_value")

var semanticViewDetails = g.PlainStruct("SemanticViewDetails").
	Text("ObjectKind").
	Text("ObjectName").
	Text("ParentEntity").
	Text("Property").
	Text("PropertyValue")

var SemanticViewsDef = g.NewInterface(
	"SemanticViews",
	"SemanticView",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-semantic-view",
	g.NewQueryStruct("CreateSemanticView").
		Create().
		OrReplace().
		SQL("SEMANTIC VIEW").
		IfNotExists().
		Name().
		SQL("TABLES").
		ListQueryStructField("logicalTables", logicalTable, g.ListOptions().Required().Parentheses()).
		OptionalSQL("RELATIONSHIPS").
		ListQueryStructField("semanticViewRelationships", semanticViewRelationship, g.ListOptions().Parentheses()).
		OptionalSQL("FACTS").
		ListQueryStructField("semanticViewFacts", semanticExpression, g.ListOptions().Parentheses()).
		OptionalSQL("DIMENSIONS").
		ListQueryStructField("semanticViewDimensions", semanticExpression, g.ListOptions().Parentheses()).
		OptionalSQL("METRICS").
		ListQueryStructField("semanticViewMetrics", metricDefinition, g.ListOptions().Parentheses()).
		OptionalComment().
		OptionalCopyGrants().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"), // both can't be used at the same time
	logicalTable,
	semanticViewRelationship,
	semanticExpression,
	metricDefinition,
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-semantic-view",
	g.NewQueryStruct("DropSemanticView").
		Drop().
		SQL("SEMANTIC VIEW").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-semantic-view",
	semanticViewDetailsDbRow,
	semanticViewDetails,
	g.NewQueryStruct("DescribeSemanticView").
		Describe().
		SQL("SEMANTIC VIEW").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
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
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-semantic-view",
	g.NewQueryStruct("AlterSemanticView").
		Alter().
		SQL("SEMANTIC VIEW").
		IfExists().
		Name().
		OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalSQL("UNSET COMMENT").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "SetComment", "UnsetComment"),
)

var primaryKey = g.NewQueryStruct("PrimaryKeys").
	ListAssignment("PRIMARY KEY", "SemanticViewColumn", g.ParameterOptions().Parentheses().NoEquals())

var uniqueKey = g.NewQueryStruct("UniqueKeys").
	ListAssignment("UNIQUE", "SemanticViewColumn", g.ParameterOptions().Parentheses().NoEquals())

var synonym = g.NewQueryStruct("Synonyms").
	ListAssignment("WITH SYNONYMS", "string", g.ParameterOptions().NoEquals().Parentheses())

var logicalTableAlias = g.NewQueryStruct("LogicalTableAlias").
	Text("LogicalTableAlias", g.KeywordOptions()).
	SQL("AS")

var semanticViewColumn = g.NewQueryStruct("SemanticViewColumn").
	Text("Name", g.KeywordOptions().Required())

var logicalTable = g.NewQueryStruct("LogicalTable").
	OptionalQueryStructField("logicalTableAlias", logicalTableAlias, g.KeywordOptions()).
	Identifier("TableName", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
	OptionalQueryStructField("primaryKeys", primaryKey, g.ParameterOptions().NoEquals()).
	ListQueryStructField("uniqueKeys", uniqueKey, g.ListOptions().NoEquals().NoComma()).
	OptionalQueryStructField("synonyms", synonym, g.ParameterOptions().NoEquals()).
	OptionalComment()

var relationshipAlias = g.NewQueryStruct("RelationshipAlias").
	Text("RelationshipAlias", g.KeywordOptions()).
	SQL("AS")

var relationshipTableAlias = g.NewQueryStruct("RelationshipTableAlias").
	Text("RelationshipTableAlias", g.KeywordOptions())

var semanticViewRelationship = g.NewQueryStruct("SemanticViewRelationship").
	OptionalQueryStructField("relationshipAlias", relationshipAlias, g.KeywordOptions()).
	OptionalQueryStructField("tableName", relationshipTableAlias, g.KeywordOptions().Required()).
	ListQueryStructField("relationshipColumnNames", semanticViewColumn, g.ListOptions().NoEquals().Parentheses().Required()).
	SQL("REFERENCES").
	OptionalQueryStructField("refTableName", relationshipTableAlias, g.KeywordOptions().Required()).
	ListQueryStructField("relationshipRefColumnNames", semanticViewColumn, g.ListOptions().NoEquals().Parentheses())

var qualifiedExpressionName = g.NewQueryStruct("QualifiedExpressionName").
	Text("QualifiedExpressionName", g.KeywordOptions())

var semanticSqlExpression = g.NewQueryStruct("SemanticSqlExpression").
	Text("SqlExpression", g.KeywordOptions().NoQuotes())

var semanticExpression = g.NewQueryStruct("SemanticExpression").
	OptionalQueryStructField("qualifiedExpressionName", qualifiedExpressionName, g.KeywordOptions().Required()).
	SQL("AS").
	OptionalQueryStructField("sqlExpression", semanticSqlExpression, g.KeywordOptions().Required()).
	OptionalQueryStructField("synonyms", synonym, g.ParameterOptions().NoEquals()).
	OptionalComment()

var windowFunctionOverClause = g.NewQueryStruct("WindowFunctionOverClause").
	SQL("PARTITION BY").
	OptionalText("PartitionByClause", g.KeywordOptions()).
	SQL("ORDER BY").
	OptionalText("OrderByClause", g.KeywordOptions()).
	OptionalText("WindowFrameClause", g.KeywordOptions())

var windowFunctionMetricDefinition = g.NewQueryStruct("WindowFunctionMetricDefinition").
	Text("WindowFunction", g.KeywordOptions().Required()).
	SQL("AS").
	Text("Metric", g.KeywordOptions().Required()).
	OptionalQueryStructField("OverClause", windowFunctionOverClause, g.ListOptions().Parentheses().NoComma().SQL("OVER"))

var metricDefinition = g.NewQueryStruct("MetricDefinition").
	OptionalQueryStructField("semanticExpression", semanticExpression, g.KeywordOptions()).
	OptionalQueryStructField("windowFunctionMetricDefinition", windowFunctionMetricDefinition, g.KeywordOptions()).
	WithValidation(g.ExactlyOneValueSet, "semanticExpression", "windowFunctionMetricDefinition")
