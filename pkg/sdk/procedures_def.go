package sdk

import g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/poc/generator"

var procedureArgument = g.NewQueryStruct("ProcedureArgument").
	Text("ArgName", g.KeywordOptions().DoubleQuotes().Required()).
	PredefinedQueryStructField("ArgDataTypeOld", "DataType", g.KeywordOptions().NoQuotes()).
	PredefinedQueryStructField("ArgDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
	PredefinedQueryStructField("DefaultValue", "*string", g.ParameterOptions().NoEquals().SQL("DEFAULT")).
	WithValidation(g.ExactlyOneValueSet, "ArgDataTypeOld", "ArgDataType")

var procedureColumn = g.NewQueryStruct("ProcedureColumn").
	Text("ColumnName", g.KeywordOptions().DoubleQuotes().Required()).
	PredefinedQueryStructField("ColumnDataTypeOld", "DataType", g.KeywordOptions().NoQuotes()).
	PredefinedQueryStructField("ColumnDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
	WithValidation(g.ExactlyOneValueSet, "ColumnDataTypeOld", "ColumnDataType")

var procedureReturns = g.NewQueryStruct("ProcedureReturns").
	OptionalQueryStructField(
		"ResultDataType",
		g.NewQueryStruct("ProcedureReturnsResultDataType").
			PredefinedQueryStructField("ResultDataTypeOld", "DataType", g.KeywordOptions().NoQuotes()).
			PredefinedQueryStructField("ResultDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
			OptionalSQL("NULL").OptionalSQL("NOT NULL").
			WithValidation(g.ExactlyOneValueSet, "ResultDataTypeOld", "ResultDataType"),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Table",
		g.NewQueryStruct("ProcedureReturnsTable").
			ListQueryStructField(
				"Columns",
				procedureColumn,
				g.ListOptions().MustParentheses(),
			),
		g.KeywordOptions().SQL("TABLE"),
	).WithValidation(g.ExactlyOneValueSet, "ResultDataType", "Table")

var procedureSQLReturns = g.NewQueryStruct("ProcedureSQLReturns").
	OptionalQueryStructField(
		"ResultDataType",
		g.NewQueryStruct("ProcedureReturnsResultDataType").
			PredefinedQueryStructField("ResultDataTypeOld", "DataType", g.KeywordOptions().NoQuotes()).
			PredefinedQueryStructField("ResultDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
			WithValidation(g.ExactlyOneValueSet, "ResultDataTypeOld", "ResultDataType"),
		g.KeywordOptions(),
	).
	OptionalQueryStructField(
		"Table",
		g.NewQueryStruct("ProcedureReturnsTable").
			ListQueryStructField(
				"Columns",
				procedureColumn,
				g.ListOptions().MustParentheses(),
			),
		g.KeywordOptions().SQL("TABLE"),
	).
	OptionalSQL("NOT NULL").
	WithValidation(g.ExactlyOneValueSet, "ResultDataType", "Table")

var (
	procedureImport  = g.NewQueryStruct("ProcedureImport").Text("ProcedureImport", g.KeywordOptions().SingleQuotes().Required())
	procedurePackage = g.NewQueryStruct("ProcedurePackage").Text("ProcedurePackage", g.KeywordOptions().SingleQuotes().Required())
)

// https://docs.snowflake.com/en/sql-reference/constructs/with and https://docs.snowflake.com/en/user-guide/queries-cte
var procedureWithClause = g.NewQueryStruct("ProcedureWithClause").
	SQLWithCustomFieldName("prefix", ",").
	Identifier("CteName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
	PredefinedQueryStructField("CteColumns", "[]string", g.KeywordOptions().Parentheses()).
	PredefinedQueryStructField("Statement", "string", g.ParameterOptions().NoEquals().NoQuotes().SQL("AS").Required())

var ProceduresDef = g.NewInterface(
	"Procedures",
	"Procedure",
	g.KindOfT[SchemaObjectIdentifierWithArguments](),
).CustomOperation(
	"CreateForJava",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#java-handler",
	g.NewQueryStruct("CreateForJava").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE JAVA").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES").Required(),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "SecretReference", g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.ParameterOptions().NoQuotes().NoEquals().SQL("EXECUTE AS")).
		PredefinedQueryStructField("ProcedureDefinition", "*string", g.ParameterOptions().NoEquals().SQL("AS")).
		WithValidation(g.ValidateValueSet, "RuntimeVersion").
		WithValidation(g.ValidateValueSet, "Packages").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateForJavaScript",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#javascript-handler",
	g.NewQueryStruct("CreateForJavaScript").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		OptionalSQL("COPY GRANTS").
		SQL("RETURNS").
		PredefinedQueryStructField("ResultDataTypeOld", "DataType", g.ParameterOptions().NoEquals()).
		PredefinedQueryStructField("ResultDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
		OptionalSQL("NOT NULL").
		SQL("LANGUAGE JAVASCRIPT").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.ParameterOptions().NoQuotes().NoEquals().SQL("EXECUTE AS")).
		PredefinedQueryStructField("ProcedureDefinition", "string", g.ParameterOptions().NoEquals().SQL("AS").Required()).
		WithValidation(g.ValidateValueSet, "ProcedureDefinition").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "ResultDataTypeOld", "ResultDataType"),
).CustomOperation(
	"CreateForPython",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#python-handler",
	g.NewQueryStruct("CreateForPython").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE PYTHON").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES").Required(),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "SecretReference", g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.ParameterOptions().NoQuotes().NoEquals().SQL("EXECUTE AS")).
		PredefinedQueryStructField("ProcedureDefinition", "*string", g.ParameterOptions().NoEquals().SQL("AS")).
		WithValidation(g.ValidateValueSet, "RuntimeVersion").
		WithValidation(g.ValidateValueSet, "Packages").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateForScala",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#scala-handler",
	g.NewQueryStruct("CreateForScala").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE SCALA").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES").Required(),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
		ListAssignment("SECRETS", "SecretReference", g.ParameterOptions().Parentheses()).
		OptionalTextAssignment("TARGET_PATH", g.ParameterOptions().SingleQuotes()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.ParameterOptions().NoQuotes().NoEquals().SQL("EXECUTE AS")).
		PredefinedQueryStructField("ProcedureDefinition", "*string", g.ParameterOptions().NoEquals().SQL("AS")).
		WithValidation(g.ValidateValueSet, "RuntimeVersion").
		WithValidation(g.ValidateValueSet, "Packages").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateForSQL",
	"https://docs.snowflake.com/en/sql-reference/sql/create-procedure#snowflake-scripting-handler",
	g.NewQueryStruct("CreateForSQL").
		Create().
		OrReplace().
		OptionalSQL("SECURE").
		SQL("PROCEDURE").
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		OptionalSQL("COPY GRANTS").
		QueryStructField(
			"Returns",
			procedureSQLReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE SQL").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ReturnResultsBehavior", "*ReturnResultsBehavior", g.KeywordOptions()).
		OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.ParameterOptions().NoQuotes().NoEquals().SQL("EXECUTE AS")).
		PredefinedQueryStructField("ProcedureDefinition", "string", g.ParameterOptions().NoEquals().SQL("AS").Required()).
		WithValidation(g.ValidateValueSet, "ProcedureDefinition").
		WithValidation(g.ValidIdentifier, "name"),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-procedure",
	g.NewQueryStruct("AlterProcedure").
		Alter().
		SQL("PROCEDURE").
		IfExists().
		Name().
		OptionalIdentifier("RenameTo", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		OptionalQueryStructField(
			"Set",
			g.NewQueryStruct("ProcedureSet").
				OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
				ListAssignment("EXTERNAL_ACCESS_INTEGRATIONS", "AccountObjectIdentifier", g.ParameterOptions().Parentheses()).
				OptionalQueryStructField("SecretsList", functionSecretsListWrapper, g.ParameterOptions().SQL("SECRETS").Parentheses()).
				OptionalAssignment("AUTO_EVENT_LOGGING", g.KindOfTPointer[AutoEventLogging](), g.ParameterOptions().SingleQuotes()).
				OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
				OptionalAssignment("LOG_LEVEL", g.KindOfTPointer[LogLevel](), g.ParameterOptions().SingleQuotes()).
				OptionalAssignment("METRIC_LEVEL", g.KindOfTPointer[MetricLevel](), g.ParameterOptions().SingleQuotes()).
				OptionalAssignment("TRACE_LEVEL", g.KindOfTPointer[TraceLevel](), g.ParameterOptions().SingleQuotes()).
				WithValidation(g.AtLeastOneValueSet, "Comment", "ExternalAccessIntegrations", "SecretsList", "AutoEventLogging", "EnableConsoleOutput", "LogLevel", "MetricLevel", "TraceLevel"),
			g.ListOptions().SQL("SET"),
		).
		OptionalQueryStructField(
			"Unset",
			g.NewQueryStruct("ProcedureUnset").
				OptionalSQL("COMMENT").
				OptionalSQL("EXTERNAL_ACCESS_INTEGRATIONS").
				OptionalSQL("AUTO_EVENT_LOGGING").
				OptionalSQL("ENABLE_CONSOLE_OUTPUT").
				OptionalSQL("LOG_LEVEL").
				OptionalSQL("METRIC_LEVEL").
				OptionalSQL("TRACE_LEVEL").
				WithValidation(g.AtLeastOneValueSet, "Comment", "ExternalAccessIntegrations", "AutoEventLogging", "EnableConsoleOutput", "LogLevel", "MetricLevel", "TraceLevel"),
			g.ListOptions().SQL("UNSET"),
		).
		OptionalSetTags().
		OptionalUnsetTags().
		PredefinedQueryStructField("ExecuteAs", "*ExecuteAs", g.ParameterOptions().NoQuotes().NoEquals().SQL("EXECUTE AS")).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "Set", "Unset", "SetTags", "UnsetTags", "ExecuteAs"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-procedure",
	g.NewQueryStruct("DropProcedure").
		Drop().
		SQL("PROCEDURE").
		IfExists().
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/show-procedures",
	g.DbStruct("procedureRow").
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
		OptionalText("external_access_integrations"),
	g.PlainStruct("Procedure").
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
		OptionalText("ExternalAccessIntegrations"),
	g.NewQueryStruct("ShowProcedures").
		Show().
		SQL("PROCEDURES").
		OptionalLike().
		OptionalExtendedIn(),
).ShowByIdOperationWithFiltering(
	g.ShowByIDInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperation(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-procedure",
	g.DbStruct("procedureDetailRow").
		Field("property", "string").
		Field("value", "sql.NullString"),
	g.PlainStruct("ProcedureDetail").
		Field("Property", "string").
		OptionalText("Value"),
	g.NewQueryStruct("DescribeProcedure").
		Describe().
		SQL("PROCEDURE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"Call",
	"https://docs.snowflake.com/en/sql-reference/sql/call",
	g.NewQueryStruct("Call").
		SQL("CALL").
		Identifier("name", g.KindOfT[SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("CallArguments", "[]string", g.KeywordOptions().MustParentheses()).
		PredefinedQueryStructField("ScriptingVariable", "*string", g.ParameterOptions().NoEquals().NoQuotes().SQL("INTO")).
		WithValidation(g.ValidIdentifier, "name"),
).CustomOperation(
	"CreateAndCallForJava",
	"https://docs.snowflake.com/en/sql-reference/sql/call-with#java-and-scala",
	g.NewQueryStruct("CreateAndCallForJava").
		SQL("WITH").
		Identifier("Name", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		SQL("AS PROCEDURE").
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE JAVA").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES").Required(),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		PredefinedQueryStructField("ProcedureDefinition", "*string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS")).
		OptionalQueryStructField(
			"WithClause",
			procedureWithClause,
			g.KeywordOptions(),
		).
		SQL("CALL").
		Identifier("ProcedureName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("CallArguments", "[]string", g.KeywordOptions().MustParentheses()).
		PredefinedQueryStructField("ScriptingVariable", "*string", g.ParameterOptions().NoEquals().NoQuotes().SQL("INTO")).
		WithValidation(g.ValidateValueSet, "RuntimeVersion").
		WithValidation(g.ValidateValueSet, "Packages").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ValidIdentifier, "ProcedureName").
		WithValidation(g.ValidIdentifier, "Name"),
).CustomOperation(
	"CreateAndCallForScala",
	"https://docs.snowflake.com/en/sql-reference/sql/call-with#java-and-scala",
	g.NewQueryStruct("CreateAndCallForScala").
		SQL("WITH").
		Identifier("Name", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		SQL("AS PROCEDURE").
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE SCALA").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES").Required(),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		PredefinedQueryStructField("ProcedureDefinition", "*string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS")).
		ListQueryStructField(
			"WithClauses",
			procedureWithClause,
			g.KeywordOptions(),
		).
		SQL("CALL").
		Identifier("ProcedureName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("CallArguments", "[]string", g.KeywordOptions().MustParentheses()).
		PredefinedQueryStructField("ScriptingVariable", "*string", g.ParameterOptions().NoEquals().NoQuotes().SQL("INTO")).
		WithValidation(g.ValidateValueSet, "RuntimeVersion").
		WithValidation(g.ValidateValueSet, "Packages").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ValidIdentifier, "ProcedureName").
		WithValidation(g.ValidIdentifier, "Name"),
).CustomOperation(
	"CreateAndCallForJavaScript",
	"https://docs.snowflake.com/en/sql-reference/sql/call-with#javascript",
	g.NewQueryStruct("CreateAndCallForJavaScript").
		SQL("WITH").
		Identifier("Name", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		SQL("AS PROCEDURE").
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		SQL("RETURNS").
		PredefinedQueryStructField("ResultDataTypeOld", "DataType", g.ParameterOptions().NoEquals()).
		PredefinedQueryStructField("ResultDataType", "datatypes.DataType", g.ParameterOptions().NoQuotes().NoEquals().Required()).
		OptionalSQL("NOT NULL").
		SQL("LANGUAGE JAVASCRIPT").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ProcedureDefinition", "string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS").Required()).
		ListQueryStructField(
			"WithClauses",
			procedureWithClause,
			g.KeywordOptions(),
		).
		SQL("CALL").
		Identifier("ProcedureName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("CallArguments", "[]string", g.KeywordOptions().MustParentheses()).
		PredefinedQueryStructField("ScriptingVariable", "*string", g.ParameterOptions().NoEquals().NoQuotes().SQL("INTO")).
		WithValidation(g.ValidateValueSet, "ProcedureDefinition").
		WithValidation(g.ExactlyOneValueSet, "ResultDataTypeOld", "ResultDataType").
		WithValidation(g.ValidIdentifier, "ProcedureName").
		WithValidation(g.ValidIdentifier, "Name"),
).CustomOperation(
	"CreateAndCallForPython",
	"https://docs.snowflake.com/en/sql-reference/sql/call-with#python",
	g.NewQueryStruct("CreateAndCallForPython").
		SQL("WITH").
		Identifier("Name", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		SQL("AS PROCEDURE").
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE PYTHON").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		TextAssignment("RUNTIME_VERSION", g.ParameterOptions().SingleQuotes().Required()).
		ListQueryStructField(
			"Packages",
			procedurePackage,
			g.ParameterOptions().Parentheses().SQL("PACKAGES").Required(),
		).
		ListQueryStructField(
			"Imports",
			procedureImport,
			g.ParameterOptions().Parentheses().SQL("IMPORTS"),
		).
		TextAssignment("HANDLER", g.ParameterOptions().SingleQuotes().Required()).
		PredefinedQueryStructField("ProcedureDefinition", "*string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS")).
		ListQueryStructField(
			"WithClauses",
			procedureWithClause,
			g.KeywordOptions(),
		).
		SQL("CALL").
		Identifier("ProcedureName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("CallArguments", "[]string", g.KeywordOptions().MustParentheses()).
		PredefinedQueryStructField("ScriptingVariable", "*string", g.ParameterOptions().NoEquals().NoQuotes().SQL("INTO")).
		WithValidation(g.ValidateValueSet, "RuntimeVersion").
		WithValidation(g.ValidateValueSet, "Packages").
		WithValidation(g.ValidateValueSet, "Handler").
		WithValidation(g.ValidIdentifier, "ProcedureName").
		WithValidation(g.ValidIdentifier, "Name"),
).CustomOperation(
	"CreateAndCallForSQL",
	"https://docs.snowflake.com/en/sql-reference/sql/call-with#snowflake-scripting",
	g.NewQueryStruct("CreateAndCallForSQL").
		SQL("WITH").
		Identifier("Name", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		SQL("AS PROCEDURE").
		ListQueryStructField(
			"Arguments",
			procedureArgument,
			g.ListOptions().MustParentheses(),
		).
		QueryStructField(
			"Returns",
			procedureReturns,
			g.KeywordOptions().SQL("RETURNS").Required(),
		).
		SQL("LANGUAGE SQL").
		PredefinedQueryStructField("NullInputBehavior", "*NullInputBehavior", g.KeywordOptions()).
		PredefinedQueryStructField("ProcedureDefinition", "string", g.ParameterOptions().NoEquals().SingleQuotes().SQL("AS").Required()).
		ListQueryStructField(
			"WithClauses",
			procedureWithClause,
			g.KeywordOptions(),
		).
		SQL("CALL").
		Identifier("ProcedureName", g.KindOfT[AccountObjectIdentifier](), g.IdentifierOptions().Required()).
		PredefinedQueryStructField("CallArguments", "[]string", g.KeywordOptions().MustParentheses()).
		PredefinedQueryStructField("ScriptingVariable", "*string", g.ParameterOptions().NoEquals().NoQuotes().SQL("INTO")).
		WithValidation(g.ValidateValueSet, "ProcedureDefinition").
		WithValidation(g.ValidIdentifier, "ProcedureName").
		WithValidation(g.ValidIdentifier, "Name"),
)
