package sdk

import (
	"context"
)

type SemanticViews interface {
	Create(ctx context.Context, request *CreateSemanticViewRequest) error
	Drop(ctx context.Context, request *DropSemanticViewRequest) error
}

type CreateSemanticViewOptions struct {
	create       bool                   `ddl:"static" sql:"CREATE"`
	OrReplace    *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	semanticView bool                   `ddl:"static" sql:"SEMANTIC VIEW"`
	IfNotExists  *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
	Comment      *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
	CopyGrants   *bool                  `ddl:"keyword" sql:"COPY GRANTS"`
}

type DropSemanticViewOptions struct {
	drop         bool                   `ddl:"static" sql:"DROP"`
	semanticView bool                   `ddl:"static" sql:"SEMANTIC VIEW"`
	IfExists     *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name         SchemaObjectIdentifier `ddl:"identifier"`
}
