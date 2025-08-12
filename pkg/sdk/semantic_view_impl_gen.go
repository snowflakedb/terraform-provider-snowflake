package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
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

func (v *semanticViews) Describe(ctx context.Context, id SchemaObjectIdentifier) ([]semanticView, error) {
	opts := &DescribeSemanticViewOptions{
		name: id,
	}
	rows, err := validateAndQuery[semanticViewsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[semanticViewsRow, semanticView](rows), nil
}

func (v *semanticViews) Show(ctx context.Context, request *ShowSemanticViewRequest) ([]semanticView, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[semanticViewsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[semanticViewsRow, semanticView](dbRows)
	return resultList, nil
}

func (v *semanticViews) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*semanticView, error) {
	request := NewShowSemanticViewRequest().
		WithLike(Like{Pattern: String(id.Name())}).
		WithIn(In{Schema: id.SchemaId()})
	semanticViews, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(semanticViews, func(r semanticView) bool { return r.Name == id.Name() })
}

func (v *semanticViews) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*semanticView, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *CreateSemanticViewRequest) toOpts() *CreateSemanticViewOptions {
	opts := &CreateSemanticViewOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Tables:      r.Tables,
		Comment:     r.Comment,
		CopyGrants:  r.CopyGrants,
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

func (r semanticViewsRow) convert() *semanticView {
	// TODO: Mapping
	return &semanticView{}
}

func (r *ShowSemanticViewRequest) toOpts() *ShowSemanticViewsOptions {
	opts := &ShowSemanticViewsOptions{
		Terse:      r.Terse,
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}
