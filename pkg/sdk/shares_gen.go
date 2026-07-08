package sdk

import (
	"context"
	"time"
)

type Shares interface {
	Create(ctx context.Context, id AccountObjectIdentifier, opts *CreateShareOptions) error
	Alter(ctx context.Context, id AccountObjectIdentifier, opts *AlterShareOptions) error
	Drop(ctx context.Context, id AccountObjectIdentifier, opts *DropShareOptions) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, opts *ShowShareOptions) ([]Share, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Share, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Share, error)
	Describe(ctx context.Context, id AccountObjectIdentifier) ([]ShareInfo, error)
	DescribeProvider(ctx context.Context, id AccountObjectIdentifier) ([]ShareInfo, error)
	DescribeConsumer(ctx context.Context, id ExternalObjectIdentifier) ([]ShareInfo, error)
}

type Share struct {
	CreatedOn    time.Time
	Kind         ShareKind
	Name         ExternalObjectIdentifier
	DatabaseName AccountObjectIdentifier
	To           []AccountIdentifier
	Owner        string
	Comment      string
}

func (v *Share) ID() AccountObjectIdentifier {
	return v.Name.objectIdentifier.(AccountObjectIdentifier)
}

func (v *Share) ExternalID() ExternalObjectIdentifier {
	return v.Name
}

func (v *Share) ObjectType() ObjectType {
	return ObjectTypeShare
}

type shareRow struct {
	CreatedOn    time.Time `db:"created_on"`
	Kind         string    `db:"kind"`
	OwnerAccount string    `db:"owner_account"`
	Name         string    `db:"name"`
	DatabaseName string    `db:"database_name"`
	To           string    `db:"to"`
	Owner        string    `db:"owner"`
	Comment      string    `db:"comment"`
}

// CreateShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/create-share.
type CreateShareOptions struct {
	create    bool                    `ddl:"static" sql:"CREATE"`
	OrReplace *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	share     bool                    `ddl:"static" sql:"SHARE"`
	name      AccountObjectIdentifier `ddl:"identifier"`
	Comment   *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

// DropShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/drop-share.
type DropShareOptions struct {
	drop     bool                    `ddl:"static" sql:"DROP"`
	share    bool                    `ddl:"static" sql:"SHARE"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

// AlterShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/alter-share.
type AlterShareOptions struct {
	alter    bool                    `ddl:"static" sql:"ALTER"`
	share    bool                    `ddl:"static" sql:"SHARE"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
	Add      *ShareAdd               `ddl:"keyword" sql:"ADD"`
	Remove   *ShareRemove            `ddl:"keyword" sql:"REMOVE"`
	Set      *ShareSet               `ddl:"keyword" sql:"SET"`
	Unset    *ShareUnset             `ddl:"keyword" sql:"UNSET"`
	SetTag   []TagAssociation        `ddl:"keyword" sql:"SET TAG"`
	UnsetTag []ObjectIdentifier      `ddl:"keyword" sql:"UNSET TAG"`
}

type ShareAdd struct {
	Accounts          []AccountIdentifier `ddl:"parameter" sql:"ACCOUNTS"`
	ShareRestrictions *bool               `ddl:"parameter" sql:"SHARE_RESTRICTIONS"`
}

type ShareRemove struct {
	Accounts []AccountIdentifier `ddl:"parameter" sql:"ACCOUNTS"`
}

type ShareSet struct {
	Accounts []AccountIdentifier `ddl:"parameter" sql:"ACCOUNTS"`
	Comment  *string             `ddl:"parameter,single_quotes" sql:"COMMENT"`
}

type ShareUnset struct {
	Comment *bool `ddl:"keyword" sql:"COMMENT"`
}

// ShowShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/show-shares.
type ShowShareOptions struct {
	show       bool       `ddl:"static" sql:"SHOW"`
	shares     bool       `ddl:"static" sql:"SHARES"`
	Like       *Like      `ddl:"keyword" sql:"LIKE"`
	StartsWith *string    `ddl:"parameter,single_quotes,no_equals" sql:"STARTS WITH"`
	Limit      *LimitFrom `ddl:"keyword" sql:"LIMIT"`
}

// describeShareOptions is based on https://docs.snowflake.com/en/sql-reference/sql/desc-share.
type describeShareOptions struct {
	describe bool             `ddl:"static" sql:"DESCRIBE"`
	share    bool             `ddl:"static" sql:"SHARE"`
	name     ObjectIdentifier `ddl:"identifier"`
}

type ShareInfo struct {
	Kind     ObjectType
	Name     ObjectIdentifier
	SharedOn time.Time
}

type shareDetailsRow struct {
	Kind     string    `db:"kind"`
	Name     string    `db:"name"`
	SharedOn time.Time `db:"shared_on"`
}
