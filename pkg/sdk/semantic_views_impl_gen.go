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

func (v *semanticViews) Describe(ctx context.Context, request *DescribeSemanticViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *semanticViews) Show(ctx context.Context, request *ShowSemanticViewsRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateSemanticViewRequest) toOpts() *CreateSemanticViewOptions {
	opts := &CreateSemanticViewOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Comment:     r.Comment,
		CopyGrants:  r.CopyGrants,
	}
	return opts
}

func (r *DropSemanticViewRequest) toOpts() *DropSemanticViewOptions {
	opts := &DropSemanticViewOptions{
		name:     r.name,
		IfExists: r.IfExists,
	}
	return opts
}

func (r *DescribeSemanticViewRequest) toOpts() *DescribeSemanticViewOptions {
	opts := &DescribeSemanticViewOptions{
		name: r.name,
	}
	return opts
}

func (r *ShowSemanticViewsRequest) toOpts() *ShowSemanticViewsOptions {
	opts := &ShowSemanticViewsOptions{
		Like:       r.Like,
		In:         r.In,
		Terse:      r.Terse,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}
