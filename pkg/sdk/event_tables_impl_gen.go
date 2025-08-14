package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ EventTables = (*eventTables)(nil)
var _ convertibleRow[EventTable] = new(eventTableRow)
var _ convertibleRow[EventTableDetails] = new(eventTableDetailsRow)

type eventTables struct {
	client *Client
}

func (v *eventTables) Create(ctx context.Context, request *CreateEventTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *eventTables) Show(ctx context.Context, request *ShowEventTableRequest) ([]EventTable, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[eventTableRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[eventTableRow, EventTable](dbRows)
}

func (v *eventTables) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*EventTable, error) {
	request := NewShowEventTableRequest().
		WithIn(In{Schema: id.SchemaId()}).
		WithLike(Like{Pattern: String(id.Name())})
	eventTables, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(eventTables, func(r EventTable) bool { return r.Name == id.Name() })
}

func (v *eventTables) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*EventTable, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *eventTables) Describe(ctx context.Context, id SchemaObjectIdentifier) (*EventTableDetails, error) {
	opts := &DescribeEventTableOptions{
		name: id,
	}
	result, err := validateAndQueryOne[eventTableDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convertErr()
}

func (v *eventTables) Drop(ctx context.Context, request *DropEventTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *eventTables) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropEventTableRequest(id).WithIfExists(Bool(true))) }, ctx, id)
}

func (v *eventTables) Alter(ctx context.Context, request *AlterEventTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (r *CreateEventTableRequest) toOpts() *CreateEventTableOptions {
	opts := &CreateEventTableOptions{
		OrReplace:                  r.OrReplace,
		IfNotExists:                r.IfNotExists,
		name:                       r.name,
		ClusterBy:                  r.ClusterBy,
		DataRetentionTimeInDays:    r.DataRetentionTimeInDays,
		MaxDataExtensionTimeInDays: r.MaxDataExtensionTimeInDays,
		ChangeTracking:             r.ChangeTracking,
		DefaultDdlCollation:        r.DefaultDdlCollation,
		CopyGrants:                 r.CopyGrants,
		Comment:                    r.Comment,
		RowAccessPolicy:            r.RowAccessPolicy,
		Tag:                        r.Tag,
	}
	return opts
}

func (r *ShowEventTableRequest) toOpts() *ShowEventTableOptions {
	opts := &ShowEventTableOptions{
		Terse:      r.Terse,
		Like:       r.Like,
		In:         r.In,
		StartsWith: r.StartsWith,
		Limit:      r.Limit,
	}
	return opts
}

func (r eventTableRow) convertErr() (*EventTable, error) {
	t := &EventTable{
		CreatedOn:    r.CreatedOn,
		Name:         r.Name,
		DatabaseName: r.DatabaseName,
		SchemaName:   r.SchemaName,
	}
	if r.Owner.Valid {
		t.Owner = r.Owner.String
	}
	if r.Comment.Valid {
		t.Comment = r.Comment.String
	}
	if r.OwnerRoleType.Valid {
		t.OwnerRoleType = r.OwnerRoleType.String
	}
	return t, nil
}

func (r *DescribeEventTableRequest) toOpts() *DescribeEventTableOptions {
	opts := &DescribeEventTableOptions{
		name: r.name,
	}
	return opts
}

func (r eventTableDetailsRow) convertErr() (*EventTableDetails, error) {
	return &EventTableDetails{
		Name:    r.Name,
		Kind:    r.Kind,
		Comment: r.Comment,
	}, nil
}

func (r *DropEventTableRequest) toOpts() *DropEventTableOptions {
	opts := &DropEventTableOptions{
		IfExists: r.IfExists,
		name:     r.name,
		Restrict: r.Restrict,
	}
	return opts
}

func (r *AlterEventTableRequest) toOpts() *AlterEventTableOptions {
	opts := &AlterEventTableOptions{
		IfNotExists: r.IfNotExists,
		name:        r.name,

		DropAllRowAccessPolicies: r.DropAllRowAccessPolicies,

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
		RenameTo:  r.RenameTo,
	}
	if r.Set != nil {
		opts.Set = &EventTableSet{
			DataRetentionTimeInDays:    r.Set.DataRetentionTimeInDays,
			MaxDataExtensionTimeInDays: r.Set.MaxDataExtensionTimeInDays,
			ChangeTracking:             r.Set.ChangeTracking,
			Comment:                    r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &EventTableUnset{
			DataRetentionTimeInDays:    r.Unset.DataRetentionTimeInDays,
			MaxDataExtensionTimeInDays: r.Unset.MaxDataExtensionTimeInDays,
			ChangeTracking:             r.Unset.ChangeTracking,
			Comment:                    r.Unset.Comment,
		}
	}
	if r.AddRowAccessPolicy != nil {
		opts.AddRowAccessPolicy = &EventTableAddRowAccessPolicy{
			RowAccessPolicy: r.AddRowAccessPolicy.RowAccessPolicy,
			On:              r.AddRowAccessPolicy.On,
		}
	}
	if r.DropRowAccessPolicy != nil {
		opts.DropRowAccessPolicy = &EventTableDropRowAccessPolicy{
			RowAccessPolicy: r.DropRowAccessPolicy.RowAccessPolicy,
		}
	}
	if r.DropAndAddRowAccessPolicy != nil {
		opts.DropAndAddRowAccessPolicy = &EventTableDropAndAddRowAccessPolicy{}
		opts.DropAndAddRowAccessPolicy.Drop = EventTableDropRowAccessPolicy{
			RowAccessPolicy: r.DropAndAddRowAccessPolicy.Drop.RowAccessPolicy,
		}
		opts.DropAndAddRowAccessPolicy.Add = EventTableAddRowAccessPolicy{
			RowAccessPolicy: r.DropAndAddRowAccessPolicy.Add.RowAccessPolicy,
			On:              r.DropAndAddRowAccessPolicy.Add.On,
		}
	}
	if r.ClusteringAction != nil {
		opts.ClusteringAction = &EventTableClusteringAction{
			ClusterBy:         r.ClusteringAction.ClusterBy,
			SuspendRecluster:  r.ClusteringAction.SuspendRecluster,
			ResumeRecluster:   r.ClusteringAction.ResumeRecluster,
			DropClusteringKey: r.ClusteringAction.DropClusteringKey,
		}
	}
	if r.SearchOptimizationAction != nil {
		opts.SearchOptimizationAction = &EventTableSearchOptimizationAction{}
		if r.SearchOptimizationAction.Add != nil {
			opts.SearchOptimizationAction.Add = &SearchOptimization{
				On: r.SearchOptimizationAction.Add.On,
			}
		}
		if r.SearchOptimizationAction.Drop != nil {
			opts.SearchOptimizationAction.Drop = &SearchOptimization{
				On: r.SearchOptimizationAction.Drop.On,
			}
		}
	}
	return opts
}
