package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

var externalAccessIntegrations = g.NewQueryStruct("ExternalAccessIntegrations").
	List("ExternalAccessIntegrations", "AccountObjectIdentifier", g.ListOptions().Required().MustParentheses())

var streamlitSet = g.NewQueryStruct("StreamlitSet").
	OptionalTextAssignment("ROOT_LOCATION", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("MAIN_FILE", g.ParameterOptions().SingleQuotes()).
	OptionalIdentifier("QueryWarehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("QUERY_WAREHOUSE")).
	OptionalQueryStructField("ExternalAccessIntegrations", externalAccessIntegrations, g.ParameterOptions().SQL("EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
	OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("TITLE", g.ParameterOptions().SingleQuotes()).
	WithValidation(g.ValidIdentifierIfSet, "QueryWarehouse").
	WithValidation(g.AtLeastOneValueSet, "RootLocation", "MainFile", "QueryWarehouse", "ExternalAccessIntegrations", "Comment", "Title")

var streamlitUnset = g.NewQueryStruct("StreamlitUnset").
	OptionalSQL("QUERY_WAREHOUSE").
	OptionalSQL("COMMENT").
	OptionalSQL("TITLE").
	WithValidation(g.AtLeastOneValueSet, "QueryWarehouse", "Title", "Comment")

var StreamlitsDef = g.NewInterface(
	"Streamlits",
	"Streamlit",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-streamlit",
	g.NewQueryStruct("CreateStreamlit").
		Create().
		OrReplace().
		SQL("STREAMLIT").
		IfNotExists().
		Name().
		TextAssignment("ROOT_LOCATION", g.ParameterOptions().SingleQuotes().Required()).
		TextAssignment("MAIN_FILE", g.ParameterOptions().SingleQuotes().Required()).
		OptionalIdentifier("QueryWarehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("QUERY_WAREHOUSE")).
		OptionalQueryStructField("ExternalAccessIntegrations", externalAccessIntegrations, g.ParameterOptions().SQL("EXTERNAL_ACCESS_INTEGRATIONS").Parentheses()).
		OptionalTextAssignment("TITLE", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "QueryWarehouse").
		WithValidation(g.ConflictingFields, "IfNotExists", "OrReplace"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-streamlit",
	g.NewQueryStruct("AlterStreamlit").
		Alter().
		SQL("STREAMLIT").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			streamlitSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			streamlitUnset,
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set", "Unset"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-streamlit",
	g.NewQueryStruct("DropStreamlit").
		Drop().
		SQL("STREAMLIT").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-streamlits",
	g.DbStruct("streamlitsRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("database_name", "string").
		Field("schema_name", "string").
		Field("title", "sql.NullString").
		Field("owner", "string").
		Field("comment", "sql.NullString").
		Field("query_warehouse", "sql.NullString").
		Field("url_id", "string").
		Field("owner_role_type", "string"),
	g.PlainStruct("Streamlit").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("DatabaseName", "string").
		Field("SchemaName", "string").
		Field("Title", "string").
		Field("Owner", "string").
		Field("Comment", "string").
		Field("QueryWarehouse", "string").
		Field("UrlId", "string").
		Field("OwnerRoleType", "string"),
	g.NewQueryStruct("ShowStreamlits").
		Show().
		Terse().
		SQL("STREAMLITS").
		OptionalLike().
		OptionalIn().
		OptionalLimit(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-streamlit",
	g.DbStruct("streamlitsDetailRow").
		Field("name", "string").
		Field("title", "sql.NullString").
		Field("root_location", "string").
		Field("main_file", "string").
		Field("query_warehouse", "sql.NullString").
		Field("url_id", "string").
		Field("default_packages", "string").
		Field("user_packages", "string").
		Field("import_urls", "string").
		Field("external_access_integrations", "string").
		Field("external_access_secrets", "string"),
	g.PlainStruct("StreamlitDetail").
		Field("Name", "string").
		Field("Title", "string").
		Field("RootLocation", "string").
		Field("MainFile", "string").
		Field("QueryWarehouse", "string").
		Field("UrlId", "string").
		Field("DefaultPackages", "string").
		Field("UserPackages", "[]string").
		Field("ImportUrls", "[]string").
		Field("ExternalAccessIntegrations", "[]string").
		Field("ExternalAccessSecrets", "string"),
	g.NewQueryStruct("DescribeStreamlit").
		Describe().
		SQL("STREAMLIT").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
