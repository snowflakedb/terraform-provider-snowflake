package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var connectionsDef = g.NewInterface(
	"Connections",
	"Connection",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-connection",
	g.NewQueryStruct("CreateConnection").
		Create().
		SQL("CONNECTION").
		IfNotExists().
		Name().
		OptionalIdentifier(
			"AsReplicaOf",
			g.KindOfT[sdkcommons.ExternalObjectIdentifier](),
			g.IdentifierOptions().SQL("AS REPLICA OF"),
		).
		OptionalComment().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "AsReplicaOf"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-connection",
	g.NewQueryStruct("Alter").
		Alter().
		SQL("CONNECTION").
		IfExists().
		Name().
		OptionalQueryStructField(
			"EnableConnectionFailover",
			g.NewQueryStruct("EnableConnectionFailover").
				List("ToAccounts", "AccountIdentifier", g.ListOptions().NoParentheses().Required()).
				WithValidation(g.AtLeastOneValueSet, "ToAccounts"),
			g.KeywordOptions().SQL("ENABLE FAILOVER TO ACCOUNTS"),
		).
		OptionalQueryStructField(
			"DisableConnectionFailover",
			g.NewQueryStruct("DisableConnectionFailover").
				OptionalQueryStructField(
					"ToAccounts",
					g.NewQueryStruct("ToAccounts").
						List("Accounts", "AccountIdentifier", g.ListOptions().NoParentheses().Required()),
					g.KeywordOptions().SQL("TO ACCOUNTS"),
				),
			g.KeywordOptions().SQL("DISABLE FAILOVER"),
		).
		OptionalSQL("PRIMARY").
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("ConnectionSet").
				OptionalComment().
				WithValidation(g.AtLeastOneValueSet, "Comment"),
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("ConnectionUnset").
				OptionalSQL("COMMENT").
				WithValidation(g.AtLeastOneValueSet, "Comment"),
			g.KeywordOptions().SQL("UNSET"),
		).
		WithValidation(g.ExactlyOneValueSet, "EnableConnectionFailover", "DisableConnectionFailover", "Primary", "Set", "Unset"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-connection",
	g.NewQueryStruct("DropConnection").
		Drop().
		SQL("CONNECTION").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-connections",
	g.StructPair("connectionRow", "Connection").
		OptionalText("region_group").
		Text("snowflake_region").
		Time("created_on").
		Text("account_name").
		Text("name").
		OptionalText("comment").
		Field("is_primary", "string", "bool", g.WithBoolParsed()).
		Field("primary", "string", "ExternalObjectIdentifier").
		AccountIdentifierArray("failover_allowed_to_accounts").
		Text("connection_url").
		Text("organization_name").
		Text("account_locator"),
	g.NewQueryStruct("ShowConnections").
		Show().
		SQL("CONNECTIONS").
		OptionalLike(),
).
	WithShowByIDFindPredicateKind(g.ShowByIDFindPredicateNameAndLocator)
