package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var secondaryRoleOptionEnum = g.NewEnum("SecondaryRoleOption", "SecondaryRoleOptions", "ALL", "NONE")

func sessionSet() *g.QueryStruct {
	return g.NewQueryStruct("SessionSet").
		PredefinedQueryStructField("SessionParameters", "*SessionParameters", g.ListOptions()).
		WithValidation(g.ValidateValue, "SessionParameters")
}

func sessionUnset() *g.QueryStruct {
	return g.NewQueryStruct("SessionUnset").
		PredefinedQueryStructField("SessionParametersUnset", "*SessionParametersUnset", g.ListOptions()).
		WithValidation(g.ValidateValue, "SessionParametersUnset")
}

var sessionsDef = g.NewInterface(
	"Sessions",
	"Session",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](), // Sessions has no named object; required placeholder.
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-session",
	g.NewQueryStruct("AlterSession").
		Alter().
		SQL("SESSION").
		// No Name() — ALTER SESSION acts on the current session, not a named object.
		OptionalQueryStructField("Set", sessionSet(), g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("Unset", sessionUnset(), g.KeywordOptions().SQL("UNSET")).
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset"),
).CustomOperation(
	"UseWarehouse",
	"https://docs.snowflake.com/en/sql-reference/sql/use-warehouse",
	g.NewQueryStruct("UseWarehouse").
		SQL("USE WAREHOUSE").
		Identifier("name", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"UseDatabase",
	"https://docs.snowflake.com/en/sql-reference/sql/use-database",
	g.NewQueryStruct("UseDatabase").
		SQL("USE DATABASE").
		Identifier("name", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"UseSchema",
	"https://docs.snowflake.com/en/sql-reference/sql/use-schema",
	g.NewQueryStruct("UseSchema").
		SQL("USE SCHEMA").
		Identifier("name", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().Required()).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"UseRole",
	"https://docs.snowflake.com/en/sql-reference/sql/use-role",
	g.NewQueryStruct("UseRole").
		SQL("USE ROLE").
		Identifier("name", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"UseSecondaryRoles",
	"https://docs.snowflake.com/en/sql-reference/sql/use-secondary-roles",
	g.NewQueryStruct("UseSecondaryRoles").
		SQL("USE SECONDARY ROLES").
		PredefinedQueryStructField("Option", "SecondaryRoleOption", g.KeywordOptions().Required()),
).WithCustomInterfaceMethod(
	"ShowParameters",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("opts", "*ShowParametersOptions")},
	"[]*Parameter", "error",
).WithEnums(secondaryRoleOptionEnum)
