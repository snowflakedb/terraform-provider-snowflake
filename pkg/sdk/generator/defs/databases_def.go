package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var databasePairs = g.StructPair("databaseRow", "Database").
	Time("created_on").
	Text("name").
	OptionalBoolFromText("is_default", g.WithRequiredInPlain()).
	OptionalBoolFromText("is_current", g.WithRequiredInPlain()).
	Field("origin", "sql.NullString", "ObjectIdentifier", g.WithManualConvert()).
	OptionalText("owner", g.WithRequiredInPlain()).
	OptionalText("comment", g.WithRequiredInPlain()).
	OptionalText("options", g.WithRequiredInPlain()).
	Field("retention_time", "sql.NullString", "int", g.WithManualConvert()).
	OptionalText("resource_group", g.WithRequiredInPlain()).
	OptionalTime("dropped_on", g.WithRequiredInPlain()).
	PlainOnlyField("Transient", "bool").
	OptionalText("kind", g.WithRequiredInPlain()).
	OptionalText("owner_role_type", g.WithRequiredInPlain())

var databaseSetStruct = g.NewQueryStruct("DatabaseSet").
	OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
	OptionalIdentifier("ExternalVolume", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXTERNAL_VOLUME").Equals()).
	OptionalIdentifier("Catalog", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("CATALOG").Equals()).
	OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", nil).
	OptionalTextAssignment("DEFAULT_DDL_COLLATION", g.ParameterOptions().SingleQuotes()).
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
		"ReplaceInvalidCharacters", "DefaultDdlCollation", "StorageSerializationPolicy",
		"LogLevel", "LogEventLevel", "TraceLevel",
		"SuspendTaskAfterNumFailures", "TaskAutoRetryAttempts", "UserTaskManagedInitialWarehouseSize",
		"UserTaskTimeoutMs", "UserTaskMinimumTriggerIntervalInSeconds",
		"QuotedIdentifiersIgnoreCase", "EnableConsoleOutput", "Comment")

var databaseUnsetStruct = g.NewQueryStruct("DatabaseUnset").
	OptionalSQL("DATA_RETENTION_TIME_IN_DAYS").
	OptionalSQL("MAX_DATA_EXTENSION_TIME_IN_DAYS").
	OptionalSQL("EXTERNAL_VOLUME").
	OptionalSQL("CATALOG").
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
		"ReplaceInvalidCharacters", "DefaultDdlCollation", "StorageSerializationPolicy",
		"LogLevel", "LogEventLevel", "TraceLevel",
		"SuspendTaskAfterNumFailures", "TaskAutoRetryAttempts", "UserTaskManagedInitialWarehouseSize",
		"UserTaskTimeoutMs", "UserTaskMinimumTriggerIntervalInSeconds",
		"QuotedIdentifiersIgnoreCase", "EnableConsoleOutput", "Comment")

var databasesDef = g.NewInterface(
	"Databases",
	"Database",
	g.KindOfT[sdkcommons.AccountObjectIdentifier](),
).CreateOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/create-database",
	g.NewQueryStruct("CreateDatabase").
		Create().
		OrReplace().
		OptionalSQL("TRANSIENT").
		SQL("DATABASE").
		IfNotExists().
		Name().
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalIdentifier("ExternalVolume", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXTERNAL_VOLUME").Equals()).
		OptionalIdentifier("Catalog", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("CATALOG").Equals()).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", nil).
		OptionalTextAssignment("DEFAULT_DDL_COLLATION", g.ParameterOptions().SingleQuotes()).
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
	"https://docs.snowflake.com/en/sql-reference/sql/create-database",
	g.NewQueryStruct("CloneDatabase").
		Create().
		OrReplace().
		OptionalSQL("TRANSIENT").
		SQL("DATABASE").
		IfNotExists().
		Name().
		PredefinedQueryStructField("Clone", "Clone", g.KeywordOptions()).
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithAdditionalValidations(),
).CustomOperation(
	"CreateShared",
	"https://docs.snowflake.com/en/sql-reference/sql/create-database",
	g.NewQueryStruct("CreateSharedDatabase").
		Create().
		OrReplace().
		OptionalSQL("TRANSIENT").
		SQL("DATABASE").
		IfNotExists().
		Name().
		Identifier("FromShare", g.KindOfT[sdkcommons.ExternalObjectIdentifier](), g.IdentifierOptions().SQL("FROM SHARE").Required()).
		OptionalIdentifier("ExternalVolume", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXTERNAL_VOLUME").Equals()).
		OptionalIdentifier("Catalog", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("CATALOG").Equals()).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", nil).
		OptionalTextAssignment("DEFAULT_DDL_COLLATION", g.ParameterOptions().SingleQuotes()).
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
		WithValidation(g.ValidIdentifier, "FromShare").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithValidation(g.ValidIdentifierIfSet, "ExternalVolume").
		WithValidation(g.ValidIdentifierIfSet, "Catalog"),
).CustomOperation(
	"CreateSecondary",
	"https://docs.snowflake.com/en/sql-reference/sql/create-database",
	g.NewQueryStruct("CreateSecondaryDatabase").
		Create().
		OrReplace().
		OptionalSQL("TRANSIENT").
		SQL("DATABASE").
		IfNotExists().
		Name().
		Identifier("PrimaryDatabase", g.KindOfT[sdkcommons.ExternalObjectIdentifier](), g.IdentifierOptions().SQL("AS REPLICA OF").Required()).
		OptionalNumberAssignment("DATA_RETENTION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalNumberAssignment("MAX_DATA_EXTENSION_TIME_IN_DAYS", g.ParameterOptions()).
		OptionalIdentifier("ExternalVolume", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("EXTERNAL_VOLUME").Equals()).
		OptionalIdentifier("Catalog", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("CATALOG").Equals()).
		OptionalBooleanAssignment("REPLACE_INVALID_CHARACTERS", nil).
		OptionalTextAssignment("DEFAULT_DDL_COLLATION", g.ParameterOptions().SingleQuotes()).
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
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ValidIdentifier, "PrimaryDatabase").
		WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists").
		WithValidation(g.ValidIdentifierIfSet, "ExternalVolume").
		WithValidation(g.ValidIdentifierIfSet, "Catalog"),
).CustomOperation(
	"CreateFromListing",
	"https://docs.snowflake.com/en/sql-reference/sql/create-database",
	g.NewQueryStruct("CreateDatabaseFromListing").
		Create().
		SQL("DATABASE").
		Name().
		TextAssignment("FROM LISTING", g.ParameterOptions().SingleQuotes().NoEquals()).
		WithValidation(g.ValidIdentifier, "name").
		WithAdditionalValidations(),
).AlterOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/alter-database",
	g.NewQueryStruct("AlterDatabase").
		Alter().
		SQL("DATABASE").
		IfExists().
		Name().
		Identifier("NewName", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("RENAME TO")).
		Identifier("SwapWith", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().SQL("SWAP WITH")).
		OptionalQueryStructField("Set", databaseSetStruct, g.ListOptions().NoParentheses().SQL("SET")).
		OptionalQueryStructField("Unset", databaseUnsetStruct, g.ListOptions().NoParentheses().SQL("UNSET")).
		OptionalSetTags().
		OptionalUnsetTags().
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "NewName", "Set", "Unset", "SwapWith", "SetTags", "UnsetTags").
		WithValidation(g.ValidIdentifierIfSet, "NewName").
		WithValidation(g.ValidIdentifierIfSet, "SwapWith"),
).CustomOperation(
	"AlterReplication",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-database",
	g.NewQueryStruct("AlterDatabaseReplication").
		Alter().
		SQL("DATABASE").
		Name().
		OptionalQueryStructField("EnableReplication",
			g.NewQueryStruct("EnableReplication").
				WithField(g.NewField("ToAccounts", g.KindOfSlice("AccountIdentifier"), g.Tags().Keyword().NoParentheses().SQL("TO ACCOUNTS"), nil)).
				OptionalSQL("IGNORE EDITION CHECK"),
			g.KeywordOptions().SQL("ENABLE REPLICATION")).
		OptionalQueryStructField("DisableReplication",
			g.NewQueryStruct("DisableReplication").
				WithField(g.NewField("ToAccounts", g.KindOfSlice("AccountIdentifier"), g.Tags().Keyword().NoParentheses().SQL("TO ACCOUNTS"), nil)),
			g.KeywordOptions().SQL("DISABLE REPLICATION")).
		OptionalSQL("REFRESH").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "EnableReplication", "DisableReplication", "Refresh"),
).CustomOperation(
	"AlterFailover",
	"https://docs.snowflake.com/en/sql-reference/sql/alter-database",
	g.NewQueryStruct("AlterDatabaseFailover").
		Alter().
		SQL("DATABASE").
		Name().
		OptionalQueryStructField("EnableFailover",
			g.NewQueryStruct("EnableFailover").
				WithField(g.NewField("ToAccounts", g.KindOfSlice("AccountIdentifier"), g.Tags().Keyword().NoParentheses().SQL("TO ACCOUNTS"), nil)),
			g.KeywordOptions().SQL("ENABLE FAILOVER")).
		OptionalQueryStructField("DisableFailover",
			g.NewQueryStruct("DisableFailover").
				WithField(g.NewField("ToAccounts", g.KindOfSlice("AccountIdentifier"), g.Tags().Keyword().NoParentheses().SQL("TO ACCOUNTS"), nil)),
			g.KeywordOptions().SQL("DISABLE FAILOVER")).
		OptionalSQL("PRIMARY").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ExactlyOneValueSet, "EnableFailover", "DisableFailover", "Primary"),
).DropOperation(
	"https://docs.snowflake.com/en/sql-reference/sql/drop-database",
	g.NewQueryStruct("DropDatabase").
		Drop().
		SQL("DATABASE").
		IfExists().
		Name().
		OptionalSQL("CASCADE").
		OptionalSQL("RESTRICT").
		WithValidation(g.ValidIdentifier, "name").
		WithValidation(g.ConflictingFields, "Cascade", "Restrict"),
).CustomOperation(
	"Undrop",
	"https://docs.snowflake.com/en/sql-reference/sql/undrop-database",
	g.NewQueryStruct("UndropDatabase").
		SQL("UNDROP").
		SQL("DATABASE").
		Name().
		WithValidation(g.ValidIdentifier, "name"),
).ShowOperationWithPairedStructs(
	"https://docs.snowflake.com/en/sql-reference/sql/show-databases",
	databasePairs,
	g.NewQueryStruct("ShowDatabases").
		Show().
		OptionalSQL("TERSE").
		SQL("DATABASES").
		OptionalSQL("HISTORY").
		OptionalLike().
		OptionalStartsWith().
		OptionalLimitFrom(),
	g.ShowByIDLikeFiltering,
).ShowParameters("AccountObjectIdentifier").
	WithCustomInterfaceMethod(
		"Use", "Use is based on https://docs.snowflake.com/en/sql-reference/sql/use-database",
		[]*g.MethodParameter{g.NewMethodParameter("id", "AccountObjectIdentifier")},
		"error",
	).WithCustomInterfaceMethod(
	"Describe", "Describe is based on https://docs.snowflake.com/en/sql-reference/sql/desc-database",
	[]*g.MethodParameter{g.NewMethodParameter("id", "AccountObjectIdentifier")},
	"*DatabaseDetails", "error",
)
