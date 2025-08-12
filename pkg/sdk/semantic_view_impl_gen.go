package sdk

import (
	"context"
)

var _ SemanticViews = (*semanticViews)(nil)

type semanticViews struct {
	client *Client
}

func (v *semanticViews) Create(ctx context.Context, request *CreateSemanticViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *semanticViews) Drop(ctx context.Context, request *DropSemanticViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *semanticViews) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropSemanticViewRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *semanticViews) Describe(ctx context.Context, id SchemaObjectIdentifier) ([]SemanticViewDetails, error) {
	opts := &DescribeSemanticViewOptions{
		name: id,
	}
	rows, err := validateAndQuery[semanticViewDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[semanticViewDetailsRow, SemanticViewDetails](rows), nil
}

func (v *semanticViews) Show(ctx context.Context, request *ShowSemanticViewRequest) ([]SemanticView, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[semanticViewDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[semanticViewDBRow, SemanticView](dbRows)
	return resultList, nil
}

func (r *CreateSemanticViewRequest) toOpts() *CreateSemanticViewOptions {
	opts := &CreateSemanticViewOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		Comment:    r.Comment,
		CopyGrants: r.CopyGrants,
	}
	if r.tables != nil {
		s := make([]LogicalTable, len(r.tables))
		for i, v := range r.tables {
			s[i] = LogicalTable{
				logicalTableName: v.logicalTableName,
			}
		}
		opts.tables = s
	}
	return opts
}

func (r *DropSemanticViewRequest) toOpts() *DropSemanticViewOptions {
	opts := &DropSemanticViewOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeSemanticViewRequest) toOpts() *DescribeSemanticViewOptions {
	opts := &DescribeSemanticViewOptions{
		name: r.name,
	}
	return opts
}

func (r semanticViewDetailsRow) convert() *SemanticViewDetails {
	// TODO: Mapping
	return &SemanticViewDetails{}
}

func (r *ShowSemanticViewRequest) toOpts() *ShowSemanticViewOptions {
	opts := &ShowSemanticViewOptions{
		Terse:      r.Terse,
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r semanticViewDBRow) convert() *SemanticView {
	// TODO: Mapping
	return &SemanticView{}
}
