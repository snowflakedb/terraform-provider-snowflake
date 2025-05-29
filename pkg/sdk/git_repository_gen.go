package sdk

import (
	"context"
	"database/sql"
	"time"
)

type GitRepositories interface {
	Create(ctx context.Context, request *CreateGitRepositoryRequest) error
	Alter(ctx context.Context, request *AlterGitRepositoryRequest) error
	Drop(ctx context.Context, request *DropGitRepositoryRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Describe(ctx context.Context, id SchemaObjectIdentifier) ([]GitRepository, error)
	Show(ctx context.Context, request *ShowGitRepositoryRequest) ([]GitRepository, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*GitRepository, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*GitRepository, error)
	ShowGitBranches(ctx context.Context, request *ShowGitBranchesGitRepositoryRequest) ([]GitBranch, error)
	ShowGitTags(ctx context.Context, request *ShowGitTagsGitRepositoryRequest) ([]GitTag, error)
}

// CreateGitRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-git-repository.
type CreateGitRepositoryOptions struct {
	create         bool                    `ddl:"static" sql:"CREATE"`
	OrReplace      *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	gitRepository  bool                    `ddl:"static" sql:"GIT REPOSITORY"`
	IfNotExists    *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name           SchemaObjectIdentifier  `ddl:"identifier"`
	Origin         string                  `ddl:"parameter,single_quotes" sql:"ORIGIN"`
	ApiIntegration AccountObjectIdentifier `ddl:"identifier,equals" sql:"API_INTEGRATION"`
	GitCredentials *SchemaObjectIdentifier `ddl:"identifier,equals" sql:"GIT_CREDENTIALS"`
	Comment        *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag            []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
}

// AlterGitRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-git-repository.
type AlterGitRepositoryOptions struct {
	alter         bool                   `ddl:"static" sql:"ALTER"`
	gitRepository bool                   `ddl:"static" sql:"GIT REPOSITORY"`
	IfExists      *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
	Set           *GitRepositorySet      `ddl:"keyword" sql:"SET"`
	Unset         *GitRepositoryUnset    `ddl:"list,no_parentheses" sql:"UNSET"`
	Fetch         *bool                  `ddl:"keyword" sql:"FETCH"`
	SetTags       []TagAssociation       `ddl:"keyword" sql:"SET TAG"`
	UnsetTags     []ObjectIdentifier     `ddl:"keyword" sql:"UNSET TAG"`
}

type GitRepositorySet struct {
	ApiIntegration *AccountObjectIdentifier `ddl:"identifier,equals" sql:"API_INTEGRATION"`
	GitCredentials *SchemaObjectIdentifier  `ddl:"identifier,equals" sql:"GIT_CREDENTIALS"`
	Comment        *string                  `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type GitRepositoryUnset struct {
	GitCredentials *bool `ddl:"keyword" sql:"GIT_CREDENTIALS"`
	Comment        *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropGitRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-git-repository.
type DropGitRepositoryOptions struct {
	drop          bool                   `ddl:"static" sql:"DROP"`
	gitRepository bool                   `ddl:"static" sql:"GIT REPOSITORY"`
	IfExists      *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

// DescribeGitRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-git-repository.
type DescribeGitRepositoryOptions struct {
	describe      bool                   `ddl:"static" sql:"DESCRIBE"`
	gitRepository bool                   `ddl:"static" sql:"GIT REPOSITORY"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

type gitRepositoriesRow struct {
	CreatedOn      time.Time      `db:"created_on"`
	Name           string         `db:"name"`
	DatabaseName   string         `db:"database_name"`
	SchemaName     string         `db:"schema_name"`
	Origin         string         `db:"origin"`
	ApiIntegration string         `db:"api_integration"`
	GitCredentials sql.NullString `db:"git_credentials"`
	Owner          string         `db:"owner"`
	OwnerRoleType  string         `db:"owner_role_type"`
	Comment        sql.NullString `db:"comment"`
	LastFetchedAt  sql.NullTime   `db:"last_fetched_at"`
}

type GitRepository struct {
	CreatedOn      time.Time
	Name           string
	DatabaseName   string
	SchemaName     string
	Origin         string
	ApiIntegration *AccountObjectIdentifier
	GitCredentials *SchemaObjectIdentifier
	Owner          string
	OwnerRoleType  string
	Comment        *string
	LastFetchedAt  *time.Time
}

// ShowGitRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-git-repositories.
type ShowGitRepositoryOptions struct {
	show            bool       `ddl:"static" sql:"SHOW"`
	gitRepositories bool       `ddl:"static" sql:"GIT REPOSITORIES"`
	Like            *Like      `ddl:"keyword" sql:"LIKE"`
	In              *In        `ddl:"keyword" sql:"IN"`
	Limit           *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

func (v *GitRepository) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}
func (v *GitRepository) ObjectType() ObjectType {
	return ObjectTypeGitRepository
}

// ShowGitBranchesGitRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-git-branches.
type ShowGitBranchesGitRepositoryOptions struct {
	showGitBranches bool                   `ddl:"static" sql:"SHOW GIT BRANCHES"`
	Like            *Like                  `ddl:"keyword" sql:"LIKE"`
	in              bool                   `ddl:"static" sql:"IN"`
	GitRepository   *bool                  `ddl:"keyword" sql:"GIT REPOSITORY"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
}

type gitBranchesRow struct {
	Name       string `db:"name"`
	Path       string `db:"path"`
	Checkouts  string `db:"checkouts"`
	CommitHash string `db:"commit_hash"`
}

type GitBranch struct {
	Name       string
	Path       string
	Checkouts  string
	CommitHash string
}

// ShowGitTagsGitRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-git-tags.
type ShowGitTagsGitRepositoryOptions struct {
	showGitTags   bool                   `ddl:"static" sql:"SHOW GIT TAGS"`
	Like          *Like                  `ddl:"keyword" sql:"LIKE"`
	in            bool                   `ddl:"static" sql:"IN"`
	GitRepository *bool                  `ddl:"keyword" sql:"GIT REPOSITORY"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
}

type gitTagsRow struct {
	Name       string `db:"name"`
	Path       string `db:"path"`
	CommitHash string `db:"commit_hash"`
	Author     string `db:"author"`
	Message    string `db:"message"`
}

type GitTag struct {
	Name       string
	Path       string
	CommitHash string
	Author     string
	Message    string
}
