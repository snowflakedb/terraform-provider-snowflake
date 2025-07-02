package sdk

import (
	"fmt"
	"slices"
	"strings"

	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

type OrganizationAccountEdition string

var (
	OrganizationAccountEditionEnterprise       OrganizationAccountEdition = "ENTERPRISE"
	OrganizationAccountEditionBusinessCritical OrganizationAccountEdition = "BUSINESS_CRITICAL"
)

var AllOrganizationAccountEditions = []OrganizationAccountEdition{
	OrganizationAccountEditionEnterprise,
	OrganizationAccountEditionBusinessCritical,
}

func ToOrganizationAccountEdition(s string) (OrganizationAccountEdition, error) {
	s = strings.ToUpper(s)
	if !slices.Contains(AllOrganizationAccountEditions, OrganizationAccountEdition(s)) {
		return "", fmt.Errorf("invalid organization account edition: %s", s)
	}
	return OrganizationAccountEdition(s), nil
}

var OrganizationAccountsDef = g.NewInterface(
	"OrganizationAccounts",
	"OrganizationAccount",
	g.KindOfT[AccountObjectIdentifier](),
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
			Assignment("EDITION", g.KindOfT[OrganizationAccountEdition](), g.ParameterOptions().Required().NoQuotes()).
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
			OptionalIdentifier("Name", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions()).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("OrganizationAccountSet").
					// Currently, Organization Accounts use the same set of parameters as regular accounts
					PredefinedQueryStructField("Parameters", g.KindOfTPointer[AccountParameters](), g.ListOptions().NoParentheses()).
					OptionalIdentifier("ResourceMonitor", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("RESOURCE_MONITOR")).
					OptionalIdentifier("PasswordPolicy", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("PASSWORD POLICY")).
					OptionalIdentifier("SessionPolicy", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("SESSION POLICY")).
					WithValidation(g.ExactlyOneValueSet, "Parameters", "ResourceMonitor", "PasswordPolicy", "SessionPolicy"),
				g.KeywordOptions().SQL("SET"),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("OrganizationAccountUnset").
					PredefinedQueryStructField("Parameters", g.KindOfTPointer[AccountParametersUnset](), g.ListOptions().NoParentheses()).
					OptionalSQL("RESOURCE_MONITOR").
					OptionalSQL("PASSWORD POLICY").
					OptionalSQL("SESSION POLICY").
					WithValidation(g.ExactlyOneValueSet, "Parameters", "ResourceMonitor", "PasswordPolicy", "SessionPolicy"),
				g.KeywordOptions().SQL("UNSET"),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalQueryStructField(
				"RenameTo",
				g.NewQueryStruct("OrganizationAccountRename").
					Identifier("NewName", g.KindOfTPointer[AccountObjectIdentifier](), g.IdentifierOptions().Required().SQL("RENAME TO")).
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
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-organization-accounts",
		g.DbStruct("organizationAccountDbRow").
			Text("organization_name").
			Text("account_name").
			Text("snowflake_region").
			Text("edition").
			Text("account_url").
			Text("created_on").
			Text("comment").
			Text("account_locator").
			Text("account_locator_url").
			Number("managed_accounts").
			Text("consumption_billing_entity_name").
			OptionalText("marketplace_consumer_billing_entity_name").
			Text("marketplace_provider_billing_entity_name").
			OptionalText("old_account_url").
			Bool("is_org_admin").
			OptionalText("account_old_url_saved_on").
			OptionalText("account_old_url_last_used").
			OptionalText("organization_old_url").
			OptionalText("organization_old_url_saved_on").
			OptionalText("organization_old_url_last_used").
			Bool("is_events_account").
			Bool("is_organization_account"),
		g.PlainStruct("OrganizationAccount").
			Text("OrganizationName").
			Text("AccountName").
			Text("SnowflakeRegion").
			Field("Edition", g.KindOfT[OrganizationAccountEdition]()).
			Text("AccountUrl").
			Text("CreatedOn").
			Text("Comment").
			Text("AccountLocator").
			Text("AccountLocatorUrl").
			Number("ManagedAccounts").
			Text("ConsumptionBillingEntityName").
			OptionalText("MarketplaceConsumerBillingEntityName").
			Text("MarketplaceProviderBillingEntityName").
			OptionalText("OldAccountUrl").
			Bool("IsOrgAdmin").
			OptionalText("AccountOldUrlSavedOn").
			OptionalText("AccountOldUrlLastUsed").
			OptionalText("OrganizationOldUrl").
			OptionalText("OrganizationOldUrlSavedOn").
			OptionalText("OrganizationOldUrlLastUsed").
			Bool("IsEventsAccount").
			Bool("IsOrganizationAccount"),
		g.NewQueryStruct("ShowOrganizationAccounts").
			Show().
			SQL("ORGANIZATION ACCOUNTS").
			OptionalLike(),
	)
