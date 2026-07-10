package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var schemaPairs = g.StructPair("schemaRow", "Schema").
	Time("created_on").
	OptionalTime("dropped_on", g.WithRequiredInPlain()).
	Text("name").
	BoolFromText("is_default").
	BoolFromText("is_current").
	Text("database_name").
	Text("owner").
	OptionalText("comment", g.WithRequiredInPlain()).
	OptionalText("options").
	Text("retention_time").
	Text("owner_role_type")

var schemaSetStruct = g.NewQueryStruct("SchemaSet").
	OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalIdentifier("ExternalVolume", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXTERNAL_VOLUME").Equals()).
	OptionalIdentifier("Catalog", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("CATALOG").Equals()).
	OptionalBooleanAssignment("PIPE_EXECUTION_PAUSED", nil).
	OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", nil).
	OptionalAssignment("DEFAULT_DDL_COLLATION", "StringAllowEmpty", g.ParameterOptions()).
	OptionalAssignment("STORAGE_SERIALIZATION_POLICY", "StorageSerializationPolicy", g.ParameterOptions()).
	WithField(g.OptionalEnumLegacy[sdkcommons.LogLevel]("LogLevel", g.ParameterOptions().SingleQuotes().SQL("LOG_LEVEL"))).
	WithField(g.OptionalEnumLegacy[sdkcommons.LogLevel]("LogEventLevel", g.ParameterOptions().SingleQuotes().SQL("LOG_EVENT_LEVEL"))).
	WithField(g.OptionalEnumLegacy[sdkcommons.TraceLevel]("TraceLevel", g.ParameterOptions().SingleQuotes().SQL("TRACE_LEVEL"))).
	OptionalNumberAssignment("SUSPEND_TASK_AFTER_NUM_FAILURES", g.ParameterOptions()).
	OptionalNumberAssignment("TASK_AUTO_RETRY_ATTEMPTS", g.ParameterOptions()).
	OptionalAssignment("USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE", "WarehouseSize", g.ParameterOptions()).
	OptionalNumberAssignment("USER_TASK_TIMEOUT_MS", g.ParameterOptions()).
	OptionalNumberAssignment("USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS", g.ParameterOptions()).
	OptionalBooleanAssignment("QUOTED_IDENTIFIERS_IGNORE_CASE", nil).
	OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
	OptionalComment().
	WithValidation(g.ValidIdentifierIfSet, "ExternalVolume").
	WithValidation(g.ValidIdentifierIfSet, "Catalog").
	WithValidation(g.AtLeastOneValueSet,
		"DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ExternalVolume", "Catalog",
		"PipeExecutionPaused", "ReplaceInvalidCharacters", "DefaultDdlCollation", "StorageSerializationPolicy",
		"LogLevel", "LogEventLevel", "TraceLevel",
		"SuspendTaskAfterNumFailures", "TaskAutoRetryAttempts", "UserTaskManagedInitialWarehouseSize",
		"UserTaskTimeoutMs", "UserTaskMinimumTriggerIntervalInSeconds",
		"QuotedIdentifiersIgnoreCase", "EnableConsoleOutput", "Comment")

var schemaUnsetStruct = g.NewQueryStruct("SchemaUnset").
	OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
	OptionalSQL("MAX_DATA_EXTENSION_TIME_IN_DAYS").
	OptionalSQL("EXTERNAL_VOLUME").
	OptionalSQL("CATALOG").
	OptionalSQL("PIPE_EXECUTION_PAUSED").
	OptionalSQL("REPLACE_INVALID_CHARACTERS").
	OptionalSQL("DEFAULT_DDL_COLLATION").
	OptionalSQL("STORAGE_SERIALIZATION_POLICY").
	OptionalSQL("LOG_LEVEL").
	OptionalSQL("LOG_EVENT_LEVEL").
	OptionalSQL("TRACE_LEVEL").
	OptionalSQL("SUSPEND_TASK_AFTER_NUM_FAILURES").
	OptionalSQL("TASK_AUTO_RETRY_ATTEMPTS").
	OptionalSQL("USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE").
	OptionalSQL("USER_TASK_TIMEOUT_MS").
	OptionalSQL("USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS").
	OptionalSQL("QUOTED_IDENTIFIERS_IGNORE_CASE").
	OptionalSQL("ENABLE_CONSOLE_OUTPUT").
	OptionalSQL("COMMENT").
	WithValidation(g.AtLeastOneValueSet,
		"DataRetentionTimeInDays", "MaxDataExtensionTimeInDays", "ExternalVolume", "Catalog",
		"PipeExecutionPaused", "ReplaceInvalidCharacters", "DefaultDdlCollation", "StorageSerializationPolicy",
		"LogLevel", "LogEventLevel", "TraceLevel",
		"SuspendTaskAfterNumFailures", "TaskAutoRetryAttempts", "UserTaskManagedInitialWarehouseSize",
		"UserTaskTimeoutMs", "UserTaskMinimumTriggerIntervalInSeconds",
		"QuotedIdentifiersIgnoreCase", "EnableConsoleOutput", "Comment")

var schemasDef = g.NewInterface(
	"Schemas",
	"Schema",
	g.KindOfT[sdkcommons.DatabaseObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-schema",
	g.NewQueryStruct("CreateSchema").
		Create().
		OrReplace().
		OptionalSQL("TRANSIENT").
		SQL("SCHEMA").
		IfNotExists().
		Name().
		OptionalSQL("WITH MANAGED ACCESS").
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalIdentifier("ExternalVolume", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXTERNAL_VOLUME").Equals()).
		OptionalIdentifier("Catalog", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("CATALOG").Equals()).
		OptionalBooleanAssignment("PIPE_EXECUTION_PAUSED", nil).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", nil).
		OptionalAssignment("DEFAULT_DDL_COLLATION", "StringAllowEmpty", g.ParameterOptions()).
		OptionalAssignment("STORAGE_SERIALIZATION_POLICY", "StorageSerializationPolicy", g.ParameterOptions()).
		WithField(g.OptionalEnumLegacy[sdkcommons.LogLevel]("LogLevel", g.ParameterOptions().SingleQuotes().SQL("LOG_LEVEL"))).
		WithField(g.OptionalEnumLegacy[sdkcommons.LogLevel]("LogEventLevel", g.ParameterOptions().SingleQuotes().SQL("LOG_EVENT_LEVEL"))).
		WithField(g.OptionalEnumLegacy[sdkcommons.TraceLevel]("TraceLevel", g.ParameterOptions().SingleQuotes().SQL("TRACE_LEVEL"))).
		OptionalNumberAssignment("SUSPEND_TASK_AFTER_NUM_FAILURES", g.ParameterOptions()).
		OptionalNumberAssignment("TASK_AUTO_RETRY_ATTEMPTS", g.ParameterOptions()).
		OptionalAssignment("USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE", "WarehouseSize", g.ParameterOptions()).
		OptionalNumberAssignment("USER_TASK_TIMEOUT_MS", g.ParameterOptions()).
		OptionalNumberAssignment("USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS", g.ParameterOptions()).
		OptionalBooleanAssignment("QUOTED_IDENTIFIERS_IGNORE_CASE", nil).
		OptionalBooleanAssignment("ENABLE_CONSOLE_OUTPUT", nil).
		OptionalComment().
		OptionalTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithValidation(g.ValidIdentifierIfSet, "ExternalVolume").
		WithValidation(g.ValidIdentifierIfSet, "Catalog"),
).CustomOperation(
	"Clone",
	"https://docs.snowflake.com/en/sql-reference/sql/create-schema",
	g.NewQueryStruct("CloneSchema").
		Create().
		OrReplace().
		OptionalSQL("TRANSIENT").
		SQL("SCHEMA").
		IfNotExists().
		Name().
		PredefinedQueryStructField("Clone", "Clone", g.KeywordOptions()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-schema",
	g.NewQueryStruct("AlterSchema").
		Alter().
		SQL("SCHEMA").
		IfExists().
		Name().
		Identifier("RenameTo", g.KindOfTPointer[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		Identifier("SwapWith", g.KindOfTPointer[sdkcommons.DatabaseObjectIdentifier](), g.IdentifierOptions().SQL("SWAP WITH")).
		OptionalQueryStructField("Set", schemaSetStruct, g.ListOptions().NoParentheses().SQL("SET")).
		OptionalQueryStructField("Unset", schemaUnsetStruct, g.ListOptions().NoParentheses().SQL("UNSET")).
		OptionalSetTags().
		OptionalUnsetTags().
		OptionalSQL("ENABLE MANAGED ACCESS").
		OptionalSQL("DISABLE MANAGED ACCESS").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "RenameTo", "SwapWith", "Set", "Unset", "SetTags", "UnsetTags", "EnableManagedAccess", "DisableManagedAccess").
		WithValidation(g.ValidIdentifierIfSet, "RenameTo").
		WithValidation(g.ValidIdentifierIfSet, "SwapWith"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-schema",
	g.NewQueryStruct("DropSchema").
		Drop().
		SQL("SCHEMA").
		IfExists().
		Name().
		OptionalSQL("CASCADE").
		OptionalSQL("RESTRICT").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "Cascade", "Restrict"),
).CustomOperation(
	"Undrop",
	"https://docs.snowflake.com/en/sql-reference/sql/undrop-schema",
	g.NewQueryStruct("UndropSchema").
		SQL("UNDROP").
		SQL("SCHEMA").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-schemas",
	schemaPairs,
	g.NewQueryStruct("ShowSchemas").
		Show().
		Terse().
		SQL("SCHEMAS").
		OptionalSQL("HISTORY").
		OptionalLike().
		OptionalExtendedIn().
		OptionalStartsWith().
		OptionalLimit(),
	g.ShowByIDExtendedInFiltering,
	g.ShowByIDLikeFiltering,
).DescribeOperationWithPairedStructs(
	g.DescriptionMappingKindSlice,
	"https://docs.snowflake.com/en/sql-reference/sql/desc-schema",
	g.StructPair("schemaDetailRow", "SchemaDetails").
		Time("created_on").
		Text("name").
		Text("kind"),
	g.NewQueryStruct("DescribeSchema").
		Describe().
		SQL("SCHEMA").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowParameters("DatabaseObjectIdentifier").
	WithCustomInterfaceMethod(
		"Use", "Use is based on https://docs.snowflake.com/en/sql-reference/sql/use-schema",
		[]*g.MethodParameter{g.NewMethodParameter("id", "DatabaseObjectIdentifier")},
		"error",
	)
