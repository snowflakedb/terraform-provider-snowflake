package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var externalFunctionArgument = g.NewQueryStruct("ExternalFunctionArgument").
	Text("ArgName", g.KeywordOptions().NoQuotes().Required()).
	PredefinedQueryStructField("ArgDataType", g.KindOfT[sdkcommons.DataType](), g.KeywordOptions().NoQuotes().Required())

var externalFunctionHeader = g.NewQueryStruct("ExternalFunctionHeader").
	Text("Name", g.KeywordOptions().SingleQuotes().Required()).
	PredefinedQueryStructField("Value", g.KindOfT[string](), g.ParameterOptions().SingleQuotes().Required())

var externalFunctionContextHeader = g.NewQueryStruct("ExternalFunctionContextHeader").Text("ContextFunction", g.KeywordOptions().NoQuotes().Required())

var externalFunctionSet = g.NewQueryStruct("ExternalFunctionSet").
	OptionalIdentifier("ApiIntegration", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("API_INTEGRATION =")).
	ListQueryStructField(
		"Headers",
		externalFunctionHeader,
		g.ParameterOptions().Parentheses().SQL("HEADERS"),
	).
	ListQueryStructField(
		"ContextHeaders",
		externalFunctionContextHeader,
		g.ParameterOptions().Parentheses().SQL("CONTEXT_HEADERS"),
	).
	OptionalNumberAssignment("MAX_BATCH_ROWS", g.ParameterOptions()).
	OptionalTextAssignment("COMPRESSION", g.ParameterOptions()).
	OptionalIdentifier("RequestTranslator", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("REQUEST_TRANSLATOR =")).
	OptionalIdentifier("ResponseTranslator", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RESPONSE_TRANSLATOR =")).
	WithValidation(g.ExactlyOneValueSet, "ApiIntegration", "Headers", "ContextHeaders", "MaxBatchRows", "Compression", "RequestTranslator", "ResponseTranslator")

var externalFunctionUnset = g.NewQueryStruct("ExternalFunctionUnset").
	OptionalSQL("COMMENT").
	OptionalSQL("HEADERS").
	OptionalSQL("CONTEXT_HEADERS").
	OptionalSQL("MAX_BATCH_ROWS").
	OptionalSQL("COMPRESSION").
	OptionalSQL("SECURE").
	OptionalSQL("REQUEST_TRANSLATOR").
	OptionalSQL("RESPONSE_TRANSLATOR").
	WithValidation(g.AtLeastOneValueSet, "Comment", "Headers", "ContextHeaders", "MaxBatchRows", "Compression", "Secure", "RequestTranslator", "ResponseTranslator")

var externalFunctionPairs = g.StructPair("externalFunctionRow", "ExternalFunction").
	Text("created_on").
	Text("name").
	OptionalText("schema_name", g.WithManualConvert(), g.WithRequiredInPlain()).
	BoolFromText("is_builtin").
	BoolFromText("is_aggregate").
	BoolFromText("is_ansi").
	Number("min_num_arguments").
	Number("max_num_arguments").
	Text("arguments", g.WithPlainFieldName("ArgumentsRaw")).
	PlainOnlyField("Arguments", "[]DataType").
	Text("description").
	Field("catalog_name", "sql.NullString", "string", g.WithManualConvert()).
	BoolFromText("is_table_function").
	BoolFromText("valid_for_clustering").
	OptionalBoolFromText("is_secure", g.WithRequiredInPlain()).
	BoolFromText("is_external_function").
	Text("language").
	OptionalBoolFromText("is_memoizable", g.WithRequiredInPlain()).
	OptionalBoolFromText("is_data_metric", g.WithRequiredInPlain())

var externalFunctionPropertyPairs = g.StructPair("externalFunctionPropertyRow", "ExternalFunctionProperty").
	Text("property").
	Text("value")

// TODO [SNOW-2048276]: Add dedicated external Drop and DropSafely functions
var externalFunctionsDef = g.NewInterface(
	"ExternalFunctions",
	"ExternalFunction",
	g.KindOfT[sdkcommons.SchemaObjectIdentifierWithArguments](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-external-function",
	g.NewQueryStruct("CreateExternalFunction").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("EXTERNAL FUNCTION").
		Identifier("name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			externalFunctionArgument,
			g.ListOptions().MustParentheses()).
		PredefinedQueryStructField("ResultDataType", g.KindOfT[sdkcommons.DataType](), g.ParameterOptions().NoEquals().SQL("RETURNS").Required()).
		WithField(g.OptionalEnumLegacy[sdkcommons.ReturnNullValues]("ReturnNullValues", g.KeywordOptions())).
		WithField(g.OptionalEnumLegacy[sdkcommons.NullInputBehavior]("NullInputBehavior", g.KeywordOptions())).
		WithField(g.OptionalEnumLegacy[sdkcommons.ReturnResultsBehavior]("ReturnResultsBehavior", g.KeywordOptions())).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		Identifier("ApiIntegration", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("API_INTEGRATION =").Required()).
		ListQueryStructField(
			"Headers",
			externalFunctionHeader,
			g.ParameterOptions().Parentheses().SQL("HEADERS"),
		).
		ListQueryStructField(
			"ContextHeaders",
			externalFunctionContextHeader,
			g.ParameterOptions().Parentheses().SQL("CONTEXT_HEADERS"),
		).
		OptionalNumberAssignment("MAX_BATCH_ROWS", g.ParameterOptions()).
		OptionalTextAssignment("COMPRESSION", g.ParameterOptions()).
		OptionalIdentifier("RequestTranslator", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("REQUEST_TRANSLATOR =")).
		OptionalIdentifier("ResponseTranslator", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RESPONSE_TRANSLATOR =")).
		TextAssignment("AS", g.ParameterOptions().NoEquals().SingleQuotes().Required()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "ApiIntegration").
		WithValidation(g.ValidIdentifierIfSet, "RequestTranslator").
		WithValidation(g.ValidateValueSet, "As").
		WithValidation(g.ValidIdentifierIfSet, "ResponseTranslator"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-function",
	g.NewQueryStruct("AlterExternalFunction").
		Alter().
		SQL("FUNCTION").
		IfExists().
		Name().
		OptionalQueryStructField(
			"Set",
			externalFunctionSet,
			g.KeywordOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			externalFunctionUnset,
			g.ListOptions().NoParentheses().SQL("UNSET"),
		).
		WithValidation(g.ExactlyOneValueSet, "Set", "Unset").
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-external-functions",
	externalFunctionPairs,
	g.NewQueryStruct("ShowFunctions").
		Show().
		SQL("EXTERNAL FUNCTIONS").
		OptionalLike().
		OptionalIn(),
	g.ShowByIDInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-function",
	externalFunctionPropertyPairs,
	g.NewQueryStruct("DescribeExternalFunction").
		Describe().
		SQL("FUNCTION").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
