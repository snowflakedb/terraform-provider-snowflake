package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var applicationRoleKindOfRole = g.NewQueryStruct("KindOfRole").
	OptionalIdentifier("RoleName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("ROLE")).
	OptionalIdentifier("ApplicationRoleName", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION ROLE")).
	OptionalIdentifier("ApplicationName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION")).
	WithValidation(g.ExactlyOneValueSet, "RoleName", "ApplicationRoleName", "ApplicationName")

var applicationRolesDef = g.NewInterface(
	"ApplicationRoles",
	"ApplicationRole",
	g.KindOfT[sdkcommons.DatabaseObjectIdentifier](),
).CustomOperation(
	"Grant",
	"https://docs.snowflake.com/en/sql-reference/sql/grant-application-role",
	g.NewQueryStruct("GrantApplicationRole").
		Grant().
		SQL("APPLICATION ROLE").
		Name().
		QueryStructField(
			"To",
			applicationRoleKindOfRole,
			g.KeywordOptions().SQL("TO"),
		).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"Revoke",
	"https://docs.snowflake.com/en/sql-reference/sql/revoke-application-role",
	g.NewQueryStruct("RevokeApplicationRole").
		Revoke().
		SQL("APPLICATION ROLE").
		Name().
		QueryStructField(
			"From",
			applicationRoleKindOfRole,
			g.KeywordOptions().SQL("FROM"),
		).
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-application-roles",
	g.StructPair("applicationRoleDbRow", "ApplicationRole").
		Time("created_on").
		Text("name").
		Text("owner").
		Text("comment").
		Text("owner_role_type"),
	g.NewQueryStruct("ShowApplicationRoles").
		Show().
		SQL("APPLICATION ROLES IN APPLICATION").
		Identifier("ApplicationName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions()).
		OptionalLimitFrom().
		WithValidation(g.ValidIdentifier, "ApplicationName"),
).ShowByIdOperationWithFiltering(
	g.ShowByIDApplicationNameFiltering,
).
	WithCustomInterfaceMethod(
		"RevokeSafely",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokeApplicationRoleRequest")},
		"error",
	)
