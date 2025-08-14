package sdk

import (
	"context"
	"database/sql"
	"time"
)

type Roles interface {
	Create(ctx context.Context, req *CreateRoleRequest) error
	Alter(ctx context.Context, req *AlterRoleRequest) error
	Drop(ctx context.Context, req *DropRoleRequest) error
	DropSafely(ctx context.Context, id AccountObjectIdentifier) error
	Show(ctx context.Context, req *ShowRoleRequest) ([]Role, error)
	ShowByID(ctx context.Context, id AccountObjectIdentifier) (*Role, error)
	ShowByIDSafely(ctx context.Context, id AccountObjectIdentifier) (*Role, error)
	Grant(ctx context.Context, req *GrantRoleRequest) error
	Revoke(ctx context.Context, req *RevokeRoleRequest) error
	Use(ctx context.Context, req *UseRoleRequest) error
	UseSecondary(ctx context.Context, req *UseSecondaryRolesRequest) error
}

type Role struct {
	CreatedOn       time.Time
	Name            string
	IsDefault       bool
	IsCurrent       bool
	IsInherited     bool
	AssignedToUsers int
	GrantedToRoles  int
	GrantedRoles    int
	Owner           string
	Comment         string
}

func (v *Role) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

func (v *Role) ObjectType() ObjectType {
	return ObjectTypeRole
}

type roleDBRow struct {
	CreatedOn       time.Time      `db:"created_on"`
	Name            string         `db:"name"`
	IsDefault       sql.NullString `db:"is_default"`
	IsCurrent       sql.NullString `db:"is_current"`
	IsInherited     sql.NullString `db:"is_inherited"`
	AssignedToUsers int            `db:"assigned_to_users"`
	GrantedToRoles  int            `db:"granted_to_roles"`
	GrantedRoles    int            `db:"granted_roles"`
	Owner           sql.NullString `db:"owner"`
	Comment         sql.NullString `db:"comment"`
}

func (row roleDBRow) convert() (*Role, error) {
	role := &Role{
		CreatedOn:       row.CreatedOn,
		Name:            row.Name,
		AssignedToUsers: row.AssignedToUsers,
		GrantedToRoles:  row.GrantedToRoles,
		GrantedRoles:    row.GrantedRoles,
	}
	if row.IsDefault.Valid {
		role.IsDefault = row.IsDefault.String == "Y"
	}
	if row.IsCurrent.Valid {
		role.IsCurrent = row.IsCurrent.String == "Y"
	}
	if row.IsInherited.Valid {
		role.IsInherited = row.IsInherited.String == "Y"
	}
	if row.Owner.Valid {
		role.Owner = row.Owner.String
	}
	if row.Comment.Valid {
		role.Comment = row.Comment.String
	}
	return role, nil
}

// CreateRoleOptions based on https://docs.snowflake.com/en/sql-reference/sql/create-role
type CreateRoleOptions struct {
	create      bool                    `ddl:"static" sql:"CREATE"`
	OrReplace   *bool                   `ddl:"keyword" sql:"OR REPLACE"`
	role        bool                    `ddl:"static" sql:"ROLE"`
	IfNotExists *bool                   `ddl:"keyword" sql:"IF NOT EXISTS"`
	name        AccountObjectIdentifier `ddl:"identifier"`
	Comment     *string                 `ddl:"parameter,single_quotes" sql:"COMMENT"`
	Tag         []TagAssociation        `ddl:"keyword,parentheses" sql:"TAG"`
}

// AlterRoleOptions based on https://docs.snowflake.com/en/sql-reference/sql/alter-role
type AlterRoleOptions struct {
	alter    bool                    `ddl:"static" sql:"ALTER"`
	role     bool                    `ddl:"static" sql:"ROLE"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`

	// One of
	RenameTo     *AccountObjectIdentifier `ddl:"identifier" sql:"RENAME TO"`
	SetComment   *string                  `ddl:"parameter,single_quotes" sql:"SET COMMENT"`
	SetTags      []TagAssociation         `ddl:"keyword" sql:"SET TAG"`
	UnsetComment *bool                    `ddl:"keyword" sql:"UNSET COMMENT"`
	UnsetTags    []ObjectIdentifier       `ddl:"keyword" sql:"UNSET TAG"`
}

// DropRoleOptions based on https://docs.snowflake.com/en/sql-reference/sql/drop-role
type DropRoleOptions struct {
	drop     bool                    `ddl:"static" sql:"DROP"`
	roles    bool                    `ddl:"static" sql:"ROLE"`
	IfExists *bool                   `ddl:"keyword" sql:"IF EXISTS"`
	name     AccountObjectIdentifier `ddl:"identifier"`
}

// ShowRoleOptions based on https://docs.snowflake.com/en/sql-reference/sql/show-roles
type ShowRoleOptions struct {
	show    bool          `ddl:"static" sql:"SHOW"`
	roles   bool          `ddl:"static" sql:"ROLES"`
	Like    *Like         `ddl:"keyword" sql:"LIKE"`
	InClass *RolesInClass `ddl:"keyword" sql:"IN CLASS"`
}

type RolesInClass struct {
	Class ObjectIdentifier `ddl:"identifier"`
}

// GrantRoleOptions based on https://docs.snowflake.com/en/sql-reference/sql/grant-role
type GrantRoleOptions struct {
	grant bool                    `ddl:"static" sql:"GRANT"`
	role  bool                    `ddl:"static" sql:"ROLE"`
	name  AccountObjectIdentifier `ddl:"identifier"`
	Grant GrantRole               `ddl:"keyword,no_parentheses" sql:"TO"`
}

type GrantRole struct {
	// one of
	Role *AccountObjectIdentifier `ddl:"identifier" sql:"ROLE"`
	User *AccountObjectIdentifier `ddl:"identifier" sql:"USER"`
}

// RevokeRoleOptions based on https://docs.snowflake.com/en/sql-reference/sql/revoke-role
type RevokeRoleOptions struct {
	revoke bool                    `ddl:"static" sql:"REVOKE"`
	role   bool                    `ddl:"static" sql:"ROLE"`
	name   AccountObjectIdentifier `ddl:"identifier"`
	Revoke RevokeRole              `ddl:"keyword,no_parentheses" sql:"FROM"`
}

type RevokeRole struct {
	// one of
	User *AccountObjectIdentifier `ddl:"identifier" sql:"USER"`
	Role *AccountObjectIdentifier `ddl:"identifier" sql:"ROLE"`
}
