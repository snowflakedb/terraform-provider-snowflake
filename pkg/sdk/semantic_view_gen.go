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
	Describe(ctx context.Context, id SchemaObjectIdentifier) ([]semanticView, error)
	Show(ctx context.Context, request *ShowSemanticViewsRequest) ([]semanticView, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*semanticView, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*semanticView, error)
}

// CreateSemanticViewOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-semantic-view.
type CreateSemanticViewOptions struct {
	create       bool                   `ddl:"static" sql:"CREATE"`
	OrReplace    *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	semanticView bool                   `ddl:"static" sql:"SEMANTIC VIEW"`
	IfNotExists  *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
	Tables       []LogicalTable         `ddl:"parameter,parentheses" sql:"TABLES"`
	Comment      *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	CopyGrants   *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

type LogicalTable struct {
	logicalTableAlias *LogicalTableAlias     `ddl:"identifier"`
	logicalTableName  SchemaObjectIdentifier `ddl:"identifier"`
}

type LogicalTableAlias struct {
	logicalTableAlias StringProperty `ddl:"identifier"`
	As                string         `ddl:"parameter,single_quotes" sql:"AS"`
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

type semanticViewsRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	DatabaseName  string         `db:"database_name"`
	SchemaName    string         `db:"schema_name"`
	Owner         string         `db:"owner"`
	OwnerRoleType string         `db:"owner_role_type"`
	Comment       sql.NullString `db:"comment"`
}

type semanticView struct {
	CreatedOn     time.Time
	Name          string
	DatabaseName  string
	SchemaName    string
	Owner         string
	OwnerRoleType string
	Comment       *string
}

// ShowSemanticViewsOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-semantic-views.
type ShowSemanticViewsOptions struct {
	show          bool       `ddl:"static" sql:"SHOW"`
	Terse         *bool      `ddl:"keyword" sql:"TERSE"`
	semanticViews bool       `ddl:"static" sql:"SEMANTIC VIEWS"`
	Like          *Like      `ddl:"keyword" sql:"LIKE"`
	In            *In        `ddl:"keyword" sql:"IN"`
	StartsWith    *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit         *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}
