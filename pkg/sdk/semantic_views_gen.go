package sdk

import (
	"context"
)

type SemanticViews interface {
	Create(ctx context.Context, request *CreateSemanticViewRequest) error
	Drop(ctx context.Context, request *DropSemanticViewRequest) error
}

type CreateSemanticViewOptions struct {
	create       bool                     `ddl:"static" sql:"CREATE"`
	OrReplace    *bool                    `ddl:"keyword" sql:"OR REPLACE"`
	semanticView bool                     `ddl:"static" sql:"SEMANTIC VIEW"`
	IfNotExists  *bool                    `ddl:"keyword" sql:"IF NOT EXISTS"`
	name         SchemaObjectIdentifier   `ddl:"identifier"`
	tables       []LogicalTableIdentifier `ddl:"tables"`
	Comment      *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
	CopyGrants   *bool                    `ddl:"keyword" sql:"COPY GRANTS"`
}

type DropSemanticViewOptions struct {
	drop         bool                   `ddl:"static" sql:"DROP"`
	semanticView bool                   `ddl:"static" sql:"SEMANTIC VIEW"`
	IfExists     *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
}

type DescribeSemanticViewOptions struct {
	describe     bool                   `ddl:"static" sql:"DESCRIBE"`
	semanticView bool                   `ddl:"static" sql:"SEMANTIC VIEW"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
}

type ShowSemanticViewsOptions struct {
	show          bool       `ddl:"static" sql:"SHOW"`
	Terse         *bool      `ddl:"keyword" sql:"TERSE"`
	semanticViews bool       `ddl:"static" sql:"SEMANTIC VIEWS"`
	Like          *Like      `ddl:"keyword" sql:"LIKE"`
	In            *In        `ddl:"keyword" sql:"IN"`
	StartsWith    *string    `ddl:"parameter,no_equals,single_quotes" sql:"STARTS WITH"`
	Limit         *LimitFrom `ddl:"parameter,no_equals" sql:"LIMIT"`
}

type semanticViewsRow struct {
	ObjectKind    string `db:"object_kind"`
	ObjectName    string `db:"object_name"`
	ParentEntity  string `db:"parent_entity"`
	Property      string `db:"property"`
	PropertyValue string `db:"property_value"`
}

type SemanticView struct {
	ObjectKind    string
	ObjectName    string
	ParentEntity  string
	Property      string
	PropertyValue string
}
