package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Tags interface {
	Create(ctx context.Context, request *CreateTagRequest) error
	Alter(ctx context.Context, request *AlterTagRequest) error
	Show(ctx context.Context, opts *ShowTagRequest) ([]Tag, error)
	ShowByID(ctx context.Context, id SchemaObjectIdentifier) (*Tag, error)
	ShowByIDSafely(ctx context.Context, id SchemaObjectIdentifier) (*Tag, error)
	Drop(ctx context.Context, request *DropTagRequest) error
	DropSafely(ctx context.Context, id SchemaObjectIdentifier) error
	Undrop(ctx context.Context, request *UndropTagRequest) error
	Set(ctx context.Context, request *SetTagRequest) error
	Unset(ctx context.Context, request *UnsetTagRequest) error
	SetOnCurrentAccount(ctx context.Context, request *SetTagOnCurrentAccountRequest) error
	UnsetOnCurrentAccount(ctx context.Context, request *UnsetTagOnCurrentAccountRequest) error
}

type setTagOptions struct {
	alter      bool             `ddl:"static" sql:"ALTER"`
	objectType ObjectType       `ddl:"keyword"`
	objectName ObjectIdentifier `ddl:"identifier"`
	column     *string          `ddl:"parameter,no_equals,double_quotes" sql:"MODIFY COLUMN"`
	SetTags    []TagAssociation `ddl:"keyword" sql:"SET TAG"`
}

type unsetTagOptions struct {
	alter      bool               `ddl:"static" sql:"ALTER"`
	objectType ObjectType         `ddl:"keyword"`
	IfExists   *bool              `ddl:"keyword" sql:"IF EXISTS"`
	objectName ObjectIdentifier   `ddl:"identifier"`
	column     *string            `ddl:"parameter,no_equals,double_quotes" sql:"MODIFY COLUMN"`
	UnsetTags  []ObjectIdentifier `ddl:"keyword" sql:"UNSET TAG"`
}

// createTagOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-tag
type createTagOptions struct {
	create        bool                   `ddl:"static" sql:"CREATE"`
	OrReplace     *bool                  `ddl:"keyword" sql:"OR REPLACE"`
	tag           string                 `ddl:"static" sql:"TAG"`
	IfNotExists   *bool                  `ddl:"keyword" sql:"IF NOT EXISTS"`
	name          SchemaObjectIdentifier `ddl:"identifier"`
	AllowedValues *AllowedValues         `ddl:"keyword" sql:"ALLOWED_VALUES"`
	Comment       *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type AllowedValues struct {
	Values []AllowedValue `ddl:"list,comma"`
}

type AllowedValue struct {
	Value string `ddl:"keyword,single_quotes"`
}

// showTagOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-tags
type showTagOptions struct {
	show bool        `ddl:"static" sql:"SHOW"`
	tag  bool        `ddl:"static" sql:"TAGS"`
	Like *Like       `ddl:"keyword" sql:"LIKE"`
	In   *ExtendedIn `ddl:"keyword" sql:"IN"`
}

type Tag struct {
	CreatedOn     time.Time
	Name          string
	DatabaseName  string
	SchemaName    string
	Owner         string
	Comment       string
	AllowedValues []string
	OwnerRoleType string
}

func (v *Tag) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
}

type tagRow struct {
	CreatedOn     time.Time      `db:"created_on"`
	Name          string         `db:"name"`
	DatabaseName  string         `db:"database_name"`
	SchemaName    string         `db:"schema_name"`
	Owner         string         `db:"owner"`
	Comment       string         `db:"comment"`
	AllowedValues sql.NullString `db:"allowed_values"`
	OwnerRoleType string         `db:"owner_role_type"`
}

func (tr tagRow) convertErr() (*Tag, error) {
	t := &Tag{
		CreatedOn:     tr.CreatedOn,
		Name:          tr.Name,
		DatabaseName:  tr.DatabaseName,
		SchemaName:    tr.SchemaName,
		Owner:         tr.Owner,
		Comment:       tr.Comment,
		OwnerRoleType: tr.OwnerRoleType,
	}
	if tr.AllowedValues.Valid {
		t.AllowedValues = ParseCommaSeparatedStringArray(tr.AllowedValues.String, true)
	}
	return t, nil
}

type TagSetMaskingPolicies struct {
	MaskingPolicies []TagMaskingPolicy `ddl:"list,comma"`
	Force           *bool              `ddl:"keyword" sql:"FORCE"`
}

type TagUnsetMaskingPolicies struct {
	MaskingPolicies []TagMaskingPolicy `ddl:"list,comma,no_quotes"`
}

type TagMaskingPolicy struct {
	Name SchemaObjectIdentifier `ddl:"identifier" sql:"MASKING POLICY"`
}

type TagSet struct {
	MaskingPolicies *TagSetMaskingPolicies `ddl:"keyword"`
	Comment         *string                `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type TagUnset struct {
	MaskingPolicies *TagUnsetMaskingPolicies `ddl:"keyword"`
	AllowedValues   *bool                    `ddl:"keyword" sql:"ALLOWED_VALUES"`
	Comment         *bool                    `ddl:"keyword" sql:"COMMENT"`
}

type TagAdd struct {
	AllowedValues *AllowedValues `ddl:"keyword" sql:"ALLOWED_VALUES"`
}

type TagDrop struct {
	AllowedValues *AllowedValues `ddl:"keyword" sql:"ALLOWED_VALUES"`
}

type TagRename struct {
	Name SchemaObjectIdentifier `ddl:"identifier"`
}

// alterTagOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-tag
type alterTagOptions struct {
	alter    bool                   `ddl:"static" sql:"ALTER"`
	tag      string                 `ddl:"static" sql:"TAG"`
	ifExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`

	// One of
	Add    *TagAdd    `ddl:"keyword" sql:"ADD"`
	Drop   *TagDrop   `ddl:"keyword" sql:"DROP"`
	Set    *TagSet    `ddl:"keyword" sql:"SET"`
	Unset  *TagUnset  `ddl:"keyword" sql:"UNSET"`
	Rename *TagRename `ddl:"keyword" sql:"RENAME TO"`
}

// dropTagOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-tag
type dropTagOptions struct {
	drop     bool                   `ddl:"static" sql:"DROP"`
	tag      string                 `ddl:"static" sql:"TAG"`
	IfExists *bool                  `ddl:"keyword" sql:"IF EXISTS"`
	name     SchemaObjectIdentifier `ddl:"identifier"`
}

// undropTagOptions is based on https://docs.snowflake.com/en/sql-reference/sql/undrop-tag
type undropTagOptions struct {
	undrop bool                   `ddl:"static" sql:"UNDROP"`
	tag    string                 `ddl:"static" sql:"TAG"`
	name   SchemaObjectIdentifier `ddl:"identifier"`
}
