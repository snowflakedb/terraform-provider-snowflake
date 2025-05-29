package sdk

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"
)

//go:generate go run ./poc/main.go

var alterServiceSet = g.NewQueryStruct("CortexSearchServiceSet").
	// Fields
	OptionalTextAssignment("TARGET_LAG", g.ParameterOptions().SingleQuotes()).
	OptionalIdentifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE")).
	OptionalComment().
	WithValidation(g.AtLeastOneValueSet, "TargetLag", "Warehouse", "Comment")

var CortexSearchServiceDef = g.NewInterface(
	"CortexSearchServices",
	"CortexSearchService",
	g.KindOfT[SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/create-cortex-search",
	g.NewQueryStruct("CreateCortexSearchService").
		// Fields
		Create().
		OrReplace().
		SQL("CORTEX SEARCH SERVICE").
		IfNotExists().
		Name().
		TextAssignment("ON", g.ParameterOptions().NoEquals().NoQuotes().Required().SQL("ON")).
		OptionalQueryStructField(
			"Attributes",
			g.NewQueryStruct("Attributes").
				SQL("ATTRIBUTES").
				List("Columns", "string", g.ListOptions().NoEquals().NoParentheses()),
			g.KeywordOptions(),
		).
		Identifier("Warehouse", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Equals().Required().SQL("WAREHOUSE")).
		TextAssignment("TARGET_LAG", g.ParameterOptions().SingleQuotes().Required()).
		OptionalTextAssignment("EMBEDDING_MODEL", g.ParameterOptions().SingleQuotes()).
		OptionalComment().
		PredefinedQueryStructField("QueryDefinition", "string", g.ParameterOptions().NoEquals().NoQuotes().Required().SQL("AS")).
		// Validations
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "On").
		WithValidation(g.ValidateValueSet, "TargetLag").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/alter-cortex-search",
	g.NewQueryStruct("AlterCortexSearchService").
		// Fields
		Alter().
		SQL("CORTEX SEARCH SERVICE").
		IfExists().
		Name().
		OptionalQueryStructField("Set", alterServiceSet, g.KeywordOptions().SQL("SET")).
		// Validations
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Set"),
).ShowOperation(
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/show-cortex-search",
	// Fields
	g.DbStruct("cortexSearchServiceRow").
		Field("created_on", "time.Time").
		Field("name", "string").
		Field("database_name", "string").
		Field("schema_name", "string").
		OptionalText("comment"),
	g.PlainStruct("CortexSearchService").
		Field("CreatedOn", "time.Time").
		Field("Name", "string").
		Field("DatabaseName", "string").
		Field("SchemaName", "string").
		Field("Comment", "string"),
	g.NewQueryStruct("ShowCortexSearchService").
		Show().
		SQL("CORTEX SEARCH SERVICES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/desc-cortex-search",
	g.DbStruct("cortexSearchServiceDetailsRow").
		Text("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		Text("target_lag").
		Text("warehouse").
		OptionalText("search_column").
		OptionalText("attribute_columns").
		OptionalText("columns").
		OptionalText("definition").
		OptionalText("comment").
		Text("service_query_url").
		Text("data_timestamp").
		Number("source_data_num_rows").
		Text("indexing_state").
		OptionalText("indexing_error").
		OptionalText("embedding_model"),
	g.PlainStruct("CortexSearchServiceDetails").
		Text("CreatedOn").
		Text("Name").
		Text("DatabaseName").
		Text("SchemaName").
		Text("TargetLag").
		Text("Warehouse").
		OptionalText("SearchColumn").
		Field("AttributeColumns", "[]string").
		Field("Columns", "[]string").
		OptionalText("Definition").
		OptionalText("Comment").
		Text("ServiceQueryUrl").
		Text("DataTimestamp").
		Number("SourceDataNumRows").
		Text("IndexingState").
		OptionalText("IndexingError").
		OptionalText("EmbeddingModel"),
	g.NewQueryStruct("DescribeCortexSearchService").
		Describe().
		SQL("CORTEX SEARCH SERVICE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DropOperation(
	"https://docs.snowflake.com/LIMITEDACCESS/cortex-search/sql/drop-cortex-search",
	g.NewQueryStruct("DropCortexSearchService").
		// Fields
		Drop().
		SQL("CORTEX SEARCH SERVICE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
