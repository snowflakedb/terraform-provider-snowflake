package sdk

import (
	"context"
	"database/sql"
	"time"
)

type SemanticViews interface {
	Create(ctx context.Context, request *CreateSemanticViewRequest) error
	Drop(ctx context.Context, request *DropSemanticViewRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Describe(ctx context.Context, id SchemaObjectIdentifier) ([]SemanticViewDetails, error)
	Show(ctx context.Context, request *ShowSemanticViewRequest) ([]SemanticView, error)
	Alter(ctx context.Context, request *AlterSemanticViewRequest) error
}

// CreateSemanticViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-semantic-view.
type CreateSemanticViewOptions struct {
	create                    bool                       `ddl:"static" sql:"CREATE"`
	OrReplace                 *bool                      `ddl:"keyword" sql:"OR REPLACE"`
	semanticView              bool                       `ddl:"static" sql:"SEMANTIC VIEW"`
	IfNotExists               *bool                      `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                      SchemaObjectIdentifier     `ddl:"identifier"`
	tables                    bool                       `ddl:"static" sql:"TABLES"`
	logicalTables             []LogicalTable             `ddl:"list,parentheses"`
	Relationships             *bool                      `ddl:"keyword" sql:"RELATIONSHIPS"`
	semanticViewRelationships []SemanticViewRelationship `ddl:"list,parentheses"`
	Facts                     *bool                      `ddl:"keyword" sql:"FACTS"`
	semanticViewFacts         []SemanticExpression       `ddl:"list,parentheses"`
	Dimensions                *bool                      `ddl:"keyword" sql:"DIMENSIONS"`
	semanticViewDimensions    []SemanticExpression       `ddl:"list,parentheses"`
	Metrics                   *bool                      `ddl:"keyword" sql:"METRICS"`
	semanticViewMetrics       []MetricDefinition         `ddl:"list,parentheses"`
	Comment                   *string                    `ddl:"parameter,single_quotes" sql:"COMMENT"`
	CopyGrants                *bool                      `ddl:"keyword" sql:"COPY GRANTS"`
}

type LogicalTable struct {
	logicalTableAlias *LogicalTableAlias     `ddl:"keyword"`
	TableName         SchemaObjectIdentifier `ddl:"identifier"`
	primaryKeys       *PrimaryKeys           `ddl:"parameter,no_equals"`
	uniqueKeys        []UniqueKeys           `ddl:"list,no_equals,no_comma"`
	synonyms          *Synonyms              `ddl:"parameter,no_equals"`
	Comment           *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type LogicalTableAlias struct {
	LogicalTableAlias string `ddl:"keyword"`
	as                bool   `ddl:"static" sql:"AS"`
}

type PrimaryKeys struct {
	PrimaryKey []SemanticViewColumn `ddl:"parameter,parentheses,no_equals" sql:"PRIMARY KEY"`
}

type UniqueKeys struct {
	Unique []SemanticViewColumn `ddl:"parameter,parentheses,no_equals" sql:"UNIQUE"`
}

type Synonyms struct {
	WithSynonyms []string `ddl:"parameter,parentheses,no_equals" sql:"WITH SYNONYMS"`
}

type SemanticViewRelationship struct {
	relationshipAlias          *RelationshipAlias      `ddl:"keyword"`
	tableName                  *RelationshipTableAlias `ddl:"keyword"`
	relationshipColumnNames    []SemanticViewColumn    `ddl:"list,parentheses,no_equals"`
	references                 bool                    `ddl:"static" sql:"REFERENCES"`
	refTableName               *RelationshipTableAlias `ddl:"keyword"`
	relationshipRefColumnNames []SemanticViewColumn    `ddl:"list,parentheses,no_equals"`
}

type RelationshipAlias struct {
	RelationshipAlias string `ddl:"keyword"`
	as                bool   `ddl:"static" sql:"AS"`
}

type RelationshipTableAlias struct {
	RelationshipTableAlias string `ddl:"keyword"`
}

type SemanticViewColumn struct {
	Name string `ddl:"keyword"`
}

type SemanticExpression struct {
	qualifiedExpressionName *QualifiedExpressionName `ddl:"keyword"`
	as                      bool                     `ddl:"static" sql:"AS"`
	sqlExpression           *SemanticSqlExpression   `ddl:"keyword"`
	synonyms                *Synonyms                `ddl:"parameter,no_equals"`
	Comment                 *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type QualifiedExpressionName struct {
	QualifiedExpressionName string `ddl:"keyword"`
}

type SemanticSqlExpression struct {
	SqlExpression string `ddl:"keyword,no_quotes"`
}

type MetricDefinition struct {
	semanticExpression             *SemanticExpression             `ddl:"keyword"`
	windowFunctionMetricDefinition *WindowFunctionMetricDefinition `ddl:"keyword"`
}

type WindowFunctionMetricDefinition struct {
	WindowFunction string                    `ddl:"keyword"`
	as             bool                      `ddl:"static" sql:"AS"`
	Metric         string                    `ddl:"keyword"`
	OverClause     *WindowFunctionOverClause `ddl:"list,parentheses,no_comma" sql:"OVER"`
}

type WindowFunctionOverClause struct {
	PartitionBy       *bool   `ddl:"keyword" sql:"PARTITION BY"`
	PartitionByClause *string `ddl:"keyword"`
	OrderBy           *bool   `ddl:"keyword" sql:"ORDER BY"`
	OrderByClause     *string `ddl:"keyword"`
	WindowFrameClause *string `ddl:"keyword"`
}

// DropSemanticViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-semantic-view.
type DropSemanticViewOptions struct {
	drop         bool                   `ddl:"static" sql:"DROP"`
	semanticView bool                   `ddl:"static" sql:"SEMANTIC VIEW"`
	IfExists     *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
}

// DescribeSemanticViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-semantic-view.
type DescribeSemanticViewOptions struct {
	describe     bool                   `ddl:"static" sql:"DESCRIBE"`
	semanticView bool                   `ddl:"static" sql:"SEMANTIC VIEW"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
}

type semanticViewDetailsRow struct {
	ObjectKind    string `db:"object_kind"`
	ObjectName    string `db:"object_name"`
	ParentEntity  string `db:"parent_entity"`
	Property      string `db:"property"`
	PropertyValue string `db:"property_value"`
}

type SemanticViewDetails struct {
	ObjectKind    string
	ObjectName    string
	ParentEntity  string
	Property      string
	PropertyValue string
}

// ShowSemanticViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views.
type ShowSemanticViewOptions struct {
	show          bool       `ddl:"static" sql:"SHOW"`
	Terse         *bool      `ddl:"keyword" sql:"TERSE"`
	semanticViews bool       `ddl:"static" sql:"SEMANTIC VIEWS"`
	Like          *Like      `ddl:"keyword" sql:"LIKE"`
	In            *In        `ddl:"keyword" sql:"IN"`
	StartsWith    *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit         *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

type semanticViewDBRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	DatabaseName  string         `db:"database_name"`
	SchemaName    string         `db:"schema_name"`
	Comment       sql.NullString `db:"comment"`
	Owner         string         `db:"owner"`
	OwnerRoleType string         `db:"owner_role_type"`
	Extension     sql.NullString `db:"extension"`
}

type SemanticView struct {
	CreatedOn     time.Time
	Name          string
	DatabaseName  string
	SchemaName    string
	Comment       *string
	Owner         string
	OwnerRoleType string
	Extension     *string
}

func (v *SemanticView) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}
func (v *SemanticView) ObjectType() ObjectType {
	return ObjectTypeSemanticView
}

// AlterSemanticViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-semantic-view.
type AlterSemanticViewOptions struct {
	alter        bool                   `ddl:"static" sql:"ALTER"`
	semanticView bool                   `ddl:"static" sql:"SEMANTIC VIEW"`
	IfExists     *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
	SetComment   *string                `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	UnsetComment *bool                  `ddl:"keyword" sql:"UNSET COMMENT"`
}
