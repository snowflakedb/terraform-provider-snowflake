package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

var applicationRoleKindOfRole = g.NewQueryStruct("KindOfRole").
	OptionalIdentifier("RoleName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("ROLE")).
	OptionalIdentifier("ApplicationRoleName", g.KindOfT[DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION ROLE")).
	OptionalIdentifier("ApplicationName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION")).
	WithValidation(g.ExactlyOneValueSet, "RoleName", "ApplicationRoleName", "ApplicationName")

var ApplicationRolesDef = g.NewInterface(
	"ApplicationRoles",
	"ApplicationRole",
	g.KindOfT[DatabaseObjectIdentifier](),
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
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-application-roles",
	g.DbStruct("applicationRoleDbRow").
		Field("created_on", "time.Time").
		Field("name", "string").
		Field("owner", "string").
		Field("comment", "string").
		Field("owner_role_type", "string"),
	g.PlainStruct("ApplicationRole").
		Field("CreatedOn", "time.Time").
		Field("Name", "string").
		Field("Owner", "string").
		Field("Comment", "string").
		Field("OwnerRoleType", "string"),
	g.NewQueryStruct("ShowApplicationRoles").
		Show().
		SQL("APPLICATION ROLES IN APPLICATION").
		Identifier("ApplicationName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions()).
		OptionalLimitFrom().
		WithValidation(g.ValidIdentifier, "ApplicationName"),
).ShowByIdOperationWithFiltering(
	g.ShowByIDApplicationNameFiltering,
)
