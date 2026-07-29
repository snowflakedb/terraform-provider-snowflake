package sdk

import (
	"context"
	"errors"
)

var (
	_ validatable = new(SetCommentOptions)
	_ validatable = new(SetColumnCommentOptions)
)

type Comments interface {
	Set(ctx context.Context, opts *SetCommentOptions) error
	SetColumn(ctx context.Context, opts *SetColumnCommentOptions) error
}

type comments struct {
	client *Client
}

var _ Comments = (*comments)(nil)

// SetCommentOptions is based on https://docs.snowflake.com/en/sql-reference/sql/comment.
type SetCommentOptions struct {
	comment    bool             `ddl:"static" sql:"COMMENT"`
	IfExists   *bool            `ddl:"keyword" sql:"IF EXISTS"`
	on         bool             `ddl:"static" sql:"ON"`
	ObjectType ObjectType       `ddl:"keyword"`
	ObjectName ObjectIdentifier `ddl:"identifier"`
	Value      *string          `ddl:"parameter,single_quotes,no_equals" sql:"IS"`
}

func (opts *SetCommentOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (c *comments) Set(ctx context.Context, opts *SetCommentOptions) error {
	if opts == nil {
		opts = &SetCommentOptions{}
	}
	// opts.name = name
	if err := opts.validate(); err != nil {
		return err
	}
	stmt, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, stmt)
	return err
}

// SetColumnCommentOptions is based on https://docs.snowflake.com/en/sql-reference/sql/comment.
type SetColumnCommentOptions struct {
	comment  bool             `ddl:"static" sql:"COMMENT"`
	IfExists *bool            `ddl:"keyword" sql:"IF EXISTS"`
	on       bool             `ddl:"static" sql:"ON"`
	Column   ObjectIdentifier `ddl:"identifier" sql:"COLUMN"`
	Value    *string          `ddl:"parameter,single_quotes,no_equals" sql:"IS"`
}

func (opts *SetColumnCommentOptions) validate() error {
	if opts == nil {
		return errors.Join(ErrNilOptions)
	}
	return nil
}

func (c *comments) SetColumn(ctx context.Context, opts *SetColumnCommentOptions) error {
	if opts == nil {
		opts = &SetColumnCommentOptions{}
	}
	// We only want to render table.column, not the fully qualified name with database and schema.
	if v, ok := opts.Column.(TableColumnIdentifier); ok {
		opts.Column = NewDatabaseObjectIdentifier(v.tableName, v.columnName)
	}
	if err := opts.validate(); err != nil {
		return err
	}
	stmt, err := structToSQL(opts)
	if err != nil {
		return err
	}
	_, err = c.client.exec(ctx, stmt)
	return err
}
