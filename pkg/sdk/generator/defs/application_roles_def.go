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

// applicationRolesDef creates an interface that allows for querying application roles.
// It does not allow for other DDL queries (CREATE, ALTER, DROP, ...) to be called, because they are not possible
// to be called from the program level. Application roles are a special case where they're only usable
// inside application context (e.g. setup.sql). Right now, they can be only manipulated from the program context
// by applying debug_mode parameter to the application, but it's a hacky solution and even with that you're limited with GRANT and REVOKE options.
// That's why we're only exposing SHOW operations, because only they are the only allowed operations to be called from the program context.
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
	g.ShowByIDApplicationNameFiltering,
).WithCustomInterfaceMethod(
	"RevokeSafely",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokeApplicationRoleRequest")},
	"error",
)
