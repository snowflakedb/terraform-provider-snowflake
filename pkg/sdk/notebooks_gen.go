package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Notebooks interface {
	Create(ctx context.Context, request *CreateNotebookRequest) error
	Alter(ctx context.Context, request *AlterNotebookRequest) error
	Drop(ctx context.Context, request *DropNotebookRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*NotebookDetails, error)
	Show(ctx context.Context, request *ShowNotebookRequest) ([]Notebook, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Notebook, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Notebook, error)
}

// CreateNotebookOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-notebook.
type CreateNotebookOptions struct {
	create                      bool                      `ddl:"static" sql:"CREATE"`
	OrReplace                   *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	notebook                    bool                      `ddl:"static" sql:"NOTEBOOK"`
	IfNotExists                 *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name                        SchemaObjectIdentifier    `ddl:"identifier"`
	From                        *Location                 `ddl:"parameter,single_quotes,no_equals" sql:"FROM"`
	Title                       *string                   `ddl:"parameter,single_quotes" sql:"TITLE"`
	MainFile                    *string                   `ddl:"parameter,single_quotes" sql:"MAIN_FILE"`
	Comment                     *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	QueryWarehouse              *AccountObjectIdentifier  `ddl:"identifier,equals" sql:"QUERY_WAREHOUSE"`
	IdleAutoShutdownTimeSeconds *int                      `ddl:"parameter,no_quotes" sql:"IDLE_AUTO_SHUTDOWN_TIME_SECONDS"`
	Warehouse                   *AccountObjectIdentifier  `ddl:"identifier,equals" sql:"WAREHOUSE"`
	RuntimeName                 *string                   `ddl:"parameter,single_quotes" sql:"RUNTIME_NAME"`
	ComputePool                 *AccountObjectIdentifier  `ddl:"identifier,equals" sql:"COMPUTE_POOL"`
	ExternalAccessIntegrations  []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	RuntimeEnvironmentVersion   *string                   `ddl:"parameter,single_quotes" sql:"RUNTIME_ENVIRONMENT_VERSION"`
	DefaultVersion              *string                   `ddl:"parameter,no_quotes" sql:"DEFAULT_VERSION"`
}

// AlterNotebookOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-notebook.
type AlterNotebookOptions struct {
	alter     bool                    `ddl:"static" sql:"ALTER"`
	notebook  bool                    `ddl:"static" sql:"NOTEBOOK"`
	IfExists  *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name      SchemaObjectIdentifier  `ddl:"identifier"`
	RenameTo  *SchemaObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	Set       *NotebookSet            `ddl:"keyword" sql:"SET"`
	Unset     *NotebookUnset          `ddl:"list,no_parentheses" sql:"UNSET"`
	SetTags   []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTags []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

type NotebookSet struct {
	Comment                     *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
	QueryWarehouse              *AccountObjectIdentifier  `ddl:"identifier,equals" sql:"QUERY_WAREHOUSE"`
	IdleAutoShutdownTimeSeconds *int                      `ddl:"parameter,no_quotes" sql:"IDLE_AUTO_SHUTDOWN_TIME_SECONDS"`
	Secrets                     *SecretsList              `ddl:"parameter,parentheses" sql:"SECRETS"`
	MainFile                    *string                   `ddl:"parameter,single_quotes" sql:"MAIN_FILE"`
	Warehouse                   *AccountObjectIdentifier  `ddl:"identifier,equals" sql:"WAREHOUSE"`
	RuntimeName                 *string                   `ddl:"parameter,single_quotes" sql:"RUNTIME_NAME"`
	ComputePool                 *AccountObjectIdentifier  `ddl:"identifier,equals" sql:"COMPUTE_POOL"`
	ExternalAccessIntegrations  []AccountObjectIdentifier `ddl:"parameter,parentheses" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	RuntimeEnvironmentVersion   *string                   `ddl:"parameter,single_quotes" sql:"RUNTIME_ENVIRONMENT_VERSION"`
}

type NotebookUnset struct {
	Comment                    *bool `ddl:"keyword" sql:"COMMENT"`
	QueryWarehouse             *bool `ddl:"keyword" sql:"QUERY_WAREHOUSE"`
	Secrets                    *bool `ddl:"keyword" sql:"SECRETS"`
	Warehouse                  *bool `ddl:"keyword" sql:"WAREHOUSE"`
	RuntimeName                *bool `ddl:"keyword" sql:"RUNTIME_NAME"`
	ComputePool                *bool `ddl:"keyword" sql:"COMPUTE_POOL"`
	ExternalAccessIntegrations *bool `ddl:"keyword" sql:"EXTERNAL_ACCESS_INTEGRATIONS"`
	RuntimeEnvironmentVersion  *bool `ddl:"keyword" sql:"RUNTIME_ENVIRONMENT_VERSION"`
}

// DropNotebookOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-notebook.
type DropNotebookOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	notebook bool                   `ddl:"static" sql:"NOTEBOOK"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// DescribeNotebookOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-notebook.
type DescribeNotebookOptions struct {
	describe bool                   `ddl:"static" sql:"DESCRIBE"`
	notebook bool                   `ddl:"static" sql:"NOTEBOOK"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

type NotebookDetailsRow struct {
	Title                           sql.NullString `db:"title"`
	MainFile                        string         `db:"main_file"`
	QueryWarehouse                  sql.NullString `db:"query_warehouse"`
	UrlId                           string         `db:"url_id"`
	DefaultPackages                 string         `db:"default_packages"`
	UserPackages                    sql.NullString `db:"user_packages"`
	RuntimeName                     sql.NullString `db:"runtime_name"`
	ComputePool                     sql.NullString `db:"compute_pool"`
	Owner                           string         `db:"owner"`
	ImportUrls                      string         `db:"import_urls"`
	ExternalAccessIntegrations      string         `db:"external_access_integrations"`
	ExternalAccessSecrets           string         `db:"external_access_secrets"`
	CodeWarehouse                   string         `db:"code_warehouse"`
	IdleAutoShutdownTimeSeconds     int            `db:"idle_auto_shutdown_time_seconds"`
	RuntimeEnvironmentVersion       string         `db:"runtime_environment_version"`
	Name                            string         `db:"name"`
	Comment                         sql.NullString `db:"comment"`
	DefaultVersion                  string         `db:"default_version"`
	DefaultVersionName              string         `db:"default_version_name"`
	DefaultVersionAlias             sql.NullString `db:"default_version_alias"`
	DefaultVersionLocationUri       string         `db:"default_version_location_uri"`
	DefaultVersionSourceLocationUri sql.NullString `db:"default_version_source_location_uri"`
	DefaultVersionGitCommitHash     sql.NullString `db:"default_version_git_commit_hash"`
	LastVersionName                 string         `db:"last_version_name"`
	LastVersionAlias                sql.NullString `db:"last_version_alias"`
	LastVersionLocationUri          string         `db:"last_version_location_uri"`
	LastVersionSourceLocationUri    sql.NullString `db:"last_version_source_location_uri"`
	LastVersionGitCommitHash        sql.NullString `db:"last_version_git_commit_hash"`
	LiveVersionLocationUri          sql.NullString `db:"live_version_location_uri"`
}

type NotebookDetails struct {
	Title                           *string
	MainFile                        string
	QueryWarehouse                  *AccountObjectIdentifier
	UrlId                           string
	DefaultPackages                 string
	UserPackages                    *string
	RuntimeName                     *string
	ComputePool                     *AccountObjectIdentifier
	Owner                           string
	ImportUrls                      string
	ExternalAccessIntegrations      string
	ExternalAccessSecrets           string
	CodeWarehouse                   string
	IdleAutoShutdownTimeSeconds     int
	RuntimeEnvironmentVersion       string
	Name                            string
	Comment                         *string
	DefaultVersion                  string
	DefaultVersionName              string
	DefaultVersionAlias             *string
	DefaultVersionLocationUri       string
	DefaultVersionSourceLocationUri *string
	DefaultVersionGitCommitHash     *string
	LastVersionName                 string
	LastVersionAlias                *string
	LastVersionLocationUri          string
	LastVersionSourceLocationUri    *string
	LastVersionGitCommitHash        *string
	LiveVersionLocationUri          *string
}

// ShowNotebookOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-notebooks.
type ShowNotebookOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	notebooks  bool       `ddl:"static" sql:"NOTEBOOKS"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	In         *In        `ddl:"keyword" sql:"IN"`
	Limit      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
}

type notebookRow struct {
	CreatedOn      time.Time      `db:"created_on"`
	Name           string         `db:"name"`
	DatabaseName   string         `db:"database_name"`
	SchemaName     string         `db:"schema_name"`
	Comment        sql.NullString `db:"comment"`
	Owner          string         `db:"owner"`
	QueryWarehouse sql.NullString `db:"query_warehouse"`
	UrlId          string         `db:"url_id"`
	OwnerRoleType  string         `db:"owner_role_type"`
	CodeWarehouse  string         `db:"code_warehouse"`
}

type Notebook struct {
	CreatedOn      time.Time
	Name           string
	DatabaseName   string
	SchemaName     string
	Comment        *string
	Owner          string
	QueryWarehouse *AccountObjectIdentifier
	UrlId          string
	OwnerRoleType  string
	CodeWarehouse  AccountObjectIdentifier
}

func (v *Notebook) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *Notebook) ObjectType() ObjectType {
	return ObjectTypeNotebook
}
