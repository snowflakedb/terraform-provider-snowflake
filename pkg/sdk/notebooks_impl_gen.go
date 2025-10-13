package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Notebooks = (*notebooks)(nil)

var (
	_ convertibleRow[Notebook] = new(notebooksRow)
	_ convertibleRow[Notebook] = new(notebooksRow)
)

type notebooks struct {
	client *Client
}

func (v *notebooks) Create(ctx context.Context, request *CreateNotebookRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *notebooks) Alter(ctx context.Context, request *AlterNotebookRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *notebooks) Drop(ctx context.Context, request *DropNotebookRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *notebooks) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropNotebookRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *notebooks) Describe(ctx context.Context, id SchemaObjectIdentifier) (*Notebook, error) {
	opts := &DescribeNotebookOptions{
		name: id,
	}
	result, err := validateAndQueryOne[notebooksRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert()
}

func (v *notebooks) Show(ctx context.Context, request *ShowNotebookRequest) ([]Notebook, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[notebooksRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[notebooksRow, Notebook](dbRows)
}

func (v *notebooks) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Notebook, error) {
	request := NewShowNotebookRequest().
		WithIn(In{Schema: id.SchemaId()}).
		WithLike(Like{Pattern: String(id.Name())})
	notebooks, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(notebooks, func(r Notebook) bool { return r.Name == id.Name() })
}

func (v *notebooks) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Notebook, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *CreateNotebookRequest) toOpts() *CreateNotebookOptions {
	opts := &CreateNotebookOptions{
		OrReplace:                   r.OrReplace,
		IfNotExists:                 r.IfNotExists,
		name:                        r.name,
		From:                        r.From,
		Mainfile:                    r.Mainfile,
		Comment:                     r.Comment,
		QueryWarehouse:              r.QueryWarehouse,
		IdleAutoShutdownTimeSeconds: r.IdleAutoShutdownTimeSeconds,
		Warehouse:                   r.Warehouse,
		RuntimeName:                 r.RuntimeName,
		ComputePool:                 r.ComputePool,
		Externalaccessintegrations:  r.Externalaccessintegrations,
		RuntimeEnvironmentVersion:   r.RuntimeEnvironmentVersion,
		DefaultVersion:              r.DefaultVersion,
	}
	return opts
}

func (r *AlterNotebookRequest) toOpts() *AlterNotebookOptions {
	opts := &AlterNotebookOptions{
		IfExists: r.IfExists,
		name:     r.name,
		RenameTo: r.RenameTo,
	}
	if r.Set != nil {
		opts.Set = &NotebookSet{
			Comment:                     r.Set.Comment,
			QueryWarehouse:              r.Set.QueryWarehouse,
			IdleAutoShutdownTimeSeconds: r.Set.IdleAutoShutdownTimeSeconds,

			Mainfile:                   r.Set.Mainfile,
			Warehouse:                  r.Set.Warehouse,
			RuntimeName:                r.Set.RuntimeName,
			ComputePool:                r.Set.ComputePool,
			Externalaccessintegrations: r.Set.Externalaccessintegrations,
			RuntimeEnvironmentVersion:  r.Set.RuntimeEnvironmentVersion,
		}
		if r.Set.SecretsList != nil {
			opts.Set.SecretsList = &SecretsList{
				SecretsList: r.Set.SecretsList.SecretsList,
			}
		}
	}
	if r.Unset != nil {
		opts.Unset = &NotebookUnset{
			Comment:                    r.Unset.Comment,
			QueryWarehouse:             r.Unset.QueryWarehouse,
			Secrets:                    r.Unset.Secrets,
			Warehouse:                  r.Unset.Warehouse,
			RuntimeName:                r.Unset.RuntimeName,
			ComputePool:                r.Unset.ComputePool,
			ExternalAccessIntegrations: r.Unset.ExternalAccessIntegrations,
			RuntimeEnvironmentVersion:  r.Unset.RuntimeEnvironmentVersion,
		}
	}
	return opts
}

func (r *DropNotebookRequest) toOpts() *DropNotebookOptions {
	opts := &DropNotebookOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeNotebookRequest) toOpts() *DescribeNotebookOptions {
	opts := &DescribeNotebookOptions{
		name: r.name,
	}
	return opts
}

func (r notebooksRow) convert() (*Notebook, error) {
	return &Notebook{
		CreatedOn:      r.CreatedOn,
		Name:           r.Name,
		DatabaseName:   r.DatabaseName,
		SchemaName:     r.SchemaName,
		Comment:        &r.Comment.String,
		Owner:          r.Owner,
		QueryWarehouse: &AccountObjectIdentifier{r.QueryWarehouse.String},
		UrlId:          r.UrlId,
		OwnerRoleType:  r.OwnerRoleType,
		CodeWarehouse:  AccountObjectIdentifier{r.CodeWarehouse},
	}, nil
}

func (r *ShowNotebookRequest) toOpts() *ShowNotebookOptions {
	opts := &ShowNotebookOptions{
		Like:       r.Like,
		In:         r.In,
		Limit:      r.Limit,
		StartsWith: r.StartsWith,
	}
	return opts
}
