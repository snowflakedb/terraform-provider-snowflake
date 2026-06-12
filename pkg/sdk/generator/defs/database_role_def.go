//go:build sdk_generation

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

func newDatabaseRoleKindOfRole() *g.QueryStruct {
	return g.NewQueryStruct("DatabaseRoleKindOfRole").
		OptionalIdentifier("DatabaseRoleName", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("DATABASE ROLE")).
		OptionalIdentifier("AccountRoleName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("ROLE")).
		WithValidation(g.ExactlyOneValueSet, "DatabaseRoleName", "AccountRoleName")
}

var databaseRolesDef = g.NewInterface(
	"DatabaseRoles",
	"DatabaseRole",
	g.KindOfT[sdkcommons.DatabaseObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-database-role",
	g.NewQueryStruct("CreateDatabaseRole").
		Create().
		OrReplace().
		SQL("DATABASE ROLE").
		IfNotExists().
		Name().
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-database-role",
	g.NewQueryStruct("AlterDatabaseRole").
		Alter().
		SQL("DATABASE ROLE").
		IfExists().
		Name().
		OptionalIdentifier("Rename", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField("Set", g.NewQueryStruct("DatabaseRoleSet").
			OptionalComment().
			WithValidation(g.AtLeastOneValueSet, "Comment"),
			g.ListOptions().NoParentheses().SQL("SET"),
		).
		OptionalQueryStructField("Unset", g.NewQueryStruct("DatabaseRoleUnset").
			OptionalSQL("COMMENT").
			WithValidation(g.AtLeastOneValueSet, "Comment"),
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Rename", "Set", "Unset", "SetTags", "UnsetTags").
		WithValidation(g.ValidIdentifierIfSet, "Rename").
		WithValidation(g.AdditionalValidations),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-database-role",
	g.NewQueryStruct("DropDatabaseRole").
		Drop().
		SQL("DATABASE ROLE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-database-roles",
	g.StructPair("databaseRoleDBRow", "DatabaseRole").
		Text("created_on").
		Text("name").
		OptionalBoolFromText("is_default", g.WithRequiredInPlain()).
		OptionalBoolFromText("is_current", g.WithRequiredInPlain()).
		OptionalBoolFromText("is_inherited", g.WithRequiredInPlain()).
		Number("granted_to_roles").
		Number("granted_to_database_roles").
		Number("granted_database_roles").
		Text("owner").
		OptionalText("comment", g.WithRequiredInPlain()).
		OptionalText("owner_role_type", g.WithRequiredInPlain()).
		PlainOnlyField("DatabaseName", "string"),
	g.NewQueryStruct("ShowDatabaseRoles").
		Show().
		SQL("DATABASE ROLES").
		OptionalLike().
		SQL("IN DATABASE").
		Identifier("Database", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions()).
		OptionalLimitFrom().
		WithValidation(g.ValidIdentifier, "Database"),
	g.ShowByIDSuppressed,
).CustomOperation(
	"Grant",
	"https://docs.snowflake.com/en/sql-reference/sql/grant-database-role",
	g.NewQueryStruct("GrantDatabaseRole").
		Grant().
		SQL("DATABASE ROLE").
		Name().
		QueryStructField("To", newDatabaseRoleKindOfRole(), g.KeywordOptions().SQL("TO")).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"Revoke",
	"https://docs.snowflake.com/en/sql-reference/sql/revoke-database-role",
	g.NewQueryStruct("RevokeDatabaseRole").
		Revoke().
		SQL("DATABASE ROLE").
		Name().
		QueryStructField("From", newDatabaseRoleKindOfRole(), g.KeywordOptions().SQL("FROM")).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"GrantToShare",
	"https://docs.snowflake.com/en/sql-reference/sql/grant-database-role-share",
	g.NewQueryStruct("GrantDatabaseRoleToShare").
		Grant().
		SQL("DATABASE ROLE").
		Name().
		SQL("TO SHARE").
		Identifier("Share", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "Share"),
).CustomOperation(
	"RevokeFromShare",
	"https://docs.snowflake.com/en/sql-reference/sql/grant-database-role-share",
	g.NewQueryStruct("RevokeDatabaseRoleFromShare").
		Revoke().
		SQL("DATABASE ROLE").
		Name().
		SQL("FROM SHARE").
		Identifier("Share", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "Share"),
).WithCustomInterfaceMethod(
	"ShowByID", "",
	[]*g.MethodParameter{g.NewMethodParameter("id", "DatabaseObjectIdentifier")},
	"*DatabaseRole", "error",
).WithCustomInterfaceMethod(
	"ShowByIDSafely", "",
	[]*g.MethodParameter{g.NewMethodParameter("id", "DatabaseObjectIdentifier")},
	"*DatabaseRole", "error",
).WithCustomInterfaceMethod(
	"RevokeSafely", "",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokeDatabaseRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"RevokeFromShareSafely", "",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokeFromShareDatabaseRoleRequest")},
	"error",
)
