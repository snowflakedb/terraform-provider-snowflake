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
}

// CreateSemanticViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-semantic-view.
type CreateSemanticViewOptions struct {
	create        bool                   `ddl:"static" sql:"CREATE"`
	OrReplace     *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	semanticView  bool                   `ddl:"static" sql:"SEMANTIC VIEW"`
	IfNotExists   *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
	tables        bool                   `ddl:"static" sql:"TABLES"`
	logicalTables []LogicalTable         `ddl:"list,parentheses"`
	Comment       *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	CopyGrants    *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

type LogicalTable struct {
	logicalTableAlias *LogicalTableAlias     `ddl:"keyword"`
	TableName         SchemaObjectIdentifier `ddl:"identifier"`
	primaryKeys       *PrimaryKeys           `ddl:"parameter,no_equals"`
	uniqueKeys        []UniqueKeys           `ddl:"list,no_equals"`
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

type SemanticViewColumn struct {
	Name string `ddl:"keyword"`
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
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	DatabaseName  string         `db:"database_name"`
	SchemaName    string         `db:"schema_name"`
	Owner         string         `db:"owner"`
	OwnerRoleType string         `db:"owner_role_type"`
	Comment       sql.NullString `db:"comment"`
}

type SemanticViewDetails struct {
	CreatedOn     time.Time
	Name          string
	DatabaseName  string
	SchemaName    string
	Owner         string
	OwnerRoleType string
	Comment       *string
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
	Owner         string         `db:"owner"`
	OwnerRoleType string         `db:"owner_role_type"`
	Comment       sql.NullString `db:"comment"`
}

type SemanticView struct {
	CreatedOn     time.Time
	Name          string
	DatabaseName  string
	SchemaName    string
	Owner         string
	OwnerRoleType string
	Comment       *string
}

func (v *SemanticView) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}
