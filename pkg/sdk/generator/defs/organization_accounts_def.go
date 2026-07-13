package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var OrganizationAccountEditionEnumDef = g.NewEnum(
	"OrganizationAccountEdition", "OrganizationAccountEditions",
	"ENTERPRISE", "BUSINESS_CRITICAL",
)

var organizationAccountsDef = g.NewInterface(
	"OrganizationAccounts",
	"OrganizationAccount",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-organization-account",
		g.NewQueryStruct("CreateOrganizationAccount").
			Create().
			SQL("ORGANIZATION ACCOUNT").
			Name().
			TextAssignment("ADMIN_NAME", g.ParameterOptions().Required().NoQuotes()).
			OptionalTextAssignment("ADMIN_PASSWORD", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("ADMIN_RSA_PUBLIC_KEY", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("FIRST_NAME", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("LAST_NAME", g.ParameterOptions().SingleQuotes()).
			TextAssignment("EMAIL", g.ParameterOptions().Required().SingleQuotes()).
			OptionalBooleanAssignment("MUST_CHANGE_PASSWORD", g.ParameterOptions()).
			EnumAssignment("EDITION", OrganizationAccountEditionEnumDef, g.ParameterOptions().Required().NoQuotes()).
			OptionalTextAssignment("REGION_GROUP", g.ParameterOptions().DoubleQuotes()).
			OptionalTextAssignment("REGION", g.ParameterOptions().DoubleQuotes()).
			OptionalComment().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.AtLeastOneValueSet, "AdminPassword", "AdminRsaPublicKey"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-organization-account",
		g.NewQueryStruct("AlterOrganizationAccount").
			Alter().
			SQL("ORGANIZATION ACCOUNT").
			OptionalIdentifier("Name", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions()).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("OrganizationAccountSet").
					// Currently, Organization Accounts use the same set of parameters as regular accounts
					PredefinedQueryStructField("Parameters", g.KindOfTPointer[sdkcommons.AccountParameters](), g.ListOptions().NoParentheses()).
					OptionalIdentifier("ResourceMonitor", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("RESOURCE_MONITOR")).
					OptionalIdentifier("PasswordPolicy", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PASSWORD POLICY")).
					OptionalIdentifier("SessionPolicy", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("SESSION POLICY")).
					OptionalComment().
					WithValidation(g.ExactlyOneValueSet, "Parameters", "ResourceMonitor", "PasswordPolicy", "SessionPolicy", "Comment"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("OrganizationAccountUnset").
					PredefinedQueryStructField("Parameters", g.KindOfTPointer[sdkcommons.AccountParametersUnset](), g.ListOptions().NoParentheses()).
					OptionalSQL("RESOURCE_MONITOR").
					OptionalSQL("PASSWORD POLICY").
					OptionalSQL("SESSION POLICY").
					OptionalSQL("COMMENT").
					WithValidation(g.ExactlyOneValueSet, "Parameters", "ResourceMonitor", "PasswordPolicy", "SessionPolicy", "Comment"),
				g.KeywordOptions().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalQueryStructField(
				"RenameTo",
				g.NewQueryStruct("OrganizationAccountRename").
					Identifier("RenameTo", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Required().SQL("RENAME TO")).
					OptionalBooleanAssignment("SAVE_OLD_URL", g.ParameterOptions()),
				g.KeywordOptions(),
			).
			OptionalSQL("DROP OLD URL").
			WithValidation(g.ValidIdentifierIfSet, "Name").
			WithValidation(g.ConflictingFields, "Name", "Set").
			WithValidation(g.ConflictingFields, "Name", "Unset").
			WithValidation(g.ConflictingFields, "Name", "SetTags").
			WithValidation(g.ConflictingFields, "Name", "UnsetTags").
			WithValidation(g.ExactlyOneValueSet, "Set", "Unset", "SetTags", "UnsetTags", "RenameTo", "DropOldUrl"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-organization-accounts",
		g.StructPair("organizationAccountDbRow", "OrganizationAccount").
			Text("organization_name").
			Text("account_name").
			Text("snowflake_region").
			Enum("edition", OrganizationAccountEditionEnumDef).
			Text("account_url").
			Text("created_on").
			OptionalText("comment").
			Text("account_locator").
			Text("account_locator_url").
			Number("managed_accounts").
			Text("consumption_billing_entity_name").
			OptionalText("marketplace_consumer_billing_entity_name").
			OptionalText("marketplace_provider_billing_entity_name").
			OptionalText("old_account_url").
			Bool("is_org_admin").
			OptionalText("account_old_url_saved_on").
			OptionalText("account_old_url_last_used").
			OptionalText("organization_old_url").
			OptionalText("organization_old_url_saved_on").
			OptionalText("organization_old_url_last_used").
			Bool("is_events_account").
			Bool("is_organization_account"),
		g.NewQueryStruct("ShowOrganizationAccounts").
			Show().
			SQL("ORGANIZATION ACCOUNTS").
			OptionalLike(),
	).
	WithCustomInterfaceMethod("ShowParameters", "", nil, "[]*Parameter", "error").
	WithCustomInterfaceMethod("UnsetAllParameters", "", nil, "error").
	WithCustomInterfaceMethod(
		"UnsetPolicySafely",
		"UnsetPolicySafely unsets a policy on the current account by a given supported kind.\nIt ignores an error that occurs on the Snowflake side whenever you try to unset policy which is already unset.",
		[]*g.MethodParameter{g.NewMethodParameter("kind", "PolicyKind")},
		"error",
	).
	WithCustomInterfaceMethod(
		"SetPolicySafely",
		"SetPolicySafely sets a policy on the current account by a given supported kind.\nIt firstly tries to unset the policy with UnsetPolicySafely method to make sure that the policy is not set,\nthen proceeds by setting the passed policy on the organization account.",
		[]*g.MethodParameter{
			g.NewMethodParameter("kind", "PolicyKind"),
			g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]()),
		},
		"error",
	).
	WithCustomInterfaceMethod("UnsetAll", "", nil, "error").
	WithEnums(
		OrganizationAccountEditionEnumDef,
	).
	WithShowObjectType("Account").
	WithShowByIDFindPredicateKind(g.ShowByIDFindPredicateAccountName)
