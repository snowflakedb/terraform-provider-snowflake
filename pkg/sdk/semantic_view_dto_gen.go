package sdk

//go:generate go run ./dto-builder-generator/main.go

var (
	_ optionsProvider[CreateSemanticViewOptions]   = new(CreateSemanticViewRequest)
	_ optionsProvider[DropSemanticViewOptions]     = new(DropSemanticViewRequest)
	_ optionsProvider[DescribeSemanticViewOptions] = new(DescribeSemanticViewRequest)
	_ optionsProvider[ShowSemanticViewOptions]     = new(ShowSemanticViewRequest)
	_ optionsProvider[AlterSemanticViewOptions]    = new(AlterSemanticViewRequest)
)

type CreateSemanticViewRequest struct {
	OrReplace                 *bool
	IfNotExists               *bool
	name                      SchemaObjectIdentifier // required
	logicalTables             []LogicalTableRequest  // required
	Relationships             *bool
	semanticViewRelationships []SemanticViewRelationshipRequest
	Facts                     *bool
	semanticViewFacts         []SemanticExpressionRequest
	Dimensions                *bool
	semanticViewDimensions    []SemanticExpressionRequest
	Metrics                   *bool
	semanticViewMetrics       []MetricDefinitionRequest
	Comment                   *string
	CopyGrants                *bool
}

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
	WithSynonyms []string
}

type SemanticViewRelationshipRequest struct {
	relationshipAlias          *RelationshipAliasRequest
	tableName                  *RelationshipTableAliasRequest // required
	relationshipColumnNames    []SemanticViewColumnRequest    // required
	refTableName               *RelationshipTableAliasRequest // required
	relationshipRefColumnNames []SemanticViewColumnRequest
}

type RelationshipAliasRequest struct {
	RelationshipAlias string
}

type RelationshipTableAliasRequest struct {
	RelationshipTableAlias string
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
	PartitionByClause *string
	OrderByClause     *string
	WindowFrameClause *string
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

type AlterSemanticViewRequest struct {
	IfExists     *bool
	name         SchemaObjectIdentifier // required
	SetComment   *string
	UnsetComment *bool
}
