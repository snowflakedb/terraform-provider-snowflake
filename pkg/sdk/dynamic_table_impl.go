package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ DynamicTables = (*dynamicTables)(nil)
var _ convertibleRow[DynamicTable] = new(dynamicTableRow)
var _ convertibleRow[DynamicTableDetails] = new(dynamicTableDetailsRow)

type dynamicTables struct {
	client *Client
}

func (v *dynamicTables) Create(ctx context.Context, request *CreateDynamicTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *dynamicTables) Alter(ctx context.Context, request *AlterDynamicTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *dynamicTables) Drop(ctx context.Context, request *DropDynamicTableRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *dynamicTables) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropDynamicTableRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *dynamicTables) Describe(ctx context.Context, request *DescribeDynamicTableRequest) (*DynamicTableDetails, error) {
	opts := request.toOpts()
	row, err := validateAndQueryOne[dynamicTableDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return conversionErrorWrapped(row.convert())
}

func (v *dynamicTables) Show(ctx context.Context, request *ShowDynamicTableRequest) ([]DynamicTable, error) {
	opts := request.toOpts()
	rows, err := validateAndQuery[dynamicTableRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[dynamicTableRow, DynamicTable](rows)
}

func (v *dynamicTables) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*DynamicTable, error) {
	request := NewShowDynamicTableRequest().WithIn(&In{Schema: id.SchemaId()}).WithLike(&Like{String(id.Name())})
	dynamicTables, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(dynamicTables, func(r DynamicTable) bool { return r.Name == id.Name() })
}

func (v *dynamicTables) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*DynamicTable, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (s *CreateDynamicTableRequest) toOpts() *createDynamicTableOptions {
	return &createDynamicTableOptions{
		OrReplace:   Bool(s.orReplace),
		name:        s.name,
		warehouse:   s.warehouse,
		targetLag:   s.targetLag,
		query:       s.query,
		Comment:     s.comment,
		RefreshMode: s.refreshMode,
		Initialize:  s.initialize,
	}
}

func (s *AlterDynamicTableRequest) toOpts() *alterDynamicTableOptions {
	opts := alterDynamicTableOptions{
		name: s.name,
	}
	if s.suspend != nil {
		opts.Suspend = s.suspend
	}
	if s.resume != nil {
		opts.Resume = s.resume
	}
	if s.refresh != nil {
		opts.Refresh = s.refresh
	}
	if s.set != nil {
		opts.Set = &DynamicTableSet{s.set.targetLag, s.set.warehouse}
	}
	return &opts
}

func (s *DropDynamicTableRequest) toOpts() *dropDynamicTableOptions {
	return &dropDynamicTableOptions{
		name:     s.name,
		IfExists: s.IfExists,
	}
}

func (s *DescribeDynamicTableRequest) toOpts() *describeDynamicTableOptions {
	return &describeDynamicTableOptions{
		name: s.name,
	}
}

func (s *ShowDynamicTableRequest) toOpts() *showDynamicTableOptions {
	opts := showDynamicTableOptions{}
	if s.like != nil {
		opts.Like = s.like
	}
	if s.in != nil {
		opts.In = s.in
	}
	if s.startsWith != nil {
		opts.StartsWith = s.startsWith
	}
	if s.limit != nil {
		opts.Limit = s.limit
	}
	return &opts
}
