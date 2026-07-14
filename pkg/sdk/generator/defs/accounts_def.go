//go:build sdk_generation

package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var accountEditionDef = g.NewEnum("AccountEdition", "AccountEditions", "STANDARD", "ENTERPRISE", "BUSINESS_CRITICAL")

func accountRename() *g.QueryStruct {
	return g.NewQueryStruct("AccountRename").
		Identifier("NewName", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO").Required()).
		OptionalAssignmentWithFieldName("SAVE_OLD_URL", "bool", g.ParameterOptions(), "SaveOldURL")
}

func accountDrop() *g.QueryStruct {
	return g.NewQueryStruct("AccountDrop").
		OptionalSQLWithCustomFieldName("OldUrl", "DROP OLD URL").
		OptionalSQLWithCustomFieldName("OldOrganizationUrl", "DROP OLD ORGANIZATION URL")
}

func accountLevelParameters() *g.QueryStruct {
	return g.NewQueryStruct("AccountLevelParameters").
		PredefinedQueryStructField("AccountParameters", "*LegacyAccountParameters", g.ListOptions().NoParentheses()).
		PredefinedQueryStructField("SessionParameters", "*SessionParameters", g.ListOptions().NoParentheses()).
		PredefinedQueryStructField("ObjectParameters", "*ObjectParameters", g.ListOptions().NoParentheses()).
		PredefinedQueryStructField("UserParameters", "*UserParameters", g.ListOptions().NoParentheses())
}

func accountLevelParametersUnset() *g.QueryStruct {
	return g.NewQueryStruct("AccountLevelParametersUnset").
		PredefinedQueryStructField("AccountParameters", "*LegacyAccountParametersUnset", g.ListOptions().NoParentheses()).
		PredefinedQueryStructField("SessionParameters", "*SessionParametersUnset", g.ListOptions().NoParentheses()).
		PredefinedQueryStructField("ObjectParameters", "*ObjectParametersUnset", g.ListOptions().NoParentheses()).
		PredefinedQueryStructField("UserParameters", "*UserParametersUnset", g.ListOptions().NoParentheses())
}

func accountFeaturePolicySet() *g.QueryStruct {
	return g.NewQueryStruct("AccountFeaturePolicySet").
		OptionalIdentifier("FeaturePolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("FEATURE POLICY")).
		SQL("FOR ALL APPLICATIONS")
}

func accountFeaturePolicyUnset() *g.QueryStruct {
	return g.NewQueryStruct("AccountFeaturePolicyUnset").
		OptionalSQL("FEATURE POLICY").
		SQL("FOR ALL APPLICATIONS")
}

func accountSet() *g.QueryStruct {
	return g.NewQueryStruct("AccountSet").
		PredefinedQueryStructField("Parameters", "*AccountParameters", g.ListOptions().NoParentheses()).
		OptionalQueryStructField("LegacyParameters", accountLevelParameters(), g.ListOptions().NoParentheses()).
		OptionalIdentifier("ResourceMonitor", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("RESOURCE_MONITOR")).
		OptionalIdentifier("PackagesPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PACKAGES POLICY")).
		OptionalIdentifier("PasswordPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PASSWORD POLICY")).
		OptionalIdentifier("SessionPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("SESSION POLICY")).
		OptionalIdentifier("AuthenticationPolicy", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("AUTHENTICATION POLICY")).
		OptionalQueryStructField("FeaturePolicySet", accountFeaturePolicySet(), g.KeywordOptions()).
		OptionalTextAssignment("CONSUMPTION_BILLING_ENTITY", g.ParameterOptions().DoubleQuotes()).
		OptionalAssignmentWithFieldName("IS_ORG_ADMIN", "bool", g.ParameterOptions(), "OrgAdmin").
		OptionalSQL("FORCE").
		WithValidation(g.ExactlyOneValueSet, "Parameters", "LegacyParameters", "ResourceMonitor", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "FeaturePolicySet", "OrgAdmin", "ConsumptionBillingEntity").
		WithAdditionalValidations()
}

func accountUnset() *g.QueryStruct {
	return g.NewQueryStruct("AccountUnset").
		PredefinedQueryStructField("Parameters", "*AccountParametersUnset", g.ListOptions().NoParentheses()).
		OptionalQueryStructField("LegacyParameters", accountLevelParametersUnset(), g.ListOptions().NoParentheses()).
		OptionalSQL("AUTHENTICATION POLICY").
		OptionalQueryStructField("FeaturePolicyUnset", accountFeaturePolicyUnset(), g.KeywordOptions()).
		OptionalSQL("PACKAGES POLICY").
		OptionalSQL("PASSWORD POLICY").
		OptionalSQL("SESSION POLICY").
		OptionalSQL("RESOURCE_MONITOR").
		OptionalSQL("CONSUMPTION_BILLING_ENTITY").
		WithValidation(g.ExactlyOneValueSet, "Parameters", "LegacyParameters", "PackagesPolicy", "PasswordPolicy", "SessionPolicy", "AuthenticationPolicy", "ResourceMonitor", "FeaturePolicyUnset", "ConsumptionBillingEntity").
		WithAdditionalValidations()
}

var accountPairs = g.StructPair("accountDBRow", "Account").
	Text("organization_name").
	Text("account_name").
	OptionalText("region_group").
	Text("snowflake_region").
	OptionalEnum("edition", accountEditionDef).
	OptionalText("account_url", g.WithDbFieldName("AccountURL"), g.WithPlainFieldName("AccountURL")).
	OptionalTime("created_on").
	OptionalText("comment").
	Text("account_locator").
	OptionalText("account_locator_url", g.WithDbFieldName("AccountLocatorURL"), g.WithPlainFieldName("AccountLocatorUrl")).
	Field("managed_accounts", "sql.NullInt32", "*int", g.WithManualConvert()).
	OptionalText("consumption_billing_entity_name").
	OptionalText("marketplace_consumer_billing_entity_name").
	OptionalText("marketplace_provider_billing_entity_name").
	OptionalText("old_account_url", g.WithDbFieldName("OldAccountURL"), g.WithPlainFieldName("OldAccountURL")).
	OptionalBool("is_org_admin").
	OptionalTime("account_old_url_saved_on").
	OptionalTime("account_old_url_last_used").
	OptionalText("organization_old_url").
	OptionalTime("organization_old_url_saved_on").
	OptionalTime("organization_old_url_last_used").
	OptionalBool("is_events_account").
	Bool("is_organization_account").
	OptionalTime("dropped_on").
	OptionalTime("scheduled_deletion_time").
	OptionalTime("restored_on").
	OptionalText("moved_to_organization").
	OptionalText("moved_on").
	OptionalTime("organization_URL_expiration_on", g.WithDbFieldName("OrganizationUrlExpirationOn"), g.WithPlainFieldName("OrganizationUrlExpirationOn"))

var accountsDef = g.NewInterface(
	"Accounts",
	"Account",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-account",
	g.NewQueryStruct("CreateAccount").
		Create().
		SQL("ACCOUNT").
		Name().
		TextAssignment("ADMIN_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("ADMIN_PASSWORD", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("ADMIN_RSA_PUBLIC_KEY", g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("ADMIN_USER_TYPE", "UserType", g.ParameterOptions()).
		OptionalTextAssignment("FIRST_NAME", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("LAST_NAME", g.ParameterOptions().SingleQuotes()).
		TextAssignment("EMAIL", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("MUST_CHANGE_PASSWORD", g.ParameterOptions()).
		EnumAssignment("EDITION", accountEditionDef, g.ParameterOptions()).
		OptionalTextAssignment("REGION_GROUP", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("REGION", g.ParameterOptions().SingleQuotes()).
		OptionalComment().
		OptionalTextAssignment("CONSUMPTION_BILLING_ENTITY", g.ParameterOptions().DoubleQuotes()).
		OptionalBooleanAssignment("POLARIS", g.ParameterOptions()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.AtLeastOneValueSet, "AdminPassword", "AdminRsaPublicKey"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-account",
	g.NewQueryStruct("AlterAccount").
		Alter().
		SQL("ACCOUNT").
		OptionalIdentifier("Name", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions()).
		OptionalQueryStructField("Set", accountSet(), g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("Unset", accountUnset(), g.ListOptions().NoParentheses().SQL("UNSET")).
		PredefinedQueryStructField("SetTag", "[]TagAssociation", g.KeywordOptions().SQL("SET TAG")).
		PredefinedQueryStructField("UnsetTag", "[]ObjectIdentifier", g.KeywordOptions().SQL("UNSET TAG")).
		OptionalInlineQueryStructField("Rename", accountRename()).
		OptionalInlineQueryStructField("Drop", accountDrop()).
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTag", "UnsetTag", "Drop", "Rename").
		WithAdditionalValidations(),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-account",
	g.NewQueryStruct("DropAccount").
		Drop().
		SQL("ACCOUNT").
		IfExists().
		Name().
		OptionalNumberAssignment("GRACE_PERIOD_IN_DAYS", g.ParameterOptions()).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"Undrop",
	"https://docs.snowflake.com/en/sql-reference/sql/undrop-account",
	g.NewQueryStruct("UndropAccount").
		SQL("UNDROP").
		SQL("ACCOUNT").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-organisation-accounts",
	accountPairs,
	g.NewQueryStruct("ShowAccount").
		Show().
		SQL("ACCOUNTS").
		OptionalSQL("HISTORY").
		OptionalLike(),
	g.ShowByIDLikeFiltering,
).
	WithShowByIDFindPredicateKind(g.ShowByIDFindPredicateAccountName).
	WithCustomInterfaceMethod(
		"ShowParameters",
		"",
		[]*g.MethodParameter{},
		"[]*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"UnsetAllParameters",
		"",
		[]*g.MethodParameter{},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetAllPoliciesSafely",
		"UnsetAllPoliciesSafely calls UnsetPolicySafely for every policy that can be unset from the current account.",
		[]*g.MethodParameter{},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetPolicySafely",
		"UnsetPolicySafely unsets a policy on the current account by a given supported kind.\nIt ignores an error that occurs on the Snowflake side whenever you try to unset policy which is already unset.",
		[]*g.MethodParameter{
			g.NewMethodParameter("kind", "PolicyKind"),
		},
		"error",
	).
	WithCustomInterfaceMethod(
		"UnsetAll",
		"UnsetAll unsets all policies and parameters that can be attached to the current account.",
		[]*g.MethodParameter{},
		"error",
	).
	WithEnums(accountEditionDef)
