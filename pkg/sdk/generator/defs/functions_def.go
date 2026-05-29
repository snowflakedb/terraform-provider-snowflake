package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var functionArgument = func() *g.QueryStruct {
	return g.NewQueryStruct("FunctionArgument").
		Text("ArgName", g.KeywordOptions().DoubleQuotes().Required()).
		PredefinedQueryStructField("ArgDataTypeOld", "DataType", g.KeywordOptions().NoQuotes()).
		PredefinedQueryStructField("ArgDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
		PredefinedQueryStructField("DefaultValue", "*string", g.ParameterOptions().NoEquals().SQL("DEFAULT"))
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// .WithValidation(g.ExactlyOneValueSet, "ArgDataTypeOld", "ArgDataType")
}

var functionColumn = func() *g.QueryStruct {
	return g.NewQueryStruct("FunctionColumn").
		Text("ColumnName", g.KeywordOptions().DoubleQuotes().Required()).
		PredefinedQueryStructField("ColumnDataTypeOld", "DataType", g.KeywordOptions().NoQuotes()).
		PredefinedQueryStructField("ColumnDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required())
	// TODO [next PR]: validation is not generated properly as this is used as an array; using the additionalValidations above for now
	// .WithValidation(g.ExactlyOneValueSet, "ColumnDataTypeOld", "ColumnDataType")
}

var functionReturns = func() *g.QueryStruct {
	return g.NewQueryStruct("FunctionReturns").
		OptionalQueryStructField(
			"ResultDataType",
			g.NewQueryStruct("FunctionReturnsResultDataType").
				PredefinedQueryStructField("ResultDataTypeOld", "DataType", g.KeywordOptions().NoQuotes()).
				PredefinedQueryStructField("ResultDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
				WithValidation(g.ExactlyOneValueSet, "ResultDataTypeOld", "ResultDataType"),
			g.KeywordOptions(),
		).
		OptionalQueryStructField(
			"Table",
			g.NewQueryStruct("FunctionReturnsTable").
				ListQueryStructField(
					"Columns",
					functionColumn(),
					g.ParameterOptions().Parentheses().NoEquals(),
				).
				WithAdditionalValidations(),
			g.KeywordOptions().SQL("TABLE"),
		).WithValidation(g.ExactlyOneValueSet, "ResultDataType", "Table")
}

var (
	functionImports            = g.NewQueryStruct("FunctionImport").Text("FunctionImport", g.KeywordOptions().SingleQuotes())
	functionPackages           = g.NewQueryStruct("FunctionPackage").Text("FunctionPackage", g.KeywordOptions().SingleQuotes())
	functionSecretsListWrapper = g.NewQueryStruct("SecretsList").
					List("SecretsList", "SecretReference", g.ListOptions().Required().MustParentheses())
)

// TODO [next PRs]: support adding a field only in plain struct (in this case: `ArgumentsOld` and `ReturnTypeOld`)
var functionPairs = g.StructPair("functionRow", "Function").
	Text("created_on").
	Text("name").
	Text("schema_name", g.WithManualConvert()).
	BoolFromText("is_builtin").
	BoolFromText("is_aggregate").
	BoolFromText("is_ansi").
	Number("min_num_arguments").
	Number("max_num_arguments").
	Text("arguments", g.WithPlainFieldName("ArgumentsRaw")).
	Text("description").
	Text("catalog_name", g.WithManualConvert()).
	BoolFromText("is_table_function").
	BoolFromText("valid_for_clustering").
	OptionalBoolFromText("is_secure", g.WithRequiredInPlain()).
	OptionalText("secrets").
	OptionalText("external_access_integrations").
	BoolFromText("is_external_function").
	Text("language").
	OptionalBoolFromText("is_memoizable", g.WithRequiredInPlain()).
	OptionalBoolFromText("is_data_metric", g.WithRequiredInPlain()).
	WithConvertGeneration()

var functionDetailPairs = g.StructPair("functionDetailRow", "FunctionDetail").
	Text("property").
	OptionalText("value", g.WithManualConvert()).
	WithConvertGeneration()

var functionsDef = g.NewInterface(
	"Functions",
	"Function",
	g.KindOfT[sdkcommons.SchemaObjectIdentifierWithArguments](),
).CustomOperation(
	"CreateForJava",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#java-handler",
	g.NewQueryStruct("CreateForJava").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		IfNotExists().
		Identifier("name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument(),
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns(),
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		SQL("LANGUAGE JAVA").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalTextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		ListQueryStructField(
			"Imports",
			functionImports,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		ListQueryStructField(
			"Packages",
			functionPackages,
			g.ParameterOptions().Parentheses().SQL("PACKAGES"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "SecretReference", g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[sdkcommons.LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[sdkcommons.MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[sdkcommons.TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "*string", g.ParameterOptions().NoEquals().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).CustomOperation(
	"CreateForJavascript",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#javascript-handler",
	g.NewQueryStruct("CreateForJavascript").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		Identifier("name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument(),
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns(),
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		SQL("LANGUAGE JAVASCRIPT").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[sdkcommons.LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[sdkcommons.MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[sdkcommons.TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "string", g.ParameterOptions().NoEquals().SQL("AS").Required()).
		WithValidation(g.ValidateValueSet, "FunctionDefinition").
		WithValidation(g.ValidIdentifier, "name").
		WithAdditionalValidations(),
).CustomOperation(
	"CreateForPython",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#python-handler",
	g.NewQueryStruct("CreateForPython").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		OptionalSQL("AGGREGATE").
		SQL("FUNCTION").
		IfNotExists().
		Identifier("name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument(),
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns(),
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		SQL("LANGUAGE PYTHON").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		ListQueryStructField(
			"Imports",
			functionImports,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		ListQueryStructField(
			"Packages",
			functionPackages,
			g.ParameterOptions().Parentheses().SQL("PACKAGES"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "SecretReference", g.ParameterOptions().Parentheses()).
		OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[sdkcommons.LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[sdkcommons.MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[sdkcommons.TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "*string", g.ParameterOptions().NoEquals().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "RuntimeVersion").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).CustomOperation(
	"CreateForScala",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#scala-handler",
	g.NewQueryStruct("CreateForScala").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		IfNotExists().
		Identifier("name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument(),
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		SQL("RETURNS").
		PredefinedQueryStructField("ResultDataTypeOld", "DataType", g.ParameterOptions().NoEquals()).
		PredefinedQueryStructField("ResultDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		SQL("LANGUAGE SCALA").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		ListQueryStructField(
			"Imports",
			functionImports,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		ListQueryStructField(
			"Packages",
			functionPackages,
			g.ParameterOptions().Parentheses().SQL("PACKAGES"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "SecretReference", g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[sdkcommons.LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[sdkcommons.MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[sdkcommons.TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "*string", g.ParameterOptions().NoEquals().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithValidation(g.ExactlyOneValueSet, "ResultDataTypeOld", "ResultDataType").
		WithAdditionalValidations(),
).CustomOperation(
	"CreateForSQL",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#sql-handler",
	g.NewQueryStruct("CreateForSQL").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		Identifier("name", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument(),
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns(),
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalSQL("MEMOIZABLE").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[sdkcommons.LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[sdkcommons.MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[sdkcommons.TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "string", g.ParameterOptions().NoEquals().SQL("AS").Required()).
		WithValidation(g.ValidateValueSet, "FunctionDefinition").
		WithValidation(g.ValidIdentifier, "name").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-function",
	g.NewQueryStruct("AlterFunction").
		Alter().
		SQL("FUNCTION").
		IfExists().
		Name().
		Identifier("RenameTo", g.KindOfTPointer[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("FunctionSet").
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
				OptionalQueryStructField("SecretsList", functionSecretsListWrapper, g.ParameterOptions().SQL("SECRETS").Parentheses()).
				OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
				OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[sdkcommons.LogLevel](), g.ParameterOptions().SingleQuotes()).
				OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[sdkcommons.MetricLevel](), g.ParameterOptions().SingleQuotes()).
				OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[sdkcommons.TraceLevel](), g.ParameterOptions().SingleQuotes()).
				WithValidation(g.AtLeastOneValueSet, "Comment", "ExternalAccessIntegrations", "SecretsList", "EnableConsoleOutput", "LogLevel", "MetricLevel", "TraceLevel"),
			g.ListOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("FunctionUnset").
				OptionalSQL("COMMENT").
				OptionalSQL("EXTERNAL_ACCESS_INTEGRATIONS").
				OptionalSQL("ENABLE_CONSOLE_OUTPUT").
				OptionalSQL("LOG_LEVEL").
				OptionalSQL("METRIC_LEVEL").
				OptionalSQL("TRACE_LEVEL").
				WithValidation(g.AtLeastOneValueSet, "Comment", "ExternalAccessIntegrations", "EnableConsoleOutput", "LogLevel", "MetricLevel", "TraceLevel"),
			g.ListOptions().SQL("UNSET"),
		).
		OptionalSQL("SET SECURE").
		OptionalSQL("UNSET SECURE").
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set", "Unset", "SetSecure", "UnsetSecure", "SetTags", "UnsetTags"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-function",
	g.NewQueryStruct("DropFunction").
		Drop().
		SQL("FUNCTION").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-user-functions",
	functionPairs,
	g.NewQueryStruct("ShowFunctions").
		Show().
		SQL("USER FUNCTIONS").
		OptionalLike().
		OptionalExtendedIn(),
	g.ShowByIDExtendedInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-function",
	functionDetailPairs,
	g.NewQueryStruct("DescribeFunction").
		Describe().
		SQL("FUNCTION").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).WithCustomInterfaceMethod(
	"DescribeDetails",
	"DescribeDetails returns aggregated describe results for the given function.",
	[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifierWithArguments]())},
	"*FunctionDetails", "error",
).WithCustomInterfaceMethod(
	"ShowParameters",
	"",
	[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifierWithArguments]())},
	"[]*Parameter", "error",
)
