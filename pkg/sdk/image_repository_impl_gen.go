package sdk

import (
	"context"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ ImageRepositories = (*imageRepositories)(nil)

type imageRepositories struct {
	client *Client
}

func (v *imageRepositories) Create(ctx context.Context, request *CreateImageRepositoryRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *imageRepositories) Alter(ctx context.Context, request *AlterImageRepositoryRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *imageRepositories) Drop(ctx context.Context, request *DropImageRepositoryRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *imageRepositories) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropImageRepositoryRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *imageRepositories) Show(ctx context.Context, request *ShowImageRepositoryRequest) ([]ImageRepository, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[imageRepositoriesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[imageRepositoriesRow, ImageRepository](dbRows)
	return resultList, nil
}

func (v *imageRepositories) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*ImageRepository, error) {
	request := NewShowImageRepositoryRequest().
		WithIn(In{Schema: id.SchemaId()}).
		WithLike(Like{Pattern: String(id.Name())})
	imageRepositories, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(imageRepositories, func(r ImageRepository) bool { return r.Name == id.Name() })
}

func (v *imageRepositories) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*ImageRepository, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (r *CreateImageRepositoryRequest) toOpts() *CreateImageRepositoryOptions {
	opts := &CreateImageRepositoryOptions{
		OrReplace:   r.OrReplace,
		IfNotExists: r.IfNotExists,
		name:        r.name,
		Comment:     r.Comment,
	}
	return opts
}

func (r *AlterImageRepositoryRequest) toOpts() *AlterImageRepositoryOptions {
	opts := &AlterImageRepositoryOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	if r.Set != nil {
		opts.Set = &ImageRepositorySet{
			Comment: r.Set.Comment,
		}
	}
	return opts
}

func (r *DropImageRepositoryRequest) toOpts() *DropImageRepositoryOptions {
	opts := &DropImageRepositoryOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *ShowImageRepositoryRequest) toOpts() *ShowImageRepositoryOptions {
	opts := &ShowImageRepositoryOptions{
		Like: r.Like,
		In:   r.In,
	}
	return opts
}

func (r imageRepositoriesRow) convert() *ImageRepository {
	return &ImageRepository{
		CreatedOn:                r.CreatedOn,
		Name:                     r.Name,
		DatabaseName:             r.DatabaseName,
		SchemaName:               r.SchemaName,
		RepositoryUrl:            r.RepositoryUrl,
		Owner:                    r.Owner,
		OwnerRoleType:            r.OwnerRoleType,
		Comment:                  r.Comment,
		PrivatelinkRepositoryUrl: r.PrivatelinkRepositoryUrl,
	}
}
