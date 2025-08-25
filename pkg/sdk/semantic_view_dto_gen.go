package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSemanticViewOptions]   = new(CreateSemanticViewRequest)
	_ optionsProvider[AlterSemanticViewOptions]    = new(AlterSemanticViewRequest)
	_ optionsProvider[DropSemanticViewOptions]     = new(DropSemanticViewRequest)
	_ optionsProvider[DescribeSemanticViewOptions] = new(DescribeSemanticViewRequest)
	_ optionsProvider[ShowSemanticViewOptions]     = new(ShowSemanticViewRequest)
)

type CreateSemanticViewRequest struct {
	OrReplace   *bool
	IfNotExists *bool
	name        SchemaObjectIdentifier // required
	// Adjusted manually (changed from opt structs to Requests)
	logicalTables             []LogicalTableRequest // required
	semanticViewRelationships []SemanticViewRelationshipRequest
	semanticViewFacts         []SemanticExpressionRequest
	semanticViewDimensions    []SemanticExpressionRequest
	semanticViewMetrics       []MetricDefinitionRequest
	Comment                   *string
	CopyGrants                *bool
}

type AlterSemanticViewRequest struct {
	IfExists     *bool
	name         SchemaObjectIdentifier // required
	SetComment   *string
	UnsetComment *bool
}

type DropSemanticViewRequest struct {
	IfExists *bool
	name     SchemaObjectIdentifier // required
}

type DescribeSemanticViewRequest struct {
	name SchemaObjectIdentifier // required
}

type ShowSemanticViewRequest struct {
	Terse      *bool
	Like       *Like
	In         *In
	StartsWith *string
	Limit      *LimitFrom
}

// The requests below added manually

type LogicalTableRequest struct {
	logicalTableAlias *LogicalTableAliasRequest
	TableName         SchemaObjectIdentifier // required
	primaryKeys       *PrimaryKeysRequest
	uniqueKeys        []UniqueKeysRequest
	synonyms          *SynonymsRequest
	Comment           *string
}

type LogicalTableAliasRequest struct {
	LogicalTableAlias string
}

type PrimaryKeysRequest struct {
	PrimaryKey []SemanticViewColumn
}

type UniqueKeysRequest struct {
	Unique []SemanticViewColumn
}

type SynonymsRequest struct {
	WithSynonyms []Synonym
}

type SemanticViewRelationshipRequest struct {
	relationshipAlias          *RelationshipAliasRequest
	tableNameOrAlias           *RelationshipTableAliasRequest // required
	relationshipColumnNames    []SemanticViewColumnRequest    // required
	refTableNameOrAlias        *RelationshipTableAliasRequest // required
	relationshipRefColumnNames []SemanticViewColumnRequest
}

type RelationshipAliasRequest struct {
	RelationshipAlias string
}

type RelationshipTableAliasRequest struct {
	RelationshipTableName  *SchemaObjectIdentifier
	RelationshipTableAlias *string
}

type SemanticViewColumnRequest struct {
	Name string // required
}

type SemanticExpressionRequest struct {
	qualifiedExpressionName *QualifiedExpressionNameRequest // required
	sqlExpression           *SemanticSqlExpressionRequest   // required
	synonyms                *SynonymsRequest
	Comment                 *string
}

type QualifiedExpressionNameRequest struct {
	QualifiedExpressionName string
}

type SemanticSqlExpressionRequest struct {
	SqlExpression string
}

type MetricDefinitionRequest struct {
	semanticExpression             *SemanticExpressionRequest
	windowFunctionMetricDefinition *WindowFunctionMetricDefinitionRequest
}

type WindowFunctionMetricDefinitionRequest struct {
	WindowFunction string // required
	Metric         string // required
	OverClause     *WindowFunctionOverClauseRequest
}

type WindowFunctionOverClauseRequest struct {
	PartitionBy       *bool
	PartitionByClause *string
	OrderBy           *bool
	OrderByClause     *string
	WindowFrameClause *string
}
