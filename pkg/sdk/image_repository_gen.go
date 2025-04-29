package sdk

import (
	"context"
	"time"
)

type ImageRepositories interface {
	Create(ctx context.Context, request *CreateImageRepositoryRequest) error
	Alter(ctx context.Context, request *AlterImageRepositoryRequest) error
	Drop(ctx context.Context, request *DropImageRepositoryRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Show(ctx context.Context, request *ShowImageRepositoryRequest) ([]ImageRepository, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*ImageRepository, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*ImageRepository, error)
}

// CreateImageRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-image-repository.
type CreateImageRepositoryOptions struct {
	create          bool                   `ddl:"static" sql:"CREATE"`
	OrReplace       *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	imageRepository bool                   `ddl:"static" sql:"IMAGE REPOSITORY"`
	IfNotExists     *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
	Comment         *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// AlterImageRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-image-repository.
type AlterImageRepositoryOptions struct {
	alter           bool                   `ddl:"static" sql:"ALTER"`
	imageRepository bool                   `ddl:"static" sql:"IMAGE REPOSITORY"`
	IfExists        *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
	Set             *ImageRepositorySet    `ddl:"keyword" sql:"SET"`
}

type ImageRepositorySet struct {
	Comment *StringAllowEmpty `ddl:"parameter" sql:"COMMENT"`
}

// DropImageRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-image-repository.
type DropImageRepositoryOptions struct {
	drop            bool                   `ddl:"static" sql:"DROP"`
	imageRepository bool                   `ddl:"static" sql:"IMAGE REPOSITORY"`
	IfExists        *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name            SchemaObjectIdentifier `ddl:"identifier"`
}

// ShowImageRepositoryOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-image-repositories.
type ShowImageRepositoryOptions struct {
	show              bool  `ddl:"static" sql:"SHOW"`
	imageRepositories bool  `ddl:"static" sql:"IMAGE REPOSITORIES"`
	Like              *Like `ddl:"keyword" sql:"LIKE"`
	In                *In   `ddl:"keyword" sql:"IN"`
}

type imageRepositoriesRow struct {
	CreatedOn                time.Time `db:"created_on"`
	Name                     string    `db:"name"`
	DatabaseName             string    `db:"database_name"`
	SchemaName               string    `db:"schema_name"`
	RepositoryUrl            string    `db:"repository_url"`
	Owner                    string    `db:"owner"`
	OwnerRoleType            string    `db:"owner_role_type"`
	Comment                  string    `db:"comment"`
	PrivatelinkRepositoryUrl string    `db:"privatelink_repository_url"`
}

type ImageRepository struct {
	CreatedOn                time.Time
	Name                     string
	DatabaseName             string
	SchemaName               string
	RepositoryUrl            string
	Owner                    string
	OwnerRoleType            string
	Comment                  string
	PrivatelinkRepositoryUrl string
}

func (v *ImageRepository) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

func (v *ImageRepository) ObjectType() ObjectType {
	return ObjectTypeImageRepository
}
