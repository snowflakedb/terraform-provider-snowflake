package defs

import "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/ssot"

// DatabaseRoleAttrs holds the database role fields that CREATE, ALTER, and SHOW reuse.
//
// The SHOW columns keep the types of the handwritten SDK definition (pkg/sdk/generator/defs/database_role_def.go), so
// the migration does not change sdk.DatabaseRole or the resource show_output: created_on stays a string, and is_default,
// is_current, and is_inherited are nullable 'Y'/'N' text mapped to a plain bool.
var DatabaseRoleAttrs = struct {
	Comment                ssot.Attribute
	CreatedOn              ssot.Attribute
	Name                   ssot.Attribute
	IsDefault              ssot.Attribute
	IsCurrent              ssot.Attribute
	IsInherited            ssot.Attribute
	GrantedToRoles         ssot.Attribute
	GrantedToDatabaseRoles ssot.Attribute
	GrantedDatabaseRoles   ssot.Attribute
	Owner                  ssot.Attribute
	OwnerRoleType          ssot.Attribute
}{
	Comment: ssot.Attribute{
		Name:         "comment",
		Type:         ssot.AttrString,
		Description:  "Specifies a comment for the database role.",
		OptionalInDb: true,
	},
	CreatedOn:              ssot.Attribute{Name: "created_on", Type: ssot.AttrString, Readonly: true},
	Name:                   ssot.Attribute{Name: "name", Type: ssot.AttrString, Readonly: true},
	IsDefault:              ssot.Attribute{Name: "is_default", Type: ssot.AttrBool, StoredAsText: true, OptionalInDb: true, BoolTrueValue: "Y", Readonly: true},
	IsCurrent:              ssot.Attribute{Name: "is_current", Type: ssot.AttrBool, StoredAsText: true, OptionalInDb: true, BoolTrueValue: "Y", Readonly: true},
	IsInherited:            ssot.Attribute{Name: "is_inherited", Type: ssot.AttrBool, StoredAsText: true, OptionalInDb: true, BoolTrueValue: "Y", Readonly: true},
	GrantedToRoles:         ssot.Attribute{Name: "granted_to_roles", Type: ssot.AttrInt, Readonly: true},
	GrantedToDatabaseRoles: ssot.Attribute{Name: "granted_to_database_roles", Type: ssot.AttrInt, Readonly: true},
	GrantedDatabaseRoles:   ssot.Attribute{Name: "granted_database_roles", Type: ssot.AttrInt, Readonly: true},
	Owner:                  ssot.Attribute{Name: "owner", Type: ssot.AttrString, Readonly: true},
	OwnerRoleType:          ssot.Attribute{Name: "owner_role_type", Type: ssot.AttrString, OptionalInDb: true, Readonly: true},
}

// DatabaseRole is migrated from the handwritten SDK definition one operation at a time, so SdkDefineFromSsot() stays off
// until every operation, custom operation, and custom interface method has moved here.
//
// Checked against create-database-role, alter-database-role, drop-database-role, and show-database-roles.
//
// CREATE OR ALTER DATABASE ROLE is documented. It is never added, because our resources do not use CREATE OR ALTER.
//
// SHOW DATABASE ROLES documents no column list. The columns are the ones of the handwritten SDK definition, kept as they
// are so the stable resource does not change. An uncommitted probe confirmed the same 11 columns in the same order.
// created_on is TIMESTAMP_LTZ in Snowflake and stays a string, as in the handwritten definition. A role without a
// comment returns an empty comment, not NULL.
//
// The page also does not show LIKE, which the SDK and ShowByID use. The probe showed that LIKE before IN DATABASE and
// LIMIT ... FROM are applied, not only accepted: each filtered a database with two roles to the expected row, and a LIKE
// that matches nothing returned no rows. STARTS WITH and SHOW TERSE are rejected, so they are not added.
//
// The database is not a SHOW column. The SDK fills DatabaseName from the request.
//
// Object parameters: none is documented for a database role.
var DatabaseRole = ssot.Object{
	Name:       "DatabaseRole",
	NamePlural: "DatabaseRoles",
	IdType:     "DatabaseObjectIdentifier",
	Attributes: []ssot.Attribute{
		DatabaseRoleAttrs.Comment,
	},
	Create: ssot.CreateSpec{
		Clauses: []*ssot.Attribute{
			&DatabaseRoleAttrs.Comment,
		},
	},
	Alter: ssot.AlterSpec{Clauses: []ssot.AlterClause{
		ssot.AlterRenameTo(),
		ssot.AlterSet(
			&DatabaseRoleAttrs.Comment,
		),
		ssot.AlterUnset(
			&DatabaseRoleAttrs.Comment,
		),
		ssot.AlterSetTags(),
		ssot.AlterUnsetTags(),
	}},
	Show: ssot.OutputSpec{
		Columns: []*ssot.Attribute{
			&DatabaseRoleAttrs.CreatedOn,
			&DatabaseRoleAttrs.Name,
			&DatabaseRoleAttrs.IsDefault,
			&DatabaseRoleAttrs.IsCurrent,
			&DatabaseRoleAttrs.IsInherited,
			&DatabaseRoleAttrs.GrantedToRoles,
			&DatabaseRoleAttrs.GrantedToDatabaseRoles,
			&DatabaseRoleAttrs.GrantedDatabaseRoles,
			&DatabaseRoleAttrs.Owner,
			&DatabaseRoleAttrs.Comment,
			&DatabaseRoleAttrs.OwnerRoleType,
		},
	},
}
