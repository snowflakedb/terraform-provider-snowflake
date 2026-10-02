package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

func contextFunctionSelect(name, sql string) *g.QueryStruct {
	return g.NewQueryStruct(name).
		SQLWithCustomFieldName("selectExpr", sql)
}

// NameSingular is empty on purpose. A non-empty singular ("ContextFunction") would produce
// CurrentAccountContextFunctionOptions etc.
// Identifier kind is a required placeholder — context functions are not a named Snowflake object.
var contextFunctionsDef = g.NewInterface(
	"ContextFunctions",
	"",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentAccount",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_account",
	g.StructPair("currentAccountDBRow", "CurrentAccount").
		Text("CURRENT_ACCOUNT", g.WithPlainFieldName("Value")),
	contextFunctionSelect("CurrentAccount", "SELECT CURRENT_ACCOUNT() as CURRENT_ACCOUNT"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentOrganizationName",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_organization_name",
	g.StructPair("currentOrganizationNameDBRow", "CurrentOrganizationName").
		Text("CURRENT_ORGANIZATION_NAME", g.WithPlainFieldName("Value")),
	contextFunctionSelect("CurrentOrganizationName", "SELECT CURRENT_ORGANIZATION_NAME() as CURRENT_ORGANIZATION_NAME"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentAccountName",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_account_name",
	g.StructPair("currentAccountNameDBRow", "CurrentAccountName").
		Text("CURRENT_ACCOUNT_NAME", g.WithPlainFieldName("Value")),
	contextFunctionSelect("CurrentAccountName", "SELECT CURRENT_ACCOUNT_NAME() as CURRENT_ACCOUNT_NAME"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentRole",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_role",
	g.StructPair("currentRoleDBRow", "CurrentRole").
		AccountObjectIdentifier("CURRENT_ROLE", g.WithPlainFieldName("Value")),
	contextFunctionSelect("CurrentRole", "SELECT CURRENT_ROLE() as CURRENT_ROLE"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentSecondaryRoles",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_secondary_roles",
	g.StructPair("currentSecondaryRolesDBRow", "CurrentSecondaryRoles").
		Field("CURRENT_ROLES", "string", "[]AccountObjectIdentifier", g.WithPlainFieldName("Roles")).
		PlainOnlyField("Value", "SecondaryRoleOption"),
	contextFunctionSelect("CurrentSecondaryRoles", "SELECT CURRENT_SECONDARY_ROLES() as CURRENT_ROLES"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentRegion",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_region",
	g.StructPair("currentRegionDBRow", "CurrentRegion").
		Text("CURRENT_REGION", g.WithPlainFieldName("Value")),
	contextFunctionSelect("CurrentRegion", "SELECT CURRENT_REGION() AS CURRENT_REGION"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentSession",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_session",
	g.StructPair("currentSessionDBRow", "CurrentSession").
		Text("CURRENT_SESSION", g.WithPlainFieldName("Value")),
	contextFunctionSelect("CurrentSession", "SELECT CURRENT_SESSION() as CURRENT_SESSION"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentUser",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_user",
	g.StructPair("currentUserDBRow", "CurrentUser").
		AccountObjectIdentifier("CURRENT_USER", g.WithPlainFieldName("Value")),
	contextFunctionSelect("CurrentUser", "SELECT CURRENT_USER() as CURRENT_USER"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentHost",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/system_allowlist",
	g.StructPair("currentHostDBRow", "CurrentHost").
		Text("HOST", g.WithPlainFieldName("Value"), g.WithValueAdjuster("prefixHTTPS")),
	contextFunctionSelect("CurrentHost", "SELECT VALUE:host::VARCHAR AS HOST FROM TABLE(FLATTEN(INPUT => PARSE_JSON(SYSTEM$ALLOWLIST()))) WHERE VALUE:type::VARCHAR = 'SNOWFLAKE_DEPLOYMENT'"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentSessionDetails",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions-context",
	g.StructPair("currentSessionDetailsDBRow", "CurrentSessionDetails").
		Text("CURRENT_ACCOUNT", g.WithPlainFieldName("Account")).
		Text("CURRENT_ACCOUNT_NAME", g.WithPlainFieldName("AccountName")).
		Text("CURRENT_ORGANIZATION_NAME", g.WithPlainFieldName("OrganizationName")).
		Text("CURRENT_ROLE", g.WithPlainFieldName("Role")).
		Text("CURRENT_REGION", g.WithPlainFieldName("Region")).
		Text("CURRENT_SESSION", g.WithPlainFieldName("Session")).
		Text("CURRENT_USER", g.WithPlainFieldName("User")),
	contextFunctionSelect("CurrentSessionDetails", "SELECT CURRENT_ACCOUNT() as CURRENT_ACCOUNT, CURRENT_ROLE() as CURRENT_ROLE, CURRENT_REGION() AS CURRENT_REGION, CURRENT_SESSION() as CURRENT_SESSION, CURRENT_USER() as CURRENT_USER, CURRENT_ACCOUNT_NAME() as CURRENT_ACCOUNT_NAME, CURRENT_ORGANIZATION_NAME() as CURRENT_ORGANIZATION_NAME"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	// TODO(SNOW-1805152): Remove LastQueryId and utilize gosnowflake.WithQueryIDChan instead whenever query id is needed
	"LastQueryId",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/last_query_id",
	g.StructPair("lastQueryIdDBRow", "LastQueryId").
		OptionalText("LAST_QUERY_ID", g.WithPlainFieldName("Value"), g.WithRequiredInPlain()),
	contextFunctionSelect("LastQueryId", "SELECT LAST_QUERY_ID() as LAST_QUERY_ID"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentDatabase",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_database",
	g.StructPair("currentDatabaseDBRow", "CurrentDatabase").
		OptionalText("CURRENT_DATABASE", g.WithPlainFieldName("Value"), g.WithRequiredInPlain()),
	contextFunctionSelect("CurrentDatabase", "SELECT CURRENT_DATABASE() as CURRENT_DATABASE"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentSchema",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_schema",
	g.StructPair("currentSchemaDBRow", "CurrentSchema").
		OptionalText("CURRENT_SCHEMA", g.WithPlainFieldName("Value"), g.WithRequiredInPlain()),
	contextFunctionSelect("CurrentSchema", "SELECT CURRENT_SCHEMA() as CURRENT_SCHEMA"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructsAndOpts(
	"CurrentWarehouse",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/current_warehouse",
	g.StructPair("currentWarehouseDBRow", "CurrentWarehouse").
		OptionalText("CURRENT_WAREHOUSE", g.WithPlainFieldName("Value"), g.WithRequiredInPlain()),
	contextFunctionSelect("CurrentWarehouse", "SELECT CURRENT_WAREHOUSE() as CURRENT_WAREHOUSE"),
	[]g.CustomOperationOption{g.WithNoRequest()},
).CustomShowOperationWithPairedStructs(
	"IsRoleInSession",
	g.ShowMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/functions/is_role_in_session",
	g.StructPair("isRoleInSessionDBRow", "IsRoleInSession").
		Bool("IS_ROLE_IN_SESSION", g.WithPlainFieldName("Value")),
	g.NewQueryStruct("IsRoleInSession").
		SQLWithCustomFieldName("selectIsRoleInSession", "SELECT IS_ROLE_IN_SESSION").
		QueryStructField(
			"Arguments",
			g.NewQueryStruct("IsRoleInSessionArguments").
				Identifier("Role", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SingleQuotes().Required()).
				WithValidation(g.ValidIdentifier, "Role"),
			g.ListOptions().MustParentheses().Required(),
		).
		SQL("AS IS_ROLE_IN_SESSION"),
)
