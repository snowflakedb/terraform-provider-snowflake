package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

func grantRoleTo() *g.QueryStruct {
	return g.NewQueryStruct("GrantRoleTo").
		OptionalIdentifier("Role", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("ROLE")).
		OptionalIdentifier("User", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("USER")).
		WithValidation(g.ExactlyOneValueSet, "Role", "User")
}

func revokeRoleFrom() *g.QueryStruct {
	return g.NewQueryStruct("RevokeRoleFrom").
		OptionalIdentifier("User", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("USER")).
		OptionalIdentifier("Role", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("ROLE")).
		WithValidation(g.ExactlyOneValueSet, "Role", "User")
}

var rolesDef = g.NewInterface(
	"Roles",
	"Role",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-role",
	g.NewQueryStruct("CreateRole").
		Create().
		OrReplace().
		SQL("ROLE").
		IfNotExists().
		Name().
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-role",
	g.NewQueryStruct("AlterRole").
		Alter().
		SQL("ROLE").
		IfExists().
		Name().
		RenameTo().
		OptionalTextAssignment("SET COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalSetTags().
		OptionalSQLWithCustomFieldName("UnsetComment", "UNSET COMMENT").
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "SetComment", "SetTags", "UnsetComment", "UnsetTags").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-role",
	g.NewQueryStruct("DropRole").
		Drop().
		SQL("ROLE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-roles",
	g.StructPair("roleDBRow", "Role").
		Time("created_on").
		Text("name").
		OptionalBoolFromText("is_default", g.WithRequiredInPlain()).
		OptionalBoolFromText("is_current", g.WithRequiredInPlain()).
		OptionalBoolFromText("is_inherited", g.WithRequiredInPlain()).
		Number("assigned_to_users").
		Number("granted_to_roles").
		Number("granted_roles").
		OptionalText("owner", g.WithRequiredInPlain()).
		OptionalText("comment", g.WithRequiredInPlain()),
	g.NewQueryStruct("ShowRoles").
		Show().
		SQL("ROLES").
		OptionalLike().
		OptionalQueryStructField(
			"InClass",
			g.NewQueryStruct("RolesInClass").
				Identifier("Class", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions()).
				WithValidation(g.ValidIdentifier, "Class"),
			g.KeywordOptions().SQL("IN CLASS"),
		).
		WithAdditionalValidations(),
	g.ShowByIDLikeFiltering,
).CustomOperation(
	"Grant",
	"https://docs.snowflake.com/en/sql-reference/sql/grant-role",
	g.NewQueryStruct("GrantRole").
		Grant().
		SQL("ROLE").
		Name().
		QueryStructField("Grant", grantRoleTo(), g.KeywordOptions().SQL("TO").Required()).
		WithValidation(g.ValidIdentifier, "name").
		WithAdditionalValidations(),
).CustomOperation(
	"Revoke",
	"https://docs.snowflake.com/en/sql-reference/sql/revoke-role",
	g.NewQueryStruct("RevokeRole").
		Revoke().
		SQL("ROLE").
		Name().
		QueryStructField("Revoke", revokeRoleFrom(), g.KeywordOptions().SQL("FROM").Required()).
		WithValidation(g.ValidIdentifier, "name"),
).WithCustomInterfaceMethod(
	"RevokeSafely", "",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokeRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"Use", "",
	[]*g.MethodParameter{g.NewMethodParameter("req", "*UseRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"UseSecondary", "",
	[]*g.MethodParameter{g.NewMethodParameter("req", "*UseSecondaryRolesRequest")},
	"error",
).WithShowByIDFindPredicateKind(g.ShowByIDFindPredicateFullID)
