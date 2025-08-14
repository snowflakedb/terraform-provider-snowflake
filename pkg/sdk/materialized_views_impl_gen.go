package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/internal/tracking"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ MaterializedViews = (*materializedViews)(nil)

type materializedViews struct {
	client *Client
}

func (v *materializedViews) Create(ctx context.Context, request *CreateMaterializedViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *materializedViews) Alter(ctx context.Context, request *AlterMaterializedViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *materializedViews) Drop(ctx context.Context, request *DropMaterializedViewRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *materializedViews) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropMaterializedViewRequest(id).WithIfExists(Bool(true))) }, ctx, id)
}

func (v *materializedViews) Show(ctx context.Context, request *ShowMaterializedViewRequest) ([]MaterializedView, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[materializedViewDBRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[materializedViewDBRow, MaterializedView](dbRows)
}

func (v *materializedViews) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*MaterializedView, error) {
	request := NewShowMaterializedViewRequest().
		WithIn(In{Schema: id.SchemaId()}).
		WithLike(Like{Pattern: String(id.Name())})
	materializedViews, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(materializedViews, func(r MaterializedView) bool { return r.Name == id.Name() })
}

func (v *materializedViews) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*MaterializedView, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *materializedViews) Describe(ctx context.Context, id SchemaObjectIdentifier) ([]MaterializedViewDetails, error) {
	opts := &DescribeMaterializedViewOptions{
		name: id,
	}
	rows, err := validateAndQuery[materializedViewDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRowsErr[materializedViewDetailsRow, MaterializedViewDetails](rows)
}

func (r *CreateMaterializedViewRequest) toOpts() *CreateMaterializedViewOptions {
	opts := &CreateMaterializedViewOptions{
		OrReplace:   r.OrReplace,
		Secure:      r.Secure,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		CopyGrants:  r.CopyGrants,
		Comment:     r.Comment,
		Tag:         r.Tag,
		sql:         r.sql,
	}
	if r.Columns != nil {
		s := make([]MaterializedViewColumn, len(r.Columns))
		for i, v := range r.Columns {
			s[i] = MaterializedViewColumn(v)
		}
		opts.Columns = s
	}
	if r.ColumnsMaskingPolicies != nil {
		s := make([]MaterializedViewColumnMaskingPolicy, len(r.ColumnsMaskingPolicies))
		for i, v := range r.ColumnsMaskingPolicies {
			s[i] = MaterializedViewColumnMaskingPolicy(v)
		}
		opts.ColumnsMaskingPolicies = s
	}
	if r.RowAccessPolicy != nil {
		opts.RowAccessPolicy = &MaterializedViewRowAccessPolicy{
			RowAccessPolicy: r.RowAccessPolicy.RowAccessPolicy,
			On:              r.RowAccessPolicy.On,
		}
	}
	if r.ClusterBy != nil {
		opts.ClusterBy = &MaterializedViewClusterBy{}
		if r.ClusterBy.Expressions != nil {
			s := make([]MaterializedViewClusterByExpression, len(r.ClusterBy.Expressions))
			for i, v := range r.ClusterBy.Expressions {
				s[i] = MaterializedViewClusterByExpression(v)
			}
			opts.ClusterBy.Expressions = s
		}
	}
	return opts
}

func (r *AlterMaterializedViewRequest) toOpts() *AlterMaterializedViewOptions {
	opts := &AlterMaterializedViewOptions{
		name:              r.name,
		RenameTo:          r.RenameTo,
		DropClusteringKey: r.DropClusteringKey,
		SuspendRecluster:  r.SuspendRecluster,
		ResumeRecluster:   r.ResumeRecluster,
		Suspend:           r.Suspend,
		Resume:            r.Resume,
	}
	if r.ClusterBy != nil {
		opts.ClusterBy = &MaterializedViewClusterBy{}
		if r.ClusterBy.Expressions != nil {
			s := make([]MaterializedViewClusterByExpression, len(r.ClusterBy.Expressions))
			for i, v := range r.ClusterBy.Expressions {
				s[i] = MaterializedViewClusterByExpression(v)
			}
			opts.ClusterBy.Expressions = s
		}
	}
	if r.Set != nil {
		opts.Set = &MaterializedViewSet{
			Secure:  r.Set.Secure,
			Comment: r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &MaterializedViewUnset{
			Secure:  r.Unset.Secure,
			Comment: r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropMaterializedViewRequest) toOpts() *DropMaterializedViewOptions {
	opts := &DropMaterializedViewOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowMaterializedViewRequest) toOpts() *ShowMaterializedViewOptions {
	opts := &ShowMaterializedViewOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r materializedViewDBRow) convert() *MaterializedView {
	materializedView := MaterializedView{
		CreatedOn:          r.CreatedOn,
		Name:               r.Name,
		DatabaseName:       r.DatabaseName,
		SchemaName:         r.SchemaName,
		Rows:               r.Rows,
		Bytes:              r.Bytes,
		SourceDatabaseName: r.SourceDatabaseName,
		SourceSchemaName:   r.SourceSchemaName,
		SourceTableName:    r.SourceTableName,
		RefreshedOn:        r.RefreshedOn,
		CompactedOn:        r.CompactedOn,
		Owner:              r.Owner,
		Invalid:            r.Invalid,
		BehindBy:           r.BehindBy,
		Text:               tracking.TrimMetadata(r.Text),
		IsSecure:           r.IsSecure,
	}
	if r.Reserved.Valid {
		materializedView.Reserved = &r.Reserved.String
	}
	if r.ClusterBy.Valid {
		materializedView.ClusterBy = r.ClusterBy.String
	}
	if r.InvalidReason.Valid {
		materializedView.InvalidReason = r.InvalidReason.String
	}
	if r.Comment.Valid {
		materializedView.Comment = r.Comment.String
	}
	materializedView.AutomaticClustering = r.AutomaticClustering == "ON"
	if r.OwnerRoleType.Valid {
		materializedView.OwnerRoleType = r.OwnerRoleType.String
	}
	if r.Budget.Valid {
		materializedView.Budget = r.Budget.String
	}
	return &materializedView
}

func (r *DescribeMaterializedViewRequest) toOpts() *DescribeMaterializedViewOptions {
	opts := &DescribeMaterializedViewOptions{
		name: r.name,
	}
	return opts
}

func (r materializedViewDetailsRow) convert() *MaterializedViewDetails {
	details := &MaterializedViewDetails{
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
	return details
}
