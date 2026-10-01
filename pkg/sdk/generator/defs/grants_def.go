package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var grantInheritedFromEnumDef = g.NewEnum("GrantInheritedFrom", "GrantInheritedFroms", "ACCOUNT", "DATABASE", "SCHEMA")

func accountRoleGrantPrivileges() *g.QueryStruct {
	return g.NewQueryStruct("AccountRoleGrantPrivileges").
		PredefinedQueryStructField("GlobalPrivileges", "[]GlobalPrivilege", g.InlineOptions()).
		PredefinedQueryStructField("AccountObjectPrivileges", "[]AccountObjectPrivilege", g.InlineOptions()).
		PredefinedQueryStructField("SchemaPrivileges", "[]SchemaPrivilege", g.InlineOptions()).
		PredefinedQueryStructField("SchemaObjectPrivileges", "[]SchemaObjectPrivilege", g.InlineOptions()).
		OptionalSQL("ALL PRIVILEGES").
		WithValidation(g.ExactlyOneValueSet, "AllPrivileges", "GlobalPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges").
		WithAdditionalValidations() // validatePrivileges on each privilege slice
}

func grantOnAccountObject() *g.QueryStruct {
	return g.NewQueryStruct("GrantOnAccountObject").
		PredefinedQueryStructField("Object", "*Object", g.InlineOptions()).
		WithValidation(g.ValidateValueSet, "Object").
		WithAdditionalValidations() // validateUserInput on Object.ObjectType
}

func grantOnSchema() *g.QueryStruct {
	return g.NewQueryStruct("GrantOnSchema").
		OptionalIdentifier("Schema", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("SCHEMA")).
		OptionalIdentifier("AllSchemasInDatabase", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("ALL SCHEMAS IN DATABASE")).
		OptionalIdentifier("FutureSchemasInDatabase", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FUTURE SCHEMAS IN DATABASE")).
		WithValidation(g.ExactlyOneValueSet, "Schema", "AllSchemasInDatabase", "FutureSchemasInDatabase")
}

func grantOnSchemaObjectIn() *g.QueryStruct {
	return g.NewQueryStruct("GrantOnSchemaObjectIn").
		PredefinedQueryStructField("PluralObjectType", "PluralObjectType", g.KeywordOptions().SQL("ALL").Required()).
		OptionalIdentifier("InDatabase", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("IN DATABASE")).
		OptionalIdentifier("InSchema", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("IN SCHEMA")).
		WithValidation(g.ExactlyOneValueSet, "InDatabase", "InSchema")
}

func grantOnSchemaObject() *g.QueryStruct {
	return g.NewQueryStruct("GrantOnSchemaObject").
		PredefinedQueryStructField("SchemaObject", "*Object", g.InlineOptions()).
		OptionalQueryStructField("All", grantOnSchemaObjectIn(), g.KeywordOptions().SQL("ALL")).
		OptionalQueryStructField("Future", grantOnSchemaObjectIn(), g.KeywordOptions().SQL("FUTURE")).
		WithValidation(g.ExactlyOneValueSet, "SchemaObject", "All", "Future")
}

func accountRoleGrantOn() *g.QueryStruct {
	return g.NewQueryStruct("AccountRoleGrantOn").
		OptionalSQL("ACCOUNT").
		OptionalInlineQueryStructField("AccountObject", grantOnAccountObject()).
		OptionalInlineQueryStructField("Schema", grantOnSchema()).
		OptionalInlineQueryStructField("SchemaObject", grantOnSchemaObject()).
		WithValidation(g.ExactlyOneValueSet, "Account", "AccountObject", "Schema", "SchemaObject")
}

func inheritedAccountRoleGrantPrivileges() *g.QueryStruct {
	return g.NewQueryStruct("InheritedAccountRoleGrantPrivileges").
		PredefinedQueryStructField("AccountObjectPrivileges", "[]AccountObjectPrivilege", g.InlineOptions()).
		PredefinedQueryStructField("SchemaPrivileges", "[]SchemaPrivilege", g.InlineOptions()).
		PredefinedQueryStructField("SchemaObjectPrivileges", "[]SchemaObjectPrivilege", g.InlineOptions()).
		OptionalSQL("ALL PRIVILEGES").
		WithValidation(g.ExactlyOneValueSet, "AllPrivileges", "AccountObjectPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges").
		WithAdditionalValidations() // validatePrivileges on each privilege slice
}

func inheritedAccountRoleGrantIn() *g.QueryStruct {
	return g.NewQueryStruct("InheritedAccountRoleGrantIn").
		OptionalSQL("ACCOUNT").
		OptionalIdentifier("Database", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("DATABASE")).
		OptionalIdentifier("Schema", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("SCHEMA")).
		WithValidation(g.ExactlyOneValueSet, "Account", "Database", "Schema").
		WithValidation(g.ValidIdentifierIfSet, "Database").
		WithValidation(g.ValidIdentifierIfSet, "Schema")
}

func databaseRoleGrantPrivileges() *g.QueryStruct {
	return g.NewQueryStruct("DatabaseRoleGrantPrivileges").
		PredefinedQueryStructField("DatabasePrivileges", "[]AccountObjectPrivilege", g.InlineOptions()).
		PredefinedQueryStructField("SchemaPrivileges", "[]SchemaPrivilege", g.InlineOptions()).
		PredefinedQueryStructField("SchemaObjectPrivileges", "[]SchemaObjectPrivilege", g.InlineOptions()).
		OptionalSQL("ALL PRIVILEGES").
		WithValidation(g.ExactlyOneValueSet, "DatabasePrivileges", "SchemaPrivileges", "SchemaObjectPrivileges", "AllPrivileges").
		WithAdditionalValidations() // validatePrivileges on each privilege slice
}

func databaseRoleGrantOn() *g.QueryStruct {
	return g.NewQueryStruct("DatabaseRoleGrantOn").
		OptionalIdentifier("Database", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("DATABASE")).
		OptionalInlineQueryStructField("Schema", grantOnSchema()).
		OptionalInlineQueryStructField("SchemaObject", grantOnSchemaObject()).
		WithValidation(g.ExactlyOneValueSet, "Database", "Schema", "SchemaObject")
}

func inheritedDatabaseRoleGrantPrivileges() *g.QueryStruct {
	return g.NewQueryStruct("InheritedDatabaseRoleGrantPrivileges").
		PredefinedQueryStructField("SchemaPrivileges", "[]SchemaPrivilege", g.InlineOptions()).
		PredefinedQueryStructField("SchemaObjectPrivileges", "[]SchemaObjectPrivilege", g.InlineOptions()).
		OptionalSQL("ALL PRIVILEGES").
		WithValidation(g.ExactlyOneValueSet, "AllPrivileges", "SchemaPrivileges", "SchemaObjectPrivileges").
		WithAdditionalValidations() // validatePrivileges on each privilege slice
}

func inheritedDatabaseRoleGrantIn() *g.QueryStruct {
	return g.NewQueryStruct("InheritedDatabaseRoleGrantIn").
		OptionalIdentifier("Database", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("DATABASE")).
		OptionalIdentifier("Schema", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("SCHEMA")).
		WithValidation(g.ExactlyOneValueSet, "Database", "Schema").
		WithValidation(g.ValidIdentifierIfSet, "Database").
		WithValidation(g.ValidIdentifierIfSet, "Schema")
}

func shareGrantOnTable() *g.QueryStruct {
	return g.NewQueryStruct("OnTable").
		Identifier("Name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("TABLE")).
		Identifier("AllInSchema", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("ALL TABLES IN SCHEMA")).
		WithValidation(g.ExactlyOneValueSet, "Name", "AllInSchema")
}

func shareGrantOn() *g.QueryStruct {
	return g.NewQueryStruct("ShareGrantOn").
		Identifier("Database", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("DATABASE")).
		Identifier("Schema", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("SCHEMA")).
		Identifier("Function", g.KindOfT[sdkcommons.SchemaObjectIdentifierWithArguments](), g.IdentifierOptions().SQL("FUNCTION")).
		OptionalInlineQueryStructField("Table", shareGrantOnTable()).
		Identifier("Tag", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("TAG")).
		Identifier("View", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("VIEW")).
		WithValidation(g.ExactlyOneValueSet, "Database", "Schema", "Function", "Table", "Tag", "View")
}

func ownershipGrantOn() *g.QueryStruct {
	return g.NewQueryStruct("OwnershipGrantOn").
		PredefinedQueryStructField("Object", "*Object", g.InlineOptions()).
		OptionalQueryStructField("All", grantOnSchemaObjectIn(), g.KeywordOptions().SQL("ALL")).
		OptionalQueryStructField("Future", grantOnSchemaObjectIn(), g.KeywordOptions().SQL("FUTURE")).
		WithValidation(g.ExactlyOneValueSet, "Object", "All", "Future")
}

func ownershipGrantTo() *g.QueryStruct {
	return g.NewQueryStruct("OwnershipGrantTo").
		OptionalIdentifier("DatabaseRoleName", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("DATABASE ROLE")).
		OptionalIdentifier("AccountRoleName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("ROLE")).
		WithValidation(g.ExactlyOneValueSet, "DatabaseRoleName", "AccountRoleName")
}

func ownershipCurrentGrants() *g.QueryStruct {
	return g.NewQueryStruct("OwnershipCurrentGrants").
		PredefinedQueryStructField("OutboundPrivileges", "OwnershipCurrentGrantsOutboundPrivileges", g.KeywordOptions().Required()).
		SQL("CURRENT GRANTS")
}

func revokeOwnershipGrantOn() *g.QueryStruct {
	return g.NewQueryStruct("RevokeOwnershipGrantOn").
		OptionalQueryStructField("Future", grantOnSchemaObjectIn(), g.KeywordOptions().SQL("FUTURE")).
		WithValidation(g.ValidateValueSet, "Future")
}

func showGrantsIn() *g.QueryStruct {
	return g.NewQueryStruct("ShowGrantsIn").
		OptionalSQL("ACCOUNT").
		OptionalIdentifier("Schema", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("SCHEMA")).
		OptionalIdentifier("Database", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("DATABASE"))
}

func showGrantsOn() *g.QueryStruct {
	return g.NewQueryStruct("ShowGrantsOn").
		OptionalSQL("ACCOUNT").
		PredefinedQueryStructField("Object", "*Object", g.InlineOptions())
}

func showGrantsToShare() *g.QueryStruct {
	return g.NewQueryStruct("ShowGrantsToShare").
		Identifier("Name", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("SHARE")).
		OptionalIdentifier("InApplicationPackage", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("IN APPLICATION PACKAGE"))
}

func showGrantsTo() *g.QueryStruct {
	return g.NewQueryStruct("ShowGrantsTo").
		Identifier("Application", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION")).
		Identifier("ApplicationRole", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION ROLE")).
		Identifier("Role", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("ROLE")).
		Identifier("User", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("USER")).
		OptionalInlineQueryStructField("Share", showGrantsToShare()).
		Identifier("DatabaseRole", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("DATABASE ROLE"))
}

func showGrantsOf() *g.QueryStruct {
	return g.NewQueryStruct("ShowGrantsOf").
		Identifier("ApplicationRole", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("APPLICATION ROLE")).
		Identifier("Role", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("ROLE")).
		Identifier("DatabaseRole", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("DATABASE ROLE")).
		Identifier("Share", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("SHARE"))
}

var grantPairs = g.StructPair("grantRow", "Grant").
	Time("created_on").
	Text("privilege").
	Field("granted_on", "string", "ObjectType", g.WithCustomParser("ObjectTypeFromShowGrants")).
	Field("grant_on", "string", "ObjectType", g.WithCustomParser("ObjectTypeFromShowGrants")).
	Field("name", "string", "ObjectIdentifier", g.WithManualConvert()).
	Field("granted_to", "string", "ObjectType", g.WithCustomParser("ObjectTypeFromShowGrants")).
	Field("grant_to", "string", "ObjectType", g.WithCustomParser("ObjectTypeFromShowGrants")).
	Field("grantee_name", "string", "ObjectIdentifier", g.WithManualConvert()).
	Bool("grant_option").
	AccountObjectIdentifier("granted_by", g.WithPlainFieldName("GrantedBy")).
	OptionalBoolFromText("is_inherited", g.WithBoolTrueValue("true")).
	OptionalEnum("inherited_from", grantInheritedFromEnumDef).
	OptionalText("inherited_from_database", g.WithCustomParser("trimQuotes")).
	OptionalText("inherited_from_schema", g.WithCustomParser("trimQuotes"))

// NameSingular is empty on purpose. The generator names Options `{kind}{NameSingular}Options`;
// a non-empty singular ("Grant") would produce GrantPrivilegesToAccountRoleGrantOptions etc.
// Identifier kind is a required placeholder — grants are not a named Snowflake object.
var grantsDef = g.NewInterface(
	"Grants",
	"",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CustomOperation(
	// Unexported unit op: single-statement GRANT PRIVILEGES TO ROLE. Public orchestrator is WithCustomInterfaceMethod.
	"grantPrivilegesToAccountRole",
	"https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#syntax",
	g.NewQueryStruct("GrantPrivilegesToAccountRole").
		Grant().
		OptionalInlineQueryStructField("Privileges", accountRoleGrantPrivileges()).
		OptionalQueryStructField("On", accountRoleGrantOn(), g.KeywordOptions().SQL("ON")).
		Identifier("AccountRole", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("TO ROLE").Required()).
		OptionalSQL("WITH GRANT OPTION").
		WithValidation(g.ValidateValueSet, "Privileges").
		WithValidation(g.ValidateValueSet, "On"),
).CustomOperation(
	"revokePrivilegesFromAccountRole",
	"https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege#syntax",
	g.NewQueryStruct("RevokePrivilegesFromAccountRole").
		Revoke().
		OptionalSQL("GRANT OPTION FOR").
		OptionalInlineQueryStructField("Privileges", accountRoleGrantPrivileges()).
		OptionalQueryStructField("On", accountRoleGrantOn(), g.KeywordOptions().SQL("ON")).
		Identifier("AccountRole", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FROM ROLE").Required()).
		OptionalSQL("RESTRICT").
		OptionalSQL("CASCADE").
		WithValidation(g.ValidateValueSet, "Privileges").
		WithValidation(g.ValidateValueSet, "On").
		WithValidation(g.ValidIdentifier, "AccountRole").
		WithValidation(g.ConflictingFields, "Restrict", "Cascade"),
).CustomOperation(
	"grantPrivilegesToDatabaseRole",
	"https://docs.snowflake.com/en/sql-reference/sql/grant-privilege#syntax",
	g.NewQueryStruct("GrantPrivilegesToDatabaseRole").
		Grant().
		OptionalInlineQueryStructField("Privileges", databaseRoleGrantPrivileges()).
		OptionalQueryStructField("On", databaseRoleGrantOn(), g.KeywordOptions().SQL("ON")).
		Identifier("DatabaseRole", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("TO DATABASE ROLE").Required()).
		OptionalSQL("WITH GRANT OPTION").
		WithValidation(g.ValidateValueSet, "Privileges").
		WithValidation(g.ValidateValueSet, "On"),
).CustomOperation(
	"revokePrivilegesFromDatabaseRole",
	"https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege#syntax",
	g.NewQueryStruct("RevokePrivilegesFromDatabaseRole").
		Revoke().
		OptionalSQL("GRANT OPTION FOR").
		OptionalInlineQueryStructField("Privileges", databaseRoleGrantPrivileges()).
		OptionalQueryStructField("On", databaseRoleGrantOn(), g.KeywordOptions().SQL("ON")).
		Identifier("DatabaseRole", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("FROM DATABASE ROLE").Required()).
		OptionalSQL("RESTRICT").
		OptionalSQL("CASCADE").
		WithValidation(g.ValidateValueSet, "Privileges").
		WithValidation(g.ValidateValueSet, "On").
		WithValidation(g.ValidIdentifier, "DatabaseRole").
		WithValidation(g.ConflictingFields, "Restrict", "Cascade"),
).CustomOperation(
	"grantOwnership",
	"https://docs.snowflake.com/en/sql-reference/sql/grant-ownership#syntax",
	g.NewQueryStruct("GrantOwnership").
		SQL("GRANT OWNERSHIP").
		QueryStructField("On", ownershipGrantOn(), g.KeywordOptions().SQL("ON").Required()).
		QueryStructField("To", ownershipGrantTo(), g.KeywordOptions().SQL("TO").Required()).
		OptionalInlineQueryStructField("CurrentGrants", ownershipCurrentGrants()),
).CustomOperation(
	"GrantInheritedPrivilegesToAccountRole",
	"https://docs.snowflake.com/en/user-guide/inherited-grants-using#syntax",
	g.NewQueryStruct("GrantInheritedPrivilegesToAccountRole").
		SQL("GRANT INHERITED").
		QueryStructField("Privileges", inheritedAccountRoleGrantPrivileges(), g.InlineOptions()).
		PredefinedQueryStructField("OnAll", "PluralObjectType", g.ParameterOptions().NoEquals().SQL("ON ALL").Required()).
		QueryStructField("In", inheritedAccountRoleGrantIn(), g.KeywordOptions().SQL("IN").Required()).
		Identifier("AccountRole", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("TO ROLE").Required()).
		WithValidation(g.ValidateValueSet, "OnAll").
		WithValidation(g.ValidIdentifier, "AccountRole"),
).CustomOperation(
	"RevokeInheritedPrivilegesFromAccountRole",
	"https://docs.snowflake.com/en/user-guide/inherited-grants-using#syntax",
	g.NewQueryStruct("RevokeInheritedPrivilegesFromAccountRole").
		SQL("REVOKE INHERITED").
		QueryStructField("Privileges", inheritedAccountRoleGrantPrivileges(), g.InlineOptions()).
		PredefinedQueryStructField("OnAll", "PluralObjectType", g.ParameterOptions().NoEquals().SQL("ON ALL").Required()).
		QueryStructField("In", inheritedAccountRoleGrantIn(), g.KeywordOptions().SQL("IN").Required()).
		Identifier("AccountRole", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FROM ROLE").Required()).
		WithValidation(g.ValidateValueSet, "OnAll").
		WithValidation(g.ValidIdentifier, "AccountRole"),
).CustomOperation(
	"GrantInheritedPrivilegesToDatabaseRole",
	"https://docs.snowflake.com/en/user-guide/inherited-grants-using#syntax",
	g.NewQueryStruct("GrantInheritedPrivilegesToDatabaseRole").
		SQL("GRANT INHERITED").
		QueryStructField("Privileges", inheritedDatabaseRoleGrantPrivileges(), g.InlineOptions()).
		PredefinedQueryStructField("OnAll", "PluralObjectType", g.ParameterOptions().NoEquals().SQL("ON ALL").Required()).
		QueryStructField("In", inheritedDatabaseRoleGrantIn(), g.KeywordOptions().SQL("IN").Required()).
		Identifier("DatabaseRole", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("TO DATABASE ROLE").Required()).
		WithValidation(g.ValidateValueSet, "OnAll").
		WithValidation(g.ValidIdentifier, "DatabaseRole"),
).CustomOperation(
	"RevokeInheritedPrivilegesFromDatabaseRole",
	"https://docs.snowflake.com/en/user-guide/inherited-grants-using#syntax",
	g.NewQueryStruct("RevokeInheritedPrivilegesFromDatabaseRole").
		SQL("REVOKE INHERITED").
		QueryStructField("Privileges", inheritedDatabaseRoleGrantPrivileges(), g.InlineOptions()).
		PredefinedQueryStructField("OnAll", "PluralObjectType", g.ParameterOptions().NoEquals().SQL("ON ALL").Required()).
		QueryStructField("In", inheritedDatabaseRoleGrantIn(), g.KeywordOptions().SQL("IN").Required()).
		Identifier("DatabaseRole", g.KindOfT[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("FROM DATABASE ROLE").Required()).
		WithValidation(g.ValidateValueSet, "OnAll").
		WithValidation(g.ValidIdentifier, "DatabaseRole"),
).CustomOperation(
	"GrantPrivilegeToShare",
	"https://docs.snowflake.com/en/sql-reference/sql/grant-privilege-share",
	g.NewQueryStruct("GrantPrivilegeToShare").
		Grant().
		PredefinedQueryStructField("Privileges", "[]ObjectPrivilege", g.InlineOptions()).
		OptionalQueryStructField("On", shareGrantOn(), g.KeywordOptions().SQL("ON")).
		Identifier("To", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("TO SHARE").Required()).
		WithValidation(g.ValidIdentifier, "To").
		WithValidation(g.ValidateValueSet, "On").
		WithValidation(g.ValidateValueSet, "Privileges").
		WithAdditionalValidations(), // validatePrivileges on each privilege
).CustomOperation(
	"RevokePrivilegeFromShare",
	"https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege-share",
	g.NewQueryStruct("RevokePrivilegeFromShare").
		Revoke().
		PredefinedQueryStructField("Privileges", "[]ObjectPrivilege", g.InlineOptions()).
		OptionalQueryStructField("On", shareGrantOn(), g.KeywordOptions().SQL("ON")).
		Identifier("From", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("FROM SHARE").Required()).
		WithValidation(g.ValidIdentifier, "From").
		WithValidation(g.ValidateValueSet, "On").
		WithValidation(g.ValidateValueSet, "Privileges").
		WithAdditionalValidations(), // validatePrivileges on each privilege
).CustomOperation(
	"RevokeOwnership",
	"https://docs.snowflake.com/en/sql-reference/sql/revoke-privilege#syntax",
	g.NewQueryStruct("RevokeOwnership").
		SQL("REVOKE OWNERSHIP").
		QueryStructField("On", revokeOwnershipGrantOn(), g.KeywordOptions().SQL("ON").Required()).
		QueryStructField("From", ownershipGrantTo(), g.KeywordOptions().SQL("FROM").Required()).
		OptionalSQL("RESTRICT").
		OptionalSQL("CASCADE").
		WithValidation(g.ValidateValueSet, "On").
		WithValidation(g.ValidateValueSet, "From").
		WithValidation(g.ConflictingFields, "Restrict", "Cascade"),
).CustomShowOperationWithPairedStructs(
	"showGrants",
	g.ShowMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/show-grants",
	grantPairs,
	g.NewQueryStruct("ShowGrant").
		Show().
		OptionalSQL("INHERITED").
		OptionalSQL("FUTURE").
		SQL("GRANTS").
		OptionalQueryStructField("On", showGrantsOn(), g.KeywordOptions().SQL("ON")).
		OptionalQueryStructField("To", showGrantsTo(), g.KeywordOptions().SQL("TO")).
		OptionalQueryStructField("Of", showGrantsOf(), g.KeywordOptions().SQL("OF")).
		OptionalQueryStructField("In", showGrantsIn(), g.KeywordOptions().SQL("IN")).
		WithValidation(g.MoreThanOneValueSet, "On", "To", "Of", "In").
		WithValidation(g.ConflictingFields, "Inherited", "Future"),
).WithCustomInterfaceMethod(
	"GrantPrivilegesToAccountRole",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*GrantPrivilegesToAccountRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"RevokePrivilegesFromAccountRole",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokePrivilegesFromAccountRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"GrantPrivilegesToDatabaseRole",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*GrantPrivilegesToDatabaseRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"RevokePrivilegesFromDatabaseRole",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokePrivilegesFromDatabaseRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"GrantOwnership",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*GrantOwnershipRequest")},
	"error",
).WithCustomInterfaceMethod(
	"RevokePrivilegesFromAccountRoleSafely",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokePrivilegesFromAccountRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"RevokeInheritedPrivilegesFromAccountRoleSafely",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokeInheritedPrivilegesFromAccountRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"RevokePrivilegesFromDatabaseRoleSafely",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokePrivilegesFromDatabaseRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"RevokeInheritedPrivilegesFromDatabaseRoleSafely",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokeInheritedPrivilegesFromDatabaseRoleRequest")},
	"error",
).WithCustomInterfaceMethod(
	"RevokePrivilegeFromShareSafely",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*RevokePrivilegeFromShareRequest")},
	"error",
).WithCustomInterfaceMethod(
	"Show",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("request", "*ShowGrantsRequest")},
	"[]Grant", "error",
).WithEnums(grantInheritedFromEnumDef)
