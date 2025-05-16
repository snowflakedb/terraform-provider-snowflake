package sdk

//go:generate go run ./poc/main.go

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/v2/pkg/sdk/poc/generator"
)

var ConnectionDef = g.NewInterface(
	"Connections",
	"Connection",
	g.KindOfT[AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-connection",
	g.NewQueryStruct("CreateConnection").
		Create().
		SQL("CONNECTION").
		IfNotExists().
		Name().
		OptionalIdentifier(
			"AsReplicaOf",
			g.KindOfT[ExternalObjectIdentifier](),
			g.IdentifierOptions().SQL("AS REPLICA OF")).
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
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-connections",
	g.DbStruct("connectionRow").
		OptionalText("region_group").
		Text("snowflake_region").
		Field("created_on", "time.Time").
		Text("account_name").
		Text("name").
		Field("comment", "sql.NullString").
		Text("is_primary").
		Text("primary").
		Text("failover_allowed_to_accounts").
		Text("connection_url").
		Text("organization_name").
		Text("account_locator"),
	g.PlainStruct("Connection").
		OptionalText("RegionGroup").
		Text("SnowflakeRegion").
		Field("CreatedOn", "time.Time").
		Text("AccountName").
		Text("Name").
		OptionalText("Comment").
		Bool("IsPrimary").
		Field("Primary", "ExternalObjectIdentifier").
		Field("FailoverAllowedToAccounts", "[]AccountIdentifier").
		Text("ConnectionUrl").
		Text("OrganizationName").
		Text("AccountLocator"),
	g.NewQueryStruct("ShowConnections").
		Show().
		SQL("CONNECTIONS").
		OptionalLike(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
)
