package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var (
	_ Views                       = (*views)(nil)
	_ convertibleRow[ViewDetails] = new(viewDetailsRow)
	_ convertibleRow[View]        = new(viewDBRow)
)

type views struct {
	client *Client
}

func (v *views) Create(ctx context.Context, request *CreateViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *views) Alter(ctx context.Context, request *AlterViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *views) Drop(ctx context.Context, request *DropViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *views) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropViewRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *views) Show(ctx context.Context, request *ShowViewRequest) ([]View, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[viewDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[viewDBRow, View](dbRows)
}

func (v *views) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*View, error) {
	request := NewShowViewRequest().
		WithIn(ExtendedIn{In: In{Schema: id.SchemaId()}}).
		WithLike(Like{Pattern: String(id.Name())})
	views, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(views, func(r View) bool { return r.Name == id.Name() })
}

func (v *views) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*View, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *views) Describe(ctx context.Context, id SchemaObjectIdentifier) ([]ViewDetails, error) {
	opts := &DescribeViewOptions{
		name: id,
	}
	rows, err := validateAndQuery[viewDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[viewDetailsRow, ViewDetails](rows)
}

func (r *CreateViewRequest) toOpts() *CreateViewOptions {
	opts := &CreateViewOptions{
		OrReplace:   r.OrReplace,
		Secure:      r.Secure,
		Temporary:   r.Temporary,
		Recursive:   r.Recursive,
		IfNotExists: r.IfNotExists,
		name:        r.name,

		CopyGrants: r.CopyGrants,
		Comment:    r.Comment,

		Tag: r.Tag,
		sql: r.sql,
	}

	if r.Columns != nil {
		s := make([]ViewColumn, len(r.Columns))
		for i, v := range r.Columns {
			s[i] = ViewColumn{
				Name:    v.Name,
				Tag:     v.Tag,
				Comment: v.Comment,
			}
			if v.ProjectionPolicy != nil {
				s[i].ProjectionPolicy = &ViewColumnProjectionPolicy{
					ProjectionPolicy: v.ProjectionPolicy.ProjectionPolicy,
				}
			}
			if v.MaskingPolicy != nil {
				s[i].MaskingPolicy = &ViewColumnMaskingPolicy{
					MaskingPolicy: v.MaskingPolicy.MaskingPolicy,
					Using:         v.MaskingPolicy.Using,
				}
			}
		}
		opts.Columns = s
	}

	if r.RowAccessPolicy != nil {
		opts.RowAccessPolicy = &ViewRowAccessPolicy{
			RowAccessPolicy: r.RowAccessPolicy.RowAccessPolicy,
			On:              r.RowAccessPolicy.On,
		}
	}

	if r.AggregationPolicy != nil {
		opts.AggregationPolicy = &ViewAggregationPolicy{
			AggregationPolicy: r.AggregationPolicy.AggregationPolicy,
			EntityKey:         r.AggregationPolicy.EntityKey,
		}
	}

	return opts
}

func (r *AlterViewRequest) toOpts() *AlterViewOptions {
	opts := &AlterViewOptions{
		IfExists:          r.IfExists,
		name:              r.name,
		RenameTo:          r.RenameTo,
		SetComment:        r.SetComment,
		UnsetComment:      r.UnsetComment,
		SetSecure:         r.SetSecure,
		SetChangeTracking: r.SetChangeTracking,
		UnsetSecure:       r.UnsetSecure,
		SetTags:           r.SetTags,
		UnsetTags:         r.UnsetTags,

		DropAllRowAccessPolicies: r.DropAllRowAccessPolicies,
	}

	if r.AddDataMetricFunction != nil {
		opts.AddDataMetricFunction = &ViewAddDataMetricFunction{
			DataMetricFunction: r.AddDataMetricFunction.DataMetricFunction,
		}
	}

	if r.DropDataMetricFunction != nil {
		opts.DropDataMetricFunction = &ViewDropDataMetricFunction{
			DataMetricFunction: r.DropDataMetricFunction.DataMetricFunction,
		}
	}

	if r.ModifyDataMetricFunction != nil {
		opts.ModifyDataMetricFunction = &ViewModifyDataMetricFunctions{
			DataMetricFunction: r.ModifyDataMetricFunction.DataMetricFunction,
		}
	}

	if r.SetDataMetricSchedule != nil {
		opts.SetDataMetricSchedule = &ViewSetDataMetricSchedule{
			DataMetricSchedule: r.SetDataMetricSchedule.DataMetricSchedule,
		}
	}

	if r.UnsetDataMetricSchedule != nil {
		opts.UnsetDataMetricSchedule = &ViewUnsetDataMetricSchedule{}
	}

	if r.AddRowAccessPolicy != nil {
		opts.AddRowAccessPolicy = &ViewAddRowAccessPolicy{
			RowAccessPolicy: r.AddRowAccessPolicy.RowAccessPolicy,
			On:              r.AddRowAccessPolicy.On,
		}
	}

	if r.DropRowAccessPolicy != nil {
		opts.DropRowAccessPolicy = &ViewDropRowAccessPolicy{
			RowAccessPolicy: r.DropRowAccessPolicy.RowAccessPolicy,
		}
	}

	if r.DropAndAddRowAccessPolicy != nil {
		opts.DropAndAddRowAccessPolicy = &ViewDropAndAddRowAccessPolicy{}

		opts.DropAndAddRowAccessPolicy.Drop = ViewDropRowAccessPolicy{
			RowAccessPolicy: r.DropAndAddRowAccessPolicy.Drop.RowAccessPolicy,
		}

		opts.DropAndAddRowAccessPolicy.Add = ViewAddRowAccessPolicy{
			RowAccessPolicy: r.DropAndAddRowAccessPolicy.Add.RowAccessPolicy,
			On:              r.DropAndAddRowAccessPolicy.Add.On,
		}
	}

	if r.SetAggregationPolicy != nil {
		opts.SetAggregationPolicy = &ViewSetAggregationPolicy{
			AggregationPolicy: r.SetAggregationPolicy.AggregationPolicy,
			EntityKey:         r.SetAggregationPolicy.EntityKey,
			Force:             r.SetAggregationPolicy.Force,
		}
	}

	if r.UnsetAggregationPolicy != nil {
		opts.UnsetAggregationPolicy = &ViewUnsetAggregationPolicy{}
	}

	if r.SetMaskingPolicyOnColumn != nil {
		opts.SetMaskingPolicyOnColumn = &ViewSetColumnMaskingPolicy{
			Name:          r.SetMaskingPolicyOnColumn.Name,
			MaskingPolicy: r.SetMaskingPolicyOnColumn.MaskingPolicy,
			Using:         r.SetMaskingPolicyOnColumn.Using,
			Force:         r.SetMaskingPolicyOnColumn.Force,
		}
	}

	if r.UnsetMaskingPolicyOnColumn != nil {
		opts.UnsetMaskingPolicyOnColumn = &ViewUnsetColumnMaskingPolicy{
			Name: r.UnsetMaskingPolicyOnColumn.Name,
		}
	}

	if r.SetProjectionPolicyOnColumn != nil {
		opts.SetProjectionPolicyOnColumn = &ViewSetProjectionPolicy{
			Name:             r.SetProjectionPolicyOnColumn.Name,
			ProjectionPolicy: r.SetProjectionPolicyOnColumn.ProjectionPolicy,
			Force:            r.SetProjectionPolicyOnColumn.Force,
		}
	}

	if r.UnsetProjectionPolicyOnColumn != nil {
		opts.UnsetProjectionPolicyOnColumn = &ViewUnsetProjectionPolicy{
			Name: r.UnsetProjectionPolicyOnColumn.Name,
		}
	}

	if r.SetTagsOnColumn != nil {
		opts.SetTagsOnColumn = &ViewSetColumnTags{
			Name:    r.SetTagsOnColumn.Name,
			SetTags: r.SetTagsOnColumn.SetTags,
		}
	}

	if r.UnsetTagsOnColumn != nil {
		opts.UnsetTagsOnColumn = &ViewUnsetColumnTags{
			Name:      r.UnsetTagsOnColumn.Name,
			UnsetTags: r.UnsetTagsOnColumn.UnsetTags,
		}
	}

	return opts
}

func (r *DropViewRequest) toOpts() *DropViewOptions {
	opts := &DropViewOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowViewRequest) toOpts() *ShowViewOptions {
	opts := &ShowViewOptions{
		Terse:      r.Terse,
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r viewDBRow) convert() (*View, error) {
	view := View{
		CreatedOn:    r.CreatedOn,
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
	}
	if r.Reserved.Valid {
		view.Reserved = r.Reserved.String
	}
	if r.Owner.Valid {
		view.Owner = r.Owner.String
	}
	if r.Comment.Valid {
		view.Comment = r.Comment.String
	}
	if r.Text.Valid {
		view.Text = tracking.TrimMetadata(r.Text.String)
	}
	if r.Kind.Valid {
		view.Kind = r.Kind.String
	}
	if r.IsSecure.Valid {
		view.IsSecure = r.IsSecure.Bool
	}
	if r.IsMaterialized.Valid {
		view.IsMaterialized = r.IsMaterialized.Bool
	}
	if r.OwnerRoleType.Valid {
		view.OwnerRoleType = r.OwnerRoleType.String
	}
	if r.ChangeTracking.Valid {
		view.ChangeTracking = r.ChangeTracking.String
	}
	return &view, nil
}

func (r *DescribeViewRequest) toOpts() *DescribeViewOptions {
	opts := &DescribeViewOptions{
		name: r.name,
	}
	return opts
}

func (r viewDetailsRow) convert() (*ViewDetails, error) {
	details := &ViewDetails{
		Name:       r.Name,
		Type:       r.Type,
		Kind:       r.Kind,
		IsNullable: r.Null == "Y",
		IsPrimary:  r.PrimaryKey == "Y",
		IsUnique:   r.UniqueKey == "Y",
	}
	if r.Default.Valid {
		details.Default = String(r.Default.String)
	}
	if r.Check.Valid {
		details.Check = Bool(r.Check.String == "Y")
	}
	if r.Expression.Valid {
		details.Expression = String(r.Expression.String)
	}
	if r.Comment.Valid {
		details.Comment = String(r.Comment.String)
	}
	if r.PolicyName.Valid {
		details.PolicyName = String(r.PolicyName.String)
	}
	return details, nil
}
