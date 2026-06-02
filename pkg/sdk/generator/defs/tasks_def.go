package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var taskPairs = g.StructPair("taskDBRow", "Task").
	Text("created_on").
	Text("name").
	Text("id").
	Text("database_name").
	Text("schema_name").
	Text("owner").
	OptionalText("comment", g.WithRequiredInPlain()).
	Field("warehouse", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("Warehouse"), g.WithManualConvert()).
	OptionalText("schedule", g.WithRequiredInPlain()).
	Field("predecessors", "string", "[]SchemaObjectIdentifier", g.WithManualConvert()).
	PlainField("state", "TaskState", g.WithManualConvert()).
	Text("definition").
	OptionalText("condition", g.WithRequiredInPlain()).
	Field("allow_overlapping_execution", "string", "bool", g.WithBoolTrueValue("true")).
	Field("error_integration", "sql.NullString", "*AccountObjectIdentifier", g.WithPlainFieldName("ErrorIntegration"), g.WithManualConvert()).
	OptionalText("last_committed_on", g.WithRequiredInPlain()).
	OptionalText("last_suspended_on", g.WithRequiredInPlain()).
	Text("owner_role_type").
	OptionalText("config", g.WithRequiredInPlain()).
	OptionalText("budget", g.WithRequiredInPlain()).
	PlainField("task_relations", "TaskRelations", g.WithCustomParser("ToTaskRelations")).
	OptionalText("last_suspended_reason", g.WithRequiredInPlain()).
	Field("target_completion_interval", "sql.NullString", "*TaskTargetCompletionInterval", g.WithPlainFieldName("TargetCompletionInterval"), g.WithManualConvert())

var taskCreateWarehouse = g.NewQueryStruct("CreateTaskWarehouse").
	OptionalIdentifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE")).
	OptionalAssignment("USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
	WithValidation(g.ExactlyOneValueSet, "Warehouse", "UserTaskManagedInitialWarehouseSize")

var tasksDef = g.NewInterface(
	"Tasks",
	"Task",
	g.KindOfT[sdkcommons.SchemaObjectIdentifier](),
).
	CreateOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/create-task",
		g.NewQueryStruct("CreateTask").
			Create().
			OrReplace().
			SQL("TASK").
			IfNotExists().
			Name().
			PredefinedQueryStructField("Warehouse", "*CreateTaskWarehouse", g.KeywordOptions()).
			OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("CONFIG", g.ParameterOptions().NoQuotes()).
			OptionalBooleanAssignment("ALLOW_OVERLAPPING_EXECUTION", nil).
			PredefinedQueryStructField("SessionParameters", "*SessionParameters", g.ListOptions().NoParentheses()).
			OptionalNumberAssignment("USER_TASK_TIMEOUT_MS", nil).
			OptionalNumberAssignment("SUSPEND_TASK_AFTER_NUM_FAILURES", nil).
			OptionalIdentifier("ErrorIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("ERROR_INTEGRATION")).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalIdentifier("Finalize", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Equals().SQL("FINALIZE")).
			OptionalNumberAssignment("TASK_AUTO_RETRY_ATTEMPTS", g.ParameterOptions()).
			OptionalTags().
			OptionalNumberAssignment("USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS", g.ParameterOptions()).
			OptionalTextAssignment("TARGET_COMPLETION_INTERVAL", g.ParameterOptions().SingleQuotes()).
			OptionalAssignment("SERVERLESS_TASK_MIN_STATEMENT_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
			OptionalAssignment("SERVERLESS_TASK_MAX_STATEMENT_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
			ListAssignment("AFTER", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().NoEquals()).
			OptionalTextAssignment("WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithAdditionalValidations().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifierIfSet, "ErrorIntegration").
			WithValidation(g.ConflictingFields, "OrReplace", "IfNotExists"),
		taskCreateWarehouse,
	).
	CustomOperation(
		"CreateOrAlter",
		"https://docs.snowflake.com/en/sql-reference/sql/create-task#create-or-alter-task",
		g.NewQueryStruct("CloneTask").
			CreateOrAlter().
			SQL("TASK").
			Name().
			PredefinedQueryStructField("Warehouse", "*CreateTaskWarehouse", g.KeywordOptions()).
			OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
			OptionalTextAssignment("CONFIG", g.ParameterOptions().NoQuotes()).
			OptionalBooleanAssignment("ALLOW_OVERLAPPING_EXECUTION", nil).
			OptionalNumberAssignment("USER_TASK_TIMEOUT_MS", nil).
			PredefinedQueryStructField("SessionParameters", "*SessionParameters", g.ListOptions().NoParentheses()).
			OptionalNumberAssignment("SUSPEND_TASK_AFTER_NUM_FAILURES", nil).
			OptionalIdentifier("ErrorIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("ERROR_INTEGRATION")).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalIdentifier("Finalize", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Equals().SQL("FINALIZE")).
			OptionalNumberAssignment("TASK_AUTO_RETRY_ATTEMPTS", g.ParameterOptions()).
			ListAssignment("AFTER", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().NoEquals()).
			OptionalTextAssignment("WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
			WithAdditionalValidations().
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifierIfSet, "ErrorIntegration"),
	).
	CustomOperation(
		"Clone",
		"https://docs.snowflake.com/en/sql-reference/sql/create-task#create-task-clone",
		g.NewQueryStruct("CloneTask").
			Create().
			OrReplace().
			SQL("TASK").
			Name().
			SQL("CLONE").
			Identifier("sourceTask", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Required()).
			OptionalSQL("COPY GRANTS").
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ValidIdentifier, "sourceTask"),
	).
	AlterOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/alter-task",
		g.NewQueryStruct("AlterTask").
			Alter().
			SQL("TASK").
			IfExists().
			Name().
			OptionalSQL("RESUME").
			OptionalSQL("SUSPEND").
			ListAssignment("REMOVE AFTER", "SchemaObjectIdentifier", g.ParameterOptions().NoEquals()).
			ListAssignment("ADD AFTER", "SchemaObjectIdentifier", g.ParameterOptions().NoEquals()).
			OptionalQueryStructField(
				"Set",
				g.NewQueryStruct("TaskSet").
					OptionalIdentifier("Warehouse", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("WAREHOUSE")).
					OptionalAssignment("USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("SCHEDULE", g.ParameterOptions().SingleQuotes()).
					OptionalTextAssignment("CONFIG", g.ParameterOptions().NoQuotes()).
					OptionalBooleanAssignment("ALLOW_OVERLAPPING_EXECUTION", nil).
					OptionalNumberAssignment("USER_TASK_TIMEOUT_MS", nil).
					OptionalNumberAssignment("SUSPEND_TASK_AFTER_NUM_FAILURES", nil).
					OptionalIdentifier("ErrorIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("ERROR_INTEGRATION")).
					OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
					PredefinedQueryStructField("SessionParameters", "*SessionParameters", g.ListOptions().NoParentheses()).
					OptionalNumberAssignment("TASK_AUTO_RETRY_ATTEMPTS", nil).
					OptionalNumberAssignment("USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS", nil).
					OptionalTextAssignment("TARGET_COMPLETION_INTERVAL", g.ParameterOptions().SingleQuotes()).
					OptionalAssignment("SERVERLESS_TASK_MIN_STATEMENT_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
					OptionalAssignment("SERVERLESS_TASK_MAX_STATEMENT_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
					WithAdditionalValidations().
					WithValidation(g.AtLeastOneValueSet, "Warehouse", "UserTaskManagedInitialWarehouseSize", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "ErrorIntegration", "Comment", "SessionParameters", "TaskAutoRetryAttempts", "UserTaskMinimumTriggerIntervalInSeconds", "TargetCompletionInterval", "ServerlessTaskMinStatementSize", "ServerlessTaskMaxStatementSize").
					WithValidation(g.ConflictingFields, "Warehouse", "UserTaskManagedInitialWarehouseSize").
					WithValidation(g.ValidIdentifierIfSet, "ErrorIntegration"),
				g.ListOptions().SQL("SET").NoParentheses(),
			).
			OptionalQueryStructField(
				"Unset",
				g.NewQueryStruct("TaskUnset").
					OptionalSQL("WAREHOUSE").
					OptionalSQL("USER_TASK_MANAGED_INITIAL_WAREHOUSE_SIZE").
					OptionalSQL("SCHEDULE").
					OptionalSQL("CONFIG").
					OptionalSQL("ALLOW_OVERLAPPING_EXECUTION").
					OptionalSQL("USER_TASK_TIMEOUT_MS").
					OptionalSQL("SUSPEND_TASK_AFTER_NUM_FAILURES").
					OptionalSQL("ERROR_INTEGRATION").
					OptionalSQL("COMMENT").
					OptionalSQL("TASK_AUTO_RETRY_ATTEMPTS").
					OptionalSQL("USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS").
					OptionalSQL("TARGET_COMPLETION_INTERVAL").
					OptionalSQL("SERVERLESS_TASK_MIN_STATEMENT_SIZE").
					OptionalSQL("SERVERLESS_TASK_MAX_STATEMENT_SIZE").
					PredefinedQueryStructField("SessionParametersUnset", "*SessionParametersUnset", g.ListOptions().NoParentheses()).
					WithAdditionalValidations().
					WithValidation(g.AtLeastOneValueSet, "Warehouse", "UserTaskManagedInitialWarehouseSize", "Schedule", "Config", "AllowOverlappingExecution", "UserTaskTimeoutMs", "SuspendTaskAfterNumFailures", "ErrorIntegration", "Comment", "SessionParametersUnset", "TaskAutoRetryAttempts", "UserTaskMinimumTriggerIntervalInSeconds", "TargetCompletionInterval", "ServerlessTaskMinStatementSize", "ServerlessTaskMaxStatementSize"),
				g.ListOptions().SQL("UNSET").NoParentheses(),
			).
			OptionalSetTags().
			OptionalUnsetTags().
			OptionalIdentifier("SetFinalize", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Equals().SQL("SET FINALIZE")).
			OptionalSQL("UNSET FINALIZE").
			OptionalTextAssignment("MODIFY AS", g.ParameterOptions().NoQuotes().NoEquals()).
			OptionalTextAssignment("MODIFY WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			OptionalSQL("REMOVE WHEN").
			WithValidation(g.ValidIdentifier, "name").
			WithValidation(g.ExactlyOneValueSet, "Resume", "Suspend", "RemoveAfter", "AddAfter", "Set", "Unset", "SetTags", "UnsetTags", "SetFinalize", "UnsetFinalize", "ModifyAs", "ModifyWhen", "RemoveWhen"),
	).
	DropOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/drop-task",
		g.NewQueryStruct("DropTask").
			Drop().
			SQL("TASK").
			IfExists().
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	ShowOperationWithPairedStructs(
		"https://docs.snowflake.com/en/sql-reference/sql/show-tasks",
		taskPairs,
		g.NewQueryStruct("ShowTasks").
			Show().
			Terse().
			SQL("TASKS").
			OptionalLike().
			OptionalExtendedIn().
			OptionalStartsWith().
			OptionalSQL("ROOT ONLY").
			OptionalLimit(),
		g.ShowByIDExtendedInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperationWithPairedStructs(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-task",
		taskPairs,
		g.NewQueryStruct("DescribeTask").
			Describe().
			SQL("TASK").
			Name().
			WithValidation(g.ValidIdentifier, "name"),
	).
	CustomOperation(
		"Execute",
		"https://docs.snowflake.com/en/sql-reference/sql/execute-task",
		g.NewQueryStruct("ExecuteTask").
			SQL("EXECUTE").
			SQL("TASK").
			Name().
			OptionalSQL("RETRY LAST").
			WithValidation(g.ValidIdentifier, "name"),
	).
	WithCustomInterfaceMethod(
		"ShowParameters",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]())},
		"[]*Parameter", "error",
	).
	WithCustomInterfaceMethod(
		"SuspendRootTasks",
		"",
		[]*g.MethodParameter{
			g.NewMethodParameter("taskId", g.KindOfT[sdkcommons.SchemaObjectIdentifier]()),
			g.NewMethodParameter("id", g.KindOfT[sdkcommons.SchemaObjectIdentifier]()),
		},
		"[]SchemaObjectIdentifier", "error",
	).
	WithCustomInterfaceMethod(
		"ResumeTasks",
		"",
		[]*g.MethodParameter{g.NewMethodParameter("ids", "[]SchemaObjectIdentifier")},
		"error",
	)
