package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ Notebooks = (*notebooks)(nil)

var (
	_ convertibleRow[NotebookDetails] = new(NotebookDetailsRow)
	_ convertibleRow[Notebook]        = new(notebookRow)
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

func (v *notebooks) Describe(ctx context.Context, id SchemaObjectIdentifier) (*NotebookDetails, error) {
	opts := &DescribeNotebookOptions{
		name: id,
	}
	result, err := validateAndQueryOne[NotebookDetailsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert()
}

func (v *notebooks) Show(ctx context.Context, request *ShowNotebookRequest) ([]Notebook, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[notebookRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return convertRows[notebookRow, Notebook](dbRows)
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
		Title:                       r.Title,
		MainFile:                    r.MainFile,
		Comment:                     r.Comment,
		QueryWarehouse:              r.QueryWarehouse,
		IdleAutoShutdownTimeSeconds: r.IdleAutoShutdownTimeSeconds,
		Warehouse:                   r.Warehouse,
		RuntimeName:                 r.RuntimeName,
		ComputePool:                 r.ComputePool,
		ExternalAccessIntegrations:  r.ExternalAccessIntegrations,
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

		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &NotebookSet{
			Comment:                     r.Set.Comment,
			QueryWarehouse:              r.Set.QueryWarehouse,
			IdleAutoShutdownTimeSeconds: r.Set.IdleAutoShutdownTimeSeconds,

			MainFile:                   r.Set.MainFile,
			Warehouse:                  r.Set.Warehouse,
			RuntimeName:                r.Set.RuntimeName,
			ComputePool:                r.Set.ComputePool,
			ExternalAccessIntegrations: r.Set.ExternalAccessIntegrations,
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

func (r NotebookDetailsRow) convert() (*NotebookDetails, error) {
	n := &NotebookDetails{
		MainFile:                    r.MainFile,
		UrlId:                       r.UrlId,
		DefaultPackages:             r.DefaultPackages,
		Owner:                       r.Owner,
		ImportUrls:                  r.ImportUrls,
		ExternalAccessIntegrations:  r.ExternalAccessIntegrations,
		ExternalAccessSecrets:       r.ExternalAccessSecrets,
		CodeWarehouse:               r.CodeWarehouse,
		IdleAutoShutdownTimeSeconds: r.IdleAutoShutdownTimeSeconds,
		RuntimeEnvironmentVersion:   r.RuntimeEnvironmentVersion,
		Name:                        r.Name,
		DefaultVersion:              r.DefaultVersion,
		DefaultVersionName:          r.DefaultVersionName,
		DefaultVersionLocationUri:   r.DefaultVersionLocationUri,
		LastVersionName:             r.LastVersionName,
		LastVersionLocationUri:      r.LastVersionLocationUri,
	}

	// Optionals.
	mapNullString(&n.Title, r.Title)
	mapNullStringWithMapping(&n.QueryWarehouse, r.QueryWarehouse, ParseAccountObjectIdentifier)
	mapNullString(&n.UserPackages, r.UserPackages)
	mapNullString(&n.RuntimeName, r.RuntimeName)
	mapNullString(&n.ComputePool, r.ComputePool)
	mapNullString(&n.Comment, r.Comment)
	mapNullString(&n.DefaultVersionAlias, r.DefaultVersionAlias)
	mapNullString(&n.DefaultVersionSourceLocationUri, r.DefaultVersionSourceLocationUri)
	mapNullString(&n.DefaultVersionGitCommitHash, r.DefaultVersionGitCommitHash)
	mapNullString(&n.LastVersionAlias, r.LastVersionAlias)
	mapNullString(&n.LastVersionSourceLocationUri, r.LastVersionSourceLocationUri)
	mapNullString(&n.LastVersionGitCommitHash, r.LastVersionGitCommitHash)
	mapNullString(&n.LiveVersionLocationUri, r.LiveVersionLocationUri)

	return n, nil
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

func (r notebookRow) convert() (*Notebook, error) {
	n := &Notebook{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		DatabaseName:  r.DatabaseName,
		SchemaName:    r.SchemaName,
		Owner:         r.Owner,
		UrlId:         r.UrlId,
		OwnerRoleType: r.OwnerRoleType,
		CodeWarehouse: AccountObjectIdentifier{r.CodeWarehouse},
	}

	mapNullString(&n.Comment, r.Comment)
	mapNullStringWithMapping(&n.QueryWarehouse, r.QueryWarehouse, ParseAccountObjectIdentifier)

	return n, nil
}
