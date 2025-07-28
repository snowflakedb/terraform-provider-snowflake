package sdk

import (
	"context"
	"log"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
)

var _ GitRepositories = (*gitRepositories)(nil)

type gitRepositories struct {
	client *Client
}

func (v *gitRepositories) Create(ctx context.Context, request *CreateGitRepositoryRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *gitRepositories) Alter(ctx context.Context, request *AlterGitRepositoryRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *gitRepositories) Drop(ctx context.Context, request *DropGitRepositoryRequest) error {
	opts := request.toOpts()
	return validateAndExec(v.client, ctx, opts)
}

func (v *gitRepositories) DropSafely(ctx context.Context, id SchemaObjectIdentifier) error {
	return SafeDrop(v.client, func() error { return v.Drop(ctx, NewDropGitRepositoryRequest(id).WithIfExists(true)) }, ctx, id)
}

func (v *gitRepositories) Describe(ctx context.Context, id SchemaObjectIdentifier) (*GitRepository, error) {
	opts := &DescribeGitRepositoryOptions{
		name: id,
	}
	result, err := validateAndQueryOne[gitRepositoriesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	return result.convert(), nil
}

func (v *gitRepositories) Show(ctx context.Context, request *ShowGitRepositoryRequest) ([]GitRepository, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[gitRepositoriesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[gitRepositoriesRow, GitRepository](dbRows)
	return resultList, nil
}

func (v *gitRepositories) ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*GitRepository, error) {
	request := NewShowGitRepositoryRequest().
		WithIn(In{Schema: id.SchemaId()}).
		WithLike(Like{Pattern: String(id.Name())})
	gitRepositories, err := v.Show(ctx, request)
	if err != nil {
		return nil, err
	}
	return collections.FindFirst(gitRepositories, func(r GitRepository) bool { return r.Name == id.Name() })
}

func (v *gitRepositories) ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*GitRepository, error) {
	return SafeShowById(v.client, v.ShowByID, ctx, id)
}

func (v *gitRepositories) ShowGitBranches(ctx context.Context, request *ShowGitBranchesGitRepositoryRequest) ([]GitBranch, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[gitBranchesRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[gitBranchesRow, GitBranch](dbRows)
	return resultList, nil
}

func (v *gitRepositories) ShowGitTags(ctx context.Context, request *ShowGitTagsGitRepositoryRequest) ([]GitTag, error) {
	opts := request.toOpts()
	dbRows, err := validateAndQuery[gitTagsRow](v.client, ctx, opts)
	if err != nil {
		return nil, err
	}
	resultList := convertRows[gitTagsRow, GitTag](dbRows)
	return resultList, nil
}

func (r *CreateGitRepositoryRequest) toOpts() *CreateGitRepositoryOptions {
	opts := &CreateGitRepositoryOptions{
		OrReplace:      r.OrReplace,
		IfNotExists:    r.IfNotExists,
		name:           r.name,
		Origin:         r.Origin,
		ApiIntegration: r.ApiIntegration,
		GitCredentials: r.GitCredentials,
		Comment:        r.Comment,
		Tag:            r.Tag,
	}
	return opts
}

func (r *AlterGitRepositoryRequest) toOpts() *AlterGitRepositoryOptions {
	opts := &AlterGitRepositoryOptions{
		IfExists: r.IfExists,
		name:     r.name,

		Fetch:     r.Fetch,
		SetTags:   r.SetTags,
		UnsetTags: r.UnsetTags,
	}
	if r.Set != nil {
		opts.Set = &GitRepositorySet{
			ApiIntegration: r.Set.ApiIntegration,
			GitCredentials: r.Set.GitCredentials,
			Comment:        r.Set.Comment,
		}
	}
	if r.Unset != nil {
		opts.Unset = &GitRepositoryUnset{
			GitCredentials: r.Unset.GitCredentials,
			Comment:        r.Unset.Comment,
		}
	}
	return opts
}

func (r *DropGitRepositoryRequest) toOpts() *DropGitRepositoryOptions {
	opts := &DropGitRepositoryOptions{
		IfExists: r.IfExists,
		name:     r.name,
	}
	return opts
}

func (r *DescribeGitRepositoryRequest) toOpts() *DescribeGitRepositoryOptions {
	opts := &DescribeGitRepositoryOptions{
		name: r.name,
	}
	return opts
}

func (r gitRepositoriesRow) convert() *GitRepository {
	gitRepository := &GitRepository{
		CreatedOn:     r.CreatedOn,
		Name:          r.Name,
		DatabaseName:  r.DatabaseName,
		SchemaName:    r.SchemaName,
		Origin:        r.Origin,
		Owner:         r.Owner,
		OwnerRoleType: r.OwnerRoleType,
	}
	id, err := ParseAccountObjectIdentifier(r.ApiIntegration)
	if err != nil {
		log.Printf("[DEBUG] failed to parse api integration in git repository: %v", err)
	} else {
		gitRepository.ApiIntegration = &id
	}

	if r.GitCredentials.Valid {
		id, err := ParseSchemaObjectIdentifier(r.GitCredentials.String)
		if err != nil {
			log.Printf("[DEBUG] failed to parse git credentials in git repository: %v", err)
		} else {
			gitRepository.GitCredentials = &id
		}
	}

	if r.Comment.Valid {
		gitRepository.Comment = &r.Comment.String
	}

	if r.LastFetchedAt.Valid {
		gitRepository.LastFetchedAt = &r.LastFetchedAt.Time
	}

	return gitRepository
}

func (r *ShowGitRepositoryRequest) toOpts() *ShowGitRepositoryOptions {
	opts := &ShowGitRepositoryOptions{
		Like:  r.Like,
		In:    r.In,
		Limit: r.Limit,
	}
	return opts
}

func (r *ShowGitBranchesGitRepositoryRequest) toOpts() *ShowGitBranchesGitRepositoryOptions {
	opts := &ShowGitBranchesGitRepositoryOptions{
		Like:          r.Like,
		GitRepository: r.GitRepository,
		name:          r.name,
	}
	return opts
}

func (r gitBranchesRow) convert() *GitBranch {
	return &GitBranch{
		Name:       r.Name,
		Path:       r.Path,
		Checkouts:  r.Checkouts,
		CommitHash: r.CommitHash,
	}
}

func (r *ShowGitTagsGitRepositoryRequest) toOpts() *ShowGitTagsGitRepositoryOptions {
	opts := &ShowGitTagsGitRepositoryOptions{
		Like:          r.Like,
		GitRepository: r.GitRepository,
		name:          r.name,
	}
	return opts
}

func (r gitTagsRow) convert() *GitTag {
	return &GitTag{
		Name:       r.Name,
		Path:       r.Path,
		CommitHash: r.CommitHash,
		Author:     r.Author,
		Message:    r.Message,
	}
}
