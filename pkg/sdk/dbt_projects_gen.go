package sdk

import (
	"context"
	"database/sql"
	"time"
)

type DbtProjects interface {
	Create(ctx context.Context, request *CreateDbtProjectRequest) error
	Alter(ctx context.Context, request *AlterDbtProjectRequest) error
	Drop(ctx context.Context, request *DropDbtProjectRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, request *ShowDbtProjectRequest) ([]DbtProject, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*DbtProject, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*DbtProject, error)
	Describe(ctx context.Context, id SchemaObjectIdentifier) (*DbtProjectDetails, error)
}

// CreateDbtProjectOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-dbt-project.
type CreateDbtProjectOptions struct {
	create         bool                      `ddl:"static" sql:"CREATE"`
	OrReplace      *bool                     `ddl:"keyword" sql:"OR REPLACE"`
	dbtProject     bool                      `ddl:"static" sql:"DBT PROJECT"`
	IfNotExists    *bool                     `ddl:"keyword" sql:"IF NOT EXISTS"`
	name           SchemaObjectIdentifier    `ddl:"identifier"`
	From           *string                   `ddl:"parameter,single_quotes,no_equals" sql:"FROM"`
	DefaultArgs    *string                   `ddl:"parameter,single_quotes" sql:"DEFAULT_ARGS"`
	DefaultVersion *DbtProjectDefaultVersion `ddl:"parameter,no_quotes" sql:"DEFAULT_VERSION"`
	Comment        *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterDbtProjectOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-dbt-project.
type AlterDbtProjectOptions struct {
	alter      bool                   `ddl:"static" sql:"ALTER"`
	dbtProject bool                   `ddl:"static" sql:"DBT PROJECT"`
	IfExists   *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
	Set        *DbtProjectSet         `ddl:"keyword" sql:"SET"`
	Unset      *DbtProjectUnset       `ddl:"keyword" sql:"UNSET"`
}

type DbtProjectSet struct {
	DefaultArgs    *string                   `ddl:"parameter,single_quotes" sql:"DEFAULT_ARGS"`
	DefaultVersion *DbtProjectDefaultVersion `ddl:"parameter,no_quotes" sql:"DEFAULT_VERSION"`
	Comment        *string                   `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type DbtProjectUnset struct {
	DefaultArgs    *bool `ddl:"keyword" sql:"DEFAULT_ARGS"`
	DefaultVersion *bool `ddl:"keyword" sql:"DEFAULT_VERSION"`
	Comment        *bool `ddl:"keyword" sql:"COMMENT"`
}

// DropDbtProjectOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-dbt-project.
type DropDbtProjectOptions struct {
	drop       bool                   `ddl:"static" sql:"DROP"`
	dbtProject bool                   `ddl:"static" sql:"DBT PROJECT"`
	IfExists   *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowDbtProjectOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-dbt-projects.
type ShowDbtProjectOptions struct {
	show        bool  `ddl:"static" sql:"SHOW"`
	dbtProjects bool  `ddl:"static" sql:"DBT PROJECTS"`
	Like        *Like `ddl:"keyword" sql:"LIKE"`
	In          *In   `ddl:"keyword" sql:"IN"`
}

type dbtProjectDBRow struct {
	CreatedOn      time.Time      `db:"created_on"`
	Name           string         `db:"name"`
	DatabaseName   string         `db:"database_name"`
	SchemaName     string         `db:"schema_name"`
	SourceLocation sql.NullString `db:"source_location"`
	DefaultArgs    sql.NullString `db:"default_args"`
	DefaultVersion sql.NullString `db:"default_version"`
	Owner          string         `db:"owner"`
	OwnerRoleType  string         `db:"owner_role_type"`
	Comment        sql.NullString `db:"comment"`
}

type DbtProject struct {
	CreatedOn      time.Time
	Name           string
	DatabaseName   string
	SchemaName     string
	SourceLocation *string
	DefaultArgs    *string
	DefaultVersion *string
	Owner          string
	OwnerRoleType  string
	Comment        *string
}

func (v *DbtProject) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *DbtProject) ObjectType() ObjectType {
	return ObjectTypeDbtProject
}

// DescribeDbtProjectOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-dbt-project.
type DescribeDbtProjectOptions struct {
	describe   bool                   `ddl:"static" sql:"DESCRIBE"`
	dbtProject bool                   `ddl:"static" sql:"DBT PROJECT"`
	name       SchemaObjectIdentifier `ddl:"identifier"`
}

type dbtProjectDetailsRow struct {
	Property string `db:"property"`
	Value    string `db:"value"`
}

type DbtProjectDetails struct {
	Property string
	Value    string
}
