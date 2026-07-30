package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var CortexSearchServiceRefreshModeEnumDef = g.NewEnum(
	"CortexSearchServiceRefreshMode", "CortexSearchServiceRefreshModes",
	"FULL", "INCREMENTAL",
)

var CortexSearchServiceInitializeEnumDef = g.NewEnum(
	"CortexSearchServiceInitialize", "CortexSearchServiceInitializes",
	"ON_CREATE", "ON_SCHEDULE",
)

func alterCortexSearchServiceSet() *g.QueryStruct {
	return g.NewQueryStruct("CortexSearchServiceSet").
		OptionalTextAssignment("TARGET_LAG", g.ParameterOptions().SingleQuotes()).
		OptionalIdentifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE")).
		OptionalNumberAssignment("FULL_INDEX_BUILD_INTERVAL_DAYS", nil).
		OptionalBooleanAssignment("REQUEST_LOGGING", nil).
		OptionalNumberAssignment("AUTO_SUSPEND", nil).
		OptionalComment().
		WithValidation(g.AtLeastOneValueSet, "TargetLag", "Warehouse", "FullIndexBuildIntervalDays", "RequestLogging", "AutoSuspend", "Comment")
}

// This is a workaround for non-existing UNSETs and the fact that AUTO_SUSPEND is either int or NULL (not expressible directly by the sql builder).
func alterCortexSearchServiceSetDefaults() *g.QueryStruct {
	return g.NewQueryStruct("CortexSearchServiceSetDefaults").
		OptionalSQLWithCustomFieldName("AutoSuspend", "AUTO_SUSPEND = NULL").
		WithValidation(g.AtLeastOneValueSet, "AutoSuspend")
}

func alterCortexSearchServiceSetPrimaryKey() *g.QueryStruct {
	return g.NewQueryStruct("CortexSearchServiceSetPrimaryKey").
		ListAssignment("PRIMARY KEY", "string", g.ParameterOptions().Parentheses()).
		WithValidation(g.ValidateValueSet, "PrimaryKey")
}

func alterCortexSearchServiceSetAttributes() *g.QueryStruct {
	return g.NewQueryStruct("CortexSearchServiceSetAttributes").
		SQL("ATTRIBUTES").
		List("Columns", "string", g.ListOptions().NoEquals().Parentheses()).
		WithValidation(g.ValidateValueSet, "Columns")
}

// TODO [SNOW-3863599]: add support for scoring profiles
// TODO [SNOW-3863599]: add support for multi-column variant
// TODO [SNOW-3863599]: All columns are used without double quotes, so are case-insensitive; it's cause by Snowflake limitation: Invalid source query: quoted identifier or reserved word "<column name>" not allowed.
// TODO [SNOW-3863599]: REQUEST_LOGGING and AUTO_SUSPEND are not exposed in SHOW or DESCRIBE output.
var cortexSearchServicesDef = g.NewInterface(
	"CortexSearchServices",
	"CortexSearchService",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-cortex-search",
	g.NewQueryStruct("CreateCortexSearchService").
		Create().
		OrReplace().
		SQL("CORTEX SEARCH SERVICE").
		IfNotExists().
		Name().
		TextAssignment("ON", g.ParameterOptions().NoEquals().NoQuotes().Required().SQL("ON")).
		ListAssignment("PRIMARY KEY", "string", g.ParameterOptions().NoEquals().Parentheses()).
		OptionalQueryStructField(
			"Attributes",
			g.NewQueryStruct("Attributes").
				SQL("ATTRIBUTES").
				List("Columns", "string", g.ListOptions().NoEquals().NoParentheses()),
			g.KeywordOptions(),
		).
		Identifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().Required().SQL("WAREHOUSE")).
		TextAssignment("TARGET_LAG", g.ParameterOptions().SingleQuotes().Required()).
		OptionalTextAssignment("EMBEDDING_MODEL", g.ParameterOptions().SingleQuotes()).
		OptionalEnumAssignment("REFRESH_MODE", CortexSearchServiceRefreshModeEnumDef, g.ParameterOptions().NoQuotes()).
		OptionalEnumAssignment("INITIALIZE", CortexSearchServiceInitializeEnumDef, g.ParameterOptions().NoQuotes()).
		OptionalNumberAssignment("FULL_INDEX_BUILD_INTERVAL_DAYS", nil).
		OptionalBooleanAssignment("REQUEST_LOGGING", nil).
		OptionalNumberAssignment("AUTO_SUSPEND", nil).
		OptionalComment().
		PredefinedQueryStructField("QueryDefinition", "string", g.ParameterOptions().NoEquals().NoQuotes().Required().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "On").
		WithValidation(g.ValidateValueSet, "TargetLag").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-cortex-search",
	g.NewQueryStruct("AlterCortexSearchService").
		Alter().
		SQL("CORTEX SEARCH SERVICE").
		IfExists().
		Name().
		OptionalSQL("SUSPEND").
		OptionalSQL("RESUME").
		OptionalSQL("REFRESH").
		OptionalQueryStructField("Set", alterCortexSearchServiceSet(), g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("SetDefaults", alterCortexSearchServiceSetDefaults(), g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("SetPrimaryKey", alterCortexSearchServiceSetPrimaryKey(), g.KeywordOptions().SQL("SET")).
		OptionalQueryStructField("SetAttributes", alterCortexSearchServiceSetAttributes(), g.KeywordOptions().SQL("SET")).
		OptionalSQL("UNSET PRIMARY KEY").
		OptionalSQL("UNSET ATTRIBUTES").
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "Suspend", "Resume", "Refresh", "Set", "SetDefaults", "SetPrimaryKey", "SetAttributes", "UnsetPrimaryKey", "UnsetAttributes", "SetTags", "UnsetTags"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-cortex-search",
	g.StructPair("cortexSearchServiceRow", "CortexSearchService").
		Time("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		Text("warehouse").
		Text("target_lag").
		OptionalText("definition").
		OptionalText("search_column").
		OptionalPlainField("attribute_columns", "[]string").
		OptionalPlainField("columns", "[]string").
		OptionalPlainField("primary_key_columns", "[]string").
		Number("scoring_profile_count").
		OptionalText("comment", g.WithRequiredInPlain()),
	g.NewQueryStruct("ShowCortexSearchService").
		Show().
		SQL("CORTEX SEARCH SERVICES").
		OptionalLike().
		OptionalIn().
		OptionalStartsWith().
		OptionalLimitFrom(),
	g.ShowByIDLikeFiltering,
	g.ShowByIDInFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSingleValue,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-cortex-search",
	g.StructPair("cortexSearchServiceDetailsRow", "CortexSearchServiceDetails").
		Text("created_on").
		Text("name").
		Text("database_name").
		Text("schema_name").
		Text("target_lag").
		Text("warehouse").
		OptionalText("search_column").
		OptionalPlainField("attribute_columns", "[]string").
		OptionalPlainField("columns", "[]string").
		OptionalText("definition").
		OptionalText("comment").
		Text("service_query_url").
		Text("data_timestamp").
		Number("source_data_num_rows").
		Text("indexing_state").
		OptionalText("indexing_error").
		Text("serving_state").
		OptionalText("embedding_model").
		OptionalPlainField("primary_key_columns", "[]string").
		Number("scoring_profile_count").
		OptionalNumber("full_index_build_interval_days"),
	g.NewQueryStruct("DescribeCortexSearchService").
		Describe().
		SQL("CORTEX SEARCH SERVICE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-cortex-search",
	g.NewQueryStruct("DropCortexSearchService").
		Drop().
		SQL("CORTEX SEARCH SERVICE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).WithEnums(CortexSearchServiceRefreshModeEnumDef, CortexSearchServiceInitializeEnumDef)
