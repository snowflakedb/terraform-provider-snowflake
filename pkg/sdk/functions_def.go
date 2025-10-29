package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

var functionArgument = g.NewQueryStruct("FunctionArgument").
	Text("ArgName", g.KeywordOptions().DoubleQuotes().Required()).
	PredefinedQueryStructField("ArgDataTypeOld", "DataType", g.KeywordOptions().NoQuotes()).
	PredefinedQueryStructField("ArgDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
	PredefinedQueryStructField("DefaultValue", "*string", g.ParameterOptions().NoEquals().SQL("DEFAULT")).
	WithValidation(g.ExactlyOneValueSet, "ArgDataTypeOld", "ArgDataType")

var functionColumn = g.NewQueryStruct("FunctionColumn").
	Text("ColumnName", g.KeywordOptions().DoubleQuotes().Required()).
	PredefinedQueryStructField("ColumnDataTypeOld", "DataType", g.KeywordOptions().NoQuotes()).
	PredefinedQueryStructField("ColumnDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
	WithValidation(g.ExactlyOneValueSet, "ColumnDataTypeOld", "ColumnDataType")

var functionReturns = g.NewQueryStruct("FunctionReturns").
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
				functionColumn,
				g.ParameterOptions().Parentheses().NoEquals(),
			),
		g.KeywordOptions().SQL("TABLE"),
	).WithValidation(g.ExactlyOneValueSet, "ResultDataType", "Table")

var (
	functionImports            = g.NewQueryStruct("FunctionImport").Text("FunctionImport", g.KeywordOptions().SingleQuotes())
	functionPackages           = g.NewQueryStruct("FunctionPackage").Text("FunctionPackage", g.KeywordOptions().SingleQuotes())
	functionSecretsListWrapper = g.NewQueryStruct("SecretsList").
					List("SecretsList", "SecretReference", g.ListOptions().Required().MustParentheses())
)

var FunctionsDef = g.NewInterface(
	"Functions",
	"Function",
	g.KindOfT[SchemaObjectIdentifierWithArguments](),
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
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns,
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
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "*string", g.ParameterOptions().NoEquals().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
).CustomOperation(
	"CreateForJavascript",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#javascript-handler",
	g.NewQueryStruct("CreateForJavascript").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		SQL("LANGUAGE JAVASCRIPT").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "string", g.ParameterOptions().NoEquals().SQL("AS").Required()).
		WithValidation(g.ValidateValueSet, "FunctionDefinition").
		WithValidation(g.ValidIdentifier, "name"),
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
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns,
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
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "*string", g.ParameterOptions().NoEquals().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "RuntimeVersion").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
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
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument,
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
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "*string", g.ParameterOptions().NoEquals().SQL("AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithValidation(g.ExactlyOneValueSet, "ResultDataTypeOld", "ResultDataType"),
).CustomOperation(
	"CreateForSQL",
	"https://docs.snowflake.com/en/sql-reference/sql/create-function#sql-handler",
	g.NewQueryStruct("CreateForSQL").
		Create().
		OrReplace().
		OptionalSQL("TEMPORARY").
		OptionalSQL("SECURE").
		SQL("FUNCTION").
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			functionArgument,
			g.ListOptions().MustParentheses()).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			functionReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		PredefinedQueryStructField("ReturnNullValues", "*ReturnNullValues", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalSQL("MEMOIZABLE").
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
		OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[LogLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[MetricLevel](), g.ParameterOptions().SingleQuotes()).
		OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[TraceLevel](), g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("FunctionDefinition", "string", g.ParameterOptions().NoEquals().SQL("AS").Required()).
		WithValidation(g.ValidateValueSet, "FunctionDefinition").
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-function",
	g.NewQueryStruct("AlterFunction").
		Alter().
		SQL("FUNCTION").
		IfExists().
		Name().
		Identifier("RenameTo", g.KindOfTPointer[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("FunctionSet").
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
				OptionalQueryStructField("SecretsList", functionSecretsListWrapper, g.ParameterOptions().SQL("SECRETS").Parentheses()).
				OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
				OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[LogLevel](), g.ParameterOptions().SingleQuotes()).
				OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[MetricLevel](), g.ParameterOptions().SingleQuotes()).
				OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[TraceLevel](), g.ParameterOptions().SingleQuotes()).
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
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-user-functions",
	g.DbStruct("functionRow").
		Field("created_on", "string").
		Field("name", "string").
		Field("schema_name", "string").
		Field("is_builtin", "string").
		Field("is_aggregate", "string").
		Field("is_ansi", "string").
		Field("min_num_arguments", "int").
		Field("max_num_arguments", "int").
		Field("arguments", "string").
		Field("description", "string").
		Field("catalog_name", "string").
		Field("is_table_function", "string").
		Field("valid_for_clustering", "string").
		Field("is_secure", "sql.NullString").
		OptionalText("secrets").
		OptionalText("external_access_integrations").
		Field("is_external_function", "string").
		Field("language", "string").
		Field("is_memoizable", "sql.NullString").
		Field("is_data_metric", "sql.NullString"),
	g.PlainStruct("Function").
		Field("CreatedOn", "string").
		Field("Name", "string").
		Field("SchemaName", "string").
		Field("IsBuiltin", "bool").
		Field("IsAggregate", "bool").
		Field("IsAnsi", "bool").
		Field("MinNumArguments", "int").
		Field("MaxNumArguments", "int").
		Field("ArgumentsOld", "[]DataType").
		Field("ReturnTypeOld", "DataType").
		Field("ArgumentsRaw", "string").
		Field("Description", "string").
		Field("CatalogName", "string").
		Field("IsTableFunction", "bool").
		Field("ValidForClustering", "bool").
		Field("IsSecure", "bool").
		OptionalText("Secrets").
		OptionalText("ExternalAccessIntegrations").
		Field("IsExternalFunction", "bool").
		Field("Language", "string").
		Field("IsMemoizable", "bool").
		Field("IsDataMetric", "bool"),
	g.NewQueryStruct("ShowFunctions").
		Show().
		SQL("USER FUNCTIONS").
		OptionalLike().
		OptionalExtendedIn(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDExtendedInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-function",
	g.DbStruct("functionDetailRow").
		Field("property", "string").
		Field("value", "sql.NullString"),
	g.PlainStruct("FunctionDetail").
		Text("Property").
		OptionalText("Value"),
	g.NewQueryStruct("DescribeFunction").
		Describe().
		SQL("FUNCTION").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
)
