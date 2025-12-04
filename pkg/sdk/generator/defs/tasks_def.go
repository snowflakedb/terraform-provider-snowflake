package defs

import (
	g "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/generator/gen/sdkcommons"
)

var taskDbRow = g.DbStruct("taskDBRow").
	Text("created_on").
	Text("name").
	Text("id").
	Text("database_name").
	Text("schema_name").
	Text("owner").
	OptionalText("comment").
	OptionalText("warehouse").
	OptionalText("schedule").
	Text("predecessors").
	Text("state").
	Text("definition").
	OptionalText("condition").
	Text("allow_overlapping_execution").
	OptionalText("error_integration").
	OptionalText("last_committed_on").
	OptionalText("last_suspended_on").
	Text("owner_role_type").
	OptionalText("config").
	OptionalText("budget").
	Text("task_relations").
	OptionalText("last_suspended_reason").
	OptionalText("target_completion_interval")

var task = g.PlainStruct("Task").
	Text("CreatedOn").
	Text("Name").
	Text("Id").
	Text("DatabaseName").
	Text("SchemaName").
	Text("Owner").
	Text("Comment").
	Field("Warehouse", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier]()).
	Text("Schedule").
	Field("Predecessors", g.KindOfTSlice[sdkcommons.SchemaObjectIdentifier]()).
	Field("State", g.KindOfT[sdkcommons.TaskState]()).
	Text("Definition").
	Text("Condition").
	Bool("AllowOverlappingExecution").
	Field("ErrorIntegration", g.KindOfTPointer[sdkcommons.AccountObjectIdentifier]()).
	Text("LastCommittedOn").
	Text("LastSuspendedOn").
	Text("OwnerRoleType").
	Text("Config").
	Text("Budget").
	Field("TaskRelations", "TaskRelations").
	Text("LastSuspendedReason").
	Field("TargetCompletionInterval", "*TaskTargetCompletionInterval")

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
			OptionalSessionParameters().
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
			OptionalSessionParameters().
			OptionalNumberAssignment("SUSPEND_TASK_AFTER_NUM_FAILURES", nil).
			OptionalIdentifier("ErrorIntegration", g.KindOfT[sdkcommons.AccountObjectIdentifier](), g.IdentifierOptions().Equals().SQL("ERROR_INTEGRATION")).
			OptionalTextAssignment("COMMENT", g.ParameterOptions().SingleQuotes()).
			OptionalIdentifier("Finalize", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.IdentifierOptions().Equals().SQL("FINALIZE")).
			OptionalNumberAssignment("TASK_AUTO_RETRY_ATTEMPTS", g.ParameterOptions()).
			ListAssignment("AFTER", g.KindOfT[sdkcommons.SchemaObjectIdentifier](), g.ParameterOptions().NoEquals()).
			OptionalTextAssignment("WHEN", g.ParameterOptions().NoQuotes().NoEquals()).
			SQL("AS").
			Text("sql", g.KeywordOptions().NoQuotes().Required()).
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
					OptionalSessionParameters().
					OptionalNumberAssignment("TASK_AUTO_RETRY_ATTEMPTS", nil).
					OptionalNumberAssignment("USER_TASK_MINIMUM_TRIGGER_INTERVAL_IN_SECONDS", nil).
					OptionalTextAssignment("TARGET_COMPLETION_INTERVAL", g.ParameterOptions().SingleQuotes()).
					OptionalAssignment("SERVERLESS_TASK_MIN_STATEMENT_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
					OptionalAssignment("SERVERLESS_TASK_MAX_STATEMENT_SIZE", "WarehouseSize", g.ParameterOptions().SingleQuotes()).
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
					OptionalSessionParametersUnset().
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
	ShowOperation(
		"https://docs.snowflake.com/en/sql-reference/sql/show-tasks",
		taskDbRow,
		task,
		g.NewQueryStruct("ShowTasks").
			Show().
			Terse().
			SQL("TASKS").
			OptionalLike().
			OptionalExtendedIn().
			OptionalStartsWith().
			OptionalSQL("ROOT ONLY").
			OptionalLimit(),
	).
	ShowByIdOperationWithFiltering(
		g.ShowByIDExtendedInFiltering,
		g.ShowByIDLikeFiltering,
	).
	DescribeOperation(
		g.DescriptionMappingKindSingleValue,
		"https://docs.snowflake.com/en/sql-reference/sql/desc-task",
		taskDbRow,
		task,
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
	)
