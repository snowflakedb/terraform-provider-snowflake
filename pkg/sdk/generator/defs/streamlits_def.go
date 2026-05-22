package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var externalAccessIntegrations = g.NewQueryStruct("ExternalAccessIntegrations").
	List("ExternalAccessIntegrations", "AccountObjectIdentifier", g.ListOptions().Required().MustParentheses())

var streamlitSet = g.NewQueryStruct("StreamlitSet").
	OptionalTextAssignment("ROOT_LOCATION", g.ParameterOptions().SingleQuotes()).
	OptionalTextAssignment("MAIN_FILE", g.ParameterOptions().SingleQuotes()).
	OptionalIdentifier("QueryWarehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("QUERY_WAREHOUSE")).
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

var streamlitsDef = g.NewInterface(
	"Streamlits",
	"Streamlit",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
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
		OptionalIdentifier("QueryWarehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("QUERY_WAREHOUSE")).
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
		Identifier("RenameTo", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
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
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-streamlits",
	g.StructPair("streamlitsRow", "Streamlit").
		Text("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		OptionalText("title", g.WithRequiredInPlain()).
		Text("owner").
		OptionalText("comment", g.WithRequiredInPlain()).
		OptionalText("query_warehouse", g.WithRequiredInPlain()).
		Text("url_id").
		Text("owner_role_type").
		WithConvertGeneration(),
	g.NewQueryStruct("ShowStreamlits").
		Show().
		Terse().
		SQL("STREAMLITS").
		OptionalLike().
		OptionalIn().
		OptionalLimit(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-streamlit",
	g.StructPair("streamlitsDetailRow", "StreamlitDetail").
		Text("name").
		OptionalText("title", g.WithRequiredInPlain()).
		Text("root_location").
		Text("main_file").
		OptionalText("query_warehouse", g.WithRequiredInPlain()).
		Text("url_id").
		Text("default_packages").
		StringList("user_packages").
		StringList("import_urls").
		StringList("external_access_integrations").
		Text("external_access_secrets").
		WithConvertGeneration(),
	g.NewQueryStruct("DescribeStreamlit").
		Describe().
		SQL("STREAMLIT").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
