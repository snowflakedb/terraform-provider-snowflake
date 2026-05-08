package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var managedAccountPairs = g.StructPair("managedAccountDBRow", "ManagedAccount").
	OptionalText("account_name", g.WithPlainFieldName("Name"), g.WithRequiredInPlain()).
	Text("cloud").
	Text("region").
	OptionalText("account_locator", g.WithPlainFieldName("Locator"), g.WithRequiredInPlain()).
	Text("created_on").
	OptionalText("account_url", g.WithPlainFieldName("URL"), g.WithRequiredInPlain()).
	Text("account_locator_url", g.WithPlainFieldName("AccountLocatorURL")).
	Bool("is_reader").
	OptionalText("comment")

var managedAccountsDef = g.NewInterface(
	"ManagedAccounts",
	"ManagedAccount",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-managed-account",
		g.NewQueryStruct("CreateManagedAccount").
			Create().
			SQL("MANAGED ACCOUNT").
			Name().
			QueryStructField(
				"CreateManagedAccountParams",
				g.NewQueryStruct("CreateManagedAccountParams").
					TextAssignment("ADMIN_NAME", g.ParameterOptions().SingleQuotes().Required()).
					TextAssignment("ADMIN_PASSWORD", g.ParameterOptions().SingleQuotes().Required()).
					PredefinedQueryStructField("typeProvider", "string", g.StaticOptions().SQL("TYPE = READER")).
					OptionalComment().
					WithValidation(g.ValidateValueSet, "AdminName").
					WithValidation(g.ValidateValueSet, "AdminPassword"),
				g.ListOptions().NoParentheses().Required(),
			).
			WithValidation(g.ValidIdentifier, "name"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-managed-account",
		g.NewQueryStruct("DropManagedAccount").
			Drop().
			SQL("MANAGED ACCOUNT").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-managed-accounts",
		managedAccountPairs,
		g.NewQueryStruct("ShowManagedAccounts").
			Show().
			SQL("MANAGED ACCOUNTS").
			OptionalLike(),
	)
